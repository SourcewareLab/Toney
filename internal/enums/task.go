package enums

type TaskStatus int

const (
	Pending = iota
	Started
	Abandoned
	Complete
)

type TaskTabs string

const (
	Tasks  TaskTabs = "Tasks"
	Github TaskTabs = "Github"
)
