# Spec — PremiumList Life / New Business (migrasi Pega → Go + React + Oracle)

Status: ready-for-agent
Konteks: `premiumlist-life` — "Life — Penawaran & Premium List (new business)" (context-map §2.8)
Modul: **PremiumList Life** (124 berkas). ⛔ **Endorsement Life BUKAN bagian spec ini** —
konteks terpisah, lihat `.scratch/endorsement-life/`.
Tanggal: 2026-09-15 · **direvisi 2026-09-16 — bagian penyimpanan**

> ⚠️ **Revisi penyimpanan 2026-09-16** `[keputusan work owner]`. **Seluruh JSON dibuang**; polis
> pindah ke **tujuh tabel relasional baru** + **`T_WORK_POLIS`** (lintas-lini).
> Yang berubah: **§6** (batas transaksi — titik potong lenyap), **§12** (baru — bentuk
> penyimpanan), **AC 32–50** (baru), §Out of Scope, dan §Pertanyaan terbuka.
> **Yang sengaja TIDAK berubah**: mesin alur (§3), periode tutup buku (§4), uang di batas (§5),
> penomoran (§7), unggah CSV (§8), efek keluar (§9), kode mati (§10), penulisan detail peserta
> inline (§11), dan **AC 1–31**.
> Sumber: `revisi-penyimpanan-premiumlist.md` — ditutup dengan **tiga contoh `DATA_JSON` nyata**
> (TP/EDM, QP/NB, QR/NB).

