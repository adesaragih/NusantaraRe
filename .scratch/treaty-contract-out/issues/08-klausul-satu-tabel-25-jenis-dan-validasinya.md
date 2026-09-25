# 08: Klausul — satu tabel, 25 jenis, validasi per jenis

**Status:** ready-for-agent

**Blocked by:** 04 (kombinasi tahun/grup/jenis dibuka oleh kontrak)

⚠️ **Inti konteks ini.** Tiket paling berisi: 13 acceptance criteria, dan satu-satunya tempat di
mana penyimpangan sadar 2 diterapkan.

## Hasil & nilai pengguna

Sebagai **underwriter**, saya ingin mengelola **dua puluh lima jenis klausul** kontrak treaty di
satu tempat — sebagian berlapis dengan baris rincian di bawahnya — dan saya ingin **setiap jenis
memeriksa field yang memang relevan baginya**, bukan satu daftar wajib yang seragam.
*(User story 17–22 di spec)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Satu bentuk baris klausul; jenis sebagai atribut |
| `internal/repository` | Tulis/baca **satu tabel** `proportionalarrg` |
| `internal/services` | Pemetaan jenis klausul → aturan wajib-isi; anti-dobel; hitungan turunan |
| `internal/handlers` | Endpoint daftar + simpan per jenis klausul |
| `frontend/` | Grid per jenis klausul; grid rincian untuk jenis berlapis |

## Rule Pega sumber

`[terverifikasi]` Sensus penuh atas **25** activity `SaveTreatyArr*` di
`Treaty Contract Out/Activity/` — seluruhnya berujung pada **dua** rule Connect-SQL:

| Menulis lewat | Jml | Jenis klausul |
| --- | ---: | --- |
| `SaveMasterProportionalArrg` → `POOLDATA.PEGA_PROPORTIONALARRG` (**35 kolom**) | **18** | `BordereAux`, `CashLossLimit`, `ClaimCoorp`, `CoinsPanel`, `EPI`, `ExGratia`, `ExclutionTreaty`, `FacIn`, `LimitMB`, `MaxCoinsPanel`, `MinLOL`, `MinLOLMB`, `PLA`, `Portfolio`, `ProfitComm`, `Ricomm`, `TerrLimit`, `TreatyLimit` |
| `SaveMasterProportionalArrgChild` → `POOLDATA.PEGA_M_PROPORTIONALARRG_CHILD` (**26 kolom**) | **7** | `CashLossLimitList`, `ClaimCoorpChild`, `EpiList`, `ExGratiaChildList`, `FacInList`, `PLAList`, `TreatyLimitChild` |

