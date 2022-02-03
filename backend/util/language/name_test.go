package language

import (
	"evelp/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	name := &model.Name{
		De: "Deutsch",
		En: "English",
		Fr: "Français",
		Ja: "日本語",
		Ru: "Pусский",
		Zh: "中文",
	}

	assert.Equal(t, "Deutsch", Name("de", *name))
	assert.Equal(t, "English", Name("en", *name))
	assert.Equal(t, "Français", Name("fr", *name))
	assert.Equal(t, "日本語", Name("ja", *name))
	assert.Equal(t, "Pусский", Name("ru", *name))
	assert.Equal(t, "中文", Name("zh", *name))
}