Sumber: `.scratch/premiumlist-life/grilling-ronde-1.md`, `CONTEXT.md`,
`docs/adr/ADR-0001`–`ADR-0015`, `discovery/modules|flows/PremiumList Life.md`,
`discovery/context-map.md` §2.8, `discovery/open-questions.md`
Skill: `/mattpocock-skills:to-spec`
Revisi: **2026-09-15 (b)** — `[keputusan work owner]` **Endorsement Life dipisah menjadi konteks
sendiri.** Materi endorsement dikeluarkan dari spec ini (US, §7, AC 12, tabel uji, §1 batas konteks,
Out of Scope); **mesin bersama tetap** sebagai komponen. Tiket 07 dipindah ke
`.scratch/endorsement-life/BAHAN-tiket-pl-number-edm.md`.
Revisi: **2026-09-15 (a)** — OQ-068 ditutup; `InsertLifePremiumDetail_act`
(`ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT`) dan `SaveMasterLPDet`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET`) masuk `PremiumList Life/`;
§11 ditambahkan (penulisan detail **inline, bukan job**); AC 28a/28b ditambahkan.

> **Konvensi penandaan.** `[terverifikasi]` = terbukti korpus dengan **class + nama + path**;
> `[keputusan work owner]`; `[data DBA]`; `[terbuka]` = OQ. **Identitas rule wajib menyertakan
> class** — nama sama di class berbeda = rule berbeda.

> **Sumber tunggal.** Korpus Pega `D:\XML\RNM_BRD\` (READ-ONLY) dan artefak di
> `OUTPUT_HASIL_RNM\`, sekaligus repo target tunggal. ⛔ `D:\XML\nusantara-re\` di-blacklist.

---

## Problem Statement

**Premium List Life adalah hulu domain jiwa** — ia yang menghasilkan `PremiumListSummary` dan
`PremiumListDetail` yang kemudian dikonsumsi Claim Life, Komite Claim Life, dan Endorsement.
`[terverifikasi]` Class integrasi `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` dirujuk **46 rule** di Claim
Life, 25 di Endorsement, 19 di PremiumList, 11 di Komite — ia tulang punggung domain Life.

Masalah yang dihadapi bila status quo dipertahankan:

1. **Kode mati bercampur kode hidup, dan penandanya berbohong dua arah.** `[terverifikasi]` Di satu
   activity: dua langkah mati ditandai remark dengan benar, empat gerbang polis lama yang **juga
   mati** masih ditandai **aktif**. Membedakan yang berlaku menuntut penelusuran langkah demi
   langkah — dan konfirmasi manusia (**OQ-066**).
2. **Dua versi aturan tutup buku hidup berdampingan** — satu membaca dari tabel, satu menanam angka
   `25`. Yang **lebih baru** justru yang menanam.
3. **Uang menyeberang batas sebagai teks.** `[data DBA]` Seluruh parameter procedure penulis premium
   summary bertipe `VARCHAR2`, termasuk kolom uang.
4. **Batas transaksi tidak seragam.** `[data DBA]` Dua procedure commit sendiri, dua tidak — sehingga
   satu operasi bisnis dapat terpotong di tengah tanpa ada yang menyatakannya.

## Solution

Membangun konteks **Life — Penawaran & Premium List** sebagai layanan Go + antarmuka React yang
menulis ke Oracle `POOLDATA` yang sama, dengan:

- **Alur penawaran → premium list yang dinyatakan eksplisit**, termasuk arti tiap keputusan.
- **Ambang tutup buku dibaca dari tabel**, tidak pernah ditanam.
- **Uang sebagai desimal presisi arbitrer**, dengan konversi teks ↔ desimal di **satu batas**
  (**ADR-0003**).
- **Batas transaksi dipegang Go secara sadar**, mengetahui procedure mana yang memutusnya.
- **Kode mati tidak ikut pindah.**

Perilaku dibawa apa adanya (paritas), kecuali **tiga penyimpangan sadar**: email diperlakukan
sebagai **jalur alarm** bukan notifikasi bisnis; penulisan detail peserta **inline saat simpan,
bukan job**; dan Go memegang batas transaksi menggantikan commit implisit Pega.

### Kontrak hilir — ke Claim Life

| Arah | Kontrak |
| --- | --- |
| **Keluar** | Konteks ini **menghasilkan** baris `M_LIFE_PREMIUM_SUMMARY` + `LIFE_PREMIUM_DETAIL`, berkunci `PL_NUMBER` (NB) atau `PL_NUMBER_EDM` (endorsement) |
| **Dikonsumsi** | Claim Life membacanya sebagai `PremiumListSummary` / `PremiumListDetail` — struktur yang dipakai mesin status klaim (**ADR-0011**: `PremiumListDetail` adalah cerminan baris `AdjustmentList`) |

**Perubahan bentuk `M_LIFE_PREMIUM_SUMMARY` atau `LIFE_PREMIUM_DETAIL` adalah perubahan kontrak
lintas konteks**, bukan perubahan internal.

---

## User Stories

### Penawaran (Offer)

1. Sebagai **inputor**, saya ingin mendaftarkan penawaran polis jiwa baru, supaya penawaran masuk
   siklus. `[terverifikasi]` `PremiumList Life/InputPolicyHolder.xml`
   (`ASM-FW-GISFW-WORK-LIFE` / `INPUTPOLICYHOLDER` / `RULE-OBJ-FLOW`), shape `Assignment2`
   "Input Offer".
2. Sebagai **inputor**, saya ingin memilih **Confirm / Reject / Decline** atas penawaran, supaya
   keputusan saya menentukan langkah berikutnya. `[terverifikasi]` nilai keluaran
   `ASM-FW-GISFW-WORK-LIFE` / `ISLIFEACCEPTED` / `RULE-DECLARE-DECISIONTABLE`.
3. Sebagai **inputor**, saya ingin **Reject** mengembalikan saya ke layar input, supaya saya dapat
   memperbaiki isian. `[keputusan work owner]`
4. Sebagai **inputor**, saya ingin **Decline** menutup case, supaya penawaran yang tidak dilanjutkan
   berhenti bersih. `[keputusan work owner]`
5. Sebagai **organisasi**, saya ingin keputusan accept/reject/decline **diambil manusia**, bukan
   dihitung formula, supaya penilaian tetap pada yang berwenang. `[keputusan work owner]` — ini
   sebabnya kriteria masukan DecisionTable **tidak perlu direplikasi**.
6. Sebagai **inputor**, saya ingin penawaran yang berhenti di tahap penawaran tetap tersimpan,
   supaya riwayatnya ada. `[terverifikasi]` keluaran `Offer` pada
   `ASM-FW-GISFW-WORK-LIFE` / `ISFLAGONGOINGPOLICY` / `RULE-DECLARE-DECISIONTABLE`.
7. Sebagai **inputor**, saya ingin melanjutkan ke pengisian premium list bila penawaran memang
   berlanjut. `[keputusan work owner]` keluaran `Premium` = `2`.

### Premium List Detail

8. Sebagai **inputor**, saya ingin mengisi rincian premium list per peserta, supaya dasar
   perhitungan premi lengkap. `[terverifikasi]` shape `ASSIGNMENT63` "Input Premium List Detail".
9. Sebagai **inputor**, saya ingin mengunggah daftar peserta lewat **CSV**, supaya tidak mengetik
   ratusan baris. `[terverifikasi]` `PremiumList Life/Activity/UploadCSVLifePremium_Act.xml`
   (`@BASECLASS` / `UPLOADCSVLIFEPREMIUM_ACT` / `RULE-OBJ-ACTIVITY`) — **rule base-class, dipakai
   bersama modul lain**.
10. Sebagai **inputor**, saya ingin unggahan divalidasi sebelum masuk, supaya kesalahan angka
    ketahuan sejak awal. `[terverifikasi]` `PremiumList Life/Activity/ValidasiUploadPL_act.xml`
    (`ASM-FW-GISFW-WORK-LIFE` / `VALIDASIUPLOADPL_ACT` / `RULE-OBJ-ACTIVITY`).
11. Sebagai **inputor**, saya ingin melihat pesan kesalahan yang menyebut **kolom mana** yang salah,
    supaya saya tahu apa yang diperbaiki. `[terverifikasi]` langkah validasi menyebut
    `NET_PREMIUM`, `GROSS_PREMIUM`, `SHARE_NUSANTARA_RE`, `SUM_INSURED`, `CEDING_RETENTION`,
    `SUM_REASURED`, `CLAIM` — **seluruhnya kolom uang**.
12. Sebagai **inputor**, saya ingin hasil unggahan dapat saya tinjau sebelum disimpan permanen,
    supaya saya tidak menyimpan data salah. `[terverifikasi]` staging
    `INSERT INTO POOLDATA.M_TEMPUPLOADLIFE` (`PremiumList Life/RDBList/InsertDataUploadLife.xml`);
    harness peninjau `ViewCSVResult_LifePremium` dan varian `_QP` / `_TP` / `_TR`.
13. Sebagai **inputor**, saya ingin memutuskan Confirm / Reject / Decline lagi setelah premium list
    terisi. `[terverifikasi]` `Decision3` memakai `IsLifeAccepted`.

### Penomoran dan periode

14. Sebagai **organisasi**, saya ingin nomor PL dibuat otomatis, supaya tidak bentrok.
    `[terverifikasi]` `PremiumList Life/RDBList/GetSequenceNumber_SQL.xml` →
    `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` (**ADR-0006**).
15. Sebagai **organisasi**, saya ingin prefix nomor **di-lookup**, bukan ditanam.
    `[terverifikasi]` `PremiumList Life/RDBList/GetKodeProdLife_SQL.xml` → `POOLDATA.KODE_PRODUKSI`.
16. Sebagai **Finance**, saya ingin transaksi setelah tanggal tutup buku dibukukan ke **periode
    bulan berikutnya**, supaya pembukuan bulan berjalan tidak tercemar. `[terverifikasi]`
    `PremiumList Life/Activity/SubmitPremiumList_Act.xml`
    (`ASM-FW-GISFW-WORK-LIFE` / `SUBMITPREMIUMLIST_ACT` / `RULE-OBJ-ACTIVITY`).
17. Sebagai **Finance**, saya ingin **tanggal tutup buku dibaca dari tabel**, supaya dapat diubah
    tanpa rilis kode. `[terverifikasi]` `GETTanggalClosing_SQL` → `POOLDATA.TANGGAL_CLOSING`;
    `25` hanya **fallback** bila tabel kosong.

### Penyimpanan polis

18. Sebagai **organisasi**, saya ingin polis yang dikonfirmasi tersimpan lengkap, supaya menjadi
    dasar klaim kelak. `[terverifikasi]` `PremiumList Life/Activity/InsertJsonPolisLife_Act.xml`
    (`ASM-FW-GISFW-WORK-LIFE` / `INSERTJSONPOLISLIFE_ACT` / `RULE-OBJ-ACTIVITY`, **17 langkah**).
19. Sebagai **organisasi**, saya ingin ringkasan premium tersimpan dengan seluruh komponen uangnya.
    `[data DBA]` `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY` — **37 kolom**.
20. Sebagai **Finance**, saya ingin komponen uang tersimpan terpisah — premi, komisi, brokerage,
    pajak, klaim, saldo, biaya admin, potongan — beserta varian **refund** dan **retro**, supaya
    rekonsiliasi mungkin. `[data DBA]`
21. Sebagai **Finance**, saya ingin **tidak ada nilai uang yang berubah** saat menyeberang batas,
    supaya angka dapat dipertanggungjawabkan. `[data DBA]` parameter procedure seluruhnya
    `VARCHAR2` termasuk uang → konversi eksplisit (**ADR-0003**).
22. ~~Sebagai **organisasi**, saya ingin data polis tersimpan sebagai JSON yang dapat di-upsert.~~
    ⚠️ **DIGANTI 2026-09-16** → *Sebagai **organisasi**, saya ingin setiap atribut polis menjadi
    **kolom bernama** di tabel relasional, supaya bentuk polis dapat diperiksa, dicari, dan
    divalidasi — bukan tersembunyi di dalam satu dokumen.* `[keputusan work owner]` (§12)
23. Sebagai **organisasi**, saya ingin data penawaran tersimpan dengan **delapan tanggal** siklusnya,
    supaya lini masa penawaran terbaca. `[data DBA]` `INSERTJSONOFFERLIFE` — `RECEIVED`, `OFFERING`,
    `RESPONSE`, `CONFIRMATION`, `RECONFIRMATION`, `REALIZATION`, `BINDING_DATE`, `MAX_TBC`.
24. Sebagai **organisasi**, saya ingin jenis treaty terbaca dari data, bukan konstanta.
    `[terverifikasi]` `TypeCeding`: `1` = QS, `2` = SURPLUS, `3` = QS + SURPLUS, `4` = XOL.

#### Penyimpanan relasional ⚠️ BARU 2026-09-16 — §12

24a. Sebagai **organisasi**, saya ingin satu polis tersimpan **utuh atau tidak sama sekali**,
     supaya tidak pernah ada polis yang tersimpan separuh. `[keputusan work owner]`
24b. Sebagai **Finance**, saya ingin **rekap uang per mata uang tersimpan**, supaya angka ringkasan
     polis dapat dibaca kembali tanpa menghitung ulang seluruh peserta. `[keputusan work owner]`
24c. Sebagai **underwriter**, saya ingin **hasil spreading dibekukan** pada saat polis dibuat,
     supaya perubahan master treaty kemudian tidak mengubah angka polis yang sudah berjalan.
     `[keputusan work owner]`
24d. Sebagai **tim operasi**, saya ingin **keadaan tangga proses** tersimpan di tempatnya sendiri
     yang dipakai bersama seluruh lini, supaya polis Life dan Non-Life dapat dipantau dengan satu
     cara. `[keputusan work owner]`
24e. Sebagai **tim operasi**, saya ingin sistem **tetap cepat pada jutaan baris peserta**, sehingga
     setiap kunci asing ber-index dan penghapusan besar punya jalur yang aman. `[keputusan work
     owner]`
24f. Sebagai **tim migrasi**, saya ingin polis lama pindah dari JSON ke tabel **tanpa kehilangan
     satu nilai pun**.

### Komponen bersama dengan konteks lain

25. Sebagai **organisasi**, saya ingin rekam summary menyediakan **dua kolom nomor terpisah**
    (`PL_NUMBER` dan `PL_NUMBER_EDM`), supaya jalur new business dan jalur endorsement tidak saling
    timpa. Konteks ini mengisi **yang pertama**. `[terverifikasi]` keduanya adalah **dua parameter
    terpisah** pada `PEGA_M_LIFE_PREMIUM_SUMMARY`.
26. Sebagai **organisasi**, saya ingin mesin lapisan bawah dipakai **bersama, bukan digandakan**,
    supaya perilaku kedua jalur tidak menyimpang diam-diam. `[terverifikasi]` `SaveMasterLPDet`
    (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET`) adalah **satu rule yang sama**
    di kedua modul.

> `[keputusan work owner]` **Endorsement Life adalah konteks terpisah** — alur "pilih polis new
> business yang sudah ada lalu muat data lamanya untuk di-endorse" tidak ada di konteks ini.
> Lihat `.scratch/endorsement-life/`.

