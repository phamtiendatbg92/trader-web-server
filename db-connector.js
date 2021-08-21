// mongoo db
var mongoose = require('mongoose');
var Schema = mongoose.Schema;
var {UrlMap} = require('./Config/Constants');

////////////////// AccountModel ///////////////////
const accountSchema = new Schema({
    name: String,
    passWord: String,
});

const AccountModel = mongoose.model('account', accountSchema);

////////////////// TutorialModel Model ///////////////////
const tutorialSchema = new Schema({
    title: String,
    content: String, // content
    tags: Array,
    url : String,
});

const TutorialModel = mongoose.model('Tutorials', tutorialSchema);

////////////////// HashTag Model ///////////////////
const hashTagSchema = new Schema({
    tags: Array
});

const HashTagModel = mongoose.model('HashTags', hashTagSchema);

mongoose.connect('mongodb://localhost:27017/trader', { useNewUrlParser: true, useUnifiedTopology: true });
const db = mongoose.connection;
db.on('error', console.error.bind(console, 'connection error:'));
db.once('open', function () {
    HashTagModel.findOne({}, function (err, doc) {
        console.log(doc);
      });
});




module.exports = {AccountModel , TutorialModel, HashTagModel};