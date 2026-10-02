# Papan tiket modul `treatyinadjustment` — empat dari 13

Papan **lengkap** modul ini ada di
`D:\XML_NURE\_migration-docs\treaty-in-adjustment\5-tiket\issues\README.md`. Folder ini memuat tiket
yang sudah disentuh saja, disalin apa adanya supaya kode dan tiketnya ditinjau bersama.

| # | Tiket | Keadaan |
| --- | --- | --- |
| `01` | Versi baru lahir dari versi berlaku terakhir, rujukan dasarnya disimpan eksplisit | lapisan skema selesai (migrasi `440`); **keempat penolakannya menunggu jalur simpan** |
| `02` | Materialitas mengunci ruas yang boleh disunting | kolomnya ada sejak `401`; **dua nilainya belum dinyatakan** + penegakan belum |
| `03` | Tanggal berlaku ditolak bila di luar periode kontrak | kolomnya ada sejak `401`; **bawaannya belum ada** + penegakan belum |
| `05` | *Perluas* — kolom nomor urut versi berdiri berdampingan | selesai (migrasi `441`) |

⛔ **`02` dan `03` BELUM selesai, dan sisanya BUKAN hanya penegakan.** Kolomnya memang ada sejak
migrasi `401` modul `treatyin` — `KAMUS-KOLOM.md` satu model untuk kedua modul — tetapi daftar
periksanya menuntut dua hal lagi yang bersifat **skema**: *"`SIFAT_MATERIAL_ADDENDUM` berdiri **dengan
dua nilainya**"* (`02`) dan *"atribut tanggal berlaku berdiri … **dengan bawaan tanggal mulai
kontrak**"* (`03`). Keduanya belum ada. Menandai tiket ini selesai karena kolomnya ada akan mengulang
persis cacat yang `TDA-10` catat: *"materialitas hanya ditegakkan di layar; nol penegakan di sisi
simpan."* Uraian dan dua pertanyaan yang menutupnya di
[`../KEPUTUSAN-TIKET-02-03.md`](../KEPUTUSAN-TIKET-02-03.md).

## Dua pertentangan yang diadili, bukan disamarkan

| Pertentangan | Putusan | Di mana tertulis |
| --- | --- | --- |
| Tiket `05` menuntut `NOMOR_URUT_VERSI` **boleh kosong**; `KAMUS-KOLOM.md` dan tiket `14` menetapkannya `NOT NULL` | `05` menang, **sebagai langkah migrasi tersendiri** (`441`) — bukan dengan menyunting berkas tiket `14`, supaya jejak pertentangannya tidak hilang | `backend/migrations/441_*.sql` |
| Tiket `03` menuntut **trigger** `INV-54`; ADR-0056 (K-4) melarang trigger pembawa aturan bisnis | ditegakkan di **services**, sejalan dengan `INV-53` yang sebentuk | [`../KEPUTUSAN-TIKET-02-03.md`](../KEPUTUSAN-TIKET-02-03.md) §2 |

## Yang terbuka berikutnya

Dengan `01` dan `05` mendarat, `06` (dari `01`) dan `10` (dari `05`) tidak lagi tertahan tiket.
Tiga tiket tertahan **pertanyaan bisnis yang belum dikirim**: `04` ← `DB-16a`, `07` ← `DB-20`,
`08` ← `DB-16b`.

Dan satu tepi melintas ke papan induk: tiket **`40`** Treaty In (*nilai versi sebelumnya
berdampingan, di-SELECT bukan disalin*) diblokir `01` — dan `01` kini selesai.
