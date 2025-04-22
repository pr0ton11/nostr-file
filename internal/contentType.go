package internal

// ContentType is a map of file extensions to their corresponding MIME types.
var ContentType = map[string]string{
	"txt":  "text/plain",
	"htm":  "text/html",
	"html": "text/html",
	"css":  "text/css",
	"js":   "application/javascript",
	"json": "application/json",
	"xml":  "application/xml",
	"csv":  "text/csv",
	"jpg":  "image/jpeg",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"gif":  "image/gif",
	"bmp":  "image/bmp",
	"ico":  "image/x-icon",
	"svg":  "image/svg+xml",
	"webp": "image/webp",
	"mp4":  "video/mp4",
	"asf":  "video/x-ms-asf",
	"mov":  "video/quicktime",
	"avi":  "video/x-msvideo",
	"ogv":  "video/ogg",
	"mkv":  "video/x-matroska",
	"wmv":  "video/x-ms-wmv",
	"flv":  "video/x-flv",
	"mp3":  "audio/mpeg",
	"wav":  "audio/wav",
	"ogg":  "audio/ogg",
	"flac": "audio/flac",
	"pls":  "audio/x-scpls",
	"m3u":  "application/x-mpegURL",
	"m3u8": "application/vnd.apple.mpegurl",
	"xspf": "application/xspf+xml",
	"md":   "text/markdown",
	"pdf":  "application/pdf",
	"zip":  "application/zip",
}

// GetContentType returns the content type for a given file extension.
// If the extension is not found, it returns "text/html" as a default.
func GetContentType(ext string) string {
	if contentType, ok := ContentType[ext]; ok {
		return contentType
	}
	return "text/html"
}
