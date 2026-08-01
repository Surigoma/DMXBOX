# DMXBOX

[English](README.md)

DMXBOXは、ブラウザから設定・操作できるDMX照明制御システムです。GoバックエンドがHTTP APIと本番用フロントエンドを配信するため、リリース版は1台のPC上で単一のアプリケーションとして動作します。

HTTP/TCP制御入力、FTDI互換DMXハードウェアとArt-Net出力、グループ単位のフェード、音響ミキサー向けOSCミュート制御に対応しています。

## 主な機能

- 設定可能なフレームレートによる512チャンネルDMX出力
- Fade、Cut、Mute、Unmute操作
- FTDI互換シリアルDMX、Art-Net、OSC、コンソール出力
- Dimmerとホワイトカラー照明のデバイスモデル
- ブラウザベースの設定・操作画面
- Windows x86-64、Linux x86-64、Linux ARM64向けビルド
- フロントエンド・バックエンドの自動テスト

## クイックスタート

リリース版では、実行ファイルと`static/`ディレクトリを同じ場所に置きます。既存の`config.json`がある場合は、それも同じ場所へ配置します。そのディレクトリから実行ファイルを起動し、次のURLを開きます。

```text
http://localhost:8080/gui/
```

`config.json`がない場合はデフォルト設定が生成されます。本番の照明機器へ接続する前に、出力先とシリアルポートの設定を確認してください。

ローカル開発では次を実行します。

```sh
task dev
```

起動後、`http://localhost:5173`を開きます。

## ドキュメント

- [日本語ドキュメント一覧](docs/ja/README.md)
- [利用者ガイド](docs/ja/user-guide.md) — 導入、操作、更新、トラブルシューティング
- [設定リファレンス](docs/ja/configuration.md) — `config.json`の全体構造と記述例
- [TCPプロトコル仕様](docs/ja/tcp-protocol.md)
- [照明機材・出力仕様](docs/ja/lighting-specifications.md)
- [開発者ガイド](docs/ja/developer-guide.md) — アーキテクチャ、開発、テスト、ビルド、API、リリース
- [English documentation](docs/en/README.md)

バックエンド起動中は、次のURLで対話形式のAPIドキュメントも確認できます。

```text
http://localhost:8080/docs/index.html
```

## スクリーンショット

### コントロール画面

![コントロール](images/control.png)

### 設定画面

![設定](images/settings.png)

## ライセンス

DMXBOXは[MIT License](LICENSE)で提供されます。
