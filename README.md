# BTC跨链交易构造工具

​	btctool用于构造跨链用的BTC特定格式的交易。

## 编译

​	进入到项目的cmd目录下，运行下列命令，需要安装golang和相关依赖

```go
go build -o btctool main.go
```

## 运行

​	btctool可以针对比特币测试网和本地仿真网络，如果使用跨链生态提供的测试网，那么btctool选择测试网即可，如果是在本地运行跨链测试网络，则选择仿真网络。

测试网络：

```shell
./btctool -tool=cctx -idxes=1 -utxovals=0.01 -txids=c09d7d7a321d025ac0cad75855b1b0313e55660a5a77b9d038bb7f606be6a744 -value=0.008 -fee=0.00001 -targetaddr=AdzZ2VKufdJWeB8t9a8biXoHbbMe2kZeyH -privkb58=cRRMYvoHPNQu1tCz4ajPxytBVc2SN6GWLAVuyjzm4MVwyqZVrAcX -contract=b6bf9abf29ee6b8c9828a48b499ad667da1ad003
```

仿真网络：

```shell
./btctool -tool=regauto -fee=0.001 -privkb58=cRRMYvoHPNQu1tCz4ajPxytBVc2SN6GWLAVuyjzm4MVwyqZVrAcX -pwd=test -user=test -targetaddr=AdzZ2VKufdJWeB8t9a8biXoHbbMe2kZeyH -url=http://172.168.3.77:18443 -value=0.01 -contract=b6bf9abf29ee6b8c9828a48b499ad667da1ad003 -tochain=2
```

## 参数

测试网络

|    FLAGS    |                            USAGE                             |
| :---------: | :----------------------------------------------------------: |
|    -tool    |   选择对应的工具，cctx是测试网工具，regauto是仿真网络工具    |
|   -idxes    | UTXO在交易中的位置，即第几个输出，例如0、1等，可多个，用","隔开 |
|  -utxovals  |              每个UTXO的金额，可多个，用","隔开               |
|  -privkb58  | 签名用的base58形式私钥，如果不设置，会返回未签名的交易，用户可以自行签名 |
|   -value    | 跨链交易金额，即锁定到联盟链多签地址的金额，将value个BTC转移到目标链，默认1万聪 |
|    -fee     |               该比特币交易的手续费，默认1000聪               |
| -targetaddr |                       用户的目标链地址                       |
|  -contract  |                    目标链代币智能合约地址                    |
|  -tochain   |            目标链ID，联盟链用来确定BTC跨链目的地             |

仿真网络

|    FLAGS    |                            USAGE                             |
| :---------: | :----------------------------------------------------------: |
|    -tool    |   选择对应的工具，cctx是测试网工具，regauto是仿真网络工具    |
|  -privkb58  | 签名用的base58形式私钥，如果不设置，会返回未签名的交易，用户可以自行签名 |
|    -pwd     |                     比特币客户端rpc密码                      |
|    -user    |                    比特币客户端rpc用户名                     |
| -targetaddr |                       用户的目标链地址                       |
|    -url     |                     比特币客户端rpc地址                      |
|   -value    | 跨链交易金额，即锁定到联盟链多签地址的金额，将value个BTC转移到目标链，默认1万聪 |
|  -contract  |                    目标链代币智能合约地址                    |
|  -tochain   |            目标链ID，联盟链用来确定BTC跨链目的地             |

## 运行实例

```shell
./btctool -tool=cctx -idxes=1 -utxovals=0.01 -txids=c09d7d7a321d025ac0cad75855b1b0313e55660a5a77b9d038bb7f606be6a744 -value=0.008 -fee=0.00001  -targetaddr=AdzZ2VKufdJWeB8t9a8biXoHbbMe2kZeyH -privkb58=cRRMYvoHPNQu1tCz4ajPxytBVc2SN6GWLAVuyjzm4MVwyqZVrAcX -contract=56faac6081cd320fab3347c62faea86344a8aece 
2019/11/04 14:34:10.220569 [INFO ] GID 1, Signed cross chain transaction with your private key
2019/11/04 14:34:10.220716 [INFO ] GID 1, ------------------------Your signed cross chain transaction------------------------
010000000144a7e66b607fbb38d0b9775a0a66553e31b0b15558d7cac05a021d327a7d9dc0010000006a473044022061723a7ba8d6c07cd2cf53ea6211ead6eea300cb4ea6bfb8e7ea6160f0424700022016e264422798e1efe5d2a5d54e21aee9aae42d9d41d6976d2760363d7c846779012103128a2c4525179e47f38cf3fefca37a61548ca4610255b3fb4ee86de2d3e80c0fffffffff0300350c000000000017a91487a9652e9b396545598c0fc72cb5a98848bf93d38700000000000000003d6a3b660200000000000000000000000000000014ceaea84463a8ae2fc64733ab0f32cd8160acfa5614f3b8a17f1f957f60c88f105e32ebff3f022e56a458090300000000001976a91428d2e8cee08857f569e5a1b147c5d5e87339e08188ac00000000

2019/11/04 14:34:10.220742 [INFO ] GID 1, spv addr not set, you need to broadcast tx by yourself
```

​	如上，结果中的十六进制字符串即为签名后的交易，用户可通过其他工具自行广播，例如全节点rpc命令sendrawtransaction或者一些[网站](https://tbtc.bitaps.com/broadcast)，或者通过设置-spvaddr让[轻客户端](https://github.com/Zou-XueYan/spvwallet/tree/nowallet)去广播交易。