package main

import (
	"context"
	"errors"
	"log"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

func Search(searchString string, searchType string) (result *spotify.SearchResult, err error) {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     config.SPOTIFY_ID,
		ClientSecret: config.SPOTIFY_SECRET,
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)
	switch searchType {
	case "artist":
		return client.Search(ctx, searchString, spotify.SearchTypeArtist)
	case "track":
		return client.Search(ctx, searchString, spotify.SearchTypeTrack)
	case "album":
		return client.Search(ctx, searchString, spotify.SearchTypeAlbum)
	default:
		err = errors.New("Invalid search type: Please use 'artist', 'track', or 'album'")
		return nil, err
	}
}
