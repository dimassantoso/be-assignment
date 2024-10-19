package shared

import (
	"context"
	"math"
	"strings"

	"github.com/golangid/candi/candihelper"
	"github.com/golangid/candi/candishared"
	"gorm.io/gorm"
)

type BaseResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Meta model
type Meta struct {
	Page         int `json:"page"`
	Limit        int `json:"limit"`
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
}

// NewMeta create new meta for slice data
func NewMeta(page, limit, totalRecords int) (m Meta) {
	m.Page, m.Limit, m.TotalRecords = page, limit, totalRecords
	m.CalculatePages()
	return m
}

// CalculatePages meta method
func (m *Meta) CalculatePages() {
	m.TotalPages = int(math.Ceil(float64(m.TotalRecords) / float64(m.Limit)))
}

// Filter data
type Filter struct {
	Limit              int      `json:"limit" default:"10"`
	Page               int      `json:"page" default:"1"`
	Offset             int      `json:"-"`
	Search             string   `json:"search,omitempty"`
	OrderBy            string   `json:"order_by,omitempty"`
	Sort               string   `json:"sort,omitempty" default:"desc" lower:"true"`
	ShowAll            bool     `json:"show_all"`
	AllowEmptyFilter   bool     `json:"-"`
	CompID             *string  `json:"-"`
	LangID             *string  `json:"-"`
	Relations          []string `json:"-"`
	GroupIDs           []string `json:"group_ids"`
	UserGroupIDs       []string `json:"user_group_ids"`
	IsViewAllUserGroup bool     `json:"is_view_all_user_group"`
}

func (f *Filter) ParseToQuery(db *gorm.DB) *gorm.DB {
	for _, preload := range f.Relations {
		var args []any
		if strings.HasSuffix(preload, "Language") && f.LangID != nil {
			args = []any{"language_id = ?", *f.LangID}
		}
		db = db.Preload(preload, args...)
	}
	return db
}

func (f *Filter) ParseFromTokenClaim(ctx context.Context) {
	stdClaim := candishared.ParseTokenClaimFromContext(ctx)
	additional, ok := stdClaim.Additional.(map[string]string)
	if ok {
		f.CompID = candihelper.WrapPtr(additional["company_id"])
		f.LangID = candihelper.WrapPtr(additional["language_id"])
		if additional["group_ids"] != "" {
			f.GroupIDs = strings.Split(additional["group_ids"], ",")
		}
		if additional["user_group_ids"] != "" {
			f.UserGroupIDs = strings.Split(additional["user_group_ids"], ",")
		}
		f.IsViewAllUserGroup = additional["is_view_all_user_group"] == "true"
	}
}

// CalculateOffset method
func (f *Filter) CalculateOffset() int {
	f.Offset = (f.Page - 1) * f.Limit
	return f.Offset
}

// GetPage method
func (f *Filter) GetPage() int {
	return f.Page
}

// IncrPage method
func (f *Filter) IncrPage() {
	f.Page++
}

// GetLimit method
func (f *Filter) GetLimit() int {
	return f.Limit
}
