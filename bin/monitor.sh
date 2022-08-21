cd "$(dirname "${BASH_SOURCE[0]}")"
./danmu_sender &
./danmu_geter | while read line
do
    echo "$line" >> ~/bililog
    notify-send "$line"
done
