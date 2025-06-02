function ensureElements(config) {
    var results = {};
    for (var _i = 0, config_1 = config; _i < config_1.length; _i++) {
        var _a = config_1[_i], id = _a.id, _b = _a.type, type = _b === void 0 ? HTMLElement : _b, _c = _a.required, required = _c === void 0 ? true : _c;
        var element = document.getElementById(id);
        if (!element) {
            if (required)
                throw new Error("Missing required element #".concat(id));
            results[id] = null;
            continue;
        }
        if (!(element instanceof type)) {
            throw new Error("Element #".concat(id, " is not ").concat(type.name));
        }
        results[id] = element;
    }
    return results;
}
var ws = new WebSocket("ws://192.168.1.115:8080/api/chat");
function connect() {
    try {
        ws.onopen = function () {
            console.log("Connected to WebSocket server");
        };
        ws.onmessage = function (event) {
            var messages = ensureElements([{ id: "messages", type: HTMLDivElement }]);
            var messageDisplay = messages["messages"];
            messageDisplay.innerHTML += "<p>".concat(event.data, "</p>");
        };
        ws.onclose = function () {
            console.log("Websocket connection closed, Retrying...");
            setTimeout(connect, 1000);
        };
        ws.onerror = function (error) {
            console.error("WebSocket error:", error);
        };
    }
    catch (err) {
        if (err instanceof Error) {
            console.error("error : ", err);
        }
        else {
            console.error("Unknown error");
        }
    }
}
function sendMessage() {
    var input = ensureElements([{ id: "messageInput", type: HTMLInputElement }]);
    var messageInput = input["messageInput"];
    var message = messageInput.value;
    ws.send(message);
    messageInput.value = "";
}
function button(buttonId, buttonHandler) {
    var button = ensureElements([{ id: buttonId, type: HTMLButtonElement }]);
    var messageButton = button[buttonId];
    messageButton.addEventListener("click", buttonHandler);
}
connect();
button("messageButton", sendMessage);
