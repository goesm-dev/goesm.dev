// Command syncdocs copies the documentation of goesm and gosfc into the
// content tree, at the versions go.mod pins, so the site documents the
// compiler it is built with.
//
//	go run ./cmd/syncdocs -goesm ../goesm -gosfc third_party/gosfc
//	go run ./cmd/syncdocs -goesm ../goesm -gosfc third_party/gosfc -check
//
// The -goesm and -gosfc directories are git checkouts of the repositories;
// files are read with `git show <rev>:<path>`, so the checked-out branch does
// not matter, only that the pinned revision has been fetched.
//
// goesm's documents go to the goesm site (content/, routes at /), gosfc's
// to the gosfc site (content/gosfc/, routes at /gosfc/). It writes, in the
// site's content directory:
//
//   - {en,ja}/reference/*.md: whole documents (ARCHITECTURE.md,
//     docs/*.md, ...), with `source:` frontmatter for the "edit" link
//   - {en,ja}/_<repo>/*.md: README sections, included by guide pages
//   - public/repo/<repo>/...: the images those files show
//   - site/sources_gen.go: the revisions and module versions synced
//
// Relative links are rewritten: to the site's own page when the target is a
// synced document (or a README section a guide page includes), else to the
// file on GitHub at the synced revision.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"goesm.dev/press"
)

type repo struct {
	name   string // source name in the site config
	module string // Go module path
	web    string
	dir    string // git checkout
	// content is the content directory of the site documenting the
	// repository ("content/"), and base the route the site is under ("").
	content, base string
	rev           string // git revision (tag or commit)
	ver           string // module version from go.mod
	// docs maps repository paths of the English documents to site routes
	// (without base and locale prefix). The Japanese document is the same path with
	// .ja.md.
	docs map[string]string
	// sections maps README section slugs (English heading) to the route of
	// the guide page that includes them.
	sections map[string]string
	// later maps documents that the version in go.mod does not have yet
	// (they are in docs) to the commit they are synced from instead. Drop
	// an entry once go.mod selects a version that has the document.
	later map[string]string
}

var locales = []struct{ code, prefix, suffix string }{
	{"en", "", ".md"},
	{"ja", "/ja", ".ja.md"},
}

