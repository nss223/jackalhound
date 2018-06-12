# multiAsset
一共有三个链码，`regcc`、`mapcc`、`pointcc`。分别运行在三条channel上，这三条channel对应为`regchannel`、`mapchannel`、`pointchannel`。

配置好的环境见e2e_cli文件夹下的配置文件。三个组织，四条链，每个组织两个节点，若干个用户。所有链码的权限设置都是member级别。不过在链码实现中加入了权限管理。

以下所有链码的调用方式都是invoke，通过传入的第一个参数决定函数名，第二个通常为用户名，第三个通常为账户名，其余为参数。调用命令类似如下格式（本链码在e2e_cli下开发，环境可能不同，请按照环境调整参数）
```bash
peer chaincode invoke -o orderer.example.com:7050  --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C regchannel -n regcc -c '{"Args":["queryAccount","User1@org1.example.com","all","all"]}' -v 1.0
```

在调用完一条非查询类命令后，最好`sleep`几秒等待数据写入账本，否则可能会报错。

通过解析当前用户证书可以得到当前用户名，借助这个信息可以实现一部分细化的权限管理。本例在某些函数中实现了这一方法达成权限管理的目的。
这一块也可以替换成为验证公钥签名算法，不过过于麻烦，并没有在本例中实现。

### regchannel

regcc负责管理记录用户注册的账户名，Userid相当于身份证号，Userid对应的值为User结构，User结构下维护一个用户当前所拥有的账户列表Accounts。Accounts
为键值对的一个map，键为Accountid，值是一个Account结构的数据，用于存储该账户的一些基本信息。包括：
1. ChannelID：该账户注册在哪条channel上
2. AccountType：该账户的类型，决定了调用账户的方法和逻辑，不过在本例中没有用到
3. Issuer：该项虚拟资产的发行方

注意，这里因为跨链调用能力限制，对于虚拟映射类资产，我们每一条channel在这里注册一个账户，AccountId直接就是该条channel的ID。对于账户类资产，一个
channel上可以有多个账户在User名下。
regcc被视为账户管理最高权限机构，因此，凡是想要在其他链上注册新账户的用户都必须先向regcc提交申请，格式为：createAccout（Userid，Accoutid，Account属性）
当用户在regcc之外的channel上调用链码进行注册时，如果没有实现向regcc备案注册，则无法注册成功。
regcc下有两种对象，User和Account，每种对象有三种通用调用方法：create、query、delete

createUser：生成一个新的用户，并给其一个空账户列表,这里建议Userid使用当前用户的domain，形如：User1@org2.example.com, 以免在例子接下来的演示中造成bug。
应当传入参数: `{"Args":["createUser","Userid"]}`

queryUser：查询用户名下所有账户，只返回Accountid，详细信息请使用queryAccount。
应当传入参数：`{"Args":["queryUser","Userid"]}`

deleteUser:删除账户，该操作将在regcc上删除Userid对应的这一键值对。包括其名下的所有账户列表。
应当传入参数: `{"Args":["deleteUser","Userid"]}`

createAccout：在指定用户名下添加一个新的账户，Userid未注册时将报错。因查询函数使用了“all”字段，故账户名不能为“all”，否则无法注册。
应当传入参数：`{"Args":["createAccount","Userid","Accountid","ChannelID","AccountType","Issur"]}`

queryAccount: 查询指定账户的基本信息，需要指定该账户所属用户名。Key可选参数为“all”（返回全部信息）“ChannelID”“AccountType”“Issuer”
应当传入参数：`{"Args":["queryAccount","Userid","Accountid","Key"]}`

deleteAccount: 删除指定账户，需指定在那个用户名下。
应当传入参数：`{"Args":["deleteAccount","Userid","Accountid"]}`

setAssetByAccount： 修改指定账户信息，相当于重新建一个名字一样的账户。
应当传入参数：`{"Args":["serAssetByAccount","Userid","Accountid","ChannelID","AccountType","Issur"]}`

