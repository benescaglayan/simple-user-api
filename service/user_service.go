package service

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	errs "user-service/error"
	"user-service/model"
	"user-service/repository"
)

type UserService struct {
	userRepository repository.UserRepositoryInterface
}

func NewUserService(userRepository repository.UserRepositoryInterface) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

type UserServiceInterface interface {
	Create(*gin.Context, model.CreateUserDomainModel) (*model.UserDomainModel, error)
	GetById(*gin.Context, string) (*model.UserDomainModel, error)
	GetAll(*gin.Context) ([]*model.UserViewModel, error)
	DeleteById(*gin.Context, string) error
	UpdateById(*gin.Context, string, model.UpdateUserDomainModel) (*model.UserViewModel, error)
}

func (s *UserService) Create(ctx *gin.Context, createModel model.CreateUserDomainModel) (*model.UserDomainModel, error) {
	isEmailInUse, err := s.userRepository.CheckIfEmailAlreadyInUse(ctx, createModel.Email)
	if isEmailInUse {
		return nil, errs.EmailAlreadyInUseError
	} else if err != nil {
		return nil, err
	}

	hashedPasswordInBytes, err := bcrypt.GenerateFromPassword([]byte(createModel.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		return nil, errs.ServerError
	}

	entity := model.UserEntity{
		Id:       primitive.NewObjectID(),
		Name:     createModel.Name,
		Email:    createModel.Email,
		Password: string(hashedPasswordInBytes),
	}

	err = s.userRepository.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	user := model.UserDomainModel{
		Id:    entity.Id.Hex(),
		Email: entity.Email,
		Name:  entity.Name,
	}

	return &user, nil
}

func (s *UserService) GetById(ctx *gin.Context, id string) (user *model.UserDomainModel, err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return nil, errs.BadRequestError
	}

	userEntity, err := s.userRepository.GetById(ctx, objectId)
	if err != nil {
		return nil, err
	}

	user = &model.UserDomainModel{
		Id:    userEntity.Id.Hex(),
		Name:  userEntity.Name,
		Email: userEntity.Email,
	}

	return
}

func (s *UserService) GetAll(ctx *gin.Context) (userViews []*model.UserViewModel, err error) {
	userEntities, err := s.userRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(userEntities) == 0 {
		return nil, errs.NotFoundError
	}

	for i := 0; i < len(userEntities); i++ {
		var userView model.UserViewModel

		userView.Id = userEntities[i].Id.Hex()
		userView.Name = userEntities[i].Name
		userView.Email = userEntities[i].Email

		userViews = append(userViews, &userView)
	}

	return
}

func (s *UserService) DeleteById(ctx *gin.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return errs.BadRequestError
	}

	err = s.userRepository.DeleteById(ctx, objectId)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateById(ctx *gin.Context, id string, updateModel model.UpdateUserDomainModel) (*model.UserViewModel, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err)
		return nil, errs.BadRequestError
	}

	if updateModel.Email != nil {
		isEmailInUse, err := s.userRepository.CheckIfEmailAlreadyInUse(ctx, *(updateModel.Email))
		if isEmailInUse {
			return nil, errs.EmailAlreadyInUseError
		} else if err != nil {
			return nil, err
		}
	}

	if updateModel.Password != nil {
		hashedPasswordInBytes, err := bcrypt.GenerateFromPassword([]byte(*updateModel.Password), bcrypt.MinCost)
		if err != nil {
			log.Println(err)
			return nil, errs.ServerError
		}

		var password = string(hashedPasswordInBytes)
		updateModel.Password = &password
	}

	userEntity, err := s.userRepository.UpdateById(ctx, objectId, updateModel)
	if err != nil {
		return nil, err
	}

	uv := &model.UserViewModel{
		Id:    userEntity.Id.Hex(),
		Name:  userEntity.Name,
		Email: userEntity.Email,
	}

	return uv, nil
}
