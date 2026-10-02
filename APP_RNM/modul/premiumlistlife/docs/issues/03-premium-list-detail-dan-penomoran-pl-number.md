# 03: Premium List Detail dan penomoran `PL_NUMBER`

**Status:** sebagian — uji penomoran paralel, spreading beku, `FACTOR` + kolom uang penuh, dan uji dua jalur valuasi belum ada

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**, 01 (penawaran — tahap Premium hanya terbuka setelah `Confirm` → `Premium`), 02
(periode tutup buku — periode produksi adalah **masukan** procedure penomoran)

## Hasil & nilai pengguna

Sebagai **inputor Life**, saya ingin mengisi rincian premium list untuk penawaran yang sudah
dikonfirmasi dan memperoleh **satu `PL_NUMBER` resmi** untuknya, supaya premium list itu punya
identitas tunggal yang dapat dirujuk seluruh perusahaan — termasuk oleh Claim Life di hilir.
*(User story 9–16 di spec)*

## Area codebase

`internal/handlers` (endpoint isi detail + endpoint submit), `internal/services` (perakitan
`PremiumListDetail`; gerbang "nomor lahir sekali"), `internal/repository` (rantai `KODE_PRODUKSI` →
`PROC_GENERATE_SEQUENCE_NUMBER`), `frontend/` (grid Premium List Detail, tampilan `PL_NUMBER`).

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SubmitPremiumList_Act` | `ASM-FW-GISFW-WORK-LIFE` / `SUBMITPREMIUMLIST_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/SubmitPremiumList_Act.xml` | rantai penomoran (17 langkah) |
| `SavePremiumList_Act` | `ASM-FW-GISFW-WORK-LIFE` / `SAVEPREMIUMLIST_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/SavePremiumList_Act.xml` | perakitan detail (451.849 byte, **15 langkah**) |
| `GetKodeProdLife_SQL` | `ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETKODEPRODLIFE_SQL` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetKodeProdLife_SQL.xml` | `SELECT KODE … FROM POOLDATA.KODE_PRODUKSI WHERE TYPE='LIFE'` |
| `GetSequenceNumber_SQL` | `ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETSEQUENCENUMBER_SQL` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetSequenceNumber_SQL.xml` | `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER(…)` + `COMMIT;` (baris 88) |
| `GetPLNumber_Act` | `ASM-FW-GISFW-WORK-LIFE` / `GETPLNUMBER_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/GetPLNumber_Act.xml` | ⚠️ **RALAT 28-09-2026** — BUKAN "baca balik nomor". Ia mengisi **autocomplete nomor polis yang sudah ada** (`JSON_POLIS.NOPOLIS AS CARI1`), bergerbang `.Type = "TR" \|\| .Type = "TP"` (b450), satu langkah `RDB-List` ber-`REPEAT` sekali. Nol kaitan dengan penerbitan `PL_NUMBER`; rantai penomoran seluruhnya di `SubmitPremiumList_Act.xml` |
| `GetProductDtlPL`, `GetRateProductLife`, `GetRateLifePM` | `ASM-FW-GISFW-…` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/` | data produk & rate |

`[terverifikasi]` **Rantai penomoran** di `SubmitPremiumList_Act` (baris `<pyStepPageReference>`):

| Step | Baris | Deskripsi |
| ---: | ---: | --- |
| 5 | 1376 | Reset value (Temporary) |
| **6–9** | 1783 / 1959 / 2127 / 2295 | **Generate PL Number QR / QP / TP / TR** — empat cabang per `Type` |
| **10** | 2463 | **AMBIL KODE PROD** → `GetKodeProdLife_SQL` |
| **12** | 2810 | **generate MM.YYYY DAN SEQUENCE** → `GetSequenceNumber_SQL` |
| 14 | 3219 | Set COB & PL Number to Pega |
| 17 | 3705 | `Obj-Save` |

`[terverifikasi]` **Perakitan detail** di `SavePremiumList_Act` step **8** "Insert to table detail",
sembilan sub-langkah:

| Sub-step | Deskripsi | Status |
| --- | --- | --- |
| 8.1 | Protect Age | aktif |
| 8.2 | Protect Sum Insured | aktif |
| ~~8.3~~ | ~~Protect Gross Premium~~ | **REMARK** (baris 1891) |
| **8.4** | **Set property `PremiumListDetail`** | aktif |
| ~~8.5–8.7~~ | ~~if proratetype 1 / 2 / 3~~ | **REMARK** (3704 / 3863 / 4050) |
| 8.8 | insert nilai dari data-batch → int | aktif |
| ~~8.9~~ | ~~insert ke tabel detail premium list back up~~ → `SaveMasterLPBackUp` | **REMARK** (5887) |

