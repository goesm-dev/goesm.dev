package site

import (
	"strings"

	"goesm.dev/press"
)

// Version says which release of goesm (or which commit of gosfc) the docs
// of a site describe: the version go.mod selects, which cmd/syncdocs
// synced the docs from.
type Version struct {
	// Label is the version as the navigation bar shows it: the module
	// version, as pkg.go.dev shows it, prereleases and pseudo-versions
	// included ("v0.0.1-beta.3", "v0.0.0-20261005171715-4f0689be9e96").
	Label string `json:"label"`
	// Title says what Label is ("Docs for goesm v0.0.1-beta.3").
	Title string `json:"title"`
	// Heading names the menu's links ("Version"), in the narrow-screen menu.
	Heading string          `json:"heading"`
	Links   []press.NavLink `json:"links"`
}

// versionOf describes the version the docs of site s describe, in the
// language of locale code.
func versionOf(s *press.Site, code string) Version {
	ja := code == "ja"
	pick := func(en, jaText string) string {
		if ja {
			return jaText
		}
		return en
	}
	heading := pick("Version", "バージョン")
	name, repo, version, rev := "goesm", "https://github.com/goesm-dev/goesm", goesmVersion, goesmRef
	if s == gosfcSite {
		name, repo, version, rev = "gosfc", "https://github.com/goesm-dev/gosfc", gosfcVersion, gosfcRef
	}
	pkg := press.NavLink{Text: "pkg.go.dev", Link: "https://pkg.go.dev/" + strings.TrimPrefix(repo, "https://") + "@" + version, External: true}
	v := Version{
		Label:   version,
		Title:   pick("Docs for "+name+" "+version, name+" "+version+" のドキュメント"),
		Heading: heading,
	}
	if commit, date := pseudo(version); date != "" {
		// No release yet (gosfc): go.mod selects a commit.
		v.Links = []press.NavLink{
			{Text: pick("Commit "+commit+" ("+date+")", "コミット "+commit+"（"+date+"）"), Link: repo + "/commit/" + rev, External: true},
			pkg,
			{Text: pick("Commit history", "コミット履歴"), Link: repo + "/commits/main", External: true},
		}
		return v
	}
	v.Links = []press.NavLink{
		{Text: pick("Release notes", "リリースノート"), Link: repo + "/releases/tag/" + version, External: true},
		pkg,
		{Text: pick("All releases", "すべてのリリース"), Link: repo + "/releases", External: true},
	}
	return v
}

// pseudo returns the short commit and the date (YYYY-MM-DD, UTC) of a
// pseudo-version ("v0.0.0-20261005171715-4f0689be9e96"), or the version
// itself when it is not one.
func pseudo(v string) (commit, date string) {
	parts := strings.Split(v, "-")
	if len(parts) < 3 {
		return v, ""
	}
	rev, ts := parts[len(parts)-1], parts[len(parts)-2]
	// A prerelease base adds a ".0." before the timestamp
	// (v1.2.4-pre.0.20261005171715-4f0689be9e96).
	ts = ts[strings.LastIndex(ts, ".")+1:]
	if len(rev) != 12 || len(ts) != 14 {
		return v, ""
	}
	return rev[:7], ts[:4] + "-" + ts[4:6] + "-" + ts[6:8]
}
