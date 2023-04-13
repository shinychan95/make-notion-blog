package notion

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsePropTitleItalic(t *testing.T) {
	// Arrange
	properties := "{\"title\":[[\"흥미롭고 유익하고 관심을 끌만한 주제는 무엇이 있을까?\",[[\"i\"]]]]}"
	expected := "_흥미롭고 유익하고 관심을 끌만한 주제는 무엇이 있을까?_"

	// Act
	actual := ParsePropTitle(properties)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestParsePropTitleBoldAndBackground(t *testing.T) {
	// Arrange
	properties := "{\"title\":[[\"Notion 에 편하게 글 적고 알아서 블로그에 반영이 된다면?\",[[\"b\"],[\"h\",\"red_background\"]]]]}"
	expected := "**Notion 에 편하게 글 적고 알아서 블로그에 반영이 된다면?**"

	// Act
	actual := ParsePropTitle(properties)

	// Assert
	assert.Equal(t, expected, actual)
}
