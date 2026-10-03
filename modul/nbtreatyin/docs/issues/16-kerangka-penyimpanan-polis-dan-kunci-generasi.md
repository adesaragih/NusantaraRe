# 16: Kerangka penyimpanan polis dan kunci generasi

**Status:** ready-for-agent
**Blocked by:** — *(dapat mulai segera)*
**Menutup:** NB AC **1–7** *(7 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-5..ID-10

## Hasil & nilai pengguna

Polis treaty tersimpan sebagai baris di dalam tabel relasional, bukan sebagai dokumen teks.
Satu polis baru menempati **satu generasi**, dan generasi itu punya kunci yang mencegah dua berkas
menempati tempat yang sama.

⭐ **Ini tiket pertama rantai penyimpanan** — seluruh tiket berikutnya menulis ke kerangka ini.

## Yang dibangun

Kerangka tabel akar dan tabel generasi, berbagi kunci utama, beserta ketiga aturan keunikannya:

| Aturan | Yang dicegahnya |
| --- | --- |
| kunci alami `(NOPOLIS, PRODKE)` unik | dua berkas menempati generasi yang sama |
| `OLD_POLIS_ID` unik | satu generasi punya dua penerus — rantai menjadi pohon |
| baris generasi tertutup tidak dapat disunting | nilai yang sudah jadi dasar selisih berubah diam-diam |

⭐ Pada polis baru, `PRODKE = 0` dan `OLD_POLIS_ID` kosong. Tabel akar **lintas-lini** — dipakai
bersama modul lain, dan baris antar modul **tidak saling menunjuk**.

## Batas — yang TIDAK termasuk

⛔ Tipe dan presisi kolom — tiket **18**.
⛔ Isi tabel anak — tiket **19**.
⛔ Baris endorsemen `PRODKE ≥ 1` — tiket **24**.
⛔ `CREATE TABLE` dan presisi fisik — dicocokkan DBA **di dalam** tiket ini, bukan prasyaratnya.

## Cara mengujinya

Lewat seam `repository`. Menyimpan dua polis dengan kunci alami sama harus **ditolak oleh basis
data**, bukan oleh pemeriksaan di aplikasi — test memverifikasi penolakan itu datang dari lapisan
penyimpanan.

## Acceptance criteria

- [ ] **AC 1** — menyimpan dua polis dengan `NOPOLIS` dan `PRODKE` sama **ditolak**
- [ ] **AC 2** — polis baru tersimpan dengan `PRODKE = 0`
- [ ] **AC 3** — `OLD_POLIS_ID` pada polis baru bernilai kosong
- [ ] **AC 4** — dua baris tidak boleh berbagi `OLD_POLIS_ID` yang sama
- [ ] **AC 5** — tabel generasi dan tabel akar **berbagi kunci utama**, tanpa kolom kunci tamu terpisah
- [ ] **AC 6** — menyunting baris generasi yang sudah ditutup **ditolak**
- [ ] **AC 7** — baris antar modul pada tabel akar **tidak saling menunjuk**
