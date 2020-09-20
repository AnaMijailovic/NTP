package model

type ReorganizeData struct {
	Src string
	Dest string
	Recursive bool
	FileType bool
	FileSize int64
	CreatedDate string
}
