# ecs-tag-shift

AWS ECSのタスク定義やコンテナ定義のイメージタグを更新するCLIツール

## 概要

`ecs-tag-shift` は、ECSタスク定義やコンテナ定義のJSONファイルを読み込み、コンテナイメージのタグを効率的に更新するツールです。JSONC（コメント付きJSON）形式の入力に対応し、パイプラインでの利用を想定した設計になっています。

## 主な機能

- ECSタスク定義JSONの読み込みと表示
- コンテナ定義（containerDefinitions）の読み込みと表示
- コンテナイメージタグの一括更新・個別更新
- JSONC（コメント付きJSON）入力のサポート
- 標準入力・ファイル指定の両方に対応

## インストール

```bash
git clone https://github.com/dev-shimada/ecs-tag-shift.git
cd ecs-tag-shift
go build -o ecs-tag-shift ./cmd/ecs-tag-shift
```

## 使用方法

### 基本構文

```bash
ecs-tag-shift [グローバルオプション] <サブコマンド> [引数] [オプション]
```

### グローバルオプション

| オプション | 短縮形 | 説明 | デフォルト値 |
|-----------|--------|------|-------------|
| `--mode` | `-m` | 入力形式を指定 (`task` または `container`) | `task` |
| `--help` | `-h` | ヘルプを表示 | - |
| `--version` | `-v` | バージョン情報を表示 | - |

#### --mode オプションの詳細

- **`task`**: ECSタスク定義JSON全体を処理します
- **`container`**: containerDefinitions セクションのみを処理します（配列形式のみ許可）

**入力形式:**
- JSONC（コメント付きJSON）をサポート
- 出力は常にJSON形式（コメントは削除されます）

**エラー処理:**
- エラーメッセージは標準エラー出力（stderr）に出力されます
- エラー時の終了コードは `1` です

---

## サブコマンド

### show

タスク定義またはコンテナ定義の内容を表示します。

#### 構文

```bash
ecs-tag-shift [--mode <mode>] show [file] [options]
```

#### 引数

| 引数 | 説明 | 必須 |
|-----|------|------|
| `file` | 入力ファイルのパス（省略時は標準入力から読み込み） | ❌ |

#### オプション

| オプション | 短縮形 | 説明 | デフォルト値 |
|-----------|--------|------|-------------|
| `--output` | `-o` | 出力形式 (`json`, `yaml`, `text`) | `json` |
| `--all` | | タスク定義またはコンテナ定義の全フィールドを表示 | `false` |

#### 使用例

**タスク定義モード（`--mode task`）:**

```bash
# JSON形式で表示（デフォルト）
ecs-tag-shift show task-definition.json

# 標準入力から読み込み
cat task-definition.json | ecs-tag-shift show

# YAML形式で表示
ecs-tag-shift show task-definition.json --output yaml

# TEXT形式で表示
ecs-tag-shift show task-definition.json --output text

# 全フィールドを表示
ecs-tag-shift show task-definition.json --all
```

**コンテナ定義モード（`--mode container`）:**

```bash
# コンテナ定義を表示
ecs-tag-shift --mode container show container-definitions.json

# TEXT形式で表示
ecs-tag-shift -m container show container-definitions.json -o text
```

#### 出力例

**タスク定義モード - JSON形式（デフォルト）:**

```json
{
  "family": "my-app",
  "revision": 15,
  "containers": {
    "web": "123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.2.2",
    "nginx": "nginx:latest"
  }
}
```

**タスク定義モード - TEXT形式:**

```text
Family: my-app
Revision: 15

Containers:
  - web: 123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.2.2
  - nginx: nginx:latest
```

**コンテナ定義モード - JSON形式:**

```json
{
  "containers": {
    "web": "123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.2.2",
    "nginx": "nginx:latest"
  }
}
```

**コンテナ定義モード - TEXT形式:**

```text
Containers:
  - web: 123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.2.2
  - nginx: nginx:latest
```

---

### shift

コンテナイメージのタグを更新します。更新結果は標準出力に出力されるため、リダイレクトでファイルに保存できます。

