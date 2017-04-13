package kademlia

import (
	"bufio"
	"fmt"
	"kademlia/internal/dht"
	"os"
	"regexp"
	"strings"
)

func Start() {
	go dht.StartServer()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Starting DHT interface")

	const usage = `Operations:
    GET <key>
    SET <key> <value>
    OWNERS <key>
    LIST_LOCAL
    BATCH <file_name>`
	fmt.Println(usage)
	r, _ := regexp.Compile("(GET) (.*)|(SET) (.*) (.*)|(LIST_LOCAL)|(OWNERS) (.*)|(BATCH (.*))")

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
					go runCommand(res, i)
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
		err := dht.Set(cmds[i+1], cmds[i+2])
		if err != nil {
			fmt.Println(err)
		}

	} else if cmds[i] == "GET" {
		err := dht.Get(cmds[i+1])
		if err != nil {
			fmt.Println("ERROR: ", err)
		}

	} else if cmds[i] == "OWNERS" {
		err := dht.Owners(cmds[i+1])
		if err != nil {
			fmt.Println("ERROR: ", err)
		}

	} else if cmds[i] == "LIST_LOCAL" {
		err := dht.ListLocal()
		if err != nil {
			fmt.Println("ERROR: ", err)
		}

	} else {
		err := Batch(cmds[i+1])
		if err != nil {
			fmt.Println("ERROR: ", err)
		}
	}
}

func Batch(fname string) error {
	r, _ := regexp.Compile("(GET) (.*)|(SET) (.*) (.*)|(LIST_LOCAL)|(OWNERS) (.*)|(BATCH (.*))")
	if file, err := os.Open(fname); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if r.MatchString(line) {
				res := r.FindStringSubmatch(line)
				for i := range res {
					if i > 0 && res[i] != "" {
						runCommand(res, i)
						break
					}
				}
			}
		}
	} else {
		return err
	}
	return nil
}
