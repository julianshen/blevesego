package blevesego

import (
	"errors"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
	"github.com/huichen/sego"
)

type SegoTokenizer struct {
	handle *sego.Segmenter
}

func NewSegoTokenizer(dictpath string) *SegoTokenizer {
	segmenter := sego.Segmenter{}
	segmenter.LoadDictionary(dictpath)
	return &SegoTokenizer{&segmenter}
}

func convertToken(s *sego.Segment, pos int) *analysis.Token {
	token := analysis.Token{
		Term:     []byte(s.Token().Text()),
		Start:    s.Start(),
		End:      s.End(),
		Position: pos,
		Type:     analysis.Ideographic,
	}
	return &token
}

func (x *SegoTokenizer) Tokenize(sentence []byte) analysis.TokenStream {
	result := make(analysis.TokenStream, 0)
	pos := 1
	words := x.handle.Segment(sentence)

	for _, word := range words {
		result = append(result, convertToken(&word, pos))
		pos++

		for _, seg := range word.Token().Segments() {
			result = append(result, convertToken(seg, pos))
			pos++
		}
	}
	return result
}

func NewTokenizer(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	dictpath, ok := config["dictpath"].(string)
	if !ok {
		return nil, errors.New("config dictpath not found")
	}

	return NewSegoTokenizer(dictpath), nil
}

func init() {
	registry.RegisterTokenizer(Name, NewTokenizer)
}
