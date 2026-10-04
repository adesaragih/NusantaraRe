# Struktur Tabel — Treaty In Adjustment

Acuan bentuk tabel untuk aplikasi Go. Dibuat 1 Oktober 2026, bersamaan migrasi `440`–`441`
(tiket `01` dan `05`).

## Modul ini tidak membuat satu tabel pun

Dan itu **bukan kelalaian**. Papan tiket modul ini dipisahkan dari papan Treaty In 25 September 2026,
tetapi **model datanya tetap satu**: modul ini tidak punya §10 sendiri, dan
`treaty-in/SPEC-MODEL-DATA.md` memerikan `KONTRAK` dan `VERSI_KONTRAK` untuk keduanya.

| Yang modul ini lakukan atas tabel milik `treatyin` | Migrasi | Tiket |
| --- | --- | --- |
| menambah kolom `VERSI_KONTRAK.ID_VERSI_KONTRAK_DASAR` + `FK_VERSI_KONTRAK_DASAR` + `IX_VERSI_KONTRAK_DASAR` | `440` | `01` |
| melonggarkan `VERSI_KONTRAK.NOMOR_URUT_VERSI` menjadi **boleh kosong** | `441` | `05` |

⛔ **Kolom `VERSI_KONTRAK` digambarkan di
[`../../treatyin/docs/STRUKTUR-TABEL-TREATY-IN.md`](../../treatyin/docs/STRUKTUR-TABEL-TREATY-IN.md),
dan SENGAJA tidak diulang di sini.** `TestDokumenSTRUKTURSepakatAtasTabelBersama` menuntut dua dokumen
yang menggambarkan tabel yang sama **sepakat kolom demi kolom**. Menyalin ke-49 barisnya ke sini
berarti menanam salinan yang wajib identik — dan salinan yang wajib identik adalah salinan yang akan
menyimpang, lalu bedanya baru terlihat saat satu sisi menulis kolom yang sisi lain tidak baca. Satu
dokumen, satu pemilik.

## Kolom yang modul ini bawa artinya

Ketiganya berdiri di DDL modul `treatyin`; yang **artinya** milik modul ini.

| Kolom | Tipe | Null | Arti | Tiket |
| --- | --- | --- | --- | --- |
| `ID_VERSI_KONTRAK_DASAR` | bilangan bulat | ya | rujukan eksplisit ke **versi berlaku terakhir** saat versi ini dibuat. Kosong = versi **pertama** kontrak itu — keadaan yang benar, bukan data yang hilang | `01` |
| `NOMOR_URUT_VERSI` | bilangan bulat | **ya** sejak `441` | kosong = baris warisan yang belum dinomori ulang. Bukan nol | `05` |
| `SIFAT_MATERIAL_ADDENDUM` | teks | ya | `MATERIAL` / `TIDAK_MATERIAL` — menentukan ruas mana yang boleh disunting | `02` |
| `TANGGAL_BERLAKU_ADDENDUM` | DATE | ya | sejak kapan versi berlaku; bawaannya tanggal mulai kontrak | `03` |
| `JENIS_ADDENDUM` | teks | ya | dari `EDMState` sistem lama | `07` |
| `ID_DOKUMEN_ADDENDUM` | bilangan bulat | ya | **belum berkunci asing** — `DOKUMEN_ADDENDUM` dibuat tiket `04`, yang tertahan `DB-16a` | `04` |

## Relasi yang modul ini pasang

| Anak | Kolom | Induk | ON DELETE | Invarian |
| --- | --- | --- | --- | --- |
| `VERSI_KONTRAK` | `ID_VERSI_KONTRAK_DASAR` | `VERSI_KONTRAK.ID_VERSI_KONTRAK` | **tanpa** — Oracle MENOLAK | INV-18 |

Kunci asing yang menunjuk tabelnya sendiri, dan itu sah di Oracle. **Ditetapkan sadar**: kolom ini
penunjuk ke baris **sejarah**, bukan kepemilikan, sehingga menghapus versi yang masih menjadi dasar
versi lain harus **gagal** — bukan diam-diam memutus rantainya.

## Yang BELUM ada, dan tiket mana membuatnya

| Yang belum ada | Tiket | Penahan |
| --- | --- | --- |
| `DOKUMEN_ADDENDUM` + FK `VERSI_KONTRAK.ID_DOKUMEN_ADDENDUM` | `04` | **`DB-16a`** — belum dikirim |
| `NILAI_SELISIH` berkunci bisnis | `06` | `01` *(selesai)* |
| Titik beku jenis dan materialitas | `07` | **`DB-20`** — belum dikirim |
| Tanggal berlaku pada **dokumen** | `08` | **`DB-16b`** — belum dikirim |
| Pengisian `NOMOR_URUT_VERSI` menurut kronologi | `10` | `05` *(selesai)* |
| `INV-69` dan `INV-70` + pemantau kebasiannya | `11` | `02` · `06` |
| Pencabutan pembacaan urutan dari pengenal | `12` | `10` |
