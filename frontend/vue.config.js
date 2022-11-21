const { defineConfig } = require("@vue/cli-service");
module.exports = {
  transpileDependencies: true,
  devServer: {
    allowedHosts: "all",
  },
  css: {
    extract: false,
  },
  configureWebpack: {
    optimization: {
      splitChunks: false,
    },
  },
};
