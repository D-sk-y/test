package entity

type Task struct {
	ID        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Status    string `db:"status" json:"status"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

type TaskList []Task

func NewTask(name string, status string) *Task {

	return &Task{
		Name:   name,
		Status: status,
	}
}

func NewTaskList(tasks []Task) TaskList {
	return TaskList(tasks)
}
