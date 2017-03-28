let express = require('express');
let app = express();
let spawn = require('child_process').spawn;
let ffmpeg = require('fluent-ffmpeg');
let fs = require('fs');
let path = require('path');
let cameraName = 'camera';
let port = 8080;

// Create the camera output directory if it doesn't already exist
if (fs.existsSync(cameraName) === false) {
  fs.mkdirSync(cameraName);
}

// Start the camera stream
// Have to do a smaller size otherwise FPS takes a massive hit :(
let cameraStream = spawn('raspivid', ['-o', '-', '-t', '0', '-n', '-h', '360', '-w', '640']);

// Convert the camera stream to hls
let conversion = new ffmpeg(cameraStream.stdout).noAudio().format('hls').inputOptions(['-re']).output(`${cameraName}/${cameraName}.m3u8`);

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

// Express middleware
app.use(function(req, res, next) {
  // Allow CORs
  res.header("Access-Control-Allow-Origin", "*");
  res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
  next();
});

// Essentially create a file server on the camera directory
app.get('/camera/:id', (req, res) => {
  res.set('Content-Type', 'application/x-mpegURL');
  let filepath = path.join(__dirname, cameraName, req.params.id);
  let readStream = fs.createReadStream(filepath);

  readStream.on('open', () => {
    readStream.pipe(res);
  });

  readStream.on('error', (err) => {
    res.status(400).json({'message': 'not found'});
  });
});

console.log(`STARTING CAMERA STREAM SERVER AT PORT ${port}`);
app.listen(port);
