package main

import (
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Rate struct {
	TWD float32 `json:"TWD"`
}

type ApiResponse struct {
	Base  string `json:"base"`
	Date  string `json:"date"`
	Rates Rate   `json:"rates"`
}

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("amount is required.")
		return
	}
	amount, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		fmt.Println("amount is invalid.")
		return
	}
	req, err := http.NewRequest("GET", "https://api.exchangerate.host/latest", nil)
	if err != nil {
		fmt.Println("failed to build request.")
		return
	}
	q := req.URL.Query()
	q.Add("base", "USD")
	q.Add("symbols", "TWD")
	q.Add("amount", strconv.FormatFloat(amount, 'f', 6, 64))
	req.URL.RawQuery = q.Encode()

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	client := http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	s.Start()
	s.FinalMSG = "Done! \n"
	resp, err := client.Do(req)
	if err != nil {
		s.Stop()
		fmt.Println("failed to get a response.")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.Stop()
		fmt.Println("failed to read body.")
		return
	}

	var data ApiResponse
	json.Unmarshal(body, &data)
	s.Stop()

	fmt.Printf("%f (USD) = %f (TWD)\n", amount, data.Rates.TWD)
}
