package ast

// Create a new document with initialized slices.
// In case you're on a hot path you always want to use a pre-initialized Document.
func ExampleNewDocument() {

	json := []byte(`{
"data": {
  "array": ["hello", "world"],
  "string": "mystring",
  "null", null,
  "integer": 12345,
  "float": 1.2345
}
	`)

	doc := NewDocument()
	doc.Input.ResetInputBytes(json)
}
