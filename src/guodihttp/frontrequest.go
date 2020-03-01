package guodihttp

import (
	"encoding/json"
	"fmt"
	"guodi/src/guodiauthenticate"
	"guodi/src/guodicommon"
	"guodi/src/guodiredis"
	guodisql "guodi/src/guodisql"
	"guodi/src/guodizap"
	"net/http"
	"strconv"
)

type DateInteface interface {
}
type RequestResult struct {
	Status int64        `json:"status"`
	Data   DateInteface `json:"data"`
}
type SuccessResult struct {
	Success int64  `json:"success"`
	Msg     string `json:"msg"`
}
type LoginResult struct {
	Succeess int64  `json:"success"`
	Msg      string `json:"msg"`
	UserID   string `json:"userID"`
	Token    string `json:"token"`
}

type HistoryResult struct {
	Status int64                   `json:"status"`
	Data   []guodisql.HistoryOrder `json:"data"`
}

func successMsg(msg string, status int64, success int64) *RequestResult {
	requestResult := &RequestResult{
		Status: status,
		Data: &SuccessResult{
			Success: success,
			Msg:     msg,
		},
	}
	return requestResult
}

/**
 * @description: 用户进入首页
 * @param {type}
 * @return:
 */
func RequestIndex(w http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()
	// 跨域解决
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
	// 获得请求地址ip
	requestIP := r.RemoteAddr

	// 获取参数
	token := r.FormValue("token")
	email := r.FormValue("email")
	if len(token) > 0 && len(email) > 0 {
		// 开始token验证
		// 应该调用的函数guodiredis.CheckToken(email,token)
		if guodiredis.CheckToken(email, token) {
			indexResult := successMsg("", 200, 1)

			body, err := json.Marshal(indexResult)
			if err != nil {
				guodizap.Errorf("无法正常转换json", err)
			}
			w.Write(body)
			guodizap.Infof("%s", requestIP+"成功登录")
		} else {
			indexResult := successMsg("身份过期,请重新登录", 200, 0)

			body, err := json.Marshal(indexResult)
			if err != nil {
				guodizap.Errorf("无法正常转换json", err)
			}
			w.Write(body)
			guodizap.Infof(requestIP + "进入游客模式")
		}

	} else {
		indexResult := successMsg("游客模式", 200, 0)

		body, err := json.Marshal(indexResult)
		if err != nil {
			guodizap.Errorf("无法正常转换json", err)
		}
		w.Write(body)
		guodizap.Infof(requestIP + "进入游客模式")
	}
}

/**
 * @description: 用户请求登录
 * @param {type}
 * @return:
 */
func RequestLogin(w http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()
	// 跨域解决
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
	// 获得请求地址ip
	requestIP := r.RemoteAddr
	appID := r.FormValue("appID")
	md5ID := r.FormValue("md5ID")
	date := r.FormValue("date")
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	fmt.Println(appID, md5ID, date, email, password)
	// 判断参数是否有空值
	if len(appID) > 0 && len(md5ID) > 0 && len(date) > 0 && len(email) > 0 && len(password) > 0 {
		// 开始验证前端调用前端验证模块
		// func FrontAuthentice(appID string , md5ID string , date string)bool{}
		if guodiauthenticate.FrontAuthentice(appID, md5ID, date) {
			// 开始验证用户信息
			// func CheckUserByEmailAndPassword(email string,password string)(isexist bool,userID string){}

			isexist, userID := guodisql.CheckUserByEmailAndPassword(email, password)

			if isexist && len(userID) > 0 {
				// 登录成功开始调用公共服务为登录生成token、func GetFontToken() string{} guodicommon
				token := guodicommon.GetFontToken()
				// token := "asdfgjikolji"
				if len(token) > 0 {
					// 成功生成token
					loginResult := &RequestResult{
						Status: 200,
						Data: &LoginResult{
							Succeess: 1,
							Msg:      "",
							UserID:   userID,
							Token:    token,
						},
					}
					body, err := json.Marshal(loginResult)
					if err != nil {
						guodizap.Errorf("无法正常转换json", err)
					}
					w.Write(body)
					guodizap.Infof(requestIP + email + "用户成功登录")
					guodiredis.SaveToken(email, token)
					//
				} else {
					loginResult := &RequestResult{
						Status: 500,
						Data: &LoginResult{
							Succeess: 0,
							Msg:      "请求失败",
							UserID:   "",
							Token:    "",
						},
					}
					body, err := json.Marshal(loginResult)
					if err != nil {
						guodizap.Errorf("无法正常转换json", err)
					}
					w.Write(body)
					guodizap.Infof(requestIP + email + "用户登录时token值出现了空值，公共服务出现异常，无效的请求")
				}
			} else {
				loginResult := &RequestResult{
					Status: 200,
					Data: &LoginResult{
						Succeess: 0,
						Msg:      "用户名或密码错误",
						UserID:   "",
						Token:    "",
					},
				}
				body, err := json.Marshal(loginResult)
				if err != nil {
					guodizap.Errorf("无法正常转换json", err)
				}
				w.Write(body)
				guodizap.Infof("%s", email+"用户登录时用户名或密码出现了错误")
			}
		} else {
			loginResult := &RequestResult{
				Status: 500,
				Data: &LoginResult{
					Succeess: 0,
					Msg:      "请求失败",
					UserID:   "",
					Token:    "",
				},
			}
			body, err := json.Marshal(loginResult)
			if err != nil {
				guodizap.Errorf("无法正常转换json", err)
			}
			w.Write(body)
			guodizap.Infof(requestIP + "无效的请求")
		}
	} else {
		loginResult := &RequestResult{
			Status: 500,
			Data: &LoginResult{
				Succeess: 0,
				Msg:      "请求失败",
				UserID:   "",
				Token:    "",
			},
		}
		body, err := json.Marshal(loginResult)
		if err != nil {
			guodizap.Errorf("无法正常转换json", err)
		}
		w.Write(body)
		guodizap.Infof(requestIP + "无效的请求")
	}
}

