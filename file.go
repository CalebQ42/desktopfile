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
	rawMap        map[string]*Group
	Options       Options
	EndingComment string
	groupOrder    []string
}

//Options are options for when reading or writing a File.
type Options struct {
	DefaultLocale            *Locale //The default locale used when getting values.
	AllowDuplicateKeysJoin   bool    //When encountering duplicate keys in a desktop file, join their values and comments instead of returning an error. Has precedence.
	AllowDuplicateKeysIgnore bool    //When encountering duplicate keys in a desktop file, ignore duplicate entries.
	AllowDuplicateGroups     bool    //When encountering duplicate group headers, join thier keys instead of returning an error.
}

//DefaultOptions returns the default Options used when reading.
//Currently just returns `Options{}` but if further options are added that need to be non-default values, they will be set here.
func DefaultOptions() Options {
	return Options{}
}

//Open reads the .desktop file from an io.Reader.
func Open(reader io.Reader) (*File, error) {
	return OpenWithOptions(reader, Options{})
}

//OpenWithOptions reads the .desktop file from the io.Reader, using the given Options.
func OpenWithOptions(reader io.Reader, op Options) (*File, error) {
	file := File{
		groupOrder: make([]string, 0),
		Options:    op,
		rawMap:     make(map[string]*Group),
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
		if err != nil {
			if err == io.EOF && line == "" {
				break
			} else if err != io.EOF {
				return nil, err
			}
		}
		line = strings.TrimSuffix(strings.TrimSpace(line), "\n")
		if strings.HasPrefix(line, "#") || line == "" {
			if strings.HasSuffix(commentTemp, "\n") {
				commentTemp += line
			} else {
				commentTemp += "\n" + line
			}
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			group := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
			if _, ok := file.rawMap[group]; !ok {
				file.rawMap[group] = &Group{
					parent:  &file,
					entries: make(map[string]*Entry),
				}
			} else if !op.AllowDuplicateGroups {
				return nil, errors.New("Line " + strconv.Itoa(lineNum) + " has a dpulicate group header: " + group)
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
		value := strings.TrimSpace(line[equLoc+1:])
		var locale *Locale
		if strings.HasSuffix(key, "]") && strings.Contains(key, "]") {
			locale = LocaleFromString(key[strings.Index(key, "[")+1 : len(key)-1])
			key = strings.TrimSpace(key[:strings.Index(key, "[")])
		}
		if k, ok := curGroup.entries[key]; !ok {
			curGroup.entries[key] = &Entry{
				parent:      curGroup,
				locales:     make(map[Locale]*LocaleValue),
				localeOrder: make([]Locale, 0),
			}
			curGroup.entryOrder = append(curGroup.entryOrder, key)
			if locale == nil {
				curGroup.entries[key].Value = Value(value)
				curGroup.entries[key].Comment = commentTemp
				commentTemp = ""
			}
		} else if ok && locale == nil {
			if op.AllowDuplicateKeysJoin {
				k.Value += Value(value)
				k.Comment += commentTemp
				commentTemp = ""
			} else if !op.AllowDuplicateKeysIgnore {
				return nil, errors.New("Line " + strconv.Itoa(lineNum) + " has a dupicate key: " + key)
			}
		}
		if locale != nil {
			if lv, ok := curGroup.entries[key].locales[*locale]; ok {
				if op.AllowDuplicateKeysJoin {
					lv.Value += Value(value)
					lv.Comment += commentTemp
					commentTemp = ""
				} else if !op.AllowDuplicateKeysIgnore {
					return nil, errors.New("Line " + strconv.Itoa(lineNum) + " has a duplicate key: " + strings.TrimSpace(line[:equLoc]))
				}
			} else {
				curGroup.entries[key].locales[*locale] = &LocaleValue{
					Value:   Value(value),
					Comment: commentTemp,
				}
				commentTemp = ""
				curGroup.entries[key].localeOrder = append(curGroup.entries[key].localeOrder, *locale)
			}
		}
	}
	if commentTemp != "" {
		file.EndingComment = commentTemp
	}
	return &file, nil
}

//DefaultGroup attempts to return the default Desktop Entry Group.
//If the Desktop Entry Group is not present, it is created.
func (f File) DefaultGroup() *Group {
	if _, ok := f.rawMap["Desktop Entry"]; ok {
		return f.rawMap["Desktop Entry"]
	}
	f.rawMap["Desktop Entry"] = &Group{
		parent:  &f,
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
		parent:  &f,
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
		parent:  f,
		entries: make(map[string]*Entry),
	}
	return f.rawMap[name]
}

//RemoveGroup removes the group from the File.
func (f *File) RemoveGroup(name string) {
	delete(f.rawMap, name)
	for i, v := range f.groupOrder {
		if v == name {
			f.groupOrder = append(f.groupOrder[:i], f.groupOrder[i+1:]...)
			return
		}
	}
}
