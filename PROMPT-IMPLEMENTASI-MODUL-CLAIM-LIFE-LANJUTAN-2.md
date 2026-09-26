# PROMPT — lanjutan 2 modul Claim Life: sesudah tiket 06 (`e2190bb`), **sisa modul dalam SATU giliran**

> **Skill:** ketik `/mattpocock-skills:implement` sebagai manusia, lalu tempel berkas ini **utuh**.
> Brief modul **`PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE.md`** §1–§11 dan brief lanjutan 1 §3–§4
> *(tiket 04, 05, 15, 07)* berlaku **seluruhnya**. Berkas ini menambah tiga hal: **aturan satu modul
> per giliran** *(§1)*, **paket keputusan** yang membuka tiket 09–11 *(§2)*, dan fakta XML baru
> untuk tiket 04, 05, 08 *(§5)*.
>
> **Urutan baca:** §1–§2 berkas ini → brief modul §1.2 dan §3 → tiket yang sedang dikerjakan utuh →
> §5 berkas ini untuk tiket itu → XML-nya *(brief modul §2)*.
>
> **GILIRAN INI:** Langkah 0 → **04 → 05 → 15 → 07 → 08 → 09 → 10 → 12 → 11 → 13** — **seluruhnya,
> tanpa mengakhiri giliran di antara dua tiket.** Laporan hanya di akhir *(§6)*.

---

## 0. KEADAAN AWAL — 26 September 2026 malam, sesudah `e2190bb`

| | Keadaan |
| --- | --- |
| `HEAD` | `e2190bb` *claim-life: tiket 06 — validasi DOL per Type, ContentNote, pintu tanggal kejadian*; sebelumnya `2a99841` *(docs)*, `7aa0e94` *(tiket 03)*. Working tree bersih |
| Uji yang lulus | vet · vet db · gofmt nol · build · **158 PASS · 0 FAIL · 24 SKIP** · `tsc` · **5** JS · **88** modul |
| Verifikasi independen `e2190bb` | seluruh angka tereproduksi *(10 berkas +1.109/−5; 14 test baru; rute `PUT …/tanggal-kejadian`; 422 dengan kalimat XML)*; dibaca ulang sendiri dari korpus: pemanggil `ValidasiDOL_Act` memang **dua Section** *(`ClaimLifeDetailGCNM`, `EditDateClaimLife_Section`)*, `ValidasiSTNC_Act` memang salinan *(`pzOriginalInstanceKey` → `VALIDASIDOL_ACT`)*, langkah 4 memang `Property-Set-Messages` pada `.DATE_OF_LOSS`. Satu ketidaktepatan kecil: `utils.ParseTanggal` masih dipanggil di **dua** tempat di luar `repository` *(`handlers/dol.go` untuk mengurai masukan HTTP — wajar di batas; `models.Peserta.JendelaValuasi`)*; klaim "satu-satunya … dipindah" tidak tepat, akibatnya nol; tetap `[terbuka — tiket 14]` sampai tanggal peserta berhenti bertipe teks |
| Tiket | 01 `claimed` 4/14 · 14 `claimed` 40/53 · 02 `claimed` 7/26 · 03 `claimed` 11/28 · 06 `claimed` 4/5 · **04, 05, 15, 07, 08, 09, 10, 12, 11, 13 `ready-for-agent`** |
| Penjaga | `CREATE` = 20 · `Buka()` = 8 · `kolomSalin` = 24 · penyebut `EDMSTATUS` = 2 · `Type` satu rumah tersimpan *(`satutype_test.go`)* |

---

## 1. ATURAN SATU MODUL PER GILIRAN `[DIPUTUSKAN — work owner, 26 September 2026 malam]`

Kata work owner: *"Jika mengerjakan tiket jangan hanya 1 tiket. Kalau bisa mengerjakan 1 modul.
Yang cepat dan jelas hasilnya, dan terlebih hasilnya tepat dan sesuai."* Akibatnya:

1. **Giliran tidak berakhir sesudah satu tiket.** Sesudah commit tiket, tulis **satu baris** kemajuan
   *(`tiket NN — <SHA> — AC x/y — test a PASS · b SKIP`)* lalu **langsung** mulai tiket berikutnya. Nol
   pertanyaan ke manusia di tengah giliran kecuali §7 *(persetujuan)*. Laporan lengkap **satu kali** di
   akhir giliran *(§6)*.
