# 如果有老进程存在的话，杀掉
pid=$(pgrep -f server_proxy)
if [ -n "$pid" ]; then
  kill -9 $pid
  sleep 1
fi

nohup ./server_proxy > nohup.out 2>&1 &