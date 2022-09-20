package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	errs "user-service/error"
	"user-service/model"
	serviceMock "user-service/service/mock"
)

func Test_Create_Should_Return_200_And_User_When_Nothing_Fails(t *testing.T) {
	var createViewModel = model.CreateUserViewModel{
		Email:    "batuhan@site.com",
		Name:     "Batuhan",
		Password: "123456",
	}

	var createDomainModel = model.CreateUserDomainModel{
		Email:    createViewModel.Email,
		Name:     createViewModel.Name,
		Password: createViewModel.Password,
	}

	var domainModel = model.UserDomainModel{
		Id:    primitive.NewObjectID().Hex(),
		Email: createDomainModel.Email,
		Name:  createDomainModel.Name,
	}

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("Create", mock.Anything, createDomainModel).Return(&domainModel, nil).Once()

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	requestBody, _ := json.Marshal(createViewModel)
	ctx.Request = &http.Request{Body: io.NopCloser(bytes.NewBuffer(requestBody))}

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.Create(ctx)

	var user model.UserDomainModel
	json.NewDecoder(responseRecorder.Result().Body).Decode(&user)

	assert.Equal(t, ctx.Writer.Status(), 200)
	assert.Equal(t, user.Email, domainModel.Email)
	assert.Equal(t, user.Name, domainModel.Name)
	userServiceMock.AssertExpectations(t)

}

func Test_Create_Should_Return_400_And_BadRequestError_When_Email_Is_Invalid(t *testing.T) {
	var email = "not an email"

	var createViewModel = model.CreateUserViewModel{
		Name:     "something",
		Email:    email,
		Password: "something",
	}

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("Create", mock.Anything, mock.Anything, mock.Anything).Maybe().Times(0)

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	requestBody, _ := json.Marshal(createViewModel)
	ctx.Request = &http.Request{Body: io.NopCloser(bytes.NewBuffer(requestBody))}

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.Create(ctx)

	var err map[string]string
	json.NewDecoder(responseRecorder.Result().Body).Decode(&err)

	assert.Equal(t, ctx.Writer.Status(), 400)
	userServiceMock.AssertExpectations(t)
}

func Test_Create_Should_Return_409_And_EmailAlreadyInUseError_When_Email_Already_Exists(t *testing.T) {
	var email = "actual@email.com"

	var createViewModel = model.CreateUserViewModel{
		Name:     "something",
		Email:    email,
		Password: "something",
	}

	var createDomainModel = model.CreateUserDomainModel{
		Name:     createViewModel.Name,
		Email:    createViewModel.Email,
		Password: createViewModel.Password,
	}

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("Create", mock.Anything, createDomainModel).Return(nil, errs.EmailAlreadyInUseError).Once()

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	requestBody, _ := json.Marshal(createViewModel)
	ctx.Request = &http.Request{Body: io.NopCloser(bytes.NewBuffer(requestBody))}

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.Create(ctx)

	var err map[string]string
	json.NewDecoder(responseRecorder.Result().Body).Decode(&err)

	assert.Equal(t, ctx.Writer.Status(), 409)
	assert.Equal(t, err["error"], errs.EmailAlreadyInUseError.Error())
	userServiceMock.AssertExpectations(t)
}

func Test_GetById_Should_Return_200_And_User_When_Nothing_Fails(t *testing.T) {
	var id = primitive.NewObjectID()

	var userDomainModel = model.UserDomainModel{
		Id: id.Hex(),
	}

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("GetById", mock.Anything, id.Hex()).Return(&userDomainModel, nil).Once()

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	ctx.AddParam("id", id.Hex())

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.GetById(ctx)

	var user model.UserDomainModel
	json.NewDecoder(responseRecorder.Result().Body).Decode(&user)

	assert.Equal(t, ctx.Writer.Status(), 200)
	assert.Equal(t, user.Id, id.Hex())
	userServiceMock.AssertExpectations(t)

}

