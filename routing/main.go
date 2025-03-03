package routing

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "plugin"
    "strings"
    "time"
)

const (
    DefaultAPIFolder = "api"
    PluginExtension = ".so"
    CompiledFolder  = ".GocaGola"
)

var pluginCache = make(map[string]time.Time)

var supportedMethods = []string{"GET", "PUT", "POST", "PATCH", "DELETE"}

func Initialize(apiFolder string) {
    basePath, err := getBasePath(apiFolder)
    if err != nil {
        log.Fatal("Initialization failed:", err)
        return
    }

    router := setupRouter()
    if err = loadAPIHandlers(router, basePath); err != nil {
        log.Fatal("Error loading handlers:", err)
        return
    }

    router.Run(":8080")
}

func getBasePath(apiFolder string) (string, error) {
    basePath, err := os.Getwd()
    if err != nil {
        return "", fmt.Errorf("failed to get working directory: %v", err)
    }

    if err := os.MkdirAll(filepath.Join(basePath, CompiledFolder), 0755); err != nil {
        return "", fmt.Errorf("failed to create compiled API folder: %v", err)
    }

    if apiFolder == "" {
        apiFolder = DefaultAPIFolder
    }

    if !strings.HasPrefix(apiFolder, "/") {
        apiFolder = filepath.Join(basePath, apiFolder)
    }

    if _, err := os.Stat(apiFolder); os.IsNotExist(err) {
        return "", fmt.Errorf("API folder not found at %s", apiFolder)
    }

    return apiFolder, nil
}

func setupRouter() *gin.Engine {
    r := gin.Default()
    r.SetTrustedProxies([]string{"127.0.0.1"})
    return r
}

func shouldCompilePlugin(sourcePath, compiledPath string) bool {
    sourceInfo, err := os.Stat(sourcePath)
    if err != nil {
        return true
    }

    compiledInfo, err := os.Stat(compiledPath)
    if err != nil {
        return true
    }

    lastCompiled, exists := pluginCache[sourcePath]
    return !exists || sourceInfo.ModTime().After(lastCompiled) || sourceInfo.ModTime().After(compiledInfo.ModTime())
}

func compilePlugin(filePath string) (map[string]interface{}, error) {
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return nil, fmt.Errorf("file %s does not exist", filePath)
    }

    outputPath, err := prepareOutputPath(filePath)
    if err != nil {
        return nil, err
    }

    if !shouldCompilePlugin(filePath, outputPath) {
        return loadPackage(outputPath)
    }

    if err := buildPlugin(filePath, outputPath); err != nil {
        return nil, err
    }

    pluginCache[filePath] = time.Now()
    return loadPackage(outputPath)
}

func prepareOutputPath(filePath string) (string, error) {
    pwd, err := os.Getwd()
    if err != nil {
        return "", fmt.Errorf("failed to get working directory: %v", err)
    }

    relPath, err := filepath.Rel(pwd, filePath)
    if err != nil {
        return "", fmt.Errorf("failed to get relative path: %v", err)
    }

    outputPath := filepath.Join(CompiledFolder, relPath+PluginExtension)
    outputDir := filepath.Dir(outputPath)

    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create output directory: %v", err)
    }

    return outputPath, nil
}

func buildPlugin(sourcePath, outputPath string) error {
    log.Printf("Compiling plugin: %s -> %s\n", sourcePath, outputPath)
    
    cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outputPath, sourcePath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("plugin compilation failed: %v", err)
    }

    log.Printf("Successfully compiled plugin: %s\n", outputPath)
    return nil
}

func loadPackage(path string) (map[string]interface{}, error) {
    plg, err := plugin.Open(path)
    if err != nil {
        return nil, fmt.Errorf("unable to open plugin: %v", err)
    }

    routeHandler := make(map[string]interface{})
    for _, method := range supportedMethods {
        if handler, err := plg.Lookup(method); err == nil {
            routeHandler[method] = handler
        }
    }

    return routeHandler, nil
}

func loadAPIHandlers(r *gin.Engine, basePath string) error {
    basePathClean := filepath.Clean(basePath)
    
    return filepath.Walk(basePathClean, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return fmt.Errorf("walk error at %s: %v", path, err)
        }

        if shouldSkipFile(info) {
            return skipIfNeeded(info)
        }

        return handleGoFile(r, basePathClean, path)
    })
}

func shouldSkipFile(info os.FileInfo) bool {
    return info.IsDir() || !strings.HasSuffix(info.Name(), ".go")
}

func skipIfNeeded(info os.FileInfo) error {
    if info.IsDir() && info.Name() == CompiledFolder {
        return filepath.SkipDir
    }
    return nil
}

func handleGoFile(r *gin.Engine, basePath, filePath string) error {
    routeHandler, err := compilePlugin(filePath)
    if err != nil {
        return fmt.Errorf("failed to compile plugin %s: %v", filePath, err)
    }

    if routeHandler == nil {
        return nil
    }

    routePath := buildRoutePath(basePath, filePath)
    return registerHandlers(r, routePath, routeHandler)
}

func buildRoutePath(basePath, filePath string) string {
    relPath, _ := filepath.Rel(basePath, filePath)
    relPath = strings.TrimSuffix(relPath, ".go")
    
    if strings.HasSuffix(relPath, "/main") {
        relPath = strings.TrimSuffix(relPath, "/main")
    }
    
    return filepath.Join("/api", relPath)
}

func registerHandlers(r *gin.Engine, routePath string, handlers map[string]interface{}) error {
    for protocol, handlerFunc := range handlers {
        if handler, ok := handlerFunc.(func(*gin.Context)); ok {
            log.Printf("Registering route: %s %s\n", protocol, routePath)
            r.Handle(protocol, routePath, handler)
        }
    }
    return nil
}