# esp32-cam
基于esp32-cam的一个项目

## 目录结构
+ 3dmodel
  + 3d模型，用来固定开发板 
+ board
  + 开板相关代码
+ servers
  + 服务器相关代码
    
## 材料
+ ESP32-CAM 开发版
+ USB TO TTL下载器

## 开发环境
+ PlatformIO
  + https://docs.platformio.org/en/latest/core/installation.html
+ CLion 2020.3
  + https://www.jetbrains.com/clion/
  + 安装插件 platformio

## 功能
+ 本地监控 ✅
+ 远程监控 ❌
+ 摄像头旋转 ❌
+ 远程操控 ❌

## 使用
+ clone项目到本地
+ src目录下新建config.h文件，文件内容如下
```c
//配置你的wifi名称和密码，配置2.4G的WIFI
#define WIFI_NAME "xxx"
#define WIFI_PASSWORD "xxx"
```
+ 编译项目
+ 接线
  + ESP32-CAM 引脚  
    ![image1](./image/ESP32-CAM.png)
  + ESP32-CAM 接线  
    ![image2](./image/ESP32-CAM-CONNECTION.png)
+ 上传
+ 查看日志
```
// terminal 执行
platformio device monitor
```
如果看到下面日志，则说明以上流程走通了
```
WiFi connected
Starting web server on port: '80'
Starting stream server on port: '81'
Camera Ready! Use 'http://x.x.x.x' to connect
```
+ 浏览器访问上面日志中的地址 http://x.x.x.x
```
除了查看日志，也可以打开你家路由器的管理后台，看路由器给esp32-arduino分配的是哪个地址
```  
