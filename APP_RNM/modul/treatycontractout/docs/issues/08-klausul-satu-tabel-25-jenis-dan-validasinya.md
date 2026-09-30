# 08: Klausul — satu tabel, 25 jenis, validasi per jenis

**Status:** selesai (29-09-2026)

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

## Pembacaan ulang XML — 29-09-2026 (sesi modul, lanjutan 1)

Nomor baris = baris mentah berkas korpus `Treaty Contract Out/`; langkah aktivitas dibaca lengkap, **hanya langkah
hidup** (bukan `pyStepsBlockName = //`) dengan prasyarat **aktif** (`pyStepsPreCondition = true`).

| Unsur | Bukti | Dibawa sebagai |
| --- | --- | --- |
| jalan masuk | `Section/InputTreatyContract.xml` b22196 `List Description` → harness `Harness/InboxTreatyContractDescription.xml` b359 | tombol baris tahun membuka `PanelKlausulTahun`; menu `InboxTreatyContractDescription` memilih tahun lebih dulu |
| kepala | harness b2003 `TreatyGroupID`, b2489 `Underwriting Year` (`.TreatyYear`), b2665 `Transaction Year` (`.UnderwritingYear`), b2839/b3025 tanggal, b3209 `Treaty Description`, b3382 `Proportion Type` | tujuh medan baca-saja; label bersilang OQ-TCO-05 dibawa apa adanya |
| grid jenis | harness b4880 `For Non XOL` (saringan `IsXOL` b5324/b5335 = `0`), b7971 `For XOL` (b8415/b8426 = `1`); kolom b5440 `ID`, b5549 `Description Name`; tombol b6059 `Show` | `GET /api/treaty-contract-out/jenis-klausul?isXol=` dari master `TREATYDESC` (baca-saja, AC 26/27) |
| DescID → jenis | `Activity/SetKirimIDDesc.xml` — 54 rujukan `InputData.CARIDESC` (mis. b1017/b1142 `10001` → `InputTreatyArrTreatyLimit`, b1196/b1288 `10001` → `InputTreatyArrTreatyLimitChild`, b1342 `10010` → `InputTreatyRicomm`) | induk dan anak **berbagi** DescID, dibedakan `PARENTREINSTYPEID` (`"00"` = induk) |
| aturan wajib-isi | sensus 25 aktivitas `Activity/SaveTreatyArr*` (langkah `Property-Set-Messages` berprasyarat `<medan>==""`, WhenTrue 2) | `models.AturanKlausulTCO` — satu tabel di kode (AC 34), pesan menyebut medan (AC 35) |
| Rp/Usd anak | `Activity/HitungRpUsd.xml` b252/b341 (TreatyLimitChild), b370/b397 (PLA), b418/b439 (CashLoss), b460/b481 (EPI): `toDecimal(Pct) * toDecimal(<induk>.Usd\|Rp) / 100` | `RpUsdAnakTCO` — `apd`, kuantisasi 8 desimal; klien tidak mengirimnya |
| total Pct anak | `Activity/SaveTreatyArrEpiList_Act.xml` b1055 `Menjumlahkan Total PCT`, b1300 `Valdiasi Total agar tidak boleh lebih dari 100` (tujuh aktivitas anak memuatnya) | total anak BARU > 100 ditolak 422 (baris dikunci `FOR UPDATE`) |
| spreading | `Activity/TreatyTestChildTotal_Act.xml` b892 `ERRMSG + " \n Please make sure spreading is 100%"`; dipanggil `SaveTreatyArrTreatyLimitChild_Act.xml` b3820 dan tombol `GridTreatyArrangementTreatyLimit.xml` b11687/b11883 | **peringatan** sesudah simpan TreatyLimitChild — baris tetap tersimpan |
| anti-dobel | `SaveTreatyArrEPI_Act.xml` b2011 `"Data sudah pernah di Input"` berprasyarat b2089 `@equals(OutputData.HASILD7,1)`; `RDBList/SaveMasterProportionalArrg.xml` b120–b121 hanya mengeluarkan `HASIL1`, `HASIL2` | pesan VERBATIM, kunci dobel per jenis `[keputusan kami]` (OQ-TCO-14), 409 menyebut baris lain |
| exclusion | empat sub-bagian `GridTreatyArrangementExclutionTreaty{Occupation,Clausule,Object,Periode}.xml` mengirim `Param.Type` (b5082/b5168, b4305/b4392, b2086/b2202, b1798/b1916); nilai `Occupation` b7718, `Clause` b6943 | satu jenis, empat subjenis; nama occupation/clause dari master |
| form per jenis | label medan b2777–b4487 `GridTreatyArrangementEpi.xml`; `Save` b5501, `Add` b8980, `Edit` b10917, `Show Child` b11200; `GridTreatyArrTreatyEpiList.xml` b3029 `Pct`, b8657 `Close Child`; label khusus MinLOL b1540, MinLOLMB b1548, MaxCoinsPanel b1525, Coins b3394/b4193/b4880, Exclusion b2007/b2296/b2745/b2935/b3222, b2030/b2320, b566, b500 | `PanelJenisKlausul` dirakit dari aturan server; ±50 baris label diuji ke korpus |

