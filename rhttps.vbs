set ws=WScript.CreateObject("WScript.Shell")
ws.Run "rhttps.exe -listen 127.0.0.1:443 -proxy socks5://127.0.0.1:1080",0,true