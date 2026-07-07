package controller //用户登录注册

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TypeInUserInfo struct {
	Name     string `form:"name"`
	Password string `form:"password"`
}

type UserInfo struct {
	ID       uint `gorm:"primarykey"`
	Name     string
	Password string
}

func InitUserDB() *gorm.DB {
	DSN := "root:Johntang2005@tcp(127.0.0.1:3306)/users?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println("用户数据库连接失败！")
		return nil
	}
	fmt.Println("用户数据库连接成功！")
	var user UserInfo
	db.AutoMigrate(&user)
	return db
}

func LoginFunc(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginUser TypeInUserInfo // 用户输入的用户名和密码
		var user UserInfo
		ctx.ShouldBind(&loginUser) // 获取注册表单参数
		// 判断用户输入参数是否正确
		if len(loginUser.Name) == 0 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户名不能为空",
			})
			return
		}
		if len(loginUser.Password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码不能少于6位",
			})
			return
		}
		// 判断用户是否存在
		db.Where("name=?", loginUser.Name).First(&user)
		if user.Name == "" { // 该用户不存在
			fmt.Println("用户不存在")
			// ctx.Redirect(http.StatusUnprocessableEntity, "/register")
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": http.StatusUnprocessableEntity,
				"msg":  "该用户不存在",
			})
			return
		}
		// 用户存在，核对密码是否正确
		if user.Password == loginUser.Password {
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"msg":  "登陆成功！",
			})
		} else {
			ctx.JSON(422, gin.H{
				"code": 422,
				"msg":  "密码错误，登录失败！",
			})
		}
	}
}

func RegisterFunc(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginUser TypeInUserInfo // 用户输入的用户名和密码
		var user UserInfo
		ctx.ShouldBind(&loginUser) // 获取注册表单参数
		// 判断用户输入参数是否正确
		if len(loginUser.Name) == 0 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "用户名不能为空",
			})
			return
		}
		if len(loginUser.Password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code":    422,
				"message": "密码不能少于6位",
			})
			return
		}
		// 检查用户名是否存在
		db.Where("name=?", loginUser.Name).First(&user)
		if user.Name != "" { // 该用户存在
			fmt.Println("用户存在")
			// ctx.Redirect(http.StatusUnprocessableEntity, "/login")
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": http.StatusUnprocessableEntity,
				"msg":  "该用户存在",
			})
			return
		}
		//用户不存在，可以注册
		cur_user := UserInfo{Name: loginUser.Name, Password: loginUser.Password}
		fmt.Println(cur_user)
		result := db.Create(&cur_user)
		if result.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "注册失败！",
			})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusOK,
				"msg":  "注册成功！",
				"data": cur_user,
			})
		}
	}
}
