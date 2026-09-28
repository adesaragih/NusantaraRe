# 04: Mesin status per baris + aturan turunan "klaim selesai"

**Status:** sebagian — header tidak selalu mencerminkan baris terakhir (jalur aksep Claim Life dan Komite); `ACCEPTED_NO` tidak dicerminkan ke peserta

**Blocked by:** 03 (baris `AdjustmentList` + Save ke Outstanding)

## Hasil & nilai pengguna

Sebagai **pengguna mana pun**, saya melihat status tiap baris adjustment secara terpisah dan tahu
persis apa yang sudah diputuskan dan apa yang belum — dan sebuah klaim yang seluruh barisnya sudah
diputus tampak "selesai" sehingga antrean kerja saya bersih. *(User story 25 dan 28 di spec)*

Ini **inti spesifikasi**: tiket yang menetapkan bahwa yang diputuskan adalah baris, bukan klaim.

## Area codebase

`internal/models` (status baris sebagai nilai tertutup), `internal/services` (transisi + aturan
turunan), `internal/handlers` (endpoint transisi), `frontend/` (tampilan status klaim dan baris).

Status klaim **dihitung**, bukan disimpan sebagai kolom mandiri — meskipun sistem lama
menyimpannya sebagai kolom.

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/SaveOutStandingLife_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY` | penulis `0` |
| `Komite Claim Life/Activity/KomitePostAdjustment.xml` | `ASM-FW-GCNMFW-WORK-KOMITELIFE` / `KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` penulis `1` (6 `Property-Set`) dan `2` (2 `Property-Set`) |
| `Claim Life/Activity/RejectOSClaimLife_Act.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `REJECTOSCLAIMLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` penulis `2` dari sisi Admin |
| `Claim Life/Activity/SetSTS_Reject.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SETSTS_REJECT` / `RULE-OBJ-ACTIVITY` | ⚠️ **DIRALAT 26-09-2026 menurut XML** — bukan "penurunan status klaim → baris". `[terverifikasi]` ia `Property-Set` atas **`.DiagnoseList`** (kelas `Data-DiagnoseLife`, pecahan 236–258): status **peserta** disalin ke tiap baris **diagnosa**, bukan ke baris adjustment. Pemanggilnya hanya `Section/ClaimLifeDetailGCNM.xml` |
| `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SERVICEINSERTARASAPASCLAIMLIFE_ACT` | ⭐ **BARU 26-09-2026** — `[terverifikasi]` satu-satunya penulis tingkat **header**: `pyWorkPage.ClaimData.STS_REJECT` dan `.ACCEPTEDNO` dari baris adjustment, di putaran bersarang tanpa henti (pecahan 371–418) |

`[terverifikasi]` **Sensus penulis lengkap** ada di register **OQ-061**: tidak ada satu pun rule di
`Claim Life` maupun `Komite Claim Life` yang menulis `0` setelah `1` atau `2`.

## ADR terkait

**ADR-0011** (unit = baris; kefinalan; arti tunggal nilai `2`), **ADR-0001** (jalur balik Komite
bekerja pada tingkat baris).

## Acceptance criteria

- [x] Baris yang sudah bernilai Aksep atau Ditolak **tidak dapat berubah lagi** melalui jalur mana
      pun. *(AC 2 spec)* — bukti: `APP_RNM/internal/services/statusbaris.go:Transisi`, `APP_RNM/internal/repository/klaimlife.go:KlaimLife.PerbaruiStatusBaris` (`WHERE … AND STS_REJECT = :kodeLama`); uji `TestBarisFinalTidakDapatBerubah`, `TestTransisiHanyaDariOutstanding`
- [x] Menolak sebuah baris **tidak** menutup klaim; klaim tetap dapat menerima baris baru.
      *(AC 4 spec)* — bukti: uji `TestMenolakSatuBarisTidakMenutupKlaim`, `TestKlaimTidakTerminalSetelahPenolakan`
- [x] Status "Ditolak" pada sebuah baris **selalu** berarti baris itu ditolak — **tidak pernah**
      berarti klaim selesai, apa pun sumber penolakannya. — bukti: `APP_RNM/internal/models/statusklaim.go:Klaim.StatusTurunan`; uji `TestStatusKlaimTurunan`
