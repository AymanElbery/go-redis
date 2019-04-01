package main

import (
	"fmt"
	"strconv"

	// Import the Radix.v2 redis package.
	"log"

	"github.com/mediocregopher/radix.v2/redis"
)

// Define a custom struct to hold Album data.
type Album struct {
	Title  string
	Artist string
	Price  float64
	Likes  int
}

func main() {
	// Establish a connection to the Redis server listening on port 6379 of the local machine ...
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}

	// Importantly, use defer to ensure the connection is always properly closed before exiting the main() function ...
	defer conn.Close()

	// Send our command across the connection ...
	resp := conn.Cmd("HMSET", "album:1", "title", "Let's Go", "artist", "Ayman Elbery", "price", 25.95, "likes", 25)

	// Check the Err field of the *Resp object for any errors.
	if resp.Err != nil {
		log.Fatal(resp.Err)
	}

	// Fetch all album fields with the HGETALL command ...
	allAlbum, err := conn.Cmd("HGETALL", "album:1").Map()
	if err != nil {
		log.Fatal(err)
	}

	ab, err := populateAlbum(allAlbum)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s by %s: $%.2f (%d likes)\n", ab.Title, ab.Artist, ab.Price, ab.Likes)

}

func populateAlbum(allAlbum map[string]string) (*Album, error) {
	var err error
	ab := new(Album)
	ab.Title = allAlbum["title"]
	ab.Artist = allAlbum["artist"]

	ab.Price, err = strconv.ParseFloat(allAlbum["price"], 64)
	if err != nil {
		return nil, err
	}

	ab.Likes, err = strconv.Atoi(allAlbum["likes"])
	if err != nil {
		return nil, err
	}
	return ab, nil
}
