package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	stories := getNewStories()
	for _, storyID := range stories {
		story := getStoryDetails(storyID)
		fmt.Println(story.Title)
		countdown(story.Time)
	}
}

type storyDetail struct {
	Title string `json:"title"`
	ID    string `json:"id"`
	Time  int    `json:"time"`
}

func countdown(number int) {
	counter := number
	for counter > 0 {
		counter = counter - 1
	}
	println(fmt.Sprintf("Countdown done for %d", number))
}

func getNewStories() []int {
	stories := []int{}
	resp, _ := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&stories)
	return stories
}

func getStoryDetails(storyID int) storyDetail {
	story := storyDetail{}
	resp, _ := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json?print=pretty", storyID))
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&story)
	return story
}