/**
 * @description: 用户注册
 * @param {type}
 * @return:
 */
func RequestRegister(w http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()
	// 跨域解决
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
	// 获得请求地址ip
	requestIP := r.RemoteAddr
	// 获取参数
	appID := r.FormValue("appID")
	md5ID := r.FormValue("md5ID")
	date := r.FormValue("date")
	email := r.PostFormValue("email")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	authenticate := r.PostFormValue("authenticate")
	if len(appID) > 0 && len(md5ID) > 0 && len(date) > 0 && len(email) > 0 && len(password) > 0 && len(username) > 0 {
		// 开始前端验证
		// 开始验证前端调用前端验证模块
		// func FrontAuthentice(appID string , md5ID string , date string)bool{}
		if guodiauthenticate.FrontAuthentice(appID, md5ID, date) {
			// 调用redis相关服务检验验证码是否正确
			// unc CheckEmailTempAuthentic(email string,authentic string)bool{}
			// guodiredis.CheckEmailTempAuthentic(email,authenticate)
			if guodiredis.CheckEmailTempAuthentic(email, authenticate) {
				// 开始保存用户信息调用mysql相关服务
				issuccess := guodisql.SaveUser(email, username, password)
				// issuccess := true
				if issuccess {
					registerResult := successMsg("", 200, 1)

					body, err := json.Marshal(registerResult)
					if err != nil {
						guodizap.Errorf("%s", "json转换失败", err)
					}
					w.Write(body)
					guodizap.Infof("%s", email+"用户注册成功")
				} else {
					registerResult := successMsg("请求失败", 500, 0)

					body, err := json.Marshal(registerResult)
					if err != nil {
						guodizap.Errorf("%s", "json转换失败", err)
					}
					w.Write(body)
					guodizap.Errorf("%s", requestIP+"在进行注册的时候数据库存储不通过")
				}
			} else {
				registerResult := successMsg("验证码失效请重试", 200, 0)

				body, err := json.Marshal(registerResult)
				if err != nil {
					guodizap.Errorf("%s", "json转换失败", err)
				}
				w.Write(body)
				guodizap.Errorf("%s", requestIP+authenticate+"在进行注册的时候验证码不通过")
			}
		} else {
			registerResult := successMsg("请求失败", 500, 0)

			body, err := json.Marshal(registerResult)
			if err != nil {
				guodizap.Errorf("%s", "json转换失败", err)
			}
			w.Write(body)
			guodizap.Errorf("%s", requestIP+"在进行注册的时候前端验证不通过")
		}
	}
}

/**
 * @description: 用户查询历史记录
 * @param {type}
 * @return:
 */