`[terverifikasi]` Langkah lain yang **REMARK**: step 11 `Call SetRateLIfePremium_Act` "Perhitungan
rate - net premium" (6411) dan step 13 `Call SpreadingLife_Act` (6642).

`[data DBA]` `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER(p_class, p_jenis, p_proddate, OUT p_bulan,
OUT p_seq_number)` — sequence per `(class, jenis, tahun)` di `GENERATE_SEQUENCE_NUMBER` (PK
komposit), memakai `SELECT … FOR UPDATE`; keluaran `LPAD(seq, 5, '0')`; periode digulir lewat
`POOLDATA.TANGGAL_CLOSING`.

## ADR terkait

**ADR-0006** (penomoran **wajib** lewat `PROC_GENERATE_SEQUENCE_NUMBER`; aplikasi tidak menyusun
format sendiri), **ADR-0003** (uang non-float), **ADR-0007** (jejak audit), **ADR-0015** (batas
transaksi dipegang Go; commit segera setelah nomor terbentuk agar lock `FOR UPDATE` lekas lepas).

## Acceptance criteria

- [x] ⚠️ **RALAT 28-09-2026 — AC ini keliru terhadap XML, dan diperbaiki di sini.** Bunyi lamanya:
      *"`PL_NUMBER` diperoleh dari `PROC_GENERATE_SEQUENCE_NUMBER`; aplikasi **tidak** menyusun format
      nomor sendiri."* `SubmitPremiumList_Act` baris **3126** membuktikan sebaliknya — yang menyusun
      bentuk nomor adalah **activity**-nya, bukan procedure-nya:

      ```
      .PremiumListSummary.PL_NUMBER =
          ParamSeq.HASIL3 + InputData.CARI20 + pyWorkPage.BusinessCode
          + "." + ParamSeq.HASIL1 + "." + ParamSeq.HASIL2
      ```

      Procedure hanya mengembalikan **dua** nilai: `HASIL1` (periode `MM.YYYY`) dan `HASIL2` (urut).
      Awalan datang dari `KODE_PRODUKSI` (langkah 10) dan huruf tipe dari langkah 6–9.

      **Bunyi yang benar:** urut dan periode diperoleh dari penghitung
      `GENERATE_SEQUENCE_NUMBER`; bentuk nomornya dirakit di lapisan **`models`** (murni, dapat diuji
      tanpa Oracle) — **tidak** di `services`, **tidak** di handler, dan **tidak** di layar.
      *(**ADR-0006** tetap: penomoran tidak boleh dikarang di tempat yang tersebar.)* — bukti: `models/polis_nomor.go:NomorPL`; uji `TestRakitNomorPLBerbentukSepertiB3126`, `TestNolPembentukBentukNomorDiLapisanLayanan`
- [x] Prefix diperoleh lewat **lookup** ke `POOLDATA.KODE_PRODUKSI` (`TYPE='LIFE'`), tidak ditanam
      sebagai konstanta. *(AC 11 spec)* — bukti: `repository/penomor.go:AwalanProduksi`
- [x] ⚠️ **RALAT 28-09-2026 — periode TIDAK dikirim ke procedure.** `pyMemo` check-in terakhir
      `SubmitPremiumList_Act` (20260122) berbunyi harfiah **"buang ParamSeq.CARI3"**, dan memang
      `ParamSeq.CARI3` tidak diset di satu langkah pun; `TO_DATE({ParamSeq.CARI3},'DD/MM/YYYY')` di
      `GetSequenceNumber_SQL` karena itu menerima **NULL**, dan procedure menentukan periodenya
      sendiri dari `SYSDATE`.

      **Bunyi yang benar:** karena procedure tidak dipanggil, periodenya dihitung di Go — dengan
      aturan tiket **02** (`models.PeriodeProduksi`), bukan `time.Now()` mentah. Aturannya **satu**,
      dipakai penomoran klaim maupun premium list. — bukti: `repository/penomor.go:HitungPeriodeNomor`; uji `TestPeriodeNomorMengikutiHariTutupBuku`, `TestPeriodeNomorTidakMelompatiBulanPendek`
