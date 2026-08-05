package utils

import "github.com/gosimple/slug"

// ToSlug converts a string to a URL-friendly slug
func ToSlug(text string) string {
	return slug.Make(text)
}
