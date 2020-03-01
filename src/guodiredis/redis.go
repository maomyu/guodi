package guodiredis

import (
	"fmt"
	"guodi/src/gore"
	"guodi/src/guodizap"

	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)

		guodizap.Panicf("%s:%s", msg, err)
	}
}

//appID验证
func CheckAppID(appID string) (isexist bool, appsecretID string) {
	conn, err := gore.Dial("192.168.10.252:16380")
	failOnError(err, "redis connection failed")
	defer conn.Close()
	rep, err := gore.NewCommand("HGET", "AppAuth", appID).Run(conn)
	s, _ := rep.String()
	if s != "" {
		return true, s
	}
	return false, ""

}

//Token验证
func CheckToken(email string, token string) bool {
	conn, err := gore.Dial("192.168.10.252:16380")
	failOnError(err, "redis connection failed")

	defer conn.Close()

	rep, err := gore.NewCommand("GET", email).Run(conn)
	s, _ := rep.String()
	if token == s {
		reptime, err := gore.NewCommand("TTL", email).Run(conn) //返回给定的key剩余的生存时间，返回单位为秒
		time, _ := reptime.Int()                                //将返回的gore.Reply类型转换为int类型

		fmt.Println(time)
		_, err = gore.NewCommand("EXPIRE", email, time+20).Run(conn) //设置过期时间
		failOnError(err, "redis time set up failed")
		return true
	}
	return false
}

//邮箱验证码验证
//邮箱验证码验证
func CheckEmailTempAuthentic(email string, authentic string) bool {
	conn, err := gore.Dial("192.168.10.252:16380")
	failOnError(err, "redis connection failed")

	defer conn.Close()

	rep, err := gore.NewCommand("GET", email+"temp").Run(conn)
	s, _ := rep.String()
	if s != "" && s == authentic {
		return true
	}
	return false
}

//保存token
func SaveToken(email string, token string) bool {
	conn, err := gore.Dial("192.168.10.252:16380")
	failOnError(err, "redis connection failed")
	defer conn.Close()

	_, err = gore.NewCommand("SET", email, token).Run(conn)
	if err != nil {
		return false
	}
	_, err = gore.NewCommand("EXPIRE", email, 10*24*60).Run(conn) //设置过期时间
	return true

}

//保存验证码
func SaveAuthenticate(email string, authenticate string) bool {
	conn, err := gore.Dial("192.168.10.252:16380")
	failOnError(err, "redis connection failed")
	defer conn.Close()

	_, err = gore.NewCommand("SET", email+"temp", authenticate).Run(conn)
	if err != nil {
		return false
	}
	_, err = gore.NewCommand("EXPIRE", email+"temp", 60).Run(conn) //设置过期时间
	return true

}

//验证客服人员信息是否正确
func CheckCustomerService(customerservicename string) (customerID string) {
	conn, err := gore.Dial("192.168.10.252:16380")
	failOnError(err, "redis connection failed")
	defer conn.Close()
	rep, err := gore.NewCommand("HGET", "customerservice", customerservicename).Run(conn)
	if err != nil {
		log.Println("查找人员失败")
		return ""
	}
	s, _ := rep.String()
	return s
}

//返回所有的客户ID
func SelectCustomer() []string {
	conn, err := gore.Dial("192.168.10.252:16380")
	failOnError(err, "redis connection failed")
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
	return total
}
func SelectCustomerID(customername string) string {
	conn, err := gore.Dial("192.168.10.252:16380")
	failOnError(err, "redis connection failed")
	defer conn.Close()
	rep, err := gore.NewCommand("HKEYS", "customerservice").Run(conn)
	a, _ := rep.Array()
	for _, value := range a {
		s, _ := value.String()
		repvalue, _ := gore.NewCommand("HGET", "customerservice", s).Run(conn)
		r, _ := repvalue.String()
		if r == customername {
			return s
		}
	}
	return ""
}
