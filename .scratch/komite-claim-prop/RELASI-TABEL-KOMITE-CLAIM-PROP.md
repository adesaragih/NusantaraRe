# Relasi Tabel — Komite Claim Prop

> ⭐ `[keputusan work owner]` **2026-09-20.** **Awalan tabel klaim non-life yang lama — berhuruf `P` sesudah `CLAIM` — diganti menjadi `T_CLAIM_`**, sebab tabelnya kini **dipakai bersama lini FAC dan PROP** sehingga huruf **P** pada awalan menyesatkan. ⭐ **NONPROP nanti ikut tabel yang sama.**
>
> ⛔ **Nama awalan lamanya sengaja TIDAK dikutip harfiah** di berkas mana pun di luar lini Life — ⭐ supaya pencarian atas awalan lama itu **hanya** menemukan berkas yang memang belum diselaraskan, bukan kalimat yang menerangkan penggantiannya. ⭐ **Tujuh tabel dipakai bersama; lima di antaranya berkunci asing GANDA.** ⛔ **Tabel lini Life tidak disentuh** — awalannya berbeda dan tidak ikut terganti.
>
> ⚠️ Akibatnya **dua berkas lini Life masih menyebut nama tabel penyesuaian yang lama**. ⭐ Itu **dicatat sebagai `[terbuka]`**, ⛔ **bukan diperbaiki**.


**Tanggal:** 2026-09-19 · **Modul:** `Komite Claim Prop` · Sisi induk dilacak di `Claim Prop`
**Acuan nama kolom:** `STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md` — berkas itu **tidak disunting** di sini.

> ⛔ **Berkas lama NOL disunting, tanpa kecuali.** ⛔ Modul lain NOL. ⛔ Kode NOL ·
> `CREATE TABLE` NOL · DDL NOL · nilai rahasia NOL.
> ⛔ **Tidak satu pun butir `[terbuka]` dinyatakan tertutup.**
> ⛔ **Nol nomor baris XML dikutip** — bukti memakai path berkas + `pxInsName` + nomor langkah Pega.

---

## §A — Lima acuan yang mengikat

| # | Acuan | Dipakai di |
| --- | --- | --- |
| **A1** | Nama kolom diambil dari berkas struktur, **bukan** dari ejaan korpus. Ejaan korpus = **bukti asal**, bukan nama | seluruh §B · §C |
| **A2** | Dua kolom usul bernama **`CLOSE_FILE`** dan **`RESERVED_CLAIM`** `[keputusan work owner 2026-09-19]` | §C3 |
| **A3** | `T_GENERAL_KOMITE` = **7** kolom · `T_KOMITE_KOMITELIST` = **9** kolom · keduanya **lintas-lini** | §B · §D |
| **A4** | **Delapan awalan** `T_WORK_CLAIM`, sepasang per lini; modul ini memakai **`CLMP-` / `TKMT-`** | §B relasi 1 · 2 |
| **A5** | `ADJUSTMENT_ID` menunjuk **dua tabel** menurut lini → **tanpa `REFERENCES`**, ditegakkan di Go | §B relasi 4 |

### ⚠️ Ralat penamaan — satu baris, sesuai A2

> `STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md` masih menulis usulan lama **`KOMITE_USUL_TUTUP`** dan
> **`KOMITE_USUL_CADANG`**; yang **BERLAKU** adalah **`CLOSE_FILE`** dan **`RESERVED_CLAIM`**
> `[keputusan work owner 2026-09-19]`. ⛔ Berkas strukturnya **tidak disunting** di tugas ini.

⭐ **Catatan yang menguatkan A2, bukan mengubahnya:** nama baru itu **cocok dengan sasaran di
korpus**. `KomitePostAdjustment` langkah **11** menulis ke `…ClaimData.IsCloseFile` dan
`…ClaimData.IsReservedClaim` — jadi `CLOSE_FILE` / `RESERVED_CLAIM` mengikuti ejaan sisi klaim,
bukan ejaan sisi komite.

---

## §B — Daftar relasi

### Ringkasan

| # | Induk | Anak | Kunci tamu | Kardinalitas | ON DELETE | Index | Status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| **1** | `T_WORK_CLAIM` *(baris klaim)* | `T_WORK_CLAIM` *(baris komite)* | `COVER_KEY` | **1:N** | **di Go** | **biasa** pada `COVER_KEY` | `[keputusan work owner 2026-09-18]` — ⛔ **tidak terbukti dari korpus sebagai kolom** |
| **2** | `T_WORK_CLAIM` *(baris komite)* | `T_GENERAL_KOMITE` | **tanpa kolom** — shared PK (`ID` = `ID`) | **1:1** | **di Go** | PK saja | `[keputusan work owner 2026-09-18]` — ⛔ **tidak terbukti dari korpus** |
| **3** | `T_GENERAL_KOMITE` | `T_KOMITE_KOMITELIST` | `DATA_KOMITE_ID` | **1:N** | **CASCADE** | **biasa** pada `DATA_KOMITE_ID` | ✅ **terbukti korpus** sebagai daftar bersarang |
| **4** | `T_CLAIM_ADJUSTMENT` | `T_GENERAL_KOMITE` | `ADJUSTMENT_ID` | **1:1** | **di Go** | ⭐ **UNIK** | ✅ **terbukti korpus** — menggantikan kunci posisional |
| **5** | `T_CLAIM_ADJUSTMENT` | `T_WORK_CLAIM` *(baris komite)* | `KOMITE_ID` | **1:1** | **di Go** *(penunjuk)* | ⭐ **UNIK**, nullable | ⛔ **TIDAK terbukti dari korpus** — lihat B2-5 |
| **6** ⭐ | `T_WORK_CLAIM` *(baris komite)* | `HISTORYAKSEPTASIPEGA` | `ID_KOMITE` | **1:N** | **JANGAN cascade** | **biasa** pada `ID_KOMITE` | ⭐ **RELASI BARU** — terbukti korpus, belum pernah terdaftar |

