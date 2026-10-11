# Paritas aksi — PremiumList Life

> Berkas ini lahir 28-09-2026 untuk sensus remark (GILIRAN-12 paket 2); modul ini sebelumnya tanpa
> berkas paritas. Paritas layar per tiket tetap di `issues/`.

## Sensus remark 28-09-2026 (GILIRAN-12 paket 2)

Setiap activity PremiumList Life dicetak beserta `pyStepsBlockName`; `//` = ter-remark (beserta
anaknya). 38 activity, **9** ber-remark, **36** langkah `//` (label `SHOWFAILURE`, `AA`, `Err`, `REM`,
`COMMITSTEP` dll. bukan remark).

**Cara hitung** `[terverifikasi]` — dua cara, jendela = seluruh `PremiumList Life/Activity/*.xml`:

| Cara | Perintah | activity | ber-remark | langkah `//` |
| --- | --- | ---: | ---: | ---: |
| A — pengurai XML (expat), per `rowdata` ber-`pyStepPageReference` yang `pyStepsBlockName` tepat `//` | skrip sensus sesi (cetak activity, nomor langkah, bNNN) | 38 | 9 | 36 |
| B — grep mentah per berkas | `grep -o "<pyStepsBlockName>//</pyStepsBlockName>" "<berkas>" \| wc -l`, dijumlah | 38 | 9 | 36 |

Alat diuji lebih dulu pada item yang sudah diketahui: `Claim Life/Activity/SaveOutStandingLife_Act.xml`
— 8 langkah `//` (dibaca ulang pemeriksa independen GILIRAN-11) — kedua cara menjawab 8. Label
lompatan lain (`EXIT`, `send`, `AA`, …) bukan remark dan tidak terhitung.

⚠️ **Cakupan celah langkah hidup DINYATAKAN:** sensus ini mendaftar langkah hidup yang tidak ditiru
untuk activity **ber-remark** saja. Activity tanpa remark tidak diaudit ulang langkah demi langkah di
giliran ini; statusnya tetap seperti baris paritasnya (bila ada) — **belum** audit penuh.


| Activity | Langkah `//` (b) | Ditiru? | Keputusan |
| --- | --- | --- | --- |
| `ProtectAccept` | 4.1 b720 (NoOffer) · 6 b2340 (+6.1–6.5; 6.3 b2756) | **ya — dibuang** (`models/polis_validasi.go`) | tiket 01 ralat; gerbang posisi dibetulkan |
| `ValidasiUploadPL_act` | 9.2 b8025 · 9.3 b8198 · 9.4 b8368 · 9.5 b8517 · 9.11 b9667 · 10 b13018 · 11 b13171 · 12 b13322 · 13 b13475 · 19 b14393 · 33 b16535 | **ya — dibuang** (sertifikat terpakai/ganda, `MEDICAL_STATUS`, lampiran; + duplikat nama+DOB) | tiket 04 ralat, AC disunting di tempat |
| `SubmitPremiumList_Act` | 15 b3398 (geser `MM.YYYY.` bila di atas tanggal closing) | tidak — tetapi `models/polis_periode.go` mengutip b3491 | kutipan diralat; premis "dua versi hidup" → OQ-PL-13 |
| `InsertJsonPolisLife_Act` | 16 b5294 `Commit` · 17 b5383 `Connect-REST InsertLifePremiumDetail` | tidak (Go memegang transaksi; salinan peserta lewat `repository/polis_warisan.go`, keputusan pl2) | nama layanan step 17 diralat di spec, 05b, 06 → OQ-PL-14 |
| `SavePremiumList_Act` | 8.3 b1891 · 8.5 b3704 · 8.6 b3863 · 8.7 b4050 · 8.9 b5887 · 11 b6411 · 13 b6642 | tidak | tercatat tiket 03/05a. Langkah HIDUP 6–8.2 (batas umur / sum insured) kini **penolakan Validate CSV** — lihat bagian di bawah (keputusan work owner 05-10-2026) |
| `UploadCSVLifePremium_Act` | 5 b889 · 6 b1042 · 7 b1188 (+7.2 b2517, 7.4 b2797) · 8 b2922 | tidak (tinjauan tidak menulis) | — · ⚠️ beda dari Pega (05-10-2026): `POLICY_HOLDER`/`PLAN` peserta diisi dari polis (`POLICY_HOLDER_NAME`, `BUSINESS_NAME`), nilai CSV diabaikan — tiket 04 |
| `InputOfferLife_ACT` | 1 b354 (`ASMForceCaseClose`) | tidak — Decline menutup lewat flow | tercatat tiket 01 |
| `LoadAttachmentData` | 1 b276 · 2 b493 · 3 b661 · 4 b771 | tidak | — |
| `serviceInsertArasapasLife_act` | 9 b1597 (`UpdateErrorNoteJsonPolis`) | tidak (efek stub) | — |

