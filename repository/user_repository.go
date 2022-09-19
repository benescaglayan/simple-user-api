package repository

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	errs "user-service/error"
	"user-service/model"
)

const (
	CollectionName = "User"
)

type UserRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		userCollection: database.Collection(CollectionName),
	}
}

type UserRepositoryInterface interface {
	Create(*gin.Context, model.UserEntity) error
	GetById(*gin.Context, primitive.ObjectID) (*model.UserEntity, error)
	CheckIfEmailAlreadyInUse(*gin.Context, string) (bool, error)
	GetAll(*gin.Context) ([]*model.UserEntity, error)
	DeleteById(*gin.Context, primitive.ObjectID) error
	UpdateById(*gin.Context, primitive.ObjectID, model.UpdateUserDomainModel) (*model.UserEntity, error)
}

func (r *UserRepository) Create(ctx *gin.Context, user model.UserEntity) error {
	_, err := r.userCollection.InsertOne(ctx, user)
	if err != nil {
		log.Println(err)
		return errs.ServerError
	}

	return nil
}

func (r *UserRepository) GetById(ctx *gin.Context, id primitive.ObjectID) (user *model.UserEntity, err error) {
	filter := bson.D{{"_id", id}}

	result := r.userCollection.FindOne(ctx, filter).Decode(&user)
	if result == mongo.ErrNoDocuments {
		return nil, errs.NotFoundError
	} else if result != nil {
		log.Println(err)
		return nil, errs.ServerError
	}

	return
}

func (r *UserRepository) GetAll(ctx *gin.Context) (users []*model.UserEntity, err error) {
	cur, err := r.userCollection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Println(err)
		return users, errs.ServerError
	}

	for cur.Next(ctx) {
		var user model.UserEntity
		err := cur.Decode(&user)
		if err != nil {
			log.Println(err)
		}

		users = append(users, &user)
	}

	return
}

func (r *UserRepository) DeleteById(ctx *gin.Context, id primitive.ObjectID) error {
	filter := bson.D{{"_id", id}}

	result, err := r.userCollection.DeleteOne(ctx, filter)
	if err != nil {
		return errs.ServerError
	} else if result.DeletedCount == 0 {
		return errs.NotFoundError
	}

	return nil
}

func (r *UserRepository) CheckIfEmailAlreadyInUse(ctx *gin.Context, email string) (bool, error) {
	filter := bson.D{{"email", email}}

	count, err := r.userCollection.CountDocuments(ctx, filter)
	if err != nil {
		log.Println(err)
		return false, errs.ServerError
	}

	if count == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (r *UserRepository) UpdateById(ctx *gin.Context, id primitive.ObjectID, domainModel model.UpdateUserDomainModel) (*model.UserEntity, error) {
	var fieldsToUpdate bson.D

	if domainModel.Name != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.E{Key: "name", Value: domainModel.Name})
	}
	if domainModel.Email != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.E{Key: "email", Value: domainModel.Email})
	}
	if domainModel.Password != nil {
		fieldsToUpdate = append(fieldsToUpdate, bson.E{Key: "password", Value: domainModel.Email})
	}
	
	filter := bson.D{{"_id", id}}

	result, err := r.userCollection.UpdateOne(ctx, filter, bson.D{{"$set", fieldsToUpdate}})
	if err != nil {
		log.Println(err)
		return nil, errs.ServerError
	} else if result.ModifiedCount == 0 {
		return nil, errs.NotFoundError
	}

	user, err := r.GetById(ctx, id)

	return user, nil
}
