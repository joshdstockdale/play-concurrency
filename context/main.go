package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {

	start := time.Now()
	ctx := context.WithValue(context.Background(), "username", "josh")
	userID, err := fetchUserID(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("the response took %v -> %+v\n", time.Since(start), userID)
}

func fetchUserID(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()

	val := ctx.Value("username")
	fmt.Printf("username is %s\n", val)

	type result struct {
		userId string
		err    error
	}
	resultch := make(chan result, 1)

	go func() {
		res, err := thirdPartyHTTPCall()
		resultch <- result{
			userId: res,
			err:    err,
		}
	}()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-resultch:
		return res.userId, res.err
	}
}

func thirdPartyHTTPCall() (string, error) {
	time.Sleep(time.Millisecond * 90)
	return "user id 1", nil
}
