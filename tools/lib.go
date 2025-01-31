package jsr_tools

import jsr "github.com/teamdunno/go-jsr-registry"

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
