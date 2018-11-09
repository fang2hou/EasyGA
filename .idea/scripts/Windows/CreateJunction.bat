@echo off
cd /d %~dp0..\..\..\..\
if exist "%USERPROFILE%\go\src\github.com\fang2hou\easyga" (
    echo Folder is exist.
    echo You can use another bat to delete.
) else (
    if not exist "%USERPROFILE%\go\src\github.com\fang2hou\" (
        mkdir "%USERPROFILE%\go\src\github.com\fang2hou\"
    )
    mklink /j "%USERPROFILE%\go\src\github.com\fang2hou\easyga" EasyGA
    echo Done.
)
pause