- [ ] ⚠️ **Diselaraskan 2026-09-16:** pencerminan terjadi pada kolom `STS_REJECT` /
      `ACCEPTED_NO` di **`T_CLAIMLF_PREMIUMLIST_DETAIL`** dan **`T_GENERAL_CLAIM`** (spec §2b) — unit
      keputusannya tetap baris `T_CLAIMLF_ADJUSTMENT` (**ADR-0011**), yang kini menggantung pada
      **peserta**. *(AC 33 spec; penyimpangan sadar 2)* — belum: `STS_REJECT` dicerminkan ke peserta dan header (`PerbaruiStatusBaris`, `CerminkanHeader`), tetapi `ACCEPTED_NO` hanya ke header — kolomnya tidak ada di `T_CLAIMLF_PREMIUMLIST_DETAIL`
- [ ] `PremiumListDetail` dan header klaim **selalu mencerminkan** baris adjustment terakhir, dan
      **tidak** ditulis sebagai status mandiri. *(AC 7 spec)* — belum: `Status.ubah` memakai `BarisTerakhir`, tetapi `Akseptasi.SimpanAdjustment` (`akseptasi.go`) dan jalur Komite (`komite_akseptasi.go`) mencerminkan header dengan status baris yang diputus, bukan baris terakhir klaim
- [x] Klaim dilaporkan "selesai" **hanya** bila tidak ada baris berstatus Outstanding **dan** ada
      sekurangnya satu baris berstatus Aksep. *(AC 8 spec)* — bukti: `APP_RNM/internal/models/statusklaim.go:Klaim.StatusTurunan`; uji `TestStatusKlaimTurunan`
- [x] Bila tidak ada baris Outstanding dan tidak ada pula yang Aksep, klaim berada dalam keadaan
      **ditolak seluruhnya** — dan tetap dapat dilanjutkan dengan baris baru. — bukti: `APP_RNM/internal/models/statusklaim.go:Klaim.StatusTurunan` (`KlaimDitolakSeluruhnya`); uji `TestStatusKlaimTurunan`, `TestKlaimTidakTerminalSetelahPenolakan`
- [x] Riwayat lengkap seluruh baris pada satu klaim dapat dilihat, sehingga putaran Komite terbaca.
      *(User story 27 spec)* — bukti: `APP_RNM/internal/repository/klaimlife.go:KlaimLife.AmbilBaris` (seluruh baris, urut `ID`), `APP_RNM/frontend/src/pages/claimlife/KlaimLife.tsx:KlaimLife`

## Catatan penutupan (2026-09-14)

**Aturan turunan "klaim selesai" disetujui work owner** `[keputusan work owner]` dan tercatat di
`CONTEXT.md` (Lampiran, butir 1):

- Status klaim adalah **turunan**, **bukan** kolom tersimpan.
- Unit keputusan = **baris `AdjustmentList`** (**ADR-0011**).
- **Header klaim mengikuti adjustment terakhir.**

`[terverifikasi 2026-09-14]` **Di Pega**, `AdjustmentList` **tidak punya tabel fisik** (kelas
`ASM-FW-GISFW-Data-AdjustmentLife`, berawalan `Data-` = embedded), dan barisnya di-persist ke
`POOLDATA.OS_AKSEPTASI_KLAIM_LIFE` yang kolomnya identik. Bukti:
`Claim Life/Activity/SaveOutStandingLife_Act.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` /
`SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`) — `pyStepsObjectName = .AdjustmentList`,
`pyStepsClassName = ASM-FW-GISFW-Data-AdjustmentLife`.

⚠️ **Diselaraskan 2026-09-16 — itu keadaan Pega, bukan keadaan sistem baru.** `[keputusan work
owner]` Di sistem baru baris adjustment punya **tabelnya sendiri**, `T_CLAIMLF_ADJUSTMENT`, yang
menggantung pada **peserta** (`T_CLAIMLF_PREMIUMLIST_DETAIL`) — bukan pada header klaim, dan **bukan**
pada `OS_AKSEPTASI_KLAIM_LIFE` (spec §2b, AC 33).

