// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package transifex

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var pullTransifexCommand = &cobra.Command{
	Use:   "pull -l pt_BR [catalog folder]",
	Short: "pulls the translation files from transifex",
	Args:  cobra.ExactArgs(1),
	Run:   pullCatalog,
}

var languages = []string{}

func init() {
	pullTransifexCommand.Flags().StringSliceVarP(&languages, "languages", "l", nil, "languages")
	pullTransifexCommand.MarkFlagRequired("languages")
}

func pullCatalog(cmd *cobra.Command, args []string) {
	folder := args[0]

	for _, lang := range languages {

		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf(
				"https://www.transifex.com/api/2/project/%s/resource/%s/translation/%s/?mode=reviewed&file=po",
				project, resource, lang,
			), nil)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		req.SetBasicAuth("api", apiKey)

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		os.Remove(path.Join(folder, lang+".po"))
		file, err := os.OpenFile(path.Join(folder, lang+".po"), os.O_CREATE|os.O_RDWR, 0644)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		_, err = file.Write(b)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}