- [x] Nomor lahir **sekali** per premium list: submit kedua atas premium list yang sudah bernomor
      **tidak** menggerakkan sequence dan **tidak** mengubah nomor. Gerbangnya berdiri di **dua**
      tempat — di layanan (dibaca sebelum penghitung disentuh) dan di kalimat `WHERE … PL_NUMBER IS
      NULL` query penulisnya. — bukti: `services/polis_nomor.go:terbitkanDalam`, `repository/polis_nomor.go:sqlTulisNomorPL`; uji `TestPenghitungTidakTersentuhBilaSudahBernomor`, `TestTulisNomorPLTidakDapatMenimpaNomorYangSudahAda`
- [x] Empat cabang per `Type` (`QR`, `QP`, `TP`, `TR`) menentukan skema nomor yang dipakai; cabang
      dipilih dari data, bukan dari urutan langkah.

      ⚠️ **TAMBAHAN 28-09-2026, dan ini yang paling mudah dirusak "perbaikan".** Keempat tipe itu
      berbagi **SATU** penghitung. `SubmitPremiumList_Act.xml` baris **2717** memakai teks
      **harfiah** `"QR/QP/TP/TR"` sebagai
      bagian kunci `JENIS` (`ParamSeq.CARI2 = ParamSeq.HASIL3+"QR/QP/TP/TR"`); tipe yang sedang
      berjalan hanya masuk ke **nomornya**, tidak ke kuncinya. Memecahnya menjadi empat penghitung
      menerbitkan empat deret yang masing-masing mulai dari 1 — dan setiap nomor baru bertabrakan
      dengan nomor lama. — bukti: `models/polis_nomor.go:KodeTipePL`, `models/polis_nomor.go:JenisPenghitungPL`; uji `TestKodeTipePLEmpatCabangDariData`, `TestSatuPenghitungUntukEmpatTipe`
- [ ] Commit terjadi **segera setelah** nomor terbentuk, sehingga lock `SELECT … FOR UPDATE` tidak
      menahan pengguna lain. (**ADR-0015**) Transaksinya memuat **hanya** pengambilan nomor dan
      penyimpanannya — sengaja berbeda dari pendaftaran klaim, yang memegang kunci sampai seluruh
      pendaftaran selesai. ⚠️ Nol `COMMIT` di teks SQL (**ADR-U-0029**): batasnya dipegang Go,
      meski `GetSequenceNumber_SQL` baris 88 punya satu. — belum: benar untuk `Terbitkan` (uji `TestNolCommitDiQueryNomor`), tetapi jalur Confirm/Submit memanggil `terbitkanDalam` di transaksi `simpanDalam` yang juga memuat rekap + salinan warisan + penutupan, jadi kunci penghitung ditahan lebih lama
- [ ] Dua submit berurutan menghasilkan dua nomor **berbeda dan berurutan** di bawah beban paralel. — belum: nol uji beban paralel terhadap Oracle
- [ ] Baris `PremiumListDetail` tersimpan utuh dengan seluruh kolom uang, dan **tidak satu pun**
      melewati `float`. *(AC 14 spec; **ADR-0003**)* — belum: unggahan mengisi 32 dari 44 kolom uang `052` (`FACTOR`, `RATE`, `SUM_AT_RISK_*`, `DEDUCTION`, `RI_ADMIN_FEE*` dll. tidak diisi); uang memang dikirim sebagai teks (uji `TestUangDikirimSebagaiTeks`)
- [ ] `PL_NUMBER` yang sudah terbit **terlihat** pengguna dan dapat dibaca kembali lewat API
      (`GET /api/polis-life/{id}` dan `GET /api/polis-life/{id}/nomor`). — belum: `GET /api/polis-life/{id}/nomor` dibuang (nol pemanggil); nomor terbaca lewat `GET /api/polis-life/{id}` dan tampil di `PremiumListDetail.tsx`

### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §12

- [ ] ⚠️ Peserta tersimpan di **`T_PREMIUM_LIST_DETAIL`**; `M_LIFE_PREMIUM_DETAIL` **tidak ditulis**.
      *(AC 33 spec; penyimpangan sadar 1)* — belum: peserta memang di `T_PREMIUM_LIST_DETAIL` (`repository/polis_unggah.go:SisipPeserta`), tetapi pl2 membalik larangannya — `M_LIFE_PREMIUM_DETAIL` ditulis (`repository/polis_warisan.go:Ganti`)
- [ ] Peserta menyimpan **kedua jendela valuasi** — `GROSS_VALUATION_*` (jalur `Q*`) **dan**
      `RETROCESSION_VALUATION_*` (jalur `T*`) — beserta `EFFECTIVE_DATE`, `LAPSE_DATE`, `PERIOD_MM`.
      Test wajib memuat **kedua jalur**. *(AC 38 spec)* — belum: kedua jendela ditulis unggahan (`repository/polis_unggah.go:kolomTanggalPeserta`), tetapi nol uji yang memuat jalur `Q*` dan `T*` terpisah
