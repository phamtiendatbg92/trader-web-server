var express = require('express');
var path = require('path');
var cookieParser = require('cookie-parser');
var logger = require('morgan');
var cors = require('cors');
var indexRouter = require('./routes/index');

var app = express();
// Add headers
app.use(cors());
app.use(logger('dev'));
// app.use(express.json());
// app.use(express.urlencoded({ extended: false }));
app.use(express.json({limit: '50mb'}));
app.use(express.urlencoded({extended: false}));

app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

app.use('/api/v1', indexRouter);

module.exports = app;
