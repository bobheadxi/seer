const path = require('path');

module.exports = {
  outputDir: path.resolve(__dirname, '../docs'),
  configureWebpack: (config) => {
    config.devtool = 'source-map'; // eslint-disable-line
  },
};
