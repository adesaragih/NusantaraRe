# Tiket - Penawaran dan Endorsemen Jiwa

> Dokumen ini memuat **badan tiket lengkap**, disusun per modul lalu per nomor.
> Disusun 25 September 2026 dari berkas tiket proyek migrasi Nusantara Re.

## Matriks status

| Modul | Tiket | Siap | Tertahan | needs-info | wontfix | Lain |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| PremiumList Life | **9** | 7 | 0 | 0 | 2 | 0 |
| Endorsement Life | **12** | 11 | 0 | 0 | 1 | 0 |
| **Jumlah** | **21** | **18** | **0** | **0** | **3** | **0** |

---

# PremiumList Life

Jumlah tiket: **9**

## PremiumList Life - 01 - Penawaran Life — Confirm / Reject / Decline, dan percabangan Offer / Premium

**Status:** ready-for-agent

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**

#### Hasil & nilai pengguna

Sebagai **inputor Life**, saya ingin mencatat penawaran dari ceding dan menyatakan keputusan
**Confirm**, **Reject**, atau **Decline** atasnya, supaya setiap penawaran punya jejak keputusan yang
jelas dan alur berikutnya (berhenti di penawaran, atau lanjut ke premium list) ditentukan secara
sadar — bukan oleh aturan tersembunyi. *(User story 1–8 di spec)*

#### Area codebase

`internal/handlers` (endpoint buat/ubah penawaran + endpoint keputusan), `internal/services`
(transisi tahap; penentuan Offer/Premium), `internal/repository` (tulis rekam offer JSON; baca balik
id offer), `frontend/` (layar Input Offer, tombol keputusan, tampilan tahap berjalan).

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InputPolicyHolder` | `ASM-FW-GISFW-WORK-LIFE` / `INPUTPOLICYHOLDER` / `RULE-OBJ-FLOW` | `PremiumList Life/InputPolicyHolder.xml` | titik masuk flow (106.880 byte) |
| `InputOfferLife_ACT` | `ASM-FW-GISFW-WORK-LIFE` / `INPUTOFFERLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/InputOfferLife_ACT.xml` | simpan penawaran (131.691 byte, **8 langkah**) |
| `SaveOfferJsonLife_SQL` | `ASM-FW-GISFW-INT-OFFERJSON` / `ASM!SAVEOFFERJSONLIFE_SQL` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/SaveOfferJsonLife_SQL.xml` | `POOLDATA.INSERTJSONOFFERLIFE(...)` + `COMMIT;` |
| `GetIdOffer_SQL` | `ASM-FW-GISFW-INT-OFFERJSON` / `ASM!GETIDOFFER_SQL` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetIdOffer_SQL.xml` | baca balik id offer |
| `IsLifeAccepted` | `ASM-FW-GISFW-WORK-LIFE` / `ISLIFEACCEPTED` / `RULE-OBJ-DECISIONTABLE` | `PremiumList Life/DecisionTable/IsLifeAccepted.xml` | nilai keluaran `Confirm` / `Decline` / `Reject` |
| `IsFlagOnGoingPolicy` | `ASM-FW-GISFW-WORK-LIFE` / `ISFLAGONGOINGPOLICY` / `RULE-OBJ-DECISIONTABLE` | `PremiumList Life/DecisionTable/IsFlagOnGoingPolicy.xml` | nilai keluaran `Decline` / `Offer` / `Premium` |
| `setNoOffer_Act`, `setCeding_act`, `setPolicyHolder_act`, `setSOB_act`, `setSecurityReinsurer_act` | `ASM-FW-GISFW-WORK-LIFE` / `…` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/` | pengisian field penawaran |

`[terverifikasi]` **Peta 8 langkah `InputOfferLife_ACT`** (nomor dari `<pyStepPageReference>`, yang
muncul **sesudah** `<pyStepsActivityName>` langkahnya):

| Step | Langkah | Catatan |
| ---: | --- | --- |
| ~~1~~ | ~~`Call ASMForceCaseClose`~~ "Jika status Decline / Closed" | **REMARK** (`<pyStepsBlockName>//`, baris 354) |
| 2 | `Property-Set` | |
| 3 | `Property-Set` "Pega to json_offer_life" | rakit payload |
| **4** | `RDB-List` "Insert to table json_offer_life" | → `SaveOfferJsonLife_SQL` |
| 5 | `RDB-List` "Get ID from json_offer_life" | → `GetIdOffer_SQL` |
| 6, 7 | `Property-Set` | |
| 8 | `Obj-Save` | |

`[keputusan work owner]` Keputusan `Confirm`/`Reject`/`Decline` dibuat **manual** oleh inputor/admin.
Kedua DecisionTable mengekspor **nol baris keputusan** — yang direplikasi adalah **akibat**
keputusan, bukan formula yang memilihnya.

`[keputusan work owner]` `IsFlagOnGoingPolicy`: **`1` = Offer** (berhenti di tahap penawaran),
**`2` = Premium** (lanjut Input Premium List Detail).

#### ADR terkait

**ADR-U-0007** (jejak audit setiap transisi), **ADR-U-0003** (uang non-float — nilai penawaran),
**ADR-U-0009** (migrasi penuh), **ADR-U-0001** (batas konteks — penawaran adalah hulu Claim Life).

#### Acceptance criteria

- [ ] `Confirm` pada tahap penawaran melanjutkan case ke penentuan Offer/Premium. *(AC 1 spec)*
- [ ] `Reject` pada tahap mana pun **mengembalikan** case ke layar Input Offer — bukan menutupnya,
      bukan memajukannya. *(AC 2 spec)*
- [ ] `Decline` pada tahap mana pun **menutup** case; case tertutup tidak dapat dilanjutkan maupun
      diputuskan ulang. *(AC 3 spec)*
- [ ] Keluaran **Offer** menghentikan siklus di tahap penawaran, dan penawaran **tetap tersimpan**
      serta dapat dibaca kembali. *(AC 4 spec)*
- [ ] Keluaran **Premium** membuka tahap Input Premium List Detail. *(AC 5 spec)*
- [ ] Tidak ada aturan otomatis yang menetapkan keputusan; keputusan selalu datang dari tindakan
      pengguna. *(AC 6 spec)*
- [ ] Setiap transisi tahap menulis jejak audit: siapa, kapan, dari tahap apa ke tahap apa.
      (**ADR-U-0007**)
- [ ] Rekam penawaran ditulis lewat `INSERTJSONOFFERLIFE` (upsert berkunci `IDPEGA` + `STATUS`),
      lalu id offer dibaca kembali — keputusan berikutnya memakai id itu.
- [ ] ⚠️ Pemanggilan `INSERTJSONOFFERLIFE` **commit sendiri**. Kegagalan sesudahnya tidak boleh
      menghapus rekam offer; pemanggilan ulang **aman** karena upsert. *(lihat catatan di bawah)*
- [ ] Tidak ada nilai uang pada penawaran yang melewati `float`, termasuk di JSON API.

##### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §12

- [ ] ⚠️ Penawaran tersimpan **relasional**; `JSON_OFFER_LIFE` **tidak ditulis** dan
      `INSERTJSONOFFERLIFE` **tidak dipanggil**. AC di atas tentang "commit sendiri" karena itu
      **tidak berlaku**. *(AC 32, 33, 34 spec; penyimpangan sadar 1)*
- [ ] ⚠️ Riwayat penawaran/konfirmasi ceding tersimpan di **`T_VIEW_SUGGEST`**, anak **langsung**
      header polis — tidak ada simpul `OfferFacIn` di skema baru. Kolomnya: nomor urut, tanggal, PIC,
      hasil (`Accept`/`Reject`/`Decline`), komentar, dan tahap (`Offer`/`Bind`).
      *(AC 44 spec; `[terverifikasi]` `PremiumList Life/Activity/AddHistorySuggest.xml`,
      `ASM-FW-GISFW-WORK-LIFE!ADDHISTORYSUGGEST`)*
      *(AC 14 spec; **ADR-U-0003**)*

#### Blocker

**Tidak ada pemblokir.** Satu catatan terbuka yang **tidak** memblokir: **OQ-067** — lihat di bawah.

#### Catatan — `INSERTJSONOFFERLIFE` ada di tahap ini, bukan di rantai simpan premium list

⚠️ `[terverifikasi]` Aturan urutan procedure `[keputusan desain]` menempatkan `INSERTJSONOFFERLIFE`
**paling akhir**, sesudah `INSERTJSONPOLISLIFE`. Di korpus, ia **tidak berada di rantai itu**: satu-
satunya perujuk `SaveOfferJsonLife_SQL` di kedua modul adalah `InputOfferLife_ACT` step 4 — **satu
tahap lebih awal**. Rantai simpan premium list (`InsertJsonPolisLife_Act`) merujuk
`GetJsonProductLife`, `InsertPLSummary`, `InsertJsonPolis`, `SaveLifeinProduction_SQL`,
`GetNopolisByIDPega` — bukan `SaveOfferJsonLife_SQL`.

