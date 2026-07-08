# 备忘录后端API列表

## 登录页面 "/login"

```
/login
```

* 登录页面URL，需要前端以**form表单**提供用户输入的用户名name与密码password

```
{
	"name": "John",
	"password": "12345678"
}
```

* 如果用户名为空，返回422错误：

```
{
	"code": 422,
	"msg": "用户名不能为空"
}
```

* 如果密码小于6位，返回422错误：

```
{
	"code": 422,
	"msg": "密码不能少于6位"
}
```

* 如果用户名未注册，返回422错误，并重定向到注册URL：

```
{
	"code": 422,
	"msg": "该用户不存在"
}
```

* 如果注册token异常，返回500错误：

```
{
	"code": 500,
	"msg": "token发放异常"
}
```

* 如果密码错误，返回422错误：

```
{
	"code": 422,
	"msg": "密码错误，登录失败！"
}
```

* 如果密码正确，返回200，且返回token字符串

```
{
	"code": 200,
	"msg": "登陆成功！",
	"token": token
}
```

## 注册页面 "/register"

```
/register
```

注册页面，

* 注册页面URL，需要前端以**form表单**提供用户输入的用户名name与密码password

```
{
	"name": "John",
	"password": "12345678"
}
```

* 如果用户名为空，返回422错误：

```
{
	"code": 422,
	"msg": "用户名不能为空"
}
```

* 如果密码小于6位，返回422错误：

```
{
	"code": 422,
	"msg": "密码不能少于6位"
}
```

* 如果用户名已经注册，返回422错误，并重定向到登录页面

```
{
	"code": 422,
	"msg": "该用户存在"
}
```

* 如果注册新用户异常，返回500错误：

```
{
	"code": 400,
	"msg": "注册失败！"
}
```

* 如果注册成功，返回200，且返回用户信息，请跳转至登录页面

```
{
    "code": 200,
    "data": {
        "ID": 10,
        "Name": "Dox",
        "Password": "12345678"
    },
    "msg": "注册成功！"
}
```

## 备忘录主页面 "/index"

以下内容，需要前端提供authorization，bearer token格式，token在登录时由后端返回

备忘录结构：每一条备忘录包含四个字段：

* MemoID：备忘录ID（主键）
* UserID：书写本条本条备忘录的用户ID
* Content：备忘录内容
* Status：当前备忘录状态：
  * 0：处于“已完成”状态，用户显示该备忘录为灰色
  * 1：处于“未完成”状态，用户显示该备忘录为正常颜色

对于post、put与delete方法，如果token过期或无效，统一返回：

```
{
	"code": 401, 
	"msg": "权限不足"
}
```

### GET请求：获取当前用户记录的所有备忘录：

```
/index
```

只需要提供认证token即可，返回格式如下：代码200，data包含当前用户所有备忘录的状态

```
{
    "code": 200,
    "data": [
        {
            "MemoID": 6,
            "UserID": 8,
            "Content": "唱歌",
            "Status": 1
        },
        {
            "MemoID": 8,
            "UserID": 8,
            "Content": "洗澡",
            "Status": 1
        }
    ],
    "msg": "获取备忘录信息成功"
}
```

### POST请求：书写新备忘录

```
/index
```

需要提供认证token，并通过**form表单**，提供新加入的备忘录的content

```
{
	"content": "打羽毛球"
}
```

返回格式：代码200，并直接返回经过post操作后该用户的所有备忘录，用于刷新页面使用

```
{
    "code": 200,
    "data": [
        {
            "MemoID": 6,
            "UserID": 8,
            "Content": "唱歌",
            "Status": 1
        },
        {
            "MemoID": 8,
            "UserID": 8,
            "Content": "洗澡",
            "Status": 1
        },
        {
            "MemoID": 9,
            "UserID": 8,
            "Content": "打羽毛球",
            "Status": 1
        }
    ],
    "msg": "获取备忘录信息成功"
}
```

### PUT请求：修改某一条备忘录的状态

```
/index
```

需要提供认证token，并通过**json格式**向后端传递希望切换状态（“已完成”或“未完成”）的备忘录ID

返回值：代码200，并直接返回经过put操作后该用户的所有备忘录，用于刷新页面使用

```
{
    "code": 200,
    "data": [
        {
            "MemoID": 6,
            "UserID": 8,
            "Content": "唱歌",
            "Status": 0
        },
        {
            "MemoID": 8,
            "UserID": 8,
            "Content": "洗澡",
            "Status": 1
        },
        {
            "MemoID": 9,
            "UserID": 8,
            "Content": "打羽毛球",
            "Status": 1
        }
    ],
    "msg": "获取备忘录信息成功"
}
```

### DELETE请求：彻底删除某一条备忘录

```
/index
```

需要提供认证token，并通过**json格式**向后端传递希望删除的备忘录ID

返回值：代码200，并直接返回经过delete操作后该用户的所有备忘录，用于刷新页面使用

```
{
    "code": 200,
    "data": [
        {
            "MemoID": 8,
            "UserID": 8,
            "Content": "洗澡",
            "Status": 1
        },
        {
            "MemoID": 9,
            "UserID": 8,
            "Content": "打羽毛球",
            "Status": 1
        }
    ],
    "msg": "获取备忘录信息成功"
}
```

