async log component
=====================
this log component is based on the original golang package "log", and some more features was added, including:
1. asynchronous log
2. support log level
3. support rolling log file 
4. implement io.Writer interface

usage
=====================
create a logger
--------------------
> logger := asyncLog.NewLogger(fileName, logLevel, logBufInCnt, logTreadCnt)

setup logger, for example, max file cnt, the max log file size and the io.Writer level
-------------------
> logger.SetLogCnt(maxFileCnt)
> logger.SetFileSize(logFileMaxSize)
> logger.SetWriterLv(level)


