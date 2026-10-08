package tests

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"path"
	"testing"

	"gin/app/facade"
	"gin/pkg/serviceprovider/filesystem"

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

// TestDiskNotFound 文件磁盘不存在测试
func TestDiskNotFound(t *testing.T) {
	require.NotNil(t, facade.File())

	_, err := facade.File().Disk("not-exists").Upload([]byte("hello"), "tests/filesystem/not-exists.txt")
	require.ErrorIs(t, err, filesystem.ErrDiskNotFound)
}
