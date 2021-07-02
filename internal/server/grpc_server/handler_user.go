package grpc_server

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/service"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/server/errors"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/server/request"
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

func (s *ServerGRPC) GetUser(ctx context.Context, req *UserUUIDRequest) (*UserInfoResponse, error) {
	// Validate request
	if err := req.validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong get user request")
		return nil, getErrBadRequest()
	}

	// Get user
	userInfo, err := s.domain.User.GetUser(ctx, string(req.UUID))
	if err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			return nil, getErrNoutFound(errors.ErrResouceUser)
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user")
		return nil, getErrServerError()
	}

	// Return user info
	return &UserInfoResponse{
		UUID:  userInfo.UUID,
		ID:    userInfo.ID,
		Phone: userInfo.Phone,
		Email: userInfo.Email,
	}, nil
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

	return nil, nil
}

func (s *ServerGRPC) DeleteUser(ctx context.Context, req *UserUUIDRequest) (*Empty, error) {
	// Validate request
	if err := req.validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong delete user request")
		return nil, getErrBadRequest()
	}

	// Delete user
	if err := s.domain.User.DeleteUser(ctx, req.UUID); err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			return nil, getErrNoutFound(errors.ErrResouceUser)
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user")
		return nil, getErrServerError()
	}

	return nil, nil
}

// Request validate
func (u *UserListRequest) validate() error {
	return nil
}

func (u *UserUUIDRequest) validate() error {
	return request.ValidateUserUUID(u.UUID)
}

func (u *UserCreateRequest) validate() error {
	return request.ValidateUserCreate(u.ID, u.Password, u.Phone, u.Email)
}

func (u *UserUpdateRequest) validate() error {
	return request.ValidateUserUpdate(u.UUID, u.Password, u.Phone, u.Email)
}

// DTO <-> Model
func userCreateToUserInfoModel(userCreate *UserCreateRequest) *model.UserInfo {
	return &model.UserInfo{
		ID:    userCreate.ID,
		Phone: userCreate.Phone,
		Email: userCreate.Email,
	}
}

func userUpdateToUserInfoModel(userUpdate *UserUpdateRequest) *model.UserInfo {
	return &model.UserInfo{
		UUID:  userUpdate.UUID,
		Phone: userUpdate.Phone,
		Email: userUpdate.Email,
	}
}

func UserModelToUserInfo(userModel *model.UserInfo) *UserInfoResponse {
	return &UserInfoResponse{
		UUID:  userModel.UUID,
		ID:    userModel.ID,
		Phone: userModel.Phone,
		Email: userModel.Email,
	}
}

func UserModelListToUserInfoList(userModelList []model.UserInfo) *UserListResponse {
	userInfos := []*UserListResponse_User{}
	for _, userModel := range userModelList {
		tmp := UserListResponse_User{
			UUID:  userModel.UUID,
			ID:    userModel.ID,
			Phone: userModel.Phone,
			Email: userModel.Email,
		}
		userInfos = append(userInfos, &tmp)
	}
	return &UserListResponse{
		Uesrs: userInfos,
	}
}
