# Protocol Information

## Introduction

### Receiving Client Requests

All client requests received by the server are handled with a four-way handshake. The client sends a Scout frame 
which, if all is OK, is acknowledged by the server returning an acknowledgement frame. The client then sends 
the data frame which is again acknowledged by the server.

The Piconet firmware handles this protocol and simply presents the original Scout and Data frames to 
PiconetSFS via the Piconet RX_TRANSMIT event. This is represented in PiconetSFS by the _RxTransmit_ structure.

Full details of the Piconet commands and events can be found within the PiconetSFS documentation.

### Replying to the Client

Data is sent back to the client by using the Piconet's TX command. The Piconet firmware again handles the lower
level four-way handshake.

## Protocol Details

Each client request includes a function code which tells the server what action to take. Function code 0 for example
is a request to decode a command typed by the user e.g. *CAT, *I AM etc. As a further example, function code 9 is a BPUT request.

A full list of function codes is shown below.

    Function Code    Description
          0          Command line decoding
          1          Save
          2          Load
          3          Examine
          4          Catalogue header (Acorn only)
          5          Load as command
          6          Open file
          7          Close file
          8          Get byte
          9          Put byte
         10          Get bytes
         11          Put bytes
         12          Read random access information
         13          Set random access information
         14          Read disc name information
         15          Read logged on users
         16          Read date/time
         17          Read EOF (end of file) information
         18          Read object information
         19          Set object information
         20          Delete object
         21          Read user environment
         22          Set user's boot option
         23          Logoff
         24          Read user information
         25          Read file server version number
         26          Read file server free space
         27          Create directory, specifying size
         28          Set time/date
         29          Create file of specified size
         30          Read user free Space (Acorn only)
         31          Set user free Space (Acorn only)
         32          Read client user identifier
         33          Read Users Extended
         34          User Info Extended
         35          Copy Data
         36          Server Management (Acorn only)
         37
         38          Save file 32-bit
         39          Create file 32-bit
         40          Load file 32-bit
         41          Get 32-bit random access information
         42          Set 32-bit random access information
         43          Get Bytes 32-bit
         44          Put Bytes 32-bit
         45          Examine 32-bit
         46          Open Object 32-bit

         64          Read account information (SJ Research only)
         65          Read/write system information (SJ Research only)
         66          Read encryption key (SJ Research only)
         67          Write Backup (SJ Research only)

## Standard Requests and Replies

### Standard Request

Requests received by server are all exhibit the basic structure as shown below. 
All requests are represented by the _RxTransmit_ structure in PiconetSFS and contain two 
parts: a Scout frame and a Data frame.

Scout Frame:

    Byte 0 -   Destination Station
    Byte 0 -   Destination Station
    Byte 1 -   Destination Network
    Byte 2 -   Source Station
    Byte 3 -   Source Network
    Byte 4 -   Control Byte
    Byte 5 -   Port
    Byte 6-n - Data (scouts do not normally contain data except for broadcasts)

Data Frame:

    Byte 0 -   Destination Station
    Byte 1 -   Destination Network
    Byte 2 -   Source Station
    Byte 3 -   Source Network
    Byte 4 -   Reply Port
    Byte 5 -   Function Code
    Byte 6-n - Data (function code specific, see below)

### Standard Reply

Replies sent by the server to the client exhibit a basic structure as follows:

    Byte 0 - Source Station
    Byte 1 - Source Network
    Byte 2 - Control Byte
    Byte 3 - Port
    Byte 4 - Command Code
    Byte 5 - Return Code
    Byte 6-n - Data (function code specific, see below)

Bytes 4-n are represented with the _FSReply_ structure within PiconetSFS.