⚠️ **Arah relasi 5 saya balik** dari yang tertulis di brief. Brief menulis
`T_WORK_CLAIM -> T_CLAIM_ADJUSTMENT KOMITE_ID`; berkas Life relasi **11** menulis
`T_CLAIMLF_ADJUSTMENT -> T_WORK_CLAIM`. Yang kedua benar secara bentuk: **kolom `KOMITE_ID` duduk
di tabel penyesuaian**, menunjuk ke baris komite. Induk = pemilik nilai yang ditunjuk. Dilaporkan,
tidak dibetulkan diam-diam.

### B2 · B3 · B4 · B5 — satu per satu

#### Relasi 1 — `COVER_KEY` : baris klaim → baris komite

| | |
| --- | --- |
| **Bukti** | `[keputusan work owner 2026-09-18]`. ⛔ **Tidak ada kolom `COVER_KEY` di korpus.** Di Pega hubungan induk-anak ditegakkan mekanisme *cover/covered* bawaan, bukan kolom. Jejak yang ada: `Claim Prop/Activity/AddKomiteTreatyChild_ACT.xml` (`ASM-FW-GCNMFW-DATA-ADJUSTMENT!ADDKOMITETREATYCHILD_ACT`) langkah **12** menyalin halaman kasus anak, dan `Komite Claim Prop/Activity/KomitePostAdjustment.xml` (`ASM-FW-GCNMFW-WORK-KOMITETREATY!KOMITEPOSTADJUSTMENT`) langkah **4** membuka kasus induk lewat *handle* |
| **DIISI** | **AddKomiteChild langkah 12** — saat kasus anak komite dibuat |
| **DIBACA** | **KomitePost langkah 4** (`Obj-Open-By-Handle`), bergerbang `pyWorkPage.TransferType=="2"` |
| **Bila kosong** | ⛔ **Tidak terbaca dari korpus** — langkah 4 tidak punya jalur kegagalan; gerbangnya menguji `TransferType`, bukan keberadaan induk. `[terbuka]` |
| **ON DELETE** | **di Go** — menghapus baris klaim harus menghapus baris komitenya, tetapi keduanya tabel yang sama sehingga `CASCADE` self-referencing berisiko; ditegakkan eksplisit |
| **Index** | **biasa** pada `COVER_KEY` — dipakai mencari seluruh baris komite milik satu klaim |

#### Relasi 2 — shared PK : baris komite → `T_GENERAL_KOMITE`

| | |
| --- | --- |
| **Bukti** | `[keputusan work owner 2026-09-18]`. ⛔ **Tidak ada bukti korpus** — Pega tidak punya tabel `T_GENERAL_KOMITE`; datanya hidup sebagai properti pada halaman kasus (`pyWorkPage.KomiteLoop`, `.KomiteCount`, `.AcceptStatus`). Relasi ini **lahir dari keputusan**, bukan dari korpus |
| **DIISI** | tidak berlaku — tidak ada kolom penyambung |
| **DIBACA** | tidak berlaku |
| **Bila kosong** | tidak mungkin kosong — `ID`-nya PK kedua tabel |
| **ON DELETE** | **di Go** — pasangan shared-PK wajib lahir dan mati bersama |
| **Index** | PK saja; **tidak perlu index tambahan** |

#### Relasi 3 — `DATA_KOMITE_ID` : header komite → daftar penyetuju

| | |
| --- | --- |
| **Bukti** | ✅ **terbukti korpus** sebagai **daftar bersarang** di dalam kasus komite. `Claim Prop/Activity/AddKomiteTreatyChild_ACT.xml` langkah **22.1** menambah baris penyetuju satu per satu (`KomiteList(<APPEND>)`), langkah **23** menghitung panjangnya. Nama kolomnya `[keputusan work owner]` — di Pega tidak ada kolom penyambung karena daftarnya bersarang |
| **DIISI** | **AddKomiteChild langkah 22.1** — `KomiteID` · `KomiteAproval` = `0` · `KomiteEmail` · `IDKomite`; dan langkah **23** mengisi `KomiteLoop` = panjang daftar, `KomiteCount` = `1` |
| **DIBACA** | **KomitePost langkah 6 · 7 · 8 · 9 · 26 · 26.1** — menulis keputusan penyetuju dan menyusun teks kronologi |
| **Bila kosong** | ⛔ **Tidak terbaca dari korpus.** Langkah 23 mengisi `KomiteLoop` dari panjang daftar tanpa memeriksa apakah daftarnya terisi; bila roster kosong, `KomiteLoop` = 0. **Tidak ada pemeriksaan.** `[terbuka]` |
| **ON DELETE** | ⭐ **CASCADE** — baris penyetuju **tidak punya arti** tanpa header kasusnya; di Pega ia bagian dari halaman kasus yang sama |
| **Index** | **biasa** pada `DATA_KOMITE_ID` |

#### Relasi 4 — `ADJUSTMENT_ID` : baris penyesuaian → header komite

