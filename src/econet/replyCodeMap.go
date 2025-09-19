package econet

/*
ReplyCodeMap The BBC Microcomputer can only cope with error numbers from
the Econet file server in the range &A8 to &C0 (decimal 168 to
192), but the file server can generate many more errors than this
range allows. To overcome this problem, &A8 is used as a
composite error number so that it covers every error with a
number less than &A0.
*/
var ReplyCodeMap = map[ReturnCode][]byte{
	RCBadUserName:           []byte("BAD USERNAME\r"),
	RCNotLoggedIn:           []byte("NOT LOGGED ON\r"),
	RCTypesDontMatch:        []byte("TYPES DON'T MATCH\r"),
	RCBadRename:             []byte("BAD RENAME\r"),
	RCAlreadyAUser:          []byte("ALREADY A USER\r"),
	RCPasswordFileFull:      []byte("PASSWORD FILE FULL\r"),
	RCDirectoryFull:         []byte("DIRECTORY FULL\r"),
	RCDirectoryNotEmpty:     []byte("DIRECTORY NOT EMPTY\r"),
	RCIsADirectory:          []byte("IS A DIRECTORY\r"),
	RCMapFault:              []byte("MAP FAULT\r"),
	RCOutsideFile:           []byte("OUTSIDE FILE\r"),
	RCTooManyUsers:          []byte("TOO MANY USERS\r"),
	RCBadPassword:           []byte("BAD PASSWORD\r"),
	RCInsufficientPrivilege: []byte("INSUFFICIENT PRIVILEGE\r"),
	RCWrongPassword:         []byte("WRONG PASSWORD\r"),
	RCUserNotKnown:          []byte("USER NOT KNOWN\r"),
	RCInsufficientAccess:    []byte("INSUFFICIENT ACCESS\r"),
	RCNotADirectory:         []byte("NOT A DIRECTORY\r"),
	RCWhoAreYou:             []byte("WHO ARE YOU\r"),
	RCTooManyOpenFiles:      []byte("TOO MANY OPEN FILES\r"),
	RCFileReadOnly:          []byte("FILE READ ONLY\r"),
	RCObjectInUse:           []byte("OBJECT IN USE\r"),
	RCEntryLocked:           []byte("ENTRY LOCKED\r"),
	RCDiskFull:              []byte("DISK FULL\r"),
	RCDiscFault:             []byte("DISC FAULT\r"),
	RCDiscChanged:           []byte("DISC CHANGED\r"),
	RCDiskReadOnly:          []byte("DISC READ ONLY\r"),
	RCBadFileName:           []byte("BAD FILE NAME\r"),
	RCDriveDoorOpen:         []byte("DRIVE DOOR OPEN\r"),
	RCBadAttempts:           []byte("BAD ATTRIBUTE\r"),
	RCWriteOnly:             []byte("WRITE ONLY\r"),
	RCNotFound:              []byte("NOT FOUND\r"),
	RCSyntax:                []byte("SYNTAX\r"),
	RCChannel:               []byte("CHANNEL\r"),
	RCEOF:                   []byte("EOF\r"),
	RCBadNumber:             []byte("BAD NUMBER\r"),
	RCBadString:             []byte("BAD STRING\r"),
	RCBadCommand:            []byte("BAD COMMAND\r"),
}

/*
0X14 OBJECT NOT A DIRECTORY
    0X16 M/C NUMBER EQUALS ZERO
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

*/

