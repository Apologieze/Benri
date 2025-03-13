set RESHACKER_PATH="icon/ResourceHacker.exe"
set INPUT_EXE=".\src\Benri.exe"
set OUTPUT_EXE="Benri.exe"
set ICON_FILE="icon/IconW.ico"

echo Replacing icon...
%RESHACKER_PATH% -open %INPUT_EXE% -save %OUTPUT_EXE% -action delete -mask ICONGROUP,
%RESHACKER_PATH% -open %OUTPUT_EXE% -save %OUTPUT_EXE% -action addoverwrite -res %ICON_FILE% -mask ICONGROUP,MAINICON,

del %INPUT_EXE%
echo Done.