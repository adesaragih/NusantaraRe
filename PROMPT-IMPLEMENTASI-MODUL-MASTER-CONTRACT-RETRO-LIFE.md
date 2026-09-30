# PROMPT — IMPLEMENTASI MODUL **MASTER CONTRACT RETRO LIFE** *(sesi baru, folder `OUTPUT_HASIL_RNM`, cabang **`dev`**)*

> Permintaan work owner 30-09-2026: *"sekarang aku mau migrasi yang Treaty Contract Retro Life"* — di korpus namanya **Master Contract Retro Life**
> *(folder `D:\XML\RNM_BRD\Master Contract Retro Life\`)*. Satu sesi, satu modul, satu folder.
>
> **XML adalah patokan dasar. Tiket adalah hasil grilling.** Bila tiket tidak sesuai XML atau katalog DEV, **perbaiki tiketnya** *(catatan bertanggal,
> kalimat lama dikutip)* — jangan diam-diam menyimpang. Jangan membuat menu, layar, tombol, atau aksi yang **tidak ada** di Pega.
> **Setiap pembacaan activity mencetak `pyStepsBlockName`** *(`//` = langkah ter-remark, tidak pernah jalan)*; baca WHEN dari urutan aksi; precondition
> `WhenTrue 2` = lanjut, `3` = lewati; cari **penulis** sebuah kolom, bukan hanya pembacanya.

## 0. BATAS KERJA

| Hal | Isi |
| --- | --- |
| Cabang | **`dev`**. Commit dengan jalur eksplisit `git commit -o -- <jalur>`. **Nol `git push`, nol `git pull --rebase`** |
| Folder yang boleh disunting | **`APP_RNM/modul/mastercontractretrolife/**`** saja, ditambah `inti/backend/daftar/modul_mastercontractretrolife_gen.go` hasil `go generate` |
| Dilarang | menyunting `inti/`, modul lain, konfigurasi root *(go.mod, package.json, vite.config.ts, tsconfig.json)*. Pola dari `modul/treatycontractout` **boleh ditiru**, **tidak boleh diimpor**. Kebutuhan fungsi bersama → catat di laporan, jangan tambahkan ke `inti` |
| Nama | backend `mastercontractretrolife`, frontend dan dokumen di folder modul yang sama; nilai `MODUL_AKTIF` `mastercontractretrolife` |
| Nomor | rentang migrasi **100–139**, slot menu **958–959** *(`MODUL.md`)* |
| Oracle | skema `POOLDATA`. `-migrate` dijalankan **work owner**. Uji tag `db` tanpa `ORACLE_DSN` = SKIP; **POOLDATA bukan target uji `db`**. Nol penulisan ke DEV oleh sesi ini |
| Rahasia | nol nama orang, nomor polis, kredensial, host/IP, alias tnsnames di berkas apa pun. Fixture berawalan `UJI-` |
| Prosedur Oracle | **tidak dipanggil**; logikanya ditiru di Go *(keputusan o)*. Nol `COMMIT` di teks SQL |

## 1. SUMBER

| Sumber | Letak |
| --- | --- |
| Korpus Pega *(READ-ONLY)* | `Harness` 4 *(InboxRetroLifeReinsurersList, InboxRetroLimitReinsurers, InboxSecurityReinsurerLife, InboxBusinessLifeReinsurers)*; `Section` 8; `FlowAction` 2 *(ViewRate, ViewRateTable)*; `DataTransform` 1; `RDBList` 12; `ReportDefinition` 12; `Activity` 27 |
| Bukan sumber | folder korpus `Claude outputs\` *(xlsx dan diff hasil sesi lama)* — boleh dibaca sebagai petunjuk, **tidak** sebagai bukti |
| Hasil grilling | `modul/mastercontractretrolife/docs/`: `spec.md`, `TAMBAHAN-SPEC-ronde-2.md` *(T1–T9)*, `grilling-ronde-1*.md`, `grilling-ronde-2.md`, `issues/01`–`12`, `ddl-tables-from-dba.md`, `procedure-bodies-from-dba.md` |
| Dokumen tersegel | `grilling-ronde-*` — **tidak disunting**. Ralat ditulis di tiket, `spec.md` bagian ralat bertanggal, atau berkas `RALAT-DEV-30-09-2026.md` baru |

## 2. RALAT DARI KATALOG DEV *(diukur asisten 30-09-2026, agregat, baca-saja)* — **mengalahkan spec dan tiket bila bertentangan**

| # | Klaim di dokumen | Fakta DEV | Menjadi |
| ---: | --- | --- | --- |
| K1 | spec penyimpangan **7** dan **8**, `ddl-tables-from-dba.md`, tiket 12, T5: *"PK `ID` dan FK `ON DELETE CASCADE` **sudah dieksekusi di basis data**"*, *"basis data kini mengaskade sendiri"* | kelima tabel `TREATYYEAR_LIFE`, `TREATYCONTRACT_LIFE`, `TREATYREINSURER_LIFE`, `TREATYSECURITYREINSURER_LIFE`, `TREATYBUSINESS_LIFE` punya **nol** constraint P, R, maupun U | **tiket 12 dicabut: nol migrasi, nol DDL, nol tabel baru** — preseden **tco4** Treaty Contract Out `[keputusan asisten dari preseden tco4; veto work owner]`. Rentang 100–139 tetap kosong; penjaga modul: nol berkas migrasi |
| K2 | kaskade hapus mengandalkan FK DB *(T5)* | tidak ada FK | kaskade **di Go**, satu transaksi, anak lebih dulu, urutannya: (1) security, (2) reinsurer dan business, (3) kontrak, (4) tahun; **sesudah** popup konfirmasi *(penyimpangan 3 tetap)* |
| K3 | keunikan `ID` dijamin PK | tidak ada PK; data DEV: nol `ID` ganda, nol `ID` kosong | `ID` baru = `'1' \|\| LPAD(<seq>.NEXTVAL, 6, '0')` persis RDB/prosedur; sequence per tabel **dari XML** *(RDB `SaveMaster*_Life_SQL`, `SaveTreatyBusinessAll_Life_SQL`)* — DEV juga punya `TREATYSECURITY_LIFE_SEQ` dan `TREATY_REINSURER_LIFE_SEQ`; pakai yang dirujuk XML, sebut buktinya |
| K4 | penyimpangan **1** normalisasi: `REINSTYPEID` sekali, `TREATYYEAR` sekali | kolom salinan tetap ada di tabel warisan dan dibaca prosedur lain *(`PEGA_M_TREATYYEAR_LIFE`, `PEGA_M_TREATYBUSINESS_LIFE`, `PEGA_M_REINSURANCETYPE_LIFE`)* | skema **tidak** diubah; layanan **menulis semua salinan** dari induknya dalam transaksi yang sama, sehingga salinan selalu sama dengan induk; uji: salinan = induk |
| K5 | penyimpangan **2** selisih dihitung | kolom `IDR_SELISIH`, `USD_SELISIH` tetap ada dan dibaca | selisih **dihitung di Go lalu tetap ditulis** ke kolomnya |
| K6 | penyimpangan **5** dan tiket 08: jejak audit penerapan massal | nol tabel jejak *(K1)* | preseden **OQ-TCO-25**: `USERID` dan `TGLUPDATE` diisi layanan di setiap baris yang berubah, ditambah satu baris log server berisi cacah baris berubah *(tanpa nama orang)* |
| K7 | penyimpangan **4** + tiket 03: gerbang tahun *"aturannya dipertahankan"* | T1 ronde 2: gerbang itu **MATI** di Pega | premis keputusan 4 keliru → **OQ-MCRL-01** untuk work owner. Bawaan sampai dijawab: **ikut XML** *(gerbang tidak ditegakkan)*, tiket 03 ditandai `[menunggu OQ-MCRL-01]`, pelajaran GILIRAN-11 *(dua gerbang mati sempat ditegakkan)* |
| K8 | penyimpangan **6** `HASIL1`/`o_message` | prosedur tidak dipanggil *(keputusan o)* | galat Go sendiri wajib sampai ke layar berkata-kata; tidak ada yang ditelan |

Data DEV sebagai acuan uji manual: `TREATYYEAR_LIFE` 2 baris, `TREATYCONTRACT_LIFE` 9, `TREATYREINSURER_LIFE` 29, `TREATYSECURITYREINSURER_LIFE` 0,
`TREATYBUSINESS_LIFE` 16; `ID` tujuh digit berawalan 1; sequence berikutnya lebih besar dari angka urut di `ID` terbesar.

## 3. MENU — satu modul, satu menu

Satu tombol **"Master Contract Retro Life"** di bawah `GROUPMENU` **MASTER** *(baris tabel sudah ada di isi awal 900)*. Klik membuka **halaman awal**
= harness yang di Pega menjadi pintu masuk modul — tentukan dari XML dan sebut buktinya. Tiga harness lain dan `ViewRate`/`ViewRateTable` dibuka
**dari dalam** layar sesuai aksi Pega, **bukan** tombol menu.

Slot `958_menu_mastercontractretrolife.sql` hanya `UPDATE {skema}.M_NAV_MENU SET DIMIGRASI = '1', TGL_UBAH = SYSDATE WHERE KODE = 'mastercontractretrolife'`
*(+ `_down` ke `'0'`)*; **nol `INSERT`**. `menu.ts` menyatakan `HALAMAN_AWAL`. ⚠️ Bentuk ini mengikuti **`PROMPT-MENU-DATAR-PER-GROUPMENU.md`**: bila
commit brief itu belum ada di `dev` saat paket menu, **kerjakan semua paket lain dulu** dan laporkan paket menu tertunda — jangan menulis butir anak model lama.

## 4. PARITAS LEBIH DULU

`docs/PARITAS-LAYAR-DAN-AKSI.md`: setiap harness, section, tombol, dan aksi di XML → rule yang dipanggil *(activity + nomor baris `bNNN` +
`pyStepsBlockName`)* → RDB/RD → tabel/kolom → rute API → komponen React → status. **Nol baris tanpa bukti XML.** Aksi yang di XML ter-remark atau tidak
terjangkau ditandai, tidak dibangun. Sebelum menulis kode, cocokkan tiket 01–11 dengan tabel paritas; ralat tiket yang meleset.

## 5. PAKET — satu commit per paket, uji hijau di setiap commit

| # | Paket | Tiket |
| ---: | --- | --- |
| 0 | `PARITAS-LAYAR-DAN-AKSI.md` + `RALAT-DEV-30-09-2026.md` *(K1–K8)* + ralat bertanggal di tiket 03, 09, 12 dan spec §penyimpangan 1–8 + `OQ-MASTER-CONTRACT-RETRO-LIFE.md` | — |
| 1 | kerangka modul: `backend/modul.go`, models, repository baca kelima tabel + lookup *(RD `Browse*`)*, rute baca, penjaga modul *(nol migrasi, nol tabel baru, kata cadangan Oracle)* | 12 *(versi ralat)* |
| 2 | tahun treaty: grid + tambah/ubah + ID dari sequence | 01, 03 |
| 3 | kontrak dan batas proteksi *(Treaty Limit)*: grid, simpan, selisih dihitung | 02, 04 |
| 4 | reinsurer: share, total, pesan galat seperti `SetErrorMessageReinsurer` | 05, 11 |
| 5 | security reinsurer berjenjang *(eksposur = share × share induk)* | 06 |
| 6 | business + tampilan rate *(`ViewRate`, `ViewRateTable`)* + "terapkan ke semua" dengan pratinjau | 07, 08 |
| 7 | hapus berjenjang dengan popup konfirmasi, satu transaksi | 09 |
| 8 | penegakan galat di kelima jalur simpan | 10 |
| 9 | frontend semua layar di `modul/mastercontractretrolife/frontend/` *(label VERBATIM dari XML, kelas CSS berawalan modul di berkas CSS modul sendiri)* | 01–11 |
| 10 | menu *(§3)* + `MODUL.md` Status → dimigrasi | — |
| 11 | dokumen: status tiket, PARITAS, OQ, `MODUL.md` *(prefix rute API, kontrak)* | — |

Backend dan frontend boleh dikerjakan per paket bersama *(vertikal)*. Setiap commit: `go build`, `go vet` dan `go test` *(dengan dan tanpa tag `db`)*,
`gofmt`, `tsc`, `vitest`, `vite build`. Uji di pohon kerja utama, bukan ekspor sebagian.

## 6. BERHENTI BILA

- XML menunjukkan layar atau aksi yang tidak tercakup tiket dan tidak jelas perilakunya — catat sebagai OQ, lanjutkan paket lain.
- Sebuah aturan hanya bisa dibangun dengan DDL pada tabel warisan — catat sebagai OQ, **jangan** menulis migrasi.
- Paket menu sementara brief menu datar belum masuk `dev` *(§3)*.

## 7. LAPORAN

Satu pesan: tabel **paket → commit → tiket → bukti XML → angka uji**; ralat K1–K8 yang diterapkan; OQ baru; halaman awal beserta buktinya;
contoh `GET` tiap rute baca *(tanpa data orang)*; pengingat `-migrate` untuk slot 958; bab **TELEMETRI EKSEKUSI**.

## 8. OQ UNTUK WORK OWNER

| OQ | Pertanyaan | Bawaan sampai dijawab |
| --- | --- | --- |
| OQ-MCRL-01 | Gerbang tahun di Pega ternyata mati *(T1)*. Tetap ditegakkan sebagai penyimpangan sadar, atau ikut Pega? | ikut Pega, tidak ditegakkan |
| OQ-MCRL-02 | PK dan FK di kelima tabel warisan tidak ada di DEV, padahal dokumen DBA menyebutnya sudah dipasang. Apakah DDL itu dipasang di lingkungan lain, atau memang belum pernah? | nol DDL, kaskade dan keunikan di Go *(K1–K3)* |

---

*Disusun 30 September 2026 dari katalog DEV (objek, constraint, cacah baris, `ID` ganda/kosong, sequence, prosedur yang merujuk kelima tabel), daftar
rule korpus `Master Contract Retro Life`, `spec.md`, `TAMBAHAN-SPEC-ronde-2.md`, `procedure-bodies-from-dba.md`, tiket 01–12, dan preseden tco4/OQ-TCO-25
Treaty Contract Out.*
