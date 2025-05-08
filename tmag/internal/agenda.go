package internal

import (
	"fmt"
	"math/rand"
	"strings"
)

type Agenda struct {
	Toastmaster       string `yaml:"toastmaster,omitempty"`
	Speaker1          string `yaml:"speaker 1,omitempty"`
	Speaker2          string `yaml:"speaker 2,omitempty"`
	TableTopicsMaster string `yaml:"tabletopics master,omitempty"`
	GeneralEvaluator  string `yaml:"general evaluator,omitempty"`
	Evaluator1        string `yaml:"evaluator 1,omitempty"`
	Evaluator2        string `yaml:"evaluator 2,omitempty"`
	Grammarian        string `yaml:"grammarian,omitempty"`
	AhCounter         string `yaml:"ah counter,omitempty"`
	Timer             string `yaml:"timer,omitempty"`
}

var supportedRoles []string = []string{
	"toastmaster",
	"speaker",
	"speaker 1",
	"speaker 2",
	"tabletopics master",
	"general evaluator",
	"evaluator 1",
	"evaluator 2",
	"grammarian",
	"ah counter",
	"timer",
}

func NewAgenda() Agenda {
	return Agenda{}
}

func IsSpeakerRole(role string) bool {
	return strings.Contains(role, "speaker")
}

func (a *Agenda) GetSpeakers() []string {
	speakers := []string{}
	if a.Speaker1 != "" {
		speakers = append(speakers, a.Speaker1)
	}
	if a.Speaker2 != "" {
		speakers = append(speakers, a.Speaker2)
	}

	return speakers
}

func (a *Agenda) GetMemberForRole(role string) (string, error) {
	roleLower := strings.ReplaceAll(strings.ToLower(role), " ", "")
	var member string

	switch roleLower {
	case "toastmaster":
		member = a.Toastmaster
	case "speaker":
		member = a.Speaker1
	case "speaker1":
		member = a.Speaker1
	case "speaker2":
		member = a.Speaker2
	case "tabletopicsmaster":
		member = a.TableTopicsMaster
	case "generalevaluator":
		member = a.GeneralEvaluator
	case "evaluator":
		member = a.Evaluator1
	case "evaluator1":
		member = a.Evaluator1
	case "evaluator2":
		member = a.Evaluator2
	case "grammarian":
		member = a.Grammarian
	case "ahcounter":
		member = a.AhCounter
	case "timer":
		member = a.Timer
	default:
		return "", fmt.Errorf("unexpected role: %s. supported roles: %s", role, strings.Join(supportedRoles, ", "))
	}

	return member, nil
}

func (a *Agenda) AddAssignment(role string, member string) error {
	roleLower := strings.ReplaceAll(strings.ToLower(role), " ", "")

	switch roleLower {
	case "toastmaster":
		a.Toastmaster = member
	case "speaker":
		a.Speaker1 = member
	case "speaker1":
		a.Speaker1 = member
	case "speaker2":
		a.Speaker2 = member
	case "tabletopicsmaster":
		a.TableTopicsMaster = member
	case "generalevaluator":
		a.GeneralEvaluator = member
	case "evaluator":
		a.Evaluator1 = member
	case "evaluator1":
		a.Evaluator1 = member
	case "evaluator2":
		a.Evaluator2 = member
	case "grammarian":
		a.Grammarian = member
	case "ahcounter":
		a.AhCounter = member
	case "timer":
		a.Timer = member
	default:
		return fmt.Errorf("unexpected role: %s. supported roles: %s", role, strings.Join(supportedRoles, ", "))
	}

	return nil
}

func (a *Agenda) RepeatsRole(other Agenda) bool {
	return a.Toastmaster == other.Toastmaster ||
		a.Speaker1 == other.Speaker1 ||
		a.Speaker2 == other.Speaker2 ||
		a.TableTopicsMaster == other.TableTopicsMaster ||
		a.GeneralEvaluator == other.GeneralEvaluator ||
		a.Grammarian == other.Grammarian ||
		a.AhCounter == other.AhCounter ||
		a.Timer == other.Timer
}

func Shuffle(a []string) {
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
}

func RandomAgenda(roles []string, members []string) (*Agenda, error) {
	if len(roles) != len(members) {
		return nil, fmt.Errorf("number of roles, %d, did not equal number of members, %d", len(roles), len(members))
	}

	Shuffle(members)

	var a = NewAgenda()

	for i, v := range roles {
		err := a.AddAssignment(v, members[i])
		if err != nil {
			return nil, err
		}
	}

	return &a, nil
}
