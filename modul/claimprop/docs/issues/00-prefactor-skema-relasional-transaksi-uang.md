# 00: PREFACTOR — skema relasional, batas transaksi penomoran, dan satu jalur uang

**Status:** dibangun 07-10-2026 — migrasi 520–531 dan 533 (532 dicabut, OQ-CP-17); procedure `PEGA_JSON_*` tidak dipanggil (RALAT di bawah); pemuat data lama `backend/alat/pemuatlama` diuji-kering di DEV — RALAT 07-10-2026 (semula `ready-for-agent`)
**AC yang dipegang (RALAT 07-10-2026, prompt §7 butir 4):** AC 132 (ADR-0009, migrasi penuh tanpa koeksistensi) — dipegang pemuat data lama `backend/alat/pemuatlama` (uji-kering DEV 07-10-2026: 2.451 kasus, 2.032 siap, 419 menunggu OQ-CP-18).
**Blocked by:** None (can start immediately)
**Menutup:** AC 1–28 · AC 105 · AC 106 · AC 107 · AC 123 · AC 124 · **AC 128** · **AC 129** *(35 AC)*

## Hasil & nilai pengguna

Hari ini **data klaim Claim Prop tidak dapat ditanyakan**. `[terverifikasi]` Peserta, baris
adjustment, estimasi, spreading, dan interest hidup sebagai *page list* di dalam work object Pega dan
dipersist lewat `Obj-Save` ke blob; yang menyentuh Oracle hanya **proyeksi** `OS_AKSEPTASI_KLAIM` dan
`JSON_KLAIM`. Tidak ada tabel yang bisa di-`SELECT` untuk menjawab "berapa outstanding per treaty".

Setelah tiket ini, Finance dapat menanyakan angka klaim lewat query biasa, nomor dokumen **tidak
terbakar** saat pembuatan gagal, dan nilai uang berhenti berubah karena urutan pemanggilan.

⚠️ **Tiket ini memblokir seluruh tiket lain.** Tidak ada perilaku lain yang dapat selesai
sebelumnya, karena semuanya menulis ke tabel yang belum ada.

## ⚠️ Dua hal baru — kolom yang dibuang, dan kaskade hapus yang ditolak

`[keputusan work owner]` 2026-09-19. **Keduanya perilaku BARU**, bukan paritas.

**① Kolom `FLAG_ON_GOING_COMMITTEE` pada `T_GENERAL_CLAIM` DIBUANG.** Ia **tidak dibuat** di skema
baru, dan **nilai lamanya tidak dimigrasikan**. Layar akseptasi **menurunkan sendiri** *"komite
sedang berjalan"* dari kasus komite yang berjalan. ⚠️ `[terverifikasi]` Kolom itu **ada di Pega dan
dipakai** — ditulis `Activity/AddKomiteTreatyChild_ACT.xml` langkah **3**, juga
`KomitePostAdjustment` langkah **26.1** dan **27**, dan **dibaca layar akseptasi** — karena itu ia
**penyimpangan sadar, bukan jalur mati yang dibuang**. Barisnya di
`STRUKTUR-TABEL-CLAIM-PROP.md` dan `RELASI-TABEL-CLAIM-PROP.md` **ditandai DIBUANG, tidak dihapus**,
supaya jejaknya tetap terbaca.

**② Jalur kaskade hapus klaim memeriksa lebih dulu.** Sebelum satu baris pun disentuh, jalur itu
memeriksa apakah klaim tersebut **PERNAH** punya kasus komite. **Bila pernah — seluruh penghapusan
DITOLAK**, dan kaskadenya **tidak pernah dijalankan**, walau komitenya **sudah selesai maupun
menolak**, dan **jalur komitenya mana pun — penyerahan baris penyesuaian maupun penutupan tanpa
pembayaran** `[keputusan work owner]` 2026-09-19. Pesannya **menyebut berapa kasus komite yang
pernah ada**. **Bila tidak pernah** — kaskade berjalan **penuh seperti biasa**, sampai tingkat
terdalam.

**③ Klaim lama hasil migrasi tidak dapat dihapus sama sekali**, apa pun riwayat komitenya
`[keputusan work owner — atas rekomendasi asisten]` 2026-09-19. Sebabnya: bukti *"pernah punya
kasus komite"* **diturunkan dari adanya kasus komite**, sedangkan untuk data lama keberadaan kasus
itu bergantung pada apa yang berhasil dimigrasikan — dan kolom `FLAG_ON_GOING_COMMITTEE` yang dulu
membuktikannya **sudah dibuang**. ⚠️ `[keputusan work owner]` 2026-09-19 — **migrasi tidak membaca
kolom itu sama sekali**; nilainya **tidak dipakai sebagai bukti** apa pun. ⛔ **Mudah dicabut** bila
yang dimaksud ternyata *semua klaim tanpa kecuali* — lihat ralat di **AC 128** `spec.md`.

⚠️ **Ini mempersempit AC 5**, yang semula tidak mengenal klaim yang pernah masuk komite. AC 5
**tetap bernomor 5** dan **ralatnya ada di `spec.md`** — baris **Menutup:** di atas **tidak berubah
karenanya**; yang bertambah hanya AC **128** dan **129**.

⚠️ `[terverifikasi]` **Pega tidak punya penjaga ini** — nol pemeriksaan, nol pesan, nol
penanganan. **Dibangun, bukan dimigrasikan.** Lingkup: **Claim Prop** dan **Claim Non Prop**
*(Non Prop menyusul)*; **Claim — Life tidak termasuk**.

## Area codebase

