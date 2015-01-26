#!/usr/bin/env python
import socket
import struct
import pylab
from pylab import *
import time
import array
import matplotlib.pyplot as plt
import numpy as np

def main():
	i=0
	UDP_IP = '127.0.0.1'
	UDP_PORT = 8080
	BUFFER_SIZE = 2100  # Normally 1024, but we want fast response
	sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
	sock.bind((UDP_IP, UDP_PORT))
	head=struct.Struct('< 10s d d q')
	pklen=0
	fig, ax = plt.subplots()
	line, = ax.plot(np.random.randn(256))
	plt.ion()
	plt.ylim([0,1])
	plt.ion()
	plt.show(block=False)
	plt.ylabel('Amplitude')
	plt.xlabel('Sample Number')
	while True:
		packet, addr = sock.recvfrom(BUFFER_SIZE) # buffer size is 1024 bytes
		pklen=len(packet)
		header=list(head.unpack(packet[:34]))
		Data_format=struct.Struct('<%dd' % header[-1])
		data=array.array('d',packet[34:])
		val=list(Data_format.unpack_from(data))
		#pklen,header,val=get_udp_packets()
		print '-'*100
		print 'Packet number=',i
		print 'Packet length=',pklen
		print 'Header=',header
		axleng = header[-1]
		xmin=i*axleng
		xmax=(i+1)*axleng
		xax=range(xmin,xmax)
		line.set_ydata(val)
		line.set_xdata(xax)
    	#line2.set_ydata(np.random.randn(leng))
    #line2.set_xdata(xax)
		plt.xlim(xmin,xmax)
		
    	#plt.hold(True)
    #fig.canvas.draw()
    #fig.canvas.flush_events()
		ax.draw_artist(ax.patch)
		ax.draw_artist(line)
    #ax.draw_artist(line2,'r')
		fig.canvas.update()
		fig.canvas.flush_events()
		#print 'COOOOO'
		#print 'Values=',val
		i+=1

main()