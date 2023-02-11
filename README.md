# 直播流媒体服务系统
免费的流媒体服务平台及视觉分析平台
## 简介
1. 包括屏幕管理，摄像头管理
2. 支持国标28181，Onvif协议
3. 支持实况，云台操作和录像功能
## 安装
目前只支持linux平台
1. 安装mysql和redis服务
2. 修改conf/app.yaml配置文件
3. 启动./lsm run
4. 管理地址http://localhost:3082
## 国标
### 编码
#### 规则（20位）
省级（2位）+ 市级（2位）+ 区级（2位）+ 单位（2位）+ 行业（2位）+ 类型（3位）+ 网络（1位）+ 序号（6位），详见《GBT 28181-2016 公共安全视频监控联网系统信息传输、交换、控制技术要求.pdf》P88 附录D.1 编码规则A
#### 示例
| 服务器域 | 服务器编码 | 摄像头1编码 | 摄像头2编码 |
| :-- | --- | --- | --- |
| 2101120190 | 21011201902000000001 | 21011201901320000001 | 21011201901320000002 |
## 录像
容量(G) = 码流 * 60 * 60 * 24 / 8 / 1024 / 1024

# 部署
## Docker
### Mysql
```
mkdir -p /home/mysql/data
docker rm -f mysql
docker run -d -p 3306:3306 -v /etc/localtime:/etc/localtime -v /home/mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 --name=mysql --privileged=true --restart=unless-stopped mysql:5.7
```
### Redis
```
mkdir -p /home/redis/{conf,data}
docker rm -f redis
docker run -d -p 6379:6379 -v /etc/localtime:/etc/localtime -v /home/redis/data:/data --name=redis --privileged=true --restart=unless-stopped redis:6-alpine --requirepass 123456
```
### Lal
```
mkdir -p /home/lal/conf
docker rm -f lal
docker run -d -p 1935:1935 -p 8080:8080 -p 4433:4433 -p 5544:5544 -p 8883:8083 -p 8084:8084 -p 30000-30100:30000-30100/udp -v /etc/localtime:/etc/localtime -v /home/lal/conf/lalserver.conf.json:/lal/conf/lalserver.conf.json --name=lal --log-opt=max-size=100m q191201771/lal
```
### Zlm
```
mkdir -p /home/zlm/conf
docker run -d -p 1935:1935 -p 8080:80 -p 8443:443 -p 8554:554 -p 10000:10000 -p 10000:10000/udp -p 8000:8000/udp -p 9000:9000/udp -p 30000-30100:30000-30100/udp -v /etc/localtime:/etc/localtime -v /home/zlm/conf/config.ini:/opt/media/conf/config.ini --name=zlm --privileged=true --restart=unless-stopped --log-opt=max-size=100m zlmediakit/zlmediakit:master
```
### lsm
```
mkdir -p /home/lsm
bash/sh svc.sh -h
```
## Swarm
### Nvidia Docker
#### /etc/docker/daemon.json
```
{
  "runtimes": {
    "nvidia": {
      "path": "/usr/bin/nvidia-container-runtime",
      "runtimeArgs": []
    }
  },
  "default-runtime": "nvidia",
  "node-generic-resources": [
    "NVIDIA-GPU=GPU-a829aa4d-5766-8713-c2c9-a2e7a4662298"
  ]
}
```
#### /etc/nvidia-container-runtime/config.toml
```
swarm-resource = "DOCKER_RESOURCE_GPU"-
```
### 初始化
```
docker swarm init
```
## 准备
```bash
go install github.com/swaggo/swag/cmd/swag@v1.8.7
```
## 编译
```bash
bash build.sh
```
