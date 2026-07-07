package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"memo/controller"
)

type Memo struct {
	ID      uint `gorm:"primarykey"`
	Content string
	Status  int
}
type postForm struct {
	Content string `json:"content" form:"content"`
}
type putJson struct {
	ID uint `json:"id"`
}

var (
	db    *gorm.DB
	err   error
	memos []Memo
)

// 返回全部数据
func allRows(ctx *gin.Context) {
	var mms []Memo
	result := db.Find(&mms)
	if result.Error != nil {
		fmt.Println("错误：", result.Error.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "获取备忘录信息失败",
			"data": nil,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "获取备忘录信息成功",
			"data": mms,
		})
	}
}

// GET请求
func getFunc(ctx *gin.Context) {
	allRows(ctx)
}

func postFunc(ctx *gin.Context) {
	var form postForm
	err := ctx.ShouldBind(&form)
	if err != nil {
		fmt.Println("错误：", err.Error())
		return
	}
	tmp_memo := Memo{Content: form.Content, Status: 1}
	db.Create(&tmp_memo)
	allRows(ctx)
}

func putFunc(ctx *gin.Context) {
	var put putJson
	err := ctx.ShouldBind(&put)
	if err != nil {
		fmt.Println("错误：", err.Error())
		return
	}
	var tmp_memo Memo
	db.First(&tmp_memo, put.ID)
	// 翻转状态
	if tmp_memo.Status == 1 {
		tmp_memo.Status = 0
	} else {
		tmp_memo.Status = 1
	}
	db.Save(&tmp_memo)
	// 返回新的全部数表
	allRows(ctx)
}
func deleteFunc(ctx *gin.Context) {
	var delete putJson
	err := ctx.ShouldBind(&delete)
	if err != nil {
		fmt.Println("错误：", err.Error())
		return
	}
	// 根据ID删除对应行
	var tmp_memo Memo
	db.Delete(&tmp_memo, delete.ID)
	// 返回新的全部数表
	allRows(ctx)
}

func main() {
	// 默认路由
	route := gin.Default()
	// 连接数据库
	DSN := "root:Johntang2005@tcp(127.0.0.1:3306)/memo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败，错误：", err.Error())
	} else {
		fmt.Println("数据库连接成功")
	}
	// 自动迁移
	db.AutoMigrate(&Memo{})

	// 用户数据库的初始化
	userDB := controller.InitUserDB()
	// 登录注册操作
	route.POST("/login", controller.LoginFunc(userDB))
	route.POST("/register", controller.RegisterFunc(userDB))

	// 对于备忘录的CRUD操作
	route.GET("/index", getFunc)
	route.POST("/index", postFunc)
	route.PUT("/index", putFunc)
	route.DELETE("/index", deleteFunc)

	route.Run(":8080")
}
