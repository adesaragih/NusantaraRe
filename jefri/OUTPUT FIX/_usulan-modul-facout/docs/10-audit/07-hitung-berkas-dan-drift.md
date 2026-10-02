# Audit 07 — Hitungan berkas per subfolder, kelas rule, dan drift antar folder

> **Diukur ulang 21 September 2026.** Sumber angka bagi `05-tickets\facout\F11-jalur-facout-endorsement.md`
> dan `F12-pendukung-layar-facout.md`, serta bagi **DRAF K-058**.

---

## 1. Jumlah `.xml` per subfolder

⚠️ Dihitung **per subfolder**, satu enumerasi masing-masing. Listing rekursif sekali jalan terbukti
terpotong di batas entri dan menghasilkan angka salah.

`[terverifikasi]`

| Subfolder | NB FacIn | RNW Fac In | Endorsment | Total |
| --- | ---: | ---: | ---: | ---: |
| Activity | 609 | 556 | 582 | 1.747 |
| ConnectREST | 3 | 3 | 3 | 9 |
| DataPage | 30 | 30 | 41 | 101 |
| DataTransform | 148 | 138 | 157 | 443 |
| DecisionTable | 12 | 10 | 10 | 32 |
| DecisionTree | 1 | 1 | 1 | 3 |
| Flow | 6 | 4 | 2 | 12 |
| FlowAction | 250 | 239 | 265 | 754 |
| Harness | 41 | 41 | 37 | 119 |
| RDBList | 217 | 200 | 188 | 605 |
| ReportDefinition | 123 | 118 | 138 | 379 |
| Section | 432 | 397 | 434 | 1.263 |
| SystemSettings | 1 | 1 | 1 | 3 |
| When | **210** | **189** | **202** | **601** |
| **TOTAL** | **2.083** | **1.927** | **2.061** | **6.071** |

✅ Cocok `CLAUDE.md` §2 (6.071 = 2083+1927+2061) dan §4.5 (`When` = 601).
✅ Cocok `00-RINGKASAN-EKSEKUTIF.md` baris 182 ("605 berkas SQL") — itu subfolder `RDBList`.

---

## 2. Distribusi `<pxObjClass>` per subfolder

`[terverifikasi]` **Seragam mutlak** — setiap subfolder memuat **satu kelas rule saja**, nol
pengecualian:

| Subfolder | Berkas | `<pxObjClass>` |
| --- | ---: | --- |
| Activity | 1.747 | `Rule-Obj-Activity` ×1747 |
| ConnectREST | 9 | `Rule-Connect-REST` ×9 |
| DataPage | 101 | `Rule-Declare-Pages` ×101 |
| **DataTransform** | 443 | **`Rule-Obj-Model` ×443** |
| DecisionTable | 32 | `Rule-Declare-DecisionTable` ×32 |
| DecisionTree | 3 | `Rule-Declare-DecisionTree` ×3 |
| Flow | 12 | `Rule-Obj-Flow` ×12 |
| FlowAction | 754 | `Rule-Obj-FlowAction` ×754 |
| Harness | 119 | `Rule-HTML-Harness` ×119 |
| **RDBList** | 605 | **`Rule-Connect-SQL` ×605** |
| ReportDefinition | 379 | `Rule-Obj-Report-Definition` ×379 |
| Section | 1.263 | `Rule-HTML-Section` ×1263 |
| SystemSettings | 3 | `Rule-Admin-System-Settings` ×3 |
| When | 601 | `Rule-Obj-When` ×601 |

⛔ **Dilaporkan apa adanya. TIDAK disimpulkan** apakah `Rule-Obj-Model` wajar untuk folder bernama
`DataTransform`. Itu **`[di luar korpus]`**.

