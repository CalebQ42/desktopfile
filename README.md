# desktopfile

Go library to parse/generate .desktop files. This library (at least tries to) follows the freedesktop.org specifications. Functions will return zero values instead of nil so you can safely string together multiple functions without worry. Ex: File.Group("Non-existant").Key("Not-here").ToArray().

When reading, formatting and comments are NOT preserved.
