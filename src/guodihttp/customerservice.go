package guodihttp

import (
	"encoding/json"
	"fmt"
	"guodi/src/guodiredis"
	"guodi/src/guodizap"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// 配置升级程序
var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 存放用户发送的消息
type CustomerMessagee struct {
	Username   string `json:"username"`
	Message    string `json:"message"`
	CustomerID string `json:"customerID"`
}

//存放客服发送的消息
type ClientMessage struct {
	CustomerID string `json:"customerID"`
	Message    string `json:"message"`
	Username   string `json:"username"`
}

// 所有的在线客服
var customer = make(map[string]*websocket.Conn)

// 所有的客户
var clients = make(map[string]*websocket.Conn)

//客服所对应的用户
var cst_clt = make(map[string]string)

// var clients = make(map[cst_clt]*websocket.Conn)

//客服通道
var broadcast = make(chan CustomerMessagee)

// 用户通道
var userbroadcast = make(chan ClientMessage)

// 初始化所有的客服
func InitCustomer() {
	customers := guodiredis.SelectCustomer()
	for _, c := range customers {
		customer[c] = nil
		cst_clt[c] = ""
	}

}

func CustomerLogin(w http.ResponseWriter, r *http.Request) {
	// 获取参数
	r.ParseForm()
	// 获得客服名
	customername := r.FormValue("customername")

	customerID := guodiredis.SelectCustomerID(customername)
	// 检验客服是否处于在线
	if customer[customerID] == nil {
		// 客服不在线

		ws, err := upgrader.Upgrade(w, r, nil)
		customer[customerID] = ws
		if err != nil {
			fmt.Println("连接失败")
			log.Fatal(err)
		}
		loginResult := &RequestResult{
			Status: 200,
			Data: &SuccessResult{
				Success: 1,
				Msg:     "",
			},
		}
		body, err := json.Marshal(loginResult)
		if err != nil {
			guodizap.Errorf("%s", "json转换失败")
		}
		w.Write(body)

		for {
			// 检测被绑定的用户
			// 客服读取消息
			// isbind, clientname := checkBindClient(customerID)

			// if isbind {
			var clientmessage ClientMessage
			//读取用户信息
			guodizap.Infof("%s", customerID+"客服准备读取消息")
			_, msg, _ := customer[customerID].ReadMessage()
			if msg != nil { //空消息不进行发送

				clientmessage.CustomerID = customerID
				clientmessage.Message = string(msg)
				// clientmessage.Username = clientname
				customer[customerID].WriteJSON(clientmessage)
			}
			// userbroadcast <- clientmessage
			// }
			clients[customerID].WriteJSON(clientmessage)
		}

	} else {
		//客服已登陆
		loginResult := &RequestResult{
			Status: 200,
			Data: &SuccessResult{
				Success: 0,
				Msg:     "请不要重复登录",
			},
		}
		body, err := json.Marshal(loginResult)
		if err != nil {
			guodizap.Errorf("%s", "json转换失败")
		}
		w.Write(body)
	}
}
func UserFindCustomer(w http.ResponseWriter, r *http.Request) {
	// 获取参数
	r.ParseForm()
	// 获得用户名
	username := r.FormValue("username")
	// token := r.FormValue("token")
	// 判断是否有在线客服并返回一个客服
	online, customerID := checkCustomer()
	guodizap.Infof("%s", username+"找到了一个在线客服："+customerID)
	if online {
		// 找到了在线用户并且空闲
		ws, err := upgrader.Upgrade(w, r, nil)
		// 用户成功绑定客服
		clients[customerID] = ws
		guodizap.Infof("%s", username+"成功绑定客服")
		//表示客服处于忙碌状态
		cst_clt[customerID] = username
		guodizap.Infof("%s", customerID+"客服正处于忙碌状态")
		if err != nil {
			fmt.Println("连接失败")
			log.Fatal(err)
		}
		// 用户读取消息

		for {
			var customermessagee CustomerMessagee
			//读取客服信息
			guodizap.Infof("%s", customerID+"用户准备读取消息")
			if customer[customerID] == nil {
				guodizap.Infof("%s", customerID+"客服连接出现异常")
			}
			// 用户读取消息
			_, msg, _ := clients[customerID].ReadMessage()
			guodizap.Infof("%s", "绑定"+customerID+"的用户收到消息")

			// _, msg, _ := customer[customerID].ReadMessage()
			// guodizap.Infof("%s", customerID+"收到消息")
			if msg != nil {
				// guodizap.Info(username + "收到消息：" + string(msg))
				customermessagee.CustomerID = customerID
				customermessagee.Message = string(msg)
				customermessagee.Username = username
				clients[customerID].WriteJSON(customermessagee)
			}
			// broadcast <- customermessagee
			customer[customerID].WriteJSON(customermessagee)

		}

	} else {
		// 客服正忙
		loginResult := &RequestResult{
			Status: 200,
			Data: &SuccessResult{
				Success: 0,
				Msg:     "客服正忙,请稍后重试",
			},
		}
		body, err := json.Marshal(loginResult)
		if err != nil {
			guodizap.Errorf("%s", "json转换失败")
		}
		w.Write(body)
	}
}
func checkCustomer() (isonline bool, customerID string) {
	// 循环所有的在线客服
	for customerID, c := range customer {
		if c != nil {
			// 判断客服是否处于空闲等待状态
			if cst_clt[customerID] == "" {
				guodizap.Info("成功找到空闲等待的用户" + customerID)
				return true, customerID
			}
		}
	}
	return false, ""
}
func checkBindClient(customerID string) (isbind bool, clientname string) {
	if cst_clt[customerID] != "" {
		guodizap.Info("该客服已经成功绑定用户:" + cst_clt[customerID])
		return true, cst_clt[customerID]
	}
	// PrintSuccessLog("该客服处于空闲等待状态")
	return false, ""
}
