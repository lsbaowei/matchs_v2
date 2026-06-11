## 实施任务

- [x] 修复 `DFAMatcher.Match` 的替换文本返回逻辑。
- [x] 修复 `DFAMatcher.Match` 中 `onlyOne` 不生效的问题。
- [x] 修复 `DFANode.AddWord` 对空关键词的处理。
- [x] 修复 `AssembleMatcher` 的 `A|B|C` 顺序匹配逻辑。
- [x] 修复不支持脱敏的 matcher 未返回原文的问题。
- [x] 修复 `MatchService.Match` 的 matcher 执行顺序和替换文本汇总逻辑。
- [x] 新增单元测试覆盖 DFA、组合规则、正则规则和服务层聚合逻辑。
- [x] 增加 `go.mod`，使用 Go 标准库 `regexp` 提供正则能力。

## 验证任务

- [x] 执行 `gofmt` 格式化变更文件。
- [x] 执行 `go test ./...`。
- [x] 检查编辑文件无 linter 诊断。
