package command

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/google/subcommands"
	"github.com/morikuni/aec"
)

var imageDir string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	imageDir = filepath.Join(homeDir, ".whc", "images")
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		log.Fatal(err)
	}
}

func getImagePath(tokenID string) (string, error) {
	filePath := filepath.Join(imageDir, tokenID+".jpg")
	stat, err := os.Stat(filePath)
	if err == nil && stat.Size() > 0 {
		return filePath, nil
	}
	res, err := http.Get(fmt.Sprintf("https://hackers-metadata.herokuapp.com/image/%s.jpg", tokenID))
	if err != nil {
		return "", fmt.Errorf("failed to fetch hacker: %s", err)
	}
	defer res.Body.Close()
	image, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read hacker: %s", err)
	}
	if err := ioutil.WriteFile(filePath, image, 0644); err != nil {
		return "", fmt.Errorf("failed to keep hacker: %s", err)
	}
	return filePath, nil
}

type AsciiCmd struct {
	width int
}

func (*AsciiCmd) Name() string     { return "ascii" }
func (*AsciiCmd) Synopsis() string { return "show hacker" }
func (*AsciiCmd) Usage() string {
	return "whc ascii [-w <width>] token_id"
}
func (c *AsciiCmd) SetFlags(f *flag.FlagSet) {
	f.IntVar(&c.width, "width", 0, "image width")
	f.IntVar(&c.width, "w", 0, "image width (short hand)")
}

func (c *AsciiCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	tokenIDs := f.Args()
	if len(tokenIDs) < 1 {
		f.Usage()
		return subcommands.ExitFailure
	}
	flags := aic_package.DefaultFlags()
	flags.Colored = true
	if c.width > 0 {
		flags.Width = c.width
	}

	height := 0
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		for _, tokenID := range tokenIDs {
			imagePath, err := getImagePath(tokenID)
			if err != nil {
				fmt.Println(err)
				return subcommands.ExitFailure
			}
			art, err := aic_package.Convert(imagePath, flags)
			if err != nil {
				fmt.Println("\033[0;31mDon't fake the token_id!!\033[0m")
				return subcommands.ExitFailure
			}
			fmt.Println(art)
			height = strings.Count(art, "\n") + 1
			fmt.Print(aec.Up(uint(height)))
			select {
			case <-t.C:
			case <-ctx.Done():
				return subcommands.ExitSuccess
			}
		}
	}
}
