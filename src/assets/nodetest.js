const http = require('http');
const fs = require('fs');

http.createServer( (req, res) => {
  res.end('Hello from NodeJS!');
}).listen(3000, () => console.log('Server listen in 3000') );
