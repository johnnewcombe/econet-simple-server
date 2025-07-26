# Piconet File Server

This is a work in progress!

A really simple Econet File Server based around the Piconet device.


## Protocol


Client creates a transmit block to send to the server which includes...

* Control Byte - The top bit must be se and the other seven bits are optional, typically this would be 80h
* Destination Port - Tells the server what the purpose of the message is. Ports 90-9f are used by the fileserver
* Destination Station - The number of the station to which the transmission is intended e.g. 254 for the server.
* Network Number - The network number of the station to which the transmission is intended e.g. 0 for the local network.

For example the command _I AM SYST SYST_  from station 100 (64h) results in the following scout and data frames received at the server

    scout-dst-stn=FE, scout-dst-net=00, scout-src-stn=64, scout-scr-net=00, scout-ctrl-byte=80, scout-port=99, scout-port-desc=FileServerCommand, data-dst-stn=FE, data-dst-net=00, data-src-stn=64, data-scr-net=00, data-ctrl-byte=90, data-port=00, data-port-desc=Immediate Operation, data-bytes=[00 00 00 49 20 41 4D 20 53 59 53 54 20 53 59 53 54 0D]

The reply for this is as follows;