queryHistory: 查询制定用户的历史交易记录，返回结果需要另外解码，没有集成到链码中。
应当传入参数：`{"Args":["queryHistory","Userid"]}`

queryall：对用户名下所有资产进行归集，返回归集信息，没有集成处理信息的模块，可以另行实现
应当传入参数：`{"Args":["queryall","Userid"]}`

### mapchannel

mapcc负责管理记录虚拟化资产，因为链码调用逻辑问题，mapcc下对每个Userid再另外维护一个User结构，包含该channel内用户所拥有的所有虚拟化资产的账户id。
每个键Accountid对应值为Account结构，一个Accountid对应一项虚拟化资产。比如说你有两栋房子并且虚拟化上链了，那么你在mapcc下的账户列表中就拥有两个房
产属性的Accountid。每个Account结构下包含了如下信息：
1. Type：该项虚拟资产的属性
2. Owner：该项虚拟资产的所有人
3. Issuer：该项虚拟资产的发行方
4. Other：备注字段

在mapcc上注册用户Userid时，会检查用户是否在regcc上注册，如果没有提前注册申请，那么无法在mapcc上注册。
mapcc下也有两种对象，User和Account，User结构下的Accounts是用户在mapcc上所拥有资产的Accountid的一个列表（以map方式实现便于查询），并且只记录id，
而资产的id另行作为mapcc上的一个键值对来存储资产的信息。每个Accountid作为键对应一个值Account。Account结构上面已经提到。
这两种对象也同样有三种通用的方法：create、query、delete

createUser：生成一个新的用户，并给其一个空账户列表，为了区分不同的键，这里的userid应当以大写字母“U”开头，否则不合法
应当传入参数: `{"Args":["createUser","Userid"]}`

queryUser：查询用户名下所有账户，只返回Accountid，详细信息请使用queryAccount。
应当传入参数：`{"Args":["queryUser","Userid"]}`

deleteUser:删除账户，该操作将在regcc上删除Userid对应的这一键值对。包括其名下的所有账户对应的键值对也会被删除。
应当传入参数：`{"Args":["deleteUser","Userid"]}`

createAccout：在指定用户名下添加一个新的账户，Userid未注册时将报错。不用另行输入Owner字段。需要以小写字母"a"开头
应当传入参数：`{"Args":["createAccount","Userid","Accountid","Type","Issuer","Other"]}`

queryAccount: 查询指定账户的基本信息，需要指定该账户所属用户名。返回所有信息
应当传入参数：`{"Args":["queryAccount","Userid","Accountid"]}`

deleteAccount: 删除指定账户，不需指定在那个用户名下，因此也无法删除用户名下列表中的账户id，这属于未修复的bug，对accountdelete函数简单修复即可。
应当传入参数：`{"Args":["deleteAccount","Accountid"]}`

queryHistory: 查询指定键的历史交易记录，返回结果需要另外解码，没有集成到链码中。键可以是Userid，也可以是Accountid。
应当传入参数：`{"Args":["queryHistory","Userid"/"Accountid"]}`

trade: 实现了用户A将自己名下的虚拟财产转移至B的名下。需要B先在mapcc上成功注册一个账户。
应当传入参数：`{"Args":["trade","Aid","Bid","assetid"]}`

### pointchannel

pointcc负责管理用户的账户类资产，也就是所谓的积分。pointcc相对于其他两条链更为简单，只有一种对象Account，Account作为数据结构包含如下信息：
1. Balance：当前账户积分余额，限定为非负整数
2. Owner：该项资产的所有人
3. Issuer：该项资产的发行方
4. Other：备注字段

在pointcc上注册账户时，会检查该账户是否已在regcc上申请注册到用户名下过，否则无法注册。
进行转账交易时，需要时同一发行方发行的积分才能相互交易，如果想要向某人交易，需要确认他拥有同一发行方的pointcc上的账户。
Account对象也有三种通用方法：create、query、delete。此外还有一项trade方法，这里的trade方法与mapcc略有不同。

