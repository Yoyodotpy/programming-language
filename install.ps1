$Repo = "Yoyodotpy/programming-language"

New-Item -ItemType Directory -Force -Path ".\gfpl" | Out-Null
Set-Location -Path ".\gfpl"

$Url = "https://github.com/$Repo/releases/latest/download/gfpl.exe"
Write-Host "downloading gfpl from $Url..."
Invoke-WebRequest -Uri $Url -OutFile "gfpl.exe"

Write-Host "downloading examples scripts..."
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/$Repo/refs/heads/main/test/fibonacci.lamb" -OutFile "fibonacci.lamb"
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/$Repo/refs/heads/main/test/hello_world_complicated.lamb" -OutFile "hello_world_complicated.lamb"

Write-Host "successfully downloaded gfpl and example scripts."
Write-Host "to test gfpl, you can run:"
Write-Host ".\gfpl.exe fibonacci.lamb"
Write-Host "please look at the git repo if you need help: github.com/$Repo"
