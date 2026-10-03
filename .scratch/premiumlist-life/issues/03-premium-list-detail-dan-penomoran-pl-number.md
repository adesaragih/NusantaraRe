# 03: Premium List Detail dan penomoran `PL_NUMBER`

**Status:** ready-for-agent

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
| `GetPLNumber_Act` | `ASM-FW-GISFW-WORK-LIFE` / `GETPLNUMBER_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/GetPLNumber_Act.xml` | baca balik nomor via `GetPLandNopolis_sql` |
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

- [ ] `PL_NUMBER` diperoleh dari `PROC_GENERATE_SEQUENCE_NUMBER`; aplikasi **tidak** menyusun format
      nomor sendiri, dan tidak ada pembentuk format di lapisan services. *(AC 10 spec; **ADR-0006**)*
- [ ] Prefix diperoleh lewat **lookup** ke `POOLDATA.KODE_PRODUKSI` (`TYPE='LIFE'`), tidak ditanam
      sebagai konstanta. *(AC 11 spec)*
- [ ] Periode yang dikirim ke procedure berasal dari tiket **02**, bukan dari `time.Now()` mentah.
- [ ] Nomor lahir **sekali** per premium list: submit kedua atas premium list yang sudah bernomor
      **tidak** menggerakkan sequence dan **tidak** mengubah nomor.
- [ ] Empat cabang per `Type` (`QR`, `QP`, `TP`, `TR`) menentukan skema nomor yang dipakai; cabang
      dipilih dari data, bukan dari urutan langkah.
- [ ] Commit terjadi **segera setelah** nomor terbentuk, sehingga lock `SELECT … FOR UPDATE` tidak
      menahan pengguna lain. (**ADR-0015**)
- [ ] Dua submit berurutan menghasilkan dua nomor **berbeda dan berurutan** di bawah beban paralel.
- [ ] Baris `PremiumListDetail` tersimpan utuh dengan seluruh kolom uang, dan **tidak satu pun**
      melewati `float`. *(AC 14 spec; **ADR-0003**)*
- [ ] `PL_NUMBER` yang sudah terbit **terlihat** pengguna dan dapat dibaca kembali lewat API.

### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §12

- [ ] ⚠️ Peserta tersimpan di **`T_PREMIUM_LIST_DETAIL`**; `M_LIFE_PREMIUM_DETAIL` **tidak ditulis**.
      *(AC 33 spec; penyimpangan sadar 1)*
- [ ] Peserta menyimpan **kedua jendela valuasi** — `GROSS_VALUATION_*` (jalur `Q*`) **dan**
      `RETROCESSION_VALUATION_*` (jalur `T*`) — beserta `EFFECTIVE_DATE`, `LAPSE_DATE`, `PERIOD_MM`.
      Test wajib memuat **kedua jalur**. *(AC 38 spec)*
- [ ] `FACTOR` diperlakukan sebagai **desimal** tujuh angka, bukan bilangan bulat. *(AC 39 spec)*
- [ ] ⚠️ Hasil spreading per peserta **dibekukan dan disimpan** di `T_PREMIUM_LIST_SPREADING`, dan
      rincian per reinsurer di `T_PREMIUM_LIST_SPREADING_RETRO`. Perubahan master treaty sesudahnya
      **tidak mengubah** angka yang sudah tersimpan. *(AC 41 spec; penyimpangan sadar 3)*
- [ ] Polis tanpa retrosesi tersimpan dengan **nol baris** spreading dan spreading retro — **bukan**
      kegagalan. *(AC 42 spec)*

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
