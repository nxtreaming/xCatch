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

1. 安装 Go 1.23+
2. 获取 uTools API Key:
   - 方式一: 在 https://twitter.good6.top 登录 Twitter，通过 [用户中心](https://twitter.utools.me/userCenter) 查看 apiKey
   - 方式二: 在商城直接购买 https://www.idatariver.com/zh-cn/project/twitterplankey-32e6

### 安装

```bash
cd xCatch
go mod tidy
go build -o xcatch.exe ./cmd/
# macOS/Linux 可使用: go build -o xcatch ./cmd/
```

### 配置

#### 方式一：配置文件（推荐）

复制 `config.ini.example` 为 `config.ini`，填入你的 API Key：

```bash
cp config.ini.example config.ini
# PowerShell: Copy-Item config.ini.example config.ini
```

```ini
[xcatch]
api_key = your_api_key_here
# auth_token = your_auth_token_here
# ct0 = your_ct0_cookie_here
# XCATCH_TEST_USER_ID = 44196397
# XCATCH_TEST_SCREEN_NAME = elonmusk
# base_url = https://fapi.uk
# timeout_sec = 30
# max_retries = 3
# rate_limit = 5
```

#### 方式二：环境变量

环境变量会覆盖配置文件中的同名配置：

| 变量 | 必填 | 说明 | 默认值 |
|------|------|------|--------|
| `XCATCH_API_KEY` | ✅ | uTools API Key | - |
| `XCATCH_AUTH_TOKEN` | ❌ | Twitter auth_token（部分接口需要） | - |
| `XCATCH_CT0` | ❌ | Twitter ct0（鉴权接口建议与 auth_token 一起设置） | - |
| `XCATCH_BASE_URL` | ❌ | API 基础 URL | `https://fapi.uk` |
| `XCATCH_TIMEOUT_SEC` | ❌ | HTTP 超时（秒） | `30` |
| `XCATCH_MAX_RETRIES` | ❌ | 最大重试次数 | `3` |
| `XCATCH_RATE_LIMIT` | ❌ | QPS 限制 | `5` |

配置优先级：环境变量 > config.ini > 默认值

如果 `config.ini` 存在但格式错误，程序会打印 warning 并回退到默认值 + 环境变量。

### auth_token 说明

以下接口需要提供 `auth_token`，否则会直接返回错误：

- `GetHomeTimeline`
- `GetMentionsTimeline`
- `GetAccountAnalytics`

可通过 `config.ini` 的 `auth_token` 字段或环境变量 `XCATCH_AUTH_TOKEN` 设置。

根据官方 `go-client-generated` 参考实现，`GetHomeTimeline` / `GetMentionsTimeline` 这类接口通常还会携带 `ct0`。本项目会在配置了 `ct0` 时自动透传（`config.ini` 的 `ct0` 字段或环境变量 `XCATCH_CT0`）。

## 集成测试（真实 API）

项目包含两类测试：

- 单元测试（mock server）：`go test ./...`
- 集成测试（真实请求）：`go test -tags integration ...`

### User 真实集成测试前置条件

1. 设置运行开关（必须）
   - PowerShell: `$env:XCATCH_RUN_INTEGRATION="1"`
2. 提供 API Key（必须）
   - `config.ini` 的 `api_key`，或环境变量 `XCATCH_API_KEY`
3. 提供测试用户（建议）
   - `XCATCH_TEST_USER_ID`（推荐在 `config.ini` 的 `[xcatch]` 下配置）
   - `XCATCH_TEST_SCREEN_NAME`（用于 screenName 相关接口）
4. 鉴权接口（可选）
   - `GetAccountAnalytics` 需要 `auth_token`（建议同时设置 `ct0`）

### 运行命令

```powershell
# 全量 user 真实集成测试
$env:XCATCH_RUN_INTEGRATION="1"
go test -tags integration ./pkg/utools -run TestUserIntegration_RealAPI -v

# 仅测账号分析（需要 auth_token）
$env:XCATCH_RUN_INTEGRATION="1"
go test -tags integration ./pkg/utools -run TestUserIntegration_RealAPI/GetAccountAnalytics -v
```

