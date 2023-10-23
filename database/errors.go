package database

import "errors"

var ErrDataExists = errors.New("表中已存在数据")
var ErrDataNotUnique = errors.New("数据不唯一")
var ErrDataExceeded = errors.New("数据超限")

// var errDataNotFound = errors.New("未找到数据")

// var errBackupFailed = errors.New("归档数据库失败")
