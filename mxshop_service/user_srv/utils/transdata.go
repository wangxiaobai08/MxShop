package utils

import (
	"user_srv/model"
	"user_srv/proto"
)

// ModelToResponse
// @Description: 将model.user转换成 UserInfoResponse
// @param user
// @return *proto.UserInfoResponse
func ModelToResponse(user model.User) *proto.UserInfoResponse {
	userInfoResponse := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoResponse.Birthday = uint64(user.Birthday.Unix())
	}
	return &userInfoResponse
}
