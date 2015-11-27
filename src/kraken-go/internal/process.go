package internal

import (
	"os"
	"os/exec"
	"fmt"
	"time"
	"strconv"
	"strings"
	"../storage"
	"../errors"
)

const (
	OS_WIN		string = "win"
	OS_UNIX		string = "unix"
)

/**
 * Process class
 */
type ProcessInstance struct {
	env			*Environment
}

/**
 * Process constructor
 */
func CreateProcess(env *Environment) *ProcessInstance {
	process := &ProcessInstance{}

	process.env = env

	return process
}

/**
 * Process.PrepareWindowsCommand([]string) (string, []string)
 */
func (p *ProcessInstance) PrepareWindowsCommand(params []string) (string, []string) {
	return "cmd", append([]string{"/C","start"}, params...)
}

/**
 * Process.PrepareUnixCommand([]string) (string, []string)
 */
func (p *ProcessInstance) PrepareUnixCommand(params []string) (string, []string) {
	return "nohup", append(params, "&")
}

/**
 * Process.PrepareCommand([]string) (string, []string)
 */
func (p *ProcessInstance) PrepareCommand(params []string) (string, []string) {
	os,  _ := p.env.GetConfig().Get("env").Get("os").String()
	exe, _ := p.env.GetConfig().Get("env").Get("exe").String()
	exearr := strings.Split(exe, " ")

	params = append(exearr, params...)

	switch os {
		case OS_WIN:
			return p.PrepareWindowsCommand(params)
		case OS_UNIX:
			return p.PrepareUnixCommand(params)
		default:
			return "", []string{}
	}
}

/**
 * Process.Start(string,string,string,string) errors.Error
 */
func (p *ProcessInstance) Start(alias string, projectName string, componentName string, processName string) errors.Error {
	params := []string{}
	params = append(params, alias, projectName, componentName, processName)

	// Prepare command
	exe, params := p.PrepareCommand(params)
	if exe == "" {
		return errors.New(30, "Wrong configuration specified.")
	}

	cmd := exec.Command(exe, params...)

	// Start the process
	if err := cmd.Start(); err != nil {
		return errors.New(15, err.Error())
	}

	// Don't let function exit before our command has finished running
	cmd.Wait()

	return nil
}

/**
 * ProcList struct
 */
type ProcList map[string]string

/**
 * ProcManager class
 */
type ProcManager struct {
	env				*Environment
	storage			*storage.FileStorage
	timeOut         int
	timeInterval    int
}

/**
 * ProcManager constructor
 */
func CreateProcManager(env *Environment) *ProcManager {
	pm := &ProcManager{}

	var err errors.Error
	pm.storage, err = storage.NewFileStorage("kraken")

	if err != nil {
		return nil
	}

	pm.env			= env
	pm.timeOut		= 1000
	pm.timeInterval = 50

	return pm
}

func (pm *ProcManager) CreateProcess(alias string, projectName string, componentName string, processName string, force bool) (int, errors.Error) {
	if !force && pm.ExistsProcess(alias) {
		return 0, errors.New(2, "Process already exists.")
	}

	// clean polluted data
	pm.CleanAfterProcess(alias)

	// create process
	process := CreateProcess(pm.env)
	err  := process.Start(alias, projectName, componentName, processName)
	if err != nil {
		return 0, err
	}

	// find pid
	pid := 0
	for i := 0; i < pm.timeOut; i = i + pm.timeInterval {
		pid = pm.GetPid(alias)
		if pid != 0 {
			break
		}
		time.Sleep(time.Duration(pm.timeInterval) * time.Millisecond)
	}

	return pid, nil
}

/**
 * ProcManager.DestroyProcess(string, bool) errors.Error
 */
func (pm *ProcManager) DestroyProcess(alias string, force bool) errors.Error {
	if !force {

	}

	pid := pm.GetPid(alias)
	if pid == 0 {
		return errors.New(28, "Couldnt get process pid.")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return errors.New(29, "Couldnt find process.")
	}

	process.Kill()

	// clean polluted data
	pm.CleanAfterProcess(alias)

	return nil
}

/**
 * ProcManager.ExistsProcess(string) bool
 */
func (pm *ProcManager) ExistsProcess(alias string) bool {
	pid := pm.GetPid(alias)
	if pid == 0 {
		return false
	}

	_, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	return true
}

/**
 * ProcManager.CleanAfterProcess(string) bool
 */
func (pm *ProcManager) CleanAfterProcess(alias string) bool {
	st := pm.storage

	st.Open()

	record := storage.CreateDataRecord()
	record.Set("alias", alias)

	_, err := st.Remove(record)

	st.Close()

	if err != nil {
		return false
	} else {
		return true
	}
}

/**
 * ProcManager.GetPid(string) int
 */
func (pm *ProcManager) GetPid(alias string) int {
	st  := pm.storage
	pid := 0

	st.Open()

	record := storage.CreateDataRecord()
	record.Set("alias", alias)

	res, err := st.Get(record)

	st.Close()

	if err != nil {
		return 0
	}
	if len(res) > 0 {
		pid, _ = strconv.Atoi(res[0].Get("pid"))
	}

	return pid
}

/**
 *
 */
func fmtDummy() {
	fmt.Printf("%s\n","")
}
