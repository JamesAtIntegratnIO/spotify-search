package main

import (
	"context"
	"log"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

func getTopTracksByArtist(artistName string) ([]Track, error) {
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

	artist, err := GetArtistID(artistName)

	topTracks, err := client.GetArtistsTopTracks(ctx, artist.ID, "US")
	if err != nil {
		return nil, err
	}
	var tt []Track
	for _, t := range topTracks {
		tt = append(tt, Track{
			ID:         string(t.ID),
			Name:       t.Name,
			PreviewURL: t.PreviewURL,
			Artists:    sArtistsToArtists(t.Artists),
			Album: Album{
				Name: t.Album.Name,
				ID:   string(t.Album.ID),
			},
		})
	}

	return tt, nil
}

func GetArtistID(artistName string) (spotify.FullArtist, error) {

	res, err := Search(artistName, "artist")
	if err != nil {
		return spotify.FullArtist{}, err
	}

	artist := res.Artists.Artists[0]

	return artist, nil
}
