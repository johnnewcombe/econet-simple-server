# Piconet Simple File Server

## Introduction

This is a work in progress!

The Piconet Simple Fileserver (PSFS) is a cross-platform Econet fileserver designed around the Piconet Econet-USB device 
designed by JimR (Jim Raynor?) and implemented and further developed by Ken Lowe. 

This software is provided as a single binary file with versions for Linux (AMD64/Arm), MacOS (Intel/Arm) and Windows 
being available. The software talks directly to the Piconet device and does not require drivers to be installed.

## The Piconet Device

The original Piconet design allowed modern computers to talk to Acorn Econet networks using a board which provides ad
interface between an Acorn ADF10 Econet module and a Raspberry Pi Pico. Later boards designed by Ken Lowe negated the 
the ADF10. The RPi Pico handles the low level Econet communications and when connected to a computer via USB provides a
simple to use interface for higher level software.

See https://github.com/jprayner/piconet for details of the device and how to obtain one.

## Installing the Fileserver

Download the latest version of the compiled binaries. Extract the archive and copy the appropriate platform version
to a suitable location.

The software can be run using the command line as follows...

    PiconetSFS --port <port-name> --root-folder <fileserver strage location>

An example of running the server on Linux is shown below.

    server --port /dev/ttyUSB0 --root-folder ./filestore

The above command assumes that the Piconet device when plugged in appears as _/dev/ttyUSB0_. When run this will create 
the _filestore_ directory if it doesn't already exist.

## Installing the Piconet Device

Simply plug the device into a suitable USB port. The device will show up as a tty e.g. _/dev/ttyUSB0_ device in Linux a
nd MacOS, in windows it will appear as a COM: port.

## Setting the Piconet Device Using a UDEV Rule

On Linux and similar architectures, creating a _udev_ rules file and a small script can ensure that the Piconet device 
is always available on the same device e.g. _/dev/econet_/.

Place the following text in a file called _60-piconet.rules_ and placed in the _/etc/udev/rules.d_ directory.  

    SUBSYSTEM=="tty", ATTRS{idVendor}=="2e8a", ATTRS{idProduct}=="000a", MODE="0660", SYMLINK+="econet", GROUP="wheel", RUN+="/usr/local/bin/piconet"

Then create a small script called _piconet_ with the following statement and place it in _/usr/local/bin_. 
Set the script to be executable, e.g. _chmod 755 /usr/local/bin/piconet_.

    /usr/bin/stty -F /dev/econet 115200




## TODO

Implement all CLI commands (Function Code 0)

Implement Remaining Function Codes (OSWORD, OSBPUT etc)
Create a client and or tests that can be used to test the above Function Codes
Implement Find Server call (see findserver.txt, test with findserver.bas)
Ensure all primitives are handled. (See The Econet Micro Guide P.32)


## Random Protocol Stuff

The client contacts the server specifying a function code (FC), this determines the fileservers actions.e.g. 
FC=0 means that the data is a CLI command to be decoded. Each supported CLI command returns a reply indicating
the result and also returns a Command Code (CC) indicating the command that was executed.

See FSOps.txt for details of the CLI Decode and all other Function code operations.

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

## File Handles

/*
Only 255 file handles a allowed per session, as file handle is identified by a single byte on the client
File servers often only allowed 255 file handles total for the server, on this system we have 255 handles
per user session e.g. logged on user at a specific machine.

* can support as many clients as we want
* @param \HomeLan\FileStore\Authentication\User $oUser
* @return int
  */

// Get the next free id for a file handle for the given user
