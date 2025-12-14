package cmdCli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/ystyle/jvms/internal/entity"
	"github.com/ystyle/jvms/utils/file"
)

func init_(config *entity.Config) *cli.Command {
	return &cli.Command{
		Name:        "init",
		Usage:       "Initialize config file",
		Description: `before init you should clear JAVA_HOME, PATH Environment variable。`,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "java_home",
				Usage: "the JAVA_HOME location",
				Value: filepath.Join(os.Getenv("ProgramFiles"), "jdk"),
			},
			cli.StringFlag{
				Name:  "originalpath",
				Usage: "the jdk download index file url.",
				Value: DefaultOriginalpath,
			},
		},
		Action: func(c *cli.Context) error {
			if c.IsSet("java_home") || config.JavaHome == "" {
				config.JavaHome = c.String("java_home")
			}
			cmd := exec.Command("cmd", "/C", "setx", "JAVA_HOME", config.JavaHome, "/M")
			err := cmd.Run()
			if err != nil {
				return errors.New("set Environment variable `JAVA_HOME` failure: Please run as admin user")
			}
			fmt.Println("set `JAVA_HOME` Environment variable to ", config.JavaHome)

			if c.IsSet("originalpath") || config.Originalpath == "" {
				config.Originalpath = c.String("originalpath")
			}
			path := fmt.Sprintf(`%s/bin;%s;%s`, config.JavaHome, os.Getenv("PATH"), file.GetCurrentPath())
			cmd = exec.Command("cmd", "/C", "setx", "path", path, "/m")
			err = cmd.Run()
			if err != nil {
				return errors.New("set Environment variable `PATH` failure: Please run as admin user")
			}
			fmt.Println("add jvms.exe to `path` Environment variable")
			return nil
		},
	}

}
