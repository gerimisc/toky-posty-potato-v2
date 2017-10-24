package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func NewClient() *SlidoHttp {

	c := &http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}
	return &SlidoHttp{
		Client: c,
	}
}

func (c *SlidoHttp) ObtainToken(eventID string, count int) []string {
	slidoURL := "https://app.sli.do/api/v0.4/events/" + eventID + "/auth"
	var tokens []string
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
	slidoVoteURL := "https://app.sli.do/api/v0.4/events/" + eventID + "/questions/" + questionID + "/like"

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

	fmt.Printf("Voted on %v of the event %v by %v times", questionID, eventID, len(tokens))
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

type SlidoHttp struct {
	Client *http.Client
}

type SlidoJSONRes struct {
	Access_token string `json:"access_token"`
}
