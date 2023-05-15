package cli

import (
	"io/fs"
	"time"
    _ "path/filepath"
    "github.com/hornedout-law/tmkdir/tcron"
)


type Flags struct {
	Day   *int
	Week  *int
	Month *int
	Hour  *int
	Date  *string
	// experimental
	Time *string
}



type Job struct {
    name string
    task tcron.Task
    created_at time.Time
    last_exec time.Duration
    next_exec time.Time
}
type Tmkdir interface {

    createFsEntry(filename string) (fs.FileInfo, error)
    removeFsEntry(filepath string) error
    listTcronFiles() ([]fs.FileInfo, error)
}

