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
	createCSS         bool
	featureFolderName string
	typescript        bool

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
			extension := ".jsx"
			target := "./src/features/" + featureFolderName + "/"
			if typescript == true {
				extension = ".tsx"
			}

			componentTarget := target + componentName + extension

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
			propsName := "props"

			for index, r := range componentName {
				if r >= 'A' && r <= 'Z' && index > 0 {
					className = append(className, '-', unicode.ToLower(r))
				} else {
					className = append(className, unicode.ToLower(r))
				}
			}

			if typescript == true {
				propsName = componentName + "Props"
			}

			reactTemplate, err := template.New("reactComponent").Parse(`{{if .includePropImport}}import { {{.propsName}} } from './types';
{{end}}{{if .includeCSS}}import './{{.componentName}}.css';{{end}}{{if .spacing}}

{{end}}export function {{.componentName}}({{if .includePropImport}}props: {{end}}{{.propsName}}) {
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
				"componentName":     componentName,
				"className":         string(className),
				"includeCSS":        createCSS,
				"includePropImport": typescript,
				"spacing":           typescript || createCSS,
				"propsName":         propsName,
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

			indexTarget := target + "index.js"
			if typescript == true {
				indexTarget = target + "index.ts"
			}

			if indexFileInfo, err := os.Stat(indexTarget); err == nil {
				indexFile, err := os.OpenFile(indexTarget, os.O_APPEND|os.O_WRONLY, 644)
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

			if typesFileInfo, err := os.Stat(target + "types.ts"); err == nil && typescript == true {
				typesFile, err := os.OpenFile(target+"types.ts", os.O_APPEND|os.O_WRONLY, 644)
				defer typesFile.Close()

				if err != nil {
					logFatal("Could not open index file. Aborting.", err)
					return
				}

				exportString, err := template.New("interfaceString").Parse(`{{if .notEmpty}}
{{end}}export interface {{.interfaceName}} {
  // properties go here
}`)

				if err != nil {
					logFatal("Could not create export string template. Aborting.", err)
					return
				}

				err = exportString.Execute(typesFile, map[string]interface{}{
					"interfaceName": propsName,
					"notEmpty":      typesFileInfo.Size() != 0,
				})

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
	addCmd.PersistentFlags().BoolVarP(&typescript, "typescript", "t", false, "the repository uses TypeScript")
}
