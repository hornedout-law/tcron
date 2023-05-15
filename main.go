package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	_ "log"
	"os"
	"path/filepath"
	_ "strings"

	"github.com/hornedout-law/tmkdir/cli"
)

// get home directory

var flags cli.Flags = cli.Flags{}
var commands []string = []string{"list", "remove", "create"}

func IsValid(fp string) bool {
	// Check if file already exists
    absolutePath, err := filepath.Abs(fp)
    if err!=nil {
        log.Fatal("failed to process relative path to absolute at func IsValid")
    }
	if _, err := os.Stat(absolutePath); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(absolutePath, d, 0644); err == nil {
		os.Remove(fp) // And delete it
		return true
	}

	return false
}

func getCommandAndDirectory(args []string) (string, string) {
	if len(args) > 2 {
		log.Fatal("too much arguments")
	}
	for i, arg := range args {
		for _, command := range commands {
			if arg == command {
				if i <= len(args)-2 {
                    
                    if IsValid(args[i+1]) {
                        return arg, args[i+1]
                    }else {
                        log.Fatal("invalid path ", args[i+1])
                    }
				} else {
					return arg, ""
				}
			}
            
	}
	return "create", args[0]
}

func init() {

	flags.Day = flag.Int("d", 1, "days until next exec")
	flags.Week = flag.Int("w", 1, "weeks until next exec")
	flags.Month = flag.Int("m", 1, "months until next exec")
	flags.Hour = flag.Int("h", 1, "hours until next exec")
	flags.Date = flag.String("D", "", "")
}

func main() {
	flag.Parse()
	// get directories
	fmt.Println("days until exec: ", *(flags.Day), "\nweeks until exec: ", flags.Week, "\nmonths until exec: ", flags.Month, "\nhours until exec: ", flags.Hour, "\ndata of exec: ", flags.Date)
	args := flag.Args()
	fmt.Println("left over argumets: ", args)
	cmd, dir := getCommandAndDirectory(args)
	fmt.Println("command to be executed: ", cmd, "\ndirectory to be used with", dir)
}
