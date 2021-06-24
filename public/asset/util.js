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

/* compoents */

function CreateLoading() {
  const self = cc('div');
  self.view = () => m('div').attr({id:self.raw_id}).addClass('text-center').append([
    m('div').addClass('spinner-border').attr({role:'status'}).append(
      m('span').addClass('visually-hidden').text('Loading...')
    ),
  ]);
  return self;
}

function CreateInfoPair(name, messages) {
  
  const infoMsg = cc('div', 'abount'+name+'msg');
  infoMsg.view = () => m('div').attr({id:InfoMsg.raw_id}).addClass('card text-dark bg-light my-3').append([
    m('div').text(name).addClass('card-header'),
    m('div').addClass('card-body text-secondary').append(
      m('div').addClass('card-text').append(messages),
    ),
  ]);
  infoMsg.setMsg = (messages) => {
    $(infoMsg.id + ' .card-text').html('').append(messages);
  };
  const infoBtn = {
    view: () => m('i').addClass('bi bi-info-circle').css({cursor:'pointer'})
    .attr({title:'显示/隐藏'+name}).click(() => { $(infoMsg.id).toggle() }),
  }
  return [infoBtn, infoMsg];
}

function CreateAlerts() {
  const alerts = cc('div');

  alerts.insertElem = (elem) => {
    $(alerts.id).prepend(elem);
  };

  alerts.insert = (msgType, msg) => {
    const time = dayjs().format('HH:mm:ss');
    const time_and_msg = `${time} ${msg}`;
    if (msgType == 'danger') {
      console.log(time_and_msg);
    }
    const elem = m('div')
      .addClass(`alert alert-${msgType} alert-dismissible fade show mt-1 mb-0`)
      .attr({role:'alert'})
      .append([
        m('span').text(time_and_msg),
        m('button').attr({type: 'button', class: "btn-close", 'data-bs-dismiss': "alert", 'aria-label':"Close"}),
      ]);
    alerts.insertElem(elem);
  };

  alerts.clear = () => {
    $(alerts.id).html('');
  };

  return alerts;
}

function StatusElem(status) {
  if (status == 'alive') {
    return m('span').text('Alive').addClass('badge bg-success');
  }
}
