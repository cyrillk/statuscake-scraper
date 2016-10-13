package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

// curl -H "API: YYY" -H "Username: ZZZ" -X GET https://www.statuscake.com/API/Tests/
// curl -H "API: YYY" -H "Username: ZZZ" -X GET "https://www.statuscake.com/API/Tests/Details?TestID=1408767"

type ShortTest struct {
	TestID int `json:"TestID"`
}

type DetailedTest struct {
	TestID       int      `json:"TestID"`
	WebsiteName  string   `json:"WebsiteName"`
	URI          string   `json:"URI"`
	ContactGroup string   `json:"ContactGroup"`
	Status       string   `json:"Status"`
	Tags         []string `json:"Tags"`
	Uptime       float32  `json:"Uptime"`
	CheckRate    int      `json:"CheckRate"`
}

func retrieveShortTests(username string, apikey string) []ShortTest {

	rawurl := "https://www.statuscake.com/API/Tests/"

	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("API", apikey)
	req.Header.Add("Username", username)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(resp.Status)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var shortTests []ShortTest
	err = json.Unmarshal(body, &shortTests)
	if err != nil {
		panic(err)
	}

	return shortTests
}

func retrieveDetailedTest(username string, apikey string, testID int) DetailedTest {

	rawurl := "https://www.statuscake.com/API/Tests/Details"

	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	q := u.Query()
	q.Set("TestID", strconv.Itoa(testID))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("API", apikey)
	req.Header.Add("Username", username)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic(resp.Status)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var detailedTest DetailedTest
	err = json.Unmarshal(body, &detailedTest)
	if err != nil {
		panic(err)
	}

	return detailedTest
}

func main() {

	var username string
	flag.StringVar(&username, "username", "", "Username")

	var apikey string
	flag.StringVar(&apikey, "apikey", "", "API Key")

	flag.Parse()

	shortTests := retrieveShortTests(username, apikey)

	var resultsWriter = tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)

	for _, item := range shortTests {
		detailed := retrieveDetailedTest(username, apikey, item.TestID)

		record := []string{
			strings.Join(detailed.Tags, ","),
			detailed.WebsiteName,
			detailed.URI,
			// detailed.Status,
			// strconv.Itoa(detailed.CheckRate),
		}

		fmt.Fprintln(resultsWriter, strings.Join(record, "\t")+"\t")
	}

	resultsWriter.Flush()
}