Skema Oracle baru · lapisan persistensi · tipe uang · batas transaksi · migrasi data warisan ·
model penanda kesalahan · konvensi penamaan kolom.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/DeleteAjsutment_Act.xml` | — | nol `RDB-List`; penghapusan hanya manipulasi page |
| `Activity/AddAdjustment_Act.xml` | 5 | baris adjustment dibentuk di clipboard, bukan di tabel |
| `Activity/SaveOutstanding_Act.xml` | — | persistensi lewat `Obj-Save` atas `pyWorkPage` |
| `RDBList/GetSequenceNumber_SQL.xml` | — | `COMMIT` di dalam SQL-nya sendiri |
| `RDBList/GenerateNoCLMTreatyIn.xml` | — | generator nomor ber-`COMMIT` |
| `RDBList/GenerateNoPLATreatyIn.xml` | — | generator nomor ber-`COMMIT` |
| `RDBList/GenerateNoDlaTreatyIn.xml` | — | generator nomor ber-`COMMIT` |
| `RDBList/GenerateNoTRTInTemp.xml` | — | generator nomor *temp*; menulis ke halaman salah ketik |
| `Activity/TryMakePLA_Act.xml` | 11 · 14 · 26 | sequence diambil, hasilnya hanya ke clipboard, step lanjut dapat melempar exception — **nol `Obj-Save` di rule ini** |
| `Activity/CurencyEstimation_Act.xml` | — | mata uang ditetapkan **per baris** |
| `Activity/CountEstimation_Act.xml` | — | bercabang eksplisit satu-mata-uang versus multi-mata-uang |
| `RDBList/CekLunasPremi_Sql.xml` | — | `arasapas.invoice` diprefiks, `detail_invoice` tidak |
| `RDBList/GetEmailCeding_SQL.xml` | — | schema keempat: fungsi di schema `gl` |
| `RDBList/CariHistoryClaim_SQL.xml` | — | enam alias, lima berbohong; tambalan koma-ke-titik atas nilai uang |
| `Activity/CountValueADJTreaty_Act.xml` | 2 · 14 · 16 · 19 · 20 | lima titik tulis penanda kesalahan dengan nilai berbeda |
| `Activity/ProteksiInitialandDate_Act.xml` | 4.1 · 5 | satu tulis dan satu baca berambang `>1` |
| `Activity/CheckDateDOL_Act.xml` | 25 · 26 | tulis `1` dan `0` — **`1` tidak pernah lolos `>1`** |
| `Section/AdjustmentDetail.xml` | `pyDisabledWhen` | pembaca kedua penanda kesalahan (`>1`) |
| `Activity/AddListClaimAmount.xml` | 1 · 7.1 | ⭐ **BARU 2026-09-19 (ronde 4)** — share ceding **masuk ke `.Value` di sini**; langkah 1 memakunya `"100"` bila kosong/nol |
| `Activity/CountListClaimAmountIDR.xml` | 1.1 | ⭐ **BARU** — jalur kedua yang memasukkan share ceding ke `.Value` |
| `Activity/CountPersen_act.xml` | 4.3.1 · 4.4 · **5.1** | ⭐ **BARU** — 4.3.1 menjumlah `.Value` yang **sudah** berisi share ceding; **5.1 mengalikannya SEKALI LAGI** |
| `Activity/CountPersen_act.xml` | 6.1 · 6.2.1 | ⭐ **BARU** — blok kembar yang **sudah diperbaiki**: rumus yang sama **tanpa** share ceding |
| `Activity/CountGrossAdjTreaty_Act.xml` | 3 · 4 · 5.2 | ⭐ **BARU** — tiga perkalian share ceding **sekali saja**; ⚠️ `[terbuka]` belum terbukti cacat |
| `Activity/CountValueADJTreaty_Act.xml` | 12 | ⭐ **BARU** — dua perkalian share ceding **sekali saja**; ⚠️ `[terbuka]` belum terbukti cacat |


## Keadaan sekarang di Oracle

⛔ **Bab ini memuat FAKTA, bukan rancangan.** Nama tabel baru, pembagian tabel, dan kolomnya
**belum diputuskan work owner**. Tiket ini tetap berbunyi *rancang skema baru*; yang ditulis di bawah
hanyalah keadaan yang selama ini tidak tertulis, supaya perancangnya tahu apa yang sedang ia ganti.

> ⚠️ **RALAT 07-10-2026** — kalimat lama: *"Nama tabel baru, pembagian tabel, dan kolomnya belum diputuskan work owner."*
> Sudah diputuskan: diagram sheet Claim Prop (18 dan 20-09-2026: sepuluh tabel `T_CLAIM_*` + lima tabel bersama,
> `T_WORK_CLAIM.LINI = 'PROP'`, ID `CLMP-`) dan keputusan 07-10-2026 (kolom PROP di `T_GENERAL_CLAIM`,
> `NUMBER(38,10)`). Bentuk mengikat = migrasi 520–531, 533 dan lampiran pengikat STRUKTUR.

### Kolom nyata kedua tabel klaim

`[data DBA]` DDL diberikan work owner **2026-09-18**.

```
POOLDATA.JSON_KLAIM          MNK_NO_KLAIM · DATA_JSON (CLOB) · IDPEGA · TGL_INPUT ·
                             TGL_KONVERSI · NOPOLIS · IDPROD · STS_KONVERSI

POOLDATA.OS_AKSEPTASI_KLAIM  CASEID · NOCLAIM · DATA_JSON (CLOB) · TANGGAL · NOPOLIS ·
                             STS_REJECT · STS_KONVERSI · STS_DLA · MASTERID · CLAIMOLD
```

**Seluruh data klaim ada di dalam satu kolom `DATA_JSON` bertipe CLOB.** Sisanya hanya kunci dan
penanda status.

`[terverifikasi]` Pembacaannya memakai **notasi titik ke dalam CLOB itu**:
`RDBList/DataOutstandingTreatyin.xml` mengambil `a.DATA_JSON.Currency` · `.GrossValue` · `.Value` ·
`.AcceptedNo` · `.Type` · `.KursValue` · `.TotalGross`; `RDBList/CariHistoryClaim_SQL.xml` mengambil
`a.DATA_JSON.DateOfLoss` · `.ClaimNo` · `.CauseOfLoss` · `.ClaimEstimate`.

**Inilah sebab kalimat pembuka tiket ini:** tidak ada kolom yang bisa di-`SELECT`, tidak ada tipe,
dan tidak ada index atas nilai klaim.

> ⚠️ **RALAT 2026-09-18 — ketujuh kolom SUDAH TERJAWAB; `[terbuka]` dicabut.** Teks lama:
> *"⚠️ `[terbuka]` **Arti enam kolom penanda belum dijelaskan DBA** dan tidak boleh disimpulkan dari
> namanya: `STS_KONVERSI` (di kedua tabel) · `STS_REJECT` · `STS_DLA` · `IDPROD` · `MASTERID` ·
> `CLAIMOLD`. `TGL_KONVERSI` ikut menggantung … tidak dapat dipetakan sebelum dijawab."*
> Jawabannya di dua sub-bab berikut.

### Arti ketujuh kolom penanda, dan tanda tangan procedure penulisnya

`[terverifikasi]` **Jendela sensus:** **329 berkas**, seluruh teks, dan penelusuran **diikat ke
halaman pemanggilnya** — `InputData.CARI<n>`, bukan `CARI<n>` telanjang. ⚠️ `TempKasir.CARI16` dan
`TempCFS.CARI20` milik halaman lain dan memberi **jawaban palsu**; jebakan lingkup yang sama seperti
`DataChronology.CARI1` pada jejak audit.

`[data DBA]` Tanda tangan `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM`, dipetakan **posisional** ke
pemanggilnya `RDBList/SaveOSClaim_SQL.xml`:

```
(NoClaim, PolisNO, PegaID, DataPega, STATUS_RJT, KONVERSI, STS_PLA, IDMASTER, OldClaim,
 ErrMsg OUT, StsSimpan OUT)
   ↑CARI29  ↑CARI3  ↑CARI2  ↑CARI1   ↑CARI10    ↑CARI16   ↑CARI17  ↑CARI18  ↑CARI20
