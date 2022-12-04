package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "sort"
    "strings"
)

// Subdomain represents a subdomain found using crt.sh
type Subdomain struct {
    Name string `json:"name_value"`
}

func main() {
    // Check if the user provided a domain name
    if len(os.Args) < 2 {
        // If no domain was provided, print an error message and exit
        fmt.Println("Error: No domain provided.")
        os.Exit(1)
    }

    // Parse the domain name from the command-line arguments
    domain := os.Args[1]

    // Use crt.sh to find subdomains for the given domain
    resp, err := http.Get(fmt.Sprintf("https://crt.sh/?q=%.%s&output=json", domain))
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer resp.Body.Close()

    // Check the server response
    if resp.StatusCode != http.StatusOK {
        fmt.Printf("Error: Server returned HTTP %d\n", resp.StatusCode)
        os.Exit(1)
    }
    contentType := resp.Header.Get("Content-Type")
    if !strings.Contains(contentType, "application/json") {
        fmt.Printf("Error: Invalid Content-Type: %s\n", contentType)
        os.Exit(1)
    }

    // Parse the JSON response from crt.sh
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    var subdomains []Subdomain
    err = json.Unmarshal(body, &subdomains)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // Extract the subdomain names and sort them
    names := make([]string, len(subdomains))
    for i, subdomain := range subdomains {
        names[i] = strings.TrimPrefix(subdomain.Name, "*.")
    }
    sort.Strings(names)

    // Print the subdomains
    fmt.Printf("Subdomains for %s:\n", domain)
    for _, name := range names {
        fmt.Println(name)
    }
}
