package stream

import (
	"github.com/Jizzberry/Jizzberry_go/pkg/models/files"
)

var mimeTypes = map[string]string{
	"asf":                     "application/vnd.ms-asf",
	"avi":                     "video/x-msvideo",
	"flv":                     "video/x-flv",
	"matroska,webm":           "video/webm",
	"m4v":                     "video/x-m4v",
	"mov,mp4,m4a,3gp,3g2,mj2": "video/mp4",
	"mpeg":                    "video/mpeg",
	"mpegts":                  "video/mpeg",
	"mpegvideo":               "video/mpeg",
	"ogg":                     "video/ogg",
	"matroska":                "video/x-matroska",
	"webm":                    "video/webm",
}

func MimeTypeFromFormat(sceneId int64) string {
	model := files.Initialize()
	defer model.Close()

	if val := model.Get(files.Files{GeneratedID: sceneId}); len(val) > 0 {
		if mime, ok := mimeTypes[val[0].Format]; ok {
			return mime
		}
	}
	return "none"
}

func IsSupportedCodec(codec string) bool {
	switch codec {
	case
		"h264",
		"vp8",
		"vp9",
		"av1":
		return true
	}
	return false
}
