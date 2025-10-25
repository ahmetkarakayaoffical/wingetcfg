# ==============================================
# SCN Orion Plus - WingetCFG Go Module Fix Script
# Author: Ahmet Karakaya
# ==============================================

# Proje klasörü
$ProjectPath = "D:\scnorionplus\wingetcfg"

# Versiyon bilgisi
$OldTag = "v0.0.03"
$NewTag = "v0.0.3"
$ModulePath = "github.com/ahmetkarakayaoffical/wingetcfg"

Write-Host "➡️  WingetCFG module version fix starting..." -ForegroundColor Cyan

# 1️⃣ Proje dizinine git
Set-Location $ProjectPath

# 2️⃣ Eski tag varsa sil
$tags = git tag
if ($tags -contains $OldTag) {
    Write-Host "🗑️  Removing old tag: $OldTag"
    git tag -d $OldTag
    git push origin ":refs/tags/$OldTag"
}

# 3️⃣ Yeni tag oluştur
Write-Host "🏷️  Creating new tag: $NewTag"
git tag $NewTag
git push origin $NewTag

# 4️⃣ go.mod içeriğini kontrol et
$goModPath = Join-Path $ProjectPath "go.mod"
$goModContent = Get-Content $goModPath

# Eğer eski versiyon varsa düzelt
if ($goModContent -match $OldTag) {
    Write-Host "✏️  Updating go.mod version to $NewTag"
    $goModContent = $goModContent -replace $OldTag, $NewTag
    $goModContent | Set-Content $goModPath
}

# 5️⃣ module satırı doğru mu?
if ($goModContent -notmatch $ModulePath) {
    Write-Host "⚠️  Module path incorrect, fixing..."
    $goModContent = $goModContent -replace "module .*", "module $ModulePath"
    $goModContent | Set-Content $goModPath
}

# 6️⃣ Modülleri düzenle
Write-Host "🧹  Running go clean & tidy..."
go clean -modcache
go mod tidy

# 7️⃣ Commit ve push
Write-Host "💾  Committing changes..."
git add go.mod go.sum
git commit -m "Fix: corrected module path and version ($NewTag)"
git push origin main

Write-Host "✅ All done! Module successfully tagged and fixed." -ForegroundColor Green
