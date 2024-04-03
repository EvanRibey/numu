package cmd

import (
	"fmt"
	"log"
	"os"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
)

func logFatal(logString string, logError error) {
	log.Fatal(logString, logError)
	fmt.Println("An error occurred when creating the file. Please try again.")
}

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

			componentFile, err := os.Create(componentTarget)
			defer componentFile.Close()

			if err != nil {
				logFatal("Could not create component file:"+componentName+".jsx. Aborting.", err)
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

			reactTemplate, err := template.New("reactComponent").Parse(`{{if .includeCSS}}import './{{.componentName}}.css';

{{end}}export function {{.componentName}}(props) {
  return (
    <div{{if .includeCSS}} className="{{.className}}"{{end}}>
    </div>
  );
}`)

			if err != nil {
				logFatal("Could not create component template. Aborting.", err)
				return
			}

			err = reactTemplate.Execute(componentFile, map[string]interface{}{
				"componentName": componentName,
				"className":     string(className),
				"includeCSS":    createCSS,
			})

			if err != nil {
				logFatal("Could not execute component template. Aborting.", err)
				return
			}

			if createCSS == true {
				cssFile, err := os.Create(target + componentName + ".css")
				defer cssFile.Close()

				if err != nil {
					logFatal("Could not create associated CSS file. Aborting.", err)
					return
				}

				cssTemplate, err := template.New("cssModule").Parse(`.{{.className}} {
  /* class properties go here */
}`)

				if err != nil {
					logFatal("Could not create CSS template string. Aborting.", err)
					return
				}

				err = cssTemplate.Execute(cssFile, map[string]interface{}{
					"className": string(className),
				})

				if err != nil {
					logFatal("Could not execute CSS template string. Aborting.", err)
					return
				}
			}

			if indexFileInfo, indexErr := os.Stat(target + "index.js"); indexErr == nil {
				indexFile, err := os.OpenFile(target+"index.js", os.O_APPEND|os.O_WRONLY, 644)
				defer indexFile.Close()

				if err != nil {
					logFatal("Could not open index file. Aborting.", err)
					return
				}

				if indexFileInfo.Size() == 0 {
					_, err = indexFile.Write([]byte("import './" + componentName + "';"))
				} else {
					_, err = indexFile.Write([]byte("\nimport './" + componentName + "';"))
				}

				if err != nil {
					logFatal("Could not write to index file. Aborting.", err)
					return
				}
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
