package trie

import (
	json "encoding/json"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonDecodeGithubComWolanmSearchEnginePkgTrie(in *jlexer.Lexer, out *Node) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "is_end":
			out.IsEnd = bool(in.Bool())
		case "children":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Children).UnmarshalJSON(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonEncodeGithubComWolanmSearchEnginePkgTrie(out *jwriter.Writer, in Node) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"is_end\":"
		out.RawString(prefix[1:])
		out.Bool(bool(in.IsEnd))
	}
	{
		const prefix string = ",\"children\":"
		out.RawString(prefix)
		out.Raw((in.Children).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Node) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEncodeGithubComWolanmSearchEnginePkgTrie(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Node) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEncodeGithubComWolanmSearchEnginePkgTrie(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Node) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDecodeGithubComWolanmSearchEnginePkgTrie(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Node) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDecodeGithubComWolanmSearchEnginePkgTrie(l, v)
}
func easyjsonDecodeGithubComWolanmSearchEnginePkgTrie1(in *jlexer.Lexer, out *Trie) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "root":
			if in.IsNull() {
				in.Skip()
				out.Root = nil
			} else {
				if out.Root == nil {
					out.Root = new(Node)
				}
				(*out.Root).UnmarshalEasyJSON(in)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonEncodeGithubComWolanmSearchEnginePkgTrie1(out *jwriter.Writer, in Trie) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"root\":"
		out.RawString(prefix[1:])
		if in.Root == nil {
			out.RawString("null")
		} else {
			(*in.Root).MarshalEasyJSON(out)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Trie) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEncodeGithubComWolanmSearchEnginePkgTrie1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Trie) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEncodeGithubComWolanmSearchEnginePkgTrie1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Trie) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonDecodeGithubComWolanmSearchEnginePkgTrie1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Trie) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonDecodeGithubComWolanmSearchEnginePkgTrie1(l, v)
}
