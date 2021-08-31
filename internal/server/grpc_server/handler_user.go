package grpc_server

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/service"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/server/errors"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/server/request"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/uuid"
)

func (s *ServerGRPC) ListUser(ctx context.Context, req *UserListRequest) (*UserListResponse, error) {
	// Validate request
	if err := req.validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong list user request")
		return nil, getErrBadRequest()
	}

	// List user
	userModels, err := s.domain.User.ListUser(ctx, int(req.Offset), int(req.Limit))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to list user")
		return nil, getErrServerError()
	}

	return UserModelListToUserInfoList(userModels), nil
}

func (s *ServerGRPC) CreateUser(ctx context.Context, req *UserCreateRequest) (*UserInfoResponse, error) {
	// Validate request
	if err := req.validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong create user request")
		return nil, getErrBadRequest()
	}

	// Create user
	user, err := s.domain.User.CreateUser(ctx, userCreateToUserInfoModel(req), req.Password)
	if err != nil {
		if err == service.ErrRepoConflict {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to create user becase of duplication")
			return nil, getErrConflict(errors.ErrResouceUser)
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create user")
		return nil, getErrBadRequest()
	}

	return UserModelToUserInfo(user), nil
}

func (s *ServerGRPC) GetUser(ctx context.Context, req *UserIDRequest) (*UserInfoResponse, error) {
	// Validate request
	if err := req.validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong get user request")
		return nil, getErrBadRequest()
	}

	// Get user
	userInfo, err := s.domain.User.GetUser(ctx, uuid.FromStringOrNil(string(req.Id)))
	if err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			return nil, getErrNoutFound(errors.ErrResouceUser)
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user")
		return nil, getErrServerError()
	}

	// Return user info
	return UserModelToUserInfo(userInfo), nil
}

func (s *ServerGRPC) UpdateUser(ctx context.Context, req *UserUpdateRequest) (*Empty, error) {
	// Validate request
	if err := req.validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong update user request")
		return nil, getErrBadRequest()
	}

	// Update user
	if err := s.domain.User.UpdateUser(ctx, userUpdateToUserInfoModel(req), req.Password); err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			return nil, getErrNoutFound(errors.ErrResouceUser)
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user")
		return nil, getErrBadRequest()
	}

	return &Empty{}, nil
}

func (s *ServerGRPC) DeleteUser(ctx context.Context, req *UserIDRequest) (*Empty, error) {
	// Validate request
	if err := req.validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong delete user request")
		return nil, getErrBadRequest()
	}

	// Delete user
	if err := s.domain.User.DeleteUser(ctx, uuid.FromStringOrNil(req.Id)); err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			return nil, getErrNoutFound(errors.ErrResouceUser)
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user")
		return nil, getErrServerError()
	}

	return &Empty{}, nil
}

// Request validate
func (u *UserListRequest) validate() error {
	return nil
}

func (u *UserIDRequest) validate() error {
	return request.ValidateUserUUID(u.Id)
}

func (u *UserCreateRequest) validate() error {
	return request.ValidateUserCreate(u.LoginID, u.Password, u.Phone, u.Email)
}

func (u *UserUpdateRequest) validate() error {
	return request.ValidateUserUpdate(u.Id, u.Password, u.Phone, u.Email)
}

// DTO <-> Model
func userCreateToUserInfoModel(userCreate *UserCreateRequest) *model.UserInfo {
	return &model.UserInfo{
		LoginID: userCreate.LoginID,
		Phone:   userCreate.Phone,
		Email:   userCreate.Email,
	}
}

func userUpdateToUserInfoModel(userUpdate *UserUpdateRequest) *model.UserInfo {
	return &model.UserInfo{
		ID:    uuid.FromStringOrNil(userUpdate.Id),
		Phone: userUpdate.Phone,
		Email: userUpdate.Email,
	}
}

func UserModelToUserInfo(userModel *model.UserInfo) *UserInfoResponse {
	return &UserInfoResponse{
		Id:      userModel.ID.String(),
		LoginId: userModel.LoginID,
		Phone:   userModel.Phone,
		Email:   userModel.Email,
	}
}

func UserModelListToUserInfoList(userModelList []model.UserInfo) *UserListResponse {
	userInfos := []*UserInfoResponse{}
	for _, userModel := range userModelList {
		tmp := UserInfoResponse{
			Id:      userModel.ID.String(),
			LoginId: userModel.LoginID,
			Phone:   userModel.Phone,
			Email:   userModel.Email,
		}
		userInfos = append(userInfos, &tmp)
	}
	return &UserListResponse{
		Uesrs: userInfos,
	}
}
