package database

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// 将在用的数据库文件复制(归档)至归档目录下
func archiveDBFile(grpid int) error {
	// check and create dir
	if _, err := os.Stat(dbArchDir); os.IsNotExist(err) {
		err = os.MkdirAll(dbArchDir, 0755)
		if err != nil {
			return err
		}
	}

	// get and check active db file full path
	src := getDBFilePath(grpid)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dst := filepath.Join(dbArchDir, getArchName(strconv.Itoa(grpid)))
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// {yyyymm}-{grpid}(_{seq}).db
func getArchName(grpid string) string {
	// 数据库文件归档前处理文件名，标明年-月份并(最大努力)防止覆盖
	currentTime := time.Now()
	yearMonth := currentTime.Format("200601")
	newName := fmt.Sprintf("%s-%s%s", grpid, yearMonth, dbFileExtension)
	counter := 1

	for {
		_, err := os.Stat(filepath.Join(dbArchDir, newName)) // /archived/1234.go
		if os.IsNotExist(err) {
			return newName
		} else if err != nil {
			return fmt.Sprintf("%s-%s%s", grpid, yearMonth, dbFileExtension)
		} else {
			newName = fmt.Sprintf("%s-%s_%d%s", grpid, yearMonth, counter, dbFileExtension)
			counter++
		}
	}
}
