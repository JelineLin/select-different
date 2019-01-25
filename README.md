# selectDiff
select a different file from 100000000 files

//按照1~1亿数字命名文件

1k文件 --- 1级文件
第一次hash换算后的文件 3M ---- 2级文件

1.先将A、B设备的1亿个1K文件的内容通过MD5 hash算法 换算成 32Byte的字符串 存于文件中
  1K / 32B = 32倍
  10^9 / 32 = 3125000K = 3052M = 3G

2.将这些字符串按存于1000个文件中，一个文件约3M（一个2级文件内存放3125个1级文件的内容hash字符串）

3.再2级文件的内容（3M）换成成一个32B的字符串

4.设备B将这10000个 32B的字符串 存放在内存中 （共计 3M）  ---- 使用golang的map数据类型

4.设备A发送一个32B字符串至设备B
     设备B根据字符串在map中做对比 找出不同的字符串 --- 找出内容不同的1级文件所在范围（100000个1级文件）
     if false ==》 设备A继续发送下一个字符串
     if true == 》 设备A 将字符串对应 2级文件内容 加载到内存
                    设备 将字符串对应 2级文件内容 加载到内存
                      重复步骤4；直至找到 1级文件对应的 32B字符串

5.找到不同的1k小文件

MD5值需要计算 一亿一千次
最大比较次数 20000次

=====
或者分为四级文件  每个文件按照100个1级文件 累积hash换算

1亿个MD5值存到1万个文件A中   每个A中存有10000个MD5值对应具体的1k文件
1万个A文件存到100个文件B中  每个B中存有100个MD5值对应具体的A文件
100个B文件存到10个文件C中   每个C中存有10个MD5值对应具体的B文件
10个C文件存到1个文件D中   唯一的文件D中存有10个MD5值对应具体的C文件

MD5值需要计算 一亿零一万零一百一十次
对比次数最大为 10120次

空间换时间

如果按照两个文件 生成一个高一级文件 2的26次 至 2的27次方 每次查找为27次
但是需要的存储空间过大 而且cpu计算MD5值开销也不小

============
所以按照 四级划分 每100个文件生成一个高级文件 的方式最佳

通信协议：

[0-4)：prefix = 0xAAAA
[4]:reqType
[5-7):dataLen
[7:7+dataLen):data

reqType:
1: reqString //请求字符串
2: rspBool //返回对比结果
3: reqStringKey //对应的k值
