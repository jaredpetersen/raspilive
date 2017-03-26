let express = require('express');
let app = express();
let path = require('path');
let spawn = require('child_process').spawn;
let ffmpeg = require('fluent-ffmpeg');
let streamPass = require('stream').PassThrough;
let port = 8080;

// Start the camera stream
// Have to do a smaller size otherwise FPS takes a massive hit :(
let cameraStream = spawn('raspivid', ['-o', '-', '-t', '0', '-n', '-h', '360', '-w', '640']);

// Convert the camera stream to hls
let conversion = new ffmpeg(cameraStream.stdout).noAudio().format('hls').inputOptions(['-re']);

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

// Output the conversion to a stream
// Can't just pipe directly off of conversionStream for every request, since ffmpeg only allows 1 stream output
let conversionStream = new streamPass;
conversion.pipe(conversionStream);

// Express middleware
app.use(function(req, res, next) {
  // Allow CORs
  res.header("Access-Control-Allow-Origin", "*");
  res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
  next();
});

app.get('/camera/1', (req, res) => {
  res.set('Content-Type', 'application/x-mpegURL');
  conversionStream.pipe(res);
});

console.log(`STARTING CAMERA STREAM SERVER AT PORT ${port}`);
app.listen(port);
