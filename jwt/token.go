package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : token.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/1 10:21
* 修改历史 : 1. [2022/5/1 10:21] 创建文件 by LongYong
*/

func GenToken(encode TokenSubjectEncode, issuer string, secretKey string, expire time.Duration) string {

	var subject string

	if encode != nil {
		subject = encode()
	}

	signKey := []byte(secretKey)
	claims := &jwt.StandardClaims{
		Subject:   subject,
		ExpiresAt: time.Now().Add(expire).Unix(),
		Issuer:    issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signKey)
	if err != nil {
		return ""
	}
	return ss
}

type TokenSubjectEncode func() string
type TokenSubjectDecode func(subject string) (any, error)

type TokenInvalidError struct {
}

func (t TokenInvalidError) Error() string {
	return "Token is invalid"
}

func CheckToken(secretKey string, tokenStr string, decode TokenSubjectDecode) (any, error) {
	if token, err := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	}); err == nil {
		claims := token.Claims.(*jwt.StandardClaims)

		if token.Valid {

			if decode == nil {
				return nil, nil
			} else {
				subject := claims.Subject
				return decode(subject)
			}
		} else {
			return nil, TokenInvalidError{}
		}
	} else {
		return nil, err
	}

}
