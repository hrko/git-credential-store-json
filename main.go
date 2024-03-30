package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Accoordig to the man page of the `git-credential`,
// > In all cases, all bytes are treated as-is (i.e., there is no quoting, and one cannot transmit a value with newline or NUL in it)
// However, I want JSON to be human-readable, so I will use string instead of bytes.
type Credential struct {
	Protocol          string   `json:"protocol"`
	Host              string   `json:"host"`
	Path              string   `json:"path,omitempty"`
	Username          string   `json:"username,omitempty"`
	Password          string   `json:"password,omitempty"`
	PasswordExpiryUTC string   `json:"password_expiry_utc,omitempty"`
	OAuthRefreshToken string   `json:"oauth_refresh_token,omitempty"`
	WWWAuth           []string `json:"wwwauth,omitempty"`
}

var verbose bool

func main() {
	log.SetPrefix("git-credential-store-json: ")

	var credsFilePath string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get the home directory: %s\n", err)
	}
	defaultCredentialsFile := filepath.Join(homeDir, ".git-credentials.json")

	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.StringVar(&credsFilePath, "f", defaultCredentialsFile, "path to the credentials file")
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		log.Fatalf("action (get, store, erase) is required")
	}

	action := args[0]
	switch action {
	case "get":
		getCredential(credsFilePath)
	case "store":
		storeCredential(credsFilePath)
	case "erase":
		eraseCredential(credsFilePath)
	default:
		log.Fatalf("Invalid action: %s\n", action)
	}
}

func (c *Credential) Print() {
	fmt.Printf("protocol=%s\n", c.Protocol)
	fmt.Printf("host=%s\n", c.Host)
	if c.Path != "" {
		fmt.Printf("path=%s\n", c.Path)
	}
	if c.Username != "" {
		fmt.Printf("username=%s\n", c.Username)
	}
	if c.Password != "" {
		fmt.Printf("password=%s\n", c.Password)
	}
	if c.PasswordExpiryUTC != "" {
		fmt.Printf("password_expiry_utc=%s\n", c.PasswordExpiryUTC)
	}
	if c.OAuthRefreshToken != "" {
		fmt.Printf("oauth_refresh_token=%s\n", c.OAuthRefreshToken)
	}
	for _, wwwauth := range c.WWWAuth {
		fmt.Printf("wwwauth[]=%s\n", wwwauth)
	}
}

func readCredentialFromStdin() Credential {
	var credential Credential

	// Read and parse stdin until blank line
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		var key, value string
		parts := strings.SplitN(line, "=", 2)
		key = string(parts[0])
		if len(parts) != 2 {
			value = ""
		} else {
			value = parts[1]
		}

		switch key {
		case "protocol":
			credential.Protocol = value
		case "host":
			credential.Host = value
		case "path":
			credential.Path = value
		case "username":
			credential.Username = value
		case "password":
			credential.Password = value
		case "password_expiry_utc":
			credential.PasswordExpiryUTC = value
		case "oauth_refresh_token":
			credential.OAuthRefreshToken = value
		case "wwwauth[]":
			credential.WWWAuth = append(credential.WWWAuth, value)
		default:
			log.Printf("Unknown key: %s\n", key)
		}
	}

	if verbose {
		log.Println("Read credential from stdin:")
		log.Println(credential)
	}

	return credential
}

func readCredentialsFromFile(path string) []Credential {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("readCredentialsFromFile: Failed to open the credentials file: %s\n", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var credentials []Credential
	err = decoder.Decode(&credentials)
	if err != nil && io.EOF != err {
		log.Fatalf("readCredentialsFromFile: Failed to decode the credentials file: %s\n", err)
	}

	if verbose {
		log.Println("Read credentials from file:")
		for _, c := range credentials {
			log.Println(c)
		}
	}

	return credentials
}

func saveCredentialsToFile(path string, credentials []Credential) {
	if verbose {
		log.Println("Save credentials to file:")
		for _, c := range credentials {
			log.Println(c)
		}
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("saveCredentialsToFile: Failed to open the credentials file: %s\n", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(credentials)
	if err != nil {
		log.Fatalf("saveCredentialsToFile: Failed to encode the credentials file: %s\n", err)
	}
}

func getCredential(path string) {
	creds := readCredentialsFromFile(path)
	credInput := readCredentialFromStdin()

	for _, c := range creds {
		if c.Protocol == credInput.Protocol && c.Host == credInput.Host {
			c.Print()
			if verbose {
				log.Println("Found credential:")
				log.Println(c)
			}
			return
		}
	}

	if verbose {
		log.Println("Credential not found")
	}
}

func storeCredential(path string) {
	creds := readCredentialsFromFile(path)
	credInput := readCredentialFromStdin()

	found := false
	for i, c := range creds {
		if c.Protocol == credInput.Protocol && c.Host == credInput.Host {
			creds[i] = credInput
			found = true
			if verbose {
				log.Println("Updated credential:")
				log.Println(creds[i])
			}
			break
		}
	}
	if !found {
		creds = append(creds, credInput)
		if verbose {
			log.Println("Added credential:")
			log.Println(credInput)
		}
	}

	saveCredentialsToFile(path, creds)
}

func eraseCredential(path string) {
	creds := readCredentialsFromFile(path)
	credInput := readCredentialFromStdin()

	for i, c := range creds {
		if c.Protocol == credInput.Protocol && c.Host == credInput.Host {
			creds = append(creds[:i], creds[i+1:]...)
			if verbose {
				log.Println("Erased credential:")
				log.Println(c)
			}
			break
		}
	}

	saveCredentialsToFile(path, creds)
}
