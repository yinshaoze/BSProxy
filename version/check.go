package version

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func printErr(err error) {
	log.Printf("Error to check for update, caution: %v.", err.Error())
	log.Println(`You can check it yourself at https://github.com/yinshaoze/BSProxy/releases`)
}

func CheckUpdate() {
	resp, err := http.Get(`https://cdn.jsdelivr.net/gh/yinshaoze/BSProxy@master/version/version.go`)
	if err != nil {
		printErr(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		printErr(err)
		return
	}
	if strings.Contains(string(body), Version) {
		fmt.Println("Your BSProxy is up-to-date. Have fun!")
	} else {
		fmt.Println("Your BSProxy is out of date! Check for the latest version at https://github.com/yinshaoze/BSProxy/releases")
	}
}
