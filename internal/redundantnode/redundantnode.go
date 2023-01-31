package redundantnode

import (
	"fmt"
	"github.com/Moorelife/WebMind/internal/webmind1/groupnodes"
	"github.com/Moorelife/WebMind/internal/webmind1/localnode"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

type RedundantNode struct {
	Monitor    localnode.LocalNode   `json:"Monitor"`
	GroupNodes groupnodes.GroupNodes `json:"GroupNodes"`
}

func NewRedundantNode(monitor localnode.LocalNode, trio []int) *RedundantNode {
	var creation = RedundantNode{
		Monitor: monitor,
	}
	for key, value := range trio {
		if value == 0 {
			trio[key] = monitor.LocalPort + 1 + key
		}
	}
	creation.GroupNodes = *groupnodes.NewGroupNodes(monitor.LocalAddress, trio)

	return &creation
}

// CreateNodes starts up the four LocalNodes based on the RedundantNode configuration.
func (r *RedundantNode) CreateNodes() {
	startLocalNode(r.Monitor.LocalAddress, r.Monitor.LocalPort)
}

func startLocalNode(address string, port int) {
	execpath := "start"
	parameters := []string{execpath, "TESTTITLE", "B:\\webnode.cmd", address, strconv.Itoa(port)}
	sysproc := &syscall.SysProcAttr{}
	attr := os.ProcAttr{
		Dir: ".",
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Sys: sysproc,
	}
	log.Printf("%v", parameters)
	process, err := os.StartProcess(execpath, parameters, &attr)
	if err != nil {
		panic(fmt.Sprintf("StartProcess for local node failed: %v", err))
	}
	// It is not clear from docs, but Release actually detaches the process
	err = process.Release()
	if err != nil {
		panic(fmt.Sprintf("Release for local node failed: %v", err))
	}

	exePath = "_path_to_the_background_process"
	cmd := exec.Command("/usr/bin/nohup", exePath, "parameter1", "parameter2")
	if err := cmd.Start(); err != nil {
		fmt.Println("There was a problem running ", exePath, ":", err)
	} else {
		cmd.Process.Release()
		fmt.Println(exePath, " has been started.")
	}
}
