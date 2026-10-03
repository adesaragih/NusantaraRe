# Spec — Master Contract Retro Life (migrasi Pega → Go + React + Oracle)

Status: **ready-for-agent** — ✅ **OQ-001 dan OQ-002 DITUTUP** `[data DBA]`
Konteks: `master-contract-retro-life` — **konteks/menu sendiri** `[keputusan work owner]`
Modul: **Master Contract Retro Life** (66 berkas — **modul terkecil di korpus**)
Tanggal: 2026-09-15
Sumber: `.scratch/master-contract-retro-life/grilling-ronde-1-jawaban.md` (**11 verdict final**),
`procedure-bodies-from-dba.md` `[data DBA]`, `ddl-tables-from-dba.md` `[data DBA]`,
`grilling-ronde-1.md` (temuan Bagian A–D), `CONTEXT.md`, `docs/adr/ADR-0001`–`ADR-0015`,
`discovery/modules/Master Contract Retro Life.md`, `discovery/context-map.md` §2.9
Skill: `/mattpocock-skills:to-spec`

> **Konvensi penandaan.** `[terverifikasi]` = terbukti korpus dengan **class + nama + path**;
> `[keputusan work owner]`; `[fakta bisnis — work owner]`; `[data DBA]`; `[terbuka]` = OQ.
> **Identitas rule wajib menyertakan class** — nama sama di class berbeda = rule berbeda.

