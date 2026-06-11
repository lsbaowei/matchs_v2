# matcher-performance 规格变更

## ADDED Requirements

### Requirement: 普通关键词匹配必须使用单次文本扫描

`DFAMatcher.Match` 在普通关键词匹配阶段必须避免从文本每个位置重新开始完整 trie 扫描，应该通过 `Aho-Corasick` failure 指针复用已扫描状态。

#### Scenario: 大量普通关键词匹配

- **GIVEN** 已构建 `10k` 条普通关键词
- **WHEN** 对长文本调用 `DFAMatcher.Match(text, false, repl)`
- **THEN** 匹配阶段应以单次文本扫描为主
- **AND** 不应随最长关键词长度呈近似乘法增长

### Requirement: 普通关键词替换语义必须保持兼容

`DFAMatcher.Match` 升级后必须继续返回命中的关键词列表和替换后的文本。

#### Scenario: 前缀词同时命中

- **GIVEN** 已构建普通关键词 `坏` 和 `坏人`
- **WHEN** 匹配文本 `坏人来了`
- **THEN** `sensitiveWords` 必须包含 `坏` 和 `坏人`
- **AND** `replaceText` 必须返回 `**来了`

#### Scenario: `onlyOne=true` 时首次命中立即返回

- **GIVEN** 已构建普通关键词 `敏感词`
- **WHEN** 调用 `DFAMatcher.Match("敏感词和敏感词", true, '*')`
- **THEN** `sensitiveWords` 必须只包含一个 `敏感词`
- **AND** `replaceText` 必须返回 `***和敏感词`

### Requirement: 组合规则必须支持候选预筛且不得漏报

`AssembleMatcher` 可以使用 anchor 倒排索引减少候选规则，但预筛不得改变 `A|B|C` 和 `A#B` 的命中语义。

#### Scenario: anchor 不存在时跳过规则检查

- **GIVEN** 已构建组合规则 `A|B|C`
- **WHEN** 匹配文本 `xxBxxC`
- **THEN** 不应执行该规则的完整顺序检查
- **AND** 不得返回 `A|B|C`

#### Scenario: anchor 存在时执行完整规则检查

- **GIVEN** 已构建组合规则 `A|B|C`
- **WHEN** 匹配文本 `AxxBxxC`
- **THEN** 必须返回规则 `A|B|C`

### Requirement: 正则规则必须支持轻量预筛且保留 fallback

`RegexpMatcher` 可以为简单正则提取固定 literal 作为 anchor，但无法提取 anchor 的正则必须保留在 fallback 列表中继续匹配。

#### Scenario: literal anchor 不存在时跳过正则执行

- **GIVEN** 已构建可提取 anchor 的正则规则
- **WHEN** 文本不包含该 anchor
- **THEN** 可以跳过该正则的 `MatchString` 执行

#### Scenario: 无法提取 anchor 的正则仍可命中

- **GIVEN** 已构建无法提取 anchor 的正则规则 `\d{11}`
- **WHEN** 匹配文本 `手机号13800138000`
- **THEN** 必须返回规则 `reg@\d{11}`

### Requirement: 必须提供性能基准测试

本次性能优化必须新增 benchmark，用于观察大量规则下 `Build` 和 `Match` 的性能趋势。

#### Scenario: 运行 benchmark

- **GIVEN** 已完成性能优化
- **WHEN** 执行 `go test -bench=. -benchmem ./...`
- **THEN** 必须输出普通关键词、组合规则和正则规则相关 benchmark
- **AND** benchmark 名称必须能体现规则规模或匹配场景
