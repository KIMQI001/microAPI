### 服务器Post容量测试
    1.git fc 0.5.8版 并编译
（参考：https://github.com/KIMQI001/FC/blob/master/Ops/scripts/deploy.sh）
    
    2.cd go-filecoin/tools/fast/bin/localnet
    ./localnet -small-sectors=false -miner-count=2 -blocktime=30s > localnet.log 2>&1 &
    等待 localnet.log 最后出现ctrl-c 字样，通知我

