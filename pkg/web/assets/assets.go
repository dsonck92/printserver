package assets

import _ "embed"

//go:embed bootstrap.min.css
var Bootstrap []byte

//go:embed htmx.min.js
var HTMX []byte
