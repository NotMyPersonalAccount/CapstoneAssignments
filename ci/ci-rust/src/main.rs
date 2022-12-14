use chrono::{DateTime, TimeZone, Utc};
use git2::{Oid, Repository, RepositoryOpenFlags};
use std::collections::HashMap;
use std::fs::File;
use std::io::Write;

// FileInfo holds information relevant information on a file for generating the directory.
struct FileInfo {
    name: String,
    created_time: DateTime<Utc>,
    updated_time: DateTime<Utc>,
}

fn main() {
    // Open the git repository.
    let repo = Repository::open_ext(".", RepositoryOpenFlags::CROSS_FS, vec!["/"]).unwrap();

    // Get the index to list out files in the repository.
    let index = repo.index().unwrap();

    let mut initial_commits: HashMap<String, Oid> = HashMap::new();
    let mut latest_commits: HashMap<String, Oid> = HashMap::new();
    // Iterate over all commits in the repository, latest commits first.
    let mut revwalk = repo.revwalk().unwrap();
    revwalk.push_head().unwrap();
    revwalk.set_sorting(git2::Sort::TIME).unwrap();
    for commit in revwalk {
        // Get current tree.
        let commit = commit.unwrap();
        let commit = repo.find_commit(commit).unwrap();
        let tree = commit.tree().unwrap();

        // Get the diff between the current tree and the parent tree.
        let parent = commit.parent(0);
        let diff = match parent {
            Ok(parent) => repo.diff_tree_to_tree(Some(&parent.tree().unwrap()), Some(&tree), None),
            Err(_) => repo.diff_tree_to_tree(None, Some(&tree), None),
        };
        let diff = diff.unwrap();

        // Iterate over each diff.
        diff.foreach(
            &mut |delta, _| {
                let path = delta.new_file().path().unwrap();
                // Set initial commit for the changed file to the current commit.
                initial_commits.insert(path.display().to_string(), commit.id());
                // Set latest commit for the changed file to the current commit if not already set.
                latest_commits
                    .entry(path.display().to_string())
                    .or_insert(commit.id());
                true
            },
            None,
            None,
            None,
        )
        .unwrap();
    }

    let mut files: Vec<FileInfo> = Vec::new();
    // Iterate over files in the repository.
    index.iter().for_each(|file| {
        let name = String::from_utf8(file.path).unwrap();
        let initial_commit = repo
            .find_commit(*initial_commits.get(&name).unwrap())
            .unwrap();
        let latest_commit = repo
            .find_commit(*latest_commits.get(&name).unwrap())
            .unwrap();
        // Add file info to files vec.
        files.push(FileInfo {
            name,
            created_time: Utc.timestamp(initial_commit.committer().when().seconds(), 0),
            updated_time: Utc.timestamp(latest_commit.committer().when().seconds(), 0),
        });
    });

    // Sort `files` vec.
    files.sort_by(|a, b| {
        // If both were introduced in the same commit, sort by alphabetical order.
        if a.created_time.eq(&b.created_time) {
            a.name.cmp(&b.name)
        } else {
            // Sort by creation time.
            a.created_time.cmp(&b.created_time).reverse()
        }
    });

    // Get template from file.
    let template = std::fs::read_to_string("ci/directory_template.html").unwrap();
    let mut pages = String::from("\n");
    // Iterate over each file.
    files.iter().for_each(|file| {
        // Filter out non-HTML files, templates, and the directory file.
        if !file.name.ends_with(".html")
            || file.name.ends_with("directory_template.html")
            || file.name.eq("index.html")
        {
            return;
        }
        // Append HTML.
        pages += &format!(
            "\t<li><a href=\"{}\">{}</a> (Created {}, Updated {})</li>\n",
            file.name,
            file.name.replace("pages/", ""),
            file.created_time.format("%B %d %Y"),
            file.updated_time.format("%B %d %Y")
        )
        .to_string();
    });
    // Write to index.html with the generated HTML.
    File::create("index.html")
        .unwrap()
        .write_all(
            template
                .replace("<!-- Pages go here! -->", &pages)
                .as_bytes(),
        )
        .unwrap();
}
