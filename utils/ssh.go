package utils

import (
    "bytes"
    dCfg "deploy/config"
    "fmt"
    "github.com/pkg/sftp"
    "golang.org/x/crypto/ssh"
    "io"
    "io/ioutil"
    "net"
    "os"
    "path"
    "time"
)

type ServerConn struct {
    addr           string
    user           string
    privateKey     string
    privateKeyPath string
    sshClient      *ssh.Client
    sftpClient     *sftp.Client
}

func NewServerConn(addr, user, privateKeyPath string) *ServerConn {
    return &ServerConn{
        addr:           addr,
        user:           user,
        privateKeyPath: privateKeyPath,
    }
}

// 连接ssh服务器
func (s *ServerConn) getSshConnect() (sshClient *ssh.Client, err error) {
    var (
        keys   []ssh.Signer
        config ssh.ClientConfig
        signer ssh.Signer
        buffer []byte
    )
    if s.sshClient != nil {
        sshClient = s.sshClient
        return
    }
    config = ssh.ClientConfig{
        User: s.user,
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
        Timeout: time.Duration(dCfg.GConfig.SshConnectTimeout) * time.Second, //连接超时时间
    }
    if buffer, err = ioutil.ReadFile(s.privateKeyPath); err != nil {
        err = fmt.Errorf("IP: %s 读取私钥失败: %v", s.addr, err)
        return
    }
    if signer, err = ssh.ParsePrivateKey(buffer); err != nil {
        err = fmt.Errorf("连接服务器失败[%s] : %v", s.addr, err.Error())
        return
    }
    keys = append(keys, signer)
    config.Auth = append(config.Auth, ssh.PublicKeys(keys...))
    if s.sshClient, err = ssh.Dial("tcp", s.addr, &config); err != nil {
        err = fmt.Errorf("无法连接到服务器[%s] : %v", s.addr, err.Error())
        return
    }
    sshClient = s.sshClient
    return
}

// 返回sftp连接
func (s *ServerConn) getSftpConnect() (sftpClient *sftp.Client, err error) {
    var sshClient *ssh.Client
    if s.sftpClient != nil {
        sftpClient = s.sftpClient
        return
    }
    if sshClient, err = s.getSshConnect(); err != nil {
        return
    }
    s.sftpClient, err = sftp.NewClient(sshClient, sftp.MaxPacket(1<<15))
    sftpClient = s.sftpClient
    return
}

// 关闭连接
func (s *ServerConn) Close() {
    if s.sshClient != nil {
        _ = s.sshClient.Close()
        s.sshClient = nil
    }
    if s.sftpClient != nil {
        _ = s.sftpClient.Close()
        s.sftpClient = nil
    }
}

// 尝试连接服务器
func (s *ServerConn) TryConnect() (err error) {
    if _, err = s.getSshConnect(); err != nil {
        return
    }
    s.Close()
    return
}

// 在远程服务器执行命令
func (s *ServerConn) RunCmd(cmd string) (output string, err error) {
    var (
        sshClient *ssh.Client
        session   *ssh.Session
    )
    if sshClient, err = s.getSshConnect(); err != nil {
        return
    }
    if session, err = sshClient.NewSession(); err != nil {
        return
    }
    defer session.Close()
    var buf bytes.Buffer
    session.Stdout = &buf
    session.Stdin = &buf
    if err = session.Run(cmd); err != nil {
        err = fmt.Errorf("执行命令失败: %v", err.Error())
        return
    }
    output = buf.String()
    return
}

// 拷贝本机文件到远程服务器
func (s *ServerConn) CopyFile(srcFile, dstFile string) (err error) {
    var (
        client    *sftp.Client
        dstPath   string
        f         *os.File
        w         *sftp.File
        writeSize int64
        fileInfo  os.FileInfo
    )
    if client, err = s.getSftpConnect(); err != nil {
        return
    }
    dstPath = path.Dir(dstFile)
    if _, err = s.RunCmd("mkdir -p " + dstPath); err != nil {
        err = fmt.Errorf("创建目录失败：%v", err)
        return
    }

    if f, err = os.Open(srcFile); err != nil {
        err = fmt.Errorf("打开本地文件失败: %v", err)
        return
    }
    defer f.Close()

    if w, err = client.Create(dstFile); err != nil {
        err = fmt.Errorf("创建文件失败 [%s]: %v", dstFile, err)
        return
    }
    defer w.Close()

    if writeSize, err = io.Copy(w, f); err != nil {
        err = fmt.Errorf("拷贝文件失败: %v", err)
        return
    }

    fileInfo, _ = f.Stat()
    if fileInfo.Size() != writeSize {
        err = fmt.Errorf("写入文件大小错误，源文件大小：%d, 写入大小：%d", fileInfo.Size(), writeSize)
        return
    }
    return
}
