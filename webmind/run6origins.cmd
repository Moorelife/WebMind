cls
start "11111" webmind --port 11111
start "7777" webmind --port 7777 --origin 192.168.2.111:11111
start "8888" webmind --port 8888 --origin 192.168.2.111:19999
start "9999" webmind --port 9999 --origin 192.168.2.111:17777
start "17777" webmind --port 17777 --origin 192.168.2.111:8888
start "18888" webmind --port 18888 --origin 192.168.2.111:9999
start "19999" webmind --port 19999 --origin 192.168.2.111:11111
pause