**Langkah HIDUP yang tidak ditiru dan belum tercatat** (kini dicatat sebagai celah):

- `ProtectAccept` 3 (b624, `EmailTypePL==1` atau KELUAR) dan 8 (b3792, daftar mata uang kosong di
  Premium) — tiket 01. `ValidasiPenawaran` sendiri belum tersambung ke rute mana pun.
- `ValidasiUploadPL_act` 2 — mengisi 0 kolom uang yang kosong (sub-langkah 2.1–2.32) → **OQ-PL-12**;
  3–8 (b6874–b7767) — keberadaan medan header; pesannya ada di `models/polis_unggah.go`, belum dipakai.
- ~~`SavePremiumList_Act` 8.1 (b1635, batas umur) dan 8.2 (b1833, batas sum insured; dilewati TP/TR
  berkode RNML-FL b1803), masing-masing dengan keluar di langkah 9 (b6146).~~ — ✅ ditiru 05-10-2026
  sebagai penolakan Validate CSV (lihat bagian `SavePremiumList_Act` 6–8.2 di bawah).
- `Calculate1_Act` 7.5–7.9 — `CekDoubleInsured` untuk akumulasi retensi ceding (`SumCRetensi` lawan
  `MaxCRetensi`).

**OQ baru (untuk work owner):**

- **OQ-PL-12** — langkah 2 `ValidasiUploadPL_act` menjadikan kolom uang yang kosong **0** sebelum
  pemeriksaan "HARUS ADA"; aplikasi menolaknya (ADR-U-0027: kosong bukan nol). Ikuti Pega (terima
  sebagai 0) atau pertahankan penolakan sebagai penyimpangan sadar? — ✅ **ditutup 29-09-2026 (GILIRAN-17): ikuti
  Pega, `IsiNolUangKosong` sebelum validasi (tiket 04).**
- **OQ-PL-13** — "ikuti yang dari DB" berdiri di atas premis bahwa `SubmitPremiumList_Act` membaca
  ambang secara hidup; nyatanya ia hanya mengalir ke langkah 15 yang ter-remark. Pembaca hidup tabel
  itu adalah `PROC_GENERATE_SEQUENCE_NUMBER` (nomor PL); `ProdDateTime` di Pega memakai `>25`
  tertanam. Apakah ambang DB juga berlaku untuk `ProdDateTime`? — ✅ **ditutup 29-09-2026 (GILIRAN-17): ikut XML, 25
  tertanam (b1170); `ProdDateTime` tidak dihitung aplikasi (tiket 02).**
- **OQ-PL-14** — cerita spec 34 membuang `ConvertJsonNusareToProduction` dengan alasan step 17
  ter-remark; step 17 memanggil `InsertLifePremiumDetail`, dan `convertJsonNusareToProduction` hidup
  di `serviceInsertArasapasLife_act` langkah 5 (b963). Tetap dibuang? — ✅ **ditutup 29-09-2026 (GILIRAN-17):
  ditiru lewat outbox, pelaksana stub `PelaksanaPremiumList` (tiket 06).**

## Pertanyaan terbuka GILIRAN-13/14 — 29-09-2026

