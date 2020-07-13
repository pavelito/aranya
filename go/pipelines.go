package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

//generate a channel of StoryIDs
//generate a channel of DoneIDs
//print doneIDs

type storyDetail struct {
	Title   string `json:"title"`
	ID      int    `json:"id"`
	Time    int    `json:"time"`
	IsPrime bool   `json:"is_prime"`
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

func getStoryDetail(in <-chan int) <-chan storyDetail {
	out := make(chan storyDetail)
	var wg sync.WaitGroup
	const concurrency = 25 //Bounded Concurrency
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			getStoryDetailConcurrently(in, out)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func getStoryDetailConcurrently(in <-chan int, output chan<- storyDetail) {
	for storyID := range in {
		story := storyDetail{}
		resp, err := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", storyID))
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		json.NewDecoder(resp.Body).Decode(&story)
		output <- story
	}
}

func isPrime(number int) bool {
	for i := 2; i <= number/2; i++ {
		if number%i == 0 {
			return false
		}
	}
	return true
}

func processStory(in <-chan storyDetail) <-chan storyDetail {
	out := make(chan storyDetail)
	var wg sync.WaitGroup
	const concurrency = 10 //Bounded Concurrency
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			processStoryConcurrently(in, out)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func processStoryConcurrently(in <-chan storyDetail, output chan<- storyDetail) {
	for story := range in {
		story.IsPrime = isPrime(story.ID)
		output <- story
	}
}

func main() {
	start := time.Now()
	storyIDs := generateStoryIDs()
	storyDetails := getStoryDetail(storyIDs)
	processedStories := processStory(storyDetails)

	for story := range processedStories {
		log.Println("%v", story)
	}
	elapsed := time.Since(start)
	log.Printf("Total seconds to finish - %s", elapsed)
}
