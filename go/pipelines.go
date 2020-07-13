package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
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

func isPrime(number int) bool {
	return big.NewInt(int64(number)).ProbablyPrime(number)
}

func processStory(in <-chan storyDetail) <-chan storyDetail {
	out := make(chan storyDetail)
	go func() {
		wg := &sync.WaitGroup{}
		for story := range in {
			wg.Add(1)
			go processStoryConcurrently(story, out, wg)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func processStoryConcurrently(story storyDetail, output chan<- storyDetail, wg *sync.WaitGroup) {
	story.IsPrime = isPrime(story.ID)
	output <- story
	wg.Done()
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
