package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/pflag"
	"golang.org/x/net/html"
)

var useDump bool
var dumpState spew.ConfigState
var firstN int
var printNumber bool

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

func dump(w io.Writer, node *html.Node) error {
	dumpState.Fdump(w, node)
	_, err := fmt.Println(w)
	if err != nil {
		return err
	}
	return nil
}

func printNode(w io.Writer, node *html.Node) error {
	if useDump {
		return dump(w, node)
	}
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

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func parse() error {
	pflag.StringVar(&dumpState.Indent, "dump-indent", "  ",
		"spew.ConfigState.Indent")
	pflag.IntVar(&dumpState.MaxDepth, "dump-max-depth", 0,
		"spew.ConfigState.MaxDepth")
	pflag.BoolVar(&dumpState.DisableMethods, "dump-disable-methods", false,
		"spew.ConfigState.DisableMethods")
	pflag.BoolVar(&dumpState.DisablePointerMethods, "dump-disable-pointer-methods", false,
		"spew.ConfigState.DisablePointerMethods")
	pflag.BoolVar(&dumpState.DisablePointerAddresses, "dump-disable-pointer-addresses", false,
		"spew.ConfigState.DisablePointerAddresses")
	pflag.BoolVar(&dumpState.DisableCapacities, "dump-disable-capacities", false,
		"spew.ConfigState.DisableCapacities")
	pflag.BoolVar(&dumpState.ContinueOnMethod, "dump-continue-on-method", false,
		"spew.ConfigState.ContinueOnMethod")
	pflag.BoolVar(&dumpState.SortKeys, "dump-sort-keys", false,
		"spew.ConfigState.SortKeys")
	pflag.BoolVar(&dumpState.SpewKeys, "dump-spew-keys", false,
		"spew.ConfigState.SpewKeys")
	pflag.BoolVarP(&useDump, "dump", "d", false, "dump node")
	pflag.IntVarP(&firstN, "n", "n", -1, "select first n node(s)")
	pflag.BoolVar(&printNumber, "print-number", false,
		"print the number of the matched nodes")
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
	if firstN >= 0 {
		exprs = exprs[:min(firstN, len(exprs))]
	}
	if printNumber {
		fmt.Println(len(exprs))
		return nil
	}
	for _, expr := range exprs {
		nodes := find(docNode, expr)
		err = printNodes(os.Stdout, nodes)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	err := parse()
	if err != nil {
		log.Fatal(err)
	}
}
