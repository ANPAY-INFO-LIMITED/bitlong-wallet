# wallet

## pack api

```bash
pwsh ./pack.api.ps1
```

## Update log

### v0.1.1
- 增加查询宇宙资产记录的接口返回数据的最大值

### v0.1.0
- 升级lit,lnd,tapd版本
- 编译命令请增加 -tags "signrpc walletrpc chainrpc invoicesrpc autopilotrpc btlapi"
example:
```
1.gomobile.exe bind -v -target=android -tags "signrpc walletrpc chainrpc invoicesrpc autopilotrpc btlapi"
2.go build  -tags "signrpc walletrpc chainrpc invoicesrpc autopilotrpc btlapi" -o main.exe 
```

### v0.0.6
- login接口增加超时
- [获取所有BTC和资产价值](https://bitlong.gitbook.io/api-doc#id-202489-112258-huo-qu-suo-you-btc-he-zi-chan-jia-zhi)

### v0.0.5

- [托管账户通道交易记录查询新增费用字段](https://app.gitbook.com/o/YKwEeYmlIe2KJGpL2TiX/s/nQHIKw00nyi8RS189Cfp/trade/custodyaccount)
- [上传本地资产UTXO](https://bitlong.gitbook.io/api-doc#id-202482-172842-shang-chuan-ben-di-zi-chan-utxo)
- [上传本地资产发行成功记录](https://bitlong.gitbook.io/api-doc#id-202482-172829-shang-chuan-ben-di-zi-chan-fa-xing-cheng-gong-ji-lu)

### v0.0.4

- [新增从种子恢复钱包的接口](https://seven-liquors-doc.gitbook.io/btlapi-1/wallet/recoverwallet)


### v0.0.3

- [关注和取消关注公平发射资产](https://bitlong.gitbook.io/api-doc#id-202482-084610-guan-zhu-he-qu-xiao-guan-zhu-gong-ping-fa-she-zi-chan)

### v0.0.2 

- [查询资产发行和铸造费用参数更新](https://bitlong.gitbook.io/api-doc#id-202481-141909-cha-xun-zi-chan-fa-xing-he-zhu-zao-fei-yong-can-shu-geng-xin)