func RequestHistory(w http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()
	// 跨域解决
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
	// 获得请求地址ip
	requestIP := r.RemoteAddr
	// 获取参数
	userID := r.FormValue("userID")
	if len(userID) > 0 {
		// 调用mysql相关服务、验证userID是否存在

		if guodisql.CheckUserId(userID) {
			// 调用查询历史记录
			historyorder := guodisql.SelectHistoryByID(userID)
			// var historyorder []HistoryOrder
			historyResult := &HistoryResult{
				Status: 200,
				Data:   historyorder,
			}
			body, err := json.Marshal(historyResult)
			if err != nil {
				guodizap.Infof("%s", requestIP+"json转换出现错误")
			}
			w.Write(body)
			guodizap.Infof("%s", userID+"："+requestIP+"查询了一次历史记录")
		} else {
			historyResult := successMsg("请求失败", 500, 0)

			body, err := json.Marshal(historyResult)
			if err != nil {
				guodizap.Errorf("%s", "json转换失败", err)
			}
			w.Write(body)
			guodizap.Errorf("%s", requestIP+"查询历史订单出现错误")
		}
	} else {
		historyResult := successMsg("请求失败", 500, 0)

		body, err := json.Marshal(historyResult)
		if err != nil {
			guodizap.Errorf("%s", "json转换失败", err)
		}
		w.Write(body)
		guodizap.Errorf("%s", requestIP+"查询历史订单出现错误")
	}
}

/**
 * @description: 保存快递信息
 * @param {type}
 * @return:
 */
func RequestSave(w http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()
	// 跨域解决
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
	// 获得请求地址ip
	requestIP := r.RemoteAddr
	userID := r.Form["userID"][0]
	orderresult := r.Form["orderresult"][0]
	orderID := r.Form["orderID"][0]
	//string类型转换为int64
	i, err := strconv.ParseInt(orderresult, 10, 64)
	if err != nil {
		guodizap.Errorf("%s", err)
	}
	if guodisql.SaveOrder(orderID, i, userID) {
		saveResult := successMsg("保存快递信息成功", 200, 1)
		body, err := json.Marshal(saveResult)
		if err != nil {
			guodizap.Infof("%s", requestIP+"保存快递信息成功")
		}
		w.Write(body)
	} else {
		saveResult := successMsg("保存快递信息失败", 500, 0)
		body, err := json.Marshal(saveResult)
		if err != nil {
			guodizap.Infof("%s", requestIP+"保存快递信息成功")
		}
		w.Write(body)
	}
}
func RequestGetEmail(w http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()
	// 跨域解决
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
	// 获得请求地址ip
	requestIP := r.RemoteAddr
	// email := r.Form["email"][0]
	// result := fmt.Sprintf("%s", email)
	// exist:=guodiemail.SendEmail(result)
	if true {
		emailResult := successMsg("", 200, 1)
		body, err := json.Marshal(emailResult)
		if err != nil {
			guodizap.Infof("%s", requestIP+"获取验证码成功")
		}
		w.Write(body)
	} else {
		emailResult := successMsg("获取验证码失败", 200, 0)
		body, err := json.Marshal(emailResult)
		if err != nil {
			guodizap.Infof("%s", requestIP+"获取验证码失败")
		}
		w.Write(body)
	}
}
func RequestAuthenticeEmail(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// 跨域解决
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
	// 获得请求地址ip
	requestIP := r.RemoteAddr
	// 获取参数
	email := r.PostFormValue("email")
	if len(email) > 0 {
		// 调用sql验证是否已经注册
		if guodisql.CheckUserEmail(email) {
			historyResult := successMsg("", 200, 1)
			body, err := json.Marshal(historyResult)
			if err != nil {
				guodizap.Infof("%s", requestIP+"json转换失败")
			}
			w.Write(body)
		} else {
			historyResult := successMsg("邮箱已经被注册", 200, 0)
			body, err := json.Marshal(historyResult)
			if err != nil {
				guodizap.Infof("%s", requestIP+"json转换失败")
			}
			w.Write(body)
		}
	} else {
		historyResult := successMsg("邮箱不能为空", 500, 0)
		body, err := json.Marshal(historyResult)
		if err != nil {
			guodizap.Infof("%s", requestIP+"json转换失败")
		}
		w.Write(body)
	}
}
