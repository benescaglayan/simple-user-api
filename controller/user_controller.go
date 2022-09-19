package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	errs "user-service/error"
	"user-service/model"
	"user-service/service"
)

type UserController struct {
	userService service.UserServiceInterface
	validator   *validator.Validate
}

func NewUserController(userService service.UserServiceInterface, validator *validator.Validate) *UserController {
	return &UserController{
		userService: userService,
		validator:   validator,
	}
}

func (c *UserController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.userService.GetById(ctx, id)
	if err != nil {
		configureErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, user)
}

func (c *UserController) Create(ctx *gin.Context) {
	var createUserRequest model.CreateUserViewModel

	err := ctx.BindJSON(&createUserRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, map[string]interface{}{"error": "Bad request"})
		return
	}

	err = c.validator.Struct(createUserRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	createUserModel := model.CreateUserDomainModel{
		Name:     createUserRequest.Name,
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
	}

	user, err := c.userService.Create(ctx, createUserModel)
	if err != nil {
		configureErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, user)
}

func (c *UserController) GetUsers(ctx *gin.Context) {
	users, err := c.userService.GetAll(ctx)
	if err != nil {
		configureErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, users)
}

func (c *UserController) DeleteById(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.userService.DeleteById(ctx, id)
	if err != nil {
		configureErrorResponse(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *UserController) UpdateById(ctx *gin.Context) {
	id := ctx.Param("id")

	var updateUserRequest model.UpdateUserViewModel

	err := ctx.BindJSON(&updateUserRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, map[string]interface{}{"error": "Bad request"})
		return
	}

	err = c.validator.Struct(updateUserRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
		return
	}

	updateUserModel := model.UpdateUserDomainModel{
		Name:     updateUserRequest.Name,
		Email:    updateUserRequest.Email,
		Password: updateUserRequest.Password,
	}

	user, err := c.userService.UpdateById(ctx, id, updateUserModel)
	if err != nil {
		configureErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, user)
}

func configureErrorResponse(ctx *gin.Context, err error) {
	if errors.Is(err, errs.BadRequestError) {
		ctx.IndentedJSON(http.StatusBadRequest, map[string]interface{}{"error": errs.BadRequestError.Error()})
		return
	} else if errors.Is(err, errs.NotFoundError) {
		ctx.IndentedJSON(http.StatusNotFound, map[string]interface{}{"error": errs.NotFoundError.Error()})
		return
	} else if errors.Is(err, errs.EmailAlreadyInUseError) {
		ctx.IndentedJSON(http.StatusConflict, map[string]interface{}{"error": errs.EmailAlreadyInUseError.Error()})
		return
	} else if errors.Is(err, errs.ServerError) {
		ctx.IndentedJSON(http.StatusInternalServerError, map[string]interface{}{"error": errs.ServerError.Error()})
		return
	}
}
