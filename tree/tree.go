package tree

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/josa42/js-refactor-cli/utils"
)

var (
	importExp = regexp.MustCompile(`(?m)^import\s+([^;]*\s+from\s+)?["'](?P<path>\..+)["'];$`)
)

type Node struct {
	Path    string
	Nodes   []*Node
	content string
}

func (n Node) Rel(filePath string) string {
	base := filepath.Dir(filePath)
	rel, _ := filepath.Rel(base, n.Path)

	if !strings.HasPrefix(rel, ".") {
		return "./" + rel
	}
	return rel
}

func (n Node) RelTo(filePath string) string {
	base := filepath.Dir(n.Path)
	rel, _ := filepath.Rel(base, filePath)

	if !strings.HasPrefix(rel, ".") {
		return "./" + rel
	}
	return rel
}

func (n *Node) Move(filePath string, nodes map[string]*Node) {

	n.read()
	for _, c := range n.Nodes {
		find := stripExt(c.Rel(n.Path), ".js")
		replace := stripExt(c.Rel(filePath), ".js")

		n.replace(find, replace)
	}
	n.write()

	// Move file
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	os.Rename(n.Path, filePath)

	// update references from other files
	for _, dn := range nodes {
		if dn.DependsOn(n) {
			dn.read()

			find := stripExt(dn.RelTo(n.Path), ".js")
			replace := stripExt(dn.RelTo(filePath), ".js")

			dn.replace(find, replace)
			dn.write()
		}
	}
	n.Path = filePath
}

func (n Node) updateDependency(source, target string) {
	find := stripExt(n.RelTo(source), ".js")
	replace := stripExt(n.RelTo(target), ".js")

	n.replace(find, replace)
}

func (n *Node) read() {
	content, _ := ioutil.ReadFile(n.Path)
	n.content = string(content)
}

func (n *Node) write() {
	ioutil.WriteFile(n.Path, []byte(n.content), 0644)
	n.content = ""
}

func (n *Node) replace(find, replace string) {
	n.content = strings.ReplaceAll(n.content, find, replace)
}

func (n Node) Content() string {
	content, _ := ioutil.ReadFile(n.Path)
	return string(content)
}

func (n Node) DependsOn(node *Node) bool {
	for _, c := range n.Nodes {

		if c.Path == node.Path {
			return true
		}
	}

	return false
}

var ignoreDirs = []string{"node_modules", ".git"}

func GetNodes(rootPath string) map[string]*Node {

	nodes := map[string]*Node{}

	utils.WalkFiles(rootPath, func(path string) {
		GetNode(path, nodes)
	})

	return nodes
}

func GetNode(path string, nodes map[string]*Node) *Node {
	path, _ = filepath.Abs(path)

	if node, ok := nodes[path]; ok {
		return node
	}

	n := &Node{Path: path}
	nodes[path] = n

	dir := filepath.Dir(path)

	if filepath.Ext(path) == ".js" {
		text := n.Content()

		m := [][]string{}

		m = append(m, importExp.FindAllStringSubmatch(string(text), -1)...)

		for _, g := range m {
			ipath := addExt(filepath.Join(dir, g[2]), ".js")

			n.Nodes = append(n.Nodes, GetNode(ipath, nodes))
		}
	}

	return n
}

func stripExt(path, ext string) string {
	if filepath.Ext(path) == ext {
		return regexp.MustCompile(regexp.QuoteMeta(ext)+"$").ReplaceAllString(path, "")
	}
	return path
}

func addExt(path, ext string) string {
	if filepath.Ext(path) == "" {
		return path + ext
	}
	return path
}
