package main

import (
	"exp/slido_upvote/handler"
	"fmt"
	"os"
	"strconv"
)

func main() {

	if len(os.Args[1:]) < 3 {
		panic("Please supply the arguments: ./main <eventID> <questionID> <votes>")
	} else {

		var slidoClient *handler.SlidoHttp
		slidoClient = handler.NewClient()
		eventID := os.Args[1]
		questionID := os.Args[2]
		upvoteCount, _ := strconv.Atoi(os.Args[3])
		//upVoteCount, _ := strconv.Atoi(os.Args[2])
		tokenSlice := slidoClient.ObtainToken(eventID, upvoteCount)
		fmt.Println(tokenSlice)
		slidoClient.UpVote(eventID, questionID, tokenSlice...)
	}
}
