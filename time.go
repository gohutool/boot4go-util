package util4go

import (
	"strconv"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : time.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/9 19:51
* 修改历史 : 1. [2022/5/9 19:51] 创建文件 by LongYong
*/

type Dimension struct {
	M1  string
	M15 string
	M30 string
	M60 string

	LM1  string
	LM15 string
	LM30 string
	LM60 string
}

func Time2Dimension(t time.Time) (*Dimension, error) {
	d := &Dimension{}
	/*
		if s, err := strconv.Atoi(t.Format("200601021504")); err != nil {
				return nil, err
			} else {
				d.M1 = t.Format("2006/01/02 15:04")
				d.M15 = t.Format("2006/01/02 15:") + strconv.Itoa((s%100+16)/15) + "/4"
				d.M30 = t.Format("2006/01/02 15:") + strconv.Itoa((s%100+31)/30) + "/2"
				d.M60 = t.Format("2006/01/02 15")

				d.LM1 = t.Add(-1 * time.Minute).Format("2006/01/02 15:04")
				d.LM15 = t.Add(-15*time.Minute).Format("2006/01/02 15:04") + strconv.Itoa((s%100+16)/15) + "/4"
			}
	*/

	d.M1 = t.Format("2006/01/02 15:04")
	d.M15 = t.Format("2006/01/02 15:") + strconv.Itoa((t.Minute()+16)/15) + "/4"
	d.M30 = t.Format("2006/01/02 15:") + strconv.Itoa((t.Minute()+31)/30) + "/2"
	d.M60 = t.Format("2006/01/02 15")

	d.LM1 = t.Add(-1 * time.Minute).Format("2006/01/02 15:04")
	t15 := t.Add(-15 * time.Minute)
	d.LM15 = t15.Format("2006/01/02 15:04") + strconv.Itoa((t15.Minute()+16)/15) + "/4"
	t30 := t.Add(-30 * time.Minute)
	d.LM30 = t30.Format("2006/01/02 15:04") + strconv.Itoa((t30.Minute()+16)/30) + "/2"
	t60 := t.Add(-60 * time.Minute)
	d.LM60 = t60.Format("2006/01/02 15")

	return d, nil
}