- **OQ-PL-15** — ✅ **DITUTUP 29-09-2026** (work owner, "ikuti rekomendasi"): migrasi **058** memulai ulang
  `SEQ_WORK_POLIS` pada **22374** (nomor `NBLF-` tertinggi terlihat 22373 + 1). `[sementara — DBA memastikan
  pyLastReservedID awalan NBLF- di PC_DATA_UNIQUEID sebelum data nyata]`. *(Tiket 00, GILIRAN-15.)*
- **OQ-PL-17** *(untuk DBA — lahir dari tinjauan GILIRAN-15)* — pastikan `pyLastReservedID` awalan `NBLF-` di
  `PC_DATA_UNIQUEID` (penghitung Pega, tidak terlihat dari akun `POOLDATA`) tidak melampaui 22373 sebelum data nyata masuk
  `T_WORK_POLIS`. Bila melampaui, `SEQ_WORK_POLIS` dimajukan lagi. Label `[sementara]` 058 bergantung padanya.
- **OQ-PL-16** — ✅ **DITUTUP 29-09-2026, butir bq** `[DIPUTUSKAN — XML; veto work owner]`: `Decision3` dirutekan dari
  `FLAG_ONGOING_POLICY` (decision table `IsFlagOnGoingPolicy`: "0" → Offer, "1" → Premium, otherwise `Decline` tanpa
  konektor → 409). *(Tiket 01.)*

## `SetRateLIfePremium_Act` — ditiru untuk Type QR (keputusan work owner 05-10-2026)

| Activity | Langkah | Ditiru? | Di mana |
| --- | --- | --- | --- |
| `SetRateLIfePremium_Act` | 1–10.4 (produk `PRODUCTINWARD_LIFE`/`PRODUCT_LIFE`, Find RIRATE / RIRISK, set Rate, set Sum at Risk) | **ya — hanya Type QR**, ditulis ulang menurut spesifikasi work owner | `models/polis_hitungqr.go` (rumus murni), `repository/polis_hitungqr.go` (bahan), `services/polis_unggah.go` `periksaDanHitung` (**Calculate CSV saja**, sebelum transaksi; Validate CSV hanya memeriksa bentuk dan batas usia / sum insured — keputusan work owner 05-10-2026) |
| `SetRateLIfePremium_Act` | 10.5 dst. (spreading retro, summary retro) dan kolom `*_REFUND` | tidak | — |

Type QP/TP/TR tidak dihitung: seluruh kolom peserta tetap dari CSV.

**Sumber:** parameter produk dari tabel flat Master Product Name Life `M_PRODUCTNAME_LIFE`
(`CEDINGRETENTIONNUM`, `CEDINGLIMIT`, `RNMSHARE`, `RICOMM`, `BROKERAGE`, `RIRISKID`), plan
`M_PRODUCTNAME_LIFE_PLAN` (`NAME` = Class of Business polis → `RIRATEID`), tabel `M_RATE_LIFE` dan
`RIRISK_LIFE` dibaca **utuh** per `IDUSEDBY` (tanpa batas 500 baris milik tampilan) — bukan view lama
`PRODUCTINWARD_LIFE`/`PRODUCT_LIFE` seperti XML.

