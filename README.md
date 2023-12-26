# rhttps

## 修复 chrome 浏览器内置翻译功能

chrome 浏览器内置翻译已经很久不能用了，常规的改 hosts 需要经常更新 ip，很麻烦，要不然就是全局代理也不完美，因此有了这个项目

## 使用方法

1. 用命令在本地启动程序 `rhttps.exe -listen 127.0.0.1:443 -proxy socks5://127.0.0.1:1080`, 代理地址改成你自己的

2. 修改 `C:\Windows\System32\drivers\etc\hosts`, 增加如下内容

```
127.0.0.1 translate.googleapis.com
127.0.0.1 translate.google.com
```

3. chrome 打开网页，内置翻译功能就能用了

## 设置无黑窗开机启动

1. 修改 rhttps.vbs 脚本中的代理参数
2. 将 rhttps.vbs 放到开机启动目录中 `C:\Users\<你的用户名>\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`

## 原理

该工具在本地监听 443 端口的流量，hosts 中将翻译相关的地址指向了 127.0.0.1，浏览器翻译接口的时候，会将翻译请求发送到 127.0.0.1 的 443 端口，此时工具会从流量中解析出要访问的域名，然后通过代理与域名服务器建立隧道，转发流量，给翻译域名挂上了代理
