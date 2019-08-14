package main

import (
	"fmt"
	"strconv"
)

//string与int类型的转换  Atoi
func strToInt() {
	str := "10"
	n, err := strconv.Atoi(str)
	if err != nil {
		panic("str to int error")
	}
	fmt.Println(n)
}

//Itoa  int与string类型相互转换
func intToStr() {
	n := 11
	s := strconv.Itoa(n)
	fmt.Println(s)
}

//parse系列  转换字符串给指定类型
func parseXXX() {
	//ParseBool
	b, err := strconv.ParseBool("true")
	if err != nil {
		panic("parse bool err")
	}
	fmt.Println(b)

	//ParseInt
	n, err := strconv.ParseInt("100", 10, 64)
	if err != nil {
		panic("parse int err")
	}
	fmt.Println(n)
}

//format 将给定类型格式化为字符串
func formatXXX() {
	s := strconv.FormatBool(true)
	fmt.Println(s)
}

//参数 r必须是：字母（广义）、数字、标点、符号、ASCII空格
func isPrintTest() {
	//r := []rune("abc")
	b := strconv.IsPrint('d')
	fmt.Println(b)
}

func canBackquoteTest() {
	s := `n
	hao`
	b := strconv.CanBackquote(s) //s只要是多行就为false
	fmt.Println(b)
}

//strconv 包含了基本数据类型与其字符串表示的转换
//主要常用函数：Atoi  Itia   parse系列  format系列  append系列
func main() {
	//strToInt()
	//intToStr()
	//parseXXX()
	//formatXXX()
	//isPrintTest()
	canBackquoteTest()
}
