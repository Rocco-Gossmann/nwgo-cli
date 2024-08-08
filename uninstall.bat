@echo off
FOR /F "tokens=*" %%F IN ('"go env GOPATH"') DO SET gopath=%%F

IF %gopath% == "" GOTO no_gopath

@copy %gopath%\bin\nwgo %gopath%\bin\nwgo.exe

%gopath%\bin\nwgo.exe uninstall
del %gopath%\bin\nwgo.exe

GOTO end

:no_gopath
echo "kein GOPATH"

GOTO end

:end
@echo on
