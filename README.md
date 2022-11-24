# Spotify Friend Activity

This program serves a web client that calls a backend proxy on your server which talks to Spotify's unofficial activity API and displays it in a webpage with similar styling to Spotify.

# Usage

---

### Manual:

Clone the repo:

```
git clone https://github.com/benonymity/spotify-buddylist.git
```

Enter the `activity` folder and run

```
npm run build
```

Then fill in your sp_dc token in `sp_dc.txt`.

Enter the `api` folder and run

```
go build server.go
./server
```

### Docker:

Pull `benonymity/spotify-buddylist` and run it with an environment variable SP_DC containing your token. Bind port 10000 to wherever you want to serve the webpage.

Example Docker command:

```
docker run -d -p 10000:10000 --name spotify-buddylist benonymity/spotify-buddylist:latest
```

# Getting your sp_dc token

---

Login to the [web player](https://open.spotify.com/) in an incognito window and Inspect Element. Open the network tab and search `clienttoken`, then copy the token from the response. This is your one-year-expiry sp_dc token.

# Todo:

---

- [ ] Cache API results in backend
- [ ] Save history to DB and allow seeing past activity?
- [x] Add last listened time to frontend
- [x] Github Actions to build Docker image
