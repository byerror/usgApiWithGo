package utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"yongyou/global"
)

var myClient *http.Client

// func Execute2(url string) {
// 	http.DefaultClient.Transport = &http.Transport{
// 		TLSClientConfig: &tls.Config{
// 			InsecureSkipVerify: true,
// 			//MaxVersion:         tls.VersionTLS12,
// 			CipherSuites: []uint16{
// 				tls.TLS_AES_128_GCM_SHA256,
// 				tls.TLS_CHACHA20_POLY1305_SHA256,
// 				tls.TLS_AES_256_GCM_SHA384,
// 				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
// 				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
// 				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
// 				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
// 				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
// 				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
// 				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
// 				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
// 				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
// 				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
// 				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
// 				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
// 			},
// 			PreferServerCipherSuites: true,
// 		},
// 	}
// 	// var header *http.Header

//		// http.Header = &header
//		res, err := http.Get(url)
//		if err != nil {
//			fmt.Print("err>>", err.Error())
//			return
//		}
//		defer res.Body.Close()
//		buf := &bytes.Buffer{}
//		buf.ReadFrom(res.Body)
//		fmt.Printf("ret is>>%s", buf.String())
//	}
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

func Execute(url, method string, pData []byte, headers map[string]string) (result interface{}, err error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(pData))
	if err != nil {
		return "NewRequest err>>", err
	}
	if len(headers) > 0 {
		for key, value := range headers {
			request.Header.Add(key, value)
		}
	}

	request.SetBasicAuth(global.Cfg.Api.User, reverse(global.Cfg.Api.Psw))

	response, err := GetClient().Do(request)
	if err != nil {
		return "client.Do err>>", err
	}
	defer response.Body.Close()

	fmt.Println("response code>>", response.StatusCode)
	buff := &bytes.Buffer{}
	_, err = buff.ReadFrom(response.Body)
	if err != nil {
		return "readData err>>", err
	}
	return buff.Bytes(), nil

}

func Execute1xxx(url, method string, pData []byte, headers map[string]string) (result interface{}, err error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(pData))
	if err != nil {
		return "NewRequest err>>", err
	}
	if len(headers) > 0 {
		for key, value := range headers {
			request.Header.Add(key, value)
		}
	}

	// certPool := x509.NewCertPool()
	// bCert, _ := os.ReadFile("demo2.pem")
	// bOk := certPool.AppendCertsFromPEM(bCert)
	// fmt.Println("read cert>>", bOk)
	// items := tls.InsecureCipherSuites()
	// suits := make([]uint16, len(items))
	// for index := range items {
	// 	item := items[index]
	// 	suits[index] = item.ID
	// }

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			MaxVersion:         tls.VersionTLS12, // MinVersion: tls.VersionTLS12,
			// CipherSuites: []uint16{
			// 	//tls.DHE - RSA - AES256 - GCM - SHA384,
			// 	0x009f,

			// 	 tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			// 	 		TLS_DHE_RSA_WITH_AES_256_GCM_SHA384
			// 	//Cipher Suite: TLS_DHE_RSA_WITH_AES_256_GCM_SHA384 (0x009f)
			// 	// 0xc02c, 0xc030, 0x009f, 0xcca9, 0xcca8, 0xccaa, 0xc02b, 0xc02f, 0x009e,
			// 	// 0xc024, 0xc028, 0x006b, 0xc023, 0xc027, 0x0067, 0xc00a, 0xc014, 0x0039,
			// 	// 0xc009, 0xc013, 0x0033, 0x009d, 0x009c, 0x003d, 0x003c, 0x0035, 0x002f, 0x00ff,
			// },

			// RootCAs:            certPool,
			// VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			// 	// f, _ := os.OpenFile("cer2.cer", os.O_CREATE|os.O_WRONLY, 0o777)
			// 	// f.Write(rawCerts[0])
			// 	// f.Close()
			// 	var cert, err = x509.ParseCertificate(rawCerts[0])

			// 	if err != nil {
			// 		fmt.Println("err>>" + err.Error())
			// 	}
			// 	fmt.Println("test>>", cert.NotBefore, cert.NotAfter)
			// 	return nil
			// },
		},
	}}

	response, err := client.Do(request)
	if err != nil {
		return "client.Do err>>", err
	}
	defer response.Body.Close()

	fmt.Println("response code>>", response.StatusCode)
	buff := &bytes.Buffer{}
	_, err = buff.ReadFrom(response.Body)
	if err != nil {
		return "readData err>>", err
	}
	return buff.Bytes(), nil

}

func Executexxx(url, method string, pData []byte, headers map[string]string) (result string, err error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(pData))
	if err != nil {
		return "newRequest出错", err
	}

	if len(headers) > 0 {
		for key, value := range headers {
			request.Header.Add(key, value)
		}
	}
	// var conn *tls.Conn
	// tlsConfig := http.DefaultTransport.(*http.Transport).TLSClientConfig
	tlsConfig := &tls.Config{InsecureSkipVerify: true} //  MaxVersion: tls.VersionTLS12, MinVersion: tls.VersionTLS12,
	// VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	// 	var cert, err = x509.ParseCertificate(rawCerts[0])
	// 	if err != nil {
	// 		fmt.Println("err>>" + err.Error())
	// 	}
	// 	fmt.Println("test>>", cert.NotBefore, cert.NotAfter)
	// 	return nil
	// },
	// VerifyConnection: func(cs tls.ConnectionState) error {
	// 	return nil
	// },

	// tr := http.DefaultTransport.(*http.Transport).Clone()
	// tr.TLSClientConfig = tlsConfig
	// tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {

	// 	conn, err = tls.Dial(network, addr, tlsConfig)
	// 	return conn, err
	// }
	// tr.DialTLS = func(network, addr string) (net.Conn, error) {

	// }

	transport := &http.Transport{TLSClientConfig: tlsConfig} // DialTLS: func(network, addr string) (net.Conn, error) {
	// 	conn, err = tls.Dial(network, addr, tlsConfig)
	// 	return conn, err
	// },

	client := &http.Client{Transport: transport}
	resp, err := client.Do(request)
	// version := conn.ConnectionState().Version
	// fmt.Print(version)
	if err != nil {
		return "请求出错", err
	}
	defer resp.Body.Close()
	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "readFrom出错", err
	}
	return buf.String(), nil
}

