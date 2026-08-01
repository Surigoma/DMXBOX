# DMXBOX開発者ガイド

[English](../en/developer-guide.md) | [ドキュメント一覧](README.md)

## ソース構成

```text
backend/                 Goアプリケーションとテスト
  artnet/                Art-Net通信
  config/                設定モデルと永続化
  dmxServer/             DMXグループ、デバイス、タイミング、出力
  httpServer/            HTTPサーバーとコントローラー
  oscServer/             OSCメッセージ出力
  packageModule/         モジュールのライフサイクルとメッセージルーティング
  tcpServer/             TCPコマンド入力
frontend/                React・TypeScriptアプリケーション
  src/component/         操作・設定コンポーネント
  src/contexts/          フロントエンド実行時設定
  src/routes/            TanStack Routerのルート
  src/test/              ブラウザテスト共通ヘルパー
.github/workflows/       タグで起動するリリースワークフロー
Taskfile.yml             リポジトリ全体のタスク入口
```

## アーキテクチャ

`backend/main.go`が設定を読み込み、有効なモジュールを登録し、モジュールマネージャーを初期化して各モジュールを起動します。モジュール間では、`packageModule.ModuleManagerType`が`message.Message`をルーティングします。

主なモジュール：

- `http`：HTTP APIと静的フロントエンド
- `tcp`：オプションのTCPコマンド待受
- `dmx`：DMX状態、フェード、フレーム制御、レンダラー出力
- `osc`：オプションのOSC Mute出力

DMXサーバーは512バイトのユニバースを保持します。デバイスモデルがユニバース内の担当範囲を更新し、設定されたレンダラーがコンソール、シリアルDMX、Art-Netへ結果を送信します。

フロントエンドはTanStack Router、Material UI、React Hook Form、SWR、Zodを使用します。開発時はViteで独立して動作し、本番時は生成された`static/`をGoのHTTPサーバーが配信します。

## 必要な環境

- [Task](https://taskfile.dev/)
- `backend/go.mod`と互換性のあるGo
- Vite 8と互換性のあるNode.js
- CorepackとYarn 4
- フロントエンドのブラウザテスト用Chromium

開発タスクやレポート生成では、必要に応じてAir、swag、gotestsum、gocover-cobertura、gopoghも使用します。

## セットアップと開発

Corepackを有効化し、ローカル環境に必要な依存関係を用意してから次を実行します。

```sh
task dev
```

開発時のURL：

- フロントエンド：`http://localhost:5173`
- バックエンド：`http://localhost:8080`
- Swagger UI：`http://localhost:8080/docs/index.html`

フロントエンドのバックエンドポート設定は`frontend/public/config.json`、バックエンドの開発用設定は`backend/config.json`です。

## テスト

ブラウザテストを直接実行せず、リポジトリのタスクを使用します。

```sh
task test
```

フロントエンドのブラウザテストをヘッドレス・カバレッジ付きで実行した後、バックエンドのテストとレポート生成を行います。レポート出力先：

```text
frontend/test/
backend/test/
```

ウォッチモード：

```sh
task test_watch
```

フロントエンドのLint：

```sh
cd frontend
yarn lint
```

バックエンドの静的解析：

```sh
cd backend
go vet ./...
```

バックエンドのレポート生成タスクは、ローカルテストが失敗した場合でも結果を確認できるよう、テストコマンド後の処理を継続する設計です。タスク全体の終了状態だけではなく、テスト出力と生成されたレポートを確認してください。

## API開発

ルートは`backend/httpServer/httpServer.go`で登録し、コントローラーは`backend/httpServer/controller/`以下に配置します。

Swaggerアノテーションを変更した後は、APIドキュメントを再生成します。

```sh
cd backend
task update_swag
```

現在のAPIには、ヘルスチェック、設定、Fade制御、OSC Mute制御、機能一覧、エンドポイント一覧、バージョン情報があります。リクエストとレスポンスの詳細は、起動中のSwagger UIを参照してください。

## デバイスモデルの追加

1. `backend/dmxServer/devices/spec/`以下にデバイスコンストラクターを実装します。
2. モデル名と使用チャンネル数を設定します。
3. `dmxServer.DeviceTypes`へコンストラクターを登録します。
4. `frontend/src/types.ts`へモデルと検証を追加します。
5. `frontend/src/component/settings/device/model/`以下へ設定UIを追加します。
6. チャンネル境界を含むバックエンド・フロントエンドテストを追加します。
7. 設定リファレンスを更新します。

## 出力レンダラーの追加

1. `dmxServer/controller.Controller`を通してレンダラーを実装します。
2. `dmxServer.RenderTypes`へ登録します。
3. バックエンド設定モデルを拡張します。
4. フロントエンドのスキーマと設定UIを追加します。
5. 初期化、出力、失敗、再読み込み、終了をテストします。
6. 設定項目とプラットフォーム要件を記載します。

## ビルド

現在の対応プラットフォーム向けにビルドします。

```sh
task build
```

Windows x86-64、Linux x86-64、Linux ARM64をまとめてビルドします。

```sh
task build_all
```

配布ディレクトリの作成、全ターゲットのビルド、現在のバックエンド設定のコピーを実行するには、デフォルトタスクを使います。

```sh
task
```

成果物は`dist/`以下へ出力され、フロントエンドは`dist/static/`へビルドされます。

現在はTaskfileのバージョン変数と埋め込み用の`backend/.version`が、成果物名と実行時バージョンを制御します。リリース準備時はこれらを同期してください。

## リリース手順

このプロジェクトではローカル検証を主なリリースゲートとし、GitHub Actionsは手動で作成した`v*`タグに限定しています。

推奨チェックリスト：

1. 作業ツリーとリリース対象の変更を確認します。
2. `task test`を実行し、フロントエンド・バックエンド両方の結果を確認します。
3. `frontend/`で`yarn lint`、`backend/`で`go vet ./...`を実行します。
4. 空の配布ディレクトリから`task build_all`を実行します。
5. 配布版フロントエンドと`/api/v1/health`を確認します。
6. 必要に応じてFTDI、Art-Net、Fade、Cut、TCP、OSCを実機確認します。
7. バージョンとリリースノートを更新します。
8. リリース状態をコミットします。
9. 任意の`v*`タグを作成してpushします。
10. GitHub ReleaseとダウンロードしたZIPを確認します。

## コントリビューション上の注意

- 無関係なローカル変更を保持します。
- Goコードは`gofmt`で整形します。
- TypeScriptのLintを通過させます。
- 振る舞いを変更した場合は回帰テストを追加します。
- 公開仕様を変更した場合は日本語・英語ドキュメントを揃えて更新します。
- 生成ファイル`frontend/src/routeTree.gen.ts`は手作業で編集しません。

## 関連ドキュメント

- [利用者ガイド](user-guide.md)
- [設定リファレンス](configuration.md)
- [TCPプロトコル仕様](tcp-protocol.md)
- [照明機材・出力仕様](lighting-specifications.md)
- [ルートREADME](../../README.ja.md)
