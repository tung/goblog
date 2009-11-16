package models

import (
	"container/vector";
	"strconv";
)

type Entries struct {
	entries vector.Vector;
	nextID int;
}

type Entry struct {
	ID int;
	Heading, Body string;
	Comments vector.StringVector;
}

func (es *Entries) String() string {
	out := "[";
	i := 0;
	for entry := range es.entries.Iter() {
		if i > 0 {
			out += ", "
		}
		out += entry.(*Entry).String();
		i++;
	}
	return out + "]";
}

func (es *Entries) AddEntry(heading, body string) *Entry {
	entry := new(Entry);
	entry.ID = es.nextID;
	es.nextID++;
	entry.Heading = heading;
	entry.Body = body;
	es.entries.Push(entry);
	return entry;
}

func (es *Entries) EntriesReversed() []*Entry {
	rev := make([]*Entry, es.entries.Len());
	i := 0;
	for e := range es.entries.Iter() {
		rev[len(rev) - i - 1] = e.(*Entry);
		i++;
	}
	return rev;
}

func (es *Entries) FindEntry(id int) *Entry {
	// Linear search, nothing fancy yet.
	for e := range es.entries.Iter() {
		if e.(*Entry).ID == id {
			return e.(*Entry)
		}
	}
	return nil
}

func (e *Entry) String() string {
	return "\"" + e.Heading + "\" (" + strconv.Itoa(e.Comments.Len()) + ")";
}

func (e *Entry) AddComment(comment string) string {
	e.Comments.Push(comment);
	return comment;
}
