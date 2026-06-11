# matchs

`matchs` 是一个用于文本关键词快速查找与替换的 Go 包。它可以对输入文本批量匹配普通关键词、组合规则和正则规则，并返回命中的规则列表以及替换后的文本。

## 功能特性

- 支持普通关键词匹配，适合敏感词、黑名单词等固定文本场景。
- `DFAMatcher` 内部使用 `Aho-Corasick` 自动机匹配普通关键词，并按 rune 级替换命中文本，适配中文等多字节字符。
- 支持组合规则：
  - `A|B|C`：表示 `A`、`B`、`C` 必须按顺序出现在文本中。
  - `A#B`：表示文本必须包含 `A`，且不能包含 `B`；也支持 `A#B#C` 表示包含 `A` 且不能包含 `B`、`C`。
- 支持正则规则，使用 `reg@` 前缀声明。
- 支持 `onlyOne` 参数控制是否命中一次后立即返回。

## 安装

```bash
go get git.faceue.com/matchs
```

如果是在当前仓库内开发，可以直接运行：

```bash
go test ./...
```

## 快速开始

```go
package main

import (
	"fmt"

	"git.faceue.com/matchs"
)

func main() {
	service := matchs.NewMatchService()
	service.Build([]string{
		"敏感词",
		"A|B",
		`reg@1\d{10}`,
	})

	words, replaced := service.Match("敏感词 AxxB 手机号13800138000", false, '*')

	fmt.Println(words)
	fmt.Println(replaced)
}
```

输出示例：

```text
[敏感词 A|B reg@1\d{10}]
*** AxxB 手机号13800138000
```

## 规则说明

### 普通关键词

不包含 `|`、`#`，且不以 `reg@` 开头的规则会作为普通关键词处理。

```go
service.Build([]string{"敏感词", "违规词"})
words, replaced := service.Match("这里有敏感词", false, '*')
```

普通关键词命中后会参与替换，`replaced` 会返回替换后的文本。

普通关键词由 `DFAMatcher` 处理。虽然保留了 `DFA` 这个 matcher 类型名，内部实现已使用 `Aho-Corasick` 自动机，以便更好地支持大量关键词匹配。

### 顺序组合规则

使用 `|` 分隔多个词，表示这些词必须按顺序出现在文本中。

```go
service.Build([]string{"A|B|C"})
```

- `AxxBxxC`：命中。
- `BxxAxxC`：不命中。

组合规则只返回命中的规则，不执行文本替换。

### 排除组合规则

使用 `#` 表示排除关系。第一个词是必须包含的主词，后续词都是排除词。

```go
service.Build([]string{"A#B"})
```

- `xxAxx`：命中。
- `xxAxxBxx`：不命中。

该规则表示文本必须包含 `A`，且不能包含 `B`。

`A#B#C` 表示文本必须包含 `A`，且不能包含 `B` 或 `C`。

注意：组合规则当前不建议混用 `|` 和 `#`。如果规则同时包含两者，会按 `|` 规则处理。

### 正则规则

使用 `reg@` 前缀声明正则规则。

```go
service.Build([]string{`reg@1\d{10}`})
```

正则规则命中后返回完整规则 tag，例如 `reg@1\d{10}`。当前正则规则只返回命中规则，不执行文本替换。

正则语法使用 Go 标准库 `regexp`，不支持 lookbehind、backreference 等回溯型正则特性。无法编译的正则规则会在 `Build` 阶段被忽略。

## API 说明

### `NewMatchService`

创建匹配服务。

```go
service := matchs.NewMatchService()
```

### `Build`

构建规则列表。

```go
service.Build([]string{"敏感词", "A|B", `reg@1\d{10}`})
```

规则会按以下方式分类：

- 以 `reg@` 开头：正则规则。
- 包含 `|` 或 `#`：组合规则。
- 其他：普通关键词。

### `Match`

执行匹配。

```go
words, replaced := service.Match(text, onlyOne, repl)
```

参数说明：

- `text`：待匹配文本。
- `onlyOne`：是否只返回第一次命中的结果；服务层会按 `DFA`、`ASSEMBLE`、`REGEXP` 的顺序判断。
- `repl`：普通关键词命中后的替换字符。

返回值说明：

- `words`：命中的规则或关键词列表。
- `replaced`：替换后的文本；组合规则和正则规则不执行替换，会保持原文或上一轮替换结果。

## 匹配顺序

`MatchService` 按固定顺序执行 matcher：

1. `DFA`：普通关键词，内部使用 `Aho-Corasick` 自动机。
2. `ASSEMBLE`：组合规则。
3. `REGEXP`：正则规则。

当 `onlyOne=true` 时，命中任意一类规则后会立即返回。普通关键词内部会返回起点最靠前的命中；同一起点下，短词会先于长词返回。

## 性能说明

- 普通关键词由 `DFAMatcher` 使用 `Aho-Corasick` 自动机匹配，适合较大规模的固定关键词集合。
- 组合规则会基于首个非空词建立 anchor 倒排索引，匹配时先筛选候选规则，再执行完整的 `A|B|C` 或 `A#B` 校验。
- 正则规则会使用 Go 标准库 `regexp` 的 `LiteralPrefix` 提取固定前缀；无法提取前缀的规则会进入 fallback 列表逐条匹配。
- 如果规则数量很大，建议通过 `go test -bench=. -benchmem ./...` 观察本地机器上的构建和匹配成本。

## 注意事项

- 空字符串关键词会被忽略。
- 只有普通关键词支持文本替换。
- 组合规则和正则规则目前只返回命中的规则标识。
- 正则能力由 Go 标准库 `regexp` 提供。
- `reg@` 前缀优先级最高；以 `reg@` 开头的规则会作为正则处理，即使后续表达式中包含 `|` 或 `#`。

