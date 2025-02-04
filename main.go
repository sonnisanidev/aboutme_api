package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// GitHub raw content URL for the JSON file
	url := "https://raw.githubusercontent.com/sonnisanidev/aboutme_api/main/data.json"

	// Send GET request to GitHub
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching JSON:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Parse JSON into a map
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Print the parsed JSON data
	fmt.Println(data)
}
