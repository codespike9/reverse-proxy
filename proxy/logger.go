package proxy

import (
	"log"
	"os"
)

var logger *log.Logger

func InitLogger() {
	logFile, err := os.OpenFile("proxy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Could not open log file: ", err)
	}
	logger = log.New(logFile, "", log.LstdFlags)
}

func Log(message string) {
	logger.Println(message)

}
