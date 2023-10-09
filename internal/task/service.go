package task

import "time"

type TaskService struct {
	repo ITaskRepo
}

func NewTaskService(repo ITaskRepo) *TaskService {
	return &TaskService{repo: repo}
}

func (ts *TaskService) CreateTask(task *Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.Status = NotStarted
	err := ts.repo.CreateTask(task)
	if err != nil {
		return err
	}
	return nil
}
