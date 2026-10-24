---
title: 設定リファレンス
nav_order: 1
parent: 設定・機材
---

# DMXBOX設定リファレンス

[ドキュメント一覧](README.md)

DMXBOXはバックエンドのカレントディレクトリにある`config.json`を読み込みます。ソースツリー内の開発用ファイルは`backend/config.json`です。ブラウザから保存すると設定全体が置き換えられ、バックエンドモジュールが再読み込みされます。

入力モジュールの追加・削除やOSC出力の有効・無効を変更した場合は、保存後にDMXBOXを再起動してください。画面に再起動が必要な旨を表示し、この保存に含まれる設定は再起動後に適用します。

手作業で編集する前に、正常に動作している設定をバックアップしてください。

## 全体例

```json
{
    "output": {
        "target": ["console"],
        "ftdi": { "port": "COM3" },
        "artnet": {
            "addr": "127.0.0.1",
            "universe": 0,
            "subuni": 0,
            "net": 0
        },
        "osc": {
            "ip": "127.0.0.1",
            "port": 8765,
            "format": "/yosc:req/set/MIXER:Current/InCh/Fader/On/{}/1",
            "type": "int",
            "inverse": true,
            "channels": [1, 2, 3, 4]
        }
    },
    "input": {
        "modules": ["http"],
        "http": {
            "ip": "127.0.0.1",
            "port": 8080,
            "accepts": [
                "http://localhost:8080",
                "http://127.0.0.1:8080"
            ]
        },
        "tcp": {
            "ip": "127.0.0.1",
            "port": 50000
        }
    },
    "dmx": {
        "groups": {
            "stage": {
                "name": "Stage",
                "devices": [
                    { "model": "dimmer", "channel": 1, "max": [255] },
                    { "model": "wclight", "channel": 10, "max": [255, 255, 0] }
                ]
            }
        },
        "fadeInterval": 0.7,
        "delay": 0,
        "fps": 30
    }
}
```

## `input`

### `input.modules`

有効にする入力モジュールの一覧です。

- `http`：HTTP APIと本番用フロントエンド
- `tcp`：プレーンTCPコマンドサーバー

通常のローカル構成では`http`を含めてください。

### `input.http`

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `ip` | string | 待受アドレス。ローカル限定では`127.0.0.1`を使用します。 |
| `port` | integer | HTTPポート。通常は`8080`です。 |
| `accepts` | string array | CORSで許可するOriginです。 |

ループバック以外のアドレスへバインドすると、APIと設定画面がネットワークへ公開されます。DMXBOXは現在アプリケーションレベルの認証を提供していないため、信頼できるネットワークと適切なホストファイアウォールがある場合に限って使用してください。

### `input.tcp`

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `ip` | string | TCP待受アドレスです。 |
| `port` | integer | TCP待受ポート。通常は`50000`です。 |

コマンドの区切り、応答、コマンド一覧、互換エイリアスについては[TCPプロトコル仕様](tcp-protocol.md)を参照してください。

## `output`

### `output.target`

有効にする出力の一覧です。

- `console`：ハードウェアを使わずDMX出力をログ表示
- `ftdi`：シリアルDMX出力
- `artnet`：Art-Net出力
- `osc`：OSCのMute/Unmute出力

複数の出力を同時に有効化できます。DMXモジュールには少なくとも1つのDMXレンダラー（`console`、`ftdi`、`artnet`）が必要です。

### `output.ftdi`

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `port` | string | `COM3`や`/dev/ttyUSB0`などのシリアルデバイスです。 |

### `output.artnet`

| 項目 | 型 | 範囲 | 説明 |
| --- | --- | --- | --- |
| `addr` | string | — | 宛先IPv4アドレスです。ローカルインターフェースから到達可能である必要があります。 |
| `port` | integer | 省略可 | UDPポート。省略時は6454です。 |
| `net` | integer | 0～127 | Art-NetのNetフィールドです。 |
| `subuni` | integer | 0～15 | Art-NetのSub-Netフィールドです。 |
| `universe` | integer | 0～15 | Art-NetのUniverseフィールドです。 |

