"use strict"

// make a new vnode by name, or return its view.
function m(name) {
  if (jQuery.type(name) == 'string') {
    return $(document.createElement(name));
  }
  return name.view();
}

// cc creates a component with an id.
function cc(name, id, elements) {
  if (!id) id = 'r' + Math.round(Math.random() * 100000000);
  const vnode = m(name).attr('id', id);
  if (elements) vnode.append(elements);
  return {id: '#'+id, raw_id: id, view: () => vnode};
}

function hide(id) {
  $(id).addClass('d-none');
}

function show(id) {
  $(id).removeClass('d-none');
}

function toggle(id) {
  $(id).toggleClass('d-none');
}

function disable(id) {
  const nodeName = $(id).prop('nodeName');
  if (nodeName == 'BUTTON' || nodeName == 'INPUT') {
    $(id).prop('disabled', true); 
  } else {
    $(id).css('pointer-events', 'none');
  }
}

function enable(id) {
  const nodeName = $(id).prop('nodeName');
  if (nodeName == 'BUTTON' || nodeName == 'INPUT') {
    $(id).prop('disabled', false);
  } else {
    $(id).css('pointer-events', 'auto');
  }
}

// options = { method, url, body, alerts, buttonID, responseType }
function ajax(options, onSuccess, onFail, onAlways) {
  if (options.buttonID) disable(options.buttonID);
  const xhr = new XMLHttpRequest();
  xhr.open(options.method, options.url);
  xhr.onerror = () => {
    window.alert('An error occurred during the transaction');
  };
  xhr.addEventListener('load', function() {
    if (this.status == 200) {
      if (onSuccess) {
        if (options.responseType && options.responseType == 'text') {
          onSuccess(this.responseText);
          return;
        }
        const resp = this.responseText ? JSON.parse(this.responseText) : null;
        onSuccess(resp);
      }
    } else {
      let msg;
      try {
        const resp = JSON.parse(this.responseText);
        msg = resp.message ? resp.message : `${this.status} ${this.responseText}`;
      } catch {
        msg = `${this.status} ${this.responseText}`;
      }
      if (options.alerts) {
        options.alerts.insert('danger', msg);
      } else {
        console.log(msg);
      }
      if (onFail) onFail(this);
    }
  });
  xhr.addEventListener('loadend', function() {
    if (options.buttonID) enable(options.buttonID);
    if (onAlways) onAlways(this);
  });
  xhr.send(options.body);
}

// 获取地址栏的参数。
function getUrlParam(param) {
  let loc = new URL(document.location);
  return loc.searchParams.get(param);
}