⚠️ **Koreksi 2026-09-16** `[keputusan work owner]`: `OS_AKSEPTASI_KLAIM_LIFE` **tetap di-`INSERT`
flat**, berdampingan dengan tabel relasional — hilir masih membaca dari sana. Yang **dibuang hanya
JSON**. Jadi baris adjustment tersimpan di **`T_CLAIMLF_ADJUSTMENT`** *dan* ikut terbawa ke rekam flat
itu; yang berubah adalah **di mana unit keputusan tinggal**, bukan hilangnya tabel akseptasi.
*(AC 32 spec)*

**Ini tidak mengubah ADR-0011** — unit keputusan tetap **baris**. Yang berubah hanya di mana baris
itu tinggal: dari tabel akseptasi bersama menjadi tabel adjungan per peserta.

`[data DBA]` `STS_REJECT` pada tabel **warisan** bertipe `NUMBER(38)`. ⚠️ Di skema baru tipenya
ditetapkan sendiri di tiket **14**; nilainya tetap `0`/`1`/`2` dan diisi **menurut aksi**, bukan
di-hardcode (AC 42 spec).

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Pembacaan ulang XML — 26 September 2026 malam

Yang sudah dibaca di tiket 03 tidak dibaca ulang; dirujuk saja. Yang **baru** dibaca giliran ini:
sensus penulis `STS_REJECT` di **seluruh** korpus `Claim Life`, dan `serviceInsertArasapasClaimLife_act`
(81.079 byte → 1.893 baris pecahan).

**Sensus penulis `STS_REJECT`** — ⚠️ **DIRALAT di bab hasil review di bawah: yang benar ENAM
`Property-Set` di LIMA rule.** Tabel di bawah ini adalah sensus ronde pertama yang **kurang satu
baris**; ia dibiarkan terbaca apa adanya supaya kekeliruannya tidak hilang:

| Rule | Properti | Tingkat | Nilai |
| --- | --- | --- | --- |
| `SaveOutStandingLife_Act` | `.STS_REJECT` | baris adjustment | `0` *(tiket 03, bergerbang `PrintFaceClaim`)* |
| `RejectOSClaimLife_Act` | `.STS_REJECT` | baris adjustment | `2` |
| `RejectOSClaimLife_Act` | `PremiumListDetail(local.IndexPremium).STS_REJECT` | **peserta** | `2` |
| `SetSTS_Reject` | `.STS_REJECT` atas `.DiagnoseList` | **diagnosa** | `Primary.STS_REJECT` |
| `serviceInsertArasapasClaimLife_act` | `pyWorkPage.ClaimData.STS_REJECT` | **header klaim** | `.STS_REJECT` |

### ⭐ Tiga hal yang mengubah tiket

**1. Header klaim MEMANG dicerminkan — dan sumbernya baris TERAKHIR.** `[terverifikasi]`
`serviceInsertArasapasClaimLife_act` langkah 1 mengulang `PremiumListDetail`, langkah 1.1 mengulang
`.AdjustmentList`, dan langkah 1.1.1 *(pecahan baris 371–418)* menyetel
`pyWorkPage.ClaimData.ACCEPTEDNO = .ACCEPTEDNO` dan `pyWorkPage.ClaimData.STS_REJECT = .STS_REJECT`
**tanpa precondition** dan tanpa henti. Nilai yang bertahan milik baris yang diulang terakhir — persis
AC *"header klaim selalu mencerminkan baris adjustment terakhir"*, yang sampai kini hanya berdiri di
atas catatan work owner. Kini `[terverifikasi]`.

**2. Kedua tingkat TIDAK mencerminkan baris yang sama.** Peserta mengikuti baris yang **berubah**
(`local.IndexPremium` = peserta pemilik baris itu); header mengikuti baris **terakhir**. Karena itu
`PerbaruiStatusBaris` dan `CerminkanHeader` dipisah, bukan disatukan demi kerapian.

