# 项目配置
APP_NAME = test
MAIN_PATH = ./cmd
BUILD_DIR = build

# Go 工具链
GOCMD     = go
GOBUILD   = $(GOCMD) build
GOCLEAN   = $(GOCMD) clean
GOTEST    = $(GOCMD) test
GOMOD     = $(GOCMD) mod
GOFMT     = $(GOCMD) fmt

# 平台检测：Windows 下自动加 .exe 后缀
ifeq ($(OS),Windows_NT)
    BIN_EXT = .exe
    RM      = cmd /c rmdir /s /q
    MKDIR   = cmd /c if not exist
else
    BIN_EXT =
    RM      = rm -rf
    MKDIR   = mkdir -p
endif

BINARY = $(BUILD_DIR)/$(APP_NAME)$(BIN_EXT)

.PHONY: all build clean test run tidy fmt vet

# 默认目标：整理依赖 → 构建
all: tidy build

# 构建可执行文件
build:
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BINARY) $(MAIN_PATH)

# 清理构建产物
clean:
	$(GOCLEAN)
	$(RM) $(BUILD_DIR)

# 运行测试
test:
	$(GOTEST) ./...

# 构建并运行
run: build
	./$(BINARY)

# 整理 go.mod / go.sum
tidy:
	$(GOMOD) tidy

# 格式化代码
fmt:
	$(GOFMT) ./...

# 静态检查
vet:
	$(GOCMD) vet ./...