func main() {
	goesmDir := flag.String("goesm", "../goesm", "git checkout of github.com/goesm-dev/goesm")
	gosfcDir := flag.String("gosfc", "third_party/gosfc", "git checkout of github.com/goesm-dev/gosfc")
	check := flag.Bool("check", false, "report files that are out of date instead of writing them")
	flag.Parse()

	repos := []*repo{
		{
			name: "goesm", module: "github.com/goesm-dev/goesm", web: "https://github.com/goesm-dev/goesm", dir: *goesmDir,
			content: "content/", base: "",
			docs: map[string]string{
				"README.md":                   "/guide/",
				"ARCHITECTURE.md":             "/reference/architecture/",
				"docs/example-output.md":      "/reference/example-output/",
				"docs/conformance.md":         "/reference/conformance/",
				"docs/otelc.md":               "/reference/otelc/",
				"docs/gopherjs-comparison.md": "/reference/gopherjs-comparison/",
				"bench/README.md":             "/reference/benchmark/",
				"CONTRIBUTING.md":             "/reference/contributing/",
				"docs/releasing.md":           "/reference/releasing/",
				"docs/js-imports.md":          "/reference/js-imports/",
				"docs/js-exports.md":          "/reference/js-exports/",
				"docs/dom.md":                 "/reference/dom/",
				"compare/README.md":           "/reference/compare/",
				"docs/concurrency.md":         "/reference/concurrency/",
				"docs/use-cases.md":           "/reference/use-cases/",
				"CHANGELOG.md":                "/reference/changelog/",
			},
			// The release notes of v0.0.1-beta.0 to beta.3 were written
			// after beta.3 (goesm #91, #92).
			later: map[string]string{
				"CHANGELOG.md": "65fa6ac4fb8b350ddd82b296fdc02f8d3c0c87af",
			},
			sections: map[string]string{
				"intro":        "/guide/",
				"why-goesm":    "/guide/",
				"how-it-works": "/guide/",
				"status":       "/guide/status/",
				"install":      "/guide/getting-started/",
				"usage":        "/guide/getting-started/",
				"performance":  "/guide/performance/",
			},
		},
		{
			name: "gosfc", module: "github.com/goesm-dev/gosfc", web: "https://github.com/goesm-dev/gosfc", dir: *gosfcDir,
			content: "content/gosfc/", base: "/gosfc",
			docs: map[string]string{
				"README.md":       "/guide/",
				"ARCHITECTURE.md": "/reference/architecture/",
				"CONTRIBUTING.md": "/reference/contributing/",
			},
			sections: map[string]string{
				"intro":                        "/guide/",
				"usage-astro":                  "/guide/getting-started/",
				"writing-the-go-block":         "/guide/go-block/",
				"go-in-astro-files":            "/guide/astro/",
				"importing-go-from-javascript": "/guide/importing-go/",
				"benchmark":                    "/guide/benchmark/",
			},
		},
	}

	out := map[string][]byte{}
	for _, r := range repos {
		if err := r.resolve(); err != nil {
			fatal(err)
		}
		if err := r.sync(out); err != nil {
			fatal(err)
		}
	}
	var gen bytes.Buffer
	gen.WriteString("// Code generated by cmd/syncdocs. DO NOT EDIT.\n\npackage site\n\n// Revisions and module versions of the synced documentation.\nconst (\n")
	for _, r := range repos {
		fmt.Fprintf(&gen, "\t%sRef = %q\n", r.name, r.rev)
		fmt.Fprintf(&gen, "\t%sVersion = %q // %s\n", r.name, r.ver, r.module)
	}
	gen.WriteString(")\n")
	src, err := format.Source(gen.Bytes())
	if err != nil {
		fatal(err)
	}
	out["site/sources_gen.go"] = src

	// Files this tool owns that are no longer produced are removed.
	owned := []string{"public/repo"}
	for _, r := range repos {
		for _, l := range locales {
			owned = append(owned, r.content+l.code+"/_"+r.name)
		}
	}
	var stale []string
	for _, dir := range owned {
		filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
			if err == nil && !d.IsDir() && out[filepath.ToSlash(p)] == nil {
				stale = append(stale, filepath.ToSlash(p))
			}
			return nil
		})
	}
	for _, l := range locales {
		files, _ := filepath.Glob("content/" + l.code + "/reference/*.md")
		more, _ := filepath.Glob("content/*/" + l.code + "/reference/*.md")
		files = append(files, more...)
		for _, f := range files {
			data, _ := os.ReadFile(f)
			if bytes.Contains(data, []byte(generatedMarker)) && out[filepath.ToSlash(f)] == nil {
				stale = append(stale, filepath.ToSlash(f))
			}
		}
	}

	var changed []string
	for _, name := range sortedKeys(out) {
		old, err := os.ReadFile(name)
		if err == nil && bytes.Equal(old, out[name]) {
			continue
		}
		changed = append(changed, name)
		if !*check {
			os.MkdirAll(filepath.Dir(name), 0o755)
			if err := os.WriteFile(name, out[name], 0o644); err != nil {
				fatal(err)
			}
		}
	}
	for _, name := range stale {
		changed = append(changed, name+" (stale)")
		if !*check {
			os.Remove(name)
		}
	}
	if *check && len(changed) > 0 {
		fmt.Fprintf(os.Stderr, "syncdocs: out of date (run go run ./cmd/syncdocs):\n  %s\n", strings.Join(changed, "\n  "))
		os.Exit(1)
	}
	for _, c := range changed {
		fmt.Println(c)
	}
}

const generatedMarker = "<!-- Synced by cmd/syncdocs. Edit the source instead. -->"

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "syncdocs:", err)
	os.Exit(1)
}

func sortedKeys(m map[string][]byte) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// resolve finds the version go.mod selects and the git revision for it.
func (r *repo) resolve() error {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Version}}", r.module)
	b, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("go list -m %s: %v", r.module, err)
	}
	r.ver = strings.TrimSpace(string(b))
	r.rev = r.ver
	// v0.0.0-20261004181437-b7d80aeafa77: a pseudo-version names a commit.
	if parts := strings.Split(r.ver, "-"); len(parts) >= 3 && len(parts[len(parts)-1]) == 12 {
		r.rev = parts[len(parts)-1]
	}
	if _, err := r.show("go.mod"); err != nil {
		return fmt.Errorf("%s: revision %s not found in %s (fetch it first): %v", r.name, r.rev, r.dir, err)
	}
	if len(r.rev) == 12 {
		full, err := exec.Command("git", "-C", r.dir, "rev-parse", r.rev).Output()
		if err == nil {
			r.rev = strings.TrimSpace(string(full))
		}
	}
	return nil
}

