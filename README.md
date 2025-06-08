# gha_openfga

GitHub Actions で OpenFGA を含めたテストを実行する実験

## ローカルでのテスト

```
make test
```

## OpenFGA 動作確認

```sh
export FGA_STORE_ID=$(fga store create --name "FGA Demo Store" | jq -r .store.id)
fga model write --store-id $FGA_STORE_ID --file=./model.fga
fga tuple write --store-id $FGA_STORE_ID user:anne reader document:doc1
fga query check --store-id $FGA_STORE_ID user:anne reader document:doc1
fga query check --store-id $FGA_STORE_ID user:bob reader document:doc1
```