**3. `SetSTS_Reject` bukan bagian mesin status — ralat tabel "Rule Pega sumber".** Teks lama: *"penurunan
status klaim → baris"*. Yang sebenarnya: ia `Property-Set` atas **`.DiagnoseList`** *(kelas
`Data-DiagnoseLife`, pecahan 236–258)*, menyalin status **peserta** ke tiap baris **diagnosa** — bukan
ke baris adjustment. Pemanggilnya hanya `Section/ClaimLifeDetailGCNM.xml`, jadi ia tindakan layar.
⛔ `DiagnoseList` tidak ada di skema kita; pencerminan ke sana **tidak ditiru** sampai butir **al**
diputuskan.

### Pertanyaan yang XML tidak jawab

| # | Pertanyaan | Pemilik |
| ---: | --- | --- |
| 1 | Apakah `DiagnoseList` menjadi tabel — butir **al** | work owner *(bukti dari tiket 08)* |
| 2 | Arti kode `4` di data warisan | work owner |

---

## Implementasi — 26 September 2026 malam (tiket 04)

**Status: `claimed`** — **6 dari 8 AC tertutup.** Titik tetap `e1bd6e3`.
Angka verifikasi ada di bab hasil review di bawah, sesudah perbaikannya.

| Berkas | Isi |
| --- | --- |
| `models/statusklaim.go` *(baru)* | `StatusKlaim` tertutup + `Klaim.StatusTurunan()` |
| `models/klaimlife.go` | `statusTurunan` masuk kontrak API, berdampingan dengan `kodeStatus` mentah |
| `services/statusbaris.go` *(baru)* | `Transisi`, `TandaiOutstandingKlaim`, `BarisTerakhir`, layanan `Status.Ubah` |
| `repository/klaimlife.go` | `PerbaruiStatusBaris` *(baris + peserta)*, `CerminkanHeader` |
| `models/kodestatus_test.go` *(baru)* | penjaga statik: literal kode status hanya di `models` |

**Penjaga aturan bisnis, masing-masing dibuktikan dapat gagal:** kefinalan dilepas → 2 test merah ·
gerbang sudah-berstatus dilepas → 1 test merah · syarat *"ada Aksep"* dilepas → 2 test merah ·
literal kode lewat variabel lokal → penjaga statik merah.

⚠️ Penjaga literal kode **percobaan pertama lebih sempit daripada kalimatnya sendiri**: ia hanya
mencocokkan `KodeStatus = "n"`, sehingga mutasi lewat variabel lokal lolos dan hijaunya tak berarti.
Polanya diperlebar dan dibuktikan ulang. Test `TestMenolakSatuBarisTidakMenutupKlaim` juga diperkuat
di ronde yang sama: ia semula tidak dapat gagal ketika kefinalan dilepas dari jalur `TandaiOutstanding`.

### ⛔ Yang BELUM ditutup — 2 AC

| AC | Sebab |
| --- | --- |
| *"pencerminan terjadi pada `STS_REJECT`/`ACCEPTED_NO` di peserta dan header"* | mekanismenya ada, teruji, dan dipanggil `Status.Ubah` — tetapi **pintu HTTP-nya milik tiket 05**. Ditutup di sana, bukan dicentang di sini dengan pemanggil yang belum ada |
| *"`PremiumListDetail` dan header selalu mencerminkan baris terakhir"* | sama |

⚠️ Pencerminan ke **diagnosa** tidak ditiru: butir **al** masih `[USULAN]`, tabelnya tidak dibuat, dan
menebaknya berarti mengarang tabel. Bila **al** disahkan, pencerminan itu ditambahkan di **tiket 08**
bersama tabelnya — bukan di sini.

### Hasil `/code-review` atas titik tetap `e1bd6e3`

⛔ **Sensus di bab di atas SALAH, dan labelnya `[terverifikasi]`.** Yang benar **enam** `Property-Set`
di **lima** rule. Yang terlewat: `Activity/SaveAdjustment_Act.xml` *(pecahan 1833–1899)*, yang
menyetel **`.AdjustmentList(<LAST>).STS_REJECT = 1`** bersama `.ACCEPTEDNO` dan
`.ACCEPTATION_DATE = @CurrentDateTime()` dalam **satu** `Property-Set`, tanpa precondition.

