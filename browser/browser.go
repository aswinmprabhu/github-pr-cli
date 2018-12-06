package browser

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenURLInBrowser opens the given url in a browser regardless
// of the platform
func OpenURLInBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return fmt.Errorf("Failed to open the browser : %v", err)
	}
	fmt.Println("Opening in browser.....")
	return nil
}
