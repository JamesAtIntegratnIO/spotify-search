package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	cli "github.com/urfave/cli/v2"

	c "github.com/jamesatintegratnio/spotify-search/config"
)

type Form struct {
	Song   string
	Artist string
	Format string
}

var config = *c.SetupConfig()
var auth = spotifyauth.New(spotifyauth.WithRedirectURL(config.AUTH_URL))

// var tmpl = func() *template.Template {
// 	t := template.New("")
// 	err := filepath.Walk("templates/", func(path string, info os.FileInfo, err error) error {
// 		if strings.Contains(path, ".html") {
// 			fmt.Println(path)
// 			_, err = t.ParseFiles(path)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 		}
// 		return err
// 	})

// 	if err != nil {
// 		panic(err)
// 	}
// 	return t
// }()

func main() {

	app := &cli.App{
		Name:        "Spotify Search",
		HelpName:    "spotify-search",
		Description: "A tool for searching for spotify track info.",

		Commands: []*cli.Command{
			{
				Name:    "track",
				Aliases: []string{"t"},
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
						tracks, err := TrackSearch(strings.Join(cCtx.Args().Slice(), " "))
						if err != nil {
							return err
						}
						msgs := TrackView(tracks, cCtx.Int("qty"), cCtx.String("format"))
						for _, msg := range msgs {
							fmt.Println(msg)
						}
					}
					return nil
				},
			},
			{
				Name:    "track_id",
				Aliases: []string{"tid"},
				Usage:   "Search for track details of a `TRACK ID`",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "format",
						Aliases:     []string{"f"},
						Usage:       "the format you want the results in, either `JSON` or `pretty`",
						DefaultText: "pretty",
					},
				},
				Action: func(cCtx *cli.Context) error {
					track, err := TrackIDSearch(cCtx.Args().First())
					if err != nil {
						fmt.Println(err)
					}
					msg := TrackView([]Track{track}, 1, cCtx.String("format"))
					for _, msg := range msg {
						fmt.Println(msg)
					}
					return nil
				},
			},
			{
				Name:    "toptracks",
				Aliases: []string{"tt"},
				Usage:   "Get the top tracks of `ARTISTS NAME`",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "format",
						Aliases:     []string{"f"},
						Usage:       "the format you want the results in, either `JSON` or `pretty`",
						DefaultText: "pretty",
					},
				},
				Action: func(cCtx *cli.Context) error {
					format := "pretty"
					if cCtx.String("format") != "" {
						format = cCtx.String("format")
					}
					topTracks, err := getTopTracksByArtist(strings.Join(cCtx.Args().Slice(), " "))
					if err != nil {
						fmt.Println(err)
					}
					msg := TrackView(topTracks, 10, format)
					for _, msg := range msg {
						fmt.Println(msg)
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
					port := 8080
					if cCtx.NumFlags() > 0 {
						port = cCtx.Int("port")
					}
					tmpl := ParseTemplates()
					fmt.Printf("Web server running at: http://localhost:%d", port)
					http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

						if r.Method != http.MethodPost {
							tmpl.ExecuteTemplate(w, "main", nil)
							return
						}

						details := Form{
							Song:   r.FormValue("track"),
							Format: r.FormValue("format"),
						}
						tracks, err := TrackSearch(details.Song)
						if err != nil {
							fmt.Println(err)
						}
						format, _ := strconv.ParseBool(details.Format)

						tmpl.ExecuteTemplate(w, "main", struct {
							Success bool
							Format  bool
							Tracks  []Track
						}{true, format, tracks})
					})
					http.HandleFunc("/toptracks", func(w http.ResponseWriter, r *http.Request) {
						if r.Method != http.MethodPost {
							tmpl.ExecuteTemplate(w, "top_tracks", nil)
							return
						}
						details := Form{
							Artist: r.FormValue("artist"),
							Format: r.FormValue("format"),
						}
						topTracks, err := getTopTracksByArtist(details.Artist)
						if err != nil {
							fmt.Println(err)
						}
						format, _ := strconv.ParseBool(details.Format)

						tmpl.ExecuteTemplate(w, "top_tracks", struct {
							Success bool
							Format  bool
							Tracks  []Track
						}{true, format, topTracks})

					})
					http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

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

type Page struct {
	Title string
}

func ParseTemplates() *template.Template {
	templ := template.New("")
	err := filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				log.Println("Parse error", err)
			}
		}

		return err
	})

	if err != nil {
		panic(err)
	}

	return templ
}
