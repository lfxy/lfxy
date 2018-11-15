package main

import (
    "bytes"
    r "crypto/rand"
    "fmt"
    "math/rand"
    "strconv"
    "strings"
    "time"
    "os/exec"
    log "github.com/sirupsen/logrus"
)

/**
*生成随机字符
**/
func RandString(length int) string {
    rand.Seed(time.Now().UnixNano())
    rs := make([]string, length)
    for start := 0; start < length; start++ {
        t := rand.Intn(3)
        if t == 0 {
            rs = append(rs, strconv.Itoa(rand.Intn(10)))
        } else if t == 1 {
            rs = append(rs, string(rand.Intn(26)+65))
        } else {
            rs = append(rs, string(rand.Intn(26)+97))
        }
    }
    return strings.Join(rs, "")
}

func Rs2(length int) []byte {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    buffer := bytes.NewBufferString("")
    for i := 0; i < length; i++ {
        isLetter := r.Intn(2)
        if isLetter > 0 {
            letter := r.Intn(52)
            if letter < 26 {
                letter += 97
            } else {
                letter += 65 - 26
            }
            buffer.WriteString(string(letter))
            //buffer.WriteString(fmt.Sprintf("%c", letter))
        } else {
            buffer.WriteString(strconv.Itoa(r.Intn(10)))
        }
    }
    return buffer.Bytes()
}

func RandomCreateBytes(n int, alphabets ...byte) []byte {
    //const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    const alphanum = "0123456789abcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    var randby bool
    if num, err := r.Read(bytes); num != n || err != nil {
        rand.Seed(time.Now().UnixNano())
        randby = true
    }
    for i, b := range bytes {
        if len(alphabets) == 0 {
            if randby {
                bytes[i] = alphanum[rand.Intn(len(alphanum))]
            } else {
                bytes[i] = alphanum[b%byte(len(alphanum))]
            }
        } else {
            if randby {
                bytes[i] = alphabets[rand.Intn(len(alphabets))]
            } else {
                bytes[i] = alphabets[b%byte(len(alphabets))]
            }
        }
    }
    return bytes
}

func Test_Randstring() {
    fmt.Println(string(Rs2(10)))
    st1 := time.Now().UnixNano()
    for i := 0; i < 1000; i++ {
        Rs2(10)
    }
    et1 := time.Now().UnixNano()
    st2 := time.Now().UnixNano()
    var k []byte
    for i := 0; i < 1000; i++ {
        k = RandomCreateBytes(10)
    }
    et2 := time.Now().UnixNano()

    st := time.Now().UnixNano()
    s := ""
    for i := 0; i < 1000; i++ {
        s = RandString(10)
    }
    et := time.Now().UnixNano()

    fmt.Println(string(s), string(k), st, et, et1-st1, et2-st2, et-st)
}


func TestScale(){
    out, err := exec.Command("bash", "-c", "oc get statefulsets kafka-94x1").Output()
    if err != nil {
        fmt.Println("err:", err)
    }else {
        fmt.Println("err is nil")
    }
    fmt.Println("out:", string(out))
}

func test_floattoin(){
    a := "ddd"
    log.Info("aaaa:", a)
    log.Infof("bbb:%s", a)
    var ti int64
    t := time.Now().Month()
    ti = int64(t)
    fmt.Printf("%d\n",ti)
    f1 := 15.453
    i1 := int(f1)
    fmt.Printf("%d\n",i1)

}
func test_map(m map[string][]*string){
    var v []*string
    str1 := "111"
    str2 := "222"
    v = append(v, &str1)
    v = append(v, &str2)
    m["aa"] = v
    if v1, ok := m["aa"]; ok {
        str3 := "333"
        v1 = append(v1, &str3)
        m["aa"] = v1
    }

}

func test_slice(sl []int){
    sl = append(sl, 1)
    sl = append(sl, 2)
    sl = append(sl, 3)
}
func test_map2(m map[string]int){
    m["aaa"] = 111
    m["bbb"] = 223
}
func test_func_param(){
    m := make(map[string][]*string)
    test_map(m)
    for _, in := range m["aa"] {
        fmt.Println(*in)
    }
    var sl []int
    test_slice(sl)
    fmt.Println("slice---------------")
    for _, ts := range sl {
        fmt.Println(ts)
    }
    m2 := make(map[string]int)
    test_map2(m2)
    fmt.Println("m2---------------")
    for k, v := range m2 {
        fmt.Println(k)
        fmt.Println(v)
    }

}
func main(){
    /*var k []byte
    for i := 0; i < 100; i++ {
        k = RandomCreateBytes(3)
        fmt.Println(string(k))
    }*/
    //TestScale()

//    test_func_param()

    first_run := "2017-11-10 10:10:11"
    timeLayout := "2006-01-02 15:04:05"
    loc, _ := time.LoadLocation("Local")
    theTime, _ := time.ParseInLocation(timeLayout, first_run, loc)
    fmt.Print(theTime)
    fmt.Print("\n")
    fmt.Println(int(theTime.Weekday()))
}
