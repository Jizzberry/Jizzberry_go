package helpers

import "path/filepath"

var (
	configPath    string
	ThumbnailPath string
	DatabasePath  string
	JsonPath      string
	FFMPEGPath    string
	LogsPath      string
	StaticPath    string
	TemplatePath  string
)

const (
	UsernameKey = "username"
	PasswordKey = "password"
	SessionsKey = "sessions"
	LoginURL    = "/auth/login/"
	PrevURLKey  = "prevurl"

	RowStructTag = "row"
	PKStructTag  = "pk"
	DateLayout   = "01-02-06"

	ScraperWebsite     = "website"
	ScraperActor       = "actor"
	ScraperActorList   = "actor_list"
	ScraperStudioList  = "studio_list"
	ScraperSingleVideo = "single_video"
	ScraperVideos      = "videos"
	ScraperImage       = "image"

	ActorName   = "name"
	ActorBday   = "birthdate"
	ActorBplace = "birthplace"
	ActorHeight = "height"
	ActorWeight = "weight"

	ActorListName  = "name"
	ActorListURLID = "url_id"

	StudioListName = "name"

	VideosName = "name"
	VideosLink = "link"

	VideoTitle  = "title"
	VideoActors = "actors"
	VideoTags   = "tags"

	ImageLink = "link"

	YamlForEach       = "selector"
	YamlSelector      = "subselector"
	YamlURL           = "url"
	YamlLastPage      = "last_page"
	YamlUrlRegex      = "url_regex"
	YamlForEachAttr   = "attr"
	YamlStringRegex   = "regex"
	YamlStringReplace = "replace"

	Art = ` 
$$$$$\ $$\                     $$\                                               
   \__$$ |\__|                    $$ |                                              
      $$ |$$\ $$$$$$$$\ $$$$$$$$\ $$$$$$$\   $$$$$$\   $$$$$$\   $$$$$$\  $$\   $$\ 
      $$ |$$ |\____$$  |\____$$  |$$  __$$\ $$  __$$\ $$  __$$\ $$  __$$\ $$ |  $$ |
$$\   $$ |$$ |  $$$$ _/   $$$$ _/ $$ |  $$ |$$$$$$$$ |$$ |  \__|$$ |  \__|$$ |  $$ |
$$ |  $$ |$$ | $$  _/    $$  _/   $$ |  $$ |$$   ____|$$ |      $$ |      $$ |  $$ |
\$$$$$$  |$$ |$$$$$$$$\ $$$$$$$$\ $$$$$$$  |\$$$$$$$\ $$ |      $$ |      \$$$$$$$ |
 \______/ \__|\________|\________|\_______/  \_______|\__|      \__|       \____$$ |
                                                                          $$\   $$ |
                                                                          \$$$$$$  |
                                                                           \______/`
)

func initPaths() {
	configPath = GetWorkingDirectory()
	ThumbnailPath = filepath.Join(GetWorkingDirectory(), "assets/thumbnails")
	DatabasePath = filepath.Join(GetWorkingDirectory(), "assets/database")
	JsonPath = filepath.Join(GetWorkingDirectory(), "assets/json")
	FFMPEGPath = filepath.Join(GetWorkingDirectory(), "assets/ffmpeg")
	LogsPath = filepath.Join(GetWorkingDirectory(), "logs")
	StaticPath = filepath.Join(GetWorkingDirectory(), "web/templates/static")
	TemplatePath = filepath.Join(GetWorkingDirectory(), "web/templates/Components")
}
