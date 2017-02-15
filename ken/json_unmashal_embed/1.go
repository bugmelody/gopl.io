// json.Unmarshal 使用匿名字段不会设置该字段.要么有名字,要么通过tag给名字
package main

import (
	"encoding/json"
	"log"
)

type Login_info struct {
	Ticket           string `json:"ticket"`
	Ticket_expire_at string `json:"ticket_expire_at"`
	Tmp_ret_key      string `json:"tmp_ret_key"`
	Device_token     string `json:"device_token"`
}

type Basic_user_info struct {
	Id         int
	Username   string
	Phone      string
	Mail       string
	Nickname   string
	Logo       string
	Bind       string
	Is_regular int
}

var d struct {
	Errno        int
	Msg          string
	Device_token string
	Data         struct {
		Basic_user_info
		Login_info
	}
}

var e struct {
	Errno        int
	Msg          string
	Device_token string
	Data         struct {
		Basic_user_info Basic_user_info
		Login_info      Login_info
	}
}

var f struct {
	Errno        int
	Msg          string
	Device_token string
	Data         struct {
		Basic_user_info `json:"basic_user_info"`
		Login_info      `json:"login_info"`
	}
}

func main() {
	s := `{
		"errno": 5,
		"msg": "msg",
		"data": {
			"login_info": {
				"ticket": "I4Untv14863461351",
				"ticket_expire_at": "2017-03-08 09:55:35",
				"device_token": "74wq1486346135"
			},
		"basic_user_info": {
			"id": 1,
			"username": "testuser2",
			"phone": "",
			"mail": "",
			"nickname": "",
			"logo": "",
			"bind": "00000000",
			"isRegular": 1
		}
	}
	}`
	json.Unmarshal([]byte(s), &d)
	log.Println("=======================")
	log.Printf("%+v", d) // {Errno:5 Msg:msg Device_token: Data:{Basic_user_info:{Id:0 Username: Phone: Mail: Nickname: Logo: Bind: Is_regular:0} Login_info:{Ticket: Ticket_expire_at: Tmp_ret_key: Device_token:}}}
	log.Println("=======================")
	json.Unmarshal([]byte(s), &e)
	log.Println("=======================")
	log.Printf("%+v", e) // {Errno:5 Msg:msg Device_token: Data:{Basic_user_info:{Id:1 Username:testuser2 Phone: Mail: Nickname: Logo: Bind:00000000 Is_regular:0} Login_info:{Ticket:I4Untv14863461351 Ticket_expire_at:2017-03-08 09:55:35 Tmp_ret_key: Device_token:74wq1486346135}}}
	log.Println("=======================")
	json.Unmarshal([]byte(s), &f)
	log.Println("=======================")
	log.Printf("%+v", f) // {Errno:5 Msg:msg Device_token: Data:{Basic_user_info:{Id:1 Username:testuser2 Phone: Mail: Nickname: Logo: Bind:00000000 Is_regular:0} Login_info:{Ticket:I4Untv14863461351 Ticket_expire_at:2017-03-08 09:55:35 Tmp_ret_key: Device_token:74wq1486346135}}}
	log.Println("=======================")
}
