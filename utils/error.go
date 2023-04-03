package utils

import (
	"log"
	"runtime"
)

func CheckError(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Fatalf("Error in %s:%d: %v", file, line, err)
	}
}
