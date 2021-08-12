/**
 * @Author: Anpw
 * @Description:
 * @File:  uploadFile
 * @Version: 1.0.0
 * @Date: 2021/7/27 1:40
 */

package upload

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name      string
	AccessUrl string
	Hash      string
	Dst       string
}

func Upload(fileType FileType, file multipart.File, fileHeader *multipart.FileHeader, path string) (*FileInfo, error) {
	fileName := GetFileName(fileHeader.Filename)
	uploadSavePath := GetSavePath()
	uploadSavePath = uploadSavePath + "/" + path
	dst := uploadSavePath + "/" + fileName
	if !CheckContainExt(fileType, fileName) {
		return nil, errors.New("不支持该文件")
	}
	if CheckSavePath(uploadSavePath) {
		if err := CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
			return nil, errors.New("创建文件保存目录失败 ")
		}
	}
	if CheckPermission(uploadSavePath) {
		return nil, errors.New("文件权限不足 ")
	}
	if err := SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}
	hash, err := FileNameHash(dst)
	if err != nil {
		logs.Error("计算文件哈希值失败 err:", err)
		return nil, err
	}
	accessUrl := beego.AppConfig.String("UploadServerUrl") + "/" + fileName
	return &FileInfo{
		Name:      fileName,
		AccessUrl: accessUrl,
		Hash:      hash,
		Dst:       dst,
	}, nil
}

func Update(fileType FileType, file multipart.File, fileHeader *multipart.FileHeader, filePath string) (*FileInfo, error) {
	fileName := GetFileName(fileHeader.Filename)
	path, oldFileName := filepath.Split(filePath)
	uploadSavePath := path
	dst := uploadSavePath + oldFileName
	newdst := uploadSavePath + fileName
	if !CheckContainExt(fileType, fileName) {
		return nil, errors.New("不支持该文件")
	}
	if CheckSavePath(uploadSavePath) {
		if err := CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
			return nil, errors.New("创建文件保存目录失败 ")
		}
	}
	if CheckPermission(uploadSavePath) {
		return nil, errors.New("文件权限不足 ")
	}
	if err := UpdateFile(fileHeader, dst, newdst); err != nil {
		return nil, err
	}
	hash, err := FileNameHash(newdst)
	if err != nil {
		logs.Error("计算文件哈希值失败 err:", err)
		return nil, err
	}
	accessUrl := beego.AppConfig.String("UploadServerUrl") + "/" + fileName
	return &FileInfo{
		Name:      fileName,
		AccessUrl: accessUrl,
		Hash:      hash,
		Dst:       newdst,
	}, nil
}

func Delete(dst string) error {
	if err := DeleteFile(dst); err != nil {
		return err
	}
	return nil
}
