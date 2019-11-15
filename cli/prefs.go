package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	addRegistrar(registerPref)
}

func registerPref(app *cobra.Command) {
	tree := addPrefCommands()
	app.AddCommand(tree)
}

func addPrefCommands() (res *cobra.Command) {
	res = &cobra.Command{
		Use:   "prefs",
		Short: "List and set DigitalRebar Provision operational preferences",
	}

	commands := []*cobra.Command{}
	commands = append(commands, &cobra.Command{
		Use:   "list",
		Short: "List all preferences",
		RunE: func(c *cobra.Command, args []string) error {
			prefs := map[string]string{}
			if err := Session.Req().UrlFor("prefs").Do(&prefs); err != nil {
				return generateError(err, "Error listing prefs")
			}
			return prettyPrint(prefs)
		},
	})
	prefsMap := map[string]string{}
	commands = append(commands, &cobra.Command{
		Use:   "set [- | JSON or YAML Map of strings | pairs of string args]",
		Short: "Set preferences",
		Long: `WARNING! Changing baseTokenSecret or systemGrantorSecret
will invalidate existing sessions and tokens.  This will require users to
reauthenticate and may cause active WSS sessions to hang.  Do not change
these values unless you are able to reset all client sessions.`,
		Args: func(c *cobra.Command, args []string) error {
			prefsMap = map[string]string{}
			if len(args) == 1 {
				if err := into(args[0], &prefsMap); err != nil {
					return fmt.Errorf("Invalid prefs: %v\n", err)
				}
				return nil
			}
			if len(args) != 0 && len(args)%2 == 0 {
				for i := 0; i < len(args); i += 2 {
					prefsMap[args[i]] = args[i+1]
				}
				return nil
			}
			return fmt.Errorf("prefs set either takes a single argument or a multiple of two, not %d", len(args))
		},
		RunE: func(c *cobra.Command, args []string) error {
			prefs := map[string]string{}
			if err := Session.Req().Post(prefsMap).UrlFor("prefs").Do(&prefs); err != nil {
				return generateError(err, "Error setting prefs")
			}
			return prettyPrint(prefs)
		},
	})
	res.AddCommand(commands...)
	return res
}