createAccout：添加一个新的账户，Userid未注册时将报错。如果不是管理员添加，那么Balance会自动调整为0.
应当传入参数：`{"Args":["createAccount","Userid","Accountid","Balance","Issuer","Other"]}`

queryAccount: 查询指定账户的基本信息，Key可选参数为“all”（返回全部信息）“Balance”“Other”“Issuer”
应当传入参数：`{"Args":["queryAccount","Accountid","Key"]}`

deleteAccount:删除账户，该操作将在pointcc上删除对应的键值对。
应当传入参数：`{"Args":["deleteAccount","Accountid"]}`

setAccount: 修改账户信息，这个函数只有管理员才有权力使用。传入的参数与createAccount一致
应当传入参数：`{"Args":["setAccount","Userid","Accountid","Balance","Issuer","Other"]}`

queryHistory：查询账户历史交易信息
应当传入参数：`{"Args":["queryHistory","Accountid"]}`

trade：实现了账户A到账户B的积分交易，两个账户应当是同一发行方账户。A向B转X积分，如果当前用户不是管理员或账户A所有者，那么交易失败。
应当传入参数：`{"Args":["Aid","Bid","X"]}`

### pointchannel

如果两个channel都是按照pointcc启用的话，那么这两条channel可以实现以下跨channel逻辑

在pointchannel上实现了一个跨channel转移账户类积分的函数crosstrade，可以在两条channel间转移积分类资产。这个过程中，用户调用crosstrade函数进行跨channel申请，
接下来再由管理员账户进行操作，在指定channel上生成对应的账户类资产，并将Owner设为指定用户。之后用户可以自行在对应链上正常交易积分。

extrade：同一channel上不同类型积分资产交易互换，设计方案中此函数的实现应当是收集交易双方签名后才可以调用，在实现中，为简便起见，这个命令目前只能由管理员调用。
收集签名才能使用在逻辑上是可行的，但工作量过大，暂未实现签名模块。
应当传入参数：`{"Args":["extrade","User1's Account of a_bank","User2's Account of a_bank","Xa","User1's Account of b_bank","User2's Account of b_bank","Xb"]}`
一共有六个参数，前三个参数描述用户User1的a银行应当向User2的a银行账户转Xa积分，后三个参数描述用户User2的b银行账户应当向User1的b银行账户转Xb积分。

crosstrade: 这个函数实现了用户向管理员申请跨链的信息发送，实质上是将用户的积分资产转移至管理员指定账户并锁定。之后由管理员在指定链上调用creatAccount进行操作。
应当传入参数: `{"Args":["crosstrade","UserAccount","AdminAccount","Value","Aimchannel"]}`

管理员监听区块链系统，在收到crosstrade交易在Aimchannel上生成Owner为UserAccount.Owner的账户，管理员需要先在regchannel上注册要生成的账户名，之后如下调用createAccout：
应当传入参数: `{"Args":["createAccount","UserAccount.Owner","newAccountID","UserAccount.Balance","newIssuerID","UserAccount.Other","From","OriIssuer"]}`
其中newAccountID由管理员事先申请注册，newIssuerID是旧有Issuer在新链上的标识，应当不与新链其他Issuer的ID相同，并且保证同一OriIssuer的newIssuerID应当一致。From代表该积分原本属于哪一个channel

### mapchannel

在mapchannel上实现了一个新的资产类型，名为AccountCL，意为Account Credit Lines，该资产拥有Balance，可以对Balance进行拆分，拆分后的资产可以查询到是由哪个父资产拆分而来，
每个资产可以查询到其拆分的一级子资产。除去其拆分特性，其余特性继承虚拟类资产的一切方法。
```go
type AccountCL {
	AccountMV	//实现继承特性的Account minimum version
	Balance		//信用额度
	Parent		//信用担保来源
	Children	//信用担保去向
}
```

