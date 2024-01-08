# Spotify Friend Activity

This program serves a webpage that imitates Spotify's official friend activity list for the account token you supply.

It serves a web client that calls a backend proxy on your server which talks to Spotify's unofficial activity API using a privledged access token and displays it in a webpage with similar styling to Spotify.

Example:
![Example](https://github.com/benonymity/spotify-buddylist/assets/62854267/a071aa7b-0f34-49d4-b0db-43847497adf8)


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

Login to the [web player](https://open.spotify.com/) and Inspect Element. Open the network tab and then paste [https://open.spotify.com/get_access_token?reason=transport&productType=web_player](https://open.spotify.com/get_access_token?reason=transport&productType=web_player) into the normal URL bar, press enter to load the tab, then open the first item in the network page:

![Screenshot 2024-01-08 at 4 08 21 PM|300](https://github.com/benonymity/spotify-buddylist/assets/62854267/fc3171d1-0eac-4eef-b94a-dfc76e40f3f1)

Scroll down to cookie in the Request Headers, then hunt for `sp_dc=...` and copy what come after until the semicolon. This is your one-year-expiry sp_dc token.

# Todo:

- [x] Cache API results in backend to avoid DDOS
- [x] Save history to DB
- [x] Create UI for viewing past activity
- [ ] Allow following of new users from webpage
