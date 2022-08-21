cd "$(dirname "${BASH_SOURCE[0]}")"

sender_pid=`ps -aux | grep danmu_sender | grep -v grep | awk '{print $2}'`
geter_pid=`ps -aux | grep danmu_geter | grep -v grep | awk '{print $2}'`

if [ "$sender_pid$geter_pid" ]; then
    [ "$sender_pid" ] && kill -9 $sender_pid
    [ "$geter_pid" ] && kill -9 $geter_pid
    notify-send "弹幕服务器已停止"
    ~/scripts/dwm-status.sh
else
    bash monitor.sh 2&>1 >/dev/null &
fi