```

| Kolom | Parameter | Diisi dari | Nilai yang benar-benar masuk | Tanda |
| --- | --- | --- | --- | --- |
| `STS_REJECT` | `STATUS_RJT` | `InputData.CARI10` ← `TempOSAkseptasi.Type` | **0** Save Outstanding · **2** Acceptation · **4** Close Claim | `[terverifikasi]` |
| `MASTERID` | `IDMASTER` | `InputData.CARI18` | ⚠️ **dua sumber** — `ClaimData.IDMaster` dari `SaveOutstanding_Act`, `TreatyInMaster.ID` dari `CloseClaimProp` | `[terverifikasi]` |
| `STS_KONVERSI` | `KONVERSI` | `InputData.CARI16` | ⚠️ **NOL PENULIS** di 329 berkas — masuk **kosong** | `[terverifikasi]` |
| `STS_DLA` | `STS_PLA` | `InputData.CARI17` | ⚠️ **NOL PENULIS** — masuk **kosong** | `[terverifikasi]` |
| `CLAIMOLD` | `OldClaim` | `InputData.CARI20` | ⚠️ **NOL PENULIS** — masuk **kosong** | `[terverifikasi]` |
| `IDPROD` | — | hardcode di `PEGA_JSON_KLAIM_PNC` | **1** | `[data DBA]` |
| `TGL_KONVERSI` | — | hardcode | **NULL** | `[data DBA]` |

**Empat temuan yang menempel:**

1. ⚠️ `[terverifikasi]` **Status konversi tidak dibaca dari kolomnya.**
   `RDBList/getStatusKonversi_SQL.xml` menghitung baris di **sistem tujuan**
   (`COUNT(1) … from reinsurance.trloss_detail_t WHERE NO_AKSEP = …`), bukan membaca
   `STS_KONVERSI`. Kolomnya bukan hanya kosong — ia juga **tidak dipakai**.
2. ⚠️ `[data DBA]` **Parameter bernama `STS_PLA` mengisi kolom `STS_DLA`.** PLA dan DLA **dua
   dokumen berbeda**; nama parameternya menyebut dokumen yang salah.
3. ⚠️ `[terverifikasi]` **`NOCLAIM` diisi dua kali di activity yang sama.**
   `Activity/SaveOutstanding_Act.xml` menulis `InputData.CARI29` dari `ClaimData.NoClaim` **lalu**
   dari `ClaimData.ClaimNo`. **Yang terakhir menang.**
4. `[terverifikasi]` **Dua jalur, satu tabel.** `SaveOSClaim_SQL` hanya dipanggil dari
   `Activity/SaveOutstanding_Act.xml`; `SaveDataToOsAkseptasiNP` hanya dari
   `Activity/CloseClaimProp.xml`.

### Tabel proyeksi diikuti APA ADANYA — bukan penyimpangan

`[keputusan work owner]` **2026-09-18.** Ketujuh penanda **ditiru apa adanya. Tidak ada yang
dianggap cacat:**

- **`STS_KONVERSI` · `STS_DLA` · `CLAIMOLD`** — kolomnya **tetap ada** dan **tetap ditulis kosong**.
  Nol penulis di Pega, nol penulis di Go. **Bukan cacat; jangan diisi, jangan dibuang.**
- **`IDPROD = 1`** dan **`TGL_KONVERSI = NULL`** — hardcode **ditiru**.
- Parameter **`STS_PLA` yang mengisi kolom `STS_DLA`** — nama **dibiarkan**.
- **`NOCLAIM` ditulis dua kali**, yang terakhir menang — **ditiru**.
- **`MASTERID` bersumber beda menurut jalurnya** — **ditiru**.

⚠️ **BUKAN penyimpangan sadar.** Tidak ada tanda penyimpangan dan tidak ada AC baru untuk butir-butir
ini. Rinciannya di spec **§1a** dan **§1b**.

⚠️ `[terbuka]` **Sisa untuk DBA hanya dua, keduanya TIDAK memblokir tiket ini:** domain nilai
`STS_KONVERSI` **bila kolomnya dipakai modul lain**, dan arti `IDPROD` di `POOLDATA.KODE_PRODUKSI` —
**bukan di sini**, karena di sini nilainya dipaku **1**.

### ⚠️ Batas `COMMIT` — procedure yang sama, `COMMIT` dijalankan aplikasi

`[keputusan work owner]` **2026-09-18.** Go memanggil **stored procedure yang sama** —
`POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM` · `POOLDATA.PEGA_JSON_KLAIM_PNC` ·
`POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` — dan menjalankan `COMMIT` **sendiri sesudah panggilan itu**.
**Logika penulisan tidak direplikasi di aplikasi.**

> ⚠️ **RALAT 07-10-2026** — kalimat lama: *"Go memanggil stored procedure yang sama … dan menjalankan COMMIT sendiri sesudah
> panggilan itu."* Diganti keputusan 07-10-2026 (prompt implementasi §3): procedure `PEGA_JSON_*` **tidak**
> dipanggil; isinya, dibaca dari `ALL_SOURCE` DEV, ditulis ulang sebagai SQL langsung di transaksi aplikasi
> (`repository.SisipOS`, `repository.SalinJSONKlaim`), hanya kolom selain `DATA_JSON` — `DATA_JSON` tidak
> diisi. `COMMIT` oleh aplikasi; nol `COMMIT` di teks SQL.

`[terverifikasi]` **Sejalan dengan AC 7, bukan melawannya.** AC 7 menuntut **nol `COMMIT` di dalam
teks SQL aplikasi**. Di Pega `COMMIT` tertanam di blok anonim rule — `BEGIN … ; COMMIT; END;` pada
`RDBList/InsertClaimPNC.xml`. Di Go ia **dikeluarkan aplikasi**, bukan ditanam di teks SQL. Objek yang
dipanggil tetap sama.

⚠️ `[data DBA]` **Akibat yang mengubah cara memanggil.** Ketiga procedure memasang
`EXCEPTION WHEN OTHERS THEN ROLLBACK` **di dalam dirinya sendiri**. **`ROLLBACK` telanjang di Oracle
membatalkan SELURUH transaksi** — termasuk apa pun yang aplikasi tulis **sebelumnya** dalam transaksi
yang sama — dan `ROLLBACK TO SAVEPOINT` **tidak menolong** karena yang dipasang bukan itu. Karena itu:

1. Panggilan procedure ini adalah **langkah terakhir sebelum `COMMIT`**.
2. **Tidak boleh ada pekerjaan penting yang masih menggantung** sebelumnya dalam transaksi yang sama.
3. Aplikasi **memeriksa `StsSimpan`** sesudah panggilan — **1 berhasil, 0 gagal**.
4. Pada **0**, aplikasi **tidak boleh melaporkan berhasil**: transaksinya **sudah dibatalkan oleh
   procedure**, sehingga `COMMIT` sesudahnya **tidak menyimpan apa pun**.

⚠️ **AC 6 belum menampung urutan ini** — bunyinya menuntut hasil (*satu aksi = satu transaksi; nol
keadaan separuh tersimpan*), bukan urutan. Usulan AC tambahan dilaporkan ke work owner; **tidak**
dituliskan sendiri ke spec.

### Inventaris objek Oracle yang disentuh Claim Prop

`[terverifikasi]` **Jendela sensus:** teks SQL di dalam tag `py*SQL*` pada **58 dari 58** berkas
`RDBList/`. Prosa keterangan rule **tidak** ikut dibaca — kata *"from"*/*"into"* dalam kalimat bahasa
Inggris akan terjaring sebagai nama tabel bila seluruh XML dibaca mentah. Pemanggilan stored procedure
**tidak** muncul di daftar ini karena bentuknya `BEGIN … END`, bukan `FROM`/`INTO`; keempatnya
didaftar tersendiri di bab berikut.

**36 baris ejaan — 34 objek Oracle yang berbeda**, dengan aksi dan rule yang menyentuhnya. Selisih
duanya: `TREATYINPRODUCTION` dan `TREATYBUSINESS` muncul **dua kali**, sekali telanjang dan sekali
ber-prefiks. **Yang menyiapkan pemetaan skema menyiapkan 34**, bukan 36.

> ⚠️ **RALAT 2026-09-18** — kalimat lama berbunyi *"**36 objek Oracle**, dengan aksi dan rule yang
> menyentuhnya"*. Angka 36 menghitung kedua objek itu dua kali. Ini **bukan** soal huruf
> besar-kecil — ini akibat langsung temuan **A** (telanjang versus ber-prefiks). Daftar tabelnya
> **tidak diubah**: tetap 36 baris, supaya tiap ejaan tetap terlihat beserta rule-nya.

| Objek | Aksi | Rule |
| --- | --- | --- |
| `AGENT` | baca | `GetAddressCeding` · `GetAddressTreatyIn` · `GetLeaderReport` |
| `ARASAPAS.INVOICE` | baca | `CekLunasPremi_Sql` |
| `BANKACCOUNT` | baca | `GetDataBankAccount2_sql` · `GetDataBankAccount_sql` · `GetDatabyClientName` |
| `BUSINESS` | baca | `GetOldIDBusiness` |
| `CITY` | baca | `BrowseRW_SQL` |
| `CURRENCY` | baca | `GetCurrency` |
| `JSON_KLAIM` | baca | `CariHistoryClaim_SQL` |
| `JSON_POLIS` | baca | `GetCaseIDNBTretyIn` · `GetPolicyData` · `GetYearofQuartal` |
| `MARKETINGOFFICER` | baca | `GetMObyNopol_SQL` |
| `M_CLIENT` | baca | `GetAddressCeding` · `GetAddressTreatyIn` |
| `M_TREATYGROUP` | **hapus** | `DeleteDataTreatyGroup_SQL` |
| `M_TREATY_IN` | baca | `GetLimitsTreatyIn_SQL` |
| `M_TREATY_IN_EDM` | baca | `GetLimitsTreatyIn_SQL` |
| `OS_AKSEPTASI_KLAIM` | baca | `DataOutstandingTreatyin` · `GetDataKlaimTreatyin` |
| `POLICYJSON` | baca | `GetDataNopolisTreatyin` |
| `POOLDATA.DIRECTTOKASIR_LOG` | baca+tulis | `GetStatusKasir_SQL` · `InsertLOGDirectKasir_SQL` |
| `POOLDATA.EMAILKOMITE` | baca | `GetLimitDirekturUtama_SQL` |
| `POOLDATA.KODE_PRODUKSI` | baca | `GetKodeProdNonLife_SQL` |
| `POOLDATA.MONITORING_KLAIM_LOG` | tulis | `InsertLogServiceClaim` |
| `POOLDATA.M_SITE_DATABASE` | baca | `GetIDConsultanAdj_SQL` |
| `POOLDATA.OPENPROTEKSI_EDM` | baca | `CekProteksiKlaim` |
| `POOLDATA.TANGGAL_CLOSING` | baca | `GETTanggalClosing_SQL` |
| `POOLDATA.TREATYBUSINESS` | baca | `GetTreatyInMasterProp_SQL` |
| `POOLDATA.TREATYINDETAIL` | baca | `GetTreatyInMasterProp_SQL` · `GetTreatyInMaster_SQL` |
| `POOLDATA.TREATYINDETAILEDM` | baca | `GetTreatyInMasterProp_SQL` · `GetTreatyInMaster_SQL` |
| `POOLDATA.TREATYINPRODUCTION` | baca | `CekPolicyNumber_SQL` |
| `POOLDATA.T_FOLDER_IMAGE` | baca | `GetAppName_SQL` |
| `PROPORTIONALARRG` | baca | `GetLimitPLATreatyin` |
| `REINSURANCE.TRLOSS_DETAIL_T` | baca | `getStatusKonversi_SQL` |
| `RW` | baca | `BrowseRW_SQL` |
| `TREATYBUSINESS` | baca | `GetTreatyGroupID` |
| `TREATYINPRODUCTION` | baca | `GetMObyNopol_SQL` · `GetNopolis_SQL` · `SetPolicyTreatyProp` |
| `TREATYREINSURER` | baca | `GetListRetro_Sql` |
| `TREATYYEAR` | baca | `TreatyYearTreatyin_SQL` |
| `T_STORAGE_IMAGE` | baca+tulis | `GetLinkStorage_SQL` · `Insert_T_Storage_SQL` · `Update_T_Storage_SQL` |
| `V_D_CAUSE_OF_LOSS_BUSINESS` | baca | `GetLBUID_SQL` |

⚠️ **Daftar ini adalah bukti untuk AC 8.** Dua objek muncul **dua kali** di daftar, karena objek yang
sama ditulis **lebih dari satu cara**:

| Objek | Telanjang di | Berprefiks di |
| --- | --- | --- |
| treaty-in production | **3 rule** — `GetMObyNopol_SQL` · `GetNopolis_SQL` · `SetPolicyTreatyProp` | 1 rule — `CekPolicyNumber_SQL` |
| treaty business | 1 rule — `GetTreatyGroupID` | 1 rule — `GetTreatyInMasterProp_SQL` |

`[terverifikasi]` **Penguatnya — huruf besar-kecilnya juga tidak konsisten, dan bukan menurut pola
apa pun.** Ejaan apa adanya dari teks SQL:

| Objek | Ejaan yang muncul | Rule |
| --- | --- | --- |
| treaty-in production | `TREATYINPRODUCTION` (telanjang, huruf besar) | `SetPolicyTreatyProp` |
| | `treatyinproduction` (telanjang, huruf kecil) | `GetMObyNopol_SQL` · `GetNopolis_SQL` |
| | `pooldata.TREATYINPRODUCTION` (prefiks huruf kecil) | `CekPolicyNumber_SQL` |
| treaty business | `treatybusiness` (telanjang, huruf kecil) | `GetTreatyGroupID` |
| | `POOLDATA.TREATYBUSINESS` (prefiks huruf besar) | `GetTreatyInMasterProp_SQL` |

Jadi `TREATYINPRODUCTION` muncul dalam **tiga** ejaan dan `TREATYBUSINESS` dalam **dua** — dan
**prefiks `pooldata` sendiri ditulis dua cara**, huruf kecil di satu rule dan huruf besar di rule
lain. Bukan "telanjang selalu huruf kecil": `SetPolicyTreatyProp` menulisnya telanjang **dan** huruf
besar.

⚠️ **Akibat untuk siapa pun yang menyensus ulang:** sensus nama tabel **wajib tidak peka huruf
besar-kecil**. Yang peka akan melaporkan lima tabel berbeda padahal hanya dua objek.

> ⚠️ **RALAT 2026-09-18 — klaim "lima objek" DICABUT.** Teks lama menghitung **beda huruf
> besar-kecil** sebagai temuan dan menjumlahkannya jadi *"Total LIMA objek dieja lebih dari satu
> cara"*, termasuk menyebut `jsoN_polis` sebagai hal yang tidak akan tertebak. Itu **menggelembungkan
> dua objek menjadi lima**. `[keputusan work owner]` Oracle **melipat identifier tanpa kutip menjadi
> huruf besar**, sehingga `treatyinproduction` dan `TREATYINPRODUCTION` adalah **objek yang sama** —
> beda huruf **bukan cacat dan tidak berakibat apa pun**. Pemisahan yang benar di bawah.

⚠️ **Yang berakibat hanya DUA objek di atas.** `[terverifikasi]` Pemeriksaan ulang `sensus.py --tabel`
atas teks SQL di **58 dari 58** berkas `RDBList/` memisahkan dua hal yang sempat tercampur:

> ⚠️ **RALAT 07-10-2026** — alat `sensus.py` sudah tidak ada dan kewajibannya dicabut (spec, bagian awal). Angka di bab ini fakta
> historis sapuan, bukan langkah kerja.

| | Objek | Akibat |
| --- | --- | --- |
| **A. Ditulis dua bentuk** — telanjang di satu rule, ber-prefiks di rule lain | `TREATYINPRODUCTION` · `TREATYBUSINESS` | ⚠️ **NYATA.** Nama telanjang bergantung pada **schema bawaan koneksi**; ia dapat menunjuk objek lain bila Go menyambung sebagai pengguna yang berbeda. **Inilah isi AC 8.** |
| **B. Hanya beda huruf besar-kecil** | `JSON_POLIS` · `DIRECTTOKASIR_LOG` · `T_STORAGE_IMAGE` | **NOL akibat.** Oracle **melipat** identifier tanpa kutip menjadi huruf besar, jadi `treatyinproduction` dan `TREATYINPRODUCTION` adalah objek yang sama. Bukan cacat, bukan temuan. |

`[keputusan work owner]` **2026-09-18 — beda huruf besar-kecil bukan masalah.** Di Go seluruh nama
ditulis satu cara, konsisten. Yang wajib diperbaiki hanya **prefiks schema**, sesuai AC 8.

⚠️ **Satu akibat yang tetap berlaku, tetapi hanya bagi pembaca korpus:** sensus nama tabel **wajib
tidak peka huruf besar-kecil**. Yang peka akan melaporkan objek yang sama sebagai beberapa tabel
berbeda. ⭐ **Aturannya: kunci sensus pada NAMA DASAR**, dengan **prefiks schema dipisahkan lebih
dulu** — barulah golongan A dan B dapat dibedakan.

> ⛔ **RALAT 2026-09-19** — kalimat lamanya **dikutip, tidak dihapus**: *"Pakai `sensus.py --tabel`,
> yang mengunci pada **nama dasar** dan memisahkan A dari B."* Alat itu **sudah tidak ada** dan
> kewajibannya **dicabut** `[keputusan work owner]` 2026-09-19; ⭐ **aturannya tetap berlaku** dan
> kini ditulis langsung, tidak lagi menumpang pada sebuah skrip. ⛔ **Nol AC berubah.**

**Empat schema** muncul: `pooldata` · `arasapas` · `reinsurance` (pada `TRLOSS_DETAIL_T`), dan
sisanya **telanjang**.

`[terverifikasi]` **Bukti untuk AC 2:** `OS_AKSEPTASI_KLAIM` hanya **dibaca** di seluruh `RDBList/`;
penulisannya lewat stored procedure — lihat bab berikut.

### Empat stored procedure, dan yang satu tidak ada

`[terverifikasi]` **Nol rule `RDBList` meng-`INSERT` langsung** ke `JSON_KLAIM` atau
`OS_AKSEPTASI_KLAIM`. Penulisannya lewat stored procedure:

| Procedure | Rule pemanggil | Ada di DB? |
| --- | --- | --- |
| `POOLDATA.PEGA_JSON_KLAIM_PNC` | `InsertClaimPNC` ← `Activity/InsertJsonClaimTreaty_act.xml` | ✅ ada `[data DBA]` |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM` | `SaveOSClaim_SQL` ← `Activity/SaveOutstanding_Act.xml` | ✅ ada `[data DBA]` |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTNP` | `SaveDataToOsAkseptasiNP` ← `Activity/CloseClaimProp.xml` | ✅ ada `[data DBA]` |
| `POOLDATA.PEGA_JSON_OS_AKSEP_KLAIMTRT` | `SaveOSKlaimTreaty_SQL` ← `Activity/SetOutstanding_Act.xml` step **4** · `Activity/UpdateTableOS.xml` step **7** | ❌ **TIDAK ADA** `[keputusan work owner]` |

⚠️ `[keputusan work owner]` **2026-09-18 — jalur KLAIMTRT tidak ada lagi dan TIDAK dimigrasikan.**
`RDBList/SaveOSKlaimTreaty_SQL.xml` beserta **kedua** titik pemanggilnya — `Activity/SetOutstanding_Act.xml`
step **4** dan `Activity/UpdateTableOS.xml` step **7** — dibuang. Ini penerapan **preseden Q1**:
tertulis tetapi tidak digunakan → jangan ikut migrasi.

**Jalur simpan yang berlaku:** `Activity/SaveOutstanding_Act.xml` → `SaveOSClaim_SQL` →
`POOLDATA.PEGA_JSON_OS_AKSEP_KLAIM`.

⚠️ `[terverifikasi]` **Nama parameter berbohong, nilainya benar — bukan bug.** `InsertClaimPNC`
mengirim `TempPNC.BUSINESS_CODE` ke parameter bernama `PegaID`, dan `TempPNC.TSI` ke parameter
bernama `DataPega` (CLOB). Yang sebenarnya diisi terbaca di `Activity/InsertJsonClaimTreaty_act.xml`:
`TempPNC.BUSINESS_CODE` diisi `TempOpenPage.pzInsKey` dan `TempPNC.TSI` diisi
`@GCNM.GetPageJSONString()`. Jadi nilainya memang **kunci kasus** dan **dokumen JSON** — hanya
namanya yang menyesatkan. **Dicatat supaya pembaca berikutnya tidak melaporkannya sebagai bug.**

### Empat temuan dari isi stored procedure

`[data DBA]` Keempatnya fakta dari DDL, bukan tafsiran.

1. ⚠️ **Baris `OS_AKSEPTASI_KLAIM` tidak pernah ditimpa — selalu bertambah.** Kedua procedure
   (`PEGA_JSON_OS_AKSEP_KLAIM` dan `…KLAIMTNP`) menjalankan `SELECT count(1) INTO id_count` lalu
   **tidak pernah memakai** `id_count`; di `…KLAIMTNP` cabang `UPDATE`-nya **dikomentari habis**.
   Jadi **satu `CASEID` dapat punya banyak baris**, dan mana yang berlaku hanya terbaca dari
   `TANGGAL` ditambah `DATA_JSON.AcceptedNo` — persis urutan yang dipakai
   `RDBList/DataOutstandingTreatyin.xml` (`order by a.TANGGAL asc, a.DATA_JSON.AcceptedNo desc`).
   ⚠️ **Akibat ke migrasi: migrasi wajib memutuskan baris mana yang berlaku per klaim.** AC 9 belum
   menuntut itu — lihat catatan di akhir bab.
2. ⚠️ **Batas transaksi ada di dalam database, bukan di aplikasi.** Tiap procedure memasang
   `EXCEPTION WHEN OTHERS THEN ROLLBACK` **di dalam dirinya**, dan pemanggilnya menutup dengan
   `COMMIT;` di blok anonim (`BEGIN … ; COMMIT; END;` pada `InsertClaimPNC` dan
   `SaveOSKlaimTreaty_SQL`). **Aplikasi tidak pernah memegang kendali transaksi.** Ini bukti tambahan
   untuk **AC 7** dan **AC 13** — bunyi kedua AC itu tidak berubah.
3. ⚠️ **`TANGGAL` dipotong ke tanggal; jamnya dibuang** —
   `TO_DATE(to_char(sysdate,'dd/MM/yyyy'), 'dd/MM/yyyy')`. Akibatnya **urutan dalam satu hari tidak
   kronologis** dan hanya dapat dipilah oleh `AcceptedNo`.
4. ⚠️ **`PEGA_JSON_KLAIM_PNC` berpola `UPDATE`-lalu-`INSERT` atas `IDPEGA`.** Jadi `JSON_KLAIM`
   **satu baris per kasus Pega** — **berlawanan** dengan `OS_AKSEPTASI_KLAIM` yang menumpuk. Dua
   tabel proyeksi, dua perilaku baris yang berbeda. `INSERT`-nya juga memasang `IDPROD = 1` dan
   `TGL_KONVERSI = NULL` · `STS_KONVERSI = NULL` secara **hardcode**.

> ⚠️ **Spec perlu satu AC tambahan — TIDAK ditulis di sini.** AC 9 hanya berbunyi *"Migrasi
> memindahkan data klaim lama dari bentuk JSON ke bentuk relasional"*; ia **tidak** menuntut
> pemilihan baris yang berlaku ketika satu `CASEID` punya banyak baris. Usulan bunyinya dilaporkan
> ke work owner, bukan disisipkan sendiri ke spec.

## ADR terkait

**ADR-0003** (uang non-float) · **ADR-0006** (penomoran lewat stored procedure) ·
**ADR-0011** (unit keputusan = baris `AdjustmentList`).

## Acceptance criteria

### Skema relasional

- [ ] ⚠️ Skema relasional Claim Prop **dirancang baru**; test yang mengandalkan tabel klaim warisan **gagal** *(AC 1 spec)*
- [ ] Pohon klaim tersimpan relasional dan tiap baris dapat di-`SELECT` tanpa membongkar JSON *(AC 2 spec)*
- [ ] ⚠️ Enam properti spreading menjadi **tabel berbeda** meski di Pega berbagi satu class. **Jebakan:** class yang sama bukan bukti peran yang sama — jangan menyatukannya jadi satu tabel. Ini **bukan penyimpangan**: skema relasionalnya memang dirancang baru (AC 1), dan berbagi definisi class di Pega tidak punya akibat perilaku apa pun *(AC 3 spec)*
- [ ] Jejak audit menjadi tabel milik Claim Prop sendiri, bukan menumpang struktur modul lain *(AC 4 spec)*
- [ ] Menghapus klaim mengkaskade sampai tingkat terdalam; test memeriksa cicit *(AC 5 spec)*
      ⚠️ **RALAT 2026-09-19** — butir di atas **dikutip apa adanya, tidak dihapus**. AC 5 kini punya
      **satu pengecualian**: klaim yang **pernah** punya kasus komite **tidak dapat dihapus sama
      sekali**, dan **kaskadenya tidak pernah dijalankan**. Diuji oleh butir **AC 128** di bawah.
- [ ] Satu aksi pengguna = satu transaksi; keadaan separuh tersimpan membuat test **gagal** *(AC 6 spec)*
- [ ] ⚠️ Nol `COMMIT` di dalam SQL aplikasi — 15 dari 58 rule yang `COMMIT` sendiri **tidak direplikasi**. **Alasan menyimpang:** batas transaksi menjadi urusan aplikasi *(AC 7 spec)*
- [ ] ⚠️ Setiap objek Oracle diprefiks schema eksplisit; empat schema dikenal `pooldata` · `arasapas` · `reinsurance` · `gl`. **Alasan menyimpang:** mayoritas tabel di Pega tanpa prefiks sama sekali, sehingga resolusi bergantung schema bawaan koneksi *(AC 8 spec)*
- [ ] Migrasi memindahkan data klaim lama dari bentuk JSON ke bentuk relasional *(AC 9 spec)*
- [ ] ⚠️ Parser JSON warisan membuang seluruh kunci berawalan `px`/`py`/`pz` dan tidak mengandalkan kunci tertentu selalu hadir *(AC 10 spec)*
- [ ] ⚠️ Status tampilan (`pyExpanded`) **tidak** ikut disimpan. **Alasan menyimpang:** kebocoran lapisan UI ke penyimpanan tidak layak ditiru *(AC 11 spec)*
- [ ] `[data DBA]` Pembaca JSON warisan menerima **dua format tanggal** berdampingan dalam satu dokumen *(AC 12 spec)*
- [ ] ⚠️ `[data DBA]` **Migrasi memilih satu baris yang berlaku per klaim** dari `OS_AKSEPTASI_KLAIM` yang menumpuk, memakai urutan `TANGGAL` **naik** lalu `DATA_JSON.AcceptedNo` **turun** — urutan yang sama dipakai `RDBList/DataOutstandingTreatyin.xml`; baris sisanya **disimpan sebagai riwayat**, tidak dibuang. **Alasan menyimpang:** kedua stored procedure penulis menjalankan `SELECT count(1) INTO id_count` lalu **tidak pernah memakainya**, dan di `PEGA_JSON_OS_AKSEP_KLAIMTNP` cabang `UPDATE`-nya dikomentari habis — sehingga satu `CASEID` punya banyak baris **tanpa penanda mana yang berlaku** *(AC 123 spec)*

### Batas transaksi penomoran

- [ ] ⚠️ Pengambilan nomor dan penyimpanan hasilnya berada dalam **satu transaksi**; penyimpanan gagal → nomor tidak terpakai. **Alasan menyimpang:** di Pega lima generator `COMMIT` sendiri, sehingga nomor terbakar begitu di-generate *(AC 13 spec)*
- [ ] Test: gagalkan langkah sesudah pengambilan nomor → nomor berikutnya **berurutan rapat** *(AC 14 spec)*
- [ ] Logika pembentukan nomor tetap di stored procedure, tidak direplikasi di aplikasi *(AC 15 spec)*
- [ ] ⚠️ Nomor **temp** tidak dapat dipromosikan menjadi nomor final — ia kekurangan kode bisnis dan bulan *(AC 16 spec)*
- [ ] Nomor klaim, PLA, dan DLA memakai pola tanda tangan yang sama — lima parameter masuk, dua keluar *(AC 17 spec)*
- [ ] ⚠️ `[data DBA]` **Panggilan stored procedure penulis adalah langkah TERAKHIR sebelum `COMMIT`** dalam satu transaksi, tanpa pekerjaan penting yang masih menggantung sebelumnya; aplikasi memeriksa `StsSimpan` (**`1`** berhasil, **`0`** gagal) dan pada `0` melaporkan **gagal**. **Alasan menyimpang:** ketiga procedure memasang `EXCEPTION WHEN OTHERS THEN ROLLBACK` di dalam dirinya sendiri, dan `ROLLBACK` telanjang di Oracle membatalkan **seluruh transaksi** — termasuk pekerjaan aplikasi sebelumnya — sehingga pada `StsSimpan = 0` transaksinya sudah dibatalkan dan `COMMIT` sesudahnya tidak menyimpan apa pun. AC 6 menuntut **hasil**, bukan **urutan** *(AC 124 spec)*

### Uang

- [ ] ⚠️ Nilai uang bertipe desimal presisi arbitrer; nol *binary floating point* di lapisan mana pun maupun di kontrak API *(AC 18 spec)*
- [ ] `[keputusan work owner]` Nilai disimpan penuh tanpa pembulatan; pembulatan hanya untuk tampilan *(AC 19 spec)*
- [ ] ⚠️ Perhitungan spreading dilakukan **satu fungsi**. **Alasan menyimpang:** di Pega nilai spreading ditulis di 20 titik dengan 4 perlakuan pembulatan berbeda, sehingga hasilnya bergantung urutan pemanggilan *(AC 20 spec)*
- [ ] ⚠️ Urutan operasi **kali dulu, bagi terakhir**; hasil tidak berubah oleh urutan pemanggilan *(AC 21 spec)*
- [ ] ⭐ **BARU 2026-09-19 (ronde 4)** — `[keputusan work owner]` **Share ceding dikalikan TEPAT SEKALI** pada rantai estimasi klaim. **Alasan menyimpang — perubahan sadar KETIGA:** di Pega `CountPersen_act` langkah **5.1** mengalikan share ceding pada nilai yang **sudah** dikalikan share ceding di hulu (`AddListClaimAmount` 7.1 atau `CountListClaimAmountIDR` 1.1 → dijumlah di `CountPersen_act` 4.3.1 → `.ClaimEstimation` di 4.4). Faktor kesalahannya **ShareCeding/100**, arahnya **terlalu kecil**; pada ShareCeding = 100 selisihnya **nol**, sehingga cacat ini tidak terlihat pada nilai bawaan. Catatan pengembang di rule itu sendiri melarang perkalian kedua tersebut. **Rumusnya DIPERBAIKI, bukan ditiru** *(⛔ belum ada nomor AC spec — `spec.md` belum ditambal)*
  - ⚠️ `[terbuka]` **Nasib nilai lama belum diputuskan** — angka yang sudah terlanjur tersimpan dengan perkalian ganda. **Menunggu work owner.**
  - ⚠️ `[terbuka]` **Cakupannya belum ditetapkan.** Perbaikan ini diterapkan pada `CountPersen_act` **5.1 saja**; lima titik perkalian lain (`CountGrossAdjTreaty_Act` 3 · 4 · 5.2 dan `CountValueADJTreaty_Act` 12) mengalikan **sekali saja** dan **tidak ada catatan yang melarangnya** — keduanya **disalin apa adanya** sampai work owner memutuskan sebaliknya.
  - ⛔ `[terbuka]` **Bergantung pada arti label `//`.** `CountPersen_act` langkah **5** berlabel `//` tetapi gerbangnya `2/2` — **tanpa syarat**. Ekspor Pega **tidak memuat medan apa pun yang mematikan langkah**. Kalau `//` bukan saklar, perkalian ganda itu benar-benar berjalan di produksi. Lihat `grilling-ronde-4.md` §B3-RALAT.
