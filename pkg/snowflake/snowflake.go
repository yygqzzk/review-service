package snowflake

import (
	"errors"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var (
	InvalidInitParamErr  = errors.New("snowflake初始化失败, 无效的startTime或machineID")
	InvalidTimeFormatErr = errors.New("snowflake初始化失败, 无效的startTime格式")
)

var node *sf.Node

// 初始化雪花算法
func Init(startTime string, machineID int64) (err error) {
	if len(startTime) == 0 || machineID <= 0 {
		return InvalidInitParamErr
	}

	var st time.Time

	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return InvalidTimeFormatErr
	}

	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID)
	return
}

// 生成雪花ID
func GenID() int64 {
	return node.Generate().Int64()
}