### Ralat bertanggal 29-09-2026

1. **`TerritorialLimit` bukan gerbang** CoinsPanel, MaxCoinsPanel, MinLOL, MinLOLMB, ExclutionTreaty: prasyaratnya
   NONAKTIF (`SaveTreatyArrCoinsPanel_Act.xml` b862/b1045/b1232 `pyStepsPreCondition=false`) dan merujuk halaman lain
   (b984/b1172 `InputTreatyExclutionTreaty`, b1338 `InputTreatyTerr`) — residu salin-tempel. Tabel tiket diralat.
2. **CoinsPanel** mewajibkan `CoIns_Min`, `CoIns_Max`, `TreatyLimit` (form `From`/`To`/`Treaty Limit`); **MaxCoinsPanel**
   `CoIns_Max`; **MinLOL/MinLOLMB** `Pct`.
3. **ExclutionTreaty** = empat subjenis (`Param.Type`): Occupation (`ID_Occupation`, `Occupation`, `Line`, `Usd`, `Rp`
   wajib), Clause (`ID_Clause`, `Clause` wajib). Object dan Periode mengirim `Type` kosong — di Pega nol validasi
   berjalan; satu-satunya medan formnya diwajibkan `[keputusan kami]` (**OQ-TCO-14**).
4. **Total Pct anak "wajib 100%"** hanya **peringatan** dan hanya pada TreatyLimitChild; yang MENOLAK adalah total > 100,
   dan itu berlaku pada ketujuh jenis anak.
5. **`Data sudah pernah di Input` mati di Pega**: prasyarat `HASILD7` tidak pernah diisi prosedur. AC 30 ditegakkan
   dengan kunci dobel per jenis (`KunciDobel` di tabel aturan) — **OQ-TCO-14**.
6. **Pemilih `ReinsTypeID`** di form memakai `D_EnumerationList` yang tidak diekspor; sistem baru memakai daftar jenis
   reasuransi tiket 02 (**OQ-TCO-15**).
7. **Nama kolom master** `TREATYDESC` (`DESCNAME`, `ISXOL`, `STATUSAKTIF`), `OCCUPATION` (`NAME`), `CLAUSE` (`INFO`,
   `TYPE = 'FIRE'`) diturunkan dari properti RD, bukan dari SQL korpus `[dugaan kuat]` (**OQ-TCO-16**).
8. **LimitMB dan Portfolio ditahan** (AC 36): simpan ditolak 422 dengan alasannya; layar menampilkan alasan itu, bukan form.
9. **Ruang lingkup baris** = `TREATYYEARID` (+ `TREATYDESCID`, `PARENTREINSTYPEID`) — bukan kombinasi teks.
   `TREATYYEAR`/`TREATYGROUPID` tetap ditulis untuk hilir.

