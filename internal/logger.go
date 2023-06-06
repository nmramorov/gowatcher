package metrics

import (
	"log"
	"os"
)

// Логгер, используемый для логирования событий.
var InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

// Логгер, используемый для логирования ошибок.
var ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