| | |
| --- | --- |
| **Bukti** | ✅ **terbukti korpus**, tetapi **dalam bentuk lain**. Di Pega penunjuknya **posisional** — `Claim Prop/Activity/AddKomiteTreatyChild_ACT.xml` langkah **13** mengisi `IndexAdjustment` dari **nomor urut baris**. `ADJUSTMENT_ID` adalah **penggantinya** `[keputusan work owner 2026-09-18]` |
| **DIISI** | di Pega: **AddKomiteChild langkah 13**, saat kasus komite dibuat |
| **DIBACA** | **KomitePost langkah 3** menyalinnya ke variabel kerja; dipakai langkah **6 · 7 · 15 · 16.9 · 23 · 24 · 25 · 27** |
| **Bila kosong** | ⛔ **Tidak terbaca dari korpus** — kedelapan langkah itu menulis ke posisi tersebut **tanpa memeriksa lebih dulu** apakah barisnya ada. `[terbuka]` — ini butir 2 di berkas struktur |
| **ON DELETE** | **di Go** — A5: dua tabel tujuan menurut lini, jadi **tanpa `REFERENCES`**; basis data tidak dapat menegakkannya |
| **Index** | ⭐ **UNIK** — menegakkan **1:1**: satu baris penyesuaian **paling banyak satu** kasus komite yang hidup |

#### Relasi 5 — `KOMITE_ID` : ⛔ **TIDAK terbukti dari korpus**

| | |
| --- | --- |
| **Bukti** | ⛔ **Tidak ada.** Yang ada di korpus hanyalah **penanda boolean**, bukan penunjuk: `Claim Prop/Activity/AddKomiteTreatyChild_ACT.xml` langkah **1** menulis `.IsKomite := 1`, dan `KomitePostAdjustment` langkah **23** menulis `…AdjustmentList(…).IsKomite := 0`. **Tidak ada properti apa pun pada baris penyesuaian yang memuat kunci kasus komite.** Sudah disisir di `pyParamArray`, `pyStepsCallParams`, dan `pyStepsJavaSource` kedua rule — aturan **I2** dipenuhi |
| **DIISI** | ⛔ tidak ada di korpus. Sebagai kolom baru ia `[keputusan work owner]`, tanggal **belum tercatat** |
| **DIBACA** | ⛔ tidak ada di korpus |
| **Bila kosong** | tidak berlaku — nullable menurut berkas Life |
| **ON DELETE** | **di Go** — ia penunjuk balik, bukan kepemilikan |
| **Index** | ⭐ **UNIK**, nullable — menegakkan **1:1** arah balik relasi 4 |

> ### ⚠️ Relasi 4 dan 5 adalah **lingkar**
>
> `T_GENERAL_KOMITE.ADJUSTMENT_ID` menunjuk baris penyesuaian, dan
> `T_CLAIM_ADJUSTMENT.KOMITE_ID` menunjuk balik ke baris komite. Keduanya ber-index **UNIK**.
> ⚠️ **Tidak ada apa pun yang menjaga keduanya konsisten** — basis data tidak bisa, karena relasi 4
> tanpa `REFERENCES`. `[terbuka]` siapa yang menjaga lingkar itu di Go.

#### Relasi 6 ⭐ — `ID_KOMITE` : baris komite → `HISTORYAKSEPTASIPEGA` — **RELASI BARU**

| | |
| --- | --- |
| **Bukti** | ✅ **terbukti korpus, bukti tiga bagian lengkap.** `Komite Claim Prop/RDBList/InsertHistoryAkseptasiPega_Sql.xml` (`ASM-FW-GISFW-INT-POLICYJSON!ASM!INSERTHISTORYAKSEPTASIPEGA_SQL`) menyisipkan ke `HISTORYAKSEPTASIPEGA` dengan kolom **`ID_PEGA` · `Tgl_Transfer` · `Status` · `Username` · `Workbasket` · `ID_KOMITE`**. Dipanggil `KomitePostAdjustment` langkah **33**, parameternya diisi langkah **32** |
| **DIISI** | **KomitePost langkah 32** — `ID_PEGA` ← kunci kasus **klaim** · `ID_KOMITE` ← kunci kasus **komite** · `Status` ← `ACCEPT`/`REJECT` · `Username` ← operator · `Workbasket` ← `"KLAIM"`. Dijalankan langkah **33**, **tanpa gerbang** |
| **DIBACA** | ⛔ **nol pembaca di kedua modul** — disisir seluruh `RDBList`, `ReportDefinition`, `Activity`, dan `Section` di `Komite Claim Prop` dan `Claim Prop`. Ia **tulis-saja** di dalam lingkup ini |
| **Bila kosong** | ⛔ tidak terbaca — `INSERT`-nya tidak punya jalur kegagalan, dan membawa **`COMMIT` sendiri** di dalam SQL |
| **ON DELETE** | ⭐ **JANGAN cascade.** Ia **jejak audit** — menghapus kasus komite **tidak boleh** menghapus riwayatnya |
| **Index** | **biasa** pada `ID_KOMITE` |

> ⚠️ **Kuncinya kunci internal Pega, bukan ID bisnis.** `ID_PEGA` dan `ID_KOMITE` diisi dari
> `pzInsKey` — pengenal internal Pega. Sesudah migrasi, **pengenal itu tidak ada lagi**.
> `[terbuka]` apa yang menggantikannya, dan apakah riwayat lama ikut dipetakan.

### Rekap §B

| | |
| --- | --- |
| Jumlah relasi | **6** |
| Relasi **kelima yang belum terdaftar** | ⭐ **ADA** — relasi **6**, `HISTORYAKSEPTASIPEGA.ID_KOMITE` |
| Relasi yang **TIDAK terbukti dari korpus** | **3** — relasi **1** (`COVER_KEY`), **2** (shared PK), **5** (`KOMITE_ID`) |
| Bukti **tiga bagian lengkap** | **3 dari 6** — relasi **3 · 4 · 6** |

---

## §C — Relasi yang bukan kunci tamu

### C1 — Tujuh belas properti baris penyesuaian yang DIBACA

