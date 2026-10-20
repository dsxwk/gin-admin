package errcode

import (
	errCode "gin/pkg/errcode"
	"net/http"
)

const (
	// FilesystemErrCodePrefix 文件系统错误码前缀
	FilesystemErrCodePrefix = 900
)

// FilesystemErrCode 文件系统错误码
type FilesystemErrCode struct{}

// NotInitialized 文件系统服务未初始化
func (s FilesystemErrCode) NotInitialized() errCode.ErrorCode {
	return errCode.NewError(1, "文件系统服务未初始化").
		WithHttpCode(http.StatusInternalServerError).
		WithPrefix(FilesystemErrCodePrefix)
}

// DiskNotFound 文件磁盘不存在
func (s FilesystemErrCode) DiskNotFound() errCode.ErrorCode {
	return errCode.NewError(2, "文件磁盘不存在").
		WithHttpCode(http.StatusNotFound).
		WithPrefix(FilesystemErrCodePrefix)
}

// InvalidPath 文件路径不合法
func (s FilesystemErrCode) InvalidPath() errCode.ErrorCode {
	return errCode.NewError(3, "文件路径不合法").
		WithHttpCode(http.StatusBadRequest).
		WithPrefix(FilesystemErrCodePrefix)
}

// UnsupportedSource 不支持的上传来源
func (s FilesystemErrCode) UnsupportedSource() errCode.ErrorCode {
	return errCode.NewError(4, "不支持的上传来源").
		WithHttpCode(http.StatusBadRequest).
		WithPrefix(FilesystemErrCodePrefix)
}
