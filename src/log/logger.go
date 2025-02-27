package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log = logrus.New()

func InitLogger() {
	// Gunakan lumberjack untuk log rotation
	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log", // Lokasi penyimpanan file log
		MaxSize:    10,             // Maksimum ukuran file log (MB)
		MaxBackups: 5,              // Maksimum jumlah file log lama yang disimpan
		MaxAge:     30,             // Berapa hari sebelum log dihapus
		Compress:   true,           // Kompres file log lama (gzip)
	}

	// Set output log ke file dengan rotation
	Log.SetOutput(logFile)

	// Set format log ke JSON
	Log.SetFormatter(&logrus.JSONFormatter{})

	// Set level log dari ENV, default INFO
	LOG_LEVEL := os.Getenv("LOG_LEVEL") // Default ke INFO jika tidak ada env LOG_LEVEL
	if LOG_LEVEL == "" {
		LOG_LEVEL = "TraceLevel"
	}
	// fmt.Println("Env LOG_LEVEL:", LOG_LEVEL)

	level, err := logrus.ParseLevel(LOG_LEVEL)
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)
	fmt.Println("real LOG_LEVEL:", Log.GetLevel())
}