- [ ] ⚠️ Nilai turunan dihitung **sekali jalan dari basis asli**. **Alasan menyimpang:** `[data DBA]` satu baris contoh membuktikan nilai turunan dihitung dari angka yang sudah dipotong 4 desimal di hulu *(AC 22 spec)*
- [ ] ⚠️ Angka yang **sudah terbit** — nomor akseptasi terbit, DLA dicetak, data terkirim ke Kasir — dibekukan dan tidak ikut dihitung ulang *(AC 23 spec)*
- [ ] Nilai spreading diperlakukan sebagai turunan, boleh dihitung ulang dari sumber *(AC 24 spec)*
- [ ] ⚠️ Nilai uang tidak pernah disimpan sebagai teks; **nol tambalan pemisah desimal** di jalur baca mana pun. **Alasan menyimpang:** dua jalur Pega menambal koma-ke-titik — `RDBList/CariHistoryClaim_SQL.xml` dan `Activity/CekPremiLunas_Act.xml` step 5 *(AC 107 spec)*

> ⚠️ **RALAT 07-10-2026** — kalimat lama: *"Share ceding dikalikan TEPAT SEKALI … perubahan sadar KETIGA"*. **Dibatalkan**:
> spec §4 (19-09-2026) memutuskan *"BUKAN perubahan sadar"* — perkalian ganda share ceding tidak pernah
> berjalan, jadi tidak ada yang diperbaiki. Implementasi meniru rantai yang berjalan (`models.HitungTurunan`).

