package main

import (
	"fmt"
	"github.com/sugiantoaudi/toky-posty-potato-v2/handler"
	"os"
	"strconv"
)

func main() {

	if len(os.Args[1:]) < 3 {
		panic("Please supply the arguments: ./main <eventID> <questionID> <votes>")
	} else {

		var slidoClient *handler.SlidoHttp
		slidoClient = handler.NewClient()
		eventCode := os.Args[1]
		questionID := os.Args[2]

		event_id, uuid := slidoClient.ObtainIDs(eventCode)

		upvoteCount, _ := strconv.Atoi(os.Args[3])

		//upVoteCount, _ := strconv.Atoi(os.Args[2])

		tokenSlice := slidoClient.ObtainToken(uuid, upvoteCount)
		fmt.Println(tokenSlice)
		slidoClient.UpVote(event_id, questionID, tokenSlice...)
	}
}
