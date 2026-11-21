package core_test

import (
	"testing"

	"github.com/Wsine/feishu2md/core"
	"github.com/chyroc/lark"
	"github.com/stretchr/testify/assert"
)

func TestParseDocxCodeWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		content  []string
		expected string
	}{
		{
			name:     "leading spaces",
			content:  []string{"    print('nested')"},
			expected: "```python\n    print('nested')\n```\n",
		},
		{
			name:     "leading tab",
			content:  []string{"\treturn value"},
			expected: "```python\n\treturn value\n```\n",
		},
		{
			name:     "blank lines and trailing spaces",
			content:  []string{"\n  value = 1  \n\n"},
			expected: "```python\n\n  value = 1  \n\n\n```\n",
		},
		{
			name:     "multiple text runs",
			content:  []string{"    if ready:\n", "        run()"},
			expected: "```python\n    if ready:\n        run()\n```\n",
		},
		{
			name:     "unindented code",
			content:  []string{"print('plain')"},
			expected: "```python\nprint('plain')\n```\n",
		},
		{
			name:     "empty code",
			content:  []string{""},
			expected: "```python\n\n```\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := make([]*lark.DocxTextElement, 0, len(tt.content))
			for _, content := range tt.content {
				elements = append(elements, &lark.DocxTextElement{
					TextRun: &lark.DocxTextElementTextRun{Content: content},
				})
			}
			block := &lark.DocxBlock{
				BlockType: lark.DocxBlockTypeCode,
				Code: &lark.DocxBlockText{
					Style:    &lark.DocxTextStyle{Language: lark.DocxCodeLanguagePython},
					Elements: elements,
				},
			}
			parser := core.NewParser(core.NewConfig("", "").Output)
			assert.Equal(t, tt.expected, parser.ParseDocxBlock(block, 0))
		})
	}
}
