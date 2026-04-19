package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/david-truong/liferay-issues-cli/internal/ui"
	"github.com/spf13/cobra"
)

var boardCmd = &cobra.Command{
	Use:   "board",
	Short: "Manage boards",
}

var boardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List boards",
	Long: `List Jira Agile boards visible to you.

Examples:
  issues board list                          # all boards for default project
  issues board list --project LPD            # boards for a specific project
  issues board list --type scrum             # only scrum boards
  issues board list --name "Site Management" # search by name`,
	RunE: boardListRun,
}

var boardViewCmd = &cobra.Command{
	Use:   "view <ID>",
	Short: "View board details and its active sprints",
	Args:  cobra.ExactArgs(1),
	RunE:  boardViewRun,
}

func init() {
	boardListCmd.Flags().StringP("project", "p", "", "filter by project key (defaults to jira.default_project)")
	boardListCmd.Flags().String("type", "", "filter by board type (scrum, kanban)")
	boardListCmd.Flags().String("name", "", "search by board name")

	boardCmd.AddCommand(boardListCmd)
	boardCmd.AddCommand(boardViewCmd)
}

func boardListRun(cmd *cobra.Command, args []string) error {
	if err := initClient(); err != nil {
		return err
	}

	project, _ := cmd.Flags().GetString("project")
	if project == "" {
		project = cfg.Jira.DefaultProject
	}
	boardType, _ := cmd.Flags().GetString("type")
	name, _ := cmd.Flags().GetString("name")

	boards, err := client.GetBoards(boardType, name, project)
	if err != nil {
		return err
	}

	if len(boards) == 0 {
		fmt.Println("No boards found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "\033[1mID\tNAME\tTYPE\033[0m\n")
	for _, b := range boards {
		fmt.Fprintf(w, "%d\t%s\t%s\n", b.ID, b.Name, b.Type)
	}
	w.Flush()

	return nil
}

func boardViewRun(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("board ID must be a number: %s", args[0])
	}

	if err := initClient(); err != nil {
		return err
	}

	board, err := client.GetBoard(id)
	if err != nil {
		return err
	}

	fmt.Printf("\033[1m%s\033[0m  (ID: %d, Type: %s)\n", board.Name, board.ID, board.Type)

	if board.Type == "kanban" {
		fmt.Println("  This is a Kanban board (no sprints).")
		return nil
	}

	sprints, err := client.GetSprints(id, "active,future")
	if err != nil {
		fmt.Printf("  Could not fetch sprints: %v\n", err)
		return nil
	}

	if len(sprints) == 0 {
		fmt.Println("  No active/future sprints.")
	} else {
		fmt.Println()
		ui.PrintSprintTable(sprints)
	}

	return nil
}
