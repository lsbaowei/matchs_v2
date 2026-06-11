# matcher-service 规格变更

## ADDED Requirements

### Requirement: `MatchService` 必须按稳定顺序执行 matcher

`MatchService.Match` 必须按 `DFA`、`ASSEMBLE`、`REGEXP` 的顺序执行已构建的 matcher。

#### Scenario: `onlyOne=true` 时返回稳定结果

- **GIVEN** 同时构建普通关键词、组合规则和正则规则
- **WHEN** 调用 `MatchService.Match(text, true, repl)`
- **THEN** 必须优先返回 `DFA` 命中的结果
- **AND** 不得依赖 `map` 遍历顺序

### Requirement: `MatchService` 必须返回最终替换文本

`MatchService.Match` 必须把每个 matcher 返回的 `replaceText` 传递给下一轮 matcher，并返回最终文本。

#### Scenario: DFA 命中后返回脱敏文本

- **GIVEN** 已构建普通关键词 `敏感词`
- **WHEN** 调用 `MatchService.Match("这里有敏感词", false, '*')`
- **THEN** `sensitiveWords` 必须包含 `敏感词`
- **AND** `replaceText` 必须返回 `这里有***`

### Requirement: DFA matcher 必须支持替换和 `onlyOne`

`DFAMatcher.Match` 必须对命中的普通关键词区间执行 rune 级替换，并在 `onlyOne=true` 时首次命中后立即返回。

#### Scenario: 只替换首次命中

- **GIVEN** 已构建普通关键词 `敏感词`
- **WHEN** 调用 `DFAMatcher.Match("敏感词和敏感词", true, '*')`
- **THEN** `sensitiveWords` 必须只包含一个 `敏感词`
- **AND** `replaceText` 必须返回 `***和敏感词`

### Requirement: DFA matcher 必须忽略空关键词

`DFANode.AddWord` 必须忽略空字符串，避免根节点被标记为命中节点。

#### Scenario: 构建空关键词

- **GIVEN** 构建词表包含空字符串
- **WHEN** 匹配任意文本
- **THEN** 不得因为空字符串产生命中结果

### Requirement: 组合规则必须按顺序匹配 `|`

`A|B|C` 规则必须表示 `A`、`B`、`C` 按文本出现顺序依次命中，且允许第一个词从位置 `0` 开始命中。

#### Scenario: 从文本开头开始命中

- **GIVEN** 已构建组合规则 `A|B|C`
- **WHEN** 匹配文本 `AxxBxxC`
- **THEN** 必须返回规则 `A|B|C`

#### Scenario: 顺序不满足时不命中

- **GIVEN** 已构建组合规则 `A|B|C`
- **WHEN** 匹配文本 `BxxAxxC`
- **THEN** 不得返回规则 `A|B|C`

### Requirement: 不支持脱敏的 matcher 必须返回原文

`AssembleMatcher` 和 `RegexpMatcher` 当前不执行脱敏替换，但必须把输入文本作为 `desensitization` 返回。

#### Scenario: 组合规则命中但不脱敏

- **GIVEN** 已构建组合规则 `A#B`
- **WHEN** 匹配文本 `xxAxx`
- **THEN** 必须返回规则 `A#B`
- **AND** `desensitization` 必须等于 `xxAxx`