`[terverifikasi]` Dihitung ulang dari `Komite Claim Prop/Section/ShowTransfer.xml`
(`ASM-FW-GCNMFW-WORK-KOMITETREATY!SHOWTRANSFER`): **19 sel** membaca `.Adjustment.*` —
**17 hanya-baca** + **2 yang disunting penyetuju**. Cocok dengan pembagian di berkas struktur.

| # | Properti korpus | Kolom/arti | Dibaca di | Gerbangnya |
| --- | --- | --- | --- | --- |
| 1 | `.Adjustment.AdjustmentValue` | nilai penyesuaian | ShowTransfer, sel nilai | jalur penyesuaian `.TransferType = 2` |
| 2 | `.Adjustment.GrossAdjustment` | nilai kotor penyesuaian | ShowTransfer | jalur penyesuaian |
| 3 | `.Adjustment.GrossValue` | nilai kotor | ShowTransfer | jalur penyesuaian |
| 4 | `.Adjustment.Currency` | mata uang | ShowTransfer | jalur penyesuaian |
| 5 | `.Adjustment.Type` | jenis penyesuaian | ShowTransfer | jalur penyesuaian |
| 6 | `.Adjustment.PersenRNM` | persen RNM | ShowTransfer | jalur penyesuaian |
| 7 | `.Adjustment.NameOfBank` | nama bank | ShowTransfer | jenis pembayaran `1`/`2`/`5` |
| 8 | `.Adjustment.BranchOfBank` | cabang bank | ShowTransfer | jenis pembayaran `1`/`2`/`5` |
| 9 | `.Adjustment.NoAccount` | nomor rekening | ShowTransfer | jenis pembayaran `1`/`2`/`5` |
| 10 | `.Adjustment.SwiftCode` | kode swift | ShowTransfer | ⭐ **`.Adjustment.SwiftCode != ''`** — gerbang isi-sendiri |
| 11 | `.Adjustment.PayableTo` | dibayarkan kepada | ShowTransfer | jenis pembayaran |
| 12 | `.Adjustment.Payable` | dapat dibayar | ShowTransfer | jenis pembayaran |
| 13 | `.Adjustment.IndividualRiskType` | jenis risiko perorangan | ShowTransfer | katastrofa / *Big Claim* |
| 14 | `.Adjustment.IndividualRiskPercentage` | persen risiko perorangan | ShowTransfer | katastrofa / *Big Claim* |
| 15 | `.Adjustment.IndividualRiskValue` | nilai risiko perorangan | ShowTransfer | katastrofa / *Big Claim* |
| 16 | `.Adjustment.IndividualRiskRNM` | nilai risiko perorangan bagian RNM | ShowTransfer | katastrofa / *Big Claim* |
| 17 | `.Adjustment.ProposeAdjustmentValue` | nilai penyesuaian yang diusulkan | ShowTransfer | jalur penyesuaian |

Dua yang **disunting**, bukan dibaca: `.Adjustment.IsProposeClose` · `.Adjustment.IsPropReserved`.

> ### ⚠️ Ralat kecil pada daftar di berkas struktur — satu tukar
>
> Daftar 17 di `STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md` menyebut **"kode mata uang"**
> (`.Adjustment.CurrencyID`), yang **tidak muncul sebagai sel** di ShowTransfer; dan **tidak
> menyebut** `.Adjustment.IndividualRiskRNM`, yang **muncul**. ⛔ Jumlahnya tetap **17**.
> Berkas struktur **tidak disunting**; selisihnya dicatat di sini saja.

#### Jalur bacanya sebagai rantai — **terbukti ulang**, dan kuncinya ternyata **DUA bagian**

```
1.  kasus komite  ->  kunci kasus klaim induk
       bukti: KomitePostAdjustment langkah 4 — Obj-Open-By-Handle,
              bergerbang pyWorkPage.TransferType=="2"

2a. kasus komite  ->  POSISI baris penyesuaian
       bukti: KomitePostAdjustment langkah 3 — Local.AdjusmentID := .IndexAdjustment
              asal:  AddKomiteTreatyChild_ACT langkah 13 — IndexAdjustment := .pxListSubscript

2b. kasus komite  ->  POSISI objek pertanggungan            <-- ⭐ BELUM PERNAH TERCATAT
       bukti: KomitePostAdjustment langkah 5 — Local.CoverageID := .Adjustment.CoverageSubscript
              asal:  AddKomiteTreatyChild_ACT langkah 13 — CoverageSubscript := .IndexCoverage

3.  baris penyesuaian = daftar penyesuaian klaim itu, pada posisi tersebut
       bukti: KomitePostAdjustment langkah 6 · 7 · 15 · 16.9 · 23 · 24 · 25 · 27
```

### C2 — ⭐ KUNCINYA POSISIONAL — terbukti ulang, dan **ADA PENUNJUK POSISIONAL KEDUA**

> ## ⭐⭐ TEMUAN BESAR — dilaporkan keras-keras, bukan dikubur
>
> `[terverifikasi]` **Penghubungnya memang nomor urut baris, bukan pengenal baris** — dan
> buktinya kini **lebih tajam dari sebelumnya**:
>
> ```
> AddKomiteTreatyChild_ACT langkah 13:
>     childPageKomite.IndexAdjustment             :=  .pxListSubscript
>     childPageKomite.Adjustment.CoverageSubscript :=  .IndexCoverage
> ```
>
> `.pxListSubscript` adalah **nomor urut baris** yang sedang diulang — properti bawaan Pega,
> bukan pengenal tersimpan. ⛔ **Tidak ada keraguan lagi.**
>
> ### ⚠️ Dan penunjuk posisionalnya **DUA**, bukan satu
>
> Selain posisi **baris penyesuaian**, kasus komite juga membawa posisi **objek
> pertanggungan** (`CoverageSubscript`), dibaca `KomitePostAdjustment` langkah **5**.
> **Penunjuk kedua ini belum pernah tercatat di berkas mana pun**, dan ia punya kerapuhan
> yang sama persis: menyisipkan atau menghapus objek pertanggungan di klaim induk membuat
> penunjuk komite menunjuk objek yang salah.
>
> ⭐ **Akibatnya untuk rancangan:** `ADJUSTMENT_ID` menggantikan penunjuk **pertama**.
> ⛔ **Tidak ada kolom apa pun yang menggantikan penunjuk KEDUA.** `[terbuka]` baru.

