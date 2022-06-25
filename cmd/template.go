/*
Copyright © 2022 Aurelio Calegari, et al.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"github.com/aurc/loggo/internal/config"
	"github.com/aurc/loggo/internal/loggo"
	"github.com/spf13/cobra"
	"log"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Starts the loggo template manager app only",
	Long: `Starts the loggo template manager app only so that
you can edit or create new templates. For example:

To start from a blank canvas:
	loggo template
To start from an existing file and update or save new from it:
	loggo template --file <some existing template>
To start from an example template:
	loggo template --example=true
`,
	Run: func(cmd *cobra.Command, args []string) {
		templateFile := cmd.Flag("file").Value.String()
		example := cmd.Flag("example").Value.String() == "true"
		var cfg *config.Config
		var err error
		if len(templateFile) == 0 {
			if example {
				cfg, err = config.MakeConfig("")
			} else {
				cfg = &config.Config{
					Keys:          make([]config.Key, 0),
					LastSavedName: "",
				}
			}
		} else {
			cfg, err = config.MakeConfig(templateFile)
		}
		if err != nil {
			log.Fatalln("Unable to start app: ", err)
		}
		app := loggo.NewAppWithConfig(cfg)
		view := loggo.NewTemplateView(app, true, nil, nil)
		app.Run(view)

	},
}

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.Flags().
		StringP("file", "f", "", "Input Template File")
	templateCmd.Flags().
		StringP("example", "e", "", "Load example log template. "+
			"If `file` flag provided this flag is ignored.")
}