### 为什么会看到 SKIP

以下情况会被标记为 `SKIP`（这是设计行为，不是本地代码失败）：

- 缺少前置配置（如 `XCATCH_TEST_USER_ID`、`XCATCH_TEST_SCREEN_NAME`、`auth_token`）
- 上游接口返回 `5xx`（例如 500/502）

只有非 `5xx` 的真实调用错误才会导致测试失败。

### 使用示例

```bash
# 设置 API Key
export XCATCH_API_KEY="your_api_key_here"   # PowerShell: $env:XCATCH_API_KEY="your_api_key_here"

# 查询用户信息
./xcatch.exe user elonmusk

# 获取用户推文（默认1页）
./xcatch.exe tweets 44196397

# 获取用户推文（最多3页）
./xcatch.exe tweets 44196397 3

# 注意：max_pages 必须是正整数

# 查看推文详情及回复
./xcatch.exe tweet 1234567890

# 搜索推文
./xcatch.exe search "bitcoin" Latest

# 获取粉丝列表
./xcatch.exe followers 44196397

# 获取关注列表
./xcatch.exe followings 44196397

# 获取用户点赞
./xcatch.exe likes 44196397

# 获取热门趋势
./xcatch.exe trending
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

## 接口能力矩阵（快速索引）

### CLI 命令与 SDK 方法映射

| CLI 命令 | SDK 方法 | 说明 |
|---|---|---|
| `user <screen_name>` | `GetUserByScreenNameV2` | 用户资料查询 |
| `tweets <user_id> [max_pages]` | `GetUserTweets` / `NewPageIterator` | 用户推文分页 |
| `tweet <tweet_id>` | `GetTweetDetail` | 推文详情与回复线程 |
| `search <query> [type]` | `Search` | 高级搜索 |
| `followers <user_id>` | `GetFollowers` | 粉丝列表 |
| `followings <user_id>` | `GetFollowings` | 关注列表 |
| `likes <user_id>` | `GetUserLikes` | 点赞列表 |
| `trending` | `GetTrending` | 热门趋势 |

### 常用接口能力

| 能力 | 代表方法 | 需要 `auth_token` | 支持 `cursor` |
|---|---|---|---|
| 用户资料 | `GetUserByScreenNameV2` / `GetUserByIDV2` | 否 | 否 |
| 用户推文 | `GetUserTweets` | 否 | 是 |
| 推文详情 | `GetTweetDetail` | 否 | 是 |
| 搜索 | `Search` | 否 | 是 |
| 粉丝/关注 | `GetFollowers` / `GetFollowings` | 否 | 是 |
| 点赞列表 | `GetUserLikes` | 否 | 是 |
| Home 时间线 | `GetHomeTimeline` | 是 | 是 |
| Mentions 时间线 | `GetMentionsTimeline` | 是 | 是 |
| 账号分析 | `GetAccountAnalytics` | 是 | 否 |

## Endpoint 路径对照表（方法 -> Path）

> 说明：以下路径来自当前 SDK 实现，便于与你的 uTools 文档逐项核对。

### User

| SDK 方法 | Path |
|---|---|
| `GetUserByScreenName` | `/api/base/apitools/getUserByIdOrNameShow` |
| `GetUserByID` | `/api/base/apitools/usersByIdRestIds` |
| `GetUsersByIDs` | `/api/base/apitools/usersByIdRestIds` |
| `GetUsernameChanges` | `/api/base/apitools/usernameChanges` |
| `LookupUser` | `/api/base/apitools/getUserByIdOrNameLookup` |
| `GetUserByScreenNameV2` | `/api/base/apitools/userByScreenNameV2` |
| `GetUserByIDV2` | `/api/base/apitools/uerByIdRestIdV2` |
| `GetUsersByIDsV2` | `/api/base/apitools/usersByIdRestIds` |
| `GetAccountAnalytics` | `/api/base/apitools/accountAnalytics` |

### Tweet

| SDK 方法 | Path |
|---|---|
| `GetUserTweets` | `/api/base/apitools/userTweetsV2` |
| `GetUserTimeline` | `/api/base/apitools/userTimeline` |
| `GetTweetDetail` | `/api/base/apitools/tweetTimeline` |
| `GetTweetSimple` | `/api/base/apitools/tweetSimple` |
| `GetTweetsByIDs` | `/api/base/apitools/tweetResultsByRestIds` |
| `GetUserReplies` | `/api/base/apitools/userTweetReply` |
| `GetUserLikes` | `/api/base/apitools/userLikeV2` |
| `GetUserHighlights` | `/api/base/apitools/highlightsV2` |
| `GetUserArticlesTweets` | `/api/base/apitools/userArticlesTweets` |
| `GetHomeTimeline` | `/api/base/apitools/homeTimeline` |
| `GetMentionsTimeline` | `/api/base/apitools/mentionsTimeline` |
| `GetRetweeters` | `/api/base/apitools/retweetersV2` |
| `GetFavoriters` | `/api/base/apitools/favoritersV2` |
| `GetQuotes` | `/api/base/apitools/quotesV2` |

### Search

| SDK 方法 | Path |
|---|---|
| `Search` | `/api/base/apitools/search` |
| `SearchBox` | `/api/base/apitools/searchBox` |
| `GetTrends` | `/api/base/apitools/trends` |
| `GetTrending` | `/api/base/apitools/trending` |
| `GetNews` | `/api/base/apitools/news` |
| `GetExplorePage` | `/api/base/apitools/explore` |
| `GetSports` | `/api/base/apitools/sports` |
| `GetEntertainment` | `/api/base/apitools/entertainment` |

### Social / List / Communities

| SDK 方法 | Path |
|---|---|
| `GetFollowers` | `/api/base/apitools/followersListV2` |
| `GetFollowings` | `/api/base/apitools/followingsListV2` |
| `GetFollowerIDs` | `/api/base/apitools/followersIds` |
| `GetFollowingIDs` | `/api/base/apitools/followingsIds` |
| `GetRelationship` | `/api/base/apitools/getFriendshipsShow` |
| `GetFollowersYouKnow` | `/api/base/apitools/followersYouKnowV2` |
| `GetBlueVerifiedFollowers` | `/api/base/apitools/blueVerifiedFollowersV2` |
| `GetListByUser` | `/api/base/apitools/getListByUserIdOrScreenName` |
| `GetListMembers` | `/api/base/apitools/listMembersByListIdV2` |
| `GetListTimeline` | `/api/base/apitools/listLatestTweetsTimeline` |
| `GetCommunitiesByScreenName` | `/api/base/apitools/getCommunitiesByScreenName` |
| `GetCommunityInfo` | `/api/base/apitools/communitiesFetchOneQuery` |
| `GetCommunityTweets` | `/api/base/apitools/communitiesTweetsTimelineV2` |
| `GetCommunityMembers` | `/api/base/apitools/communitiesMemberV2` |

### Client Utilities

| SDK 方法 | Path |
|---|---|
| `TokenSync` | `/api/base/apitools/tokenSync` |

## 接口版本与兼容建议（Legacy vs V2）

为减少上游接口变更带来的影响，建议优先使用 V2 命名接口。

### 优先级建议

1. 优先使用带 `V2` 的方法（如 `GetUserByScreenNameV2`、`GetUserTweets`）。
2. Legacy 方法仅在 V2 不满足需求时使用。
3. 新增功能默认接入 V2 路径，并在 PR 中记录对应 path。

### 常见对应关系

| 场景 | 推荐方法（优先） | 兼容方法（次选） |
|---|---|---|
| 用户查询（用户名） | `GetUserByScreenNameV2` | `GetUserByScreenName` |
| 用户查询（ID） | `GetUserByIDV2` | `GetUserByID` |
| 批量用户查询 | `GetUsersByIDsV2` | `GetUsersByIDs` |
| 用户推文时间线 | `GetUserTweets` | `GetUserTimeline` |

### 维护建议

- 当文档更新接口路径时，优先同步本 README 的“Endpoint 路径对照表”。
- 变更路径后，建议至少验证：`user`、`tweets`、`search`、`followers` 四类命令。
- 如果出现 `403` / `code=88` 高频错误，先排除账号与频控问题，再判断是否为接口变更。

### 版本迁移 Checklist（建议用于每次接口升级）

- [ ] 对照 uTools 文档确认 endpoint path 与参数名
- [ ] 更新 SDK 方法中的 path 常量
- [ ] 更新 README 的“Endpoint 路径对照表”
- [ ] 本地回归验证：`user` / `tweets` / `search` / `followers`
- [ ] 核查鉴权接口：`GetHomeTimeline` / `GetMentionsTimeline` / `GetAccountAnalytics`
- [ ] 运行测试：`go test ./...`

### PR 描述模板（接口路径/版本变更）

```md
## Why
- [背景] 本次同步 uTools 文档中的接口变更

