package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"sort"
	"strings"
	"time"
)

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
	repo, err := git.PlainOpen(".")
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
			UpdatedTime:  time.Time{},
		})
		return nil
	})

	// Sort `files` slice.
	sort.Slice(files, func(i, j int) bool {
		fileI := files[i]
		fileJ := files[j]
		// If both were introduced in the same commit, sort by alphabetical order.
		if fileI.CreationTime == fileJ.CreationTime {
			return strings.Compare(fileI.Name, fileJ.Name) > 1
		}
		// Sort by creation time.
		return fileI.CreationTime.After(fileJ.CreationTime)
	})
}