"use strict";

const information = document.getElementById("information");
const room = document.getElementById("room");
const roomLabel = document.getElementById("room-label");
const name = document.getElementById("name");
const nameLabel = document.getElementById("name-label");
const post = document.getElementById("post");
const postLabel = document.getElementById("post-label");
const output = document.getElementById("output");
const openButton = document.getElementById("open-button");
const sendButton = document.getElementById("send-button");
const closeButton = document.getElementById("close-button");

let socket = null;

window.onbeforeunload = closeWs;

function openWs() {
    socket = new WebSocket("ws://" + window.location.host + "/chat?room=" + room.value + "&user=" + name.value);

    socket.onopen = function() {
        output.innerHTML = "Connected\n";
        sendButton.disabled = false;
        closeButton.disabled = false;
    };

    socket.onerror = function() {
        output.innerHTML = "Connection Error\n";
    }

    socket.onmessage = function(e) {
        const data = JSON.parse(e.data);
        if(data.system) {
            output.innerHTML = data.user_name + " " + data.message + "\n" + output.innerHTML;
        } else {
            output.innerHTML = data.user_name + ": " + data.message + "\n" + output.innerHTML;
        }
    };

    post.value = "";
    room.hidden = true;
    roomLabel.hidden = true;
    name.hidden = true;
    nameLabel.hidden = true;
    post.hidden = false;
    postLabel.hidden = false;
    openButton.hidden = true;

    const roomDisplay = document.createElement("span");
    roomDisplay.style.marginRight = "5px";
    roomDisplay.appendChild(document.createTextNode("room: " + room.value));
    const nameDisplay = document.createElement("span");
    nameDisplay.appendChild(document.createTextNode("name: " + name.value));
    information.appendChild(roomDisplay);
    information.appendChild(nameDisplay);
}

function send() {
    socket.send(JSON.stringify(
        {
            message: post.value
        }
    ));
    post.value = "";
}

function closeWs() {
    sendButton.disabled = true;
    closeButton.disabled = true;
    post.hidden = true;
    postLabel.hidden = true;
    socket.close();
    openButton.hidden = false;
    room.hidden = false;
    roomLabel.hidden = false;
    name.hidden = false;
    nameLabel.hidden = false;
    information.textContent = "";
    output.innerHTML = "Chat finished\n" + output.innerHTML;
}
