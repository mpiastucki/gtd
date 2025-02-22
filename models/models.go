package models

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

type Status int

const (
	INBOX Status = iota
	ACTION
	LATER
	WAITING
	DONE
)

func (s Status) String() string {
	switch s {
	case INBOX:
		return "INBOX"
	case ACTION:
		return "ACTION"
	case LATER:
		return "LATER"
	case WAITING:
		return "WAITING"
	case DONE:
		return "DONE"
	default:
		return "INBOX"
	}
}

func newStatus(newStatus string) (Status, error) {
	switch newStatus {
	case "INBOX":
		return INBOX, nil
	case "ACTION":
		return ACTION, nil
	case "LATER":
		return LATER, nil
	case "WAITING":
		return WAITING, nil
	case "DONE":
		return DONE, nil
	default:
		return INBOX, fmt.Errorf("error converting %s to status", newStatus)
	}
}

type Task struct {
	Status      Status
	Description string
	Note        string
	URL         string
	Project     string
}

func NewTask() Task {
	t := Task{
		Status:  INBOX,
		Project: "noProject",
	}
	return t
}

func (t *Task) String() string {
	return fmt.Sprintf("%s;%s;%s;%s;%s\n", t.Status.String(), t.Description, t.Note, t.URL, t.Project)
}

func (t *Task) fromString(str string) error {
	pattern := `^(.+);(.+);(.*);(.*);(.*)$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(str)
	if len(matches) < 6 {
		return fmt.Errorf("error parsing text data into a new task")
	}

	status, err := newStatus(matches[1])
	if err != nil {
		return err
	}

	t.Status = status
	t.Description = matches[2]
	t.Note = matches[3]
	t.URL = matches[4]
	t.Project = matches[5]

	return nil
}

type TaskManager struct {
	Tasks        []Task
	StatusIndex  map[Status][]int
	ProjectIndex map[string][]int
}

func (tm *TaskManager) UpdateIndexes() {
	tm.StatusIndex = make(map[Status][]int)
	tm.ProjectIndex = make(map[string][]int)

	for i := 0; i < len(tm.Tasks); i++ {
		t := tm.Tasks[i]
		tm.StatusIndex[t.Status] = append(tm.StatusIndex[t.Status], i)
		tm.ProjectIndex[t.Project] = append(tm.ProjectIndex[t.Project], i)
	}
}

func (tm *TaskManager) Load(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	for sc.Scan() {
		line := sc.Text()
		t := Task{}

		err = t.fromString(line)
		if err != nil {
			log.Printf("error parsing line from file: %s, %v\n", line, err)
			continue
		}
		tm.Tasks = append(tm.Tasks, t)
	}
	if sc.Err() != nil {
		return sc.Err()
	}

	if len(tm.Tasks) > 0 {
		tm.UpdateIndexes()
	}

	return nil
}

func (tm *TaskManager) Save(filepath string) (int, error) {
	f, err := os.Create(filepath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	output := ""

	for _, task := range tm.Tasks {
		output += task.String()
	}

	bw := bufio.NewWriter(f)
	n, err := bw.Write([]byte(output))
	if err != nil {
		return n, err
	}

	err = bw.Flush()
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (tm *TaskManager) AddTask(t Task) error {
	if t.Description == "" {
		return fmt.Errorf("error adding task: cannot have empty task")
	}

	tm.Tasks = append(tm.Tasks, t)
	tm.UpdateIndexes()
	return nil
}

func (tm *TaskManager) DeleteTask(index int) error {
	if index < 0 || index >= len(tm.Tasks) {
		return fmt.Errorf("error deleting task: index not found")
	}

	tm.Tasks = append(tm.Tasks[:index], tm.Tasks[index+1:]...)
	tm.UpdateIndexes()

	return nil
}

func (tm *TaskManager) GetTask(index int) (Task, error) {
	if index < 0 || index >= len(tm.Tasks) {
		return Task{}, fmt.Errorf("error getting task: index not found")
	}

	t := tm.Tasks[index]

	return t, nil
}

func (tm *TaskManager) GetTaskIndexesByStatus(status Status) []int {
	return tm.StatusIndex[status]
}

func (tm *TaskManager) GetTaskIndexesByProject(project string) []int {
	return tm.ProjectIndex[project]
}

func (tm *TaskManager) GetTasksByIndex(indexes []int) []Task {
	t := make([]Task, 0)

	for _, taskIndex := range indexes {
		t = append(t, tm.Tasks[taskIndex])
	}
	return t
}

// TODO
func (tm *TaskManager) DeleteByStatus(status Status) error {
	tasksToDelete, ok := tm.StatusIndex[status]
	if !ok {
		return fmt.Errorf("error deleting tasks in status %s: no tasks found", status)
	}

	filteredTasks := make([]Task, 0)

	for key, task := range tm.Tasks {
		saveThisTask := true
		for _, taskIndex := range tasksToDelete {
			if key != taskIndex {
				saveThisTask = false
			}
		}
		if saveThisTask {
			filteredTasks = append(filteredTasks, task)
		}
	}

	tm.Tasks = filteredTasks
	tm.UpdateIndexes()
	return nil
}

func (tm *TaskManager) DeleteByProject(project string) error {
	tasksToDelete, ok := tm.ProjectIndex[project]
	if !ok {
		return fmt.Errorf("error finding project: %s", project)
	}

	filteredTasks := make([]Task, 0)

	for key, task := range tm.Tasks {
		saveThisTask := true
		for _, taskIndex := range tasksToDelete {
			if key != taskIndex {
				saveThisTask = false
			}
		}

		if saveThisTask {
			filteredTasks = append(filteredTasks, task)
		}

	}

	tm.Tasks = filteredTasks
	tm.UpdateIndexes()
	return nil
}