- [ ] `FACTOR` diperlakukan sebagai **desimal** tujuh angka, bukan bilangan bulat. *(AC 39 spec)* — belum: dibaca sebagai desimal (`models/polis_detail.go:KolomGridPeserta`, kolom `NUMBER(38,8)`), tetapi tidak ada jalur yang mengisi `FACTOR` dan nol uji nilai tujuh desimal
- [ ] ⚠️ Hasil spreading per peserta **dibekukan dan disimpan** di `T_PREMIUM_LIST_SPREADING`, dan
      rincian per reinsurer di `T_PREMIUM_LIST_SPREADING_RETRO`. Perubahan master treaty sesudahnya
      **tidak mengubah** angka yang sudah tersimpan. *(AC 41 spec; penyimpangan sadar 3)* — belum: nol penulis `T_PREMIUM_LIST_SPREADING` / `T_PREMIUM_LIST_SPREADING_RETRO`
- [ ] Polis tanpa retrosesi tersimpan dengan **nol baris** spreading dan spreading retro — **bukan**
      kegagalan. *(AC 42 spec)* — belum: penulis spreading belum ada, jadi nol baris belum dibedakan dari belum ditulis

## Blocker

**Tidak ada pemblokir.** Satu catatan yang **tidak** memblokir tiket ini: **OQ-068** (di mana baris
detail new business mendarat di `M_LIFE_PREMIUM_DETAIL`) — memblokir **tiket 08**, bukan tiket ini.

## Catatan — kode mati yang tidak dimigrasikan

`[terverifikasi]` **Jangan** replikasi: `Protect Gross Premium` (8.3), ketiga cabang `proratetype`
(8.5–8.7), backup detail (8.9), `SetRateLIfePremium_Act` (11), `SpreadingLife_Act` (13) — seluruhnya
`<pyStepsBlockName>//`.

⚠️ `[terverifikasi]` **Asimetri rujukan**: `SaveMasterLPBackUp` dirujuk `<RequestType>` di
`SavePremiumList_Act` baris 5931, tetapi **tidak ada berkas rule dengan nama itu di korpus** —
rujukan mati, sejalan dengan langkahnya yang sudah REMARK.

⚠️ **OQ-066**: penanda `<pyStepsBlockName>` **tidak dapat dipercaya sendirian** di modul Life.
Konfirmasikan ke work owner sebelum menyimpulkan langkah lain hidup atau mati.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

## Implementasi

**Dikerjakan 28-09-2026.** `main`, sesudah `8f4df09` (tiket 02).

### Yang dikirim

| Lapisan | Berkas | Isi |
| --- | --- | --- |
| models | `polis_nomor.go` | bentuk `PL_NUMBER` — murni, nol Oracle: `KodeTipePL` (empat cabang dari data), `JenisPenghitungPL` (teks harfiah `QR/QP/TP/TR`), `ClassPenghitungPL`, `PeriodeNomorPL` (`MM.YYYY` → `MM.YY`), `RakitNomorPL`, `NomorPL` |
| models | `polis_detail.go` | `KolomGridPeserta` — 38 kolom urut `PL_Detail_Sec`, beserta `MedanGridTanpaKolom` |
| models | `polis_periode.go` | `DiJakarta` — satu sumber zona untuk seluruh repo |
| repository | `penomor.go` | `Penomor` (tipe baru); `HitungPeriodeNomor` kini mendelegasikan pergeseran bulan ke `models.PeriodeProduksi` |
| repository | `polis_nomor.go` | `Ringkas` (satu query, `LEFT JOIN`), `Keadaan`, `Identitas`, `TulisNomor` |
| repository | `polis_detail.go` | grid peserta, SQL dirakit dari daftar kolom, seluruh angka lewat `TM9` |
| services | `polis_nomor.go` | `NomorPremiumList.Terbitkan` / `.Baca` — gerbang lahir-sekali, gerbang kasus tertutup, transaksi pendek |
| services | `polis_detail.go` | `DetailPolis.Kepala` / `.Peserta` — bentuk JSON milik lapisan ini |
| handlers | `rute_premiumlist.go` | `GET {id}`, `GET {id}/peserta`, `GET {id}/nomor`, `POST {id}/nomor` |
| frontend | `PremiumListDetail.tsx` | grid + nomor; daftar kolom datang dari server |
| frontend | `InputOffer.tsx` | `judulKeputusan` — judul mengikuti tahapnya |

