<template>
  <div className="friendactivity-content-bottom">
    <a
      :href="'https://open.spotify.com/track/' + friend.track.uri.split(':')[2]"
      target="_blank"
      ><div className="friendactivity-content-bottom-left">
        <img
          v-if="friend.user.imageUrl"
          :src="friend.user.imageUrl"
          alt="Profile"
        />
        <div v-else class="no-photo">
          <svg viewBox="0 0 80 90" width="18" height="20">
            <path
              d="M67.74 61.78l-14.5-8.334c-.735-.422-1.24-1.145-1.385-1.98-.145-.835.088-1.685.638-2.33l5.912-6.93c3.747-4.378 5.81-9.967 5.81-15.737v-2.256c0-6.668-2.792-13.108-7.658-17.67C51.622 1.92 45.17-.386 38.392.054c-12.677.82-22.607 11.772-22.607 24.934v1.483c0 5.77 2.063 11.36 5.81 15.736l5.912 6.933c.55.644.783 1.493.638 2.33-.143.834-.648 1.556-1.383 1.98l-14.494 8.33C4.7 66.077 0 74.15 0 82.844v6.76h3.333v-6.76c0-7.5 4.055-14.46 10.59-18.174l14.5-8.334c1.597-.918 2.692-2.487 3.007-4.302.315-1.815-.19-3.66-1.387-5.06l-5.913-6.936c-3.23-3.775-5.01-8.594-5.01-13.57v-1.484c0-11.41 8.562-20.9 19.488-21.608 5.85-.377 11.415 1.61 15.67 5.598 4.26 3.992 6.605 9.404 6.605 15.24v2.254c0 4.976-1.778 9.796-5.01 13.57l-5.915 6.935c-1.195 1.4-1.7 3.246-1.386 5.06.313 1.816 1.41 3.385 3.008 4.303l14.507 8.338c6.525 3.71 10.58 10.67 10.58 18.17v6.76H80v-6.76c0-8.695-4.7-16.768-12.26-21.063z"
              fill="#b3b3b3"
              fill-rule="evenodd"
            ></path>
          </svg>
        </div>
        <svg class="play" style="width: 24px; height: 24px" viewBox="0 0 24 24">
          <path fill="currentColor" d="M8,5.14V19.14L19,12.14L8,5.14Z" />
        </svg>
        <svg
          v-show="friend.timestamp > new Date().getTime() - 300000"
          class="presence"
          viewBox="0 0 20 20"
        >
          <circle fill="#000" cx="10" cy="10" r="5"></circle>
          <circle fill="#4275CA" cx="10" cy="10" r="4"></circle>
        </svg>
      </div>
    </a>
    <div className="friendactivity-content-bottom-right">
      <div className="friendactivity-content-bottom-right-user">
        <div style="width: 100%; display: flex; justify-content: space-between">
          <span
            ><a
              :href="
                'https://open.spotify.com/user/' + friend.user.uri.split(':')[2]
              "
              target="_blank"
              >{{ friend.user.name }}
            </a>
          </span>
          <div className="friendactivity-content-bottom-right-user-time">
            <span v-if="new Date().getTime() - friend.timestamp > 300000">{{
              difference
            }}</span>
            <svg
              v-else
              width="14"
              height="14"
              viewBox="0 0 14 14"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fill="none"
                stroke="#b3b3b3"
                stroke-width="2"
                d="M1 14V4m4 10V0m4 14v-4m4 4V7"
              />
            </svg>
          </div>
        </div>
      </div>
      <div className="friendactivity-content-bottom-right-song">
        <a
          :href="
            'https://open.spotify.com/track/' + friend.track.uri.split(':')[2]
          "
          target="_blank"
          >{{ friend.track.name }}</a
        >
        <span> • </span>
        <a
          :href="
            'https://open.spotify.com/artist/' +
            friend.track.artist.uri.split(':')[2]
          "
          target="_blank"
          >{{ friend.track.artist.name }}</a
        >
      </div>
      <div className="friendactivity-content-bottom-right-line">
        <div className="friendactivity-content-bottom-right-album">
          <div className="friendactivity-content-bottom-right-album-icon">
            <svg
              v-if="friend.track.context.uri.split(':')[1] == 'album'"
              width="14"
              height="14"
              viewBox="0 0 22 22"
              xmlns="http://www.w3.org/2000/svg"
            >
              <circle
                fill="none"
                stroke="#b3b3b3"
                stroke-width="2"
                cx="11"
                cy="11"
                r="10"
              />
              <circle
                fill="none"
                stroke="#b3b3b3"
                stroke-width="2"
                cx="11"
                cy="11"
                r="3"
              />
            </svg>
            <svg
              v-if="friend.track.context.uri.split(':')[1] == 'playlist'"
              style="width: 12px; height: 12px; margin-top: 2px"
              viewBox="0 0 32 32"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fill="none"
                stroke="#b3b3b3"
                stroke-width="3"
                d="M11 25a5 5 0 1 0-10 0 5 5 0 1 0 10 0V5l20-5v20a5 5 0 1 0-10 0 5 5 0 1 0 10 0"
              />
            </svg>
            <svg
              v-if="friend.track.context.uri.split(':')[1] == 'artist'"
              fill="#b3b3b3"
              width="12"
              height="13.3333"
              viewBox="0 0 18 20"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fill="none"
                stroke="#b3b3b3"
                stroke-width="1"
                d="M15.216 13.717 12 11.869a.492.492 0 0 1-.243-.348.496.496 0 0 1 .112-.41l1.311-1.537A5.498 5.498 0 0 0 14.5 6v-.5a5.524 5.524 0 0 0-1.739-4.014A5.46 5.46 0 0 0 8.636.011c-2.88.187-5.135 2.673-5.135 5.66V6c0 1.311.469 2.58 1.319 3.574l1.311 1.537a.49.49 0 0 1 .112.41.49.49 0 0 1-.244.348l-3.213 1.847A5.513 5.513 0 0 0 0 18.501V20h1v-1.499c0-1.616.874-3.116 2.283-3.917l3.215-1.848c.388-.223.654-.604.73-1.045a1.494 1.494 0 0 0-.337-1.229L5.579 8.925A4.505 4.505 0 0 1 4.499 6v-.329c0-2.461 1.845-4.509 4.2-4.662a4.468 4.468 0 0 1 3.377 1.206A4.461 4.461 0 0 1 13.5 5.5V6a4.5 4.5 0 0 1-1.08 2.925l-1.311 1.537a1.499 1.499 0 0 0 .394 2.274l3.218 1.849a4.513 4.513 0 0 1 2.28 3.916V20h1v-1.499a5.517 5.517 0 0 0-2.785-4.784Z"
              />
            </svg>
          </div>
          <div className="friendactivity-content-bottom-right-album-name">
            <a
              :href="
                'https://open.spotify.com/' +
                friend.track.context.uri.split(':')[1] +
                '/' +
                friend.track.context.uri.split(':')[2]
              "
              target="_blank"
              >{{ friend.track.context.name }}</a
            >
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    friend: {
      required: true,
      type: Object,
    },
  },
  computed: {
    difference() {
      let difference = Math.round(
        (new Date().getTime() - this.friend.timestamp) / 60000
      );
      if (difference < 60) {
        return String(difference) + " min";
      } else if (difference < 1440) {
        return String(Math.round(difference / 60)) + " hr";
      } else if (difference < 10080) {
        return String(Math.round(difference / 1440)) + " d";
      } else {
        return String(Math.round(difference / 10080)) + " w";
      }
    },
  },
};
</script>

