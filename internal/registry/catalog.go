package registry

import (
	"sort"
	"strings"
)

type CatalogItem struct {
	Code, Name, Category string
	Enabled              bool
}
type Catalog struct{ items map[string]CatalogItem }

func NewCatalog() *Catalog { return &Catalog{items: map[string]CatalogItem{}} }
func (c *Catalog) Add(v CatalogItem) {
	if v.Code != "" {
		c.items[v.Code] = v
	}
}
func (c *Catalog) Find(code string) (CatalogItem, bool) { v, ok := c.items[code]; return v, ok }
func (c *Catalog) Search(term string) []CatalogItem {
	term = strings.ToLower(term)
	out := []CatalogItem{}
	for _, v := range c.items {
		if v.Enabled && (term == "" || strings.Contains(strings.ToLower(v.Name), term)) {
			out = append(out, v)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out
}
func (c *Catalog) Disable(code string) bool {
	v, ok := c.items[code]
	if !ok {
		return false
	}
	v.Enabled = false
	c.items[code] = v
	return true
}
