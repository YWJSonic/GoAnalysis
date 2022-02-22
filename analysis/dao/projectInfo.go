package dao

import (
	"codeanalysis/analysis/constant"
	"codeanalysis/load/project/goloader"
	"codeanalysis/types"
)

func NewProjectInfo(node *goloader.GoFileNode) *ProjectInfo {
	info := &ProjectInfo{
		LocalPackageMap:  make(map[string]*PackageInfo),
		VendorPackageMap: make(map[string]*PackageInfo),
		GolangPackageMap: make(map[string]*PackageInfo),
	}
	info.name = node.Path()
	info.ProjectRoot = node
	return info
}

type ProjectInfo struct {
	PointBase
	ModuleInfo       *ModuleInfo
	ProjectRoot      FileDataNode
	LocalPackageMap  map[string]*PackageInfo // 內部實做 package <packagePath, *PackageInfo>
	VendorPackageMap map[string]*PackageInfo // 外部引用 package <packagePath, *PackageInfo>
	GolangPackageMap map[string]*PackageInfo // 系統自帶 package <packagePath, *PackageInfo>
}

// 讀寫此 package 關連資料
//
// @params string		package 路徑
// @params *PackageInfo	預設關聯結構
//
// @return *PackageInfo	回傳的關聯資料
// @return bool			是否已存在資料
func (self *ProjectInfo) LoadOrStoryPackage(packageType types.TypeFrom, pwd string, info *PackageInfo) (*PackageInfo, bool) {

	switch packageType {
	case constant.From_Golang:
		if packageInfo, ok := self.GolangPackageMap[pwd]; ok {
			return packageInfo, true
		}
		self.GolangPackageMap[pwd] = info

	case constant.From_Local:
		if packageInfo, ok := self.LocalPackageMap[pwd]; ok {
			return packageInfo, true
		}
		self.LocalPackageMap[pwd] = info

	case constant.From_Vendor:
		if packageInfo, ok := self.VendorPackageMap[pwd]; ok {
			return packageInfo, true
		}
		self.VendorPackageMap[pwd] = info
	}

	return info, false
}

func (self *ProjectInfo) Output() string {
	data := ""
	packageOutput := packageOutput{}
	for _, packageInfo := range self.LocalPackageMap {

		// import data
		for _, importData := range packageInfo.AllImportLink {
			packageOutput.ImportInfo = append(packageOutput.ImportInfo, importInfo{
				RefPath: importData.Path,
			})
		}

		// const data
		for _, constData := range packageInfo.AllConstInfos {
			var outputData = outputData(constData.(*ConstInfo).ContentTypeInfo)
			packageOutput.ConstInfo = append(packageOutput.ConstInfo, outputData)
		}

		// var data
		for _, varData := range packageInfo.AllVarInfos {
			var outputData = outputData(varData.(*VarInfo).ContentTypeInfo)
			packageOutput.VarInfo = append(packageOutput.VarInfo, outputData)
		}
	}

	return data
}
func outputData(itype ITypeInfo) (outputData output) {
	switch info := itype.(type) {
	case *TypeInfoBool, *TypeInfoNumeric, *TypeInfoString:
		outputData = output{
			RefPaths: []string{info.GetTypeName()},
		}

	case *TypeInfoArray:
		outputData = output{
			RefPaths: []string{info.GetTypeName()},
		}

	case *TypeInfoChannel:
		outputData = output{
			RefPaths: []string{info.GetTypeName(), info.ContentTypeInfo.GetTypeName()},
		}

	case *TypeInfoFunction:
		outputData = output{
			RefPaths: []string{info.GetTypeName()},
		}

		for _, inparams := range info.ParamsInPoint {
			outputData.RefPaths = append(outputData.RefPaths, inparams.GetName())
		}

		for _, outparams := range info.ParamsOutPoint {
			outputData.RefPaths = append(outputData.RefPaths, outparams.GetName())
		}

	case *TypeInfoInterface:
		outputData = output{
			RefPaths: []string{info.GetTypeName()},
		}

		for _, imatchInfo := range info.MatchInfos {
			matchInfo := imatchInfo.(*FuncInfo)

			for _, inparams := range matchInfo.ParamsInPoint {
				outputData.RefPaths = append(outputData.RefPaths, inparams.GetName())
			}

			for _, outparams := range matchInfo.ParamsOutPoint {
				outputData.RefPaths = append(outputData.RefPaths, outparams.GetName())
			}
		}

	case *TypeInfoMap:
		outputData = output{
			RefPaths: []string{info.GetTypeName(), info.KeyType.GetTypeName(), info.ValueType.GetTypeName()},
		}
	case *TypeInfoPointer:
		outputData = output{
			RefPaths: []string{info.GetTypeName()},
		}

	case *TypeInfoQualifiedIdent:
		outputData = output{
			RefPaths:     []string{info.GetTypeName()},
			ReceiverPath: info.ImportLink.Path,
		}

	case *TypeInfoSlice:
		outputData = output{
			RefPaths: []string{info.ContentTypeInfo.GetTypeName()},
		}
	case *TypeInfoStruct:
		outputData = output{
			RefPaths: []string{info.GetTypeName()},
		}

		for _, implicitlyVarInfo := range info.ImplicitlyVarInfos {
			outputData.RefPaths = append(outputData.RefPaths, implicitlyVarInfo.ContentTypeInfo.GetTypeName())
		}

		for _, varInfo := range info.VarInfos {
			outputData.RefPaths = append(outputData.RefPaths, varInfo.ContentTypeInfo.GetTypeName())
		}
	}
	return
}

type projectOutput struct {
	ModuleName         string
	GoVersion          string
	VendorMap          map[string]struct{}
	VendorPackageInfo  map[string]packageOutput
	ProjectPackageInfo map[string]packageOutput
}
type packageOutput struct {
	ImportInfo []importInfo
	ConstInfo  []output // <name, ConstInfo>
	VarInfo    []output // <name, VarInfo>
	TypeInfo   []output // <name, TypeInfo>
}

type importInfo struct {
	RefPath string // 其他包引用路徑
}

type typeOutput struct {
	VarInfo map[string]varOutput // <name, TypeOutput>
}

type output struct {
	ReceiverPath string
	RefPaths     []string
}

type constOutput struct {
	RefPath  string // 其他包引用路徑
	TypeName string // 參數類型名稱
	Name     string // 參數名稱
}

type varOutput struct {
	RefPath  string // 其他包引用路徑
	TypeName string // 變數類型名稱
	Name     string // 參數名稱
}