**Rumus** (presisi penuh apd, kolom hasil dibulatkan half-up 4 desimal saat disimpan, `RISK` 7):
CONTRACT = 1 bila `PERIOD_MM` < 12, selain itu round(`PERIOD_MM`/12) · `CEDING_RETENTION` = SI ×
Ceding Retention % (maks Ceding's Limit) · `SUM_REASURED` = SI − retensi · `SHARE_NUSANTARA_RE` =
SUM_REASURED × RNM Share % · tahun pertama (`BEGIN_DATE` = `GROSS_VALUATION_BEGIN_DATE`):
`SUM_AT_RISK_GROSS` = share, `RISK` kosong; selain itu `RISK` = risk master / 1000 dan
`SUM_AT_RISK_GROSS` = share × risk master / 1000 · `GROSS_PREMIUM` = RATE/1000 × share × FACTOR
(kosong/0 → 1) × (EM_PERCENT + 1) · `DEDUCTION` = gross × Deduction % · `BROKERAGE_FEE` = gross ×
Brokerage % · `NET_PREMIUM` = gross − deduction − brokerage · `RI_ADMIN_FEE` dari CSV (kosong → 0, tidak
mengurangi net) · langkah 10.3.3: `SUM_AT_RISK_RETRO` = `RETROCEDED_SHARE` = `SUM_AT_RISK_GROSS` (keputusan
work owner 06-10-2026; **beda dari Pega**: Pega = SHARE × `local.Risk` dan hanya bila produk punya Outward List).
Type lain: kedua kolom dari CSV apa adanya (kosong → NULL). Kolom hasil yang terisi di CSV diabaikan dan ditimpa.

**Perbedaan dari XML (disengaja):**

- `CURRENT_AGE` tidak dipakai — AGE rate = `ENTRY_AGE`.
- YEAR risk = tahun(`GROSS_VALUATION_BEGIN_DATE`) − tahun(`BEGIN_DATE`), **tanpa +1**.
- `EM_PERCENT` dibaca sebagai **pecahan** (0.5 = 50%).
- Rate / risk **tidak ditemukan atau lebih dari satu** baris cocok → baris **ditolak** dengan pesan yang
  menyebut age/contract/sex/year yang dicari (Pega diam-diam memberi 0 atau memakai nilai peserta
  sebelumnya).
- Master rate yang **mencampur** GENDER U dengan F/M → seluruh unggahan ditolak.
- **Ceding's Limit wajib** terisi dan > 0; tepat satu plan harus cocok dengan Class of Business;
  `RIRATEID` plan dan `RIRISKID` produk wajib terisi (syarat P1–P4, galat 409 berpesan, nol tulisan).
- Pencocokan rate: master U → AGE + CONTRACT (SEX tidak dibaca); master F/M → Pro Rate Type 1:
  AGE + CONTRACT + SEX, Pro Rate Type 2/3: AGE + SEX (Pro Rate Type kosong ditolak).
- Judul CSV Type QR wajib memuat `PERIOD_MM`; isinya (bilangan bulat ≥ 1) dan `SEX` (F/M untuk master
  ber-gender) diperiksa saat hitung.

## `SavePremiumList_Act` 6–8.2 — batas produk di Validate CSV (keputusan work owner 05-10-2026)

| Langkah | Pega | Aplikasi |
| --- | --- | --- |
| 6–7 batas produk | `GetRateProductLife`: `SELECT * FROM PRODUCTINWARD_LIFE WHERE ID = ProductNameID` | `M_PRODUCTNAME_LIFE` (`MINAGE`, `MAXAGE`, `MINSUMINSURED`, `MAXSUMINSURED`), kolom disebut satu per satu — `repository/polis_hitungqr.go` `BatasProduk` |
| 8.1 Protect Age | `ENTRY_AGE` < MINAGE atau > MAXAGE → `<NAME_OF_INSURED> Age exceeds the limit, at list <idx>` | sama, pesan VERBATIM, `at list` = nomor baris CSV, kolom `ENTRY_AGE` |
| 8.2 Protect Sum Insured | `SUM_INSURED` di luar MIN/MAXSUMINSURED; dilewati TP/TR ber-R/I SLIP `RNML-FL` | sama, kolom `SUM_INSURED` |
| 9 / 15 | `Page-Set-Messages` lalu `Obj-Save WithErrors=true` — **pesan saja, data tetap tersimpan** | **penolakan** `models.PeriksaBatasProduk` di `periksaDanHitung` (Validate CSV dan Calculate CSV, semua Type): Validate tidak lolos → Calculate CSV tidak dapat diklik; nol tulisan |

Perbedaan dari Pega (disengaja): sumber batas `M_PRODUCTNAME_LIFE` (bukan `PRODUCTINWARD_LIFE`); nilai
peserta dibaca dari baris CSV (bukan tabel peserta tersimpan); melewati batas **memblokir** penyimpanan.
Batas kosong dan nilai peserta yang kosong / sudah ditolak validasi bentuk tidak diperiksa; produk kosong
atau tidak ada di `M_PRODUCTNAME_LIFE` tidak diperiksa (Type QR sudah ditolak syarat hitung). Langkah cek
batas terpisah sesudah Calculate CSV dihapus.

## `SetRateLIfePremium_Act` 10.8–10.12 — spreading retro Type QR (keputusan work owner 06-10-2026)

Ditiru untuk **Type QR** di Calculate CSV (`periksaDanHitung` sesudah hitung premi → `models/polis_spreading.go`,
bahan `repository/polis_spreading.go`), disimpan di transaksi yang sama dengan peserta ke
`T_PREMIUM_LIST_SPREADING` (053, satu baris per lapisan berbagian > 0) dan `T_PREMIUM_LIST_SPREADING_RETRO` (054, satu
baris per reinsurer). Rekap **Retro by Reinsurer** tetap di tab Summary.

**Tampilan per peserta (keputusan work owner 07-10-2026):** tab **Spreading** (semua peserta, berhalaman) DIHAPUS beserta
rutenya. Spreading dilihat dengan mengklik (atau Enter/Space) satu baris grid **Participant Details** → popup
`Spreading - <Name Of Insured>` (`components/ModalSpreadingPeserta.tsx`): kepala Name Of Insured, Certificate No, Currency,
Retroceded Share peserta, treaty year (label `TREATYYEAR` + periode `TREATYYEAR_LIFE`); isi lapisan → reinsurer. Sumber:
`GET /api/polis-life/{id}/peserta/{detailId}/spreading` (baca-saja; peserta bukan milik polis → 404; tanpa spreading →
lapisan kosong, pesan "No retro spreading for this participant."). Semua angka lewat satu formatter grid (`pemisahRibuan`,
koma desimal untuk bentuk teks apa pun dari server).

**Perbedaan dari Pega (disengaja):**

- Treaty year dipilih dari **periode** `TREATYYEAR_LIFE` yang memuat `BEGIN_DATE` produk (`M_PRODUCTNAME_LIFE`), bukan
  dari Outward List / ORS; Underwriting/Transaction Year tidak dipakai. Nol atau lebih dari satu → 409.
- Lapisan = seluruh kontrak treaty year itu, **kaskade Maximum Limit** (IDR/USD sesuai mata uang peserta) menaik;
  nilai yang di-spreading = `RETROCEDED_SHARE` presisi penuh.
- **Penolakan**: melebihi Max lapisan terakhir, mata uang selain IDR/USD, rate retro tidak ada/ganda (baris); Business
  tidak ada di Business List, total share ≠ 100%, Max Limit kosong, master rate tidak sah (409) — hanya untuk lapisan
  yang terpakai.
- Rate retro dicari dengan fungsi yang SAMA dengan rate premi (`MasterRate.Cari`); GROSS **tanpa FACTOR**.
- Tahun pertama: Discount = GROSS × COMMISION%, Comm = (GROSS − Discount) × OVR_COMM%; bukan tahun pertama: Comm =
  GROSS × OVR_COMM%. Nominal Discount/Comm tidak disimpan (tidak ada kolomnya).
- Berhenti di tingkat reinsurer (`TREATYSECURITYREINSURER_LIFE` tidak dipakai).
- Rekap retro per reinsurer per mata uang dari data — **tanpa** empat nama tetap Pega (NUSARE/HANNOVER/TOA/GIBRALTAR).
- Validate CSV **tidak** menjalankan spreading (keputusan work owner 05-10-2026: Validate hanya bentuk + batas usia/SI).

## `InsertJsonPolisLife_Act` langkah 9/11 — `LIFEINPRODUCTION` (keputusan work owner 06-10-2026)

Ditiru dengan **SQL langsung** di transaksi Confirm (`SummaryPremiumList.simpanDalam` langkah 6, sesudah PL Number, WPC
final, summary dan peserta warisan) — `repository/polis_produksi.go`. Tabel WARISAN POOLDATA, nol DDL, nol COMMIT.

- Kolom = daftar INSERT `SaveLifeinProduction_SQL` (36) **tanpa** `LAYER_1..LAYER_4` (NULL — "tidak perlu digunakan
  lagi"); nilai dari baris utama `T_PREMIUM_LIST` (`ID = ID_PEGA`). `IDPEGA` = pengenal kasus (pola tabel warisan lain di
  modul ini), `NOPOLIS` = PL Number, `TYPECEDINGNAME` = QS/SURPLUS/QS + SURPLUS/XOL dari `TYPE_CEDING`,
  `DATERECEIVED`/`WPC` tanpa jam, `TGL_INPUT` = `SYSDATE`.
- Sekali per kasus: baris lama `IDPEGA` yang sama dihapus dulu — Confirm ulang tidak menggandakan.
- `JSON_OFFER_LIFE` ditiru saat **Save Offer** dan `JSON_POLIS` saat Confirm — dua bagian di bawah.

## `InputOfferLife_ACT` langkah 3-4 — `JSON_OFFER_LIFE` (keputusan work owner 06-10-2026)

Ditiru dengan **SQL langsung** di transaksi **Save Offer** (`FormPenawaran.Simpan`, sesudah isian tersimpan ke baris utama
`T_PREMIUM_LIST`, sebelum riwayat `T_VIEW_SUGGEST`) — `repository/polis_jsonoffer.go`. Prosedur
`INSERTJSONOFFERLIFE` **tidak** dipanggil. Kolom = `ALL_TAB_COLUMNS` POOLDATA (06-10-2026).

- Satu baris per kasus: baca ID lama → `DELETE WHERE IDPEGA` → `INSERT … SELECT` dari baris utama (isian TERAKHIR).
  **ID lama dipertahankan**; kasus baru: `TO_CHAR(JSON_OFFER_SEQ.NEXTVAL)` — sequence warisan yang SUDAH ada (PK prosedur,
  `[data DBA]` tiket 09), nol sequence baru.
- `IDPEGA` = pengenal kasus (sama dengan `LIFEINPRODUCTION`/`JSON_POLIS`); `REINSURANCETYPE` ← `JENIS_ASURANSI`;
  `TYPECEDING` = QS/SURPLUS/QS + SURPLUS/XOL (CARI8); `STATUS` ← `STATUS_PENAWARAN`; `MARKETING_NOTE` ←
  `KETERANGAN_MARKETING`; `OPERATORNAME` ← `CREATE_OP_NAME`; `UW_NOTE` ← `KETENTUAN_UNDERWRITING`; delapan tanggal
  **tanpa jam** (Pega mengirim `DD/MM/YYYY HH24:MI:SS`); `TGL_INPUT` = `SYSDATE`.
- ⚠️ `JSONDATA` dan `OLDID` (`p_OFFERID`) **NULL**. `NO_OFFER` (`GetIdOffer_SQL` → `'OFR-'||ID`, langkah 5-6) tetap
  **tidak** diisi.
- Lebar kolom diperiksa di Go (`LENGTHB` baris utama) sebelum menulis: melebihi (mis. `STATUS`/`BUSINESSCODE`
  VARCHAR2(10)) → **400** berpesan, bukan ORA-12899. Gagal tulis apa pun → seluruh Save Offer batal.

## Save Offer Status `Bind` / `Closed` menutup kasus (keputusan work owner 08-10-2026)

⚠️ **Menyimpang dari Pega, bukan korpus.** Hanya kasus tombol "Input Offer" (bendera `"0"`, `FlagPolisPenawaran`);
kasus "Input Premium" (bendera `"1"`, radio EmailTypePL) tidak berubah.

- `Accept` / `Pending` / `Decline` → Save Offer seperti sebelumnya (baris utama + baris status + `T_VIEW_SUGGEST` +
  `JSON_OFFER_LIFE`); kasus tetap terbuka.
- `Bind` → simpan seperti di atas, lalu **ditutup Resolved-Completed** = akibat `Confirm` jalur Offer (Decision3
  "Offer" → END52). **Tanpa** `SimpanPolis`: nomor PL, `LIFEINPRODUCTION`, `JSON_POLIS`, `M_LIFE_PREMIUM_*` tidak
  terbit (jalur Offer tidak melewati `Utility1`), dan tanpa efek keluar. Gerbang Confirm tahap penawaran (isian wajib)
  berlaku atas isian yang sedang disimpan.
- `Closed` → simpan, lalu **ditutup Resolved-Rejected** = akibat `Decline` (End1).
- Aturannya `models.AkibatStatusOffer` (diturunkan dari `TransisiPenawaran` + `PenggolongOtomatis`); penutupan memakai
  `Penawaran.terapkanDalam` — jalur Confirm/Decline yang SAMA (`TutupKasus` + jejak `T_CLAIMLF_JEJAK`, sebab
  `Save Offer Bind -> Confirm -> Offer` / `Save Offer Closed -> Decline`) — di transaksi Save Offer: gagal menutup =
  seluruh Save Offer batal. Kasus tertutup / bukan tahap penawaran → 409 sebelum tulisan apa pun (gerbang lama Save Offer).
- Layar: Bind/Closed → popup konfirmasi ("Save the offer and close the case as Resolved-Completed?" /
  "…Resolved-Rejected?", Cancel / Yes); sesudah tutup kalimat "Case closed — <status>." lalu kembali ke kotak masuk.
  ~~Tombol Confirm/Decline tidak berubah~~ — **sejak 08-10-2026 (keputusan work owner) tombol `Confirm` TIDAK dirender di
  tahap penawaran untuk kasus bendera `"0"`** (`InputOffer.tsx` `confirmTampil`; bendera dari data penawaran) — yang tersisa
  di panel Decision hanya `Decline`. Server menolak `Confirm` tahap penawaran bendera `"0"` dengan **409** "Confirm is not
  available for offer cases. Choose Status Bind and press Save Offer." (`models.ConfirmPenawaranTersedia`,
  `ErrConfirmPenawaranTidakTersedia`) sebelum gerbang lain dan nol tulisan; Save Offer Bind tidak lewat jalur itu. Bendera
  `"1"` dan tahap Input Premium Detail tidak berubah.