**Pengenal baris yang terlewat: ⛔ TIDAK ADA.** Disisir `pyParamArray`, `pyStepsCallParams`, dan
`pyStepsJavaSource` pada `AddKomiteTreatyChild_ACT` dan `KomitePostAdjustment`, serta seluruh sel
`ShowTransfer` — **tidak ketemu di medan-medan itu** satu pun properti pada baris penyesuaian yang
berfungsi sebagai pengenal baris. Sesuai aturan **I2**, kalimat ini dibatasi pada medan yang
benar-benar disisir.

### C3 — Lintas modul lewat activity: ⭐ **DELAPAN penulisan, bukan satu**

`[terverifikasi]` **Sensus, bukan contoh.** Seluruh `pyParamArray` `KomitePostAdjustment`
(**51 langkah**) disisir untuk sasaran berawalan `TempOpenPage` — yaitu **kasus klaim induk** yang
dibuka langkah 4.

**Hasilnya dua golongan:**

| Golongan | Jumlah penulisan | Sifatnya |
| --- | --- | --- |
| **A.** menulis ke **baris penyesuaian** `…ClaimData.AdjustmentList(…)` | **17** | jalur yang sudah terdokumentasi |
| **B.** menulis **langsung ke kasus klaim** | ⭐ **8** | **belum pernah disensus** |

#### Golongan B — kedelapan penulisan lintas modul tingkat klaim

| Langkah | Sasaran di kasus klaim | Sumber di kasus komite | Gerbang |
| --- | --- | --- | --- |
| **11** | `.ClaimData.IsCloseFile` | `.Adjustment.IsProposeClose` | ⭐ **TANPA GERBANG** |
| **11** | `.ClaimData.IsReservedClaim` | `.Adjustment.IsPropReserved` | ⭐ **TANPA GERBANG** |
| **14** | `.IsAnyAcceptation` | `1` | penyetuju terakhir **dan** hasil setuju |
| **22** | `.AktifButton` | `0` | penyetuju terakhir |
| **24** | `.ClaimData.IsSubjectivity` | `.IsSubjectivity` | penyetuju terakhir **dan** setuju |
| **25** | `.AktifButton` | `0` | hasil tolak |
| **26.1** | `.ClaimData.FlagOnGoingCommitte` | `"Adjustment"` | penyetuju belum memutus |
| **27** | `.ClaimData.FlagOnGoingCommitte` | `"Adjustment"` | penyetuju terakhir |

**Kedua pasangan nama yang diminta brief** (langkah 11, tanpa gerbang):

```
kasus komite . IsProposeClose   ->  kasus klaim . IsCloseFile
kasus komite . IsPropReserved   ->  kasus klaim . IsReservedClaim
```

⭐ Inilah asal nama kolom **`CLOSE_FILE`** dan **`RESERVED_CLAIM`** (A2) — keduanya mengikuti
ejaan **sisi klaim**, bukan sisi komite.

> ⚠️ **Yang berubah dari catatan lama:** catatan sebelumnya menyebut langkah **11** sebagai
> penulisan lintas modul. ⭐ **Ia bukan satu-satunya — ia satu dari delapan, dan satu-satunya yang
> TANPA GERBANG.** Enam sasaran berbeda ditulis seluruhnya. Ditulis sebagai **relasi data**, bukan
> kunci tamu: nilai yang sama hidup di dua tempat dengan dua nama.

⚠️ `[terbuka]` **baru** — enam sasaran golongan B selain kedua penanda usul
(`IsAnyAcceptation` · `AktifButton` · `ClaimData.IsSubjectivity` · `ClaimData.FlagOnGoingCommitte`)
**belum punya rumah kolom** di berkas struktur mana pun. Apakah keempatnya kolom sisi klaim,
adalah urusan modul Claim Prop — **tidak diisi sendiri di sini**.

#### Arah sebaliknya — Claim Prop menulis ke kasus komite

`[terverifikasi]` `AddKomiteTreatyChild_ACT` langkah **13 · 14 · 15 · 16 · 22.1 · 23** menyalin
⭐ **49 nilai** ke kasus komite saat ia dibuat — **24** di antaranya `Adjustment.*`.

> ⚠️ **Selisih yang wajib dicatat:** berkas struktur menetapkan **17 dibaca, tidak disalin** —
> itu **rancangan untuk Go**. **Di Pega ketujuh belasnya MEMANG DISALIN**, bersama **enam
> properti lagi** yang tidak termasuk daftar 17: `KursIDR` · `EstimationValue` ·
> `ValueAdjustment` · `SpreadingAdjustment` · `SpreadingQuotaShare` · `IndividualRiskRNM`.
> ⛔ Bukan pertentangan — rancangan Go memang mengganti salinan dengan bacaan. Dicatat supaya
> selisihnya tidak dikira kelalaian.

### C4 — Dua tabel yang belum masuk pohon relasi

#### `HISTORYAKSEPTASIPEGA` — ⭐ **TERBACA PENUH**

