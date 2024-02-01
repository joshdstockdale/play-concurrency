package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type UserProfile struct {
	ID       int
	Comments []string
	Likes    int
	Friends  []int
}

func main() {
	start := time.Now()
	userProfile, err := handleGetUserProfile(10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Profile is %+v\n", userProfile)
	fmt.Println("Fetching took...", time.Since(start))
}

type Response struct {
	data any
	err  error
}

func handleGetUserProfile(id int) (*UserProfile, error) {

	var (
		respch = make(chan Response, 3)
		wg     = &sync.WaitGroup{}
	)

	go getComments(id, respch, wg)
	go getLikes(id, respch, wg)
	go getFriends(id, respch, wg)
	wg.Add(3)
	wg.Wait()
	close(respch)

	userProfile := &UserProfile{}
	for resp := range respch {
		if resp.err != nil {
			return nil, resp.err
		}
		switch data := resp.data.(type) {
		case int:
			userProfile.Likes = data
		case []int:
			userProfile.Friends = data
		case []string:
			userProfile.Comments = data
		}
	}

	return userProfile, nil
}

func getComments(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 200)
	comments := []string{
		"Good job!", "Interesting",
	}
	respch <- Response{
		data: comments,
		err:  nil,
	}
	wg.Done()
}

func getLikes(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 200)
	respch <- Response{
		data: 55,
		err:  nil,
	}
	wg.Done()
}

func getFriends(id int, respch chan Response, wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 100)
	friendsIds := []int{11, 33, 124}
	respch <- Response{
		data: friendsIds,
		err:  nil,
	}
	wg.Done()
}
