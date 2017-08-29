let express = require('express');
let app = express();
let cors = require('cors');
let spawn = require('child_process').spawn;
let ffmpeg = require('fluent-ffmpeg');
let fs = require('fs');
let crypto = require('crypto');

// Config information
const config = require('./config.json');
const port = config.port;
const cameraName = config.cameraName;
const hlsListSize = config.hlsListSize;
const hlsEncryptionEnabled = config.hlsEncryption.enabled;

// Camera stream options
const raspividOptions = ['-o', '-', '-t', '0', '-vf', '-w', '1280', '-h', '720', '-fps', '30']; 
const ffmpegInputOptions = ['-re'];
//const ffmpegOutputOptions = ['-vcodec copy', '-g 50', `-hls_wrap ${hlsWrapLength}`];
const ffmpegOutputOptions = ['-vcodec copy', '-g 50', `-hls_list_size ${hlsListSize}`, '-hls_flags delete_segments'];

// Create the camera output directory if it doesn't already exist
// Directory contains all of the streaming video files
// We don't want the async version since this only is run once at startup and the directory needs to be created
// before we can really do anything else
if (fs.existsSync(cameraName) === false) {
  fs.mkdirSync(cameraName);
}

// Encrypt HLS stream?
if (hlsEncryptionEnabled) {
  // Config information for hls encryption
  const publicBaseURL = config.hlsEncryption.publicBaseURL;
  const keyFileName = config.hlsEncryption.keyFileName;
  const keyInfoFileName = config.hlsEncryption.keyInfoFileName;

  // Setup encryption
  let keyFileContents = crypto.randomBytes(16);
  let initializationVector = crypto.randomBytes(16).toString('hex');
  let keyInfoFileContents = `${publicBaseURL}/${keyFileName}\n./${cameraName}/${keyFileName}\n${initializationVector}`;

  // Populate the encryption files, overwrite them if necessary
  fs.writeFileSync(`./${cameraName}/${keyFileName}`, keyFileContents);
  fs.writeFileSync(keyInfoFileName, keyInfoFileContents);

  // Add an option to the output stream to include the key info file in the livestream playlist
  ffmpegOutputOptions.push(`-hls_key_info_file ${keyInfoFileName}`);
}

// Start the camera stream
let cameraStream = spawn('raspivid', raspividOptions);

// Convert the camera stream to hls
let conversion = new ffmpeg(cameraStream.stdout).noAudio().format('hls').inputOptions(ffmpegInputOptions).outputOptions(ffmpegOutputOptions).output(`${cameraName}/livestream.m3u8`);

// Set up stream conversion listeners
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
app.use(cors());

// Set up a fileserver for the streaming video files
app.use(`/${cameraName}`, express.static(cameraName));

app.listen(port);
console.log(`STARTING CAMERA STREAM SERVER AT PORT ${port}`);