> **Sumber tunggal.** Korpus Pega `D:\XML\RNM_BRD\` (READ-ONLY) dan artefak di
> `OUTPUT_HASIL_RNM\`, sekaligus repo target tunggal. ⛔ `D:\XML\nusantara-re\` di-blacklist.

> `[keputusan work owner]` **Master Product Name Life adalah konteks TERPISAH** — tidak dicakup spec
> ini. Pelajaran dari NB/EDM: menggabung menyembunyikan perbedaan.

---

## Problem Statement

Retrosesi lini life bersandar pada **data master** yang menentukan berapa risiko ditahan dan berapa
yang dilepas: kontrak treaty per tahun, batas proteksi per mata uang, siapa reinsurernya dan berapa
persen bagiannya, siapa yang meretrosesi bagian itu lagi, dan jenis bisnis apa yang tercakup. Hari
ini seluruhnya dikelola di empat layar grid Pega di atas lima tabel Oracle.

Lima hal membuat modul ini berisiko dipindahkan secara salah:

1. **Angkanya berjenjang, dan itu tidak terlihat dari datanya.** Share reinsurer tingkat kedua adalah
   persentase **dari share induknya**, bukan dari treaty. `10%` di tingkat dua atas induk `40%`
   berarti eksposur `4%` — bukan `10%`. Salah membacanya salah menghitung seluruh eksposur.
2. **Batas proteksi disimpan bersama selisihnya.** Kontrak menyimpan batas bawah, batas atas, **dan**
   selisih keduanya sebagai kolom terpisah — tiga angka untuk dua fakta, dan `[data DBA]` procedure
   **menulis selisih apa adanya dari parameter**; basis data tidak menghitungnya. Tidak ada penjaga
   agar ketiganya tetap sejalan.
3. **Menghapus induk meninggalkan anak yatim.** Keempat penghapus di korpus adalah
   `DELETE … WHERE ID = …` datar, tanpa kaskade.
4. **Kegagalan penyimpanan jatuh diam-diam.** Kelima procedure penulis mengembalikan pesan galat, dan
   **tidak satu pun** jalur simpan membacanya. Penyimpanan yang ditolak basis data tampak berhasil.
5. **Tidak ada yang menjaga aturan — di mana pun.** `[data DBA]` Body procedure terbukti berisi
   **upsert mentah + audit, nol logika bisnis**; DDL lama terbukti **nol `NOT NULL`, nol `UNIQUE`**.
   Aturan bisnis tidak ada di Pega, tidak ada di basis data. Ia harus dibangun.

## Solution

Membangun ulang Master Contract Retro Life sebagai **menu tersendiri** di atas Go + React + Oracle,
dengan perilaku dibawa apa adanya (paritas) kecuali **delapan penyimpangan sadar** yang sudah
diputus.

Modul ini **bukan proses berjenjang**: tidak ada persetujuan, tidak ada Submit/Decline, tidak ada
status akseptasi. Ia **editor master murni** — buka grid, tambah/ubah/hapus baris, simpan.

⚠️ **Konsekuensi paling mengikat dari data DBA:** karena procedure **nol validasi** dan DDL lama
**nol constraint**, **seluruh aturan bisnis wajib ditegakkan di lapisan Go** — total share, gerbang
tahun, wajib-isi, normalisasi, selisih dihitung, kaskade hapus. Basis data adalah **lapis pertahanan
kedua** (PK + FK yang kini sudah ada), **bukan** penegak aturan.

### Delapan penyimpangan sadar

| # | Penyimpangan | Alasan | Sumber |
| --- | --- | --- | --- |
| 1 | **Normalisasi** — `REINSTYPEID` sekali di kontrak, `TREATYYEAR` sekali di tahun; anak mewarisi | Pega menyimpan `REINSTYPEID` **3×** dan `TREATYYEAR` **2×** tanpa penjaga konsistensi; procedure pun tetap menulisnya ke tabel anak | `[keputusan work owner]` |
| 2 | **`*_SELISIH` dihitung, bukan disimpan mentah** | selisih adalah turunan dua batas; `[data DBA]` DB menulisnya apa adanya dari parameter — tidak ada yang menjaganya sinkron | `[keputusan work owner]` |
| 3 | **Kaskade hapus induk → seluruh anak, dengan popup konfirmasi Ya/Batal SEBELUM hapus** | Pega hapus datar → yatim. Work owner menetapkan hapus induk **memang** dimaksudkan menghapus sub-pohonnya, tetapi **tidak boleh tak sengaja** | `[keputusan work owner]` |
| 4 | **Gerbang tahun membandingkan nilai tanggal**, bukan potongan teks | `@substring(…,6,10)` rapuh; `[data DBA]` kolom tanggal memang bertipe `DATE`, jadi pemotongan teks memang keliru sejak awal | `[keputusan work owner]` |
| 5 | **"Terapkan ke semua" + pratinjau/konfirmasi + jejak audit** | Pega mengubah banyak baris tanpa konfirmasi dan tanpa jejak | `[keputusan work owner]` |
| 6 | **`HASIL1`/`o_message` wajib diperiksa — gagal terang-terangan, HTML dibersihkan** | Pega mengabaikannya; `[data DBA]` isinya pesan teks **bergaya HTML Pega** | `[keputusan work owner]` |
| 7 | **`ID` menjadi PRIMARY KEY di kelima tabel** — ✅ **sudah dieksekusi di basis data** | DDL lama: `ID VARCHAR2(100)` tanpa PK di empat tabel; keunikan tak dijamin padahal upsert berkunci `ID` | `[data DBA]` + `[keputusan work owner]` |
| 8 | **FK antar tabel** — ✅ **sudah dieksekusi di basis data**, mode **`ON DELETE CASCADE`** | DDL lama nol FK — inilah sebab risiko yatim; `CASCADE` selaras dengan penyimpangan 3 | `[data DBA]` + `[keputusan work owner]` |

### Kontrak batas

| Arah | Isi |
| --- | --- |
| **Keluar** | Modul ini **menghasilkan** data master retrosesi life di lima tabel `POOLDATA.*_LIFE`, dikonsumsi **Claim Life** (membaca `RETROCESSIONLIFE`, `TREATYYEAR_LIFE`), **Komite Claim Life** (class `ASM-FW-GISFW-INT-TREATYYEAR_LIFE`), dan **Master Product Name Life** (`TREATYCONTRACT_LIFE`, `TREATYYEAR_LIFE`) |
| **Masuk** | Enumerasi jenis reasuransi dari class `ASM-FW-GISFW-INT-REINSURANCETYPE`; daftar business dari `ASM-FW-GISFW-INT-BUSINESS`; tabel rate dari `ASM-FW-GISFW-INT-M_RATE_LIFE` dan `ASM-FW-GISFW-INT-RATE_LIFE_SUMMARY` |
| **Tidak ada** | `[terverifikasi]` **Nol integrasi luar** — tanpa `ConnectREST`, tanpa Google Storage, tanpa `M_LINK_SERVICE`, tanpa Arasapas, tanpa db-link, tanpa kolom JSON. Seluruh interaksi luar = **12 Connect-SQL ke `POOLDATA`** |

---

## User Stories

### Tahun treaty

1. Sebagai **admin master retro life**, saya ingin membuat **tahun treaty** baru beserta tahun
   underwriting dan rentang tanggalnya, supaya kontrak dapat digantungkan padanya.
2. Sebagai **admin master**, saya ingin tahun treaty **tidak dapat dihapus** setelah dibuat, supaya
   riwayat kontrak tahun-tahun lampau tidak pernah hilang. `[fakta bisnis — work owner]` — tahun
   treaty **abadi**, dan itu disengaja.
3. Sebagai **admin master**, saya ingin tahun treaty **wajib** punya tahun underwriting, tahun
   treaty, tanggal mulai, dan tanggal akhir, supaya tidak ada tahun setengah jadi.
4. Sebagai **admin master**, saya ingin **tidak perlu menentukan nomor identitas** baris baru, supaya
   penomoran tidak pernah bentrok antar petugas.

### Kontrak treaty dan batas proteksi

5. Sebagai **admin master**, saya ingin membuat **kontrak treaty** di bawah satu tahun treaty dengan
   **jenis reasuransi** tertentu, supaya proteksi per jenis terpisah rapi.
6. Sebagai **underwriter**, saya ingin menetapkan **batas bawah dan batas atas** proteksi per mata
   uang (IDR dan USD), supaya lebar layer proteksi jelas.
7. Sebagai **underwriter**, saya ingin **lebar layer dihitung sistem** dari kedua batas itu, supaya
   angkanya tidak pernah menyimpang dari batasnya sendiri.
8. Sebagai **underwriter**, saya ingin layer tersusun **menaik** — QS → 2nd QS → Surplus → 2nd
   Surplus — dengan batas atas satu layer menjadi batas bawah layer berikutnya, supaya proteksi
   bersambung tanpa celah maupun tumpang tindih. `[fakta bisnis — work owner]`
9. Sebagai **admin master**, saya ingin tahun pada **tanggal mulai kontrak** harus sama dengan tahun
   treaty induknya, supaya kontrak tidak nyasar tahun.
10. Sebagai **admin master**, saya ingin dapat mengisi batas **hanya dalam IDR** bila transaksinya
    memang rupiah, supaya saya tidak dipaksa mengarang angka USD. `[data DBA]` `USD` nullable = sah.
11. Sebagai **Finance**, saya ingin **tidak ada nilai batas yang berubah** saat menyeberang batas
    penyimpanan, supaya angka proteksi persis seperti yang saya tetapkan.

### Reinsurer dan share

12. Sebagai **admin master**, saya ingin menambahkan **reinsurer** ke sebuah kontrak beserta
    **persentase share**, komisi, dan overriding commission-nya, supaya pembagian risiko tercatat.
13. Sebagai **admin master**, saya ingin persentase share dan komisi **ditolak bila di luar 0–100**,
    supaya salah ketik tertangkap saat itu juga.
14. Sebagai **underwriter**, saya ingin melihat **total share** seluruh reinsurer pada satu kontrak,
    supaya saya tahu berapa yang sudah teralokasi.
15. Sebagai **underwriter**, saya ingin **tetap dapat menyimpan** meski total share belum 100%,
    supaya saya dapat menyusun kontrak bertahap. `[keputusan work owner]`
16. Sebagai **underwriter**, saya ingin kontrak yang totalnya **belum 100%** ditandai mencolok,
    supaya celahnya tidak luput dari perhatian.
17. Sebagai **manajemen**, saya ingin **laporan kontrak dengan total share ≠ 100%**, supaya seluruh
    celah terlihat dalam satu tempat.

### Retrosesi atas retrosesi

18. Sebagai **admin master**, saya ingin menambahkan **security reinsurer** di bawah seorang
    reinsurer beserta persentase sharenya, supaya retrosesi atas retrosesi tercatat.
19. Sebagai **underwriter**, saya ingin share security reinsurer dibaca sebagai **persentase dari
    share reinsurer induknya**, bukan dari treaty, supaya eksposurnya tidak terbaca terlalu besar.
    `[fakta bisnis — work owner]`
20. Sebagai **underwriter**, saya ingin sistem menghitung **eksposur efektif** sebagai share anak ×
    share induk, supaya angka yang saya lihat langsung bermakna terhadap treaty.

### Business

21. Sebagai **admin master**, saya ingin menambahkan **jenis business** ke sebuah kontrak beserta
    **rate reasuransinya**, supaya cakupan treaty jelas.
22. Sebagai **admin master**, saya ingin dapat **melihat tabel rate** sebelum memilih, supaya saya
    memilih rate yang benar.
23. Sebagai **admin master**, saya ingin **rate tersimpan apa adanya** seperti saya menuliskannya,
    supaya format rate bertingkat tidak rusak. `[data DBA]` `RIRATE` bertipe teks.
24. Sebagai **admin master**, saya ingin dapat **menerapkan satu business ke seluruh kontrak
    berjenis reasuransi sama** sekaligus, supaya saya tidak mengulang pekerjaan yang sama puluhan
    kali.
25. Sebagai **admin master**, saya ingin **diberi tahu berapa baris akan terpengaruh** sebelum
    penerapan massal itu dijalankan, supaya saya tidak mengubah lebih dari yang saya maksud.
26. Sebagai **auditor**, saya ingin penerapan massal **tercatat di jejak audit**, supaya perubahan
    puluhan baris dapat ditelusuri.

### Menghapus

27. Sebagai **admin master**, saya ingin menghapus induk **ikut menghapus seluruh anaknya**, supaya
    tidak ada baris yatim yang tertinggal. `[keputusan work owner]`
28. Sebagai **admin master**, saya ingin **diberi peringatan lebih dulu** — menyebut apa dan berapa
    yang akan ikut terhapus — supaya saya tidak menghapus sub-pohon secara tak sengaja.
29. Sebagai **admin master**, saya ingin **membatalkan** pada peringatan itu membuat **tidak ada
    satu pun** baris terhapus, supaya batal benar-benar berarti batal.

### Jenis reasuransi

30. Sebagai **admin master**, saya ingin memilih jenis reasuransi dari **daftar yang sah untuk lini
    life**, supaya tidak ada jenis dari lini lain yang masuk.
31. Sebagai **admin master**, saya ingin melihat **nama** jenis reasuransi, bukan kodenya, supaya
    layar terbaca manusia.

### Penyimpanan dan kegagalan

32. Sebagai **admin master**, saya ingin **diberi tahu bila penyimpanan ditolak** basis data, supaya
    saya tidak mengira data tersimpan padahal tidak. `[keputusan work owner]`
33. Sebagai **admin master**, saya ingin pesan galat itu **terbaca manusia**, bukan potongan markup,
    supaya saya tahu apa yang harus diperbaiki. `[data DBA]` galat DB mengandung HTML.
34. Sebagai **auditor**, saya ingin setiap penyimpanan mencatat **siapa** yang melakukannya, supaya
    perubahan master dapat ditelusuri.
35. Sebagai **organisasi**, saya ingin aturan penyimpanan tetap memanggil **procedure yang ada**,
    supaya tidak ada dua sumber kebenaran yang dapat menyimpang. `[keputusan work owner]`

### Skema dan migrasi

36. Sebagai **tim migrasi**, saya ingin skema target memuat **kunci primer dan kunci asing**, supaya
    keunikan dan keterhubungan dijamin basis data, bukan hanya oleh aplikasi.
37. Sebagai **tim migrasi**, saya ingin nilai uang dan share pindah **tanpa berubah satu digit pun**,
    supaya rekonsiliasi tidak menemukan selisih.
38. Sebagai **tim migrasi**, saya ingin **sequence** pindah dengan nilai berjalan yang benar, supaya
    identitas baru tidak bertabrakan dengan yang lama.

### Yang sengaja tidak dibawa

39. Sebagai **tim migrasi**, saya ingin **kolom selisih tidak ditulis mentah** dari masukan pengguna,
    karena ia turunan.
40. Sebagai **tim migrasi**, saya ingin **`REINSTYPEID` dan `TREATYYEAR` tidak digandakan** ke tabel
    anak.
41. Sebagai **tim migrasi**, saya ingin **perbandingan tahun lewat potongan teks tidak ikut pindah**.
42. Sebagai **tim migrasi**, saya ingin berkas non-Pega di dalam ekspor korpus **tidak dijadikan
    sumber perilaku**. `[terbuka]` OQ-054.

---

## Implementation Decisions

### 1. Batas konteks

`[keputusan work owner]` **Konteks ini = Master Contract Retro Life saja.** Master Product Name Life
(114 berkas) punya spec dan tiketnya sendiri. Keduanya menyentuh `TREATYCONTRACT_LIFE` dan
`TREATYYEAR_LIFE`, tetapi itu **tabel bersama**, bukan alasan menyatukan menu.

`[terverifikasi]` Modul ini **bukan proses berjenjang**: nol rujukan `StatusAkseptasi`, tanpa
`Akseptasi_DT`, tanpa tombol Submit/Decline; hanya **1 DataTransform** di seluruh modul.

`[terverifikasi]` **Modul paling terisolasi di korpus**: tanpa `ConnectREST`, tanpa `When`, tanpa
`SystemSettings`, tanpa db-link, tanpa kolom JSON, **tanpa jejak guard identitas sama sekali**.
Karena itu **ADR-0013** (resolusi endpoint) dan **ADR-0005** (flag lingkungan) **tidak berlaku** di
sini — tidak ada efek keluar untuk digerbangi.

Konteks hilir: **Claim Life**, **Komite Claim Life**, **Master Product Name Life**.

### 2. Penempatan modul

`[terverifikasi]` Struktur mengikat `CLAUDE.md` §5: kode di-scaffold **di dalam
`OUTPUT_HASIL_RNM\`**, arah dependency `handlers` → `services` → `repository`.

| Lapisan | Tanggung jawab |
| --- | --- |
| `models` | Tahun treaty, kontrak, reinsurer, security reinsurer, business; jenis reasuransi |
| `repository` | Pemanggilan lima procedure penulis; empat penghapus; dua belas kueri baca; **pemeriksaan `o_message`** |
| `services` | **Seluruh aturan bisnis** — kaskade hapus berkonfirmasi, gerbang tahun, wajib-isi, hitung selisih, hitung total share, hitung eksposur berjenjang, penerapan massal |
| `handlers` | Endpoint CRUD keempat grid; endpoint pratinjau hapus dan pratinjau penerapan massal; laporan share ≠ 100% |
| `frontend` | Empat layar grid; dua layar tampilan rate; penanda "belum 100%"; dua dialog konfirmasi |

### 3. Entitas dan hierarki

`[terverifikasi]` Terbaca penuh dari tanda tangan kelima procedure penulis, dikuatkan `[data DBA]`
oleh DDL dan keempat FK:

```
TREATYYEAR_LIFE            (ID, TREATYYEAR, UNDERWRITINGYEAR, STARTDATE, ENDDATE)
  └─ TREATYCONTRACT_LIFE   (ID, IDTREATYYEAR → tahun, REINSTYPEID,
                            TREATYSTARTDATE, TREATYENDDATE,
                            IDR, USD, B_IDR, B_USD, IDR_SELISIH, USD_SELISIH)
       ├─ TREATYREINSURER_LIFE         (ID, TREATYCONTRACTID → kontrak,
       │                                REINSURERID, REINSURERNAME,
       │                                PCTSHARE, COMMISION, OVR_COMM)
       │    └─ TREATYSECURITYREINSURER_LIFE (ID, TREATYREINSURERID → reinsurer,
       │                                     REINSURERID, REINSURERNAME, PCTSHARE)
       └─ TREATYBUSINESS_LIFE          (ID, TREATYCONTRACTID → kontrak,
                                        BIZCODE, BIZNAME, RIRATEID, RIRATE)
