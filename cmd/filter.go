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

// filterCmd represents the filter command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Starts the loggo filter manager app only",
	Long: `Starts the loggo filter manager app only so that
you can edit or create new filters. For example:

To start from a blank canvas:
	loggo filter
To start from an existing file and update or save new from it:
	loggo filter --file <some existing filter>
To start from an example filter:
	loggo filter --example=true
`,
	Run: func(cmd *cobra.Command, args []string) {
		filterFile := cmd.Flag("file").Value.String()
		example := cmd.Flag("example").Value.String() == "true"
		var cfg *config.Config
		var err error
		if len(filterFile) == 0 {
			if example {
				cfg, err = config.MakeConfig("")
			} else {
				cfg = &config.Config{
					Keys:          make([]config.Key, 0),
					LastSavedName: "",
				}
			}
		} else {
			cfg, err = config.MakeConfig(filterFile)
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
	rootCmd.AddCommand(filterCmd)

	filterCmd.Flags().
		StringP("file", "f", "", "Input Template File")
	filterCmd.Flags().
		StringP("example", "e", "", "Load example log filter. "+
			"If `file` flag provided this flag is ignored.")
}