### Yang dibangun

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `tco_klausul.go` (+uji) | `AturanKlausulTCO` (25 jenis + 3 subjenis exclusion tambahan), `PeriksaKlausulTCO`, `RpUsdAnakTCO`, `PeriksaTotalAnakTCO`, `PeringatanSpreadingTCO`, `UraiDesimalMasukTCO` |
| repository | `tco_klausul.go` (+uji) | `MasterKlausulTCO` satu tabel `T_PROPORTIONALARRG` (35 kolom, 9 kolom induk NULL untuk anak), master `TREATYDESC`/`OCCUPATION`/`CLAUSE` baca-saja, `MasterTahunTreaty.Kunci` |
| services | `tco_klausul.go` (+uji) | `KlausulTCO`: jenis+aturan, daftar (total Pct + peringatan), simpan satu transaksi {kunci tahun, dobel, total anak, tulis, jejak} |
| handlers | `tco_klausul.go` (+uji, +uji `db`) | 5 rute; tanpa jalur hapus |
| frontend | `PanelKlausulTahun.tsx`, `PanelJenisKlausul.tsx`, `InboxTreatyContractDescription.tsx` (+uji), `KLAUSUL_TCO`/`LABEL_MEDAN_*`, `api.ts` (+4), tombol `List Description` hidup, butir menu ketiga | |

**Status:** selesai 29-09-2026 — commit `treaty-contract-out: tiket 08 — klausul satu tabel, 25 jenis, validasinya`.

## Keputusan work owner 29-09-2026

- **OQ-TCO-14 — ditutup.** Jawaban: *"setuju"*. Kunci anti-dobel per jenis (`KunciDobel`) dan wajib-isi satu-satunya
  medan subjenis Object/Periode dikonfirmasi `[keputusan work owner 29-09-2026]`.
- **OQ-TCO-15 — ditutup.** Jawaban: *"benar"*. Pilihan `ReinsTypeID` form klausul = daftar jenis reasuransi tersaring
  tiket 02.
- **OQ-TCO-16 — ditutup.** Jawaban: *"benar"*. Nama kolom master `TREATYDESC` (`DESCNAME`, `ISXOL`, `STATUSAKTIF`),
  `OCCUPATION` (`NAME`), `CLAUSE` (`INFO`, `TYPE = 'FIRE'`) dikonfirmasi.

### Ralat bertanggal 29-09-2026 — tco5 menu satu butir

**Keputusan work owner tco5**: butir menu ketiga `InboxTreatyContractDescription` **dibuang** — harness itu popup dari
form kontrak (`Section/InputTreatyContract.xml` tombol `List Description` b22196 → harness b22323/b23088). Layar
klausul (`PanelKlausulTahun`) tetap, dibuka tombol `List Description` baris tahun treaty. Halaman pembungkus
`InboxTreatyContractDescription.tsx` tidak lagi dirujuk menu; rutenya di `App.tsx` (suntingan work owner) menunggu.

## Ralat bertanggal 29-09-2026 — tco4 (nol tabel baru) `[keputusan work owner]`

- Tabel `PROPORTIONALARRG` warisan (tipe CAMPUR `[data DBA]`), sequence `PROPORTIONALARRG_SEQ` (lebar 7).
  `RP`/`USD`/`PCT`/`PCTME`/`KURS` teks → dibaca apa adanya dan diurai di Go (titik atau koma), ditulis desimal teks
  bertitik (**OQ-TCO-23**); `TREATYLIMIT`/`COINS_*`/`MORE*` NUMBER.
- Anti-dobel AC 30 dibandingkan di Go: `12.5` dan `12,50` sama nilainya walau teksnya beda.
- **Jejak gugur**.

## ⛔ Ralat bertanggal — 29-09-2026 (penyisiran layar: `Limit MB` 500 karena titik ribuan)

