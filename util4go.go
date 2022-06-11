package util4go

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : util4go.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/27 09:51
* 修改历史 : 1. [2022/4/27 09:51] 创建文件 by LongYong
*/

func Null2Default[T any](v1 *T, v2 T) T {
	if v1 == nil {
		return v2
	}

	return *v1
}

func IsNull[T any](v1 *T, fn func() T) T {
	if v1 == nil {
		return fn()
	}

	return *v1
}

var Pid = os.Getpid()
var Counter uint32
var MachineID = ReadMachineID()

func ReadMachineID() []byte {
	id := make([]byte, 3)
	hostname, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}

func GetRandomUUID() string {
	var b [12]byte
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	b[4] = MachineID[0]
	b[5] = MachineID[1]
	b[6] = MachineID[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	b[7] = byte(Pid >> 8)
	b[8] = byte(Pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&Counter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return fmt.Sprintf(`%x`, string(b[:]))
}
