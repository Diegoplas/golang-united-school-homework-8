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

// func ValidateOperation(op string) error {
// 	if op != "add" && op != "remove" && op != "findById" && op != "list" {
// 		return fmt.Errorf("Operation %s not allowed!", op)

// 	}
// 	return nil
// }

func List(args Arguments, buff io.Writer) error {
	file, err := os.Open(args["fileName"])
	if err != nil {
		return err
	}
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("error reading file on List.. %v", err)
	}
	buff.Write(byteValue)
	return nil
}

func Add(args Arguments, buff io.Writer) error {
	if args["item"] == "" {
		return errors.New("-item flag has to be specified")
	}
	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
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
	for _, user := range items.Users {
		userID := fmt.Sprint(user["id"])
		newUserID := fmt.Sprint(newUser["id"])
		if userID == newUserID {
			buff.Write([]byte(fmt.Sprintf("Item with id %s already exists", newUserID)))
			return fmt.Errorf("Item with id %s already exists", newUserID)
		}
	}
	newUsers := append(items.Users, newUser)
	marsh, err := json.Marshal(newUsers)
	if err != nil {
		log.Printf("error marshaling. %v", err)
	}
	err = ioutil.WriteFile(args["fileName"], marsh, 0777)
	if err != nil {
		log.Printf("error while writing the file. %v", err)
	}
	return nil
}

func RemoveById(args Arguments, buff io.Writer) error {
	if args["id"] == "" {
		return errors.New("-id flag has to be specified")
	}
	file, err := os.Open(args["fileName"])
	if err != nil {
		return err
	}
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("error reading file. %v", err)
	}
	newUsers := make([]map[string]interface{}, 0)
	var items Items
	if len(byteValue) != 0 {
		err = json.Unmarshal(byteValue, &items.Users)
		if err != nil {
			log.Printf("error unmarshal1. %v", err)
			return err
		}
	}
	for index, user := range items.Users {
		inputID := fmt.Sprint(args["id"])
		savedID := fmt.Sprint(user["id"])
		if inputID == savedID {
			//newUsers = append(items.Users[:index], items.Users[index+1:]...)
			newUsers = append(items.Users[:index], items.Users[index+1:]...)
			os.Remove(args["fileName"])
			newFile, _ := os.Create(args["fileName"])
			items.Users = append(items.Users[:index], items.Users[index+1:]...)
			dataForFile, _ := json.Marshal(items.Users)
			newFile.Write(dataForFile)
		}
	}
	marsh, err := json.Marshal(newUsers)
	if err != nil {
		log.Printf("error marshaling. %v", err)
	}
	err = ioutil.WriteFile(args["fileName"], marsh, 0777)
	if err != nil {
		log.Printf("error while writing the file. %v", err)
	}
	buff.Write([]byte(fmt.Sprintf("Item with id %s not found", args["id"])))
	return nil
}

func FindByID(args Arguments, buff io.Writer) error {
	if args["id"] == "" {
		idFlagErr := "-id flag has to be specified"
		fmt.Println(idFlagErr)
		return errors.New(idFlagErr)
	}
	file, err := os.Open(args["fileName"])
	if err != nil {
		return err
	}
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("error reading file. %v", err)
	}
	var items Items
	if len(byteValue) == 0 {
		buff.Write([]byte(""))
		return nil
	}
	err = json.Unmarshal(byteValue, &items.Users)
	if err != nil {
		log.Printf("error unmarshal FindByID. %v", err)
	}
	for _, user := range items.Users {
		inputID := fmt.Sprintln(args["id"])
		savedID := fmt.Sprintln(user["id"])
		if inputID == savedID {
			marshaledUser, err := json.Marshal(user)
			if err != nil {
				log.Printf("error marshaling FindByID. %v", err)
			}
			fmt.Println(user)
			buff.Write(marshaledUser)
			return nil
		}
	}
	fmt.Printf("Item with id %s not found", args["id"])

	buff.Write([]byte(""))
	return nil
}

func Perform(args Arguments, writer io.Writer) error {
	if args["operation"] == "" {
		return errors.New("-operation flag has to be specified")
	}
	if args["fileName"] == "" {
		return errors.New("-fileName flag has to be specified")
	}
	switch args["operation"] {
	case "list":
		fmt.Println("list case")
		err := List(args, writer)
		if err != nil {
			fmt.Println(err)
			return err
		}
	case "add":
		fmt.Println("add case")
		err := Add(args, writer)
		if err != nil {
			fmt.Println(err)
			return err
		}
	case "remove":
		fmt.Println("remove case")
		err := RemoveById(args, writer)
		if err != nil {
			fmt.Println("Missing::", err)
			return err
		}
	case "findById":
		fmt.Println("finding case")
		err := FindByID(args, writer)
		if err != nil {
			fmt.Println(err)
			return err
		}
	default:
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
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
	newArgument := Arguments{"id": id, "operation": operation, "item": item, "fileName": jsonFileName}
	fmt.Println("id:", id, "operation:", operation, "item:", item, "FileName:", jsonFileName)
	err := (Perform(newArgument, &buffer))
	if err != nil {
		log.Println("Err on Perform Operation: ", err)
	}
}
