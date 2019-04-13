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
	"fmt"
	"github.com/jparound30/goboxer"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
	"os"
)

// getinfoCmd represents the getinfo command
var getinfoCmd = &cobra.Command{
	Use:   "getinfo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var getinfoFolderCmd = &cobra.Command{
	Use:   "folder",
	Short: "get folder's all information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// initialization
		err := createGoboxerApiConn()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		id := cmd.Flag("id").Value.String()
		fmt.Printf("get folder info abount folderId:%s\n", id)

		folder := goboxer.NewFolder(apiConn)

		info, err := folder.GetInfo(id, goboxer.FolderAllFields)
		if err != nil {
			// error type check
			var t *goboxer.ApiStatusError
			if xerrors.As(err, &t) {
				fmt.Printf("ApiStatusError: %+v\n", t)
			} else {
				fmt.Printf("otherError: %+v\n", err)
			}
			os.Exit(1)
		}
		fmt.Printf("%+v\n", info)
	},
}

func init() {
	rootCmd.AddCommand(getinfoCmd)
	getinfoCmd.AddCommand(getinfoFolderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getinfoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getinfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	getinfoFolderCmd.Flags().StringP("id", "i", "", "folder id")
	getinfoFolderCmd.MarkFlagRequired("id")
}
