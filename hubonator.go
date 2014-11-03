package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var cmdFetch = &cobra.Command{
		Use:   "fetch",
		Short: "fetch",
		Long:  "fetch",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("fetch not implemented yet.")
		},
	}

	var cmdImportTweets = &cobra.Command{
		Use:   "import_tweets",
		Short: "import",
		Long:  "import",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("import_tweets not implemented yet.")
		},
	}

	var cmdInit = &cobra.Command{
		Use:   "init",
		Short: "init",
		Long:  "init",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("init not implemented yet.")
		},
	}

	var cmdServer = &cobra.Command{
		Use:   "server",
		Short: "server",
		Long:  "server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("server not implemented yet.")
		},
	}

	var rootCmd = &cobra.Command{Use: "hubonator"}
	rootCmd.AddCommand(cmdFetch, cmdImportTweets, cmdInit, cmdServer)
	rootCmd.Execute()
}
