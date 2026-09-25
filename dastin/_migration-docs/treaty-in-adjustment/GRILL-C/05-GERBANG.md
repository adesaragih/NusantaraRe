> Modul  : Treaty In Adjustment · ronde C
> Dibuat : 2026-09-24
> Sifat  : kesiapan untuk to-spec. **Bukan pernyataan selesai.**
> Kriteria: `METODE` §7.2 — grilling selesai bila **semua bisa di-spec**, bukan bila semua terjawab

# Gerbang sesudah ronde C

## 1. Anggaran GRILL

| | Sebelum ronde C | **Sesudah** |
|---|---|---|
| Butir GRILL tersisa | C3, E3a, E3b, E3c, I2 | **nol** |
| — dijawab ronde ini | — | C3 (`GRL-14`), E3a (`GRL-15`), E3b (`GRL-16`), I2 residu (`GRL-17`) |
| — **larut**, bukan dijawab | — | **E3c** — `GRL-14` menghapus `ActualValue`, sehingga "ringkasan yang ditulis ke `ActualValue` lalu ditimpa" tidak punya tempat untuk terjadi. Nasib TDA-06: **diperbaiki** |
| — turun ke KONFIRMASI | — | **payung I2** — ADR-0042 butir 2 dan 3, ADR-0043, ADR-0054 (NC-06) |

**Anggaran tujuh butir yang ditetapkan di muka habis.** Itu bukan pernyataan bahwa grilling selesai.

## 2. Yang menghalangi pernyataan "siap to-spec"

Tiga, dan hanya satu di tangan sesi ini:

| # | Penahan | Pemilik | Dapat ditutup dengan bekerja? |
|---|---|---|---|
| 1 | **Tujuh TDA tanpa nasib**, cabangnya sudah ditutup — `03-LUBANG.md` §4 | sesi ini, **ronde D** | **ya** |
| 2 | **Paket `REV-1`…`REV-6`** belum diserahkan dan belum ditanggapi pemilik ADR | pemilik ADR induk | **tidak** |
| 3 | `PP-1` — siapa mengerjakan to-spec induk | pemilik proses | tidak memblokir; hanya arah dampak GRL-01 |

**Penahan 1 adalah alasan gerbang ini belum dibuka**, dan ia ditemukan di ronde ini, bukan diwarisi:
`METODE` §6.3 menuntut setiap TDA membawa nasib, dan TDA tanpa nasib **tidak bisa di-spec** — ia
sampai ke sesi to-spec ditugaskan kepada cabang yang tidak akan pernah bersidang.

## 3. Prasyarat cabang K, keadaannya sekarang

Dari `CABANG-K-PEMETAAN-TO-SPEC.md` §4:

| # | Prasyarat | Keadaan |
|---|---|---|
| 1 | ketujuh butir GRILL terkunci | **terpenuhi** — C1, C2, C3, E1, E3, I1, I2 |
| 2 | paket `REV` diserahkan dan ditanggapi | **belum** — satu-satunya di luar tangan sesi mana pun |
| 3 | `EXP-1` ditutup | **terpenuhi** — 26 September |
| 4 | setiap KONFIRMASI, REKOMENDASI, DITUNDA punya tempat | **hampir** — tujuh TDA di §2 butir 1 adalah sisanya |

## 4. Yang sudah dapat di-spec hari ini, dan tidak perlu menunggu

Dinyatakan supaya penahan di §2 tidak terbaca sebagai penahan seluruhnya:

- seluruh model versi dan identitas (GRL-09, GRL-10, GRL-11);
- daftar keadaan dan perpindahan untuk versi addendum (GRL-08);
- materialitas sebagai turunan (GRL-12) dan himpunan jenis (GRL-13);
- bentuk dan isi `NILAI_SELISIH`, termasuk kunci padanannya (E1 KONFIRMASI, `SEAM-ADJUSTMENT` §3);
- keempat putusan ronde ini.

Yang **tidak** dapat di-spec sampai penahan 1 tertutup adalah bagian-bagian yang bersentuhan dengan
ketujuh TDA itu — terutama picker (`TDA-11`), penomoran (`TDA-01`, `TDA-12`), dan pilihan jenis di
layar (`TDA-13`).

## 5. Titik periksa

`METODE` §7.3 menuntut titik periksa berkala. Titik periksa berikutnya: **sesudah ronde D**, dan
pertanyaannya satu — *"adakah TDA yang masih tanpa nasib, dan adakah cabang yang masih menunggu
tanpa ronde?"*

Pertanyaan itu dijawab dengan **mencocokkan dua tabel di indeks**, bukan dengan mengingat. Ia lolos
sekali; ia tidak boleh lolos dua kali.
