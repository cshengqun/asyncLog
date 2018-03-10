async log component
=====================
this log component is based on the original golang package "log", and some more features was added, including:
1. asynchronous log
2. support log level
3. support rolling log file 

usage
=====================
create a logger
--------------------
>logger := asyncLog.NewLogger(fileName, logLevel, logBufInCnt, logTreadCnt)

setup logger, for example, max file cnt and the max log file size
-------------------
>logger.SetLogCnt(maxFileCnt)
>logger.SetFileSize(logFileMaxSize)


