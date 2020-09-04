@echo off
setlocal enabledelayedexpansion

for %%i in ("%~dp0.") do SET "script_dir=%%~fi"
cd "%script_dir%"
for %%i in ("%~dp0.") do SET "dirname=%%~ni"

rem https://medium.com/@joshroppo/setting-go-1-5-variables-at-compile-time-for-versioning-5b30a965d33e
for /f %%i in ('git describe --long --tags --dirty --always') do set gitver=%%i
set VERSION=v0.0.1
for /f %%i in ('git config user.name') do set usern=%%i

rem https://stackoverflow.com/a/1445724/6309
: Sets the proper date and time stamp with 24Hr Time for log file naming
: convention

SET HOUR=%time:~0,2%
if "%HOUR:~0,1%" == " " (SET HOUR=0%HOUR:~1,1%) 
rem SET dtStamp9=%date:~-4%%date:~4,2%%date:~7,2%_0%time:~1,1%%time:~3,2%%time:~6,2% 
SET dtStamp24=%date:~-4%%date:~3,2%%date:~0,2%-%HOUR%%time:~3,2%%time:~6,2%

rem if "%HOUR:~0,1%" == " " (SET dtStamp=%dtStamp9%) else (SET dtStamp=%dtStamp24%)
SET dtStamp=%dtStamp24%

if not "%1" == "amd" (
go build -race -ldflags "-X %dirname%/version.GitTag=%gitver% -X %dirname%/version.BuildUser=%usern% -X %dirname%/version.Version=%VERSION% -X %dirname%/version.BuildDate=%dtStamp%" .
) else (
    cmd /V /C "set GOOS=linux&& set GOARCH=amd64&& go build -ldflags ^"-X %dirname%/version.GitTag=%gitver% -X %dirname%/version.BuildUser=%usern% -X %dirname%/version.Version=%VERSION% -X %dirname%/version.BuildDate=%dtStamp%^" ."
    if not "%2" == "" (
        scp %dirname% %2:~/bin/%dirname%
    )
)