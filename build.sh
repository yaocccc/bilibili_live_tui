_dir=$(pwd)

cd $_dir/danmu_geter
go build -o $_dir/bin/danmu_geter main.go

cd $_dir/danmu_sender
go build -o $_dir/bin/danmu_sender main.go

cp $_dir/config.toml $_dir/bin/config.toml
