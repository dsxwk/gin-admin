package filesystem

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"testing"
	"uuid"

	"gin/app/errcode"
	"gin/app/facade"
	"gin/app/model"
	"gin/pkg/pathx"

	"github.com/stretchr/testify/require"
)

// TestLocalUpload 本地文件上传测试
func TestLocalUpload(t *testing.T) {
	require.NotNil(t, facade.File())

	disk := facade.File().Disk("local")
	objectPath := "tests/filesystem/local-upload.txt"
	t.Cleanup(func() {
		_ = disk.Delete(objectPath)
	})

	result, err := disk.Upload([]byte("hello"), objectPath)
	require.NoError(t, err)
	require.Equal(t, "local", result.Disk)
	require.Equal(t, path.Join(facade.Config().Filesystem.Local.Root, objectPath), result.Path)
	require.Contains(t, result.URL, path.Join("/", facade.Config().Filesystem.Local.Root, objectPath))
	require.Equal(t, int64(5), result.Size)

	exists, err := disk.Exists(result.Path)
	require.NoError(t, err)
	require.True(t, exists)

	data, err := disk.Read(result.Path)
	require.NoError(t, err)
	require.Equal(t, []byte("hello"), data)
}

// TestMultipartUpload multipart文件上传测试
func TestMultipartUpload(t *testing.T) {
	require.NotNil(t, facade.File())

	disk := facade.File().Disk()
	objectPath := "tests/filesystem/multipart-upload.png"
	t.Cleanup(func() {
		_ = disk.Delete(objectPath)
	})

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "avatar.png")
	require.NoError(t, err)
	_, err = part.Write([]byte("image-data"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	request := httptest.NewRequest("POST", "/upload", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	require.NoError(t, request.ParseMultipartForm(1<<20))

	fileHeader := request.MultipartForm.File["file"][0]
	result, err := disk.Upload(fileHeader, objectPath)
	require.NoError(t, err)
	require.Equal(t, int64(len("image-data")), result.Size)

	data, err := disk.Read(result.Path)
	require.NoError(t, err)
	require.Equal(t, []byte("image-data"), data)
}

// TestFileRecordLifecycle 文件表记录生命周期测试
func TestFileRecordLifecycle(t *testing.T) {
	disk := facade.File().WithUploader(9).WithGroup(7).Disk()
	objectPath := "tests/filesystem/file-record.txt"
	require.NoError(t, disk.Delete(objectPath))
	t.Cleanup(func() {
		_ = disk.Delete(objectPath)
	})

	result, err := disk.Upload([]byte("record-data"), objectPath)
	require.NoError(t, err)
	require.NotEmpty(t, result.ObjectKey)
	require.NotEmpty(t, result.OriginName)

	db := facade.DB("mysql")
	require.NotNil(t, db)
	var count int64
	require.NoError(t, db.Model(&model.File{}).
		Where("storage = ? AND object_key = ?", result.Disk, result.ObjectKey).
		Count(&count).Error)
	require.Equal(t, int64(1), count)

	var record model.File
	require.NoError(t, db.Where("storage = ? AND object_key = ?", result.Disk, result.ObjectKey).
		First(&record).Error)
	require.Equal(t, int64(9), record.UploaderId)
	require.Equal(t, int64(7), record.GroupId)
	require.Equal(t, "file-record.txt", record.OriginName)

	require.NoError(t, disk.Delete(objectPath))
	count = 0
	require.NoError(t, db.Model(&model.File{}).
		Where("storage = ? AND object_key = ?", result.Disk, result.ObjectKey).
		Count(&count).Error)
	require.Zero(t, count)
}

// TestDeleteRemovesEmptyDirectories 删除文件后回收空目录测试
func TestDeleteRemovesEmptyDirectories(t *testing.T) {
	disk := facade.File().Disk()
	objectDir := path.Join("tests/filesystem", "cleanup-"+uuid.New().String())
	objectPath := path.Join(objectDir, "empty-dir.txt")
	t.Cleanup(func() {
		_ = disk.Delete(objectPath)
	})

	result, err := disk.Upload([]byte("cleanup"), objectPath)
	require.NoError(t, err)

	dirPath := filepath.Join(pathx.RootPath(), filepath.FromSlash(result.Path))
	dirPath = filepath.Dir(dirPath)
	require.DirExists(t, dirPath)

	require.NoError(t, disk.Delete(objectPath))
	_, err = os.Stat(dirPath)
	require.ErrorIs(t, err, os.ErrNotExist)
}

// TestDiskNotFound 文件磁盘不存在测试
func TestDiskNotFound(t *testing.T) {
	require.NotNil(t, facade.File())

	_, err := facade.File().Disk("not-exists").Upload([]byte("hello"), "tests/filesystem/not-exists.txt")
	require.ErrorIs(t, err, (errcode.FilesystemErrCode{}).DiskNotFound())
}
