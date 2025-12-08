package cmdCli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/ystyle/jvms/internal/entity"
	"github.com/ystyle/jvms/utils/file"
	"github.com/ystyle/jvms/utils/jdk"
)

// Shared flags between switch and use commands
var switchFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "as_path",
		Usage: "as_path.",
	},
}

func switch_(cfx *entity.Config) *cli.Command {
	cmd := &cli.Command{
		Name:      "switch",
		ShortName: "s",
		Usage:     "Switch to use the specified version or index number.",
		Flags:     switchFlags,
		Action:    switchFunc(cfx),
	}
	return cmd
}

// SwitchFunc is used by both switch and use commands
func switchFunc(cfx *entity.Config) func(*cli.Context) error {
	return func(c *cli.Context) error {
		v := strings.TrimSpace(c.Args().Get(0))
		if v == "" {
			return errors.New("you should input a version or index number, Type \"jvms list\" to see what is installed")
		}

		// Check if input is a number (index)
		index, err := strconv.Atoi(v)
		if err == nil && index > 0 {
			asPath := c.Bool("as_path")

			// If not as_path, try index expansion
			if !asPath {
				index, err := strconv.Atoi(v)
				if err == nil && index > 0 {
					installed := jdk.GetInstalled(cfx.Store)
					if len(installed) == 0 {
						return errors.New("no JDK installations found")
					}
					if index > len(installed) {
						return fmt.Errorf("invalid index: %d, should be between 1 and %d", index, len(installed))
					}
					v = installed[index-1]
					fmt.Printf("Using index %d to select JDK %s\n", index, v)
				}
			}
		}
		if !jdk.IsVersionInstalled(cfx.Store, v) {
			fmt.Printf("jdk %s is not installed. ", v)
			return nil
		}
		// Create or update the symlink
		if file.Exists(cfx.JavaHome) {
			err := os.Remove(cfx.JavaHome)
			if err != nil {
				return errors.New("Switch jdk failed, please manually remove " + cfx.JavaHome)
			}
		}
		cmd := exec.Command("cmd", "/C", "setx", "JAVA_HOME", cfx.JavaHome, "/M")
		err = cmd.Run()
		if err != nil {
			return errors.New("set Environment variable `JAVA_HOME` failure: Please run as admin user")
		}
		err = os.Symlink(filepath.Join(cfx.Store, v), cfx.JavaHome)
		if err != nil {
			return errors.New("Switch jdk failed, " + err.Error())
		}
		fmt.Println("Switch success.\nNow using JDK " + v)
		cfx.CurrentJDKVersion = v
		return nil
	}
}
