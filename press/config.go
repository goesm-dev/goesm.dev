package press

// Config describes a documentation site, in the spirit of VitePress's
// .vitepress/config: title, locales with their navigation, sidebar and UI
// strings, and where pages come from.
type Config struct {
	// SiteURL is the canonical origin, without a trailing slash
	// ("https://goesm.dev"); it is used for the sitemap and canonical links.
	SiteURL string
	// Base is the route every page of the site is under, without a trailing
	// slash: "" for a site at the root of SiteURL, "/gosfc" for one below
	// it. Links in the configuration (nav, sidebar, hero and feature links)
	// are relative to Base and the locale prefix; links in Markdown that
	// start with "/" are paths from the origin.
	Base string
	// EditBase is prepended to a page's content path for its "Edit this
	// page" link ("https://github.com/org/site/edit/main/content/").
	EditBase string
	// Sources are other repositories whose documents are synced into the
	// content tree. A page whose frontmatter says `source: goesm:docs/x.md`
	// resolves its relative links against docs/ of source "goesm".
	Sources map[string]Source
	// Social links shown in the navigation bar.
	Social []Social
	// Locales; the first is the root locale (Prefix "").
	Locales []*Locale
}

// Source is a repository documents are synced from.
type Source struct {
	// Repo is the repository's web URL ("https://github.com/goesm-dev/goesm").
	Repo string
	// Ref is the commit or tag the synced files come from.
	Ref string
	// Assets is the site path under which the source's images are copied
	// ("/repo/goesm/").
	Assets string
}

type Social struct {
	Icon string `json:"icon"` // "github"
	Link string `json:"link"`
}

// Locale is one language of the site.
type Locale struct {
	// Code is the content directory of the locale ("en", "ja").
	Code string
	// Prefix is the route prefix: "" for the root locale, "/ja" otherwise.
	Prefix string
	// Lang is the BCP 47 tag for <html lang> and hreflang.
	Lang        string
	Label       string // in its own language: "English", "日本語"
	Title       string
	Description string
	Nav         []NavItem
	// Sidebar maps a route prefix (without the locale prefix, e.g.
	// "/guide/") to the sidebar shown on pages below it.
	Sidebar map[string][]SidebarGroup
	UI      UI
}

type NavItem struct {
	Text string
	Link string // route without the locale prefix, or an absolute URL
	// Match is the route prefix that makes the item active (defaults to Link).
	Match string
	// Root marks Link as a path from the origin, for a link to another site
	// on the same domain: neither Base nor the locale prefix is added.
	Root bool
}

type SidebarGroup struct {
	Text  string
	Items []SidebarItem
}

type SidebarItem struct {
	Text string // defaults to the page title
	Link string // route without the locale prefix
}

// UI holds the theme's strings in one language.
type UI struct {
	OnThisPage    string `json:"onThisPage"`
	Prev          string `json:"prev"`
	Next          string `json:"next"`
	EditPage      string `json:"editPage"`
	SyncedFrom    string `json:"syncedFrom"` // "Synced from" + link to the source file
	Search        string `json:"search"`
	SearchHint    string `json:"searchHint"`
	NoResults     string `json:"noResults"`
	Menu          string `json:"menu"`
	ToggleTheme   string `json:"toggleTheme"`
	SkipToContent string `json:"skipToContent"`
	Language      string `json:"language"`
	NotFound      string `json:"notFound"`
	NotFoundText  string `json:"notFoundText"`
	BackHome      string `json:"backHome"`
	CopyCode      string `json:"copyCode"`
	Copied        string `json:"copied"`
	Footer        string `json:"footer"`
	SiteSource    string `json:"siteSource"`   // link text to the site's own repository
	ViewMarkdown  string `json:"viewMarkdown"` // link to the page's Markdown version
	// OtherPages heads the pages outside the sidebar in llms.txt.
	OtherPages string `json:"-"`
	// Containers are the default titles of ::: blocks and GitHub alerts.
	Containers map[string]string `json:"-"`
}
