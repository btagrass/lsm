# [直播流媒体服务系统](https://github.com/btagrass/lsm)
免费的基于 [Lal](https://pengrl.com/) 直播流媒体服务管理系统，支持Onvif协议，海康Isc平台（开发中），华为S200平台（开发中）和国标GB28181（待开发）。
## 简介
### 业务系统
1. 视频墙管理
2. 摄像头管理
    * 同步
    * 实况
    * 录像
    * 快照
    * 云台
    * 预置位
3. 视频管理
    * 虚拟流
4. 实时流管理
### 系统设置
1. 资源管理
2. 角色管理
3. 用户管理
4. 部门管理
5. 字典管理
## 编译
```bash
bash build.sh
```
## 部署
目前只支持linux平台
1. 安装mysql
2. 修改conf/app.yaml配置文件
3. 启动./lsm run
4. [管理地址](http://localhost:3082)
    * 用户名：admin
    * 密码：admin
## 联系
vx: btagrass
