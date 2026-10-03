# Permintaan untuk tim inti — modul `nbtreatyin`

> Wilayah sunting implementasi ini hanya `modul/nbtreatyin/` dan
> `inti/backend/daftar/modul_nbtreatyin_gen.go`. Setiap perubahan di luar itu yang dibutuhkan modul
> ini ditulis di sini, beserta sebabnya. ⛔ Tidak satu pun dikerjakan sendiri.
>
> Cabang: `modul/nbtreatyin/implementasi` (dari `dev` 51d9b9cc). Tanggal: 2026-10-03.

## A · Uji di luar wilayah yang berubah merah karena modul ini menyala

| # | Berkas | Yang terjadi | Permintaan |
| ---: | --- | --- | --- |
| A1 | `frontend/daftar.menuTabel.test.ts` baris 90-92 | Uji gigit "modul frontend berbaris `DIMIGRASI='0'`" memakai **nbtreatyin** sebagai contoh modul yang belum dimigrasi. Modul ini sekarang menyalakan barisnya (slot menu `968`), sehingga hasilnya bukan lagi selisih yang diharapkan. | Ganti contohnya dengan modul yang barisnya masih `DIMIGRASI='0'`. |
| A2 | `modul/claimlife/backend/repository/batasanpemakaian_test.go:250` (`TestSetiapPemanggilBukaMemeriksaBolehDilewati`) | Sensus pemanggil `skemauji.Buka()` di seluruh repo: **17**, angka tertanam **16**. Pemanggil baru: `modul/nbtreatyin/backend/repository/polis_db_test.go` (uji repository lawan Oracle yang diminta brief; ia MEMERIKSA `BolehDilewati`). | Naikkan angkanya menjadi 17 (pemilik claimlife / tim inti). |

## B · Konstanta dan kepemilikan bersama

| # | Hal | Sebab | Permintaan |
| ---: | --- | --- | --- |
| B1 | `inti/backend/lini.go` hanya mengenal `LiniLife = "LIFE"` | Kasus NB Treaty In ditulis ke `T_WORK_POLIS.LINI`. Modul memakai konstanta lokal `models.LiniKasus = "NONLIFE"`, sejalan `KODE_PRODUKSI.TYPE = 'NONLIFE'` yang dibaca `GeneratePolicyNoTreaty_Act` (RDB `GetKodeProdNonLife_SQL`). | Tetapkan konstanta lini Non-Life bersama di `inti`; modul akan beralih memakainya. |
| B2 | `T_WORK_POLIS` + `SEQ_WORK_POLIS` milik **premiumlistlife** (migrasi 050/057/059) | Tabel kasus lintas-lini; migrasi 320 `nbtreatyin` membuat kunci tamu ke sana. Modul ini karenanya **membutuhkan premiumlistlife aktif** (MODUL_AKTIF) dan migrasinya berjalan lebih dulu (rentang 050 < 320). | Pertimbangkan memindah `T_WORK_POLIS` ke inti, atau nyatakan ketergantungan antar-modul di perakit. |
| B3 | `inti.Pelaku` tidak membawa nama tampilan | AC 40/42 menuntut `OperatorName`/`USERNAME` = **nama tampilan** (`OperatorID.pyUserName`). Modul membaca `M_LOGIN_GO.NAME` sendiri (pola modul `marketingofficer`). | Bila diinginkan satu sumber, sediakan nama tampilan di `inti` (mis. `Pelaku.Nama` dari sesi). |

## C · Sambungan dan data yang menunggu pihak lain

| # | Hal | Pihak | Catatan |
| ---: | --- | --- | --- |
| C1 | Sambungan konversi Arasapas treaty (`serviceInsertArasapas_act`) | tim inti + `[pemilik export Pega]` | Muatan 4 medan dibangun (KEPUTUSAN-RONDE-12 butir 7). Kunci `M_LINK_SERVICE` rule treaty tidak terekspor; pemanggilan layanan luar menuntut persetujuan (pola `outbox.ErrArasapasBelumDisetujui`). Pengirim bawaan modul GAGAL TERANG di produksi, dilewati di luar produksi. Antre-ulang (`outbox.Antrean`) dan email kegagalan belum dipasang — keduanya `BelumDiputuskan` di inti. |
| C2 | `M_NBTRIN_PERAN_TEMPAT` (migrasi 330) | `[IAM]` + work owner | Nol baris dari migrasi. Isian per tempat (saat ini satu tempat terpakai: `LISTSUGGEST_PRODUCTIONDATE`) berikut ARAH `MUNCUL`/`KECUALI` — tiket 05. |
| C3 | Tipe kolom view `TREATYINDETAILJOINEDM`, tabel `TREATYINDETAIL`, kolom `AGENT.STS_PKP` / `STATUSACTIVE` | `[DBA]` | Repository membaca tipe dari `SYS.ALL_TAB_COLUMNS` saat berjalan; nama kolom AGENT diambil dari nama properti RD (belum dikonfirmasi DBA). |
| C4 | `make test-db` | operator lingkungan | Uji `-tags db` modul ini MELEWATI tanpa `ORACLE_SCHEMA` + `ORACLE_SKEMA_UJI=true`; belum pernah dijalankan lawan skema uji. |

## D · Berkas di luar wilayah yang termodifikasi oleh pihak lain

`package.json` dan `package-lock.json` sudah termodifikasi di salinan kerja sebelum implementasi ini
dimulai (vite `^8.3.1` → `^7.3.6`, `@vitejs/plugin-react` `^4.7.0`). ⛔ Tidak di-commit dan tidak
dikembalikan oleh implementasi ini.
