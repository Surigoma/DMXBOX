---
title: 外部連携
nav_order: 3
has_toc: false
---

# 外部連携

制御ソフトウェアからDMXBOXを操作するためのインターフェースです。
グループIDと出力先を設定し、まずconsole出力で接続を確認します。

| インターフェース | 選ぶ場面 | 詳細 |
| --- | --- | --- |
| HTTP | 制御に加え、設定・状態・操作ログも取得したい | [HTTP API](http-api.md) |
| TCP | 改行区切りの短い制御コマンドを送りたい | [TCPプロトコル](tcp-protocol.md) |

TCPは入力モジュールの追加と再起動が必要です。
外部操作が制限される場合は、[Browser onlyの状態](operations-guide.md)も確認してください。
