package guodihttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("该测试被跳过")
	}
	// 创建一个用于运行测试的多路复用器
	mux := http.NewServeMux()
	// 绑定测试的处理器
	mux.HandleFunc("/post", RequestIndex)
	// 创建一个记录器，用于获取服务器返回的HTTP响应
	writer := httptest.NewRecorder()

	request, _ := http.NewRequest("GET", "/post?email=sdfasd&&token=sdfaf", nil)
	//向被测试处理器发送请求
	mux.ServeHTTP(writer, request)
	if writer.Code != 200 {
		t.Errorf("response code id %v", writer.Code)
	}
	re := new(RequestResult)
	json.Unmarshal(writer.Body.Bytes(), re)

	fmt.Println(re)

}

func TestLogin(t *testing.T) {
	if testing.Short() {
		t.Skip("该测试被跳过")
	}
	// 创建一个用于运行测试的多路复用器
	mux := http.NewServeMux()
	// 绑定测试的处理器
	mux.HandleFunc("/login", RequestLogin)

	writer := httptest.NewRecorder()

	request, _ := http.NewRequest("POST", "/login?appID=asdf&&md5ID=sadfas&&date=sdfsd", nil)
	// request.Header.Set("Content-Type", writer.FormDataContentType())
	//向被测试处理器发送请求
	mux.ServeHTTP(writer, request)

	re := new(RequestResult)
	json.Unmarshal(writer.Body.Bytes(), re)

	fmt.Println(re)
}