```

`[data DBA]` **Keempat FK sudah terpasang di basis data**, mode **`ON DELETE CASCADE`**:

| FK | Kolom anak | → induk |
| --- | --- | --- |
| 1 | `TREATYCONTRACT_LIFE.IDTREATYYEAR` | `TREATYYEAR_LIFE.ID` |
| 2 | `TREATYREINSURER_LIFE.TREATYCONTRACTID` | `TREATYCONTRACT_LIFE.ID` |
| 3 | `TREATYSECURITYREINSURER_LIFE.TREATYREINSURERID` | `TREATYREINSURER_LIFE.ID` |
| 4 | `TREATYBUSINESS_LIFE.TREATYCONTRACTID` | `TREATYCONTRACT_LIFE.ID` |

⚠️ **Penyimpangan sadar 1 — normalisasi.** `[keputusan work owner]` Di Pega,
`REINSTYPEID`/`REINSTYPENAME` tersimpan **tiga kali** (kontrak, reinsurer, business) dan
`TREATYYEAR` **dua kali** (tahun, business), **tanpa penjaga konsistensi** — dan `[data DBA]`
procedure **tetap menulisnya** ke tabel anak. Sistem baru:

- `REINSTYPEID` **hanya di kontrak**; reinsurer dan business mewarisinya lewat `TREATYCONTRACTID`.
- `TREATYYEAR` **hanya di tahun treaty**; business mewarisinya.
- `REINSTYPENAME` boleh disimpan sebagai **nilai tampilan turunan/di-cache**, ditandai jelas sebagai
  turunan — **bukan sumber kebenaran**, dan selalu dapat dibangun ulang dari `REINSTYPEID`.

⚠️ Karena procedure tetap menerima parameter-parameter itu, **lapisan repository yang memutuskan apa
yang dikirim** — normalisasi adalah keputusan skema baru, bukan tiruan perilaku procedure.

### 4. Empat grid

`[terverifikasi]` Empat Harness, seluruhnya class `DATA-PORTAL` — rasio Harness tertinggi di korpus:

| Harness | Identitas | Mengelola |
| --- | --- | --- |
| `InboxRetroLifeReinsurersList` | `DATA-PORTAL` / `INBOXRETROLIFEREINSURERSLIST` (530.052 byte) | **titik masuk** — daftar retrosesi life |
| `InboxRetroLimitReinsurers` | `DATA-PORTAL` / `INBOXRETROLIMITREINSURERS` (575.179 byte) | kontrak + batas proteksi |
| `InboxBusinessLifeReinsurers` | `DATA-PORTAL` / `INBOXBUSINESSLIFEREINSURERS` (486.911 byte) | business per kontrak |
| `InboxSecurityReinsurerLife` | `DATA-PORTAL` / `INBOXSECURITYREINSURERLIFE` (444.874 byte) | security reinsurer |

`[terverifikasi]` Delapan Section pendukung, seluruhnya `@BASECLASS`: `GRIDRETROCESSIONLIFE`,
`INPUTRETROCESSIONLIFE`, `INPUTDTLRETROCESSIONLIFE`, `INPUTRETROLIMITREINSURERS`,
`INPUTBUSINESSLIFEREINSURERS`, `INPUTSECURITYLIFEREINSURERS`, `INPUTSECURITYREINSURERLIFE`,
`VIEWRATE`.

`[terverifikasi]` Dua FlowAction, keduanya **tampilan saja**: `@BASECLASS` / `VIEWRATE` dan
`@BASECLASS` / `VIEWRATETABLE`, disiapkan `SetParamRate` (`ParamID.RIRATEID` ←
`InputBusinessLife.RIRATEID`) dan `SetParamRateTable`, membaca `ASM-FW-GISFW-INT-M_RATE_LIFE` /
`BROWSERATELIFE_RD` dan `ASM-FW-GISFW-INT-RATE_LIFE_SUMMARY` / `BROWSERATELIFESUMMARY`.

⚠️ `[data DBA]` **Jebakan memo.** `SetOutputParam_DT` (`@BASECLASS` / `SETOUTPUTPARAM_DT` /
`RULE-OBJ-MODEL`) bermemo **`not used`**, **tetapi masih dirujuk tiga Harness dan dua Section**.
Ini keluarga **OQ-066**: **jangan membuang perilaku hanya karena memo mengatakan tidak dipakai.**
Telusuri rujukannya sebelum memutuskan.

### 5. Batas proteksi dan lebar layer

`[fakta bisnis — work owner]` Arti keenam kolom limit kontrak:

| Kolom | Arti | Tipe `[data DBA]` |
| --- | --- | --- |
| `B_IDR` / `B_USD` | **batas bawah** (prefix `B_` = bawah) | `NUMBER` |
| `IDR` / `USD` | **batas atas** | `NUMBER` |
| `IDR_SELISIH` / `USD_SELISIH` | **atas − bawah** = **lebar layer** | `NUMBER` |

`[fakta bisnis — work owner]` **Layer tersusun menaik**, batas atas satu layer menjadi batas bawah
layer berikutnya:

```
QS  →  2nd QS  →  Surplus  →  2nd Surplus
```

⚠️ **Penyimpangan sadar 2 — selisih dihitung.** `[keputusan work owner]` Selisih adalah **turunan**
dari kedua batas. `[data DBA]` `INSERTTREATYCONTRACT_LIFE` **menulis `IDR_SELISIH`/`USD_SELISIH` apa
adanya dari parameter** — basis data **tidak** menghitungnya. Karena itu **Go yang menghitung**
sebelum mengirim, dan **tidak ada** jalur yang membiarkan pengguna mengetik selisih langsung.

`[data DBA]` ✅ **ADR-0003 aman**: kolom uang bertipe **`NUMBER` tanpa presisi** — angka presisi
penuh Oracle, bukan teks. Go memakai **desimal presisi arbitrer**; **tidak** lewat `float`.

`[data DBA]` **`USD` nullable = aturan sah**, bukan kelalaian — konsisten dengan `SaveTreatyLimit_Act`
`[terverifikasi]` yang mewajibkan `IDR`, `B_IDR`, `B_USD` tetapi **tidak** `USD`.

### 6. Dua tingkat share dan eksposur berjenjang

`[fakta bisnis — work owner]` `TREATYSECURITYREINSURER_LIFE.PCTSHARE` adalah **persentase dari share
reinsurer induknya**, **bukan** dari keseluruhan treaty. Ini **retrosesi atas retrosesi**.

> Reinsurer A memperoleh **40%** treaty. Security Reinsurer X di bawah A memperoleh **10%**.
> **Eksposur X terhadap treaty = 10% × 40% = 4%.**

`[terverifikasi]` Bukti hierarki: `INSERTSECURITYREINSURER_LIFE` menerima **`p_TREATYREINSURERID`**
yang menunjuk baris reinsurer induk. `[data DBA]` Kini ditegakkan **FK 3** di basis data.

**Konsekuensi mengikat:** setiap tampilan atau laporan eksposur **wajib mengalikan berjenjang**.
Menampilkan `10%` tanpa konteks induknya adalah kesalahan.

`[data DBA]` `PCTSHARE` bertipe `NUMBER` di kedua tabel — desimal, bukan teks.

### 7. Total share dan validasi

`[terverifikasi]` `CountingPercentShare_Act` (`@BASECLASS` / `COUNTINGPERCENTSHARE_ACT` /
`RULE-OBJ-ACTIVITY`) **hanya menjumlahkan** `PCTSHARE` seluruh reinsurer pada satu
(`TREATYYEARID`, `TREATYCONTRACTID`) lewat `GetMasterReinsurerLifeList_SQl`
(`ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!GETMASTERREINSURERLIFELIST_SQL`), lalu menaruh
hasilnya di field `STDRATING` untuk **ditampilkan**. **Tidak ada perbandingan terhadap 100 di
seluruh modul**, dan `[data DBA]` **tidak ada pula di procedure**.

`[keputusan work owner]` **Total share TIDAK memblokir penyimpanan.** Yang berlaku:

- total ditampilkan **mencolok**;
- bila ≠ 100%, diberi **tanda "belum 100%"**;
- tersedia **laporan kontrak dengan total share ≠ 100%**;
- penyimpanan **tetap diperbolehkan** — master lazim disusun bertahap.

`[terverifikasi]` **Validasi per baris tetap berlaku** — `SetErrorMessageReinsurer` (`@BASECLASS` /
`SETERRORMESSAGEREINSURER` / `RULE-OBJ-ACTIVITY`):

```
@toDecimal(PctShare) > 100 || @toDecimal(PctShare) < 0   → galat
@toDecimal(Ricomm)   > 100 || @toDecimal(Ricomm)   < 0   → galat
```

`[terverifikasi]` **Field wajib** saat simpan:

| Entitas | Wajib |
| --- | --- |
| Tahun treaty | `UNDERWRITINGYEAR`, `TREATYYEAR`, `STARTDATE`, `ENDDATE` |
| Kontrak | `REINSTYPEID`, `TREATYSTARTDATE`, `TREATYENDDATE`, `B_IDR`, `IDR`, `B_USD` |

⚠️ `[data DBA]` **Basis data tidak menegakkan satu pun** — **nol `NOT NULL`** di kelima tabel.
**Wajib-isi seluruhnya ditegakkan di Go.**

### 8. Gerbang konsistensi tahun

`[terverifikasi]` Tiga activity simpan berbagi gerbang:

```
@substring(InputRetrocessionLife.TREATYSTARTDATE,6,10) <> InputRetrocessionLifeTreatyType.TREATYYEAR_LIFE
```

⚠️ **Penyimpangan sadar 4.** `[keputusan work owner]` **Aturannya dipertahankan** — tahun pada
tanggal mulai kontrak harus sama dengan tahun treaty induknya. **Caranya diperbaiki**: bandingkan
**tahun dari nilai bertipe tanggal**.

`[data DBA]` **DDL menguatkannya**: `TREATYSTARTDATE`, `TREATYENDDATE`, `STARTDATE`, `ENDDATE`, dan
`TGLUPDATE` seluruhnya bertipe **`DATE`**. Pemotongan teks posisi 6–10 atas nilai `DATE` memang
keliru sejak awal.

### 9. Menghapus — kaskade dengan konfirmasi

`[terverifikasi]` Keempat penghapus di korpus, seluruhnya **`DELETE … WHERE ID = …` datar, tanpa
kaskade**:

| Rule | Class / Nama | Sasaran |
| --- | --- | --- |
| `DeleteTreatyLimit_SQL` | `ASM-FW-GISFW-INT-RETROCESSIONLIFE` / `ASM!DELETETREATYLIMIT_SQL` | `treatycontract_life` |
| `DeleteSecurityReinsurer_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!DELETESECURITYREINSURER_SQL` | `TREATYREINSURER_LIFE` |
| `DeleteSecurityReinsurerLife_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!DELETESECURITYREINSURERLIFE_SQL` | `TREATYSECURITYREINSURER_LIFE` |
| `DeleteRowBusinessList` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` / `ASM!DELETEROWBUSINESSLIST` | `treatybusiness_life` |

