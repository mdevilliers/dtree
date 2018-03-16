package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/go-yaml/yaml"
)

/*
go build && ./dtree > example.dot && dtree dot -Tsvg example.dot > example.svg
*/
func main() {

	root := "/home/mark/go/src/github.com/karhoo/"

	repos := []Repository{
		{Path: "api-v1-bookings"},
		{Path: "svc-driver-tracking"},
		{Path: "svc-trips"},
		{Path: "svc-estimates"},
		/*		{Path: "svc-addresses"},
				{Path: "svc-availability"},
				{Path: "svc-quotes"},
				{Path: "svc-supply"},
				{Path: "svc-supply-configuration"},
				{Path: "svc-timezones"},
				{Path: "svc-zones"},
				{Path: "svc-users"},
				{Path: "svc-trips"},
				{Path: "svc-qta"},
				{Path: "svc-zones"},
				{Path: "svc-predictions"},
				{Path: "svc-poi"},
				{Path: "svc-roads"},*/
		//{Path: "svc-driver-photos"},// godep
	}

	// api-v1-quotes
	/*
		repos := []Repository{
			{Path: "api-v1-quotes"},
			{Path: "svc-addresses"},
			{Path: "svc-availability"},
			{Path: "svc-estimates"},
			{Path: "svc-quotes"},
			{Path: "svc-supply"},
			{Path: "svc-supply-configuration"},
			{Path: "svc-timezones"},
			{Path: "svc-zones"},
			{Path: "svc-users"},
			{Path: "svc-trips"},
			{Path: "svc-qta"},
			{Path: "svc-zones"},
			{Path: "svc-predictions"},
			{Path: "svc-poi"},
			{Path: "svc-roads"},
			//{Path: "svc-driver-photos"},// godep
		}*/
	//all
	/*
		repos := []Repository{
			{Path: "svc-users"},
			{Path: "svc-addresses"},
			{Path: "svc-driver-tracking"},
			{Path: "svc-estimates"},
			//{Path: "svc-money"},// python
			{Path: "svc-quotes"},
			{Path: "svc-supply"},
			{Path: "svc-supply-configuration"},
			{Path: "svc-timezones"},
			{Path: "svc-trips"},
			{Path: "svc-trips-scheduler"},
			{Path: "api-v1-bookings"},
			{Path: "api-v1-directory"},
			{Path: "api-v1-quotes"},
			{Path: "lib-api"},
			{Path: "lib-common"},
			//{Path: "lib-permissions"},// godep
			//{Path: "lib-proto"}, // no deps
		}
	*/
	allNodes := map[string]Node{}
	allEdges := []Edge{}

	for _, r := range repos {
		nodes, edges, err := glide{}.Parse(path.Join(root, r.Path))

		if err != nil {
			panic(err)
		}

		allEdges = append(allEdges, edges...)

		for _, n := range nodes {
			allNodes[n.Name] = n // WARNING : this just overrides Nodes - LWW
		}

	}

	nodesArr := []Node{}

	for _, n := range allNodes {
		nodesArr = append(nodesArr, n)
	}

	postProcessNodes(nodesArr)

	err := writeDot(nodesArr, allEdges, os.Stdout)

	if err != nil {
		panic(err)
	}
}

type Repository struct {
	Path string
}

type Node struct {
	Name   string
	Labels Attributes
}

type Edge struct {
	Source string
	Target string
	Labels Attributes
}

type Attributes map[string]interface{}

func writeDot(nodes []Node, edges []Edge, writer io.Writer) error {

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("digraph G {\n")

	for _, edge := range edges {
		buf.WriteString(fmt.Sprintf(`"%s"->"%s";`, edge.Source, edge.Target))
		buf.WriteByte('\n')
	}
	for _, node := range nodes {

		fillcolor := "white"

		typez := node.Labels[typeStr]

		switch typez {
		case KarhooAPI:
			fillcolor = "blue"
		case KarhooSvc:
			fillcolor = "red"
		case KarhooLibrary:
			fillcolor = "green"

		}

		buf.WriteString(fmt.Sprintf(`"%s" [fillcolor=%s style=filled];`, node.Name, fillcolor))
		buf.WriteByte('\n')
	}
	buf.WriteString("}\n")
	_, err := writer.Write(buf.Bytes())
	return err
}

type glideFile struct {
	Package string        `yaml:"package"`
	Imports []glideImport `yaml:"import"`
}
type glideImport struct {
	Package string `json:"package"`
	Version string `json:"version"`
}

var (
	versionStr           = "version"
	relationshipStr      = "relationship"
	typeStr              = "type"
	versionUnknown       = "UNKNOWN - NOT STATED"
	relationshipDirect   = "direct"
	relationshipInDirect = "indirect"
)

type glide struct{}

func (g glide) Parse(pth string) ([]Node, []Edge, error) {

	pathToGlideFile := path.Join(pth, "glide.yaml")

	data, err := ioutil.ReadFile(pathToGlideFile)

	if err != nil {
		return nil, nil, err
	}

	glideFile := glideFile{}

	err = yaml.Unmarshal([]byte(data), &glideFile)

	if err != nil {
		return nil, nil, fmt.Errorf("cannot unmarshal data: %v", err)
	}

	nodes := []Node{
		{Name: glideFile.Package, Labels: Attributes{}},
	}

	edges := []Edge{}

	for _, p := range glideFile.Imports {

		version := p.Version

		if version == "" {
			version = versionUnknown
		}

		node := Node{
			Name:   p.Package,
			Labels: Attributes{},
		}

		edge := Edge{
			Source: glideFile.Package,
			Target: p.Package,
			Labels: Attributes{relationshipStr: relationshipDirect, versionStr: version},
		}

		nodes = append(nodes, node)
		edges = append(edges, edge)
	}

	return nodes, edges, nil
}

type dependancyType string

var (
	Default       = dependancyType("DEFAULT")
	KarhooLibrary = dependancyType("LIB")
	KarhooSvc     = dependancyType("SVC")
	KarhooAPI     = dependancyType("API")
)

func isKarhooSvc(name string) bool {
	return strings.Contains(name, "svc") && strings.Contains(name, "karhoo")
}

func isKarhooLib(name string) bool {
	return strings.Contains(name, "lib") && strings.Contains(name, "karhoo")
}

func isKarhooAPI(name string) bool {
	return strings.Contains(name, "api-v1") && strings.Contains(name, "karhoo")
}

func postProcessNodes(nodes []Node) {

	for _, n := range nodes {
		if isKarhooSvc(n.Name) {
			n.Labels[typeStr] = KarhooSvc
		} else if isKarhooLib(n.Name) {
			n.Labels[typeStr] = KarhooLibrary
		} else if isKarhooAPI(n.Name) {
			n.Labels[typeStr] = KarhooAPI
		} else {
			n.Labels[typeStr] = Default
		}
	}
}

type dep struct{}

func (d dep) Parse(path string) ([]Node, []Edge, error) {
	return nil, nil, errors.New("not implemented")
}
