Go Chatsnap
========
See a running version at [chat.willgoldie.me](http://chat.willgoldie.me/).

Chatsnap is an application that creates emoji from ngram elements (currently words) and uses them for chat purposes.
There is an old version of it somewhere that [arjunyg](https://github.com/arjunyg), [talbenari1](https://github.com/talbenari1), and I made for a hackathon in NodeJS (which I can't look at without vomiting). Some of that has been ported over to this version, but most of the code is original.

go-chatsnap runs pretty well on heroku. 
You need Yahoo BOSS api keys, Pubnub API keys, and a Redis addon set as the following respective environment variables:


- BING_APP_ID (note: this is a misnomer, it should actually be an azure account key. naming will be fixed in future versions)
- PUBNUB_PUBLISH_KEY
- PUBNUB_SUBSCRIBE_KEY
- PUBNUB_SECRET_KEY
- REDIS_URL
- SMS_TARGET (optional space-separated list of phone numbers in the format 555-555-5555 to send admin logs to - you will probably have to modify the app to use non ca numbers)
- PORT


Feel free to fork or pull request.
