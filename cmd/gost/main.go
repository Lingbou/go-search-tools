package main

import (
    "context"
    // "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"

    "github.com/spf13/cobra"
    "github.com/fatih/color"
    "github.com/schollz/progressbar/v3"
)

var (
    rootCmd = &cobra.Command{
        Use:   "fsearch [command]",
        Short: "文件搜索工具，支持按文件名和内容搜索",
        Long:  "fsearch 是一个强大的文件搜索工具，可以快速在目录树中查找文件或内容",
    }

    // 搜索文件名的命令
    searchNameCmd = &cobra.Command{
        Use:   "name [flags] <pattern>",
        Short: "按文件名搜索",
        Long:  "按文件名搜索，支持通配符(*和?)",
        Args:  cobra.ExactArgs(1),
        Run:   runSearchName,
    }

    // 搜索文件内容的命令
    searchContentCmd = &cobra.Command{
        Use:   "content [flags] <pattern>",
        Short: "按文件内容搜索",
        Long:  "按文件内容搜索，使用正则表达式匹配",
        Args:  cobra.ExactArgs(1),
        Run:   runSearchContent,
    }

    // 命令行参数
    searchPath   string
    ignoreCase   bool
    recursive    bool
    maxDepth     int
    excludeDirs  []string
    includeExts  []string
    excludeExts  []string
    numWorkers   int
    timeout      time.Duration
    colorOutput  bool
    showProgress bool
)

func init() {
    // 全局参数
    rootCmd.PersistentFlags().StringVarP(&searchPath, "path", "p", ".", "搜索路径")
    rootCmd.PersistentFlags().BoolVarP(&ignoreCase, "ignore-case", "i", false, "忽略大小写")
    rootCmd.PersistentFlags().BoolVarP(&colorOutput, "color", "c", true, "启用颜色输出")
    rootCmd.PersistentFlags().BoolVarP(&showProgress, "progress", "P", false, "显示进度")
    
    // 文件名搜索参数
    searchNameCmd.Flags().BoolVarP(&recursive, "recursive", "r", true, "递归搜索子目录")
    searchNameCmd.Flags().IntVarP(&maxDepth, "max-depth", "d", -1, "最大递归深度，-1表示不限制")
    searchNameCmd.Flags().StringSliceVarP(&excludeDirs, "exclude-dir", "e", []string{}, "排除的目录")
    searchNameCmd.Flags().StringSliceVarP(&includeExts, "include-ext", "I", []string{}, "只包含的文件扩展名")
    searchNameCmd.Flags().StringSliceVarP(&excludeExts, "exclude-ext", "E", []string{}, "排除的文件扩展名")
    
    // 内容搜索参数
    searchContentCmd.Flags().BoolVarP(&recursive, "recursive", "r", true, "递归搜索子目录")
    searchContentCmd.Flags().IntVarP(&maxDepth, "max-depth", "d", -1, "最大递归深度，-1表示不限制")
    searchContentCmd.Flags().StringSliceVarP(&excludeDirs, "exclude-dir", "e", []string{}, "排除的目录")
    searchContentCmd.Flags().StringSliceVarP(&includeExts, "include-ext", "I", []string{}, "只包含的文件扩展名")
    searchContentCmd.Flags().StringSliceVarP(&excludeExts, "exclude-ext", "E", []string{}, "排除的文件扩展名")
    searchContentCmd.Flags().IntVarP(&numWorkers, "workers", "w", 4, "并行工作线程数")
    searchContentCmd.Flags().DurationVarP(&timeout, "timeout", "t", 0, "搜索超时时间，例如10s, 2m等")
    
    // 将子命令添加到根命令
    rootCmd.AddCommand(searchNameCmd, searchContentCmd)
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

// 按文件名搜索的执行函数
func runSearchName(cmd *cobra.Command, args []string) {
    pattern := args[0]
    
    // 检查路径是否存在
    if _, err := os.Stat(searchPath); os.IsNotExist(err) {
        color.Red("错误: 搜索路径不存在: %s", searchPath)
        os.Exit(1)
    }
    
    // 编译通配符模式
    if !ignoreCase {
        pattern = strings.ToLower(pattern)
    }
    
    // 用于存储匹配的文件
    var matches []string
    var mu sync.Mutex
    
    // 计算文件总数用于进度条
    var totalFiles int
    if showProgress {
        totalFiles = countFiles(searchPath)
        if totalFiles == 0 {
            color.Yellow("没有找到文件")
            return
        }
    }
    
    // 创建进度条
    var bar *progressbar.ProgressBar
    if showProgress {
        bar = progressbar.Default(int64(totalFiles), "搜索中")
    }
    
    // 递归搜索文件
    err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        // 更新进度条
        if showProgress {
            bar.Add(1)
        }
        
        // 排除目录检查
        if info.IsDir() {
            // 检查是否是要排除的目录
            for _, excludeDir := range excludeDirs {
                if info.Name() == excludeDir {
                    return filepath.SkipDir
                }
            }
            
            // 检查深度
            if maxDepth >= 0 {
                depth := strings.Count(strings.TrimPrefix(path, searchPath), string(os.PathSeparator))
                if depth > maxDepth {
                    return filepath.SkipDir
                }
            }
            
            return nil
        }
        
        // 扩展名过滤
        ext := filepath.Ext(info.Name())
        if len(includeExts) > 0 {
            found := false
            for _, includeExt := range includeExts {
                if ext == includeExt {
                    found = true
                    break
                }
            }
            if !found {
                return nil
            }
        }
        
        if len(excludeExts) > 0 {
            for _, excludeExt := range excludeExts {
                if ext == excludeExt {
                    return nil
                }
            }
        }
        
        // 文件名匹配
        nameToMatch := info.Name()
        if ignoreCase {
            nameToMatch = strings.ToLower(nameToMatch)
        }
        
        // 简单的通配符匹配实现
        if matchPattern(nameToMatch, pattern) {
            mu.Lock()
            matches = append(matches, path)
            mu.Unlock()
            
            // 打印匹配结果
            printMatch(path, info)
        }
        
        return nil
    })
    
    if err != nil {
        color.Red("搜索过程中出错: %v", err)
        os.Exit(1)
    }
    
    // 打印结果摘要
    if len(matches) == 0 {
        color.Yellow("没有找到匹配的文件")
    } else {
        color.Green("共找到 %d 个匹配的文件", len(matches))
    }
}

