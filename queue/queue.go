package queue

type JobQueue struct {
	Name string
	channel chan Job	
}

type Job struct {
	ID string
	Run func()
}

func NewQueue(name string) *JobQueue {
	queue := make(chan Job)
	go func() {
		for job := range queue {
			job.Run()
		}
	}()

	return &JobQueue{
		Name: name,
		channel: queue,
	}	
}

func (q *JobQueue) SendToQueue(j Job) {
	q.channel <- j
}

