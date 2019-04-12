// Copyright © 2019 Nobuhiro Tabuki
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jparound30/goboxer"
	"github.com/spf13/cobra"
	"net/http"
	"os"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create folder ?")
	},
}

var createFolderCmd = &cobra.Command{
	Use:   "folder",
	Short: "create folders from file (UTF-8 only)",
	Long:  "create folders from file (UTF-8 only)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("folder called")
		err := createGobxerApiConn()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		inf := cmd.Flag("infile").Value.String()
		fmt.Printf("read from %s\n", inf)
		parent := cmd.Flag("parent").Value.String()
		fmt.Printf("create under the folder (id= %s)\n", parent)

		file, err := os.Open(inf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		type Task struct {
			name string
			id   string
			req  *goboxer.Request
		}

		folder := goboxer.NewFolder(apiConn)

		var folderNames []*Task
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineBytes := scanner.Bytes()
			lineBytes = bytes.TrimPrefix(lineBytes, UTF8_BOM)
			line := string(lineBytes)
			if line == "" {
				continue
			}
			f := &Task{
				name: line,
				req:  folder.CreateReq(parent, line, nil),
			}
			folderNames = append(folderNames, f)
		}

		if err = scanner.Err(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		batchReq := goboxer.NewBatchRequest(apiConn)

		reqs := []*goboxer.Request{}
		for i, t := range folderNames {
			if len(reqs) < 20 {
				reqs = append(reqs, t.req)
			}
			if len(reqs) == 20 || len(folderNames)-1 == i {
				response, err := batchReq.ExecuteBatch(reqs)
				if err != nil {
					fmt.Println("103")
					fmt.Println(err)
					os.Exit(1)
				}

				for _, resp := range response.Responses {
					for _, v := range folderNames {
						if v.req == resp.Request {
							if resp.ResponseCode == http.StatusCreated {
								f := &goboxer.Folder{}
								err = json.Unmarshal(resp.Body, f)
								if err != nil {
									fmt.Println("115")
									fmt.Println(err)
									os.Exit(1)
								}
								v.id = *f.ID
								break
							} else {
								v.id = "failed"
							}
						}
					}
				}
				reqs = []*goboxer.Request{}
			}

		}

		for _, v := range folderNames {
			fmt.Printf("%s,%s\n", v.name, v.id)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createFolderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createFolderCmd.Flags().StringP("infile", "i", "folder.csv", "input file path")
	createFolderCmd.MarkFlagRequired("infile")
	createFolderCmd.Flags().StringP("parent", "p", "", "parent folder id")
	createFolderCmd.MarkFlagRequired("parent")
}
