package cmd

import "github.com/mpiastucki/gtd/models"

type command interface {
	execute() error
}

type addTaskCommand struct {
	tm *models.TaskManager
	t  models.Task
}

func (c *addTaskCommand) execute() error {
	c.tm.AddTask(c.t)
	return nil
}

type deleteTaskCommand struct {
	tm    *models.TaskManager
	index int
}

func (c *deleteTaskCommand) execute() error {
	err := c.tm.DeleteTask(c.index)
	return err
}

type getTaskCommand struct {
	tm    *models.TaskManager
	index int
	task  models.Task
}

func (c *getTaskCommand) execute() error {
	t, err := c.tm.GetTask(c.index)
	if err != nil {
		c.task = t
	}
	return err
}
