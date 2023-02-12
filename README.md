# 直播流媒体服务系统
免费的流媒体服务平台及视觉分析平台
## 简介
1. 包括屏幕管理，摄像头管理
2. 支持国标28181，Onvif协议
3. 支持实况，云台操作和录像功能
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
4. 管理地址http://localhost:3082
## 联系
vx: btagrass
