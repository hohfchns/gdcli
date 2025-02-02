[Setup]
AppName=Godot Project CLI
AppVersion=1.0
DefaultDirName={commonpf}\Godot Project CLI
DefaultGroupName=Godot Project CLI
OutputDir=.
OutputBaseFilename=gdcliSetup
SetupIconFile=icon.ico

[Files]
Source: "bin/gdcli.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "icon.ico"; DestDir: "{app}"; Flags: ignoreversion

[Registry]
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; \
ValueType: expandsz; ValueName: "Path"; \
ValueData: "{olddata};{app}"

[Icons]
Name: "{group}\Godot Project CLI"; Filename: "{app}\gdcli.exe"; IconFilename: "{app}\icon.ico"
