package models

type Task struct {
	Title  string
	Descr  string
	Status string
}

func NewTask(title string, descr string, status string) Task {
	return Task{
		Title:  title,
		Descr:  descr,
		Status: status,
	}
}