### Rantai penomoran, sebagaimana ditiru

```
gerbang kasus tertutup (T_WORK_POLIS)      <- di LUAR transaksi
gerbang lahir-sekali   (MIN/MAX PL_NUMBER) <- sebelum penghitung disentuh
KODE_PRODUKSI TYPE='LIFE'      -> awalan            (langkah 10)
TANGGAL_CLOSING                -> hari tutup buku
aturan tiket 02                -> periode MM.YYYY   (langkah 12, keluaran proc)
GENERATE_SEQUENCE_NUMBER       -> urut              (langkah 12, FOR UPDATE)
  kunci (CLASS, JENIS, TAHUN) = (pxObjClass, awalan+"QR/QP/TP/TR", tahun)
models.NomorPL                 -> awalan+tipe+COB+"."+MM.YY+"."+urut5  (b3126)
UPDATE ... WHERE PL_NUMBER IS NULL
```

### Cacat yang ditemukan dan diperbaiki dalam kode yang SUDAH ada

`repository.HitungPeriodeNomor` (penomoran klaim, butir o1) menggeser bulan dengan
`saat.AddDate(0,1,0)`. Go **melimpahkan** tanggal yang tidak ada: **31 Januari + 1 bulan = 3 Maret**,
sehingga periodenya `03.2026` dan **Februari terlewat sama sekali**. Oracle `ADD_MONTHS` justru
**menjepit** ke akhir bulan (29 Februari) — `02.2026`. Ia menyala pada tanggal **29–31** bulan yang
penggantinya lebih pendek, dan hasilnya nomor yang terbukukan ke bulan yang salah **tanpa satu pun
galat**. Terbukti merah lebih dahulu atas implementasi lama
(`TestPeriodeNomorTidakMelompatiBulanPendek`, tujuh kasus), lalu hijau.

Sekaligus: harinya kini dibaca di zona **Jakarta** (sama dengan `SYSDATE` server) — sebelumnya di
zona `saat` sendiri, sehingga cutover dan pergeseran dapat berselisih sehari.

### Yang TIDAK dikerjakan di tiket ini, dan sebabnya

| Butir | Sebab |
| --- | --- |
| Perakitan baris `PremiumListDetail` (`SavePremiumList_Act` step 8) | jalur **tulis** peserta adalah unggahan CSV — **tiket 04** |
| Kedua jendela valuasi, `FACTOR` desimal, spreading dibekukan | ditulis saat peserta **disimpan** — **tiket 04/05a**; tiket ini hanya membacanya |
| `ProtectProductName_Act`, `ChooseProdName`, `SOB_Harness`, `Retro_Section`, `SecurityReinsurer_Section` | lima pemilih pada kepala polis, bukan pada penomoran maupun grid — menyusul bersama layar kepala polis |

### Dua permintaan bersamaan atas polis yang SAMA

Gerbang lahir-sekali dibaca tanpa kunci baris, jadi dua permintaan dapat sama-sama melewatinya.
Yang kedua lalu menulis **nol** baris (`WHERE … PL_NUMBER IS NULL` sudah tidak cocok), transaksinya
**dibatalkan**, dan kenaikan penghitung ikut batal — **nol nomor terbakar**.

⛔ Sebab nol barisnya **dibedakan**, tidak ditebak: nol baris peserta menjawab
`ErrPolisTanpaPeserta` (*"unggah rincian peserta lebih dahulu"*), sedangkan baris yang ada tetapi
sudah bernomor menjawab `ErrNomorPLTerbitBersamaan` (*"baca ulang nomornya"*), keduanya **409**.
Menjawab keduanya dengan kalimat yang sama mengirim orang mencari peserta yang sebenarnya ada.
Dikunci `TestDuaSebabNolBarisDibedakan`.

### Butir `[terbuka]` yang LAHIR di tiket ini

