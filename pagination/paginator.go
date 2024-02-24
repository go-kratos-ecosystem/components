package pagination

import "encoding/json"

type Paginator struct {
	Page    int `json:"page"`
	PrePage int `json:"pre_page"`
	Total   int `json:"total"`
}

func NewPaginator(page, pageSize, total int) *Paginator {
	return &Paginator{
		Page:    page,
		PrePage: pageSize,
		Total:   total,
	}
}

func (p *Paginator) GetPage() int {
	return p.Page
}

func (p *Paginator) GetPerPage() int {
	return p.PrePage
}

func (p *Paginator) GetTotal() int {
	return p.Total
}

func (p *Paginator) GetLastPage() int {
	return (p.Total + p.PrePage - 1) / p.PrePage
}

func (p *Paginator) GetOffset() int {
	return (p.Page - 1) * p.PrePage
}

func (p *Paginator) GetLimit() int {
	return p.PrePage
}

func (p *Paginator) HasMore() bool {
	return p.Page*p.PrePage < p.Total
}

func (p *Paginator) GetNextPage() int {
	if p.HasMore() {
		return p.Page + 1
	}
	return p.Page
}

func (p *Paginator) GetPrevPage() int {
	if p.Page > 1 {
		return p.Page - 1
	}
	return p.Page
}

func (p *Paginator) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"page":      p.Page,
		"pre_page":  p.PrePage,
		"total":     p.Total,
		"next_page": p.GetNextPage(),
		"prev_page": p.GetPrevPage(),
		"last_page": p.GetLastPage(),
		"has_more":  p.HasMore(),
	}
}

func (p *Paginator) ToJSON() ([]byte, error) {
	return json.Marshal(p.ToMap())
}
