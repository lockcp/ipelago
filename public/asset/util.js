"use strict"

const httpRegex = /https?:\/\/[^\s,()!]+/;

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
  const handleErr = (errMsg) => {
    if (onFail) {
      onFail(errMsg);
      return;
    }
    if (options.alerts) {
      options.alerts.insert('danger', errMsg);
    } else {
      console.log(errMsg);
    }
  }

  if (options.buttonID) disable(options.buttonID);

  const xhr = new XMLHttpRequest();

  if (options.responseType) {
    xhr.responseType = options.responseType;
  } else {
    xhr.responseType = 'json';
  }

  xhr.open(options.method, options.url);

  xhr.onerror = () => {
    const errMsg = 'An error occurred during the transaction';
    handleErr(errMsg);
  };

  xhr.addEventListener('load', function() {
    if (this.status == 200) {
      onSuccess(this.response);
    } else {
      let errMsg;
      if (this.response && this.response.message) {
        errMsg = `${this.status} ${this.response.message}`;
      } else {
        errMsg = `${this.status} ${this.responseText}`
      }
      handleErr(errMsg);
    }
  });

  xhr.addEventListener('loadend', function() {
    if (options.buttonID) enable(options.buttonID);
    if (onAlways) onAlways(this);
  });

  xhr.send(options.body);
}

function ajaxPromise(options, n) {
  const second = 1000;
  return new Promise((resolve, reject) => {
    const timeout = window.setTimeout(() => {
      reject('timeout');
    }, n*second);

    ajax(options, (result) => {
      // onSuccess
      resolve(result);
    }, (errMsg) => {
      // onError
      reject(errMsg);
    }, () => {
      // onAlways
      window.clearTimeout(timeout);
    });
  });
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

function CreateSmallLoading() {
  const self = cc('div');
  self.view = () => m('div').attr({id:self.raw_id})
    .addClass('spinner-border spinner-border-sm').attr({role:'status'}).append(
      m('span').addClass('visually-hidden').text('Loading...')
    );
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

function CreateLogs() {  
  const self = cc('ul');
  self.insert = (msgType, msg) => {
    const time = dayjs().format('HH:mm:ss');
    const item = m('li').addClass(`text-${msgType}`).text(`${time} ${msg}`);
    $(self.id).prepend(item);
  }
  self.clear = () => {
    $(self.id).html('');
  };
  return self;
}

function changeStatus(id, status) {
  if (status == 'alive' || status == 'alive-but-no-news') {
    $(id).attr('class', 'IslandStatus badge rounded-pill bg-success').text('Alive');
  } else if (status == 'timeout') {
    $(id).attr('class', 'IslandStatus badge rounded-pill bg-waring text-dark').text('Timeout');
  } else if (status == 'down') {
    $(id).attr('class', 'IslandStatus badge rounded-pill bg-dark').text('Down');
  } else if (status == 'unfollowed') {
    $(id).attr('class', 'IslandStatus badge rounded-pill bg-secondary').text('Unfollowed');
  }
}
