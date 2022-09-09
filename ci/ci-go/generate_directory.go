package main

import (
	_ "embed"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
	"sort"
	"strings"
	"time"
)

//go:embed directory_template.html
var template []byte

// file holds information relevant information on a file for generating the directory.
type file struct {
	// Name is the name of the file.
	Name string
	// CreationTime is the time the file was first added to the git repository.
	CreationTime time.Time
	// UpdatedTime is the time of the last update to the file.
	UpdatedTime time.Time
}

func main() {
	// Open the git repository.
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		panic("no git repository in current working directory")
	}

	// Get the current HEAD.
	head, err := repo.Head()
	if err != nil {
		panic(err)
	}

	// Get HEAD commit.
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		panic(err)
	}

	var files []file
	// Iterate over the files in the tree.
	iter, err := commit.Files()
	iter.ForEach(func(f *object.File) error {
		// Get commit history of the file.
		commits, _ := repo.Log(&git.LogOptions{FileName: &f.Name})

		var updatedTime time.Time
		var creationTime time.Time
		// repo.Log returns a CommitIter; we must iterate over all commits to get the earliest.
		commits.ForEach(func(commit *object.Commit) error {
			// Only set updatedTime if not already set.
			if updatedTime.Equal(time.Time{}) {
				updatedTime = commit.Committer.When
			}
			creationTime = commit.Committer.When
			return nil
		})

		// Add file info to `files` slice.
		files = append(files, file{
			Name:         f.Name,
			CreationTime: creationTime,
			UpdatedTime:  updatedTime,
		})
		return nil
	})

	// Sort `files` slice.
	sort.SliceStable(files, func(i, j int) bool {
		fileI := files[i]
		fileJ := files[j]
		// If both were introduced in the same commit, sort by alphabetical order.
		if fileI.CreationTime == fileJ.CreationTime {
			return strings.Compare(fileI.Name, fileJ.Name) > 1
		}
		// Sort by creation time.
		return fileI.CreationTime.After(fileJ.CreationTime)
	})

	pages := "\n"
	// Iterate over each file.
	for _, f := range files {
		// Filter out non-HTML files, templates, and the directory file.
		if !strings.HasSuffix(f.Name, ".html") || strings.Contains(f.Name, "_template") || f.Name == "index.html" {
			continue
		}
		// Append HTML.
		pages += "\t<li><a href=\"" + f.Name + "\">" + f.Name + "</a> (Created " + f.CreationTime.Format("January 02 2006") + ", Updated " + f.UpdatedTime.Format("January 02 2006") + ")</li>\n"
	}
	// Replace placeholder with HTML.
	content := strings.Replace(string(template), "<!-- Pages go here! -->", pages, 1)
	// Write to index.html
	os.WriteFile("index.html", []byte(content), 06666)
}
