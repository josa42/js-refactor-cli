/*
Copyright © 2020 Josa Gesell <josa@gesell.me>

*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/josa42/js-refactor-cli/tree"
	"github.com/josa42/js-refactor-cli/utils"
	"github.com/spf13/cobra"
)

// moveCmd represents the move command
var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		rootPath := cmd.Flags().Lookup("root").Value.String()
		if rootPath == "" {
			rootPath, _ = os.Getwd()
		}

		rootPath, _ = filepath.Abs(rootPath)

		nodes := tree.GetNodes(rootPath)

		source := filepath.Join(rootPath, args[0])
		target := filepath.Join(rootPath, args[1])

		if isDir(source) {
			utils.WalkFiles(source, func(s string) {
				printRel(rootPath, s)
				nodes[s].Move(filepath.Join(target, filepath.Base(s)), nodes)
			})
		} else {
			if isDir(target) {
				target = filepath.Join(target, filepath.Base(source))
			}
			if node, ok := nodes[source]; ok {
				printRel(rootPath, target)
				node.Move(target, nodes)
			}
		}
	},
}

func printRel(base, path string) {
	rp, _ := filepath.Rel(base, path)
	fmt.Println(rp)
}

func isDir(p string) bool {
	return strings.HasSuffix(p, "/") || filepath.Ext(p) == ""
}

func init() {
	rootCmd.AddCommand(moveCmd)

	moveCmd.Flags().StringP("root", "r", "", "Root path")
}