<style>
.friendactivity-content-bottom {
  display: flex;
  align-items: center;
  margin-top: 6px;
  margin-bottom: 6px;
}
.friendactivity-content-bottom-left {
  width: 42px;
  height: 42px;
  margin-right: 12px;
  position: relative;
}
.friendactivity-content-bottom-left img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
}
.friendactivity-content-bottom-left:hover img {
  opacity: 50%;
}
.friendactivity-content-bottom-left .play {
  position: absolute;
  top: 9px;
  left: 9px;
  opacity: 0;
}
.friendactivity-content-bottom-left .presence {
  opacity: 100%;
  position: absolute;
  width: 24px;
  height: 24px;
  top: -7px;
  left: 25px;
  z-index: 3;
}
.friendactivity-content-bottom-left:hover .play {
  opacity: 100%;
}
.friendactivity-content-bottom-left .no-photo {
  background-color: #282828;
  width: 42px;
  height: 42px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.friendactivity-content-bottom-right {
  color: #9a9a9a;
}
.friendactivity-content-bottom-right-user {
  font-size: 1em;
  font-weight: 700;
  margin-bottom: 4px;
  display: flex;
  color: #fff;
  justify-content: space-between;
}
.friendactivity-content-bottom-right-user-time {
  color: #9a9a9a;
  font-size: 0.7em;
  font-weight: 500;
}
.friendactivity-content-bottom-right-user div {
  margin-right: 10px;
  justify-content: space-between;
}
.friendactivity-content-bottom-right-user-icon {
  margin-left: auto;
  margin-right: -20px;
  width: 16px;
  height: 16px;
  color: #9a9a9a;
}
.friendactivity-content-bottom-right-song,
.friendactivity-content-bottom-right-artist,
.friendactivity-content-bottom-right-album {
  font-size: 0.82em;
}
.friendactivity-content-bottom-right-song {
  margin-bottom: 3px;
  width: 260px;
  height: 20px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  text-align: left;
}
.friendactivity-content-bottom-right-artist {
  margin-bottom: 8px;
}
.friendactivity-content-bottom-right-album {
  padding-right: 24px;
  width: 255px;
  display: flex;
  justify-content: left;
}
.ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.friendactivity-content-bottom-right-album-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.friendactivity-content-bottom-right-album-icon {
  margin-right: 6px;
  width: 16px;
  height: 16px;
}
</style>
