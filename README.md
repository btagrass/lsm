# 直播流媒体服务系统
免费的基于 [Lal](https://pengrl.com/) 直播流媒体服务管理系统，支持支持国标GB28181（待实现），Onvif协议。
## 简介
### 业务系统
1. 视频墙管理
2. 摄像头管理
    * 实况
    * 云台
    * 录像（待实现）
### 系统设置
1. 资源管理
2. 角色管理
3. 用户管理
4. 部门管理
5. 字典管理
## 准备
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```
## 编译
```bash
bash build.sh
```
## 部署
目前只支持linux平台
1. 安装mysql和redis服务
2. 修改conf/app.yaml配置文件
3. 启动./lsm run
4. [管理地址](http://localhost:3082)
    * 用户名：admin
    * 密码：admin
## 联系
vx: btagrass
