# iPelago beta

iPelago 是一种不同于微博、长毛象、独立博客的分布式消息分享方案。特点：分布式、免费、免备案、易用、可订阅。名称 iPelago 源于群岛的英文 archipelago, 中文名 "群岛", 请看 https://ipelago.org


## 安装运行

### 直接下载可执行文件

- Windows 用户可直接下载 https://github.com/ahui2016/ipelago/releases 或 https://gitee.com/ipelago/ipelago/releases
- 下载解压缩后，双击 ipelago.exe 运行，然后用浏览器访问 `http://127.0.0.1:996` 即可。

### 手动编译

- 先正确安装 git 和 [Go 语言](https://golang.google.cn/)。
- 由于采用了 go-sqlite3, 因此如果在 Windows 里编译, 需要先安装 [TDM-GCC Compiler](https://sourceforge.net/projects/tdm-gcc/)

```
$ cd ~
$ git clone https://github.com/ahui2016/ipelago.git
$ cd ipelago
$ go build
$ ./ipelago
```

然后用浏览器访问 `http://127.0.0.1:996` 即可。如果有端口冲突，可使用参数 `-addr` 更改端口，比如:

```
$ ./ipelago -addr 127.0.0.1:955
```