// 获取策略状态，查看是否已启用 返回包含 <enable>true</enable>表示成功，<enable>true</enable>为关闭，其它为失败，
func GetPolicy() (isEnable bool, err error) {
	ip := global.Cfg.Api.Ip
	///restconf/data/huawei-security-policy:sec-policy/vsys=public/static-policy/rule=test
	url := fmt.Sprintf("https://%s:8447/restconf/data/huawei-security-policy:sec-policy/vsys=public/static-policy/rule=%s", ip, global.Cfg.Api.Policy.Rule)
	// url := fmt.Sprintf("https://%s:8447/restconf/data/huawei-security-policy:sec-policy/vsys=public/static-policy/rule=167", ip)
	// url = "https://192.168.1.254:8447/restconf/data/huawei-security-policy:sec-policy"

	// utils.Execute2(url)
	// return
	//data := fmt.Sprintf("<rule><source-ip><address-ipv4>192.168.0.167/32</address-ipv4></source-ip><action>true</action><enable>%t</enable></rule>", isEnable)
	headers := make(map[string]string)
	headers["Accept"] = "*/*"
	headers["User-Agent"] = "curl/7.58.0"
	headers["Connection"] = "keep-alive"

	// Cache-Control: max-age=0
	// Authorization: Basic YXBpMTpwYWY2MDg5MDFBQQ==
	headers["Upgrade-Insecure-Requests"] = "1"
	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36"
	headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"
	headers["Accept-Language"] = "zh-CN,zh;q=0.9,en;q=0.8"
	headers["Accept-Encoding"] = "gzip, deflate"

	tmpStr, err := Execute(url, "GET", nil, headers)

	if err != nil {
		fmt.Println(tmpStr, "err>>"+err.Error())
		return false, err
	}
	ret := (string)(tmpStr.([]byte))
	fmt.Println("result is>>", ret)
	if strings.Contains(ret, "<enable>true</enable>") {
		return true, nil
	} else if strings.Contains(ret, "<enable>false</enable>") {
		return false, nil
	}
	return false, fmt.Errorf("无法判断状态 %s", ret)
}

// 开启或关闭策略
func SetPolicy(isEnable bool) (err error) {
	ip := global.Cfg.Api.Ip
	url := fmt.Sprintf("https://%s:8447/restconf/data/huawei-security-policy:sec-policy/vsys=public/static-policy/rule=%s", ip, global.Cfg.Api.Policy.Rule)
	// url = "https://192.168.1.254:8447/restconf/data/huawei-security-policy:sec-policy"

	// utils.Execute2(url)
	// return
	//data := fmt.Sprintf("<rule><source-ip><address-ipv4>192.168.0.167/32</address-ipv4></source-ip><action>false</action><enable>%t</enable></rule>", isEnable)
	data := fmt.Sprintf(global.Cfg.Api.Policy.Template, isEnable)
	headers := make(map[string]string)
	headers["Accept"] = "*/*"
	headers["User-Agent"] = "curl/7.58.0"

	tmpStr, err := Execute(url, "PUT", []byte(data), headers)
	if err != nil {
		fmt.Println(tmpStr, "err>>"+err.Error())
		return err
	}
	b1 := tmpStr.([]byte)
	if len(b1) == 0 { //操作成功不返回消息的
		fmt.Println("操作成功！")
		return nil
	}
	return errors.New(string(b1))

}

// 支持同时检测多个端口，当其中一个端口失败，返回错误
func CheckPorts(ip string, ports []int) (isOpen bool, err error) {
	for _, port := range ports {
		iTimeOut := time.Second
		hostPort := net.JoinHostPort(ip, strconv.Itoa(port))
		conn, err := net.DialTimeout("tcp", hostPort, iTimeOut)
		if err != nil {
			return false, fmt.Errorf("端口%d连接失败!", port)
		}
		defer conn.Close()

	}
	return true, nil
}

func FormatAsDate(t *time.Time) (s string) {

	s = t.Format("2006/01/02 15:04")
	return s
}

// bXml := []byte(`
// <rule>
// <source-ip>
// 	<address-ipv4>192.168.0.167/32</address-ipv4>
// </source-ip>
// <action>true</action>
// <enable>false</enable>
// </rule>
// `)

// 定义函数, 接受一个字符串作为参数, 返回翻转后的字符串.
func reverse(s string) string {
	rns := []rune(s) // 转换为rune类型
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		// 交换字符串中的字母, 如第一个与最后一个交换等等.
		rns[i], rns[j] = rns[j], rns[i]
	}
	// 返回翻转后的字符串.
	return string(rns)
}

// 获取当前执行文件所在目录
func GetExeDir() (path string) {
	fPath, _ := os.Executable()
	return filepath.Dir(fPath)
}
