package jsr_tools

import (
	"regexp"

	jsr "github.com/teamdunno/go-jsr-registry"
)

func GetImportsFromDependencies(deps jsr.ModuleGraph2Dependencies) (res jsr.ModuleGraph2Dependencies) {
	if deps == nil {
		return jsr.ModuleGraph2Dependencies{}
	}
	for _, dep := range deps {
		if dep.Kind != "import" {
			continue
		}
		res = append(res, dep)
	}
	return res
}

func GetStaticImportsFromDependencies(deps jsr.ModuleGraph2Dependencies) (res jsr.ModuleGraph2Dependencies) {
	if deps == nil {
		return jsr.ModuleGraph2Dependencies{}
	}
	for _, dep := range deps {
		if !(dep.Kind == "import" && dep.Type == "static") {
			continue
		}
		res = append(res, dep)
	}
	return res
}

func GetDynamicImportsFromDependencies(deps jsr.ModuleGraph2Dependencies) (res jsr.ModuleGraph2Dependencies) {
	if deps == nil {
		return jsr.ModuleGraph2Dependencies{}
	}
	for _, dep := range deps {
		if !(dep.Kind == "import" && dep.Type == "dynamic") {
			continue
		}
		res = append(res, dep)
	}
	return res
}

func GetExportsFromDependencies(deps jsr.ModuleGraph2Dependencies) (res jsr.ModuleGraph2Dependencies) {
	if deps == nil {
		return jsr.ModuleGraph2Dependencies{}
	}
	for _, dep := range deps {
		if dep.Kind != "export" {
			continue
		}
		res = append(res, dep)
	}
	return res
}

func GetStaticExportsFromDependencies(deps jsr.ModuleGraph2Dependencies) (res jsr.ModuleGraph2Dependencies) {
	if deps == nil {
		return jsr.ModuleGraph2Dependencies{}
	}
	for _, dep := range deps {
		if !(dep.Kind == "export" && dep.Type == "static") {
			continue
		}
		res = append(res, dep)
	}
	return res
}

func GetDynamicExportsFromDependencies(deps jsr.ModuleGraph2Dependencies) (res jsr.ModuleGraph2Dependencies) {
	if deps == nil {
		return jsr.ModuleGraph2Dependencies{}
	}
	for _, dep := range deps {
		if !(dep.Kind == "export" && dep.Type == "dynamic") {
			continue
		}
		res = append(res, dep)
	}
	return res
}

func GetYankedVersionsFromPackageMeta(meta jsr.PackageMetaVersions) (res jsr.PackageMetaVersions) {
	if meta == nil {
		return jsr.PackageMetaVersions{}
	}
	res = make(jsr.PackageMetaVersions)
	for name, ver := range meta {
		if !ver.Yanked {
			continue
		}
		res[name] = ver
	}
	return res
}
func GetUnyankedVersionsFromPackageMeta(meta jsr.PackageMetaVersions) (res jsr.PackageMetaVersions) {
	if meta == nil {
		return jsr.PackageMetaVersions{}
	}
	res = make(jsr.PackageMetaVersions)
	for name, ver := range meta {
		if ver.Yanked {
			continue
		}
		res[name] = ver
	}
	return res
}

func HideHiddenObjectsFromManifest(manifest jsr.PackageManifest) (res jsr.PackageManifest, regexError error) {
	if manifest == nil {
		return jsr.PackageManifest{}, nil
	}
	res = make(jsr.PackageManifest)
	hiddenObjRegex, err := regexp.Compile(`^.*(\/(\.[^\/]+)|(_[^\/]*))$`)
	if err != nil {
		return nil, err
	}
	for name, ver := range manifest {
		if !hiddenObjRegex.Match([]byte(name)) {
			continue
		}
		res[name] = ver
	}
	return res, nil
}

func HideNormalObjectsFromManifest(manifest jsr.PackageManifest) (res jsr.PackageManifest, regexError error) {
	if manifest == nil {
		return jsr.PackageManifest{}, nil
	}
	res = make(jsr.PackageManifest)
	normalObjRegex, err := regexp.Compile(`^(\/([^\._][^\/]*)?)*$`)
	if err != nil {
		return nil, err
	}
	for name, ver := range manifest {
		if !normalObjRegex.Match([]byte(name)) {
			continue
		}
		res[name] = ver
	}
	return res, nil
}

func HasJSRJsonInManifest(manifest jsr.PackageManifest) bool {
	_, ok := manifest["/jsr.json"]
	return ok
}

func HasPackageJsonInManifest(manifest jsr.PackageManifest) bool {
	_, ok := manifest["/package.json"]
	return ok
}

func HasDenoJsonInManifest(manifest jsr.PackageManifest) bool {
	_, ok := manifest["/deno.json"]
	if ok {
		return true
	}
	_, ok = manifest["/deno.jsonc"]
	return ok
}

func HasBunfigTomlInManifest(manifest jsr.PackageManifest) bool {
	_, ok := manifest["/bunfig.toml"]
	return ok
}
