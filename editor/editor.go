package editor

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

// OpenEditor creates a new file with the input and returns the contents
// after closing and deleting the new file
func OpenEditor(input string) ([]byte, error) {
	editor := os.Getenv("EDITOR")
	// create a temp file
	tmpDir := os.TempDir()
	tmpFile, tmpFileErr := ioutil.TempFile(tmpDir, "prtitle")
	if tmpFileErr != nil {
		color.Red("Error %s while creating tempFile", tmpFileErr)
		os.Exit(0)
	}
	// see if the editor exists
	path, err := exec.LookPath(editor)
	if err != nil {
		color.Red("Error %s while looking for %s\n", path, editor)
		os.Exit(0)
	}
	// write the input to the file
	inputBytes := []byte(input)
	if err := ioutil.WriteFile(tmpFile.Name(), inputBytes, 0644); err != nil {
		color.Red("Error while writing to file : %s\n", err)
		os.Exit(0)
	}

	cmd := exec.Command(path, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// open the file in the editor
	err = cmd.Start()
	if err != nil {
		return []byte(``), fmt.Errorf("Editor execution failed: %s", err)
	}
	err = cmd.Wait()
	if err != nil {
		color.Red("Command finished with error: %v\n", err)
		os.Exit(0)
	}

	// read from file
	fileContent, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		color.Red("Error while Reading: %s\n", err)
		os.Exit(0)

	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			color.Red("Error while deleting the tmp file")
			os.Exit(0)

		}
	}()
	return fileContent, nil
}
