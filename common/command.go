package common

// command
const (
	CdStr    = "cd"
	LsStr    = "ls"
	ExitStr  = "exit"
	MkdirStr = "mkdir"
	PutStr   = "put"
	GetStr   = "get"
)

// commandId
const (
	empty = uint8(iota)
	CdId
	LsId
	ExitId
	MkdirId
	PutId
	GetId
)

// command -> commandId
var commands = map[string]uint8{
	CdStr:    CdId,
	LsStr:    LsId,
	ExitStr:  ExitId,
	MkdirStr: MkdirId,
	PutStr:   PutId,
	GetStr:   GetId,
}

// Get commandId by command
func GetCommandId(command string) (uint8, error) {
	if commandId, ok := commands[command]; ok {
		return commandId, nil
	} else {
		return empty, InvalidCommandErr
	}
}
