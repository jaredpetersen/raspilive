const express = require('express');
const app = express();
const compression = require('compression');
const cors = require('cors');
const fs = require('fs');
const spawn = require('child_process').spawn;
const ffmpeg = require('fluent-ffmpeg');

module.exports = (directory, format, port) => {
  // Create the camera output directory if it doesn't already exist
  // Sync, because this is only run once at startup and everything depends on it
  if (fs.existsSync(directory) === false) fs.mkdirSync(directory);

  // Start the camera stream
  const raspividOptions = ['-o', '-', '-t', '0', '-vf', '-w', '1280', '-h', '720', '-fps', '25'];
  const cameraStream = spawn('raspivid', raspividOptions);

  // Setup up a special shutdown function that's called when encountering an error
  // so that we always shut down the camera stream properly
  const kill = (err) => {
    cameraStream.kill();
    throw err;
  };

  // Set up camera stream conversion
  let conversionStream = ffmpeg(cameraStream.stdout)
    .noAudio();

  if (format === 'hls') {
    conversionStream
      .videoCodec('copy')
      .format('hls')
      .inputOptions(['-re'])
      .outputOptions(['-hls_flags delete_segments'])
      .output(`${directory}/livestream.m3u8`);
  }
  else if (format === 'mpeg-dash') {
    conversionStream
      .videoCodec('copy')
      .format('dash')
      .inputOptions(['-re'])
      .outputOptions(['-seg_duration', '2', '-window_size', '10', '-extra_window_size', '10'])
      .output(`${directory}/livestream.mpd`);
  }
  else {
    kill(Error('unsupported format'));
  }

  // Start stream processing
  conversionStream
    .on('error', (err, stdout, stderr) => kill(err))
    .on('start', (commandLine) => console.log('started video processing: ' + commandLine))
    .on('stderr', (stderrLine) => console.log('conversion: ' + stderrLine))
    .run();

  // Endpoint the streaming files will be available on
  const endpoint = '/camera';

  // Setup express server
  app.use(cors());
  app.use(compression({ level: 9 }));
  app.use(endpoint, express.static(directory));
  app.listen(port);

  console.log('camera stream server started');
  console.log(`format: ${format}`);
  console.log(`directory: ${directory}`);
  console.log(`endpoint: ${endpoint}`) 
  console.log(`port: ${port}`);
};
