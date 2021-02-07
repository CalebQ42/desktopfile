package desktopfile

import (
	"strings"
)

//Locale is the information about a locale given for a Key. A Locale should ALWAYS have a Language value.
//In a Key-Value pair, the locale is given as lang_COUNTRY@MODIFIER.
//If encoding is specified (syntax of COUNTRY.ENCODING) it is ignored.
type Locale struct {
	Language string
	Country  string
	Modifier string
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

//Comment represents some comment in the .desktop file.
//This records the actual actual comment and which line it goes on.
//Inline comments are stored with their associated Key.
type Comment struct {
	Comment    string
	LineNumber int
}

func (c *Comment) String() string {
	return c.Comment
}

//Key represents a key value pair.
type Key struct {
	Key           string
	RawValue      string
	InlineComment string
	Locales       []LocaleValue
}

//LocaleValue is a value for a Key with a specific locale.
type LocaleValue struct {
	Locale
	RawValue      string
	InlineComment string
}