func (r *repo) show(p string) ([]byte, error) {
	cmd := exec.Command("git", "-C", r.dir, "show", r.rev+":"+p)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	b, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git show %s:%s: %v %s", r.rev, p, err, strings.TrimSpace(stderr.String()))
	}
	return b, nil
}

func (r *repo) sync(out map[string][]byte) error {
	for _, l := range locales {
		for src, route := range r.docs {
			file := strings.TrimSuffix(src, ".md") + l.suffix
			r := r
			ref := ""
			if rev, ok := r.later[src]; ok {
				at := *r
				at.rev, r, ref = rev, &at, "ref: "+rev+"\n"
			}
			data, err := r.show(file)
			if err != nil {
				return err
			}
			text := string(data)
			if src == "README.md" {
				if err := r.syncReadme(out, l.code, l.prefix, file, text); err != nil {
					return err
				}
				continue
			}
			body, err := r.rewrite(out, l.prefix, file, dropLangSwitch(text))
			if err != nil {
				return err
			}
			name := strings.Trim(strings.TrimPrefix(route, "/reference/"), "/")
			doc := "---\nsource: " + r.name + ":" + file + "\n" + ref + "---\n\n" + generatedMarker + "\n\n" + body
			out[r.content+l.code+"/reference/"+name+".md"] = []byte(doc)
		}
	}
	return nil
}

// syncReadme splits a README into its intro (after the language switch line)
// and ## sections, named by the English heading's slug. The Japanese README
// has the same sections in the same order.
func (r *repo) syncReadme(out map[string][]byte, code, prefix, file, text string) error {
	en, err := r.show("README.md")
	if err != nil {
		return err
	}
	names := []string{"intro"}
	for _, h := range headings(string(en), "## ") {
		names = append(names, press.Slug(h))
	}
	parts := splitSections(text)
	if len(parts) != len(names) {
		return fmt.Errorf("%s: %d sections, README.md has %d", file, len(parts), len(names))
	}
	for i, body := range parts {
		if _, ok := r.sections[names[i]]; !ok {
			continue
		}
		if names[i] == "intro" {
			body = readmeIntro(body)
		}
		body, err := r.rewrite(out, prefix, file, body)
		if err != nil {
			return err
		}
		out[r.content+code+"/_"+r.name+"/readme-"+names[i]+".md"] = []byte(generatedMarker + "\n\n" + strings.TrimSpace(body) + "\n")
	}
	return nil
}

// readmeIntro drops the title, logo, badges and language switch above the
// intro text.
func readmeIntro(s string) string {
	lines := strings.Split(s, "\n")
	i := 0
	for i < len(lines) {
		l := strings.TrimSpace(lines[i])
		if l == "" || strings.HasPrefix(l, "# ") || strings.HasPrefix(l, "<h1") || strings.HasPrefix(l, "[![") || isLangSwitch(l) {
			i++
			continue
		}
		break
	}
	return strings.Join(lines[i:], "\n")
}

// isLangSwitch reports whether a line is a document's link to its other
// language ("[日本語](README.ja.md)", "English | [日本語](...)"). The site has
// its own language menu.
func isLangSwitch(l string) bool {
	if !mdLink.MatchString(l) {
		return false
	}
	l = mdLink.ReplaceAllString(strings.TrimSpace(l), "$2")
	return strings.NewReplacer("English", "", "日本語", "", "|", "", " ", "").Replace(l) == ""
}

// dropLangSwitch removes the language switch line from the top of a document
// (before its first ## heading).
func dropLangSwitch(s string) string {
	lines := strings.SplitAfter(s, "\n")
	for i, l := range lines {
		if strings.HasPrefix(l, "## ") {
			break
		}
		if isLangSwitch(l) {
			return strings.Join(lines[:i], "") + strings.Join(lines[i+1:], "")
		}
	}
	return s
}

func headings(s, prefix string) []string {
	var out []string
	inFence := false
	for _, l := range strings.Split(s, "\n") {
		if strings.HasPrefix(l, "```") {
			inFence = !inFence
		}
		if !inFence && strings.HasPrefix(l, prefix) {
			out = append(out, strings.TrimSpace(l[len(prefix):]))
		}
	}
	return out
}

