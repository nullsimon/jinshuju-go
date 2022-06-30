package jinshuju

import (
	"strings"
)

type Choice struct {
	Name string `json:"name"`

	Value string `json:"value"`

	Hidden bool `json:"hidden"`

	IsOther bool `json:"is_other"`

	OtherValue string
}

type Validation struct {
	Required bool `json:"required"`
}

type MediaType struct {
	Type string `json:"type"`

	Value []string `json:"value"`
}

type Field struct {
	FieldId string

	Label string `json:"label"`

	Type string `json:"type"`

	Notes string `json:"notes"`

	Private bool `json:"private"`

	Validation Validation `json:"validation"`

	MaxFileQuantity int `json:"max_file_quantity"`

	MediaType MediaType `json:"media_type"`

	Precision string `json:"precision"`

	Choices []Choice `json:"choices"`

	ChooseChoices []Choice

	AllowOther bool `json:"allow_other"`
}

type FieldMap map[string]Field

type FieldType1 interface {
	isSingleChoice() bool
	isMultipleChoice() bool
	isAllowOther() bool
	isAttachment() bool
	isDateTime() bool
	isInput() bool
	isSectionBreak() bool
	trimLabel()
}

const (
	//SingleChoiceField 单选
	SingleChoiceField FieldType = iota + 1

	DropDownField
	//MultipleChoiceField 多选
	MultipleChoiceField

	//ParagraphTextField 多行文本
	ParagraphTextField
	//SingleLineTextField 单行文本
	SingleLineTextField
	//MultipleLineTextField 多行文本
	MultipleLineTextField

	DateTimeField
	AttachmentField

	//SectionBreakField 分割线
	SectionBreakField
)

type FieldType int

func (f *Field) isSingleChoice() bool {
	return f.Type == "single_choice" || f.Type == "drop_down"
}

func (f *Field) isMultipleChoice() bool {
	return f.Type == "multiple_choice"
}

func (f *Field) isAllowOther() bool {
	return f.AllowOther
}

func (f *Field) isAttachment() bool {
	return f.Type == "attachment"
}

func (f *Field) isDateTime() bool {
	return f.Type == "date_time"
}

func (f *Field) isInput() bool {
	return f.Type == "single_line_text" || f.Type == "paragraph_text" || f.Type == "multiple_line_text"
}

func (f *Field) isSectionBreak() bool {
	return f.Type == "section_break"
}

func (f *Field) trimLabel(label string) {
	f.Label = strings.TrimSpace(f.Label)
}

type Entry struct {
	SerialNumber        int     `json:"serial_number"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
	InfoFillingDuration float64 `json:"info_filling_duration"`
	CreatorName         string  `json:"creator_name"`
	Fields              []FieldMap
}
