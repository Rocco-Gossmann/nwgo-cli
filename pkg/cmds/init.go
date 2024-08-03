package cmds

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/rocco-gossmann/go_utils"
	"github.com/rocco-gossmann/nwgo-cli/pkg"
	"github.com/spf13/cobra"
)

type t_TemplateFile struct {
	targetDir    string
	Template     string
	TargetOutPut string
	Replacements map[string]string
}

func newTemplateFile() t_TemplateFile {
	var ret t_TemplateFile
	ret.Replacements = make(map[string]string)
	return ret
}

func (t *t_TemplateFile) copy(templateFile string, outputFile string) {

	var makeFile = strings.Join([]string{t.targetDir, outputFile}, "/")
	data, err := pkg.Templates.ReadFile(templateFile)
	go_utils.Err(err)

	var output string = string(data)

	for placeholder, replacement := range t.Replacements {
		output = strings.ReplaceAll(output, placeholder, string(replacement))
	}

	err = os.WriteFile(makeFile, []byte(output), 0777)
	go_utils.Err(err)

}

var initCommand cobra.Command = cobra.Command{
	Use:   "init [destinationPath]",
	Short: "initializes a new Project",
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		var templateFile = newTemplateFile()
		var targetDir = "."

		// Target dir
		if len(args) > 0 {
			targetDir = args[0]
		}
		go_utils.MkDir(targetDir)

		templateFile.targetDir = targetDir

		// Package Name
		packageName, err := requestPackageName()
		templateFile.Replacements["%%PackageName%%"] = packageName
		go_utils.Err(err)

		// Copy all the Stuff
		templateFile.copy("tmpls/package.json.tmpl", "package.json")
		templateFile.copy("tmpls/index.html.tmpl", "index.html")
		templateFile.copy("tmpls/main.js.tmpl", "main.js")

	},
}

func copyTemplateFile(src string, dst string) {

}

func requestPackageName() (string, error) {
	var tries = 0
	var err error
	var reader = bufio.NewReader(os.Stdin)
	var packageName = ""

	fmt.Println("Enter the package name (only a-z 0-9 -_ and . characters first letter must be a-z )\n-----------")
	for tries < 3 {

		packageName, err = reader.ReadString(byte('\n'))
		if err != nil {
			return "", err
		}

		packageName = strings.Trim(packageName, "\r\n")

		match, err := regexp.MatchString("^[a-z][a-z0-9-_.]*$", packageName)
		if err != nil {
			return "", err
		}
		go_utils.Err(err)

		if !match {
			fmt.Println("\nyour input contains invalid characters (only use a-z 0-9 - _ and .  first letter must be a-z)\n-----------")
			tries += 1
		} else {
			break
		}

	}

	if tries == 3 {
		fmt.Println("something seems to be wrong with your inputs => Abort")
		return "", errors.New("try limit reached")
	}

	return packageName, nil
}
