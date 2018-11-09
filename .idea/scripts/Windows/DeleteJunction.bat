@echo off
if exist "%USERPROFILE%\go\src\github.com\fang2hou\easyga" (
    rd "%USERPROFILE%\go\src\github.com\fang2hou\easyga"
    echo "Deleted."
)
pause