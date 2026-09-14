// Package apicontract embeds the public HTTP contract so the running server
// always exposes the same document used to generate the frontend client.
package apicontract

import _ "embed"

// OpenAPI is the canonical Resonance OpenAPI document.
//
//go:embed openapi.json
var OpenAPI []byte