## `InsertJsonPolisLife_Act` langkah 10 — `JSON_POLIS` (keputusan work owner 06-10-2026)

Ditiru dengan **SQL langsung** di transaksi Confirm (`simpanDalam` langkah 7, sesudah PL Number dan `LIFEINPRODUCTION`)
— `repository/polis_jsonpolis.go`. Prosedur `INSERTJSONPOLISLIFE` **tidak** dipanggil.

- Diisi HANYA `IDPEGA` (= `LIFEINPRODUCTION`), `NOPOLIS` = PL Number, `TGL_PROD`, `USERNAME` ← `CREATE_OP_NAME`
  (= `CREATEOPNAME`), `TGL_INPUT` = `SYSDATE`, dan **`PRODKE` = `'0'`** (teks, VARCHAR2(5) — keputusan work owner
  08-10-2026; sebelumnya NULL). ⚠️ **`DATA_JSON` NULL**; `NOENDORS`, `ERR_NOTE`, `TGL_KONVERSI`, `STS_KONVERSI`,
  `TGL_KONVERSI_RETRO`, `STS_KONVERSI_RETRO`, `OLDNOPOLIS` NULL (layanan hilir).
- Dampak `PRODKE '0'` ke Endorsement Life (pembaca `JSON_POLIS`, kodenya tidak diubah): versi warisan `NVL(PRODKE, 1)`
  = 0 < versi aplikasi 1 (baris aplikasi tetap menang); nomor endorsement memperlakukan kosong = 0, jadi `'0'` sama hasilnya.
- `TGL_PROD` (`models.TanggalProduksi`) — aturan periode PL Number (`penomor.HitungPeriodeNomor`, TANGGAL_CLOSING), bukan
  ambang 25 Pega: periode = bulan transaksi (Jakarta) → tanggal transaksi; periode bergeser → tanggal 1 bulan periode.
  Nomor yang sudah terbit sebelumnya: TGL_PROD lama dipertahankan, bila tidak ada dihitung ulang pada saat Confirm.
- `DELETE WHERE IDPEGA` lalu `INSERT` — Confirm ulang tidak menggandakan; gagal tulis → seluruh Confirm batal.
- ⚠️ `convertJsonNusareToProduction` tetap **stub** (OQ-PL-11): `DATA_JSON` kosong harus dikonfirmasi ke pemilik
  Arasapas sebelum stub dilepas.
