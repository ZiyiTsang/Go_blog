package typesTool

import (
	"Go_blog/pkg/logTool"
	"strconv"
)

func Int64ToString(a int64) string {
	return strconv.FormatInt(a, 10)
}

func Uint64ToString(num uint64) string {
	return strconv.FormatUint(num, 10)
}
func StringToUint64(str string) uint64 {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		logTool.CheckError(err)
	}
	return i
}
func StringToint64(str string) int64 {
	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		logTool.CheckError(err)
	}
	return int64(i)
}