### `output.osc`

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `ip` | string | OSC受信先アドレスです。 |
| `port` | integer | OSC受信先UDPポートです。 |
| `format` | string | OSCアドレスのテンプレートです。すべての`{}`がチャンネル値へ置換されます。 |
| `type` | string | `int`は`int32`、`float`は`float32`として送信します。 |
| `inverse` | boolean | 受信先へ送るMute値を反転します。 |
| `channels` | integer array | アドレステンプレートへ代入する値です。 |

## `dmx`

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `groups` | object | 固定のグループIDからグループ定義へのマップです。 |
| `fadeInterval` | number | デフォルトのフェード時間（秒）です。 |
| `delay` | number | 設定上のエフェクト遅延（秒）です。 |
| `fps` | number | DMX更新レートです。フロントエンドは0～99を許容しますが、正の値を使用してください。 |

### グループ

`dmx.groups`の各キーはHTTP・TCP制御APIから使用するグループIDです。各グループは次を持ちます。

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `name` | string | ブラウザに表示する名前です。 |
| `devices` | array | グループとして同時に制御するデバイスです。 |

### デバイス

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `model` | string | `dimmer`または`wclight`です。 |
| `channel` | integer | 1始まりのDMX開始チャンネル。1～512です。 |
| `max` | integer array | Fade In時のチャンネル別最大値。通常は0～255です。 |

デバイスのチャンネル構成：

- `dimmer`：1チャンネル。`max`は1要素必要です。
- `wclight`：寒色、暖色、フラッシュの3チャンネルです。通常、フロントエンドではフラッシュを0に保ちます。

デバイス全体が512チャンネルの範囲内に収まる必要があります。例えば`wclight`を開始できるのは510チャンネルまでです。

シリアルDMX、Art-Net、デバイス動作、フェード、終了時の挙動については[照明機材・出力仕様](lighting-specifications.md)を参照してください。

## 設定を組み立てる例

以下は全体例を基に変更する項目です。JSON断片だけで設定全体を置き換えないでください。

| 用途 | 変更する項目 | 保存後の確認 |
| --- | --- | --- |
| 機器なしで確認 | `output.target: ["console"]` | ゲージとログ、グループID |
| シリアルDMX | `output.target: ["ftdi"]`、`output.ftdi.port` | OS上のポート名と実機の応答 |
| Art-Net | `output.target: ["artnet"]`、`output.artnet` | 送信PCのNIC、宛先、3つのアドレスフィールド |
| TCPから制御 | `input.modules: ["http", "tcp"]`、`input.tcp` | 再起動後の待受と改行付きコマンド |
| OSCを併用 | DMX出力に加えて`output.target`へ`osc`を追加 | 再起動後のMute／Unmuteとミキサー |

`console`、`ftdi`、`artnet`は同時に選択できます。
同じDMXユニバースを各出力へ渡すため、出力先ごとに異なるチャンネル配置にはなりません。

### チャンネル配置の考え方

例の`dimmer`がチャンネル1を、`wclight`が10～12を使用する場合、2～9は未使用です。
次のデバイスは使用中の範囲と重ならない位置へ配置します。
同じ灯体を複数のグループへ登録すると、そのチャンネルを複数の操作が更新するため、
グループ単位で独立して動くとは限りません。

`max`は割合ではなく0～255の値です。例えば`[128]`は255を上限とした約半分の値ですが、
画面の出力割合では設定済み最大値として100%になります。
`wclight`は3要素が必要で、ストロボを使わない構成では3つ目を0にします。

### 保存時に拒否される値

バックエンドは`dmx.fps`の0以下・1000超、負のフェード時間・遅延、
未対応のモデル名、1～512に収まらないチャンネル、足りない`max`、255超の最大値を拒否します。
グループIDには`.`、`[`、`]`を使用できず、`__proto__`、`constructor`、`prototype`も使用できません。
TCP連携でも扱いやすいよう、IDには`stage`など空白のない英数字を使うことを推奨します。
チャンネル重複と物理機器のモードの不一致は検証されないため、パッチ表で確認してください。

## 開発用フロントエンド設定

`frontend/public/config.json`はバックエンド設定とは別のファイルです。

```json
{
    "backendPort": 8080
}
```

Viteで起動する開発用フロントエンドが、ローカルバックエンドのポートを知るために使用します。本番リリースでは、フロントエンドとバックエンドの両方をGoのHTTPサーバーが配信します。
