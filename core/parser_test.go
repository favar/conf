package core

import (
	"fmt"
	"testing"
)

type Account struct {
	Id      int
	Name    string
	Setting Setting
	Dbs     []string
}

type Setting struct {
	User string
	Ip   string
}

var s = 0.0

func GetId() float64 {
	s++
	return s
}

var source = map[string]interface{}{
	"Id":   GetId(),
	"Name": "UserName",
	"Age":  22,
	"Account": map[string]interface{}{
		"Id":   GetId(),
		"Name": "AccountName",
		"Setting": map[string]interface{}{
			"User": "zzz",
			"Ip":   "192.168.1.1",
		},
		"Dbs": []interface{}{"root", "123456", "3306", "192.168.1.1"},
	},
	"Dbs": []interface{}{"root", "123456", "3306", "192.168.1.1"},
}

var ps = JsonParser(source)

func Test_parser_Deep1Basic(t *testing.T) {
	var id1 int
	ps.Get("Id", &id1)
	println(id1)
}

func Test_parser_Deep2Basic(t *testing.T) {
	var id3 uint64
	ps.Get("Account.Id", &id3)
	println(id3)

}

func Test_parser_Deep1Struct(t *testing.T) {
	var account Account
	ps.Get("Account", &account)
	fmt.Println(account)
}

func Test_parser_Deep1Slice(t *testing.T) {
	var dbs []string = make([]string, 0)
	ps.Get("Dbs", &dbs)
	fmt.Println(dbs)
}

func Test_parser_Deep2Struct(t *testing.T) {
	var setting Setting
	ps.Get("Account.Setting", &setting)
	fmt.Println(setting)
	ps.Get("Account.Setting", &setting)
	fmt.Println(setting)
}

func Test_parser_None(t *testing.T) {
	var id1 int
	err := ps.Get("NoneId", &id1)
	if err == nil {
		t.Fail()
	}
	println(id1)
}
