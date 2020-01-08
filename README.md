# go-esacli

esa (https://esa.io) の posts に読み書きを行うツール

## Usage

1. esa から Read/Write できる Personal Token を発行します

`https://${team}.esa.io/user/applications`

2. 環境変数に PersonalToken, チーム名を登録します

```
export ESA_TOKEN=${personalToken}
export ESA_TEAM=${team}
```

3. post を読み込みます

```
go-esacli --path ${category}/${postName}
# e.g. go-esacli --path 日報/2018/12/31/bhi-kojimat
```

4. post を書き込みます

```
go-esacli --path ${category}/${postName} --input ${updatePostMd}
# e.g. go-esacli --path 日報/2018/12/31/bhi-kojimat --input today.md

# OR

cat ${updatePostMd} | go-esacli --path ${category}/${postName} --input -
# e.g. pbpaste | go-esacli --path 日報/2018/12/31/bhi-kojimat --input -
```
