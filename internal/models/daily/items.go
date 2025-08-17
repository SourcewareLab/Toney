package daily

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SourcewareLab/Toney/internal/config"
	"github.com/SourcewareLab/Toney/internal/enums"
	"github.com/charmbracelet/bubbles/list"
)

type Tasks struct {
	Recurring []Task       `json:"recurring"`
	Unique    []Task       `json:"unique"`
	Github    []GithubTask `json:"github"`
}

type GithubTask struct {
	TaskTitle string           `json:"title"`
	TaskDesc  string           `json:"desc"`
	Status    enums.TaskStatus `json:"status"`
	Ref       string           `json:"ref"`
	Repo      string           `json:"repo"`
	Owner     string           `json:"owner"`
	// We will not store these data in the file
	Link     string   `json:"-"`
	Labels   []string `json:"-"`
	Assignee []string `json:"-"`
}

func (m GithubTask) Title() string       { return m.TaskTitle }
func (m GithubTask) Description() string { return m.TaskDesc }
func (m GithubTask) FilterValue() string { return m.TaskTitle }

type Task struct {
	TaskTitle string           `json:"title"`
	TaskDesc  string           `json:"desc"`
	Status    enums.TaskStatus `json:"status"`
	// We will not store these data in the file
	Index    int            `json:"-"` // Point to index in the respective type array
	TaskType enums.TaskType `json:"-"`
}

func (m Task) Title() string       { return m.TaskTitle }
func (m Task) Description() string { return m.TaskDesc }
func (m Task) FilterValue() string { return m.TaskTitle }

func (m Tasks) ItemsAsList() []list.Item {
	lst1 := TaskToItems(m.Recurring)
	lst2 := TaskToItems(m.Unique)

	return append(lst1, lst2...)
}

func TaskToItems(tasks []Task) []list.Item {
	list := make([]list.Item, 0)
	for i, v := range tasks {
		v.Index = i
		v.TaskType = enums.RecurringTask
		list = append(list, v)
	}
	return list
}

func GetItems() Tasks {
	path := GetPath()

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		f, err2 := os.Create(path)
		if err2 != nil {
			fmt.Println(err2.Error())
		}

		content, err2 := os.ReadFile(GetYesterdayPath())
		if err2 != nil {
			fmt.Println(err2.Error())
		}

		tasks := Tasks{}
		json.Unmarshal(content, &tasks)

		tasks.Unique = make([]Task, 0)

		data, err2 := json.Marshal(tasks)
		if err2 != nil {
			fmt.Println(err2.Error())
		}

		f.Write(data)
	} else if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	content, _ := os.ReadFile(path)

	tasks := Tasks{}

	err = json.Unmarshal(content, &tasks)
	if err != nil {
		fmt.Println(err.Error())
	}

	return tasks
}

func WriteItems(tasks Tasks) {
	path := GetPath()

	data, _ := json.Marshal(tasks)

	os.WriteFile(path, data, 0o644)
}

func GetPath() string {
	home, _ := os.UserHomeDir()
	date := time.Now().Format("2006-01-02")

	return filepath.Join(home, config.AppConfig.General.NotesDir, ".daily", date)
}

func GetYesterdayPath() string {
	home, _ := os.UserHomeDir()
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	return filepath.Join(home, config.AppConfig.General.NotesDir, ".daily", yesterday)
}