⚠️ **Penyimpangan sadar 3 — kaskade + konfirmasi.** `[keputusan work owner]` **Menghapus induk ikut
menghapus seluruh anaknya:**

| Induk dihapus | Ikut terhapus |
| --- | --- |
| **Kontrak** | seluruh reinsurer + seluruh security reinsurer di bawahnya + seluruh business |
| **Reinsurer** | seluruh security reinsurer di bawahnya |

**Popup konfirmasi muncul LEBIH DULU** — sebelum apa pun terhapus — menyebut **apa dan berapa** yang
akan ikut terhapus, misalnya:

> *"Menghapus kontrak ini akan ikut menghapus 3 reinsurer dan 5 business. Lanjut?"* **[Ya] [Batal]**

**Ya → kaskade berjalan. Batal → tidak ada satu pun baris terhapus.**

`[data DBA]` **FK di basis data dipasang `ON DELETE CASCADE`**, selaras. Basis data adalah **lapis
kedua**; pesan yang menyebut jumlah anak tetap **tanggung jawab Go**, karena FK tidak dapat
menjelaskan apa pun kepada pengguna.

`[fakta bisnis — work owner]` **Tahun treaty abadi** — `TREATYYEAR_LIFE` **tidak punya jalur hapus**,
dan itu **disengaja**. `[terverifikasi]` Korpus: lima penulis, hanya **empat** penghapus.

