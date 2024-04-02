#实现读取当前目录下的go文件，并统计go代码行数。
#！/bin/bash

#统计go代码行数
function count_go_code(){
    #获取当前目录下的go文件
    go_files=`ls *.go`
    #统计go代码行数
    go_code_count=0
    for file in $go_files
    do
        go_code_count=`expr $go_code_count + $(cat $file | wc -l)`
    done
    echo "go代码行数：$go_code_count"
}

#调用统计go代码行数函数
count_go_code

