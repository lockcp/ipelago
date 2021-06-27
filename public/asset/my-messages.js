
const MsgInput = cc('textarea');
const PostBtn = cc('button');
const MsgPostArea = cc('div', null, [
  m(MsgInput).addClass('form-control'),
  m('p').addClass('text-end mt-2').append(
    m(PostBtn).text('Post').addClass('btn btn-outline-primary').attr({type:'button'}).click(() => {
      const msgBody = $(MsgInput.id).val().trim();
      if (!msgBody) {
        $(MsgInput.id).focus();
        return;
      }
      const bodySize = new Blob([msgBody]).size;
      if (bodySize > 1024) {
        Alerts.insert('danger', `消息体积(${bodySize} bytes) 超过上限 1024 bytes`);
        $(MsgInput.id).focus();
        return;        
      }
      const body = new FormData();
      body.set('msg-body', msgBody);
      ajax({method:'POST',url:'/api/post-message',alerts:Alerts,buttonID:PostBtn.id,body:body},
          (msgID) => {
            const msg = {
              ID: msgID,
              Time: dayjs().unix(),
              Body: msgBody,
            };
            const item = MsgItem(msg);
            $(MsgList.id).prepend(m(item));
            item.init();
            $(MsgInput.id).val('').focus();
          });
    }),
  ),
]);
