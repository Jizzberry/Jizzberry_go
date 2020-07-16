package files

import (
	"bytes"
	"encoding/json"
	"github.com/Jizzberry/Jizzberry_go/pkg/database/router"
	"github.com/Jizzberry/Jizzberry_go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/actor"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/studios"
	"github.com/Jizzberry/Jizzberry_go/pkg/models/tags"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Store relation of actors, studios, tags with files in JSON format
func setAllRelations(genID int64, actors string, studios string, tags string) {
	setActorRelation(genID, actors)
	setStudioRelation(genID, studios)
	setTagRelation(genID, tags)
}

func setActorRelation(genId int64, actors string) {

	actorsSli := strings.Split(actors, ", ")

	jsonFile := readJson(router.GetJson("actorsRelation"))
	defer jsonFile.Close()

	relation := parseJson(jsonFile)
	deleteRelation(&relation, strconv.FormatInt(genId, 10))
	if relation != nil {
		if actors != "" {
			actorsModel := actor.Initialize()
			defer actorsModel.Close()

			for _, a := range actorsSli {
				tmp := actorsModel.Get(actor.Actor{Name: a})
				if len(tmp) > 0 {
					relation[strconv.FormatInt(tmp[0].GeneratedID, 10)] = append(relation[strconv.FormatInt(tmp[0].GeneratedID, 10)], strconv.FormatInt(genId, 10))
				}
			}
		}
		writeJson(jsonFile, relation)
	}
}

func GetActorRelations(ActorID string) []string {
	jsonFile := readJson(router.GetJson("actorsRelation"))
	defer jsonFile.Close()

	relation := parseJson(jsonFile)
	if val, ok := relation[ActorID]; ok {
		return val
	}
	return nil
}

func GetUsedActors() []string {
	jsonFile := readJson(router.GetJson("actorsRelation"))
	defer jsonFile.Close()

	relation := parseJson(jsonFile)

	keys := make([]string, 0, len(relation))
	for k := range relation {
		keys = append(keys, k)
	}
	return keys
}

func setStudioRelation(genId int64, studio string) {
	split := strings.Split(studio, ", ")

	jsonFile := readJson(router.GetJson("studiosRelation"))
	defer jsonFile.Close()

	relation := parseJson(jsonFile)
	deleteRelation(&relation, strconv.FormatInt(genId, 10))

	if relation != nil {
		if studio != "" {
			studiosModel := studios.Initialize()
			defer studiosModel.Close()
			for _, s := range split {
				tmp := studiosModel.Get(studios.Studio{Name: s})
				if len(tmp) > 0 {
					relation[strconv.FormatInt(tmp[0].GeneratedID, 10)] = append(relation[strconv.FormatInt(tmp[0].GeneratedID, 10)], strconv.FormatInt(genId, 10))
				}
			}
		}
		writeJson(jsonFile, relation)
	}
}

func GetStudioRelations(studioId string) []string {
	jsonFile := readJson(router.GetJson("studiosRelation"))
	defer jsonFile.Close()

	relation := parseJson(jsonFile)
	if val, ok := relation[studioId]; ok {
		return val
	}
	return nil
}

// Returns only studios which have a valid relation
func GetUsedStudios() []string {
	jsonFile := readJson(router.GetJson("studiosRelation"))
	defer jsonFile.Close()

	relation := parseJson(jsonFile)

	keys := make([]string, 0, len(relation))
	for k := range relation {
		keys = append(keys, k)
	}
	return keys
}

func setTagRelation(genId int64, tag string) {
	split := strings.Split(tag, ", ")

	jsonFile := readJson(router.GetJson("tagsRelation"))
	defer jsonFile.Close()

	relation := parseJson(jsonFile)
	deleteRelation(&relation, strconv.FormatInt(genId, 10))

	if relation != nil {
		if tag != "" {
			tagsModel := tags.Initialize()
			defer tagsModel.Close()

			for _, t := range split {
				tmp := tagsModel.Get(tags.Tag{Name: t})
				if len(tmp) > 0 {
					relation[strconv.FormatInt(tmp[0].GeneratedID, 10)] = append(relation[strconv.FormatInt(tmp[0].GeneratedID, 10)], strconv.FormatInt(genId, 10))
				}
			}
		}

		writeJson(jsonFile, relation)
	}

}

func GetTagRelations(tagId string) []string {
	jsonFile := readJson(router.GetJson("tagsRelation"))
	defer jsonFile.Close()

	relation := parseJson(jsonFile)
	if val, ok := relation[tagId]; ok {
		return val
	}
	return nil
}

func deleteRelation(relations *map[string][]string, genID string) {
	if relations != nil {
		for key, value := range *relations {
			for _, v := range value {
				if v == genID {
					delete(*relations, key)
				}
			}
		}
	}
}

func readJson(filename string) *os.File {
	jsonFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		helpers.LogError(err.Error(), component)
		return nil
	}
	return jsonFile
}

// Parse JSON to map
func parseJson(file *os.File) map[string][]string {
	byteValue, _ := ioutil.ReadAll(file)
	byteValue = bytes.Trim(byteValue, "\x00")
	relation := make(map[string][]string)

	if len(byteValue) > 0 {
		err := json.Unmarshal(byteValue, &relation)
		if err != nil {
			helpers.LogError(err.Error(), component)
		}
	}
	return relation
}

// Truncate and write JSON
func writeJson(file *os.File, relation map[string][]string) {
	marshal, err := json.Marshal(relation)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}

	err = file.Truncate(0)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}

	_, err = file.WriteAt(marshal, 0)
	if err != nil {
		helpers.LogError(err.Error(), component)
	}
}
