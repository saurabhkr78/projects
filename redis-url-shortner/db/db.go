package db

import (
	urlModel "redis-url-shortner/model"
)

var URLs = []urlModel.URL{
	{
		ShortID:     "aB92x",
		OriginalURL: "https://www.facebook.com/login/",
	},
	{
		ShortID:     "xY123",
		OriginalURL: "https://www.google.com",
	},
}
