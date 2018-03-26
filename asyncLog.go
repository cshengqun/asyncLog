package asyncLog

import (
	"log"
	"fmt"
	"os"
)


type ALog struct {
	writerLv int
	level    int
	logCnt   int
	logFileName  string
	fileSize int64
	fileStream *os.File
	err      *log.Logger
	warn     *log.Logger
	info     *log.Logger
	debug    *log.Logger
	ch  chan func()
}

const (
	ErrorLevel = iota
	WarnLevel
	InfoLevel
	DebugLevel
)

func NewLogger(logFileName string, level int, chanSize int, tCnt int, flag int) (*ALog) {
	logFile, err:= os.OpenFile(logFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("OpenFile fail")
		panic(err)
	}
	logger := new(ALog)
	logger.fileStream = logFile
	logger.level = level
	logger.writerLv = InfoLevel
	logger.logCnt = 20
	logger.logFileName = logFileName
	logger.fileSize = 100000000
	logger.err = log.New(logFile, "[ERROR]", flag)
	logger.warn = log.New(logFile, "[WARN]", flag)
	logger.info = log.New(logFile, "[INFO]", flag)
	logger.debug = log.New(logFile, "[DEBUG]", flag)
	logger.ch = make(chan func(), chanSize)
	for i:=0;i<tCnt;i++ {
		go logger.printLog()
	}
	return logger
}

func (aLog *ALog) retsetOutput () {
	aLog.fileStream.Close()
	logFile, err:= os.OpenFile(aLog.logFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("OpenFile fail")
		panic(err)
	}
	aLog.fileStream = logFile
	aLog.err.SetOutput(logFile)
	aLog.warn.SetOutput(logFile)
	aLog.info.SetOutput(logFile)
	aLog.debug.SetOutput(logFile)
}

func (aLog *ALog) rollFile () {
	var preFileName string
	for i:=aLog.logCnt;i>=1;i-- {
		j := i-1
		curFileName := fmt.Sprintf("%s_%d.log", aLog.logFileName, i)
		if j == 0 {
			preFileName = aLog.logFileName
		} else {
			preFileName = fmt.Sprintf("%s_%d.log", aLog.logFileName, j)
		}

		_, err := os.Stat(curFileName)
		if err == nil {
			os.Remove(curFileName)
		}

		_, err = os.Stat(preFileName)
		if err == nil {
			os.Rename(preFileName, curFileName)
		}
	}
}

func (aLog *ALog) printLog() {
	for function := range aLog.ch {
		fi, err := os.Stat(aLog.logFileName)
		if err == nil {
			if fi.Size() > aLog.fileSize {
				aLog.rollFile()
				aLog.retsetOutput()
			}
		}
		function()
	}
}

func (aLog *ALog) SetLogCnt(logCnt int) {
	aLog.logCnt = logCnt
}

func (aLog *ALog) SetFileSize(fileSize int64) {
	aLog.fileSize = fileSize
}

func (aLog *ALog) SetLevel(level int) {
	aLog.level = level
}

func (aLog *ALog) SetWriterLv(level int) {
	aLog.writerLv = level
}

func (aLog *ALog) SetPrefix(prefix string) {
	aLog.err.SetPrefix("[ERROR]" + prefix)
	aLog.warn.SetPrefix("[WARN]" + prefix)
	aLog.info.SetPrefix("[INFO]" + prefix)
	aLog.debug.SetPrefix("[DEBUG]" + prefix)
}

func (aLog *ALog) Error(format string, v ...interface{}) {
	if ErrorLevel > aLog.level {
		return
	}
	aLog.ch <- func() {
		aLog.err.Output(2, fmt.Sprintf(format, v...))
	}
}

func (aLog *ALog) Warn(format string, v ...interface{}) {
	if WarnLevel > aLog.level {
		return
	}
	aLog.ch <- func() {
		aLog.warn.Output(2, fmt.Sprintf(format, v...))
	}
}

func (aLog *ALog) Info(format string, v ...interface{}) {
	if InfoLevel > aLog.level {
		return
	}
	aLog.ch <- func() {
		aLog.info.Output(2, fmt.Sprintf(format, v...)) 
	}
}

func (aLog *ALog) Debug(format string, v ...interface{}) {
	if DebugLevel > aLog.level {
		return
	} 
	aLog.ch <- func() {
		aLog.debug.Output(2, fmt.Sprintf(format, v...))
	}
}


func (aLog *ALog) Write(p []byte) (n int, err error) {
	if aLog.writerLv > aLog.level {
		return
	}
	switch aLog.writerLv {
		case ErrorLevel:
			aLog.ch <- func() {
				aLog.err.Output(2, string(p))
			}
		case WarnLevel:
			aLog.ch <- func() {
				aLog.warn.Output(2, string(p))
			}
		case InfoLevel:
			aLog.ch <- func() {
				aLog.info.Output(2, string(p))
			}
		case DebugLevel:
			aLog.ch <- func() {
				aLog.debug.Output(2, string(p))
			}
		default:
			aLog.ch <- func() {
				aLog.debug.Output(2, string(p))
			}
	}
	n = len(p)
	err = nil
	return
}

