# 00: Skema tujuh tabel + `T_WORK_POLIS` + migrasi — **PREFACTOR**

**Status:** sebagian terkerjakan — DDL + penjaga ada; `-migrate` dan rekonsiliasi menunggu skema uji dan persetujuan

**Blocked by:** CL-01 (kerangka aplikasi + seam API — scaffolding lintas konteks)

⚠️ **PREFACTOR dan tiket PERTAMA.** Diberi nomor `00` supaya berada di depan tanpa menomori ulang
sembilan tiket yang sudah terbit. **Seluruh tiket 01–09 kini memblokir pada tiket ini.**

Alasannya: bentuk penyimpanan berubah total — dokumen JSON + tabel flat warisan → **tujuh tabel
relasional empat tingkat**. Tidak ada irisan lain yang dapat berdiri sebelum bentuk barunya ada.
*"Make the change easy, then make the easy change."*

## Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin setiap atribut polis menjadi **kolom bernama** bertipe benar, dan
seluruh polis lama pindah **tanpa kehilangan satu nilai pun** — termasuk seluruh peserta, seluruh
rekap mata uang, seluruh spreading, dan seluruh riwayat penawaran. Sebagai **organisasi**, saya ingin
bentuk polis dapat diperiksa, dicari, dan divalidasi. *(Spec §12, US 24a–24f)*

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `migrations/` | DDL tujuh tabel + `T_WORK_POLIS` + sequence + **index pada setiap FK**; skrip migrasi & rekonsiliasi |
| `internal/models` | Agregat polis: header → {rekap mata uang, peserta → spreading → spreading retro, riwayat penawaran} |
| — | Skrip rekonsiliasi jumlah baris, nilai uang, dan tanggal |

## Bentuk baru — spec §12

```
T_PREMIUM_LIST (PK ID)
  ├─ T_PREMIUM_LIST_SUMMARY        1:N  FK PREMIUM_LIST_ID   ON DELETE CASCADE
  ├─ T_PREMIUM_LIST_DETAIL         1:N  FK PREMIUM_LIST_ID   ON DELETE CASCADE
  │     └─ T_PREMIUM_LIST_SPREADING       1:N  FK DETAIL_ID          ON DELETE CASCADE
  │            └─ T_PREMIUM_LIST_SPREADING_RETRO  1:N  FK SPREADING_ID  ON DELETE CASCADE
  └─ T_VIEW_SUGGEST                1:N  FK PREMIUM_LIST_ID   ON DELETE CASCADE
T_WORK_POLIS (PK ID)   ⬅ mandiri, LINTAS-LINI — bukan anak T_PREMIUM_LIST
```

Daftar kolom lengkap tiap tabel ada di **spec §12** dan
`.scratch/premiumlist-life/revisi-penyimpanan-premiumlist.md`.

## Sumber migrasi

