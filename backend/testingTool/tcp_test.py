import socket
import sys

client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

if len(sys.argv) > 1:
    msg = " ".join(sys.argv[1:])
else:
    msg = input("cmd>")
client.connect(("127.0.0.1", 50000))
client.send(msg.encode("utf-8"))
resp = client.recv(512)
print(resp)
client.close()
