package ui

import (
	"math/rand"

	"github.com/Varun0856/keyzen/internal/mode"
	"github.com/Varun0856/keyzen/internal/storage"
	"github.com/Varun0856/keyzen/internal/typing"
	"github.com/Varun0856/keyzen/internal/words"
)

func buildText(md mode.Mode, keys *storage.KeyStore, opts words.Options) string {
	text := ""

	if md.Kind == mode.Practice && keys != nil {
		if totals, err := keys.Load(); err == nil {
			if w := typing.Weights(totals.Runes()); len(w) > 0 {
				text = words.Weighted(md.TextLength(), w, rand.Float64)
			}
		}
	}
	if text == "" {
		text = words.Random(md.TextLength())
	}
	return words.Decorate(text, opts, rand.Float64)
}
