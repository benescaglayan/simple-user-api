package mock

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"user-service/model"
)

type UserRepositoryInterface struct {
	mock.Mock
}

func (_m *UserRepositoryInterface) GetById(ctx *gin.Context, id primitive.ObjectID) (*model.UserEntity, error) {
	args := _m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.UserEntity), args.Error(1)
}

func (_m *UserRepositoryInterface) Create(ctx *gin.Context, user model.UserEntity) error {
	args := _m.Called(ctx, user)

	return args.Error(0)
}

func (_m *UserRepositoryInterface) CheckIfEmailAlreadyInUse(ctx *gin.Context, email string) (bool, error) {
	args := _m.Called(ctx, email)

	return args.Bool(0), args.Error(1)
}

func (_m *UserRepositoryInterface) GetAll(ctx *gin.Context) ([]*model.UserEntity, error) {
	args := _m.Called(ctx)

	return args.Get(0).([]*model.UserEntity), args.Error(1)
}

func (_m *UserRepositoryInterface) DeleteById(ctx *gin.Context, id primitive.ObjectID) error {
	args := _m.Called(ctx, id)

	return args.Error(0)
}

func (_m *UserRepositoryInterface) UpdateById(ctx *gin.Context, id primitive.ObjectID, updateModel model.UpdateUserDomainModel) (*model.UserEntity, error) {
	args := _m.Called(ctx, id, updateModel)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.UserEntity), args.Error(1)
}
