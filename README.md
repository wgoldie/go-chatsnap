Go Chatsnap
========
See a running version at [chat.willgoldie.me](http://chat.willgoldie.me/).

Chatsnap is an application that creates emoji from ngram elements (currently words) and uses them for chat purposes.
There is an old version of it somewhere that [arjunyg](https://github.com/arjunyg), [talbenari1](https://github.com/talbenari1), and I made for a hackathon in NodeJS (which I can't look at without vomiting). Some of that has been ported over to this version, but most of the code is original.

go-chatsnap runs pretty well on heroku. You need Yahoo BOSS api keys, Pubnub API keys, and a Redis addon set as environment vars as used in server.go.

Feel free to fork or pull request.
