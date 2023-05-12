package tcron

type CronJob struct {
    Name string
    Task Task
}

type Response = int

type Tasker struct {

}

type tasker interface {
    Add(t Task, r *Response) error
    Init() error
}