- ~~**OQ-PL-01 `CLASS` penghitung.**~~ ✅ **DITUTUP 28-09-2026 — lihat blok di bawah.**

  Bunyi aslinya dipertahankan utuh, tidak dihapus:

  > **OQ-PL-01 `CLASS` penghitung.** `ParamSeq.CARI1 = pyWorkPage.pxObjClass` (b2670) adalah kelas
  **konkret saat berjalan**. Korpus memuat dua: `ASM-FW-GISFW-Work-LIFE` (tempat seluruh rule
  didefinisikan; 107 rujukan `pyActivityClass` di `ShowLifePremiumDetail`) dan
  `RNM-FW-LIFEFW-Work-LIFE` (8 rujukan di section yang sama, 394 di korpus). Yang dipakai
  `ASM-FW-GISFW-Work-LIFE`. **Bila baris `GENERATE_SEQUENCE_NUMBER` nyata bertuliskan yang lain,
  penomoran mulai dari satu dan setiap nomor baru bertabrakan dengan nomor lama.** Dikunci uji
  (`TestClassPenghitungPLTidakBergeserDiam`). **Minta konfirmasi DBA sebelum dipakai di lingkungan
  mana pun yang datanya nyata.**
  **Penutupan 28-09-2026 — `[data DEV, agregat]`.** Katalog instance **pengembangan** dibaca
  (agregat per kelas dan jenis, nol baris disalin ke artefak mana pun):

  | Modul | `CLASS` | `JENIS` | Catatan |
  | --- | --- | --- | --- |
  | **PremiumList Life** | `ASM-FW-GISFW-Work-LIFE` | `RNML-QR/QP/TP/TR` | dua baris — tahun 2025 dan 2026; `NO_SEQ` tertinggi 39 |
  | Claim Life | — | `RNML-K` | `NO_SEQ` 31 |
  | Komite | — | `RNML-A` | `NO_SEQ` 3 |
  | NB | — | `RNM-QR/QP/TP` | pembanding |

  ⛔ **`RNM-FW-LIFEFW-Work-LIFE` TIDAK ADA di tabel itu sama sekali.** Kelas konkret saat berjalan
  karena itu `ASM-FW-GISFW-Work-LIFE` — persis yang sudah dipakai kode, dan kekhawatiran "penomoran
  mulai dari satu" tidak terwujud.

  ⛔ **`JENIS` ikut TERBUKTI, dan bukan hanya kelasnya.** Nilai nyatanya `RNML-QR/QP/TP/TR` —
  yaitu `HASIL3` (awalan `RNML-` dari `KODE_PRODUKSI`) ditempel di depan teks harfiah
  `"QR/QP/TP/TR"`, persis susunan `SubmitPremiumList_Act.xml` b2717. **Satu baris penghitung untuk
  keempat tipe**, dan `NO_SEQ` 39 membuktikannya: empat deret terpisah akan menampakkan empat baris,
  bukan dua (satu per tahun).

  ⚠️ Yang TIDAK berubah: konstanta `models.ClassPenghitungPL` tetap `ASM-FW-GISFW-Work-LIFE`, dan
  ujinya (`TestClassPenghitungPLTidakBergeserDiam`) tetap menguncinya. Yang berubah hanya
  **dasarnya** — dari dugaan berbasis cacah rujukan korpus menjadi bacaan dari tabelnya sendiri.

  ⚠️ `[data DEV]`, bukan produksi. Nol baris, nol nama, nol nomor polis disalin ke repositori ini.

- **OQ-PL-02 `RetrocadedShare` vs `RETROCEDED_SHARE`.** Keduanya tampil di grid `PL_Detail_Sec` yang
  sama, jadi keduanya medan berbeda; hanya yang kedua punya kolom. Mana yang mana — belum terjawab.
- **OQ-PL-03 `REINSTYPENAME`.** Tampil di `PL_Detail_Sec`, nol kolom di migrasi 050–056, dan nol
  rule lain di korpus PremiumList Life yang menyebutnya.
- **OQ-PL-04 tempat tinggal `PL_NUMBER`.** Kolomnya hanya ada di `T_PREMIUM_LIST_DETAIL` (052);
  `_SUMMARY` (055) dan header (051) tidak punya. Akibatnya polis **tanpa peserta** belum dapat
  dinomori. Itu sejalan dengan urutan aslinya (submit sesudah rincian terisi), tetapi bila polis
  kelak perlu bernomor lebih dahulu, kolom `PL_NUMBER` di header-lah jawabannya — **keputusan
  skema**, dan migrasi baru hanya dari keputusan yang tercatat.

### Langkah 15 `SubmitPremiumList_Act` — `[terbuka]`, tidak ditiru

b3413 `@replaceAll(.PremiumListSummary.PL_NUMBER, Local.CurrentMMYY, Local.NextMMYY)`.
`Local.CurrentMMYY` berbentuk `.MM.YYYY.` (b1142, tahun **empat** angka), sedangkan `PL_NUMBER`
memuat `.MM.YY.` (tahun **dua** angka, b3021/b3075) — teks yang dicari **tidak pernah ada** di
dalam nomornya, sehingga langkah itu tidak mengubah apa pun. `Local.NextMMYY` (b1163) memperkuatnya:
ia memformat `\"dd\"`, yaitu **hari**, bukan bulan. Tidak ditiru; dilaporkan.

