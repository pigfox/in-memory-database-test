package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Memory struct {
	collection  map[string]string
	transaction bool
}

type UserInput struct {
	command string
	key     string
	value   string
}

var memory Memory
var transaction Memory
var commands map[string]int

func newMemory() {
	memory.collection = make(map[string]string)
	transaction.collection = make(map[string]string)
}

/*
Sets the cmd name and including it's number of valid args
GET x
SET x 6
*/
func newCommandsMap() {
	commands = make(map[string]int)
	commands["GET"] = 2
	commands["SET"] = 3
	commands["UNSET"] = 2
	commands["BEGIN"] = 1
	commands["COMMIT"] = 1
	commands["ROLLBACK"] = 1
	commands["NUMEQUALTO"] = 2
	commands["END"] = 1
}

func init() {
	newMemory()
	newCommandsMap()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Starting...")
	for {
		inputStr, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}

		ui := parseInput(inputStr)

		if ui.command == "SET" {
			set(ui)
		}

		if ui.command == "GET" {
			if memoryKeyExist(ui.key) {
				fmt.Println(get(ui.key))
			} else {
				fmt.Println("Nil")
			}
		}

		if ui.command == "UNSET" {
			if memoryKeyExist(ui.key) {
				unSet(ui.key)
			} else {
				fmt.Println("Invalid key")
			}
		}

		if ui.command == "BEGIN" {
			begin()
		}

		if ui.command == "COMMIT" {
			commit()
		}

		if ui.command == "ROLLBACK" {
			rollback()
		}

		if ui.command == "NUMEQUALTO" {
			fmt.Println(numEqualTo(ui))
		}

		if ui.command == "END" {
			break
		}
	}
	fmt.Println("Bye")
}

func parseInput(s string) UserInput {
	s = strings.TrimSuffix(s, "\n")
	arr := strings.Split(s, " ")
	length := len(arr)
	ui := UserInput{}

	if length == 0 {
		fmt.Println("Invalid input")
		return ui
	}

	if !validCommand(arr[0]) {
		fmt.Println("Invalid command")
		return ui
	}

	if commands[arr[0]] != length {
		fmt.Println(arr[0] + " needs " + strconv.Itoa(commands[arr[0]]) + " total args. Has " + strconv.Itoa(length))
		return ui
	}

	if arr[0] == "GET" || arr[0] == "UNSET" {
		ui.command = arr[0]
		ui.key = arr[1]
	}

	if arr[0] == "NUMEQUALTO" {
		ui.command = arr[0]
		ui.value = arr[1]
	}

	if arr[0] == "SET" {
		ui.command = arr[0]
		ui.key = arr[1]
		ui.value = arr[2]
	}

	if arr[0] == "BEGIN" || arr[0] == "ROLLBACK" || arr[0] == "COMMIT" || arr[0] == "END" {
		ui.command = arr[0]
	}

	return ui
}

func validCommand(c string) bool {
	for k, _ := range commands {
		if k == c {
			return true
		}
	}
	return false
}

func memoryKeyExist(key string) bool {
	if _, ok := memory.collection[key]; ok {
		return true
	}
	return false
}

func set(ui UserInput) {
	if memory.transaction {
		transaction.collection[ui.key] = ui.value
	} else {
		memory.collection[ui.key] = ui.value
	}
}

func unSet(key string) {
	memory.collection[key] = "Nil"
}

func get(key string) string {
	return memory.collection[key]
}

func numEqualTo(ui UserInput) int {
	occurances := 0
	for _, v := range memory.collection {
		if v == ui.value {
			occurances++
		}
	}
	return occurances
}

func begin() {
	memory.transaction = true
}

func commit() {
	memory.transaction = false
	if len(transaction.collection) == 0 {
		fmt.Println("NO TRANSACTION")
		return
	}
	memory.collection = transaction.collection
	resetTransaction()
}

func rollback() {
	memory.transaction = false
	if len(transaction.collection) == 0 {
		fmt.Println("NO TRANSACTION")
		return
	}

	resetTransaction()
}

func resetTransaction() {
	transaction.collection = map[string]string{}
}