| | |
| --- | --- |
| **Rule yang menyentuhnya** | `Komite Claim Prop/RDBList/InsertHistoryAkseptasiPega_Sql.xml` — `ASM-FW-GISFW-INT-POLICYJSON!ASM!INSERTHISTORYAKSEPTASIPEGA_SQL` |
| **Dipanggil** | `KomitePostAdjustment` langkah **33**; parameternya diisi langkah **32** |
| **Kolomnya** | `ID_PEGA` · `Tgl_Transfer` · `Status` · `Username` · `Workbasket` · `ID_KOMITE` |
| **Kunci ke data komite** | ⭐ **ADA — `ID_KOMITE`**, diisi dari kunci kasus komite |
| **Milik siapa** | ⚠️ kelasnya `ASM-FW-GISFW-INT-POLICYJSON` — **antarmuka**, bukan kelas kerja modul ini. `[terbuka]` apakah ia milik migrasi ini atau milik sistem lain |

⛔ **Kolomnya tidak ditebak** — keenamnya terbaca langsung dari pernyataan SQL rule itu.
⛔ **Nol DDL ditulis.**

#### `T_STORAGE_IMAGE` — ⚠️ **menyentuh modul, tetapi kuncinya ke komite TIDAK TERBACA**

| | |
| --- | --- |
| **Rule yang menyentuhnya** | **12 rule** di `Komite Claim Prop`, dan **12 rule berpasangan** di `Claim Prop` — `InsertGoogleStorage_Act` · `GetUrlGoogleStorage_Act` · `GetBase64Attachment` · `InsertDocument_Act` · `ServiceGoogle` · `GetMimeType` · `GenerateImageID_SQL` · `GetAppName_SQL` · `GetLinkStorage_SQL` · `GetTokenStorage_SQL` · `Insert_T_Storage_SQL` · `Update_T_Storage_SQL` |
| **Kunci ke data komite** | ⛔ **tidak terbaca dari korpus** |
| **Milik siapa** | ⚠️ kelasnya `ASM-FW-GISFW-INT-T_STORAGE_IMAGE` — **antarmuka penyimpanan berkas bersama**, dipakai banyak modul. `[terbuka]` |

⛔ **BERHENTI di sini**, sesuai perintah brief. Kolomnya tidak ditebak, DDL-nya tidak ditulis.

### C5 — `T_VIEW_SUGGEST` dan `DOCUMENT_CLAIM`

| Tabel | Menyentuh sisi komite lini PROP? | Bukti |
| --- | --- | --- |
| **`T_VIEW_SUGGEST`** | ⛔ **TIDAK** | disisir seluruh berkas `Komite Claim Prop` **dan** `Claim Prop`, penyaringan **tidak peka huruf besar-kecil** (aturan I1) — **nol kemunculan di kedua modul** |
| **`DOCUMENT_CLAIM`** | ✅ **YA** | **5 rule** di `Komite Claim Prop`: `GetBase64Attachment` · `InsertDocument_Act` (`ASM-FW-GCNMFW-INT-DOCUMENT_CLAIM!INSERTDOCUMENT_ACT`) · **`KomitePost_Close`** · **`KomitePost_Reject`** · `PrintFileAcceptance_TKMT`. ⭐ Dua di antaranya rule **penutupan dan penolakan komite** |

⚠️ `[terbuka]` **Kunci `DOCUMENT_CLAIM` ke data komite tidak ditelusuri** di tugas ini — brief
hanya meminta *ya atau tidak*. Jawabannya **ya**, dan penelusuran kuncinya **belum dikerjakan**.

---

## §D — Periksa silang dengan Komite Claim Life

### D1 — Bentuk relasi tabel lintas-lini

| Relasi | Prop *(berkas ini)* | Life *(`STRUKTUR-TABEL-KOMITE-CLAIM-LIFE.md`)* | Sama? |
| --- | --- | --- | --- |
| baris komite → `T_GENERAL_KOMITE` | shared PK, 1:1, di Go | relasi **8** — shared PK, 1:1, di Go | ✅ **SAMA** |
| `T_GENERAL_KOMITE` → `T_KOMITE_KOMITELIST` | `DATA_KOMITE_ID`, 1:N, **CASCADE** | relasi **9** — `DATA_KOMITE_ID`, 1:N, **CASCADE** | ✅ **SAMA** |
| tabel penyesuaian → `T_GENERAL_KOMITE` | `ADJUSTMENT_ID`, 1:1, NOT NULL, UNIK, tanpa `REFERENCES` | relasi **10** — identik, dua tabel tujuan menurut `LINI` | ✅ **SAMA** |
| tabel penyesuaian → baris komite | `KOMITE_ID`, 1:1, nullable, UNIK | relasi **11** — identik | ✅ **SAMA** |
| baris klaim → baris komite (`COVER_KEY`) | 1:N, di Go | ⚠️ **tidak terdaftar** sebagai relasi bernomor di berkas Life | ⚠️ **BEDA** |
| baris komite → `HISTORYAKSEPTASIPEGA` | ⭐ relasi **6** baru | ⛔ **tidak ada** di berkas Life | ⚠️ **BEDA** |

**Mana yang benar untuk kedua selisih itu:**

1. **`COVER_KEY`** — bentuk **Prop** yang lebih lengkap. Berkas Life **memang memuat** aturan
   `COVER_KEY` di bagian `T_WORK_CLAIM`-nya, hanya **tidak mengangkatnya** ke tabel relasi
   bernomor. ⛔ Bukan pertentangan, melainkan **kelengkapan pencatatan**.
