# PROMPT — SERAGAMKAN KOLOM **`T_WORK_POLIS`** DENGAN **`T_WORK_CLAIM`** *(sesi baru, folder `D:\XML\RNM_BRD\OUTPUT_HASIL_RNM`, cabang **`dev`**)*

> Keputusan work owner 01-10-2026: *"samakan kolom antara T_WORK_POLIS dan T_WORK_CLAIM: CREATE_OP, CREATE_OP_NAME, TGL_UPDATE, TGL_CREATE. Ikuti saran
> asisten. Nama sama antara T_WORK_POLIS dan T_WORK_CLAIM. COVER_KEY tetap ada, jangan dihapus, di dua tabel itu."*
>
> ⛔ **Jalankan SEBELUM `PROMPT-IMPLEMENTASI-TIGA-MODUL-LIFE-GELOMBANG-2.md`**, atau tunggu sampai sesi itu selesai — modul Endorsement di gelombang itu
> menulis kasus polis, dan harus memakai nama kolom yang baru. Commit dengan jalur eksplisit `git commit -o -- <jalur>`. **Nol `git push`, nol
> `git pull`, nol `-migrate`** *(dijalankan work owner)*.

## 0. KEADAAN DEV *(asisten 01-10-2026, katalog baca-saja)*

| # | `T_WORK_CLAIM` | Tipe | `T_WORK_POLIS` | Tipe |
| ---: | --- | --- | --- | --- |
| 1 | `ID` | VARCHAR2(32) NOT NULL | `ID` | VARCHAR2(32) NOT NULL |
| 2 | `COVER_KEY` | VARCHAR2(32), self-FK | — | — |
| 3 | `LINI` | VARCHAR2(16) | `LINI` | VARCHAR2(255) |
| 4 | `PY_POSITION` | VARCHAR2(64) | `POSITION` | VARCHAR2(255) |
| 5 | `STATUS_WORK` | VARCHAR2(32) | `STATUS` | VARCHAR2(255) |
| 6 | `CREATE_OP` | VARCHAR2(64) | — | — |
| 7 | `CREATE_OP_NAME` | VARCHAR2(128) | — | — |
| 8 | `TGL_CREATE` | DATE | — | — |
| 9 | `TGL_UPDATE` | DATE | — | — |
| — | `SENDTO_ADMIN`, `SENDTO_MEDICAL`, `TYPE`, `CASE_ID`, `TAHAP` *(khusus klaim)* | | `FLAG_ONGOING_POLICY` *(khusus polis)* | VARCHAR2(1) |

## 1. ARTI KOLOM MENURUT XML — dasar penamaan

| Kolom | Properti Pega | Bukti | Keputusan |
| --- | --- | --- | --- |
| klaim `STATUS_WORK` | `pyWorkStatus` *(`Resolved-Completed`)* | `migrations/017_kolom_status_work.sql` | — |
| polis `STATUS` | **`pyWorkStatus`** *(`Input Offer Life`, `Input Premium Detail`, `Input Premium Summary`, `Resolved-Rejected`, `Resolved-Completed`)* | `repository/polis_work.go` ralat 28-09-2026 | **sama artinya → diganti nama menjadi `STATUS_WORK`** |
| klaim `PY_POSITION` | **`pyPosition`** — peran, mis. `ReasLifeMedicalAdvisor` | `handlers/tahap.go` | — |
| polis `POSITION` | **`Position`** — layar `Offer`/`Premium`, properti **berbeda** | `Activity/ProtectAccept.xml` b1207, b2288 | **tidak diganti nama** — menyamakannya dengan `PY_POSITION` akan menyatakan dua properti Pega berbeda sebagai satu |
| `COVER_KEY` | penunjuk kasus induk *(klaim: kasus Komite `KMTLF-` → klaim)* | `migrations/001_t_work_claim.sql` butir d | **ditambahkan ke polis**, tetap di klaim. XML PremiumList dan Endorsement **tidak** memakai penunjuk induk *(nol `pxCoverInsKey`)*, jadi kolomnya **tetap kosong** sampai ada modul yang terbukti mengisinya |
| `FLAG_ONGOING_POLICY` | `curWorkPage.FlagOnGoingPolicy` di kelas kasus `ASM-FW-GISFW-WORK-LIFE` | `CreateInputLife` b618, `IsFlagOnGoingPolicy` | **tetap di `T_WORK_POLIS`** — sifat kasus, bukan data dokumen |

## 2. BENTUK TUJUAN `T_WORK_POLIS`

| # | Kolom | Tipe | Kosong | Perubahan |
| ---: | --- | --- | --- | --- |
| 1 | `ID` | VARCHAR2(32) | tidak | — |
| 2 | `LINI` | VARCHAR2(255) | ya | — |
| 3 | `POSITION` | VARCHAR2(255) | ya | — |
| 4 | **`STATUS_WORK`** | VARCHAR2(255) | ya | **ganti nama** dari `STATUS`; panjang **tidak** diperkecil *(memperkecil dapat gagal pada data panjang)* |
| 5 | `FLAG_ONGOING_POLICY` | VARCHAR2(1) | ya | — |
| 6 | **`COVER_KEY`** | VARCHAR2(32) | ya | **baru**; `REFERENCES T_WORK_POLIS(ID)`, **tanpa** `ON DELETE` — sama persis dengan klaim butir d |
| 7 | **`CREATE_OP`** | VARCHAR2(64) | ya | **baru** |
| 8 | **`CREATE_OP_NAME`** | VARCHAR2(128) | ya | **baru** |
| 9 | **`TGL_CREATE`** | DATE | ya | **baru** |
| 10 | **`TGL_UPDATE`** | DATE | ya | **baru** |

`T_WORK_CLAIM` **tidak diubah**.

