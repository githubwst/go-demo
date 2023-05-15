package alfred

import (
	ga "github.com/jason0x43/go-alfred"
	"os"
)

var gwf ga.Workflow

func init() {
	workDir, _ := os.Getwd()
	gwf, _ = ga.OpenWorkflow(workDir, false)
}

// Your workflow starts here
func run() {
	// Add a "Script Filter" result

	// Send results to Alfred
	//gwf.SendToAlfred()
}

func main() {
	// Wrap your entry point with Run() to catch and log panics and
	// show an error in Alfred instead of silently dying

}
