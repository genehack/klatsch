package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var cmdFetch = &cobra.Command{
		Use:   "fetch",
		Short: "fetch",
		Long:  "long fetch",
		Run:   fetch,
	}

	var cmdImportTweets = &cobra.Command{
		Use:   "import_tweets",
		Short: "import",
		Long:  "long import",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("import_tweets not implemented yet.")
		},
	}

	var cmdInit = &cobra.Command{
		Use:   "init",
		Short: "init",
		Long:  "long init",
		Run:   initStuff,
	}

	var cmdServer = &cobra.Command{
		Use:   "server",
		Short: "server",
		Long:  "long server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("server not implemented yet.")
		},
	}

	var rootCmd = &cobra.Command{Use: "klatsch"}
	rootCmd.AddCommand(cmdFetch, cmdImportTweets, cmdInit, cmdServer)
	rootCmd.Execute()
}
