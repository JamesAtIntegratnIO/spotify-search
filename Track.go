package main

import (
	"fmt"
	"strings"

	"github.com/zmb3/spotify/v2"
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
		artists := sArtistsToArtists(strack.Artists)

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

func TrackIDSearch(trackID string) (Track, error) {
	res, err := Search(trackID, "track_id")
	if err != nil {
		return Track{}, err
	}
	artists := sArtistsToArtists(res.Tracks.Tracks[0].Artists)
	t := Track{
		ID:         res.Tracks.Tracks[0].ID.String(),
		Name:       res.Tracks.Tracks[0].Name,
		PreviewURL: res.Tracks.Tracks[0].PreviewURL,
		Artists:    artists,
		Album: Album{
			Name:    res.Tracks.Tracks[0].Album.Name,
			ID:      res.Tracks.Tracks[0].Album.ID.String(),
			Artists: artists,
		},
	}
	return t, nil
}

type sArtists []spotify.SimpleArtist

func sArtistsToArtists(sArtists sArtists) (artists []Artist) {
	for _, a := range sArtists {
		b := Artist{
			ID:   a.ID.String(),
			Name: a.Name,
		}
		artists = append(artists, b)
	}
	return
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
