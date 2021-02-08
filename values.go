package desktopfile

import (
	"strconv"
	"strings"
)

var escapeReplacer = strings.NewReplacer("\\s", " ", "\\n", "\n", "\\t", "	", "\\r", "\r", "\\\\", "\\")

//Locale is the information about a locale given for a Entry. A Locale should ALWAYS have a Language value.
type Locale struct {
	Language string
	Country  string
	Modifier string
}

//LocaleFromString creates a Locale from a string (in the format lang_COUNTRY@MODIFIER).
//If encoding is specified (in the form COUNTRY.ENCODING), it is ignored.
func LocaleFromString(locale string) (out Locale) {
	underInd := strings.Index(locale, "_")
	atInd := strings.Index(locale, "@")
	if underInd == -1 && atInd == -1 {
		out.Language = locale
	} else if underInd == -1 {
		out.Language = locale[:atInd]
		out.Modifier = locale[atInd+1:]
	} else if atInd == -1 {
		out.Language = locale[:underInd]
		out.Country = locale[underInd+1:]
	} else {
		out.Language = locale[:underInd]
		out.Country = locale[underInd+1 : atInd]
		out.Modifier = locale[atInd:]
	}
	if dotInd := strings.Index(out.Country, "."); dotInd != -1 {
		out.Country = out.Country[:dotInd]
	}
	return
}

func (l Locale) String() (out string) {
	out = strings.ToLower(l.Language)
	if l.Country != "" {
		out += "_" + strings.ToUpper(l.Country)
	}
	if l.Modifier != "" {
		out += "@" + strings.ToUpper(l.Modifier)
	}
	return
}

//Value is a value for a given Entry
type Value string

//String returns the Value as a string. Works on all types.
//If of a string type, replacing escaped characters.
func (v Value) String() string {
	return escapeReplacer.Replace(string(v))
}

//IsBool returns whether or not the Value is a bool type.
func (v Value) IsBool() bool {
	tmp := strings.ToLower(string(v))
	if tmp == "true" || tmp == "false" {
		return true
	}
	return false
}

//AsBool tries to convert the Value to a bool. If not a bool, returns false.
func (v Value) AsBool() bool {
	return strings.ToLower(string(v)) == "true"
}

//IsFloat returns whether or not the Value is a float type.
func (v Value) IsFloat() bool {
	_, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		return false
	}
	return true
}

//AsFloat attempts to convert the Value to a float64. If not a float, returns 0.0.
func (v Value) AsFloat() float64 {
	f, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		return 0
	}
	return f
}

//IsInt returns whether or not the Value is a int type.
func (v Value) IsInt() bool {
	_, err := strconv.Atoi(string(v))
	if err != nil {
		return false
	}
	return true
}

//AsInt attempts to convert the Value to an int. If not an int, returns 0.
func (v Value) AsInt() int {
	i, err := strconv.Atoi(string(v))
	if err != nil {
		return 0
	}
	return i
}

//IsArray returns whether or not the value is an array.
//Value Arrays have values seperated by semicolons.
func (v Value) IsArray() bool {
	if strings.Count(string(v), "\\;") < strings.Count(string(v), ";") {
		return true
	}
	return false
}

//AsArray attempts to convert the Value to a Value slice.
//If it is not an array, it returns a zero length slice.
func (v Value) AsArray() []Value {
	split := strings.Split(string(v), ";")
	if strings.HasSuffix(string(v), ";") {
		split = split[:len(split)-1]
	}
	for i := 0; i < len(split); i++ {
		if strings.HasSuffix(split[i], "\\") {
			split[i] = split[i][:len(split[i])-1] + ";" + split[i+1]
			for j := i + 1; j < len(split)-1; j++ {
				split[j] = split[j+1]
			}
			split = split[:len(split)-1]
			i--
		}
	}
	out := make([]Value, len(split))
	for i := range out {
		out[i] = Value(split[i])
	}
	return out
}

//LocaleValue is a value for a Entry with a specific locale.
type LocaleValue struct {
	Locale
	Value   Value
	Comment string
}

//Entry represents a entry value pair.
type Entry struct {
	Comment string
	Value   Value
	Locales []*LocaleValue
}

//ValueAtLocale gets the
func (e Entry) ValueAtLocale(l Locale) Value {
	matchInd := -1
	matchedLang, matchedCountry := false, false
	for i, lv := range e.Locales {
		if l.Language != lv.Language {
			continue
		}
		if !matchedLang {
			matchedLang = true
			matchInd = i
		}
		if l.Country != lv.Country {
			continue
		}
		if !matchedCountry {
			matchedCountry = true
			matchInd = i
		}
		if l.Modifier == lv.Modifier {
			return lv.Value
		}
	}
	if matchInd == -1 {
		return e.Value
	}
	return e.Locales[matchInd].Value
}

//Group is a set of Entries under a group header.
type Group struct {
	entries map[string]*Entry
	Comment string
}

//HasEntry returns if the Group has the given Entry
func (g Group) HasEntry(key string) bool {
	for name := range g.entries {
		if name == key {
			return true
		}
	}
	return false
}

//GetEntry returns the Entry at the given key.
//If one is not found, this returns a zero valued Entry that's still safe to use.
func (g Group) GetEntry(key string) *Entry {
	if e, ok := g.entries[key]; ok {
		return e
	}
	return &Entry{
		Locales: make([]*LocaleValue, 0),
	}
}

//AddEntry adds a Entry with the given key.
//It returns the pointer so you can edit it's value.
//
//If adding a localized entry, do it via it's main Entry.
func (g *Group) AddEntry(key string) *Entry {
	if e, ok := g.entries[key]; ok {
		return e
	}
	g.entries[key] = &Entry{
		Locales: make([]*LocaleValue, 0),
	}
	return g.entries[key]
}

//RemoveEntry removes the entry with the given key.
func (g *Group) RemoveEntry(key string) {
	delete(g.entries, key)
}
