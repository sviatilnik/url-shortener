package pkg1

import (
	"log"
	"os"
)

func Panic() {
	panic("panic") // want "panic detected"
}

func Exit() {
	os.Exit(1) // want "os.Exit detected outside of main function"
}

func Fatal() {
	log.Fatal("fatal") // want "log.Fatal detected outside of main function"
}
