cd ~/workspace/go/src/bilibili/bin
./danmu_sender &
~/scripts/dwm-status.sh
./danmu_geter | while read line
do
    notify-send "$line"
done
