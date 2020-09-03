@echo off
setlocal enabledelayedexpansion
call build.bat
if not "%ERRORLEVEL%"=="0" ( goto:eof )
call run.bat %*
if not "%ERRORLEVEL%"=="0" ( goto:eof )
