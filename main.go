package main

import (
    "flag"
    "os"
    "os/exec"
    "fmt"
    "log"
)

var timeFlag *string
var fileFlag *string

var commands []string = []string{"create", "remove", "list"}

var root string
func init(){
    root, err = os.Getwd()
    if err!=nil {
    fmt.Println("failed to get current directory .")
}


    timeFlag := flag.String("t", "1d", "specifies the time the directory/file should exist before being automaticaly deleted.")
    fileFlag := flag.Bool("f", false, "when specified the program will create a file at the specified path.")
    
}

func main(){
    fmt.Println("")
}