One of the following Return Codes is returned in byte 5 of the Standard Reply. If a value other than zero is
to be returned to the client, then a message followed by a CR indicating the reason for the error is returned
in bytes 6-n.

    0x00 COMMAND COMPLETE
    0x14 OBJECT NOT A DIRECTORY
    0x16 M/C NUMBER EQUALS ZERO
    0x21 CANNOT FIND PASSWORD FILE
    0x29 OBJECT '$.PASSWORDS' HAS WRONG TYPE
    0x32 SIN = 0
    0x34 REF COUNT = $00
    0x35 SIZE TOO BIG OR SIZE=0 !
    0x36 INVALID WINDOW ADDRESS
    0x37 NO FREE CACHE DESCRIPTORS
    0x38 WINDOW REF COUNT > 0
    0x3B REF COUNT = $FF
    0x3C STORE DEADLOCK
    0x3D ARITH OVERFLOW IN TSTGAP
    0x41 CDIR TOO BIG
    0x42 BROKEN DIR 
    0x46 WRONG ARG TO SET/READ OBJECT ATTRIBUTES
    0x4C NO WRITE ACCESS
    0x4E CLIENT ASKS FOR TOO MANY ENTRIES
    0x4F BAD ARG. TO EXAMINE
    0x53 SIN NOT FOR START OF CHAIN
    0x54 INSERT A FILE SERVER DISC 
    0x54 INSERT A FILESERVER DISC 
    0x56 ILLEGAL DRIVE NUMBER
    0x59 NEW MAP DOESN'T FIT IN OLD SPACE
    0x5A DISC OF SAME NAME ALREADY IN USE
    0x5B NO MORE SPACE IN MAP DESCRIPTORS
    0x5C INSUFFICIENT SPACE
    0x61 RNDMAN.RESTART CALLED TWICE
    0x61 OBJECT NOT OPEN
    0x64 HANDTB FULL
    0x66 RNDMAN.COPY NOT FOR FILE OBJECTS
    0x67 RANDTB FULL
    0x69 OBJECT NOT FILE
    0x6D INVALID ARG TO RDSTAR
    0x71 INVALID NUMBER OF SECTORS
    0x72 STORE ADDRESS OVERFLOW
    0x73 ACCESSING BEYOND END OF FILE
    0x83 TOO MUCH DATA SENT FROM CLIENT
    0x84 WAIT BOMBS OUT
    0x85 INVALID FUNCTION CODE
    0x8A FILE TOO BIG
    0x8C BAD PRIVILEGE LETTER 
    0x8D EXCESS DATA IN PUTBYTES
    0x8E BAD INFO ARGUMENT.
    0x8F BAD ARG TO RDAR
    0x90 BAD DATE AND TIME
    0xAC BAD USER NAME 
    0xAE NOT LOGGED ON 
    0xAF TYPES DON'T MATCH
    0xB0 BAD RENAME 
    0xB1 ALREADY A USER 
    0xB2 PASSWORD FILE FULL UP
    0xB3 DIR. FULL 
    0xB4 DIR. NOT EMPTY 
    0xB5 IS A DIRECTORY 
    0xB6 MAP FAULT 
    0xB7 OUTSIDE FILE 
    0xB8 TOO MANY USERS 
    0xB9 BAD PASSWORD 
    0xBA INSUFFICIENT PRIVILEGE 
    0xBB WRONG PASSWORD 
    0xBC USER NOT KNOWN 
    0xBD INSUFFICIENT ACCESS 
    0xBD INSUFFICIENT ACCESS 
    0xBD INSUFFICIENT ACCESS 
    0xBE NOT A DIRECTORY 
    0xBF WHO ARE YOU 
    0xC0 TOO MANY OPEN FILES 
    0xC1 FILE READ ONLY 
    0xC2 OBJECT IN USE (I.E. OPEN)
    0xC2 ALREADY OPEN AT STATION NET.STN 
    0xC3 ENTRY LOCKED 
    0xC6 DISC FULL 
    0xC7 DISC FAULT
    0xC8 DISC CHANGED 
    0xC9 DISC READ ONLY 
    0xCC BAD FILE NAME 
    0xCC PRINTER NAME TOO LONG
    0xCD DRIVE DOOR OPEN 
    0xCF BAD ATTRIBUTE 
    0xD4 WRITE ONLY 
    0xD6 NOT FOUND 
    0xD6 DISC NAME NOT FOUND
    0xDC SYNTAX 
    0xDE CHANNEL 
    0xDF EOF 
    0xF0 BAD NUMBER 
    0xFD BAD STRING 
    0xFE BAD COMMAND 


## Function Codes

Within the request to the server is a function code. This is used to determine what action to take.
The following sections describe the various function codes and communication details involved.

### Function Code 0

Function code 0 is used to decode a command typed by the user. The bytes 6-n of the request Data Frame 
(see standard Request above) are defined as shown below.

    Byte 6 - User Root Directory (URD)
    Byte 7 - Current Selected Directory (CSD)
    Byte 8 - Current Selected Library (CSL)
    Byte 9-n - Command to be decoded in ASCII followed by CR.

The server will try and match the ASCII command with one of the following Command Codes. The Command 
Code is returned in byte 4 of the Standard Reply and is used to determine the action, if any, for the
client to take.

    0    No Action, command complete
    1    *Save
    2    *Load
    3    *Cat
    4    *Info, *Printer, *Printout
    5    *I AM
    6    *SDisc (Acorn only)
    7    *Dir, *SDisc (SJ Research only)
    8    Unrecognised command
    9    *Lib


#### Function Code 0, Command Code 1 (*Save)

The data (bytes 6-n) of the Standard Reply that is sent back to the client is as follows.

    Byte 4 -     1
    Byte 5 -     Return Code
    Byte 6-10  - 32-bit Load Address
    Byte 11-14 - 32-bit Execute Address
    Byte 15-17 - 24-bit File Size
    Byte 18-n  - File Name in ASCII followed by CR


#### Function Code 0, Command Code 2 (*Load)

TBA.

#### Function Code 0, Command Code 3 (*Cat)

TBA.

#### Function Code 0, Command Code 4 (*Info)

TBA.

#### Function Code 0, Command Code 5 (*I AM)

The data (bytes 6-n) of the Standard Reply that is sent back to the client is as follows.

    Byte 4 - 5
    Byte 5 - Return Code
    Byte 6 - User Root Directory (URD)
    Byte 7 - Current Selected Directory (CSD)
    Byte 8 - Current Selected Library (CSL)
    Byte 9 - Boot Option

### Function Code 1

A requests with the function code set to 1 is a request to save a file. With System and Atom computers
this call is made following a *SAVE command line request (Function Code 0).
The BBC and later computers interpret the parameters to a *SAVE command internally and will enter the
protocol by issuing a save with function code se to 1.

    Byte 6 -   Data Acknowledge Port
    Byte 7 -   Current Selected Directory (CSD)
    Byte 8 -   Current Selected Library (CSL)
    Byte 9-12  32-bit Start Address                        <----- is this little endian?
    Byte 13-16 32-bit Execute Address                      <----- is this little endian?
    Byte 17-19 24-bit file size
    Byte 20-n  File Name in ASCII followed by CR

The reply from the server to the client is as follows.

    Byte 6 - Data Port
    Byte 7 - Maximum Block Size
    Byte 8 - File Leaf Name

If everything has been successful and the , the client and server will move into the data exchange phase at
which point file data will be received in blocks of size determined by the maximum block size value.



