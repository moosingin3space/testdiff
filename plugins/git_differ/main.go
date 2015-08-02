package git_differ

import (
	git "github.com/libgit2/git2go"
	"github.com/natefinch/pie"
	"log"
	"net/rpc/jsonrpc"
)

type Differ struct{}

func (d Differ) Diff(workingDir string) []string {
	var repo *git.Respository
	var index *git.Index
	var diff *git.Diff
	var differentFiles []string
	var err error

	if repo, err = git.OpenRepository(workingDir); err != nil {
		log.Fatal(err)
		return nil
	}

	if index, err = git.OpenIndex(workingDir); err != nil {
		log.Fatal(err)
		return nil
	}

	if diff, err = repo.DiffIndexToWorkdir(index, nil); err != nil {
		log.Fatal(err)
		return nil
	}

	// TODO actually do diff
}

func main() {
	var differ Differ
	server := pie.NewProvider()
	server.Register(differ)
	server.Serve()
}
