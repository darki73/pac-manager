package prompt

import (
	"github.com/chzyer/readline"
	"os"
)

type noBellStdout struct{}

func (nbs *noBellStdout) Write(p []byte) (n int, err error) {
	if len(p) == 1 && p[0] == readline.CharBell {
		return 0, nil
	}
	return os.Stdout.Write(p)
}

func (nbs *noBellStdout) Close() error {
	return readline.Stdout.Close()
}

var noBellStdoutInstance = &noBellStdout{}
