package test

import (
	"fmt"
	"runtime"
	"testing"
)

func HandleError(affect int) error {
	pc, _, _, ok := runtime.Caller(1) // 获取调用 getCallerInfo 的函数的信息
	if !ok {
		fmt.Println("unknown", 0)
	}
	fun := runtime.FuncForPC(pc)
	if fun == nil {
		fmt.Println("unknown", 0)
	}
	// fun.Name() 会返回完整的函数名，包括包名，例如 "main.myFunc"
	fmt.Println(fun.Name()) // 获取文件路径和行号
	fmt.Println(fun.FileLine(pc))

	return nil
}

func a() {
	HandleError(1)
}

func TestHandle(t *testing.T) {
	a()
}
