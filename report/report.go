package report

import (
	"os"
	"strconv"
	"time"
)

var curr int64 = 0

func Report(rq []byte, res []byte, dir string) {
	curr += 1
	fname := dir + "/" + strconv.FormatInt(curr, 10) + ".md"
	file, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write([]byte("# Request\r\n"))
	file.Write([]byte("```\r\n"))
	file.Write(rq)
	file.Write([]byte("```\r\n"))
	file.Write([]byte("\r\n"))
	file.Write([]byte("# Response\r\n"))
	file.Write([]byte("```\r\n"))
	file.Write(res)
	file.Write([]byte("\r\n```\r\n"))
}

func MakeReportDir() string {
	dir := time.Now().Format("20060102_150405")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}
	return dir
}
