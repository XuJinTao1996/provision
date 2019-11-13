package cli

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/spf13/cobra"
)

func blobCommands(bt string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   bt,
		Short: fmt.Sprintf("Access CLI commands relating to %v", bt),
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "list [path]",
		Short: fmt.Sprintf("List all %v", bt),
		Long:  fmt.Sprintf("You can pass an optional path parameter to show just part of the %s", bt),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) <= 1 {
				return nil
			}
			return fmt.Errorf("%v: Expected 0 or 1 argument", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			req := Session.Req().List(bt)
			if len(args) == 1 {
				req.Params("path", args[0])
			}
			data := []interface{}{}
			err := req.Do(&data)
			if err != nil {
				return generateError(err, "listing %v", bt)
			} else {
				return prettyPrint(data)
			}
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:     "download [item] to [dest]",
		Aliases: []string{"show", "get"},
		Short:   fmt.Sprintf("Download the %v named [item] to [dest]", bt),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 1 || len(args) == 3 {
				return nil
			}
			return fmt.Errorf("%v requires 1 or 2 arguments", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			dest := os.Stdout
			if len(args) == 2 && args[1] != "-" {
				var err error
				dest, err = os.Create(args[1])
				if err != nil {
					return fmt.Errorf("Error opening dest file %s: %v", args[1], err)
				}
			}
			if err := Session.GetBlob(dest, bt, args[0]); err != nil {
				return generateError(err, "Failed to fetch %v: %v", bt, args[0])
			}
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "exists [item]",
		Short: fmt.Sprintf("Checks to see if [item] %s exists and prints its checksum", bt),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 1 {
				return nil
			}
			return fmt.Errorf("%v requires 1", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			sum, err := Session.GetBlobSum(bt, args[0])
			if err != nil {
				return generateError(err, "Failed to exists %v: %v", bt, args[0])
			}
			fmt.Printf("%s: %s\n", args[0], sum)
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:    "static [item]",
		Hidden: true,
		Short:  "Download [item] from the static file server. They will always go to stdout.",
		Args: func(c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("%v requires 1 argument", c.UseLine())
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			rd, err := Session.File(args[0])
			if rd != nil {
				defer rd.Close()
			}
			if err != nil {
				return err
			}
			_, err = io.Copy(os.Stdout, rd)
			return err
		},
	})
	explode := false
	upload := &cobra.Command{
		Use:   "upload [src] as [dest]",
		Short: fmt.Sprintf("Upload the %v [src] as [dest]", bt),
		Long: `The DRP files API allows exploding a compressed file, using
bsdtar, after it has been uploaded.  This can be very
helpful when multiple files or a full directory tree
are being uploaded.

This is a two stage process enabled by the --explode
flag.  The first stage simply uploads the compressed
file to the target path and location.  The second stage
explodes the file in that path.

For example: _drpcli files upload my.zip as mypath/my.zip --explode_

The above command will upload the _my.zip_ file to the
files _mypath_ location.  It will also expand all
the files in _my.zip_ into _/mypath_ after upload.
All paths in _my.zip_ will be preserved and created
relative to _/mypath_.`,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 1 || len(args) == 3 {
				return nil
			}
			return fmt.Errorf("%v requires 1 or 2 arguments", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			item := args[0]
			dest := path.Base(item)
			if len(args) == 3 {
				dest = args[2]
			}
			data, err := urlOrFileAsReadCloser(item)
			if err != nil {
				return fmt.Errorf("Error opening src file %s: %v", item, err)
			}
			defer data.Close()
			if info, err := Session.PostBlobExplode(data, explode, bt, dest); err != nil {
				return generateError(err, "Failed to post %v: %v", bt, dest)
			} else {
				return prettyPrint(info)
			}
		},
	}
	upload.Flags().BoolVar(&explode, "explode", false, "After upload, file will be untarred/unzipped in file's local path")
	cmd.AddCommand(upload)

	cmd.AddCommand(&cobra.Command{
		Use:   "destroy [item]",
		Short: fmt.Sprintf("Delete the %v [item] on the DRP server", bt),
		Args: func(c *cobra.Command, args []string) error {
			if len(args) == 1 {
				return nil
			}
			return fmt.Errorf("%v requires 1 argument", c.UseLine())
		},
		RunE: func(c *cobra.Command, args []string) error {
			if err := Session.DeleteBlob(bt, args[0]); err != nil {
				return generateError(err, "Failed to delete %v: %v", bt, args[0])
			}
			fmt.Printf("Deleted %s", args[0])
			return nil
		},
	})
	return cmd
}

func init() {
	addRegistrar(registerFile)
}

func registerFile(app *cobra.Command) {
	app.AddCommand(blobCommands("files"))
}
