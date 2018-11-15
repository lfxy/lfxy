package main

import (
	"fmt"
	"modules"
	"github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"
)

/*type Profile struct {
    Id          int
    Age         int16
    User        *User   `orm:"reverse(one)"` // 设置一对一反向关系(可选)
}*/

/*type Tag struct {
    Id    int
    Name  string
    Posts []*Post `orm:"reverse(many)"`
}*/

/*type Post struct {
    Id    int
    Title string
    User  *User  `orm:"rel(fk)"`    //设置一对多关系
    //Tags  []*Tag `orm:"rel(m2m)"`
}

type User struct {
    Id          int
    Name        string
    //Profile     *Profile   `orm:"rel(one)"` // OneToOne relation
    Post        []*Post `orm:"reverse(many)"` // 设置一对多的反向关系
}*/

type Post struct {
    Id    int `orm:"pk;auto"`
    Title string
    User  *User  `orm:"null;rel(fk);on_delete(set_null)"`    //设置一对多关系
}

type User struct {
    Id          int `orm:"pk;auto"`
    Name        string `orm:"index"`
    Post        []*Post `orm:"reverse(many)"` // 设置一对多的反向关系
}

func init() {
    // 需要在init中注册定义的model
    orm.RegisterDriver("mysql", orm.DRMySQL)

    orm.RegisterDataBase("default", "mysql", "root:restfulapi123@tcp(10.209.224.161:10022)/czq?charset=utf8")
    //orm.RegisterDataBase("czq", "mysql", "root:restfulapi123@/10.209.224.161:10022/czq?charset=utf8")
    orm.RegisterModel(new(User), new(Post))
    orm.RunSyncdb("default", false, true)
}

func test_user() {
    o := orm.NewOrm()

    user := new(User)
    user.Name = "slene1"
	post := new(Post)
	post.Title = "slene1_t"
	post.User = user

    fmt.Println(o.Insert(user))
    fmt.Println(o.Insert(post))
	user.Id = 2
    user.Name = "slene2"
    fmt.Println(o.Insert(user))

	var posts []*Post
	num, err := o.QueryTable("post").Filter("User", 1).RelatedSel().All(&posts)
	if err == nil {
		fmt.Printf("%d posts read\n", num)
		for _, post := range posts {
			fmt.Printf("Id: %d, UserName: %d, Title: %s\n", post.Id, post.User.Name, post.Title)
		}
	}

	user.Id = 3
	user.Name = "aaa"
	post.Id = 2;
    post.Title = "slene3_t"
	post.User = user
    fmt.Println(o.Insert(post))

}
func main() {
}
