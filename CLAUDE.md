# CLAUDE.md

このファイルは、このリポジトリでコードを扱う際のClaude Code (claude.ai/code) へのガイダンスを提供します。

## プロジェクト概要

これはCodeCraftersチャレンジ用にGoで書かれたLox言語インタープリタです。Loxは「Crafting Interpreters」という本で紹介されている動的型付けスクリプト言語です。

## 開発コマンド

### インタープリタの実行
```bash
# フルインタープリタの実行
make test_run
# または直接: go run ./app/main.go run <file>

# 特定フェーズの実行
make test_tokenize   # トークン化のみ
make test_parse      # ASTへのパース
make test_evaluate   # 式の評価
```

### テストと提出
```bash
# CodeCraftersテストの実行
make test

# ソリューションの提出
make submit
```

## アーキテクチャ

インタープリタは古典的な解釈フェーズに従います：

1. **トークン化** (`app/token/`) - ソースコードをトークンに変換
2. **パース** (`app/parse/`, `app/evaluate/parse.go`, `app/run/parse.go`) - 抽象構文木（AST）を構築
3. **評価/実行** (`app/evaluate/`, `app/run/`) - ASTを解釈

主要コンポーネント：
- **環境** (`app/run/env.go`) - 変数のスコープと格納
- **文** (`app/run/statements.go`) - 文の型と実行ロジック
- **メインエントリ** (`app/main.go`) - CLI（コマンド: tokenize, parse, evaluate, run）

## 言語機能

インタープリタがサポートする機能：
- 変数 (`var x = 1;`)
- 関数 (`fun name() { ... }`)
- 制御フロー (`if`, `else`, `while`)
- 算術・論理演算子
- `{}`によるブロックスコープ
- print文
- コメント (`//`)

## テスト方法

ルートディレクトリの`testfile`を手動テストに使用してください。サポートされているすべての言語機能の例が含まれています。