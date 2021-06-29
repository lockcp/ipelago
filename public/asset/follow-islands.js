
async function followIslands(islands) {
  for (let i=0; i<islands.length; i++) {
    Logs.insert('dark', '正在尝试订阅: ' + islands[i]);
    try {
      const name = await follow(islands[i]);
      Logs.insert('success', '订阅成功: ' + name);
    } catch (error) {
      if (error.indexOf("UNIQUE constraint failed: island.address") >= 0) {
        // 忽略已订阅的小岛
        Logs.insert('dark', '自动忽略: ' + islands[i]);
        continue;
      }
      if (error.indexOf("DENY") >= 0) {
        // 忽略已屏蔽的小岛
        Logs.insert('dark', '自动忽略: ' + islands[i]);
        continue;
      }
      Logs.insert('danger', error);
    }
  }
  Logs.insert('primary', `处理完成，结果如下所示：`);
  enable(SubmitBtn.id);
}

function follow(islandAddr) {
  return new Promise((resolve, reject) => {
    const timeout = window.setTimeout(() => {
      reject(Error('timeout'));
    }, 10*1000);

    const body = new FormData();
    body.set('address', islandAddr);
    ajax({method:'POST',url:'/api/follow-island',body:body}, (name) => {
      // onSuccess
      resolve(name);
    }, (that) => {
      // onError
      const errMsg = `${that.status} ${that.responseText}`;
      reject(errMsg);
    }, () => {
      // onAlways
      window.clearTimeout(timeout);
    });
  });
}
