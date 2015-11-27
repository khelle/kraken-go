package wrapper

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"bufio"
	"io"
	"bytes"
	"strconv"
	"../../storage"
	"../../errors"
)

/**
 * ProcessWrapper class
 */
type ProcessWrapper struct {}

/**
 * ProcessWrapper constructor
 */
func New() *ProcessWrapper {
	wrapper := &ProcessWrapper{}

	return wrapper
}

/**
 * ProcessWrapper.Start()
 */
func (wrapper *ProcessWrapper) Start(args []string) errors.Error {
	env := []string{"E:\\Programy\\WebServ2.1\\httpd-users\\Kraken-standalone\\src\\kraken\\kraken-foundation\\procrun"}
	env = append(env, args...)

	// prepare php process
	cmd := exec.Command("php", env...)

	// Capture the output
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return errors.New(2, err.Error())
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.New(3, err.Error())
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.New(4, err.Error())
	}

	lsout := bufio.NewReader(stdout)
	lserr := bufio.NewReader(stderr)

	// Make process responsive for kill signal
	go func(ls io.WriteCloser) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Printf("input=%s\n", text)
			io.CopyN(ls, bytes.NewBufferString(text + "\n"), 4096)
		}
		ls.Close()
	}(stdin)

	// Fetch stdout & stderr of process
	go func(ls *bufio.Reader) {
		out := bufio.NewWriter(os.Stdout)
		for {
			line, _, _ := ls.ReadLine()

			out.WriteString(string(line) + "\n")
			out.Flush()
		}
	}(lsout)

	go func(ls *bufio.Reader) {
		err := bufio.NewWriter(os.Stderr)
		for {
			line, _, _ := ls.ReadLine()

			err.WriteString(string(line) + "\n")
			err.Flush()
		}
	}(lserr)

	// start the process
	if err := cmd.Start(); err != nil {
		return errors.New(1, err.Error())
	}

	// register process
	if cerr := wrapper.Register(args); cerr != nil {
		return cerr
	}

	// CTRL+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(wrapper *ProcessWrapper, args []string){
		<-c
		wrapper.Unregister(args)
		os.Exit(0)
	}(wrapper, args)

	// Don't let function exit before our command has finished running
	cmd.Wait()

	// unregister process
	if cerr := wrapper.Unregister(args); cerr != nil {
		return cerr
	}

	return nil
}

/**
 * ProcessWrapper.Register([]string) errors.Error
 */
func (wrapper *ProcessWrapper) Register(args []string) errors.Error {
	// update storage
	st, err := storage.NewFileStorage("kraken")
	if err != nil {
		return err
	}
	st.Open()

	data := map[string]string{}
	data["alias"]		= args[0]
	data["project"] 	= args[1]
	data["component"] 	= args[2]
	data["process"] 	= args[3]
	data["pid"] 		= strconv.Itoa(os.Getpid())
	record := storage.CreateDataRecord().FromMap(data)

	_, err = st.Add(record)

	st.Close()

	if err != nil {
		return err
	}

	return nil
}

/**
 * ProcesWrapper.Unregister([]string) errors.Error
 */
func (wrapper *ProcessWrapper) Unregister(args []string) errors.Error {
	// update storage
	st, err := storage.NewFileStorage("kraken")
	if err != nil {
		return err
	}
	st.Open()

	data := map[string]string{}
	data["alias"]		= args[0]
	data["project"] 	= args[1]
	data["component"] 	= args[2]
	data["process"] 	= args[3]
	data["pid"] 		= strconv.Itoa(os.Getpid())
	record := storage.CreateDataRecord().FromMap(data)

	_, err = st.Remove(record)
	st.Close()

	if err != nil {
		return err
	}

	return nil
}

func fmtDummy() {
	fmt.Printf("%s\n", "")
}
