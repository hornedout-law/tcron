package tcron

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Task struct {
    Schedule Schedule `json:"schedule"`
    Command  string `json:"command"`
}


type Flags struct {
	Day   *int
	Week  *int
	Month *int
	Hour  *int
    SetAt time.Time
}


// a Schedule needs to be retrieved from a source when a the process shuts down

type Schedule struct {
    StartedAt time.Time
    Phase time.Duration
}


func (schd Schedule ) Next() int64{
    // calculate the time until the next execution 
    until_next_exec := (time.Now().UnixMilli()-schd.StartedAt.UnixMilli())%int64(schd.Phase)
    return until_next_exec
}


func (flags Flags)parseSchedule() Schedule{
    hourInMilli := 60*60*1000
    phase := (*(flags.Day)*24+*(flags.Month)*24*30+*(flags.Week)*24*7+hourInMilli)*hourInMilli
    return Schedule{flags.SetAt, time.Duration(phase)}
}


type stack interface {
    run() 
    runOnce(j Job) (string, error)
    runTask(t Task) (string, error)
    append(j Job) 
    pop(id string) error
}

type Job struct {
    Id       string `json:"id"`
    Task       Task `json:"task"`
    RunOnce bool `json:"runOnce"`
}

// Stack represents the cronjobs registered

type Stack struct {
    Stk []Job `json:"stack"`
}

type Tcron struct {
    Listener net.Listener
    Stack *Stack
    Logger Logger
}

type Logger struct {
    ErrorLog *fs.FileInfo
    OutPutLog *fs.FileInfo
}

type logger interface {
    LogError (err error) 
    LogOutput (output string)
}

// generage Job id
func generateId() string {
    randBytes := make([]byte, 20)
    _, err:= rand.Read(randBytes)
    if err!=nil {
        log.Fatal("error genrating random slice of bytes : ",err)
    }
    return string(randBytes)
}

// Read ~/.tcron.json and parse json output into a Stack struct 

func (tc Tcron) init() (Stack, error) {
	HOME := os.Getenv("HOME")
    var newStack Stack
    absolutePath := filepath.Clean(HOME+"/.tcron.json")
    _, err:= os.Stat(absolutePath)
    if err!=nil {
        initBytes, err := json.Marshal(newStack)
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
        return newStack, err
    }
    err = json.Unmarshal(cronFile,newStack)
    
	if err!=nil {
        fmt.Println("failer to parse ~/.tcron.json")
        return newStack, err
    }
    tc.Stack = &newStack
    return newStack, nil
	
}

func (s *Stack) append(task Task) *Stack{
    newJob := Job{}
    newJob.Task = task
    newJob.Id = generateId()
    s.Stk = append(s.Stk, newJob)
    return s
}

func (t *Stack) pop(id string) error {
    i:=-1
    for j, job := range t.Stk{
        if id == job.Id {
            i=j
        }
    }
    if i<0{
        return errors.New(fmt.Sprint("cronjob with id: ", id, "does not exist."))
    }else {
        t.Stk = append(t.Stk[:i], t.Stk[:i+1]...)
        return nil
    }
}

func (s *Stack)run(){
    for _, job := range s.Stk {
        if job.RunOnce == true {
            go s.runOnce(job)
        }else {
            go s.runTask(job.Task)
        }
    }
    select {}
}


func (s *Stack) runTask (t Task ) {
    sleeptime := t.Schedule.Next()

    time.Sleep(time.Duration(sleeptime))


    cmd := exec.Command("sh", "-c", t.Command)
    cmd.Run()
    
}

func (s *Stack) runOnce(j Job) {
    s.runTask(j.Task)
    s.pop(j.Id)
}

func (tc Tcron) start(){
    stack, err := tc.init()
    tc.Stack = &stack
    if err!=nil {
        fmt.Println("error initializing Stack")
        log.Fatal(err)
    }
    go tc.Stack.run()
    rpc.Register(Stack)
    l, err := net.Listen("tcp", ":6450")
    
    tc.Listener = l
    go http.Serve(l)
}

func (tc Tcron) reload(){
    tc.start()
}

func (tc Tcron) stop(){

    tc.Listener.Close()
}