`[terverifikasi]` Catatan penamaan: `DeleteTreatyLimit_SQL` ber-class
`ASM-FW-GISFW-INT-RETROCESSIONLIFE` tetapi menghapus `treatycontract_life` — class tidak sejalan
dengan tabel sasarannya. Kejanggalan penamaan, bukan cacat perilaku.

### 10. Jenis reasuransi

`[terverifikasi]` Enumerasi disimpan di class **`ASM-FW-GISFW-INT-REINSURANCETYPE`**, dibaca lewat
`BrowseReinsuranceTypeLimit_RD` (`ASM-FW-GISFW-INT-REINSURANCETYPE` / `BROWSEREINSURANCETYPELIMIT_RD`
/ `RULE-OBJ-REPORT-DEFINITION`) dengan penyaring **`.Flag = 1`**, dan **nama ditampilkan dari kolom
`.Note`**.

`[terverifikasi]` `TreatyLimit_TypeProtect` (`@BASECLASS` / `TREATYLIMIT_TYPEPROTECT` /
`RULE-OBJ-ACTIVITY`) adalah **lookup**, bukan rumus: ia mencocokkan `REINSTYPEID` dengan `.ID` lalu
menyalin `.Note` ke `REINSTYPENAME`.

`[fakta bisnis — work owner]` **Lima jenis untuk lini life:**

| ID | Nama (`.Note`) |
| --- | --- |
| 10196 | QS |
| 10197 | 2ND QS |
| 10198 | SURPLUS |
| 10199 | 2ND SURPLUS |
| 10200 | OR |

`[fakta bisnis — work owner]` **`Flag = 1` berarti "for life"** — penanda lini, **bukan**
aktif/nonaktif. Sejalan dengan memo korpus `[terverifikasi]` `pyUsage: "Parameter Flag, 1 for life"`.

`[fakta bisnis — work owner]` Kolom `.Code` dan `.Type` **tidak dipakai** sistem.

⚠️ `"OR"` dicatat **apa adanya**; artinya tidak ditafsirkan.

`[terverifikasi]` RD serupa (`BrowseReinsuranceType_RD`, class sama) ada di **sembilan modul lain** —
enumerasi ini dipakai lintas lini. Menutup pertanyaan ini menutup **OQ-057** untuk life.

### 11. Jalur simpan — lima procedure, dipanggil apa adanya

`[keputusan work owner]` **Go memanggil kelima stored procedure Oracle apa adanya; tidak ditulis
ulang.** Konsisten dengan empat konteks Life sebelumnya dan **ADR-0006**.

`[terverifikasi]` Tanda tangannya:

| Rule Connect-SQL | Class / Nama | Procedure |
| --- | --- | --- |
| `SaveMasterTreatyYear_Life_SQL` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `ASM!SAVEMASTERTREATYYEAR_LIFE_SQL` | `INSERTTREATYYEAR_LIFE` |
| `SaveMasterTreatyContract_Life_SQL` | `ASM-FW-GISFW-INT-TREATYCONTRACT_LIFE` / `ASM!SAVEMASTERTREATYCONTRACT_LIFE_SQL` | `INSERTTREATYCONTRACT_LIFE` |
| `SaveMasterTreatyReinsurer_Life_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!SAVEMASTERTREATYREINSURER_LIFE_SQL` | `INSERTREINSURER_LIFE` |
| `SaveMasterTreatySecurityReinsurer_Life_SQL` | `ASM-FW-GISFW-INT-TREATYREINSURER_LIFE` / `ASM!SAVEMASTERTREATYSECURITYREINSURER_LIFE_SQL` | `INSERTSECURITYREINSURER_LIFE` |
| `SaveMasterTreatyBusiness_Life_SQL` | `ASM-FW-GISFW-INT-TREATYBUSINESS_LIFE` / `ASM!SAVEMASTERTREATYBUSINESS_LIFE_SQL` | `INSERTBUSINESS_LIFE` |

`[data DBA]` **Perilaku kelimanya IDENTIK** — kini diketahui penuh (OQ-002 ditutup):

| Aspek | Perilaku |
| --- | --- |
| **Mode simpan** | **UPSERT dikunci `ID`** — `SELECT COUNT(1) … WHERE ID = p_ID`; ada → `UPDATE`, tidak ada → `INSERT` |
| **Identitas baris baru** | **dibuat basis data**: `'1' \|\| lpad(<sequence>.nextval, 6, '0')` → mis. `1000044`. `p_ID` **diabaikan** saat INSERT, dipakai hanya untuk mencari saat UPDATE |
| **`TGLUPDATE`** | selalu **`SYSDATE`** — parameter `p_TGLUPDATE` **diabaikan** |
| **`USERID`** | disimpan apa adanya dari `p_USERID` |
| **Transaksi** | **`COMMIT` di dalam tiap procedure**; `ROLLBACK` pada exception terluar |
| **Validasi bisnis** | **NOL** — tidak ada cek 100%, tidak ada cek anak, tidak ada unique selain PK `ID` |

**Konsekuensi mengikat:**

1. Repository "save" adalah **upsert**, bukan insert murni.
2. **Aplikasi tidak menetapkan identitas baris baru** — kirim `ID` kosong; basis data yang
   membuatnya. Konsisten **ADR-0006**.
3. **Jangan andalkan cap waktu aplikasi** — basis data yang menetapkannya.
4. **Setiap simpan adalah transaksi mandiri**; tidak ada transaksi lintas-baris di sisi basis data.
   Penyimpanan banyak baris **tidak atomik** — dan itu harus disadari, bukan disembunyikan.

⚠️ `[data DBA]` **Anti-dobel logis tidak dijamin.** Upsert dikunci `ID` mencegah dua baris ber-`ID`
sama — **tidak** mencegah dua reinsurer sama di satu kontrak. Bila anti-dobel logis dikehendaki, itu
aturan Go dan idealnya `UNIQUE` di basis data. → keputusan tiket.

### 12. Penegakan `HASIL1` / `o_message`

`[terverifikasi]` Kelima procedure mengembalikan `{OutputData.HASIL1 out}`. Sensus 27 Activity modul
ini: **hanya `DeleteRowBusiness.xml`** yang menyebut `HASIL1`; **tidak satu pun activity `Save*`
membacanya**. Galat procedure **jatuh diam-diam**.

`[data DBA]` **Semantiknya kini diketahui:**

| Nilai | Arti |
| --- | --- |
| **kosong / `NULL`** | **sukses** |
| **berisi teks** | **gagal** |

⚠️ `[data DBA]` Isinya **pesan bergaya UI Pega yang mengandung HTML**
(`<span style="color:red">…</span>`) ditambah `SQLERRM`.

⚠️ **Penyimpangan sadar 6 — gagal terang-terangan, HTML dibersihkan.** `[keputusan work owner]`

- `o_message` **wajib diperiksa** setelah **setiap** pemanggilan procedure.
- **Non-kosong = gagal.** Penyimpanan yang ditolak basis data **tidak boleh** tampak berhasil.
- **HTML dari Pega JANGAN diteruskan mentah** ke API maupun layar — ekstrak maknanya, sajikan pesan
  bersih yang terbaca manusia.

Sekeluarga dengan keputusan "gagal terang-terangan" pada ambang tutup buku di PremiumList Life.

### 13. Terapkan ke semua

`[terverifikasi]` `SaveBusinessToAllLife_Act` menyimpan business ke **seluruh** baris berjenis
reasuransi sama (`.REINSTYPEID == Param.REINSTYPEID`), lewat Connect-SQL
`SaveTreatyBusinessAll_Life_SQL`.

⚠️ **Penyimpangan sadar 5.** `[keputusan work owner]` Fitur **dibawa**, dengan dua pengaman yang
tidak ada di Pega: **pratinjau jumlah baris + konfirmasi** sebelum dijalankan, dan hasilnya
**tercatat di jejak audit**.

⚠️ `[data DBA]` Karena tiap procedure **commit sendiri**, penerapan massal **tidak atomik**: bila
gagal di tengah, sebagian baris sudah tersimpan. Perilaku itu harus **terlihat** — laporkan berapa
berhasil dan berapa gagal, jangan tampilkan sukses tunggal yang menyesatkan.

### 14. Skema Oracle target

`[data DBA]` DDL kelima tabel diterima; **OQ-001 ditutup**.

