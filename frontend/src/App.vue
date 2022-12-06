<template>
  <div className="friendactivity">
    <div ref="unfocus" className="friendactivity-content">
      <div className="friendactivity-content-top">Spotify Friend Activity</div>
      <div v-if="listeningActivity.length > 0">
        <TransitionGroup name="fade">
          <div v-for="friend in listeningActivity" :key="friend.user.uid">
            <button
              v-if="!focus"
              @click="
                refresh = false;
                focus = true;
                listeningActivity = [friend];
              "
              className="friendactivity-content-button"
            >
              <user-card :friend="friend" />
            </button>
            <user-card v-if="focus" :friend="friend" />
          </div>
        </TransitionGroup>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import UserCard from "./components/UserCard.vue";

export default {
  data() {
    return {
      listeningActivity: [],
      refresh: true,
      search: "",
      searching: false,
      focus: false,
    };
  },
  components: {
    UserCard,
  },
  methods: {
    loadStats() {
      // let url = window.location + "api/0";
      if (this.refresh == false) {
        return;
      }
      let url = "http://192.168.0.30:10000/api/0";
      axios
        .get(url)
        .then((response) => response.data)
        .then((data) => {
          this.listeningActivity = data.friends;
          this.listeningActivity
            .sort((a, b) => {
              return a.timestamp - b.timestamp;
            })
            .reverse();
          this.loaded = true;
        });
    },
    autoReload() {
      if (!document.hidden) {
        this.loadStats();
      }
    },
    focusSearchbar() {
      this.$refs.searchbox.focus();
    },
    unfocusSearchbar() {
      this.$refs.searchbox.blur();
      this.search = "";
    },
    trySearch() {
      if (this.search == "") {
        this.searching = false;
      } else {
        this.searching = true;
        // Search logic here
      }
    },
  },
  mounted() {
    document.title = "Spotify Friend Activity";
    this.loadStats();
    addEventListener("visibilitychange", this.autoReload);
    addEventListener("keydown", (e) => {
      let key = e.keyCode;
      if (key == 27) {
        this.unfocusSearchbar();
      } else if (key > 65 && key < 90) {
        this.focusSearchbar();
      }
    });
    setInterval(
      function () {
        this.autoReload();
      }.bind(this),
      5000
    );
  },
};
</script>

<style>
html {
  background-color: #121212;
  font-family: spotify-font;
}

@font-face {
  font-family: spotify-font;
  src: url(https://encore.scdn.co/fonts/CircularSp-Book-4eaffdf96f4c6f984686e93d5d9cb325.woff2);
}

a {
  color: #9a9a9a;
  text-decoration: none;
}

a:hover {
  color: #fff;
  text-decoration: underline;
}

.friendactivity {
  color: #fff;
  margin-left: auto;
  margin-right: auto;
  display: flex;
  padding-top: 20px;
  text-align: center;
}

.friendactivity-content {
  margin-left: auto;
  margin-right: auto;
  width: 320px;
}

.friendactivity-content-top {
  font-size: 1.16em;
  font-weight: bold;
  padding-bottom: 12px;
}

.friendactivity-content-button {
  background-color: #121212;
  border: none;
  cursor: pointer;
  margin-top: 8px;
  border-radius: 1rem;
}

.friendactivity-content-button:hover {
  background-color: #282828;
}

.friendactivity-content-line {
  border-bottom: 0.5px solid #282828;
}

.friendactivity-content-input {
  padding-bottom: 4px;
  display: flex;
  position: relative;
}

.friendactivity-content-input input {
  width: 100%;
  padding-top: 12 px;
  padding-bottom: 6px;
  border: none;
  height: 19px;
  font-size: 18pt;
  font-family: spotify-font;
  border-bottom: 0.5px solid #282828;
  background-color: #121212;
  color: #fff;
  text-align: center;
  background: linear-gradient(#383838 0 0) bottom left/
    var(--underline-width, 0%) 0.1em no-repeat;
  transition: background-size 0.5s;
}

.friendactivity-content-input input:focus {
  outline: none;
  --underline-width: 100%;
}

.friendactivity-content-input svg {
  padding-bottom: 4px;
  position: absolute;
  right: 6px;
  bottom: 6px;
}

.fade-move,
.fade-enter-active,
.fade-leave-active {
  transition: all 0.2s cubic-bezier(0.55, 0, 0.1, 1);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: scaleY(0.01) translate(30px, 0);
}

.fade-leave-active {
  position: absolute;
}
</style>
