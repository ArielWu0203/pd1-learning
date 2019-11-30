// require : import module
var http = require('http');

/*
@param : 
    request : including url, http header, data
    response : respond HTTP Request
*/
function onRequest(request,  response) {

    response.writeHead(200, {'Content-Type': 'text/plain'});
    response.write('Hello World');
    // end : inform client that messages have been sent.
    response.end();
    
}
http.createServer(onRequest).listen(3000);