| Kelompok | Tipe | Catatan |
| --- | --- | --- |
| Uang: `IDR`, `USD`, `B_IDR`, `B_USD`, `IDR_SELISIH`, `USD_SELISIH` | **`NUMBER`** (tanpa presisi) | ✅ **ADR-0003 aman** — angka presisi penuh, bukan teks |
| Share/komisi: `PCTSHARE`, `COMMISION`, `OVR_COMM` | **`NUMBER`** | desimal |
| ⚠️ `RIRATE` | **`VARCHAR2(1000)`** | **teks** — bawa apa adanya; **jangan** paksa jadi angka |
| Tanggal: `TGLUPDATE`, `STARTDATE`, `ENDDATE`, `TREATYSTARTDATE`, `TREATYENDDATE` | **`DATE`** | menguatkan §8 |
| Identitas: `ID`, `*ID` | `VARCHAR2(100)` | |
| Nama: `REINSTYPENAME`, `REINSURERNAME`, `BIZNAME` | `VARCHAR2(1000)` | lebar untuk cache nama |
| Nullability | **semua kolom nullable, nol `NOT NULL`** | wajib-isi **di Go** |

⚠️ **Penyimpangan sadar 7 — `ID` menjadi PRIMARY KEY di kelima tabel.** ✅ `[data DBA]` +
`[keputusan work owner]` **sudah dieksekusi di basis data.** Semula hanya `TREATYBUSINESS_LIFE` yang
punya PK (`TREATYBUSINESS_LIFE_PK`); empat lainnya tidak — padahal upsert berkunci `ID`.

⚠️ **Penyimpangan sadar 8 — FK antar tabel, `ON DELETE CASCADE`.** ✅ `[data DBA]` +
`[keputusan work owner]` **sudah dieksekusi di basis data.** Keempatnya tercantum di §3.

`[data DBA]` **Sequence**: `TREATYYEAR_LIFE_SEQ` `START WITH 44` → identitas tahun berikutnya
`1000044`. Empat sequence lain dirujuk procedure (`TREATYCONTRACT_LIFE_seq`,
`TREATYREINSURER_LIFE_SEQ`, `TREATYSECURITYREINSURER_LIFE_SEQ`, `TREATYBUSINESS_LIFE_SEQ`).
Tablespace `TBS_POOLDATA`.

### 15. Tiga activity yang ternyata bukan rumus

`[terverifikasi]` Batas telusur D3 menandai tiga activity `@BASECLASS` sebagai "rumus belum dibaca".
Ketiganya sudah dibuka pada grilling Ronde 1 — dan **dua di antaranya bukan rumus**:

| Activity | Identitas | Kenyataan |
| --- | --- | --- |
| `CountingPercentShare_Act` | `@BASECLASS` / `COUNTINGPERCENTSHARE_ACT` | **penjumlah** `PCTSHARE`, hasilnya ke field `STDRATING` |
| `TreatyLimit_TypeProtect` | `@BASECLASS` / `TREATYLIMIT_TYPEPROTECT` | **lookup** `REINSTYPEID` → `.Note` |
| `SetValueRetroLimit_TreatyYearLife` | `@BASECLASS` / `SETVALUERETROLIMIT_TREATYYEARLIFE` | **penyalur konteks** + kendali tampilan |

Digabung dengan `[data DBA]` bahwa procedure pun **nol logika bisnis**, kesimpulannya tegas:
**tidak ada perhitungan retro life di mana pun pada sistem lama.** Aturan yang dispesifikasikan di
sini **dibangun**, bukan dimigrasikan.

---

## Testing Decisions

### Apa yang membuat test baik di sini

Test memeriksa **perilaku yang terlihat dari luar** — apa yang tersimpan, apa yang ditolak, apa yang
tampil — bukan susunan internal.

Tiga hal di konteks ini **wajib** diuji terhadap data nyata: **angka** (batas, selisih, share,
eksposur berjenjang), **penolakan** (validasi, `o_message`), dan **kaskade** (apa yang ikut terhapus,
dan apa yang tidak terhapus saat dibatalkan). Ketiganya adalah inti kebenaran master.

### Seam — **memakai ulang seam yang sudah ada**

`[terverifikasi]` Repo target belum di-scaffold. Seam yang ditetapkan spec Claim — Life adalah
**API HTTP**, dan konteks ini memakainya kembali — **tidak menambah seam**:

> **Seam utama: API HTTP Master Contract Retro Life.** Test menggerakkan CRUD keempat grid lewat
> endpoint REST dan memeriksa hasilnya lewat endpoint REST, dengan `handlers → services → repository`
> terpasang sungguhan, terhadap skema uji Oracle.

**Batas proses difake:**

| Batas | Perlakuan |
| --- | --- |
| Oracle | **skema uji nyata, bukan mock** — kelima procedure **adalah** perilaku simpan (upsert, ID dari sequence, commit internal), dan kaskade FK hanya berperilaku benar pada basis data sungguhan |
| **Tidak ada yang lain** | `[terverifikasi]` modul ini **tanpa integrasi luar sama sekali** — tidak ada Arasapas, email, Google Storage, maupun REST. **Tidak ada yang perlu difake.** |

**Tidak ada seam kedua.** Konteks ini tidak punya worker asinkron, tidak punya efek keluar.

### Modul yang diuji

| Yang diuji | Lewat seam |
| --- | --- |
| CRUD keempat grid, satu per satu | API HTTP + skema uji |
| Identitas baris baru **dibuat basis data**, bukan aplikasi | API HTTP + skema uji |
| Simpan ulang baris yang sama = **upsert**, bukan baris kedua | API HTTP + skema uji |
| Hierarki: anak tidak dapat lahir tanpa induk (ditegakkan FK) | API HTTP + skema uji |
| Lebar layer **dihitung** dari kedua batas; tidak dapat diketik langsung | API HTTP |
| Layer menaik: batas atas layer = batas bawah layer berikutnya | API HTTP |
| Total share dihitung dan ditampilkan; **tidak** memblokir | API HTTP |
| Validasi per baris `0..100`; wajib-isi ditegakkan **Go** | API HTTP |
| **Eksposur berjenjang** = share anak × share induk | API HTTP |
| Gerbang tahun lewat **nilai tanggal** | API HTTP + jam terkendali |
| **Kaskade hapus**: apa yang ikut terhapus | API HTTP + skema uji |
| **Batal pada konfirmasi**: tidak ada satu pun baris terhapus | API HTTP + skema uji |
| Tahun treaty **tidak dapat dihapus** | API HTTP |
| Enumerasi jenis reasuransi disaring `Flag = 1` | API HTTP + skema uji |
| **`o_message` non-kosong = gagal**, dan HTML tidak bocor ke API | API HTTP + skema uji |
| Terapkan-ke-semua: pratinjau, konfirmasi, audit, **laporan sebagian-gagal** | API HTTP + skema uji |
| Uang tidak berubah menyeberang batas | API HTTP |
| `RIRATE` tersimpan **sebagai teks**, apa adanya | API HTTP + skema uji |

### Prior art

`[terverifikasi]` **Tidak ada** — nol kode, nol test di `OUTPUT_HASIL_RNM`. Spec Claim — Life,
Komite Claim Life, PremiumList Life, dan Endorsement Life menetapkan bentuknya; konteks ini mengikuti
bentuk yang sama.

Perintah verifikasi wajib ditulis eksplisit di tiap tiket selama `Makefile` belum ada. Target:
`go test ./internal/...` dan `cd frontend && npm test`.

---

## Acceptance Criteria

**Tahun treaty**

1. Tahun treaty dapat dibuat dengan tahun underwriting, tahun treaty, tanggal mulai, dan tanggal
   akhir; keempatnya **wajib** — ditegakkan **di Go**, karena basis data tidak punya `NOT NULL`.
2. ⚠️ Tahun treaty **tidak dapat dihapus** lewat jalur mana pun. Test yang menemukan endpoint hapus
   tahun **gagal**. `[fakta bisnis — work owner]`
3. Tahun treaty dapat diubah, dan perubahannya mencatat **siapa** pelakunya.
4. **Aplikasi tidak menetapkan identitas baris baru** — ia mengirim identitas kosong dan basis data
   yang membuatnya. Test yang menemukan pembentukan identitas di sisi aplikasi **gagal**.
   *(**ADR-0006**)*

