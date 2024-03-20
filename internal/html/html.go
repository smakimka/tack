package html

import (
	"fmt"
)

var template = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>%s</title>
</head>
<body>
  %s
</body>
</html>`

type HTML struct {
	html  string
	title string
	body  string
}

func New() *HTML {
	return &HTML{html: template}
}

func (h *HTML) SetTitle(title string) {
	h.title = title
}

func (h *HTML) SetBody(body string) {
	h.body = body
}

func (h *HTML) String() string {
	return fmt.Sprintf(h.html, h.title, h.body)
}
