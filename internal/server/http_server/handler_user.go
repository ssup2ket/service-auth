package http_server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/service"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/server/errors"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/server/middleware"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/server/request"
	modeluuid "github.com/ssup2ket/ssup2ket-auth-service/pkg/model/uuid"
)

// List users
func (s *ServerHTTP) GetUsers(w http.ResponseWriter, r *http.Request, params GetUsersParams) {
	ctx := r.Context()

	// Set offsetm, limit
	offset := 0
	if params.Offset != nil {
		offset = int(*params.Offset)
	}
	limit := 0
	if params.Limit != nil {
		limit = int(*params.Limit)
	}

	// List user
	userModels, err := s.domain.User.ListUser(ctx, offset, limit)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to list user")
		render.Render(w, r, getErrRendererServerError())
		return
	}

	render.JSON(w, r, UserModelListToUserInfoList(userModels))
}

// Create a user
func (s *ServerHTTP) PostUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userCreate := UserCreate{}

	// Unmarshal request
	if err := render.Bind(r, &userCreate); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong create user request")
		render.Render(w, r, getErrRendererBadRequest())
		return
	}

	// Create user
	user, err := s.domain.User.CreateUser(ctx, userCreateToUserInfoModel(&userCreate), userCreate.Password)
	if err != nil {
		if err == service.ErrRepoConflict {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to create user becase of duplication")
			render.Render(w, r, getErrRendererConflict(errors.ErrResouceUser))
			return
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create user")
		render.Render(w, r, getErrRendererServerError())
		return
	}

	render.JSON(w, r, UserModelToUserInfo(user))
}

// Get a user
func (s *ServerHTTP) GetUsersUserID(w http.ResponseWriter, r *http.Request, userID UserID) {
	ctx := r.Context()

	// Validate request
	if err := userID.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong user ID")
		render.Render(w, r, getErrRendererBadRequest())
		return
	}

	// Get user
	userInfo, err := s.domain.User.GetUser(ctx, modeluuid.FromStringOrNil(string(userID)))
	if err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			render.Render(w, r, getErrRendererNotFound(errors.ErrResouceUser))
			return
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user")
		render.Render(w, r, getErrRendererServerError())
		return
	}

	render.JSON(w, r, UserModelToUserInfo(userInfo))
}

// Update a user
func (s *ServerHTTP) PutUsersUserID(w http.ResponseWriter, r *http.Request, userID UserID) {
	ctx := r.Context()
	userUpdate := UserUpdate{}

	// Unmarshal request
	if err := render.Bind(r, &userUpdate); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong update user request")
		render.Render(w, r, getErrRendererBadRequest())
		return
	}

	// Validate request
	if err := userID.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong user ID")
		render.Render(w, r, getErrRendererBadRequest())
		return
	}

	// Update user
	if err := s.domain.User.UpdateUser(ctx, userUpdateToUserInfoModel(string(userID), &userUpdate), userUpdate.Password); err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			render.Render(w, r, getErrRendererNotFound(errors.ErrResouceUser))
			return
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user")
		render.Render(w, r, getErrRendererServerError())
		return
	}

	render.JSON(w, r, nil)
}

// Delete a user
func (s *ServerHTTP) DeleteUsersUserID(w http.ResponseWriter, r *http.Request, userID UserID) {
	ctx := r.Context()

	// Validate request
	if err := userID.Validate(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong user ID")
		render.Render(w, r, getErrRendererBadRequest())
		return
	}

	// Delete user
	if err := s.domain.User.DeleteUser(ctx, modeluuid.FromStringOrNil(string(userID))); err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			render.Render(w, r, getErrRendererNotFound(errors.ErrResouceUser))
			return
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user")
		render.Render(w, r, getErrRendererServerError())
		return
	}

	render.JSON(w, r, nil)
}