**Kontrak dan batas proteksi**

5. Kontrak lahir **di bawah** satu tahun treaty; kontrak tanpa tahun induk **ditolak**.
6. Kontrak wajib memuat `REINSTYPEID`, tanggal mulai, tanggal akhir, **batas bawah IDR**, **batas
   atas IDR**, dan **batas bawah USD**.
7. **`USD` boleh kosong** — kontrak ber-IDR saja tersimpan tanpa keluhan. `[data DBA]` aturan sah.
8. ⚠️ **Lebar layer DIHITUNG** dari `batas atas − batas bawah`, per mata uang. **Tidak ada** jalur
   yang membiarkan pengguna mengetiknya, dan tidak ada jalur yang mengirimkannya mentah dari
   masukan. `[keputusan work owner]`
9. Lebar layer yang tersimpan **selalu** sama dengan selisih kedua batas, termasuk setelah salah satu
   batas diubah.
10. Batas bawah yang **lebih besar** dari batas atas **ditolak**, dengan pesan yang menyebut mata
    uangnya.
11. Susunan layer menaik dapat direkam: batas atas satu layer boleh menjadi batas bawah layer
    berikutnya tanpa dianggap tumpang tindih.
12. Seluruh nilai batas diperlakukan sebagai **desimal presisi arbitrer**; **tidak ada** yang
    melewati `float` di lapisan mana pun maupun di JSON API. *(**ADR-0003**; `[data DBA]` kolom
    bertipe `NUMBER`)*
13. Nilai batas yang ditulis dan dibaca kembali **identik** — tidak ada pembulatan diam.
14. ⚠️ **Gerbang tahun**: tahun pada tanggal mulai kontrak **harus sama** dengan tahun treaty
    induknya. Perbandingan memakai **tahun dari nilai tanggal**; test yang menemukan pemotongan teks
    posisi tetap **gagal**. `[keputusan work owner]`

**Reinsurer dan share**

15. Reinsurer lahir **di bawah** satu kontrak, dengan `PCTSHARE`, komisi, dan overriding commission.
16. `PCTSHARE` dan komisi di luar rentang **0–100 ditolak**, dengan pesan yang menyebut kolomnya.
17. Total share seluruh reinsurer pada satu kontrak **dihitung dan ditampilkan**.
18. ⚠️ Total share **TIDAK memblokir penyimpanan** — kontrak dengan total ≠ 100% **tetap tersimpan**.
    `[keputusan work owner]`
19. Kontrak dengan total ≠ 100% **ditandai mencolok** di layar.
20. Tersedia **laporan kontrak dengan total share ≠ 100%**.

**Retrosesi atas retrosesi**

21. Security reinsurer lahir **di bawah** satu reinsurer; tanpa reinsurer induk **ditolak**.
22. ⚠️ `PCTSHARE` security reinsurer dibaca sebagai **persentase dari share induknya**, bukan dari
    treaty. `[fakta bisnis — work owner]`
23. **Eksposur efektif** terhadap treaty dihitung **share anak × share induk**, dan itulah angka yang
    ditampilkan sebagai eksposur. Test memuat kasus `10% × 40% = 4%`.
24. Mengubah share **induk** mengubah eksposur efektif seluruh anaknya.

**Business**

25. Business lahir **di bawah** satu kontrak, dengan kode business, nama, dan rate reasuransi.
26. **`RIRATE` tersimpan sebagai teks apa adanya** — tidak dikonversi, tidak diformat ulang, tidak
    dipaksa menjadi angka. `[data DBA]` kolomnya `VARCHAR2(1000)`.
27. Tabel rate dapat **dilihat** sebelum memilih, tanpa mengubah apa pun.
28. ⚠️ **Terapkan ke semua** menampilkan **jumlah baris yang akan terpengaruh** dan **menunggu
    konfirmasi** sebelum dijalankan. `[keputusan work owner]`
29. ⚠️ Hasil penerapan massal **tercatat di jejak audit**, termasuk berapa baris berubah.
    `[keputusan work owner]`
30. Penerapan massal yang **dibatalkan** pada dialog konfirmasi **tidak mengubah apa pun**.
31. Bila penerapan massal gagal di tengah, sistem melaporkan **berapa berhasil dan berapa gagal** —
    tidak menampilkan sukses tunggal yang menyesatkan. `[data DBA]` tiap procedure commit sendiri.

**Menghapus**

32. ⚠️ Menghapus **kontrak** memunculkan **konfirmasi lebih dulu** yang menyebut **apa dan berapa**
    yang akan ikut terhapus — reinsurer, security reinsurer, dan business di bawahnya.
    `[keputusan work owner]`
33. ⚠️ Menghapus **reinsurer** memunculkan konfirmasi yang menyebut berapa **security reinsurer**
    akan ikut terhapus. `[keputusan work owner]`
34. ⚠️ Menekan **Ya** menghapus induk **beserta seluruh sub-pohonnya**; tidak ada baris yatim yang
    tertinggal. `[keputusan work owner]`
35. ⚠️ Menekan **Batal** membuat **tidak ada satu pun** baris terhapus — induk maupun anak.
    `[keputusan work owner]`
36. Menghapus baris **tanpa anak** tetap memerlukan konfirmasi, dengan pesan yang menyatakan tidak
    ada anak yang terpengaruh.
37. Menghapus security reinsurer atau business (baris daun) menghapus **hanya baris itu**.

**Jenis reasuransi**

38. Daftar jenis reasuransi yang ditawarkan hanya yang **`Flag = 1`** (for life) — lima jenis, ID
    **10196**–**10200**.
39. Layar menampilkan **nama** dari kolom `.Note`, bukan kode.
40. `REINSTYPEID` di luar kelima nilai itu **ditolak**.
41. ⚠️ **Normalisasi**: `REINSTYPEID` tersimpan **hanya pada kontrak**; reinsurer dan business
    mewarisinya lewat `TREATYCONTRACTID`. Test yang menemukan `REINSTYPEID` yang dapat ditulis pada
    tabel anak **gagal**. `[keputusan work owner]`
42. ⚠️ `TREATYYEAR` tersimpan **hanya pada tahun treaty**; business mewarisinya.
    `[keputusan work owner]`
43. `REINSTYPENAME` yang disimpan ditandai **turunan/cache**, dan selalu dapat dibangun ulang dari
    `REINSTYPEID`.

**Penyimpanan**

44. Kelima jalur simpan memanggil **stored procedure Oracle yang ada**. Aplikasi **tidak** menulis
    ulang logikanya. `[keputusan work owner]`
45. Menyimpan baris yang **sudah ada** memperbarui baris itu — **upsert dikunci identitas**, bukan
    baris kedua. Dibuktikan dengan menyimpan dua kali lalu menghitung baris. `[data DBA]`
46. ⚠️ **`o_message` DIPERIKSA** setelah **setiap** pemanggilan procedure: **kosong/NULL = sukses**,
    **berisi teks = gagal**. Kegagalan **ditampilkan**; penyimpanan yang ditolak **tidak pernah**
    tampak berhasil. Test yang menemukan pemanggilan tanpa pemeriksaan hasil **gagal**.
    `[keputusan work owner]`
47. ⚠️ **HTML dari pesan galat basis data tidak bocor** ke API maupun layar; pesan yang sampai ke
    pengguna **bersih dan terbaca**. `[keputusan work owner]`
48. Setiap penyimpanan mengirim **identitas pengguna**; cap waktu **tidak** dikirim aplikasi — basis
    data yang menetapkannya. `[data DBA]`
49. Uang dan persentase dikirim sebagai **desimal presisi arbitrer**, dan dibaca kembali **identik**.
    *(**ADR-0003**)*
50. Penyimpanan banyak baris **tidak diklaim atomik** — tiap baris adalah transaksi mandiri, dan
    perilaku itu terlihat pada antarmuka. `[data DBA]` `COMMIT` di dalam tiap procedure.

**Skema**

51. ⚠️ **`ID` adalah PRIMARY KEY** di kelima tabel. `[data DBA]` sudah terpasang.
52. ⚠️ **Keempat FK terpasang** dengan mode **`ON DELETE CASCADE`**, selaras dengan kaskade di Go.
    `[data DBA]` sudah terpasang.
