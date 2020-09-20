package model

import "time"

type DeleteData struct {
	Path string
	Recursive bool
	Empty bool
	CreatedBefore time.Time
	NotAccessedAfter time.Time
}