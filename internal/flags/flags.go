package flags

import "github.com/urfave/cli/v2"

var Test = []cli.Flag{
	&cli.StringFlag{
		Name: "test.testlogfile",
	},
	&cli.StringFlag{
		Name: "test.paniconexit0",
	},
	&cli.StringFlag{
		Name: "test.timeout",
	},
	&cli.StringFlag{
		Name: "test.run",
	},
	&cli.StringFlag{
		Name: "test.coverprofile",
	},
	&cli.StringFlag{
		Name: "test.v",
	},
}
