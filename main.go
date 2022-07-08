package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Arguments map[string]string

// "id":        "",
// "operation": "findById",
// "item":      "",
// "fileName":  fileName,

type Items struct {
	Users []map[string]interface{}
}

func ValidateArguments(args Arguments) error {
	if args["operation"] == "" {
		return errors.New("-operation flag has to be specified")
	}
	if args["item"] == "" {
		return errors.New("`-item flag has to be specified")
	}
	if args["id"] == "" {
		return errors.New("`-id flag has to be specified")
	}
	if args["fileName"] == "" {
		return errors.New("`-fileName flag has to be specified")
	}
	return nil
}

func ValidateOperation(op string) error {
	if op != "add" && op != "remove" && op != "findById" && op != "list" {
		return errors.New("Operation abcd not allowed!")
	}
	return nil
}

func List(file *os.File) {
	fmt.Println("Listing users...")
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("error reading file on List.. %v", err)
	}

	fmt.Println(string(byteValue))
}

func Add(file *os.File, args Arguments) {
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("error reading file. %v", err)
	}
	var newUser map[string]interface{}
	var items Items
	if len(byteValue) != 0 {
		err = json.Unmarshal(byteValue, &items.Users)
		if err != nil {
			log.Printf("error unmarshal1. %v", err)
		}
	}

	err = json.Unmarshal([]byte(args["item"]), &newUser)
	if err != nil {
		log.Printf("error unmarshal2. %v", err)
	}
	newUsers := append(items.Users, newUser)
	fmt.Println(newUsers)
	marsh, err := json.Marshal(newUsers)
	if err != nil {
		log.Printf("error marshaling. %v", err)
	}
	err = ioutil.WriteFile(args["fileName"], marsh, 0777)
	if err != nil {
		log.Printf("error while writing the file. %v", err)
	}
}

func FindByID(file *os.File, args Arguments, writer io.Writer) {
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("error reading file. %v", err)
	}
	var items Items
	if len(byteValue) == 0 {
		writer.Write([]byte(""))
		return
	}
	err = json.Unmarshal(byteValue, &items.Users)
	if err != nil {
		log.Printf("error unmarshal FindByID. %v", err)
	}
	for _, user := range items.Users {
		inputID := fmt.Sprintln(args["id"])
		if inputID == user["id"] {
			marshaledUser, err := json.Marshal(user)
			if err != nil {
				log.Printf("error marshaling FindByID. %v", err)
			}
			writer.Write(marshaledUser)
			return
		}
	}
	writer.Write([]byte(""))
}

//// Meter las validaciones a cada funcion

func Perform(args Arguments, writer io.Writer) error {
	fmt.Println("ON PERFORM!")
	err := ValidateArguments(args)
	if err != nil {
		return err
	}
	err = ValidateOperation(args["operation"])
	if err != nil {
		return err
	}
	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return errors.New("error opening or creating file")
	}
	defer file.Close()

	switch args["operation"] {
	case "list":
		fmt.Println("list case")
		List(file)
	case "add":
		fmt.Println("add case")
		Add(file, args)
	default:
		fmt.Println("something wrong")
	}

	return nil
}

func main() {
	var buffer bytes.Buffer
	var id, operation, item, jsonFileName string

	flag.StringVar(&id, "id", "", "id of a user")
	flag.StringVar(&operation, "operation", "", "operation to be performed")
	flag.StringVar(&item, "item", "", "list of users")
	flag.StringVar(&jsonFileName, "fileName", "", "filename of the json to be created")
	flag.Parse()
	fmt.Println(item)
	newArgument := Arguments{"id": "", "operation": operation, "item": item, "fileName": jsonFileName}
	fmt.Println("operation", operation, "item", item, "FileName", jsonFileName)
	err := (Perform(newArgument, &buffer))
	if err != nil {
		log.Println("Err on Perform Operation: ", err)
	}
}
