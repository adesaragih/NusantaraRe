> Modul  : Treaty In Adjustment · **ronde C** — `ActualValue`, kemampuan mati, dan warisan
> Dibuat : 2026-09-24
> Peran  : penggrill menanyakan; AI grilling menjawab dengan bukti; putusan milik pemilik proses
> Masukan: `PENGETAHUAN.md` (AS-IS, menang atas berkas penggrill) · `GRILL-A/06-PUTUSAN.md` ·
>          `GRILL-B/06-PUTUSAN.md` · `KEPUTUSAN-GRILLING-ADJUSTMENT.md` · `PETA-SUMBER-INDUK.md` ·
>          `../METODE-GRILLING.md` · ADR induk 0037, 0040, 0042, 0043, 0049, 0053, 0054, 0055
> Status : **TERBUKA**
> Sifat  : grilling. Tidak ada spesifikasi, DDL, struct Go, komponen React, endpoint, atau tiket.

# Ronde C — lingkup

## 1. Butir yang diadili

Anggaran GRILL tersisa **empat** menurut `PEMILAHAN-SISA-GRILLING.md` dan §7 berkas penggrill:

| Butir | Isi | Keadaan masuk ronde |
|---|---|---|
| **C3** | arti bisnis `ActualValue` | terbuka; prasyaratnya C1 (GRL-12) dan C2 (GRL-13) **sudah terkunci** |
| **E3a** | pro rata | terbuka |
| **E3b** | share fakultatif | terbuka |
| **I2** | aturan baru atas baris warisan | terbuka; prasyarat bacanya **ADR-0043, dibaca lebih dulu** — lihat §4 |

**E3c** (ringkasan yang tidak pernah tersimpan, TDA-06) **tidak** diadili di ronde ini: jawabannya
ditentukan C3. Ia masuk ronde berikutnya, dan itu bukan penundaan melainkan urutan.

## 2. Bukti yang dipakai

Seluruh pembacaan ekspor lewat perkakas, **tidak satu berkas XML pun dibaca langsung** — 185 MB
berdua. Perkakas yang dijalankan ronde ini, beserta apa yang ditolaknya:

| Perkakas | Dijalankan atas | Ditolak | Titik buta yang ia nyatakan sendiri |
|---|---|---|---|
| `tools/tulis.py sapu` | `.IsProRate`, `.EDMEffective`, `.ProRatePercent` | **0** simpul salinan terbungkus; **0** berkas gagal urai | 32 langkah Java, 191 langkah SQL/REST |
| `tools/pre.py` | `TreatyEDMProRateCalculation`, `TreatyEDMDifferenceDeduction` | — | — |
| `tools/dt.py` | `TreatyCalculateProratePct` | — | — |
| `tools/cat56.py` | katalog khas Adjustment | — | **56 aturan**; penentu kolom sisi |
| pencarian teks pemanggil | `TreatyCalculateProratePct`, `TreatyInDifferenceFacShare`, `TreatyEDMCalculateDifference`, `TreatyInActualUpdateValueShare` | — | pencarian teks **kebal terhadap nama tag**, dan itu kekuatannya di sini |

## 3. Bukti yang sengaja TIDAK dipakai

- **Berkas `PENGETAHUAN-PENGGRILL-ADJUSTMENT.md`** dipakai hanya sebagai **peta**, tidak sebagai
  bukti. Ia menyatakan sendiri bahwa ia bukan sumber kebenaran. Setiap fakta yang dipakai di
  `01-TEMUAN.md` disapu ulang atas ekspor.
- **`TREATYINOFFER` dan jalur penerbitan hilir** — di luar lingkup ronde ini.
- **Data produksi** — tidak ada yang terjangkau; klaim "pernah terjadi" tidak dibuat.

## 4. Prasyarat baca yang ditagih `MA-10`, dan hasilnya

`KEPUTUSAN-GRILLING-ADJUSTMENT.md` menutup daftar sumber induk dengan kalimat: *"ADR-0043 belum
pernah saya baca dan namanya menjanjikan jawaban bagi I2. Ia dibaca **sebelum** I2 diajukan."*

**Ketiganya dibaca sebelum ronde ini disusun** — ADR-0042, ADR-0043, ADR-0054 — dan hasilnya
**mengubah bentuk I2**, bukan hanya mengisinya. Lihat `01-TEMUAN.md` NC-06.

## 5. Aturan berhenti ronde ini

Tiga, dan hanya tiga:

1. sebuah jawaban **membatalkan** GRL yang sudah terkunci — bukan menambah, membatalkan;
2. sebuah butir tidak dapat diputuskan tanpa wewenang pemilik proses, **dan** pemecahannya tidak
   menolong karena yang belum diputuskan duduk tepat di tengahnya;
3. sebuah fakta yang dibutuhkan ternyata **tidak terbaca dari ekspor**, dan tidak ada perkakas yang
   dapat membacanya.

> Selain ketiganya: kerjakan, catat, lanjut.

Butir yang menunggu **data produksi** atau **jawaban bisnis** tidak menghentikan apa pun: ia menjadi
`UA-` atau `DB-` bernomor dan rancangannya tetap diputuskan (`METODE` §7.2).

## 6. Satu tabrakan bentuk, dilaporkan alih-alih dipilih diam-diam

`METODE` §4.1 menuntut **satu pertanyaan per giliran**. Bentuk ronde yang dipakai hari ini
mengajukan **seluruh frontier sekaligus** — empat pertanyaan bernomor, masing-masing dengan
rekomendasi.

Keduanya milik pemilik proses, dan yang lebih baru dipakai. Perbedaannya dicatat di sini supaya
dapat dikembalikan dalam satu kalimat: bila bentuk satu-per-giliran yang dikehendaki, ronde ini
dibaca berurutan C3 → E3a → E3b → I2, dan sisanya ditahan.
