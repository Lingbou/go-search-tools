
# Go Search Tools

## 简介
Go Search Tools 是一款基于Go语言开发的强大命令行文件搜索工具。它支持按文件名和文件内容进行搜索，具备灵活的过滤选项和高效的搜索性能，能够帮助开发者、系统管理员等快速定位所需文件。

## 功能特性
### 搜索功能
- **文件名搜索**：支持使用通配符 `*` 和 `?` 进行匹配，可通过 `--ignore-case` 参数忽略大小写，实现对文件名的精准或模糊查找。
- **文件内容搜索**：基于字符串匹配查找文件内容，支持多线程并行搜索（可通过 `--workers` 参数调整并发数），并提供超时控制（`--timeout` 参数防止长时间搜索）。

### 过滤选项
- **目录与深度控制**：可递归搜索子目录（默认开启），并通过 `--max-depth` 参数限制最大递归深度；也能通过 `--exclude-dir` 参数排除特定目录。
- **文件扩展名过滤**：支持通过 `--include-ext` 参数指定只搜索特定扩展名的文件，或通过 `--exclude-ext` 参数排除特定扩展名的文件。

### 输出与体验优化
- **彩色输出**：在ANSI终端中提供彩色输出，清晰区分不同类型的信息，增强可读性。
- **进度条显示**：可通过 `--progress` 参数选择是否显示进度条，实时掌握搜索进度。
- **人性化格式**：对文件大小和时间进行人性化格式处理，便于查看。

## 安装方法
### 方式一：使用Go工具链安装（推荐）
确保已安装Go环境（版本需在1.20及以上），在终端执行以下命令：
```bash
go install github.com/Lingbou/go-search-tools@latest
```
### 方式二：从源码构建
1. 克隆项目仓库：
```bash
git clone https://github.com/Lingbou/go-search-tools.git
```
2. 进入项目目录：
```bash
cd go-search-tools
```
3. 构建项目：
```bash
go build -o gost
```
或者
```bash
make build
```
构建成功后，会在当前目录生成名为 `gost` 的可执行文件（在Windows下为 `gost.exe` ）。

## 使用示例
### 按文件名搜索
- **搜索当前目录及其子目录下所有的go文件**：
```bash
gost name --include-ext .go "*.go"
```
- **在指定目录递归搜索，但只检查深度不超过3的目录，查找文件名包含config的文件**：
```bash
gost name --path /指定路径 --max-depth 3 "config.*"
```

### 按文件内容搜索
- **在当前目录及其子目录中，使用8个线程并行搜索包含error字符串的文件，设置搜索超时时间为30秒**：
```bash
gost content --workers 8 --timeout 30s "error"
```
- **在指定目录搜索文件内容，排除node_modules目录和.jpg文件，查找包含critical bug字符串的文件**：
```bash
gost content --path /指定路径 --exclude-dir node_modules --exclude-ext .jpg "critical bug"
```

## 贡献指南
非常欢迎大家为Go Search Tools贡献代码！如果你想参与项目开发，可以参考以下步骤：
1.  Fork本仓库到你自己的GitHub账号。
2.  克隆你Fork后的仓库到本地：
```bash
git clone https://github.com/Lingbou/go-search-tools.git
```
3.  创建一个新的分支：
```bash
git checkout -b [分支名称]
```
4.  在新分支上进行代码修改和功能开发。请确保你的代码遵循项目现有的代码风格和规范。
5.  完成修改后，提交你的代码并推送到远程仓库：
```bash
git add .
git commit -m "阿巴阿巴"
git push origin [分支名称]
```
6.  回到GitHub，在你的仓库页面发起Pull Request，详细描述你的改动内容和目的。我们会尽快对你的PR进行审核和反馈。

在贡献代码之前，建议先阅读项目的 [行为准则](此处若有相关准则文档链接可补充) ，确保你的贡献符合项目的整体理念和社区规范。

## 许可证
本项目采用 [MIT许可证](LICENSE) ，详情请查看项目根目录下的 `LICENSE` 文件。这意味着你可以自由地使用、修改和分发本项目代码，只需在衍生作品中保留原许可证和版权声明。 