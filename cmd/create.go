package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Tanuki application or service",
	Long:  `Create a Tanuki applicaion or service and sets up the necessary structure and libraries`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		create(args[0])
	},
}

func create(path string) {
	fmt.Println("Creating Tanuki application at", path)
	executable, _ := os.Executable()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.FileMode(0000755))
		os.Mkdir(join(path, "/bin"), os.FileMode(0000755))
		os.Mkdir(join(path, "/listener"), os.FileMode(0000755))
		os.Mkdir(join(path, "/static"), os.FileMode(0000755))
		copy(executable, join(path, "/tanuki"))
	}
}

func copy(src, dst string) (int64, error) {
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	if err != nil {
		return 0, err
	}

	err = os.Chmod(dst, os.FileMode(0000755))
	return nBytes, err
}
