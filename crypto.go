package util4go

import (
	"crypto/md5"
	"encoding/hex"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : crypto.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/1 10:20
* 修改历史 : 1. [2022/5/1 10:20] 创建文件 by LongYong
*/

func MD5(data string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(data))
	dataMd5 := md5Ctx.Sum(nil)
	return hex.EncodeToString(dataMd5)
}

func SaltMd5(data, salt string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(salt))
	md5Ctx.Write([]byte(data))
	md5Ctx.Write([]byte(salt))
	dataMd5 := md5Ctx.Sum(nil)
	return hex.EncodeToString(dataMd5)
}
