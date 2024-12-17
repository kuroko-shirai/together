package utils

const (
	TCP = "tcp"

	RedisKeyTrack = "track"
	MaskKeyTrack  = `{"command":%d,"track":"%s"}`

	StatusOK    = "ok"
	StatusError = "error"

	DirPlaylists = "./playlists/"
)

const (
	CmdPlay = iota + 1
	CmdPause
	CmdNext
	CmdPrev
)
