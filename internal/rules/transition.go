package rules

import (
	"fmt"
	"github.com/11DingKing/yungang-grottoes-condensation-monitoring/internal/domain/common"
)

type Graph struct{ edges map[string]map[string]bool }

func NewGraph() Graph { return Graph{edges: map[string]map[string]bool{}} }
func (g Graph) Add(from, to string) {
	if g.edges[from] == nil {
		g.edges[from] = map[string]bool{}
	}
	g.edges[from][to] = true
}
func (g Graph) Can(from, to string) bool { return g.edges[from][to] }
func (g Graph) Move(current, next string) error {
	if !g.Can(current, next) {
		return fmt.Errorf("%w: %s -> %s", common.ErrConflict, current, next)
	}
	return nil
}
func (g Graph) Reachable(start string) []string {
	seen := map[string]bool{}
	out := []string{}
	var walk func(string)
	walk = func(v string) {
		if seen[v] {
			return
		}
		seen[v] = true
		out = append(out, v)
		for n := range g.edges[v] {
			walk(n)
		}
	}
	walk(start)
	return out
}
func (g Graph) Validate() error {
	for from, targets := range g.edges {
		if from == "" {
			return common.ErrInvalid
		}
		for to := range targets {
			if to == "" {
				return common.ErrInvalid
			}
		}
	}
	return nil
}
