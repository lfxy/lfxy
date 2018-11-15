package main


import "fmt"
import "time"

type T struct {
    Id   int
    Name string
}

func Copy(a *T, b *T) error {
    b.Id = 5
    b.Name = "gert"
    a = b
    return nil
}
func CopyThatActuallyCopies(a *T, b *T) error {
    b.Id = 5
    b.Name = "gert"
    *a = *b
    return nil
}

func main() {
    /*var a = &T{1, "one"}
    var b = &T{2, "two"}

    fmt.Println(a, b)
    Copy(a, b)
    fmt.Println(a, b)
    CopyThatActuallyCopies(a, b)
    fmt.Println(a, b)
	fmt.Println("aaa:%s", *a)
	var timetest1 time.Time
	var timetest2 time.Time
	timetest1 = time.Now()
	time.Sleep(4000000000)
	timetest2 = time.Now()
	if timetest1.Add(time.Second * 3).After(timetest2) {
		fmt.Println("11111")
	} else {
		fmt.Println("22222")
	}
    fmt.Println(timetest1)
    fmt.Println(timetest2)*/

}
