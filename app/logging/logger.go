package logging

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
	TimingLogger  *log.Logger
	initialized   bool
)

func Init(baseDir string) {
	if initialized {
		return
	}
	infofile, err := os.OpenFile(baseDir+"/go_sorting_info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	warnfile, err := os.OpenFile(baseDir+"/go_sorting_warn.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	errfile, err := os.OpenFile(baseDir+"/go_sorting_error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	debugfile, err := os.OpenFile(baseDir+"/go_sorting_debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	timingfile, err := os.OpenFile(baseDir+"/go_sorting_timing.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(infofile, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(warnfile, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errfile, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(debugfile, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	TimingLogger = log.New(timingfile, "[TIMING] ", log.Ldate|log.Ltime|log.Lshortfile)
}
