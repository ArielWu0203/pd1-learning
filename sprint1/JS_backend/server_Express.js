var express = require('express');
var app = express();

app.get('/', function(req ,res) {
    res.writeHead(200, {'Content-Type': 'text/plain'});
    res.write('Hello World');
    res.end();
})

app.get('/api/problem', function(req, res) {
    res.send('Problem');
})

app.listen(3000, function() {
    console.log('App is listning on port 3000~');
})