### Mata uang

- [ ] ⚠️ Bentuk uang adalah **(nilai, mata uang, kurs) per baris**; test yang menuntut satu klaim satu mata uang **gagal**. **Peringatan lintas-modul:** ini **bukan penyimpangan** — Pega memang sudah menetapkan mata uang per baris (`Activity/CurencyEstimation_Act.xml`). Yang diperingatkan: invariant satu-mata-uang milik Claim Life **tidak boleh ikut disalin** ke sini *(AC 25 spec)*
- [ ] Subtotal per mata uang tersimpan sebagai entitas tersendiri *(AC 26 spec)*
- [ ] Tiap objek pertanggungan punya mata uang dan kursnya sendiri *(AC 27 spec)*
- [ ] Satu klaim dengan dua mata uang tersimpan dan terhitung benar dari ujung ke ujung *(AC 28 spec)*

### Model penanda kesalahan dan penamaan kolom

- [ ] ⚠️ Penanda kesalahan dipecah menjadi tiga hal bernama sendiri — *kesalahan yang memblokir* (ya/tidak) · *peringatan yang ditampilkan* (daftar pesan) · *boleh kirim ke komite* (ya/tidak). **Alasan menyimpang:** di Pega ketiganya berbagi satu properti angka bernilai `0`·`1`·`2`·`3`·kosong yang artinya tidak tertulis di mana pun; menaikkan ambang menjadi `>= 1` akan membuat saklar kirim-komite ikut memblokir hal yang bukan urusannya *(AC 105 spec)*
- [ ] ⚠️ Menghapus klaim yang **pernah** punya kasus komite **ditolak — selamanya**, walau komitenya sudah selesai maupun menolak; **kaskadenya tidak pernah dijalankan** dan **tidak satu baris pun disentuh**. Pesannya **menyebut berapa kasus komite yang pernah ada**. Klaim yang **belum pernah** ke komite **tetap terhapus** dan **mengkaskade sampai cicit**. **Alasan menyimpang:** di Pega tidak ada penjaga apa pun — klaim boleh dihapus kapan saja *(AC 128 spec)*
- [ ] ⚠️ **Tidak ada kolom `FLAG_ON_GOING_COMMITTEE`** di skema baru, dan **nilai lamanya tidak dimigrasikan**; layar akseptasi menampilkan *"komite sedang berjalan"* **dari kasus komite yang berjalan**. Test yang **mengandalkan kolom itu gagal**. **Alasan menyimpang:** kolomnya ada di Pega dan dibaca layar akseptasi, tetapi nilainya **dapat diturunkan**, sehingga menyimpannya berarti memelihara dua sumber kebenaran untuk satu fakta *(AC 129 spec)*
- [ ] ⚠️ Klaim yang kasus komitenya **hanya** lewat jalur **penutupan tanpa pembayaran** — tanpa satu pun baris penyesuaian pernah diserahkan — **tetap ditolak saat dihapus**. Test yang berhasil menghapusnya **gagal**. ⚠️ Baris penyesuaian klaim itu **tidak ikut beku**: kasus komite jalur penutupan tidak menunjuk baris mana pun *(AC 128 spec)*
- [ ] ⚠️ **Klaim lama hasil migrasi tidak dapat dihapus sama sekali**, apa pun riwayat komitenya; test yang berhasil menghapus klaim hasil migrasi **gagal**. Migrasi **tidak membaca** `FLAG_ON_GOING_COMMITTEE` — test yang menemukan migrasi membacanya, atau memakainya sebagai bukti riwayat komite, **gagal** *(AC 128 · AC 129 spec)*
- [ ] ⚠️ Nama kolom dibuat sesuai isinya, bukan disalin dari alias warisan. **Alasan menyimpang:** `RDBList/CariHistoryClaim_SQL.xml` memberi enam alias yang lima di antaranya berbohong — tanggal kejadian dinamai *start date*, nomor klaim dinamai *nama cabang*, id Pega dinamai *kode cabang*, penyebab kerugian dinamai *kode bisnis*, estimasi klaim dinamai *TSI* *(AC 106 spec)*

## Perintah verifikasi

```
# nol COMMIT di SQL aplikasi
cari "COMMIT" di lapisan persistensi                      -> nihil

# nol floating point untuk uang
cari tipe float di model uang dan kontrak API             -> nihil

# nomor tidak terbakar
jalankan test "gagal sesudah ambil nomor -> nomor berikutnya rapat"

# schema selalu eksplisit
periksa tiap FROM/JOIN/INTO membawa prefiks schema        -> tidak ada yang telanjang

# multi mata uang
jalankan test "klaim dua mata uang tersimpan dan terhitung"

# kaskade
jalankan test "hapus klaim -> cicit ikut terhapus"
```
