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
	createCSS         bool

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
			componentTarget := target + componentName + ".jsx"

			if _, err := os.Stat(target); err != nil {
				fmt.Println("Could not open feature folder. Does it exist?")
				return
			}

			if _, err := os.Stat(componentTarget); err == nil {
				fmt.Println("Component already exists. Consider a different name?")
				return
			}

			componentFile, fileErr := os.Create(componentTarget)
			defer componentFile.Close()

			if fileErr != nil {
				fmt.Println("An error occurred when creating the file. Please try again.")
				return
			}

			className := []rune(featureFolderName + "-")

			for index, r := range componentName {
				if r >= 'A' && r <= 'Z' && index > 0 {
					className = append(className, '-', unicode.ToLower(r))
				} else {
					className = append(className, unicode.ToLower(r))
				}
			}

			reactTemplate, templateErr := template.New("reactComponent").Parse(`{{if .includeCSS}}import './{{.componentName}}.css';

{{end}}export function {{.componentName}}(props) {
  return (
    <div{{if .includeCSS}} className="{{.className}}"{{end}}>
    </div>
  );
}`)

			if templateErr != nil {
				panic(templateErr)
			}

			reactTemplate.Execute(componentFile, map[string]interface{}{
				"componentName": componentName,
				"className":     string(className),
				"includeCSS":    createCSS,
			})

			if createCSS == true {
				cssFile, cssFileErr := os.Create(target + componentName + ".css")
				defer cssFile.Close()

				if cssFileErr != nil {
					fmt.Println("An error occurred when creating the file. Please try again.")
					return
				}

				cssTemplate, cssTemplateErr := template.New("cssModule").Parse(`.{{.className}} {
  /* class properties go here */
}`)

				if cssTemplateErr != nil {
					panic(templateErr)
				}

				cssTemplate.Execute(cssFile, map[string]interface{}{
					"className": string(className),
				})
			}

			fmt.Println("New component created.")
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().StringVarP(&featureFolderName, "feature", "f", "", "feature folder name (located within \"src/features/*\")")
	addCmd.MarkFlagRequired("feature")
	addCmd.PersistentFlags().BoolVarP(&createCSS, "css", "c", false, "create an associated CSS file")
}
