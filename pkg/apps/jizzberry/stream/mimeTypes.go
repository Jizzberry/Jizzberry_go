package stream

import (
	"fmt"
	"github.com/Jizzberry/Jizzberry_go/pkg/ffmpeg"
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

func MimeTypeFromFormat(format string, filepath string) string {
	if mime, ok := mimeTypes[format]; ok {
		return fmt.Sprintf("%s; %s;", mime, ffmpeg.GetCodecs(filepath))
	}
	return "unknown"
}
