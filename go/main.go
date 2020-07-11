package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	stories := getNewStories()
	const workers = 25

	wg := new(sync.WaitGroup)
	in := make(chan int, 2*workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for storyID := range in {
				story := getStoryDetails(storyID)
				log.Println(story.Title)
				go countdown(story.Time)
			}
		}()
	}

	for _, storyID := range stories {
		in <- storyID
	}

	close(in)
	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Total seconds to finish - %s", elapsed)
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
	log.Println(fmt.Sprintf("Countdown done for %d", number))
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
	resp, err := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json?print=pretty", storyID))
	if err != nil {
		log.Fatal(err)
		return storyDetail{}
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&story)
	return story

}
