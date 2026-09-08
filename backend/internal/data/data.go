package data

import "embed"

//go:embed profile.json skills.json projects.json experience.json
var files embed.FS

// Store holds each content file as raw JSON bytes so handlers can write them
// straight to the response without an unmarshal/marshal round trip.
type Store struct {
	Profile    []byte
	Skills     []byte
	Projects   []byte
	Experience []byte
}

func Load() (*Store, error) {
	profile, err := files.ReadFile("profile.json")
	if err != nil {
		return nil, err
	}
	skills, err := files.ReadFile("skills.json")
	if err != nil {
		return nil, err
	}
	projects, err := files.ReadFile("projects.json")
	if err != nil {
		return nil, err
	}
	experience, err := files.ReadFile("experience.json")
	if err != nil {
		return nil, err
	}

	return &Store{
		Profile:    profile,
		Skills:     skills,
		Projects:   projects,
		Experience: experience,
	}, nil
}