| Rule | Class / Nama / Tipe | Path |
| --- | --- | --- |
| `SaveMasterProportionalArrg` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterProportionalArrg.xml` |
| `SaveMasterProportionalArrgChild` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/SaveMasterProportionalArrgChild.xml` |
| `SaveTreatyArrEPI_Act` | `ASM-FW-GISFW-INT-PROPORTIONALARRG` / `SAVETREATYARREPI_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveTreatyArrEPI_Act.xml` (111.418 byte, 10 langkah) |
| `HitungRpUsd` | `@BASECLASS` / `HITUNGRPUSD` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/HitungRpUsd.xml` |
| `TreatyTestChildTotal_Act` | `@BASECLASS` / `TREATYTESTCHILDTOTAL_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/TreatyTestChildTotal_Act.xml` |
| `BrowseTreatyDesc_RD` | `ASM-FW-GISFW-INT-TREATYDESC` / `BROWSETREATYDESC_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/BrowseTreatyDesc_RD.xml` |
| 16 `CancelActivity*` | `Treaty Contract Out/Activity/` | pembatalan per jenis klausul |

`[terverifikasi]` **Mekanisme existing**: tiap jenis punya halaman masukannya sendiri
(`InputTreatyArrEpi`, `InputTreatyArrPLA`, `InputTreatyArrProfit`, `InputTreatyArrBord`, …), lalu
langkah **`Page-Copy`** menyalinnya ke halaman bersama `InputTreatyArrTreatyLimit` (induk) atau
`InputTreatyArrTreatyLimitChild` (anak) sebelum langkah `RDB-List` memanggil procedure.

⚠️ **Penyimpangan sadar 2 — SATU tabel untuk semua jenis klausul.** `[keputusan work owner]`
`[data DBA]` **Kedua procedure menulis ke tabel yang sama: `POOLDATA.PROPORTIONALARRG`.** Bedanya
hanya jumlah kolom. Baris "anak" mengisi **NULL** pada sembilan kolom khusus induk:
`ID_OCCUPATION`, `OCCUPATION`, `ID_CLAUSE`, `CLAUSE`, `TREATYLIMIT`, `COINS_MIN`, `COINS_MAX`,
`MORERP`, `MOREUSD`. Jenis dibedakan **`TREATYDESCID`**.

`[data DBA]` Identitas `'1' || lpad(PROPORTIONALARRG_SEQ.nextval, 7, '0')` — **7 digit**, berbeda
dari lima tabel lain yang 6 digit. **Tidak `COMMIT` sendiri**; `StsSimpan` **1 = sukses / 0 = gagal**.

## Validasi wajib-isi per jenis `[terverifikasi]`

Sensus prasyarat langkah (`<pyStepsPreCondParamsWhen> … ==""`) atas 25 activity — **ini tabel
kebenaran untuk AC validasi**:

| Jenis klausul | Field wajib |
| --- | --- |
| `CashLossLimit`, `ClaimCoorp`, `EPI`, `ExGratia`, `FacIn`, `PLA`, `TreatyLimit` | `ReinsTypeID`, `Rp`, `Usd` |
| `CashLossLimitList`, `ClaimCoorpChild`, `EpiList`, `ExGratiaChildList`, `FacInList`, `PLAList`, `TreatyLimitChild` | `Pct`, `ReinsTypeID`, `Rp`, `Usd` |
| `ProfitComm` | `Pct`, `PctMe`, `ReinsTypeID` |
| `Ricomm` | `Method`, `Pct`, `ReinsTypeID` |
| `CoinsPanel` | `TerritorialLimit`, `TreatyLimit` |
| `MaxCoinsPanel` | `CoIns_Max`, `TerritorialLimit` |
| `MinLOL`, `MinLOLMB` | `Pct`, `TerritorialLimit` |
| `TerrLimit` | `TerritorialLimit` |
| `ExclutionTreaty` | `ID_Occupation`, `TerritorialLimit` |
| `BordereAux` | `Method` |
| **`LimitMB`** | `[terbuka]` — hari ini **nol validasi** |
| **`Portfolio`** | `[terbuka]` — hari ini **nol validasi** |

## ADR terkait

**ADR-0003** (uang & persen non-float), **ADR-0006** (identitas lewat sequence),
**ADR-0007** (jejak audit), **ADR-0015** (kegagalan eksplisit).

## Acceptance criteria

- [ ] ⚠️ Seluruh **dua puluh lima** jenis klausul tersimpan dalam **satu tabel**
      `proportionalarrg`, dibedakan oleh **`TREATYDESCID`**. Test yang menemukan tabel terpisah per
      jenis klausul **gagal**. *(AC 24 spec; penyimpangan sadar 2)*
- [ ] ⚠️ Baris jenis "anak" menyimpan **sembilan kolom khusus induk sebagai NULL**
      (`ID_OCCUPATION`, `OCCUPATION`, `ID_CLAUSE`, `CLAUSE`, `TREATYLIMIT`, `COINS_MIN`,
      `COINS_MAX`, `MORERP`, `MOREUSD`). *(AC 25 spec; penyimpangan sadar 2)*
- [ ] Daftar jenis klausul dibaca dari master **`TREATYDESC`**; jenis baru dapat ditambahkan di
      master **tanpa mengubah skema**. *(AC 26 spec; User story 18)*
- [ ] Master jenis klausul **tidak ditulis** oleh konteks ini. Test yang menemukan tulisan ke
      `TREATYDESC` **gagal**. *(AC 27 spec)*
- [ ] Klausul berlapis dapat dinyatakan lewat baris "anak" pada kombinasi yang sama, dibedakan
      **`PARENTREINSTYPEID`**. *(AC 28 spec; User story 19)*
- [ ] Membatalkan pengeditan **satu** jenis klausul **tidak** membuang jenis klausul lain yang
      sedang dikerjakan. *(AC 29 spec; User story 21)*
- [ ] Baris klausul yang **sudah pernah diinput** **ditolak** dengan pesan yang jelas.
      *(AC 30 spec; User story 22)*
- [ ] Urutan baris dalam daftar klausul **dipertahankan** saat dibaca kembali. *(AC 31 spec)*
- [ ] ⚠️ Nama jenis klausul dan nama field **dipakai apa adanya** dari Pega — **tidak
      diterjemahkan, tidak ditebak**. *(AC 32 spec; `[keputusan work owner]`)*
- [ ] Setiap jenis klausul memeriksa **field yang relevan baginya sendiri**, sesuai tabel di atas.
      *(AC 33 spec; User story 20)*
- [ ] Aturan wajib-isi berada di **kode**, bukan di data master. *(AC 34 spec)*
- [ ] Pesan penolakan **menyebut field** yang kurang. *(AC 35 spec)*
- [ ] `[terbuka]` Jenis **`LimitMB`** dan **`Portfolio`** **tidak** dinyatakan selesai sebelum
      aturan wajib-isinya ditetapkan Product + UW. Keduanya **tidak** dilepas sebagai "tanpa
      validasi". *(AC 36 spec)*
- [ ] ⚠️ Nilai uang (`Rp`, `Usd`, `MoreRp`, `MoreUsd`, `TreatyLimit`, `CoIns_Min`, `CoIns_Max`) dan
      persen (`Pct`, `PctMe`) bertipe **desimal**, bukan teks. *(AC 51, 52 spec; **ADR-0003**;
      penyimpangan sadar 6)*
- [ ] Identitas baris klausul **tidak pernah diketik pengguna**; berbentuk `'1'` + **7 digit** —
      berbeda dari lima tabel lain yang 6 digit. *(AC 6 spec; `[data DBA]`; **ADR-0006**)*

## Blocker

**Tidak ada pemblokir.**

⚠️ `[terbuka]` **Dua OQ yang tidak memblokir tiket ini, tetapi memblokir dua jenis klausul:**
field wajib `LimitMB` dan `Portfolio` (AC di atas), dan **arti bisnis** setiap istilah klausul
(`EPI`, `PLA`, `Ricomm`, `MB`, `LOL`, `BordereAux`, `ExGratia`, `CashLossLimit`,
`ClaimCoorperation`, `ProfitCommision`, `pTerrLimit`, `CoinsPanel`, `ExclutionTreaty`) dan field
(`Ydcf`, `PctMe`, `LayerPartType`, `LayerType`, `SpreadingOrder`, `Method`, `Line`). Keduanya
**Product + UW**. Layar dan skema dapat dibangun penuh dengan istilah asli.

## Catatan

⚠️ **Dua perilaku existing yang mudah terlewat — keduanya nyata dan harus ikut.**

1. **`Rp` dan `Usd` pada baris anak adalah nilai TURUNAN.** `[terverifikasi]`
   `Activity/HitungRpUsd.xml` (`@BASECLASS!HITUNGRPUSD`) menghitung, untuk ketujuh jenis anak,
   `Usd = Pct × Usd_induk / 100` dan `Rp = Pct × Rp_induk / 100`. Jadi pengguna mengetik **`Pct`**,
   dan kedua nilai uang mengikuti. Perilaku ini **bukan** kosmetik layar — angka turunannya yang
   tersimpan.
2. **Total `Pct` baris anak wajib 100%.** `[terverifikasi]`
   `Activity/TreatyTestChildTotal_Act.xml` (`@BASECLASS!TREATYTESTCHILDTOTAL_ACT`) menjumlahkan
   `Pct` seluruh baris anak dan menolak bila `!= 100.00`, dengan pesan
   `"Please make sure spreading is 100%"`. Ini gerbang yang **memang ada** di existing — berbeda
   dari total share reinsurer (tiket 05) yang **tidak** ditegakkan.

⚠️ **Nama berbohong — dua kali dalam satu tiket.** `[data DBA]` Procedure
`PEGA_M_PROPORTIONALARRG_CHILD` **tidak** menulis ke tabel child; ia menulis ke tabel yang sama.
Dan `[terverifikasi]` `Activity/SetTreatyArrangementDesc_Act.xml`
(`ASM-FW-GISFW-INT-PROPORTIONALARRG`) — yang namanya paling wajar untuk "mengisi deskripsi
arrangement" — **lima dari enam langkahnya di-remark**; ia mati (**OQ-066**).

⚠️ **`ParentReinsTypeID = "00"`** adalah penanda baris tanpa induk `[terverifikasi]`
(`Activity/GetPeriode.xml`, `@BASECLASS!GETPERIODE`). Sentinel ini dibawa apa adanya.

⚠️ **Jangan memigrasikan kueri yang membaca JSON.** `GetMasterDescriptionEPIParentList`,
`GetMasterDescriptionPLAParentList`, `GetMasterDescriptionExGratiaList`,
`GetMasterDescriptionFACINParentList`, `GetMasterPortfolioListDetail`, `GetMasterPanggilID` —
seluruhnya membaca `m_PROPORTIONALARRG`. Di sistem baru jawabannya datang dari
**`proportionalarrg` relasional** (tiket 01).

⚠️ **`Portfolio` di existing berupa daftar bersarang di dalam JSON** `[terverifikasi]`
(`RDBList/GetMasterPortfolioListDetail.xml` memakai `JSON_TABLE` atas `'$.ProportionalList[*]'`
dengan kolom `PortfolioType`, `PremiLost`, `PortfolioValue`, `PortfolioDesc`). Karena JSON dibuang,
**bentuk penyimpanan Portfolio di skema baru ditetapkan bersama field wajibnya** — bagian dari OQ
`LimitMB`/`Portfolio` di atas. **Jangan tebak.**

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — satu tabel yang memuat 25 jenis dengan
pola NULL yang berbeda per jenis hanya terbukti benar pada basis data sungguhan.

```
go test ./internal/...
cd frontend && npm test
make check
```
