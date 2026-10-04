# Risiko Migrasi + Pertanyaan untuk Bisnis

> Sumber: sintesis seluruh dokumen discovery di `D:\migrasi\RNM\OUTPUT\`, dibangun dari korpus
> `D:\migrasi\RNM\{NB FacIn, RNW Fac In, Endorsment Fac In}\` (6.071 berkas `.xml`).
> Label mengikuti `CLAUDE.md` §3.
>
> Daftar **lengkap** butir yang tidak terjawab ada di `02-data/03-batas-pengetahuan.md` (33
> memblokir, 29 perlu keputusan). Dokumen ini memilih yang **paling mahal bila salah** dan
> mengurutkannya untuk dibawa ke bisnis.

---

## 1. Risiko berdasarkan dampaknya terhadap rekonsiliasi paralel run

Rekonsiliasi paralel run adalah ujian keberhasilan migrasi: sistem lama dan baru dijalankan
berdampingan, dan angkanya harus cocok. Risiko diurutkan menurut seberapa luas selisih yang
ditimbulkannya.

### R0 — ⛔ Ketiga folder korpus diekspor dari **titik waktu yang berbeda**

**Dampak:** rekonsiliasi paralel run **mustahil** sampai diperbaiki. Risiko ini mendahului semua
yang lain karena ia melemahkan dasar pembandingan itu sendiri.

`[terverifikasi]` (`01-flow/03-alur-endorsement.md`) Dari 1.571 rule bernama sama, **86 punya
`pyRuleSetVersion` berbeda** antara NB dan EDM. Arahnya **tidak konsisten**:

- `SetToInbox_ACT` di EDM lebih **lama** (2025-09 vs 2026-08 di NB/RNW) dan **kehilangan satu cabang
  routing**.
- `SetBanding_ACT` di EDM justru lebih **baru** dan **menambah** klausa persetujuan.

**Konsekuensi terhadap discovery ini:** sebagian dari 460 rule yang terhitung "bercabang antar
siklus" adalah selisih versi ekspor, bukan percabangan bisnis — dan korpus ini **tidak dapat
memisahkan keduanya**.

⚠️ Kesimpulan "tiga siklus di atas satu basis rule" **tetap berdiri** — ditopang 1.220 rule identik
dan NB=RNW 1.680 dari 1.680. Yang tidak berdiri adalah pemakaian angka 460 sebagai ukuran volume
kerja percabangan.

**Yang diminta:** satu ekspor ulang dari lingkungan **produksi**, pada satu titik waktu, untuk
ketiga siklus sekaligus. Ini juga menyelesaikan R6.

### R1 — ⛔ Arti `ProposalAcceptStatus` tidak diketahui

**Dampak:** menyalahkan arah **setiap** keputusan underwriting.

`[terverifikasi]` `DecisionTable/IsUWAccepted` mendeklarasikan enam hasil (`confirm`, `reject`,
`ask`, `banding`, `revise`, `decline`); Data Transform menulis `1`/`2`/`3`/`4`/`7`/`9`. **Baris yang
memetakan keduanya tidak ikut terekspor.**

Bukti bahwa menebak berbahaya: nilai `4` ditulis rule bernama *banding* tetapi diuji rule bernama
*fac out* (`RNW Fac In\When\IsFacout.xml:148`). Nama rule tidak dapat dipakai menebak arti.

**Mitigasi:** tidak ada yang teknis. Butuh jawaban bisnis, atau ekspor ulang DecisionTable dari Pega
**sebelum Pega dimatikan**.

### R2 — ⛔ Perbandingan uang sebagai string

**Dampak:** tangga akseptasi dan proteksi spreading mengambil cabang yang salah pada rentang nilai
tertentu.

`[terverifikasi]` Ambang **yang sama** (30 miliar) dibandingkan dua cara di korpus yang sama:
sebagai string mentah (`local.TSIinIDR > "30000000000"`) dan sebagai desimal
(`TotalTSINusaRe > toDecimal("30000000000")`).

Perbandingan string bersifat leksikografis: `"4000000000" > "30000000000"` bernilai **BENAR** karena
`'4' > '3'` — padahal 4 miliar lebih kecil dari 30 miliar.

**Keputusan yang diminta:** apakah sistem baru **mereproduksi** perilaku string (paralel run cocok,
tetapi hasilnya keliru secara aritmetika) atau **memperbaikinya** (hasil benar, paralel run berbeda
dan selisihnya harus dijelaskan)?

**Rekomendasi:** reproduksi dulu, perbaiki sebagai perubahan terpisah setelah paralel run bersih.
Dasarnya `CLAUDE.md` §1 — perbaikan dipisahkan dari migrasi.

### R3 — ⛔ Satuan `.Rate` tidak konsisten

**Dampak:** premi meleset **satu ordo besaran 10**.

`[terverifikasi]` `CountPremi_ACT` membagi dengan `1e9`; `CountProrateExtension_Act` step 7 membagi
dengan `1e4`. Keduanya memperlakukan `.Rate` sebagai masukan yang sama.

**Keputusan yang diminta:** `.Rate` disimpan dalam per mille atau persen? Bila keduanya dipakai di
tempat berbeda, apa yang menentukannya?

### R4 — ⛔ Isi tabel limit tidak ada di korpus

**Dampak:** tangga akseptasi **tidak dapat direkonsiliasi sama sekali**.

`[terverifikasi]` Mekanismenya terbaca penuh — query, urutan, cara berhenti. Yang tidak ada adalah
**datanya**: `LIMIT_BOTTOM`, `LIMIT_BOTTOM2`, `MAX_LIMIT_IDR/USD`, dan ejaan pasti nilai kolom
`JABATAN`.

⚠️ Ejaan itu bukan detail sepele: token routing di kode tanpa spasi (`DIREKTURTEKNIK`) sedangkan
ejaan kolom Oracle belum diketahui. Bila tidak cocok, **tidak ada approver yang pernah ditemukan**.

**Yang diminta:** salinan isi tabel `M_LIMIT_*` dari DBA. Bukan data pelanggan — data master
wewenang.

### R5 — ⛔ Sumber pengganti tabel internal Pega

**Dampak:** jalur banding dan flag reject kehilangan sumbernya.

`[terverifikasi]` `DATAPEGA.pc_History_ASM_FW_GISFW_Work` dibaca `GetHistoryAccPega_SQL`
**untuk mengambil keputusan alur**, bukan sekadar audit. Tabel ini internal Pega — **hilang bersama
Pega** (`CLAUDE.md` §4.3).

**Keputusan yang diminta:** riwayat akseptasi dipindahkan ke tabel aplikasi sendiri (perlu migrasi
data historis), atau keputusan alur diturunkan dari sumber lain?

### R6 — ⛔ Empat rule hilang dari ekspor

**Dampak:** jalur **Special Acceptance** dan konversi ke produksi tidak dapat dimigrasikan.

`GetLimitAkseptasi1SA_Act` · `GetLimitAkseptasi2SA_Act` · `GetLimitAkseptasi_Act2` ·
`serviceInsertArasapas_act` (kelas Work) · `UpdateErrorNoteJsonPolis` · 21 rujukan Section.

**Ini risiko dengan tenggat keras.** Ekspor ulang dari Pega menyelesaikannya — tetapi hanya selama
Pega masih hidup.

### R7 — ⚠️ 460 rule bercabang EDM belum dipetakan

**Dampak:** tidak diketahui.

`[terverifikasi]` Dari 1.680 identitas rule yang hadir di ketiga folder, **1.220 logikanya identik**
dan **460 bercabang** — seluruhnya EDM vs NB+RNW. Jumlahnya terhitung; **isinya belum dibandingkan
field-per-field**.

Sebelum dipetakan, **tidak boleh** diklaim "hanya beda kosmetik". Ini pekerjaan discovery lanjutan,
bukan pertanyaan bisnis.

### R8 — ⚠️ Nama orang di dalam logika dan di teks status

`[terverifikasi]` Enam rule memuat guard berbasis identitas orang; salah satunya melewati seluruh
proses spreading. Selain itu **669 kemunculan / 23 label status** dan pola *"IS IN … INBOX"* dengan
14 nama unik di satu flow saja — teks yang **tampil ke pengguna**.

**Keputusan yang diminta:** (a) guard dipertahankan apa adanya atau dipindah ke atribut peran?
(b) bolehkah teks status yang dilihat pengguna berubah?

Pilihan (a) "dipertahankan" membuat paralel run cocok tetapi mewariskan aturan berbasis orang ke
sistem baru. Pilihan "dipindah" mengubah hasil untuk kasus-kasus itu.

### R9 — ⚠️ Authz: `IsAdmin` fail-open

`[terverifikasi]` Operator dengan workbasket kosong dihitung sebagai admin. `pyWorkBasketList(1)`
juga hanya menguji workbasket pertama, dengan urutan tak ditentukan.

`CLAUDE.md` §6 mensyaratkan **persetujuan manusia** sebelum perubahan authn/authz — termasuk
keputusan untuk **tidak** mengubahnya.

### R10 — ⛔ Prompt AI disimpan sebagai data tabel, bukan kode

**Dampak:** perilaku model AI di alur underwriting dapat berubah **tanpa deployment dan tanpa jejak
version control**.

`[terverifikasi]` (`02-data/02-skema-oracle.md`) Tabel `M_PROMPT_AI` (kunci `KATEGORI_1` +
`KATEGORI_2`, kolom `PROMPT_AI`) memasok prompt tersebut.

**Konsekuensi:** tinjauan keamanan yang disyaratkan `CLAUDE.md` §6 tidak cukup mencakup rule — ia
harus mencakup **tabel dan siapa yang boleh menulisinya**. Di sistem baru, perubahan prompt
seharusnya melewati jalur yang meninggalkan jejak.

### R11 — ⛔ Jalur produksi endorsement menulis data pribadi

`[terverifikasi]` `FACINPRODUCTION` menulis 79 kolom di NB/RNW tetapi **82 di EDM** — dan bukan
sekadar tambahan: EDM menulis `NO_KONTRAK`, `NO_NPWP`, `NO_KTP`, `STARTDATE_DEBITUR`,
`ENDDATE_DEBITUR` **sebagai pengganti** `DEDUCTION2_MENJADI`/`_SELISIH`.

**Konsekuensi:** dua kolom memuat data pribadi. `CLAUDE.md` §6 mensyaratkan persetujuan sebelum akses
*regulated data*. Selain itu, hilangnya pasangan `DEDUCTION2_*` pada baris endorsement perlu
dikonfirmasi — apakah agregasi hilir mengandaikan kolom itu terisi?

### R12 — ⛔ Riwayat akseptasi tidak boleh dipangkas

`[terverifikasi]` `HISTORYAKSEPTASIPEGA` **di-SELECT untuk mengambil keputusan**, bukan hanya
di-INSERT: `GetAksepBanding_SQL` mengambil workbasket terakhir lalu men-join ke tabel limit untuk
memperoleh jalur banding; `GetFlagReject_SQL` menghitung baris berstatus reject.

**Konsekuensi:** kebijakan arsip/purge apa pun atas tabel ini **mengubah keputusan alur**. Ini harus
dinyatakan eksplisit sebelum go-live.

### R13 — ⚠️ Dua FlowAction bernuansa AI belum dibaca

`[terverifikasi]` `AnalysLocationbyAI` dan `AttachDoc_AI` ada di NB dan RNW, tidak ada di EDM.
**Isinya belum dibaca** — nama bukan bukti perilaku (`CLAUDE.md` §3 butir 3).

Bila benar memanggil model AI pihak ketiga di alur underwriting, `CLAUDE.md` §6 mengharuskan
**tinjauan keamanan sebelum dipindahkan**.

---

## 2. Pertanyaan untuk bisnis — diurutkan

Disusun agar dapat dibawa ke satu rapat. Kolom "bila tidak dijawab" menyatakan konsekuensinya, bukan
ancaman.

| # | Pertanyaan | Pemilik | Bila tidak dijawab |
| ---: | --- | --- | --- |
| **0** | **Dapatkah dilakukan satu ekspor ulang dari lingkungan produksi, pada satu titik waktu, untuk ketiga siklus?** | IT | **Rekonsiliasi paralel run mustahil** — mendahului semua pertanyaan lain |
| 1 | Apa arti masing-masing nilai `ProposalAcceptStatus` (`1`,`2`,`3`,`4`,`7`,`9`)? Mengapa `5`,`6`,`8` tidak dipakai? | Underwriting | Seluruh alur persetujuan tidak dapat diimplementasikan |
| 2 | Ambang uang yang di sistem lama dibandingkan **sebagai string** — direproduksi atau diperbaiki? | Underwriting + Keuangan | Tangga akseptasi salah cabang pada rentang tertentu |
| 3 | `.Rate` disimpan per mille atau persen? | Aktuaria / Underwriting | Premi meleset 10× |
| 4 | Dapatkah isi tabel `M_LIMIT_*` diberikan, termasuk ejaan pasti kolom `JABATAN`? | DBA | Tangga akseptasi tidak dapat direkonsiliasi |
| 5 | Riwayat akseptasi (kini di tabel internal Pega) dipindahkan ke tabel aplikasi? | IT + Underwriting | Jalur banding kehilangan sumber keputusan |
| 6 | Dapatkah ekspor ulang Pega dilakukan untuk 4 rule + 12 DecisionTable + 21 Section yang hilang? | IT | Special Acceptance & konversi produksi tidak dapat dimigrasikan |
| 7 | Kamus 98 `BusinessCode` dan 87 `BusinessOldId` | Product | 80 rule lini bisnis hanya dapat disalin, tidak diverifikasi |
| 8 | Guard berbasis identitas orang: dipertahankan atau dipindah ke peran? | Underwriting | Perilaku spreading untuk sebagian kasus tidak ditentukan |
| 9 | Bolehkah teks status yang memuat nama orang berubah? | Bisnis | Teks pengguna tidak dapat difinalkan |
| 10 | `IsAdmin` fail-open dipertahankan? | Keamanan | Perubahan authz tanpa persetujuan — dilarang `CLAUDE.md` §6 |
| 11 | Arti `Type` (0–9,11,12), `EdmType`, `DeductibleType` (dua ruang nilai tak beririsan) | Product | Klasifikasi produk tidak dapat diverifikasi |
| 12 | Tanda PPh/PPN dalam net payable — empat konvensi berbeda hidup berdampingan | Keuangan | Nilai bayar bersih tidak dapat difinalkan |
| 13 | `IsEdmAdjCeding` dan `IsEdmPPNPPH` berkondisi identik — apa bedanya sebenarnya? | Underwriting | Dua jenis endorsement tidak dapat dibedakan |
| 14 | Renewal dinilai atas TSI penuh, bukan selisih — benar? | Underwriting | Tangga akseptasi renewal mungkin salah dasar |
| 15 | 936 elemen UI yang dimatikan keras (`1=2`/`Never`) — disengaja atau sisa? | Bisnis | Layar baru mungkin menampilkan yang seharusnya tersembunyi |
| 16 | Dapatkah DDL (tipe kolom) 112 tabel diberikan? Terutama `FACINOFFER` | DBA | Bila `VARCHAR2`, perbandingan angka di atasnya adalah perbandingan string |
| 17 | Dapatkah isi 24 stored procedure `POOLDATA` diberikan (`ALL_SOURCE`)? | DBA | **Seluruh jalur tulis melewatinya** — tidak dapat direplikasi |
| 18 | Siapa yang boleh menulis tabel `M_PROMPT_AI`, dan apakah perubahannya perlu meninggalkan jejak? | Keamanan + IT | Perilaku model AI dapat berubah tanpa deployment |
| 19 | Kolom `NO_NPWP` / `NO_KTP` di jalur produksi endorsement — persetujuan akses *regulated data* | Kepatuhan | `CLAUDE.md` §6 melarang lanjut tanpa persetujuan |
| 20 | Apakah `HISTORYAKSEPTASIPEGA` pernah diarsip/dipangkas? | DBA + Underwriting | Pemangkasan **mengubah keputusan alur**, bukan sekadar menghapus audit |
| 21 | Arti `EdmType` `1`/`2`/`3` — dua rule memberi label "batal sejak semula" pada nilai yang **berbeda** | Product | Jenis endorsement tidak dapat dibedakan; `EdmType==3` membalik tanda seluruh selisih persentase |
| 22 | Blok `EdmType==1` tidak pernah dapat dieksekusi (kedua blok bergerbang `EdmType==2`) — cacat atau memang tidak terpakai? | Underwriting | Satu jenis endorsement kehilangan penanganannya |
| 23 | Apakah endorsement **jiwa** memang tanpa persetujuan berjenjang? | Underwriting | Tangga akseptasi jiwa tidak dapat diimplementasikan |
| 24 | Case endorsement ber-fac-retro berakhir **tanpa baris `JSON_POLIS`** — disengaja? | Underwriting + IT | Data polis endorsement retro tidak tersimpan |
| 25 | `EDMPremiMenjadi` memakai basis berbeda antar lini bisnis (FIRE vs ANEKA/GOLF) | Aktuaria | Nilai premi endorsement tidak dapat difinalkan |

---

## 3. Risiko metodologis — dan cara menghindarinya

Tiga kekeliruan berikut **terjadi** selama discovery ini dan menghasilkan angka yang salah total
sebelum dikoreksi. Siapa pun yang mengaudit ulang korpus akan menemuinya.

| # | Kekeliruan | Akibat | Koreksi |
| --- | --- | --- | --- |
| M1 | Membaca kondisi rule `When` hanya dari `<pyConditionString>` | **170 berkas (28,3 %)** tampak "tanpa kondisi" padahal kondisinya ada di `<pyConditionValue1String>`. Ini yang membuat `CLAUDE.md` §4.5 mendaftar lima rule sebagai "kosong" padahal **tidak ada satu pun** yang benar-benar kosong | Baca **kedua** tag |
| M2 | `Get-FileHash` mentah untuk membandingkan rule antar folder | Melaporkan **0** dari 1.680 identik. Penyebabnya metadata ekspor: operator impor, `pxCommitDateTime`, `pxHostId` | Buang baris `px*`/`pz*` |
| M3 | Membandingkan tanpa mengurutkan, atau dengan `Sort-Object` bawaan | Urutan tag ekspor Pega **tidak deterministik**; `Sort-Object` bawaan case-insensitive dan tidak stabil padahal korpus memuat nama rule yang hanya beda kapitalisasi. Angka berubah dari 507 menjadi **1.220** setelah dikoreksi | Urutkan dengan `[StringComparer]::Ordinal` |

⚠️ `[terverifikasi]` Nilai `pxHostId` yang berbeda (`jboss122117` dan `jboss1073`) mengonfirmasi
`CLAUDE.md` §3 butir 1: korpus dirakit dari **lebih dari satu server Pega**. Karena itu menyebut nama
rule saja tidak cukup — **path berkas wajib disebut**.

---

## 4. Urutan kerja yang diusulkan

**[usulan]** — bukan temuan korpus.

1. **Segera, selagi Pega masih hidup** — dua pekerjaan yang **tidak dapat dibatalkan** setelah Pega
   padam:
   a. Ekspor ulang 4 rule + 12 DecisionTable + 21 Section yang hilang (R6).
   b. Ekstraksi data historis riwayat akseptasi dari tabel internal Pega (R5), plus `ALL_SOURCE`
      untuk 24 stored procedure (pertanyaan 17).
2. **Sebelum baris Go pertama:** jawab pertanyaan 1–4 di §2. Keempatnya memblokir inti aplikasi.
3. **Paralel, tidak memblokir:** discovery lanjutan untuk spreading/capacity/scoring, dan pemetaan
   field-per-field atas 460 rule bercabang EDM (R7).
4. **Implementasi:** mulai dari `pkg/money` dan `internal/rules` — keduanya punya dasar terverifikasi
   dan tidak menunggu jawaban bisnis.
5. **Terakhir:** mesin tangga akseptasi, setelah pertanyaan 1, 2, dan 4 terjawab.

---

## 5. Lampiran — skrip audit percabangan antar-siklus

Menghasilkan angka yang dikutip di `05-migrasi/01-arsitektur-target.md` §1.1 dan §4.

```powershell
# Perbandingan logika rule antar-folder. Tiga koreksi wajib (lihat §3 M1-M3).
$md5 = [System.Security.Cryptography.MD5]::Create()
$ord = [StringComparer]::Ordinal
$volatil = '^\s*<(px|pz)[A-Za-z0-9_]*[ >/]|^\s*</(px|pz)[A-Za-z0-9_]*>|^\s*<pyRuleFormStatusTime>|^\s*<pyShowJavaWindowName>|^\s*<rowdata REPEATINGINDEX="(RuleReference|Warnings)">'

