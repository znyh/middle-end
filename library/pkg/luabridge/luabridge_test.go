package luabridge

import (
	"fmt"
	"testing"
)

type data struct {
	Name string
	Age  string
	Nums []int32
}

func (self *data) Say() {
	fmt.Print(self.Name, self.Age, "\n")
}

func (self *data) Set(name string, age string) {
	self.Name = name
	self.Age = age
}

//Test
func TestCall(t *testing.T) {

	L, err := SafeLoad("test.lua")
	if err != nil {
		t.Fail()
	}
	var d data
	if err = SafeCall(L, "get_table", &d); err != nil {
		t.Fail()
	}

	var s string
	if err = SafeCall(L, "get_str", &s, "cyd", 111); err != nil {
		t.Fail()
	}

	var f float32
	if err = SafeCall(L, "get_num", &f); err != nil {
		t.Fail()
	}

	var i int16
	if err = SafeCall(L, "get_num", &i); err != nil {
		t.Fail()
	}
}

func TestReg(t *testing.T) {
	L, err := SafeLoad("test.lua")
	if err != nil {
		t.Fail()
	}

	d := data{
		Name: "cyd",
		Age:  "28",
	}

	if err = SafeCall(L, "test_call_go", nil, &d); err != nil {
		fmt.Println("call reg func:" + err.Error())
		t.Fail()
	}
	fmt.Println(d)
}
