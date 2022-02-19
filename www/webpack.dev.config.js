const path = require('path');
const webpack = require('webpack');
const TerserPlugin = require('terser-webpack-plugin');

module.exports = {
  mode: 'development',
  devtool: 'eval-source-map',
  entry: {
    app: {
      import: './app/App.js',
    },
    uapp: {
      import: './uapp/App.js',
    },
  },
  watchOptions: {
    aggregateTimeout: 100,
    ignored: [
      path.resolve(__dirname, 'node_modules'),
    ],
  },
  output: {
    path: path.resolve(__dirname, 'dist-dev', 'static'),
    publicPath: '',
    filename: '[name].js',
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        enforce: 'pre',
        use: ['source-map-loader'],
      },
    ],
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env': JSON.stringify({}),
    }),
  ],
};
