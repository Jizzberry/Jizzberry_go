package helpers

var basePath = ""

func GetWorkingDirectory() string {
	return basePath
}

func SetWorkingDirectory(path string) {
	basePath = path
}
