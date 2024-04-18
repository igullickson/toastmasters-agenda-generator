package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/igullickson/toastmasters-agenda-generator/tmag/internal"
)

// maximum number of times to try generating a unique agenda compared to x previous agendas
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

		// generate new agenda
		newAgenda, err := internal.RandomAgenda(roles, members)
		if err != nil {
			log.Fatal(err)
		}

		// try to generate an agenda with unique roles
		tries := 1
		for tries < maxTries && repeatsRole(schedule, *newAgenda, maxNumberAgendasToCompare) {
			newAgenda, err = internal.RandomAgenda(roles, members)
			if err != nil {
				log.Fatal(err)
			}
			tries++
		}

		if repeatsRole(schedule, *newAgenda, maxNumberAgendasToCompare) {
			log.Printf("failed to generate agenda with unique roles compare to last %d agendas", maxNumberAgendasToCompare)
		}

		// prepend agenda to schedule
		schedule = append([]internal.Agenda{*newAgenda}, schedule...)

		marshalAndWrite(schedule, scheduleFile)

		out, err := yaml.Marshal(*newAgenda)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("new agenda added to schedule after %d tries:\n\n%s\n", tries, out)
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	scheduleCmd.PersistentFlags().StringVar(&scheduleFile, "scheduleFile", "schedule.yaml", "schedule file (default is schedule.yaml)")
}

// check whether an agenda repeats a role in the last `limit` agendas in the schedule
func repeatsRole(schedule []internal.Agenda, target internal.Agenda, limit int) bool {
	// ensure limit is within the bounds of the slice
	if limit > len(schedule) {
		limit = len(schedule)
	}

	for i := 0; i < limit; i++ {
		if target.RepeatsRole(schedule[i]) {
			return true
		}
	}
	return false
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
