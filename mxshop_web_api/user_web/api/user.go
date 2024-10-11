package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"mxshop_web_api/forms"
	"mxshop_web_api/global"
	"mxshop_web_api/global/response"
	"mxshop_web_api/proto"
	"mxshop_web_api/utils"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Description: 获取用户列表
// @param c
func GetUserList(c *gin.Context) {
	//获取参数
	pagenums := c.DefaultQuery("page", "0")
	pagenumsInt, _ := strconv.Atoi(pagenums)
	pagesize := c.DefaultQuery("size", "20")
	pagesizeInt, _ := strconv.Atoi(pagesize)
	//调用grpc中的服务
	GrpcData, err := global.UserClient.GetUserList(context.WithValue(context.Background(), "gincontt", c), &proto.PageInfoRequest{Pagenums: uint32(pagenumsInt),
		Pagesize: uint32(pagesizeInt),
	})
	if err != nil {
		zap.S().Error("[GetUserList] 查询 【用户列表】 失败", "msg", err.Error())
		utils.HandleValidatorError(c, err)
		return
	}
	//将数据包装成响应发送给前端
	result := make([]interface{}, 0)
	for _, value := range GrpcData.Data {
		user := response.UserResponse{
			Id:       value.Id,
			Name:     value.NickName,
			Gender:   value.Gender,
			Mobile:   value.Mobile,
			Birthday: time.Time(time.Unix(int64(value.Birthday), 0)),
		}
		result = append(result, user)
	}
	c.JSON(http.StatusOK, result)
}

func Register(c *gin.Context) {
	// 1.表单认证
	registerForm := forms.RegisterForm{}
	err := c.ShouldBind(&registerForm)
	if err != nil {

		fmt.Println("c.ShouldBind error", err.Error())
		utils.HandleValidatorError(c, err)
		return
	}
	// 2.通过redis 验证 验证码是否正确
	connectRedis()
	value, err := redisclient.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil { // redis中没有验证码
		zap.S().Warnw("验证码发送/redis存储失败", "用户手机号", registerForm.Mobile)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	} else { // 验证码错误
		if value != registerForm.Code {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "验证码错误",
			})
			return
		}
	}
	userResponse, err := global.UserClient.CreateUser(context.Background(), &proto.CreateUserInfoRequest{
		NickName: registerForm.Mobile,
		Password: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorw("[CreateUser] 失败", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, c)
		return
	}
	token, err := utils.GenerateToken(uint(userResponse.Id), userResponse.NickName, uint(userResponse.Role))
	if err != nil {
		zap.S().Errorw("生成token失败", "err:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":        userResponse.Id,
		"nickName":  userResponse.NickName,
		"token":     token,
		"expiresAt": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}

// PasswordLogin
// @Description: 手机密码登录
// @param c
func PasswordLogin(c *gin.Context) {
	//1.实例化自定义表单验证结构并绑定表单结构
	passwordloginform := forms.PasswordLoginForm{}
	err := c.Bind(&passwordloginform)
	if err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	//2.验证图形验证码是否正确
	verify := store.Verify(passwordloginform.CaptchaId, passwordloginform.Captcha, true)
	if !verify {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}
	//3.登录
	//3.1.获取用户加密后的密码
	userInfoResponse, err := global.UserClient.GetUserByMobile(context.WithValue(context.Background(), "ginContext", c), &proto.MobileRequest{Mobile: passwordloginform.Mobile})
	if err != nil {
		zap.S().Errorw("[GetUserByMobiles] 查询失败", "err", err.Error())
		utils.HandleGrpcErrorToHttpError(err, c)
	}
	//3.2.进行密码比对
	checkPasswordResponse, err := global.UserClient.CheckPassword(context.WithValue(context.Background(), "ginContext", c), &proto.CheckPasswordRequest{
		Password:          passwordloginform.Password,
		EncryptedPassword: userInfoResponse.Password,
	})
	if err != nil {
		zap.S().Errorw("[CheckPassword] 密码验证失败")
		utils.HandleGrpcErrorToHttpError(err, c)
	}
	//3.3.如果成功则生成token并返回响应
	if checkPasswordResponse.Success {
		token, err := utils.GenerateToken(uint(userInfoResponse.Id), userInfoResponse.NickName, uint(userInfoResponse.Role))
		if err != nil {
			zap.S().Errorw("生成token失败", "err:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "生成token失败",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":        userInfoResponse.Id,
			"nickName":  userInfoResponse.NickName,
			"token":     token,
			"expiresAt": (time.Now().Unix() + 60*60*24*30) * 1000,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "登录失败",
		})
	}
}
