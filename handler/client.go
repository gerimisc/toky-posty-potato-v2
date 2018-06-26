package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func NewClient() *SlidoHttp {

	c := &http.Client{
		Timeout: time.Second * 10, // Maximum of 2 secs
	}
	return &SlidoHttp{
		Client: c,
	}
}

func (c *SlidoHttp) ObtainIDs(eventCode string) (event_id string, uuid string) {

	slidoURL := "https://api.sli.do/v0.5/events?code=" + eventCode

	req, err := http.NewRequest(http.MethodGet, slidoURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := c.Client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, err := ioutil.ReadAll(res.Body)

	s, err := parseResponse([]byte(body))

	eid := strconv.FormatFloat(s[0]["event_id"].(float64), 'f', 0, 64)
	uuid = s[0]["uuid"].(string)
	fmt.Println(eid)
	fmt.Println(uuid)

	return eid, uuid
}

func (c *SlidoHttp) ObtainToken(eventID string, count int) []string {
	slidoURL := "https://app2.sli.do/api/v0.5/events/" + eventID + "/auth"
	var tokens []string
	fmt.Println("Obtaining " + strconv.Itoa(count) + " tokens")
	for i := 0; i < count; i++ {
		req, err := http.NewRequest(http.MethodPost, slidoURL, nil)
		if err != nil {
			log.Fatal(err)
		}

		res, getErr := c.Client.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		body, err := ioutil.ReadAll(res.Body)

		s, err := getToken([]byte(body))
		tokens = append(tokens, s.Access_token)
	}

	return tokens
}

func (c *SlidoHttp) UpVote(eventID string, questionID string, tokens ...string) *SlidoHttp {
	slidoVoteURL := "https://app2.sli.do/api/v0.5/events/" + eventID + "/questions/" + questionID + "/like"

	for _, token := range tokens {
		req, err := http.NewRequest(http.MethodPost, slidoVoteURL, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+token)

		_, getErr := c.Client.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

	}

	fmt.Printf("Voted on question %v of the event %v by %v times", questionID, eventID, len(tokens))
	return nil
}

func getToken(body []byte) (*SlidoJSONRes, error) {
	var s = new(SlidoJSONRes)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	return s, err
}

func parseResponse(body []byte) ([]map[string]interface{}, error) {
	s := []map[string]interface{}{}
	fmt.Println("Unmarshalling JSON Array response to Struct")
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("whoops:", err)
	}

	return s, err
}

type SlidoHttp struct {
	Client *http.Client
}

type SlidoJSONRes struct {
	Access_token string `json:"access_token"`
}

type SlidoIDs struct {
	Event_id string `json:"event_id"`
	Uuid     string `json:"uuid"`
}

type SlidoLoop []SlidoApi

type SlidoApi struct {
	EventID int    `json:"event_id"`
	UUID    string `json:"uuid"`
}
