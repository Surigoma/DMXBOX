import socket

client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

client.connect(("127.0.0.1", 50000))
client.send(b"test")
resp = client.recv(512)
print(resp)
client.close()
