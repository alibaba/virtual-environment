'use strict';

var express = require('express');
const request = require('request');

const envVar = process.argv.splice(2);
const envMark = envVar[0] ? envVar[0] : 'dev';
const url = envVar[1] ? envVar[1] : 'none';

console.log('envMark:', envMark);
console.log('url:' + url);

var app = express();
app.get('/demo', function (req, res) {

    const envMarkKey = 'ali-env-mark';

    const reqEnvMark = req.headers[envMarkKey] ? req.headers[envMarkKey] : 'empty';

    const headers = {};
    if (req.headers[envMarkKey]) {
        headers[envMarkKey] = reqEnvMark;
    }

    var requestText = '';
    if (url && url !== 'none' && url.length > 0) {
        request({
            url: url,
            method: 'GET',
            headers: headers
        }, function (error, response, body) {
            if (!error && response.statusCode == 200) {
                requestText = body
            } else {
                requestText = 'call ' + url + ' failed'
            }
            res.send(requestText + '[node @ ' + envMark + '] <-' + reqEnvMark + '\n');
        });
    } else {
        res.send('[node @ ' + envMark + '] <-' + reqEnvMark + '\n');
    }
});
console.log("listening to 8080")
app.listen(8080);
