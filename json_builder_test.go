package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

type user struct {
	Name string
	Age  int
}

func Test_z(t *testing.T) {
	k := 2
	v := reflect.ValueOf(&k)
	var s interface{} = 3
	v.Elem().Set(reflect.ValueOf(s))
	println(k)

}

func Test1(t *testing.T) {
	file, err := os.Open("json-conf/t1.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var c interface{}
	decoder.Decode(&c)
	fmt.Println("type:", reflect.TypeOf(c).Kind())
	fmt.Println(c)
	for k, v := range c.(map[string]interface{}) {
		fmt.Println("k:", k, ",v:", v, " type ", reflect.TypeOf(v).Kind())
		if k == "Account" {
			if m, ok := v.(map[string]interface{}); ok {
				fmt.Println("trans:", m)
			}
		}
		if k == "Port" {
			if p, ok := v.(int); ok {
				fmt.Println(p)
			} else {
				fmt.Printf("%v is not int \n", v)
			}
			c, _ := v.(float64)
			k := int(c)
			println(k)
		}

		if k == "Strings" {
			ss, ok := v.([]interface{})
			fmt.Println("t strings:", ok, " ", ss)
		}
	}
}

func Test_Conf(t *testing.T) {

	cnf, _ := Builder().
		AddJsonFiles("json-conf/t1.json", "json-conf/t2.json").
		Build()

	var s string
	err := cnf.Get("Host", &s)
	fmt.Println("value:", s, err)

	s = ""
	err = cnf.Get("zzz-Node1", &s)
	fmt.Println("value:", s, err)
}

func Test2(t *testing.T) {
	path := "json-conf/t2.json"
	file, _ := os.Open(path)
	defer file.Close()
	decoder := json.NewDecoder(file)
	var c interface{}
	decoder.Decode(&c)
	fmt.Println("type:", reflect.TypeOf(c).Kind())
}