**Sebabnya pola pencarian saya sendiri.** Sensus memakai `<PropertiesName>[^<]*STS_REJECT`;
teks yang sudah didekode memuat `<LAST>`, sehingga `[^<]*` berhenti di tanda `<` dan baris itu tak
pernah cocok. Kelas kekeliruan yang persis sama dengan yang saya buru sepanjang sesi ini — penjaga
atau pencarian yang lebih sempit daripada kalimatnya. Pola benar: `<PropertiesName>.*STS_REJECT`.

**Dua akibat nyata pada kode, bukan sekadar catatan:**

1. `1` **bukan hanya ditulis Komite** — jalur simpan Claim Life sendiri menulisnya. Dokumentasi
   `Transisi` yang menisbatkannya ke `KomitePostAdjustment` saja sudah diralat.
2. **`ACCEPTATION_DATE` lahir bersama akseptasi.** Tiket 03 AC 43 benar bahwa ia tidak distempel
   saat insert; kini terbaca **di mana** ia distempel. `Transisi` menerima jam sebagai argumen dan
   menstempelnya hanya pada Aksep; `PerbaruiStatusBaris` menulis kolomnya.

| # | Sumbu | Temuan | Tindakan |
| ---: | --- | --- | --- |
| 1 | Spec | penulis keenam terlewat; `[terverifikasi]` palsu | ✅ sensus diralat di sini dan di komentar `repository/klaimlife.go` |
| 2 | Standards | **ADR-U-0007 dilanggar** — `Ubah` menerima `pelaku` lalu membuangnya; nol jejak | ✅ antarmuka `Jejak` + `JejakBelumDiputuskan` yang **gagal terang**, direkam **di dalam** transaksi yang sama. Pola `Penomor` tiket 02. Butir **am** masih `[USULAN]`, jadi tabelnya tidak dibuat |
| 3 | Standards | **TOCTOU** — baris dibaca di luar transaksi, `UPDATE` tanpa syarat kode lama; kefinalan hanya berlaku di dalam proses | ✅ `AND STS_REJECT = :kodeLama`; dua permintaan serentak tidak lagi dapat sama-sama menang |
| 4 | Standards | komentar berkata "TIGA tingkat" padahal menulis dua; klaim "seluruh isi MURNI" padahal `Ubah` menyentuh Oracle | ✅ keduanya diralat |
| 5 | Standards + Spec | `TandaiOutstandingKlaim` salinan kembar `TandaiOutstanding` | ✅ tiket 03 kini **mendelegasi**; satu tempat aturannya |
| 6 | Standards | `BarisTerakhir` menulis lewat parameter **nilai** — kemampuan menulis yang menyelinap | ✅ parameternya penunjuk, sehingga terlihat |
| 7 | Standards | penjaga literal kode bocor *(SQL `'2'` tak tertangkap)* **dan** salah tuduh *(komentar tidak dibuang)*, tanpa pemeriksaan diri | ✅ ketiganya diperbaiki — dan penjaga yang diperlebar langsung menemukan satu pelanggaran **lama** di `migrasidata.go`, yang ikut dibetulkan |
| 8 | Standards | `StatusTurunan` menghitung kode asing `"4"` sebagai *"ditolak seluruhnya"* | ✅ keadaan baru `KlaimTidakDapatDipastikan`. Melaporkan klaim SUDAH diputus padahal keputusannya justru yang tak terbaca adalah arah kekeliruan terburuk |
| 9 | Standards | ADR-U-0029 Akibat 2 — nol test kegagalan di tengah | ✅ dua test `db` baru *(SKIP tanpa Oracle)*: pencerminan tiga tingkat, dan pembatalan yang tidak meninggalkan separuh jadi |
| 10 | Spec | AC *"klaim selesai"* tercentang padahal `statusTurunan` tidak pernah sampai ke layar | ✅ dikabelkan ke `api.ts` dan `KlaimLife.tsx`, berdampingan dengan kode mentah dan dibedakan terang |

### ⛔ Yang tetap terbuka sesudah review

