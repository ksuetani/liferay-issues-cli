package cmd

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open [TICKET]",
	Short: "Open a Jira issue in the browser",
	Args:  cobra.MaximumNArgs(1),
	RunE:  openRun,
}

func openRun(cmd *cobra.Command, args []string) error {
	ticket, err := resolveTicket(args)
	if err != nil {
		return err
	}

	url := "https://" + cfg.Jira.Instance + "/browse/" + ticket
	fmt.Println(url)
	return browser.OpenURL(url)
}
