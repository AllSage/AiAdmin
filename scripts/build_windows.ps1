#!powershell
#
# powershell -ExecutionPolicy Bypass -File .\scripts\build_windows.ps1
#
# gcloud auth application-default login

$ErrorActionPreference = "Stop"

function checkEnv() {
    $script:TARGET_ARCH=$Env:PROCESSOR_ARCHITECTURE.ToLower()
    Write-host "Building for ${script:TARGET_ARCH}"
    write-host "Locating required tools and paths"
    $script:SRC_DIR=$PWD
    if (!$env:VCToolsRedistDir) {
        $MSVC_INSTALL=(Get-CimInstance MSFT_VSInstance -Namespace root/cimv2/vs)[0].InstallLocation
        $env:VCToolsRedistDir=(get-item "${MSVC_INSTALL}\VC\Redist\MSVC\*")[0]
    }
    # Try to find the CUDA dir
    if ($null -eq $env:NVIDIA_DIR) {
        $d=(get-command -ea 'silentlycontinue' nvcc).path
        if ($d -ne $null) {
            $script:NVIDIA_DIR=($d| split-path -parent)
        } else {
            $cudaList=(get-item "C:\Program Files\NVIDIA GPU Computing Toolkit\CUDA\v*\bin\" -ea 'silentlycontinue')
            if ($cudaList.length > 0) {
                $script:NVIDIA_DIR=$cudaList[0]
            }
        }
    } else {
        $script:NVIDIA_DIR=$env:NVIDIA_DIR
    }
    
    $script:INNO_SETUP_DIR=(get-item "C:\Program Files*\Inno Setup*\")[0]

    $script:DEPS_DIR="${script:SRC_DIR}\dist\windows-${script:TARGET_ARCH}"
    $env:CGO_ENABLED="1"
    echo "Checking version"
    if (!$env:VERSION) {
        $data=(git describe --tags --first-parent --abbrev=7 --long --dirty --always)
        $pattern="v(.+)"
        if ($data -match $pattern) {
            $script:VERSION=$matches[1]
        }
    } else {
        $script:VERSION=$env:VERSION
    }
    $pattern = "(\d+[.]\d+[.]\d+).*"
    if ($script:VERSION -match $pattern) {
        $script:PKG_VERSION=$matches[1]
    } else {
        $script:PKG_VERSION="0.0.0"
    }
    write-host "Building AiAdmin $script:VERSION with package version $script:PKG_VERSION"

    # Note: Windows Kits 10 signtool crashes with GCP's plugin
    if ($null -eq $env:SIGN_TOOL) {
        ${script:SignTool}="C:\Program Files (x86)\Windows Kits\8.1\bin\x64\signtool.exe"
    } else {
        ${script:SignTool}=${env:SIGN_TOOL}
    }
    if ("${env:KEY_CONTAINER}") {
        ${script:AiAdmin_CERT}=$(resolve-path "${script:SRC_DIR}\AiAdmin_inc.crt")
        Write-host "Code signing enabled"
    } else {
        write-host "Code signing disabled - please set KEY_CONTAINERS to sign and copy AiAdmin_inc.crt to the top of the source tree"
    }

}


function buildAiAdmin() {
    write-host "Building AiAdmin CLI"
    if ($null -eq ${env:AiAdmin_SKIP_GENERATE}) {
        & go generate ./...
        if ($LASTEXITCODE -ne 0) { exit($LASTEXITCODE)}    
    } else {
        write-host "Skipping generate step with AiAdmin_SKIP_GENERATE set"
    }
    & go build -trimpath -ldflags "-s -w -X=github.com/AllSage/AiAdmin/version.Version=$script:VERSION -X=github.com/AllSage/AiAdmin/server.mode=release" .
    if ($LASTEXITCODE -ne 0) { exit($LASTEXITCODE)}
    if ("${env:KEY_CONTAINER}") {
        & "${script:SignTool}" sign /v /fd sha256 /t http://timestamp.digicert.com /f "${script:AiAdmin_CERT}" `
            /csp "Google Cloud KMS Provider" /kc ${env:KEY_CONTAINER} AiAdmin.exe
        if ($LASTEXITCODE -ne 0) { exit($LASTEXITCODE)}
    }
    New-Item -ItemType Directory -Path .\dist\windows-${script:TARGET_ARCH}\ -Force
    cp .\AiAdmin.exe .\dist\windows-${script:TARGET_ARCH}\
}

function buildApp() {
    write-host "Building AiAdmin App"
    cd "${script:SRC_DIR}\app"
    & windres -l 0 -o AiAdmin.syso AiAdmin.rc
    & go build -trimpath -ldflags "-s -w -H windowsgui -X=github.com/AllSage/AiAdmin/version.Version=$script:VERSION -X=github.com/AllSage/AiAdmin/server.mode=release" .
    if ($LASTEXITCODE -ne 0) { exit($LASTEXITCODE)}
    if ("${env:KEY_CONTAINER}") {
        & "${script:SignTool}" sign /v /fd sha256 /t http://timestamp.digicert.com /f "${script:AiAdmin_CERT}" `
            /csp "Google Cloud KMS Provider" /kc ${env:KEY_CONTAINER} app.exe
        if ($LASTEXITCODE -ne 0) { exit($LASTEXITCODE)}
    }
}

function gatherDependencies() {
    write-host "Gathering runtime dependencies"
    cd "${script:SRC_DIR}"
    md "${script:DEPS_DIR}\AiAdmin_runners" -ea 0 > $null

    # TODO - this varies based on host build system and MSVC version - drive from dumpbin output
    # currently works for Win11 + MSVC 2019 + Cuda V11
    cp "${env:VCToolsRedistDir}\x64\Microsoft.VC*.CRT\msvcp140*.dll" "${script:DEPS_DIR}\AiAdmin_runners\"
    cp "${env:VCToolsRedistDir}\x64\Microsoft.VC*.CRT\vcruntime140.dll" "${script:DEPS_DIR}\AiAdmin_runners\"
    cp "${env:VCToolsRedistDir}\x64\Microsoft.VC*.CRT\vcruntime140_1.dll" "${script:DEPS_DIR}\AiAdmin_runners\"
    foreach ($part in $("runtime", "stdio", "filesystem", "math", "convert", "heap", "string", "time", "locale", "environment")) {
        cp "$env:VCToolsRedistDir\..\..\..\Tools\Llvm\x64\bin\api-ms-win-crt-${part}*.dll" "${script:DEPS_DIR}\AiAdmin_runners\"
    }


    cp "${script:SRC_DIR}\app\AiAdmin_welcome.ps1" "${script:SRC_DIR}\dist\"
    if ("${env:KEY_CONTAINER}") {
        write-host "about to sign"
        foreach ($file in (get-childitem "${script:DEPS_DIR}\cuda\cu*.dll") + @("${script:SRC_DIR}\dist\AiAdmin_welcome.ps1")){
            write-host "signing $file"
            & "${script:SignTool}" sign /v /fd sha256 /t http://timestamp.digicert.com /f "${script:AiAdmin_CERT}" `
                /csp "Google Cloud KMS Provider" /kc ${env:KEY_CONTAINER} $file
            if ($LASTEXITCODE -ne 0) { exit($LASTEXITCODE)}
        }
    }
}

function buildInstaller() {
    write-host "Building AiAdmin Installer"
    cd "${script:SRC_DIR}\app"
    $env:PKG_VERSION=$script:PKG_VERSION
    if ("${env:KEY_CONTAINER}") {
        & "${script:INNO_SETUP_DIR}\ISCC.exe" /DARCH=$script:TARGET_ARCH /SMySignTool="${script:SignTool} sign /fd sha256 /t http://timestamp.digicert.com /f ${script:AiAdmin_CERT} /csp `$qGoogle Cloud KMS Provider`$q /kc ${env:KEY_CONTAINER} `$f" .\AiAdmin.iss
    } else {
        & "${script:INNO_SETUP_DIR}\ISCC.exe" /DARCH=$script:TARGET_ARCH .\AiAdmin.iss
    }
    if ($LASTEXITCODE -ne 0) { exit($LASTEXITCODE)}
}

function distZip() {
    write-host "Generating stand-alone distribution zip file ${script:SRC_DIR}\dist\AiAdmin-windows-${script:TARGET_ARCH}.zip"
    Compress-Archive -Path "${script:SRC_DIR}\dist\windows-${script:TARGET_ARCH}\*" -DestinationPath "${script:SRC_DIR}\dist\AiAdmin-windows-${script:TARGET_ARCH}.zip" -Force
}

try {
    checkEnv
    buildAiAdmin
    buildApp
    gatherDependencies
    buildInstaller
    distZip
} catch {
    write-host "Build Failed"
    write-host $_
} finally {
    set-location $script:SRC_DIR
    $env:PKG_VERSION=""
}
