package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	go startServer()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Starting DHT interface")

	const usage = `Operations:
    GET <key>
    SET <key> <value>
    OWNERS <key>
    LIST_LOCAL
    BATCH`
	fmt.Println(usage)
	r, _ := regexp.Compile("(GET) (.*)|(SET) (.*) (.*)|(LIST_LOCAL)|(OWNERS) (.*)|(BATCH)")

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error in reading from stdin")
			os.Exit(1)
		}

		if r.MatchString(text) {
			res := r.FindStringSubmatch(text)
			for i := range res {
				if i > 0 && res[i] != "" {
					runCommand(res, i)
					break
				}
			}
		} else {
			fmt.Println("Error: Could not interpret string")
			fmt.Println(usage)
		}
	}
}

func runCommand(cmds []string, i int) {
	if cmds[i] == "SET" {
		// SET(cmds[i + 1], cmds[i + 2])
		fmt.Println("RUNNING SET(" + cmds[i+1] + ", " + cmds[i+2] + ")")
	} else if cmds[i] == "GET" {
		// GET(cmds[i + 1])
		fmt.Println("RUNNING GET(" + cmds[i+1] + ")")
	} else if cmds[i] == "LIST_LOCAL" {
		// LIST_LOCAL()
		fmt.Println("RUNNING LIST_LOCAL")
	} else if cmds[i] == "OWNERS" {
		// OWNERS(cmds[i + 1])
		fmt.Println("RUNNING OWNERS(" + cmds[i+1] + ")")
	} else {
		// BATCH()
		fmt.Println("RUNNING BATCH")
	}
}
