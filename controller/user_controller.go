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

	domainModel, err := c.userService.GetById(ctx, id)
	if err != nil {
		configureErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, copyDomainModelToViewModel(domainModel))
}

func (c *UserController) Create(ctx *gin.Context) {
	var createViewModel model.CreateUserViewModel

	err := ctx.BindJSON(&createViewModel)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
		return
	}

	err = c.validator.Struct(createViewModel)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	domainModel, err := c.userService.Create(ctx, copyCreateViewModelToCreateDomainModel(&createViewModel))
	if err != nil {
		configureErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, copyDomainModelToViewModel(domainModel))
}

func (c *UserController) GetAll(ctx *gin.Context) {
	domainModels, err := c.userService.GetAll(ctx)
	if err != nil {
		configureErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, copyDomainModelsToViewModels(domainModels))
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

	var updateViewModel model.UpdateUserViewModel

	err := ctx.BindJSON(&updateViewModel)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})
		return
	}

	err = c.validator.Struct(updateViewModel)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	domainModel, err := c.userService.UpdateById(ctx, id, copyUpdateViewModelToUpdateDomainModel(&updateViewModel))
	if err != nil {
		configureErrorResponse(ctx, err)
		return
	}

	ctx.IndentedJSON(http.StatusOK, copyDomainModelToViewModel(domainModel))
}

func configureErrorResponse(ctx *gin.Context, err error) {
	if errors.Is(err, errs.BadRequestError) {
		ctx.IndentedJSON(http.StatusBadRequest, map[string]string{"error": errs.BadRequestError.Error()})
		return
	} else if errors.Is(err, errs.NotFoundError) {
		ctx.IndentedJSON(http.StatusNotFound, map[string]string{"error": errs.NotFoundError.Error()})
		return
	} else if errors.Is(err, errs.EmailAlreadyInUseError) {
		ctx.IndentedJSON(http.StatusConflict, map[string]string{"error": errs.EmailAlreadyInUseError.Error()})
		return
	} else if errors.Is(err, errs.ServerError) {
		ctx.IndentedJSON(http.StatusInternalServerError, map[string]string{"error": errs.ServerError.Error()})
		return
	}
}

func copyDomainModelToViewModel(domainModel *model.UserDomainModel) model.UserViewModel {
	return model.UserViewModel{
		Id:    domainModel.Id,
		Name:  domainModel.Name,
		Email: domainModel.Email,
	}
}

func copyCreateViewModelToCreateDomainModel(viewModel *model.CreateUserViewModel) model.CreateUserDomainModel {
	return model.CreateUserDomainModel{
		Name:     viewModel.Name,
		Email:    viewModel.Email,
		Password: viewModel.Password,
	}
}

func copyUpdateViewModelToUpdateDomainModel(viewModel *model.UpdateUserViewModel) model.UpdateUserDomainModel {
	return model.UpdateUserDomainModel{
		Name:     viewModel.Name,
		Email:    viewModel.Email,
		Password: viewModel.Password,
	}
}

func copyDomainModelsToViewModels(domainModels []*model.UserDomainModel) (viewModels []model.UserViewModel) {
	for i := 0; i < len(domainModels); i++ {
		var viewModel model.UserViewModel

		viewModel.Id = domainModels[i].Id
		viewModel.Name = domainModels[i].Name
		viewModel.Email = domainModels[i].Email

		viewModels = append(viewModels, viewModel)
	}

	return
}