| Sumber | Identitas | Peran |
| --- | --- | --- |
| `JSON_POLIS.DATA_JSON` (CLOB) | ditulis `PremiumList Life/RDBList/InsertJsonPolis.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTJSONPOLIS`) → `POOLDATA.INSERTJSONPOLISLIFE` | pohon polis penuh |
| `JSON_OFFER_LIFE` (CLOB + 24 kolom flat) | `RDBList/SaveOfferJsonLife_SQL.xml` (`ASM-FW-GISFW-INT-OFFERJSON` / `ASM!SAVEOFFERJSONLIFE_SQL`); dibaca `RDBList/GetOfferLife_sql.xml` (`ASM-FW-GISFW-WORK-LIFE` / `RNM!GETOFFERLIFE_SQL`) | penawaran + riwayat |
| `M_LIFE_PREMIUM_DETAIL` | `RDBList/SaveMasterLPDet.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `ASM!SAVEMASTERLPDET`) — **±81 kolom** | peserta |
| `M_LIFE_PREMIUM_SUMMARY` | `RDBList/InsertPLSummary.xml` (`ASM-FW-GISFW-INT-LIFE_PREMIUM_SUMMARY` / `ASM!INSERTPLSUMMARY`) → `POOLDATA.PEGA_M_LIFE_PREMIUM_SUMMARY` | rekap |
| `LIFEINPRODUCTION` | `RDBList/SaveLifeinProduction_SQL.xml` | produksi |

⚠️ `[terverifikasi]` `PEGA_M_LIFE_PREMIUM_SUMMARY` dipanggil **posisional** (31+ argumen); **nama
kolomnya tidak terbaca dari korpus**. Pemetaan kolom sumber → tujuan **wajib dicocokkan ke DDL
sebenarnya** saat tiket ini dikerjakan — jangan menebak dari urutan argumen.

⚠️ `[terverifikasi]` `GetOfferLife_sql` menyaring `WHERE STATUS = 'Bind' AND OLDID IS NULL`. Adanya
`OLDID` menandakan **versioning penawaran** — migrasi memutuskan apakah versi lama ikut dibawa;
temuannya **dilaporkan**, tidak ditebak.

## ADR terkait

**ADR-0003** (uang non-float), **ADR-0006** (identitas lewat sequence), **ADR-0009** (migrasi penuh).

## Acceptance criteria

- [ ] ⚠️ Skema polis **relasional penuh**: setiap atribut menjadi **kolom bernama**. Test yang
      menemukan kolom JSON menyimpan atribut polis **gagal**. *(AC 32 spec; penyimpangan sadar 1)*
- [ ] ⚠️ **Tujuh tabel** ada dengan PK sequence dan FK sesuai diagram; seluruh FK **`ON DELETE
      CASCADE`**. *(AC 46 spec)*
- [ ] ⚠️ FK spreading menunjuk **peserta**, dan FK spreading retro menunjuk **baris spreading** —
      bukan header. Test yang menemukan keduanya menggantung pada header **gagal**. *(AC 40 spec)*
- [ ] ⚠️ **`T_WORK_POLIS` ada sebagai tabel mandiri** — bukan anak `T_PREMIUM_LIST`; keadaan tangga
      **tidak** menjadi kolom header polis. Kolom minimum: identitas polis + lini, posisi/status
      tangga, audit. *(AC 45 spec; penyimpangan sadar 4)*
- [ ] ⚠️ **Setiap FK punya index.** Migrasi yang meninggalkan FK tanpa index **gagal**.
      *(AC 49 spec)*
- [ ] ⚠️ **`T_PREMIUM_LIST_DETAIL` memuat `PARENT_ID`** — FK self-reference ke `ID` peserta versi
      sebelumnya, **nullable**, ber-index. Ia dipakai konteks **Endorsement Life** untuk mencocokkan
      peserta lama↔baru. Ketiga tabel lain **tidak** memilikinya.
      *(`.scratch/endorsement-life/spec.md` §16, AC 60; `[keputusan work owner]`)*
- [ ] ⚠️ **`T_PREMIUM_LIST` memuat kolom EDM**, seluruhnya **nullable** dan kosong pada baris new
      business: `EDM_TYPE`, `OLD_POLICY_NO`, `EDM_DATE`, `EDM_NOTE`, `TYPE_CEDING`,
      `PREMI_PROPOSED`, `UANG_PERTANGGUNGAN`, `SUM_INSURED`, `JENIS_PRODUK`, `SISTEM_REASURANSI`,
      `STATUSS`, `STATUS_UPDATE`, `STATUS_SERVICE`, `START_DATE`, `END_DATE`, `NOENDORS`,
      `PL_NUMBER_EDM`, `PRODKE`, `EDMSTATUS`, `STATUSOLD`.
      *(`.scratch/endorsement-life/spec.md` §16, AC 57; `[keputusan work owner]`)*
- [ ] Versi berjalan sebuah polis dapat ditemukan sebagai baris ber-**`PRODKE` terbesar**; seluruh
      versi **hidup berdampingan**. *(Endorsement §16, AC 56)*
- [ ] ⚠️ **Properti bawaan Pega tidak menjadi kolom** — tidak ada `px*`, `py*`, `pz*`; single-page
      kosong (`Policy`, `Quotation`, `TempError`) juga tidak. *(AC 47 spec)*
- [ ] `T_PREMIUM_LIST_SUMMARY` memuat **rekap uang penuh per mata uang**, bukan hanya kode mata
      uang. *(AC 37 spec; penyimpangan sadar 2)*
- [ ] `T_PREMIUM_LIST_DETAIL` memuat **kedua jendela valuasi** — `GROSS_VALUATION_*` **dan**
      `RETROCESSION_VALUATION_*` — beserta `EFFECTIVE_DATE`, `LAPSE_DATE`, `PERIOD_MM`.
      *(AC 38 spec)*
- [ ] `FACTOR` bertipe **desimal**; nilai berdesimal tujuh angka pindah **tanpa berubah**.
      *(AC 39 spec; **ADR-0003**)*
- [ ] ⚠️ Seluruh uang dan share bertipe **desimal presisi arbitrer**; seluruh tanggal **`DATE`**;
      seluruh kolom **nullable**; identitas dari **sequence**. Test yang menemukan kolom uang
      bertipe teks atau melewati `float` **gagal**. *(AC 48 spec; penyimpangan sadar 5)*
- [ ] Seluruh polis lama pindah **tanpa kehilangan satu nilai pun**; jumlah baris per polis —
      peserta, rekap, spreading, riwayat — **sama** sebelum dan sesudah. *(AC 50 spec; **ADR-0009**)*
- [ ] Nilai uang pindah **tanpa berubah satu digit pun**; rekonsiliasi membandingkan **secara
      tepat**, bukan dengan toleransi. *(AC 50 spec; **ADR-0003**)*
- [ ] Tanggal yang berupa teks menjadi `DATE` **tanpa pergeseran zona waktu**; yang **tidak dapat
      diurai dilaporkan**, bukan didiamkan.
- [ ] Polis tanpa retrosesi pindah dengan **nol baris** spreading dan spreading retro — **bukan**
      kegagalan. *(AC 42 spec)*
- [ ] ⚠️ Migrasi **tidak mereplikasi** `@ASM.GetPageJSONString()` dan tidak menulis satu pun CLOB
      JSON. *(AC 34 spec; penyimpangan sadar 1)*
- [ ] Sequence polis pindah dengan **nilai berjalan yang benar**; penomoran **tidak melompat dan
      tidak mengulang** setelah migrasi. *(**ADR-0006**)*
- [ ] Migrasi dapat **dijalankan ulang dengan aman** dan punya **jalur mundur yang diuji**.

## Blocker

**Tidak ada pemblokir.** **OQ-001 ditutup 2026-09-16** `[keputusan work owner]` — ketujuh tabel
dirancang sendiri; presisi fisik dicocokkan DBA **di dalam tiket ini**, bukan sebagai prasyarat.

## Catatan

⚠️ **Skala: jutaan baris.** `[keputusan work owner]` `T_PREMIUM_LIST_DETAIL`, `_SPREADING`, dan
`_SPREADING_RETRO` tumbuh sangat besar — satu polis grup dapat memuat ribuan peserta, masing-masing
dengan beberapa baris spreading dan retro. Konsekuensi yang **harus** ditangani di tiket ini:
**index pada setiap FK**, **partisi per periode/tahun** dipertimbangkan, dan migrasi dijalankan
**bertahap**, bukan satu transaksi raksasa.

⚠️ **Hapus fisik polis besar itu berat** — kaskade dapat menyentuh jutaan baris.
`[keputusan work owner]` **arsip/soft-delete dipertimbangkan sebagai jalur normal**; hapus fisik
disediakan tetapi bukan jalur sehari-hari. Keputusan finalnya diambil di tiket ini dan **dicatat**.

⚠️ **Nama berbohong (OQ-066).** `[terverifikasi]` `InsertJsonPolis`, `InsertJsonPolisEDM`, dan
`SaveOfferJsonLife_SQL` benar-benar menulis **blob JSON** — ketiganya dibuang. Sebaliknya di konteks
lain nama "Json" justru menandai INSERT flat. **Baca kodenya, jangan namanya.**

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — kaskade empat tingkat, presisi desimal,
konversi tanggal, dan rekonsiliasi jumlah baris **hanya berperilaku benar pada basis data
sungguhan**; memalsukannya berarti tidak menguji apa pun yang penting.

```
go test ./internal/...
cd frontend && npm test
make check
```

## Implementasi — 27 September 2026 (migrasi `050`–`056`)

### Yang dibuat

| Migrasi | Tabel | Kolom | FK | Index |
| --- | --- | ---: | ---: | ---: |
| `050` | `T_WORK_POLIS` | 4 | 0 | 0 |
| `051` | `T_PREMIUM_LIST` | 56 | 0 | 0 |
| `052` | `T_PREMIUM_LIST_DETAIL` | 82 | 2 | 2 |
| `053` | `T_PREMIUM_LIST_SPREADING` | 14 | 1 | 1 |
| `054` | `T_PREMIUM_LIST_SPREADING_RETRO` | 16 | 1 | 1 |
| `055` | `T_PREMIUM_LIST_SUMMARY` | 39 | 1 | 1 |
| `056` | `T_VIEW_SUGGEST` | 8 | 1 | 1 |

Masing-masing berjalur mundur (`*_down.sql`, `DROP … CASCADE CONSTRAINTS`).

⛔ **Nama kolomnya DIBANGKITKAN dari `STRUKTUR-TABEL-PREMIUMLIST-LIFE.md`, tidak diketik ulang.**
219 kolom yang diketik tangan adalah 219 kesempatan salah satu huruf, dan salah satu huruf di nama
kolom baru terlihat saat migrasi data gagal mencocokkan.

### Tipe fisik — keputusan kami, dan ia menunggu DBA

STRUKTUR menyebut **kategori logis** dan menyatakan presisi fisik `[data DBA]`. Pemetaannya:

| Kategori | Tipe | Dasar |
| --- | --- | --- |
| angka desimal | `NUMBER(38,8)` | uang dan share; **ADR-0003**, `revisi-penyimpanan-premiumlist.md` |
| bilangan bulat | `NUMBER(5)` | umur, periode, `PROD_KE`, nomor baris — **konvensi yang sudah ada di repo ini** |
| DATE | `DATE` | |
| teks | `VARCHAR2(255)`; `ID` dan ber-akhiran `_ID` → `VARCHAR2(32)` | sepola `T_CLAIMLF_*` |

⚠️ **Kedua arah salahnya tidak setara**, dan itu yang menentukan mana yang harus diperiksa lebih
dulu: `VARCHAR2` terlalu pendek **menolak** data yang sah — kegagalan yang terlihat — sedangkan
`NUMBER` terlalu pendek **membulatkan uang diam-diam**.

`NUMBER(5)` dipilih **bukan** karena selera: penjaga `TestNolNumberTanpaPresisi` yang sudah ada
hanya menerima `NUMBER(19)`, `NUMBER(38,8)`, dan `NUMBER(5)`. Ronde pertama memakai `NUMBER(10)`
dan penjaga itu menolaknya — benar, sebab tipe numerik keempat berarti konvensi keempat.

### Keputusan yang dicatat

⛔ **Tidak ada constraint FK antara `T_WORK_POLIS` dan `T_PREMIUM_LIST`.** Hubungannya 1:1 lewat
**shared PK**, dan STRUKTUR menyatakannya *"tanpa kolom penyambung"*; kolom `Kunci` untuk
`T_PREMIUM_LIST.ID` berisi `PK` saja, bukan `PK, FK`. Menambahkannya berarti memutuskan arah
ketergantungan — mana yang lahir lebih dulu — yang dokumen acuan tidak putuskan.

⛔ **Nol sequence.** Pengenalnya dirakit di `repository` mengikuti pola `PengenalWorkBerikut`
(butir **pl3**), bukan `DEFAULT seq.NEXTVAL` di DDL.

⛔ **Kolom audit `T_WORK_POLIS` sengaja tidak ada.** Tiket ini menyebut "audit" **tanpa menamainya**,
dan STRUKTUR menolak menuliskannya dengan alasan yang benar: menamainya sendiri berarti mengarang.

⛔ **`M_TEMPUPLOADLIFE` tidak dibuat.** Ia tabel warisan penampung unggahan CSV, ditulis
`RDBList/InsertDataUploadLife.xml` di sistem lama — **dibaca** tiket 04, bukan dimiliki. Ada penjaga
di kedua arah: yang satu melewatinya saat membandingkan kolom, yang lain **berbunyi bila DDL kelak
membuatnya** — sebab yang berubah saat itu adalah kepemilikan, dan itu keputusan work owner.

### Empat penjaga bersama yang harus dilebarkan — aditif, dan sebabnya

| Penjaga | Sebelumnya | Sesudahnya |
| --- | --- | --- |
| `TestKaskadeHanyaPadaEmpatRelasi` | memindai SEMUA migrasi | dibatasi migrasi `001`–`049` (Claim Life) |
| `TestKolomTakDibawaHanyaAdaDiKatalog` | memindai SEMUA migrasi | dibatasi Claim Life |
| `letakStruktur` | satu dokumen | **dua** dokumen, keduanya wajib |
| cacah `CREATE` | 33 (11 tabel) | 46 (18 tabel, 19 index) |

⚠️ **Kedua pelebaran pertama itu bukan pelonggaran, melainkan penyempitan lingkup** — dan keduanya
lahir dari cacat yang sama: penjaga modul Claim Life membentang ke tabel modul lain. `LAYER_1`..`4`
*"tidak punya rumah"* di tabel klaim dan itu benar; ia **punya** rumah di `T_PREMIUM_LIST`
*(STRUKTUR b98–101, bersumber `SaveLifeinProduction_SQL`)*. Penjaga yang menuduh hal yang benar akan
dilonggarkan orang, bukan dipatuhi.

### Penjaga BARU untuk pohon polis, tiap satunya dibuat gagal lebih dulu

| Penjaga | Dibuat gagal dengan | Berbunyi |
| --- | --- | --- |
| `TestSeluruhFKPohonPolisBerkaskade` | cabut `ON DELETE CASCADE` dari 053 | menyebut berkas dan AC-nya |
| `TestSetiapFKPohonPolisBerindex` | hapus `CREATE INDEX` 053 | menyebut kolom FK-nya |
| `TestKolomDDLCocokDenganStruktur` | sisipkan `KOLOM_KARANGAN` | *"tidak diminta STRUKTUR"* |
| `TestTabelBukanMilikKitaTidakDibuat` | buat `M_TEMPUPLOADLIFE` | menyebut alasan kepemilikannya |

### ⛔ Yang BELUM dijalankan

**`-migrate` belum pernah dijalankan**, dan pengecekan tabrakan nama di katalog belum dilakukan —
keduanya menuntut Oracle nyata dan **persetujuan manusia**. Berkas migrasinya ada, jalur mundurnya
ada, dan seluruh penjaga bentuknya hijau; yang tersisa adalah menjalankannya terhadap **skema uji**.

**AC yang tercentang giliran ini:** tujuh tabel + PK + FK `ON DELETE CASCADE`; FK spreading menunjuk
peserta dan FK retro menunjuk spreading; `T_WORK_POLIS` mandiri; setiap FK ber-index; `PARENT_ID`
nullable ber-index hanya di `_DETAIL`; kolom EDM nullable; nol kolom JSON; nol properti Pega; rekap
uang penuh; kedua jendela valuasi; seluruh uang desimal, tanggal `DATE`, kolom nullable.
**Belum:** yang menuntut basis data sungguhan — rekonsiliasi, migrasi data, jalur mundur teruji.


## ⛔ RALAT 28 September 2026 — migrasi `056` GAGAL di DEV (pl6)

Work owner menjalankan `-migrate` 28-09-2026 pukul **09.33**. Hasilnya: `017`–`020`,
`030`, `050`–`055` **terpasang**; **`056_t_view_suggest` GAGAL**, dan
`T_VIEW_SUGGEST` tidak ada di DEV.

**Sebabnya bukan penanda `{{skema}}`** — sebelas langkah lain di jalan yang sama lolos dengan
penanda itu. Sebabnya kolom **`INITIAL`**, yang merupakan **kata cadangan Oracle**. Bukti yang
dijalankan work owner:

```
SELECT 1 AS INITIAL   FROM DUAL  -> ORA-00923
SELECT 1 AS "INITIAL" FROM DUAL  -> lolos (berkutip)
SELECT 1 AS NO        FROM DUAL  -> lolos
```

Nol tabel warisan memakai nama itu *(katalog DEV: nol kolom `INITIAL`, nol tabel
`%SUGGEST%`)*; sumbernya properti `.Initial` *(`Activity/AddHistorySuggest.xml`)*.

### pl6 — kolomnya diganti nama

`INITIAL` → **`INITIAL_SUGGEST`**, mengikuti pola saudaranya `DATE_SUGGEST`,
`PIC_SUGGEST`, `COMMENT_SUGGEST`.

⚠️ Berkas `056` **disunting di tempat**, bukan ditambah `057 ALTER RENAME`. Ia belum
pernah terpasang di mana pun — `T_MIGRASI` DEV tanpa `056`, skema uji belum dibuat — dan
`RENAME` untuk kolom yang belum pernah ada berarti **mewariskan riwayat yang tidak terjadi**.

### ⭐ Penjaga yang lahir dari kegagalan ini

`TestNolKataCadanganOracleSebagaiKolom`. Yang membuat `INITIAL` lolos **bukan
kecerobohan melainkan ketiadaan pemeriksa**: nol uji di repositori ini pernah menanyakan apakah
sebuah nama kolom boleh berdiri telanjang di Oracle. Kini **373 definisi kolom** diperiksa
terhadap **47 kata cadangan** setiap kali uji berjalan.

⛔ Ia memakai pengurai **produksi** `KolomCreateTable`, bukan regex kedua: pengurai kedua
adalah definisi kedua tentang *"apa itu kolom"*, dan yang kedua akan diam-diam berbeda.

⚠️ Daftarnya **sengaja tidak lengkap**, dan itu dinyatakan di berkasnya. Oracle punya ratusan kata
cadangan; yang dijaga hanya yang masuk akal muncul sebagai nama kolom di domain reasuransi. Daftar
yang berpura-pura lengkap lebih berbahaya daripada yang mengaku parsial — yang pertama membuat
orang berhenti berpikir.

Dibuktikan merah dengan menanam kembali `INITIAL` yang asli: penjaga menyebut berkas, kolom,
kode galat Oracle, dan nama penggantinya.

**Sesudah cabang ini menyatu ke `main`**, work owner menjalankan `-migrate` lagi; hanya
`056` yang tersisa dijalankan.
