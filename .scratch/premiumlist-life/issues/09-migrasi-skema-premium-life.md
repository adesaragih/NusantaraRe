# 09: Migrasi skema — `M_LIFE_PREMIUM_SUMMARY`, `M_LIFE_PREMIUM_DETAIL`, `JSON_POLIS`, `JSON_OFFER_LIFE`

**Status:** wontfix — **digantikan tiket 00** (premis salinan tabel lama dibatalkan); verifikasi 050–056 lawan STRUKTUR dikerjakan di `4559ff3`

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

## Hasil & nilai pengguna

Sebagai **tim migrasi**, saya ingin skema Oracle sistem baru menyimpan setiap nilai premium Life
dengan **presisi yang sama persis** seperti sistem berjalan, supaya tidak satu pun rupiah berubah saat
pindah dan rekonsiliasi tidak menemukan selisih yang tidak dapat dijelaskan. *(User story 39 di spec)*

## Area codebase

`migrations/` (DDL tabel target + index), `internal/repository` (pemetaan tipe kolom ↔ desimal
presisi arbitrer), skrip rekonsiliasi.

## Rule Pega sumber

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

## ADR terkait

**ADR-0003** (uang non-float; DDL `NUMBER` **tanpa presisi** → desimal presisi arbitrer di aplikasi),
**ADR-0006** (tabel sequence), **ADR-0009** (migrasi penuh, koeksistensi ditolak), **ADR-0010**
(lampiran tetap di Google Storage — **bukan** bagian migrasi tabel ini).

## Acceptance criteria

*(belum dapat difinalkan — menunggu OQ-001; disusun agar siap dijalankan begitu DDL turun)*

- [ ] DDL tabel target menyalin **tipe, presisi, PK, index, dan nullability** dari tabel sumber
      apa adanya; tidak ada kolom uang yang menjadi `FLOAT`/`BINARY_DOUBLE`. (**ADR-0003**) — belum: wontfix — premis salinan dibatalkan; tipe dirancang sendiri di tiket 00
- [ ] Nilai uang lama dibaca dan ditulis ulang **tanpa perubahan digit mana pun**; rekonsiliasi
      membandingkan nilai lama dan baru secara tepat, bukan dengan toleransi. — belum: wontfix — rekonsiliasi pindah ke tiket 00 dan belum dibangun
- [ ] Kolom `CLOB` (`DATA_JSON`, `JSONDATA`) pindah utuh, termasuk isi yang panjang. — belum: wontfix — CLOB JSON tidak dibawa (spec §12)
- [ ] Sequence (`M_LIFE_PREMIUM_SUMMARY_SEQ`, `M_LIFE_PREMIUM_DETAIL_SEQ`, `JSON_OFFER_SEQ`,
      `GENERATE_SEQUENCE_NUMBER`) dipindahkan dengan **nilai berjalan yang benar**, sehingga nomor
      pasca-migrasi tidak pernah bertabrakan dengan nomor lama. — belum: wontfix — tiket 00 memutuskan nol sequence; migrasi nilai penghitung belum ada
- [ ] Index yang menopang kueri hilir ada sejak hari pertama — khususnya
      `M_LIFE_PREMIUM_DETAIL(PL_NUMBER)`, yang dipakai Claim Life (tiket **08**). — belum: wontfix — tabel warisan milik DBA, di luar migrasi repo ini
- [ ] Migrasi dapat dijalankan ulang dengan aman dan punya jalur mundur yang diuji. — belum: wontfix — pindah ke tiket 00; jalur mundur belum diuji terhadap Oracle
- [ ] Skema uji yang dipakai seluruh tiket lain dibangun **dari DDL yang sama** — bukan dari tiruan
      yang ditulis terpisah. — belum: wontfix — skema uji belum memasang 050–056

## Blocker

🚧 **`needs-info` — OQ-001 (sisa) terbuka.** Pemilik: **DBA**.

**Yang sudah ada:** nama kolom keempat tabel (dari body procedure + `INSERT` `SaveMasterLPDet`).
**Yang belum ada:** **DDL fisik** — tipe, presisi, PK, index, nullability.

