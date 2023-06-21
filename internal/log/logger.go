package log

import (
	"log"
	"os"
)

//InfoLog: Логгер, используемый для логирования событий.
var InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

//ErrorLog: Логгер, используемый для логирования ошибок.
var ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
