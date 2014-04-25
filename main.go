package main

import "fmt"
import "io"
import "os"
import "strings"
import "net/url"
import "net/http"
import "flag"
import "log"

var instance string
var user string
var password string

func init() {
	flag.StringVar(&instance, "instance", "", "V1 Instance")
	flag.StringVar(&user, "user", "", "V1 Instance user")
	flag.StringVar(&password, "password", "", "V1 Instance password")
}

//import "github.com/joshlf13/term"

func getUrl() *url.URL {
	parsedURL, err := url.Parse(instance)
	if err != nil {
		log.Fatalf("Unable to parse \"%s\" as a url", instance)
	}

	if user != "" && password != "" {
		parsedURL.User = url.UserPassword(user, password)
	} else if user != "" {
		parsedURL.User = url.User(user)
	}

	return parsedURL
}

func main() {
	flag.Parse()

	commits := ScanGit(DefaultOptions, "head --not origin/14.1")
	fmt.Printf("%s\n", commits.AllReferences())
	V1()
}

func V1() {
	req := strings.NewReader(`{"from": "Member", }`)

	resp, err := http.Post(getUrl().String(), "application/json", req)
	if err == nil {
		io.Copy(os.Stdout, resp.Body)
		fmt.Println(resp)
	} else {
		fmt.Println(err)
	}

}