// splitSections returns the text before the first "## " heading and the body
// of each "## " section (without its heading line).
func splitSections(s string) []string {
	var out []string
	var cur strings.Builder
	inFence := false
	for _, l := range strings.Split(s, "\n") {
		if strings.HasPrefix(l, "```") {
			inFence = !inFence
		}
		if !inFence && strings.HasPrefix(l, "## ") {
			out = append(out, cur.String())
			cur.Reset()
			continue
		}
		cur.WriteString(l)
		cur.WriteByte('\n')
	}
	return append(out, cur.String())
}

var (
	mdLink  = regexp.MustCompile(`(!?)\[([^\]\n]*(?:\[[^\]\n]*\][^\]\n]*)*)\]\(([^)\s]+)\)`)
	htmlSrc = regexp.MustCompile(`(src|srcset|href)="([^"]+)"`)
)

// rewrite makes the relative links of a document of this repository
// absolute, and copies the images it references.
func (r *repo) rewrite(out map[string][]byte, prefix, file, text string) (string, error) {
	var b strings.Builder
	inFence := false
	var err error
	for _, line := range strings.SplitAfter(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inFence = !inFence
		}
		if inFence {
			b.WriteString(line)
			continue
		}
		line = mdLink.ReplaceAllStringFunc(line, func(m string) string {
			sm := mdLink.FindStringSubmatch(m)
			dest, e := r.link(out, prefix, file, sm[3], sm[1] == "!")
			if e != nil {
				err = e
			}
			return sm[1] + "[" + sm[2] + "](" + dest + ")"
		})
		line = htmlSrc.ReplaceAllStringFunc(line, func(m string) string {
			sm := htmlSrc.FindStringSubmatch(m)
			dest, e := r.link(out, prefix, file, sm[2], sm[1] != "href")
			if e != nil {
				err = e
			}
			return sm[1] + `="` + dest + `"`
		})
		b.WriteString(line)
	}
	return b.String(), err
}

func (r *repo) link(out map[string][]byte, prefix, file, dest string, image bool) (string, error) {
	if dest == "" || strings.Contains(dest, "://") || strings.HasPrefix(dest, "mailto:") {
		return dest, nil
	}
	target, frag, hasFrag := strings.Cut(dest, "#")
	if hasFrag {
		frag = "#" + frag
	}
	if target == "" {
		// An anchor in the same file. In a README, the section may be
		// on another page.
		if path.Base(file) == "README.md" || path.Base(file) == "README.ja.md" {
			return r.readmeAnchor(prefix, file, frag)
		}
		return dest, nil
	}
	rel := path.Clean(path.Join(path.Dir(file), target))
	if image {
		data, err := r.show(rel)
		if err != nil {
			return "", fmt.Errorf("%s: image %s: %v", file, dest, err)
		}
		out["public/repo/"+r.name+"/"+rel] = data
		return "/repo/" + r.name + "/" + rel, nil
	}
	en := strings.TrimSuffix(strings.TrimSuffix(rel, ".md"), ".ja") + ".md"
	if route, ok := r.docs[en]; ok {
		if en == "README.md" && frag != "" {
			return r.readmeAnchor(prefix, rel, frag)
		}
		// The document's own other language ("[日本語](CHANGELOG.ja.md)")
		// is its page in that language.
		if self := strings.TrimSuffix(strings.TrimSuffix(file, ".md"), ".ja") + ".md"; self == en && rel != file {
			prefix = ""
			if strings.HasSuffix(rel, ".ja.md") {
				prefix = "/ja"
			}
		}
		return r.base + prefix + route + frag, nil
	}
	kind := "blob"
	if path.Ext(rel) == "" {
		kind = "tree"
	}
	return r.web + "/" + kind + "/" + r.rev + "/" + rel + frag, nil
}

// readmeAnchor maps #section of a README to the guide page including it.
func (r *repo) readmeAnchor(prefix, file, frag string) (string, error) {
	slug := strings.TrimPrefix(frag, "#")
	if strings.HasSuffix(file, ".ja.md") {
		// Japanese headings: find the English section at the same position.
		ja, err := r.show(file)
		if err != nil {
			return "", err
		}
		en, err := r.show("README.md")
		if err != nil {
			return "", err
		}
		jh, eh := headings(string(ja), "## "), headings(string(en), "## ")
		for i, h := range jh {
			if press.Slug(h) == slug && i < len(eh) {
				slug = press.Slug(eh[i])
			}
		}
	}
	if route, ok := r.sections[slug]; ok {
		return r.base + prefix + route, nil
	}
	return r.base + prefix + r.docs["README.md"] + frag, nil
}
