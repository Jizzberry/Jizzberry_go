package tasks

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/helpers"
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor"
	"regexp"
	"strings"
)

const component = "TasksCommon"

func Splitter(r rune) bool {
	return r == ' ' || r == '.' || r == '-' || r == '_' || r == '[' || r == ']' || r == '(' || r == ')'
}

func MatchActorToTitle(title string) []actor.Actor {
	split := strings.FieldsFunc(title, Splitter)

	recognisedActors := make([]actor.Actor, 0)

	actorsModel := actor.Initialize()
	defer actorsModel.Close()
	words := make([]string, 0)
	for i := range split {
		// Avoid articles ig
		if len(split[i]) > 2 && strings.ToLower(split[i]) != "the" {
			words = append(words, split[i])
		}
	}
	allActors := actorsModel.GetFromTitle(words)

	for _, a := range allActors {
		regex := RegexpBuilder(a.Name)
		r, err := regexp.Compile(regex)
		if err != nil {
			helpers.LogError(err.Error(), component)
			return recognisedActors
		}

		matches := r.FindAllString(title, -1)
		if len(matches) > 0 {
			recognisedActors = append(recognisedActors, a)
		}
	}

	return recognisedActors
}

func RegexpBuilder(name string) string {
	replacer := strings.NewReplacer(" ", "\\s*", "-", "\\s*", "_", "\\s*")
	regex := replacer.Replace(name)
	regex = `(?i)` + regex
	return regex
}

func MatchActorExact(name string) *[]actor.Actor {
	actors := make([]actor.Actor, 0)
	actors = append(actors, actor.Initialize().GetExact(name))

	return &actors

}