这个资产类型可以调用两个函数进行操作：第一个是createAccountCL，第二个是splitAccount。

createAccountCL：创建新的信用担保类型资产，由银行发行到用户UserID的账户AssetID，父资产ID指向发行者，初始化子资产列表为空集
应当传入参数: `{"Args":["createAccountCL","UserID","AssetID","Type","Issuer","Other","Balance","Parent"]}`

splitAccountCL：拆分自己拥有的信用担保并交易给他人，子信用担保资产的父资产ID指向自己，将子资产ID加入Children列表。
应当传入参数: `{"Args":["splitAccountCL","AssetID","ChildrenUserID","ChildrenAssetID","Value"]}`
Value代表你要转移多少信用额度给予ChildrenUser，这些信用额度将保存在ChildrenAssetID账户下。



# DataFlow
数据流转平台

### datachannel
通道 datachannel 上的链码 data 实现了数据流转平台。
```go
type data struct {
	//ObjectType string `json:"docType"`        // 为了和状态数据库中的其他key区分，留作备用
	ID         string   `json:"id"`             // 数据 id
	URI        string   `json:"uri"`            // 数据访问地址，uri 格式
	Key        string   `json:"key"`            // 对称加密密钥
	ClearHash  string   `json:"clear_hash"`     // 明文哈希
	CipherHash string   `json:"cipher_hash"`    // 密文哈希
	Doc        string   `json:"doc"`            // 数据说明
	Creater    string   `json:"creater"`        // 创建者
	Owner      []string `json:"owner"`          // 拥有者
	Pid        string   `json:"pid"`            // 父数据 id，若为空字符串说明是根节点数据
	Timestamp  string   `json:"timestamp"`      // 创建时间
}
```

* `func dataIsValid(uri, key, clearhash, cipherhash string) bool`
检查数据是否合法
    * 需要读取数据进行验证。未实现。

### 链码函数
* `commit(id, uri, key, clearhash, cipherhash, doc)`
创建数据
    * 数据的创建者和拥有者为调用此函数的用户，父节点为空。
    + 写操作

* `share(id, user)`
分享数据
    * 将 `user` 添加到数据 `id` 的拥有者列表。
    - 如果调用者不在数据 `id` 的拥有者列表里，或者调用者不是管理员，将返回权限错误。
    + 写操作

* `branch(pid, cid, uri, key, clearhash, cipherhash, doc)`
创建数据分支
    * 原数据 `pid` 作为新数据 `cid` 的父数据，新数据 `cid` 由后边的参数生成，其创建者和拥有者为调用此函数的用户。
    - 如果调用者不在数据 `pid` 的拥有者列表里，或者调用者不是管理员，将返回权限错误。
    + 写操作

* `checkout(id)`
读数据
    - 如果调用者不在数据 `id` 的拥有者列表里，或者调用者不是管理员，将返回权限错误。
    + 只读操作，返回数据内容

* `trace(id)`
获取数据流转信息
    * 递归查找数据父节点，直到某个根节点。
    - 如果调用者不在数据 `id` 的拥有者列表里，或者调用者不是管理员，将返回权限错误。
    + 只读操作，返回从数据 `id` 到根节点之间（包含）所有数据构成的列表

* `queryByOwner()`
按所有者归集
    + 只读操作，返回调用者可读的所有数据构成的列表
    + 富查询，需要 CouchDB

* `queryByCreater()`
按创建者归集
    + 只读操作，返回调用者创建的所有数据构成的列表
    + 富查询，需要 CouchDB

* `history(id)`
查询历史
    * 特别地，只有拥有者字段会更改。
    - 如果调用者不在数据 `id` 的拥有者列表里，或者调用者不是管理员，将返回权限错误。
    + 只读操作，返回数据变更历史
