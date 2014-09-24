#!/usr/bin/env python
import socket
import struct
import pylab
from pylab import *
import time
UDP_IP = '192.168.0.24'
UDP_PORT = 8080
BUFFER_SIZE = 2048  # Normally 1024, but we want fast response
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.bind((UDP_IP, UDP_PORT))
n=1
s=struct.Struct('< 10s d d d')

def RealtimePloter(arg):
  global xval,yval
  #CurrentXAxis=pylab.arange(len(xval)-100,len(xval),1)
  CurrentXAxis=pylab.array(xval[-100:])
  line1[0].set_data(CurrentXAxis,pylab.array(yval[-100:]))
  ax.axis([CurrentXAxis.min(),CurrentXAxis.max(),-1.5,1.5])
  manager.canvas.draw()

xAchse=pylab.arange(0,100,1)
yAchse=pylab.array([0]*100)
#xAchse=0.0
#yAchse=0.0

fig = pylab.figure(1)
ax = fig.add_subplot(111)
ax.grid(True)
ax.set_title("Realtime Waveform Plot")
ax.set_xlabel("Time")
ax.set_ylabel("Amplitude")
ax.axis([0,100,-1.5,1.5])
line1=ax.plot(xAchse,yAchse,'-')

manager = pylab.get_current_fig_manager()

xval=[]
xval = [0 for x in range(100)]
yval=[]
yval = [0 for x in range(100)]



while True:
  data, addr = sock.recvfrom(BUFFER_SIZE) # buffer size is 1024 bytes
  print "Packet size:", len(data)
  pk_len=len(data)/34
  for i in range(pk_len):
    rec_data=s.unpack(data[i*34:i*34+34])
  
    print n," received message:",rec_data[1], rec_data[2],rec_data[3]
    xval.append(rec_data[2])
    yval.append(rec_data[1])
    n=n+1
  timer = fig.canvas.new_timer(interval=2)
  timer.add_callback(RealtimePloter, ())
  timer.start()
  pylab.draw()
  plt.pause(0.0001)

