function itemID(id) {
  return `i${id}`;
}

function MsgItem(msg) {
  const datetime = dayjs.unix(msg.Time).format('YYYY-MM-DD HH:mm:ss');
  const Alerts = CreateAlerts();
  const self = cc('div', itemID(msg.ID))

  self.view = () => m('li').attr({id:self.raw_id}).addClass('list-group-item d-flex justify-content-start align-items-start MsgItem mb-3').append([
    m('img').addClass('Avatar').attr({src:'/public/avatars/default.jpg'}),
    m('div').addClass('ms-3').append([
      m('div').addClass('Datetime').text(datetime),
      m('span').addClass('Contents').text(msg.Body),
      m(Alerts),
    ]),
  ]);

  return self;
}

const MsgList = cc('ul');

MsgList.view = () => m('ul')
  .attr({id:MsgList.raw_id}).addClass('list-group list-group-flush');

MsgList.append = (messages) => {
  messages.forEach(msg => {
    const item = MsgItem(msg);
    $(MsgList.id).append(m(item));
  });
};