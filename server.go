package main

import (
	"net/http"

	"github.com/spf13/cobra"
)

func server(cmd *cobra.Command, args []string) {
	// FIXME this is going to need to get more
	// complicated but it will do for now.
	http.Handle("/", http.FileServer(http.Dir("./root")))
	http.ListenAndServe(":5000", nil)
}
