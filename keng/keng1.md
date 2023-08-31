# golang操作华为防火墙ACL实现远程启停2

## 一、摸索

上文提到用go开发的程序对接防火墙api时出现`remote error: tls: handshake failure`，虽然通过调低防火墙api security加密算法强度解决了，但是总感觉怪怪的。因为用python和curl 简单测试了下，都能正常访问，难道golang还落后了？

为了搞清楚原因，先是设置了InsecureSkipVerify,测试没效果；再指定tls版本`MaxVersion:     tls.VersionTLS12, MinVersion: tls.VersionTLS12`，也没效果；改写VerifyPeerCertificate ，直接return nil，测试还是没效果。

怀疑go需要手动指定证书，把cer转成pem后==（用openssl可以转）==，加载到RootCAs，测试还是不行。

<!--more-->

```go
// certPool := x509.NewCertPool()
//pem文件转换  openssl x509 -inform der -in cer2.cer -out demo2.pem
	// bCert, _ := os.ReadFile("demo2.pem")
	// bOk := certPool.AppendCertsFromPEM(bCert)
client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
            MaxVersion:         tls.VersionTLS12, MinVersion: tls.VersionTLS12
            VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
                return nil
            }，
            //RootCAs：certPool
            
        }，
    

}
```
## 二、意外的成功

正当一愁莫展时，随手打开了fiddler抓包工具，然后抓了python demo程序发起的https请求，数据正常的返回，看不到什么特别的请求参数。随后改了防火墙的一些设置，还有go项目里的代码也不知改了哪里，反正再次run时，竟然成功了！

幸福来的太突然了，我都不知咋回事？很好奇是哪一行代码写错了，于是我把代码重写了一遍，只启用了`InsecureSkipVerify`，其他参数都没加上，测试也成功了。感觉不对劲，干脆把防火墙也重启了==（后面改动的配置没保存）==，此后再测试，请求仍然成功，感觉非常意外，但又找不到具体原因。夜已深，只好先睡。。。

## 三、罪魁祸首竟是它

第二天起来再测试，发现仍然能正常访问，心里有点慌，这神操作有点懵啊。于是再重新建一个测试项目，代码又重写了一遍，好家伙，怎么折腾都行了，此前可是怎么折腾都不行呢！正晕乎乎的不知所以时，无意间点开了fiddler，发现fiddler一直开着的，于是随手把它关了。然后继续测试代码，奇迹出现了，这回竟然出现了久违的` tls: handshake failure`报错。

此时终于知道为什么了，使用fiddler抓包，相当于设置了代理，而且fiddler能自动信任证书，所以开启抓包后go的请求能正常访问api接口。排查工作似乎又回到了原点，好在这时我已经深深的怀疑golang的net/http包有问题。

## 四、柳暗花明

即然fiddler也分析不了TLS协议握手过程，那就上Wireshark吧。启动抓包，先golang程序发起请求，然后python的程序请求,Wiresh地ark里过滤感兴趣的数据包` tcp.port==xxxx && ip.dst=xxxx`。对比两者发起的Client hello发现golang的cipher suites比python少了很多算法。

![image-20230824004348171](http://cdn.qipanet.com/blog/image-20230824004348171.png)

![image-20230824005435130](http://cdn.qipanet.com/blog/image-20230824010700069.png)

会不会是某个算法漏了，导致加密方式不一致呢？于是我在`tls.Config`里加上`CipherSuites`,把wireshare里抓取到的加密算法对应的常量值都放进数组里，测试后发现还是老样子的报错。心里很不甘心，用wireshark抓包，发现即使手动加上了28个算法常量，请求时仍然只有16个。

```go
TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MaxVersion:         tls.VersionTLS12, MinVersion: tls.VersionTLS12,
    		CipherSuites: []uint16{ .....}
}
```

到底少了哪些算法呢？我在tls包找到cipher_suites.go源码，tls加密用到的算法都在里面了，但是仔细数了下，发现确实比python的少了很多。经过一个个对比，发现少的部分算法都是TLS_DHE开头的。百度后了解到DHE算法效率低且比较占用资源，在github上也有人建议增加dhe支持，但是没被官方支持！

解决办法是用别人写的https://github.com/mordyovits/golang-crypto-tls，具体我没测试，对于我来说，我只是想找出真相而已，不想继续在这小项目上浪费太多时间了。



![image-20230824010700069](http://cdn.qipanet.com/blog/image-20230824005435130.png)

>
>
># Transport Layer Security (TLS) Parameters 传输层安全TLS参数  https://www.iana.org/assignments/tls-parameters/tls-parameters.xml
>
>
>
>

 