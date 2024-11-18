#!/bin/bash

# 指定日志文件所在的目录
log_directory="/data/will/bntradestat"

ls -t $log_directory/bntradestat* | tail -n +5 | xargs rm -f

