package logfile

import (
	"fmt"
	"os"
	"time"
)

//初始化日志文件
func InitLogfile(filename string) *os.File {

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0660)

	if err != nil {
		fmt.Println("--------------"+filename, err)
		return nil
	} else {
		return f
	}

}

//获取当前时间，格式化为"2006-01-02 03:04:05"
func GetNowtime() string {
	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	t := tm.Format("2006-01-02 03:04:05")
	return t
}

//写入日志文件
//func Write2LogFile(content string, fout *os.File) {
//	fout.WriteString(content)
//	fmt.Println("--------写入日志 end---------") //
//}
