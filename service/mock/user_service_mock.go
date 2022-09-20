package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"user-service/model"
)

type UserServiceInterface struct {
	mock.Mock
}

func (_m *UserServiceInterface) Create(ctx *gin.Context, createModel model.CreateUserDomainModel) (*model.UserDomainModel, error) {
	args := _m.Called(ctx, createModel)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.UserDomainModel), args.Error(1)
}

func (_m *UserServiceInterface) GetById(ctx *gin.Context, id string) (*model.UserDomainModel, error) {
	args := _m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.UserDomainModel), args.Error(1)
}

func (_m *UserServiceInterface) GetAll(ctx *gin.Context) ([]*model.UserViewModel, error) {
	args := _m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*model.UserViewModel), args.Error(1)
}

func (_m *UserServiceInterface) DeleteById(ctx *gin.Context, id string) error {
	args := _m.Called(ctx, id)

	return args.Error(0)
}

func (_m *UserServiceInterface) UpdateById(ctx *gin.Context, id string, updateModel model.UpdateUserDomainModel) (*model.UserViewModel, error) {
	args := _m.Called(ctx, id, updateModel)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.UserViewModel), args.Error(1)
}
