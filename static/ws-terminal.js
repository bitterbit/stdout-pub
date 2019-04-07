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
        term.writeln(" $ " + message);
    };

    var start = function() {
        if (ws) {
            return false;
        }
        ws = new WebSocket("ws://" + document.location.host+ "/ws/dashboard");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            var msg = JSON.parse(evt.data);
            var color = colors[Math.abs(msg.Source.hashCode()) % colors.length]
            term.writeln(" $ "+color + "["+msg.Date+"]" + clearColor + msg.Message.trimEnd())
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
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
