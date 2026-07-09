package controller //用户登录注册

import (
	"fmt"
	"memo/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TypeInUserInfo struct {
	Name     string `form:"name"`
	Password string `form:"password"`
}

var userDB *gorm.DB //用户数据库
var DBerr error

func InitUserDB() *gorm.DB {
	DSN := "root:Johntang2005@tcp(127.0.0.1:3306)/user_memos?charset=utf8mb4&parseTime=True&loc=Local"
	userDB, DBerr = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if DBerr != nil {
		fmt.Println("用户数据库连接失败！")
		return nil
	}
	fmt.Println("用户数据库连接成功！")
	var user common.UserInfo
	userDB.AutoMigrate(&user)
	return userDB
}

func GetUserDB() *gorm.DB {
	return userDB
}

func LoginFunc(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginUser TypeInUserInfo // 用户输入的用户名和密码
		var user common.UserInfo
		ctx.ShouldBind(&loginUser) // 获取注册表单参数
		// 判断用户输入参数是否正确
		if len(loginUser.Name) == 0 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "用户名不能为空",
			})
			return
		}
		if len(loginUser.Password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "密码不能少于6位",
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
			token, err := common.ReleaseToken(user)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "token发放异常",
				})
				return
			}
			fmt.Println("token发放成功，token是：", token)
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"msg":  "登陆成功！",
				"data": gin.H{"token": token},
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
		var user common.UserInfo
		ctx.ShouldBind(&loginUser) // 获取注册表单参数
		// 判断用户输入参数是否正确
		if len(loginUser.Name) == 0 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "用户名不能为空",
			})
			return
		}
		if len(loginUser.Password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "密码不能少于6位",
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
		cur_user := common.UserInfo{Name: loginUser.Name, Password: loginUser.Password}
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
