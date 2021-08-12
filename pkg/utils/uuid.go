/**
 * @Author: Anpw
 * @Description:
 * @File:  uuid
 * @Version: 1.0.0
 * @Date: 2021/7/13 19:17
 */

package utils

import (
	"github.com/satori/go.uuid"
	"strings"
	"sync"
)

func GenerateUID() string {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	id := uuid.NewV4()
	uid := id.String()
	uid = strings.Replace(uid, "-", "", -1)
	return uid
}
