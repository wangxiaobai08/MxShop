package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"math/rand"
	"mxshop_web_api/forms"
	"mxshop_web_api/global"
	"mxshop_web_api/utils"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	dysmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

var redisclient *redis.Client
var store = base64Captcha.DefaultMemStore

// SendNoteCode
// @Description: 发送验证码
// @param c
func SendNoteCode(c *gin.Context) {
	// 表单验证
	sendSmsForm := forms.SendSmsForm{}
	err := c.ShouldBind(&sendSmsForm)
	if err != nil {
		zap.S().Errorw("Error", "method", "SendNoteCode", "err", err.Error())
		utils.HandleValidatorError(c, err)
		return
	}

	config := sdk.NewConfig()
	credential := credentials.NewAccessKeyCredential(global.WebServiceConfig.AliSmsInfo.ApiKey, global.WebServiceConfig.AliSmsInfo.ApiSecret)
	/* use STS Token
	credential := credentials.NewStsTokenCredential("<your-access-key-id>", "<your-access-key-secret>", "<your-sts-token>")
	*/
	client, err := dysmsapi.NewClientWithOptions("cn-shenzhen", config, credential)
	if err != nil {
		panic(err)
	}
	smsCode := GenerateNoteCode(5)

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = global.WebServiceConfig.AliSmsInfo.SignName
	request.TemplateCode = global.WebServiceConfig.AliSmsInfo.TemplateCode
	request.PhoneNumbers = sendSmsForm.Mobile
	request.TemplateParam = "{\"code\":\"" + smsCode + "\"}"

	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
	ConnectRedis()
	redisclient.Set(context.WithValue(context.Background(), "ginContext", c), sendSmsForm.Mobile, smsCode, 300*time.Second)
	c.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}

// GetCaptcha
// @Description: 生成图形验证码
// @param c
func GetCaptcha(c *gin.Context) {
	//生成验证码驱动
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	//创建验证码实例
	captcha := base64Captcha.NewCaptcha(driver, store)
	//生成验证码并处理错误逻辑
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		zap.S().Error("Error", "method", "GetCaptcha", "生成验证码失败：", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"msg": "生成验证码错误"})
		return
	}
	//返回响应
	c.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})
}

// generateNoteCode
// @Description: 生成随机验证码
// @param width 验证码长度
// @return string
func GenerateNoteCode(width int) string {
	//数字数组
	number := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	len := len(number)
	//随机种子
	rand.Seed(time.Now().Unix())
	//字符串构建器
	var sb strings.Builder
	//循环生成随机数字
	for i := 0; i < len; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", number[rand.Intn(i)])
	}
	return sb.String()
}

// connectRedis
// @Description: 连接redis
func ConnectRedis() {
	redisclient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.WebServiceConfig.RedisInfo.Host, global.WebServiceConfig.RedisInfo.Port),
		Password: global.WebServiceConfig.RedisInfo.Password,
	})
}