### Efek keluar dan lingkungan

27. Sebagai **organisasi**, saya ingin polis yang tersimpan diteruskan ke Arasapas, supaya sistem
    hilir menerimanya. `[terverifikasi]` `PremiumList Life/Activity/serviceInsertArasapasLife_act.xml`.
28. Sebagai **operator**, saya ingin **diberi tahu bila penyimpanan gagal**, supaya kegagalan tidak
    tersembunyi. `[terverifikasi]` step 14 `SendEmailNotification` bergerbang
    `OutDataLife.pxResults(1).PL_NUMBER==""` — **jalur alarm**, bukan notifikasi bisnis.
29. Sebagai **operator**, saya ingin alamat layanan keluar dibaca dari tabel saat dipanggil, supaya
    pemisahan dev–prod ditentukan isi database. **ADR-0013**.
30. Sebagai **operator**, saya ingin efek keluar **tidak berjalan di lingkungan non-production**,
    supaya pengujian tidak menyentuh sistem nyata. `[terverifikasi]` keduanya bergerbang
    `IsPEGAPROD` → flag lingkungan (**ADR-0005**).

### Kontrak hilir

31. Sebagai **Claim Life**, saya ingin menerima `PremiumListSummary` dan `PremiumListDetail` yang
    lengkap, supaya klaim dapat diproses. `[terverifikasi]` class
    `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` dipakai bersama.
32. Sebagai **Claim Life**, saya ingin baris premium dapat ditemukan lewat `PL_NUMBER` **atau**
    `PL_NUMBER_EDM`, supaya klaim atas polis yang di-endorse tetap terhubung. `[terverifikasi]`

### Yang sengaja tidak dibawa

33. Sebagai **tim migrasi**, saya ingin empat gerbang treaty ID polis lama **tidak** ikut pindah.
    `[keputusan work owner]` (**OQ-031**, **OQ-066**).
34. Sebagai **tim migrasi**, saya ingin `ConvertJsonNusareToProduction` **tidak** dihidupkan kembali.
    `[terverifikasi]` step 17 ter-remark; `[keputusan work owner]` dibuang. ⛔ *Ralat 28-09-2026:*
    step 17 (`//` b5383) memanggil `InsertLifePremiumDetail` (b5422), **bukan** layanan ini;
    `convertJsonNusareToProduction` **hidup** di `serviceInsertArasapasLife_act` langkah 5 (b963) —
    premis keputusannya dikonfirmasi ulang, **OQ-PL-14**.
35. Sebagai **tim migrasi**, saya ingin commit implisit Pega **tidak** ditiru. `[terverifikasi]`
    step 16 `Commit` ter-remark.
36. Sebagai **tim migrasi**, saya ingin ambang `25` **tidak** ditanam di kode. `[keputusan work
    owner]` — ikuti versi yang membaca `TANGGAL_CLOSING`.
37. Sebagai **tim migrasi**, saya ingin shape "Input Premium List Summary" **tidak** dibangun kecuali
    terbukti terpakai. `[terverifikasi]` nol connector masuk, nol rujukan harness (**OQ-023**).

---

## Implementation Decisions

### 1. Batas konteks

`[keputusan work owner]` **Konteks ini = PremiumList Life saja (new business).** Endorsement Life
adalah **konteks/menu terpisah**, dengan spec dan tiketnya sendiri di `.scratch/endorsement-life/`.

Alasannya `[terverifikasi]`:

| Pemisah | PremiumList Life | Endorsement Life |
| --- | --- | --- |
| Class work | `ASM-FW-GISFW-WORK-LIFE` | `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE` |
| Flow masuk | `PremiumList Life/InputPolicyHolder.xml` (`INPUTPOLICYHOLDER` / `RULE-OBJ-FLOW`) | `Endorsement Life/Flow/InputEDMLife.xml` (`INPUTEDMLIFE` / `RULE-OBJ-FLOW`) |
| Alur inti | penawaran baru → premium list → simpan polis | **pilih polis NB yang sudah ada → muat data lamanya → endorse** |

Alur inti endorsement itu **tidak ada** di konteks ini.

**Yang dipakai bersama hanyalah mesin lapisan bawah** — dan itu **komponen**, bukan alasan
menyatukan menu:

| Mesin bersama | Identitas |
| --- | --- |
| Penulis detail | `SaveMasterLPDet` — `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL` (**satu rule identik** di kedua modul) |
| Penulis summary | `InsertPLSummary` — `ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY` / `RULE-CONNECT-SQL` (**satu rule identik**) |
| Tabel | `POOLDATA.M_LIFE_PREMIUM_DETAIL`, `M_LIFE_PREMIUM_SUMMARY`, `JSON_POLIS` |
| Penomoran | `PROC_GENERATE_SEQUENCE_NUMBER`, `KODE_PRODUKSI`, `TANGGAL_CLOSING` |
| Orkestrator simpan | `InsertJsonPolisLife_Act` — **nama sama, class berbeda** → **rule berbeda** |

Claim Life dan Komite Claim Life adalah **konteks hilir** — spesifikasi ini berhenti pada penulisan
`M_LIFE_PREMIUM_SUMMARY` / `M_LIFE_PREMIUM_DETAIL`.

### 2. Penempatan modul

`[terverifikasi]` Struktur mengikat `CLAUDE.md` §5: kode di-scaffold **di dalam
`OUTPUT_HASIL_RNM\`**, arah dependency `handlers` → `services` → `repository`.

| Lapisan | Tanggung jawab |
| --- | --- |
| `handlers` | Endpoint REST penawaran, premium list, unggah CSV |
| `services` | Mesin alur; periode tutup buku; penomoran; konversi uang teks ↔ desimal; orkestrasi efek keluar |
| `repository` | Pemanggilan empat procedure; staging unggah; lookup `M_LINK_SERVICE` |
| `models` | Penawaran, premium list summary + detail |
| `frontend` | Layar Input Offer, Input Premium List Detail, unggah + tinjau CSV |

### 3. Mesin alur

`[terverifikasi]` Graf `InputPolicyHolder` — 3 Assignment (**seluruhnya `WorkList`**, OQ-028),
3 Decision, 1 Utility, 11 connector:

```
Input Offer → IsLifeAccepted ─ Confirm ─→ IsFlagOnGoingPolicy ─ Offer   ─→ selesai (penawaran)
                             ├ Reject  ─→ balik ke Input Offer          │
                             └ Decline ─→ case ditutup                  └ Premium ─→ Input Premium List Detail
                                                                                       ↓
                                                                               IsLifeAccepted
                                                                        ├ Confirm ─→ simpan polis
                                                                        ├ Reject  ─→ balik ke Input Offer
                                                                        └ Decline ─→ case ditutup
```

**Aturan mengikat:**

- Ketiga keputusan **diambil manusia** — bukan formula. Yang direplikasi adalah **akibatnya**.
- **`Reject` selalu kembali ke input**; **`Decline` menutup case**; **`Confirm` melanjutkan**.
- `IsFlagOnGoingPolicy`: **`1` = Offer** (berhenti di penawaran), **`2` = Premium** (lanjut).

### 4. Periode tutup buku

**Dibaca dari `POOLDATA.TANGGAL_CLOSING`**, bukan ditanam. Bila tanggal berjalan melewati tanggal
closing, `ProdDateTime` digeser ke **tanggal 1 bulan berikutnya, `05:00 GMT` (12:00 WIB)**.

`[terverifikasi]` `25` muncul **hanya sebagai fallback** bila tabel kosong
(`@if(Local.TglProd=="",25,Local.TglProd)`).

⚠️ `[keputusan work owner]` **Fallback itu TIDAK direplikasi.** Bila `POOLDATA.TANGGAL_CLOSING`
kosong atau tidak terbaca, sistem baru **gagal terang-terangan** — transaksi ditolak dengan pesan
yang menyebut tabel sumbernya. Alasan: fallback diam membukukan transaksi ke periode yang salah
tanpa meninggalkan jejak, dan kesalahan periode baru terlihat saat tutup buku. Lihat **AC 8**.

⚠️ `[terverifikasi]` **Dua versi hidup berdampingan** — `SubmitPremiumList_Act` (`20260122`) membaca
dari tabel; `InsertJsonPolisLife_Act` (`20260728`, **lebih baru**) masih menanam `>25`.
`[keputusan work owner]` **ikuti yang dari DB.**

> ⛔ *Ralat 28-09-2026 (sensus remark GILIRAN-12):* `SubmitPremiumList_Act` langkah 2–4 membaca
> `TANGGAL_CLOSING`, tetapi hasilnya hanya mengalir ke `TempGenerate.CARI1/2` (tak dibaca rule mana
> pun) dan ke langkah 15 yang **ter-remark** (`//` b3398; b3491 miliknya). Dua aturan yang **hidup**:
> nomor PL digulir `PROC_GENERATE_SEQUENCE_NUMBER` (membaca tabel itu sendiri), dan `ProdDateTime`
> digeser `InsertJsonPolisLife_Act` langkah 4 (`>25` tertanam, b1170). Keputusan "ikuti yang dari DB"
> tetap; penerapannya pada `ProdDateTime` (yang di Pega memakai 25) dikonfirmasi ulang — **OQ-PL-13**.
>
> ✅ **OQ-PL-13 DITUTUP 29-09-2026 (GILIRAN-17)** `[keputusan work owner 29-09-2026 — lembar keputusan, "rekomendasi"]`: `ProdDateTime` **ikut XML** — ambang
> **25 tertanam** (`InsertJsonPolisLife_Act` langkah 4, gerbang b1170 `@toDecimal(Local.currentdate)>25`, nilai b1092). Keputusan
> "ikuti yang dari DB" tetap berlaku untuk periode yang ditampilkan dan untuk penomoran PL, karena pembaca hidup `TANGGAL_CLOSING`
> adalah `PROC_GENERATE_SEQUENCE_NUMBER`.

