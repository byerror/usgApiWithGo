# golang使用RESTConf操作华为USG防火墙的安全策略 

![image-20230831112552356](http://cdn.qipanet.com/blog/image-20230831112552356.png)

## 一、功能特点
本工具主要用于实现远程控制防火墙安全策略的开启与关闭，方便用户随时做上网行为管控。
* 采用B/S架构，支持PC端和移动端
* 支持安全策略模板设置
* 可随时查看策略状态
* 提供简单的操作记录报表
* windows以系统服务运行、开机自启

## 二、程序清单
程序采用gin框架，数据库用的是sqlite,日志用的是zap。
前端主要用到了weui+以及mescroll组件

## 三、程序运行
双击main.exe，自动创建系统服务，默认端口9999。
==也可以使用./main.exe install来手动安装服务==
安装服务后程序将自动启动，http://localhost:9999可看到效果。

## 四、配置文件
==安全策略模板请自行参考华为官方文件==
```yaml
server:
  #用于配置zap日志的模式
  debug: true
  #gin服务的端口
  port: 9999
database:
  path: "E:\\case\\go\\src\\yongyou\\db.db"

api: #restConf Api
  #防火墙ip
  ip: "192.168.1.1"
  user: "test"
  psw: "123fvv"
  policy:
    rule: "test"
    # 策略模板 %t将会被替换为true,false  
    template: "<rule><desc>test</desc><source-ip><address-ipv4>192.168.1.167/32</address-ipv4></source-ip><destination-ip><address-ipv4>172.16.0.223/32</address-ipv4></destination-ip><action>false</action><enable>%t</enable></rule>"


monitor:
  # 要监控的服务器IP
  IP: "172.16.0.223"
  port: 8088

```
## 五、其他
* 源码：https://github.com/root6819/usgApiWithGo
* 官网： http://www.qipanet.com 
* 微信 root6819 | QQ 302777528