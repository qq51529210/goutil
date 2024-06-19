package model

// Invite 类型
const (
	InvitePlay     = "Play"
	InvitePlayback = "Playback"
	InviteDownload = "Download"
	InviteVideo    = "video"
	InviteAudio    = "audio"
)

// sdp m 字段
const (
	SDPMediaRTPMap96      = "rtpmap:96 PS/90000"
	SDPMediaRTPMap97      = "rtpmap:97 MPEG4/90000"
	SDPMediaRTPMap98      = "rtpmap:98 H264/90000"
	SDPMediaRTPMap99      = "rtpmap:99 H265/90000"
	SDPMediaConnectionNew = "connection:new"
	SDPMediaSetupActive   = "setup:active"
	SDPMediaSetupPassive  = "setup:passive"
	SDPMediaFMT           = "96 97 98 99"
	SDPMediaFMT96         = "96"
	SDPMediaFMT97         = "97"
	SDPMediaFMT98         = "98"
	SDPMediaFMT99         = "99"
)
