package handler

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"user_srv/global"
	"user_srv/model"
	"user_srv/proto"
	"user_srv/utils"
)

var serviceName = "【User_Service】"

type UserService struct {
	proto.UnimplementedUserServer
}

// GetUserList
// @Description: 获取用户列表
// @receiver s
// @param ctx
// @param pageInfoRequest
// @return *proto.UserListResponse
// @return error
func (s *UserService) GetUserList(ctx context.Context, pageInfoRequest *proto.PageInfoRequest) (*proto.UserListResponse, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "GetUserList", "request", pageInfoRequest)
	//创建子 span
	parentSpan := opentracing.SpanFromContext(ctx)
	userListSpan := opentracing.GlobalTracer().StartSpan("GetUserList", opentracing.ChildOf(parentSpan.Context()))
	//实例化响应对象
	var response = &proto.UserListResponse{}
	//获取总行数
	var users []model.User
	result := global.DB.Find(&users)
	response.Total = int32(result.RowsAffected)
	//从数据库进行分页查询
	var pageUsers []model.User
	pagenum := pageInfoRequest.Pagenums
	pagesize := pageInfoRequest.Pagesize

	offset := utils.Paginate(int(pagenum), int(pagesize))

	global.DB.Offset(offset).Limit(int(pagesize)).Find(&pageUsers)
	//将查询的到的数据包装成回复
	for _, user := range pageUsers {
		userInfoResponse := utils.ModelToResponse(user)
		response.Data = append(response.Data, userInfoResponse)
	}
	//结束和返回
	userListSpan.Finish()
	return response, nil
}

// GetUserByMobile
// @Description: 通过电话号码获取用户信息
// @receiver s
// @param ctx
// @param mobileRequest
// @return *proto.UserInfoResponse
// @return error
func (s *UserService) GetUserByMobile(ctx context.Context, mobilerequest *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "GetUserList", "request", mobilerequest)
	//创建子 span
	parentSpan := opentracing.SpanFromContext(ctx)
	userbymobileSpan := opentracing.GlobalTracer().StartSpan("GetUserList", opentracing.ChildOf(parentSpan.Context()))
	//实例化响应对象
	var response = &proto.UserInfoResponse{}
	mobile := mobilerequest.Mobile
	//从数据库读取
	var user model.User
	result := global.DB.Where("mobile=?", mobile).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "未查找到该用户")
	}
	response = utils.ModelToResponse(user)
	//结束和返回
	userbymobileSpan.Finish()
	return response, nil

}

// GetUserById
// @Description: 通过ID获取用户信息
// @receiver s
// @param ctx
// @param idRequest
// @return *proto.UserInfoResponse
// @return error
func (s *UserService) GetUserById(ctx context.Context, idrequest *proto.IdRequest) (*proto.UserInfoResponse, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "GetUserList", "request", idrequest)
	//创建子 span
	parentSpan := opentracing.SpanFromContext(ctx)
	userbyidSpan := opentracing.GlobalTracer().StartSpan("GetUserList", opentracing.ChildOf(parentSpan.Context()))
	//实例化响应对象
	var response = &proto.UserInfoResponse{}
	id := idrequest.Id
	//从数据库读取
	var user model.User
	result := global.DB.Where("id=?", id).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "未查找到该用户")
	}
	response = utils.ModelToResponse(user)
	//结束和返回
	userbyidSpan.Finish()
	return response, nil
}

func (s *UserService) CreateUser(ctx context.Context, createuser *proto.CreateUserInfoRequest) (*proto.UserInfoResponse, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "GetUserList", "request", createuser)
	//创建子 span
	parentSpan := opentracing.SpanFromContext(ctx)
	createuserSpan := opentracing.GlobalTracer().StartSpan("GetUserList", opentracing.ChildOf(parentSpan.Context()))
	//实例化响应对象并创建用户
	var response = &proto.UserInfoResponse{}
	var user model.User
	//判断用户是否存在
	mobile := createuser.Mobile
	result := global.DB.Where("mobile=?", mobile).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	//创建用户并保存在数据库中
	inputpassword := createuser.Password
	encryptpassword := utils.EncryptPassword(inputpassword)
	user.Password = encryptpassword
	user.Mobile = mobile
	user.NickName = createuser.NickName

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	//将用户数据转化成标准响应数据
	response = utils.ModelToResponse(user)
	//结束和返回
	createuserSpan.Finish()
	return response, nil
}

func (s *UserService) UpdataUser(ctx context.Context, UpdateUserInfoRequest *proto.UpdataUserInfoRequest) (*proto.UpdateResponse, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "GetUserList", "request", UpdateUserInfoRequest)
	//创建子 span
	parentSpan := opentracing.SpanFromContext(ctx)
	UpdateUserInfoSpan := opentracing.GlobalTracer().StartSpan("GetUserList", opentracing.ChildOf(parentSpan.Context()))
	//实例化响应对象并创建用户
	var response = &proto.UpdateResponse{}
	var user model.User
	//在数据库找到对应用户
	result := global.DB.First(UpdateUserInfoRequest.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	birthDay := time.Unix(int64(UpdateUserInfoRequest.Birthday), 0)
	user.NickName = UpdateUserInfoRequest.NickName
	user.Birthday = &birthDay
	user.Gender = UpdateUserInfoRequest.Gender
	//更新信息
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	response.Success = true
	//结束和返回
	UpdateUserInfoSpan.Finish()
	return response, nil
}

func (s *UserService) CheckPassword(ctx context.Context, CheckPasswordRequest *proto.CheckPasswordRequest) (*proto.CheckPasswordResponse, error) {
	zap.S().Infow("Info", "service", serviceName, "method", "GetUserList", "request", CheckPasswordRequest)
	//创建子 span
	parentSpan := opentracing.SpanFromContext(ctx)
	CheckPasswordSpan := opentracing.GlobalTracer().StartSpan("GetUserList", opentracing.ChildOf(parentSpan.Context()))
	//实例化响应对象并创建用户
	var response = &proto.CheckPasswordResponse{}
	//检验密码
	password := CheckPasswordRequest.Password
	EncryptedPassword := CheckPasswordRequest.EncryptedPassword
	//格式化响应
	response.Success = utils.VerifyPassword(EncryptedPassword, password)
	//结束和返回
	CheckPasswordSpan.Finish()
	return response, nil
}
