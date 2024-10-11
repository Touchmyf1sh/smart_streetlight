//依赖管理文件

module smart_streetlight

go 1.23.0

require (
	github.com/goburrow/modbus v0.1.0 // indirect
	github.com/goburrow/serial v0.1.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
)

//撤回
//当前项目作为依赖被别的项目引用，如果某个版本出现了问题时
retract v1.0.0