> ✅ **DIPUTUSKAN 28-09-2026 — pl7** `[veto work owner]` *(brief GILIRAN-10 §0)*: langkah 15 adalah
> **no-op di produksi** — pola `.MM.YYYY.` (b1142) tidak pernah ada di dalam `PL_NUMBER` berformat
> `.MM.YY.` (b3021/b3075/b3126). **Tidak ditiru**; residu warisan, dengan bukti b1142/b3413/b3126.
> Bila kelak pemilik Pega menyatakan maksudnya *(menggeser periode nomor saat tutup buku)*, itu
> **keputusan baru**, bukan replikasi. Dikunci `TestLangkah15TidakDitiru` (`services/polis_summary_test.go`).

### Temuan `/code-review` yang diperbaiki sebelum commit

Tinjauan dua sumbu dijalankan atas titik tetap `8f4df09`. Yang diperbaiki:

| Sumbu | Temuan | Tindakan |
| --- | --- | --- |
| Spec (c) | ⛔ **Lubang nyata.** `MIN`/`MAX` **melewati NULL**, jadi polis yang separuh barisnya bernomor terbaca *"sudah bernomor"*, gerbang lahir-sekali kembali lebih awal, dan baris yang kosong **tidak pernah terisi**. Komentar `sqlKeadaanNomorPL` sudah menyebut selisih `COUNT(*)` vs `COUNT(PL_NUMBER)` sebagai penunjuknya — tetapi nol kode pernah **membandingkannya**. Terjangkau begitu tiket 04 menambah peserta sesudah penomoran. | `KeadaanNomorPL.Utuh()` lahir; gerbangnya kini **menyembuhkan** — baris yang tertinggal diberi nomor yang **sudah ada**, penghitung tetap tidak tersentuh. Dikunci dua uji, dan penjaganya **dibuktikan merah** lebih dahulu. |
| Standards 1 | Sensus tanpa label `[terverifikasi]`, tanpa perintah audit, tanpa cara kedua (CLAUDE.md §4 butir 9, §4a). | Ketiga sensus (40 medan · 17 langkah · rujukan kelas) kini berlabel, dengan **jendela disebut**, **perintah audit**, dan **dua cara** yang keduanya sepakat. Jebakannya ikut disebut: kemunculan **mentah** medan grid **109**, bukan 40. |
| Standards 2 | ADR-U-0027 **salah dirujuk** untuk pemetaan NULL → teks kosong; ADR itu tentang kolom nullable dan wajib-isi di kode. | Rujukan **dicabut** di kedua tempat, dan pencabutannya ditulis. Aturannya berdiri sendiri. |
| Standards 3 | `App.tsx` dan `InputOffer.tsx` tertulis ulang **LF→CRLF** — 210 dan 360 baris berubah untuk 14 dan 15 baris nyata. | Dikembalikan ke LF; seluruh berkas diperiksa ulang. Diff kini sebesar perubahannya. |
| Standards (smell) | `fmtDesimalPeserta` / `fmtTanggalPeserta` menyalin `fmtDesimal` / `fmtTanggalOracle` di paket yang **sama**. | Dibuang; yang sudah ada dipinjam. |
| Standards (smell) | `cacahPesertaDalamTx` menyalin `sqlCacahPeserta` utuh. | Dibuang; query yang sudah ada dipinjam. |
| Standards (smell) | `KolomPesertaBulat` menempuh cabang yang **persis sama** dengan `KolomPesertaUang`. | Disatukan menjadi `KolomPesertaAngka` — sejajar dengan `golonganKolom` yang sudah ada. |
| Standards + Spec | `GET /api/polis-life/{id}/nomor` → `Baca` → `bacaNomorPL` **nol pemanggil**; `DenganJam` juga. | **Dibuang seluruhnya.** Nomornya dibaca lewat `GET /api/polis-life/{id}` yang sudah membawanya. *(Rute tanpa pemanggil adalah pola cacat yang sudah berulang di repo ini.)* |
| Spec (a) | `MedanTanpaKolom()` mengaku *"DIKIRIM KE LAYAR"* padahal nol rute menyajikannya. | Kini benar-benar ikut di `KepalaPolis`, dan **tampil** di layar sebagai daftar terlipat. |

