package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	Name struct {
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"name"`
	Age         int `json:"age"`
	Preferences struct {
		Description string `json:"description"`
		Hobbies     string `json:"hobbies"`
	} `json:"preferences"`
}

type Data struct {
	User User `json:"user"`
}

func main() {
	url := "https://raw.githubusercontent.com/sonnisanidev/aboutme_api/main/data.json"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching JSON:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Printf("Name: %s %s\n", data.User.Name.First, data.User.Name.Last)
	fmt.Printf("Age: %d\n", data.User.Age)
	fmt.Printf("Description: %s\n", data.User.Preferences.Description)
	fmt.Printf("Hobbies: %s\n", data.User.Preferences.Hobbies)
}
