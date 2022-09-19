package service

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	errs "user-service/error"
	"user-service/model"
	repositoryMock "user-service/repository/mock"
)

func Test_Create_Should_Return_EmailAlreadyInUseError_When_Email_Belongs_To_A_User(t *testing.T) {
	request := model.CreateUserDomainModel{
		Email: "existing@email.com",
	}

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("CheckIfEmailAlreadyInUse", mock.Anything, request.Email).Return(true, nil).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	createdUser, err := classUnderTest.Create(&gin.Context{}, request)

	assert.Nil(t, createdUser)
	assert.NotNil(t, err)
	assert.Equal(t, errs.EmailAlreadyInUseError, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_Create_Should_Return_ServerError_When_Email_Check_Fails(t *testing.T) {
	request := model.CreateUserDomainModel{
		Email: "existing@email.com",
	}

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("CheckIfEmailAlreadyInUse", mock.Anything, request.Email).Return(false, errs.ServerError).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	createdUser, err := classUnderTest.Create(&gin.Context{}, request)

	assert.Nil(t, createdUser)
	assert.NotNil(t, err)
	assert.Equal(t, errs.ServerError, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_Create_Should_Return_ServerError_When_Database_Create_Operation_Fails(t *testing.T) {
	request := model.CreateUserDomainModel{
		Name:     "Batuhan",
		Email:    "non_existing@email.com",
		Password: "123456",
	}

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("CheckIfEmailAlreadyInUse", mock.Anything, request.Email).Return(false, nil).Once()
	userRepositoryMock.On("Create", mock.Anything, mock.MatchedBy(func(i interface{}) bool {
		return i.(model.UserEntity).Name == request.Name && i.(model.UserEntity).Email == request.Email
	})).Return(errs.ServerError).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	createdUser, err := classUnderTest.Create(&gin.Context{}, request)

	assert.Nil(t, createdUser)
	assert.NotNil(t, err)
	assert.Equal(t, errs.ServerError, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_Create_Should_Return_User_When_Nothing_Fails(t *testing.T) {
	request := model.CreateUserDomainModel{
		Name:     "Batuhan",
		Email:    "non_existing@email.com",
		Password: "123456",
	}

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("CheckIfEmailAlreadyInUse", mock.Anything, request.Email).Return(false, nil).Once()
	userRepositoryMock.On("Create", mock.Anything, mock.MatchedBy(func(i interface{}) bool {
		return i.(model.UserEntity).Name == request.Name && i.(model.UserEntity).Email == request.Email
	})).Return(nil).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	createdUser, err := classUnderTest.Create(&gin.Context{}, request)

	assert.Nil(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, createdUser.Name, request.Name)
	assert.Equal(t, createdUser.Email, request.Email)
	userRepositoryMock.AssertExpectations(t)
}

func Test_GetById_Should_Return_BadRequestError_When_Id_Is_Invalid(t *testing.T) {
	var id = "not an object id"

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("GetById", mock.Anything, mock.Anything).Maybe().Times(0)

	classUnderTest := NewUserService(userRepositoryMock)

	user, err := classUnderTest.GetById(&gin.Context{}, id)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, errs.BadRequestError, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_GetById_Should_Return_ServerError_When_Id_Is_Valid_But_Find_Operation_Fails(t *testing.T) {
	var id = primitive.NewObjectID()

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("GetById", mock.Anything, id).Return(nil, errs.ServerError).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	user, err := classUnderTest.GetById(&gin.Context{}, id.Hex())

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, errs.ServerError, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_GetById_Should_Return_User_When_Id_Is_Valid_And_Belongs_To_User(t *testing.T) {
	var id = primitive.NewObjectID()

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("GetById", mock.Anything, id).Return(&model.UserEntity{Id: id}, nil).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	user, err := classUnderTest.GetById(&gin.Context{}, id.Hex())

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, id.Hex(), user.Id)
	userRepositoryMock.AssertExpectations(t)
}

func Test_DeleteById_Should_Return_BadRequestError_When_Id_Is_Invalid(t *testing.T) {
	var id = "not an object id"

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("DeleteById", mock.Anything, mock.Anything).Maybe().Times(0)

	classUnderTest := NewUserService(userRepositoryMock)

	err := classUnderTest.DeleteById(&gin.Context{}, id)

	assert.NotNil(t, err)
	assert.Equal(t, errs.BadRequestError, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_DeleteById_Should_Return_ServerError_When_Id_Is_Valid_But_Database_Fails(t *testing.T) {
	var id = primitive.NewObjectID()

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("DeleteById", mock.Anything, id).Return(errs.ServerError).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	err := classUnderTest.DeleteById(&gin.Context{}, id.Hex())

	assert.NotNil(t, err)
	assert.Equal(t, errs.ServerError, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_DeleteById_Should_Not_Return_Error_When_Id_Is_Valid_And_Belongs_To_User(t *testing.T) {
	var id = primitive.NewObjectID()

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("DeleteById", mock.Anything, id).Return(nil).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	err := classUnderTest.DeleteById(&gin.Context{}, id.Hex())

	assert.Nil(t, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_GetAll_Should_Return_NotFoundError_When_No_Users_Exist(t *testing.T) {
	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("GetAll", mock.Anything).Return([]*model.UserEntity{}, errs.NotFoundError).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	users, err := classUnderTest.GetAll(&gin.Context{})

	assert.Nil(t, users)
	assert.NotNil(t, err)
	assert.Equal(t, errs.NotFoundError, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_GetAll_Should_Return_ServerError_When_Database_Fails(t *testing.T) {
	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("GetAll", mock.Anything).Return([]*model.UserEntity{}, errs.ServerError).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	users, err := classUnderTest.GetAll(&gin.Context{})

	assert.Nil(t, users)
	assert.NotNil(t, err)
	assert.Equal(t, errs.ServerError, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_GetAll_Should_Return_Users_When_Nothing_Fails(t *testing.T) {
	var firstUser = model.UserEntity{
		Id:   primitive.NewObjectID(),
		Name: "First User",
	}

	var secondUser = model.UserEntity{
		Id:   primitive.NewObjectID(),
		Name: "Second User",
	}

	var userEntities = []*model.UserEntity{&firstUser, &secondUser}

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("GetAll", mock.Anything).Return(userEntities, nil).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	users, err := classUnderTest.GetAll(&gin.Context{})

	assert.Nil(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, firstUser.Id.Hex(), users[0].Id)
	assert.Equal(t, firstUser.Name, users[0].Name)
	assert.Equal(t, secondUser.Id.Hex(), users[1].Id)
	assert.Equal(t, secondUser.Name, users[1].Name)
	userRepositoryMock.AssertExpectations(t)
}

func Test_UpdateById_Should_Return_BadRequestError_When_Id_Is_Invalid(t *testing.T) {
	var id = "not an object id"

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("UpdateById", mock.Anything, mock.Anything).Maybe().Times(0)

	classUnderTest := NewUserService(userRepositoryMock)

	updatedUser, err := classUnderTest.UpdateById(&gin.Context{}, id, model.UpdateUserDomainModel{})

	assert.Nil(t, updatedUser)
	assert.NotNil(t, err)
	assert.Equal(t, errs.BadRequestError, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_UpdateById_Should_Return_EmailAlreadyInUseError_When_Email_Is_Sent_And_Belongs_To_A_User(t *testing.T) {
	var id = primitive.NewObjectID()
	var email = "existing@email.com"

	var updateModel = model.UpdateUserDomainModel{
		Email: &email,
	}

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("UpdateById", mock.Anything, mock.Anything).Maybe().Times(0)
	userRepositoryMock.On("CheckIfEmailAlreadyInUse", mock.Anything, *(updateModel.Email)).Return(true, nil).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	updatedUser, err := classUnderTest.UpdateById(&gin.Context{}, id.Hex(), updateModel)

	assert.Nil(t, updatedUser)
	assert.NotNil(t, err)
	assert.Equal(t, errs.EmailAlreadyInUseError, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_UpdateById_Should_Return_ServerError_When_Email_Check_Fails(t *testing.T) {
	var id = primitive.NewObjectID()
	var email = "non_existing@email.com"

	var updateModel = model.UpdateUserDomainModel{
		Email: &email,
	}

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("UpdateById", mock.Anything, mock.Anything).Maybe().Times(0)
	userRepositoryMock.On("CheckIfEmailAlreadyInUse", mock.Anything, *(updateModel.Email)).Return(false, errs.ServerError).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	updatedUser, err := classUnderTest.UpdateById(&gin.Context{}, id.Hex(), updateModel)

	assert.Nil(t, updatedUser)
	assert.NotNil(t, err)
	assert.Equal(t, errs.ServerError, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_UpdateById_Should_Return_ServerError_When_Update_Operation_Fails(t *testing.T) {
	var id = primitive.NewObjectID()
	var email = "non_existing@email.com"

	var updateModel = model.UpdateUserDomainModel{
		Email: &email,
	}

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("UpdateById", mock.Anything, mock.Anything).Maybe().Times(0)
	userRepositoryMock.On("CheckIfEmailAlreadyInUse", mock.Anything, *(updateModel.Email)).Return(false, nil).Once()
	userRepositoryMock.On("UpdateById", mock.Anything, id, updateModel).Return(nil, errs.ServerError).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	updatedUser, err := classUnderTest.UpdateById(&gin.Context{}, id.Hex(), updateModel)

	assert.Nil(t, updatedUser)
	assert.NotNil(t, err)
	assert.Equal(t, errs.ServerError, err)
	userRepositoryMock.AssertExpectations(t)
}

func Test_UpdateById_Should_Return_ServerError_When_Nothing_Fails(t *testing.T) {
	var id = primitive.NewObjectID()
	var email = "non_existing@email.com"

	var updateModel = model.UpdateUserDomainModel{
		Email: &email,
	}

	var userEntity = &model.UserEntity{
		Email: email,
	}

	userRepositoryMock := new(repositoryMock.UserRepositoryInterface)
	userRepositoryMock.On("UpdateById", mock.Anything, mock.Anything).Maybe().Times(0)
	userRepositoryMock.On("CheckIfEmailAlreadyInUse", mock.Anything, *(updateModel.Email)).Return(false, nil).Once()
	userRepositoryMock.On("UpdateById", mock.Anything, id, updateModel).Return(userEntity, nil).Once()

	classUnderTest := NewUserService(userRepositoryMock)

	updatedUser, err := classUnderTest.UpdateById(&gin.Context{}, id.Hex(), updateModel)

	assert.Nil(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, updatedUser.Email, email)
	userRepositoryMock.AssertExpectations(t)
}
