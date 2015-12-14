(function (window, document, undefined) {
    var Chat = {
        messages: [],
        channel: 'chat',
        handles: []
    };
    
    var messagePane = document.getElementById('messages');
    var handleInput = document.getElementById('handle');
    
    var Message = function (data) {
        var scrollDown = messagePane.scrollTop === (messagePane.scrollHeight - messagePane.offsetHeight);
        var messageBox = document.createElement('div');
        var message = JSON.parse(data);
        
        messageBox.classList.add('message', message["sender"].replace(/\s/gi, '-'));
        if (message["sender"] === handleInput.value) {
            messageBox.classList.add('right');
        } else {
            messageBox.classList.add('left');
        }
        
        if (Chat.handles.indexOf(message["sender"]) === -1) newHandle(message["sender"]);
        messagePane.appendChild(messageBox);        
        var handle = document.createElement('div');
        handle.className = "handle";
        handle.innerHTML = message["sender"];
        messageBox.appendChild(handle);
        
        for(var i in message["images"]) {
            var img = document.createElement('img');
            img.src = message["images"][i];
            messageBox.appendChild(img);
        }
        
        if (scrollDown) messagePane.scrollTop = messagePane.scrollHeight;
    };
    
    function newHandle (handle) {
        Chat.handles.push(handle);
        var colors = picker.generate();
        // colorsheet.addRule('.'+handle, 'background:'+colors[0]+';border-color:'+colors[1]);
    }
    
    var ColorPicker = function () {
        this.phi = 0.61803398875;
        this.init = Math.random();
        this.i = 0;
    };
    
    ColorPicker.prototype.generate = function() {
        this.i += this.phi;
        this.i %= 1;
        return [toRGB(this.i, 0.5, 0.95), toRGB(this.i, 0.1, 0.95)];
    };
    
    
    var picker = new ColorPicker();
    
    var colorsheet = (function() {
        var style = document.createElement('style');
    
        style.appendChild(document.createTextNode(''));
    
        document.head.appendChild(style);
    
        return style.sheet;
    })();
    
    var toRGB = function(h, s, v) {
        var r, g, b, i, f, p, q, t;
        
        i = Math.floor(h * 6);
        f = h * 6 - i;
        p = v * (1 - s);
        q = v * (1 - f * s);
        t = v * (1 - (1 - f) * s);
        switch (i % 6) {
            case 0: r = v; g = t; b = p; break;
            case 1: r = q; g = v; b = p; break;
            case 2: r = p; g = v; b = t; break;
            case 3: r = p; g = q; b = v; break;
            case 4: r = t; g = p; b = v; break;
            case 5: r = v; g = p; b = q; break;
        }
        return {
            r: Math.floor(r * 255),
            g: Math.floor(g * 255),
            b: Math.floor(b * 255)
        };
    };
    
    pubnub = PUBNUB.init({
        publish_key: 'pub-c-3aa060e1-a1e3-49a2-ac3b-07d040fcf8f3',
        subscribe_key: 'sub-c-0292b910-a12d-11e5-bdb6-0619f8945a4f',
        origin : 'pubsub.pubnub.com',
        ssl : false
    });
    
    
    pubnub.subscribe({
        channel: 'chat',
        message: function (message, env, channel) {
            new Message(message);
        },
        connect: function () {
            console.info('Connection established.');
        }
    });
    
    document.getElementById('response').addEventListener('keypress', function (e) {
        if (e.keyCode === 13) submit(this);
    });
    
    document.getElementById('send-button').addEventListener('click', function (e) {
        submit(this.previousElementSibling);
    });
    
    function submit (el) {
        var val = el.value;
        if (val === '') return;
        var handle = handleInput.value;
        if (handle === '') return;
        el.value = '';
        
        if(val[0] === '#') {
            pubnub.unsubscribe({
                channel : Chat.channel,
            });
          
            Chat.channel = val.substring(1).toLowerCase();
            pubnub.subscribe({
                channel: Chat.channel,
                message: function (message, env, channel) {
                    new Message(message);
                },
                connect: function () {
                    console.info('Connection established.');
                }
            });
            document.getElementById('channel').innerHTML = "Channel: " + Chat.channel
            return;
        }
        
        var req = new XMLHttpRequest();
        req.open('POST', '/api/send', true);
        req.setRequestHeader('Content-Type', 'application/json');
        req.send(JSON.stringify({ 'message' : val, 'handle' : handle, 'channel': Chat.channel }));
    }
})(window, document);

