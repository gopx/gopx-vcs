package v1

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"io/ioutil"

	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopx.io/gopx-vcs/pkg/utils"
)

// registerPackage registers a new package to the VCS repository storage.
// It creates a new package if it does not exist.
// If the new version of target package already exists then it throws errors, otherwise
// registers the package with the new version.
func registerPackage(meta *PackageMeta, data io.Reader) error {
	switch meta.Type {
	case PackageTypePublic:
		return registerPublicPackage(meta, data)
	case PackageTypePrivate:
		return registerPrivatePackage(meta, data)
	default:
		return errors.Errorf("Package type %d is not supported yet!", int(meta.Type))
	}
}

func registerPublicPackage(meta *PackageMeta, data io.Reader) error {
	pkgName := meta.Name
	pkgVersion := meta.Version
	owner := &meta.Owner

	err := validatePackageOwnerInfo(owner)
	if err != nil {
		return err
	}

	err = ValidatePackageName(pkgName)
	if err != nil {
		return err
	}

	err = ValidatePackageVersion(pkgVersion)
	if err != nil {
		return err
	}

	err = resolvePackageRepo(pkgName)
	if err != nil {
		return errors.Wrapf(err, "Couldn't resolve the package: %s", pkgName)
	}

	verExists, err := packageVersionExists(pkgName, pkgVersion)
	if err != nil {
		return errors.Wrapf(err, "Couldn't check package version: %s", pkgName)
	}

	if verExists {
		return errors.Errorf("Package version already exists: %s@%s", pkgName, pkgVersion)
	}

	opsDir, err := tempPackageRepoOpsDir(pkgName)
	if err != nil {
		return errors.Wrapf(err, "Unable to create new temp dir for repo operations")
	}
	defer os.RemoveAll(opsDir)

	cOpt, err := packageRepoCloneOptions(pkgName)
	if err != nil {
		return errors.Wrapf(err, "Unable to create package repo clone options: %s", pkgName)
	}

	repo, err := git.PlainClone(opsDir, false, cOpt)
	if err != nil && err != transport.ErrEmptyRemoteRepository {
		return errors.Wrapf(err, "Couldn't clone the package repo: %s", pkgName)
	}

	files, err := ioutil.ReadDir(opsDir)
	if err != nil {
		return errors.Wrapf(err, "Unable to access operaion dir for package: %s", pkgName)
	}

	for _, f := range files {
		if f.Name() == ".git" {
			continue
		}
		err = os.RemoveAll(filepath.Join(opsDir, f.Name()))
		if err != nil {
			return errors.Wrapf(err, "Unable to delete package files from opsdir: %s", pkgName)
		}
	}

	wt, err := repo.Worktree()
	if err != nil {
		return errors.Wrapf(err, "Unable to access repo worktree for package: %s", pkgName)
	}

	_, err = utils.DecompressTarGZ(opsDir, data)
	if err != nil {
		return errors.Wrapf(err, "Couldn't extract package data into worktree")
	}

	files, err = ioutil.ReadDir(opsDir)
	if err != nil {
		return errors.Wrapf(err, "Unable to access operaion dir for package: %s", pkgName)
	}

	for _, f := range files {
		if f.Name() == ".git" {
			continue
		}
		_, err := wt.Add(f.Name())
		if err != nil {
			return errors.Wrapf(err, "Unable to stage files into git index: %s", pkgName)
		}
	}

	pkgTagName, err := tagNameFromVersion(pkgVersion)
	if err != nil {
		return errors.Wrapf(err, "Couldn't generate tag name form package version: %s of %s", pkgVersion, pkgName)
	}

	_, err = wt.Commit(
		vcsRepoCommitMessage(pkgTagName),
		&git.CommitOptions{
			All:       true,
			Author:    vcsRepoAuthorSignature(owner),
			Committer: vcsRepoCommitterSignature(),
		},
	)
	if err != nil {
		return errors.Wrapf(err, "Couldn't commit files of package version: %s of %s", pkgVersion, pkgName)
	}

	err = vcsRepoCreateTag(
		pkgTagName,
		vcsRepoTaggerSignature(),
		vcsRepoTagMessage(pkgTagName),
		repo,
	)
	if err != nil {
		return errors.Wrapf(err, "Unable to create tag for the package: %s of %s", pkgTagName, pkgName)
	}

	tagRefSpec := config.RefSpec(fmt.Sprintf("refs/tags/%s:refs/tags/%s", pkgTagName, pkgTagName))
	masterRefSpec := config.RefSpec("+refs/heads/master:refs/heads/master")
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{tagRefSpec, masterRefSpec},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return errors.Wrapf(err, "Unable to push package to the repo: %s", pkgName)
	}

	err = exportPackageRepo(pkgName)
	if err != nil {
		return errors.Wrapf(err, "Unable to export package: %s", pkgName)
	}

	return nil
}

func registerPrivatePackage(meta *PackageMeta, data io.Reader) error {
	return errors.Errorf("Sorry, private package is not supported yet!")
}

func splitPackageName(packageName string) (pkg, username string, isScoped bool) {
	if !strings.HasPrefix(packageName, "@") {
		return packageName, "", false
	}

	packageName = strings.TrimPrefix(packageName, "@")
	idx := strings.Index(packageName, "/")
	if idx == -1 {
		return "", packageName, true
	}
	return packageName[idx+1:], packageName[:idx], true
}
