let express = require('express');
let app = express();
let spawn = require('child_process').spawn;
let ffmpeg = require('fluent-ffmpeg');
let fs = require('fs');
let path = require('path');
let cameraName = 'camera';
let port = 8080;

// Create the camera output directory if it doesn't already exist
// Directory contains all of the streaming video files
if (fs.existsSync(cameraName) === false) {
  fs.mkdirSync(cameraName);
}

// Start the camera stream
// Have to do a smaller size otherwise FPS takes a massive hit :(
let cameraStream = spawn('raspivid', ['-o', '-', '-t', '0', '-n', '-h', '360', '-w', '640']);

// Convert the camera stream to hls
let conversion = new ffmpeg(cameraStream.stdout).noAudio().format('hls').inputOptions('-re').outputOptions(['-hls_wrap 20', '-hls_key_info_file video.keyinfo']).output(`${cameraName}/livestream.m3u8`);

// Set up listeners
conversion.on('error', function(err, stdout, stderr) {
  console.log('Cannot process video: ' + err.message);
});

conversion.on('start', function(commandLine) {
  console.log('Spawned Ffmpeg with command: ' + commandLine);
});

conversion.on('stderr', function(stderrLine) {
  console.log('Stderr output: ' + stderrLine);
});

// Start the conversion
conversion.run();

// Allows CORS
let setHeaders = (res, path) => {
  res.header("Access-Control-Allow-Origin", "*");
  res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
};

// Set up a fileserver for the streaming video files
app.use(`/${cameraName}`, express.static(cameraName, {'setHeaders': setHeaders}));

console.log(`STARTING CAMERA STREAM SERVER AT PORT ${port}`);
app.listen(port);
