# golang net/http tls客户端频繁请求 401 Unauthorized

## 一、烟零蛋

 其实这个问题是《golang操作华为防火墙ACL实现远程启停》里遇到的，项目做完准备交付，测试偶尔遇到前端一直转圈（前后端不分离项目）。经过定位发现卡在了获取防火墙策略状态那里，再fiddler抓包看到了请求发出后一直没数据返回，加上默认http.Client没设置超时处理，导致后端也无法返回数据。

为什么restConf服务不返回数据呢？ 用fiddler直接模拟请求，也卡住了，根本没返回消息;再用python来请求，也是没任何数据响应；最后重启restConf服务，`api https enable`，问题暂时解决，但是前端多刷新几次后，又出现卡顿了。

<!--more-->

## 二、渐露头角

 上面的问题不知啥情况，后面又复现不了了，但又出现了新的问题。golang项目里，只要多发起几次请求，请会出现401 Unauthorized报错。代码大概如下：

```go
func Execute(method string,url string){
	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		MaxVersion:         tls.VersionTLS12,
	}
	tr := &http.Transport{

		TLSClientConfig: tlsCfg,
	}
	client := &http.Client{
		Transport: tr,
	}
	request, _ := http.NewRequest(method, url, nil)
    client.Do(request)
	...
}
```

经过多次反复对比测试和抓包，发现了个奇怪的现象：golang项目里刷新几次不行后（401 Unauthorized），不管再用哪种语言写的客户端，都报401错误。更让人费解的是，除了重启restConf服务外，通过防火墙后台把密码用同一密码重置后，再请求，竟然又恢复成功了。此外，当出现401时，我在后台再新添加个restConf的认证帐号，也能成功请求，但换回前一个帐号，又不行了。

难道是防火墙做了api调用限制？ 想到这我立即用curl 不断请求，发现n次都正常，用python请求也正常，但换回golang才请求几次就出现401了，这下确认了是golang的问题，跟防火墙没关。于是一行行代码分析着，尽量将代码精简，只保留关键代码。起一个for循环测试，才几次就报401了。难道是golang太快了？于是加上time.Sleep 延时到5秒请求一次，再测试，一样的现象。从抓包的情况来看，正常请求和异常请求发出的数据没任何异常情况，搞不懂是啥回事。

再回看代码，感觉每次请求都创建client有点费劲，就把client提取出来，只初始化一次，再跑for循环请求测试，竟然成功了，即使没加time.sleep，也能正常返回200!

```go
func GetClient() (client *http.Client) {
	if myClient != nil {
		return myClient
	} else {

		tlsCfg := &tls.Config{
			InsecureSkipVerify: true,
			MaxVersion:         tls.VersionTLS12,
		}
		tr := &http.Transport{

			TLSClientConfig: tlsCfg,
		}
		client = &http.Client{
			Transport: tr,
		}
		myClient = client
		return client
	}
}
```

再google看下有没说明，刚输入golang http.client，后面就自动帮我补上==复用==两字了，看来有戏`https://zhuanlan.zhihu.com/p/474206147`。大概意思是http.Client会自动缓存连接，提高效率，本着不请甚解的精神 ，我默默的关闭了页面。

后来，当出现401后，我随手将golang项目停止服务后，fiddler模拟请求正常，curl、python都正常了，也就是golang程序占着茅坑导致其它应用同一个号都访问不了。这再次验证了所谓http.client复用的问题？

其实golang 的http.Client里就有说明，只怪我读书少，没见识:cry:

>// The Client's Transport typically has internal state (cached TCP
>
>// connections), so Clients should be reused instead of created as
>
>// needed. Clients are safe for concurrent use by multiple goroutines.
>
>



![image-20230829103747302](http://cdn.qipanet.com/blog/image-20230829110818098.png)

![image-20230829110818098](http://cdn.qipanet.com/blog/image-20230829103747302.png)

## 三、总结 

感觉如果golang踩坑太少，当用在关键项目上，出现问题后那场景真不堪设想！ 这篇文章写的比较乱，各位将就看吧！