### 5. Uang — teks di batas, desimal di dalam

`[data DBA]` Seluruh parameter `PEGA_M_LIFE_PREMIUM_SUMMARY` bertipe **`VARCHAR2`, termasuk kolom
uang**. **ADR-0003** karena itu diterapkan dengan bentuk khusus di sini:

- Di dalam sistem, uang adalah **desimal presisi arbitrer**.
- **Konversi teks ↔ desimal terjadi di satu batas** — lapisan repository — dan tidak di tempat lain.
- **Tidak ada tahap yang melewatkan nilai lewat `float`.**

Kelompok kolom uang: gross (`PREMIUM`, `COMMISSION`, `BROKERAGE_FEE`, `OVR_COMM`, `TAX`,
`PROF_COMM`, `CLAIM`, `CLAIM_AMOUNT`, `BALANCE`, `RI_ADMIN_FEE`, `DEDUCTION`), ditambah varian
`*_REFUND`, `*_RETRO`, dan `*_REFUND_RETRO`.

### 6. Batas transaksi — ⚠️ **DIREVISI 2026-09-16: titik potong lenyap**

`[data DBA]` Keadaan lama — dan mengapa ia rumit:

| Procedure | Commit sendiri? |
| --- | --- |
| `PEGA_M_LIFE_PREMIUM_SUMMARY` | **tidak** |
| `PROC_GENERATE_SEQUENCE_NUMBER` | **tidak** |
| ~~`INSERTJSONPOLISLIFE`~~ | ~~**ya**~~ — ⚠️ **dibuang** (§12) |
| ~~`INSERTJSONOFFERLIFE`~~ | ~~**ya**~~ — ⚠️ **dibuang** (§12) |
| `SaveMasterLPDet` (detail peserta) | **ya** — `COMMIT;` |

⚠️ **Kedua procedure yang memaksa titik potong adalah procedure JSON — dan keduanya dibuang.**
Karena polis kini disimpan ke **tujuh tabel relasional yang ditulis Go sendiri** (§12), tidak ada
lagi `COMMIT` di luar kendali Go pada jalur simpan polis.

**Aturan mengikat yang berlaku sekarang** `[keputusan work owner]`:

> **Satu polis — header, rekap mata uang, seluruh peserta, seluruh spreading, seluruh spreading
> retro, dan seluruh riwayat penawaran — ditulis dalam SATU transaksi, lalu commit sekali.**

`PROC_GENERATE_SEQUENCE_NUMBER` ikut di dalam transaksi itu, sehingga nomor dan polis lahir
**atomik**: tidak pernah ada nomor tanpa polis, tidak pernah ada polis tanpa nomor.

⚠️ **Satu titik potong tersisa** `[data DBA]`: `SaveMasterLPDet` **commit sendiri**. Ia adalah jalur
tulis ke tabel **existing** `M_LIFE_PREMIUM_DETAIL`, yang di sistem baru **tidak lagi ditulis**
(§12) — jadi titik potongnya ikut hilang. Bila jalur itu masih dipakai selama masa transisi,
ia **wajib diletakkan di luar** transaksi polis dan dapat diulang.

`[terverifikasi]` Commit eksplisit Pega (step 16) ter-remark; jangan meniru commit implisit.
Konsisten **OQ-013**.

### 7. Penomoran

`PL_NUMBER` (NB) dan `PL_NUMBER_EDM` (endorsement) adalah **dua jalur sejajar** — dibuat terpisah
agar tidak saling timpa, **bukan** relasi induk-anak. Keduanya menyeberang sebagai **parameter
terpisah**.

Prefix di-lookup dari `KODE_PRODUKSI`; sequence dari `PROC_GENERATE_SEQUENCE_NUMBER` (**ADR-0006**).

⚠️ **Itu berlaku untuk `PL_NUMBER` saja.** `[terverifikasi]` `PL_NUMBER_EDM` **tidak** lahir dari
sequence terpusat — ia dirakit dari data polis oleh `Generate_NoEndorsmentLife`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!GENERATE_NOENDORSMENTLIFE` / `RULE-CONNECT-SQL`).
**ADR-0006 tidak berlaku pada nomor EDM.** Rincian, AC, dan tiketnya **milik konteks Endorsement
Life** — lihat `.scratch/endorsement-life/`.

Konteks ini hanya berkewajiban: **menyediakan kolom `PL_NUMBER_EDM` pada rekam summary** dan
**membiarkannya kosong** untuk jalur new business.

### 8. Unggah CSV

Alur: unggah → **staging** `POOLDATA.M_TEMPUPLOADLIFE` → validasi → tinjau → simpan permanen.
Validasi menyentuh **kolom uang** (§User Stories 11), sehingga ia juga titik penegakan **ADR-0003**.

### 9. Efek keluar — dua, di belakang flag lingkungan

| Efek | Sifat |
| --- | --- |
| `serviceInsertArasapasLife_act` | **efek bisnis**; endpoint via `M_LINK_SERVICE` (**ADR-0013**) |
| `SendEmailNotification` | **jalur alarm** — dipicu bila `PL_NUMBER` kosong setelah penyimpanan |

⚠️ **Penyimpangan sadar:** email diperlakukan sebagai **alarm operasional**, bukan notifikasi
bisnis. Ia tidak boleh muncul di UI sebagai "pemberitahuan ke pengguna".

`ConvertJsonNusareToProduction` **dibuang**.

### 10. Kode mati yang tidak dimigrasikan

| Yang dibuang | Bukti |
| --- | --- |
| Empat gerbang `@contains(.ID,"1000032")`…`"1000035"` (class `ASM-FW-GISFW-Int-TREATYYEAR_LIFE`) | `[keputusan work owner]` logika polis lama. ⚠️ di korpus **masih aktif** — **OQ-066** |
| Step 16 `Commit` | `[terverifikasi]` `blockname //` |
| Step 17 `Connect-REST` ~~`ConvertJsonNusareToProduction`~~ `InsertLifePremiumDetail` *(diralat 28-09-2026, b5422)* | `[terverifikasi]` `blockname //` |
| Ambang `25` ter-hardcode | `[keputusan work owner]` |
| Shape "Input Premium List Summary" | `[terverifikasi]` nol connector masuk, nol rujukan harness — **OQ-023** |

⚠️ **Pelajaran metodologi (OQ-066):** penanda `<pyStepsBlockName>` **tidak dapat dipercaya
sendirian** — di berkas ini ia benar untuk dua langkah dan salah untuk empat. **Status hidup/mati
menuntut konfirmasi work owner**, bukan hanya pembacaan tag. Berlaku lintas modul.


### 11. Penulisan detail peserta — **inline, bukan job** `[keputusan work owner]`