**Yang TIDAK diambil:** pemindahan `PohonKlaim`→`Penomor` dan perbaikan `HitungPeriodeNomor`
disebut *scope creep* karena mengubah keluaran modul klaim yang sudah jalan. Dipertahankan:
membangun penomoran PL di atas fungsi yang **diketahui** melompati Februari berarti mengirim cacat
yang sama ke modul kedua. Perubahan keluarannya dinyatakan di atas, bukan disembunyikan.

**Yang MASIH terbuka:** AC *"dua submit berurutan menghasilkan dua nomor berbeda dan berurutan di
bawah beban paralel"* belum punya uji — ia menuntut Oracle sungguhan. Rancangannya
(`SELECT … FOR UPDATE` + transaksi pendek) masuk akal, tetapi **belum dibuktikan**; dicatat di sini
alih-alih dicentang.

### Verifikasi

```
gofmt -l .              bersih
go vet ./...            bersih
go vet -tags db ./...   bersih
go test ./... -tags db  463 PASS · 0 FAIL · 38 SKIP
npx tsc --noEmit        bersih
npx vitest run          322 PASS (27 berkas)
npm run build           bersih
```

Nol migrasi baru. Nol procedure dipanggil. Nol `COMMIT` di teks SQL.

## Implementasi — tiket 03 bagian 2 (01 Oktober 2026): data polis layar Input Premium Detail

`[keputusan work owner 01-10-2026]` bagian 1–3 layar `Section/ShowLifePremiumDetail.xml` lebih
dahulu; `Save Data` (`Calculate1_Act` + `SavePremiumList_Act`), Generate Excel, dan ringkasan per mata
uang menyusul.

| Sel Pega | Kode |
| --- | --- |
| Choose Product Name → `ChooseProdName` / `BrowseProductForNB_Life` (Ceding kasus) → `SetProdNametoPolis` | `GET /api/polis-life/{id}/cari-produk` (`PRODUCT_LIFE`) |
| Type* (QR/QP/TP/TR), Premium Payment Method* (kode 1/2/3 dari `Calculate1_Act`; teks 1 Single / 2 Annually / 3 Others — keputusan work owner 01-10-2026) | dropdown, pilihan dari server |
| R/I SLIP RNM No.* (TP/TR) → `GetPLandNopolis_sql` (NOPOLIS RNML-Q / RNML-F) | `GET /api/polis-life/cari-rislip` (`T_PREMIUM_LIST.NO_POLIS`, JSON_POLIS dibuang) |
| Marketing Officer* → `BrowseMarketingOfficer_RD` (MOStatus = 1) | `GET /api/polis-life/cari-marketing` (`MARKETINGOFFICER`) |
| Annuity Interest*, Premium Refund Factor* | desimal sebagai teks |
| Billing Name (TP/TR) / Retrocessionaire → `BrowseCedingCoLife_RD`; `setSecurityReinsurer_act` | `cari-ceding`; Retrocessionaire otomatis untuk L0000141 / L0000135 / L0000137 |
| Data penawaran, tanggal-tanggal, WPC, Status | ditampilkan read-only |

Simpan: `PUT /api/polis-life/{id}/data-polis`, hanya di tahap Input Premium Detail, menulis kolom yang
sudah ada di `T_PREMIUM_LIST` (nol migrasi). ⚠️ `[belum terverifikasi]` kolom `PRODUCT_LIFE`
(CEDINGID, SOBID, SOBNAME, POLICYHODER, POLICYHODERNAME, INWARDNAME) dan `MARKETINGOFFICER`
(CLIENTID, MOSTATUS) — nama properti RD; konfirmasi DBA sebelum DEV. Teks
dropdown Type tidak ada di korpus.

**`Save Data` = `SavePremiumList_Act` saja `[keputusan work owner 01-10-2026]`** — `Calculate1_Act`
tidak dijalankan. Langkah hidup yang dibawa: 6–7 batas produk (`PRODUCTINWARD_LIFE WHERE ID =
ProductNameID`), 8.1 Protect Age (`ENTRY_AGE` vs MINAGE/MAXAGE), 8.2 Protect Sum Insured (dilewati
TP/TR ber-R/I SLIP `RNML-FL`), 9 pesan, 15 simpan (`WithErrors=true`: data tetap tersimpan, pesan
tampil sebagai peringatan). Langkah 4–5 sudah di unggah CSV (tiket 04); 8.4/8.8 di "Simpan
permanen"; 10/12 di layar Summary; 8.3, 8.5–8.7, 8.9, 11, 13 ter-remark. ⚠️ `[belum terverifikasi]`
kolom MINSUMINSURED/MAXSUMINSURED.
