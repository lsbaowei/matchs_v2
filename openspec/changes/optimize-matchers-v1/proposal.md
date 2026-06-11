# 优化关键词匹配与替换逻辑

## 背景

当前包用于对输入文本执行一批关键词的快速查找与替换，支持普通关键词、组合规则和正则规则三类配置。现有实现存在部分行为不稳定或结果不完整的问题，会影响调用方对匹配结果和替换文本的使用。

## 问题

- `MatchService.Match` 没有正确汇总 matcher 返回的替换文本，调用方无法获得脱敏后的内容。
- `MatchService.Match` 使用 `map` 遍历 matcher，执行顺序不稳定，`onlyOne` 场景下返回结果可能不一致。
- `DFAMatcher.Match` 未启用替换逻辑，且 `onlyOne` 参数未生效。
- `DFANode.AddWord` 未忽略空关键词，可能导致根节点被标记为命中节点。
- `AssembleMatcher` 中 `A|B|C` 的顺序匹配逻辑存在位置 `0` 漏匹配和重复文本场景不准确的问题。
- `AssembleMatcher` 和 `RegexpMatcher` 不支持替换时没有稳定返回原文，可能导致服务层替换结果被清空。

## 目标

- 保持现有公开接口 `Matcher` 和 `MatchService` 不变。
- 修复普通关键词、组合规则和正则规则的匹配结果稳定性。
- 支持 DFA 匹配后的替换文本返回。
- 让 `onlyOne` 在服务层和 DFA 匹配器中行为一致。
- 增加单元测试覆盖核心匹配行为。

## 非目标

- 不改变 `reg@`、`|`、`#` 的规则配置格式。
- 不为组合规则和正则规则新增脱敏替换能力。
- 不重构为 Aho-Corasick 等新的匹配算法。
- 不改变 `Matcher.Build` 的方法签名。

## 方案

- 在 `MatchService` 中使用固定 matcher 顺序：`DFA`、`ASSEMBLE`、`REGEXP`。
- `MatchService.Match` 将每个 matcher 的输入设为上一轮返回的 `replaceText`，并把最终文本返回给调用方。
- 在 `DFAMatcher.Match` 中恢复 `replaceRune` 调用，并在 `onlyOne=true` 且首次命中时立即返回。
- 在 `DFANode.AddWord` 中忽略空字符串关键词。
- 将 `A|B|C` 的匹配逻辑改为从上一次命中结束位置继续查找。
- 让不支持脱敏的 matcher 返回原文，避免清空服务层结果。

## 风险与兼容性

- `MatchService.Match` 的 matcher 执行顺序从随机变为固定顺序，`onlyOne=true` 时结果会更稳定，但可能与依赖旧随机顺序的调用方表现不同。
- DFA 对重叠词会继续按现有扫描逻辑返回多个结果，避免扩大行为变更范围。
- 正则规则仍保持命中后返回规则 tag，不返回具体匹配内容。
