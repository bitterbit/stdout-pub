var term; // for easy debugging

var initTerminal = function(theme) {
    var clearColor = "\x1B[0m";

    var ws;
    term = new Terminal();
    term.open(document.getElementById('terminal'));
    exports.fit(term);

    if (theme != undefined) {
        term.setOption("theme", theme);
    }

    var print = function(message) {
        term.writeln(" dashboard > " + message);
    };

    var waitForConnection = function(){}; 

    var start = function() {
        console.log("start", ws);
        if (ws && ws.readyState != WebSocket.OPEN) {
            return false;
        }
        if (!ws) {
            ws = new WebSocket("ws://" + document.location.host+ "/ws/dashboard");
        } else {
            print("Connected to roddy!");
        }
        
        ws.onopen = function(evt) {
            print("Connected to roddy!");
        }
        ws.onclose = function(evt) {
            print("Disconnected from roddy");
            ws = null;
            print("Waiting for new connection...");
            waitForConnection();
        }
        ws.onmessage = function(evt) {
            var msg = JSON.parse(evt.data);
            var color = colors[Math.abs(msg.Source.hashCode()) % colors.length]
            term.writeln(" $ "+color + "["+msg.Date+"] " + clearColor + msg.Message.trimEnd()+clearColor);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    waitForConnection = function(){
        var interval = 1000;
        var tmpws = new WebSocket("ws://" + document.location.host+ "/ws/dashboard");

        var onnot = function() { setTimeout(waitForConnection, interval); };
        var onconnect = function(){
            ws = tmpws;
            start();
        };

        tmpws.onopen = onconnect; 
        tmpws.onclose = onnot;

        if (tmpws.readyState == WebSocket.CLOSED) {
            return onnot();
        } else if (tmpws.readyState == WebSocket.CONNECTED) {
            return onconnect();
        }

        console.log("other?", tmpws);
    };

    var getColors = function() {
        var colors = []
        for (var i=0; i<16; i++) {
            for (var k=0; k<16; k++) {
                var code = "\u001b[38;5;" + (i * 16 + k)+"m";
                colors.push(code)
            }
        }
        console.log(colors)
        return colors;
    }

    var colors = getColors();
    start();

    return term;
}

String.prototype.hashCode = function() {
    var hash = 0, i, chr;
    if (this.length === 0) return hash;
    for (i = 0; i < this.length; i++) {
        chr   = this.charCodeAt(i);
        hash  = ((hash << 5) - hash) + chr;
        hash |= 0; // Convert to 32bit integer
    }
    return hash;
};
