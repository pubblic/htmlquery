package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"github.com/spf13/pflag"
	"golang.org/x/net/html"
)

func compileExprs(exprStrings ...string) ([]*xpath.Expr, error) {
	exprs := make([]*xpath.Expr, len(exprStrings))
	for i, exprString := range exprStrings {
		expr, err := xpath.Compile(exprString)
		if err != nil {
			return nil, err
		}
		exprs[i] = expr
	}
	return exprs, nil
}

func find(node *html.Node, expr *xpath.Expr) []*html.Node {
	nav := htmlquery.CreateXPathNavigator(node)
	iter := expr.Select(nav)
	var founds []*html.Node
	for iter.MoveNext() {
		found := currentNode(nav)
		if len(founds) > 0 {
			first := founds[0]
			if first == found {
				continue
			}
			if nav.NodeType() == xpath.AttributeNode &&
				nav.LocalName() == first.Data &&
				nav.Value() == htmlquery.InnerText(first) {
				continue
			}
		}
		founds = append(founds, found)
	}
	return founds
}

func currentNode(nav *htmlquery.NodeNavigator) *html.Node {
	if nav.NodeType() == xpath.AttributeNode {
		textNode := &html.Node{
			Type: html.TextNode,
			Data: nav.Value(),
		}
		return &html.Node{
			Type:       html.ElementNode,
			Data:       nav.LocalName(),
			FirstChild: textNode,
			LastChild:  textNode,
		}
	}
	return nav.Current()
}

func printNodes(w io.Writer, nodes []*html.Node) error {
	for _, node := range nodes {
		err := printNode(w, node)
		if err != nil {
			return err
		}
	}
	return nil
}

func printNode(w io.Writer, node *html.Node) error {
	err := html.Render(w, node)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w)
	if err != nil {
		return err
	}
	return nil
}

func parse() error {
	pflag.Parse()
	docNode, err := htmlquery.Parse(os.Stdin)
	if err != nil {
		return err
	}
	exprStrings := pflag.Args()
	exprs, err := compileExprs(exprStrings...)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.Grow(4096)
	for _, expr := range exprs {
		nodes := find(docNode, expr)
		err = printNodes(&buf, nodes)
		if err != nil {
			return err
		}
	}
	buf.WriteTo(os.Stdout)
	return nil
}

func main() {
	err := parse()
	if err != nil {
		log.Fatal(err)
	}
}
