#!/usr/bin/env node

'use strict';

const program = require('commander');
const os = require('os');
const info = require('./package.json');
const server = require('./lib/server');

// Coercion function for number range
const range = (min, max, value, def) => {
  if (value < min || value > max) return def;
  return Number(value);
};

process.title = 'raspi-live';

program
  .name(info.name)
  .description(info.description)
  .version(info.version, '-v, --version');

program
  .command('start')
  .description('start streaming video from the raspberry pi camera module')
  .option('-d, --directory <directory>', 'streaming video file hosting location', `${os.homedir()}/camera`)
  .option('-f, --format <format>', 'video streaming format [hls, dash]', /^(hls|dash)$/i, 'hls')
  .option('-w, --width <width>', 'video resolution width', Number, 1280)
  .option('-h, --height <height>', 'video resolution height', Number, 720)
  .option('-r, --framerate <fps>', 'video frames per second', Number, 25)
  .option('-x, --horizontal-flip', 'flip the camera horizontally')
  .option('-y, --vertical-flip', 'flip the camera vertically')
  .option('-c, --compression-level <compression-level>', 'compression level [0-9]', range.bind(null, 0, 9), 9)
  .option('-t, --time <time>', 'duration of streaming files', Number, 2)
  .option('-l, --list-size <list-size>', 'number of streaming files in the playlist', Number, 10)
  .option('-s, --storage-size <storage-size>', 'number of streaming files for storage purposes', Number, 10)
  .option('-p, --port <port>', 'port number the server runs on', Number, 8080)
  .option('-S, --secure', 'run with credentials for HTTPS')
  .option('-C, --certificatePath <file>', 'path to SSL certificate', '')
  .option('-K, --keyPath <file>', 'path to private key', '')
  .action(({ directory, format, width, height, framerate, horizontalFlip = false, verticalFlip = false, compressionLevel, time, listSize, storageSize, port, secure = false, certificatePath, keyPath}) => {
    console.log('configuration:', directory, format, width, height, framerate, horizontalFlip, verticalFlip, compressionLevel, time, listSize, storageSize, port, secure, certificatePath, keyPath);
    server(directory, format, width, height, framerate, horizontalFlip, verticalFlip, compressionLevel, time, listSize, storageSize, port, secure, certificatePath, keyPath);
  });

program.helpOption('--help');
program.parse(process.argv);
