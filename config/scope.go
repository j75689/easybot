package config

import (
	"regexp"
	"strings"
)

// Scope service access scope path
var Scope = ScopeDefinition{
	"push": ScopePath{
		"/api/v1/bot/push/:userID",
	},
	"multicast": ScopePath{
		"/api/v1/bot/multicast",
	},
	"config": ScopePath{
		"/api/v1/config/:id",
	},
	"plugin": ScopePath{
		"/api/v1/plugin/:plugin",
	},
}

type ScopeTag string

func (tag *ScopeTag) Label() string {
	return strings.Title(string(*tag))
}

func (tag *ScopeTag) Value() string {
	return string(*tag)
}

type ScopePath []string

func (paths *ScopePath) Match(path string) bool {
	for _, p := range *paths {
		var replacer = regexp.MustCompile(`:[\w]*`)
		p = replacer.ReplaceAllString(p, `.*`)
		if match, _ := regexp.MatchString(p, path); match {
			return true
		}
	}
	return false
}

type ScopeDefinition map[ScopeTag]ScopePath

func (scope *ScopeDefinition) Tags() []struct{ Label, Value string } {
	tags := []struct{ Label, Value string }{}
	for key, _ := range *scope {
		tags = append(tags, struct{ Label, Value string }{Label: key.Label(), Value: key.Value()})
	}
	return tags
}

// Allow verify permission
// if scope == all , allow all path
// if scope != all , verify path
func (scope *ScopeDefinition) Allow(scopeString, path string) bool {
	targets := strings.Split(scopeString, ",")
	for _, target := range targets {
		if target == "all" {
			return true
		}
		if paths := (*scope)[ScopeTag(target)]; paths != nil {
			if paths.Match(path) {
				return true
			}
		}
	}
	return false
}
