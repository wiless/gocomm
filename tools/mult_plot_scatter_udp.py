import matplotlib.pyplot as plt
import numpy as np
import matplotlib.animation as animation
import socket
import struct
import array



#fig = plt.figure()

nrow = 2; ncol = 2;
fig,axs = plt.subplots(nrows=nrow, ncols=nrow)
CRO=[]
for ax in axs.reshape(-1):
	ax.set_xlim(-3,3)
	ax.set_ylim(-3,3)
	ax.grid(True)
	ax.spines['right'].set_color('none')
	ax.spines['top'].set_color('none')
	ax.xaxis.set_ticks_position('bottom')
	ax.spines['bottom'].set_position(('data',0))
	ax.yaxis.set_ticks_position('left')
	ax.spines['left'].set_position(('data',0))
	#for C in CRO:
	#	C = ax.scatter([],[],s=50)

CRO1=axs[0,0].scatter([], [],s=50)
CRO2=axs[0,1].scatter([], [], s=50)
CRO3=axs[1,0].scatter([], [], s=50)
CRO4=axs[1,1].scatter([], [], s=50)

Plot_keys=dict()
No_packets=dict()
packets=0
def init():
	b=np.array([])
	b=np.vstack((b,b))
	#for C in CRO:
	#	C.set_offsets(b)
	CRO1.set_offsets(b)
	CRO2.set_offsets(b)
	CRO3.set_offsets(b)
	CRO4.set_offsets(b)
	return CRO1,CRO2,CRO3,CRO4

def animate(self):
	global packets
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
	No_packets[header[0]]=No_packets.get(header[0],0)+1
	if len(val)>0:
		packets+=1
	#print '-'*100
	#print 'Packet number=',packets
	#print 'Packet length=',pklen
	#print 'Header=',header
	print No_packets
	if not Plot_keys.get(header[0],0) :
		Plot_keys[header[0]]=len(Plot_keys)+1 

	real=np.array(val[0::2])
	imag=np.array(val[1::2])
	#print real,imag
	symbols=np.vstack((real,imag))
	if Plot_keys[header[0]]==1:
		CRO1.set_offsets(symbols)
	elif Plot_keys[header[0]]==2:
		CRO2.set_offsets(symbols)
	elif Plot_keys[header[0]]==3:
		CRO3.set_offsets(symbols)
	elif Plot_keys[header[0]]==4:
		CRO4.set_offsets(symbols)
	return CRO1,CRO2,CRO3,CRO4

ani = animation.FuncAnimation(fig, animate, frames=300,interval=1,init_func=init,blit=True)
plt.show()
