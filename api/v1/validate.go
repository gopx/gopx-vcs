package v1

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"net/url"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

const (
	usernameMaxLength    = 300
	packageNameMaxLength = 300
)

// ValidatePackageVersion checks whether the package version
// is in valid format according to the semvar spec.
// Ref: https://semver.org/
func ValidatePackageVersion(version string) error {
	_, err := semver.NewVersion(version)
	return err
}

// ValidatePackageName checks whether the package name
// meets the naming constraints.
func ValidatePackageName(packageName string) error {
	ln := utf8.RuneCountInString(packageName)
	if !(ln > 0 && ln < packageNameMaxLength) {
		return errors.Errorf("Package name must be non-empty and maximum %d characters long", packageNameMaxLength)
	}

	if url.PathEscape(packageName) != packageName {
		return errors.Errorf("Package name must not contain any non-url-safe character")
	}

	if strings.ToLower(packageName) != packageName {
		return errors.Errorf("Package name must contain only lowercase charactars")
	}

	if strings.HasPrefix(packageName, ".") {
		return errors.Errorf("Package name must not start with %c character", '.')
	}

	if matched, err := regexp.MatchString(`[@~\\/!'()*\s]`, packageName); err != nil {
		return errors.Errorf("Package name couldn't be checked with regex")
	} else if matched {
		return errors.Errorf(`Package name must not contain any of these special characters: @, ~, \, /, !, ', (, ), *`)
	}

	return nil
}

func validateUsername(username string) error {
	ln := utf8.RuneCountInString(username)
	if !(ln > 0 && ln < usernameMaxLength) {
		return errors.Errorf("Username must be non-empty and maximum %d characters long", usernameMaxLength)
	}

	if url.PathEscape(username) != username {
		return errors.Errorf("Username must not contain any non-url-safe character")
	}

	if strings.ToLower(username) != username {
		return errors.Errorf("Username must contain only lowercase charactars")
	}

	if strings.HasPrefix(username, ".") {
		return errors.Errorf("Username must not start with %c character", '.')
	}

	if matched, err := regexp.MatchString(`[@~\\/!'()*\s]`, username); err != nil {
		return errors.Errorf("Username couldn't be checked with regex")
	} else if matched {
		return errors.Errorf(`Username must not contain any of these special characters: @, ~, \, /, !, ', (, ), *`)
	}

	return nil
}

func validatePackageOwnerInfo(owner *PackageOwner) error {
	if owner.Name == "" {
		return errors.Errorf("Package Owner's name can't be empty")
	}

	if owner.Username == "" {
		return errors.Errorf("Package Owner's username can't be empty")
	}
	return nil
}
