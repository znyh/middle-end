package gomsg

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime/debug"
	"time"
)

func pathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func DumpStack() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	logPath := path.Join(wd, "log")
	if exists, _ := pathExist(logPath); !exists {
		os.MkdirAll(logPath, os.ModePerm)
	}

	name := path.Join(logPath, "gomsg-stack-trace"+time.Now().Format("2006-01-02_15_04_05"))
	file, err := os.Create(name)
	if err != nil {
		log.Fatalln(err)
	}

	stack := debug.Stack()
	fmt.Fprintln(file, stack)
	file.Close()
}

// Recover recover tool function
func Recover() {
	if e := recover(); e != nil {
		log.Printf("Recover => %s\n", e)
		DumpStack()
	}
}

func RecoverFromError(cb func()) {
	if e := recover(); e != nil {
		log.Printf("Recover => %s\n", e)
		DumpStack()

		if nil != cb {
			cb()
		}
	}
}
