import matplotlib.pyplot as plt
import numpy as np
import matplotlib.animation as animation
import socket
import struct
import array

fig = plt.figure()
ax = plt.axes(xlim=(-3, 3), ylim=(-3, 3))
scat = plt.scatter([], [], s=10)
plt.grid()
def init():
	b=np.array([])
	b=np.vstack((b,b))
	scat.set_offsets(b)
	return scat,

def animate(i):
	UDP_IP = ''
	UDP_PORT = 8080
	BUFFER_SIZE = 4200  # Normally 1024, but we want fast response
	sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
	sock.bind((UDP_IP, UDP_PORT))
	head=struct.Struct('< 10s d d q')
	packet, addr = sock.recvfrom(BUFFER_SIZE) # buffer size is 1024 bytes
	pklen=len(packet)
	header=list(head.unpack(packet[:34]))
	Data_format=struct.Struct('<%dd' % header[-1])
	data=array.array('d',packet[34:])
	val=list(Data_format.unpack_from(data))
	print '-'*100
	#print 'Packet number=',i
	#print 'Packet length=',pklen
	print 'Header=',header
	real=np.array(val[0::2])
	imag=np.array(val[1::2])
	#print real,imag
	symbols=np.vstack((real,imag))
	#print symbols
	scat.set_offsets(symbols)
	return scat,


ani = animation.FuncAnimation(fig, animate, frames=100,interval=100,init_func=init,blit=True)
plt.show()