## 3. MIGRASI `059_seragam_kolom_t_work_polis.sql` *(+ `_down`)* di `modul/premiumlistlife/backend/migrations/`

| Langkah | Isi |
| --- | --- |
| 1 | `ALTER TABLE {skema}.T_WORK_POLIS RENAME COLUMN STATUS TO STATUS_WORK` |
| 2 | `ALTER TABLE {skema}.T_WORK_POLIS ADD (COVER_KEY VARCHAR2(32), CREATE_OP VARCHAR2(64), CREATE_OP_NAME VARCHAR2(128), TGL_CREATE DATE, TGL_UPDATE DATE)` |
| 3 | FK `COVER_KEY` → `T_WORK_POLIS(ID)`, nama constraint bergaya klaim, tanpa `ON DELETE` |
| 4 | isi baris lama dari `T_PREMIUM_LIST` **bila** penghubungnya terbukti *(kemungkinan `T_PREMIUM_LIST.ID_PEGA = T_WORK_POLIS.ID` — pastikan dari kode `polis_kasus.go`/`polis_nomor.go` dan sebut buktinya)*: `CREATE_OP_NAME` dari `CREATE_OP_NAME`, `TGL_CREATE` dari `TGL_INPUT`, `TGL_UPDATE` dari `TGL_UPDATE` bila ada; hanya baris dengan **tepat satu** pasangan. `CREATE_OP` baris lama **tetap kosong** — `T_PREMIUM_LIST` tidak menyimpan akun pembuat, dan mengarang nilainya dilarang *(ADR-U-0027)* |

Jalur mundur: buang FK, buang lima kolom baru, `RENAME COLUMN STATUS_WORK TO STATUS`. Nama berkas `050`–`058` **tidak diubah**. Setiap langkah dilindungi
pemeriksaan katalog supaya aman diulang. Nol `COMMIT`. Penjaga kata cadangan Oracle hijau. Penjaga rentang: `059` di rentang PremiumList 050–099.

## 4. KODE `modul/premiumlistlife/`

| Bagian | Ubah |
| --- | --- |
| Semua SQL `T_WORK_POLIS` | `STATUS` → `STATUS_WORK` *(`polis_work.go`, `polis_inbox.go`, `polis_kasus.go`, dan setiap tempat lain — cari semua)*; JSON API dan teks layar **tidak berubah** |
| Buat kasus *(`Input Offer`, `Input Premium`)* | isi `CREATE_OP` = akun pelaku, `CREATE_OP_NAME` = akun pelaku *(sama dengan Claim Life `pendaftaran.go` sampai login menyediakan nama tampilan)*, `TGL_CREATE` = `SYSDATE`, `TGL_UPDATE` = `SYSDATE`, `COVER_KEY` = NULL |
| Setiap ubah baris kasus | `TGL_UPDATE = SYSDATE` di pernyataan yang sama |
| Urutan inbox | **tidak berubah** — tetap mengikuti RD XML |
| Model dan uji | model kasus polis bertambah empat medan + `CoverKey`; uji murni SQL dan uji `db` *(SKIP tanpa DSN, `-p 1`)*; uji struktur kolom DDL = katalog |

## 5. ATURAN UNTUK MODUL BERIKUTNYA

Tulis di `docs/bersama/PANDUAN-TIM-PER-MODUL.md` dan `APP_RNM/modul/_templat/MODUL.md`: **setiap tabel kerja `T_WORK_<…>` baru memuat `ID`, `COVER_KEY`,
`LINI`, `STATUS_WORK`, `CREATE_OP`, `CREATE_OP_NAME`, `TGL_CREATE`, `TGL_UPDATE` dengan nama dan tipe seperti `T_WORK_CLAIM`**; kolom posisi dinamai
menurut properti Pega yang disimpannya *(`PY_POSITION` untuk `pyPosition`)*.

## 6. BATAS FOLDER

Boleh: `APP_RNM/modul/premiumlistlife/**`, `docs/bersama/PANDUAN-TIM-PER-MODUL.md`, `APP_RNM/modul/_templat/MODUL.md`, dan `APP_RNM/PANDUAN-MENJALANKAN.txt`
*(pengingat `-migrate` 059)*. Dilarang: `inti/`, modul lain, `T_WORK_CLAIM`, konfigurasi root.

## 7. URUTAN — satu commit per paket

| # | Paket | Commit |
| ---: | --- | --- |
| 1 | migrasi 059 + uji struktur | `premiumlist-life: 059 - T_WORK_POLIS seragam dengan T_WORK_CLAIM (STATUS_WORK, COVER_KEY, pembuat, waktu)` |
| 2 | kode + uji | `premiumlist-life: kasus polis menulis pembuat dan waktu, STATUS_WORK` |
| 3 | aturan §5 + panduan | `docs: kolom wajib tabel kerja T_WORK_*` |

Setiap commit: `go build`, `go vet` dan `go test` *(dengan dan tanpa tag `db`)*, `gofmt`, `tsc`, `vitest`, `vite build`.

## 8. LAPORAN

Satu pesan: tabel **paket → commit → berkas → angka uji**; bukti penghubung `T_PREMIUM_LIST` ↔ `T_WORK_POLIS`; daftar setiap SQL yang berganti
`STATUS` → `STATUS_WORK`; bentuk akhir kedua tabel berdampingan; pengingat `-migrate` 059 untuk work owner; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 1 Oktober 2026 dari katalog DEV `T_WORK_POLIS` dan `T_WORK_CLAIM`, migrasi 001, 017, 051, 057, kode `polis_work.go`, `polis_inbox.go`,
`polis_kasus.go`, `pendaftaran.go`, `tahap.go`, dan pencarian `pxCoverInsKey` di korpus PremiumList Life dan Endorsement Life.*
