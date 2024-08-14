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

func (t *t_TemplateFile) copyRaw(templateFile, outputFile string) {
	var makeFile = strings.Join([]string{t.targetDir, outputFile}, "/")
	data, err := pkg.Templates.ReadFile(templateFile)
	go_utils.Err(err)

	err = os.WriteFile(makeFile, data, 0777)
	go_utils.Err(err)
}

var initCommand cobra.Command = cobra.Command{
	Use:   "init [destinationPath]",
	Short: "initializes a new Project",
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		var templateFile = newTemplateFile()
		var targetDir = "."
		var platform = pkg.SetupPlatform(pkg.SetupPlatformOpts{})

		// Target dirs
		//=====================================================================
		if len(args) > 0 {
			targetDir = args[0]
		}
		go_utils.Err(go_utils.MkDir(targetDir))
		go_utils.MkDir(fmt.Sprintf("%s/goapi", targetDir))
		go_utils.MkDir(fmt.Sprintf("%s/static", targetDir))

		templateFile.targetDir = targetDir
		templateFile.Replacements["%%BackendBinary%%"] = fmt.Sprintf(
			".%s%s%s",
			platform.DirSeparator,
			platform.BackendBinarySlash,
			platform.BackendBinary,
		)

		// Package Names
		//=====================================================================
		nwJSPackageName, err := requestPackageName(
			"NWJS - Package name",
			"^[a-z][a-z0-9-_.]*$",
			"only use a-z 0-9 - _ and .  first letter must be a-z",
		)
		go_utils.Err(err)
		templateFile.Replacements["%%NWPackageName%%"] = nwJSPackageName

		goPackageName, err := requestPackageName(
			"Go - Package name",
			"^[a-z][a-z0-9-_./]*$",
			"only use a-z 0-9 - _ . and /  first letter must be a-z",
		)
		go_utils.Err(err)
		templateFile.Replacements["%%GOPackageName%%"] = goPackageName

		// Project Title
		//=====================================================================
		projectTitle, err := requestPackageName(
			"Project Title",
			".*",
			"all characters allowed",
		)
		templateFile.Replacements["%%ProjectTitle%%"] = projectTitle

		// OS specific PreRun Actions
		templateFile.Replacements["%%BuildPreRunCmd%%"] = platform.BuildPreRunJS

		// Copy all the Stuff
		//=====================================================================
		templateFile.copy("tmpls/package.json.tmpl", "package.json")
		templateFile.copy("tmpls/index.html.tmpl", "index.html")
		templateFile.copy("tmpls/go.mod.tmpl", "go.mod")
		templateFile.copy("tmpls/main.go.tmpl", "main.go")
		templateFile.copy("tmpls/server.go.tmpl", "goapi/server.go")

		templateFile.copy("tmpls/static_index.html.tmpl", "static/index.html")
		templateFile.copyRaw("tmpls/logo.png.tmpl", "static/logo.png")

	},
}

func requestPackageName(packageDescription string, validRegExp string, validChars string) (string, error) {
	var tries = 0
	var err error
	var reader = bufio.NewReader(os.Stdin)
	var packageName = ""

	fmt.Printf("Enter the %s (%s)\n-----------\n", packageDescription, validChars)
	for tries < 3 {

		packageName, err = reader.ReadString(byte('\n'))
		if err != nil {
			return "", err
		}

		packageName = strings.Trim(packageName, "\r\n")

		match, err := regexp.MatchString(validRegExp, packageName)
		if err != nil {
			return "", err
		}
		go_utils.Err(err)

		if !match {
			fmt.Printf("\nyour input contains invalid characters (%s)\n-----------\n", validChars)
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
