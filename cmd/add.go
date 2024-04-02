package cmd

import (
	"fmt"
	"os"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
)

var (
	featureFolderName string

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a new component to the project",
		Long:  "Creates a new component in provided feature dir in the project",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Too many, or not enough arguments called. You can only specify one command at a time.")
				return
			}

			componentName := args[0]
			target := "./src/features/" + featureFolderName + "/"

			if _, err := os.Stat(target); err != nil {
				fmt.Println("Could not open feature folder. Does it exist?")
				return
			}

			if _, err := os.Stat(target + componentName + ".jsx"); err == nil {
				fmt.Println("Component already exists. Consider a different name?")
				return
			}

			componentFile, fileErr := os.Create(target + componentName + ".jsx")
			className := []rune(featureFolderName + "-")

			if fileErr != nil {
				fmt.Println("An error occurred when creating the file. Please try again.")
				return
			}

			for index, r := range componentName {
				if r >= 'A' && r <= 'Z' && index > 0 {
					className = append(className, '-', unicode.ToLower(r))
				} else {
					className = append(className, unicode.ToLower(r))
				}
			}

			reactTemplate, templateErr := template.New("reactComponent").Parse(`export function {{.componentName}}(props) {
  return (
    <div className="{{.className}}">
    </div>
  );
}`)

			if templateErr != nil {
				panic(templateErr)
			}

			reactTemplate.Execute(componentFile, map[string]interface{}{
				"componentName": componentName,
				"className":     string(className),
			})

			componentFile.Close()

			fmt.Println("New component created.")
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().StringVarP(&featureFolderName, "feature", "f", "", "feature folder name (located within \"src/features/*\")")
	addCmd.MarkFlagRequired("feature")
}
