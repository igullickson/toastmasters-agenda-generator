package cmd

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/igullickson/toastmasters-agenda-generator/tmag/internal"
)

// maximum number of times to try generating a unique agenda compared to N previous agendas
var maxTries = 50

// maximum number of past agendas to compare for unique roles
var maxNumberAgendasToCompare = 2

var scheduleFile string

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Add agendas to schedule",
	Long:  `Add agenda(s) for toastmaster meetings to a schedule.`,
	Run: func(cmd *cobra.Command, args []string) {
		var roles []string = viper.GetStringSlice("roles")
		var members []string = viper.GetStringSlice("members")

		var schedule []internal.Agenda

		// read and unmarshall past agendas
		scheduleYaml, err := os.ReadFile(scheduleFile)
		if err != nil {
			log.Fatal(err)
		}

		err = yaml.Unmarshal(scheduleYaml, &schedule)
		if err != nil {
			log.Printf("error occured unmarshalling schedule, using empty one, error: %v", err)
			schedule = make([]internal.Agenda, 0)
		}

		// try to generate an agenda with unique roles
		newAgenda, err := generateUniqueAgenda(roles, members, schedule)

		tries := 1
		for tries < maxTries && err != nil {
			newAgenda, err = generateUniqueAgenda(roles, members, schedule)
			tries++
		}

		if err != nil {
			log.Fatal(err)
		}

		// prepend agenda to schedule
		schedule = append([]internal.Agenda{newAgenda}, schedule...)

		marshalAndWrite(schedule, scheduleFile)

		out, err := yaml.Marshal(newAgenda)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("new agenda added to schedule after %d attempt(s):\n\n%s\n", tries, out)
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.PersistentFlags().StringVar(&scheduleFile, "scheduleFile", "schedule.yaml", "schedule file (default is schedule.yaml)")
}

func generateUniqueAgenda(roles []string, members []string, schedule []internal.Agenda) (internal.Agenda, error) {
	newAgenda := internal.NewAgenda()

	// Initialize map to track members assigned to each role in the previous N agendas
	previousAssignments := make(map[string][]string)
	for _, role := range roles {
		previousAssignments[role] = make([]string, maxNumberAgendasToCompare+1)
		for i := 0; i < maxNumberAgendasToCompare && i < len(schedule); i++ {
			if internal.IsSpeakerRole(role) {
				previousSpeakers := schedule[i].GetSpeakers()
				previousAssignments[role] = append(previousAssignments[role], previousSpeakers...)
			} else {
				previousMember, err := schedule[i].GetMemberForRole(role)
				if err != nil {
					return newAgenda, err
				}
				previousAssignments[role] = append(previousAssignments[role], previousMember)
			}
		}
	}

	assignedMembers := make([]string, len(members))

	// Fill the agenda ensuring every member is assigned a different role
	for _, role := range roles {
		// Shuffle members to ensure randomness
		internal.Shuffle(members)

		// Find a member that hasn't been assigned the role in the last N agendas
		var member string
		for _, m := range members {
			hasBeenAssigned := slices.Contains(assignedMembers, m)
			wasPreviouslyAssigned := slices.Contains(previousAssignments[role], m)
			if !hasBeenAssigned && !wasPreviouslyAssigned {
				member = m
				break
			}
		}

		if member == "" {
			return newAgenda, fmt.Errorf("could not find unique member for role %s", role)
		}

		// Assign the member to the role in the agenda
		newAgenda.AddAssignment(role, member)

		// Update roles assigned to the member in the previous agendas
		previousAssignments[role] = append(previousAssignments[role], member)
		assignedMembers = append(assignedMembers, member)
	}

	return newAgenda, nil
}

func marshalAndWrite(schedule []internal.Agenda, fileName string) {
	out, err := yaml.Marshal(schedule)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = os.WriteFile(fileName, out, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
}
