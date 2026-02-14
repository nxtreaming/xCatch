# xCatch - X.com Content Scraper

基于 [uTools API](https://utools.readme.io/reference/getting-started-with-your-api-2) 的 X.com (Twitter) 内容抓取工具，使用 Go 语言实现。

## 功能模块

### 用户信息
- 根据用户名 (screen_name) 获取用户资料
- 根据用户 ID 获取用户资料
- 批量获取用户信息
- 用户名变更历史查询

### 推文内容
- 获取用户推文列表（支持分页）
- 获取推文详情及回复线程
- 批量获取推文
- 获取用户点赞列表
- 获取用户回复推文
- 获取用户精选推文
- 获取用户文章类推文

### 搜索
- 高级搜索（关键词、类型筛选）
- 搜索联想
- 热门趋势
- 新闻 / 体育 / 娱乐分类

### 社交关系
- 粉丝列表 / 关注列表
- 粉丝 ID / 关注 ID
- 两用户间关系查询
- 共同关注 / 蓝V粉丝

### 互动数据
- 转推者列表
- 点赞者列表
- 引用推文列表

### 社区 & 列表
- 用户社区查询
- 社区推文时间线
- Twitter 列表成员 / 时间线

## 快速开始

### 前置条件

1. 安装 Go 1.21+
2. 获取 uTools API Key:
   - 方式一: 在 https://twitter.good6.top 登录 Twitter，通过 [用户中心](https://twitter.utools.me/userCenter) 查看 apiKey
   - 方式二: 在商城直接购买 https://www.idatariver.com/zh-cn/project/twitterplankey-32e6

### 安装

```bash
cd xCatch
go mod tidy
go build -o xcatch ./cmd/
```

### 配置

#### 方式一：配置文件（推荐）

复制 `config.ini.example` 为 `config.ini`，填入你的 API Key：

```bash
cp config.ini.example config.ini
```

```ini
[xcatch]
api_key = your_api_key_here
# auth_token = your_auth_token_here
# base_url = https://fapi.uk
# timeout_sec = 30
# max_retries = 3
# rate_limit = 5
```

> ⚠️ `config.ini` 已被 `.gitignore` 排除，不会提交到 Git。

#### 方式二：环境变量

环境变量会覆盖配置文件中的同名配置：

| 变量 | 必填 | 说明 | 默认值 |
|------|------|------|--------|
| `XCATCH_API_KEY` | ✅ | uTools API Key | - |
| `XCATCH_AUTH_TOKEN` | ❌ | Twitter auth_token（部分接口需要） | - |
| `XCATCH_BASE_URL` | ❌ | API 基础 URL | `https://fapi.uk` |
| `XCATCH_TIMEOUT_SEC` | ❌ | HTTP 超时（秒） | `30` |
| `XCATCH_MAX_RETRIES` | ❌ | 最大重试次数 | `3` |
| `XCATCH_RATE_LIMIT` | ❌ | QPS 限制 | `5` |

配置优先级：环境变量 > config.ini > 默认值

### 使用示例

```bash
# 设置 API Key
export XCATCH_API_KEY="your_api_key_here"

# 查询用户信息
./xcatch user elonmusk

# 获取用户推文（默认1页）
./xcatch tweets 44196397

# 获取用户推文（最多3页）
./xcatch tweets 44196397 3

# 查看推文详情及回复
./xcatch tweet 1234567890

# 搜索推文
./xcatch search "bitcoin" Latest

# 获取粉丝列表
./xcatch followers 44196397

# 获取关注列表
./xcatch followings 44196397

# 获取用户点赞
./xcatch likes 44196397

# 获取热门趋势
./xcatch trending
```

### 作为 SDK 使用

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/xCatch/xcatch/config"
    "github.com/xCatch/xcatch/pkg/utools"
)

func main() {
    cfg := &config.Config{
        BaseURL:    "https://fapi.uk",
        APIKey:     "your_api_key",
        MaxRetries: 3,
        Timeout:    30 * time.Second,
        RateLimit:  5.0,
    }

    client, err := utools.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // 获取用户信息
    userData, err := client.GetUserByScreenNameV2(ctx, "elonmusk")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(userData))

    // 分页获取推文
    iter := client.NewPageIterator("/api/base/apitools/userTweetsV2", map[string]string{
        "userId": "44196397",
    }, 5) // 最多5页

    for iter.HasMore() {
        page, err := iter.Next(ctx)
        if err != nil {
            log.Fatal(err)
        }
        if page == nil {
            break
        }
        fmt.Printf("Page %d: %s\n", iter.PageCount(), string(page.RawData))
    }
}
```

## 项目结构

```
xCatch/
├── cmd/
│   └── main.go                  # CLI 入口
├── config/
│   ├── config.go                # 配置管理（INI 文件 + 环境变量）
│   └── errors.go                # 配置错误定义
├── pkg/
│   └── utools/
│       ├── client.go            # HTTP 客户端（认证、重试、限流、信封解包）
│       ├── cursor.go            # 分页 cursor 迭代器
│       ├── errors.go            # API 错误类型
│       ├── types.go             # 数据结构定义
│       ├── user.go              # 用户信息 API
│       ├── tweet.go             # 推文内容 API
│       ├── search.go            # 搜索 API
│       └── social.go            # 社交关系 / 列表 / 社区 API
├── config.ini.example           # 配置文件模板
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## API 节点选择

| 节点 | 地址 | 说明 |
|------|------|------|
| 默认（香港） | `https://fapi.uk` | 默认节点 |
| 美国节点 | `https://l2.fapi.uk` | 服务器在美国时推荐 |

通过 `config.ini` 中的 `base_url` 或 `XCATCH_BASE_URL` 环境变量切换节点。

## 错误处理

| 错误 | 含义 | SDK 行为 |
|------|------|----------|
| Rate limit exceeded (code 88) | 频率超限 | 指数退避重试（1s, 2s, 4s...） |
| Forbidden (403) | 机器人账号被锁 | 指数退避重试（1s, 2s, 4s...） |
| Unauthorized (401) | auth_token 缺失/无效 | 直接返回错误 |

最大重试次数通过 `XCATCH_MAX_RETRIES` 配置。

## 技术支持

- Discord: discord.gg/VEKBj3fT9G
- Telegram: https://t.me/Yang0619
