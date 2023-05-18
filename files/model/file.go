package model

import "mime/multipart"

type File struct {
	Name        string               `form:"name" binding:"required"`
	Description string               `form:"description" binding:"required"`
	Data        multipart.FileHeader `form:"file"`
}
