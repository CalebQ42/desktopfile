package desktopfile

import (
	"bufio"
	"io"
)

//File a representation of a .desktop file.
type File struct {
	rawMap        map[string]*Group
	DefaultLocale Locale
}

//Open reads the .desktop file from an io.Reader.
func Open(reader io.Reader) *File {
	file := new(File)
	rdr := bufio.NewReader(reader)
	_ = rdr
	//TODO: parse file
	return file
}

//DefaultGroup attempts to return the default Desktop Entry Group.
//If the Desktop Entry Group is not present, it is created.
func (f File) DefaultGroup() *Group {
	if _, ok := f.rawMap["Desktop Entry"]; ok {
		return f.rawMap["Desktop Entry"]
	}
	f.rawMap["Desktop Entry"] = &Group{
		entries: make(map[string]*Entry),
	}
	return f.rawMap["Desktop Entry"]
}

//HasGroup returns whether the group with the name is present.
func (f File) HasGroup(name string) bool {
	_, ok := f.rawMap[name]
	return ok
}

//GetGroup returns the Group with the given name.
//If one is not found, this returns a zero valued Group that's still safe to use.
func (f File) GetGroup(name string) *Group {
	if val, ok := f.rawMap[name]; ok {
		return val
	}
	return &Group{
		entries: make(map[string]*Entry),
	}
}

//AddGroup adds the given Group to the File and returns the Group.
//If the group is present, it returns that Group.
func (f *File) AddGroup(name string) *Group {
	_, ok := f.rawMap[name]
	if ok {
		return f.rawMap[name]
	}
	f.rawMap[name] = &Group{
		entries: make(map[string]*Entry),
	}
	return f.rawMap[name]
}

//RemoveGroup removes the group from the File.
func (f *File) RemoveGroup(name string) {
	delete(f.rawMap, name)
}
