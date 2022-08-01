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
	"context"
	"log"
	"strconv"
	"time"

	"github.com/aurc/loggo/internal/loggo"
	"github.com/aurc/loggo/internal/reader"
	"github.com/spf13/cobra"
)

// streamCmd represents the stream command
var gcpStreamCmd = &cobra.Command{
	Use:   "gcp-stream",
	Short: "Continuously stream GCP stack driver logs",
	Long: `Continuously stream Google Cloud Platform log entries
from a given selected project and GCP logging filters:

	loggo gcp-stream --project myGCPProject123 --from 1m \
            --filter 'resource.labels.namespace_name="awesome-sit" AND resource.labels.container_name="some"' \
            --template 
`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName := cmd.Flag("project").Value.String()
		from := cmd.Flag("from").Value.String()
		filter := cmd.Flag("filter").Value.String()
		templateFile := cmd.Flag("template").Value.String()
		saveParams := cmd.Flag("params-save").Value.String()
		listParams := cmd.Flag("params-list").Value.String()
		lp, _ := strconv.ParseBool(listParams)
		loadParams := cmd.Flag("params-load").Value.String()
		if len(saveParams) > 0 {
			if err := reader.Save(saveParams,
				&reader.SavedParams{
					From:     from,
					Filter:   filter,
					Project:  projectName,
					Template: templateFile,
				}); err != nil {
				log.Fatal(err)
			}
		} else if lp {
			l, err := reader.List()
			if err != nil {
				log.Fatal(err)
			}
			for _, v := range l {
				v.Print()
			}
		} else {
			if len(loadParams) > 0 {
				p, err := reader.Load(loadParams)
				if err != nil {
					log.Fatal(err)
				}
				if len(templateFile) == 0 && len(p.Template) > 0 {
					templateFile = p.Template
				}
				if len(from) == 0 && len(p.From) > 0 {
					from = p.From
				}
				if len(filter) == 0 && len(p.Filter) > 0 {
					filter = p.Filter
				}
				if len(projectName) == 0 && len(p.Project) > 0 {
					projectName = p.Project
				}
			}
			if len(projectName) == 0 {
				log.Fatal("--project flag is required.")
			}
			err := reader.CheckAuth(context.Background(), projectName)
			if err != nil {
				log.Fatal("Unable to obtain GCP credentials. ", err)
			}
			time.Sleep(time.Second)
			reader := reader.MakeGCPReader(projectName, filter, reader.ParseFrom(from), nil)
			app := loggo.NewLoggoApp(reader, templateFile)
			app.Run()
		}
	},
}

func init() {
	rootCmd.AddCommand(gcpStreamCmd)
	gcpStreamCmd.Flags().
		StringP("project", "p", "", "GCP Project ID (required)")
	//gcpStreamCmd.MarkFlagRequired("project")
	gcpStreamCmd.Flags().
		StringP("from", "d", "tail",
			`Start streaming from:
  Relative: Use format "1s", "1m", "1h" or "1d", where:
            digit followed by s, m, h, d as second, minute, hour, day.
  Fixed:    Use date format as "yyyy-MM-ddH24:mm:ss", e.g. 2022-07-30T15:00:00
  Now:      Use "tail" to start from now`)
	gcpStreamCmd.Flags().
		StringP("filter", "f", "",
			"Standard GCP filters")
	gcpStreamCmd.Flags().
		StringP("template", "t", "",
			"Rendering Template")
	gcpStreamCmd.Flags().
		StringP("params-save", "", "",
			`Save the following parameters (if provided) for reuse:
  Project:   The GCP Project ID
  Template:  The rendering template to be applied.
  From:      When to start streaming from.
  Filter:    The GCP specific filter parameters.`)
	gcpStreamCmd.Flags().
		StringP("params-load", "", "",
			`Load the parameters for reuse. If any additional parameters are 
provided, it overrides the loaded parameter with the one explicitly provided.`)
	gcpStreamCmd.Flags().
		BoolP("params-list", "", false,
			"List saved gcp connection/filtering parameters for convenient reuse.")
	gcpStreamCmd.MarkFlagsMutuallyExclusive("params-save", "params-load", "params-list")
}
