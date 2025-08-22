package fs

type FileTransfer struct {
	Filename         string
	StartAddress     uint32
	ExecuteAddress   uint32
	Size             uint32
	BytesTransferred int
	CurrentDirectory byte
	CurrentLibrary   byte
	FileData         []byte
}
