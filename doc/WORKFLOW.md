### 运行流程

已修改好脚本和配置文件，直接将压缩包内chaincode_regcc、chaincode_pointcc、chaincode_mapcc三个文件夹放入$GOPATH/src/github.com/hyperledger/fabric/examples/chaincode路径下，并将压缩包内e2e_cli文件夹放在$GOPATH/src/github.com/hyperledger/fabric/examples下。之后从头构建dockerimage和各项配置文件，再运行e2e_cli中的脚本`./network_setup.sh up`


脚本运行完成后进入cli容器bash界面，即可运行下列代码检查

单个链码调用方法:
```
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C CHANNELID -n CHAINCODEID -c '{"Args":["函数名","用户名","帐号ID",....其他参数，具体需要哪些参数可以参考示例和链码函数]}' -v 1.0
```

会检查当前用户的证书对应实体名是否和传入的用户名参数一致，如果实体名不是默认管理者且不一致时，则返回错误，调用失败。

除去query时，帐号ID不能为"all"
当query时，帐号ID为"all"则返回用户所有帐号的ID，在这种方法下的query默认返回该帐号所有信息

一个channel调用另一个channel时只能看数据，不能修改！！！

链码相互调用时不能递归，只能调用一层

对象主要有两种：User和Account。除pointcc只有Account之外，剩下两个对象均在其他链码中出现。
大致每个对象有三种操作方法：create、query、delete。函数名是方法和对象的拼接（例如：queryUser）
除此之外queryall作为归集函数比较特殊，
queryHistory是通用的方法，可以返回一个键值对的所有历史记录（其中value部分需要解码才能看）
pointcc中的trade函数和mapcc中的不同，pointcc中需要的第三个参数为积分数值，mapcc中第三个参数需要虚拟化资产的地址。

```
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["queryAccount","User1@org1.example.com","all","all"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["queryUser","User1@org1.example.com"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C pointchannel -n pointcc -c '{"Args":["createAccount","User1@org1.example.com","AccountofBank","300","A_Bank","Nothing"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C pointchannel -n pointcc -c '{"Args":["setAccount","User1@org1.example.com","AccountofBank","400","A_Bank","Nothing"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C pointchannel -n pointcc -c '{"Args":["queryAccount","AccountofBank","all"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["createAccount","User1@org1.example.com","AccountofBank2","pointchannel","points","A_Bank"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C pointchannel -n pointcc -c '{"Args":["createAccount","User1@org1.example.com","AccountofBank2","12","A_Bank","Nothing"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C pointchannel -n pointcc -c '{"Args":["queryAccount","AccountofBank2","all"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["createAccount","User1@org1.example.com","mapchannel","mapchannel","mapping","ZYP"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["createUser","User1@org1.example.com"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["queryUser","User1@org1.example.com"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["createAccount","User1@org1.example.com","a-car-BJ454852","CAR","Department of Motor Vehicles","2018"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["createAccount","User1@org1.example.com","a-house-YIHEYUAN load-5","House","Housing Authority","Peking University"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["queryAccount","a-car-BJ454852"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["createUser","User1@org3.example.com"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["createAccount","User1@org3.example.com","mapchannel","mapchannel","mapping","ZYP"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["createUser","User1@org3.example.com"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["trade","User1@org1.example.com","User1@org3.example.com","a-car-BJ454852"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["queryUser","User1@org1.example.com"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["queryUser","User1@org3.example.com"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["queryAccount","a-car-BJ454852"]}' -v 1.0
sleep 3
```

** 20180604 更新: 跨链 **
```
##############################################################################################################################################################

peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["createUser","User1@org2.example.com"]}' -v 1.0
sleep 3

peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["createAccount","User1@org2.example.com","Boba","pointchannel","points","a_Bank"]}' -v 1.0
sleep 3

peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["createAccount","User1@org2.example.com","Bobb","pointchannel","points","b_Bank"]}' -v 1.0
sleep 3

peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["createAccount","User1@org3.example.com","Carola","pointchannel","points","a_Bank"]}' -v 1.0
sleep 3

peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["createAccount","User1@org3.example.com","Carolb","pointchannel","points","b_Bank"]}' -v 1.0
sleep 3

peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C pointchannel -n pointcc -c '{"Args":["createAccount","User1@org2.example.com","Boba","300","a_Bank","Nothing"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C pointchannel -n pointcc -c '{"Args":["createAccount","User1@org2.example.com","Bobb","300","b_Bank","Nothing"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C pointchannel -n pointcc -c '{"Args":["createAccount","User1@org3.example.com","Carola","400","a_Bank","Nothing"]}' -v 1.0
sleep 3
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C pointchannel -n pointcc -c '{"Args":["createAccount","User1@org3.example.com","Carolb","400","b_Bank","Nothing"]}' -v 1.0
sleep 3

peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C pointchannel -n pointcc -c '{"Args":["extrade","Boba","Carola","7","Bobb","Carolb","6"]}' -v 1.0

#################################################################################################################################################################
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["createUser","Admin"]}' -v 1.0
sleep 3
##于reg channel上创建管理员身份
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["createAccount","Admin","Admin_a_Bank","pointchannel","points","a_Bank"]}' -v 1.0
sleep 3
##于reg channel上注册管理员在point channel上的中间人账户
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C pointchannel -n pointcc -c '{"Args":["createAccount","Admin","Admin_a_Bank","0","a_Bank","Nothing"]}' -v 1.0
sleep 3
##于point channel上创建管理员的中间人账户
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C pointchannel -n pointcc -c '{"Args":["crosstrade","Carola","Admin_a_Bank","200","point2"]}' -v 1.0
sleep 3
##User1@org3.example.com将自己名下的a_Bank的积分转移至中间人账户，发起跨链申请。
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["createAccount","User1@org3.example.com","Carola_from_another","pointchannel","points","a_Bank"]}' -v 1.0
sleep 3
##管理员在regchaneel上给用户注册一个新账户。
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C pointchannel -n pointcc -c '{"Args":["createAccount","User1@org3.example.com","Carola_from_another","200","point_a_Bank","Nothing","point channel","a_Bank"]}' -v 1.0
sleep 3
##跨链申请完成后由管理员在另一条channel上创建新的账户，这个账户的交易逻辑和积分链上的机制相同，要求Issuer相同才能trade，不同的积分交换需要调用extrade。

##################################################################################################################################################################

peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["createAccountCL","User1@org3.example.com","Carol_Credit_0","Credit","zyp","Nothing","30000000","zyp"]}' -v 1.0
sleep 3
##创建祖父信用担保
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["splitAccountCL","Carol_Credit_0","User1@org1.example.com","Alice_Credit_1","20000000"]}' -v 1.0
sleep 3
##将一部分担保额度转让给User1@org1.example.com的账户Alice_Credit_1
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["queryAccount","Carol_Credit_0","CL"]}' -v 1.0
sleep 3
##查询祖父担保的当前状况
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C mapchannel -n mapcc -c '{"Args":["queryAccount","Alice_Credit_1","CL"]}' -v 1.0
sleep 3
##查询子担保的当前状况

echo "---------------------------------Test over---------------------------------"
```

