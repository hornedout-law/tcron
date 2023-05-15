package tcron

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "io/fs"
	"io/ioutil"
	"log"
	_ "log"
	_ "net"
	_ "net/rpc"
	"os"
	"path/filepath"
	_ "strings"
	"time"

	"github.com/go-playground/locales/root"
	_ "github.com/hornedout-law/tmkdir/utils"
	"gopkg.in/ini.v1"
)

type Task struct {
    Schedule string `json:"schedule"`
    Command  string `json:"command"`
}


//func runTask(task Task) {
//	schedule, err := parseSchedule(task.Schedule)
//	if err != nil {
//		log.Printf("Failed to parse schedule for command '%s': %s", task.Command, err)
//		return
//	}
//
//	for {
//		now := time.Now().UTC()
//		nextRun := schedule.Next(now)
//
//		// Calculate the duration until the next run
//		duration := nextRun.Sub(now)
//
//		// Sleep until the next run
//		time.Sleep(duration)
//
//		// Execute the command
//		cmd := exec.Command("sh", "-c", task.Command)
//		cmd.Stdout = os.Stdout
//		cmd.Stderr = os.Stderr
//
//		err := cmd.Run()
//		if err != nil {
//			log.Printf("Failed to execute command '%s': %s", task.Command, err)
//		} else {
//			log.Printf("Executed command '%s'", task.Command)
//		}
//	}
//}
//


type Job struct {
    Id       string `json:"id"`
    Task       Task `json:"task"`
    Created_at time.Time `json:"created_at"`
    Last_exec  time.Time `json:"last_exec"`
    Next_exec  time.Time `json:"next_exec"`
}

// this is the core interface to handle the Tasker data structure
type tasker interface {
	ReadCronJobs() ([]Job, error)
	AddCronjob() (int, error)
	RemoveCronjob(id string) error
}

type Tasker struct {
	jobs []Job
}

func (t *Tasker) ReadCronjobs() ([]Job, error) {
	HOME := os.Getenv("HOME")
    var jobs []Job
    absolutePath := filepath.Clean(HOME+"/.tcron.json")
    _, err:= os.Stat(absolutePath)
    if err!=nil {
        initBytes, err := json.Marshal(jobs)
        if err!=nil {
            fmt.Println("error marsheling initial bytes.")
            log.Fatal(err)
        }

        err=ioutil.WriteFile(absolutePath, initBytes, os.ModePerm)
        
        if err!=nil {
            fmt.Println("error initializing ~/.tcron.json")
            log.Fatal(err)
        }
    }
	cronFile, err := os.ReadFile(absolutePath)
	if err!=nil {
        fmt.Println("failer to read ~/.tcron.json")
        return nil, err
    }
    err = json.Unmarshal(cronFile, jobs)
    
	if err!=nil {
        fmt.Println("failer to parse ~/.tcron.json")
        return nil, err
    }
    t.jobs = jobs
    return jobs, nil
	
}

func (t *Tasker) AddCronjob(task Task) (int, error) {
    newJob := Job{}
    newJob.Task = task
    t.jobs = append(t.jobs, newJob)
    return 0, nil
}

func (t *Tasker) RemoveCronjob(id string) error {
    i:=-1
    for j, job := range t.jobs{
        if id == job.Id {
            i=j
        }
    }
    if i<0{
        return errors.New(fmt.Sprint("cronjob with id: ", id, "does not exist."))
    }else {
        t.jobs = append(t.jobs[:i], t.jobs[:i+1]...)
        return nil
    }
}


//this is a server implement

type Scheduler struct {
	tasker Tasker
}

type scheduler interface {
    run()
    reload()
    start()
    stop()
}

func (s *Scheduler) start(){
    jobs, err := s.tasker.ReadCronjobs()
    if err!=nil{
        fmt.Print(err)
        os.Exit(1)
    }
    for _, job := range jobs {
        go s.run(job)
    }
    select {}
}

func (s *Scheduler) run(j Job){
    for {
        schedule, err:= utils.ParseSchedule(j.Task.Schedule)
    }
}

func init() {

}
