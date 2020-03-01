package sqlservice

import (
	"database/sql"
	"fmt"
	"guodi/src/guodicommon"

	_ "github.com/go-sql-driver/mysql"
)

type HistoryOrder struct {
	OrderID     string
	Orderresult int64
}

var (
	DB_connect_string string = "root:123456@tcp(192.168.10.252:3306)/guodi?charset=utf8"
)

func CheckUserByEmailAndPassword(email string, password string) (isexist bool, userID string) {
	db, err := sql.Open("mysql", DB_connect_string)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false, err.Error()
	}
	//order
	rows, err := db.Query("SELECT userid FROM user WHERE email=? and password=?", email, password)
	if rows.Next() == false {
		return false, ""
	} else {
		//var userid string
		rows.Scan(&userID)
		return true, userID
	}
}
func SaveUser(email string, username string, password string) bool {
	db, err := sql.Open("mysql", DB_connect_string)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	// randnum := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000000000)
	userID := guodicommon.GetRandomStringID()
	_, err2 := db.Exec("INSERT INTO user (userid,username,password,email) VALUES (?,?,?,?)", userID, username, password, email)
	if err2 != nil {
		print(err2.Error())
		return false
	} else {
		print("ok!")
		return true
	}
}

func SelectHistoryByID(userID string) (historyirder []HistoryOrder) {
	db, err := sql.Open("mysql", DB_connect_string)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	rows, err := db.Query("SELECT orderid, orderresult FROM `order` WHERE userid=?", userID)
	for rows.Next() {
		var history1 HistoryOrder
		rows.Scan(&history1.OrderID, &history1.Orderresult)
		//print(history1.OrderID)
		historyirder = append(historyirder, history1)
	}
	return historyirder
}
func SaveOrder(orderID string, orderresult int64, userID string) bool {
	db, err := sql.Open("mysql", DB_connect_string)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	_, err2 := db.Exec("INSERT INTO `order` (orderid,orderresult,userid) VALUES (?,?,?)", orderID, orderresult, userID)
	if err2 != nil {
		print(err2.Error())
		return false
	} else {
		print("ok!")
		return true
	}
}

func CheckUserId(userID string) bool {
	db, err := sql.Open("mysql", DB_connect_string)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	rows, err2 := db.Query("SELECT * FROM user WHERE userid=?", userID)
	if err2 != nil {
		println(err2.Error())
		return false
	} else if !rows.Next() {
		print("not found!")
		return false
	}
	println("ok!")
	return true
}
func CheckUserEmail(email string) bool {
	db, err := sql.Open("mysql", DB_connect_string)
	defer db.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	rows, err2 := db.Query("SELECT * FROM user WHERE email=?", email)
	if err2 != nil {
		println(err2.Error())
		return false
	} else if !rows.Next() {
		print("not found!")
		return false
	}
	println("ok!")
	return true
}
