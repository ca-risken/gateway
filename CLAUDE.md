# CLAUDE.md

## ルーティング・パス命名規則

- APIパス名でリソース名を結合する際は、スラッシュ区切り（`organization/alert`）ではなくハイフン区切り（`organization-alert`）を使用する
  - 理由: ActionNameの仕様に基づく。`getActionNameFromURI` は `/api/v1/{service}/{path1}` から `{service}/{path1}` をActionNameとして生成するため、サービス名自体にハイフンで結合する必要がある
  - 例: `/api/v1/organization-alert/list-notification` → ActionName: `organization-alert/list-notification`
