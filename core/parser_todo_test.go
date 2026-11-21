package core_test

import (
	"testing"

	"github.com/Wsine/feishu2md/core"
	"github.com/chyroc/lark"
	"github.com/stretchr/testify/assert"
)

func taskBlock(id, content string, done bool, children ...string) *lark.DocxBlock {
	return &lark.DocxBlock{
		BlockID:   id,
		BlockType: lark.DocxBlockTypeTodo,
		Children:  children,
		Todo: &lark.DocxBlockText{
			Style: &lark.DocxTextStyle{Done: done},
			Elements: []*lark.DocxTextElement{
				{TextRun: &lark.DocxTextElementTextRun{Content: content}},
			},
		},
	}
}

func TestParseDocxTaskList(t *testing.T) {
	tests := []struct {
		name     string
		blocks   []*lark.DocxBlock
		expected string
	}{
		{
			name: "leaf task",
			blocks: []*lark.DocxBlock{
				taskBlock("parent", "Only task", false),
			},
			expected: "- [ ] Only task\n",
		},
		{
			name: "nested checked and unchecked tasks",
			blocks: []*lark.DocxBlock{
				taskBlock("parent", "Parent", false, "first", "second"),
				taskBlock("first", "First child", true),
				taskBlock("second", "Second child", false, "grandchild"),
				taskBlock("grandchild", "Grandchild", true),
			},
			expected: "- [ ] Parent\n\t- [x] First child\n\t- [ ] Second child\n\t\t- [x] Grandchild\n",
		},
		{
			name: "checked parent preserves child order",
			blocks: []*lark.DocxBlock{
				taskBlock("parent", "Done parent", true, "second", "first"),
				taskBlock("first", "First child", false),
				taskBlock("second", "Second child", true),
			},
			expected: "- [x] Done parent\n\t- [x] Second child\n\t- [ ] First child\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Populate the parser's block map through the same entry point as an export.
			page := &lark.DocxBlock{
				BlockID:   "page",
				BlockType: lark.DocxBlockTypePage,
				Children:  []string{"parent"},
				Page: &lark.DocxBlockText{
					Elements: []*lark.DocxTextElement{
						{TextRun: &lark.DocxTextElementTextRun{Content: "Tasks"}},
					},
				},
			}
			parser := core.NewParser(core.NewConfig("", "").Output)
			got := parser.ParseDocxContent(&lark.DocxDocument{DocumentID: "page"}, append(tt.blocks, page))
			assert.Equal(t, "# Tasks\n\n"+tt.expected+"\n", got)
		})
	}
}
