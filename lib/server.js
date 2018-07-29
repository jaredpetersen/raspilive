const express = require('express');
const app = express();
const compression = require('compression');
const cors = require('cors');
const spawn = require('child_process').spawn;
const ffmpeg = require('fluent-ffmpeg');

module.exports = (directory, format = 'hls', port = 8080) => {
  // Start the camera stream
  const raspividOptions = ['-o', '-', '-t', '0', '-vf', '-w', '1280', '-h', '720', '-fps', '25'];
  const cameraStream = spawn('raspivid', raspividOptions);

  if (format === 'hls') {
    ffmpeg(cameraStream.stdout)
      .noAudio()
      .videoCodec('copy')
      .format('hls')
      .inputOptions(['-re'])
      .outputOptions(['-hls_flags delete_segments'])
      .output(`${directory}/livestream.m3u8`);
  }
  else if (format === 'dash') {
    ffmpeg(cameraStream.stdout)
      .noAudio()
      .videoCodec('copy')
      .format('dash')
      .inputOptions(['-re'])
      .outputOptions(['-seg_duration', '2', '-window_size', '10', '-extra_window_size', '10'])
      .output(`${directory}/livestream.mpd`);
  }
  else {
    throw new Error('unsupported format');
  }

  // Start stream processing
  ffmpeg
    .on('error', (err, stdout, stderr) => {
      throw err;
    })
    .on('start', (commandLine) => console.log('started video processing: ' + commandLine))
    .on('stderr', (stderrLine) => console.log('conversion: ' + stderrLine))
    .run();

  // Setup express server
  app.use(cors());
  app.use(compression({ level: 9 }));
  app.use(`/${directory}`, express.static(directory));
  app.listen(port);

  console.log(`camera stream server started: ${port}`);
};
