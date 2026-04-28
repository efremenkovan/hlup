package options

import "github.com/efremenkovan/hlup/lang"

type Options struct {
	Lang lang.Lang
}

type PatchFunc func(o *Options)

func WithLang(lang lang.Lang) PatchFunc {
	return func(o *Options) {
		o.Lang = lang
	}
}
