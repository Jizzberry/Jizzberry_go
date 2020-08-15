package jizzberry

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

func MimeTypeFromFormat(file files.Files) string {
	if val, ok := mimeTypes[file.Format]; ok {
		return val
	}
	return "none"
}

func IsSupportedCodec(codec string) bool {
	switch codec {
	case "h264":
	case "h265":
	case "vp8":
	case "vp9":
	case "av1":
		return true
	}
	return false
}

//func Avc1ToRfc6381(filepath string) string {
//	tmp := ffmpeg.ConvertFaststart(filepath)
//	fmt.Println(tmp.Name())
//	buffer := make([]byte, 8128)
//	_, err := tmp.Read(buffer)
//	if err != nil {
//		helpers.LogError(err.Error(), component)
//		return ""
//	}
//
//	var codecFlags []byte
//	// According to http://xhelmboyx.tripod.com/formats/mp4-layout.txt
//	r := regexp.MustCompile(string([]byte{97, 118, 99, 67})) // Byte code for "avcC"
//	match := r.FindStringIndex(string(buffer)) // Find index for avcC box
//	if len(match) > 1 {
//		codecFlags = buffer[match[1]+1 : match[1]+4] // Profile, compatibility and level flags
//	}
//	//err = os.Remove(tmp.Name())
//	//if err != nil {
//	//	helpers.LogError(err.Error(), component)
//	//}
//
//	return fmt.Sprintf("codecs=\"%b\"", codecFlags)
//}
