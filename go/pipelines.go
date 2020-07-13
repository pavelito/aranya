package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

//generate a channel of StoryIDs
//generate a channel of DoneIDs
//print doneIDs

type storyDetail struct {
	Title string `json:"title"`
	ID    string `json:"id"`
	Time  int    `json:"time"`
}

func generateStoryIDs() <-chan int {
	stories := []int{}
	resp, _ := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&stories)

	out := make(chan int)
	go func() {
		for _, storyID := range stories {
			out <- storyID
		}
		close(out)
	}()
	return out
}

func processStory(in <-chan int) <-chan storyDetail {
	out := make(chan storyDetail)
	go func() {
		for storyID := range in {
			story := storyDetail{}
			resp, err := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json?print=pretty", storyID))
			if err != nil {
				log.Fatal(err)
			}

			defer resp.Body.Close()

			json.NewDecoder(resp.Body).Decode(&story)
			out <- story
		}
		close(out)
	}()
	return out
}

func main() {
	start := time.Now()
	storyIDs := generateStoryIDs()
	processStory(storyIDs)

	// for story := range processedStories {
	// 	log.Println("%v", story)
	// }
	elapsed := time.Since(start)
	log.Printf("Total seconds to finish - %s", elapsed)
}
