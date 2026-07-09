package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"memo/common"
	"memo/controller"
	middleware "memo/middleWare"
)

type Memo struct {
	MemoID  uint `gorm:"primarykey"`
	UserID  uint
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
	db     *gorm.DB
	userDB *gorm.DB
	err    error
	memos  []Memo
)

// 返回该用户的全部数据
func allRows(ctx *gin.Context) {
	// 获得当前用户ID
	tmp_user, _ := ctx.Get("user")
	userID := tmp_user.(common.UserInfo).ID
	// 只获得当前用户的ID
	var mms []Memo
	result := db.Where("user_id=?", userID).Find(&mms)
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
	tmp_user, _ := ctx.Get("user")
	userID := tmp_user.(common.UserInfo).ID
	tmp_memo := Memo{UserID: userID, Content: form.Content, Status: 1}
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
	db.Where("memo_id=?", put.ID).First(&tmp_memo)

	// 检查有没有绑定前端传过来的json
	fmt.Println("前端传过来的put", put.ID)

	// 检查这条备忘录是不是当前用户写的
	tmp_user, _ := ctx.Get("user")
	userID := tmp_user.(common.UserInfo).ID
	if tmp_memo.UserID != userID {
		fmt.Printf("这条信息是%d写的，而目前操作者是%d", tmp_memo.UserID, userID)
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
		return
	}

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

	var tmp_memo Memo
	db.Where("memo_id=?", delete.ID).First(&tmp_memo)
	// 检查这条备忘录是不是当前用户写的
	tmp_user, _ := ctx.Get("user")
	userID := tmp_user.(common.UserInfo).ID
	fmt.Println("找到这条备忘录的作者是：", tmp_memo.UserID)
	if tmp_memo.UserID != userID {
		fmt.Printf("token无效，这条信息是%d写的，而目前操作者是%d", tmp_memo.UserID, userID)
		ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
		return
	}

	// 根据memoID删除对应行
	db.Delete(&tmp_memo, delete.ID)
	// 返回新的全部数表
	allRows(ctx)
}

func main() {
	// 默认路由
	route := gin.Default()
	// 允许所有CORS访问
	route.Use(cors.New(cors.Config{
		AllowAllOrigins:  true, // 允许所有源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// 连接数据库
	DSN := "root:Johntang2005@tcp(127.0.0.1:3306)/user_memos?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败，错误：", err.Error())
	} else {
		fmt.Println("数据库连接成功")
	}
	// 自动迁移
	db.AutoMigrate(&Memo{})

	// 用户数据库的初始化
	userDB = controller.InitUserDB()
	// 登录注册操作
	route.POST("/login", controller.LoginFunc(userDB))
	route.POST("/register", controller.RegisterFunc(userDB))

	// 对于备忘录的CRUD操作
	route.GET("/index", middleware.Authorization(), getFunc)
	route.POST("/index", middleware.Authorization(), postFunc)
	route.PUT("/index", middleware.Authorization(), putFunc)
	route.DELETE("/index", middleware.Authorization(), deleteFunc)

	route.Run(":8080")
}
