package cronic

import (
	"fmt"
	"io/ioutil"
	"sort"
)

type crontabCollection struct {
	crontabs map[string]*crontab
}

func crontabPaths(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to load crontabs from %s", path)
	}

	paths := make([]string, 0)
	for _, file := range files {
		paths = append(paths, path+"/"+file.Name())
	}

	return paths, nil
}

func loadCrontabs(paths []string) (*crontabCollection, error) {
	crontabs := make(map[string]*crontab, 0)

	for _, path := range paths {
		crontab, err := newCrontabFromFile(path)
		if err != nil {
			return nil, err
		}

		crontabs[path] = crontab
	}

	return &crontabCollection{
		crontabs,
	}, nil
}

func (cc crontabCollection) Paths() []string {
	var paths []string

	for _, c := range cc.crontabs {
		paths = append(paths, c.path)
	}

	sort.Strings(paths)

	return paths
}

func (cc crontabCollection) named(name string) *crontab {
	return cc.crontabs[name]
}
