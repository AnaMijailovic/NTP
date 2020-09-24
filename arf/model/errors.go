package model

type UnableToRenameFileError struct {
	Err string
}

func (e UnableToRenameFileError) Error() string {
	return e.Err
}