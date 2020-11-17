package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strings"

	"time"

	"github.com/pborman/getopt/v2"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Team []string `yaml:"team"`
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Failed to get user home dir: %v\n", err)
	}
	listenFlag := getopt.Bool('l', "listen on port 8081")
	pickOne := getopt.Bool('1', "pick one volunteer (only if -l is not specified)")
	purpose := getopt.String('p', "", "select purpose (only if -l is not specified)")
	configFile := getopt.String('c', home+"/.constand.yaml", "select config file (default: ~/.constand.yaml)")
	getopt.Parse()
	// Get the remaining positional parameters
	if *listenFlag {
		http.HandleFunc("/", standupOrderHandler)
		http.HandleFunc("/pick/one/for", pickOneHandler)
		log.Fatal(http.ListenAndServe(":8081", nil))
	}
	date := time.Now().UTC()

	configYaml, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Printf("Failed to read config file %s: %v\n", *configFile, err)
		return
	}
	config := Config{}
	err = yaml.Unmarshal(configYaml, &config)
	if err != nil {
		log.Printf("Failed to parse yaml content of %s: %v\n", *configFile, err)
		return
	}
	if *pickOne {
		fmt.Print(getPickOneResponse(config.GetTeam(), date, *purpose))
		return
	}
	fmt.Print(getStandupOrderResponse(config.GetTeam(), date))
}

func (c *Config) GetTeam() []string {
	sort.Strings(c.Team)
	return c.Team
}

func pickOneHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	team := getTeam(r)
	purpose := getPurpose(r)
	date := getDate(r)

	response := getPickOneResponse(team, date, purpose)
	html := wrapHtml("Volunteer hunting!", response, response)

	fmt.Print(response)
	_, err := w.Write([]byte(html))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func wrapHtml(title, short, content string) string {
	html := "<!doctype html> <html> <head>"
	html += "<meta property=\"og:title\" content=\"" + title + "\">"
	html += "<meta property=\"og:description\" content=\"" + short + "\">"
	html += "</head><body><pre>" + content + "</pre></body></html>"
	return html
}

func getPickOneResponse(team []string, date time.Time, purpose string) string {
	pick := getOne(team, date, purpose)
	responsePurpose := ""
	if purpose != "" {
		responsePurpose = " for " + purpose

	}
	response := fmt.Sprintf("Here is your volunteer%s: %s\n", responsePurpose, pick)
	return response
}

func standupOrderHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	team := getTeam(r)
	date := getDate(r)

	response := getStandupOrderResponse(team, date)
	short := fmt.Sprintf("Order for %s: %v", date.Format("2006-01-02"), strings.Join(getStandupOrder(team, date), ", "))
	html := wrapHtml("Your standup order", short, response)
	fmt.Print(response)

	_, err := w.Write([]byte(html))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func getStandupOrderResponse(team []string, date time.Time) string {
	standupOrder := getStandupOrder(team, date)

	//fmt.Printf("Standup order for %s: %v\n", date.Format("2006-01-02"), standupOrder)

	response := "Standup order for "
	response += date.Format("2006-01-02")
	response += "\n============================\n\n"
	response += strings.Join(standupOrder, "\n")
	response += "\n"

	return response
}

func getTeam(r *http.Request) []string {
	q := r.URL.Query()
	team := q["team"]
	sort.Strings(team)
	return team
}

func getPurpose(r *http.Request) string {
	q := r.URL.Query()
	purpose := q["purpose"]
	if len(purpose) > 0 {
		return purpose[0]
	}
	return ""
}

func getDate(r *http.Request) time.Time {
	q := r.URL.Query()
	date := time.Now().UTC()
	if len(q["date"]) == 1 {
		fmt.Println(q["date"][0])
		urldate, err := time.Parse("2006-01-02", q.Get("date"))
		if err == nil {
			date = urldate
		}
	}
	return date
}

func removeElement(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func generateSeed(team []string, date time.Time, purpose string) int64 {
	dateteam := date.Format("2006-01-02")
	dateteam += strings.Join(team, "-")
	dateteam += purpose
	sum := sha256.Sum256([]byte(dateteam))
	return int64(binary.BigEndian.Uint64(sum[:]))
}

func getOne(team []string, date time.Time, purpose string) string {
	rand.Seed(generateSeed(team, date, purpose))
	return team[rand.Intn(len(team))]
}

func getStandupOrder(team []string, date time.Time) []string {
	rand.Seed(generateSeed(team, date, ""))
	teamsize := len(team)
	teamcopy := make([]string, teamsize)
	copy(teamcopy, team)
	standupOrder := []string{}
	for i := 0; i < teamsize; i++ {
		num := rand.Intn(len(teamcopy))
		standupOrder = append(standupOrder, teamcopy[num])
		teamcopy = removeElement(teamcopy, num)
	}
	return standupOrder
}
