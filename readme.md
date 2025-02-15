## 手机控制电脑软件
### 开发语言
- 前端: vue3
- 后端: go
- 使用go:embed 打包前端资源到exe中, 避免分离运行
- 按键传输采用`http`请求, 鼠标移动采用`websocket`
```shell
go get github.com/getlantern/systray
go get github.com/spf13/viper
go get github.com/gorilla/websocket
go get github.com/kbinani/screenshot
go mod tidy
```
### 功能介绍
- go语言控制鼠标移动, go语言控制键盘输入, 让手机成为电脑的键盘鼠标
- 手机控制电脑, 包括按键, 组合键, 音量媒体控制、简易鼠标、简易26键键盘等
- web浏览器端的形式，微信里面打开局域网页面就行
- 支持手机扫描电脑端二维码直接访问局域网地址
- 支持配置自定义本地应用, 手机端一键打开
- 支持简易查看远程桌面, 触控点击远程桌面, 文本输入发送
- 支持控制网页版抖音常见快捷键, 方便刷抖音
- 支持全键盘页, 支持自定义布局页, 最大6个
- 支持切换多个服务器连接
- 支持文件上传下载(参数dir指定目录)
- 小功能: 屏幕底部中滚动, 可以控制系统音量, 底部左右滚动控制键盘上下左右
```yml
# congfig.yml
port: 666
open: true
volume: true
dir: files
apps:
  - name: 微信
    path: E:\Program Files (x86)\Tencent\WeChat\WeChat.exe
  - name: 网易云
    path: E:\Program Files (x86)\NetEase\CloudMusic\cloudmusic.exe
```

### 打包安装
- 双击build.bat 可以打包window运行文件, 双击生成的dist/dcontrol.exe, 即可运行
- 打开页面`http://localhost:666/` 访问页面
- dcontrol.exe -p 666 可以指定端口
### 其他说明
- 页面效果图见`appimg`目录
- <img src="https://gcore.jsdelivr.net/gh/dhjz/dcontrol@master/appimg/app7.jpg" style="width: 340px;"/>
- <img src="https://gcore.jsdelivr.net/gh/dhjz/dcontrol@master/appimg/app2.jpg" style="width: 340px;"/>
- <img src="https://gcore.jsdelivr.net/gh/dhjz/dcontrol@master/appimg/app3.jpg" style="width: 340px;"/>
- <img src="https://gcore.jsdelivr.net/gh/dhjz/dcontrol@master/appimg/app5.jpg" style="width: 340px;"/>
- <img src="https://gcore.jsdelivr.net/gh/dhjz/dcontrol@master/appimg/app6.jpg" style="width: 340px;"/>
- <img src="https://gcore.jsdelivr.net/gh/dhjz/dcontrol@master/appimg/app8.jpg" style="width: 340px;"/>
- <img src="https://gcore.jsdelivr.net/gh/dhjz/dcontrol@master/appimg/app9.jpg" style="width: 340px;"/>

