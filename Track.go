package main

import (
	"fmt"
	"strings"
)

type Album struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
	Tracks  []Track  `json:"tracks"`
	// define artists
}

type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Track struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Artists    []Artist `json:"artists"`
	Album      Album    `json:"album"`
	PreviewURL string   `json:"preview_url"`
	Total      int
	Err        error
}

func TrackSearch(trackName string) ([]Track, error) {
	var tracks []Track
	res, err := Search(trackName, "track")
	if err != nil {
		return nil, err
	}

	for _, strack := range res.Tracks.Tracks {
		artists := []Artist{}
		for _, artist := range strack.Artists {
			a := Artist{
				ID:   artist.ID.String(),
				Name: artist.Name,
			}
			artists = append(artists, a)
		}

		t := Track{
			ID:         strack.ID.String(),
			Name:       strack.Name,
			PreviewURL: strack.PreviewURL,
			Artists:    artists,
			Album: Album{
				Name:    strack.Album.Name,
				ID:      strack.Album.ID.String(),
				Artists: artists,
			},
		}
		tracks = append(tracks, t)
	}
	return tracks, nil
}

func TrackView(tracks []Track, qty int, format string) (msgs []string) {
	jsonfmt := "{\n  artist      = \"%s\",\n  album       = \"%s\"\n  song        = \"%s\",\n  track_id    = \"%s\"\n  preview_url = \"%s\"\n}"
	prettyfmt := "Artist: %s, Album: %s, Track: %s, TrackID: %s\nPreview: %s"
	for _, t := range tracks[:qty] {
		if strings.ToLower(format) == "json" || strings.ToLower(format) == "tf" || strings.ToLower(format) == "terraform" {
			msgs = append(msgs, fmt.Sprintf(jsonfmt, t.Artists[0].Name, t.Album.Name, t.Name, t.ID, t.PreviewURL))
		} else {
			msgs = append(msgs, fmt.Sprintf(
				prettyfmt, t.Artists[0].Name, t.Album.Name, t.Name, t.ID, t.PreviewURL))
		}
	}
	return msgs
}
