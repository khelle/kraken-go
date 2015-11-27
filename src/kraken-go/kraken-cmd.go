package main

import (
	"os"
//	"fmt"
	"./cli"
//	"./errors"
	"./internal"
//	"./storage"
//	"./lock"
)

func main() {
//	var err errors.Error
//
//	fmt.Printf("%s\n", "starting...")
//
	// console is not designed to work concurrently, so acquiring lock is required
//	sync := lock.NewLock()
//	_, err = sync.Lock()
//	errors.Log(err)
//	defer sync.Unlock()

//	fmt.Printf("%s\n", "processing...")

//	st, err := storage.NewFileStorage("../../data/test.data")
//	defer st.Close()
//	errors.Log(err)
//
//	recordp := map[string]string{}
//	recordp["project"] = "myproj"
//	recordp["component"] = "mycomp"
//	recordp["process"] = "router"
//	recordp["pid"] = "7752"
//
////	record := storage.StorageRecord(recordp)
////	_, err = st.MultiAdd([]*storage.StorageRecord{&record,&record})
////	_, err = st.Add(&record)
////	errors.Log(err)
//
////	records, err := st.GetAll()
////	fmt.Printf("%#v\n", records)
//	errors.Log(err)
//
//	recordf := map[string]string{}
////	recordf["pid"] = "7752"
//	recordf["process"] = "router1"
//	record := storage.StorageRecord(recordf)
//	find, err := st.Exclude(&record)
//
//	for _, val := range find {
//		fmt.Printf("%#v\n", val)
//	}
//
//	st.Remove(&record)


//	os.Exit(0)
//
//	pm := internal.NewProcManager()
//	if pm == nil {
//		errors.Log(errors.New(16, "Process manager couldnt been initalized."))
//	}
//
//	var pid int
//	pid, err = pm.CreateProcess("myalias", "myproject", "mycomponent", "myprocess", false)
//	errors.Log(err)
//
//	fmt.Printf("pid=%d\n", pid)

	// prepare environment
	env := internal.CreateEnvironment()
	if env == nil {
		os.Exit(1)
	}

	// parse commandLine arguments into Command object
	command := cli.CreateCommand(env, os.Args[1:])

	// execute command && check results
	if command == nil || command.Execute() != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
