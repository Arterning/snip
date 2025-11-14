# Snip - CLI Snippet Manager

一个用 Go 编写的命令行片段管理工具，帮助你保存、搜索和执行常用命令。

## 功能特性

- ✨ 添加带占位符的命令片段
- 🔍 模糊搜索已保存的命令
- 🚀 快速执行命令
- 💾 使用 YAML 文件本地存储
- 🎯 交互式占位符输入

## 安装

```bash
# 构建可执行文件
go build -o snip.exe

# （可选）将 snip.exe 添加到系统 PATH 中
```

## 使用方法

### 添加新的 snippet

```bash
snip new
```

系统会提示你输入：
- **命令**：你想保存的命令（必需）
- **描述**：命令的描述（可选）

示例：
```
Command: git commit -m "<message>"
Description: Git commit with custom message
```

### 使用占位符

在命令中使用 `<placeholder_name>` 格式添加占位符：

```bash
# 示例命令
docker run -p <port>:80 <image_name>
curl -X POST <url> -d '<data>'
ssh <user>@<host>
```

当执行这些命令时，系统会提示你输入每个占位符的值。

### 搜索和执行 snippet

```bash
snip
```

这会进入交互式搜索模式：
- 输入关键字进行模糊搜索
- 使用 ↑/↓ 箭头键选择命令
- 按 Enter 执行选中的命令
- 如果命令包含占位符，系统会提示你输入值

## 数据存储

所有 snippet 保存在 `~/.snip.yaml` 文件中。

示例文件内容：
```yaml
snippets:
  - id: 550e8400-e29b-41d4-a716-446655440000
    command: git commit -m "<message>"
    description: Git commit with custom message
    created_at: 2025-01-15T10:30:00Z
    updated_at: 2025-01-15T10:30:00Z
  - id: 660e8400-e29b-41d4-a716-446655440001
    command: docker ps -a
    description: List all containers
    created_at: 2025-01-15T10:31:00Z
    updated_at: 2025-01-15T10:31:00Z
```

## 示例使用场景

### 1. Docker 命令
```bash
# 保存命令
Command: docker run -d --name <container_name> -p <port>:80 nginx
Description: Run nginx container

# 执行时输入
container_name: my-nginx
port: 8080
```

### 2. Git 操作
```bash
# 保存命令
Command: git checkout -b <branch_name>
Description: Create and checkout new branch

# 执行时输入
branch_name: feature/new-feature
```

### 3. SSH 连接
```bash
# 保存命令
Command: ssh <user>@<host> -p <port>
Description: SSH connection

# 执行时输入
user: admin
host: example.com
port: 22
```

## 项目结构

```
snip/
├── main.go              # 程序入口
├── cmd/
│   ├── root.go         # 根命令（搜索模式）
│   └── new.go          # new 命令
├── pkg/
│   ├── snippet/
│   │   └── snippet.go  # Snippet 数据结构
│   ├── storage/
│   │   └── storage.go  # YAML 存储
│   └── executor/
│       └── executor.go # 命令执行
└── go.mod
```

## 依赖

- [cobra](https://github.com/spf13/cobra) - CLI 框架
- [promptui](https://github.com/manifoldco/promptui) - 交互式提示
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML 解析
- [uuid](https://github.com/google/uuid) - UUID 生成

## 技术细节

- **占位符格式**：`<name>` 尖括号
- **搜索方式**：模糊搜索（不区分大小写）
- **存储格式**：YAML 文件
- **命令执行**：根据操作系统使用 `cmd /C`（Windows）或 `sh -c`（Unix）

## 注意事项

- 占位符名称不能包含 `>` 字符
- 相同的占位符名称在一个命令中会被替换为同一个值
- 执行命令时会继承当前终端的环境变量
- 命令会在当前工作目录执行

## License

MIT
