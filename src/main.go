package main

import (
	"fmt"
	"guodi/src/guodihttp"
	"net/http"
)

func main() {
	// 初始化所有的客服
	guodihttp.InitCustomer()
	server := http.Server{
		Addr: "192.168.10.150:8100",
	}
	http.HandleFunc("/index", guodihttp.RequestIndex)
	http.HandleFunc("/index/login", guodihttp.RequestLogin)
	http.HandleFunc("/index/register", guodihttp.RequestRegister)
	http.HandleFunc("/index/requestAuthenticeEmail", guodihttp.RequestAuthenticeEmail)
	http.HandleFunc("/index/history", guodihttp.RequestHistory)
	http.HandleFunc("/index/get/email", guodihttp.RequestGetEmail)
	http.HandleFunc("/index/save", guodihttp.RequestSave)
	http.HandleFunc("/index/userfindcustomer", guodihttp.UserFindCustomer)
	http.HandleFunc("/index/customerlogin", guodihttp.CustomerLogin)
	err := http.ListenAndServe("0.0.0.0:8100", nil)
	fmt.Println(err)
	server.ListenAndServe()
}
