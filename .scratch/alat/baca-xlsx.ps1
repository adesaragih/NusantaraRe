# baca-xlsx.ps1 — membaca isi berkas .xlsx korpus TANPA Excel dan TANPA Python.
#
# Untuk apa: korpus `D:\XML\RNM_BRD\<modul>\Struktur_*.xlsx` adalah dokumen
# struktur rule yang ditulis pengekspor Pega (kolom: No, Jenis, Nama Rule,
# Dipanggil dari, XML, Class, Status XML). Baris berlevel `A` (No = 1, 2, …)
# adalah TITIK MASUK modul. Mesin ini tidak punya Python, jadi .xlsx dibaca
# sebagai zip lewat .NET.
#
# Pakai (PowerShell 5.1):
#   & .\.scratch\alat\baca-xlsx.ps1 -Path "D:\XML\RNM_BRD\Claim Life\Struktur_InboxClaimLife.xlsx"
#   & .\.scratch\alat\baca-xlsx.ps1 -Path "…xlsx" | Select-String '^=== SHEET|^r[0-9]+ \| [AB][0-9]+='
#
# Keluaran: satu baris per baris sheet, `r<no> | <sel>=<isi> | …`; sel kosong
# dilewati; baris baru di dalam sel diganti ` / `.
#
# ⛔ Korpus READ-ONLY: skrip ini hanya membaca (ZipFile.OpenRead).
param([Parameter(Mandatory = $true)][string]$Path)
Add-Type -AssemblyName System.IO.Compression.FileSystem
$zip = [System.IO.Compression.ZipFile]::OpenRead($Path)
try {
  function BacaEntri($z, $nama) {
    $e = $z.Entries | Where-Object { $_.FullName -eq $nama }
    if ($null -eq $e) { return $null }
    $r = New-Object System.IO.StreamReader($e.Open())
    try { return $r.ReadToEnd() } finally { $r.Dispose() }
  }
  $strings = @()
  $ssText = BacaEntri $zip 'xl/sharedStrings.xml'
  if ($ssText) {
    [xml]$ssx = $ssText
    foreach ($si in $ssx.DocumentElement.ChildNodes) { $strings += $si.InnerText }
  }
  [xml]$wb = BacaEntri $zip 'xl/workbook.xml'
  [xml]$rels = BacaEntri $zip 'xl/_rels/workbook.xml.rels'
  $relMap = @{}
  foreach ($rel in $rels.DocumentElement.ChildNodes) { $relMap[$rel.Id] = $rel.Target }
  foreach ($sh in $wb.DocumentElement.sheets.ChildNodes) {
    $rid = $sh.GetAttribute('id', 'http://schemas.openxmlformats.org/officeDocument/2006/relationships')
    $target = $relMap[$rid]
    if ($target -notmatch '^/') { $target = 'xl/' + $target } else { $target = $target.TrimStart('/') }
    Write-Output ("=== SHEET: " + $sh.name + " (" + $target + ") ===")
    [xml]$sx = BacaEntri $zip $target
    foreach ($row in $sx.DocumentElement.sheetData.ChildNodes) {
      $cells = @()
      foreach ($c in $row.ChildNodes) {
        $v = $null
        if ($c.t -eq 's') { $v = $strings[[int]$c.v] }
        elseif ($c.t -eq 'inlineStr') { $v = $c.InnerText }
        else { $v = $c.v }
        if ($null -ne $v -and "$v" -ne '') { $cells += ($c.r + '=' + ("$v" -replace "`r?`n", ' / ')) }
      }
      if ($cells.Count -gt 0) { Write-Output ("r" + $row.r + " | " + ($cells -join ' | ')) }
    }
  }
} finally { $zip.Dispose() }
