package main

// ─────────────────────────────────────────────────────────────────────────
// EDIT THIS FILE to make the site yours. Everything visitors see comes from
// here: your name, bio, projects, writing, and contact links. The values
// below are placeholders copied from the design you shared — replace them.
// ─────────────────────────────────────────────────────────────────────────

// name renders as the big figlet banner on the landing screen.
var name = "omer"

// bio is the paragraph shown to the right of your portrait. Each string is a
// line; keep them short so they fit narrow terminals.
var bio = []string{
	"is a builder & problem solver at the",
	"intersection of systems and intelligence,",
	"developing AI agents, automating complex",
	"workflows, and thinking deeply about how",
	"software can do more of the thinking for us.",
	"",
	"He also works as an AI engineer at",
	"Accenture, where he helps design and",
	"deploy intelligent agent systems for",
	"real-world enterprise problems.",
	"",
	"Previously, Omer studied Computer Science",
	"and Mathematics at the University of",
	"Michigan, where he built full-stack",
	"automation tools, conducted ML research,",
	"and taught Real Analysis — all while",
	"constructing a trebuchet on weekends.",
	"",
	"His work sits at the intersection of",
	"rigorous math, human systems, and the",
	"technology quietly reshaping both.",
}

// item is one selectable entry in a section. Star marks it as featured (✦).
type item struct {
	title string
	url   string
	star  bool
}

// group is a labeled cluster of items within a section.
type group struct {
	category string
	items    []item
}

// section is a full page reachable from the footer nav.
type section struct {
	name   string
	groups []group
}

var creations = section{
	name: "Creations",
	groups: []group{
		{category: "web", items: []item{
			{title: "Daily Integral — a daily integral website game", url: "https://dailyintegral-delta.vercel.app/"},
		}},
		{category: "graphics", items: []item{
			{title: "Mini Raytracer — a raytracer written from scratch in C++", url: "https://github.com/ojunedi/miniraytracer", star: true},
		}},
	},
}

var reflections = section{
	name: "Reflections",
	groups: []group{
		{category: "philosophy", items: []item{
			{title: "The Source of Action: Knowledge and Duty in the Bhagavad Gita", url: "https://example.com/gita"},
		}},
		{category: "technology", items: []item{
			{title: "Reimagining Human Labor in the Age of AI", url: "https://example.com/labor", star: true},
			{title: "AI as a Creative Springboard: Enhancing, Not Replacing, Human Ingenuity", url: "https://example.com/springboard"},
		}},
	},
}

// contact is one row on the Contacts page.
type contact struct {
	label  string
	handle string
	url    string
}

var contacts = []contact{
	{label: "LI", handle: "linkedin.com/in/omer-junedi", url: "https://www.linkedin.com/in/omer-junedi/"},
	{label: "GH", handle: "github.com/ojunedi", url: "https://github.com/ojunedi"},
	{label: "EM", handle: "ojunedi@umich.edu", url: "mailto:ojunedi@umich.edu"},
}

// navEntry is one footer-nav destination. Pages are driven entirely by this
// slice, so showing/hiding a page is just adding/removing an entry here — page
// indices stay consistent automatically.
type navEntry struct {
	label    string
	section  *section // page content; nil for the contacts page
	contacts bool     // true for the special contacts page
}

var nav = []navEntry{
	{label: creations.name, section: &creations},
	// Reflections is hidden for now (no writing to show yet). To bring it back,
	// just uncomment the line below — the reflections content above stays in
	// the code either way.
	// {label: reflections.name, section: &reflections},
	{label: "Contacts", contacts: true},
}
