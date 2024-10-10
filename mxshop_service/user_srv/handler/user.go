package handler

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"user_srv/global"
	"user_srv/model"
	"user_srv/proto"
	"user_srv/utils"
)

var serviceName = "【User_Service】"

type UserService struct{}

// GetUserList
// @Description: 获取用户列表
// @receiver s
// @param ctx
// @param pageInfoRequest
// @return *proto.UserListResponse
// @return error
func (s *UserService) GetUserList(ctx context.Context, pageInfoRequest *proto.PageInfoRequest, opts ...grpc.CallOption) (*proto.UserListResponse, error) {
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
func (s *UserService) GetUserByMobile(context.Context, *MobileRequest) (*UserInfoResponse, error)
