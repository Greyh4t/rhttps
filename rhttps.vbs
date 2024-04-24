set ws=WScript.CreateObject("WScript.Shell")
ws.Run "rhttps.exe -listen 127.99.99.99:443 -proxy socks5://127.0.0.1:1080 -reuseport",0,true