## What changed
- [ ] 更新 endpoint path：
  - `MethodA`: old -> new
  - `MethodB`: old -> new
- [ ] 更新 README 对照表与 FAQ

## Verification
- [ ] go test ./...
- [ ] ./xcatch.exe user elonmusk
- [ ] ./xcatch.exe tweets <user_id> 1
- [ ] ./xcatch.exe search "bitcoin" Latest
- [ ] ./xcatch.exe followers <user_id>

## Risk
- [ ] 涉及鉴权接口（auth_token）
- [ ] 涉及分页 cursor 逻辑
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

## FAQ

### 0) 最小排障流程（建议先按这个顺序检查）

1. 确认可编译：`go build -o xcatch.exe ./cmd/`
2. 确认 API Key 已生效（`config.ini` 或 `XCATCH_API_KEY`）
3. 先调用不依赖 `auth_token` 的接口（如 `user` / `search`）
4. 若报频率错误，降低并发并等待重试
5. 若仅鉴权接口失败，再检查 `auth_token`

### 0.1) 常用自检命令

```bash
# 1) 编译检查
go build -o xcatch.exe ./cmd/

# 2) 基础连通性（不依赖 auth_token）
./xcatch.exe user elonmusk

# 3) 搜索接口
./xcatch.exe search "bitcoin" Latest

# 4) 分页接口（max_pages 必须为正整数）
./xcatch.exe tweets 44196397 1
```

