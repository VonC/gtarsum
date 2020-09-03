@echo off
setlocal enabledelayedexpansion

for %%i in ("%~dp0.") do SET "script_dir=%%~fi"
cd "%script_dir%"
for %%i in ("%~dp0.") do SET "dirname=%%~ni"

set "f=%1"

if "%f%" == "" (
    set f=ex.tar
)

%dirname% "%f%"
