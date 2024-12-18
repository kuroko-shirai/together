package utils

const (
	TCP = "tcp"

	RedisKeyTrack = "track"

	MaskKeyRun  = `{"command":%d,"track":{"album":"%s","title":"%s"}}`
	MaskKeyStop = `{"command":%d}`

	StatusOK    = "ok"
	StatusError = "error"

	DirPlaylists = "./playlists/%s/%s"
)

const (
	CmdPlay = iota + 1
	CmdStop
	CmdNext
	CmdPrev
)
