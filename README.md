# Piconet File Server

This is a work in progress!

A really simple Econet File Server based around the Piconet device.


## Protocol Stuff

The client contacts the server specifying a function code (FC), this determines the fileservers actions.e.g. 
FC=0 means that the data is a CLI command to be decoded. Each supported CLI command returns a reply indicating
the result and also returns a Command Code (CC) indicating the command that was executed.

See FSOps.txt for details of the CLI Decode and all other Function code operations.


# TODO

Implement Find Server call (see findserver.txt, test with findserver.bas)


## Protocol


Client creates a transmit block to send to the server which includes...

* Control Byte - The top bit must be se and the other seven bits are optional, typically this would be 80h
* Destination Port - Tells the server what the purpose of the message is. Ports 90-9f are used by the fileserver
* Destination Station - The number of the station to which the transmission is intended e.g. 254 for the server.
* Network Number - The network number of the station to which the transmission is intended e.g. 0 for the local network.

For example the command _I AM SYST SYST_  from station 100 (64h) results in the following scout and data frames received at the server

    scout-dst-stn=FE, scout-dst-net=00, scout-src-stn=64, scout-scr-net=00, scout-ctrl-byte=80, scout-port=99, scout-port-desc=FileServerCommand, data-dst-stn=FE, data-dst-net=00, data-src-stn=64, data-scr-net=00, data-ctrl-byte=90, data-port=00, data-port-desc=Immediate Operation, data-bytes=[00 00 00 49 20 41 4D 20 53 59 53 54 20 53 59 53 54 0D]

The reply for this is as follows;


 Password File

The format of the password file is based on the Level 3 fileserver,
but with each entry on a separate line. In addition, the fields are
not fixed length.

+-----------+----------+------------+--+
| User name | Password | Free space |Op|
+-----------+----------+------------+--+
#
If the user name is longer than ten characters and does not have a '.' in
it, the first ten characters are followed by a '.', followed by the
remaining characters, in this manner:
#
   IF LEN(u$)>10 AND INSTR(u$,".")=0 THEN u$=LEFT$(u$,10)+"."+MID$(u$,11)
#
The free space is the amount of allocated space the user has in bytes.
#
    Option: b7=0  entry unused  b7=1 entry used
            b6=0  normal user   b6=1 system user
            b5=0  unlocked      b5=1 locked
            b4-b2 reserved
            b1-b0 logon option
#
An entry is occupied if the Option has bit 7 set and byte 0 is non-zero and
less than 128:
#
   occupied%=(?ptr%>32 AND ?ptr%<127 AND ptr%?opt%>127)
#
Usernames added to the password file must be a valid pathname, so must match
valid pathname syntax and have no special characters, eg " $ % & * : @ and
must not start with a digit.