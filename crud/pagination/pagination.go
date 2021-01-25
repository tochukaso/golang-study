package pagination

import (
	"fmt"
	"html/template"
)

type Paginavi struct {
	Count    int
	Page     int
	PageSize int
}

func Pagination(count int, page, pageSize int) Paginavi {
	return Paginavi{
		Count:    count,
		Page:     page,
		PageSize: pageSize,
	}
}

func (p Paginavi) CountExplanation() string {
	from := (p.Page-1)*p.PageSize + 1
	to := p.Page * p.PageSize
	return fmt.Sprintf("全 %d 件中 %d - %d 件目", p.Count, from, min(p.Count, to))
}

func (p Paginavi) Navigation() template.HTML {
	var tag string
	totalPage := (int(p.Count) + p.PageSize - 1) / p.PageSize
	tag += "<ul class='pagination'>"
	var firstDisabled string
	if p.Page <= 1 {
		firstDisabled = "btn disabled"
	}
	var lastDisabled string
	if p.Page >= totalPage {
		lastDisabled = "btn disabled"
	}
	tag += fmt.Sprintf("<li class='page-item'><button class='page-link %s' onclick=\"$('#page').val('%d');$('#searchForm').submit();\">Previous</button></li>", firstDisabled, max(1, p.Page-1))
	if firstDisabled == "" {
		bef := p.Page - 1
		tag += fmt.Sprintf("<li class='page-item'><button class='page-link' onclick=\"$('#page').val('%d');$('#searchForm').submit();\">%d</button></li>", bef, bef)
	}
	tag += fmt.Sprintf("<li class='page-item active'><button class='page-link'>%d</button></li>", p.Page)

	if lastDisabled == "" {
		af := p.Page + 1
		tag += fmt.Sprintf("<li class='page-item'><button class='page-link' onclick=\"$('#page').val('%d');$('#searchForm').submit();\">%d</button></li>", af, af)
	}

	tag += fmt.Sprintf("<li class='page-item'><button class='page-link %s' onclick=\"$('#page').val('%d');$('#searchForm').submit();\">Next</button></li>", lastDisabled, min(totalPage, p.Page+1))
	tag += "</ul>"

	return template.HTML(tag)
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