> ❓ **Pertanyaan terbuka untuk work owner.** Folder `RDBList\` berisi `Rule-Connect-SQL` — divergensi
> nama-folder vs tipe-rule yang sudah tercatat (`PANDUAN-KERJA` §3 butir 6). Apakah `DataTransform\`
> berisi `Rule-Obj-Model` merupakan **penamaan internal Pega yang wajar**, atau **divergensi sejenis**?
> Jawabannya menentukan apakah jebakan "nama folder bukan tipe rule" berlaku sekali atau dua kali.

---

## 3. Drift antar folder

⚠️ Metode: hash **kontrak 23 tag** (K-042/K-043, termasuk amandemen "isi blok `pzIndexes` diabaikan"),
dibandingkan per identitas rule (**tipe rule + nama berkas**, `HashSet` Ordinal, per folder).

`[terverifikasi]`

| Pasangan | Dibandingkan | Versi sama + isi sama | **Versi sama + isi BEDA** | Versi beda + isi sama | Versi beda + isi beda | Versi beda (total) |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| NB ↔ RNW | 1.907 | 1.907 | **0** | 0 | 0 | **0** |
| NB ↔ EDM | 1.707 | 1.353 | **267** | 0 | 87 | **87** |
| RNW ↔ EDM | 1.674 | 1.330 | **259** | 0 | 85 | **85** |
| **Total** | **5.288** | — | **526** | **0** | **172** | **172** |

**Pemeriksaan silang:** `1.353 + 267 + 87 = 1.707` ✅ dan `267 + 87 = 354` ✅ — cocok angka resmi
K-043 (1.353 identik / 354 berbeda).

### 3.1 Temuan yang menentukan DRAF K-058

⛔ **Dari 354 perbedaan NB↔EDM, 267 (75,4 %) ber-`pyRuleSetVersion` IDENTIK.** Metode deteksi drift
yang membandingkan nomor versi **melewatkan tiga dari empat perbedaan nyata**.

📌 **`versi beda + isi sama` = 0 pada ketiga pasangan.** Tidak pernah ada kasus nomor versi berubah
tanpa isi ikut berubah. Artinya: kenaikan versi **selalu** menyertai perubahan isi, tetapi perubahan
isi **sering tidak** menaikkan versi. Arah kegagalannya satu arah — dan itulah yang membuat metode
lama bias ke arah "tidak ada drift".

⚠️ `[pertanyaan terbuka]` **Angka "95 rule berbeda versi"** di `PANDUAN-KERJA` §5 dan `CLAUDE.md` §4.5
**tidak tereproduksi**. Pengukuran ini memberi **87** (NB↔EDM), **85** (RNW↔EDM), **0** (NB↔RNW),
total **172** lintas tiga pasangan. Metode yang menghasilkan 95 tidak terdokumentasi, jadi selisihnya
tidak dapat dijelaskan dari korpus.

---

## 4. Angka di dokumen `OUTPUT\` yang perlu dikoreksi

| Dokumen | Angka tertulis | Hasil ukur | Sikap |
| --- | --- | --- | --- |
| `PANDUAN-KERJA.md` §5 · `CLAUDE.md` §4.5 | "95 rule berbeda versi" | **87 / 85 / 0** (172 total) | ⚠️ tidak tereproduksi — **butuh keputusan work owner**, bukan suntingan sepihak (`CLAUDE.md` adalah berkas aturan) |
| `07-edm\02-e1-pola-perbedaan-when.md` baris 160 & 210 | `When`: 112 identik / 63 berbeda dari 175 | **119 identik / 56 berbeda** dari 175 | ⚠️ basi sejak **amandemen K-043** (`pzIndexes`). Baris 210 menyatakan "angka `When` tidak berubah oleh koreksi metode" — benar untuk metode #2/#3, **tidak bertahan** terhadap amandemen |
| Catatan lisan sesi sebelumnya | `NB\When` = 112 atau 196 | **210** | ✅ dikoreksi di dokumen ini |
| Catatan lisan sesi sebelumnya | "18 activity Fac Out", "10/18 beda" | **39 hadir di ketiganya; 9/39 beda** | ✅ dibatalkan, lihat `03-struktur-facretro-dan-populasi.md` |

📌 Angka "112" di enam dokumen lain (`00-RINGKASAN-NB.md`, `01-skema-field-nb.md`, `02-model-data.md`,
`03-spec-modul-terverifikasi.md`, ADR-0006, `00-RINGKASAN-EKSEKUTIF.md`) **bukan** hitungan berkas rule
— itu "112 Section menampilkan uang tanpa mata uang" dan "112 tabel Oracle". **Tidak perlu dikoreksi.**

---

## 5. Perintah audit

```powershell
# (1) jumlah per subfolder -> total 6071
$folders=@('NB FacIn','RNW Fac In','Endorsment Fac In')
$subs=@('Activity','ConnectREST','DataPage','DataTransform','DecisionTable','DecisionTree','Flow',
        'FlowAction','Harness','RDBList','ReportDefinition','Section','SystemSettings','When')