53. Wajib-isi **ditegakkan aplikasi**, karena basis data **nol `NOT NULL`**. Test yang mengandalkan
    basis data untuk menolak nilai kosong **gagal**.
54. Skema uji yang dipakai seluruh test dibangun **dari DDL yang sama** dengan produksi — bukan dari
    tiruan yang ditulis terpisah.

**Kejanggalan yang tidak ditiru**

55. Tidak ada padanan field `STDRATING` yang dipakai untuk menyimpan total share — beri nama yang
    jujur.
56. Berkas non-Pega di dalam ekspor korpus **tidak** dijadikan sumber perilaku. `[terbuka]` OQ-054.

---

## Pertanyaan terbuka di dalam spec

**Nol OQ pemblokir.** ✅ **OQ-001 dan OQ-002 ditutup** `[data DBA]` pada 2026-09-15.

| OQ | Status | Pemilik | Catatan |
| --- | --- | --- | --- |
| **OQ-002** | ✅ **DITUTUP** | — | Body kelima procedure diterima: **upsert dikunci `ID` + audit, nol logika bisnis**; `o_message` = pesan galat teks |
| **OQ-001** | ✅ **DITUTUP** | — | DDL kelima tabel diterima: tipe, nullability, PK, FK, sequence. **Tiket migrasi naik ke `ready`** |
| **OQ-057** | ✅ **DITUTUP untuk life** | — | Lima jenis, ID 10196–10200; `Flag = 1` = for life |
| **OQ kecil** | `[terbuka]` — **tidak memblokir** | **DBA** | Nama **tabel fisik Oracle** untuk class `REINSURANCETYPE`; isi kolom `.Code` dan `.Type` — **tidak dipakai sistem** |
| **OQ-054** | `[terbuka]` | — | Berkas non-Pega di `Claude outputs/` — **tidak dibaca**, dicatat saja |
| **OQ-066** (keluarga) | `[terbuka]` | Arsitektur Pega | ⚠️ `SetOutputParam_DT` bermemo `not used` **tetapi dirujuk 3 Harness + 2 Section** — **jangan** membuang perilaku karena memo |

⚠️ **Anti-dobel logis** (mis. dua reinsurer sama di satu kontrak) **belum diputuskan** — `[data DBA]`
basis data tidak menjaganya, dan Pega pun tidak. Bukan OQ fakta, melainkan **keputusan tiket**.

---

## Out of Scope

- **Master Product Name Life — seluruhnya.** `[keputusan work owner]` Konteks/menu terpisah (114
  berkas), dengan spec dan tiketnya sendiri.
- **Claim Life, Komite Claim Life.** Konteks hilir yang **membaca** master ini; spec ini berhenti
  pada penulisan master.
- **Lini non-Life.** Treaty, facultative, dan claim non-life punya konteks sendiri — termasuk
  enumerasi jenis reasuransi di luar `Flag = 1`.
- **Menulis ulang aturan penyimpanan ke Go.** `[keputusan work owner]` Procedure dipanggil apa
  adanya. Yang dibangun di Go adalah **aturan bisnis** yang memang tidak ada di mana pun.
- **Mengubah tipe `RIRATE`.** `[data DBA]` Ia `VARCHAR2(1000)`; mengubahnya menjadi angka menuntut
  konfirmasi format terlebih dahulu.
- **Identity & Access.** `[terverifikasi]` Modul ini **tanpa jejak guard identitas sama sekali** —
  satu-satunya di Tahap 5. Model peran dirancang terpisah (OQ-007 dan keluarganya).
- **Efek keluar dan flag lingkungan.** `[terverifikasi]` Modul ini tidak punya integrasi luar;
  **ADR-0013** dan **ADR-0005** tidak berlaku.
- **Scaffolding kode.** Pekerjaan terpisah yang mendahului tiket mana pun.

---

## Further Notes

**Ukuran pekerjaan.** `[terverifikasi]` **66 berkas — modul terkecil di korpus**: 27 Activity,
12 ReportDefinition, 12 Connect-SQL, 8 Section, 4 Harness, 2 FlowAction, 1 DataTransform.
**Nol `Flow`, nol `When`, nol `ConnectREST`.**

**Mengapa modul kecil ini tetap perlu kehati-hatian.** Ia **master data** — satu angka salah di sini
merambat ke perhitungan klaim dan eksposur di konteks lain. Tiga jebakan yang sudah terbukti:

| Jebakan | Akibat bila terlewat |
| --- | --- |
| Share tingkat dua dibaca sebagai % dari treaty | eksposur retrosesi **salah besar** |
| Selisih dikirim mentah dari masukan | lebar layer menyimpang dari batasnya sendiri |
| `o_message` diabaikan | penyimpanan gagal **tampak berhasil** |

**Temuan paling menentukan dari data DBA.** `[data DBA]` Procedure berisi **upsert mentah + audit,
nol logika bisnis**; DDL lama berisi **nol `NOT NULL`, nol `UNIQUE`, nol FK**. Digabung dengan
temuan grilling bahwa ketiga "rumus" Pega bukan rumus, kesimpulannya:

> **Tidak ada aturan bisnis di mana pun pada sistem lama — tidak di Pega, tidak di basis data.**
> Yang dispesifikasikan di sini sebagian besar **dibangun**, bukan dimigrasikan. Karena itu
> penyimpangan sadar di konteks ini berjumlah **delapan**, terbanyak di antara lima konteks Life.

**Nama yang menyesatkan di korpus** — dicatat agar tidak ditiru:

| Yang tertulis | Yang sebenarnya |
| --- | --- |
| `CountingPercentShare_Act` "rumus pembagian share" | hanya **menjumlahkan** |
| `TreatyLimit_TypeProtect` "jenis proteksi limit" | hanya **lookup nama** |
| `SetValueRetroLimit_TreatyYearLife` "nilai retro limit" | **penyalur konteks**, tanpa perhitungan |
| field `STDRATING` "standard rating" | menyimpan **total share** |
| `DeleteTreatyLimit_SQL` di class `RETROCESSIONLIFE` | menghapus `treatycontract_life` |
| `SetOutputParam_DT` bermemo `not used` | **masih dirujuk** 3 Harness + 2 Section |
| `p_TGLUPDATE` sebagai parameter | **diabaikan** — `SYSDATE` yang dipakai |
| `p_ID` saat INSERT | **diabaikan** — sequence yang dipakai |

**Baca kodenya, jangan namanya** — pelajaran yang sama dengan `"Set 0 jika EDM Batal"` di Endorsement
Life, dan kini diperluas: **baca juga body procedure-nya, jangan tanda tangannya.**

**Urutan yang saya sarankan untuk `/to-tickets`** — vertical slice:

1. Tahun treaty: CRUD + aturan abadi + identitas dari sequence.
2. Kontrak: CRUD + batas per mata uang + **lebar layer dihitung**.
3. Gerbang tahun lewat nilai tanggal.
4. Jenis reasuransi: enumerasi `Flag = 1` + normalisasi `REINSTYPEID`.
5. Reinsurer: CRUD + validasi 0–100 + total share ditampilkan.
6. Security reinsurer: CRUD + **eksposur berjenjang**.
7. Business: CRUD + tampilan rate + `RIRATE` sebagai teks.
8. Terapkan ke semua: pratinjau + konfirmasi + audit + laporan sebagian-gagal.
9. **Kaskade hapus + popup konfirmasi** di seluruh tingkat.
10. Penegakan `o_message` di kelima jalur simpan + pembersihan HTML.
11. Laporan kontrak dengan total share ≠ 100%.
12. **Migrasi skema — `ready`** (OQ-001 ditutup): DDL kelima tabel + PK + FK `ON DELETE CASCADE` +
    kelima sequence.

**Keputusan yang masih terbuka, tidak memblokir:** pembagian paket domain di dalam `internal/` —
sama seperti empat spec sebelumnya; dan **anti-dobel logis** (lihat §Pertanyaan terbuka).

**Catatan sumber.** Spec ini bersandar pada korpus Pega `D:\XML\RNM_BRD\` (READ-ONLY), artefak di
`OUTPUT_HASIL_RNM\`, dan dua berkas `[data DBA]` yang diterima 2026-09-15. Sumber ADR tunggal:
`docs/adr/ADR-0001`…`ADR-0015`.
