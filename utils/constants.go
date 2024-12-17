package utils

const (
	TCP = "tcp"

	RedisKeyTrack = "track"
	MaskKeyTrack  = `{"command":%d,"track":"%s"}`

	StatusOK    = "ok"
	StatusError = "error"

	DirPlaylist = "./playlist/"
)

const (
	CmdPlay = iota + 1
	CmdPause
	CmdNext
	CmdPrev
)
