package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/colynn/pontus/tools"
)

// StorageImportFile ..
func StorageImportFile(f multipart.File, fileHeader *multipart.FileHeader) (filePath string, err error) {
	filePath = fmt.Sprintf("assets/%s", fileHeader.Filename)
	filePath = tools.EnsureAbs(filePath)
	err = storageUploadFile(filePath, f)
	if err != nil {
		return
	}
	return
}

func storageUploadFile(path string, file multipart.File) error {
	// verify fileName exist or not
	_, fullPathStat := os.Stat(path)
	if os.IsExist(fullPathStat) {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("清理旧文件时失败，请联系管理员后重试")
		}
	}

	// storage file into disk device
	osFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer osFile.Close()
	_, err = io.Copy(osFile, file)
	if err != nil {
		return fmt.Errorf("save conf file occur error: %v", err.Error())
	}
	return nil
}