*Temuan.* Penyisiran `GET` seluruh grid klausul — 2.366 grid induk (182 tahun × 13 jenis, `induk=00`) dan 1.822 grid
anak (setiap `reinsTypeId` baris induk), 4.188 panggilan — mendapati tepat **dua** 500 *"gagal memproses permintaan
treaty contract out"*: jenis `10017` `Limit MB` tahun `1000680` dan `1000672`. Perintah audit: `GET …/jenis-klausul?isXol=0|1`
(13 jenis), `GET …/tahun?ukuran=100` sampai habis (182), lalu `GET …/tahun/{id}/klausul?descId={jenis}&induk=00` dan
`…&induk={reinsTypeId}`; kode status dicacah. Jendela: 29-09-2026 sore, backend `:8080` versi `3748a9d`. `[data DEV 29-09-2026 — dibaca executor, SELECT baca-saja atas izin work owner]`: baris induk `PROPORTIONALARRG` tahun `1000680` menyimpan `RP` dan `USD` = `"1.000.000"`,
sedangkan `PCTME` dan `MORERP` baris yang sama `1000000` (`SELECT ID, RP, USD, PCT, PCTME, KURS, TGLUPDATE, …
FROM POOLDATA.PROPORTIONALARRG WHERE TREATYYEARID='1000680' AND TREATYDESCID='10017' AND PARENTREINSTYPEID='00'` → 1
baris). ⚠️ Baris tahun `1000672` **tidak dibaca** (di luar izin baca); bahwa sebabnya sama `[dugaan]` — dipastikan
dengan penyisiran ulang sesudah backend dimuat ulang. `repository.UraiDesimalWarisanTCO` tidak mengenal titik
pemisah ribuan, dan satu baris mematikan seluruh grid jenis itu.

*Perbaikan.* (1) Titik pemisah ribuan DITERIMA hanya dalam bentuk yang tidak mungkin desimal — dua titik atau lebih,
kelompok pertama 1–3 angka tanpa nol di depan, setiap kelompok sesudahnya tepat tiga angka
(`^-?[1-9]\d{0,2}(\.\d{3}){2,}$`); `1.000` (satu titik) tetap desimal, dan `1.00.000`, `1000.000.000`, `1.000.000,5`,
`000.000.000` tetap ditolak — bukan tebakan. Jangkauannya: setiap pemakai `repository.UraiDesimalWarisanTCO` — kolom
`RP`/`USD`/`PCT`/`PCTME`/`KURS` klausul (termasuk pemeriksa dobel `medanSamaKlausulTCO` dan Pct anak), pembaca hilir
`tco_kontrak_hilir.go`, dan `PCT_SHARE` security; **bukan** `TOIDR` kurs (`models.UraiNilaiKursTCO`) dan bukan
masukan layar. Bentuk TULIS tidak berubah (OQ-TCO-23): menyimpan ulang baris seperti itu menulis `1000000` — untuk
baris Limit MB ini praktis tidak terjadi, karena jenisnya `Ditahan` (AC 36) dan tidak dapat disimpan. ⚠️ Nilai bertitik
SATU (`500.000`) tetap dibaca sebagai desimal — **OQ-TCO-27**. Uji `TestUraiDesimalWarisanTCORibuanTitik`. (2) Cabang 500 `jawabGalatTreatyContractOut` kini
**mencatat sebab aslinya** di log backend (`log.Printf`) — sebelumnya galat tak terduga hilang di layar DAN di konsol,
sehingga diagnosanya menuntut membaca DEV langsung. Pemakai tetap mendapat kalimat umum. Uji
`TestGalat500TreatyContractOutMencatatSebabnya`.

## ⛔ Keputusan work owner bertanggal — 30-09-2026 (tampilan layar klausul)

