package botlogic

type MessagePart interface {
	Build() map[string]interface{}
}

type TextPart struct {
	Text string
}

func (t *TextPart) Build() map[string]interface{} {
	return map[string]interface{}{
		"type": "text",
		"data": map[string]string{
			"text": t.Text,
		},
	}
}

type MessageBuilder struct {
	parts []MessagePart
}

func (mb *MessageBuilder) AddPart(part MessagePart) {
	mb.parts = append(mb.parts, part)
}

func (mb *MessageBuilder) Build() []map[string]interface{} {
	var message []map[string]interface{}
	for _, part := range mb.parts {
		message = append(message, part.Build())
	}
	return message
}
