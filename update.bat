@echo off

@REM 設定編碼為UTF-8
chcp 65001

cd "%2"
git fetch --all
git checkout -f %3
git pull
DEL "convertable\data\GameSetting.xlsx"
COPY "%1\GameSetting.xlsx" "%2\convertable\data\GameSetting.xlsx"
cd "convertable"
go build . && start /wait convertable.exe
git add -A
git commit -m "[TG]企劃資料更新 by %4"
git push