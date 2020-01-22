package telewatch

import (
	"fmt"
	"log"
	"os"
)

func PrintErrorAndExit(msg string, err error) {
	log.Fatal(err)
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(128)
}
