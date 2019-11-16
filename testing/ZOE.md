## 任务目标
1.测试 go_sectorbuilder.NewSortedSectorInfo()函数功能
    
    普通概念介绍：
    1. 每一次调用 addPiece() 会返回一个SectorID
    2. 然后会程序自动检测，如果满扇区会自动启动 PoRep（复制证明）
    3. PoRep 结果中的 commr 是启动 PoSt（时空证明）的req(启动参数)
    4. 即每一个 SectorID 对应一个commr
    
    测试函数概念介绍：
    1.NewSortedSectorInfo()，会将输入的 commr... 进行排序
    2.测试此函数排序出来的结果 是否也是排序 SectorID 的结果
    例： secID 1 —— commr 1
        secID 2 —— commr 2
        secID 3 —— commr 3
        secID 4 —— commr 4
    
    看是否函数输出为 commr 1,2,3,4 或者 commr 4,3,2,1
    
    测试目标流程介绍：
    1.将testFuctions包放入go-filecoin
    2.找到TestAddpiece.go中的：
           // TODO：ZOE：： 的两个标记，标记有使用说明
    
