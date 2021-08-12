/**
 * @Author: Anpw
 * @Description:
 * @File:  jwt
 * @Version: 1.0.0
 * @Date: 2021/7/18 10:42
 */

package app

import (
	"PerInfoChain/models"
	"PerInfoChain/pkg/convert"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"time"
)

type Claims struct {
	*gorm.Model
	LoginUserID string
	jwt.StandardClaims
}

func GetJWTSecret() string {
	return beego.AppConfig.String("JWTSecret")
}

func GetJWTIssuer() string {
	return beego.AppConfig.String("JWTIssuer")
}

func GetJWTExpire() string {
	return beego.AppConfig.String("JWTExpire")
}

func GenerateToken(user models.User) (string, error) {
	nowTime := time.Now()
	expire, _ := convert.StrTo(GetJWTExpire()).Int64()
	expireTime := nowTime.Add(time.Hour * time.Duration(expire))
	claims := Claims{
		LoginUserID: user.LoginUserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    GetJWTIssuer(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(GetJWTSecret()))
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetJWTSecret()), nil
	})

	if tokenClaims != nil {
		if Claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return Claims, nil
		}
	}
	return nil, err
}
