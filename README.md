# Piconet File Server

This is a work in progress!

A really simple Econet File Server based around the Piconet device.


## Protocol


Client creates a transmit block to send to the server which includes...

* Control Byte - The top bit must be se and the other seven bits are optional, typically this would be 80h
* Destination Port - Tells the server what the purpose of the message is. Ports 90-9f are used by the fileserver
* Destination Station - The number of the station to which the transmission is intended e.g. 254 for the server.
* Network Number - The network number of the station to which the transmission is intended e.g. 0 for the local network.

