package database

import (
	"path/filepath"
	"runtime"
	"strconv"
)

const (
	// baseDir         = "userdata"
	// dbSubDir        = "dbfiles"
	// dbArchSubDir    = "archives"
	dbFileExtension = ".db"
)

var (
	rootDir   string
	baseDir   string
	dbSubDir  = "dbfiles"
	dbDir     string
	dbArchDir string
	// dbDir     = filepath.Join(baseDir, dbSubDir)     // /userdata/dbfiles
	// dbArchDir = filepath.Join(baseDir, dbArchSubDir) // /userdata/archives
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	rootDir = filepath.Join(filepath.Dir(b), "..")
	baseDir = filepath.Join(rootDir, "userdata")
	dbDir = filepath.Join(baseDir, dbSubDir)
	dbArchDir = filepath.Join(baseDir, "archives")
}

// /userdata/dbfiles/{grpid}.go
func getDBFilePath(grpid int) string {
	return filepath.Join(dbDir, strconv.Itoa(grpid)+dbFileExtension)
}
