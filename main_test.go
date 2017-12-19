package main

import (
	"testing"
	"github.com/gorilla/rpc/json"
	"net/http"
	"bytes"
	"github.com/Pallinder/go-randomdata"
)

const (
	serverUrl   = "http://localhost:8080/rpc"
)

func TestUserCreate(t *testing.T) {
	//Create
	//Проверка создания пользователя
	userName := randomdata.SillyName();

	args := &User{
		Login: userName,
	}

	t.Log("1/1 Creating user")

	t.Log("User name:", args.Login)

	message, err := json.EncodeClientRequest("User.Create", args)
	if err != nil {
		t.Fatalf("%s", err)
	}
	req, err := http.NewRequest("POST", serverUrl, bytes.NewBuffer(message))
	if err != nil {
		t.Fatalf("%s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error in sending request to %s. %s", serverUrl, err)
	}
	defer resp.Body.Close()

	var resultCreate User
	err = json.DecodeClientResponse(resp.Body, &resultCreate)
	if err != nil {
		t.Fatalf("Couldn't decode response. %s", err)
	}
	t.Log("User created:", resultCreate)

	//
	if (resultCreate.Uuid != "" && resultCreate.Login == userName && resultCreate.RegDate != "" ) {
		t.Log("User created succesfully")
	} else {
		t.Fatal("User create error")
	}
}

func TestUserGet(t *testing.T) {
	//Create
	//Проверка создания пользователя
	userName := randomdata.SillyName();

	args := &User{
		Login: userName,
	}

	t.Log("1/2 Creating user")

	t.Log("User name:", args.Login)

	message, err := json.EncodeClientRequest("User.Create", args)
	if err != nil {
		t.Fatalf("%s", err)
	}
	req, err := http.NewRequest("POST", serverUrl, bytes.NewBuffer(message))
	if err != nil {
		t.Fatalf("%s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	respCreate, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error in sending request to %s. %s", serverUrl, err)
	}
	defer respCreate.Body.Close()

	var resultCreate User
	err = json.DecodeClientResponse(respCreate.Body, &resultCreate)
	if err != nil {
		t.Fatalf("Couldn't decode response. %s", err)
	}
	t.Log("User created:", resultCreate)

	if (resultCreate.Uuid != "" && resultCreate.Login == userName && resultCreate.RegDate != "" ) {
		t.Log("User created succesfully")
	} else {
		t.Fatal("User create error")
	}

	//Get
	//Проверка получения пользователя

	t.Log("2/2 Geting user")

	messageGet, err := json.EncodeClientRequest("User.Get", args)
	if err != nil {
		t.Fatalf("%s", err)
	}
	req, err = http.NewRequest("POST", serverUrl, bytes.NewBuffer(messageGet))
	if err != nil {
		t.Fatalf("%s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client = new(http.Client)
	respGet, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error in sending request to %s. %s", serverUrl, err)
	}
	defer respGet.Body.Close()

	var resultGet []User
	err = json.DecodeClientResponse(respGet.Body, &resultGet)
	if err != nil {
		t.Fatalf("Couldn't decode response. %s", err)
	}
	t.Log("User recieved:", resultGet[0])

	if (resultGet[0] == resultCreate) {
		t.Log("User recieved succesfully")
	} else {
		t.Fatal("User get error")
	}
}

func TestUserUpdate(t *testing.T) {
	//Создание пользователя
	userName := randomdata.SillyName();

	args := &User{
		Login: userName,
	}

	t.Log("1/3 Creating user")

	t.Log("User name:", args.Login)

	message, err := json.EncodeClientRequest("User.Create", args)
	if err != nil {
		t.Fatalf("%s", err)
	}
	req, err := http.NewRequest("POST", serverUrl, bytes.NewBuffer(message))
	if err != nil {
		t.Fatalf("%s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	respCreate, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error in sending request to %s. %s", serverUrl, err)
	}
	defer respCreate.Body.Close()

	var resultCreate User
	err = json.DecodeClientResponse(respCreate.Body, &resultCreate)
	if err != nil {
		t.Fatalf("Couldn't decode response. %s", err)
	}
	t.Log("User created:", resultCreate)

	if (resultCreate.Uuid != "" && resultCreate.Login != "" && resultCreate.RegDate != "" ) {
		t.Log("User created succesfully")
	} else {
		t.Fatal("User create error")
	}

	t.Log("2/3 Geting user")

	uuid := resultCreate.Uuid

	//Изменение пользователя
	newLogin := randomdata.SillyName();
	newDate := "1990-10-14";

	updArgs := &User{
		Uuid: uuid,
		Login: newLogin,
		RegDate: newDate,
	}

	t.Log("Updated user name:", updArgs.Login)

	messageUpdate, err := json.EncodeClientRequest("User.Update", updArgs)
	if err != nil {
		t.Fatalf("%s", err)
	}

	req, err = http.NewRequest("POST", serverUrl, bytes.NewBuffer(messageUpdate))
	if err != nil {
		t.Fatalf("%s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client = new(http.Client)
	respUpdate, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error in sending request to %s. %s", serverUrl, err)
	}
	defer respUpdate.Body.Close()

	var resultUpdate User
	err = json.DecodeClientResponse(respUpdate.Body, &resultUpdate)
	if err != nil {
		t.Fatalf("Couldn't decode response. %s", err)
	}

	t.Log("User updated:", resultUpdate)

	//Проверка пользователя
	argsGet := &User{
		Login: newLogin,
	}

	t.Log("3/3 Checking updated user")

	messageGet, err := json.EncodeClientRequest("User.Get", argsGet)
	if err != nil {
		t.Fatalf("%s", err)
	}
	req, err = http.NewRequest("POST", serverUrl, bytes.NewBuffer(messageGet))
	if err != nil {
		t.Fatalf("%s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client = new(http.Client)
	respGet, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error in sending request to %s. %s", serverUrl, err)
	}
	defer respGet.Body.Close()

	var resultGet []User
	err = json.DecodeClientResponse(respGet.Body, &resultGet)
	if err != nil {
		t.Fatalf("Couldn't decode response. %s", err)
	}
	t.Log("User recieved:", resultGet[0])

	if (resultGet[0].Login == newLogin && resultGet[0].RegDate == newDate) {
		t.Log("User update checkeced succesfully")
	} else {
		t.Fatal("User update error")
	}
}