package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const delayCheck = 10

func main() {
	for {
		monitor()
		time.Sleep(delayCheck * time.Second)
	}
}

func monitor() {
	listUrls := getUrls()

	for _, url := range listUrls {
		printAndLog("request: " + url)

		response, err := http.Get(url)
		if err != nil {
			printAndLog("ERROR")
			continue
		}

		if response.StatusCode != http.StatusOK {
			printAndLog("ERROR")
			continue
		}

		var health healthcheck

		body, _ := ioutil.ReadAll(response.Body)

		err = json.Unmarshal(body, &health)
		if err != nil {
			printAndLog("ERROR")
			continue
		}

		if health.IsHealthy {
			printAndLog("OK")
		} else {
			printAndLog("ERROR")
		}
	}
}

func printAndLog(text string) {
	fmt.Println(text)
	writeLog(text)
}

func writeLog(text string) {

	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + text + "\n")

	arquivo.Close()
}

func getUrls() []string {
	file, err := os.Open("healthCheckUrl.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	listUrls := []string{}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		listUrls = append(listUrls, strings.TrimSpace(line))
		if err == io.EOF {
			break
		}
	}

	file.Close()

	return listUrls
}

type healthcheck struct {
	IsHealthy bool
}
