package core

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

var Port string = "6450"

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

// add end endpoints to Stack

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


func (flags Flags)ParseSchedule() Schedule{
    hourInMilli := 60*60*1000
    phase := (*(flags.Day)*24+*(flags.Month)*24*30+*(flags.Week)*24*7+hourInMilli)*hourInMilli
    return Schedule{flags.SetAt, time.Duration(phase)}
}

// a Tasker is 
type Tasker interface {
    Run() 
    RunOnce(j Job) (string, error)
    RunTask(t Task) (string, error)
    Append(j Job) 
    Pop(id string) error
}

type Job struct {
    Id       string `json:"id"`
    Task       Task `json:"task"`
    RunOnce bool `json:"runOnce"`
}

// Stack represents the cronjobs registered

type Stack struct {
    Stack []Job `json:"stack"`
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
func GenerateId() string {
    randBytes := make([]byte, 20)
    _, err:= rand.Read(randBytes)
    if err!=nil {
        log.Fatal("error genrating random slice of bytes : ",err)
    }
    return string(randBytes)
}

// Read ~/.tcron.json and parse json output into a Stack struct 

func (tc Tcron) Init() (Stack, error) {
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
    err = json.Unmarshal(cronFile,&newStack)
    
	if err!=nil {
        fmt.Println("failer to parse ~/.tcron.json")
        return newStack, err
    }
    tc.Stack = &newStack
    return newStack, nil
	
}

func (s *Stack) Append(task Task) *Stack{
    newJob := Job{}
    newJob.Task = task
    newJob.Id = GenerateId()
    s.Stack = append(s.Stack, newJob)
    return s
}

func (t *Stack) Pop(id string) error {
    i:=-1
    for j, job := range t.Stack{
        if id == job.Id {
            i=j
        }
    }
    if i<0{
        return errors.New(fmt.Sprint("cronjob with id: ", id, "does not exist."))
    }else {
        t.Stack = append(t.Stack[:i], t.Stack[:i+1]...)
        return nil
    }
}

func (s *Stack)Run(){
    for _, job := range s.Stack {
        if job.RunOnce == true {
            go s.RunOnce(job)
        }else {
            go s.RunTask(job.Task)
        }
    }
    select {}
}


func (s *Stack) RunTask (t Task ) {
    sleeptime := t.Schedule.Next()

    time.Sleep(time.Duration(sleeptime))

    cmd := exec.Command("sh", "-c", t.Command)
    cmd.Run()
    
}

func (s *Stack) RunOnce(j Job) {
    s.RunTask(j.Task)
    s.Pop(j.Id)
}

func (tc Tcron) Start(){
    stack, err := tc.Init()
    tc.Stack = &stack
    if err!=nil {
        fmt.Println("error initializing Stack")
        log.Fatal(err)
    }
    go tc.Stack.Run()
    tcronRPC := new(TcronRPC)
    tcronRPC.Core = &tc
    rpc.Register(tcronRPC)
    rpc.HandleHTTP()
    l, err := net.Listen("tcp", ":"+Port)
    
    tc.Listener = l
    go http.Serve(l, nil)
}

func (tc Tcron) Reload(){
    tc.Stop()
    tc.Start()
}

func (tc Tcron) Stop(){

    tc.Listener.Close()
}

// type Arith int
// 
// func (t *Arith) Multiply(args *Args, reply *int) error {
// 	*reply = args.A * args.B
// 	return nil
// }
// 
// func (t *Arith) Divide(args *Args, quo *Quotient) error {
// 	if args.B == 0 {
// 		return errors.New("divide by zero")
// 	}
// 	quo.Quo = args.A / args.B
// 	quo.Rem = args.A % args.B
// 	return nil
// }

// make api endpoints for tmkdir

type TArgs struct {
    Flags Flags
    Path string
    IsFile bool
}

type TReply int64

type TcronC interface {
    CreateTcronEntry(ta *TArgs, reply *TReply)
}

type TcronRPC struct {
    Core *Tcron
}
// for this shit to work it need
func (tc *TcronRPC) CreateTcronEntry(ta *TArgs, reply *TReply) error {

    schedule := ta.Flags.ParseSchedule()
    // expect path to be a full path
    var task Task
    if ta.IsFile == true {
        task = Task{schedule, fmt.Sprintf("rm %s", ta.Path)} 
    } else {
        task = Task{schedule, fmt.Sprintf("rm -r %s", ta.Path)}
    }

    tc.Core.Stack.RunOnce(Job{GenerateId(), task, true})
    *reply = TReply(schedule.Next())
    return nil
}

func InitializeTcron() (*rpc.Client, error){
    client, error := rpc.DialHTTP("tcp", ":"+Port)
    if error!=nil {
        tcron := Tcron{}
        tcron.Start()
    }
    client,error = rpc.DialHTTP("tcp", ":"+Port)
    return client, error
}