⚠️ Alasan ini benar-benar memblokir: **seluruh parameter `PEGA_M_LIFE_PREMIUM_SUMMARY` bertipe
`VARCHAR2`, termasuk kolom uang.** Karena itu **tipe kolom sebenarnya di tabel belum diketahui** —
apakah `NUMBER` seperti pada `OS_AKSEPTASI_KLAIM_LIFE`, atau memang `VARCHAR2`. Menebak di sini akan
menentukan presisi uang seluruh domain Life, dan salah tebak baru ketahuan saat rekonsiliasi.

**JANGAN paksa `ready`.** Pola sama dengan tiket **13** di konteks Claim — Life.

**Tidak memblokir tiket lain**: kolomnya sudah cukup untuk menulis spec dan seluruh tiket 01–07.

## Perintah verifikasi

```
go test ./internal/...
```

</details>


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

## Verifikasi 050–056 lawan STRUKTUR — 28-09-2026 (giliran 10)

Tiket ini tetap `wontfix` (premisnya — salinan tabel existing — dibatalkan). Yang diminta brief
GILIRAN-10 §1 baris 3 adalah **verifikasi** migrasi 050–056 lawan
`STRUKTUR-TABEL-PREMIUMLIST-LIFE.md`. Penjaga lama `TestKolomDDLCocokDenganStruktur` mencocokkan
**nama** kolom saja; tipe, nullability, FK, dan index belum dijaga siapa pun.

### Sensus — dua cara, jendelanya disebut

Jendela: tabel kolom di bawah judul `## T_*` STRUKTUR (tujuh tabel; `DOCUMENT_POLIS` dan
`M_TEMPUPLOADLIFE` bukan milik migrasi ini) lawan `CREATE TABLE` / `FOREIGN KEY` / `CREATE INDEX`
di `migrations/05[0-6]_*.sql` tanpa berkas `_down`.

| Cara | Alat | Hasil |
| --- | --- | --- |
| A | skrip Python sekali-jalan | 7 tabel · **219** kolom · tipe (`teks` 92, `angka desimal` 99, `DATE` 20, `bilangan bulat` 8) · 6 FK · **0 selisih** |
| B | `TestMigrasi050Sampai056TipeNullFKIndexSesuaiStruktur` (pengurai Go, jalan berbeda) | sama — **0 selisih** |

⚠️ **Instrumen cara A gagal dulu, dan disebut.** Versi pertama mencocokkan sel kosong `| |` dengan
pola dua-spasi, sehingga hanya **13** dari 219 kolom terbaca — dan melaporkan "0 masalah" yang
**palsu**. Ditemukan karena cacah kolomnya dicetak; pengurai dibetulkan (`\s*`) sebelum hasilnya
dipercaya. Versi kedua melaporkan 20 "selisih tipe" yang ternyata **kosakata**: STRUKTUR menulis
`DATE` apa adanya, bukan `tanggal`. Sesudah itu nol. Uji Go karena itu **menagih 7 / 219** sebagai
jawaban yang diketahui, supaya pengurai yang rusak tidak lulus hampa.

⚠️ **Dibuktikan menggigit**: `053` dimutasi sementara (`IDR NUMBER(38,8)` → `VARCHAR2(40)`) — uji
merah pada `T_PREMIUM_LIST_SPREADING.IDR`; berkas dipulihkan (`git checkout`, status bersih).

### Yang dijaga kini

- tipe: `teks` ↔ `VARCHAR2(n)`, `angka desimal` ↔ **tepat** `NUMBER(38,8)` (ADR-U-0003),
  `bilangan bulat` ↔ `NUMBER(n)`, `DATE` ↔ `DATE`;
- `nullable = tidak` ↔ `NOT NULL`;
- `FK` di STRUKTUR ↔ `FOREIGN KEY` di DDL, **dua arah**;
- setiap FK ber-index (kolom pertama index).

### Catatan

- `T_PREMIUM_LIST.ID` *"shared PK = `T_WORK_POLIS.ID`"* — tanpa `FOREIGN KEY`, dan STRUKTUR pun tidak
  menandainya FK. Sepakat; bukan selisih.
- Kolom tabel warisan (`M_LIFE_PREMIUM_DETAIL`) di luar cakupan: DDL-nya milik DBA; kolom yang ditulis
  dijaga pemetaan `SaveMasterLPDet` (tiket 05a) dan kontrak pembaca (tiket 08).

### Angka

Go **545 PASS · 0 FAIL** tingkat atas; vet (+`-tags db`), gofmt bersih · vitest **349** · tsc bersih.
