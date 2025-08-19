package daily

import (
	"slices"

	"github.com/SourcewareLab/Toney/internal/messages"
)

func (m Daily) CreateTask(msg messages.TaskPopupMessage) {
	task := Task{
		TaskTitle: msg.Title,
		TaskDesc:  msg.Desc,
		Status:    msg.Status,
	}

	m.Tasks.All = append(m.Tasks.All, task)
	WriteItems(m.Tasks)
}

func (m Daily) DeleteTask(msg messages.TaskPopupMessage) {
	if !msg.IsDeleted {
		return
	}

	item := m.List.SelectedItem()

	task, ok := item.(Task)
	if !ok { // Making sure that item is of type Task
		// If it's a GitHub task, we can't delete it locally
		if _, isGithubTask := item.(GithubTask); isGithubTask {
			return // GitHub tasks can't be deleted from the local task list
		}
		return
	}

	m.Tasks.All = slices.Delete(m.Tasks.All, task.Index, task.Index+1)

	WriteItems(m.Tasks)
}

func (m Daily) StatusChangeTask(msg messages.TaskPopupMessage) {
	item := m.List.SelectedItem()

	task, ok := item.(Task)
	if !ok { // Making sure that item is of type Task
		// If it's a GitHub task, we can't change its status locally
		if _, isGithubTask := item.(GithubTask); isGithubTask {
			return // GitHub tasks status can't be changed from the local task list
		}
		return
	}

	task.Status = msg.Status

	m.Tasks.All[task.Index] = task

	WriteItems(m.Tasks)
}

func (m Daily) EditTask(msg messages.TaskPopupMessage) {
	task := Task{
		TaskTitle: msg.Title,
		TaskDesc:  msg.Desc,
		Status:    msg.Status,
	}

	item := m.List.SelectedItem()

	oldTask, ok := item.(Task)
	if !ok { // Making sure that item is of type Task
		// If it's a GitHub task, we can't edit it locally
		if _, isGithubTask := item.(GithubTask); isGithubTask {
			return // GitHub tasks can't be edited from the local task list
		}
		return
	}

	task.Status = oldTask.Status

	m.Tasks.All[oldTask.Index] = task

	WriteItems(m.Tasks)
}
