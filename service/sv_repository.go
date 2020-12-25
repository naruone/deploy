package service

import (
    "deploy/config"
    "deploy/model"
    "deploy/utils"
    "errors"
    "strings"
    "sync"
)

var (
    Repos = make(map[uint]*Repository)
    rLock sync.Mutex
)

type Repository struct {
    Path        string
    PackagePath string
    Project     *model.Project
    IsProcess   bool
    cLock       sync.Mutex
}

//删除项目时删除缓存
func DelRepository(projectId uint) {
    if Repos[projectId] != nil {
        delete(Repos, projectId)
    }
}

func GetRepository(project *model.Project) *Repository {
    rLock.Lock()
    defer rLock.Unlock()

    _repoPath := strings.TrimRight(config.GConfig.Repo, "/")
    if Repos[project.ProjectId] == nil {
        _packagePath := strings.TrimRight(config.GConfig.RepoPackage, "/")
        _ = utils.CreateDir(_repoPath)
        _ = utils.CreateDir(_packagePath)
        Repos[project.ProjectId] = &Repository{
            Path:        _repoPath + "/" + strings.TrimLeft(project.Dst, "/"),
            Project:     project,
            PackagePath: _packagePath + "/",
            IsProcess:   false,
        }
    } else {
        Repos[project.ProjectId].Project = project
        Repos[project.ProjectId].Path = _repoPath + "/" + strings.TrimLeft(project.Dst, "/")
    }
    return Repos[project.ProjectId]
}

//获取所有分支
func (repo *Repository) GetBranches() (result []string, err error) {
    var out string
    result = []string{}
    _, _, _ = utils.RunCmd(repo.Path, "git", "fetch", "-pq")
    if out, _, err = utils.RunCmd(repo.Path, "git", "branch", "-r"); err != nil {
        return
    }
    res := strings.Split(out, "\n")
    for _, v := range res {
        v = strings.TrimSpace(v)
        if v != "" && strings.Index(v, "origin/HEAD ->") == -1 {
            result = append(result, strings.TrimLeft(v, "origin/"))
        }
    }
    return
}

//获取当前分支版本信息
func (repo *Repository) GetVersions(branch string) (result []model.CsvVersion, err error) {
    var (
        out  string
        _res []model.CsvVersion
    )
    repo.cLock.Lock()
    defer repo.cLock.Unlock()
    if _, _, err = utils.RunCmd(repo.Path, "git", "reset", "--hard"); err != nil {
        return
    }
    if _, _, err = utils.RunCmd(repo.Path, "git", "checkout", branch); err != nil {
        return
    }
    if _, _, err = utils.RunCmd(repo.Path, "git", "pull"); err != nil {
        return
    }
    if out, _, err = utils.RunCmd(repo.Path, "git", "log", "-20",
        "--pretty=format:%h^_^%cd %cn: %s", "--date=iso"); err != nil {
        return
    }
    res := strings.Split(out, "\n")
    for _, v := range res {
        v = strings.TrimSpace(v)
        _t := strings.Split(v, "^_^")
        if v != "" {
            _res = append(_res, model.CsvVersion{
                Version: _t[0],
                Message: strings.Replace(v, "^_^", " ", 1),
            })
        }
    }
    result = _res
    return
}

func (repo *Repository) Package(startVer, endVer, name string) (filename string, delFiles []string, err error) {
    var (
        cmd         string
        _delFileLog string
        errOutput   string
    )
    filename = repo.PackagePath + name
    if startVer == "" {
        cmd = "git archive --format=tar.gz " + endVer + " -o " + filename
    } else {
        cmd = "git archive --format=tar.gz " + endVer + " $(git diff --name-only --diff-filter=ACMRTUX " + startVer + " " + endVer + ") -o " + filename
        _delFileLog, _, err = utils.RunCmd(repo.Path, "/bin/bash", "-c", "git diff --name-only --diff-filter=D "+startVer+" "+endVer)
        var _deleteFiles []string
        _delFileLog = strings.TrimSpace(_delFileLog)
        if _delFileLog != "" {
            for _, v := range strings.Split(_delFileLog, "\n") {
                _deleteFiles = append(_deleteFiles, strings.TrimSpace(v))
            }
        }
        //文件改名: 打包改名后的文件,  删除改名前的文件
        //git diff --name-status --diff-filter=R 543ccb059 467a806cf | awk '{print $2}'
        _delFileLog, _, err = utils.RunCmd(repo.Path, "/bin/bash", "-c", "git diff --name-status --diff-filter=R "+startVer+" "+endVer+"| awk -F'\t' '{print $2}'")
        if _delFileLog != "" {
            for _, v := range strings.Split(_delFileLog, "\n") {
                _deleteFiles = append(_deleteFiles, strings.TrimSpace(v))
            }
        }
        delFiles = _deleteFiles
    }
    if _, errOutput, err = utils.RunCmd(repo.Path, "/bin/bash", "-c", cmd); err != nil {
        err = errors.New(err.Error() + "\n Output: " + errOutput + "\n Cmd: " + cmd + " \n 可能原因: git archive 不允许文件名中存在空格, 请检查!")
    }
    return
}

//初始化克隆项目
func (repo *Repository) CloneRepo() (errOut string, err error, processing bool) {
    repo.cLock.Lock()
    defer repo.cLock.Unlock()
    if repo.IsProcess {
        err = errors.New("已开始初始化, 请稍后查看状态")
        processing = true
        return
    }
    repo.IsProcess = true
    _, errOut, err = utils.RunCmd(config.GConfig.Repo, "git", "clone", repo.Project.RepoUrl, repo.Project.Dst)
    repo.IsProcess = false
    return
}
