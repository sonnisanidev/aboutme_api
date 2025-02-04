package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     struct {
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"name"`
	Age          int         `json:"age"`
	IsActive     bool        `json:"isActive"`
	Roles        []string    `json:"roles"`
	Address      Address     `json:"address"`
	PhoneNumbers []Phone     `json:"phoneNumbers"`
	Preferences  Preferences `json:"preferences"`
	LastLogin    string      `json:"lastLogin"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zipCode"`
	Country string `json:"country"`
}

type Phone struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

type Preferences struct {
	Theme         string `json:"theme"`
	Notifications struct {
		Email bool `json:"email"`
		Push  bool `json:"push"`
	} `json:"notifications"`
	Description string `json:"description"`
	Hobbies     string `json:"hobbies"`
}

type Data struct {
	User User `json:"user"`
}

func main() {

	http.HandleFunc("/user", handleUser)
	http.HandleFunc("/update-user", updateUser)
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	// ... (keep this function as is)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := os.Getenv("GITHUB_TOKEN")
	print(token)
	url := "https://raw.githubusercontent.com/sonnisanidev/aboutme_api/refs/heads/main/data.json"

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read data", http.StatusInternalServerError)
		return
	}

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Failed to parse data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func updateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed. Use PUT or PATCH.", http.StatusMethodNotAllowed)
		return
	}

	// Read and log the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received request body: %s\n", string(body))

	// Parse the updated data
	var updatedData Data
	err = json.Unmarshal(body, &updatedData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch current data
	currentData, err := fetchCurrentData()
	if err != nil {
		http.Error(w, "Failed to fetch current data", http.StatusInternalServerError)
		return
	}

	// Update name, age, description, and hobbies if provided
	if updatedData.User.Name.First != "" {
		currentData.User.Name.First = updatedData.User.Name.First
	}
	if updatedData.User.Name.Last != "" {
		currentData.User.Name.Last = updatedData.User.Name.Last
	}
	if updatedData.User.Age != 0 {
		currentData.User.Age = updatedData.User.Age
	}
	if updatedData.User.Preferences.Description != "" {
		currentData.User.Preferences.Description = updatedData.User.Preferences.Description
	}
	if updatedData.User.Preferences.Hobbies != "" {
		currentData.User.Preferences.Hobbies = updatedData.User.Preferences.Hobbies
	}

	// Convert the updated data back to JSON
	jsonData, err := json.Marshal(currentData)
	if err != nil {
		http.Error(w, "Failed to process data", http.StatusInternalServerError)
		return
	}

	// Update the file in the GitHub repository
	err = updateGitHubFile(jsonData)
	if err != nil {
		http.Error(w, "Failed to update GitHub file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("JSON is updated successfully"))
}

func fetchCurrentData() (Data, error) {
	apiURL := "https://raw.githubusercontent.com/sonnisanidev/aboutme_api/main/data.json"
	resp, err := http.Get(apiURL)
	if err != nil {
		return Data{}, err
	}
	defer resp.Body.Close()

	var data Data
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return Data{}, err
	}

	return data, nil
}

func updateGitHubFile(content []byte) error {
	apiURL := "https://api.github.com/repos/sonnisanidev/aboutme_api/contents/data.json"

	token := os.Getenv("GITHUB_TOKEN")
	print(token)
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable not set")
	}

	// First, get the current file to obtain its SHA
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("Failed to create GET request: %v", err)
	}
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to execute GET request: %v", err)
	}
	defer resp.Body.Close()

	var fileInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&fileInfo)
	if err != nil {
		return fmt.Errorf("Failed to decode file info: %v", err)
	}

	sha, ok := fileInfo["sha"].(string)
	if !ok {
		return fmt.Errorf("SHA not found in file info")
	}

	// Now update the file
	requestBody := map[string]interface{}{
		"message": "Update data.json",
		"content": base64.StdEncoding.EncodeToString(content),
		"sha":     sha,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("Failed to marshal request body: %v", err)
	}

	req, err = http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("Failed to create PUT request: %v", err)
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to execute PUT request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API request failed with status: %s, body: %s", resp.Status, string(body))
	}

	return nil
}