$g=0
foreach($s in $subs){ $baris=$s.PadRight(20)
  foreach($f in $folders){
    $p="D:\migrasi\RNM\$f\$s"; $n=0
    if([IO.Directory]::Exists($p)){ foreach($x in [IO.Directory]::EnumerateFiles($p,'*.xml')){$n++} }
    $baris+="$n".PadLeft(8); $g+=$n }
  $baris }
"TOTAL=$g"
```

```powershell
# (2) Get-Sig : kontrak 23 tag + isi blok pzIndexes diabaikan + Ordinal sort
$VOLATILE=@('pxCreateDateTime','pxUpdateDateTime','pxSaveDateTime','pxCommitDateTime','pxMoveImportDateTime',
 'pxOriginalCreateDateTime','pyRuleFormStatusTime','pxWarningCreatedTime','pxCreateOperator','pxCreateOpName',
 'pxCreateSystemID','pxUpdateOperator','pxUpdateOpName','pxUpdateSystemID','pxMoveImportOperId',
 'pxMoveImportOperName','pxOriginalCreateOperator','pxOriginalCreateOpName','pxOriginalCreateSystemID',
 'pxHostId','pzChecksum','pzIndexCount','pyShowJavaWindowName')
$KEYLIKE=@('pzInsKey','pzIndexOwnerKey','pzOriginalInstanceKey','pzDocumentKey','pyJavaClassName','pxInsName')
$rxVol=[regex]('^</?('+($VOLATILE -join '|')+')[ />]|^<('+($VOLATILE -join '|')+')>')
$rxKey=[regex]('^<('+($KEYLIKE -join '|')+')>')
$rxTs=[regex]'\d{8}T\d{6}(\.\d+)?( GMT)?|_\d{8}T\d{6}_\d+'
$sha=[Security.Cryptography.SHA256]::Create()
function Get-Sig([string]$path){
  $keep=New-Object Collections.Generic.List[string]; $dalam=$false
  foreach($ln in [IO.File]::ReadAllLines($path)){
    $tr=$ln.Trim(); if($tr -eq ''){continue}
    if($dalam){ if($tr.StartsWith('</pzIndexes>')){$dalam=$false}; continue }
    if($tr.StartsWith('<pzIndexes')){ if(-not $tr.Contains('</pzIndexes>')){$dalam=$true}; continue }
    if($rxVol.IsMatch($tr)){continue}
    if($rxKey.IsMatch($tr)){ $tr=$rxTs.Replace($tr,'<TS>') }
    $keep.Add($tr) }
  $a=$keep.ToArray(); [Array]::Sort($a,[StringComparer]::Ordinal)
  [BitConverter]::ToString($sha.ComputeHash([Text.Encoding]::UTF8.GetBytes([string]::Join("`n",$a)))).Replace('-','') }
# uji silang WAJIB sebelum dipakai:
#   When\FlagOldData -> BEDA ; When\IsAdmin, IsBonding, IsAddButton, DataPage\D_AnekaList -> identik
```

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan.*
