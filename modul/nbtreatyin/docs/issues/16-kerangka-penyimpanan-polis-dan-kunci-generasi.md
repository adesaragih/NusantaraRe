# 16: Kerangka penyimpanan polis dan kunci generasi

**Status:** sebagian *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ready-for-agent)*
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

- [ ] 🟡 **AC 1** — menyimpan dua polis dengan `NOPOLIS` dan `PRODKE` sama **ditolak**
- [x] **AC 2** — polis baru tersimpan dengan `PRODKE = 0`
- [x] **AC 3** — `OLD_POLIS_ID` pada polis baru bernilai kosong
- [x] **AC 4** — dua baris tidak boleh berbagi `OLD_POLIS_ID` yang sama
- [x] **AC 5** — tabel generasi dan tabel akar **berbagi kunci utama**, tanpa kolom kunci tamu terpisah
- [ ] 🟡 **AC 6** — menyunting baris generasi yang sudah ditutup **ditolak**
- [x] **AC 7** — baris antar modul pada tabel akar **tidak saling menunjuk**

## Catatan implementasi 2026-10-03

`TGL_TUTUP` tidak diisi saat realisasi NB selesai: generasi ditutup ketika generasi penerusnya lahir
(endorsemen, di luar NB). Penyuntingan kasus yang sudah `Resolved-*` ditolak layanan
(`ErrKasusTertutup`), dan `TGL_TUTUP` terisi ditolak repository (`ErrGenerasiTertutup`, AC 6).

## ⭐ Putaran 2 — paket penyimpanan (03-10-2026)

Dasar: PROMPT-NB-TREATY-IN-PUTARAN-2 bab 0 butir 11–12, bab 2 K4/K16/K17; rincian kolom `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`.

⛔ **RALAT** atas *Catatan implementasi 2026-10-03* di atas. Bunyi lama, dikutip: *"`TGL_TUTUP` tidak
diisi saat realisasi NB selesai … dan `TGL_TUTUP` terisi ditolak repository (`ErrGenerasiTertutup`, AC 6)."*

Bunyi baru: kolom `TGL_TUTUP` **dibuang** — tidak ada di diagram grilling dan ID-10 menegakkan pembekuan
*"di services, bukan di tabel"*. Generasi **tertutup = ada baris penerus yang `OLD_POLIS_ID`-nya menunjuk
generasi ini** (`repository.syaratTerbuka`: `NOT EXISTS (… s.OLD_POLIS_ID = g.ID)` di `UPDATE` polis, nomor
polis, posisi, dan `Keadaan`). `ErrGenerasiTertutup` tetap (AC 6). `UNIQUE (OLD_POLIS_ID)` dan
`UNIQUE (NOPOLIS, PRODKE)` (diagram F16) tetap; yang kedua berupa indeks unik fungsi `CASE` supaya draf tanpa
nomor tidak bentrok. Uji `-tags db` `TestGenerasiTertutupDitolakDanPembatalanUtuh` kini menutup generasi
dengan penerus (belum dijalankan — K11). Status tetap **sebagian** (AC 1, 6 🟡: penegakan basis data).
