package econet

// TODO: Sort al this out!

type CommandCode byte

const (
	MaxBlockSize = 1280
)

// Command Codes
const (
	CCComplete CommandCode = iota
	CCSave
	CCLoad
	CCCat
	CCInfo
	CCIam
	CCSdisk
	CCDir
	CCUnrecognised
	CCLi
)

type FunctioCode byte

const (
	FCCliDecode FunctioCode = iota
	FCSave
	FCLoad
	FCExamine
	FCCatHeader
	LoadAsCommand
	FCOpenFile
	FCCloseFile
	FCGetByte
	FCPutByte
	FCGetBytes
	FCPutBytes
	FCReadRndAccInfo
	FCSetRndAccInfo
	FCReadDiskNameInfo
	FCReadLoggedOnUsers
	FCReadDateTime
	FCReadEEofInfo
	FCReadObjInfo
	FCSetObjInfo
	FCDeleteObj
	FCReadUserEnv
	FCSetUserBootOpt
	FCLogOff
	FCReadUserInfo
	FCReadFileFSVersion
	FCReadFSFreeSpace
	FCCreateDir
	FCSetTimeDate
	FCCreateFile
	FCReadUserFreeSpace
	FCSetUserFreeSpace
	FCReadClientUserId
	FCReadUsersExtended
	FCUserInfoExtended
	FCCopyData
	FCServerManagement
	_
	FCSaveFile32Bit
	FCCreateFile32Bit
	FCLoadFile32bit
	FCGetRndAccInfo32bit
	FCSetRndAccInfo32bit
	FCGetBytes32bit
	FCPutBytes32bit5
	FCExamine32bit
	FCOpenObj32bit
)

/*
const (
	FCReadAccInfoSJ      FunctioCode = 64
	FCReadWriteSysInfoSJ FunctioCode = 65
	FCReadEncrypKeySJ    FunctioCode = 66
	FCWriteBackupSJ      FunctioCode = 67
)
*/

/*

File Server Errors
==================

File servers must return consistant error numbers so that programs can tell what error has
occured. On 8-bit systems, errors <&A8 set NetError to the original error and return the error
number as &A8. On 32-bit systems, errors are returned as &0105nn. Errors without a string are
returned as FS Error XX. Some errors without strings are never returned to the caller, but
cause the server to perform different actions.

Further info...

The BBC Microcomputer can only cope with error numbers from
the Econet file server in the range &A8 to &C0 (decimal 168 to
192), but the file server can generate many more errors than this
range allows. To overcome this problem, &A8 is used as a
composite error number so that it covers every error with a
number less than &A0.

*/

type ReturnCode byte

const (
	RCOk ReturnCode = 0

	RCObjectNotADirectory          ReturnCode = 0x14
	RCMCNumberEqualsZero           ReturnCode = 0x16
	RCCannotFindPasswordFile       ReturnCode = 0x21
	RCObjectPasswordsHasWrongType  ReturnCode = 0x29
	RCSin0                         ReturnCode = 0x32
	RCRefCount0                    ReturnCode = 0x34
	RCSizeTooBigOrSize0            ReturnCode = 0x35
	RCInvalidWindowAddress         ReturnCode = 0x36
	RCNoFreeCacheDescriptions      ReturnCode = 0x37
	RCWindowRefCountGt0            ReturnCode = 0x38
	RCRefCountFF                   ReturnCode = 0x3b
	RCStoreDeadLock                ReturnCode = 0x3c
	RCArithOverflowInTstgap        ReturnCode = 0x3d
	RCCdirTooBig                   ReturnCode = 0x41
	RCBrokenDirectory              ReturnCode = 0x42
	RCWrongObjectArg               ReturnCode = 0x46
	RCNoWriteAccess                ReturnCode = 0x4c
	RCRequForTooManyEntreis        ReturnCode = 0x4e
	RCBadExamineArg                ReturnCode = 0x4f
	RCInsertFileServerDisk         ReturnCode = 0x53
	CIllegalDriveNumber            ReturnCode = 0x54
	RCNewMapTooBigForSpace         ReturnCode = 0x59
	RCDiskOfSameNameAlreadyInUse   ReturnCode = 0x5a
	RCNoMoreSpaceForMapDescriptors ReturnCode = 0x5b
	RCInsufficientSpace            ReturnCode = 0x5c
	RCRestartCalledTwice           ReturnCode = 0x61
	RCObjectNotOpen                ReturnCode = 0x61 // ??? Same return code as RCRestartCalledTwice above.
	RCHandTbFull                   ReturnCode = 0x64
	RCRndManCopyNoForFileObj       ReturnCode = 0x66
	RCRandTbFull                   ReturnCode = 0x67
	RCObjectNotFile                ReturnCode = 0x69
	RCInvalidArdToRdStar           ReturnCode = 0x6d
	RCInvalidNumberOfSectors       ReturnCode = 0x71
	RCStoreAddressOverflow         ReturnCode = 0x72
	RCAccessingBeyondEndOfFile     ReturnCode = 0x73
	RCTooMuchDataSentFromClient    ReturnCode = 0x83
	RCWaitBombsOut                 ReturnCode = 0x84
	RCInvalidFunctionCode          ReturnCode = 0x85
	RCFileTooBig                   ReturnCode = 0x8a
	RCBadPrivilegeLetter           ReturnCode = 0x8c
	RCExcessDataInPutBytes         ReturnCode = 0x8d
	RCBadInfoArgument              ReturnCode = 0x8e
	RCBadArgToRdar                 ReturnCode = 0x8f
	RCBadDateAndTime               ReturnCode = 0x90
	RCBadUserName                  ReturnCode = 0xac
	RCNotLoggedIn                  ReturnCode = 0xae
	RCTypesDontMatch               ReturnCode = 0xaf
	RCBadRename                    ReturnCode = 0xb0
	RCAlreadyAUser                 ReturnCode = 0xb1 // "BAD USERNAME",
	RCPasswordFileFull             ReturnCode = 0xb2
	RCDirectoryFull                ReturnCode = 0xb3
	RCDirectoryNotEmpty            ReturnCode = 0xb4
	RCIsADirectory                 ReturnCode = 0xb5
	RCMapFault                     ReturnCode = 0xb6
	RCOutsideFile                  ReturnCode = 0xb7
	RCTooManyUsers                 ReturnCode = 0xb8
	RCBadPassword                  ReturnCode = 0xb9
	RCInsufficientPriveledge       ReturnCode = 0xba
	RCWrongPassword                ReturnCode = 0xbb
	RCUserNotKnown                 ReturnCode = 0xbc
	RCInsufficientAccess           ReturnCode = 0xbd
	RCNotADorectoryV               ReturnCode = 0xbe
	RCWhoAreYou                    ReturnCode = 0xbf
	RCBadCommmand                  ReturnCode = 0xfe
)