2. **`HISTORYAKSEPTASIPEGA`** — ⚠️ **belum dapat diputuskan.** Rule penulisnya ada di
   `Komite Claim Prop`; apakah lini Life menulis ke tabel yang sama **tidak diperiksa** di tugas
   ini karena modul Life **di luar lingkup**. `[terbuka]`.

### D2 — `T_GENERAL_KOMITE`: **5 kolom** di Life, **7 kolom** di Prop

| | Life | Prop |
| --- | --- | --- |
| `ID` · `ADJUSTMENT_ID` · `KOMITE_LOOP` · `KOMITE_COUNT` · `ACCEPT_STATUS` | ✅ | ✅ |
| **`CLOSE_FILE`** | ⛔ **tidak ada** | ✅ |
| **`RESERVED_CLAIM`** | ⛔ **tidak ada** | ✅ |

⛔ **Berkas Life TIDAK disunting.** Dicatat sebagai **pekerjaan tertunda**: karena
`T_GENERAL_KOMITE` **lintas-lini** (A3), kedua kolom itu **ada pada tabel yang sama** yang dipakai
Life — jadi berkas Life perlu diselaraskan, bukan tabelnya yang berbeda.

⚠️ **Selisih kedua yang ikut ketemu:** `DATA_KOMITE_ID` tertulis **NOT NULL** di berkas Prop dan
**nullable** di berkas Life. Untuk satu kolom yang sama pada satu tabel lintas-lini, **hanya satu
yang benar**. ⛔ Tidak saya putuskan — `[terbuka]`.

### D3 — Kolom lintas-lini yang hanya terisi pada satu lini

| Kolom | Terisi di PROP? | Terisi di lini lain? | Maka nullable? |
| --- | --- | --- | --- |
| **`CLOSE_FILE`** | ✅ ya — `ShowTransfer`, disunting penyetuju | ⛔ **tidak terbaca** untuk Life | ⭐ **WAJIB nullable** |
| **`RESERVED_CLAIM`** | ✅ ya — `ShowTransfer` | ⛔ **tidak terbaca** untuk Life | ⭐ **WAJIB nullable** |
| `ACCEPT_STATUS` | ✅ | ✅ | nullable *(sudah)* |
| `KOMITE_LOOP` · `KOMITE_COUNT` | ✅ | ✅ | nullable *(sudah)* |
| `ADJUSTMENT_ID` | ✅ | ✅ | **NOT NULL** di kedua lini |
| `DATA_KOMITE_ID` | ✅ | ✅ | ⚠️ **berselisih** — lihat D2 |

⭐ **Dua kolom wajib nullable** karena hanya terisi pada lini PROP.

---

## §E — Apa lagi

### E1 — Yang masih kurang sebelum relasi tabel komite lengkap

| # | Yang kurang | Kenapa menahan |
| --- | --- | --- |
| **1** | **Rumah bagi penunjuk posisional KEDUA** (`CoverageSubscript`) | `ADJUSTMENT_ID` hanya mengganti penunjuk pertama; yang kedua **tidak punya pengganti sama sekali** |
| **2** | **Selisih `DATA_KOMITE_ID` NOT NULL vs nullable** antara berkas Prop dan Life | satu kolom, satu tabel lintas-lini, dua tulisan berbeda |
| **3** | **Rumah bagi enam sasaran golongan B** yang ditulis lintas modul | empat properti kasus klaim ditulis modul komite tanpa punya kolom |
| **4** | **Pemetaan `ID_PEGA`/`ID_KOMITE`** pada `HISTORYAKSEPTASIPEGA` sesudah `pzInsKey` hilang | jejak audit putus bila tidak dipetakan |
| **5** | **Kunci `DOCUMENT_CLAIM` ke data komite** | terbukti menyentuh sisi komite, kuncinya belum ditelusuri |
| **6** | **Penjaga lingkar relasi 4 ↔ 5** di Go | basis data tidak bisa menjaganya |
| **7** | **Penyelarasan berkas Life** — 5 kolom vs 7, dan `HISTORYAKSEPTASIPEGA` | berkas Life belum memuat keduanya |

### E2 — Relasi yang PALING RAWAN salah

⛔ **Ditunjuk, tidak diperbaiki.**

> ⚠️ **Paling rawan: relasi 4 — `ADJUSTMENT_ID` ber-index UNIK, kardinalitas 1:1.**

Keunikannya menyatakan: **satu baris penyesuaian paling banyak satu kasus komite, selamanya.**
Yang saya **tidak** buktikan dari korpus adalah bahwa baris penyesuaian **tidak pernah dikirim ke
komite dua kali**. Bukti tandingan justru ada:

- `KomitePostAdjustment` langkah **23** menulis `…AdjustmentList(…).IsKomite := 0` — **mengembalikan
  penanda ke nol** sesudah komite selesai. Sebuah penanda yang dikembalikan ke nol biasanya berarti
  **boleh dipakai lagi**.
- `AddKomiteTreatyChild_ACT` langkah **1** menulis `.IsKomite := 1` **tanpa memeriksa** apakah ia
  sudah `1`.

⚠️ Bila satu baris penyesuaian **boleh** masuk komite lebih dari sekali, index **UNIK** itu akan
**menolak penyerahan kedua** — dan itu perubahan perilaku yang tidak diminta siapa pun.

**Paling rawan kedua:** relasi **3** ber-`ON DELETE CASCADE`. Ia benar secara bentuk, tetapi
`T_KOMITE_KOMITELIST` memuat **keputusan penyetuju berikut tanggalnya** — data yang berwatak
**jejak**, dan jejak biasanya tidak ikut terhapus.

### E3 — Pertanyaan BARU untuk work owner — **2**

#### Pertanyaan 1 — Satu baris penyesuaian boleh masuk komite lebih dari sekali?

