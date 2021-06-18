function itemID(id) {
  return `i${id}`;
}

function MsgItem(msg) {
  const ctime = dayjs.unix(msg.CTime).format('YYYY-MM-DD HH:mm:ss');
  const Alerts = CreateAlerts();
  const self = cc('div', itemID(msg.ID))

  self.view = () => m('li').attr({id:self.raw_id}).addClass('list-group-item d-flex justify-content-between align-items-start MsgItem mb-3').append([
    m('img').addClass('me-auto').attr({src:'/public/avatar/default.jpg'}),
    m('div').append([
      m('div').addClass('')
    ]),
  ]);
}