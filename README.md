# User Management API

## Description

This project is a simple User Management API written in Go. It provides endpoints to fetch and update user information stored in a GitHub repository.

## Features

- Fetch user information from a GitHub repository
- Update user information (name, age, description, and hobbies)
- Secure updates using GitHub token authentication

## Prerequisites

- Go 1.15 or higher
- GitHub Personal Access Token with repo scope

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/user-management-api.git
   ```

2. Navigate to the project directory:
   ```
   cd user-management-api
   ```

3. Set up your GitHub token as an environment variable:
   ```
   export GITHUB_TOKEN=your_github_token_here
   ```

## Usage

1. Run the server:
   ```
   go run main.go
   ```

2. The server will start running on `http://localhost:8080`

## API Endpoints

### GET /user

Fetches the current user information.

### PUT /update-user

Updates the user information. Accepts JSON payload with the following structure:

```json
{
  "user": {
    "name": {
      "first": "New First Name",
      "last": "New Last Name"
    },
    "age": 30,
    "preferences": {
      "description": "New description",
      "hobbies": "New hobbies"
    }
  }
}
```

## Configuration

The API uses a GitHub repository to store and retrieve user data. Make sure to update the following variables in the code:

- `url` in `handleUser` function: URL of the raw JSON file in your GitHub repository
- `apiURL` in `updateGitHubFile` function: GitHub API URL for your repository and file

## Security

This API uses a GitHub Personal Access Token for authentication. Make sure to keep your token secure and never commit it to version control.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).

---
## The API-Call Examples

```
curl.exe -X PUT -H "Content-Type: application/json" --data-binary "@sentdata.json" http://localhost:8080/update-user
```