**Apa yang ditanyakan.** Rancangan sekarang memasang index **UNIK** pada `ADJUSTMENT_ID`, artinya
satu baris penyesuaian hanya boleh punya satu kasus komite selamanya. Tetapi di Pega penanda
`IsKomite` pada baris itu **dikembalikan ke nol** setelah komite selesai, dan pembuatan kasus
komite **tidak memeriksa** apakah baris itu sudah pernah dikirim.

**Kenapa muncul.** Penanda yang dikembalikan ke nol biasanya berarti "boleh dipakai lagi". Kalau
memang boleh, index UNIK itu akan menolak penyerahan kedua.

**Bedanya kalau A atau B.** **A — hanya sekali:** index UNIK dipertahankan, dan penyerahan kedua
ditolak dengan pesan yang jelas. **B — boleh berulang:** index UNIK **dicabut**, `ADJUSTMENT_ID`
menjadi biasa, kardinalitas menjadi **1:N**, dan diperlukan penanda "kasus komite mana yang sedang
aktif".

**Apa yang tertahan.** Bentuk relasi **4** dan **5** sekaligus — keduanya ber-index UNIK dan
keduanya berubah bila jawabannya B.

#### Pertanyaan 2 — Penunjuk posisional KEDUA (objek pertanggungan) diganti apa?

**Apa yang ditanyakan.** Kasus komite membawa **dua** penunjuk posisional, bukan satu: posisi baris
penyesuaian, **dan** posisi objek pertanggungan (`CoverageSubscript`, dibaca `KomitePostAdjustment`
langkah 5). `ADJUSTMENT_ID` sudah ditetapkan menggantikan yang pertama. Yang kedua **belum punya
pengganti**.

**Kenapa muncul.** Ia baru ketahuan saat menelusuri jalur baca untuk berkas ini; belum pernah
tercatat di berkas mana pun.

**Bedanya kalau A atau B.** **A — diberi kolom pengganti** (penunjuk objek pertanggungan yang
stabil): kerapuhan posisionalnya hilang, tetapi `T_GENERAL_KOMITE` bertambah satu kolom dan
menjadi **delapan**. **B — dibawa apa adanya** sebagai angka posisi: konsisten dengan keputusan
*"penunjuk posisional dipindahkan apa adanya"* 2026-09-19, tetapi kerapuhannya **ikut terbawa**,
dan kali ini **tanpa penggantian sama sekali**.

**Apa yang tertahan.** Jumlah kolom `T_GENERAL_KOMITE` — A3 menetapkannya **tujuh**, dan jawaban A
menjadikannya delapan.

### E4 — ⭐ Yang seharusnya dikerjakan tetapi TIDAK diperintahkan blok ini

⛔ **Disebutkan, tidak dikerjakan.** ⛔ **Bukan saran urutan.**

| # | Butir |
| --- | --- |
| **1** | **Menelusuri kunci `DOCUMENT_CLAIM` ke data komite.** §C5 hanya meminta *ya/tidak*; jawabannya **ya**, dan dua rule yang menyentuhnya adalah **`KomitePost_Close`** dan **`KomitePost_Reject`** — jalur penutupan dan penolakan. Kuncinya belum ditelusuri |
| **2** | **Menyelaraskan `STRUKTUR-TABEL-KOMITE-CLAIM-LIFE.md`** — 5 kolom vs 7, dan selisih `DATA_KOMITE_ID` |
| **3** | **Meralat daftar 17 di berkas struktur Prop** — satu tukar: *"kode mata uang"* keluar, `IndividualRiskRNM` masuk |
| **4** | **Membuat `RELASI-TABEL-CLAIM-PROP.md`** — sisi klaim belum punya berkas relasi, padahal relasi 1 · 4 · 5 berujung di sana |
| **5** | **Memberi rumah pada enam sasaran golongan B** (§C3) — empat properti kasus klaim yang ditulis modul komite |
| **6** | **Memeriksa apakah lini lain juga menulis ke `HISTORYAKSEPTASIPEGA`** — menentukan apakah relasi 6 lintas-lini |
| **7** | **Memeriksa `KomitePost_Close` dan `KomitePost_Reject`** sebagai jalur penulis ketiga dan keempat — keduanya menaikkan `KomiteCount` *(langkah 16 dan 15)* tetapi **tidak disebut** di berkas struktur |

---

## Lampiran — dua selisih terhadap berkas struktur, dicatat tidak diperbaiki

⛔ Berkas struktur **tidak disunting**. Keduanya dicatat di sini saja.

| # | Berkas struktur menulis | Yang terbaca dari korpus |
| --- | --- | --- |
| **1** | `KOMITE_COUNT` *"dinaikkan **KomiteRouter langkah 5** dan **KomitePost langkah 40**"* | ⚠️ **KomiteRouter langkah 5 BER-REMARK** — mati, tidak menaikkan apa pun. **KomitePost langkah 40 hidup** dan benar. ⭐ Dua penaik lain **tidak disebut**: `KomitePost_Close` langkah **16** dan `KomitePost_Reject` langkah **15**. Jadi penaik yang hidup ada **tiga**, bukan dua |
| **2** | daftar 17 properti memuat *"kode mata uang"* | `.Adjustment.CurrencyID` **tidak muncul** sebagai sel `ShowTransfer`; yang muncul dan tidak terdaftar adalah `.Adjustment.IndividualRiskRNM`. Jumlah tetap **17** |

---

## Bukti berkas lain tidak disentuh

Sidik jari MD5 diambil sebelum penulisan berkas ini dan dibandingkan sesudahnya — hasilnya
dicatat di laporan ronde. **Satu-satunya berkas baru adalah `RELASI-TABEL-KOMITE-CLAIM-PROP.md`.**
