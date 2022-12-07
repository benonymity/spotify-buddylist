# Spotify Friend Activity

This program serves a webpage that imitates Spotify's official friend activity list for the account token you supply.

It serves a web client that calls a backend proxy on your server which talks to Spotify's unofficial activity API using a privledged access token and displays it in a webpage with similar styling to Spotify.

# Usage

### Manual:

Clone the repo:

```
git clone https://github.com/benonymity/spotify-buddylist.git
```

Run

```
./build.sh
```

Then fill in your sp_dc token in `sp_dc.txt`

Enter the `api` folder and run

```
./server
```

The page is now being served on `http://localhost:10000`

### Docker:

Pull `benonymity/spotify-buddylist:latest` and run it with an environment variable SP_DC containing your token. Bind port 10000 to wherever you want to serve the webpage, and if you want to persist listening history bind `activity.db` to a local file.

Example Docker command:

```
docker run -d -p 10000:10000 -v spotify/activity.db:/activity.db --name spotify-buddylist benonymity/spotify-buddylist:latest
```

# Getting your sp_dc token

Login to the [web player](https://open.spotify.com/) in an incognito window and Inspect Element. Open the network tab and search `clienttoken`, then copy the token from the response. This is your one-year-expiry sp_dc token.

# Todo:

- [x] Cache API results in backend to avoid DDOS
- [x] Save history to DB
- [x] Create UI for viewing past activity
- [ ] Allow following of new users from webpage
