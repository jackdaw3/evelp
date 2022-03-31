pid=$(ps -ef | grep serenity | grep -v grep | awk '{print $2}')
if [ -n "$pid" ]; then
        kill -9 $pid;
fi

if [ $? -eq 0 ];then
        echo "stopped"
else
        echo "failed to stop"
fi