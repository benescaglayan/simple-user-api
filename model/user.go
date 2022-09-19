package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserEntity struct {
	Id       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type UserDomainModel struct {
	Id    string
	Name  string
	Email string
}

type UserViewModel struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserViewModel struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateUserDomainModel struct {
	Name     string
	Email    string
	Password string
}

type UpdateUserViewModel struct {
	Name     *string `json:"name"`
	Email    *string `json:"email" validate:"email"`
	Password *string `json:"password"`
}

type UpdateUserDomainModel struct {
	Name     *string
	Email    *string
	Password *string
}
