// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 110.
//!+

// Package github provides a Go API for the GitHub issue tracker.
// See https://developer.github.com/v3/search/#search-issues.
package github

import "time"

const IssuesURL = "https://api.github.com/search/issues"


/**
As before, the names of all the struct fields must be capitalized even if their JSON names are
not. However, the matching process that associates JSON names with Go struct names during
unmarshaling is case-insensitive , so it’s only necessary to use a field tag when there’s an underscore
in the JSON name but not in the Go name.
 */

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	// 返回的 json 字段 items 是小写, 这里为什么写成大写 ???
	// 把 json Unmarshal 到 struct 的时候: 忽略大小写, 参考 $ go doc json.Unmarshal
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

//!-
