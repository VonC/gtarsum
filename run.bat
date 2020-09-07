@echo off
setlocal enabledelayedexpansion

for %%i in ("%~dp0.") do SET "script_dir=%%~fi"
cd "%script_dir%"
for %%i in ("%~dp0.") do SET "dirname=%%~ni"

del *.hash*

set "f=%1"

if "%f%" == "" (
    set f=ex.tar
)
shift

%dirname% "%f%" %someargs%
set err=%ERRORLEVEL%
echo.
echo ERRORLEVEL (exit status)='%err%'
echo.
dir /a-d "*.hash*"|findstr hash
