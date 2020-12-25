# 代码发布系统

### 简介
- 支持系统`mac`,`linux`,  未在windows上尝试过
- 使用前后端分离, 前端用`vue`, 后端使用`go`
- 利用`websocket`监听支持实时发布状态监控.
- 支持直发(未测试发布机和目标机是一台的情况),  跳板发.
- 支持全量 & 增量发布.
- 仅供学习交流使用.

界面截图:
![home](./img/home.png)
![servers](./img/servers.png)
![env](./img/env.png)
![task](./img/task.png)


启动方法:
```
数据库初始化:
//1. 创建并导入数据库(略)
//2. 修改数据库配置 config/config.json

后端启动
//1. 利用 go mod 安装库
go mod tidy

//2. 运行测试
go run main.go
或者 :
fresh -c fresh.conf

前端启动:
//1. 进入web目录
//2. 安装依赖
npm install  或者 yarn
//3. 启动
yarn run serve 或者 npm run serve

------------------
//编译好之后发布自行google,  配置文件路径目前写死, 如果有兴趣可以自行修改为命令传参

```
初始用户密码: `admin / 123456`

----
##### 注意
- 因本发布系统里有`git archive`打包, 名称中有中文或者空格会导致打包失败.
- git 配置中文字符不转unicode, 不配置时有中文名文件可能报错
```
git config --global core.quotepath false
```
----
##### 参考项目:
- [gopub](https://github.com/lisijie/gopub)
- [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin)

