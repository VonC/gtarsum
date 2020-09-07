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
REM https://stackoverflow.com/questions/16354102/how-to-get-the-rest-of-arguments-in-windows-batch-file
REM https://stackoverflow.com/a/16354963/6309
SET allargs=%*
SET arg1=%1
CALL SET someargs=%%allargs:*%1=%%

%dirname% "%f%" %someargs%
set err=%ERRORLEVEL%
echo.
echo ERRORLEVEL (exit status)='%err%'
echo.
dir /a-d "*.hash*"|findstr hash
