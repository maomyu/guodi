package guodiredis

import (
	"fmt"
	"guodi/src/gore"
	"log"
	"testing"
)

func TestCheckAppID(t *testing.T) {
	// 跳过测试
	if testing.Short() {
		t.Skip("该测试被跳过")
	}
	isexist, appsecretID := CheckAppID("sadfaf")
	fmt.Println(isexist)
	fmt.Println(appsecretID)
	fmt.Println("************************")
	isexist, appsecretID = CheckAppID("2354677290")
	fmt.Println(isexist)
	fmt.Println(appsecretID)
}
func TestCheckToken(t *testing.T) {
	if testing.Short() {
		t.Skip("该测试被跳过")
	}
	fmt.Println(CheckToken("2354677290@qq.com", "asdfsafdas"))
}

func TestSaveToken(t *testing.T) {
	// if testing.Short() {
	// 	t.Skip("该测试被跳过")
	// }
	fmt.Println(SaveToken("2354677290@qq.com", "asdfsafdas"))
}
func TestSaveAuthenticate(t *testing.T) {
	if testing.Short() {
		t.Skip("该测试被跳过")
	}
	fmt.Println(SaveAuthenticate("2345678978@qq.com", "erty"))
}
func TestCheckEmailTempAuthentic(t *testing.T) {
	if testing.Short() {
		t.Skip("该测试被跳过")
	}
	fmt.Println(CheckEmailTempAuthentic("2345678978@qq.com", "ddd"))
}
func TestSelectCustomer(t *testing.T) {
	if testing.Short() {
		t.Skip("该测试被跳过")
	}
	fmt.Println(SelectCustomer())
}

func TestAddTime(t *testing.T) {
	if testing.Short() {
		t.Skip("该测试被跳过")
	}
	conn, err := gore.Dial("192.168.10.252:16380")
	if err != nil {
		log.Println("redis connection failed")
	}
	defer conn.Close()

	rep, err := gore.NewCommand("TTL", "2345678978@qq.com").Run(conn)
	s, _ := rep.Int()
	if err != nil {
		// return false
	}
	fmt.Println(s)
	_, err = gore.NewCommand("EXPIRE", "2345678978@qq.com", s+20).Run(conn) //设置过期时间
}
func TestSelectCustomerID(t *testing.T) {
	conn, err := gore.Dial("192.168.10.252:16380")
	if err != nil {
		log.Println("redis connection failed")
	}
	defer conn.Close()

	var total []string
	rep, err := gore.NewCommand("HKEYS", "customerservice").Run(conn)
	a, _ := rep.Array()
	for _, value := range a {
		s, _ := value.String()
		total = append(total, s)
		fmt.Println(s)
	}
	fmt.Println(total)
}
