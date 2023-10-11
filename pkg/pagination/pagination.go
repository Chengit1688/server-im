package pagination

type Pagination struct {
	Page     int   `json:"page" form:"page"`
	PageSize int   `json:"page_size" form:"page_size"`
	Offset   int   `json:"-"`
	Limit    int   `json:"-"`
	Count    int64 `json:"count"`
}

func (p *Pagination) Check() {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.PageSize == 0 {
		p.PageSize = 10
	}

	if p.PageSize > 100 {
		p.PageSize = 100
	}

	p.Offset = p.PageSize * (p.Page - 1)
	p.Limit = p.PageSize
}
