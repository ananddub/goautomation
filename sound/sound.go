package sound

import (
	"fmt"
	"os/exec"
	"runtime"
)

func PlayBeep() {
	return
	switch runtime.GOOS {
	case "darwin": // macOS
		path := "sound/mixkit-sport-start-bleeps-918.wav"
		exec.Command("afplay", path).Run()

	case "linux":
		exec.Command("speaker-test", "-t", "sine", "-f", "800", "-l", "1").Run()
	case "windows":
		exec.Command("rundll32", "user32.dll,MessageBeep", "0x00000040").Run()
	default:
		for i := 0; i < 4; i++ {
			fmt.Print('\a')
		}
	}
}
