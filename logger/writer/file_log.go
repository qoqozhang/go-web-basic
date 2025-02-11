package writer

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

// RotateFileLog 按照日志轮转日志文件
type RotateFileLog struct {
	FilePrefix string     // 日志文件前缀，生成的日志文件格式：XXXX_2006-01-02.log
	FilePath   string     // 日志文件的前缀路径+logs,为空则自动设置为当前目录下的logs文件夹
	MaxAge     int        // 最大保留的日志文件天数
	file       *os.File   // 文件操作符，结构体内部使用
	mu         sync.Mutex // 结构体内部使用
	currentDay string     // 当前日期，结构体内部使用
}

func (log *RotateFileLog) Rotate() error {
	log.mu.Lock()
	defer log.mu.Unlock()
	return log.rotate()
}

func (log *RotateFileLog) Close() error {
	log.mu.Lock()
	defer log.mu.Unlock()
	return log.file.Close()
}

func (log *RotateFileLog) Write(p []byte) (n int, err error) {
	//获取锁操作
	log.mu.Lock()
	defer log.mu.Unlock()

	// 如果是首次使用，则需要打开日志文件操作符
	if log.file == nil {
		if err := log.openExistingOrNew(); err != nil {
			return 0, err
		}
	}

	// 判断当前打开的日志文件是否是今天的日志文件，不是的话则轮转日志文件
	currentDay := time.Now().Format("2006-01-02")
	if log.currentDay != currentDay {
		err = log.rotate()
	}
	if err != nil {
		return 0, err
	}

	// 写入日志到文件
	n, err = log.file.Write(p)
	return n, err
}
func (log *RotateFileLog) Read(p []byte) (n int, err error) {
	log.mu.Lock()
	defer log.mu.Unlock()
	if log.file == nil {
		if err := log.openExistingOrNew(); err != nil {
			return 0, err
		}
	}
	return log.file.Read(p)
}
func (log *RotateFileLog) ReadLines(offset, size int) (p []string) {
	log.mu.Lock()
	defer log.mu.Unlock()
	if log.file == nil {
		if err := log.openExistingOrNew(); err != nil {
			return nil
		}
	}
	scanner := bufio.NewScanner(log.file)
	currentLine := 1
	lastLine := offset + size
	for scanner.Scan() {
		if currentLine >= offset && currentLine < lastLine {
			p = append(p, scanner.Text())
		}
		currentLine++
	}
	return p
}

// deleteMaxAgeBeforeLog 删除保留保留日志之前的日志
func (log *RotateFileLog) deleteMaxAgeBeforeLog() error {
	filePath, err := log.getAbsoluteFilePath()
	if err != nil {
		return err
	}
	// 遍历日志目录下面的所有文件，并判断文件名称格式是不是 xxxxx_2016-01-02的格式，如果是的话则删除
	err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		IsLogFile, _ := regexp.MatchString(log.FilePrefix+"_"+`\d{4}-\d{2}-\d{2}`, info.Name())
		if IsLogFile {
			modifiedTime := info.ModTime()
			diff := time.Now().Sub(modifiedTime)
			if log.MaxAge != 0 && diff > time.Duration(log.MaxAge*24)*time.Hour {
				err := os.Remove(path)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// rotate 关闭之前打开的文件操作符，并创建一个新的文件操作符
func (log *RotateFileLog) rotate() error {
	if err := log.file.Close(); err != nil {
		return err
	}
	logfile, err := log.GetCurrentLogName()
	if err != nil {
		return err
	}
	return log.openNew(logfile)
}

// openNew 打开一个新的文件
func (log *RotateFileLog) openNew(filename string) error {
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	log.file = f
	log.currentDay = time.Now().Format("2006-01-02")
	return nil
}

// openExistingOrNew 尝试打开文件，不存在则创建一个新的文件
func (log *RotateFileLog) openExistingOrNew() error {
	logfile, err := log.GetCurrentLogName()
	if err != nil {
		return err
	}
	_, err = os.Stat(logfile)
	if os.IsNotExist(err) {
		return log.openNew(logfile)
	}
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	log.file = file
	log.currentDay = time.Now().Format("2006-01-02")
	return nil
}

// GetAbsoluteFilePath  获取文件的绝对路径
// 1. 判断日志文件夹是 ""，则返回执行程序的当前路径下的 logs/路径,不存在则创建
// 2. 如果日志文件夹路径不为空，则判断是否存在，不存在则创建，创建失败则返回error
func (log *RotateFileLog) getAbsoluteFilePath() (fileFullName string, err error) {
	var filePath string
	filePath, err = os.Getwd()
	if err != nil {
		return "", err
	}
	if log.FilePath == "" {
		fileFullName = path.Join(filePath, "logs")
	} else {
		fileFullName = path.Join(log.FilePath, "logs")
	}

	// 判断是否存在，不存在则尝试创建
	if _, err := os.Stat(fileFullName); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(fileFullName, 0777)
			if err != nil {
				return "", fmt.Errorf("create logs dir fail: %v", err)
			}
		}
		return "", err
	}
	fmt.Printf("Error: %v\n", err)
	return fileFullName, nil
}

// GetCurrentLogName 获取当前的日志文件绝对路径
func (log *RotateFileLog) GetCurrentLogName() (logfile string, err error) {
	currentDay := time.Now().Format("2006-01-02")
	filepath, err := log.getAbsoluteFilePath()
	if err != nil {
		return "", err
	}
	logfile = path.Join(filepath, log.FilePrefix+"_"+currentDay+".log")
	return logfile, nil
}
