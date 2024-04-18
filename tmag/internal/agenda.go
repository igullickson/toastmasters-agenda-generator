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

func NewAgenda() Agenda {
	return Agenda{}
}

func RandomAgenda(roles []string, members []string) (*Agenda, error) {
	if len(roles) != len(members) {
		return nil, fmt.Errorf("number of roles, %d, did not equal number of members, %d", len(roles), len(members))
	}

	shuffle(members)

	var a = NewAgenda()

	for i, v := range roles {
		err := a.AddAssignment(v, members[i])
		if err != nil {
			return nil, err
		}
	}

	return &a, nil
}

func shuffle(a []string) {
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
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
		supportedRoles := []string{
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