// 按内容搜索的执行函数
func runSearchContent(cmd *cobra.Command, args []string) {
    pattern := args[0]
    
    // 检查路径是否存在
    if _, err := os.Stat(searchPath); os.IsNotExist(err) {
        color.Red("错误: 搜索路径不存在: %s", searchPath)
        os.Exit(1)
    }
    
    // 创建上下文用于超时控制
    var ctx context.Context
    var cancel context.CancelFunc
    if timeout > 0 {
        ctx, cancel = context.WithTimeout(context.Background(), timeout)
        defer cancel()
    } else {
        ctx = context.Background()
    }
    
    // 用于存储匹配的文件
    var matches []string
    var mu sync.Mutex
    
    // 计算文件总数用于进度条
    var totalFiles int
    if showProgress {
        totalFiles = countFiles(searchPath)
        if totalFiles == 0 {
            color.Yellow("没有找到文件")
            return
        }
    }
    
    // 创建进度条
    var bar *progressbar.ProgressBar
    if showProgress {
        bar = progressbar.Default(int64(totalFiles), "搜索中")
    }
    
    // 创建文件通道
    filesCh := make(chan string)
    resultsCh := make(chan string)
    
    // 启动工作协程
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for filePath := range filesCh {
                // 检查是否超时
                select {
                case <-ctx.Done():
                    return
                default:
                    // 搜索文件内容
                    if searchFileContent(ctx, filePath, pattern) {
                        resultsCh <- filePath
                    }
                    
                    // 更新进度条
                    if showProgress {
                        bar.Add(1)
                    }
                }
            }
        }()
    }
    
    // 收集结果的协程
    go func() {
        wg.Wait()
        close(resultsCh)
    }()
    
    // 遍历文件并发送到通道
    go func() {
        defer close(filesCh)
        
        err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            
            // 检查是否超时
            select {
            case <-ctx.Done():
                return ctx.Err()
            default:
                // 排除目录检查
                if info.IsDir() {
                    // 检查是否是要排除的目录
                    for _, excludeDir := range excludeDirs {
                        if info.Name() == excludeDir {
                            return filepath.SkipDir
                        }
                    }
                    
                    // 检查深度
                    if maxDepth >= 0 {
                        depth := strings.Count(strings.TrimPrefix(path, searchPath), string(os.PathSeparator))
                        if depth > maxDepth {
                            return filepath.SkipDir
                        }
                    }
                    
                    return nil
                }
                
                // 扩展名过滤
                ext := filepath.Ext(info.Name())
                if len(includeExts) > 0 {
                    found := false
                    for _, includeExt := range includeExts {
                        if ext == includeExt {
                            found = true
                            break
                        }
                    }
                    if !found {
                        return nil
                    }
                }
                
                if len(excludeExts) > 0 {
                    for _, excludeExt := range excludeExts {
                        if ext == excludeExt {
                            return nil
                        }
                    }
                }
                
                // 发送文件路径到通道
                filesCh <- path
                
                return nil
            }
        })
        
        if err != nil {
            color.Red("搜索过程中出错: %v", err)
        }
    }()
    
    // 处理结果
    for filePath := range resultsCh {
        // 获取文件信息
        info, err := os.Stat(filePath)
        if err != nil {
            color.Red("获取文件信息失败: %s - %v", filePath, err)
            continue
        }
        
        mu.Lock()
        matches = append(matches, filePath)
        mu.Unlock()
        
        // 打印匹配结果
        printMatch(filePath, info)
    }
    
    // 检查是否超时
    if ctx.Err() != nil {
        color.Yellow("搜索超时，已找到 %d 个匹配的文件", len(matches))
    } else if len(matches) == 0 {
        color.Yellow("没有找到匹配的文件")
    } else {
        color.Green("共找到 %d 个匹配的文件", len(matches))
    }
}

