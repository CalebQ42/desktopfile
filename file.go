package desktopfile

import (
	"bufio"
	"io"
)

//File a representation of a .desktop file.
type File struct {
	rawMap   map[string][]*Key
	options  Options
	Comments []Comment
}

//Options gives options for loading values from a .desktop file.
//DefaultLocale dictactes which locale to try to use if different locales are specified for a key.
type Options struct {
	DefaultLocale  Locale
	RemoveComments bool
	AutoEscape     bool
}

//DefaultOptions returns the default options for reading.
func DefaultOptions() Options {
	return Options{
		AutoEscape: true,
	}
}

//Open reads the .desktop file from an io.Reader.
func Open(reader io.Reader) (file *File, err error) {
	return OpenWithOptions(reader, Options{})
}

//OpenWithOptions creates a desktopfile.File with the given options.
func OpenWithOptions(reader io.Reader, op Options) (file *File, err error) {
	file = new(File)
	file.options = op
	rdr := bufio.NewReader(reader)
	_ = rdr
	return nil, nil
}
