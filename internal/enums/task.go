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
	All       TaskTabs = "All"
	Unique    TaskTabs = "Unique"
	Recurring TaskTabs = "Recurring"
	Github    TaskTabs = "Github"
)