2. **Giliran berakhir sah** hanya bila: tiket 13 sudah di-commit *(atau ditandai `claimed` dengan
   sebab)*; **atau** seluruh tiket yang tersisa terkunci gerbang di luar executor *(Oracle, keputusan
   `[USULAN]`, modul lain)* **sesudah** semua yang dapat dikerjakan tanpa gerbang itu selesai dan
   ter-commit; **atau** ada hal §7; **atau** merah yang tidak dapat dipulihkan. Berhenti sesudah commit
   hijau, tiket yang belum disentuh tetap `ready-for-agent`.
3. **Yang TIDAK dikorbankan demi cepat:** verifikasi penuh sebelum commit **dan** sesudah perbaikan
   review; `/code-review` satu kali per tiket *(dua sumbu paralel)*; satu commit per tiket; XML menang
   *(brief modul §1.2)*; klaim di tiket ditulis dari keluaran perintah.
4. **Yang dipangkas demi cepat:** *(a)* prosa tiket per bab dibatasi — pembacaan ulang XML ≤ 40 baris,
   Implementasi ≤ 50 baris, hasil review ≤ 25 baris; tabel, bukan esai; *(b)* pembuktian penjaga
   "dapat gagal" hanya untuk penjaga yang mengunci **aturan bisnis** *(rumus, gerbang, transisi)*, satu
   baris per penjaga; *(c)* XML yang **sudah dibaca dan dicatat** di tiket 03/06 **tidak dibaca ulang** —
   kutip bab dan nomor barisnya; *(d)* ronde review kedua hanya bila perbaikan menyentuh AC; *(e)* nol
   pengulangan argumen keputusan yang sudah `[USULAN]`/`[DIPUTUSKAN]` — cukup rujuk butirnya.
5. **Urutan tetap** brief modul §3: 04 → 05 → 15 → 07 → 08 → 09 → 10 → 12 → 11 → 13. Tiket yang
   sebagian terkunci tetap dikerjakan sejauh mungkin *(fungsi murni + test + antarmuka yang gagal
   terang)* dan ditandai; **tidak** menunggu.

---

## 2. PAKET KEPUTUSAN YANG MEMBUKA TIKET 09–11 — work owner mengesahkan **sebelum** menempel

Tanpa paket ini, tiket 09, 10, 11 hanya dapat dibangun sampai antarmuka. Satu kata di baris berikut
mengesahkan **seluruh** paket; butir per butir dapat ditimpa di tabel di bawahnya.

**PAKET MODUL: `[USULAN]`** ← ganti menjadi `[DIPUTUSKAN]` untuk menyetujui aj + am + af + al sekaligus.

| | Keputusan | Rekomendasi | Akibat bila `[DIPUTUSKAN]` |
| ---: | --- | --- | --- |
| **aj** | langkah migrasi **`011`**: `BRANCH_OF_BANK`, `SWIFT_CODE`, `PAYABLE_TO` pada `T_CLAIMLF_ADJUSTMENT` *(layar Pega: enam field bank)* | setujui | dikerjakan sebagai commit tambahan tiket 03 **di akhir giliran** |
| **am** *(baru)* | tabel **jejak audit** `T_CLAIMLF_JEJAK` *(per baris adjustment: `ADJUSTMENT_ID` FK, `DARI`, `KE`, `JENIS` transisi/jalur-balik, `AKUN_ID`, `WAKTU`, `CATATAN`)* lewat langkah migrasi bernomor baru — ADR-0007, tiket 09 | setujui | tiket 09 penuh; penjaga cacah `CREATE` naik dan dicatat di tiket 14 |
| **af** | tabel **Komite** `T_GENERAL_KOMITE` + `T_KOMITE_KOMITELIST` dan kolom `COVER_KEY` dibuat **di modul ini** *(tiket 10)*, memakai DDL tiket 00 Komite Claim Life *(`.scratch/komite-claim-life/issues/00-…md`, masih `ready-for-agent`, nol migrasinya di `APP_RNM`)* | setujui, dengan catatan tiket 00 Komite kelak hanya **memverifikasi**, bukan membuat ulang | tiket 10 dan 11 penuh; tiket 00 Komite diberi satu baris rujukan |
| **al** *(baru)* | `DiagnoseList` *(kelas Pega `Data-DiagnoseLife`, per peserta; sasaran `SetSTS_Reject`)* — apakah menjadi tabel `T_CLAIMLF_DIAGNOSE` *(tiket 08 Medical Check)* atau tidak disimpan | **baca XML dulu** *(§5-08)*; setujui bila `MedicalCheckClaimLife` memang mengisi daftar diagnosa | tiket 08 menyimpan diagnosa; tiket 04 mencerminkan status ke sana |
| **ag** | sumber kunci rate dan daftar treaty-year *(modul Master Product Name Life)* | ag1 — modul berikutnya | bukan giliran ini |
| **o1–o3** | penomoran | masih `[USULAN]` | bila `[DIPUTUSKAN]`: `PenomorCounter`, commit tambahan tiket 02 di akhir giliran |
| ah, ak | rate ganda; mata uang ketiga | Product+UW | — |

