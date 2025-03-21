package gravatar

import (
	"crypto/md5" //nolint:gosec
	"encoding/hex"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

// Default values for gravatar options generated by this package.
const (
	defaultSize   = 80
	defaultImage  = "robohash"
	defaultRating = "pg"
)

var (
	baseURL, _ = url.Parse("https://www.gravatar.com")
)

// New returns a Gravatar image URL for the given email address
func New(email string, opts *Options) string {
	if opts == nil {
		opts = &Options{Size: defaultSize, DefaultImage: defaultImage, Rating: defaultRating}
	}

	img := Hash(email)
	if opts.FileExtension != "" {
		img += opts.FileExtension
	}

	link := baseURL.ResolveReference(&url.URL{Path: filepath.Join("avatar", img)})

	params := make(url.Values)
	if opts.Size != 0 {
		params.Set("s", strconv.Itoa(opts.Size))
	}

	if opts.DefaultImage != "" {
		params.Set("d", opts.DefaultImage)
	}

	if opts.ForceDefault {
		params.Set("f", "y")
	}

	if opts.Rating != "" {
		params.Set("r", opts.Rating)
	}

	link.RawQuery = params.Encode()

	return link.String()
}

// Hash returns the Gravatar email MD5 hex encoded hash as defined in:
// https://en.gravatar.com/site/implement/hash/
func Hash(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	sum := md5.Sum([]byte(email)) // nolint: gosec

	return hex.EncodeToString(sum[:])
}

type Options struct {
	// Size of the image square; request images from 1px up to 2048px.
	Size int

	// DefaultImage is can be of 404, mp, identicon, monsterid, wavatar, retro, robohash, or blank.
	DefaultImage string

	// ForceDefault image to always load
	ForceDefault bool

	// Rating indicates image appropriateness, one of g, pg, r, or x.
	Rating string

	// FileExtension is optional, can be one of .png, .jpg, etc.
	FileExtension string
}
