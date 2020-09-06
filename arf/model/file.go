package model

import (
	"time"
)

type File struct {
	Name string 	   `json:"name"`
	FullPath string    `json:"fullPath"`
	IsDir bool         `json:"isDir"`
	Size float64       `json:"size"`
	FileType string    `json:"fileType"`
	Created time.Time  `json:"created"`
	Modified time.Time `json:"modified"`
	Accessed time.Time `json:"accessed"`
}