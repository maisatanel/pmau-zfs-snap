package main

import (
	"errors"
	zfs "github.com/mistifyio/go-zfs/v4"
	"github.com/spf13/pflag"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {

	// parse input

	var rds string

	pflag.StringVar(&rds, "root-dataset", "", "Specifies the parent dataset for container volume datasets")

	pflag.Parse()

	if _, err := zfs.GetDataset(rds); err != nil {
		log.Fatal("Provided parent dataset is invalid. Please specify a dataset with --root-dataset")
	}

	stdoutLogger := log.New(os.Stdout, "", log.LstdFlags) // create logger

	cmd := exec.Command("/usr/bin/podman", "auto-update", "--dry-run", "--format", "{{.ContainerName}},{{.Updated}}")
	stdout, err := cmd.Output()

	if podErr, ok := errors.AsType[*exec.ExitError](err); ok {
		if podErr.Error() == "exit status 125" {
			log.Println("podman auto-update failed for some or all containers! Please run the command yourself to debug.")
		}
	} else if err != nil {
		log.Panic(err)
	}

	lines := strings.Split(string(stdout), "\n") // create a slice of containers from stdout

	switch lines[len(lines)-1:][0] { // make sure there's no straggling items
	case "", "\n":
		lines = lines[:len(lines)-1]
	}

	var c int = 0

	for _, v := range lines {
		z := strings.Split(v, ",")
		if z[1] == "pending" {

			ds, err := zfs.GetDataset(rds + "/" + z[0])
			if err != nil {
				log.Println(err)
				continue
			}

			cmd := exec.Command("/usr/bin/podman", "container", "inspect", z[0], "--format", "{{.ImageDigest}}")
			inspectStdout, err := cmd.Output()
			if err != nil {
				log.Println(err)
				continue
			}

			id := strings.TrimRight(string(inspectStdout), "\r\n")

			snap, err := ds.Snapshot(string(id), false)

			if err != nil {
				log.Println(err)
				continue
			} else {
				stdoutLogger.Println("Created snapshot " + snap.Name)
				c++
			}
		}
	}

	stdoutLogger.Println(strconv.Itoa(c) + " snapshot(s) created")

}
