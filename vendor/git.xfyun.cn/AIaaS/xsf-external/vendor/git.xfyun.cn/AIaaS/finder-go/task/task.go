package task

type Task struct {
	taskType int32
	taskData interface{}
}

var taskChan = make(chan Task, 50)

func AddTask(task Task) {

	taskChan <- task

}

