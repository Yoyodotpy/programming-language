$Repo = "Yoyodotpy/programming-language"

New-Item -ItemType Directory -Force -Path ".\gfpl" | Out-Null
Set-Location -Path ".\gfpl"

$Url = "https://github.com/$Repo/releases/latest/download/gfpl.exe"
Write-Host "downloading gfpl from $Url..."
Invoke-WebRequest -Uri $Url -OutFile "gfpl.exe"

Write-Host "downloading examples scripts..."
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/$Repo/refs/heads/main/examples/fibonacci.gfpl" -OutFile "fibonacci.gfpl"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/$Repo/refs/heads/main/examples/hello_world_complicated.gfpl" -OutFile "hello_world_complicated.gfpl"

Write-Host "successfully downloaded gfpl and example scripts."
Write-Host ""
Write-Host "to test gfpl, you can run:"
Write-Host ".\gfpl.exe fibonacci.lamb"
Write-Host ""
Write-Host "please look at the git repo if you need help: "
Write-Host "https://www.github.com/$Repo"
