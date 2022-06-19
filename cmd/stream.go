/*
Copyright © 2022 Aurelio Calegari

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
	"github.com/aurc/loggo/internal/loggo"
	"github.com/aurc/loggo/internal/reader"
	"github.com/spf13/cobra"
	"log"
)

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Continuously stream log input source",
	Long: `Continuously stream log entries from an input stream such
as the standard input (through pipe) or a input file. Note that
if it's reading from a file, it automatically detects file 
rotation and continue to stream. For example:

	loggo stream --file <file-path>
	loggo stream | <some arbitrary input>`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName := cmd.Flag("file").Value.String()
		templateFile := cmd.Flag("template").Value.String()
		inputChan := make(chan string, 1000)
		reader := reader.MakeReader(fileName)

		if err := reader.StreamInto(inputChan); err != nil {
			log.Fatalf("unable to start app %v", err)
		}
		app := loggo.NewLoggoApp(inputChan, templateFile)
		app.Run()
	},
}

func init() {
	rootCmd.AddCommand(streamCmd)
	streamCmd.Flags().
		StringP("file", "f", "", "Input Log File")
	streamCmd.Flags().
		StringP("template", "t", "", "Rendering Template")
}
