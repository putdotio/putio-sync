var path = require('path');
var webpack = require('webpack');

module.exports = {
  entry: ['babel-polyfill', './src/app/App.jsx'],
  output: {
    path: __dirname + '/build/statics/js/',
    filename: 'index.js'
  },
  resolve: {
    extensions: ['', '.js', '.jsx']
  },
  module: {
    loaders: [
      {
        test: /.jsx?$/,
        loader: 'babel-loader',
        exclude: /(node_modules|bower_components)/,
        query: {
          presets: ['es2015', 'stage-0', 'react']
        }
      }
    ]
  }
};