*Butir 4–6 pesan work owner.* Grid `For XOL` (`IsXOL = 1`, b7971) **dibuang** — `GET …/jenis-klausul?isXol=1` di DEV
kosong (penyisiran 29-09-2026: ke-13 jenis ber-`IsXOL = 0`); rute backend tetap. Grid `IsXOL = 0` berjudul
**`Treaty Desc`** (korpus b4880: "For Non XOL"). `Show` membuka **satu** jenis dalam popup (`Modal` bersama) — menekan
jenis lain sesudah popup ditutup menampilkan jenis itu saja. **AC 29 digantikan**: jenis yang tidak sedang terbuka tidak
memegang isian; menutup popup membuang isian yang belum disimpan jenis itu. Nilai desimal grid (`Rp`, `Usd`, `Pct`,
`PctMe`, `CoIns_Min`, `CoIns_Max`, `TreatyLimit`), Total Pct, dan kurs tampil berpemisah ribuan (`formatNumber`,
`DESIMAL_TAK_DIBATASI`); catatan "dihitung server" tidak tampil. Uji: `PanelKlausulTahun.test.ts`
(`alihJenisTunggal`, popup, satu grid), `PanelJenisKlausul.test.ts` (`tampilMedanKlausul`).

## ⛔ Keputusan work owner bertanggal — 30-09-2026 (ReinsType anak Treaty Limit)

*Permintaan: "coba perbaiki reinstype pada treatylimit. ikuti xml nya aja". Jawaban: "12 jenis porsi + induknya".*

**Cacat.** Pilihan `ReinsTypeID` baris anak Treaty Limit memakai daftar induk tiket 02 (RD
`BrowseReinsuranceType_RD_Old_Ljt_id_isnotnull`), dan simpan menolak ID di luarnya. Padahal XML anak berbeda:
`GridTreatyArrTreatyLimitList.xml` b2892 `pxAutoComplete` berdaftar `ReinsTypeList.pxResults` (b2981) dari
pre-activity `TreatyContractSetReinsTypeList` (b2975) — tampil `.CARI2`, ID `.CARI1` (b3019) — **satu-satunya** grid
klausul bersumber begitu (anak jenis lain `D_EnumerationList`, induk semuanya RD tadi). Aktivitas itu **tidak diekspor**.
`[data DEV 30-09-2026 — GET baca-saja lewat backend lokal]`: 118 dari 120 baris anak Treaty Limit ber-ReinsType QS (R/I)
10004, QS (OR) 10028, SPL (OR) 10248, SPL (RI) 10249 — keempatnya di antara dua belas awalan yang RD induk singkirkan; dua
baris sisanya ORS 10007 di bawah induk ORS. Akibatnya Edit anak menampilkan ReinsType kosong dan simpan ditolak 422.

**Keputusan.** Pilihan anak Treaty Limit = baris `REINSURANCETYPE` Flag `active` yang **berawalan** salah satu dari dua
belas awalan (StartsWith — pelengkap tepat NotStartsWith RD induk, tanpa saringan Type) **ditambah ReinsType induknya
sendiri**. Induk Treaty Limit dan anak jenis lain tidak berubah (OQ-TCO-15).

**Yang dibangun.**
- models `AturanKlausul.PilihanReins` = `anak-treaty-limit` pada `TreatyLimitChild` (dikirim ke layar lewat `AturanTampil`).
- repository `LolosSaringanAnakTreatyLimitTCO` (tabel kebenaran) + `DaftarAnakTreatyLimit` (SQL bind: `FLAG = :1 AND
  (ID = :2 OR ID LIKE :3 … :14)`, urut `NOTE, ID`).
- services `JenisReasuransiTreaty.DaftarAnakTreatyLimit` (induk wajib; kosong = `ErrPilihanAnakTreatyLimitKosong`, 503);
  simpan anak Treaty Limit memeriksa ID di daftar anak.
- handler `GET /api/treaty-contract-out/jenis-reasuransi/anak-treaty-limit?induk=`.
- frontend: ReinsType induk dan anak Treaty Limit memakai `PilihSaring` (XML `pxAutoComplete`, dicari pada nama; anak
  menampilkan ID sebagai keterangan — `.CARI1` pyShow true); jenis lain tetap dropdown.
