package common

import "strings"

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func ParseTaskId(taskKey string) string {
	strArr := strings.Split(taskKey, "/")
	return strArr[len(strArr)-1]
}
