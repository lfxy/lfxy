package main

import (
	"fmt"
	"sync"
	"reflect"
//	"time"
)


func EchoNumber(i int) {
    fmt.Println(i)
}

func TestWG() {
	var wg sync.WaitGroup
    for i := 0; i < 5; i = i + 1 {
        wg.Add(1)
        go func(n int) {
             defer wg.Done()
            //defer wg.Add(-1)
            EchoNumber(n)
        }(i)
    }

    wg.Wait()
}

type People interface {
	GetRet(string) string
}

type Student struct{}

func (stu Student) GetRet(input string) string {
	fmt.Println("aaa ", input)
	return "aaa " + input
}

func (stu Student) Sval() {}
func (stu *Student) SPtr()  {}

func TestInterface() {
	var pe People = &Student{}
	pe.GetRet("ddd")
}

type Title struct {
	Name string
	Value int
}

func TestArray() {
	t1 := []Title {
		 {"aaa", 14},
		{"bbb", 15},
		{"ccc", 16},
	}

	for _, t := range t1 {
		fmt.Println(t)
	}
}

////////////
type N int
type NI interface {
	Value()
	Pointer()
}

func (n N) Value() {
	n++
	fmt.Printf("v: %p, %v\n", &n, n)
}
func (n *N) Pointer() {
	(*n)++
	fmt.Printf("v: %p, %v\n", n, *n)
}

func methedSet(a interface{}) {
	t := reflect.TypeOf(a)
	temp := t.NumMethod()
	fmt.Println("...", temp)
	for i, n := 0, t.NumMethod(); i < n; i++ {
		m := t.Method(i)
		fmt.Println(m.Name, m.Type)
	}
}

func TestPointer() {
	var a N = 25
	var i NI = &a
	i.Pointer()
	a.Pointer()
	fmt.Println(a)
	methedSet(&a)

	p := &a
	a.Value()
	a.Pointer()

	p.Value()
	p.Pointer()
	var s Student
	methedSet(s)
}

///////////
func main() {
	//TestWG()
	//TestInterface()
//	TestArray()
    TestPointer()
}
