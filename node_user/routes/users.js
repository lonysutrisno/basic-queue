require('dotenv').load();
var express = require('express');
var router = express.Router();
var redis = require('redis');
var client = redis.createClient(process.env.REDIS_PORT,process.env.REDIS_URL);
client.on('error', function(err){
  console.log('Something went wrong ', err)
});
var list_queue = process.env.LIST_QUEUE_NAME
var approved = process.env.LIST_APPROVED_NAME

/* GET users listing. */
router.get('/', function(req, res, next) {
  client.lrange(approved, 0, -1, function(error, result) {
    if (error) return res.json(error)
    if (result.length > 0){
      return res.json(result.map(function (params) {
        return JSON.parse(params);
      }))
    }else{
      return res.json(result)
    }
  });

});

router.post('/insert-queue', function(req, res, next) { 
  
  client.rpush(list_queue, JSON.stringify(req.body), function(error, result) {
    if (error) return res.json(error)

    return res.json({"number_of_queue":result})

  });
});

module.exports = router;
