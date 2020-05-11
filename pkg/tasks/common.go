package tasks

import (
	"github.com/Jizzberry/Jizzberry-go/pkg/models/actor"
	"strings"
)

func Splitter(r rune) bool {
	return r == ' ' || r == '.' || r == '-' || r == '_' || r == '[' || r == ']' || r == '(' || r == ')'
}

func MatchName(name string) actor.Actors {
	split := strings.FieldsFunc(name, Splitter)

	recognisedActors := actor.Actors{}

	actorsModel := actor.Initialize()
	defer actorsModel.Close()
	words := make([]string, 0)
	for i := range split {
		// Avoid articles ig
		if len(split[i]) > 2 && strings.ToLower(split[i]) != "the" {
			words = append(words, split[i])
		}
	}
	allActors := actorsModel.Get(words)
	for i := range words {
		for _, act := range allActors[i] {
			actorSplit := strings.FieldsFunc(act.Name, Splitter)
			if len(actorSplit) > 1 {
				if strings.ToLower(actorSplit[0]) == strings.ToLower(words[i]) {

					// Check if both words match with found name
					if len(words) > i+1 && (strings.ToLower(actorSplit[1]) == strings.ToLower(words[i+1])) {
						if !containsActors(recognisedActors.Actors, act) {
							recognisedActors.Actors = append(recognisedActors.Actors, act)
						}
					}

					// Japanese names have their last name before the first name sometimes
					if i > 0 && strings.ToLower(actorSplit[1]) == strings.ToLower(words[i-1]) {
						if !containsActors(recognisedActors.Actors, act) {
							recognisedActors.Actors = append(recognisedActors.Actors, act)
						}
					}
				}
			} else {
				// If its just one word and it matches, its good enough
				if strings.ToLower(actorSplit[0]) == strings.ToLower(words[i]) {
					if !containsActors(recognisedActors.Actors, act) {
						recognisedActors.Actors = append(recognisedActors.Actors, act)
					}
				}
			}

		}
	}
	return recognisedActors
}

func MatchActorExact(name string) *actor.Actors {
	actors := make([]actor.Actor, 0)
	actors = append(actors, actor.Initialize().GetExact(name))
	recognisedActors := actor.Actors{
		Actors: actors,
	}

	return &recognisedActors

}

func containsActors(s []actor.Actor, e actor.Actor) bool {
	for _, a := range s {
		if a.GeneratedID == e.GeneratedID {
			return true
		}
	}
	return false
}