// Get me
func (s *ServerHTTP) GetUsersMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID
	userID, err := middleware.GetUserIDFromCtx(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("No user ID in context")
		render.Render(w, r, getErrRendererServerError())
		return
	}

	// Get user
	userInfo, err := s.domain.User.GetUser(ctx, modeluuid.FromStringOrNil(string(userID)))
	if err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Msg("User doesn't exist")
			render.Render(w, r, getErrRendererNotFound(errors.ErrResouceUser))
			return
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get user")
		render.Render(w, r, getErrRendererServerError())
		return
	}

	render.JSON(w, r, UserModelToUserInfo(userInfo))
}

// Update me
func (s *ServerHTTP) PutUsersMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUpdate := UserUpdate{}

	// Unmarshal request
	if err := render.Bind(r, &userUpdate); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Wrong update user request")
		render.Render(w, r, getErrRendererBadRequest())
		return
	}

	// Get user ID
	userID, err := middleware.GetUserIDFromCtx(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("No user ID in context")
		render.Render(w, r, getErrRendererServerError())
		return
	}

	// Update user
	if err := s.domain.User.UpdateUser(ctx, userUpdateToUserInfoModel(string(userID), &userUpdate), userUpdate.Password); err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			render.Render(w, r, getErrRendererNotFound(errors.ErrResouceUser))
			return
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user")
		render.Render(w, r, getErrRendererServerError())
		return
	}

	render.JSON(w, r, nil)
}

// Delete me
func (s *ServerHTTP) DeleteUsersMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user ID
	userID, err := middleware.GetUserIDFromCtx(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("No user ID in context")
		render.Render(w, r, getErrRendererServerError())
		return
	}

	// Delete user
	if err := s.domain.User.DeleteUser(ctx, modeluuid.FromStringOrNil(string(userID))); err != nil {
		if err == service.ErrRepoNotFound {
			log.Ctx(ctx).Error().Err(err).Msg("User doesn't exist")
			render.Render(w, r, getErrRendererNotFound(errors.ErrResouceUser))
			return
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete user")
		render.Render(w, r, getErrRendererServerError())
		return
	}

	render.JSON(w, r, nil)
}

// Validate & Bind
func (u *UserID) Validate() error {
	return request.ValidateUserUUID(string(*u))
}

func (u *UserCreate) Bind(r *http.Request) error {
	return request.ValidateUserCreate(u.LoginId, u.Password, string(u.Role), u.Phone, u.Email)
}

func (u *UserUpdate) Bind(r *http.Request) error {
	return request.ValidateUserUpdate("", u.Password, string(u.Role), u.Phone, u.Email)
}

// DTO <-> Model
func userCreateToUserInfoModel(userCreate *UserCreate) *model.UserInfo {
	return &model.UserInfo{
		LoginID: userCreate.LoginId,
		Role:    model.UserRole(userCreate.Role),
		Phone:   userCreate.Phone,
		Email:   userCreate.Email,
	}
}

func userUpdateToUserInfoModel(userID string, userUpdate *UserUpdate) *model.UserInfo {
	return &model.UserInfo{
		ID:    modeluuid.FromStringOrNil(userID),
		Role:  model.UserRole(userUpdate.Role),
		Phone: userUpdate.Phone,
		Email: userUpdate.Email,
	}
}

func UserModelToUserInfo(userModel *model.UserInfo) *UserInfo {
	return &UserInfo{
		Id:      userModel.ID.String(),
		LoginId: userModel.LoginID,
		Role:    UserRole(userModel.Role),
		Phone:   userModel.Phone,
		Email:   userModel.Email,
	}
}

func UserModelListToUserInfoList(userModelList []model.UserInfo) []UserInfo {
	userInfos := []UserInfo{}
	for _, userModel := range userModelList {
		tmp := UserInfo{
			Id:      userModel.ID.String(),
			LoginId: userModel.LoginID,
			Role:    UserRole(userModel.Role),
			Phone:   userModel.Phone,
			Email:   userModel.Email,
		}
		userInfos = append(userInfos, tmp)
	}
	return userInfos
}