Yang masih `[USULAN]` saat giliran berjalan **tidak ditebak**: tabelnya tidak dibuat, penulisnya di
balik antarmuka yang gagal terang, AC-nya `[terbuka — keputusan X]`.

---

## 3. LANGKAH 0

`git add PROMPT-IMPLEMENTASI-MODUL-CLAIM-LIFE-LANJUTAN-2.md` → commit `docs: brief lanjutan 2 —
satu modul per giliran, paket keputusan` → `git status --porcelain` kosong → uji tanpa Oracle hijau
*(158 · 24 SKIP · 5 JS · 88 modul)* → SHA = titik tetap **tiket 04**.

---

## 4. YANG SUDAH ADA DAN WAJIB DIPAKAI ULANG

Brief lanjutan 1 §3 **ditambah** dari tiket 06: `services.ValidasiDOL`, `ContentNoteDari`,
`GalatDOL{PesertaID}`, `ErrDOLKosong`; `models.BusinessCodeContentNote()`, `Peserta.TanggalKejadian`,
`Peserta.JendelaValuasi`; `KlaimLife.TypeKlaim`, `PerbaruiTanggalKejadian`; rute
`PUT /api/klaim-life/{id}/peserta/{pesertaId}/tanggal-kejadian`; konstanta `TypeQR/QP/TR/TP`.

---

## 5. FAKTA XML BARU — dibaca sendiri 26 September 2026 malam, mengubah tiket 04, 05, 08

### Tiket 04 — mesin status per baris

