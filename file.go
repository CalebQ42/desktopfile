package desktopfile

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

//File a representation of a .desktop file.
type File struct {
	rawMap           map[string]*Group
	DefaultLocale    Locale
	BeginningComment string
	EndingComment    string
}

//Open reads the .desktop file from an io.Reader.
func Open(reader io.Reader) (*File, error) {
	file := File{
		rawMap: make(map[string]*Group),
	}
	rdr := bufio.NewReader(reader)
	lineNum := 0
	var curGroup *Group
	var commentTemp string
	var line string
	var err error
	for {
		lineNum++
		line, err = rdr.ReadString('\n')
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}
		line = strings.TrimSuffix(strings.TrimSpace(line), "\n")
		if strings.HasPrefix(line, "#") || line == "" {
			commentTemp += line + "\n"
			continue
		} else if strings.Contains(line, "#") {
			ind := strings.Index(line, "#")
			commentTemp += line[ind:] + "\n"
			line = strings.TrimSpace(line[:ind])
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			group := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
			if _, ok := file.rawMap[group]; !ok {
				file.rawMap[group] = &Group{
					entries: make(map[string]*Entry),
				}
			}
			file.rawMap[group].Comment += commentTemp
			commentTemp = ""
			curGroup = file.rawMap[group]
			continue
		}
		equLoc := strings.Index(line, "=")
		if equLoc == -1 {
			return nil, errors.New("Line " + strconv.Itoa(lineNum) + " is not a key, comment, group header, or whitespace.")
		} else if curGroup == nil {
			return nil, errors.New("Line " + strconv.Itoa(lineNum) + " has a key before a group heading.")
		}
		key := strings.TrimSpace(line[:equLoc])
		_ = key
	}
	if commentTemp != "" {
		file.EndingComment = commentTemp
	}
	//TODO: parse file
	return &file, nil
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
