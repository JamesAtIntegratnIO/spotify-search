package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	cli "github.com/urfave/cli/v2"

	c "github.com/jamesatintegratnio/spotify-search/config"

	"golang.org/x/oauth2/clientcredentials"
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

type Search struct {
	Song   string
	Format string
}

var config = *c.SetupConfig()
var t *template.Template

var auth = spotifyauth.New(spotifyauth.WithRedirectURL(config.AUTH_URL))

func main() {

	app := &cli.App{
		Name:        "Spotify Search",
		HelpName:    "spotify-search",
		Description: "A tool for searching for spotify track info.",

		Commands: []*cli.Command{
			{
				Name:    "track",
				Aliases: []string{"s"},
				Usage:   "Search for a track by title\nEx:spotify-search track -qty 15 -format JSON hotel california",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "qty",
						Aliases:     []string{"q"},
						Usage:       "the `NUMBER` of results you want to see",
						DefaultText: "5",
					},
					&cli.StringFlag{
						Name:        "format",
						Aliases:     []string{"f"},
						Usage:       "the format you want the results in, either `JSON` or `pretty`",
						DefaultText: "pretty",
					},
				},
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() > 0 {
						tracks, err := trackSearch(strings.Join(cCtx.Args().Slice(), " "))
						if err != nil {
							return err
						}
						msgs := trackView(tracks, cCtx.Int("qty"), cCtx.String("format"))
						for _, msg := range msgs {
							fmt.Println(msg)
						}
					}
					return nil
				},
			},
			{
				Name:    "webserver",
				Aliases: []string{"web"},
				Usage:   "search from the comfort of your browser",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "port",
						Aliases:     []string{"p"},
						Usage:       "The `PORT` to run the webserver on",
						DefaultText: "8080",
					},
				},
				Action: func(cCtx *cli.Context) error {
					tmpl := template.Must(template.ParseFiles("templates/index.html"))

					http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
						if r.Method != http.MethodPost {
							tmpl.Execute(w, nil)
							return
						}

						details := Search{
							Song:   r.FormValue("track"),
							Format: r.FormValue("format"),
						}
						tracks, err := trackSearch(details.Song)
						if err != nil {
							fmt.Println(err)
						}
						format, _ := strconv.ParseBool(details.Format)

						tmpl.Execute(w, struct {
							Success bool
							Format  bool
							Tracks  []Track
						}{true, format, tracks})
					})
					http.ListenAndServe(":8080", nil)
					return nil
				},
			},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func trackSearch(trackName string) ([]Track, error) {
	var tracks []Track
	res, err := search(trackName, "track")
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

func trackView(tracks []Track, qty int, format string) (msgs []string) {
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

func search(searchString string, searchType string) (result *spotify.SearchResult, err error) {
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