#### 構文

```bash
ecs-tag-shift [--mode <mode>] shift [file] --tag <new-tag> [options]
```

#### 引数

| 引数 | 説明 | 必須 |
|-----|------|------|
| `file` | 入力ファイルのパス（省略時は標準入力から読み込み） | ❌ |

#### オプション

| オプション | 短縮形 | 説明 | デフォルト値 |
|-----------|--------|------|-------------|
| `--tag` | `-t` | 新しいイメージタグ（例: `v1.2.3`, `latest`） | **必須** |
| `--container` | `-c` | 更新対象のコンテナ名（指定しない場合は全コンテナ） | - |
| `--image` | `-i` | 更新対象のイメージリポジトリ名（完全一致） | - |
| `--output` | `-o` | 出力形式 (`json`, `yaml`) | `json` |
| `--overwrite` | `-w` | 入力ファイルを上書き（ファイル指定時のみ有効） | `false` |

#### フィルタリング動作

- `--container` と `--image` は併用可能です
- `--container`: コンテナ名で完全一致フィルタ
- `--image`: イメージリポジトリ名で完全一致フィルタ（例: `nginx`, `my-app`）
- 両方指定した場合は AND 条件になります

#### 上書きオプション (`--overwrite`/`-w`)

- ファイル指定時のみ有効です
- `--overwrite` を指定した場合、結果を標準出力ではなく入力ファイルに上書きします
- 標準入力から読み込んだ場合は `--overwrite` を指定しても効果はありません（常に標準出力に出力）
- 上書き時は `--output` オプションで指定した形式（デフォルト: JSON）でファイルを書き込みます

#### 使用例

**タスク定義モード（`--mode task`）:**

```bash
# 全コンテナのタグを更新
ecs-tag-shift shift task-definition.json --tag v1.2.3

# 標準入力から読み込み
cat task-definition.json | ecs-tag-shift shift --tag v1.2.3

# 更新結果をファイルに保存（標準出力）
ecs-tag-shift shift task-definition.json --tag v1.2.3 > updated.json

# ファイルを直接上書き
ecs-tag-shift shift task-definition.json --tag v1.2.3 --overwrite

# 特定のコンテナのみ更新
ecs-tag-shift shift task-definition.json --container web --tag v1.2.3

# 特定のイメージリポジトリのみ更新（完全一致）
ecs-tag-shift shift task-definition.json --image my-app --tag v1.2.3

# YAML形式で出力
ecs-tag-shift shift task-definition.json --tag v1.2.3 --output yaml

# JSONCファイルを読み込んで YAML で上書き
ecs-tag-shift shift task-definition.jsonc --tag v1.2.3 -o yaml -w

# JSONCファイルを読み込んでJSONで出力
cat task-definition.jsonc | ecs-tag-shift shift --tag v1.2.3 > updated.json
```

**コンテナ定義モード（`--mode container`）:**

```bash
# コンテナ定義配列のタグを更新
ecs-tag-shift --mode container shift containers.json --tag v1.2.3

# 特定のコンテナのみ更新
ecs-tag-shift -m container shift containers.json -c web -t v1.2.3

# コンテナ定義ファイルを直接上書き
ecs-tag-shift -m container shift containers.json -t v1.2.3 -w
```

#### 出力例

**タスク定義モード - JSON形式:**

```json
{
  "family": "my-app",
  "taskRoleArn": "arn:aws:iam::123456789:role/ecsTaskRole",
  "executionRoleArn": "arn:aws:iam::123456789:role/ecsTaskExecutionRole",
  "networkMode": "awsvpc",
  "containerDefinitions": [
    {
      "name": "web",
      "image": "123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.2.3",
      "cpu": 256,
      "memory": 512,
      "essential": true
    },
    {
      "name": "nginx",
      "image": "nginx:v1.2.3",
      "cpu": 128,
      "memory": 256
    }
  ],
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512"
}
```

**コンテナ定義モード - JSON形式:**

