package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type twitterReponse struct {
	MinPosition      string `json:"min_position,omitempty"`
	MaxPosition      string `json:"max_position,omitempty"`
	HasMoreItems     bool   `json:"has_more_items,omitempty"`
	ItemsHTML        string `json:"items_html,omitempty"`
	NewLatentCount   int    `json:"new_latent_count,omitempty"`
	NewTweetsBarHTML string `json:"new_tweets_bar_html,omitempty"`
}

type twitterSummary struct {
	DateTime string        `json:"dateTime,omitempty"`
	Feeds    []twitterFeed `json:"feeds,omitempty"`
}
type twitterFeed struct {
	User   string  `json:"user,omitempty"`
	Tweets []tweet `json:"tweets,omitempty"`
}
type tweet struct {
	ID       string `json:"id,omitempty"`
	DateTime string `json:"dateTime,omitempty"`
	URL      string `json:"url,omitempty"`
	Text     string `json:"text,omitempty"`
}

func main() {
	users := getUsers()
	summary := new(twitterSummary)
	summary.DateTime = time.Now().String()
	summary.Feeds = []twitterFeed{}

	for _, element := range users {
		feed := twitterFeed{User: element, Tweets: getUserFeed(element)}
		summary.Feeds = append(summary.Feeds, feed)
	}
	outFile, err := os.Create("output" + time.Now().Format("2006_01_02_15_04_05") + ".json")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	jsonbytes, _ := json.Marshal(summary)
	_, err = outFile.Write(jsonbytes)
	if err != nil {
		log.Fatal(err)
	}
}

func getUserFeed(user string) []tweet {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	var url = "https://twitter.com/i/profiles/show/" + user + "/timeline/tweets?include_available_features=1&include_entities=1&include_new_items_bar=true"
	var headers = []struct {
		header, value string
	}{
		{"Accept", "application/json, text/javascript, */*; q=0.01"},
		{"Referer", "https://twitter.com/" + user},
		{"User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8"},
		{"X-Twitter-Active-User", "yes"},
		{"X-Requested-With", "XMLHttpRequest"},
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, element := range headers {
		request.Header.Set(element.header, element.value)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	tr := twitterReponse{}
	json.NewDecoder(response.Body).Decode(&tr)

	tweets := []tweet{}
	document, err := goquery.NewDocumentFromReader(strings.NewReader(tr.ItemsHTML))
	document.Find(".stream-item").Each(func(index int, element *goquery.Selection) {
		text := element.Find(".tweet-text").Text()
		timeMsString, exists := element.Find("._timestamp").Attr("data-time-ms")
		timeMs, err := strconv.ParseInt(timeMsString, 10, 64)
		if err == nil {
		}
		time := time.Unix(0, timeMs*int64(time.Millisecond))
		if exists {
		}
		id, exists := element.Find(".js-permalink").Attr("data-conversation-id")
		url = ""
		if exists {
			url = "https://www.twitter.com/" + user + "/status/" + id
		}
		if isGoodTweet(text) {
			tweet := tweet{DateTime: time.String(), Text: text, ID: id, URL: url}
			tweets = append(tweets, tweet)
		}
	})
	return tweets
}

func isGoodTweet(tweet string) bool {
	programmerStrings := getKeyStrings()

	for _, element := range programmerStrings {
		if strings.Contains(tweet, element) {
			return true
		}
	}
	return false
}

func getUsers() []string {
	return []string{"__apf__",
		"jessfraz",
		"TheAmyCode",
		"Nick_Craver",
		"jessitron",
		"ashleymcnamara",
		"meganfabulous",
		"toomuchpete",
		"shanselman",
		"davidfowl",
		"kevinmontrose",
		"_tessr",
		"alicegoldfuss",
		"SaraJChipps",
		"soniagupta504",
		"nnja",
		"malwareunicorn",
		"Wright2Tweet",
		"jilljubs",
		"marcgravell",
		"jessysaurusrex",
		"marthakelly",
		"keirsten",
		"likeOMGitsFEDAY",
		"Eva",
		"GabeAul",
		"juberti",
		"makinde",
		"kylerush",
		"SlexAxton",
		"addyosmani",
		"JakeWharton",
		"chrisbanes",
		"rob_pike",
		"enneff",
		"smarx",
		"bittersweetryan",
		"ElizabethN",
		"RickAndMSFT",
		"haacked"}
}

func getKeyStrings() []string {
	return []string{".NET",
		"dotnet",
		"asp.net",
		"core",
		"microsoft",
		"github",
		"javascript",
		"java",
		"sql",
		"python",
		"dev",
		"development",
		"script",
		"programming",
		"tech",
		"devops",
		"developer",
		"code",
		"git",
		"c#",
		"c sharp",
		"c-sharp",
		"csharp"}
}