func Test_GetById_Should_Return_400_When_Id_Is_Invalid(t *testing.T) {
	var id = "an invalid id"

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("GetById", mock.Anything, id).Return(nil, errs.BadRequestError).Once()

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	ctx.AddParam("id", id)

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.GetById(ctx)

	var err map[string]string
	json.NewDecoder(responseRecorder.Result().Body).Decode(&err)

	assert.Equal(t, ctx.Writer.Status(), 400)
	assert.Equal(t, err["error"], errs.BadRequestError.Error())
	userServiceMock.AssertExpectations(t)

}

func Test_GetById_Should_Return_404_When_User_Is_Not_Found(t *testing.T) {
	var id = primitive.NewObjectID()

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("GetById", mock.Anything, id.Hex()).Return(nil, errs.NotFoundError).Once()

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	ctx.AddParam("id", id.Hex())

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.GetById(ctx)

	var err map[string]string
	json.NewDecoder(responseRecorder.Result().Body).Decode(&err)

	assert.Equal(t, ctx.Writer.Status(), 404)
	assert.Equal(t, err["error"], errs.NotFoundError.Error())
	userServiceMock.AssertExpectations(t)

}

func Test_GetById_Should_Return_500_When_Fetch_Operation_Fails(t *testing.T) {
	var id = primitive.NewObjectID()

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("GetById", mock.Anything, id.Hex()).Return(nil, errs.ServerError).Once()

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	ctx.AddParam("id", id.Hex())

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.GetById(ctx)

	var err map[string]string
	json.NewDecoder(responseRecorder.Result().Body).Decode(&err)

	assert.Equal(t, ctx.Writer.Status(), 500)
	assert.Equal(t, err["error"], errs.ServerError.Error())
	userServiceMock.AssertExpectations(t)

}

func Test_GetAll_Should_Return_200_And_Users_When_Nothing_Fails(t *testing.T) {
	var firstUser = model.UserDomainModel{
		Id: primitive.NewObjectID().Hex(),
	}

	var secondUser = model.UserDomainModel{
		Id: primitive.NewObjectID().Hex(),
	}

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("GetAll", mock.Anything).Return([]*model.UserDomainModel{&firstUser, &secondUser}, nil).Once()

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.GetAll(ctx)

	var users []*model.UserViewModel
	json.NewDecoder(responseRecorder.Result().Body).Decode(&users)

	assert.Equal(t, ctx.Writer.Status(), 200)
	assert.Equal(t, users[0].Id, firstUser.Id)
	assert.Equal(t, users[1].Id, secondUser.Id)
	userServiceMock.AssertExpectations(t)
}

func Test_DeleteById_Should_Return_200_When_Nothing_Fails(t *testing.T) {
	var id = primitive.NewObjectID()

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("DeleteById", mock.Anything, id.Hex()).Return(nil).Once()

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	ctx.AddParam("id", id.Hex())

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.DeleteById(ctx)

	assert.Equal(t, ctx.Writer.Status(), 200)
	userServiceMock.AssertExpectations(t)
}

func Test_DeleteById_Should_Return_400_When_Id_Is_Invalid(t *testing.T) {
	var id = "an invalid id"

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("DeleteById", mock.Anything, id).Return(errs.BadRequestError).Once()

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	ctx.AddParam("id", id)

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.DeleteById(ctx)

	var err map[string]string
	json.NewDecoder(responseRecorder.Result().Body).Decode(&err)

	assert.Equal(t, ctx.Writer.Status(), 400)
	assert.Equal(t, err["error"], errs.BadRequestError.Error())
	userServiceMock.AssertExpectations(t)
}

func Test_DeleteById_Should_Return_404_When_User_Is_Not_Found(t *testing.T) {
	var id = primitive.NewObjectID()

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("DeleteById", mock.Anything, id.Hex()).Return(errs.NotFoundError).Once()

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	ctx.AddParam("id", id.Hex())

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.DeleteById(ctx)

	var err map[string]string
	json.NewDecoder(responseRecorder.Result().Body).Decode(&err)

	assert.Equal(t, ctx.Writer.Status(), 404)
	assert.Equal(t, err["error"], errs.NotFoundError.Error())
	userServiceMock.AssertExpectations(t)
}

