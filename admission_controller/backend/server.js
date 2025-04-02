// backend.js
const express = require('express');
const http = require('http');
const socketIo = require('socket.io');
const cors = require('cors');

const app = express();
app.use(cors());
app.use(express.json());

const server = http.createServer(app);
const io = socketIo(server, {
  cors: {
    origin: '*',
    methods: ['GET', 'POST'],
  },
});

// Updated verification endpoint
app.post('/api/verification', (req, res) => {
  const verificationData = req.body;
  console.log('Received verification event:', verificationData);

  // If the status contains "registered", treat it as success.
  if (
    verificationData.status &&
    verificationData.status.toLowerCase().includes('registered')
  ) {
    io.emit('registration', verificationData);
    res.status(200).json({ status: 'Verification event broadcasted' });
  } else {
    // For unauthorized cases, emit an admission event with failure details.
    io.emit('admission-event', verificationData);
    res
      .status(200)
      .json({ status: 'Unauthorized verification event broadcasted' });
  }
});

// Existing endpoint to receive other admission events (if needed)
app.post('/api/event', (req, res) => {
  const eventData = req.body;
  console.log('Received admission event:', eventData);
  io.emit('admission-event', eventData);
  res.status(200).json({ status: 'Event received' });
});

io.on('connection', (socket) => {
  console.log('New client connected, socket id:', socket.id);
  socket.on('disconnect', () => {
    console.log('Client disconnected');
  });
});

const PORT = 4000;
server.listen(PORT, () => {
  console.log(`Dashboard backend listening on port ${PORT}`);
});
