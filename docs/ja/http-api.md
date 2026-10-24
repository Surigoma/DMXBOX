---
title: HTTP API
nav_order: 1
parent: 外部連携
---

# HTTP API・外部連携

[ドキュメント一覧](README.md) | [TCPプロトコル仕様](tcp-protocol.md)

## 接続前の準備

ここでは`http://127.0.0.1:8080`でDMXBOXが動作し、`stage`グループがあるものとします。
例の`stage`は表示名ではなく`dmx.groups`のキーです。自分の設定のIDに置き換えてください。
操作例は実際に出力を変更するため、まず`console`出力で確認してください。
Windows PowerShellから実行するときは、`curl`ではなく`curl.exe`を使用できます。

APIのパラメーターは起動中の`http://localhost:8080/docs/index.html`でも確認できます。

## 主要エンドポイント

| Method | パス | 用途 |
| --- | --- | --- |
| GET | `/api/v1/health` | HTTPサーバーの生存確認 |
| GET | `/api/version` | 実行中のバージョン |
| GET | `/api/features` | 有効な機能 |
| GET | `/api/endpoints` | 登録されたルート |
| GET | `/api/v1/config/fade` | グループ設定 |
| GET | `/api/v1/fade-state` | グループ別の現在のソフトウェア出力状態 |
| POST | `/api/v1/fade/{group}` | Fade／Cut要求 |
| GET | `/api/v1/mute-state` | 最後に送信を完了したOSC Mute状態 |
| POST | `/api/v1/mute` | Mute／Unmute要求 |
| GET | `/api/v1/control-mode` | Browser onlyの状態 |
| GET | `/api/v1/control-state/stream` | FadeとMuteの状態ストリーム |
| GET | `/api/v1/operations?limit=100` | 最近の操作ログ |
| GET | `/api/v1/config/all` | 設定全体 |
| POST | `/api/v1/config/save` | 設定全体の保存 |

## 読み取りから確認する

```sh
curl "http://127.0.0.1:8080/api/v1/health"
curl "http://127.0.0.1:8080/api/v1/config/fade"
curl "http://127.0.0.1:8080/api/v1/control-mode"
```

healthは`status: "ok"`と時刻を返します。
グループ設定はIDをキーにしたオブジェクトです。
`browserOnly: true`の場合、外部スクリプトからの操作は制限されます。

## Fade In／OutとCut

```sh
curl -X POST "http://127.0.0.1:8080/api/v1/fade/stage?isIn=true&duration=1.5&interval=0"
curl -X POST "http://127.0.0.1:8080/api/v1/fade/stage?isIn=false&duration=1.5&interval=0"
curl -X POST "http://127.0.0.1:8080/api/v1/fade/stage?isIn=false&duration=0&interval=0"
```

| パラメーター | 意味 |
| --- | --- |
| `isIn` | `true`で最大値へ、`false`で0へ。省略時は`true` |
| `duration` | 遷移時間（秒）。省略時は`dmx.fadeInterval`。`0`はCut |
| `interval` | 開始までの待ち時間（秒）。即時開始は`0` |

パラメーターはJSON本文ではなくクエリへ指定します。
`&`をシェルが解釈しないよう、URL全体を引用符で囲んでください。
HTTP 200の`{"result":"OK"}`は内部モジュールへの要求送信を示し、機器の動作完了ではありません。
存在しないグループでも、この送信成功応答だけでは検出できないため、先にグループ一覧を確認します。

## OSC Mute

```sh
curl -X POST "http://127.0.0.1:8080/api/v1/mute?isMute=true"
curl -X POST "http://127.0.0.1:8080/api/v1/mute?isMute=false"
curl "http://127.0.0.1:8080/api/v1/mute-state"
```

`isMute`の省略時は`true`です。`osc`出力が有効で、正しい送信先が設定されている必要があります。
状態の`isMute`は`true`、`false`、未確定の`null`です。送信完了の記録であり、ミキサーからの返答ではありません。

## 状態ストリーム

```sh
curl -N "http://127.0.0.1:8080/api/v1/control-state/stream"
```

形式はSSEではなく`application/x-ndjson`です。接続直後にスナップショットを返し、
以後は250ミリ秒ごとに変更を確認して、変化があったときだけ1行のJSONを送信します。
15秒ごとの空行は接続維持用で、状態変更ではありません。

```json
{"fade":{"stage":{"level":0.5,"state":"fading","isIn":true}},"mute":{"isMute":null}}
```

これは形式を示す例です。`level`は0～1の出力割合、`state`は`idle`、`waiting`、`fading`です。
クライアントは受信チャンクではなく改行単位に分割し、空行を除いてJSON解析します。
切断時は再接続し、新しく届くスナップショットを基準にします。
プロキシを使う場合は、ストリームをまとめてバッファリングしない設定が必要です。

## 設定保存時の注意

`GET /api/v1/config/all`で取得した設定全体を編集し、`Content-Type: application/json`で
`POST /api/v1/config/save`へ送ります。部分更新用のPATCH APIではありません。
必須の構造を省略したリクエストは検証に失敗したり、意図しない値になったりするため、
単一項目だけを含むJSONを送らないでください。

成功時は`result: true`を返します。`restartRequired: true`の場合は保存後に再起動が必要です。
通常の再読み込みでも通信が切断される場合があるため、連携クライアントは再接続を扱います。
不正なJSON・設定値は400、保存先への書き込み失敗は500です。

## 接続エラーと制御制限

| 状態 | 確認すること |
| --- | --- |
| 接続拒否 | 起動状態、HTTP待受IP・ポート、ファイアウォール |
| 403で操作できない | Browser onlyが有効でないか |
| 500で操作できない | 出力モジュールが有効か、内部メッセージ送信エラーがないか |
| 200だが機器が動かない | グループID、出力先、物理機器、標準出力ログ |
| ブラウザーの別Originから失敗 | `input.http.accepts`のOrigin設定 |

Browser onlyはWeb UIで切り替えてください。Web UIを識別するヘッダーは認証情報ではありません。
CORSも認証の代わりにはならないため、公開範囲は信頼できるネットワークに限定します。
自動連携からの要求を無条件に再送すると操作が重複するため、ログと状態を確認して再試行します。
