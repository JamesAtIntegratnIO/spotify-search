# SPOTIFY SEARCH
Search for Track info for easier adding to your terraformed spotify playlists.

## Back Story
I've been playing with managing my spotify playlists in terraform. Start out as an excersize to learn more about terraform and get better at it. Turned into something I enjoy doing. The biggest pain point was finding the `track_id` to add the track. So I wrote this.

## Usage
There is only one command right now. Others could easily be added, but this is the MVP that does what I need to start.


Install it, then use the command as seen below.
> TODO: Add actual install instructions

### Help text
```
NAME:
   spotify-search track - Search for a track by title
                          Ex:spotify-search track -qty 15 -format JSON hotel california

USAGE:
   spotify-search track [command options] [arguments...]

OPTIONS:
   --format JSON, -f JSON   the format you want the results in, either JSON or `pretty` (default: pretty)
   --qty NUMBER, -q NUMBER  the NUMBER of results you want to see (default: 5)
```

### JSON Output
This matches how I add tracks to my spotify playlists. Not all of the info is used by terraform. It only needs the `track_id`. But I like knowing what song is associated with the track. I used to hand jam this mess. It was pain.
[One of my playlists](https://github.com/JamesAtIntegratnIO/spotify-playlist/blob/main/myJams.tf)
```
{
  artist      = "Eagles",
  album       = "Hotel California (2013 Remaster)"
  song        = "Hotel California - 2013 Remaster",
  track_id    = "40riOy7x9W7GXjyGp4pjAv"
  preview_url = "https://p.scdn.co/mp3-preview/412f7596ee68a616845f8b1269abaca5ad4e1b0d?cid=46aa92f8010943e6a4130cac7b47ba5d"
}
{
  artist      = "Eagles",
  album       = "Hell Freezes Over (Remaster 2018)"
  song        = "Hotel California - Live On MTV, 1994",
  track_id    = "2GpBrAoCwt48fxjgjlzMd4"
  preview_url = ""
}
{
  artist      = "Eagles",
  album       = "Hotel California (2013 Remaster)"
  song        = "Life in the Fast Lane - 2013 Remaster",
  track_id    = "6gXrEUzibufX9xYPk3HD5p"
  preview_url = "https://p.scdn.co/mp3-preview/b247da19e39dbd5338388409d22f500e8f28a847?cid=46aa92f8010943e6a4130cac7b47ba5d"
}
{
  artist      = "Eagles",
  album       = "Acoustic The Eagles & James Taylor"
  song        = "Hotel California",
  track_id    = "2ELuHKWwTMUAd9HSvMaI1j"
  preview_url = "https://p.scdn.co/mp3-preview/204a0069b62ba21415c03e89e45597586e2d7daf?cid=46aa92f8010943e6a4130cac7b47ba5d"
}
{
  artist      = "Gipsy Kings",
  album       = "!Volare! The Very Best of the Gipsy Kings"
  song        = "Hotel California (Spanish Mix)",
  track_id    = "4Rvhe8O90hFIExTJkdrRPM"
  preview_url = "https://p.scdn.co/mp3-preview/820ccab19c7232c0722defe858cd4cfe7c8acb5c?cid=46aa92f8010943e6a4130cac7b47ba5d"
}
```
### Pretty Output
Its not really pretty, but its kinda readable.
```
./spotify-search track -q 5 -f pretty hotel california
Artist: Eagles, Album: Hotel California (2013 Remaster), Track: Hotel California - 2013 Remaster, TrackID: 40riOy7x9W7GXjyGp4pjAv
Preview: https://p.scdn.co/mp3-preview/412f7596ee68a616845f8b1269abaca5ad4e1b0d?cid=46aa92f8010943e6a4130cac7b47ba5d
Artist: Eagles, Album: Hell Freezes Over (Remaster 2018), Track: Hotel California - Live On MTV, 1994, TrackID: 2GpBrAoCwt48fxjgjlzMd4
Preview: 
Artist: Eagles, Album: Hotel California (2013 Remaster), Track: Life in the Fast Lane - 2013 Remaster, TrackID: 6gXrEUzibufX9xYPk3HD5p
Preview: https://p.scdn.co/mp3-preview/b247da19e39dbd5338388409d22f500e8f28a847?cid=46aa92f8010943e6a4130cac7b47ba5d
Artist: Eagles, Album: Acoustic The Eagles & James Taylor, Track: Hotel California, TrackID: 2ELuHKWwTMUAd9HSvMaI1j
Preview: https://p.scdn.co/mp3-preview/204a0069b62ba21415c03e89e45597586e2d7daf?cid=46aa92f8010943e6a4130cac7b47ba5d
Artist: Gipsy Kings, Album: !Volare! The Very Best of the Gipsy Kings, Track: Hotel California (Spanish Mix), TrackID: 4Rvhe8O90hFIExTJkdrRPM
Preview: https://p.scdn.co/mp3-preview/820ccab19c7232c0722defe858cd4cfe7c8acb5c?cid=46aa92f8010943e6a4130cac7b47ba5d
```