## 实施任务

- [x] 为当前 `DFAMatcher`、`AssembleMatcher`、`RegexpMatcher` 增加 baseline benchmark。
- [x] 将 `DFAMatcher` 内部实现升级为 `Aho-Corasick` 自动机。
- [x] 保持 `DFAMatcher` 对外类型、构造函数和 `Matcher` 接口兼容。
- [x] 覆盖普通关键词的前缀词、重叠词、中文 rune 替换和 `onlyOne` 场景。
- [x] 为 `AssembleMatcher` 增加 anchor 倒排索引预筛。
- [x] 覆盖 `A|B|C` 顺序规则和 `A#B` 排除规则的预筛正确性。
- [x] 为 `RegexpMatcher` 增加固定 literal 预筛或分组策略。
- [x] 保留无法提取 anchor 的正则 fallback 匹配逻辑。
- [x] 更新 `README.md` 中的性能说明和适用规模建议。

## 验证任务

- [x] 执行 `gofmt` 格式化变更文件。
- [x] 执行 `go test -count=1 ./...`。
- [x] 执行 `go test -bench=. -benchmem ./...` 并记录优化后结果。
- [x] 检查编辑文件无 linter 诊断。

## Benchmark 结果

执行环境：`darwin/arm64`，`Apple M1 Pro`。

- `BenchmarkDFAMatcherBuild/words_1000`：`226286 ns/op`，`159233 B/op`，`2702 allocs/op`。
- `BenchmarkDFAMatcherBuild/words_10000`：`2726424 ns/op`，`1660969 B/op`，`26705 allocs/op`。
- `BenchmarkDFAMatcherBuild/words_100000`：`27546366 ns/op`，`18383703 B/op`，`266712 allocs/op`。
- `BenchmarkDFAMatcherMatch/words_1000`：`14977 ns/op`，`9672 B/op`，`6 allocs/op`。
- `BenchmarkDFAMatcherMatch/words_10000`：`15243 ns/op`，`9672 B/op`，`6 allocs/op`。
- `BenchmarkDFAMatcherMatchOnlyOne`：`11657 ns/op`，`9616 B/op`，`4 allocs/op`。
- `BenchmarkDFAMatcherCommonPrefix`：`9259 ns/op`，`4936 B/op`，`6 allocs/op`。
- `BenchmarkAssembleMatcherMatch/rules_10000`：`187017 ns/op`，`208 B/op`，`3 allocs/op`。
- `BenchmarkRegexpMatcherMatch/rules_10000`：`170417 ns/op`，`564 B/op`，`6 allocs/op`。
