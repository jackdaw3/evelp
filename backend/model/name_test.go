package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVal(t *testing.T) {
	name := &Name{
		De: "Deutsch",
		En: "English",
		Fr: "Français",
		Ja: "日本語",
		Ru: "Pусский",
		Zh: "中文",
	}

	assert.Equal(t, "Deutsch", name.Lang("de"))
	assert.Equal(t, "English", name.Lang("en"))
	assert.Equal(t, "Français", name.Lang("fr"))
	assert.Equal(t, "日本語", name.Lang("ja"))
	assert.Equal(t, "Pусский", name.Lang("ru"))
	assert.Equal(t, "中文", name.Lang("zh"))
}
