---
title: ドキュメント一覧
nav_order: 0
permalink: /ja/
---

# DMXBOXドキュメント

[プロジェクト概要](https://github.com/Surigoma/DMXBOX/blob/main/README.ja.md)

## はじめに

- [利用者ガイド](user-guide.md) — 導入、最初のグループ、画面の操作

## 設定・機材

[設定・機材の概要](settings.md)

- [設定リファレンス](configuration.md) — 設定例、出力先、グループ、検証条件
- [照明機材・出力仕様](lighting-specifications.md) — 対応モデル、チャンネル、DMX、Art-Net
- [Art-Net接続ガイド](artnet-guide.md) — IP設定、Universe番号、NIC選択、受信確認

## 外部連携

[HTTPとTCPの選び方](integration.md)

- [HTTP API](http-api.md) — 制御例、状態ストリーム、設定保存
- [TCPプロトコル](tcp-protocol.md) — 接続、コマンド、応答

## 運用

- [運用・トラブル対応](operations-guide.md) — 運用前確認、操作ログ、手動切り替え、復旧

## 開発

- [開発者ガイド](developer-guide.md) — 構成、開発、テスト、ビルド、拡張、リリース

## 目的から探す

| やりたいこと | 読むページ |
| --- | --- |
| 初めて起動してグループを操作する | [利用者ガイド](user-guide.md) |
| 出力先や灯体のチャンネルを設定する | [設定リファレンス](configuration.md) |
| 制御ソフトからコマンドを送る | [HTTP API](http-api.md) / [TCPプロトコル](tcp-protocol.md) |
| 本番中の不調を調べる | [運用・トラブル対応](operations-guide.md) |
| 対応機器モデルや信号断時の挙動を知る | [照明機材・出力仕様](lighting-specifications.md) |
