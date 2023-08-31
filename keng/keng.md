# golang操作华为防火墙ACL实现远程启停

## 一、需求

某个服务端应用A在地点P1，地点P2通过网关建立的VPN来访问P1（A服务器禁止上外网，外网也不能访问A的服务）。这个A服务上线运行后，P2日常访问没什么问题。某天来了个新需求，为了增加A服务的安全性，P2不能再随意访问。要求：

* 指定的管理人员可以在任意地方查看当前服务的状态==开启或关闭==。

* 指定的管理人员可以随时在任意地方开启或关闭服务

## 二、思路

由于P2是在异地通过VPN访问P1,要限制P2访问P1，可以考虑在服务A的服务器上加个防火墙策略。但是了解两边的环境后，这个方案被否了。P1处的网关没有对外除VPN外的任何端口，想要在A上加策略，感觉不好操作。那就考虑在P2的网关做限制吧！刚好P2有对外映射端口，直接拿来做个轻量的web服务，通过访问web页面操作防火墙ACL，就可以达到目的了。

> 为什么不直接开放防火墙，设置个小号，由这个小号去设置ACL?  因为考虑到一方面普通人员去管理一个专业设备，有难度也有风险（上百条ACL被删或突然被修改，风险太大）；另一方面即使日常管理是专业的IT人员，也不敢轻易开放防火墙到公网，出问题要背锅的！

既然是轻量的Web服务，那就用golang+gin来搞吧，做个简单的页面，问题不大。

## 三、准备

首先需要在USG防火墙上建立一条规则，源地址为本地需要访问服务A的地址，目的地为服务A的地址，action为deny。之后通过启停该规则，即可控制对A的访问。一开始想直接go写个客户端通过ssh到防火墙，然后执行一串命令来实现。但是这么做感觉不太好，能执行命令的帐号权限等级需要level 15了吧，那这个密码泄露可能就麻烦了。于是找了下华为的资料，果然有相关的API接口，而且还有两种方式，netconf和restconf。其中netconf是走ssh协议的，restconf走http协议。==当时考虑到ssh协议可能没那么灵活，选了restconf，结果被两坑两天==！

无论选哪种协议，都需要在USG上做好API接口配置：

* 添加API操作员帐号(在aaa下)

  ```shell
  
  manager-user username
    password cipher xxxxx
    service-type api
    level 15
  
  ```

  

* 启用api和restconf接口（在api下）

  ```shell
  api 
    api https enable
   
  ```

## 四、测试

  通过上面的设置后，可以使用curl来测试获取所有安全策略。

  ` curl --tlsv1.2 -k  -v -u uName:psw https://ip:port/restconf/data/huawei-security-policy:sec-policy`

  为了对接api，我用python requests测试了一下，可以正常连接，但是在golang net/http访问时却报错：`remote error: tls: handshake failure` 于是便有了本次踩坑之旅，==下回细说==。对于做运维的小伙伴，这里直接给出解决办法：

```shell
api 
 api https enable  
  security version tlsv1.2  #这个不是必须的
  security cipher-suit all  #如果连接不上，就把所有算法都加上
```

restconf文档里也有说明的：

![image-20230823093305944](http://cdn.qipanet.com/blog/image-20230823093305944.png)





