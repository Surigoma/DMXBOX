# GitHub Pages

## 公開設定

1. この変更を`main`へ反映します。
2. リポジトリの **Settings → Pages → Build and deployment** を開きます。
3. **Source: Deploy from a branch** を選択します。
4. **Branch: main**、フォルダー **/docs** を選択して **Save** を押します。
5. Actionsの **pages build and deployment** が成功したら、公開先を確認します。

- トップページ: <https://surigoma.github.io/DMXBOX/>
- 日本語: <https://surigoma.github.io/DMXBOX/ja/>
- English: <https://surigoma.github.io/DMXBOX/en/>

以後は`main`の`docs/`内を更新すると、自動で再公開されます。
アプリ本体のタグによるリリースとは独立して動作します。

## 更新方法

本文は既存の`docs/**/*.md`を直接編集します。公開専用の複製は不要です。
GitHub Pages標準のJekyllがMarkdownをHTMLに変換し、
Markdownファイルへの相対リンクを公開ページへのリンクへ変換します。
見た目は[Just the Docs](https://just-the-docs.com/)の`v0.12.0`を使用します。
テーマは`remote_theme`で読み込むため、独自の公開ワークフローは不要です。

各Markdownの先頭にあるYAML front matterで、左メニューの表示名と順序を指定します。
新しいガイドは同じ言語の既存ページに合わせて追加します。

```yaml
---
title: 新しいガイド
parent: 設定・機材
nav_order: 3
---
```

日本語は`docs/ja/`、英語は`docs/en/`へ同じファイル名で追加します。
左メニューは「はじめに」「設定・機材」「外部連携」「運用・トラブル対応」「開発者ガイド」の順です。
設定関連のページは`parent: 設定・機材`、連携関連は`parent: 外部連携`にします。
英語側はそれぞれ`Configuration and Hardware`、`External Integration`を指定します。
`nav_order`は同じ親の中での表示順です。カテゴリの入口は`settings.md`と`integration.md`です。
各言語のREADMEは`permalink`で`/ja/`と`/en/`に公開します。
言語切り替えは上部メニューのドロップダウンにあり、同じガイドの翻訳先へ移動します。
翻訳がないページでは、その言語のドキュメント一覧へ移動します。
左メニューには閲覧中の言語のページだけを表示します。
サイト全体の検索を使用できます。検索はテーマ標準のLunrを使うため、
日本語の単語分割に最適化した検索ではありません。

`_includes/components/sidebar.html`はテーマ`v0.12.0`を基に、
言語ごとにナビゲーションを生成するため共有キャッシュを外しています。
テーマ更新時は上流のサイドバーの変更も確認してください。

ヘッダーの太陽／月アイコンでライト／ダークを切り替えられます。
アイコンは現在の配色を示し、ツールチップと読み上げラベルは切り替え先を示します。
初回は端末の配色設定に従い、手動で選んだ配色はブラウザーに保存します。
JavaScriptが無効な場合はテーマ標準のライト表示になります。

`docs/`外のファイルは公開対象に含まれません。
プロジェクト概要やソースコードを参照するときは、GitHub上のURLを使ってください。
所有者やリポジトリ名を変える場合は、`docs/_config.yml`の`url`、`baseurl`、
`repository`とGitHubへのリンクも更新します。

公開を停止するときは **Settings → Pages** で公開元ブランチを **None** にします。

## References

- [Configuring a publishing source](https://docs.github.com/en/pages/getting-started-with-github-pages/configuring-a-publishing-source-for-your-github-pages-site)
- [About GitHub Pages and Jekyll](https://docs.github.com/en/pages/setting-up-a-github-pages-site-with-jekyll/about-github-pages-and-jekyll)
