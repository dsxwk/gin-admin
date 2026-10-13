package facade

import (
	"gin/pkg/serviceprovider/filesystem"
)

// File 文件系统管理器
// 使用示例:
//
//	facade.File().Disk("local").Upload(file, "avatar/1.png")
//	facade.File().Disk("oss").Upload(file, "avatar/1.png")
func File() *filesystem.Manager {
	return App().File()
}
