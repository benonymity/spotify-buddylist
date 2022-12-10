<template>
  <div className="friendactivity">
    <div ref="unfocus" className="friendactivity-content">
      <div className="friendactivity-content-top">Spotify Friend Activity</div>
      <transition name="name">
        <div
          v-if="userActivity.length"
          className="friendactivity-content-top-wrap"
        >
          <span>{{ user.user.name }}</span>
          <button
            @click="endFocus()"
            className="friendactivity-content-top-close"
          >
            <svg
              width="24"
              height="24"
              viewBox="0 0 30 30"
              xmlns="http://www.w3.org/2000/svg"
              style="overflow: visible"
            >
              <path
                stroke="#b3b3b3"
                d="M12 0c6.623 0 12 5.377 12 12s-5.377 12-12 12-12-5.377-12-12 5.377-12 12-12zm0 1c6.071 0 11 4.929 11 11s-4.929 11-11 11-11-4.929-11-11 4.929-11 11-11zm0 10.293l5.293-5.293.707.707-5.293 5.293 5.293 5.293-.707.707-5.293-5.293-5.293 5.293-.707-.707 5.293-5.293-5.293-5.293.707-.707 5.293 5.293z"
              />
            </svg>
          </button>
          <div
            style="border-bottom: 2px solid #282828; margin-bottom: 14px"
          ></div>
        </div>
      </transition>
      <TransitionGroup name="fade">
        <div v-for="(friend, index) in listeningActivity" :key="index + 'div'">
          <button
            :key="index + 'button'"
            @click="startFocus(friend)"
            className="friendactivity-content-button"
          >
            <user-card :key="index" :friend="friend" />
          </button>
        </div>
      </TransitionGroup>
      <TransitionGroup name="move">
        <div
          v-for="(item, index) in userActivity"
          className="friendactivity-content-items"
          :key="item.track.uri"
        >
          <item-card :style="marginStyle(index)" :key="index" :item="item" />
          <!-- <div -->
          <!--   :style="lineStyle(index)" -->
          <!--   :key="index + 'line'" -->
          <!--   className="line" -->
          <!-- /> -->
        </div>
      </TransitionGroup>
      <!-- <div v-if="!transition && focus && userActivity.length == 0"> -->
      <!--   Whoops, looks like there is no historical listening data for this user! -->
      <!-- </div> -->
    </div>
  </div>
</template>

<script>
import axios from "axios";
import UserCard from "./components/UserCard.vue";
import ItemCard from "./components/ItemCard.vue";

export default {
  data() {
    return {
      listeningActivity: [],
      userActivity: [],
      user: "",
      search: "",
      searching: false,
      focus: false,
      transition: false,
    };
  },
  components: {
    UserCard,
    ItemCard,
  },
  methods: {
    loadStats() {
      let url = window.location + "api/latest";
      if (this.focus == true) {
        return;
      }
      axios
        .get(url)
        .then((response) => response.data)
        .then((data) => {
          this.listeningActivity = data.friends;
          this.listeningActivity.sort((a, b) => {
            return b.timestamp - a.timestamp;
          });
          this.loaded = true;
        });
    },
    loadUserStats() {
      let url = window.location + "api/" + this.user.user.uri.split(":")[2];
      axios
        .get(url)
        .then((response) => response.data)
        .then((data) => {
          this.userActivity = data.activity;
          this.userActivity.sort((a, b) => {
            return b.timestamp - a.timestamp;
          });
        });
    },
    autoReload() {
      if (!document.hidden) {
        if (!this.focus) {
          this.loadStats();
        } else {
          this.loadUserStats();
        }
      }
    },
    startFocus(friend) {
      this.transition = true;
      this.focus = false;
      this.listeningActivity = [];
      this.user = friend;
      setTimeout(() => {
        this.loadUserStats();
      }, 500);
      setTimeout(() => {
        this.focus = true;
        this.transition = false;
      }, 900);
    },
    endFocus() {
      this.transition = true;
      this.focus = true;
      this.userActivity = [];
      this.friendactivity = [];
      this.user = null;
      this.loadStats();
      setTimeout(() => {
        this.focus = false;
      }, 500);
      setTimeout(() => {
        this.transition = false;
      }, 900);
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
    lineStyle(i) {
      if (this.userActivity.length == 0) {
        return "";
      } else if (i == this.userActivity.length - 1) {
        return "";
      }
      let item = this.userActivity[i];
      let diff =
        item.timestamp - this.userActivity[i + 1].timestamp - item.duration;
      if (diff > 5000) {
        return "border: dashed 1.5px #383838; height: 90px; top: 55px";
      } else if (diff < 0) {
        return "border: solid 1.5px #A52A2A	; height: 70px; top: 55px";
      } else {
        return "border: solid 1.5px #383838; height: 70px; top: 55px";
      }
    },
    marginStyle(i) {
      if (this.userActivity.length == 0) {
        return "";
      } else if (i == this.userActivity.length - 1) {
        return "";
      }
      let item = this.userActivity[i];
      let diff =
        item.timestamp - this.userActivity[i + 1].timestamp - item.duration;
      if (diff > 5000) {
        return "margin-bottom: 60px";
      } else if (diff < 0) {
        return "margin-bottom: 20px";
      } else {
        return "margin-bottom: 30px";
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
  overflow-x: hidden;
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

.friendactivity-content-top-wrap {
  position: relative;
  padding-top: 20px;
}

.friendactivity-content-top-wrap span {
  text-align: center;
  position: relative;
  top: -8px;
  color: grey;
}

.friendactivity-content-top-close {
  position: relative;
  opacity: 50%;
  background: none;
  border: none;
  cursor: pointer;
}

.friendactivity-content-top-close:hover {
  opacity: 100%;
  transition: opacity 0.1s ease-in;
}

.friendactivity-content-button {
  background-color: #121212;
  border: none;
  cursor: pointer;
  margin-top: 7px;
  margin-bottom: 7px;
  border-radius: 0.6rem;
  border: none;
}

.friendactivity-content-button:hover {
  background-color: #282828;
}

.friendactivity-content-button-secondary {
  background-color: #121212;
  border: none;
  margin-top: 7px;
  margin-bottom: 7px;
  border-radius: 0.6rem;
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
.fade-leave-active,
.fade-enter-active {
  transition: all 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-leave-active {
  position: absolute;
}

.move-move,
.move.leave.active,
.move-enter-active {
  transition: all 0.5s ease;
}

.move-enter-from,
.move-leave-to {
  opacity: 0;
  transform: translateY(30px);
}

.move-leave-active {
  position: absolute;
}

.name-move,
.name.leave.active,
.name-enter-active {
  transition: all 0.5s ease;
}

.name-enter-from,
.name-leave-to {
  opacity: 0;
  transform: translateY(30px);
}

.name-leave-active {
  position: absolute;
}

::-webkit-scrollbar {
  width: 0em;
  height: 0em;
}

.friendactivity-content-items {
  position: relative;
  margin-top: 20px;
}

.line {
  position: absolute;
  left: 20px;
}
</style>