| Fakta `[terverifikasi]` | Bukti | Akibat |
| --- | --- | --- |
| `SetSTS_Reject` langkah 1 adalah `Property-Set` **pada `.DiagnoseList`** *(kelas `ASM-FW-GISFW-Data-DiagnoseLife`)*: `.STS_REJECT = Primary.STS_REJECT` — menyalin status **peserta** ke tiap baris **diagnosa**, **bukan** ke baris adjustment | `Activity\SetSTS_Reject.xml` pecahan 236–258 *(`pyStepsObjectName .DiagnoseList`, `pyStepsClassName …Data-DiagnoseLife`)* | tabel "Rule Pega sumber" tiket 04 *(baris `SetSTS_Reject` — "penurunan status klaim → baris")* **diralat** per brief modul §1.2-b. `DiagnoseList` **tidak ada** di skema kita → pencerminan ke sana **tidak ditiru** sampai **al** diputuskan; dicatat |
| Pemanggil `SetSTS_Reject` hanya `Section\ClaimLifeDetailGCNM.xml` | grep korpus Claim Life + Komite Claim Life | ia tindakan layar, bukan bagian mesin status |
| Penulis nyata: `SaveOutStandingLife_Act` → `0` per baris adjustment bergerbang `PrintFaceClaim` *(tiket 03)*; `RejectOSClaimLife_Act` → `2` pada `.STS_REJECT` **dan** `PremiumListDetail(idx).STS_REJECT` *(baris + peserta)*; `KomitePostAdjustment` → `1`/`2` *(Komite; baca di `Komite Claim Life\Activity\`)* | tiket 03 bab pembacaan; `RejectOSClaimLife_Act.xml` pecahan 443–491 | "pencerminan dua tingkat" = **baris adjustment + peserta**. Tingkat header *(`T_WORK_CLAIM.ACCEPT_STATUS`)*: cari penulisnya di XML *(`CreateKMTLife_Act`, `KomitePostAdjustment`, flow Komite)*; bila nol penulis di Claim Life, header **tidak** dicerminkan di tiket ini dan dicatat |

Mesin: `Transisi(baris, ke)` hanya dari Outstanding; final tidak berubah; `"4"` → `TidakDiketahui`,
tidak pernah ditulis; `Klaim.StatusTurunan()` dihitung; `PohonKlaim.PerbaruiStatusBaris` mencerminkan
baris + peserta dalam satu transaksi; konstanta kode **hanya** di `models`. Bila **al** `[DIPUTUSKAN]`
dan tabel diagnosa lahir di tiket 08, pencerminan ke diagnosa ditambahkan **di tiket 08**, bukan di sini.

### Tiket 05 — Reject Outstanding oleh Admin

| Fakta `[terverifikasi]` | Bukti | Akibat |
| --- | --- | --- |
| `RejectOSClaimLife_Act` langkah 2 menyetel `.STS_REJECT = 2`, `PremiumListDetail(local.IndexPremium).STS_REJECT = 2`, **dan `PremiumListDetail(…).IsCheck = "false"`** | pecahan 443–518 | **AC baru menurut XML**: menolak baris juga **mencabut penanda dipilih** peserta *(`IS_CHECK`)*; ralat tiket per §1.2-b |
| Langkah 3 merakit `TempInputDetail.CARI1…CARI13` *(`pzInsKey`, `CLAIM_NO`, `POLICY_NO`, `POLICY_HOLDER`, `CERTIFICATE_NO`, `NAME_OF_INSURED`, `SEX`, `DOB`, `AGE`, `PLAN`, `BEGIN_DATE`, `DATE_OF_LOSS`, `EXPIRED_DATE` — tanggal `dd/MM/YYYY` `Asia/Jakarta`)* untuk Connect-SQL sesudahnya | pecahan 631–922 | baca langkah 4+ *(Connect-SQL mana; kolom apa yang ditulis ke tabel warisan)*; jalur tulis warisan kita *(`BarisLamaDari` + `PeriksaNilaiWarisan`)* mencerminkan `STS_REJECT`; ⛔ `NAME_OF_INSURED`/`POLICY_HOLDER` **tetap tidak disalin** ke tabel klaim baru *(pagar keamanan)* — bila baris warisan menuntutnya, ambil dari sumber warisan saat menulis baris warisan saja, dan catat |

### Tiket 08 — Medical Check & Claim Analis

Sebelum menulis apa pun: baca `Section\MedicalCheckClaimLife.xml` dan cari kelas
`ASM-FW-GISFW-Data-DiagnoseLife` di korpus *(`grep -l DiagnoseLife`)* — bila layar telaah medis
mengisi `.DiagnoseList` per peserta *(kode/nama diagnosa, tanggal, catatan)*, itulah bukti untuk **al**
dan bentuk tabelnya *(kolom dari properti yang diisi, path + baris per kolom)*. Tanpa **al**
`[DIPUTUSKAN]`: tahap dan jalur balik tetap dibangun *(`T_WORK_CLAIM.PY_POSITION`, `SENDTO_ADMIN`,
`SENDTO_MEDICAL` sudah ada di `001`)*, diagnosa di balik antarmuka.

### Tiket 09, 10, 12, 11, 13

Persis brief modul §4 dan brief lanjutan 1; gerbangnya paket §2. Untuk 12: nol pemanggilan
`GET_TOKEN_STORAGE` *(procedure — keputusan **o**)*, nol endpoint nyata, stub + antrean + audit.
Untuk 13: skrip + rekonsiliasi di skema uji saja; `POOLDATA` tidak disentuh.

---

## 6. LAPORAN AKHIR GILIRAN — bentuknya

1. **Tabel per tiket**: tiket · SHA titik tetap · SHA commit · AC ditutup/diralat/total · test PASS/SKIP
   sesudahnya · yang terkunci dan oleh apa.
2. **Keputusan §2** yang `[DIPUTUSKAN]` dan dipakai; yang `[USULAN]` dan akibatnya *(AC mana)*.
3. **Terbuka menurut pemilik**: work owner · Product+UW · DBA · IAM · modul lain — satu baris per butir.
4. **Telemetri** brief modul §10, ditambah: cacah tiket ter-commit di giliran ini, cacah XML dibaca
   *(berkas, byte)*, cacah ralat menurut XML, token sub-agen review per tiket.

---

## 7. PERSETUJUAN MANUSIA — berhenti dan tanya

Brief modul §9. Langkah migrasi baru **hanya** dari paket §2 yang `[DIPUTUSKAN]` *(011 bank; jejak
audit; tabel Komite; tabel diagnosa)*. Selain itu: procedure, JSON produk, layanan luar nyata,
migrasi data ke skema mana pun selain user kosong DBA, `git push`.

---

*Disusun 26 September 2026 malam sesudah verifikasi independen `e2190bb`: uji dijalankan ulang, diff
10 berkas dibaca, lima klaim XML dibaca ulang dari korpus, `SetSTS_Reject` dan `RejectOSClaimLife_Act`
dipecah dan dibaca, tiket 00 modul Komite dan daftar migrasi `APP_RNM` diperiksa.*
