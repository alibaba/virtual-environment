'use strict';

var express = require('express');
const request = require('request');
const {FORMAT_HTTP_HEADERS} = require('opentracing');
const {initTracer} = require('jaeger-client');

var envVar = process.argv.splice(2)
var envMark = envVar[0] ? envVar[0] : 'dev';
var url = envVar[1] ? envVar[1] : 'none';

console.log('envMark:', envMark);
console.log('url:' + url);

const tracer = initTracer({serviceName: 'demo'}, {});

var app = express();
app.get('/demo', function (req, res) {
  const parent = tracer.extract(FORMAT_HTTP_HEADERS, req.headers);
  const span = parent ? tracer.startSpan(`${req.host}`, {childOf: parent}) : tracer.startSpan('http_request');

  var reqEnvMark = span.getBaggageItem('ali-env-mark');

  const headers = {};
  tracer.inject(span, FORMAT_HTTP_HEADERS, headers);

  var requestText = '';
  if (url && url !== 'none' && url.length > 0) {
    request({
        url: url,
        method: 'GET',
        headers: headers
    }, function(error, response, body) {
        if (!error && response.statusCode == 200) {
            requestText = body
        }else{
            requestText = 'call ' + url + ' failed'
        }
        res.send(requestText + '\n' + '[node][request env mark is ' + (reqEnvMark?reqEnvMark:'empty') + '][my env mark is ' + envMark + ']');
    }); 
  }else{
    res.send('[node][request env mark is ' + (reqEnvMark?reqEnvMark:'empty') + '][my env mark is ' + envMark + ']');
  }
});
console.log("listening to 8080")
app.listen(8080);
