# Jackalhound

修改了mapchaincode下账户名和用户名首字母的问题，现在对首字母无要求。但是相应的delete函数拆分成了deleteUser和deleteAccount，queryall和queryHisory也做了相应拆分。

在pointchaincode中增加了新的函数extrade，该函数的作用是多资产间的相互兑换，兑换额度和比例由用户自行决定。此函数目前只能由admin权限调用，如果需要可以进一步权限的逻辑。对于demo版本来说，要实现完整的用户权限控制过于复杂，故暂不实现更完善的权限管理机制。
 Aa trans Xa token to Ba, Bb trans Xb token to Ab.
 该函数需要六个参数，参数顺序是Aa，Ba，Xa，Ab，Bb，Xb
 
