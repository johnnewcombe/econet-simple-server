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

    PiconetSFS --port <port-name> --root-folder <fileserver storage location>

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

## The Fileserver Storage Location

This can be located anywhere, however, the PiconetSFS software will need read and write access to the chosen location.

### Filenames

File can be saved as Filename with start and execute addresses and the attribute byte in hex. If 
the access byte is not specified the the file is assumed to be unlocked, writable by the owner, 
and readable by everybody else (i.e.0x13).

e.g.

    NOS_E000_E000_FFF_13.bin
    SBASIC_C000_C2B2_FFF.bin

### Option 2 - Catalogue Entry

File info could be stored in a catalogue file of some kind. Not easy for user to add files though!

## TODO

TODO how doe we manage owner of files and directories

Need to test on System 3 to see what is presented to the server when the following commands are sent

	//     *SAVE MYDATA 3000+500
	//     *SAVE MYDATA 3000 3500
	//     *SAVE BASIC C000+1000 C2B2      // adds execution address OF C2B2
	//     *SAVE PROG 3000 3500 5050 5000  // adds execution address and load address
    
Also need to determine if the following command is valid

	//     *SAVE PROG 3000+500 5050 5000  // i.e. adds execution address and load address with '+' syntax

Need a go routine that shuts down inactive sessions etc.

## Password File

The format of the password file is based on the Level 3 fileserver,
but with each entry on a separate line. In addition, the fields are
not fixed length.

    +-----------+----------+------------+--+
    | User name | Password | Free space |Op|
    +-----------+----------+------------+--+
If the user name is longer than ten characters and does not have a '.' in
it, the first ten characters are followed by a '.', followed by the
remaining characters, in this manner:

    IF LEN(u$)>10 AND INSTR(u$,".")=0 THEN u$=LEFT$(u$,10)+"."+MID$(u$,11)

The free space is the amount of allocated space the user has in bytes.

    Option: b7=0  entry unused  b7=1 entry used
            b6=0  normal user   b6=1 system user
            b5=0  unlocked      b5=1 locked
            b4-b2 reserved
            b1-b0 logon option

An entry is occupied if the Option has bit 7 set and byte 0 is non-zero and
less than 128:

    occupied%=(?ptr%>32 AND ?ptr%<127 AND ptr%?opt%>127)

Usernames added to the password file must be a valid pathname, so must match
valid pathname syntax and have no special characters, eg " $ % & * : @ and
must not start with a digit.

## File Handles

Only 255 file handles are allowed per session, e.g. logged on user at a specific machine.

