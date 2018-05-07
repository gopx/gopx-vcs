package v1

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// PackageType represents type of GoPX package.
type PackageType int

func (p PackageType) String() string {
	switch p {
	case 0:
		return "public"
	case 1:
		return "private"
	default:
		return "unknown"
	}
}

const (
	// PackageTypePublic indicates the package is public.
	PackageTypePublic PackageType = iota
	// PackageTypePrivate indicates the package is private,
	// requires Authentication to access it.
	PackageTypePrivate
)

// PackageMeta represents the metadata of a GoPX package.
type PackageMeta struct {
	Type    PackageType  `json:"type"`
	Name    string       `json:"name"`
	Version string       `json:"version"`
	Owner   PackageOwner `json:"owner"`
}

// PackageOwner represents owner of the package.
type PackageOwner struct {
	Name        string `json:"name"`
	PublicEmail string `json:"publicEmail"`
	Username    string `json:"username"`
}

// ParsePackageMeta parses the string meta and
// returns the equivalent PackageMeta representation.
func ParsePackageMeta(meta string) (*PackageMeta, error) {
	pm := new(PackageMeta)
	err := json.Unmarshal([]byte(meta), pm)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to parse package meta")
	}
	return pm, nil
}
