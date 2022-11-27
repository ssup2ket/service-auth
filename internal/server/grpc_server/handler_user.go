package grpc_server

import (
	"context"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog/log"

	"github.com/ssup2ket/service-auth/internal/domain/entity"
	"github.com/ssup2ket/service-auth/internal/domain/service"
	"github.com/ssup2ket/service-auth/internal/server/errors"
	"github.com/ssup2ket/service-auth/internal/server/middleware"
	"github.com/ssup2ket/service-auth/internal/server/request"
	"github.com/ssup2ket/service-auth/pkg/entity/uuid"
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
			return nil, getErrNotFound(errors.ErrResouceUser)
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user")
		return nil, getErrServerError()
	}

	// Return user info
	return UserModelToUserInfo(userInfo), nil
}

func (s *ServerGRPC) UpdateUser(ctx context.Context, req *UserUpdateRequest) (*empty.Empty, error) {
	// Validate request
	if err := req.validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong update user request")
		return nil, getErrBadRequest()
	}

	// Update user
	if err := s.domain.User.UpdateUser(ctx, userUpdateToUserInfoModel(req), req.Password); err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			return nil, getErrNotFound(errors.ErrResouceUser)
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user")
		return nil, getErrBadRequest()
	}

	return &empty.Empty{}, nil
}

func (s *ServerGRPC) DeleteUser(ctx context.Context, req *UserIDRequest) (*empty.Empty, error) {
	// Validate request
	if err := req.validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong delete user request")
		return nil, getErrBadRequest()
	}

	// Delete user
	if err := s.domain.User.DeleteUser(ctx, uuid.FromStringOrNil(req.Id)); err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			return nil, getErrNotFound(errors.ErrResouceUser)
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user")
		return nil, getErrServerError()
	}

	return &empty.Empty{}, nil
}

func (s *ServerGRPC) GetUserMe(ctx context.Context, req *empty.Empty) (*UserInfoResponse, error) {
	// Get user ID
	userID, err := middleware.GetUserIDFromCtx(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("No user ID in context")
		return nil, getErrServerError()
	}

	// Get user
	userInfo, err := s.domain.User.GetUser(ctx, uuid.FromStringOrNil(string(userID)))
	if err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			return nil, getErrNotFound(errors.ErrResouceUser)
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user")
		return nil, getErrServerError()
	}

	// Return user info
	return UserModelToUserInfo(userInfo), nil
}

func (s *ServerGRPC) UpdateUserMe(ctx context.Context, req *UserUpdateRequest) (*empty.Empty, error) {
	// Get and set user ID
	userID, err := middleware.GetUserIDFromCtx(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("No user ID in context")
		return nil, getErrServerError()
	}
	req.Id = userID

	// Update user
	if err := s.domain.User.UpdateUser(ctx, userUpdateToUserInfoModel(req), req.Password); err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			return nil, getErrNotFound(errors.ErrResouceUser)
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user")
		return nil, getErrBadRequest()
	}

	return &empty.Empty{}, nil
}

func (s *ServerGRPC) DeleteUserMe(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	// Get user ID
	userID, err := middleware.GetUserIDFromCtx(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("No user ID in context")
		return nil, getErrServerError()
	}

	// Delete user
	if err := s.domain.User.DeleteUser(ctx, uuid.FromStringOrNil(userID)); err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			return nil, getErrNotFound(errors.ErrResouceUser)
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user")
		return nil, getErrServerError()
	}

	return &empty.Empty{}, nil
}

// Request validate
func (u *UserListRequest) validate() error {
	return nil
}

func (u *UserIDRequest) validate() error {
	return request.ValidateUserUUID(u.Id)
}

func (u *UserCreateRequest) validate() error {
	return request.ValidateUserCreate(u.LoginId, u.Password, u.Role, u.Phone, u.Email)
}

func (u *UserUpdateRequest) validate() error {
	return request.ValidateUserUpdate(u.Id, u.Password, u.Role, u.Phone, u.Email)
}

// DTO <-> Model
func userCreateToUserInfoModel(userCreate *UserCreateRequest) *entity.UserInfo {
	return &entity.UserInfo{
		LoginID: userCreate.LoginId,
		Role:    entity.UserRole(userCreate.Role),
		Phone:   userCreate.Phone,
		Email:   userCreate.Email,
	}
}

func userUpdateToUserInfoModel(userUpdate *UserUpdateRequest) *entity.UserInfo {
	return &entity.UserInfo{
		ID:    uuid.FromStringOrNil(userUpdate.Id),
		Role:  entity.UserRole(userUpdate.Role),
		Phone: userUpdate.Phone,
		Email: userUpdate.Email,
	}
}

func UserModelToUserInfo(userModel *entity.UserInfo) *UserInfoResponse {
	return &UserInfoResponse{
		Id:      userModel.ID.String(),
		LoginId: userModel.LoginID,
		Role:    string(userModel.Role),
		Phone:   userModel.Phone,
		Email:   userModel.Email,
	}
}

func UserModelListToUserInfoList(userModelList []entity.UserInfo) *UserListResponse {
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