### 1) 报错 `config: XCATCH_API_KEY is required`

说明未读取到 API Key。请检查：

- `config.ini` 中是否已填写 `[xcatch] api_key`
- 是否通过环境变量设置了 `XCATCH_API_KEY`
- 当前执行目录下是否存在 `config.ini`

### 2) 报错 `utools: auth_token is required for this endpoint`

这是预期行为。以下接口要求提供 `auth_token`：

- Home Timeline
- Mentions Timeline
- Account Analytics

可在 `config.ini` 中设置 `auth_token`，或使用环境变量 `XCATCH_AUTH_TOKEN`。

### 3) 遇到 `Rate limit exceeded` / `code=88` 怎么办？

SDK 会自动进行指数退避重试（`1s -> 2s -> 4s ...`，上限 30s）。

建议：

- 降低并发请求数
- 适当降低 `XCATCH_RATE_LIMIT`
- 观察响应头 `x-rate-limit-reset`（接近阈值时可调用 `tokenSync`）

### 4) 遇到 `403 Forbidden` 怎么办？

根据 uTools FAQ，这通常是机器人账号临时受限或权限不足。SDK 会自动重试。

如果持续出现：

- 稍后重试
- 更换节点（`base_url`）
- 检查账号状态与接口权限

### 5) 分页 `cursor` 怎么使用？

你可以直接用 `NewPageIterator`：

- 首次请求不传 `cursor`
- 后续请求自动使用上次返回的 `NextCursor`
- `HasMore()` 为 `false` 时停止

### 6) `tweets` 命令里的 `max_pages` 有什么限制？

`max_pages` 必须是正整数，否则 CLI 会直接报错并退出。

### 7) `config.ini` 写错格式会怎样？

程序会打印 warning，并回退到“默认值 + 环境变量覆盖”的加载策略。

## 技术支持

- Discord: discord.gg/VEKBj3fT9G
- Telegram: https://t.me/Yang0619
