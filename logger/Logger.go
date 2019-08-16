package logger

import (
	"log"
	"os"
	"strings"
)

var isInTest = strings.LastIndex(os.Args[0], ".test") == len(os.Args[0])-5

// PrintError print error
func PrintError(e error) {
	// file, _ := os.OpenFile("../cache/1.log",
	// 	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// log.SetOutput(file)
	log.Printf("%v", e)
}

// PrintInfo print msg
func PrintInfo(msg string) {
	log.Println(msg)
}

// PrintInfof print msg
func PrintInfof(format string, argv ...interface{}) {
	log.Printf(format, argv...)
}

// GetLoggerPrefix get instance of LoggerPrefix
func GetLoggerPrefix(prefix string) LoggerPrefix {
	return LoggerPrefix{prefix}
}

// Logger interface
type Logger interface {
	PrintError(e error)
	PrintInfo(msg string)
	PrintInfof(format string, argv ...interface{})
}

// LoggerPrefix struct
type LoggerPrefix struct {
	prefix string
}

// PrintError print error
func (logger *LoggerPrefix) PrintError(e error) {
	log.Printf("[%s] %v", logger.prefix, e)
}

// PrintInfo print msg
func (logger *LoggerPrefix) PrintInfo(msg string) {
	if !isInTest {
		log.Printf("[%s] %s", logger.prefix, msg)
	}
}

// PrintInfof print msg
func (logger *LoggerPrefix) PrintInfof(format string, argv ...interface{}) {
	if !isInTest {
		argv = append([]interface{}{logger.prefix}, argv...)
		log.Printf("[%s] "+format, argv...)
	}
}
