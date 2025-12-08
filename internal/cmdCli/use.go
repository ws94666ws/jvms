package cmdCli

import (
	"github.com/codegangsta/cli"
	"github.com/ystyle/jvms/internal/entity"
)

func use(config *entity.Config) *cli.Command {
	cmd := &cli.Command{
		Name:      "use",
		ShortName: "u",
		Usage:     "Switch to use the specified version or index number.",
		Flags:     switchFlags,
		Action:    switchFunc(config),
	}
	return cmd
}
