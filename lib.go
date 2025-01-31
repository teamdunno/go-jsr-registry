package jsr

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"regexp"
)

const (
	// Set `Client.Protocol` to HTTPS
	ProtocolHTTPS uint8 = 0
	// Set `Client.Protocol` to HTTP
	ProtocolHTTP          uint8 = 1
	DependencyTypeStatic        = "static"
	DependencyTypeDynamic       = "dynamic"
	DependencyKindImport        = "import"
	DependencyKindExport        = "export"
)

type PackageMetaVersion struct {
	// Is this version deprecated?
	Yanked bool `json:"yanked"`
}
type PackageMetaVersions = map[string]PackageMetaVersion
type PackageMeta struct {
	// Package scope. Ex; `dunno` from "@dunno/ono"
	Scope string `json:"scope"`
	// Package name. Ex; `lexicon` from "@dunno/lexicon"
	Name string `json:"name"`
	// Latest version of the package
	Latest string `json:"latest"`
	// List of version
	Versions PackageMetaVersions `json:"versions"`
}

// Package Meta option
//
// Note: Everything is REQUIRED
type PackageMetaOption struct {
	// Package scope. Ex; `dunno` from "@dunno/ono"
	Name string
	// Package name. Ex; `lexicon` from "@dunno/lexicon"
	Scope string
}
type PackageModuleGraph1 = map[string]interface{}
type ModuleGraph2DependencySpecifierRange = [][]uint
type ModuleGraph2Dependency struct {
	// Import type. Usually this can be found if you are using the `import` keyword or the `import()` dynamic function
	//
	// or on export, this can be found on `export` keyword or `module.exports`
	//
	// Available enums
	//
	// go-jsr-registry.DependencyTypeStatic: "static"
	//
	// go-jsr-registry.DependencyTypeDynamic: "dynamic"
	Type string `json:"type"`
	// Dependency kind
	//
	// Available enums
	//
	// go-jsr-registry.DependencyKindImport: "import"
	//
	// go-jsr-registry.DependencyKindExport: "export"
	Kind string `json:"kind"`
	// Dependemcy specifier
	Specifier string `json:"specifier"`
	// Specifier range from the file
	//
	// note: the type should be like this
	//
	// [[uint, uint], [uint, uint]]
	SpecifierRange ModuleGraph2DependencySpecifierRange `json:"specifierRange"`
}
type ModuleGraph2Dependencies = []ModuleGraph2Dependency
type ModuleGraph2 struct {
	/** Dependencies that are used in the file */
	Dependencies ModuleGraph2Dependencies `json:"dependecies,omitempty"`
}

type PackageModuleGraph2 = map[string]ModuleGraph2
type Manifest struct {
	// File size, as bytes
	Size uint
	// SHA256 checksum from the file
	Checksum string
}
type PackageManifest = map[string]Manifest
type Package struct {
	// Files (including manifest) in the package
	Manifest PackageManifest `json:"manifest"`
	// Module graph for javascript-related modules
	//
	// Since `moduleGraph1` was quickly deprecated, use `moduleGraph2` instead
	ModuleGraph1 PackageModuleGraph1 `json:"moduleGraph1,omitempty"`
	// Module graph for javascript-related modules
	ModuleGraph2 PackageModuleGraph2 `json:"moduleGraph2,omitempty"`
}

// Package option
//
// Note: Everything is REQUIRED
type PackageOption struct {
	// Package scope. Ex; `dunno` from "@dunno/ono"
	Name string
	// Package name. Ex; `lexicon` from "@dunno/lexicon"
	Scope string
	// Package version
	Version string
}

// Client option
type ClientOption struct {
	// HTTP Protocol
	//
	// Available enums
	//
	// go-jsr-registry.ProtocolHTTPS: 0 (default value for Protocol)
	//
	// go-jsr-registry.ProtocolHTTP: 1
	Protocol uint8
	// Hostname (defaults to `jsr.io`).
	// You can change it to your self-hosted JSR registry
	//
	// (note: if you want to self-host, you can see the official readme on "[Self-host your own JSR registry].")
	//
	// [Self-host your own JSR registry]: https://github.com/jsr-io/jsr?tab=readme-ov-file#getting-started-entire-stack
	Hostname string
	// Http client used on the Client
	HttpClient http.Client
	// Middleware (for each request)
	Middleware func(*http.Request)
}

// Client struct. Not recommended to create new instance, directly.
//
// use `go-jsr-registry.NewClient` instead
type Client struct {
	// HTTP Protocol
	//
	// Available enums
	//
	// go-jsr-registry.ProtocolHTTPS: 0 (default value)
	//
	// go-jsr-registry.ProtocolHTTP: 1
	Protocol uint8
	// Hostname (default to `jsr.io`).
	// You can change it to your self-hosted JSR registry
	//
	// (note: if you want to self-host, you can see the official readme on "[Self-host your own JSR registry].")
	//
	// [Self-host your own JSR registry]: https://github.com/jsr-io/jsr?tab=readme-ov-file#getting-started-entire-stack
	Hostname string
	// Http client used on the Client
	httpClient http.Client
	// Middleware (for each request)
	Middleware func(*http.Request)
}

func newJSONRequest[T interface{}](c *Client, intr T, method string, path string) (*T, error) {
	if c == nil {
		return nil, errors.New("cant find pointer for go-jsr-registry.Client")
	}
	var Protocol = "https"
	if c.Protocol != 0 {
		Protocol = "http"
	}
	req, err := http.NewRequest(method, Protocol+"://"+c.Hostname+path, nil)
	if err != nil {
		return nil, err
	}
	if c.Middleware != nil {
		c.Middleware(req)
	}
	req.Header.Add("Accept", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	} else {
		defer res.Body.Close()
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data := &intr
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
func (c *Client) GetPackageMeta(option PackageMetaOption) (*PackageMeta, error) {
	res, err := newJSONRequest(c, PackageMeta{}, "GET", "/@"+option.Scope+"/"+option.Name+"/meta.json")
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (c *Client) GetPackage(option PackageOption) (*Package, error) {
	res, err := newJSONRequest(c, Package{}, "GET", "/@"+option.Scope+"/"+option.Name+"/"+option.Version+"_meta.json")
	if err != nil {
		return nil, err
	}
	return res, nil
}
func NewClient(option ...ClientOption) (*Client, error) {
	var opt *ClientOption = &ClientOption{}
	var Protocol uint8 = ProtocolHTTPS
	var Hostname string = "jsr.io"
	var HttpClient http.Client = http.Client{}
	var Middleware func(*http.Request) = nil
	if !reflect.DeepEqual(opt, ClientOption{}) {
		if opt.Protocol != 0 {
			Protocol = ProtocolHTTP
		}
		if opt.Hostname != "" {
			Hostname = opt.Hostname
		}
		if opt.Middleware != nil {
			Middleware = opt.Middleware
		}
		if !reflect.DeepEqual(opt.HttpClient, http.Client{}) {
			HttpClient = opt.HttpClient
		}
	}
	const hostnameRegex = "^([a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?(\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?)*)?(:([0-9]{1,5}))?$"
	ok, err := regexp.Match(hostnameRegex, []byte(Hostname))
	if err != nil || !ok {
		return nil, errors.New("invalid hostname")
	}
	return &Client{
		Protocol,
		Hostname,
		HttpClient,
		Middleware,
	}, nil
}
