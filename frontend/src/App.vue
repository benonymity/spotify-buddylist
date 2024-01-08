<template>
  <div className="friendactivity">
    <div ref="unfocus" className="friendactivity-content">
      <div className="friendactivity-content-top">Spotify Friend Activity</div>
      <transition name="name">
        <div
          v-if="userActivity.length"
          className="friendactivity-content-top-wrap"
        >
          <span style="left: 10px">{{ user.user.name }}</span>
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
      <TransitionGroup name="fade">
        <div v-for="(item, index) in userActivity" :key="index + 'div'">
          <item-card :style="marginStyle(index)" :key="index" :item="item" />
          <div className="friendactivity-content-items">
            <div
              :style="lineStyle(index)"
              :key="index + 'line'"
              className="line"
            />
          </div>
        </div>
        <infinite-loading
          v-if="userActivity.length"
          @infinite="infiniteHandler"
          spinner="spiral"
        >
        </infinite-loading>
      </TransitionGroup>
      <div v-if="focus && userActivity.length == 0">
        Whoops, looks like there is no historical listening data for this user!
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import UserCard from "./components/UserCard.vue";
import ItemCard from "./components/ItemCard.vue";
import InfiniteLoading from "vue-infinite-loading";

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
      page: 1,
      url: window.location,
    };
  },
  components: {
    UserCard,
    ItemCard,
    InfiniteLoading,
  },
  methods: {
    loadStats() {
      if (this.focus == true) {
        return;
      }
      let url = this.url + "api/latest";
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
      let url = this.url + "api/" + this.user.user.uri.split(":")[2] + "/0";
      axios
        .get(url)
        .then((response) => response.data)
        .then((data) => {
          this.userActivity = data.activity;
        });
    },
    infiniteHandler($state) {
      let url =
        this.url + "api/" + this.user.user.uri.split(":")[2] + "/" + this.page;
      axios
        .get(url)
        .then((response) => response.data)
        .then((data) => {
          if (data.activity.length) {
            this.userActivity = this.userActivity.concat(data.activity);
            this.page += 1;
            setTimeout(() => {
              $state.loaded();
            }, 200);
          } else {
            $state.complete();
          }
        });
    },
    autoReload() {
      if (!document.hidden) {
        if (!this.focus) {
          this.loadStats();
        }
      }
    },
    startFocus(friend) {
      this.focus = false;
      this.listeningActivity = [];
      this.userActivity = [];
      this.user = friend;
      this.loadUserStats();
      this.page = 1;
      setTimeout(() => {
        this.focus = true;
        this.listeningActivity = [];
      }, 200);
    },
    endFocus() {
      this.focus = true;
      this.userActivity.slice(0, 15);
      setTimeout(() => {
        this.userActivity = [];
      }, 300);
      this.listeningActivity = [];
      this.user = null;
      this.focus = false;
      setTimeout(() => {
        this.loadStats();
      }, 700);
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
        return "border: dashed 1.4px #383838; height: 90px; top: -72px";
      } else if (diff < 0) {
        return "border: solid 1.4px #A52A2A; height: 70px; top: -32.5px";
      } else {
        return "border: solid 1.4px #383838; height: 70px; top: -42.5px";
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
    if (String(this.url).includes("127.0.0.1") | String(this.url).includes("localhost")) {
      this.url = "http://localhost:10000/"
    }
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
  width: 320px;
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
  left: 10px;
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
  width: 330px;
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
.name-leave-active,
.name-enter-active {
  transition: all 0.5s ease;
}

.name-enter-from,
.name-leave-to {
  opacity: 0;
  transform: translateY(-30px);
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
