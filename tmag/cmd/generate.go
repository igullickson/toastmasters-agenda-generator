package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/igullickson/toastmasters-agenda-generator/tmag/internal"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate an agenda",
	Long:  `Generate an agenda for a Toastmaster meeting.`,
	Run: func(cmd *cobra.Command, args []string) {
		var roles []string = viper.GetStringSlice("roles")
		var members []string = viper.GetStringSlice("members")

		var a, err = internal.RandomAgenda(roles, members)
		if err != nil {
			log.Fatal(err)
			return
		}

		out, err := yaml.Marshal(a)
		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println(string(out))

		err = os.WriteFile("agenda.yaml", out, 0644)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