```json
[
  {
    "name": "web",
    "image": "123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.2.3",
    "cpu": 256,
    "memory": 512
  },
  {
    "name": "nginx",
    "image": "nginx:v1.2.3",
    "cpu": 128,
    "memory": 256
  }
]
```

---

## 入力ファイル形式

### タスク定義（`--mode task`）

AWS ECS タスク定義のJSON形式に準拠しています。JSONC（コメント付きJSON）もサポートしています。

**task-definition.json の例:**

```json
{
  "family": "my-app",
  "taskRoleArn": "arn:aws:iam::123456789:role/ecsTaskRole",
  "executionRoleArn": "arn:aws:iam::123456789:role/ecsTaskExecutionRole",
  "networkMode": "awsvpc",
  "containerDefinitions": [
    {
      "name": "web",
      "image": "123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.2.2",
      "cpu": 256,
      "memory": 512,
      "essential": true,
      "portMappings": [
        {
          "containerPort": 80,
          "protocol": "tcp"
        }
      ]
    }
  ],
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512"
}
```

**task-definition.jsonc の例（コメント付き）:**

```jsonc
{
  "family": "my-app",
  // IAM Roles
  "taskRoleArn": "arn:aws:iam::123456789:role/ecsTaskRole",
  "executionRoleArn": "arn:aws:iam::123456789:role/ecsTaskExecutionRole",
  "networkMode": "awsvpc",
  "containerDefinitions": [
    {
      "name": "web",
      "image": "123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.2.2", // アプリケーションイメージ
      "cpu": 256,
      "memory": 512,
      "essential": true
    }
  ],
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512"
}
```

### コンテナ定義（`--mode container`）

containerDefinitions セクションのみを含む配列形式です。単一オブジェクトはエラーになります。

**container-definitions.json の例:**

```json
[
  {
    "name": "web",
    "image": "123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.2.2",
    "cpu": 256,
    "memory": 512,
    "essential": true,
    "portMappings": [
      {
        "containerPort": 80,
        "protocol": "tcp"
      }
    ]
  },
  {
    "name": "nginx",
    "image": "nginx:latest",
    "cpu": 128,
    "memory": 256
  }
]
```

**エラーになる例（単一オブジェクト）:**

```json
{
  "name": "web",
  "image": "123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.2.2"
}
```

このような単一オブジェクトを渡すと、以下のエラーが出力されます：

```
Error: input must be an array of container definitions
```

---

## エラーハンドリング

エラーが発生した場合、エラーメッセージは標準エラー出力（stderr）に出力され、終了コード `1` で終了します。

### よくあるエラー

**ファイルが見つからない:**
```
Error: failed to read file: open task-definition.json: no such file or directory
```

**JSONパースエラー:**
```
Error: failed to parse JSON: invalid character '}' after object key
```

**コンテナ定義モードで配列以外を入力:**
```
Error: input must be an array of container definitions
```

**必須オプションが不足:**
```
Error: required flag(s) "tag" not set
```

**指定したコンテナが見つからない:**
```
Error: container 'api' not found in definitions
```

---

## 開発情報

### ディレクトリ構成

```
.
├── cmd/
│   └── ecs-tag-shift/
│       └── main.go              # エントリーポイント
├── internal/
│   ├── taskdef/
│   │   ├── loader.go            # JSON/JSONC読み込み
│   │   └── updater.go           # タグ更新ロジック
│   ├── command/
│   │   ├── show.go              # show サブコマンド
│   │   └── shift.go             # shift サブコマンド
│   └── output/
│       └── formatter.go         # JSON/YAML/TEXT出力
├── go.mod
├── go.sum
└── README.md
```

### 依存パッケージ

- `github.com/spf13/cobra` - CLIフレームワーク
- `gopkg.in/yaml.v3` - YAML出力サポート

### ビルド

```bash
go build -o ecs-tag-shift ./cmd/ecs-tag-shift
```

### テスト

```bash
go test ./...
```

---

## ライセンス

MIT License

---

## コントリビューション

Issue や Pull Request は歓迎します！
