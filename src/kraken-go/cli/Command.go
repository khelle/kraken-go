package cli

import (
	"fmt"
	"strings"
	"../errors"
	"../internal"
	"../util"
)

const (
	COMMAND_CREATE		string = "CREATE"
	COMMAND_DESTROY		string = "DESTROY"
	COMMAND_START		string = "START"
	COMMAND_STOP		string = "STOP"
)

/**
 * Command struct
 */
type Command struct {
	Env		*internal.Environment
	Action	string
	Args	map[string]string
}

/**
 * Command constructor
 */
func CreateCommand(env *internal.Environment, args []string) *Command {
	if len(args) < 1 {
		return nil
	}

	command := &Command{}

	command.Env	   = env
	command.Action = args[0]
	command.Args   = make(map[string]string)

	for i := 1; i < len(args); i++ {
		tmp := strings.Split(args[i], "=")
		command.Args[strings.Replace(tmp[0], "-", "", -1)] = tmp[1]
	}

	return command
}

/**
 * Executes given operation
 */
func (c *Command) Execute() errors.Error {
	switch c.Action {
		case COMMAND_CREATE:
			return c.Create()
		case COMMAND_DESTROY:
			return c.Destroy()
		case COMMAND_START:
			return c.Start()
		case COMMAND_STOP:
			return c.Stop()
		default:
			return errors.New(28, "Undefined command specified.")
	}
}

/**
 * Command.Create() errors.Error
 */
func (c *Command) Create() errors.Error {
	// check if arguments are valid
	args := c.Args
	if !util.KeyExists(args, "alias") || !util.KeyExists(args, "project") || !util.KeyExists(args, "component") || !util.KeyExists(args, "process") {
		return errors.New(27, "Not enough input argument.")
	}

	// create process manager
	pm := internal.CreateProcManager(c.Env)
	if pm == nil {
		return errors.New(26, "Process manager couldnt been initalized.")
	}

	// create process
	_, err := pm.CreateProcess(args["alias"], args["project"], args["component"], args["process"], false)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Command.Destroy() errors.Error
 */
func (c *Command) Destroy() errors.Error {
	// check if arguments are valid
	args := c.Args
	if !util.KeyExists(args, "alias") {
		return errors.New(27, "Not enough input argument.")
	}

	// create process manager
	pm := internal.CreateProcManager(c.Env)
	if pm == nil {
		return errors.New(26, "Process manager couldnt been initalized.")
	}

	// destroy process
	err := pm.DestroyProcess(args["alias"], false)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Command.Start() errors.Error
 */
func (c *Command) Start() errors.Error {
	return nil
}

/**
 * Command.Stop() errors.Error
 */
func (c *Command) Stop() errors.Error {
	return nil
}

func fmtDummy() {
	fmt.Printf("")
}
