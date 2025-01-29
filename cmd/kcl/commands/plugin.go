// Copyright The KCL Authors. All rights reserved.
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"kcl-lang.io/cli/pkg/plugin"
)

// executeRunCmd the run command for the root command.
func executeRunCmd(args []string) {
	cmd := NewRunCmd()
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func isHelpOrVersionFlag(flag string) bool {
	return flag == "-h" || flag == "--help" || flag == "-v" || flag == "--version"
}

func bootstrapCmdPlugin(cmd *cobra.Command, pluginHandler plugin.PluginHandler) {
	if pluginHandler == nil {
		return
	}
	if len(os.Args) > 1 {
		cmdPathPieces := os.Args[1:]

		// only look for suitable extension executables if
		// the specified command does not already exist
		// flags cannot be placed before plugin name
		if strings.HasPrefix(cmdPathPieces[0], "-") && !isHelpOrVersionFlag(cmdPathPieces[0]) {
			executeRunCmd(cmdPathPieces)
		} else if foundCmd, _, err := cmd.Find(cmdPathPieces); err != nil {
			// Also check the commands that will be added by Cobra.
			// These commands are only added once rootCmd.Execute() is called, so we
			// need to check them explicitly here.
			var cmdName string // first "non-flag" arguments
			for _, arg := range cmdPathPieces {
				if !strings.HasPrefix(arg, "-") {
					cmdName = arg
					break
				}
			}

			builtinSubCmdExist := false
			for _, cmd := range foundCmd.Commands() {
				if cmd.Name() == cmdName {
					builtinSubCmdExist = true
					break
				}
			}
			switch cmdName {
			// Don't search for a plugin
			case "help", "completion", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
			default:
				if !builtinSubCmdExist {
					fmt.Fprintf(os.Stderr, "Have not find the subCmd for %s,try to running command with plugin ...\n", cmdName)
					log.Println(cmdPathPieces)
					if err := plugin.HandlePluginCommand(pluginHandler, cmdPathPieces, false); err != nil {
						//TODO 运行 go run cmd/kcl/main.go aaa  没有进入到这里,而是把 aaa 参数认为是 main.k 这样的文件处理了,直接报错 source is nil,但是没找到是哪里
						fmt.Fprintf(os.Stderr, "Have not find the plugin!Error: %v\n,please try to run 'kcl help' to check it. \n", err)
						os.Exit(1)
					}
					// 执行了这里
					executeRunCmd(cmdPathPieces)
				}
			}
		}
	}
}
