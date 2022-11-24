package main

import (
	"bufio"
	"fmt"
	ps "github.com/mitchellh/go-ps"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var curDir string

func changeDir(path string) {
	err := os.Chdir(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	curDir = path
}

func printCurrentDir() {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)
}

func killProcess(strPid string) {
	pid, err := strconv.Atoi(strPid)
	if err != nil {
		fmt.Println("Bad pid")
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = process.Kill()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func echo(str string) {
	fmt.Println(str)
}

func printProcessList() {
	processList, err := ps.Processes()
	if err != nil {
		log.Println(err)
		return
	}

	// map ages
	for x := range processList {
		var process ps.Process
		process = processList[x]
		log.Printf("%d\t%s\n", process.Pid(), process.Executable())
		// do os.* stuff on the pid
	}
}

func execute(arg string) {
	args := strings.Split(arg, " ")

	if len(args) < 1 {
		fmt.Println("Error: Bad arguments")
	} else {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Run()
	}
}

func executeCommand(command string) bool {
	split := strings.Split(command, " ")
	switch split[0] {
	case "cd":
		changeDir(strings.Replace(command, "cd ", "", 1))
	case "pwd":
		printCurrentDir()
	case "echo":
		echo(strings.Replace(command, "echo ", "", 1))
	case "ps":
		printProcessList()
	case "exec":
		execute(strings.Replace(command, "exec ", "", 1))
	case `\quit`:
		return false
	default:
		fmt.Printf("Command '%s' not found\n", split[0])
	}

	return true
}

func main() {
	curDir, _ = os.Getwd()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\033[34m%s>\033[0m ", curDir)
	for isWorking := true; scanner.Scan() && isWorking; {
		executeCommand(scanner.Text())
		fmt.Printf("\033[34m%s>\033[0m ", curDir)
	}

}
