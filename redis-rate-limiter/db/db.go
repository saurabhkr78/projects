package db

import (
	userModel "redis-rate-limiter/model"
)

var DB = []userModel.User{
	{Id: "1", Name: "x"},
	{Id: "2", Name: "Y"},
	{Id: "3", Name: "Z"},
	{Id: "4", Name: "K"},
}
