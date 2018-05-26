package v1

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopx.io/gopx-vcs/pkg/config"
	"gopx.io/gopx-vcs/pkg/utils"
)

const (
	// repoCommitterName represents the commiter name for auto generated
	// commits.
	repoCommitterName = "GoPX"

	// repoCommitterEmail represents the commiter email for auto generated
	// commits.
	repoCommitterEmail = "gopx@gopx.io"

	// repoTaggerName represents the tagger name for auto generated
	// tags.
	repoTaggerName = "GoPX"

	// repoTaggerEmail represents the tagger email for auto generated
	// tags.
	repoTaggerEmail = "gopx@gopx.io"
)

const gitExportRepoFileName = "git-daemon-export-ok"

func packageExists(packageName string) (bool, error) {
	rPath, err := packageRepoPath(packageName)
	if err != nil {
		return false, errors.Wrapf(err, "Unable to calculate package repo path: %s", packageName)
	}
	return utils.Exists(rPath)
}

func initPackageRepo(packageName string) (*git.Repository, error) {
	rPath, err := packageRepoPath(packageName)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to calculate package repo path: %s", packageName)
	}
	return git.PlainInit(rPath, true)
}

func tempPackageRepoOpsDir(packageName string) (string, error) {
	prefix := fmt.Sprintf("gopx-package-repo-%s-", packageName)
	return ioutil.TempDir("", prefix)
}

func packageVersionExists(packageName, version string) (bool, error) {
	checkVer, err := semver.NewVersion(version)
	if err != nil {
		return false, errors.Wrapf(err, "Invalid version to check: %s", version)
	}

	rPath, err := packageRepoPath(packageName)
	if err != nil {
		return false, errors.Wrapf(err, "Unable to calculate package repo path: %s", packageName)
	}

	repo, err := git.PlainOpen(rPath)
	if err != nil {
		return false, errors.Wrapf(err, "Package repo couldn't be opened: %s", packageName)
	}

	tagIter, err := repo.TagObjects()
	if err != nil {
		return false, errors.Wrapf(err, "Couldn't access tag Objects: %s", packageName)
	}

	for {
		tag, err := tagIter.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return false, errors.Wrapf(err, "Couldn't access tag Objects: %s", packageName)
			}
		}

		tagVer, err := semver.NewVersion(tag.Name)
		if err != nil {
			return false, errors.Wrapf(err, "Invalid tag version: %s of package %s", tag.Name, packageName)
		}

		if checkVer.Equal(tagVer) {
			return true, nil
		}

	}

	return false, nil
}

func repoName(packageName string) string {
	return fmt.Sprintf("%s%s", packageName, config.VCS.RepoExt)
}

func vcsRepoAuthorSignature(pkgOwner *PackageOwner) *object.Signature {
	return &object.Signature{
		Name:  fmt.Sprintf("%s(%s)", pkgOwner.Name, pkgOwner.Username),
		Email: pkgOwner.PublicEmail,
		When:  time.Now(),
	}
}

func vcsRepoCommitterSignature() *object.Signature {
	return &object.Signature{
		Name:  repoCommitterName,
		Email: repoCommitterEmail,
		When:  time.Now(),
	}
}

func vcsRepoTaggerSignature() *object.Signature {
	return &object.Signature{
		Name:  repoTaggerName,
		Email: repoTaggerEmail,
		When:  time.Now(),
	}
}

func vcsRepoCommitMessage(tagName string) string {
	return fmt.Sprintf("Update package to version %s", tagName)
}

func vcsRepoTagMessage(tagName string) string {
	return fmt.Sprintf("Released %s", tagName)
}

func vcsRepoCreateTag(name string, tagger *object.Signature, message string, repo *git.Repository) error {
	wt, err := repo.Worktree()
	if err != nil {
		return errors.Wrapf(err, "Couldn't access the worktree")
	}

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.Master,
	})
	if err != nil {
		return errors.Wrapf(err, "Couldn't checkout master branch")
	}

	headRef, err := repo.Head()
	if err != nil {
		return errors.Wrapf(err, "Couldn't access HEAD reference")
	}

	tag := object.Tag{
		Name:       name,
		Tagger:     *tagger,
		Message:    message,
		TargetType: plumbing.CommitObject,
		Target:     headRef.Hash(),
	}

	enObj := repo.Storer.NewEncodedObject()
	tag.Encode(enObj)

	hash, err := repo.Storer.SetEncodedObject(enObj)
	if err != nil {
		return errors.Wrapf(err, "Couldn't set encoded object for tag")
	}

	tagRefName := fmt.Sprintf("refs/tags/%s", name)
	ref := plumbing.NewReferenceFromStrings(tagRefName, hash.String())

	err = repo.Storer.SetReference(ref)
	if err != nil {
		return errors.Wrapf(err, "Couldn't set reference for tag")
	}

	return nil
}

func tagNameFromVersion(pkgVersion string) (string, error) {
	ver, err := semver.NewVersion(pkgVersion)
	if err != nil {
		return "", err
	}

	tagName := fmt.Sprintf(
		"v%d.%d.%d",
		ver.Major(),
		ver.Minor(),
		ver.Patch(),
	)
	if ver.Prerelease() != "" {
		tagName = fmt.Sprintf("%s-%s", tagName, ver.Prerelease())
	}

	return tagName, nil
}

func packageRepoPath(packageName string) (string, error) {
	rr := config.VCS.RepoRoot
	repoName := repoName(packageName)
	rPath := filepath.Join(rr, repoName)

	rPath, err := filepath.Abs(rPath)
	if err != nil {
		return "", err
	}

	return rPath, nil
}

func exportPackageRepo(packageName string) error {
	rPath, err := packageRepoPath(packageName)
	if err != nil {
		return errors.Wrapf(err, "Unable to calculate package repo path: %s", packageName)
	}

	exportFile := filepath.Join(rPath, gitExportRepoFileName)
	file, err := os.Create(exportFile)
	if err != nil {
		return errors.Wrapf(err, "Unable to create %s file: %s", gitExportRepoFileName, packageName)
	}

	file.Close()

	return nil
}

func resolvePackageRepo(packageName string) error {
	rPath, err := packageRepoPath(packageName)
	if err != nil {
		return err
	}

	_, err = git.PlainOpen(rPath)
	if err == nil {
		return nil
	}

	if err != git.ErrRepositoryNotExists {
		return err
	}

	if exists, err := utils.Exists(rPath); err != nil {
		return err
	} else if exists {
		uuid := uuid.NewV4()
		repoCorrPath := fmt.Sprintf("%s-%s.corrupted", rPath, uuid.String())

		err = os.Rename(rPath, repoCorrPath)
		if err != nil {
			return errors.Wrapf(err, "Couldn't create corrupted repo backup for package: %s", packageName)
		}

		os.RemoveAll(filepath.Join(repoCorrPath, gitExportRepoFileName))
	}

	_, err = git.PlainInit(rPath, true)
	if err != nil {
		return errors.Wrapf(err, "Couldn't initialize the package repo: %s", packageName)
	}

	return nil
}

func packageRepoCloneOptions(packageName string) (*git.CloneOptions, error) {
	rPath, err := packageRepoPath(packageName)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to calculate package repo path: %s", packageName)
	}

	cOpt := &git.CloneOptions{
		URL:           rPath,
		RemoteName:    "origin",
		ReferenceName: plumbing.Master,
		SingleBranch:  true,
		Tags:          git.NoTags,
	}

	return cOpt, nil
}
