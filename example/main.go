package main

import (
        "fmt"
        "github.com/wanderer69/debug"
)

func Tst1(a, b int) int {
        c := a+b
        debug.Alias("tst1").Printf("a %v b %v\r\n", a, b)
	return c
}

func Tst2(a, b int) int {
        c := a-b
        debug.Printf("a %v b %v\r\n", a, b)
	return c
}

func Tst3(a, b int) int {
        c := a*b
        debug.Label("").Printf("a %v b %v\r\n", a, b)
	return c
}

func AllTst() {
     res1 := Tst1(1, 2)
     res2 := Tst2(res1, 1)
     res3 := Tst3(res2, res1)
     fmt.Printf("res3 %v\r\n", res3)
}

func main() {
	debug.NewDebug()
	debug.SetArea(debug.Area{Alias:"tst1"}, debug.Area{Func:"Tst2"})
//	debug.SetArea(debug.Area{Func:"Tst3"})
        AllTst()
}