// 简单的通配符匹配函数
func matchPattern(s, pattern string) bool {
    if pattern == "" {
        return s == ""
    }
    
    if pattern[0] == '*' {
        if len(pattern) == 1 {
            return true
        }
        
        for i := 0; i <= len(s); i++ {
            if matchPattern(s[i:], pattern[1:]) {
                return true
            }
        }
        return false
    }
    
    return len(s) > 0 && (pattern[0] == '?' || pattern[0] == s[0]) && matchPattern(s[1:], pattern[1:])
}

// 搜索文件内容
func searchFileContent(ctx context.Context, filePath, pattern string) bool {
    // 简单实现：这里可以使用更高级的文本搜索库如bleve
    // 为简化示例，我们只检查文件是否包含字符串
    file, err := os.Open(filePath)
    if err != nil {
        return false
    }
    defer file.Close()
    
    // 读取文件内容
    info, _ := file.Stat()
    size := info.Size()
    
    // 限制读取大小，避免处理大文件
    if size > 10*1024*1024 { // 10MB
        return false
    }
    
    content := make([]byte, size)
    _, err = file.Read(content)
    if err != nil {
        return false
    }
    
    // 检查内容是否包含模式
    if ignoreCase {
        return strings.Contains(strings.ToLower(string(content)), strings.ToLower(pattern))
    }
    
    return strings.Contains(string(content), pattern)
}

// 计算目录中的文件总数
func countFiles(path string) int {
    count := 0
    filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !info.IsDir() {
            // 扩展名过滤
            ext := filepath.Ext(info.Name())
            if len(includeExts) > 0 {
                found := false
                for _, includeExt := range includeExts {
                    if ext == includeExt {
                        found = true
                        break
                    }
                }
                if !found {
                    return nil
                }
            }
            
            if len(excludeExts) > 0 {
                for _, excludeExt := range excludeExts {
                    if ext == excludeExt {
                        return nil
                    }
                }
            }
            
            count++
        }
        return nil
    })
    return count
}

// 打印匹配结果
func printMatch(path string, info os.FileInfo) {
    if colorOutput {
        // 使用颜色区分不同部分
        fileType := ""
        if info.IsDir() {
            fileType = color.CyanString("[目录]")
        } else {
            fileType = color.GreenString("[文件]")
        }
        
        size := ""
        if !info.IsDir() {
            size = formatSize(info.Size())
        }
        
        modified := info.ModTime().Format("2006-01-02 15:04:05")
        
        fmt.Printf("%s %s %s %s\n", 
            fileType, 
            color.MagentaString(info.Name()), 
            color.BlueString(size), 
            color.YellowString(modified))
    } else {
        fmt.Printf("%s %s %d %s\n", 
            info.Mode(), 
            info.Name(), 
            info.Size(), 
            info.ModTime().Format("2006-01-02 15:04:05"))
    }
}

// 格式化文件大小
func formatSize(bytes int64) string {
    units := []string{"B", "KB", "MB", "GB", "TB"}
    var i int
    size := float64(bytes)
    
    for size >= 1024 && i < len(units)-1 {
        size /= 1024
        i++
    }
    
    return fmt.Sprintf("%.2f %s", size, units[i])
}    