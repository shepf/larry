# 使用说明
export GITHUB_ACCESS_TOKEN=your_github_access_token

go build -o mylarry  cmd/larry/main.go
./mylarry  --pub=csdn --lang=go --time=720 --cron="0 0 7 * * ?"

nohup  ./mylarry  --pub=csdn --lang=go --time=720 --cron="0 0 7 * * ?" > ~/nohup.mylarry.output 2>&1 &
tail -f ~/nohup.mylarry.output