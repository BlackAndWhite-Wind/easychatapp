package service

import (
	"EasyChatApp/common"
	"EasyChatApp/dao"
	"EasyChatApp/middleware"
	"EasyChatApp/models"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func List(ctx *gin.Context) {
	list, err := dao.GetUserList()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "获取用户列表失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, list)
}

func LoginByNameAndPassWord(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")
	data, err := dao.FindUserByName(name)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "登录失败",
		})
		return
	}
	if data.Name == "" {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名不存在",
		})
		return
	}
	ok := common.CheckPassWord(password, data.Salt, data.PassWord)
	if !ok {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "密码错误",
		})
		return
	}
	Rsp, err := dao.FindUserByNameAndPwd(name, data.PassWord)
	if err != nil {
		zap.S().Info("登录失败", err)
	}
	token, err := middleware.GenerateToken(Rsp.ID, "yk")
	if err != nil {
		zap.S().Info("生成token失败", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"tokens":  token,
		"userId":  Rsp.ID,
	})
}

func NewUser(ctx *gin.Context) {
	user := models.UserBasic{}
	user.Name = ctx.Request.FormValue("name")
	password := ctx.Request.FormValue("password")
	rePassword := ctx.Request.FormValue("Identity")

	if user.Name == "" || password == "" || rePassword == "" {
		ctx.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "用户名或密码不能为空！",
			"data":    user,
		})
		return
	}

	_, err := dao.FindUser(user.Name)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "该用户已注册",
			"data":    user,
		})
		return
	}
	if password != rePassword {
		ctx.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "两次密码不一致！",
			"data":    user,
		})
		return
	}

	salt := fmt.Sprintf("%d", rand.Int31())
	user.PassWord = common.SaltPassWord(password, salt)
	user.Salt = salt
	fmt.Println(user.PassWord)
	t := time.Now()
	user.LoginTime = &t
	user.LoginOutTime = &t
	user.HeartBeatTime = &t
	dao.CreateUser(user)
	ctx.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "新增用户成功！",
		"data":    user,
	})
}

func UpdateUser(ctx *gin.Context) {
	user := models.UserBasic{}
	id, err := strconv.Atoi(ctx.Request.FormValue("id"))
	if err != nil {
		zap.S().Info("类型转换失败", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "注销账号失败",
		})
		return
	}
	user.ID = uint(id)
	Name := ctx.Request.FormValue("name")
	PassWord := ctx.Request.FormValue("password")
	Email := ctx.Request.FormValue("email")
	Phone := ctx.Request.FormValue("phone")
	avatar := ctx.Request.FormValue("icon")
	gender := ctx.Request.FormValue("gender")
	if Name != "" {
		user.Name = Name
	}
	if PassWord != "" {
		salt := fmt.Sprintf("%d", rand.Int31())
		user.Salt = salt
		user.PassWord = common.SaltPassWord(PassWord, salt)
	}
	if Email != "" {
		user.Email = Email
	}
	if Phone != "" {
		user.Phone = Phone
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if gender != "" {
		user.Gender = gender
	}
	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		zap.S().Info("参数不匹配", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "参数不匹配",
		})
		return
	}
	Rsp, err := dao.UpdateUser(user)
	if err != nil {
		zap.S().Info("更新用户失败", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "修改信息失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "修改成功",
		"data":    Rsp.Name,
	})
}

func DeleteUser(ctx *gin.Context) {
	user := models.UserBasic{}
	id, err := strconv.Atoi(ctx.Request.FormValue("id"))
	if err != nil {
		zap.S().Info("类型转换失败", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "注销账号失败",
		})
		return
	}
	user.ID = uint(id)
	err = dao.DeleteUser(user)
	if err != nil {
		zap.S().Info("注销用户失败", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "注销账号失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "注销账号成功",
	})
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(ctx *gin.Context) {
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ctx, ws)
}

// MsgHandler 向socket连接发送消息
func MsgHandler(ctx *gin.Context, ws *websocket.Conn) {
	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(" MsgHandler 发送失败", err)
			continue
		}

		fmt.Println("接收数据：", string(data))
		err = common.Publish(ctx, common.PublishKey, string(data))
		if err != nil {
			fmt.Println("publish失败", err)
		}

		//从redis拿消息
		msg, err := common.Subscribe(ctx, common.PublishKey)
		if err != nil {
			fmt.Println(" MsgHandler 发送失败", err)
			return
		}

		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		//向连接写入消息
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			log.Fatalln("写入消息失败", err)
		}
	}
}

// SendUserMsg 发送消息
func SendUserMsg(ctx *gin.Context) {
	models.Chat(ctx.Writer, ctx.Request)
}

func ExitUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		zap.S().Info("类型转换失败", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "注销账号失败",
		})
		return
	}
	_, err = middleware.GenerateToken(uint(id), "exit")
	if err != nil {
		zap.S().Info("")
	}
}
