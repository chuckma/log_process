package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

/**
* 这里是一个引用，为什么要用，1 如果这个结构体很大的话会有性能上的优势
2 用引用的话可以用 l 这个参数来修改 LogProcess 它自身定义的参数
*/

// 定义接口，将写入和读取模块抽象出来
type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Writer(wc chan string)
}
type LogProcess struct {
	rc    chan []byte // 从读取模块到解析模块来传递数据
	wc    chan string // 从解析模块到写入模块之间传递数据
	read  Reader
	write Writer
}

type ReadFromFile struct {
	path string // 读取文件路径
}

type WriteToInfluxDB struct {
	influxDBDsn string // influx data source
}

func (r *ReadFromFile) Read(rc chan []byte) {
	// 读取模块
	// 打开文件
	file, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error:%s", err.Error()))
	}
	// 从文件末尾开始逐行读取文件内容
	// 将文件的字符指针移动到最后
	file.Seek(0,2)
	reader := bufio.NewReader(file)
	for ; ;  {
		// 返回整行内容，和一个错误信息 err
		line, err := reader.ReadBytes('\n')
		if err==io.EOF {
			// 文件末尾，等待产生新的日志
			time.Sleep(500*time.Microsecond)
			continue
		}else if err!= nil {
			panic(fmt.Sprintf("ReadBytes error:%s",err.Error()))

		}
		rc <- line[:len(line)-1]
	}

	
}

func (w *WriteToInfluxDB) Writer(wc chan string) {

	// 写入模块
	for v:=range wc {
		fmt.Println(v)
	}
}

func (l *LogProcess) Process() {
	// 解析模块
	for v:= range l.rc {
	//  放入 wc 这个 channel
		l.wc <- strings.ToUpper(string(v))

	}

}

// 这个用 & 也基于了性能上的考虑 ；lp 这个变量是引用类型的
func main() {
	s :="213"
	s1:= "2134"
	fmt.Println(strings.EqualFold(s,s1))

	r := &ReadFromFile{path: "./access.log",}
	w := &WriteToInfluxDB{influxDBDsn: "username&password",}
	lp := &LogProcess{
		rc:    make(chan []byte),
		wc:    make(chan string),
		read:  r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Writer(lp.wc)

	time.Sleep(30 * time.Second)

}
