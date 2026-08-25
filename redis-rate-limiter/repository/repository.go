package repository

import (
	userModel "redis-rate-limiter/model"
)

type UserRepository struct {
	DB *[]userModel.User //type of db is slice of user model
}

func NewUserRepository(db []userModel.User) *UserRepository {
	return &UserRepository{
		DB: &db,
	}
}
