package api

func GetDirectory(path string) (driveMap [][]string) {

	allFiles, err := getAllFolders(filepath.FromSlash(path))
	if err != nil {
		return
	}

	for _, f := range allFiles {
		split := strings.Split(f, "/")
		driveMap = append(driveMap, []string{split[len(split)-1], f})
	}

	driveMap = append(driveMap, []string{"..", filepath.Join(path, "..")})
	return
}