func Test_DeleteById_Should_Return_500_When_Delete_Operation_Fails(t *testing.T) {
	var id = primitive.NewObjectID()

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("DeleteById", mock.Anything, id.Hex()).Return(errs.ServerError).Once()

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	ctx.AddParam("id", id.Hex())

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.DeleteById(ctx)

	var err map[string]string
	json.NewDecoder(responseRecorder.Result().Body).Decode(&err)

	assert.Equal(t, ctx.Writer.Status(), 500)
	assert.Equal(t, err["error"], errs.ServerError.Error())
	userServiceMock.AssertExpectations(t)
}

func Test_UpdateById_Should_Return_200_And_User_When_Nothing_Fails(t *testing.T) {
	var id = primitive.NewObjectID()
	var email = "batuhan@site.com"

	var updateViewModel = model.UpdateUserViewModel{
		Email: &email,
	}

	var updateDomainModel = model.UpdateUserDomainModel{
		Email: updateViewModel.Email,
	}

	var domainModel = model.UserDomainModel{
		Email: *updateDomainModel.Email,
	}

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("UpdateById", mock.Anything, id.Hex(), updateDomainModel).Return(&domainModel, nil).Once()

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	ctx.AddParam("id", id.Hex())
	requestBody, _ := json.Marshal(updateViewModel)
	ctx.Request = &http.Request{Body: io.NopCloser(bytes.NewBuffer(requestBody))}

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.UpdateById(ctx)

	var user model.UserDomainModel
	json.NewDecoder(responseRecorder.Result().Body).Decode(&user)

	assert.Equal(t, ctx.Writer.Status(), 200)
	assert.Equal(t, user.Email, domainModel.Email)
	userServiceMock.AssertExpectations(t)

}

func Test_UpdateById_Should_Return_400_And_BadRequestError_When_Email_Is_Invalid(t *testing.T) {
	var id = primitive.NewObjectID()
	var email = "not an email"

	var updateViewModel = model.UpdateUserViewModel{
		Email: &email,
	}

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("UpdateById", mock.Anything, mock.Anything, mock.Anything).Maybe().Times(0)

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	ctx.AddParam("id", id.Hex())
	requestBody, _ := json.Marshal(updateViewModel)
	ctx.Request = &http.Request{Body: io.NopCloser(bytes.NewBuffer(requestBody))}

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.UpdateById(ctx)

	var err map[string]string
	json.NewDecoder(responseRecorder.Result().Body).Decode(&err)

	assert.Equal(t, ctx.Writer.Status(), 400)
	userServiceMock.AssertExpectations(t)
}

func Test_UpdateById_Should_Return_409_And_EmailAlreadyInUseError_When_Email_Already_Exists(t *testing.T) {
	var id = primitive.NewObjectID()
	var email = "actual@email.com"

	var updateViewModel = model.UpdateUserViewModel{
		Email: &email,
	}

	var updateDomainModel = model.UpdateUserDomainModel{
		Email: updateViewModel.Email,
	}

	userServiceMock := new(serviceMock.UserServiceInterface)
	userServiceMock.On("UpdateById", mock.Anything, id.Hex(), updateDomainModel).Return(nil, errs.EmailAlreadyInUseError).Once()

	responseRecorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(responseRecorder)
	ctx.AddParam("id", id.Hex())
	requestBody, _ := json.Marshal(updateViewModel)
	ctx.Request = &http.Request{Body: io.NopCloser(bytes.NewBuffer(requestBody))}

	classUnderTest := NewUserController(userServiceMock, validator.New())
	classUnderTest.UpdateById(ctx)

	var err map[string]string
	json.NewDecoder(responseRecorder.Result().Body).Decode(&err)

	assert.Equal(t, ctx.Writer.Status(), 409)
	assert.Equal(t, err["error"], errs.EmailAlreadyInUseError.Error())
	userServiceMock.AssertExpectations(t)
}
