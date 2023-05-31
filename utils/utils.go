package utils

import (
	"fmt"
	"log"

)

func HandleError(err error, message string) {
	if err != nil {
		if message == "" {
            fmt.Println("something went wrong.")
		}else {
            fmt.Println(err.Error())
        }
        log.Fatal(err)
	}
}

func parseSchedule(){

}
