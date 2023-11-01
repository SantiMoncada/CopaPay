const path = require("path");

module.exports = {
  entry: "./js/index.js",
  output: {
    path: path.resolve(__dirname, "assets/scripts"),
    filename: "index.js",
  },
  mode: "production",
};
