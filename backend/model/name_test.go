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

	assert.Equal(t, "Deutsch", name.Val("de"))
	assert.Equal(t, "English", name.Val("en"))
	assert.Equal(t, "Français", name.Val("fr"))
	assert.Equal(t, "日本語", name.Val("ja"))
	assert.Equal(t, "Pусский", name.Val("ru"))
	assert.Equal(t, "中文", name.Val("zh"))
}
