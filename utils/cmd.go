package utils

import (
    "bytes"
    "os/exec"
)

func RunCmd(dir, cmdName string, args ...string) (out string, errOut string, err error) {
    bufOut := new(bytes.Buffer)
    bufErr := new(bytes.Buffer)

    cmd := exec.Command(cmdName, args...)
    cmd.Dir = dir
    cmd.Stdout = bufOut
    cmd.Stderr = bufErr

    err = cmd.Run()
    out = string(bufOut.Bytes())
    errOut = string(bufErr.Bytes())
    return
}