| Butir | Pemilik |
| --- | --- |
| Gerbang peran: ADR-U-0011 menuntut `ReasLifeAdmin` untuk reject; `Ubah` baru menuntut identitas | tiket 07; `WajibPeran` sudah ada dan masih nol pemanggil |
| Tabel jejak audit — butir **am** | work owner |
| Penulis keenam menulis pula ke tabel **datar** warisan lewat `UpdateOsAkseptasiClaimLife_sql`; `Ubah` belum menyentuhnya | tiket 05 *(ADR-U-0042, jalur `BarisLamaDari`)* |
| *"Baris terakhir"* = ID terbesar menurut urutan baca; hubungannya dengan *"yang paling akhir diputus"* tidak dinyatakan di mana pun | work owner |

**Verifikasi sesudah perbaikan:** vet · vet db · gofmt nol · build · **167 PASS · 0 FAIL** *(dari
158)* · **26 SKIP** · `tsc` · 5 JS · 88 modul.

### Dua AC sisa ditutup oleh tiket 05 — 26 September 2026 malam

Pintu `POST /api/klaim-life/{id}/adjustment/{adjId}/tolak` (tiket 05, commit di bawah) memanggil
`Status.Ubah`, sehingga pencerminan ke peserta dan ke header benar-benar terjadi. Keduanya kini
tercentang: yang menahannya memang hanya pemanggil, bukan mekanismenya.

---

### Ralat menurut XML — 27 September 2026 (audit A0, brief lanjutan 4 bab 7)

⛔ **Akseptasi punya DUA jalur; tiket ini hanya mengenal satu.**

| Butir | Teks lama | Teks baru | Bukti |
| --- | --- | --- | --- |
| siapa menulis status Aksep | *"Aksep ditulis modul Komite, bukan Claim Life"* — `ErrAksepBukanDariModulIni` menolak tujuan Aksep bagi siapa pun | **Claim Life mengaksep sendiri** lewat `SaveAdjustment_Act`; sentinelnya **DIHAPUS** | `Claim Life/Activity/SaveAdjustment_Act.xml` pecahan baris **1833** (`ACCEPTEDNO`), **1879** (`STS_REJECT = 1`), **1899** (`ACCEPTATION_DATE = @CurrentDateTime()`) |
| gerbang peran jalur itu | — *(tidak dikenal)* | **pemegang TAHAP**, bukan daftar peran datar | `[terverifikasi — pohon XML]` tombol "Save Adjustment" (`Section/ClaimLifeDetailGCNM.xml` **22641 → 22665**) tidak dibungkus gerbang peran mana pun; seluruh leluhurnya ALWAYS. `pyCondition 1=2` yang tampak bertetangga adalah `pyContainerVisibleWhen` **layout lain** |
| prasyarat dagangnya | — | peserta `.IsCheck=true`, `.ACCEPTEDNO==""`, baris `.STS_REJECT=="0"`, dan `Type` | pecahan **854** (QP/QR) dan **1048** (TP/TR) |

⚠️ **DUA cacat rule Pega yang sengaja TIDAK ditiru — dilaporkan ke work owner:**

1. **Prasyarat tanpa kurung.** Baris **837** dan **1031** berbunyi
   `.IsCheck=true && Type=="QP" || Type=="QR"`. Karena `&&` mengikat lebih erat daripada `||`,
   bacaan harfiahnya meloloskan `QR` dan `TR` **tanpa** memeriksa `IsCheck` sama sekali. Go memakai
   bacaan yang **dimaksud** — `IsCheck && (QP||QR)`.
2. **Tahun tidak bergeser.** Baris **605** berbunyi
   `@if(MM=="12" && NextMonth=="01", @toDecimal(@CurrentDate("YY")), @toDecimal(@CurrentDate("YY")))`
   — **kedua cabangnya identik**, sehingga nomor Januari memakai tahun Desember. Go menggeser
   tahunnya.

⭐ **Temuan yang menghentikan jalur ini di satu tempat:** nomor akseptasi memuat **kode bisnis**
(`'RNML-A'||{pyWorkPage.BusinessCode}||…`), tetapi model relasional kita **tidak menyimpannya** —
`T_GENERAL_CLAIM` hanya punya `BUSINESS_NAME`, dan `BUSINESSID` bukan salah satu dari 18 kolom datar
warisan yang `Simpan` tulis. Jalurnya **gagal terang** dengan `ErrKodeBisnisBelumTersimpan` (HTTP 501)
sampai kolomnya lahir di **A1**. Tidak dikarang.