Baris peserta di `M_LIFE_PREMIUM_DETAIL` adalah yang **dibaca Claim Life** (`GetPesertaClaim_sql1`,
`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `RNM!GETPESERTACLAIM_SQL1`, berkunci `PL_NUMBER`).

`[terverifikasi]` Kedua jalur menulisnya lewat **satu rule yang sama** — `SaveMasterLPDet`
(`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET` / `RULE-CONNECT-SQL`); ekspor di
`PremiumList Life/RDBList/` dan `Endorsement Life/RDBList/` identik kecuali cap waktu ekspor.
**Yang berbeda hanyalah pemicunya.**

| Jalur | Pemicu di Pega | Pemicu di sistem baru |
| --- | --- | --- |
| **NB** | **job/batch** — `InsertLifePremiumDetail_act` (`ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT`) step 2 `pxRetrieveReportData` atas `SelectNoJsonPolis_RD`, lalu loop step 3 | **INLINE saat simpan polis** |
| **EDM** | **inline** — `Endorsement Life/Activity/InsertJsonPolisLife_Act.xml` step 11.6 | **INLINE** (tidak berubah) |

⚠️ **Penyimpangan sadar.** **Logikanya tidak diganti** — tetap `InsertLifePremiumDetail_act`. Yang
disamakan dengan EDM adalah **waktu pemicunya**. Step 2 dan loop step 3 tidak dimigrasikan (semata
mesin batch); yang dimigrasikan adalah **badan per-kasus** (3.1–3.6).

**Konsekuensi positif:** data peserta langsung tersedia untuk klaim **tanpa jeda job**; NB dan EDM
seragam.

`[terverifikasi]` Penjaga idempotensi sudah ada di korpus — langkah insert (3.3.4) dijaga
`hasilDetail.pxResults(1).PL_NUMBER==""` (baris 3820) — dan **wajib dipertahankan**.
⚠️ `SaveMasterLPDet` **commit sendiri** (`COMMIT;` baris 252) → **titik potong**, tunduk pada §6.

---

### 12. Bentuk penyimpanan — **tujuh tabel relasional + `T_WORK_POLIS`** ⚠️ BARU 2026-09-16

⚠️ **Penyimpangan sadar 1 — seluruh JSON dibuang.** `[keputusan work owner]` `JSON_POLIS.DATA_JSON`
dan `JSON_OFFER_LIFE` (keduanya CLOB) **tidak dipakai**. Perakitnya — `@ASM.GetPageJSONString()`,
`[terverifikasi]` dipanggil **211×** di seluruh korpus, **selalu tanpa argumen**, dan definisinya
(Rule-Utility-Function pustaka `ASM`) **tidak diekspor** — **tidak direplikasi**.

`[keputusan work owner]` Polis punya **tabelnya sendiri**. `M_LIFE_PREMIUM_SUMMARY`,
`M_LIFE_PREMIUM_DETAIL`, `LIFEINPRODUCTION`, `JSON_POLIS`, dan `JSON_OFFER_LIFE` **hanya dibaca saat
migrasi**.

```
T_PREMIUM_LIST (PK ID)                        header polis
  ├─ T_PREMIUM_LIST_SUMMARY        1:N  FK PREMIUM_LIST_ID   ← rekap per mata uang
  ├─ T_PREMIUM_LIST_DETAIL         1:N  FK PREMIUM_LIST_ID   ← per peserta
  │     └─ T_PREMIUM_LIST_SPREADING       1:N  FK DETAIL_ID
  │            └─ T_PREMIUM_LIST_SPREADING_RETRO  1:N  FK SPREADING_ID
  └─ T_VIEW_SUGGEST                1:N  FK PREMIUM_LIST_ID   ← riwayat penawaran
T_WORK_POLIS (PK ID)   ⬅ tabel work mandiri, LINTAS-LINI — bukan anak T_PREMIUM_LIST
```

Seluruh FK **`ON DELETE CASCADE`** + popup konfirmasi. **Empat tingkat** kedalaman.

`[terverifikasi dari contoh `DATA_JSON` nyata]` Pohon ini **sama persis** dengan bentuk payload
produksi: `PremiumListSummary → {CurrencyList[], PremiumListDetail[] → SpreadingList[] →
RetroLifeList[]}`, `OfferFacIn → ViewSuggest[]`.

#### Isi tiap tabel

**`T_PREMIUM_LIST`** — field bisnis yang di-set aplikasi ke work object: `PL_NUMBER`, `RISLIPRNM`
(opsional), identitas pihak (`POLICY_HOLDER`(+`_NAME`), `CEDING_CO`(+`_NAME`), `MARKETING_CODE`,
`MARKETING_NAME`, `MO_ID`, `BRANCH_CODE`, `BRANCH_NAME`), bisnis (`BUSINESS_CODE`, `BUSINESS_NAME`,
`SOURCE_OF_BUSINESS`, `SOB_NAME`, `JENIS_ASURANSI`, `TYPE`, `TYPE_CEDING`, `PRO_RATE_TYPE`,
`DESCRIPTION`, `WPC`), produk (`PRODUCT_NAME`, `PRODUCT_NAME_ID`), penawaran (`NO_OFFER`),
retro (`RETRO_ID`, `RETRO_NAME`, `SECURITY_REINSURER`(+`_ID`) — **opsional**, kosong pada NB murni),
angka (`ANNUITY_INTEREST`, `PREMIUM_REFUND_FACTOR`), tanggal (`DATE_RECEIVED`), dan flag proses
(`POSITION`, `PROPOSAL_ACCEPT_STATUS`, `FLAG_ON_GOING_POLICY`, `EMAIL_TYPE_PL`, `EDIT_INPUT`,
`IS_JSON_POLIS`).

⚠️ **Properti bawaan Pega dibuang** — seluruh `px*`, `py*`, `pz*` (`pxApplication`, `pyID`,
`pyStatusWork`, `pyOrigDivision`, …) beserta single-page kosong (`Policy`, `Quotation`,
`TempError`). Identitas work Pega **diganti** identitas sequence.

⚠️ `IS_JSON_POLIS`, `EMAIL_TYPE_PL`, `EDIT_INPUT` adalah **flag proses bernama menyesatkan** —
disimpan apa adanya; artinya dikonfirmasi bila perlu.

**`T_PREMIUM_LIST_SUMMARY`** — rekap **per mata uang**: `CURRENCY` + ±35 kolom uang dalam empat
kelompok: dasar (`PREMIUM`, `COMMISSION`, `BROKERAGE_FEE`, `OVR_COMM`, `TAX`, `PROF_COMM`, `CLAIM`,
`CLAIM_AMOUNT`, `BALANCE`, `DEDUCTION`, `RI_ADMIN_FEE`, `CEDING_RETENTION`, `SUM_REASURED`,
`SUM_AT_RISK_GROSS`, `SHARE_RETRO`, `SHARE_NUSANTARA_RE_GROSS`), `*_REFUND`, `*_RETRO`, dan
`*_REFUND_RETRO`. Seluruhnya **NUMBER**.

⚠️ **Penyimpangan sadar 2 — rekap uang per mata uang DISIMPAN.** `[keputusan work owner]`
Sebelumnya diduga rekap itu hanya hidup di `Param.*` dan tidak tersimpan; **contoh `DATA_JSON`
nyata membuktikan sebaliknya**. Tabel ini menyimpan rekapnya penuh, bukan sekadar daftar mata uang.

**`T_PREMIUM_LIST_DETAIL`** — satu baris per peserta, **±81 kolom** (acuan: daftar `INSERT`
`SaveMasterLPDet`, `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET`).

⚠️ **Dua jendela valuasi, keduanya wajib ada** `[terverifikasi dari dua contoh `DATA_JSON`]`:
jalur `Q*` (NB) memakai **`GROSS_VALUATION_BEGIN/EXPIRED_DATE`** beserta `EFFECTIVE_DATE`,
`LAPSE_DATE`, `PERIOD_MM`; jalur `T*` (retro/EDM) memakai **`RETROCESSION_VALUATION_*`**. Pola ini
**konsisten dengan validasi Date of Loss di Claim Life** — bukan kebetulan.

⚠️ `FACTOR` bertipe **NUMBER** dengan tujuh desimal (contoh: `1.0000000`), **bukan** bilangan bulat.

Kolom EDM — `PL_NUMBER_EDM`, `EDMSTATUS`, `STATUSOLD`, `STATUS` — ikut ada; keduanya baru terisi
pada jalur endorsement/old-data.

**`T_PREMIUM_LIST_SPREADING`** — hasil spreading per peserta, per treaty-year/jenis:
`TREATY_TYPE_ID`, `TREATY_TYPE_NAME` (QS / 2ND QS / SURPLUS / 2ND SURPLUS / OR), `TREATY_YEAR_LIFE`,
`RETROCADED_SHARE`, `IDR`, `USD`, `IDR_SELISIH`, `USD_SELISIH`, `B_IDR`, `B_USD`, `TGL_UPDATE`,
`USER_ID`.

**`T_PREMIUM_LIST_SPREADING_RETRO`** — per reinsurer: `REINSURER_NAME`, `PERCENT_SHARE`, `AMOUNT`,
`COMMISION` (sic), `OVR_COMM`, `RATE`, `PREMIUM_SPREADED_GROSS`, `PREMIUM_SPREADED_NET`,
`TREATY_TYPE_ID`, `TREATY_TYPE_NAME`, `TREATY_START_DATE`, `TREATY_END_DATE`, `TGL_UPDATE`,
`USER_ID`.

⚠️ **Penyimpangan sadar 3 — spreading DIBEKUKAN dan DISIMPAN.** `[keputusan work owner]` Pega tidak
menyimpannya di PremiumList; ia dihitung di konteks Treaty/Master Retro. Menyimpannya adalah
penyimpangan sadar **demi jejak historis** — angka spreading pada saat polis dibuat tidak berubah
ketika master treaty berubah kemudian.

**`T_VIEW_SUGGEST`** — riwayat penawaran/konfirmasi ceding: `NO`, `DATE_SUGGEST`, `PIC_SUGGEST`,
`IS_CEDING_CONFIRM` (`Accept`/`Reject`/`Decline`), `COMMENT_SUGGEST`, `INITIAL_SUGGEST` (`Offer`/`Bind`).

⛔ **RALAT 28 September 2026 (pl6).** Kolom terakhir semula bernama `INITIAL` — **kata
cadangan Oracle**. Migrasi `056` karena itu **gagal** saat work owner menjalankan
`-migrate` pukul 09.33, dan `T_VIEW_SUGGEST` tidak terbentuk di DEV sementara sebelas
langkah lain di jalan yang sama lolos. Bukti yang dijalankan work owner:

```
SELECT 1 AS INITIAL   FROM DUAL  -> ORA-00923
SELECT 1 AS "INITIAL" FROM DUAL  -> lolos (berkutip)
SELECT 1 AS NO        FROM DUAL  -> lolos
```

Berkas `056` **disunting di tempat**, bukan ditambah migrasi `057`: ia belum pernah
terpasang di mana pun *(`T_MIGRASI` DEV tanpa `056`, skema uji belum dibuat)*, dan
`ALTER RENAME` untuk kolom yang belum pernah ada berarti mewariskan riwayat yang tidak
terjadi.

⚠️ `[keputusan work owner]` `OfferFacIn` **dibuang sebagai simpul**; `ViewSuggest` **naik** menjadi
anak langsung header.

**`T_WORK_POLIS`** — ⚠️ **Penyimpangan sadar 4.** `[keputusan work owner]` Keadaan tangga proses
(`POSITION`/`STATUS` dari nilai connector `Confirm`/`Decline`/`Reject`/`Offer`/`Premium`, identitas
polis + lini, audit) berada di tabel **lintas-lini tersendiri** — **bukan** kolom `T_PREMIUM_LIST`,
dan **bukan** anaknya. Pola sama `T_WORK_CLAIM` (Claim Life). Menghapus polis menghapus baris
work-nya.

#### Tipe, kunci, dan skala

⚠️ **Penyimpangan sadar 5.** Uang/share/premi → **desimal presisi arbitrer** (**ADR-0003**, tidak
pernah `float`); tanggal → **`DATE`**; **seluruh kolom nullable** (wajib-isi di Go); identitas via
**sequence** (**ADR-0006**).

⚠️ **Kebutuhan non-fungsional — skala jutaan baris** `[keputusan work owner]`.
`T_PREMIUM_LIST_DETAIL`, `_SPREADING`, dan `_SPREADING_RETRO` tumbuh sangat besar. **Index pada
setiap FK wajib.** Partisi per periode/tahun dipertimbangkan. Dan karena hapus fisik satu polis
besar berarti kaskade jutaan baris, **arsip/soft-delete dipertimbangkan sebagai jalur normal** —
bukan hapus fisik.

---

## Testing Decisions

### Apa yang membuat test baik di sini

Test menguji **perilaku yang teramati dari luar**: apa yang terjadi pada penawaran, pada premium
list, dan pada rekam yang dihasilkan — bukan bahwa suatu procedure dipanggil. Test yang hanya
membenarkan implementasi tidak diterima.

### Seam — **memakai ulang seam yang sudah ada**

`[terverifikasi]` Repo target belum di-scaffold. Seam yang ditetapkan spec Claim — Life adalah
**API HTTP**, dan konteks ini memakainya kembali — **tidak menambah seam**:

> **Seam utama: API HTTP Life — Penawaran & Premium List.** Test menggerakkan alur lewat endpoint
> REST dan memeriksa hasilnya lewat endpoint REST, dengan `handlers → services → repository`
> terpasang sungguhan, terhadap skema uji Oracle.

**Batas proses difake:**

| Batas | Perlakuan |
| --- | --- |
| Arasapas, email | *fake* di balik interface; test memeriksa **efeknya** |
| Oracle | **skema uji nyata, bukan mock** — empat procedure adalah inti perilaku, dan sifat commit-nya tidak dapat difake dengan jujur |
| Jam | dapat dikendalikan — aturan tutup buku bergantung tanggal berjalan |

**Tidak ada seam kedua.** Berbeda dari Komite Claim Life, konteks ini tidak punya worker asinkron.

### Modul yang diuji

| Yang diuji | Lewat seam |
| --- | --- |
| Alur Confirm / Reject / Decline dan Offer / Premium | API HTTP |
| Periode tutup buku dari tabel, termasuk kegagalan terang-terangan saat tabel kosong | API HTTP + jam terkendali |
| Penomoran `PL_NUMBER` — lahir sekali, tidak bertabrakan | API HTTP |
| Unggah CSV: staging, validasi kolom uang, tinjau, simpan | API HTTP |
| Uang tidak berubah menyeberang batas teks ↔ desimal | API HTTP |
| Batas transaksi campuran — apa yang bertahan saat gagal di tengah | API HTTP + skema uji |
| Efek keluar di belakang flag lingkungan | API HTTP + fake |
| Detail peserta tertulis **inline** saat simpan — terbaca segera, tanpa job | API HTTP + skema uji |

### Prior art

`[terverifikasi]` **Tidak ada** — nol kode, nol test di `OUTPUT_HASIL_RNM`. Spec Claim — Life dan
Komite Claim Life menetapkan bentuknya; konteks ini mengikuti bentuk yang sama.

Perintah verifikasi wajib ditulis eksplisit di tiap tiket selama `Makefile` belum ada. Target:
`go test ./internal/...` dan `cd frontend && npm test`.

---

## Acceptance Criteria

**Alur**

1. `Confirm` pada tahap penawaran melanjutkan ke penentuan Offer/Premium.
2. `Reject` pada tahap mana pun **mengembalikan** case ke layar Input Offer.
3. `Decline` pada tahap mana pun **menutup** case; case tidak dapat dilanjutkan.
4. Keluaran `Offer` menghentikan siklus di tahap penawaran, dengan penawaran tetap tersimpan.
5. Keluaran `Premium` membuka tahap Input Premium List Detail.
6. Keputusan diambil pengguna; **tidak ada** aturan otomatis yang menetapkannya.

**Periode**

7. Tanggal tutup buku dibaca dari `POOLDATA.TANGGAL_CLOSING` **setiap kali**, bukan dari konstanta.
8. Bila `POOLDATA.TANGGAL_CLOSING` kosong atau tidak terbaca, sistem **gagal terang-terangan**:
   transaksi ditolak dengan pesan yang menyebut tabel sumbernya. Sistem baru **TIDAK** meniru
   fallback diam `25` milik Pega (`@if(Local.TglProd=="",25,Local.TglProd)`) — fallback diam adalah
   kegagalan tersembunyi yang membukukan transaksi ke periode yang salah tanpa jejak.
   `[keputusan work owner]`
9. Transaksi setelah tanggal closing memperoleh periode **tanggal 1 bulan berikutnya**.

**Penomoran**

10. `PL_NUMBER` diperoleh dari `PROC_GENERATE_SEQUENCE_NUMBER`; aplikasi tidak menyusun formatnya
    sendiri (**ADR-0006**).
11. Prefix diperoleh lewat lookup `KODE_PRODUKSI`, tidak ditanam.
12. Rekam summary menyediakan kolom `PL_NUMBER_EDM`; jalur new business **membiarkannya kosong** dan
    **tidak pernah** menyentuh `PL_NUMBER` milik rekam lain. *(pengisian `PL_NUMBER_EDM` = konteks
    Endorsement Life)*
13. Kedua kolom nomor hidup berdampingan pada rekam summary yang sama, sebagai dua nilai terpisah.

**Uang**

14. Tidak ada nilai uang yang melewati `float` di lapisan mana pun maupun di JSON API.
15. Konversi teks ↔ desimal terjadi **hanya di lapisan repository**.
16. Nilai uang yang dikirim ke procedure dan dibaca kembali **identik** — tidak ada pembulatan diam.
17. Validasi unggah CSV menolak nilai uang yang tidak sah dan menyebut **kolom** yang salah.

**Unggah CSV**

18. Berkas yang diunggah masuk **staging** lebih dulu; kegagalan validasi tidak mencemari tabel
    permanen.
19. Pengguna dapat meninjau hasil unggahan sebelum menyimpan permanen.

**Transaksi**

20. `PROC_GENERATE_SEQUENCE_NUMBER` dan `PEGA_M_LIFE_PREMIUM_SUMMARY` dipanggil di dalam **satu
    transaksi Go**, dan transaksi itu **commit sebelum** procedure lain dipanggil. Nomor dan rekam
    summary lahir atomik. `[keputusan desain]`
21. `INSERTJSONPOLISLIFE` dan `INSERTJSONOFFERLIFE` dipanggil **setelah** commit tersebut, dalam
    urutan itu, **di luar** transaksi Go — karena keduanya commit sendiri. `[keputusan desain]`
22. Kegagalan **sebelum** commit summary tidak meninggalkan nomor maupun rekam separuh.
23. Kegagalan **setelah** commit summary meninggalkan nomor + rekam summary **utuh**; JSON polis
    dan/atau offer boleh belum tertulis, keadaan itu **terdeteksi**, dan pemanggilan ulang aman
    karena kedua procedure adalah upsert.
24. Urutan pemanggilan keempat procedure terdokumentasi di kode sebagai bagian kebenaran, bukan
    kebetulan; ada test yang gagal bila urutannya diubah.

**Efek keluar**

25. Di lingkungan non-production, kedua efek keluar **tidak berjalan**; penyimpanan tetap berjalan.
26. Alamat Arasapas di-resolve runtime dari `M_LINK_SERVICE`; tidak ada URL sebagai literal,
    konstanta, maupun env var (**ADR-0013**).
27. Email hanya terkirim bila penyimpanan **gagal** (`PL_NUMBER` kosong) — ia alarm, bukan
    notifikasi bisnis.

**Kontrak hilir**

28. Rekam `M_LIFE_PREMIUM_SUMMARY` dan `M_LIFE_PREMIUM_DETAIL` yang dihasilkan dapat ditemukan Claim
    Life lewat `PL_NUMBER` maupun `PL_NUMBER_EDM`.
28a. Baris detail peserta **new business** ditulis **INLINE saat simpan polis** lewat logika
    `InsertLifePremiumDetail_act` → `SaveMasterLPDet` — **bukan** oleh job. Setelah respons simpan
    berhasil, baris itu **sudah ada** dan terbaca jalur baca klaim, tanpa jeda.
    `[keputusan work owner]`
28b. Penulisan detail **idempoten**: menyimpan ulang polis yang sama tidak menggandakan baris
    (setara precondition `PL_NUMBER==""` pada step 3.3.4).
29. Perubahan bentuk kedua rekam diperlakukan sebagai **perubahan kontrak lintas konteks** dan
    ditandai demikian di kode.

**Kode mati**

30. Tidak ada padanan empat gerbang treaty ID, `Commit` eksplisit, maupun
    `ConvertJsonNusareToProduction`.
31. Jenis treaty diambil dari data `TypeCeding` (`1`=QS, `2`=SURPLUS, `3`=QS+SURPLUS, `4`=XOL),
    bukan dari konstanta ID.

**Penyimpanan relasional** ⚠️ BARU 2026-09-16 — §12

32. ⚠️ **Tidak ada blob JSON** sebagai penyimpan isi polis. Test yang menemukan kolom JSON menyimpan
    atribut polis **gagal**. *(§12; penyimpangan sadar 1)*
33. ⚠️ Polis baru disimpan **hanya** ke ketujuh tabel; `M_LIFE_PREMIUM_SUMMARY`,
    `M_LIFE_PREMIUM_DETAIL`, `LIFEINPRODUCTION`, `JSON_POLIS`, dan `JSON_OFFER_LIFE` **tidak
    ditulis**. *(§12)*
34. ⚠️ **Perakit JSON tidak direplikasi** — tidak ada padanan `@ASM.GetPageJSONString()` di sistem
    baru. *(§12; penyimpangan sadar 1)*
35. ⚠️ Satu polis — header, rekap mata uang, seluruh peserta, seluruh spreading, seluruh spreading
    retro, dan seluruh riwayat penawaran — ditulis dalam **satu transaksi**; kegagalan di tingkat
    mana pun **membatalkan seluruhnya**. *(§6, §12)*
36. Nomor polis dan polisnya lahir **atomik**: tidak pernah ada nomor tanpa polis, tidak pernah ada
    polis tanpa nomor. *(§6, §7)*
37. ⚠️ **Rekap uang per mata uang tersimpan** — bukan hanya daftar mata uang. Satu polis bermata uang
    ganda menghasilkan **satu baris rekap per mata uang**, masing-masing dengan nilai uangnya.
    *(§12; penyimpangan sadar 2)*
38. Peserta menyimpan **kedua jendela valuasi**: `GROSS_VALUATION_*` **dan**
    `RETROCESSION_VALUATION_*`. Polis `Q*` mengisi jendela gross; polis `T*` mengisi jendela
    retrosesi. Test wajib memuat **kedua jalur**. *(§12)*
39. `FACTOR` diperlakukan sebagai **desimal**, bukan bilangan bulat; nilai berdesimal tujuh angka
    melewati simpan-baca **tanpa berubah**. *(§12; **ADR-0003**)*
40. ⚠️ Setiap baris spreading menunjuk **satu peserta**, dan setiap baris spreading retro menunjuk
    **satu baris spreading**. Test yang menemukan keduanya menggantung pada header **gagal**.
    *(§12)*
41. ⚠️ Nilai spreading **dibekukan** saat polis dibuat: perubahan master treaty sesudahnya **tidak
    mengubah** angka spreading polis yang sudah tersimpan. *(§12; penyimpangan sadar 3)*
42. Polis tanpa retrosesi (NB murni) tersimpan dengan **nol baris** spreading dan spreading retro —
    **bukan** kegagalan. *(§12)*
43. `RISLIPRNM`, `RETRO_ID`, `RETRO_NAME`, dan `SECURITY_REINSURER*` **boleh kosong** pada polis
    tanpa retrosesi. *(§12)*
44. ⚠️ Riwayat penawaran melekat **langsung pada header**; tidak ada simpul `OfferFacIn` di skema
    baru. *(§12)*
45. ⚠️ Keadaan tangga proses tersimpan di **`T_WORK_POLIS`**, **bukan** sebagai kolom
    `T_PREMIUM_LIST`. Menghapus polis menghapus baris work-nya. *(§12; penyimpangan sadar 4)*
46. ⚠️ Menghapus polis **mengkaskade** ke rekap mata uang, peserta, spreading, spreading retro, dan
    riwayat penawaran — didahului **popup konfirmasi Ya/Batal** yang menyebut jumlah baris tiap
    jenis. **Batal** tidak mengubah apa pun. *(§12)*
47. ⚠️ **Properti bawaan Pega tidak menjadi kolom** — tidak ada `px*`, `py*`, `pz*` di skema baru.
    *(§12)*
48. ⚠️ Seluruh uang dan share bertipe **desimal presisi arbitrer**; seluruh tanggal **`DATE`**;
    seluruh kolom **nullable**; identitas dari **sequence**. *(§12; **ADR-0003**, **ADR-0006**;
    penyimpangan sadar 5)*
49. ⚠️ **Setiap FK punya index.** Test/migrasi yang meninggalkan FK tanpa index **gagal** — ketiga
    tabel terdalam tumbuh sampai jutaan baris. *(§12)*
50. Migrasi memindahkan polis lama dari bentuk JSON + tabel flat ke ketujuh tabel **tanpa kehilangan
    satu nilai pun**; nilai uang dibandingkan **secara tepat**, bukan dengan toleransi.
    *(§12; **ADR-0003**)*

---

## Pertanyaan terbuka di dalam spec

Enam OQ. **OQ-068 ditutup 2026-09-15** — tidak ada lagi OQ yang memblokir tiket non-migrasi.
Hanya **OQ-001 (sisa)** yang menahan tiket migrasi (09); lima sisanya **non-pemblokir**.

| OQ | Pertanyaan | Pemilik | Bagian yang menunggu |
| --- | --- | --- | --- |
| ~~**OQ-001**~~ (sisa) | ~~DDL fisik `M_LIFE_PREMIUM_SUMMARY`, `JSON_POLIS`, `JSON_OFFER_LIFE`~~ | — | ✅ **DITUTUP 2026-09-16** `[keputusan work owner]` — ketujuh tabel **dirancang sendiri** (§12); tipe, kunci, dan kaskade ditetapkan di sana. Presisi fisik dicocokkan DBA **saat tiket migrasi**, bukan pemblokir. Tabel existing tinggal **dibaca** |
| **OQ-023** | Shape "Input Premium List Summary" tanpa jalur masuk terbaca — mati, atau dicapai lewat mekanisme di luar 17 tipe rule yang diekspor | Arsitektur Pega | §10 kode mati |
| **OQ-066** | Penanda `blockname` tidak dapat dipercaya sendirian — adakah langkah lain di modul Life yang aktif tetapi sebenarnya mati | Arsitektur Pega + Product+UW | §10 kode mati |
| **OQ-028** | Seluruh Assignment memakai `WorkList` (bukan `WorkBasket`) — implikasi routing/RBAC belum ditetapkan | IAM | §2 penempatan; belum menyentuh AC |
| **OQ-067** | `INSERTJSONOFFERLIFE` dipanggil di tahap **penawaran** (`InputOfferLife_ACT` step 4), bukan di rantai simpan premium list — aturan urutan "paling akhir" perlu dikonfirmasi terhadap penempatan itu | Product+UW + Arsitektur Pega | tiket 01; tidak menyentuh AC 20–24 |
| ~~**OQ-068**~~ | ~~Tidak ada penulis `M_LIFE_PREMIUM_DETAIL` di PremiumList Life~~ — **TERTUTUP 2026-09-15**: work owner menambahkan `InsertLifePremiumDetail_act` (`ASM-FW-GISFW-WORK-LIFE` / `INSERTLIFEPREMIUMDETAIL_ACT`) + `SaveMasterLPDet` ke `PremiumList Life/`; pemicu job diganti **inline** | — | **nihil — tiket 08 naik ke ready** |
| **OQ-069** | Pesan validasi menjanjikan `NET_PREMIUM > GROSS_PREMIUM`, tetapi satu-satunya cek adalah `@PropertyHasValue` | Product+UW | tiket 04; AC validasi dibatasi pada keberadaan |

---

## Out of Scope

- **Endorsement Life — SELURUHNYA.** `[keputusan work owner]` Konteks/menu terpisah: class
  `ASM-FW-GISFW-WORK-ENDORSEMENTLIFE`, flow `Endorsement Life/Flow/InputEDMLife.xml`, dan alur inti
  "pilih polis NB yang sudah ada → muat data lamanya → endorse" yang tidak ada di sini. Termasuk:
  penomoran `PL_NUMBER_EDM`, layar & unggah CSV endorsement, pembatalan endorsement, dan efek keluar
  jalur endorsement. Spec dan tiketnya di `.scratch/endorsement-life/`. **Mesin bersama**
  (`SaveMasterLPDet`, `InsertPLSummary`, penomoran, `M_LIFE_PREMIUM_DETAIL`, urutan transaksi
  campuran) tetap **di dalam** cakupan spec ini sebagai komponen.
- **Claim Life dan Komite Claim Life.** Konteks hilir; spec ini berhenti pada penulisan rekam
  premium.
- **Lini non-Life.** Treaty, facultative, dan claim non-life punya konteks sendiri.
- **Merapikan alur.** Urutan tahap dan bentuk keputusan dibawa apa adanya; tiga penyimpangan sadar
  sudah didaftar di §Solution.
- **Membangun shape "Input Premium List Summary"** sebelum OQ-023 dijawab.
- **Identity & Access.** `[terverifikasi]` ABSENT dari korpus; dibangun dari nol.
- **Scaffolding kode.** Pekerjaan terpisah yang mendahului tiket mana pun.
- ~~**Migrasi data premium list lama.** Menunggu DDL fisik (OQ-001).~~ ⚠️ **TIDAK LAGI di luar
  cakupan (2026-09-16)** — OQ-001 ditutup; migrasi menjadi **irisan PREFACTOR, tiket pertama**
  (§12).

⚠️ **Tidak dimigrasikan — JSON dan perakitnya** `[keputusan work owner]` (§12):

| Yang dibuang | Alasan |
| --- | --- |
| `JSON_POLIS.DATA_JSON` dan `POOLDATA.INSERTJSONPOLISLIFE` | polis disimpan relasional |
| `JSON_OFFER_LIFE` dan `POOLDATA.INSERTJSONOFFERLIFE` | idem |
| **`@ASM.GetPageJSONString()`** | `[terverifikasi]` perakit payload — 211 panggilan di korpus, **selalu tanpa argumen**, definisinya (Rule-Utility-Function pustaka `ASM`) **tidak diekspor**. Tidak direplikasi |
| Simpul **`OfferFacIn`** | `ViewSuggest` naik jadi anak header |
| Penulisan ke `M_LIFE_PREMIUM_SUMMARY`, `M_LIFE_PREMIUM_DETAIL`, `LIFEINPRODUCTION` | ketiganya **hanya dibaca** saat migrasi |

---

## Further Notes

**Ukuran pekerjaan.** `[terverifikasi]` PremiumList Life: 37 Activity, 25 RDBList, 19 Section,
10 FlowAction, 9 Harness, 9 ReportDefinition, 7 DataTransform, 2 DecisionTable, 2 When,
1 ConnectREST, 1 HTMLRule, 1 SystemSettings, 1 Flow di akar — **124 berkas, dan hanya itu cakupan
spec ini.** (Endorsement Life, 75 berkas, dihitung di konteksnya sendiri.)

**Urutan yang saya sarankan untuk `/to-tickets`** — vertical slice:

1. Penawaran: Input Offer + keputusan Confirm/Reject/Decline.
2. Percabangan Offer / Premium.
3. Premium List Detail manual + penomoran `PL_NUMBER`.
4. Periode tutup buku dari tabel.
5. Unggah CSV: staging, validasi kolom uang, tinjau, simpan.
6. Penyimpanan polis: keempat procedure, dengan batas transaksi campuran ditangani sadar.
7. Efek keluar + flag lingkungan.
8. Penulisan detail peserta inline saat simpan (bukan job).
9. Kontrak hilir: verifikasi Claim Life dapat menemukan rekam lewat `PL_NUMBER`.

**Keputusan yang sudah diambil sebelum tiket dibuat** (2026-09-15):

1. **Urutan pemanggilan procedure relatif terhadap commit** `[keputusan desain]` — batas transaksi
   jalur Life **campuran**. Satu transaksi Go memuat `PROC_GENERATE_SEQUENCE_NUMBER` +
   `PEGA_M_LIFE_PREMIUM_SUMMARY` lalu commit; barulah `INSERTJSONPOLISLIFE` dan
   `INSERTJSONOFFERLIFE` (yang commit sendiri) dipanggil, **setelah** commit tersebut. Lihat §6.
2. **Fallback `25`** `[keputusan work owner]` — **tidak ditiru**; gagal terang-terangan. AC 8 sudah
   diperbarui.
3. **Pembagian paket domain di dalam `internal/`** — sama seperti dua spec sebelumnya; **masih
   terbuka**, tidak memblokir tiket.

**Catatan sumber.** Spec ini bersandar **hanya** pada korpus Pega `D:\XML\RNM_BRD\` dan artefak di
`OUTPUT_HASIL_RNM\`. Sumber ADR tunggal: `docs/adr/ADR-0001`…`ADR-0015`.
