const path = require('path');
const darkTheme = require('@ant-design/dark-theme');

module.exports = {
  outputDir: path.resolve(__dirname, '../docs'),
  configureWebpack: (config) => {
    config.devtool = 'source-map'; // eslint-disable-line
  },
  css: {
    loaderOptions: {
      less: {
        modifyVars: darkTheme,
        javascriptEnabled: true,
      },
    },
  },
};
