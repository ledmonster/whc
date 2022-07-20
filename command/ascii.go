package command

import (
	"context"
	"flag"
	"fmt"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/google/subcommands"
)

type AsciiCmd struct {
	width int
}

func (*AsciiCmd) Name() string     { return "ascii" }
func (*AsciiCmd) Synopsis() string { return "show hacker" }
func (*AsciiCmd) Usage() string {
	return "whc ascii [-w <width>] token_id"
}
func (c *AsciiCmd) SetFlags(f *flag.FlagSet) {
	f.IntVar(&c.width, "width", 80, "image width")
	f.IntVar(&c.width, "w", 80, "image width (short hand)")
}

func (c *AsciiCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	args := f.Args()
	if len(args) < 1 {
		f.Usage()
		return subcommands.ExitFailure
	}
	flags := aic_package.DefaultFlags()
	flags.Colored = true
	flags.Width = c.width
	fmt.Println(c.width)
	tokenID := args[0]
	asciiArt, err := aic_package.Convert(fmt.Sprintf("https://hackers-metadata.herokuapp.com/image/%s.jpg", tokenID), flags)
	if err != nil {
		fmt.Println("\033[0;31mDon't fake the token_id!!\033[0m")
		return subcommands.ExitFailure
	}
	fmt.Println(asciiArt)
	return subcommands.ExitSuccess
}
