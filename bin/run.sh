echo 111 >> ~/log
ps -aux | grep danmu_geter | grep -v grep | awk '{print $2}' >> ~/log
ps -aux | grep danmu_sender | grep -v grep | awk '{print $2}' >> ~/log

sender_pid=`ps -aux | grep danmu_sender | grep -v grep | awk '{print $2}'`
geter_pid=`ps -aux | grep danmu_geter | grep -v grep | awk '{print $2}'`

echo "$sender_pid$geter_pid" >> ~/log

if [ "$sender_pid$geter_pid" ]; then
    [ "$sender_pid" ] && kill -9 $sender_pid
    [ "$geter_pid" ] && kill -9 $geter_pid
    notify-send "弹幕服务器已停止"
    ~/scripts/dwm-status.sh
else
    bash ~/workspace/go/src/bilibili/bin/monitor.sh 2&>1 >/dev/null &
fi
