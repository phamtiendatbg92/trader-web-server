var express = require('express');
var router = express.Router();
var { Constants } = require('../Config/Constants');
var { TutorialModel, HashTagModel } = require('../db-connector.js');
const fs = require('fs');

function onSaveError() {
  console.log("Save file error");
}

function fromTitleToUrl(title){
  return title.replace(/ /g, '-').toLowerCase()
}

function CreateNewPost(title, content, tags, url) {
  var news = new TutorialModel({
    title: title,
    content: content,
    tags: tags,
    url: url
  });
  news.save(function (err) {
    if(err)
    {
      console.log("CreateNewPost: " + err);
    }
  })
}

function handleBase64Image(content) {
  var match = content.match(/data:image\/([a-zA-Z]*);base64,([^\"]*)/g);
  if (match != null) {
    for (let i = 0; i < match.length; i++) {
      const mimeType = match[i].substring(match[i].indexOf("image") + 6, match[i].indexOf("base64") - 1);
      const imgData = match[i].substring(match[i].indexOf("base64") + 7, match[i].length);
      // Create new buffer contain data
      let imageBuffer = Buffer.from(imgData, 'base64');
      let fileName = Date.now() + "." + mimeType;
      let fileNameWithKey = Constants.SERVER_PUBLIC_KEY + Date.now() + "." + mimeType;
      // replace data to url to save to database
      content = content.replace(match[i], fileName);
      try {
        fs.writeFileSync("./public/images/" + fileName, imageBuffer, 'utf8');
      } catch (e) {
        console.log("=========Save image error: ========" + e);
        //res.sendStatus(500);
      }
    }
  }
  return content;
}

function handleHashTag(reqTags){
  HashTagModel.findOne({}, function (err, doc) {
    if (err) {
    }
    else
    {
      if(doc == null){
        var news = new HashTagModel({
          tags: reqTags,
        });
        news.save(function (err) {
          if (err) onSaveError();
        })
      }
      else
      {
        let isTagChange = false;
        // Check if need to add new hashtag to DB
        reqTags.forEach(hashTag => {
          if (!doc.tags.includes(hashTag)) {
            doc.tags.push(hashTag);
            isTagChange = true;
          }
        });
        if (isTagChange) {
          doc.save();
        }
      }
    }
  });
}

async function validateUrl(title){
  let url = fromTitleToUrl(title);
  const doc =  await TutorialModel.find({url : url}).exec();
  if(doc != "")
  {
    url = "";
  }
  return url;
}
/* Upload new post */
router.post('/upload-new-post', async function (req, res, next) {
  try {
    const url = await validateUrl(req.body.title);
    if(url != "")
    {
      let content = req.body.content;
      content = handleBase64Image(content);
      CreateNewPost(req.body.title, content, req.body.tags, url);
      handleHashTag(req.body.tags);
      res.sendStatus(200);
    }
    else
    {
      res.sendStatus(500);
    }
  }
  catch (e) {
    console.log("=========upload-new-post========" + e);
    res.sendStatus(500);
  }
});
router.get('/', function (req, res, next) {
  res.send('hello heroku !!!')
});

router.get('/get-list-tutorials', function (req, res, next) {
  const filter = {};
  TutorialModel.find(filter, function (err, doc) {
    res.json(doc);
  });
});

router.get('/detail-tutorial', function (req, res, next) {
  const filter = { url: req.query.tutorialUrl };
  TutorialModel.findOne(filter, function (err, doc) {
    if (err) {
      res.sendStatus(500);
    }
    res.json(doc);
  });
});

router.get('/get-hashtag', function (req, res, next) {
  HashTagModel.findOne({}, function (err, doc) {
    if (err) {
      res.sendStatus(500);
    }
    res.json(doc);
  });
});



module.exports = router;
