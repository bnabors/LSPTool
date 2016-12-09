/* Copyright 2016 Juniper Networks, Inc. All rights reserved.
 * Licensed under the Juniper Networks Script Software License (the "License").
 * You may not use this script file except in compliance with the License, which is located at
 * http://www.juniper.net/support/legal/scriptlicense/
 * Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 */

let CopyWebpackPlugin = require('copy-webpack-plugin');
let ExtractTextPlugin = require('extract-text-webpack-plugin');
let webpack = require('webpack');
require("babel-polyfill");

let copyrights = new webpack.BannerPlugin('Copyright 2016 Juniper Networks, Inc. All rights reserved.\r\n' +
    'Licensed under the Juniper Networks Script Software License (the "License").\r\n' +
    'You may not use this script file except in compliance with the License, which is located at\r\n' +
    'http://www.juniper.net/support/legal/scriptlicense/\r\n' +
    'Unless required by applicable law or otherwise agreed to in writing by the parties, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.');

module.exports = {
    entry: ['babel-polyfill', './src/index.js'],
    output: {
        path: '../Server/site/',
        filename: '/src/app.bundle.js',
    },
    module: {
        loaders: [{
            test: /\.jsx?$/,
            exclude: /node_modules/,
            loader: 'babel-loader',
            query: {
                presets: ['react', 'es2015', 'stage-0'],
                plugins: ['react-html-attrs', 'transform-class-properties', 'transform-decorators-legacy'],
            }
        },
            {
                test: /\.css$/,
                loader: ExtractTextPlugin.extract('style-loader', 'css-loader')
            }]
    },
    plugins: [
        copyrights,
        new webpack.DefinePlugin({
            'process.env': {
                'NODE_ENV': JSON.stringify('production')
            }
        }),
        new CopyWebpackPlugin([{from: './src/index.html'}]),
        new CopyWebpackPlugin([{from: './src/sprite.png', to: './src/sprite.png'}]),
        new ExtractTextPlugin('/src/app.bundle.css')
    ]
};
