package views

import (
	"models";
	"http";
	"io";
	"os";
	"strconv";
	"template";
)

// Persistent blog state. Switch off the server and watch this vanish.
var blog models.Entries

type frontPageEntry struct {
	*models.Entry;
	CommentCount int;
}

// url: /blog
func Index(c *http.Conn, req *http.Request) {
	revEntries := blog.EntriesReversed();
	frontPageEntries := make([]frontPageEntry, len(revEntries));
	for i, e := range revEntries {
		frontPageEntries[i].Entry = e;
		frontPageEntries[i].CommentCount = e.Comments.Len();
	}

	if tmpl, err := loadTemplate("index.html"); err == nil {
		tmpl.Execute(frontPageEntries, c);
	} else {
		c.SetHeader("Content-Type", "text/plain; charset=utf-8");
		c.WriteHeader(http.StatusInternalServerError);
		io.WriteString(c, "500 internal server error\n");
	}
}

// url: /blog/entry/add
func AddEntry(c *http.Conn, req *http.Request) {
	req.ParseForm();
	heading, body := req.FormValue("heading"), req.FormValue("body");
	if len(heading) > 0 && len(body) > 0 {
		blog.AddEntry(heading, body)
	}
	c.SetHeader("Location", "/blog");
	c.WriteHeader(http.StatusFound);
}

type singlePageEntry struct {
	*models.Entry;
	Comments []string;
	CommentCount int;
}

// url: /blog/entry/(\d+)
//   1: ID of blog post to show.
func Entry(c *http.Conn, req *http.Request) {
	idString := req.URL.String();
	idString = idString[len("/blog/entry/"):len(idString)];
	if entryID, parseErr := strconv.Atoi(idString); parseErr == nil {
		if foundEntry := blog.FindEntry(entryID); foundEntry != nil {
			var entry singlePageEntry;
			entry.Entry = foundEntry;
			entry.Comments = foundEntry.Comments.Data();
			entry.CommentCount = len(entry.Comments);

			if tmpl, err := loadTemplate("entry.html"); err == nil {
				tmpl.Execute(entry, c);
			} else {
				c.SetHeader("Content-Type", "text/plain; charset=utf-8");
				c.WriteHeader(http.StatusInternalServerError);
				io.WriteString(c, "500 internal server error\n");
			}
		} else {
			c.SetHeader("Content-Type", "text/plain; charset=utf-8");
			c.WriteHeader(http.StatusNotFound);
			io.WriteString(c, "404 not found\n");
		}
	} else {
		c.SetHeader("Content-Type", "text/plain; charset=utf-8");
		c.WriteHeader(http.StatusInternalServerError);
		io.WriteString(c, "500 internal server error\n");
	}
}

// url: /blog/comment/add
func AddComment(c *http.Conn, req *http.Request) {
	req.ParseForm();
	entryIDString, text := req.FormValue("entry_id"), req.FormValue("text");
	if len(entryIDString) > 0 && len(text) > 0 {
		if entryID, parseErr := strconv.Atoi(entryIDString); parseErr == nil {
			if foundEntry := blog.FindEntry(entryID); foundEntry != nil {
				foundEntry.AddComment(text);
			}
		}
	}
	c.SetHeader("Location", "/blog/entry/" + entryIDString);
	c.WriteHeader(http.StatusFound);
}

func loadTemplate(path string) (t *template.Template, err os.Error) {
	if buf, readErr := io.ReadFile("templates/" + path); readErr == nil {
		t = template.MustParse(string(buf), template.FormatterMap{
			"html": template.HTMLFormatter
		});
	} else {
		err = readErr
	}
	return;
}