/*
Acorn FS Errors < A8
--------------------
URERRD * &14 &10+4                                      ;OBJECT NOT A DIRECTORY
URERRF * &16 &10+6                                      ;M/C NUMBER EQUALS ZERO

ATERRA * &21 &20+1   "PW file not found"                ;CANNOT FIND PASSWORD FILE
ATERRI * &29 &20+9                                      ;OBJECT '$.PASSWORDS' HAS WRONG TYPE

STERRB * &32 &30+2                                      ;SIN = 0
STERRD * &34 &30+4                                      ;REF COUNT = $00
STERRE * &35 &30+5                                      ;SIZE TOO BIG OR SIZE=0 !
STERRF * &36 &30+6                                      ;INVALID WINDOW ADDRESS
STERRG * &37 &30+7                                      ;NO FREE CACHE DESCRIPTORS
STERRH * &38 &30+8                                      ;WINDOW REF COUNT > 0
STERRK * &3B &30+&0B                                    ;REF COUNT = $FF
STERRL * &3C &30+&0C                                    ;STORE DEADLOCK !!
STERRM * &3D &30+&0D                                    ;ARITH OVERFLOW IN TSTGAP

DRERRP * &41 &40+1                                      ;cdir too big
DRERRB * &42 &40+2   "Broken dir"                       ;BROKEN DIRECTORY
DRERRF * &46 &40+6                                      ;WRONG ARG TO SET/READ OBJECT ATTRIBUTES
DRERRL * &4C &40+12                                     ;NO WRITE ACCESS
DRERRN * &4E &40+14                                     ;CLIENT ASKS FOR TOO MANY ENTRIES
DRERRO * &4F &40+15                                     ;BAD ARG. TO EXAMINE

MPERRC * &53 &50+3                                      ;SIN NOT FOR START OF CHAIN
MPERRD * &54 &50+4   "INSERT A FILE SERVER DISC"        ;DISC NOT A FILE SERVER DISC
MPERRD * &54 &50+4   "Insert a Fileserver disc"         ;DISC NOT A FILE SERVER DISC
MPERRF * &56 &50+6                                      ;ILLEGAL DRIVE NUMBER
MPERRI * &59 &50+9                                      ;NEW MAP DOESN'T FIT IN OLD SPACE
MPERRJ * &5A &50+&0A                                    ;DISC OF SAME NAME ALREADY IN USE !
MPERRM * &5B &50+11                                     ;NO MORE SPACE IN MAP DESCRIPTORS
MPERRN * &5C &50+12  "Insufficient space"               ;Insufficient user free space (yay!)

RDERRA * &61 &60+1                                      ;RNDMAN.RESTART CALLED TWICE
RDERRE * &61 RDERRA                                     ;OBJECT NOT OPEN
RDERRD * &64 &60+4                                      ;HANDTB FULL
RDERRF * &66 &60+6                                      ;RNDMAN.COPY NOT FOR FILE OBJECTS
RDERRG * &67 &60+7                                      ;RANDTB FULL
RDERRI * &69 &60+9                                      ;Object not file
RDERRM * &6D &60+&0D                                    ;Invalid arg to RDSTAR

DCERRA * &71 &70+1                                      ;INVALID NUMBER OF SECTORS
DCERRB * &72 &70+2                                      ;STORE ADDRESS OVERFLOW
DCERRC * &73 &70+3                                      ;ACCESSING BEYOND END OF FILE

SAVERA * &83 &80+3                                      ;Too much Data sent from client
WAITER * &84 &80+4                                      ;Wait bombs out
COERRA * &85 &80+5                                      ;Invalid function code
SAERRC * &8A &80+&0A                                    ;File too big
SPERRA * &8C &80+&0C "Bad privilege letter"             ;Bad privilege letter
PBERRA * &8D &80+&0D                                    ;Excess Data in PUTBYTES
INFERA * &8E &80+&0E                                    ;Bad INFO argument.
ARGERR * &8F &80+&0F                                    ;Bad arg to RDAR

DTERR  * &90 &80+&10                                    ;Bad date and time
ATERRG * &AC         "Bad user name"; bad user name in PW file
URERRE * &AE &C0-&12 "Not logged on"; USER NOT LOGGED ON
DRERRK * &AF &C0-&11 "Types don't match"; TYPES DON'T MATCH

RNAMQQ * &B0 &C0-&10 "Bad rename"; Renaming across two discs
ATERRF * &B1 &C0-&F  "Already a user"; USERID ALREADY IN PASSWORD FILE
ATERRH * &B2 &C0-&E; PASSWORD FILE FULL UP
DRERRM * &B3 &C0-&D  "Dir. full"; MAX DIR SIZE REACHED
DRERRJ * &B4 &C0-&C  "Dir. not empty"; DIR NOT EMPTY
LODERA * &B5 &C0-&B  "Is a directory"; Trying to load a directory
MPERRL * &B6 &C0-&A  "Map fault"; DISC ERROR ON MAP READ/WRITE
RDERRL * &B7 &C0-&9  "Outside file"; Attempt to point outside file
URERRB * &B8 &C0-8   "Too many users"; USERTB FULL
ATERRE * &B9 &C0-7   "Bad password"; SYNTAX ERROR IN PASSWORD
ATERRD * &BA &C0-6   "Insufficient privilege"; INSUFFICIENT PRIVILEGE
ATERRC * &BB &C0-5   "Wrong password"; INCORRECT PASSWORD
ATERRB * &BC &C0-4   "PWEntry not known"; USERID NOT FOUND IN PW FILE
DRERRE * &BD &C0-3   "Insufficient access"; INSUFFICIENT ACCESS
RDERRK * &BD DRERRE  "Insufficient access"; Insufficient access
LODERB * &BD DRERRE  "Insufficient access"; Insufficient access
DRERRD * &BE &C0-2   "Not a directory"; OBJECT NOT A DIRECTORY
URERRA * &BF &C0-1   "Who are you"; MACHINE NUMBER NOT IN USERTB

RDERRC * &C0         "Too many open files"; HANDLE QUOTA EXHAUSTED
RDERRN * &C1         "File read only"; File not open for update
DRERRI * &C2; OBJECT IN USE (I.E.OPEN)
RDERRH * &C2         "Already open at station net.stn"; File already open
DRERRG * &C3         "Entry locked"; DIR ENTRY LOCKED
MPERRB * &C6         "Disc full"; DISC SPACE EXHAUSTED
DCERRF * &C7         "Disc fault"; UNRECOVERABLE DISC ERROR
MPERRA * &C8         "Disc changed"; DISC NUMBER NOT FOUND
DCERRE * &C9         "Disc read only"; DISC PROTECTED
DRERRA * &CC         "Bad file name"; INVALID SEPARATOR IN FILE TITLE
PRerrl * &CC; Printer name too long
Doorer * &CD         "Drive door open"; Door open in FSMODE U
SAERRA * &CF         "Bad attribute"; Invalid setaccess string

RDERRO * &D4         "Write only"; File not open for input
DRERRC * &D6         "Not found"; OBJECT NOT FOUND
MPERRK * &D6; DISC NAME NOT FOUND
SYNERR * &DC         "Syntax"; BAD COMMAND SYNTAX
RDERRB * &DE         "Channel"; INVALID HANDLE
RDERRJ * &DF         "EOF"; END OF FILE

NUMERR * &F0         "Bad number"; Bad decimal number
NAMERR * &FD         "Bad string"; Bad file name etc.
WOTERR * &FE         "Bad command"; Bad command

*/
