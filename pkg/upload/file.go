/**
 * @Author: Anpw
 * @Description:
 * @File:  file
 * @Version: 1.0.0
 * @Date: 2021/7/27 0:36
 */

package upload

import (
	"PerInfoChain/pkg/utils"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"github.com/astaxie/beego"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

type FileType int

const TypeImage FileType = iota + 1

//GetFileName
/**
 * @Author: Anpw
 * @Description: 将原始文件名加密处理
 * @param name
 * @return string
 */
func GetFileName(name string) string {
	ext := GetFileExt(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = utils.EncodeMD5(fileName)
	return fileName + ext
}

//GetFileExt
/**
 * @Author: Anpw
 * @Description: 获取文件后缀
 * @param name
 * @return string
 */
func GetFileExt(name string) string {
	return path.Ext(name)
}

//GetSavePath
/**
 * @Author: Anpw
 * @Description: 获取文件保存地址
 * @return string
 */
func GetSavePath() string {
	return beego.AppConfig.String("UploadSavePath")
}

func CheckSavePath(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsNotExist(err)
}

func CheckContainExt(t FileType, name string) bool {
	UploadFileAllowExt := beego.AppConfig.String("UploadFileAllowExt")
	ext := GetFileExt(name)
	switch t {
	case TypeImage:
		if strings.ToUpper(UploadFileAllowExt) == strings.ToUpper(ext) {
			return true
		}
	}
	return false
}

func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

func CreateSavePath(dst string, perm os.FileMode) error {
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}
	return nil
}

func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func UpdateFile(file *multipart.FileHeader, dst string, newdst string) error {
	/**
		 * TODO-Anpw: 2021/7/28 0:48 文件更新是用的删除原文件创建新文件，原本是直接修改但是会报错文件正在使用无法修改
	 	 * Description:
	*/
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	err = os.Remove(dst)
	if err != nil {
		return err
	}
	out, err := os.Create(newdst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func DeleteFile(dst string) error {
	err := os.Remove(dst)
	if err != nil {
		return err
	}
	return err
}

// FileEncrypt 文件加密
func FileEncrypt(fileName string) {
	plaintext, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)
	key, err := ioutil.ReadFile("key")
	if err != nil {
		log.Fatal(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err)
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	// Save back to file
	err = ioutil.WriteFile("ciphertext.bin", ciphertext, 0777)
	if err != nil {
		log.Panic(err)
	}
}

func FileNameHash(dst string) (string, error) {
	file, err := os.Open(dst)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	if err != nil {
		return "", err
	}
	md5h := md5.New()
	if _, err := io.Copy(md5h, file); err != nil {
		log.Fatal(err)
	}
	MD5Str := hex.EncodeToString(md5h.Sum(nil))
	return MD5Str, nil
}