Tiket ini **mengikuti korpus**: rekam offer ditulis di tahap penawaran. Aturan urutan diterapkan pada
rantai simpan premium list di tiket **05a**/**05b**. Konfirmasi dicatat sebagai **OQ-067**.

#### Catatan — kode mati yang tidak dimigrasikan

`[terverifikasi]` `Call ASMForceCaseClose` (step 1) **REMARK**. Penutupan case pada `Decline`
dikerjakan mesin alur sistem baru, bukan dengan memanggil penutup paksa.

⚠️ **OQ-066**: penanda `<pyStepsBlockName>` **tidak dapat dipercaya sendirian** di modul Life —
ada langkah aktif yang sebenarnya mati. Sebelum memigrasikan langkah mana pun dari modul ini,
konfirmasikan ke work owner; jangan menyimpulkan hidup/mati dari penanda saja.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

## PremiumList Life - 02 - Periode tutup buku dibaca dari `POOLDATA.TANGGAL_CLOSING` — dan gagal terang-terangan bila kosong

**Status:** ready-for-agent

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**

#### Hasil & nilai pengguna

Sebagai **tim akuntansi**, saya ingin setiap transaksi premium list dibukukan ke **periode produksi
yang benar** menurut tanggal tutup buku yang berlaku, dan saya ingin sistem **berhenti dan berteriak**
bila tanggal tutup buku tidak dapat dibaca — supaya kesalahan periode ketahuan saat itu juga, bukan
saat tutup buku bulan berikutnya ketika angkanya sudah terlanjur salah. *(User story 39 di spec)*

#### Area codebase

`internal/repository` (pembacaan `TANGGAL_CLOSING`), `internal/services` (kebijakan periode produksi;
jam yang dapat dikendalikan), `internal/handlers` (pesan kesalahan yang menyebut tabel sumber),
`frontend/` (menampilkan periode produksi yang dipilih pada layar premium list).

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path |
| --- | --- | --- |
| `SubmitPremiumList_Act` | `ASM-FW-GISFW-WORK-LIFE` / `SUBMITPREMIUMLIST_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/SubmitPremiumList_Act.xml` (214.154 byte, tersimpan `20260122T072213`, **17 langkah**) |
| `GETTanggalClosing_SQL` | `ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETTANGGALCLOSING_SQL` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GETTanggalClosing_SQL.xml` — `SELECT * FROM POOLDATA.TANGGAL_CLOSING` |

`[terverifikasi]` Rantai di `SubmitPremiumList_Act`:

| Step | Baris | Isi |
| ---: | ---: | --- |
| 2 | 715 | `RDB-List` "Get Tanggal Closing" → `GETTanggalClosing_SQL` |
| 3 | 892 | "Set Tanggal Closing" — `Local.TglProd = TglProd.pxResults(1).TANGGAL` (918–919) |
| 3 | **966** | ⚠️ `Local.TglProd = @if(Local.TglProd=="",25,Local.TglProd)` — **fallback diam** |
| 4 | **1211** | `Local.NextMonth = @if(@toDecimal(Local.currentdate)>Local.TglProd, @toDecimal(Local.CurrentMonth)+1, @toDecimal(Local.CurrentMonth))` |
| ~~15~~ | 3387 | ~~"Kalau acc di atas tanggal 25 akan masuk produksi bulan berikutnya"~~ — **REMARK** (`<pyStepsBlockName>//`, baris 3398); memuat precondition `@toDecimal(Local.currentdate)>Local.TglProd` (3491) |

`[terverifikasi]` Pergeseran periode: `@CurrentDate("yyyy","Asia/Jakarta") + Local.NextMonth +
"01T050000.000 GMT"` — **tanggal 1 bulan berikutnya, `05:00 GMT` = 12:00 WIB**.

⚠️ `[terverifikasi]` **Dua versi hidup berdampingan.** `SubmitPremiumList_Act` (`20260122`) membaca
ambang dari tabel; `InsertJsonPolisLife_Act` (`ASM-FW-GISFW-WORK-LIFE` / `INSERTJSONPOLISLIFE_ACT`,
`20260728` — **lebih baru**) masih menanam konstanta: step 4 dengan precondition
`@toDecimal(Local.currentdate)>25` (baris 1170). `[keputusan work owner]` **ikuti yang dari DB.**

`[data DBA]` `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` juga menggulir periode lewat
`POOLDATA.TANGGAL_CLOSING` — jadi tabel ini adalah sumber tunggal aturan periode di kedua sisi.

#### ADR terkait

**ADR-U-0006** (penomoran lewat stored procedure — periode adalah masukannya), **ADR-U-0005** (flag
lingkungan), **ADR-U-0007** (jejak audit).

#### Acceptance criteria

- [ ] Tanggal tutup buku dibaca dari `POOLDATA.TANGGAL_CLOSING` **setiap kali dibutuhkan**, bukan
      dari konstanta, bukan dari cache yang tidak pernah kedaluwarsa. *(AC 7 spec)*
- [ ] Bila `POOLDATA.TANGGAL_CLOSING` **kosong atau tidak terbaca**, transaksi **ditolak** dengan
      pesan yang **menyebut tabel sumbernya**. Sistem **tidak** memakai nilai pengganti apa pun.
      *(AC 8 spec; `[keputusan work owner]`)*
- [ ] Tidak ada konstanta `25` — maupun angka ambang lain — di kode produksi. Test yang mencari
      literal ambang di lapisan services/repository **gagal** bila ada.
- [ ] Transaksi pada tanggal **setelah** tanggal tutup buku memperoleh periode **tanggal 1 bulan
      berikutnya**, pukul `05:00 GMT` (12:00 WIB). *(AC 9 spec)*
- [ ] Transaksi pada tanggal **sama dengan** tanggal tutup buku **tetap** di periode berjalan —
      perbandingannya `>`, bukan `>=`. *(baris 1211 dan 3491 keduanya memakai `>`)*
- [ ] Pergantian tahun tertangani: tutup buku Desember menggeser ke Januari tahun berikutnya, bukan
      ke "bulan 13".
- [ ] Jam dapat dikendalikan dari test — aturan periode diuji tanpa menunggu tanggal nyata.
- [ ] Periode produksi yang terpilih **terlihat pengguna** sebelum ia menyimpan, bukan hanya
      tersimpan diam-diam.

#### Blocker

**Tidak ada.**

#### Catatan — perbedaan sengaja dari Pega

| Perilaku Pega | Sistem baru |
| --- | --- |
| `@if(Local.TglProd=="",25,Local.TglProd)` — diam-diam memakai `25` | **gagal terang-terangan**, transaksi ditolak |

Alasannya: fallback diam membukukan transaksi ke periode yang salah **tanpa meninggalkan jejak**, dan
kekeliruannya baru terlihat saat tutup buku. Ini kegagalan tersembunyi, bukan ketahanan.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

## PremiumList Life - 03 - Premium List Detail dan penomoran `PL_NUMBER`

**Status:** ready-for-agent

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**, 01 (penawaran — tahap Premium hanya terbuka setelah `Confirm` → `Premium`), 02
(periode tutup buku — periode produksi adalah **masukan** procedure penomoran)

#### Hasil & nilai pengguna

Sebagai **inputor Life**, saya ingin mengisi rincian premium list untuk penawaran yang sudah
dikonfirmasi dan memperoleh **satu `PL_NUMBER` resmi** untuknya, supaya premium list itu punya
identitas tunggal yang dapat dirujuk seluruh perusahaan — termasuk oleh Claim Life di hilir.
*(User story 9–16 di spec)*

#### Area codebase

`internal/handlers` (endpoint isi detail + endpoint submit), `internal/services` (perakitan
`PremiumListDetail`; gerbang "nomor lahir sekali"), `internal/repository` (rantai `KODE_PRODUKSI` →
`PROC_GENERATE_SEQUENCE_NUMBER`), `frontend/` (grid Premium List Detail, tampilan `PL_NUMBER`).

#### Rule Pega sumber

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

#### ADR terkait

**ADR-U-0006** (penomoran **wajib** lewat `PROC_GENERATE_SEQUENCE_NUMBER`; aplikasi tidak menyusun
format sendiri), **ADR-U-0003** (uang non-float), **ADR-U-0007** (jejak audit), **ADR-U-0015** (batas
transaksi dipegang Go; commit segera setelah nomor terbentuk agar lock `FOR UPDATE` lekas lepas).

#### Acceptance criteria

- [ ] `PL_NUMBER` diperoleh dari `PROC_GENERATE_SEQUENCE_NUMBER`; aplikasi **tidak** menyusun format
      nomor sendiri, dan tidak ada pembentuk format di lapisan services. *(AC 10 spec; **ADR-U-0006**)*
- [ ] Prefix diperoleh lewat **lookup** ke `POOLDATA.KODE_PRODUKSI` (`TYPE='LIFE'`), tidak ditanam
      sebagai konstanta. *(AC 11 spec)*
- [ ] Periode yang dikirim ke procedure berasal dari tiket **02**, bukan dari `time.Now()` mentah.
- [ ] Nomor lahir **sekali** per premium list: submit kedua atas premium list yang sudah bernomor
      **tidak** menggerakkan sequence dan **tidak** mengubah nomor.
- [ ] Empat cabang per `Type` (`QR`, `QP`, `TP`, `TR`) menentukan skema nomor yang dipakai; cabang
      dipilih dari data, bukan dari urutan langkah.
- [ ] Commit terjadi **segera setelah** nomor terbentuk, sehingga lock `SELECT … FOR UPDATE` tidak
      menahan pengguna lain. (**ADR-U-0015**)
- [ ] Dua submit berurutan menghasilkan dua nomor **berbeda dan berurutan** di bawah beban paralel.
- [ ] Baris `PremiumListDetail` tersimpan utuh dengan seluruh kolom uang, dan **tidak satu pun**
      melewati `float`. *(AC 14 spec; **ADR-U-0003**)*
- [ ] `PL_NUMBER` yang sudah terbit **terlihat** pengguna dan dapat dibaca kembali lewat API.

##### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §12

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

#### Blocker

**Tidak ada pemblokir.** Satu catatan yang **tidak** memblokir tiket ini: **OQ-068** (di mana baris
detail new business mendarat di `M_LIFE_PREMIUM_DETAIL`) — memblokir **tiket 08**, bukan tiket ini.

#### Catatan — kode mati yang tidak dimigrasikan

`[terverifikasi]` **Jangan** replikasi: `Protect Gross Premium` (8.3), ketiga cabang `proratetype`
(8.5–8.7), backup detail (8.9), `SetRateLIfePremium_Act` (11), `SpreadingLife_Act` (13) — seluruhnya
`<pyStepsBlockName>//`.

⚠️ `[terverifikasi]` **Asimetri rujukan**: `SaveMasterLPBackUp` dirujuk `<RequestType>` di
`SavePremiumList_Act` baris 5931, tetapi **tidak ada berkas rule dengan nama itu di korpus** —
rujukan mati, sejalan dengan langkahnya yang sudah REMARK.

⚠️ **OQ-066**: penanda `<pyStepsBlockName>` **tidak dapat dipercaya sendirian** di modul Life.
Konfirmasikan ke work owner sebelum menyimpulkan langkah lain hidup atau mati.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

## PremiumList Life - 04 - Unggah CSV premium list detail — staging, validasi, tinjau, simpan

**Status:** ready-for-agent

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**, 03 (premium list detail — unggahan mengisi struktur yang dibentuk di sana)

#### Hasil & nilai pengguna

Sebagai **inputor Life**, saya ingin mengunggah berkas CSV berisi ratusan baris peserta dan
**melihat dulu** apa yang lolos dan apa yang ditolak beserta **nama kolomnya**, sebelum apa pun
tersimpan permanen — supaya satu baris rusak tidak mencemari premium list dan saya tidak perlu
menebak kolom mana yang salah. *(User story 17–24 di spec)*

#### Area codebase

`internal/handlers` (endpoint unggah, endpoint tinjau, endpoint simpan permanen), `internal/services`
(mesin validasi per kolom; normalisasi uang; deteksi duplikat), `internal/repository` (tabel staging
+ pembersihannya), `frontend/` (form unggah, tabel hasil validasi berlabel kolom, tombol simpan).

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `UploadCSVLifePremium_Act` | `@BASECLASS` / `UPLOADCSVLIFEPREMIUM_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/UploadCSVLifePremium_Act.xml` | impor berkas — **rule base-class, dipakai bersama** |
| `ValidasiUploadPL_act` | `ASM-FW-GISFW-WORK-LIFE` / `VALIDASIUPLOADPL_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/ValidasiUploadPL_act.xml` | **43 langkah** validasi |
| `InsertDataUploadLife` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!INSERTDATAUPLOADLIFE` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/InsertDataUploadLife.xml` | `insert into POOLDATA.M_TEMPUPLOADLIFE (…)` — **staging** |
| `DeleteTempUploadDataLife` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!DELETETEMPUPLOADDATALIFE` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/DeleteTempUploadDataLife.xml` | `… M_TEMPUPLOADLIFE where idpega = {pyWorkPage.pyID}` |
| `CekDoubleInsured` | `ASM-FW-GISFW-INT-PRODUCT_LIFE` / `ASM!CEKDOUBLEINSURED` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/CekDoubleInsured.xml` | duplikat peserta + jumlah retensi ceding |


⚠️ **Di luar cakupan tiket ini:** jalur unggah CSV **endorsement**
(`Endorsement Life/Activity/UploadCSVEDMLifePremium_Act.xml`, `@BASECLASS` /
`UPLOADCSVEDMLIFEPREMIUM_ACT`, dan `Endorsement Life/Activity/SaveCSVEDMLife.xml`,
`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SAVECSVEDMLIFE`). `[keputusan work owner]` **Endorsement Life
adalah konteks terpisah** — lihat `.scratch/endorsement-life/`.

`[terverifikasi]` **Staging nyata**: `POOLDATA.M_TEMPUPLOADLIFE`, berkunci `idpega`, dibersihkan per
case. Validasi berjalan **atas staging**, bukan atas tabel permanen.

`[terverifikasi]` **`CekDoubleInsured`** mendeteksi duplikat dengan `upper(INSURED)` + `DOB`,
mengembalikan `min(NO)` baris pertama, dan menjumlahkan `to_number(CEDING_RETENTION)` per peserta.
⚠️ Variabel `v_Count2` dan `v_PLRetensiCeding` **dideklarasikan tetapi tidak pernah diisi** — selalu
`NULL`, dinetralkan `nvl(…,0)`. Sisa pemeriksaan kedua yang dicabut; **jangan** direplikasi.

`[terverifikasi]` **33 kolom uang** divalidasi (`<pyStepsDescription>` per langkah): `NET_PREMIUM`,
`GROSS_PREMIUM`, `SHARE_NUSANTARA_RE`, `SUM_INSURED`, `CEDING_RETENTION`, `SUM_REASURED`, `CLAIM`,
`TAX`, `BROKERAGE_FEE`, `OVR_COMM`, `PROF_COMM`, `EM_PERCENT`, `COMM`, `FLEET_DISCOUNT`, seluruh
kelompok `*_REFUND`, seluruh kelompok `*_RETRO`, `*_REFUND_RETRO`, dan `CLAIM_AMOUNT`.

`[terverifikasi]` Validasi non-uang: `CERTIFICATE_NO` tidak boleh duplikat; `NAME_OF_INSURED` dan
`POLICY_HOLDER` **harus ada di `M AGENT`**; `DOB`, `BEGIN_DATE`, `EXPIRED_DATE` format `dd/mm/yyyy`;
`PLAN`, `CURRENCY` wajib; `MEDICAL_STATUS` salah satu dari `FCL` / `M` / `NM`; lampiran wajib
("BELUM ADA LAMPIRAN").

⚠️ `[terverifikasi]` **Normalisasi desimal Pega berbahaya.** Langkah "Rubah decimal dari koma jadi
titik" (baris 12480) memakai `@replaceAll(.KOLOM, ",", ".")` atas tiap kolom uang — **mengganti
setiap koma**, tanpa membedakan pemisah desimal dari pemisah ribuan. Nilai `1,234,567.89` menjadi
`1.234.567.89`, yang bukan angka. **Jangan** direplikasi apa adanya.

#### ADR terkait

**ADR-U-0003** (uang non-float, desimal presisi arbitrer — inti tiket ini), **ADR-U-0010** (penyimpanan
berkas tetap Google Storage untuk lampiran), **ADR-U-0007** (jejak audit unggahan).

#### Acceptance criteria

- [ ] Berkas yang diunggah masuk **staging** lebih dulu; kegagalan validasi **tidak** menyentuh tabel
      permanen mana pun. *(AC 18 spec)*
- [ ] Pengguna dapat **meninjau** hasil unggahan — baris lolos dan baris ditolak — sebelum menyimpan
      permanen. *(AC 19 spec)*
- [ ] Setiap penolakan menyebut **nama kolom** yang salah dan **nomor baris** CSV-nya. *(AC 17 spec)*
- [ ] Nilai uang di-parse dengan **format yang dinyatakan eksplisit** (pemisah desimal dan pemisah
      ribuan ditentukan, bukan ditebak). `1,234,567.89` dan `1.234.567,89` **tidak** boleh keduanya
      diterima diam-diam sebagai angka yang sama. Test wajib memuat kasus pemisah ribuan.
- [ ] Nilai uang **tidak** melewati `float` pada tahap mana pun — parse langsung ke desimal presisi
      arbitrer. *(AC 14–15 spec; **ADR-U-0003**)*
- [ ] Nilai uang yang lolos validasi, dibaca kembali dari staging, **identik** dengan yang diunggah —
      tidak ada pembulatan diam. *(AC 16 spec)*
- [ ] Duplikat peserta terdeteksi dengan pembandingan nama **tanpa peduli huruf besar/kecil** +
      tanggal lahir, dan pesannya menunjuk baris pertama yang bentrok.
- [ ] `CERTIFICATE_NO` ganda di dalam satu berkas ditolak.
- [ ] `NAME_OF_INSURED` dan `POLICY_HOLDER` diverifikasi terhadap master agen; baris yang tidak
      ditemukan ditolak dengan pesan yang menyebut nilainya.
- [ ] Tanggal diverifikasi berformat `dd/mm/yyyy`; `MEDICAL_STATUS` hanya menerima `FCL`, `M`, `NM`;
      `CURRENCY` wajib ada.
- [ ] Unggahan tanpa lampiran ditolak.
- [ ] Staging dibersihkan per case setelah simpan permanen **maupun** setelah pembatalan — tidak ada
      sisa baris menggantung milik case lain.
- [ ] Mengunggah berkas kedua atas case yang sama **mengganti** isi staging, tidak menumpuk.

#### Blocker

⚠️ **OQ-069 terbuka** — pesan Pega `"NET PREMIUM HARUS ADA DAN LEBIH BESAR DARI GROSS PREMIUM"`
menjanjikan aturan yang **tidak pernah diperiksa**: satu-satunya precondition atas `NET_PREMIUM` di
seluruh berkas adalah `@PropertyHasValue(.NET_PREMIUM)`, dan **tidak ada** perbandingan terhadap
`GROSS_PREMIUM`. Arah perbandingannya pun janggal — lazimnya net lebih **kecil** dari gross.

**Selama OQ-069 terbuka:** tiket ini memeriksa **keberadaan** `NET_PREMIUM` saja, persis seperti
korpus. **JANGAN** mengarang aturan perbandingan. Pesan Pega dicatat apa adanya di kode sebagai
rujukan, dengan penunjuk ke OQ-069.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

## PremiumList Life - 05a - 05a: Simpan premium summary — konversi uang di satu batas, dan urutan procedure yang mengikat

**Status:** ready-for-agent

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**, 03 (penomoran — `PL_NUMBER` adalah masukan
rekam polis)

> ⚠️ **Diselaraskan 2026-09-16 — revisi penyimpanan.** Judul asli menyebut *"urutan procedure yang
> mengikat"*; **urutan itu lenyap**. Kedua procedure yang dulu memaksa titik potong —
> `INSERTJSONPOLISLIFE` dan `INSERTJSONOFFERLIFE` — adalah **procedure JSON**, dan keduanya
> **dibuang** (spec §6, §12). Sekarang: **satu polis = satu transaksi**. Lihat blok AC "Penyimpanan
> relasional" di bawah; AC lama tentang urutan **tidak berlaku**.

#### Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin nomor premium list dan rekam summary-nya lahir **bersama atau
tidak sama sekali**, supaya tidak pernah ada nomor yang terbit tanpa rekam, maupun rekam tanpa nomor —
dan saya ingin setiap rupiah yang dikirim ke basis data kembali **persis sama** ketika dibaca.
*(User story 25–32 di spec)*

#### Area codebase

`internal/repository` (satu transaksi meliputi penomoran + summary; konversi teks ↔ desimal **hanya di
sini**), `internal/services` (perakitan 37 kolom summary), `internal/handlers` (respons memuat
`PL_NUMBER` yang terbit), `frontend/` (tampilan ringkasan premium list).

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path |
| --- | --- | --- |
| `InsertPLSummary` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/InsertPLSummary.xml` **dan** `Endorsement Life/RDBList/InsertPLSummary.xml` |
| `GetSequenceNumber_SQL` | `ASM-FW-GISFW-INT-POLICYJSON` / `RNM!GETSEQUENCENUMBER_SQL` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetSequenceNumber_SQL.xml` |
| `InsertJsonPolisLife_Act` | `ASM-FW-GISFW-WORK-LIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml` — pemanggil, step **8** |

`[terverifikasi]` **`InsertPLSummary` adalah satu rule yang sama** di kedua modul, bukan dua rule
serupa: identitas identik, `<pxUpdateDateTime>` identik (`20241101T073100.421 GMT`),
`<pyRuleSetVersion>` identik (`01-01-81`). Diff atas dua ekspor yang dinormalisasi menyisakan
**8 baris**, seluruhnya cap waktu ekspor (`pyRuleFormStatusTime`, `pyShowJavaWindowName`).
**PremiumList Life dan Endorsement Life menulis summary lewat jalur yang sama.**

`[terverifikasi]` Bentuk panggilannya:

```
BEGIN
POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY(
  {pyWorkPage.BusinessName},
  {pyWorkPage.PremiumListSummary.PL_NUMBER},
  {pyWorkPage.PremiumListSummary.PL_NUMBER_EDM},
  {TempInputData.CARI2} … {TempInputData.CARI34},
  {pyWorkPage.pzInsKey},
  {InputParam.ERRMSG out}, {InputParam.STSSAVE out});
COMMIT;
END;
```

`[data DBA]` `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY` → `INSERT` ke `M_LIFE_PREMIUM_SUMMARY`, PK dari
`M_LIFE_PREMIUM_SUMMARY_SEQ`, **37 kolom**: identitas (`ID`, `COB`, `PL_NUMBER`, `PL_NUMBER_EDM`,
`CURRENCY`, `IDPEGA`), uang gross (`PREMIUM`, `COMMISSION`, `BROKERAGE_FEE`, `OVR_COMM`, `TAX`,
`PROF_COMM`, `CLAIM`, `CLAIM_AMOUNT`, `BALANCE`, `RI_ADMIN_FEE`, `DEDUCTION`), dan turunan
`*_REFUND` / `*_RETRO` / `*_REFUND_RETRO`. **Procedure tidak commit sendiri**; `INSERT` + rollback
saat error.

⚠️ `[data DBA]` **Seluruh parameter procedure bertipe `VARCHAR2` — termasuk kolom uang.** Uang
menyeberang batas sebagai **teks**.

#### ⚠️ Urutan procedure relatif terhadap commit `[keputusan desain]`

Batas transaksi jalur Life **campuran**:

| Procedure | Commit di dalam body? |
| --- | --- |
| `PROC_GENERATE_SEQUENCE_NUMBER` | **tidak** |
| `PEGA_M_LIFE_PREMIUM_SUMMARY` | **tidak** |
| `INSERTJSONPOLISLIFE` | **ya** |
| `INSERTJSONOFFERLIFE` | **ya** |

**Urutan yang ditetapkan:**

1. **Satu transaksi Go** memuat `PROC_GENERATE_SEQUENCE_NUMBER` **dan**
   `PEGA_M_LIFE_PREMIUM_SUMMARY`, lalu **commit**.
2. Barulah `INSERTJSONPOLISLIFE` dipanggil — di luar transaksi itu (tiket **05b**).

⚠️ `[terverifikasi]` **Ini perbaikan yang disengaja, bukan tiruan.** Di Pega, keempat pembungkus
`RULE-CONNECT-SQL` menerbitkan `COMMIT;` **tepat sesudah** panggilan procedure, di dalam blok
`<pyBrowseSQL>` yang sama:

| Pembungkus | Baris `COMMIT;` |
| --- | --- |
| `GetSequenceNumber_SQL` | 88 |
| `InsertPLSummary` | 125 |
| `InsertJsonPolis` | 102 |
| `SaveOfferJsonLife_SQL` | 141 |

Akibatnya **di Pega, nomor dan summary tidak atomik**: kegagalan di antara keduanya meninggalkan
nomor yatim yang sudah ter-commit. Sistem baru menutup celah itu.

#### ADR terkait

**ADR-U-0003** (uang non-float; DDL `NUMBER` tanpa presisi → desimal presisi arbitrer — **diperkuat**
oleh temuan bahwa parameter procedure seluruhnya `VARCHAR2`), **ADR-U-0006** (penomoran),
**ADR-U-0015** (Go memegang batas transaksi), **ADR-U-0011** (bentuk rekam premium yang dikonsumsi hilir).

#### Acceptance criteria

- [ ] `PROC_GENERATE_SEQUENCE_NUMBER` dan `PEGA_M_LIFE_PREMIUM_SUMMARY` dipanggil di dalam **satu
      transaksi Go**, dan transaksi itu **commit sebelum** procedure lain dipanggil. *(AC 20 spec)*
- [ ] ~~`INSERTJSONPOLISLIFE` dan `INSERTJSONOFFERLIFE` **tidak** dipanggil dari dalam transaksi
      itu.~~ ⚠️ **TIDAK BERLAKU 2026-09-16** — keduanya **dibuang**; lihat AC pengganti di bawah.

##### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §6, §12

- [ ] ⚠️ **Satu polis ditulis dalam SATU transaksi**: header, seluruh baris rekap mata uang, seluruh
      peserta, seluruh spreading, seluruh spreading retro, dan seluruh riwayat penawaran — lalu
      **commit sekali**. Kegagalan di tingkat mana pun **membatalkan seluruhnya**. Dibuktikan dengan
      menyuntikkan kegagalan pada baris peserta ke-N dan memastikan **tidak ada** polis tersimpan.
      *(AC 35 spec; penyimpangan sadar 1)*
- [ ] Penomoran (`PROC_GENERATE_SEQUENCE_NUMBER`) berada **di dalam** transaksi itu: tidak pernah ada
      nomor tanpa polis, tidak pernah ada polis tanpa nomor. *(AC 36 spec)*
- [ ] ⚠️ **Tidak ada procedure JSON yang dipanggil.** Test yang menemukan pemanggilan
      `INSERTJSONPOLISLIFE`, `INSERTJSONOFFERLIFE`, atau padanan `@ASM.GetPageJSONString()`
      **gagal**. *(AC 32, 34 spec; penyimpangan sadar 1)*
- [ ] ⚠️ **Rekap uang per mata uang tersimpan** di tabel rekap — bukan dihitung lalu dibuang. Polis
      bermata uang ganda menghasilkan **satu baris rekap per mata uang**, masing-masing dengan nilai
      uangnya. *(AC 37 spec; penyimpangan sadar 2)*
- [ ] ⚠️ `M_LIFE_PREMIUM_SUMMARY`, `M_LIFE_PREMIUM_DETAIL`, dan `LIFEINPRODUCTION` **tidak ditulis**.
      Test yang menemukan tulisan ke ketiganya **gagal**. *(AC 33 spec)*
- [ ] ⚠️ Bila jalur warisan `SaveMasterLPDet` masih dipakai selama transisi, ia berada **di luar**
      transaksi polis dan **dapat diulang** — ia satu-satunya titik potong yang tersisa.
      *(spec §6; `[data DBA]` `COMMIT` di dalam procedure)*
      *(AC 21 spec)*
- [ ] Kegagalan **sebelum** commit summary tidak meninggalkan nomor maupun rekam separuh: test
      menyuntikkan kegagalan di `PEGA_M_LIFE_PREMIUM_SUMMARY` dan memastikan **tidak ada** baris
      `M_LIFE_PREMIUM_SUMMARY` **dan** sequence tidak bergerak. *(AC 22 spec)*
- [ ] Kegagalan **setelah** commit summary meninggalkan nomor + rekam summary **utuh** dan keadaan itu
      **terdeteksi** — bukan senyap. *(AC 23 spec)*
- [ ] **Ada test yang gagal bila urutan pemanggilan diubah** — urutannya bagian dari kebenaran, bukan
      kebetulan. *(AC 24 spec)*
- [ ] Konversi teks ↔ desimal terjadi **hanya di lapisan repository**, di satu tempat; lapisan
      services dan handlers hanya mengenal desimal. *(AC 15 spec; **ADR-U-0003**)*
- [ ] Ke-37 kolom terisi dari sumber yang benar, dan pemetaannya diuji kolom demi kolom — **bukan**
      lewat posisi `CARI2`…`CARI34` yang tidak bernama.
- [ ] Nilai uang yang dikirim dan dibaca kembali **identik**, termasuk nilai berpecahan panjang dan
      nilai negatif. Tidak ada pembulatan diam. *(AC 16 spec)*
- [ ] Baik `PL_NUMBER` maupun `PL_NUMBER_EDM` tersimpan pada rekam summary yang sama sebagai **dua
      nilai terpisah**; jalur new business mengisi yang pertama, endorsement yang kedua.
      *(AC 13 spec)*
- [ ] Keluaran galat procedure (`ERRMSG`, `STSSAVE`) **diperiksa**; galat yang dilaporkan procedure
      tidak boleh diabaikan sehingga transaksi tampak berhasil.

#### Blocker

**Tidak ada pemblokir.** Terkait tetapi tidak memblokir: **OQ-001 sisa** (DDL fisik — memblokir
tiket **09**, bukan tiket ini; kolomnya sudah terbaca dari body procedure).

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

## PremiumList Life - 05b - 05b: ~~Simpan JSON polis dan rekam produksi~~ — ⛔ **DIBATALKAN 2026-09-16** [WONTFIX]

**Status:** wontfix — **digantikan tiket 00 + 05a**

**Blocked by:** —

> ⛔ **JANGAN KERJAKAN TIKET INI.** `[keputusan work owner]` **Seluruh JSON dibuang** (spec §12).
> Tiket ini seluruhnya tentang menulis `JSON_POLIS` dan `JSON_OFFER_LIFE` **sesudah** commit, dan
> tentang menangani titik potong yang ditimbulkan keduanya. Kedua procedure itu **tidak lagi
> dipanggil**, sehingga **seluruh alasan keberadaan tiket ini hilang**.
>
> **Ke mana isinya pindah:**
>
> | Yang dulu di sini | Sekarang |
> | --- | --- |
> | Tulis `JSON_POLIS` / `JSON_OFFER_LIFE` | **dibuang** — tidak ada padanannya (AC 32, 34 spec) |
> | Urutan "sesudah commit" & titik potong | **lenyap** — satu polis = **satu transaksi** (tiket 05a, AC 35 spec) |
> | Penulisan peserta ke `M_LIFE_PREMIUM_DETAIL` | ke **`T_PREMIUM_LIST_DETAIL`**, di dalam transaksi yang sama (tiket 03 + 05a) |
> | `SaveLifeinProduction_SQL` sebagai titik potong | `LIFEINPRODUCTION` **tidak ditulis** lagi; hanya dibaca saat migrasi (tiket 00, AC 33 spec) |
> | Deteksi keadaan separuh & pengulangan | **tidak perlu** — atomisitas menggantikannya (AC 35 spec) |
> | Efek keluar Arasapas | tetap di tiket **06** |
>
> Berkas dipertahankan sebagai jejak keputusan, **bukan** sebagai pekerjaan.

---

<details>
<summary>Isi asli (sudah tidak berlaku)</summary>

**Blocked by (asli):** 05a (transaksi summary harus commit lebih dulu — itulah inti tiket ini)

#### Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin representasi JSON polis dan rekam produksi ditulis **setelah**
premium list aman tersimpan, dan saya ingin penulisan itu **dapat diulang tanpa menggandakan data**
bila gagal di tengah — supaya kegagalan jaringan atau basis data tidak pernah meninggalkan premium
list yang benar dengan polis yang tidak pernah terbit. *(User story 33–36 di spec)*

#### Area codebase

`internal/repository` (pemanggilan procedure yang commit sendiri, di luar transaksi Go; deteksi
keadaan separuh), `internal/services` (perakitan payload JSON polis; orkestrasi urutan + pengulangan),
`internal/handlers` (status "tersimpan tetapi polis tertunda" terbaca API).

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InsertJsonPolisLife_Act` | `ASM-FW-GISFW-WORK-LIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml` (327.332 byte, tersimpan `20260728T024422`, **17 langkah**) | orkestrator |
| `InsertJsonPolis` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLIS` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/InsertJsonPolis.xml` | `POOLDATA.INSERTJSONPOLISLIFE(p_IDPEGA, p_NOPOLIS, P_JSONDATA, P_TGL_PROD, P_USERNAME, {OutputData.HASIL1 out})` + `COMMIT;` (102) |
| `SaveLifeinProduction_SQL` | `ASM-FW-GISFW-WORK-LIFE` / `ASM!SAVELIFEINPRODUCTION_SQL` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/SaveLifeinProduction_SQL.xml` | `INSERT INTO POOLDATA.LIFEINPRODUCTION` + `COMMIT;` (164) |
| `GetNopolisByIDPega` | `ASM-FW-GISFW-…` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetNopolisByIDPega.xml` | baca balik untuk verifikasi |
| `GetJsonProductLife` | `ASM-FW-GISFW-…` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetJsonProductLife.xml` | data produk untuk payload |

`[terverifikasi]` Rantai rujukan `<RequestType>` di `InsertJsonPolisLife_Act`: `GetJsonProductLife`
(1444) → `InsertPLSummary` (3112, step **8**) → `InsertJsonPolis` (4094, step **10**) →
`SaveLifeinProduction_SQL` (4271, step **11**) → `GetNopolisByIDPega` (4449, step **12**).
**`SaveOfferJsonLife_SQL` tidak ada di rantai ini** — lihat tiket 01 dan **OQ-067**.

`[data DBA]` `POOLDATA.INSERTJSONPOLISLIFE` — **upsert** ke `POOLDATA.JSON_POLIS` berkunci `IDPEGA`
(`UPDATE` bila ada, `INSERT` bila belum). Kolom: `IDPEGA`, `DATA_JSON` (**CLOB**), `TGL_INPUT`,
`NOPOLIS`, `PRODKE`, `TGL_PROD`, `USERNAME`. Gagal → jatuh ke `JSON_POLIS_ERROR`.
**COMMIT di dalam procedure.**

`[terverifikasi]` Step **12** "Cek sudah masuk atau blm datanya" membaca balik lewat
`GetNopolisByIDPega`; step **13** menyalakan penanda bila `OutDataLife.pxResults(1).PL_NUMBER==""`
(baris 5105). Inilah **deteksi keadaan separuh** milik Pega — ia dipertahankan sebagai konsep, dan
menjadi pemicu alarm di tiket **06**.

#### ⚠️ Urutan yang mengikat `[keputusan desain]`

1. Transaksi Go (penomoran + summary) **sudah commit** — tiket **05a**.
2. **Baru** `INSERTJSONPOLISLIFE` dipanggil, **di luar** transaksi Go, karena ia commit sendiri.
3. `SaveLifeinProduction_SQL` juga commit sendiri (baris 164) — ia **titik potong ketiga**.
4. **Penulisan detail peserta NB ke `M_LIFE_PREMIUM_DETAIL` termasuk di sini** — dipanggil **INLINE**
   sebagai bagian alur simpan polis, lewat logika `InsertLifePremiumDetail_act`
   (`ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT` / `RULE-OBJ-ACTIVITY`) →
   `SaveMasterLPDet` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET`), yang
   **commit sendiri** (`COMMIT;` baris 252) — **titik potong keempat**. Rincian dan AC-nya di tiket
   **08**. ⚠️ `[keputusan work owner]` **Tanpa job** — di Pega, NB dipicu batch; di sistem baru
   inline, meniru EDM.

**Jangan** menempatkan data yang harus atomik pada dua sisi procedure yang commit sendiri.

#### ADR terkait

**ADR-U-0015** (batas transaksi dipegang Go; efek yang tidak boleh hilang ditangani eksplisit),
**ADR-U-0003** (uang non-float di dalam payload JSON), **ADR-U-0013** (endpoint di-resolve runtime —
berlaku pada efek keluar di tiket 06), **ADR-U-0011**.

#### Acceptance criteria

- [ ] `INSERTJSONPOLISLIFE` dipanggil **setelah** transaksi summary commit, tidak pernah di dalamnya.
      *(AC 21 spec)*
- [ ] Kegagalan pada `INSERTJSONPOLISLIFE` meninggalkan nomor + rekam summary **utuh**; API tetap
      melaporkan premium list tersimpan, dengan penanda bahwa polis **tertunda**. *(AC 23 spec)*
- [ ] Pemanggilan ulang `INSERTJSONPOLISLIFE` untuk `IDPEGA` yang sama **tidak** menggandakan baris —
      dibuktikan dengan memanggilnya dua kali dan menghitung baris `JSON_POLIS`.
- [ ] Keadaan separuh **terdeteksi**: setelah penulisan, sistem membaca balik dan menandai kasus yang
      belum lengkap. Penanda itu terbaca lewat API dan menjadi masukan tiket 06.
- [ ] `SaveLifeinProduction_SQL` diperlakukan sebagai **titik potong** tersendiri; kegagalan
      sesudahnya tidak merusak apa yang sudah ter-commit dan dapat diulang dengan aman.
- [ ] Uang di dalam payload JSON ditulis sebagai desimal presisi arbitrer — **tidak** lewat `float`
      dan tidak lewat pembulatan diam. *(AC 14 spec; **ADR-U-0003**)*
- [ ] Ada test yang **gagal** bila urutan 05a → 05b dibalik. *(AC 24 spec)*
- [ ] Kegagalan yang jatuh ke `JSON_POLIS_ERROR` **terlihat**: sistem tidak menganggapnya sukses.
- [ ] Penulisan detail peserta NB ke `M_LIFE_PREMIUM_DETAIL` terjadi **di dalam alur simpan ini**,
      bukan dijadwalkan. Setelah respons simpan berhasil, baris untuk `PL_NUMBER` itu **sudah ada**.
      Tidak ada job/cron/worker terjadwal di jalur ini. *(AC lengkap di tiket 08;
      `[keputusan work owner]`)*

#### Blocker

**Tidak ada pemblokir.** Catatan terbuka yang tidak memblokir: **OQ-067** (penempatan
`INSERTJSONOFFERLIFE` — ditulis di tahap penawaran, tiket 01).

#### Catatan — kode mati yang tidak dimigrasikan

`[terverifikasi]` Di `InsertJsonPolisLife_Act`, hanya **dua** `<pyStepsBlockName>` berisi `//`:

| Bagian | Baris | Nasib |
| --- | ---: | --- |
| step 16 `Commit` eksplisit | 5294 | **mati** — Go memegang transaksi (**ADR-U-0015**) |
| step 17 `Connect-REST` `ConvertJsonNusareToProduction` | 5383 | **mati** — tidak dipakai |

`[keputusan work owner]` **Juga tidak dimigrasikan** meski di korpus **AKTIF**: empat gerbang treaty
ID pada step 6–7 — `@contains(.ID,"1000032")` **QS** (1714), `"1000033"` **2nd QS** (1856),
`"1000034"` **SURPLUS** (1998), `"1000035"` **2nd SURPLUS** (2140). Itu logika polis lama. Jenis
treaty diambil dari data **`TypeCeding`** (`1`=QS, `2`=SURPLUS, `3`=QS+SURPLUS, `4`=XOL — ekspresi
baris ~3787). *(AC 30–31 spec)*

`[keputusan work owner]` Step **4** dengan precondition `@toDecimal(Local.currentdate)>25` (baris
1170) **tidak** direplikasi — ambang dibaca dari tabel, lihat tiket **02**.

⚠️ **OQ-066**: penanda `<pyStepsBlockName>` **tidak dapat dipercaya sendirian** di modul ini —
berkas ini justru buktinya. Konfirmasikan ke work owner sebelum menyimpulkan langkah lain.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

</details>

## PremiumList Life - 06 - Efek keluar — kiriman Arasapas, log panggilan, dan alarm kegagalan simpan

**Status:** ready-for-agent

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**, 05a (efek keluar berjalan setelah penyimpanan satu-transaksi selesai dan keadaannya diketahui)

> ⚠️ Dikoreksi 2026-09-16: sebelumnya menunjuk 05b, yang kini **wontfix** (JSON dibuang). Simpan polis
> kini satu transaksi utuh di **05a**, jadi efek keluar bergantung pada 05a.

#### Hasil & nilai pengguna

Sebagai **tim operasi**, saya ingin premium list yang sudah tersimpan diteruskan ke Arasapas dengan
alamat yang **selalu dibaca runtime**, setiap panggilan **tercatat** beserta responsnya, dan saya
diberi tahu ketika penyimpanan ternyata **gagal** — supaya kegagalan integrasi tidak pernah menjadi
temuan tutup buku. *(User story 37–38, 40 di spec)*

#### Area codebase

`internal/services` (orkestrasi efek keluar; gerbang lingkungan), `internal/repository` (resolusi
alamat dari `M_LINK_SERVICE`; tulis log panggilan), `internal/clients` (klien Arasapas + pengirim
email di balik interface), `internal/handlers` (status efek keluar terbaca API).

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `serviceInsertArasapasLife_act` | `ASM-FW-GISFW-WORK-LIFE` / `SERVICEINSERTARASAPASLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/serviceInsertArasapasLife_act.xml` (104.302 byte, **9 langkah**) | efek bisnis |
| `GetLinkService` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/GetLinkService.xml` | resolusi alamat runtime |
| `InsertLogServiceProd` (Activity) | `ASM-FW-GISFW-WORK` / `INSERTLOGSERVICEPROD` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/InsertLogServiceProd.xml` | tulis log |
| `InsertLogServiceProd` (SQL) | `ASM-FW-GISFW-WORK` / `RNM!INSERTLOGSERVICEPROD` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/InsertLogServiceProd.xml` | `INSERT INTO pooldata.MONITORING_PROD_LOG(IDPEGA, NOPOLIS, PARAMETER, JN_SERVICE, STS_MESSAGE, RESPON_MESSAGE)` |
| `SendEmailNotification` | `@BASECLASS` / `SENDEMAILNOTIFICATION` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/SendEmailNotification.xml` | **jalur alarm** |
| `GetPolicyNoByCaseId` | `ASM-FW-GISFW-…` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/GetPolicyNoByCaseId.xml` | nomor polis dari case id |
| `serviceInsertArasapasLife_act` (Endorsement) | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SERVICEINSERTARASAPASLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/serviceInsertArasapasLife_act.xml` | **rule berbeda** — class berbeda, isi berbeda. ⚠️ **konteks lain**, disebut hanya sebagai pembanding |

`[terverifikasi]` **Peta 9 langkah `serviceInsertArasapasLife_act`**:

| Step | Baris | Langkah |
| ---: | ---: | --- |
| 1 | 304 | ambil `CaseId` |
| 2 | 438 | `RDB-List` → `GetPolicyNoByCaseId` |
| 3 | 615 | set `PolicyNo` + sysdate |
| **4** | 799 | **`Call ASM-FW-GISFW-Int-M_LINK_SERVICE.GetLinkService`** — "GET LINK SERVICE" |
| **5** | 915 | **`Connect-REST`** — "Hit service Arasapas" |
| 6 | 1066 | set param insert log |
| **7** | 1305 | **`Call InsertLogServiceProd`** — "Insert Log Service" |
| 8 | 1470 | `Page-Remove` |
| ~~9~~ | 1585 | ~~`RDB-List` "Set err_note sts_konversi 9 (UPDATE JSON_POLIS)"~~ — **REMARK** (1597) |

`[terverifikasi]` **Ini bukti langsung ADR-U-0013**: alamat Arasapas **tidak** ada sebagai literal di
`Connect-REST`; ia di-resolve satu langkah sebelumnya lewat kelas `ASM-FW-GISFW-Int-M_LINK_SERVICE`.

`[terverifikasi]` **Gerbang lingkungan.** Di `InsertJsonPolisLife_Act` (PremiumList), tiga
precondition `IsPEGAPROD` (baris 5128, 5236, 5497) menjaga step **14** (email), step **15**
(Arasapas), dan step **17** (`Connect-REST`, sudah REMARK).

`[terverifikasi]` **Email adalah alarm, bukan notifikasi bisnis.** Step 14 dijaga **dua** syarat:
`OutDataLife.pxResults(1).PL_NUMBER==""` (5105) **dan** `IsPEGAPROD` — yakni ia menyala **hanya bila
pembacaan balik menunjukkan penyimpanan gagal**. `[keputusan work owner]` dikonfirmasi.

⚠️ `[terverifikasi]` **Catatan lintas konteks — jalur endorsement berbeda, dan itu BUKAN cakupan
tiket ini.** Di `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml`
(`ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT`, 444.634 byte, **16 langkah**),
hanya dua langkah yang REMARK — dan keduanya justru **jalur deteksi dan alarm**: step **13** "Cek
sudah masuk atau blm datanya" (6399) dan step **15** `SendEmailNotification` (6752). Yang tersisa
aktif hanyalah step **16** `serviceInsertArasapasLife_act`, dijaga `IsPEGAPROD` (7230).
**Endorsement hari ini berjalan tanpa alarm.**

`[keputusan work owner]` **Endorsement Life adalah konteks terpisah.** Temuan ini dicatat di
`.scratch/endorsement-life/` sebagai bahan; keputusan apakah alarm diseragamkan **diambil di konteks
itu**, bukan di sini. Tiket ini hanya mengikat jalur **new business**.

`[terverifikasi]` `ConvertJsonNusareToProduction` (`Connect-REST`, step 17) **REMARK** — tidak
dimigrasikan. *(AC 30 spec)*

#### ADR terkait

**ADR-U-0013** (alamat di-resolve runtime dari `M_LINK_SERVICE`; **dilarang** sebagai literal,
konstanta, **maupun env var** — menggantikan **ADR-U-0004**), **ADR-U-0005** (`IsPEGAPROD` → flag
lingkungan), **ADR-U-0008** (efek keluar asinkron), **ADR-U-0007** (jejak audit).

#### Acceptance criteria

- [ ] Alamat Arasapas di-resolve **runtime** dari `M_LINK_SERVICE` pada setiap panggilan. Test yang
      memindai kode untuk URL sebagai literal, konstanta, **atau pembacaan env var** **gagal** bila
      menemukannya. *(AC 26 spec; **ADR-U-0013**)*
- [ ] Di lingkungan non-production, **kedua** efek keluar tidak berjalan — sementara penyimpanan
      premium list **tetap berjalan penuh**. *(AC 25 spec; **ADR-U-0005**)*
- [ ] Gerbang lingkungan adalah **satu** flag yang dibaca di satu tempat, bukan pemeriksaan tersebar.
- [ ] Setiap panggilan keluar menulis satu baris log berisi identitas kasus, nomor polis, parameter,
      jenis layanan, status, dan **isi respons** — baik saat berhasil maupun gagal.
- [ ] Email terkirim **hanya** bila pembacaan balik menunjukkan penyimpanan **gagal** (`PL_NUMBER`
      kosong). Penyimpanan yang berhasil **tidak** mengirim email. *(AC 27 spec)*
- [ ] Kegagalan efek keluar **tidak** membatalkan premium list yang sudah tersimpan; ia tercatat dan
      terlihat, dan dapat diulang. (**ADR-U-0008**)
- [ ] Arasapas dan pengirim email berada **di balik interface** dan difake di test; yang diperiksa

##### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §12

- [ ] ⚠️ Efek keluar membaca data polis **dari tabel relasional**, bukan dari payload JSON. Test yang
      menemukan perakitan CLOB JSON untuk dikirim keluar **gagal**. *(AC 32, 34 spec; penyimpangan
      sadar 1)*
- [ ] ⚠️ Pemicu alarm **tidak lagi** berupa "deteksi keadaan separuh" antar-procedure — keadaan itu
      mustahil karena polis ditulis **atomik**. Alarm kini menyala pada **kegagalan efek keluar**
      saja. *(AC 35 spec; spec §6)*
      test adalah **efeknya** (panggilan terjadi/tidak, log tertulis, email terpicu/tidak).

#### Blocker

**Tidak ada.**

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

## PremiumList Life - 08 - Kontrak hilir — rekam premium yang dikonsumsi Claim Life

**Status:** ready-for-agent

**Blocked by:** **00 (skema tujuh tabel — PREFACTOR)**, 05a (rekam summary), 05b (alur simpan polis — penulisan detail NB menumpang di sana)

#### Hasil & nilai pengguna

Sebagai **admin klaim Life**, saya ingin menemukan premium list dan peserta yang benar untuk klaim
yang sedang saya proses — berkunci `PL_NUMBER` — **segera setelah** polis
tersimpan, tanpa menunggu proses terjadwal apa pun, supaya klaim tidak pernah tertahan hanya karena
baris pesertanya belum sempat ditulis. *(User story 39–40 di spec; **ADR-U-0001**)*

#### Area codebase

`internal/repository` (kueri pencarian premium summary + detail; penulisan detail),
`internal/services` (pemanggilan **inline** penulisan detail di alur simpan polis; kontrak yang
dipanggil konteks Claim Life), `internal/handlers` (endpoint pencarian), penandaan **kontrak lintas
konteks** di tempat bentuk rekam didefinisikan.

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InsertLifePremiumDetail_act` | `ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT` / `RULE-OBJ-ACTIVITY` | `PremiumList Life/Activity/InsertLifePremiumDetail_act.xml` (286.827 byte, `pxUpdateDateTime` `20260211T064342.645 GMT`, ruleset `01-01-91`) | **penulis detail NB** |
| `SaveMasterLPDet` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/SaveMasterLPDet.xml` **dan** `Endorsement Life/RDBList/SaveMasterLPDet.xml` | `INSERT INTO POOLDATA.M_LIFE_PREMIUM_DETAIL` + `COMMIT;` (baris 252) |
| `InsertJsonPolisLife_Act` (EDM) | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` | pemicu **inline** detail EDM (step **11.6**) |
| `GetPesertaClaim_sql1` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `RNM!GETPESERTACLAIM_SQL1` / `RULE-CONNECT-SQL` | `Claim Life/RDBList/GetPesertaClaim_sql1.xml` | **konsumen** |
| `InsertPLSummary` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL` | `PremiumList Life/RDBList/` **dan** `Endorsement Life/RDBList/` | penulis summary |

`[terverifikasi]` **`SaveMasterLPDet` adalah satu rule yang sama** untuk NB dan EDM, bukan dua rule
serupa: identitas identik, `pxUpdateDateTime` identik (`20260211T064412.717 GMT`), ukuran identik
(15.537 byte), dan diff atas dua ekspor yang dinormalisasi menyisakan **8 baris** — seluruhnya cap
waktu ekspor. **Penulisnya sudah seragam; yang berbeda hanyalah pemicunya.**

`[terverifikasi]` Kueri konsumen `GetPesertaClaim_sql1`: `SELECT * FROM POOLDATA.M_LIFE_PREMIUM_DETAIL`
dengan `WHERE PL_NUMBER = {pyWorkPage.PolicyDataLife.PremiumListSummary.PL_NUMBER}`, ditambah
pencocokan sebagian (`LIKE`) atas `CERTIFICATE_NO` dan `UPPER(NAME_OF_INSURED)`.

`[terverifikasi]` Kolom `INSERT` `SaveMasterLPDet` memuat **`PL_NUMBER` dan `PL_NUMBER_EDM`**,
`CERTIFICATE_NO`, `NAME_OF_INSURED`, `POLICY_NO`, `CURRENCY`, seluruh kolom uang gross/`*_REFUND`/
`*_RETRO`, `RATE`, `PRORATETYPE`, `SUM_AT_RISK_GROSS`, `SUM_AT_RISK_RETRO`, `RETROCEDED_SHARE`; PK
dari `M_LIFE_PREMIUM_DETAIL_SEQ.nextval`.

`[terverifikasi]` Class integrasi `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` adalah **tulang punggung
domain Life**: Claim Life 46 rule, Endorsement Life 25, PremiumList Life 19, Komite Claim Life 11.

##### Peta langkah `InsertLifePremiumDetail_act` `[terverifikasi]`

Penomoran dari `<pyStepPageReference>` — akarnya **`RH_2`**, bukan `RH_1`:

| Step | Langkah | Catatan |
| --- | --- | --- |
| 1 | `Property-Set` | set `Param.pyReportName = "SelectNoJsonPolis_RD"`, `Param.pyReportClass = "ASM-FW-GISFW-Work-LIFE"` |
| **2** | `Call Rule-Obj-Report-Definition.pxRetrieveReportData` | ⚠️ **pemicu batch** — mengambil **daftar** kasus |
| 3 | *(loop `hasil.pxResults`)* | |
| 3.1 | `Property-Set` | |
| 3.2 | `Obj-Open-By-Handle` | buka work object per baris hasil |
| 3.3 | "Insert to table detail" | |
| 3.3.1 | `Page-Remove` | |
| 3.3.2 | `Property-Set` "get ceding co name" | |
| 3.3.3 | `Property-Set` "insert nilai dari data-batch → int" | precondition `TempError.CARIDESC==1` (baris 3634) |
| **3.3.4** | `RDB-List` → **`SaveMasterLPDet`** (baris 3730) | precondition **`hasilDetail.pxResults(1).PL_NUMBER==""`** (baris 3820) — **penjaga idempotensi** |
| 3.4 | `Property-Set` "Pega to jsondata" | |
| 3.5 | `Obj-Save` | |
| **3.6** | `Commit` | **AKTIF** |

`[terverifikasi]` **Nol `<pyStepsBlockName>` di berkas ini** — tidak ada langkah ter-remark.

`[terverifikasi]` Report Definition `SelectNoJsonPolis_RD` (class `ASM-FW-GISFW-Work-LIFE`) **dirujuk
tetapi tidak ada berkasnya di korpus** — asimetri rujukan; pemilih batch itu sendiri tidak terekspor.

#### ⚠️ Penajaman kontrak hilir — **satu tabel, dua sudut pandang** `[keputusan work owner]`

Ditambahkan 2026-09-15 dari **verdict V14** grilling Endorsement Life.

Konteks Endorsement Life menulis **baris bernilai negatif** (jurnal balik) ke tabel yang **sama**
dengan yang ditulis dan dibaca di sini. Aturannya:

| Pembaca | Melihat |
| --- | --- |
| **Akuntansi / ringkasan premium** | **seluruh** baris — positif **dan** negatif; nettonya dari penjumlahan |
| **Klaim** | **hanya peserta hidup** — yang belum dibatalkan dan belum dihapus |

`[terverifikasi]` Penandanya adalah kolom **`EDMSTATUS`** pada `M_LIFE_PREMIUM_DETAIL`, terbaca dari
daftar `INSERT` di `SaveMasterLPDet`: ia diisi dari `TempValue.EDMStatus` dengan nilai
`Old` / `New` / `Delete` / `Batal`. Kolom `STATUS` **bukan** penandanya — ia berisi `0` untuk
`QR`/`QP` dan `1` untuk `TP`/`TR` (jenis transaksi).

⚠️ `[terverifikasi]` **Jalur new business — yakni tiket ini — tidak mengisi `EDMSTATUS` sama
sekali.** Sensus `InsertLifePremiumDetail_act` (`ASM-FW-GISFW-WORK-LIFE` /
`INSERTLIFEPREMIUMDETAIL_ACT` / `RULE-OBJ-ACTIVITY`): **nol** kemunculan. Baris NB karena itu masuk
dengan `EDMSTATUS` kosong/NULL — dan **harus tetap terlihat** oleh klaim.

`[terverifikasi]` Kueri klaim hari ini **belum menegakkan** kontrak ini: `GetPesertaClaim_sql1`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `RNM!GETPESERTACLAIM_SQL1` / `RULE-CONNECT-SQL`) menyaring hanya
dengan `PL_NUMBER`, `CERTIFICATE_NO`, dan `NAME_OF_INSURED` — **tanpa** penyaring status. ⚠️ Apakah
Pega menyaring di lapisan lain **tidak terbukti dari korpus**; ini **kontrak yang wajib ditegakkan
sistem baru**.

**Penegakannya milik konteks Claim Life** (spec §16, AC 25–30; tiket 02 Claim Life). Tiket ini
mengikatnya sebagai **syarat kontrak** yang harus terbukti lewat test kontrak lintas konteks.

#### ⚠️ Penyimpangan sadar — tanpa job, penulisan INLINE `[keputusan work owner]`

| | Pega existing | Sistem baru |
| --- | --- | --- |
| **Logika penulisan NB** | `InsertLifePremiumDetail_act` | **TETAP** `InsertLifePremiumDetail_act` — tidak diganti, tidak ditulis ulang |
| **Pemicu NB** | **job/batch** (step 2 `pxRetrieveReportData` + loop step 3) | **INLINE saat proses insert/simpan polis** |
| **Pemicu EDM** | inline di alur simpan (step 11.6) | inline di alur simpan — **tidak berubah** |

**Yang disamakan adalah TIMING, bukan logikanya.** Step 2 dan loop step 3 **tidak dimigrasikan** —
keduanya semata mesin batch. Yang dimigrasikan adalah **badan per-kasus** (3.1–3.6), dipanggil sekali
untuk kasus yang sedang disimpan.

**Konsekuensi positif:** data peserta langsung tersedia untuk klaim **tanpa jeda job**, dan NB
seragam dengan EDM.

#### ADR terkait

**ADR-U-0001** (batas konteks Claim Life ↔ hulu; perubahan bentuk rekam = perubahan kontrak),
**ADR-U-0011** (bentuk `PremiumListSummary` / `PremiumListDetail` yang dipakai mesin status klaim),
**ADR-U-0003** (uang non-float menyeberang batas), **ADR-U-0015** (batas transaksi dipegang Go —
`SaveMasterLPDet` commit sendiri, jadi ia titik potong).

#### Acceptance criteria

- [ ] Baris peserta **new business** ditulis ke `M_LIFE_PREMIUM_DETAIL` lewat logika
      `InsertLifePremiumDetail_act` (`ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT`),
      dipanggil **INLINE sebagai bagian alur simpan polis** — **bukan** oleh job, cron, worker
      terjadwal, atau antrean tunda. Test yang menemukan penjadwal di jalur ini **gagal**.
      `[keputusan work owner]`
- [ ] **Setelah simpan polis NB berhasil**, baris `M_LIFE_PREMIUM_DETAIL` untuk `PL_NUMBER` itu
      **sudah ada** — dibuktikan dengan membacanya **segera** sesudah respons simpan, tanpa menunggu
      apa pun.
- [ ] Baris itu **dapat dibaca jalur baca klaim**: kueri bergaya `GetPesertaClaim_sql1` — berkunci
      `PL_NUMBER`, dengan pencocokan sebagian pada `CERTIFICATE_NO` dan `NAME_OF_INSURED` **tanpa
      peduli huruf besar/kecil** — mengembalikannya.
- [ ] `SaveMasterLPDet` dipakai **apa adanya** sebagai penulis bersama — tidak dibuatkan salinan,
      tidak divariasikan per jalur. Kolom `PL_NUMBER_EDM` tetap ada di skema dan **dibiarkan kosong**
      oleh jalur new business. *(pengisiannya milik konteks Endorsement Life)*
- [ ] Penjaga idempotensi dipertahankan: penulisan hanya terjadi bila baris untuk `PL_NUMBER` itu
      belum ada (setara precondition `PL_NUMBER==""` pada step 3.3.4). Menyimpan ulang polis yang
      sama **tidak** menggandakan baris — dibuktikan dengan menyimpan dua kali lalu menghitung baris.
- [ ] `SaveMasterLPDet` diperlakukan sebagai **titik potong transaksi** (ia `COMMIT;` sendiri, baris
      252): dipanggil **setelah** transaksi penomoran + summary commit, konsisten dengan aturan
      urutan transaksi campuran. *(AC 20–21, 24 spec)*
- [ ] Kegagalan penulisan detail **tidak** membatalkan premium list yang sudah tersimpan; keadaannya
      **terdeteksi** dan pemanggilan ulang aman berkat penjaga idempotensi. *(AC 23 spec)*
- [ ] Rekam `M_LIFE_PREMIUM_SUMMARY` dan `M_LIFE_PREMIUM_DETAIL` jalur new business dapat ditemukan
      lewat `PL_NUMBER`. *(AC 28 spec)* *(pencarian lewat `PL_NUMBER_EDM` diuji di konteks Endorsement Life)*
- [ ] Nilai uang ditulis dan dibaca sebagai **desimal presisi arbitrer**; nilai yang ditulis hulu
      dibaca hilir **identik**, tanpa pembulatan di perbatasan. *(AC 16 spec; **ADR-U-0003**)*
- [ ] Bentuk kedua rekam ditandai di kode sebagai **kontrak lintas konteks**; mengubahnya memaksa
      pembaruan sadar di sisi Claim Life. *(AC 29 spec; **ADR-U-0001**)*
- [ ] ⚠️ **Jalur baca klaim menyaring peserta batal/delete.** Pembacaan peserta bergaya
      `GetPesertaClaim_sql1` **tidak menampilkan** baris ber-`EDMSTATUS` `'Batal'` atau `'Delete'`
      — hanya **peserta hidup**. `[keputusan work owner]` *(verdict V14 grilling Endorsement Life;
      AC 25 spec Claim Life)*
- [ ] Peserta **new business** tetap muncul di jalur baca klaim meski jalur NB **tidak mengisi**
      `EDMSTATUS` (kosong/NULL). ⚠️ Penyaring naif `NOT IN ('Delete','Batal')` membuang seluruh
      peserta NB di Oracle — test wajib memuat kasus ini. *(AC 26 spec Claim Life)*
- [ ] **Jalur baca akuntansi/ringkasan premium TIDAK menyaring** — ia melihat **seluruh** baris,
      positif maupun negatif, karena nettonya diperoleh dari penjumlahan. **Satu tabel, dua sudut
      pandang**, dan itu disengaja. *(AC 49 spec Endorsement Life)*
- [ ] Ada test kontrak yang menembus dari simpan polis sampai pencarian bergaya Claim Life — **satu

##### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §12

- [ ] ⚠️ Claim Life membaca peserta dari **`T_PREMIUM_LIST_DETAIL`**, bukan dari
      `M_LIFE_PREMIUM_DETAIL` maupun dari CLOB JSON. Kontrak bacanya tetap sama bentuknya.
      *(AC 33 spec; penyimpangan sadar 1)*
- [ ] ⚠️ Aturan **peserta hidup** (`EDMSTATUS` bukan `Delete`/`Batal`, NULL tetap muncul) berlaku
      **apa adanya** pada tabel baru — pindah tabel **tidak** mengubah aturannya.
- [ ] ⚠️ Kolom `PL_NUMBER` pada tabel peserta **ber-index** sejak hari pertama — ia kunci baca Claim
      Life pada tabel berjutaan baris. *(AC 49 spec)*
      seam**, API HTTP, terhadap skema uji Oracle.

#### Blocker

**Tidak ada.** **OQ-068 ditutup 2026-09-15** — `[terverifikasi]` + `[keputusan work owner]`.

Sebelumnya tiket ini `needs-info` karena korpus tidak memperlihatkan penulis `M_LIFE_PREMIUM_DETAIL`
di jalur new business — hanya `SaveMasterLPDet` di Endorsement. Work owner kemudian menambahkan
`InsertLifePremiumDetail_act` **dan** salinan `SaveMasterLPDet` ke `PremiumList Life/`, lalu
menetapkan bahwa logikanya dipakai apa adanya sementara **pemicunya** diubah dari job menjadi inline.
Jalur NB kini terbukti dan seragam dengan EDM.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
```

## PremiumList Life - 09 - Migrasi skema — `M_LIFE_PREMIUM_SUMMARY`, `M_LIFE_PREMIUM_DETAIL`, `JSON_POLIS`, `JSON_OFFER_LIFE` [WONTFIX]

**Status:** wontfix — **digantikan tiket 00**

**Blocked by:** —

> ⛔ **JANGAN KERJAKAN TIKET INI.** `[keputusan work owner]` Tiket ini merancang skema target sebagai
> **salinan** tabel existing — *"DDL tabel target menyalin tipe, presisi, PK, index, dan nullability
> dari tabel sumber"*. Keputusan 2026-09-16 membatalkan premis itu: skema target adalah **tujuh tabel
> relasional yang dirancang sendiri**, bukan salinan; dan **CLOB JSON tidak dibawa sama sekali**
> (spec §12).
>
> **Penggantinya: tiket `00` — PREFACTOR**, yang memuat DDL ketujuh tabel + `T_WORK_POLIS`, index
> pada setiap FK, dan migrasi dengan rekonsiliasi.
>
> ⚠️ **OQ-001 ditutup** — ia yang dulu menahan tiket ini. Presisi fisik kini dicocokkan DBA **di
> dalam** tiket 00, bukan sebagai prasyarat.
>
> Berkas dipertahankan sebagai jejak keputusan, **bukan** sebagai pekerjaan.

---

<details>
<summary>Isi asli (sudah tidak berlaku)</summary>

#### Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin skema Oracle sistem baru menyimpan setiap nilai premium Life
dengan **presisi yang sama persis** seperti sistem berjalan, supaya tidak satu pun rupiah berubah saat
pindah dan rekonsiliasi tidak menemukan selisih yang tidak dapat dijelaskan. *(User story 39 di spec)*

#### Area codebase

`migrations/` (DDL tabel target + index), `internal/repository` (pemetaan tipe kolom ↔ desimal
presisi arbitrer), skrip rekonsiliasi.

#### Rule Pega sumber

| Tabel | Penulis yang terbukti di korpus | Bukti |
| --- | --- | --- |
| `POOLDATA.M_LIFE_PREMIUM_SUMMARY` | `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY` lewat `InsertPLSummary` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL`) | `PremiumList Life/RDBList/InsertPLSummary.xml`, `Endorsement Life/RDBList/InsertPLSummary.xml` — **rule yang sama** |
| `POOLDATA.M_LIFE_PREMIUM_DETAIL` | `SaveMasterLPDet` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL`), PK `M_LIFE_PREMIUM_DETAIL_SEQ.nextval` — **satu rule, dua jalur** | `PremiumList Life/RDBList/SaveMasterLPDet.xml` **dan** `Endorsement Life/RDBList/SaveMasterLPDet.xml` |
| `POOLDATA.JSON_POLIS` | `POOLDATA.INSERTJSONPOLISLIFE` lewat `InsertJsonPolis`; `InsertJsonPolisEDM` (`INSERT` langsung) | `PremiumList Life/RDBList/InsertJsonPolis.xml`, `Endorsement Life/RDBList/InsertJsonPolisEDM.xml` |
| `POOLDATA.JSON_OFFER_LIFE` | `POOLDATA.INSERTJSONOFFERLIFE` lewat `SaveOfferJsonLife_SQL` (`ASM-FW-GISFW-INT-OFFERJSON` / `ASM!SAVEOFFERJSONLIFE_SQL`) | `PremiumList Life/RDBList/SaveOfferJsonLife_SQL.xml` |
| `POOLDATA.LIFEINPRODUCTION` | `SaveLifeinProduction_SQL` (`ASM-FW-GISFW-WORK-LIFE` / `ASM!SAVELIFEINPRODUCTION_SQL`) | `PremiumList Life/RDBList/SaveLifeinProduction_SQL.xml` |
| `POOLDATA.M_TEMPUPLOADLIFE` | `InsertDataUploadLife` (staging) | `PremiumList Life/RDBList/InsertDataUploadLife.xml` |
| `POOLDATA.MONITORING_PROD_LOG` | `InsertLogServiceProd` (`ASM-FW-GISFW-WORK` / `RNM!INSERTLOGSERVICEPROD`) | `PremiumList Life/RDBList/InsertLogServiceProd.xml` |
| `POOLDATA.TANGGAL_CLOSING`, `POOLDATA.KODE_PRODUKSI`, `POOLDATA.GENERATE_SEQUENCE_NUMBER` | dibaca `GETTanggalClosing_SQL`, `GetKodeProdLife_SQL`, `PROC_GENERATE_SEQUENCE_NUMBER` | `PremiumList Life/RDBList/` |

`[data DBA]` **Nama kolom sudah diketahui** dari body tiga procedure premium (diserahkan 2026-09-15):
`M_LIFE_PREMIUM_SUMMARY` **37 kolom**; `JSON_POLIS` (`IDPEGA`, `DATA_JSON` **CLOB**, `TGL_INPUT`,
`NOPOLIS`, `PRODKE`, `TGL_PROD`, `USERNAME`, dan `NOENDORS` di jalur endorsement);
`JSON_OFFER_LIFE` (kolom relasional + **delapan tanggal** + `JSONDATA` **CLOB**, PK dari
`JSON_OFFER_SEQ`).

`[terverifikasi]` Daftar kolom `M_LIFE_PREMIUM_DETAIL` terbaca lengkap dari `INSERT` di
`SaveMasterLPDet` — termasuk `PL_NUMBER`, `PL_NUMBER_EDM`, `CERTIFICATE_NO`, `NAME_OF_INSURED`,
`CURRENCY`, `RATE`, `PRORATETYPE`, `SUM_AT_RISK_GROSS`, `SUM_AT_RISK_RETRO`, `RETROCEDED_SHARE`, dan
seluruh kelompok uang gross / `*_REFUND` / `*_RETRO` / `*_REFUND_RETRO`.

#### ADR terkait

**ADR-U-0003** (uang non-float; DDL `NUMBER` **tanpa presisi** → desimal presisi arbitrer di aplikasi),
**ADR-U-0006** (tabel sequence), **ADR-U-0009** (migrasi penuh, koeksistensi ditolak), **ADR-U-0010**
(lampiran tetap di Google Storage — **bukan** bagian migrasi tabel ini).

#### Acceptance criteria

*(belum dapat difinalkan — menunggu OQ-001; disusun agar siap dijalankan begitu DDL turun)*

- [ ] DDL tabel target menyalin **tipe, presisi, PK, index, dan nullability** dari tabel sumber
      apa adanya; tidak ada kolom uang yang menjadi `FLOAT`/`BINARY_DOUBLE`. (**ADR-U-0003**)
- [ ] Nilai uang lama dibaca dan ditulis ulang **tanpa perubahan digit mana pun**; rekonsiliasi
      membandingkan nilai lama dan baru secara tepat, bukan dengan toleransi.
- [ ] Kolom `CLOB` (`DATA_JSON`, `JSONDATA`) pindah utuh, termasuk isi yang panjang.
- [ ] Sequence (`M_LIFE_PREMIUM_SUMMARY_SEQ`, `M_LIFE_PREMIUM_DETAIL_SEQ`, `JSON_OFFER_SEQ`,
      `GENERATE_SEQUENCE_NUMBER`) dipindahkan dengan **nilai berjalan yang benar**, sehingga nomor
      pasca-migrasi tidak pernah bertabrakan dengan nomor lama.
- [ ] Index yang menopang kueri hilir ada sejak hari pertama — khususnya
      `M_LIFE_PREMIUM_DETAIL(PL_NUMBER)`, yang dipakai Claim Life (tiket **08**).
- [ ] Migrasi dapat dijalankan ulang dengan aman dan punya jalur mundur yang diuji.
- [ ] Skema uji yang dipakai seluruh tiket lain dibangun **dari DDL yang sama** — bukan dari tiruan
      yang ditulis terpisah.

#### Blocker

🚧 **`needs-info` — OQ-001 (sisa) terbuka.** Pemilik: **DBA**.

**Yang sudah ada:** nama kolom keempat tabel (dari body procedure + `INSERT` `SaveMasterLPDet`).
**Yang belum ada:** **DDL fisik** — tipe, presisi, PK, index, nullability.

⚠️ Alasan ini benar-benar memblokir: **seluruh parameter `PEGA_M_LIFE_PREMIUM_SUMMARY` bertipe
`VARCHAR2`, termasuk kolom uang.** Karena itu **tipe kolom sebenarnya di tabel belum diketahui** —
apakah `NUMBER` seperti pada `OS_AKSEPTASI_KLAIM_LIFE`, atau memang `VARCHAR2`. Menebak di sini akan
menentukan presisi uang seluruh domain Life, dan salah tebak baru ketahuan saat rekonsiliasi.

**JANGAN paksa `ready`.** Pola sama dengan tiket **13** di konteks Claim — Life.

**Tidak memblokir tiket lain**: kolomnya sudah cukup untuk menulis spec dan seluruh tiket 01–07.

#### Perintah verifikasi

```
go test ./internal/...
```

</details>

# Endorsement Life

Jumlah tiket: **12**

## Endorsement Life - 01 - Pilih polis new business + lima gerbang kelayakan endorsement

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, **CL-01** (kerangka aplikasi + seam API — scaffolding lintas konteks; tidak dibuat di
sini)

#### Hasil & nilai pengguna

Sebagai **inputor Life**, saya memasukkan **nomor polis** yang ingin saya endorse dan sistem langsung
memberi tahu apakah polis itu **boleh** di-endorse — lengkap dengan alasannya bila tidak — sehingga
saya tidak pernah membuat endorsement yang mustahil diselesaikan dan tidak ada case menggantung.
*(User story 1–8 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Permintaan kelayakan endorsement; hasil kelayakan beserta daftar alasan penolakan |
| `internal/repository` | Lookup polis di `JSON_POLIS`; pencarian endorsement berjalan; baca `EdmType` polis terakhir; **satu repository terpisah** untuk pembacaan Arasapas |
| `internal/services` | Mesin lima gerbang — berjalan berurutan, mengumpulkan seluruh alasan |
| `internal/handlers` | Endpoint cek kelayakan; endpoint pencarian polis |
| `frontend/` | Layar Input EDM: isian nomor polis, penelusuran polis, tampilan alasan penolakan |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SetErrorBatalEndorsement_Act` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETERRORBATALENDORSEMENT_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/SetErrorBatalEndorsement_Act.xml` (141.275 byte) | **gerbang kelayakan** |
| `GetPL_NumberLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETPL_NUMBERLIFE` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/GetPL_NumberLife.xml` | cari polis di `JSON_POLIS` |
| `FilterProteksiEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `FILTERPROTEKSIEDMLIFE` / `RULE-OBJ-REPORT-DEFINITION` | `Endorsement Life/ReportDefinition/FilterProteksiEDMLife.xml` | cari EDM berjalan |
| `GetEdmTypeLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETEDMTYPELIFE` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/GetEdmTypeLife.xml` | baca `EdmType` dari CLOB |
| `SearcStatusBayarArasaps_SQL` | `ASM-FW-GISFW-INT-POLICYJSON` / `ASM!SEARCSTATUSBAYARARASAPS_SQL` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/SearcStatusBayarArasaps_SQL.xml` | **cek pembayaran Arasapas** |
| `BrowsePremiumList_RD` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `BROWSEPREMIUMLIST_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Endorsement Life/ReportDefinition/BrowsePremiumList_RD.xml` | **alat bantu pencarian** |
| `InputEDMLife` (Section) | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` / `RULE-OBJ-HTML-SECTION` | `Endorsement Life/Section/InputEDMLife.xml` (1.255.106 byte) | layar |

`[terverifikasi]` **Peta lima gerbang** — sub-langkah di bawah step 3 `SetErrorBatalEndorsement_Act`:

| Sub-step | Precondition | Arti |
| --- | --- | --- |
| 3.1 | `param.Nopolis==""` | **nomor polis kosong** |
| 3.2 → 3.3 | `OutData.pxResults(1).CARI1==""` (lewat `GetPL_NumberLife`, baris 895) | **polis tidak ada di `JSON_POLIS`** |
| 3.5 → 3.6 | `@LengthOfPageList(ListEdm.pxResults)>0` (RD `FilterProteksiEDMLife`, baris 1211-1212) | **sudah ada EDM belum resolve** |
| 3.7 → 3.8 | `@contains(OutData1.pxResults(1).CARI1,"3")` (lewat `GetEdmTypeLife`, baris 1766) | **sudah pernah EDM batal** |
| 3.9 → 3.10 | `@LengthOfPageList(ListPembayaran.pxResults)>0 && TempWork.EdmType=="3"` (baris 2090) | **sudah ada pembayaran** |
| ~~4~~ | `Obj-Save` | **REMARK** (`<pyStepsBlockName>//`, baris 2562) |

`[terverifikasi]` Kunci di seluruh rule endorsement adalah **`TempWork.PolicyNo`** (`NOPOLIS`) —
bukan `PL_NUMBER`. `[keputusan work owner]` `BrowsePremiumList_RD` **alat bantu pencarian, bukan
sumber kunci**.

`[terverifikasi]` Kueri pembayaran:
`select * from ARASAPAS.DETAIL_INVOICE where inv_inv_no = {InputData.CARI18} and IVD_JR_ID = '5'`.
`[keputusan work owner]` **`IVD_JR_ID = '5'` = pembayaran/pelunasan.**

#### ADR terkait

**ADR-U-0001** (batas konteks), **ADR-U-0007** (jejak audit), **ADR-U-0009** (migrasi penuh).

#### Acceptance criteria

- [ ] Kelayakan dinilai dengan **nomor polis** sebagai kunci; penelusuran premium list hanya membantu
      menemukan nomor, **tidak** menjadi kunci. *(AC 1–2 US; §3 spec)*
- [ ] Ditolak bila **nomor polis kosong**, dengan pesan yang menyebutnya. *(AC 1 spec)*
- [ ] Ditolak bila **polis tidak ditemukan** di rekam polis. *(AC 2 spec)*
- [ ] Ditolak bila polis itu **masih punya endorsement yang belum selesai** — **satu endorsement
      terbuka per polis**. Ini **aturan integritas inti** konteks ini. *(AC 3 spec)*
- [ ] Ditolak bila polis **sudah pernah dibatalkan** (`EdmType` polis terakhir = `3`). *(AC 4 spec)*
- [ ] Endorsement **Batal** ditolak bila polis **sudah dibayar**; endorsement **Perubahan Data** atas
      polis yang sudah dibayar **tetap boleh** — gerbang ini hanya berlaku bila `EdmType=3`.
      *(AC 5 spec)*
- [ ] Kelima alasan penolakan **terbaca pengguna**, dan bila lebih dari satu gerbang gagal,
      **semuanya** dilaporkan — bukan hanya yang pertama. *(AC 8 US)*
- [ ] ⚠️ **Penyimpangan sadar — Arasapas terkurung.** Pembacaan `ARASAPAS.DETAIL_INVOICE` berada di
      **satu repository** yang ditandai **batas lintas sistem**. Test yang menemukan kueri skema
      `ARASAPAS` di luar repository itu **gagal**. *(AC 47 spec; `[keputusan desain]`)*
- [ ] ⚠️ **Rule gerbang diganti nama.** Tidak ada padanan bernama "SetErrorBatal…" di kode baru — ia
      **gerbang kelayakan endorsement**, bukan pembatalan. *(`[keputusan work owner]`)*
- [ ] **Tidak ada case endorsement yang dibuat** oleh tiket ini; kelayakan murni pemeriksaan.
      *(AC 6 spec — pembuatan case milik tiket 02)*

#### Catatan pengujian

`[keputusan desain]` **Pembacaan pembayaran Arasapas TIDAK di-fake** — ia **gerbang bisnis**, bukan
sekadar integrasi. Diuji terhadap **skema uji nyata**. Yang di-fake hanya **kiriman keluar** Arasapas
dan email (tiket 10).

#### Blocker

**Tidak ada.** Menunggu scaffolding **CL-01** yang sudah `ready-for-agent` di konteks Claim Life.

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Endorsement Life - 02 - Buat case endorsement + muat & salin data polis lama

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 01 (gerbang kelayakan — case hanya dibuat setelah kelimanya lolos)

#### Hasil & nilai pengguna

Sebagai **inputor Life**, setelah polis dinyatakan layak saya menekan "buat endorsement" dan seluruh
data polis lama — termasuk ratusan barisan peserta — **otomatis tersalin** ke endorsement baru,
sehingga saya tidak mengetik ulang apa pun dan dapat langsung bekerja pada salinannya.
*(User story 9, 12 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas endorsement; rujukan ke polis lama (ID pega lama + nomor polis lama); baris detail beserta `EDMStatus` |
| `internal/repository` | Ambil `IDPEGA` polis terakhir dari `JSON_POLIS`; baca polis NB beserta detailnya |
| `internal/services` | Orkestrasi: gerbang → buat case → muat → salin → tandai `Old` |
| `internal/handlers` | Endpoint buat endorsement |
| `frontend/` | Tombol buat endorsement; grid detail hasil salinan |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `MappingEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `MAPPINGEDMLIFE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/MappingEDMLife.xml` (194.433 byte) | **mesin inti, 14 langkah** |
| `CreateCaseEMDL` | `DATA-PORTAL` / `CREATECASEEMDL` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/CreateCaseEMDL.xml` (119.680 byte) | pembuat case |
| `GetProdkeNopolis` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETPRODKENOPOLIS` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/GetProdkeNopolis.xml` | ambil `IDPEGA` polis terakhir |
| `InputEDMLife` (Flow) | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` / `RULE-OBJ-FLOW` | `Endorsement Life/Flow/InputEDMLife.xml` (68.374 byte) | alur case |

`[terverifikasi]` **Peta 14 langkah `MappingEDMLife`** — **nol `<pyStepsBlockName>` di berkas ini**:

| Step | Langkah | Catatan |
| ---: | --- | --- |
| 1 | *(gerbang)* | "Jika ada edm blm resolve exit act" — `@LengthOfPageList(ListEdm.pxResults)>0` (baris 456) |
| **3** | `Call svcAddWorkObject` | **membuat case EDM baru** |
| **6** | `RDB-List` → `GetProdkeNopolis` (baris 1172) | "Ambil idpega prodke terakhir dari tabel json_polis" |
| **7** | `Property-Set` | **"Set Old ID Pega dan Old Policy No"** |
| **8** | `Obj-Open-By-Handle` | **"Open IDPEGA NB Life"** — memuat polis new business |
| **9** | `Property-Set` | **"Mapping detail dari Life"** |
| **10** | `Property-Set` | **"Copy page dari Life ke EDM"** |
| 11.1 | `Property-Set` | **Set `"Old"` untuk detail lama** (baris 2608) |
| 11.2 | `Property-Remove` | buang baris ber-`EDMStatus == "Delete"` |
| 12 | `Property-Set` | set `.EditInput = 1` (baris 2956) — **kunci field**, lihat tiket 08 |
| 13–14 | `Obj-Save`, **`Commit`** | **AKTIF** |

`[terverifikasi]` Kunci pembacaan: `SELECT IDPEGA … WHERE NOPOLIS={TempWork.PolicyNo} AND PRODKE IS
NOT NULL ORDER BY **PRODKE DESC**`.

#### ADR terkait

**ADR-U-0001** (batas konteks — endorsement membaca polis NB sebagai hulu), **ADR-U-0003** (uang
non-float pada detail tersalin), **ADR-U-0007** (jejak audit), **ADR-U-0015** (batas transaksi).

#### Acceptance criteria

- [ ] Case endorsement dibuat **hanya setelah kelima gerbang lolos**; gerbang yang gagal **tidak**
      meninggalkan case separuh. *(AC 6 spec)*
- [ ] Data polis new business **termuat otomatis**; seluruh baris peserta tersalin ke endorsement.
      *(AC 7 spec)*
- [ ] Baris warisan polis lama bertanda **`Old`**. *(AC 7 spec)*
- [ ] **ID Pega lama dan nomor polis lama tersimpan** pada case endorsement — keduanya dipakai tiket
      03 (popup) dan tiket 04 (penomoran). *(AC 8 spec)*
- [ ] Baris polis new business **tidak pernah dihapus fisik** oleh proses apa pun. *(AC 10 spec)*
- [ ] Nilai uang tersalin **tanpa berubah** — tidak lewat `float`, tidak ada pembulatan diam.
      *(AC 21–22 spec; **ADR-U-0003**)*
- [ ] Case yang sudah dibuat **terlihat oleh pengguna lain** sejak saat itu, sehingga gerbang "satu
      endorsement terbuka per polis" (tiket 01) menolak percobaan kedua. `[keputusan work owner]`
      — commit dini di Pega **dipertahankan dengan sengaja** karena gerbang itu membutuhkannya.
- [ ] Menyalin polis dengan ratusan peserta selesai dalam satu permintaan tanpa memotong daftar.

##### Penyalinan versi ⚠️ BARU 2026-09-16 — spec §16

- [ ] ⚠️ Submit **menyalin** header, seluruh peserta, seluruh spreading, dan seluruh spreading retro
      ke **versi baru**, masing-masing dengan **identitas baru**. *(AC 59 spec; penyimpangan sadar 9)*
- [ ] ⚠️ Setiap peserta hasil salin membawa **`PARENT_ID`** yang menunjuk peserta versi sebelumnya.
      *(AC 60 spec; penyimpangan sadar 10)*
- [ ] ⚠️ **Rekap mata uang TIDAK disalin** — ia **dihitung ulang** dari peserta versi baru. Test yang
      menemukan baris rekap tersalin **gagal**. *(AC 63 spec)*
- [ ] ⚠️ **Riwayat penawaran tidak disalin** ke versi endorsement. *(AC 64 spec)*
- [ ] ⚠️ **`OldData` tidak disimpan** — data lama dibaca dari **versi sebelumnya** (`PRODKE` lebih
      kecil, lewat `PARENT_ID`). Test yang menemukan tabel/kolom penyimpan `OldData` **gagal**.
      *(AC 69 spec)*
- [ ] ⚠️ Pencocokan peserta **tidak** bergantung urutan baris maupun `CERTIFICATE_NO`: menyisipkan
      dan menghapus peserta sehingga urutan bergeser **tetap** menghasilkan pasangan yang benar.
      *(AC 62 spec; penyimpangan sadar 10 — inilah kelas bug yang diperbaiki)*
- [ ] Versi baru ber-`PRODKE` **satu lebih besar** dari versi berjalan. *(AC 56 spec)*

#### Catatan — kode mati yang tidak dimigrasikan

`[keputusan work owner]` **`CreateCaseEMDL` step 7–11 tidak direplikasi** — penyalinan data lama di
sana **mati**: `<pyStepsBlockName>//` di baris 1084 ("-- Set data dari json_polis"), 1273 (`Java`),
1382 ("-- Copy List old"), 1532 ("-- Copy Ke new"), 1671. Deskripsi langkah mati diawali `--`.
Yang hidup adalah jalur `MappingEDMLife` di atas.

⚠️ **OQ-066 berlaku.** Penanda `<pyStepsBlockName>` **tidak dapat dipercaya sendirian** di modul
Life — penetapan "mati" di atas bersandar pada **keputusan work owner**, dengan penanda sebagai
pendukung. Sebelum menyimpulkan langkah lain hidup atau mati, konfirmasikan ke work owner.

#### Blocker

**Tidak ada.**

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Endorsement Life - 03 - Popup "lihat polis lama" — satu tampilan per jenis transaksi

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 02 (case endorsement harus sudah memegang rujukan polis lama)

#### Hasil & nilai pengguna

Sebagai **inputor Life**, sebelum mengubah apa pun saya ingin **membuka polis lama dalam tampilan
tersendiri** dan melihat kolom yang relevan bagi jenis transaksinya, sehingga saya yakin sedang
meng-endorse polis yang benar. *(User story 10–11 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Baca polis NB beserta detailnya berkunci ID pega lama / nomor polis lama |
| `internal/services` | Pemilihan tampilan menurut `Type` |
| `internal/handlers` | Endpoint "lihat polis lama" |
| `frontend/` | Popup polis lama — tampilan induk + tiga varian |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `ShowLifePremiumSummary_EDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SHOWLIFEPREMIUMSUMMARY_EDM` / `RULE-OBJ-HTML-SECTION` | `Endorsement Life/Section/ShowLifePremiumSummary_EDM.xml` (2.194.694 byte) | **pemicu, memilih varian per `Type`** |
| `ViewOldPolicy_EDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `VIEWOLDPOLICY_EDM` / `RULE-OBJ-HTML-SECTION` | `Endorsement Life/Section/ViewOldPolicy_EDM.xml` | **tampilan induk — melayani `QR`** |
| `ViewOldPolicy_EDM_QP` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `VIEWOLDPOLICY_EDM_QP` | `Endorsement Life/Section/ViewOldPolicy_EDM_QP.xml` | varian `QP` |
| `ViewOldPolicy_EDM_TP` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `VIEWOLDPOLICY_EDM_TP` | `Endorsement Life/Section/ViewOldPolicy_EDM_TP.xml` | varian `TP` |
| `ViewOldPolicy_EDM_TR` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `VIEWOLDPOLICY_EDM_TR` | `Endorsement Life/Section/ViewOldPolicy_EDM_TR.xml` | varian `TR` |
| `GetOldDetail_EDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GETOLDDETAIL_EDM` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/GetOldDetail_EDM.xml` (53.517 byte) | ⚠️ pemuat data — **rule tidak dimigrasikan, perilakunya wajib ada** |

`[terverifikasi]` **Pemasangan varian per `Type`** — di `ShowLifePremiumSummary_EDM`, tiap pasang
tombol diikuti `<pyCondition>` atas `.Type`:

| Harness (baris) | `<pyCondition>` (baris) |
| --- | --- |
| `ViewOldPolicy_EDM` (65043, 65282) | **`.Type=='QR'`** (65394) |
| `ViewOldPolicy_EDM_QP` (65575, 65816) | `.Type=='QP'` (65952) |
| `ViewOldPolicy_EDM_TR` (66136, 66374) | `.Type=='TR'` (66510) |
| `ViewOldPolicy_EDM_TP` (66694, 66932) | `.Type=='TP'` (67068) |

`[terverifikasi]` Tiap tombol adalah **pasangan aksi**: `Run Activity` → `GetOldDetail_EDM`
(`pyActivityClass` = `ASM-FW-GISFW-Work-EndorsementLife`), lalu `showHarness` dengan
`pyTarget = popup`.

#### ADR terkait

**ADR-U-0001** (batas konteks — membaca polis NB adalah pembacaan hulu), **ADR-U-0003** (uang non-float
pada nilai yang ditampilkan).

#### Acceptance criteria

- [ ] Popup "lihat polis lama" **ada dan terisi** data polis new business. *(AC 9 spec)*
- [ ] Tampilan dipilih menurut `Type`: **`QR` memakai tampilan induk**; `QP`, `TP`, `TR` memakai
      variannya masing-masing. **Keempat jenis didukung.** *(AC 9 spec)*
- [ ] Membuka popup **tidak mengubah apa pun** pada endorsement maupun pada polis lama — ia murni
      baca.
- [ ] Nilai uang yang ditampilkan **persis** seperti tersimpan; tidak ada pembulatan tampilan yang
      menyesatkan. (**ADR-U-0003**)
- [ ] Popup dapat dibuka **berulang kali** tanpa efek samping.

#### Catatan — rule mati, perilaku hidup

⚠️ `[keputusan work owner]` **Rule `GetOldDetail_EDM` tidak dimigrasikan** — 4 dari 5 langkahnya
`<pyStepsBlockName>//` (baris 244, 539, 841, 975); hanya step 2 `Obj-Open-By-Handle` "Get data Life
Old" (page `WorkLife`) yang hidup.

**Tetapi perilakunya WAJIB ADA.** `[terverifikasi]` Rule itu dijalankan oleh **keempat** tombol
sebelum popup tampil — membuangnya begitu saja **mengosongkan keempat popup**. Di sistem baru,
pemuatan polis lama memakai **jalur baca biasa** berkunci ID pega lama / nomor polis lama yang sudah
dipegang case sejak tiket 02 — tanpa activity khusus.

⚠️ **OQ-066 berlaku** — penetapan "rule mati" bersandar pada **keputusan work owner**; penanda
`<pyStepsBlockName>` hanya pendukung.

#### Blocker

**Tidak ada.**

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Endorsement Life - 04 - Penomoran `PL_NUMBER_EDM` dan kenaikan `PRODKE`

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 02 (nomor dirakit dari nomor polis lama + `PRODKE` yang dibaca saat pemetaan)

#### Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin setiap endorsement memperoleh **nomornya sendiri** yang tidak
menimpa penomoran new business, lahir **sekali saja**, dan urutannya terbaca — sehingga satu polis
dengan beberapa endorsement tetap dapat ditelusuri. *(User story 23–26 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | `PL_NUMBER_EDM` dan `PRODKE` pada entitas endorsement |
| `internal/repository` | Pembacaan `PRODKE` **satu urutan**; perakitan nomor |
| `internal/services` | Gerbang "nomor lahir sekali" |
| `internal/handlers` | Nomor tampil pada respons |
| `frontend/` | Nomor endorsement terlihat setelah terbit |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `GenerateNoEDM_Life` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `GENERATENOEDM_LIFE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/GenerateNoEDM_Life.xml` | orkestrator, **7 langkah** |
| `Generate_NoEndorsmentLife` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!GENERATE_NOENDORSMENTLIFE` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/Generate_NoEndorsmentLife.xml` | **merakit nomor** |
| `GetProdKeOldData_SQL` | `ASM-FW-GISFW-INT-OFFERJSON` / `ASM!GETPRODKEOLDDATA_SQL` / `RULE-CONNECT-SQL` | *(dirujuk `GenerateNoEDM_Life` step 3)* | baca `PRODKE` — ⚠️ urutan berbeda |
| `GetProdkeNopolis` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETPRODKENOPOLIS` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/GetProdkeNopolis.xml` | baca `IDPEGA` — urutan `PRODKE DESC` |

`[terverifikasi]` **Rantai penomoran `GenerateNoEDM_Life`:**

| Step | Baris | Isi |
| ---: | ---: | --- |
| 3 | 680 | `RDB-List` "Get ProdKe dari JSON_POLIS" → `GetProdKeOldData_SQL` |
| 4 | 872 | "Set Prodke" — `Local.Prodke + 1` (946) → `InputData.CARI14` (966) |
| **5** | 1068 | `RDB-List` "Generate No Endorsement" → `Generate_NoEndorsmentLife` — precondition `PL_NUMBER_EDM==""` (1216) |
| **6** | 1264 | "Set No Endorsement" — `PL_NUMBER_EDM = Local.Nopolis` (1337) — precondition `PL_NUMBER_EDM==""` (1390) |
| 7 | 1430 | `Obj-Save` |

`[terverifikasi]` SQL perakit:
`SELECT NOPOLIS||'/'||{InputData.CARI14} AS HASIL1 FROM POOLDATA.JSON_POLIS WHERE NOPOLIS =
{pyWorkPage.PolicyNo} ORDER BY PRODKE DESC` → **`PL_NUMBER_EDM` = `<nomor polis>/<PRODKE+1>`**.

⚠️ `[terverifikasi]` **Dua pembaca `PRODKE` dengan urutan berbeda di Pega:**

| Rule | `ORDER BY` |
| --- | --- |
| `GetProdkeNopolis` | **`PRODKE DESC`** |
| `GetProdKeOldData_SQL` | `TGL_INPUT desc` |

#### ADR terkait

**ADR-U-0006** — ⚠️ **TIDAK berlaku pada nomor EDM.** `PL_NUMBER_EDM` **bukan** dari
`PROC_GENERATE_SEQUENCE_NUMBER`; jangan memaksanya lewat sequence terpusat. ADR-U-0006 tetap berlaku
untuk `PL_NUMBER` new business (**PL-03**). Juga **ADR-U-0007** (jejak audit), **ADR-U-0015**.

#### Acceptance criteria

- [ ] `PL_NUMBER_EDM` dirakit **`<nomor polis>/<PRODKE + 1>`** — **bukan** dari
      `PROC_GENERATE_SEQUENCE_NUMBER`. Test yang menemukan pemanggilan sequence terpusat di jalur ini
      **gagal**. *(AC 23 spec; **ADR-U-0006** tidak berlaku di sini)*
- [ ] Nomor **lahir sekali**: memanggil ulang pada endorsement yang sudah bernomor **tidak** mengubah
      nomornya dan **tidak** menaikkan `PRODKE`. *(AC 24 spec)*
- [ ] Dua endorsement berurutan atas polis yang sama memperoleh **`PRODKE` berurutan**.
      *(AC 25 spec)*
- [ ] `PL_NUMBER` new business **tidak berubah** oleh endorsement apa pun. *(AC 26 spec)*
- [ ] ⚠️ **Penyimpangan sadar — satu urutan `PRODKE`.** `PRODKE` dibaca dengan
      **`ORDER BY PRODKE DESC`** di **setiap** tempat. Ada test yang **gagal** bila ada jalur yang
      memakai `TGL_INPUT` atau urutan lain. *(AC 27 spec; `[keputusan desain]`)*
- [ ] Dua pembuatan endorsement serentak atas polis yang sama tidak menghasilkan `PRODKE` kembar.
- [ ] Nomor yang terbit **terlihat pengguna** dan dapat dibaca kembali lewat API.

#### Catatan — mengapa urutannya disatukan

`PRODKE` adalah **urutan produksi** yang menjadi dasar nomor endorsement; `TGL_INPUT` adalah **waktu
pencatatan**, yang dapat menyimpang karena entri susulan atau koreksi. Bila keduanya pernah tidak
sejalan, kedua pembaca Pega memberi `PRODKE` berbeda — dan **nomor endorsement ikut salah**. Sistem
baru menutup celah itu.

#### Blocker

**Tidak ada.**

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Endorsement Life - 05 - Maksud endorsement (`EdmType`) dan status per baris (`EDMStatus`)

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 02 (baris peserta harus sudah tersalin dan bertanda `Old`)

#### Hasil & nilai pengguna

Sebagai **inputor Life**, saya menyatakan **maksud** endorsement — mengubah data atau membatalkan
polis — dan pada Perubahan Data saya dapat **menambah** peserta baru serta **menandai** peserta yang
keluar, sehingga perubahan keanggotaan tercatat per baris dan dapat ditelusuri.
*(User story 13–17 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | `EdmType` pada endorsement; `EDMStatus` pada tiap baris detail |
| `internal/repository` | Baca `EdmType` polis dari atribut di dalam CLOB |
| `internal/services` | Mesin `EDMStatus`; aturan "baris baru hanya pada `EdmType=1`" |
| `internal/handlers` | Endpoint pilih maksud; endpoint tambah/tandai-hapus peserta |
| `frontend/` | Pilihan maksud endorsement; grid detail dengan penanda status per baris |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `GetEdmTypeLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `RNM!GETEDMTYPELIFE` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/GetEdmTypeLife.xml` | `SELECT A.DATA_JSON.EdmType … FROM POOLDATA.JSON_POLIS A` |
| `MappingEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `MAPPINGEDMLIFE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/MappingEDMLife.xml` | menulis `"Old"` (baris 2608); membuang baris `"Delete"` (step 11.2) |
| `SetPremi_EDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETPREMI_EDM` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/SetPremi_EDM.xml` (345.689 byte) | menulis `"Delete"` (1384) dan `"Batal"` (1506) |
| `SaveCSVEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SAVECSVEDMLIFE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/SaveCSVEDMLife.xml` | menulis `"New"` (2714) |
| `EditDetail_Section` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `EDITDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | `Endorsement Life/Section/EditDetail_Section.xml` | layar sunting detail |

##### `EdmType` — dua maksud

`[keputusan work owner]` **`1` = Perubahan Data, `3` = Batal.** Nilai `2` dan `4` **tidak dipakai** —
enum efektif `{1, 3}`.

`[terverifikasi]` Nilai `3` terbukti **dari kode**: precondition `SetPremi_EDM` step 2.1
`.EdmBatal=="True" || pyWorkPage.EdmType==3` (baris 1273).

`[terverifikasi]` `EdmType` tersimpan **di dalam CLOB `DATA_JSON`**, **bukan** sebagai kolom.

##### `EDMStatus` — empat nilai, bedanya **cakupan minus**

`[keputusan work owner]`

| Nilai | Kapan | Akibat pada nilai uang |
| --- | --- | --- |
| `Old` | warisan polis new business | tidak diubah |
| `New` | peserta ditambah — **hanya pada `EdmType=1`** | positif, baris baru |
| `Delete` | peserta dihapus dalam Perubahan Data | **diminuskan — selektif per peserta** (tiket 06) |
| `Batal` | lewat `EdmType=3` | **seluruh peserta — diminuskan menyeluruh** (tiket 06) |

⚠️ **`Delete` bukan sekadar penanda** — ia **juga** menghasilkan nilai negatif. Perhitungannya
dikerjakan tiket **06**; tiket ini menegakkan **penandaannya**.

#### ADR terkait

**ADR-U-0011** (mesin status per baris — pola sama dengan `AdjustmentList` di Claim Life),
**ADR-U-0007** (jejak audit), **ADR-U-0003** (uang non-float).

#### Acceptance criteria

- [ ] `EdmType` hanya menerima **`1`** dan **`3`**; nilai lain **ditolak**. *(AC 11 spec)*
- [ ] `EDMStatus` hanya menerima **`Old`**, **`New`**, **`Delete`**, **`Batal`**. *(AC 12 spec)*
- [ ] Baris **`New`** hanya dapat lahir pada `EdmType=1`; pada `EdmType=3` penambahan peserta
      **ditolak**. *(AC 13 spec)*
- [ ] `EdmType=3` **otomatis** menandai **seluruh** peserta polis sebagai `Batal` — tanpa penandaan
      satu per satu. *(AC 14 spec)*
- [ ] `Delete` menandai **hanya peserta yang dipilih**. *(AC 15 spec)*
- [ ] Baris `Old` yang tidak disentuh **tetap `Old`** dan nilainya tidak berubah.
- [ ] Maksud endorsement **terkunci** setelah case dibuat — percobaan mengubahnya ditolak **di sisi
      server**, bukan hanya di layar. *(AC 33 spec; penegakan layar di tiket 08)*
- [ ] `EdmType` dibaca dari atribut di dalam JSON polis, **bukan** dari kolom tersendiri; test yang
      mengasumsikan kolom **gagal**.
- [ ] Penanda status tiap baris **terlihat pengguna** di grid detail.

##### Status sebagai turunan ⚠️ BARU 2026-09-16 — spec §16

- [ ] ⚠️ Status peserta adalah **turunan** dari `PARENT_ID` + aksi: hasil salin yang diubah →
      `Change`; peserta baru (`PARENT_ID` `NULL`) → `New`; baris salin yang ditandai keluar →
      `Delete`; pembatalan polis → `Batal`. *(AC 61, 66, 67 spec; penyimpangan sadar 11)*
- [ ] Peserta baru **tidak punya pengurang** — tidak ada baris lama untuk diselisih. *(AC 61 spec)*
- [ ] ⚠️ Selisih `new − old` dihitung dengan menyandingkan baris pada **`PARENT_ID`**-nya, bukan pada
      indeks. *(AC 65 spec; penyimpangan sadar 10)*

#### Blocker

**Tidak ada.**

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Endorsement Life - 06 - Jurnal balik — `Delete` selektif dan `Batal` menyeluruh

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 05 (penandaan `EDMStatus` harus sudah berjalan)

#### Hasil & nilai pengguna

Sebagai **Finance**, saya ingin pembatalan polis dan penghapusan peserta menghasilkan **baris
bernilai negatif** — bukan baris nol dan bukan penghapusan — sehingga jejak transaksi asli tetap ada
dan nettonya dapat dihitung dengan menjumlahkan. *(User story 18–22 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Nilai uang sebagai desimal presisi arbitrer, **bertanda** |
| `internal/repository` | Penulisan baris negatif ke tabel yang sama |
| `internal/services` | **Mesin pembalikan tanda** — cakupan selektif versus menyeluruh |
| `internal/handlers` | Ringkasan endorsement menampilkan net |
| `frontend/` | Grid detail membedakan baris positif dan negatif secara kasatmata |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path |
| --- | --- | --- |
| `SetPremi_EDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SETPREMI_EDM` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/SetPremi_EDM.xml` (345.689 byte, **9 langkah, nol `<pyStepsBlockName>`**) |

`[terverifikasi]` **Peta langkah pembalikan:**

| Step | Deskripsi | Precondition | Akibat |
| --- | --- | --- | --- |
| **2.1** | **"Set 0 jika EDM Batal"** | `.EdmBatal=="True" \|\| pyWorkPage.EdmType==3` (baris 1273) | ⚠️ **mengalikan 32 kolom uang dengan `-1`** |
| 2.2 | "flag pengurangan" | `.EdmBatal=="True"` (1441) | `.EditInput = 1`, `.EDMStatus = "Delete"` (1384) |
| 2.3 | "flag batal" | `pyWorkPage.EdmType==3` (1583) / `.EdmBatal=="True"` (1636) | `.EDMStatus = "Batal"` (1506) |
| 5–7 | "Set COB", "Set property PremiumListSummary", "Insert to summary" | — | ringkasan ikut terbentuk |
| 9 | `Obj-Save` | — | |

##### ⚠️ Deskripsi Pega **berbohong** terhadap kodenya

`[terverifikasi]` Deskripsi langkah berbunyi **"Set 0 jika EDM Batal"** — tetapi kodenya **tidak
menulis nol**. Ia menulis `<kolom> = <kolom> * -1` pada **32 kolom uang**:

```
SUM_INSURED, CEDING_RETENTION, SUM_REASURED, SHARE_NUSANTARA_RE, SHARE_NUSANTARA_RE_GROSS,
SUM_AT_RISK_GROSS, SUM_AT_RISK_RETRO, RETROCEDED_SHARE, SHARE_RETRO, RATE, FACTOR,
GROSS_PREMIUM, NET_PREMIUM, DEDUCTION, CLAIM_AMOUNT, RI_ADMIN_FEE, BROKERAGE_FEE,
beserta seluruh kelompok *_REFUND, *_RETRO, dan *_REFUND_RETRO
```

**Baca kodenya, jangan namanya.** `[keputusan work owner]` pembalikan tanda memang perilaku
akuntansi yang dikehendaki.

#### ADR terkait

**ADR-U-0003** (uang non-float, desimal presisi arbitrer — **berlaku penuh pada nilai negatif**),
**ADR-U-0011** (status per baris), **ADR-U-0001** (kontrak hilir).

#### Acceptance criteria

- [ ] ⚠️ **Penyimpangan sadar — jurnal balik.** Baris `Delete` dan `Batal` menghasilkan nilai
      **negatif**, **bukan nol** dan **bukan penghapusan**. *(AC 16 spec; `[keputusan work owner]`)*
- [ ] **Seluruh 32 kolom uang** ikut dibalik tandanya; tidak ada komponen yang tertinggal positif.
      Test memeriksa **kolom demi kolom**. *(AC 17 spec)*
- [ ] Baris positif asli **tetap ada** berdampingan dengan baris negatif di tabel yang sama.
      *(AC 18 spec)*
- [ ] **Jumlah** baris positif dan negatif untuk peserta yang dibatalkan = **nol**. *(AC 19 spec)*
- [ ] Cakupan minus benar: **`Delete` menegatifkan hanya peserta terpilih**; **`Batal` menegatifkan
      seluruh peserta**. *(AC 20 spec)*
- [ ] Nilai uang **tidak** melewati `float` di lapisan mana pun maupun di JSON API — **termasuk
      nilai negatif**. *(AC 21 spec; **ADR-U-0003**)*
- [ ] Nilai yang ditulis dan dibaca kembali **identik**; tidak ada pembulatan diam pada nilai negatif
      maupun pada nilai berpecahan panjang. *(AC 22 spec)*
- [ ] Ringkasan premium endorsement mencerminkan **net** (positif + negatif), bukan hanya salah satu
      sisi.
- [ ] Baris `Old` yang tidak ditandai `Delete` **tidak ikut dibalik** pada endorsement Perubahan

##### Hapus = flag + minus ⚠️ BARU 2026-09-16 — spec §16

- [ ] ⚠️ Peserta yang dikeluarkan ditandai **`Delete`** dan nilainya menjadi **pengurang penuh** —
      **barisnya tetap tersimpan**. Test yang menemukan penghapusan fisik baris peserta **gagal**.
      *(AC 66 spec; penyimpangan sadar 11)*
- [ ] ⚠️ `EdmType` **`3`** (Batal) menjadikan **setiap** peserta pengurang penuh dan menandainya
      **`Batal`**, **tanpa menghapus satu baris pun**. *(AC 67 spec; penyimpangan sadar 11)*
- [ ] `EdmType` `1` dan `3` memakai **jalur simpan yang sama**; yang berbeda hanya nilai dan flag.
      *(AC 68 spec)*
- [ ] Kolom uang **menerima nilai negatif** tanpa kehilangan presisi — jurnal balik memang menulis
      minus. *(**ADR-U-0003**; spec §16)*
      Data.

#### Catatan — apa yang dilihat siapa

Baris negatif ini **masuk ke tabel yang sama** yang dibaca Claim Life. Aturan hilirnya: **akuntansi
melihat seluruh baris; klaim hanya peserta hidup.** Penegakan penyaringan ada di konteks Claim Life
(**CL-02**, spec Claim Life §16 AC 25–30) dan diikat sebagai kontrak di **PL-08** serta tiket **11**
konteks ini.

#### Blocker

**Tidak ada.**

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Endorsement Life - 07 - Unggah CSV endorsement — tanpa batas baris, mengganti bukan menumpuk

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 05 (baris `New` menuntut mesin `EDMStatus` sudah berjalan)

#### Hasil & nilai pengguna

Sebagai **inputor Life**, saya mengunggah CSV berisi peserta baru untuk endorsement Perubahan Data,
**tanpa dibatasi jumlah baris**, dan bila berkasnya salah saya dapat mengunggah ulang tanpa
menggandakan apa pun. *(User story 31–36 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Baris unggahan; hasil validasi per baris |
| `internal/repository` | Penyimpanan baris hasil unggahan; pembuangan baris `New` sebelumnya |
| `internal/services` | Parsing uang; validasi konsistensi; penandaan `New`; **pemrosesan bertahap** |
| `internal/handlers` | Endpoint unggah; endpoint tinjau hasil |
| `frontend/` | Form unggah; tabel hasil validasi berlabel baris dan kolom |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `UploadCSVEDMLifePremium_Act` | `@BASECLASS` / `UPLOADCSVEDMLIFEPREMIUM_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/UploadCSVEDMLifePremium_Act.xml` (119.366 byte) | impor berkas — **rule base-class** |
| `SaveCSVEDMLife` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SAVECSVEDMLIFE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/SaveCSVEDMLife.xml` (230.581 byte) | **penyimpan, 5 langkah, nol remark** |
| `ViewCSVResult_LifeEDM` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `VIEWCSVRESULT_LIFEEDM` / `RULE-OBJ-HTML-SECTION` | `Endorsement Life/Section/ViewCSVResult_LifeEDM.xml` | layar tinjau |

`[terverifikasi]` **`UploadCSVEDMLifePremium_Act` — hanya tiga langkah pertama hidup:**

| Step | Langkah | Status |
| ---: | --- | --- |
| 1–2 | `Page-Remove`, `Page-New` "deklarasi pyWorkPage" | aktif |
| **3** | **`Call pxUploadCSVResults`** "import data life premium" | aktif |
| ~~4–6~~ | `Page-New`, `Property-Set` (+5.2, 5.4), `Obj-Save` | **REMARK** (645, 791, 2120, 2400, 2525) |

`[terverifikasi]` **`SaveCSVEDMLife` — 5 langkah:**

| Step | Langkah | Precondition |
| --- | --- | --- |
| **2.1** | `Property-Remove` **"Remove EdmStatus \"New\""** | **`.EDMStatus=="New"`** (547) |
| 3 | `call ASM-FW-GISFW-Work-LIFE.Calculate1_Act` | ⚠️ lintas class — **tidak direplikasi**, lihat catatan |
| **4.1** | "Set property PremiumListDetail" | `.PLAN = …PremiumListDetail(1).PLAN` (2466); `.POLICY_HOLDER = …PremiumListDetail(1).POLICY_HOLDER` (2495) |
| 4.2 | `Page-Set-Messages` "Set message error" | `local.errmsg==""` (2649) |
| **4.3** | **"Set \"New\" untuk detail baru"** | **`pyWorkPage.EdmType==1`** (2791) |
| 5 | `Property-Set` | `.EditInput1 = 1` (2899) |

⚠️ `[keputusan work owner]` **Batas 50.000 baris dibuang.** Gerbang Pega
`@SizeOfPropertyList(TempWorkPage.ListLifePremiumDetailUpload)>50000` (`SetPremi_EDM` baris 3373)
**tidak direplikasi**.

#### ADR terkait

**ADR-U-0003** (uang non-float — parsing CSV adalah titik masuk uang), **ADR-U-0010** (penyimpanan berkas
tetap Google Storage untuk lampiran), **ADR-U-0007** (jejak audit unggahan).

#### Acceptance criteria

- [ ] ⚠️ **Penyimpangan sadar — tanpa batas baris.** Unggahan **tidak ditolak karena jumlah baris**.
      Test memuat berkas **di atas 50.000 baris** dan berhasil. *(AC 36 spec;
      `[keputusan work owner]`)*
- [ ] Unggahan sangat besar diproses **bertahap** (streaming/batch internal) tanpa memuat seluruh
      berkas sekaligus, dan **tanpa** menolak.
- [ ] **Unggah ulang mengganti, bukan menumpuk** — seluruh baris `EDMStatus == "New"` sebelumnya
      dibuang lebih dulu. Dibuktikan dengan mengunggah dua kali lalu menghitung baris.
      *(AC 35 spec)*
- [ ] Baris `New` **hanya lahir pada `EdmType=1`**; pada endorsement **Batal**, unggahan **tidak**
      menambah peserta. *(AC 13 spec)*
- [ ] Baris dengan **`PLAN`** atau **`POLICY_HOLDER`** berbeda dari **baris pertama** ditolak, dengan
      pesan yang menyebut **nomor baris dan nama kolom**. *(AC 37 spec)*
- [ ] Nilai uang di-parse dengan **format pemisah yang dinyatakan eksplisit**, langsung ke desimal
      presisi arbitrer — **tidak** lewat `float`. Test memuat kasus **pemisah ribuan**.
      *(AC 38 spec; **ADR-U-0003**)*
- [ ] Pengguna dapat **meninjau** hasil unggahan — baris lolos dan baris ditolak — sebelum menyimpan.
- [ ] Baris `Old` dan baris bertanda `Delete` **tidak tersentuh** oleh unggah ulang.

#### Catatan — `Calculate1_Act` tidak direplikasi

⚠️ `[keputusan work owner]` Meski `SaveCSVEDMLife` step 3 memanggil `Calculate1_Act`
(`ASM-FW-GISFW-WORK-LIFE` / `CALCULATE1_ACT` / `RULE-OBJ-ACTIVITY`, 302.597 byte) lintas class, dan
berkasnya **identik** di kedua folder modul (diff ternormalisasi **nol baris**), **di jalur
endorsement ia TIDAK dijalankan.** Ia milik **PremiumList Life (new business)** saja. Perhitungan
endorsement dikerjakan **`SetPremi_EDM`** (tiket 06).

⚠️ **Pelajaran, seiring OQ-066: berkas identik ≠ dipakai.** Sebagaimana `<pyStepsBlockName>` tidak
boleh dipakai sendirian untuk menyimpulkan hidup/mati, **kesamaan berkas antar folder tidak boleh
dipakai sendirian untuk menyimpulkan "mesin bersama"** — dan rujukan `Call` hanya menunjukkan
**kemungkinan** pemanggilan.

#### Blocker

**Tidak ada.**

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Endorsement Life - 08 - Alur keputusan `Confirm`/`Decline`, kunci field permanen, dan jejak audit

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 01 (gerbang — `Decline` harus membebaskan polis untuk di-endorse ulang),
02 (case + kunci field ter-set saat pemetaan)

#### Hasil & nilai pengguna

Sebagai **atasan**, saya menyetujui atau menolak endorsement sebelum tersimpan permanen; dan sebagai
**inputor**, bila saya salah memilih polis saya menolaknya lalu memulai ulang — sementara field
penentu **terkunci** agar endorsement tidak diam-diam berubah sasaran.
*(User story 27–30 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Tahap endorsement; status akhir; rekam jejak audit |
| `internal/repository` | Penulisan jejak audit; penutupan case |
| `internal/services` | Mesin alur `Confirm`/`Decline`; **penegakan kunci field di sisi server** |
| `internal/handlers` | Endpoint keputusan; penolakan perubahan field terkunci |
| `frontend/` | Tombol keputusan; field terkunci **beserta pesan penjelas**; riwayat transisi |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `InputEDMLife` (Flow) | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` / `RULE-OBJ-FLOW` | `Endorsement Life/Flow/InputEDMLife.xml` (68.374 byte) | alur |
| `IsLifeAccepted` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `ISLIFEACCEPTED` / `RULE-OBJ-DECISIONTABLE` | `Endorsement Life/DecisionTable/IsLifeAccepted.xml` (15.096 byte) | **keluaran `Confirm` / `Decline`** |
| `InputEDMLife` (Section) | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INPUTEDMLIFE` / `RULE-OBJ-HTML-SECTION` | `Endorsement Life/Section/InputEDMLife.xml` (1.255.106 byte) | **kunci field** |
| `CancelCreateCaseEDML` | `DATA-PORTAL` / `CANCELCREATECASEEDML` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/CancelCreateCaseEDML.xml` (20.134 byte) | ⚠️ **bukan** mekanisme batal case |
| `InboxEDMLife` | `ASSIGN-WORKLIST` / `INBOXEDMLIFE` / `RULE-OBJ-REPORT-DEFINITION` | `Endorsement Life/ReportDefinition/InboxEDMLife.xml` | daftar tugas |

`[terverifikasi]` Flow memuat assignment **Input EDM Detail** dan **Input EDM Summary**, keputusan
`IsLifeAccepted`, utility simpan, routing **`WorkList`**, status akhir **`Resolved-Completed`** dan
**`Resolved-Rejected`**.

⚠️ `[terverifikasi]` **`IsLifeAccepted` versi Endorsement hanya mengeluarkan `Confirm` dan
`Decline`** — **tanpa `Reject`**, berbeda dari versi PremiumList Life (`ASM-FW-GISFW-WORK-LIFE` /
`ISLIFEACCEPTED`) yang punya tiga keluaran. **Endorsement tidak punya jalur "kembali ke input".**

##### Kunci field — mekanismenya bukan rule `When`

`[terverifikasi]` Modul ini hanya punya **dua** rule `When` (`@BASECLASS` / `ISPEGAPROD`,
`@BASECLASS` / `RECORDEVENT`), keduanya **bukan** soal editabilitas. Mekanismenya **kondisi sebaris
`<pyDisabledWhen>`** di Section:

| Baris | Kondisi | Kendali atas |
| ---: | --- | --- |
| 3166 | `.EditInput==1` | **`.PolicyNo`** |
| 7537 | `.EditInput==1` | **`.EdmTypeBatal`** |
| 15763 | `.EditInput==1` | **`.EdmBatal`** |
| 8965, 10403 | `.EditInput1=1` | kelompok field setelah CSV disimpan |

`[terverifikasi]` **`.EditInput` adalah kunci satu arah** — tiga penulis, **semuanya ke `1`**, tidak
ada yang mengembalikan ke `0`: `MappingEDMLife` (2956), `SetPremi_EDM` (1338),
`SaveCSVEDMLife` (2899, `.EditInput1`).

#### ADR terkait

**ADR-U-0007** (jejak audit setiap transisi — **inti tiket ini**), **ADR-U-0002** (RBAC — peran pemutus;
penegakannya mengikuti pola Claim Life), **ADR-U-0001**.

#### Acceptance criteria

- [ ] `Confirm` melanjutkan endorsement ke penyimpanan. *(AC 28 spec)*
- [ ] `Decline` **menutup** case; case tertutup **tidak dapat** dilanjutkan maupun diputuskan ulang.
      *(AC 29 spec)*
- [ ] **Tidak ada keluaran `Reject`** di konteks ini — test yang menemukan jalur "kembali ke input"
      **gagal**. *(AC 30 spec)*
- [ ] Setelah `Decline`, polis yang sama **dapat di-endorse ulang** — gerbang "satu EDM terbuka per
      polis" (tiket 01) **tidak lagi menolak**. *(AC 31 spec)*
- [ ] Setiap transisi tahap menulis **jejak audit**: siapa, kapan, dari tahap apa ke tahap apa.
      *(AC 32 spec; **ADR-U-0007**)*
- [ ] ⚠️ **Penyimpangan sadar — kunci field permanen + pesan.** Nomor polis dan jenis
      endorsement/batal **terkunci permanen** setelah case dibuat; percobaan mengubahnya **ditolak di
      sisi server**, bukan hanya di layar. *(AC 33 spec; `[keputusan work owner + desain]`)*
- [ ] Layar **menampilkan pesan penjelas** mengapa field terkunci — **bukan** sekadar mematikannya
      seperti Pega. *(AC 34 spec)*
- [ ] Riwayat transisi **terlihat pengguna**, bukan hanya tersimpan.

#### Catatan — `CancelCreateCaseEDML` bukan mekanisme batal

⚠️ `[keputusan work owner]` **Tidak ada tombol batal terpisah.** Membatalkan endorsement =
**`Decline`**. `[terverifikasi]` `CancelCreateCaseEDML` (`DATA-PORTAL` / `CANCELCREATECASEEDML`)
hanya **satu langkah** `Property-Set` — ia membersihkan halaman dialog, **bukan** menutup case.
**Jangan** memigrasikannya sebagai pembatalan.

#### Blocker

**Tidak ada.**

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Endorsement Life - 09 - Simpan endorsement — batas transaksi campuran dan penjaga anti-dobel

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 04 (nomor endorsement), 06 (nilai baris, termasuk yang negatif)

#### Hasil & nilai pengguna

Sebagai **organisasi**, saya ingin endorsement yang disetujui tersimpan lengkap — rekam polis,
produksi, detail peserta, dan ringkasan premium — dan bila gagal di tengah, keadaannya
**terdeteksi dan dapat diulang tanpa menggandakan apa pun**. *(User story 43 di spec; §11 spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Endorsement + baris detail versi baru (`T_PREMIUM_LIST_DETAIL`) beserta `PARENT_ID`; rekap mata uang |
| `internal/repository` | Pemanggilan keempat penulis; **penjaga anti-dobel `(NOPOLIS, PRODKE)`** |
| `internal/services` | Orkestrasi urutan; deteksi keadaan separuh; pengulangan aman |
| `internal/handlers` | Status "tersimpan" versus "tertunda" terbaca API |
| `frontend/` | Penanda endorsement tersimpan / tertunda |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Commit sendiri? |
| --- | --- | --- | --- |
| `InsertJsonPolisLife_Act` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` (444.634 byte, **16 langkah**) | orkestrator |
| `InsertJsonPolisEDM` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLISEDM` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/InsertJsonPolisEDM.xml` | **ya** — `COMMIT;` baris 107 |
| `SaveLifeinProduction_SQL` | `ASM-FW-GISFW-WORK-LIFE` / `ASM!SAVELIFEINPRODUCTION_SQL` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/SaveLifeinProduction_SQL.xml` | **ya** — `COMMIT;` baris 164 |
| `SaveMasterLPDet` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/SaveMasterLPDet.xml` | **ya** — `COMMIT;` baris 252 |
| `InsertPLSummary` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/InsertPLSummary.xml` | **ya** — `COMMIT;` baris 125 |

`[terverifikasi]` **Rantai rujukan `<RequestType>` di `InsertJsonPolisLife_Act`:**
`GetProdKeOldData_SQL` (911) → `Generate_NoEndorsmentLife` (1304) → **`InsertJsonPolisEDM` (2543)**
→ **`SaveLifeinProduction_SQL` (2720)** → **`SaveMasterLPDet` (5154, step 11.6)** →
**`InsertPLSummary` (6226, step 12.2)** → `GetNopolisByIDPega` (6444).

`[terverifikasi]` Dua langkah **REMARK**: step **13** "Cek sudah masuk atau blm datanya" (6399) dan
step **15** `SendEmailNotification` (6752) — ⚠️ **keduanya justru DIHIDUPKAN** di sistem baru, lihat
tiket **10**.

`[terverifikasi]` Ditambah **dua `Commit` aktif** yang berjalan **sebelum** pengguna menekan simpan:
`CreateCaseEMDL` step 4 dan `MappingEDMLife` step 14 (tiket 02).

⚠️ `[terverifikasi]` **`InsertJsonPolisEDM` adalah `INSERT` polos, bukan upsert:**
`INSERT INTO POOLDATA.JSON_POLIS (IDPEGA, DATA_JSON, TGL_INPUT, NOPOLIS, NOENDORS, PRODKE, TGL_PROD,
USERNAME) VALUES (…)`. Jalur new business tidak punya masalah ini karena memakai procedure **upsert**
`INSERTJSONPOLISLIFE` berkunci `IDPEGA` (**PL-05b**).

#### ADR terkait

**ADR-U-0015** (Go memegang batas transaksi; efek yang tidak boleh hilang ditangani eksplisit),
**ADR-U-0003** (uang non-float di kolom tabel relasional maupun JSON API), **ADR-U-0011**.

#### Acceptance criteria

- [ ] Setiap penulis yang **commit sendiri** diperlakukan sebagai **titik potong**; data yang harus
      atomik **tidak dipisahkan** olehnya. *(AC 39 spec; **ADR-U-0015**)*
- [ ] Kegagalan di tengah meninggalkan keadaan yang **terdeteksi dan dapat dipulihkan** — bukan
      senyap. *(AC 40 spec)*
- [ ] ⚠️ **Penyimpangan sadar — penjaga anti-dobel.** Menjalankan penyalinan dua kali untuk
      **`(NOPOLIS, PRODKE)`** yang sama **tidak** menggandakan versi — dibuktikan dengan
      memanggilnya dua kali lalu **menghitung baris versi di `T_PREMIUM_LIST`** untuk kombinasi itu.
      ⚠️ **Bukan** dengan menghitung baris `JSON_POLIS`: JSON **tidak lagi ditulis** (§16).
      *(AC 41, 71 spec; `[keputusan work owner]`)*
- [ ] Urutan pemanggilan terdokumentasi di kode **sebagai bagian kebenaran**, bukan kebetulan; ada
      test yang **gagal bila urutannya diubah**. *(AC 42 spec)*
- [ ] Commit dini di pembuatan case dan pemetaan (tiket 02) **dipertahankan** — gerbang "satu EDM
      terbuka per polis" membutuhkannya. `[keputusan work owner]`
- [ ] Baris detail **termasuk yang bernilai negatif** tertulis utuh ke tabel peserta
      (`T_PREMIUM_LIST_DETAIL`); rekap mata uang **dihitung ulang** dari peserta versi baru.
      ⚠️ Nama tabel lama (`M_LIFE_PREMIUM_DETAIL` / `_SUMMARY`) adalah **sistem lama** — target
      tulis sistem baru adalah **tujuh tabel PremiumList Life** (§16). *(AC 18, 63 spec)*
- [ ] Nilai uang di **kolom tabel** maupun di **kontrak API** ditulis sebagai **desimal presisi
      arbitrer** — tidak lewat `float`, tidak ada pembulatan diam. (**ADR-U-0003**)
- [ ] `PL_NUMBER_EDM` tertulis pada rekam summary **di samping** `PL_NUMBER`, sebagai **dua nilai
      terpisah**.
- [ ] Kegagalan setelah salah satu commit meninggalkan bagian yang sudah ter-commit **utuh**, dan

##### Satu transaksi ⚠️ BARU 2026-09-16 — spec §11, §16

- [ ] ⚠️ Satu endorsement ditulis dalam **satu transaksi** — penyalinan versi, kolom EDM pada header,
      seluruh peserta beserta `PARENT_ID`, spreading & retro, dan rekap mata uang yang dihitung
      ulang; kegagalan di mana pun **membatalkan seluruhnya**. *(AC 70 spec)*
- [ ] ⚠️ **Tidak ada procedure JSON yang dipanggil.** Test yang menemukan pemanggilan
      `InsertJsonPolisEDM` atau padanan perakit payload **gagal**. *(AC 54 spec; penyimpangan
      sadar 8)*
- [ ] ⚠️ **Penjaga anti-dobel tetap berlaku, dengan alasan baru**: menjalankan penyalinan dua kali
      untuk maksud endorsement yang sama menghasilkan **dua versi**. Pemeriksaan `(NOPOLIS, PRODKE)`
      sebelum menulis menahannya. *(AC 71 spec; spec §11)*
- [ ] Dua `Commit` dini (`CreateCaseEMDL`, `MappingEDMLife`) tetap **di luar** transaksi simpan —
      gerbang "satu endorsement terbuka per polis" membutuhkannya. *(spec §11)*
      pengulangan **aman**.

#### Catatan — kode mati yang tidak dimigrasikan

`[terverifikasi]` `Commit` eksplisit dan `Connect-REST` yang ter-remark di rantai simpan **tidak
direplikasi** — **Go memegang transaksi** (**ADR-U-0015**).

⚠️ **OQ-066 berlaku.** Berkas ini justru buktinya: dua langkah ber-`//` (13 dan 15) **bukan** kode
mati yang dibuang, melainkan jalur yang **sengaja dihidupkan kembali** (tiket 10). **Jangan**
menyimpulkan hidup/mati dari `<pyStepsBlockName>` saja.

#### Blocker

**Tidak ada.**

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Endorsement Life - 10 - Efek keluar Arasapas + alarm yang dihidupkan

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 09 (efek keluar berjalan setelah penyimpanan selesai dan keadaannya diketahui)

#### Hasil & nilai pengguna

Sebagai **tim operasi**, saya ingin endorsement yang tersimpan diteruskan ke Arasapas dengan alamat
yang **selalu dibaca runtime**, dan saya ingin **diberi tahu bila penyimpanan gagal** — karena hari
ini jalur endorsement berjalan **tanpa alarm sama sekali**. *(User story 37–40 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/services` | Orkestrasi efek keluar; gerbang lingkungan; deteksi keadaan separuh |
| `internal/repository` | Resolusi alamat dari `M_LINK_SERVICE`; pembacaan balik untuk deteksi |
| `internal/clients` | Klien Arasapas + pengirim email **di balik interface** |
| `internal/handlers` | Status efek keluar terbaca API |
| `frontend/` | Penanda "terkirim" / "tertunda" pada endorsement |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `serviceInsertArasapasLife_act` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `SERVICEINSERTARASAPASLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/serviceInsertArasapasLife_act.xml` | efek bisnis — ⚠️ **rule berbeda** dari yang ber-class `ASM-FW-GISFW-WORK-LIFE` |
| `GetLinkService` | `ASM-FW-GISFW-INT-M_LINK_SERVICE` / `GETLINKSERVICE` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/GetLinkService.xml` | **resolusi alamat runtime** |
| `InsertJsonPolisLife_Act` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` | pemicu — step 13, 15, 16 |
| `IsPEGAPROD` | `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` | `Endorsement Life/When/IsPEGAPROD.xml` | **flag lingkungan** |

`[terverifikasi]` Rantai di `InsertJsonPolisLife_Act`:

| Step | Isi | Penanda korpus | Nasib di sistem baru |
| ---: | --- | --- | --- |
| ~~13~~ | `RDB-List` "Cek sudah masuk atau blm datanya" | **REMARK** (6399) | ⚠️ **DIHIDUPKAN** |
| ~~15~~ | `call @baseclass.SendEmailNotification` | **REMARK** (6752) | ⚠️ **DIHIDUPKAN** |
| **16** | `call serviceInsertArasapasLife_act` | aktif, dijaga `IsPEGAPROD` (7230) | dipertahankan |

⚠️ `[terverifikasi]` **Endorsement hari ini berjalan tanpa alarm** — deteksi keadaan separuh dan
email kegagalan **keduanya mati**, sementara jalur new business memilikinya (**PL-06**).

#### ADR terkait

**ADR-U-0013** (alamat di-resolve **runtime** dari `M_LINK_SERVICE`; **dilarang** sebagai literal,
konstanta, **maupun env var** — menggantikan **ADR-U-0004**), **ADR-U-0005** (`IsPEGAPROD` → flag
lingkungan), **ADR-U-0008** (efek keluar asinkron dengan antre ulang), **ADR-U-0007** (jejak audit).

#### Acceptance criteria

- [ ] Di lingkungan **non-production**, efek keluar **tidak berjalan**; penyimpanan endorsement
      **tetap** berjalan penuh. *(AC 43 spec; **ADR-U-0005**)*
- [ ] Gerbang lingkungan adalah **satu** flag yang dibaca di **satu tempat**, bukan pemeriksaan
      tersebar.
- [ ] Alamat Arasapas di-resolve **runtime** dari `M_LINK_SERVICE` pada setiap panggilan. Test yang
      memindai kode untuk URL sebagai **literal, konstanta, atau pembacaan env var** **gagal** bila
      menemukannya. *(AC 44 spec; **ADR-U-0013**)*
- [ ] ⚠️ **Penyimpangan sadar — alarm dihidupkan.** Deteksi keadaan separuh **dan** email kegagalan
      **berjalan** di jalur endorsement, dengan perilaku **sama** seperti jalur new business.
      *(AC 45 spec; `[keputusan work owner]`)*
- [ ] Email terkirim **hanya** bila pembacaan balik menunjukkan penyimpanan **gagal** — ia **alarm**,
      bukan notifikasi bisnis. Penyimpanan yang berhasil **tidak** mengirim email.
- [ ] Kegagalan efek keluar **tidak** membatalkan endorsement yang sudah tersimpan; ia **tercatat,
      terlihat, dan dapat diulang**. *(AC 46 spec; **ADR-U-0008**)*
- [ ] Arasapas **kiriman keluar** dan pengirim email berada **di balik interface** dan **di-fake** di
      test; yang diperiksa adalah **efeknya** (panggilan terjadi/tidak, email terpicu/tidak).
- [ ] ⚠️ Pembacaan **pembayaran** Arasapas (tiket 01) **tidak** ikut di-fake — ia gerbang bisnis,
      diuji terhadap skema uji nyata. Keduanya **tidak** boleh tercampur di satu klien.

#### Blocker

**Tidak ada.**

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Endorsement Life - 11 - Kontrak hilir ke Claim Life — satu tabel, dua sudut pandang

**Status:** ready-for-agent

**Blocked by:** **00 (kolom EDM + PARENT_ID — PREFACTOR)**, 09 (baris — termasuk yang negatif — harus sudah tertulis)

#### Hasil & nilai pengguna

Sebagai **Finance**, saya ingin seluruh baris endorsement terlihat di laporan saya agar nettonya
benar; dan sebagai **admin klaim Life**, saya ingin peserta yang sudah dibatalkan atau dihapus
**tidak muncul** saat saya mencari peserta untuk klaim — sehingga saya tidak pernah memproses klaim
atas peserta yang tidak lagi ditanggung. *(User story 41–43 di spec)*

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/repository` | Penulisan baris ke tabel bersama; **kueri kontrak** bergaya baca-klaim |
| `internal/services` | Kontrak yang dipanggil konteks Claim Life |
| `internal/handlers` | Endpoint pencarian rekam premium endorsement |
| — | Penandaan **kontrak lintas konteks** di tempat bentuk rekam didefinisikan |

#### Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterLPDet` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL` | `Endorsement Life/RDBList/SaveMasterLPDet.xml` **dan** `PremiumList Life/RDBList/SaveMasterLPDet.xml` | **penulis bersama** — satu rule identik |
| `InsertPLSummary` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL` | kedua modul | penulis summary |
| `GetPesertaClaim_sql1` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `RNM!GETPESERTACLAIM_SQL1` / `RULE-CONNECT-SQL` | `Claim Life/RDBList/GetPesertaClaim_sql1.xml` | **konsumen** |
| `LoadDataPesertaSpesifik_Act` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `LOADDATAPESERTASPESIFIK_ACT` / `RULE-OBJ-ACTIVITY` | `Claim Life/Activity/LoadDataPesertaSpesifik_Act.xml` (67.369 byte) | pemanggil konsumen (baris 545) |

##### Kolom penanda — **terbaca dari korpus** `[terverifikasi]`

Daftar kolom `INSERT` + klausa `VALUES` di `SaveMasterLPDet`:

| Kolom | Diisi dari | Nilai | Penanda hidup/mati? |
| --- | --- | --- | --- |
| **`EDMSTATUS`** | `TempValue.EDMStatus` | `Old` / `New` / `Delete` / `Batal` | ✅ **ya** |
| `STATUS` | `CARI48` | `0` untuk `QR`/`QP`, `1` untuk `TP`/`TR` | ❌ penanda **jenis transaksi** |
| `STATUSOLD` | `CARI47` | `1`/`0` | ❌ |

⚠️ `[terverifikasi]` **Jalur new business tidak mengisi `EDMSTATUS` sama sekali** — sensus
`PremiumList Life/Activity/InsertLifePremiumDetail_act.xml` (`ASM-FW-GISFW-WORK-LIFE` /
`INSERTLIFEPREMIUMDETAIL_ACT`): **nol** kemunculan. Baris NB masuk dengan `EDMSTATUS` kosong/NULL.

**Peserta hidup** = `EDMSTATUS` **kosong/NULL**, `Old`, atau `New`. **Mati** = `Delete` atau `Batal`.

#### ADR terkait

**ADR-U-0001** (batas konteks — perubahan bentuk rekam = perubahan kontrak), **ADR-U-0011** (bentuk
`PremiumListSummary` / `PremiumListDetail` yang dipakai mesin status klaim), **ADR-U-0003**.

#### Acceptance criteria

- [ ] Baris endorsement — **termasuk yang bernilai negatif** — tertulis ke `M_LIFE_PREMIUM_DETAIL`
      dan `M_LIFE_PREMIUM_SUMMARY`, **tabel yang sama** dengan jalur new business. *(AC 18 spec)*
- [ ] Peserta ber-`EDMSTATUS` **`Batal`** atau **`Delete`** **tidak muncul** pada jalur baca klaim.
      *(AC 48 spec; `[keputusan work owner]`)*
- [ ] Peserta **new business** (`EDMSTATUS` kosong/NULL) **tetap muncul**. ⚠️ Penyaring naif
      `EDMSTATUS NOT IN ('Delete','Batal')` membuang seluruh peserta NB di Oracle — test wajib memuat
      kasus ini dan **harus gagal** bila penyaringnya naif. *(AC 48a spec)*
- [ ] `STATUS` dan `STATUSOLD` **tidak** dipakai sebagai penanda hidup/mati. *(AC 48b spec)*
- [ ] Jalur baca **akuntansi/ringkasan premium TIDAK menyaring** — ia melihat **seluruh** baris,
      positif maupun negatif. **Satu tabel, dua sudut pandang**, dan itu disengaja. *(AC 49 spec)*
- [ ] Rekam premium hasil endorsement dapat ditemukan lewat **`PL_NUMBER_EDM`** maupun
      **`PL_NUMBER`**. *(AC 50 spec)*
- [ ] Nilai uang yang ditulis hulu dibaca hilir **identik**, termasuk nilai negatif — tidak ada
      pembulatan di perbatasan. (**ADR-U-0003**)
- [ ] Bentuk kedua rekam ditandai di kode sebagai **kontrak lintas konteks**; mengubahnya memaksa
      pembaruan sadar di sisi Claim Life. *(AC 51 spec; **ADR-U-0001**)*
- [ ] Ada **test kontrak** yang menembus dari simpan endorsement sampai pencarian bergaya Claim Life

##### Hilir membaca tabel ⚠️ BARU 2026-09-16 — spec §16

- [ ] ⚠️ Claim Life membaca peserta dari **tabel relasional**, bukan dari CLOB JSON. *(AC 54 spec;
      penyimpangan sadar 8)*
- [ ] ⚠️ Baris ber-`EDMSTATUS` **`Delete`** dan **`Batal`** **tidak pernah** muncul sebagai kandidat
      peserta klaim — keduanya tetap tersimpan, tetapi tersaring di jalur baca klaim. *(AC 66, 67
      spec; kontrak §14)*
- [ ] ⚠️ Baris bernilai **negatif** hasil jurnal balik **tidak pernah** sampai ke layar klaim maupun
      ke perhitungan klaim, sementara jalur akuntansi tetap melihat **seluruh** baris. *(kontrak §14)*
      — **satu seam**, API HTTP, terhadap **skema uji Oracle nyata**.

#### Catatan — penegakan penyaring ada di konteks Claim Life

⚠️ `[terverifikasi]` `GetPesertaClaim_sql1` hari ini **tidak memuat penyaring status apa pun**;
sensus modul `Claim Life/` menemukan **nol** kemunculan `EDMSTATUS`. Apakah Pega menyaring di
lapisan lain **tidak terbukti dari korpus** — **jangan menebak mekanismenya**.

**Penegakannya milik konteks Claim Life:** spec Claim Life **§16** + **AC 25–30**, dan tiket
**CL-02**. Tiket ini mengikatnya sebagai **syarat kontrak** yang dibuktikan lewat test kontrak
lintas konteks. Diikat juga di **PL-08**.

#### Blocker

**Tidak ada pemblokir.** ⚠️ Satu catatan yang **tidak** memblokir: tipe dan nullability `EDMSTATUS`
belum terbaca (`NULL` atau string kosong) — **OQ-001 (sisa)**, pemilik **DBA**. Sampai dipastikan,
implementasi menangani **keduanya** (`IS NULL` *atau* `= ''`).

#### Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```

## Endorsement Life - 12 - Migrasi skema — tabel premium Life dan kolom penanda `EDMSTATUS` [WONTFIX]

**Status:** wontfix — **digantikan tiket 00 + PremiumList Life tiket 00**

**Blocked by:** —

> ⛔ **JANGAN KERJAKAN TIKET INI.** `[keputusan work owner]` Premisnya — merancang skema target
> sebagai salinan tabel premium existing, menunggu DDL fisik — **dibatalkan 2026-09-16**.
>
> **Ke mana isinya pindah:**
>
> | Yang dulu di sini | Sekarang |
> | --- | --- |
> | DDL tabel premium target | **PremiumList Life tiket `00`** — tujuh tabel dirancang sendiri |
> | Kolom penanda `EDMSTATUS` + kolom EDM lain | **tiket `00` konteks ini** — `ALTER` kolom nullable |
> | Pemindahan data endorsement lama | **PremiumList Life tiket `00`** — endorsement adalah **versi** polis, jadi ia pindah bersama polisnya (AC 55, 56 spec) |
> | `JSON_POLIS` / CLOB | **dibuang** — tidak dibawa sama sekali (AC 54 spec) |
>
> ⚠️ **OQ-001 ditutup** — ia yang dulu menahan tiket ini. Termasuk pertanyaan **"apakah kolom uang
> menerima nilai negatif"**: terjawab oleh rancangan sendiri — kolom uang **desimal bertanda**, dan
> jurnal balik memang menulis minus.
>
> Berkas dipertahankan sebagai jejak keputusan, **bukan** sebagai pekerjaan.

---

<details>
<summary>Isi asli (sudah tidak berlaku)</summary>

#### Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin skema Oracle sistem baru menyimpan setiap nilai endorsement
dengan presisi yang sama persis — **termasuk nilai negatif hasil jurnal balik** — dan menyimpan
penanda status peserta dengan bentuk yang sama, sehingga rekonsiliasi tidak menemukan selisih dan
penyaringan peserta hidup berperilaku identik sebelum dan sesudah pindah.

#### Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL tabel target + index |
| `internal/repository` | Pemetaan tipe kolom ↔ desimal presisi arbitrer; bentuk penyaring `EDMSTATUS` |
| — | Skrip rekonsiliasi |

#### Rule Pega sumber

| Tabel | Penulis yang terbukti | Bukti |
| --- | --- | --- |
| `POOLDATA.M_LIFE_PREMIUM_DETAIL` | `SaveMasterLPDet` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL`), PK `M_LIFE_PREMIUM_DETAIL_SEQ.nextval` — **satu rule, dua jalur** | `Endorsement Life/RDBList/SaveMasterLPDet.xml` **dan** `PremiumList Life/RDBList/SaveMasterLPDet.xml` |
| `POOLDATA.M_LIFE_PREMIUM_SUMMARY` | `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY` lewat `InsertPLSummary` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY`) | kedua modul |
| `POOLDATA.JSON_POLIS` | `InsertJsonPolisEDM` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLISEDM`) — **`INSERT` polos**; jalur NB memakai procedure upsert `INSERTJSONPOLISLIFE` | `Endorsement Life/RDBList/InsertJsonPolisEDM.xml` |
| `POOLDATA.JSON_OFFER_LIFE` | `POOLDATA.INSERTJSONOFFERLIFE` lewat `SaveOfferJsonLife_SQL` | `PremiumList Life/RDBList/SaveOfferJsonLife_SQL.xml` |
| `POOLDATA.LIFEINPRODUCTION` | `SaveLifeinProduction_SQL` (`ASM-FW-GISFW-WORK-LIFE` / `ASM!SAVELIFEINPRODUCTION_SQL`) | `Endorsement Life/RDBList/SaveLifeinProduction_SQL.xml` |
| `ARASAPAS.DETAIL_INVOICE` | **dibaca saja**, lintas skema | `Endorsement Life/RDBList/SearcStatusBayarArasaps_SQL.xml` |

`[terverifikasi]` **Kolom `M_LIFE_PREMIUM_DETAIL` terbaca lengkap** dari daftar `INSERT`
`SaveMasterLPDet` — termasuk `PL_NUMBER`, `PL_NUMBER_EDM`, `CERTIFICATE_NO`, `NAME_OF_INSURED`,
`CURRENCY`, `RATE`, `PRORATETYPE`, `SUM_AT_RISK_GROSS`, `SUM_AT_RISK_RETRO`, `RETROCEDED_SHARE`,
seluruh kelompok uang gross / `*_REFUND` / `*_RETRO` / `*_REFUND_RETRO`, **`IDPEGA`**,
**`EDMSTATUS`**, **`STATUSOLD`**, **`STATUS`**, `EM_PERCENT`, `RISK`.

`[data DBA]` Kolom `M_LIFE_PREMIUM_SUMMARY` (**37 kolom**), `JSON_POLIS`, dan `JSON_OFFER_LIFE`
sudah diketahui dari body procedure yang diserahkan 2026-09-15.

#### ADR terkait

**ADR-U-0003** (uang non-float; DDL `NUMBER` tanpa presisi → desimal presisi arbitrer di aplikasi),
**ADR-U-0009** (migrasi penuh; koeksistensi ditolak), **ADR-U-0011**, **ADR-U-0010** (lampiran tetap di
Google Storage — **bukan** bagian migrasi tabel ini).

#### Acceptance criteria

*(belum dapat difinalkan — menunggu OQ-001; disusun agar siap dijalankan begitu DDL turun)*

- [ ] DDL tabel target menyalin **tipe, presisi, PK, index, dan nullability** dari tabel sumber apa
      adanya; tidak ada kolom uang yang menjadi `FLOAT`/`BINARY_DOUBLE`. (**ADR-U-0003**)
- [ ] ⚠️ **Kolom uang menerima nilai NEGATIF tanpa kehilangan presisi** — jurnal balik endorsement
      menulis 32 kolom bertanda minus. Rekonsiliasi membandingkan nilai lama dan baru **secara
      tepat**, bukan dengan toleransi.
- [ ] ⚠️ **Tipe dan nullability `EDMSTATUS` ditetapkan** — dan bentuk penyaring "peserta hidup"
      disesuaikan: `IS NULL` atau `= ''` untuk baris new business. Test penyaring dijalankan ulang
      terhadap skema hasil migrasi.
- [ ] Kolom `CLOB` (`DATA_JSON`, `JSONDATA`) pindah utuh, termasuk isi yang panjang.
- [ ] Sequence (`M_LIFE_PREMIUM_DETAIL_SEQ`, `M_LIFE_PREMIUM_SUMMARY_SEQ`, `JSON_OFFER_SEQ`)
      dipindahkan dengan **nilai berjalan yang benar**, sehingga nomor pasca-migrasi tidak pernah
      bertabrakan dengan nomor lama.
- [ ] Index yang menopang kueri hilir ada sejak hari pertama — khususnya
      `M_LIFE_PREMIUM_DETAIL(PL_NUMBER)` yang dipakai Claim Life, dan kolom `EDMSTATUS` bila
      penyaringan menuntutnya.
- [ ] `PRODKE` pindah utuh, dan penomoran endorsement pasca-migrasi melanjutkan urutan yang benar.
- [ ] Migrasi dapat dijalankan ulang dengan aman dan punya jalur mundur yang diuji.
- [ ] Skema uji yang dipakai seluruh tiket lain dibangun **dari DDL yang sama** — bukan dari tiruan
      yang ditulis terpisah.

#### Blocker

🚧 **`needs-info` — OQ-001 (sisa) terbuka.** Pemilik: **DBA**.

**Yang sudah ada:** nama kolom keempat tabel — `M_LIFE_PREMIUM_DETAIL` terbaca lengkap dari korpus;
`M_LIFE_PREMIUM_SUMMARY`, `JSON_POLIS`, `JSON_OFFER_LIFE` dari body procedure `[data DBA]`.

**Yang belum ada:**

1. **DDL fisik** — tipe, presisi, PK, index, nullability. ⚠️ Seluruh parameter
   `PEGA_M_LIFE_PREMIUM_SUMMARY` bertipe `VARCHAR2` **termasuk kolom uang**, sehingga tipe kolom
   sebenarnya belum diketahui.
2. **Tipe dan nullability `EDMSTATUS`** — apakah baris new business menyimpan `NULL` atau string
   kosong `''`. Salah pilih akan **membuang seluruh peserta new business** dari layar klaim.

**JANGAN paksa `ready`.** Pola sama dengan **CL-13** (Claim Life) dan **PL-09** (PremiumList Life).

**Tidak memblokir tiket lain** — kolomnya sudah cukup untuk seluruh tiket 01–11.

#### Perintah verifikasi

```
go test ./internal/...
```

</details>
