package logger

import (
	"log"
	"os"
	"strings"
)

var isInTest = strings.LastIndex(os.Args[0], ".test") == len(os.Args[0])-5 && (len(os.Args) <= 1 || os.Args[1] != "console")

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

// GetPrefixLogger get instance of PrefixLogger
func GetPrefixLogger(prefix string) PrefixLogger {
	return PrefixLogger{prefix}
}

// Logger interface
type Logger interface {
	PrintError(e error)
	PrintInfo(msg string)
	PrintInfof(format string, argv ...interface{})
}

// PrefixLogger struct
type PrefixLogger struct {
	prefix string
}

// PrintError print error
func (logger *PrefixLogger) PrintError(e error) {
	log.Printf("[%s] %v", logger.prefix, e)
}

// PrintInfo print msg
func (logger *PrefixLogger) PrintInfo(msg string) {
	if !isInTest {
		log.Printf("[%s] %s", logger.prefix, msg)
	}
}

// PrintInfof print msg
func (logger *PrefixLogger) PrintInfof(format string, argv ...interface{}) {
	if !isInTest {
		argv = append([]interface{}{logger.prefix}, argv...)
		log.Printf("[%s] "+format, argv...)
	}
}