function LogicHash($p) {
    $keep = [Collections.Generic.List[string]]::new()
    foreach ($l in [IO.File]::ReadAllLines($p)) { if ($l -notmatch $volatil) { $keep.Add($l.Trim()) } }
    $arr = $keep.ToArray()
    [Array]::Sort($arr, $ord)
    return [BitConverter]::ToString($md5.ComputeHash([Text.Encoding]::UTF8.GetBytes(($arr -join "`n"))))
}

$all = foreach ($f in @("NB FacIn","RNW Fac In","Endorsment Fac In")) {
    Get-ChildItem "D:\migrasi\RNM\$f" -Recurse -Filter *.xml -File | ForEach-Object {
        [PSCustomObject]@{ Folder=$f; Key=$_.Directory.Name+"|"+$_.BaseName; Path=$_.FullName } } }

$g3 = $all | Group-Object Key | Where-Object { ($_.Group.Folder | Sort-Object -Unique).Count -eq 3 }
$same = 0; $diff = [Collections.Generic.List[string]]::new(); $nbRnwSame = 0
foreach ($grp in $g3) {
    $map = @{}
    foreach ($m in $grp.Group) { $map[$m.Folder] = (LogicHash $m.Path) }
    if ($map["NB FacIn"] -ceq $map["RNW Fac In"]) { $nbRnwSame++ }
    if (($map.Values | Sort-Object -Unique).Count -eq 1) { $same++ } else { $diff.Add($grp.Name) }
}
"identitas di 3 folder : " + $g3.Count      # -> 1680
"  logika IDENTIK      : $same"             # -> 1220
"  logika BERCABANG    : " + $diff.Count    # -> 460
"  NB == RNW           : $nbRnwSame"        # -> 1680 dari 1680
```

Hasil per tipe rule:

| Tipe | Bercabang |
| --- | ---: |
| Activity | 187 |
| Section | 93 |
| When | 65 |
| FlowAction | 39 |
| DataTransform | 21 |
| RDBList | 18 |
| ReportDefinition | 17 |
| Harness | 13 |
| DecisionTable | 5 |
| Flow | 1 |
| DataPage | 1 |
| **Total** | **460** |

⚠️ **Batas metode:** angka 460 adalah **batas atas** percabangan perilaku. Normalisasi membuang
metadata ekspor dan urutan tag, tetapi tidak membedakan perbedaan yang mengubah perilaku dari
perbedaan yang tidak (mis. teks peringatan, daftar kelas halaman). Untuk rule `When`, perbandingan
yang lebih tajam — hanya atas `pyLogic` + ekspresi kondisi — menghasilkan **39**
(`04-aturan/01-katalog-when.md` §3.2). Angka setajam itu untuk tipe rule lain **belum dihitung**.

---

*Setiap pertanyaan di dokumen ini menunggu jawaban manusia. `CLAUDE.md` §1: bila korpus tidak
menjelaskan sesuatu, itu pertanyaan terbuka — bukan ruang untuk berimprovisasi.*
