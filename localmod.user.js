// ==UserScript==
// @name        0chan localmod
// @namespace   https://0chan.hk/localmod
// @description Moderate threads via localmod API.
// @downloadURL https://raw.githubusercontent.com/Kagami/localmod/master/localmod.user.js
// @updateURL   https://raw.githubusercontent.com/Kagami/localmod/master/localmod.user.js
// @include     https://0chan.hk/*
// @include     http://nullchan7msxi257.onion/*
// @version     0.0.2
// @grant       none
// ==/UserScript==

var LM_BACKEND_URL = localStorage.getItem("LM_BACKEND_URL") ||
                     "https://lm.genshiken.org";

function mod(op, arg) {
  return new Promise(function(resolve, reject) {
    var xhr = new XMLHttpRequest();
    var url = "/api/post/" + (op === "restorePost" ? "restore" : "delete");
    var data = JSON.stringify({id: arg});
    xhr.open("POST", LM_BACKEND_URL + url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.setRequestHeader("X-Token", getToken());
    xhr.onload = function() {
      if (this.status >= 200 && this.status < 400) {
        resolve();
      } else {
        try {
          var error = JSON.parse(this.responseText).error;
          reject(new Error(error));
        } catch(e) {
          reject(new Error("bad answer"));
        }
      }
    };
    xhr.onerror = function() {
      reject(new Error("network error"));
    };
    xhr.send(data);
  });
}

function embedMod(post) {
  if (!isLogged()) return;
  if (post.classList.contains("post-op")) return;

  var idNode = post.querySelector("post-id");
  var link = post.querySelector("a:last-child");
  var postId = link.hash.slice(1);

  var btns = post.querySelector(".post-button-reply").parentNode;
  var btn = document.createElement("span");
  btn.className = "post-button";
  btn.style.color = "#16a085";
  btn.title = "Local Mod";
  var icon = document.createElement("i");
  icon.className = "fa fa-fw fa-times";

  btn.addEventListener("click", function() {
    var op = icon.classList.contains("fa-times") ? "deletePost" : "restorePost";
    var cls1 = op === "deletePost" ? "fa-times" : "fa-undo";
    var cls2 = op === "deletePost" ? "fa-undo" : "fa-times";
    if (!icon.classList.contains("fa-times") &&
        !icon.classList.contains("fa-undo")) {
      return;
    }
    icon.classList.remove(cls1);
    icon.classList.add("fa-spinner", "fa-spin");
    mod(op, postId).then(function() {
      post.classList.toggle("post-deleted");
      icon.classList.add(cls2);
    }, function(e) {
      // TODO: Use notifications.
      alert(e.message);
      icon.classList.add(cls1);
    }).then(function() {
      icon.classList.remove("fa-spinner", "fa-spin");
    });
  });

  btn.appendChild(icon);
  btns.appendChild(btn);
}

function handlePosts(container) {
  Array.prototype.forEach.call(container.querySelectorAll(".post"), embedMod);
}

function handleNavigation() {
  // TODO(Kagami): Thread listing.
  var thread = document.querySelector(".threads");
  if (!thread) return;
  if (location.pathname.split("/")[1] !== "kpop") return;
  var observer = new MutationObserver(function(mutations) {
    mutations.forEach(function(mutation) {
      Array.prototype.forEach.call(mutation.addedNodes, function(node) {
        if (node.nodeType !== Node.ELEMENT_NODE) return;
        if (node.parentNode.classList.contains("thread-tree") ||
            node.classList.contains("post-popup")) {
          embedMod(node);
        } else if (node.classList.contains("thread-tree")) {
          handlePosts(node);
        }
      });
    });
  });
  observer.observe(thread, {childList: true, subtree: true});
  handlePosts(thread);
}

function isLogged() {
  return !!localStorage.getItem("lm_token");
}

function getToken() {
  return localStorage.getItem("lm_token");
}

function login() {
  var token = prompt("Enter localmod token:");
  localStorage.setItem("lm_token", token);
}

function logout() {
  localStorage.removeItem("lm_token");
}

function embedAuth() {
  var btn = document.createElement("button");
  btn.className = "btn btn-info";
  btn.style.position = "fixed";
  btn.style.right = "0px";
  btn.style.bottom = "0px";
  btn.title = "Local Mod";
  var icon = document.createElement("i");
  icon.className = "fa";
  icon.classList.add(isLogged() ? "fa-sign-out" : "fa-sign-in");

  btn.addEventListener("click", function() {
    if (icon.classList.contains("fa-sign-in")) {
      login();
    } else {
      logout();
    }
    icon.classList.toggle("fa-sign-in");
    icon.classList.toggle("fa-sign-out");
  });

  btn.appendChild(icon);
  document.body.appendChild(btn);
}

function handleApp(container) {
  if (window.app && window.app.$bus) {
    window.app.$bus.on("refreshContentDone", handleNavigation);
    return;
  }
  var observer = new MutationObserver(function() {
    if (!window.app || !window.app.$bus) return;
    observer.disconnect();
    window.app.$bus.on("refreshContentDone", handleNavigation);
  });
  observer.observe(container, {childList: true});
}

handleApp(document.body);
embedAuth();
