# ✅ SUDAH DIJAWAB — berapa desimal untuk persen SHARE

> **Dijawab pemilik proses 4 Oktober 2026: DELAPAN desimal.**
> Keputusannya, beserta sebab penolakan kedua ujung dan syarat pembalikannya, di
> [`KEPUTUSAN-PENYELARASAN-REPO.md` §13](KEPUTUSAN-PENYELARASAN-REPO.md).
>
> ⚠️ **Baca juga §13.1.** Satu dari dua alasan penolakan 2 desimal ternyata **tidak
> terpenuhi oleh 8** pula: `PctTotal` `99.999999999999900` tetap tampil `100%` sampai 12
> desimal, dan baru terlihat pada 13 — yang melampaui presisi penyimpanan. Diukur sesudah
> keputusan diterapkan, dan ditagih balik.
>
> ⛔ **Berkas ini TIDAK dihapus.** Ukuran di bawah yang membuat jawabannya dapat diperiksa
> ulang; membuangnya meninggalkan keputusan tanpa dasar yang dapat diadu.

---

## Pertanyaan asli — berapa desimal untuk persen SHARE

**Diajukan kepada pemilik proses 3 Oktober 2026.**

⛔ **Berkas ini TIDAK memuat keputusan.** Dua aturan yang sama-sama berlaku saling bertentangan
pada satu golongan nilai, dan yang memutuskan bukan saya.

---

## Bentroknya

**Keputusan pemilik proses, 3 Oktober 2026:** persentase tampil **maksimal 2 desimal**.

**Komentar `inti/frontend/lib/format.ts` pada `formatPersen`**, dengan alasan bernomor:

> *"Jumlah desimalnya TIDAK dikurangi paksa: memaksa dua desimal menampilkan 33,333% sebagai
> '33,33%', dan tiga share 33,333 berjumlah TEPAT 100 sementara tiga share 33,33 tidak
> (SD-05, BR-01)."*

Keduanya benar untuk nilai yang berbeda, dan **tidak dapat berlaku bersamaan untuk persen SHARE** —
yaitu persen yang beberapa barisnya **dijumlahkan dan harus menghasilkan 100**.

---

## Siapa yang terkena

Diukur atas `POOLDATA.M_TREATY_IN2` (7.281 baris) dan ke-1.854 dokumen:

| Nilai | Sumber | Terisi | Rentang | Desimal terpanjang |
| --- | --- | ---: | --- | ---: |
| `CESSIONPCT` | `M_TREATY_IN2` | 3.065 | 0–100 | — |
| `RNMSHARE` | `M_TREATY_IN2` | 7.281 | 0–100 | — |
| `QSOR` | `M_TREATY_IN2` | 7.281 | 0–40 | — |
| `QSRI` | `M_TREATY_IN2` | 7.281 | 0–80 | — |
| `BROKERAGEPERCENTP` | `M_TREATY_IN2` | 7.281 | 0–12,5 | — |
| `RNMShareP` | dokumen | 1.221 | — | **30** (`2.825601535925207120348922139444`) |
| `InstallmentPct` | dokumen | 3.029 | — | 16 (`33.333333333333300`) |
| `PctTotal` | dokumen | 796 | — | 15 (`99.999999999999900`) |
| `PctLimit` | dokumen | 702 | — | 0 |
| `CoInShare` | dokumen | 702 | ⚠️ **PITA**, bukan angka: `>=30% up to < 50%` | — |

⚠️ **`PctTotal` terbesar adalah `99.999999999999900`.** Pada 2 desimal ia tampil `100%` — yang
benar sebagai ringkasan dan **salah sebagai bukti**: nilai itu memang bukan 100, dan satu-satunya
tempat perbedaannya terlihat adalah di layar.

⚠️ **`CoInShare` bukan angka sama sekali** — ia pita. Aturan desimal apa pun tidak berlaku padanya,
dan ia tampil apa adanya di kedua skenario.

---

## Yang berlaku SEMENTARA, sampai dijawab

| Golongan | Perlakuan hari ini | Di kode |
| --- | --- | --- |
| **persen share** — kelima kolom `M_TREATY_IN2` di atas, plus `PctLimit` | **desimal apa adanya** (`DESIMAL_TAK_DIBATASI`) | `JenisAngka = 'persenShare'` |
| **persen lain** — `ADJ_RATE`, `ROL`, `MDP_RATIO` | **2 desimal** | `JenisAngka = 'persen'` |
| **uang** | **4 desimal** | `JenisAngka = 'uang'` |

Pilihan sementara ini mempertahankan sifat *"jumlahnya tepat 100"*, yang bila hilang **tidak dapat
dipulihkan dari layar**. Sebaliknya, memangkas ke 2 desimal kelak hanya mengubah satu argumen.

⛔ **`ADJ_RATE` dan `ROL` sengaja BUKAN share**, dan itu diukur: keduanya melampaui 100
(`ADJ_RATE` sampai 109,6 · `ROL` sampai 199,4), jadi keduanya **laju**, bukan bagian dari sebuah
jumlah yang harus 100.

---

## Yang ditanyakan

1. **Untuk persen SHARE, mana yang menang** — 2 desimal, atau desimal apa adanya?
2. Bila 2 desimal: apakah `PctTotal` `99.999999999999900` boleh tampil `100%`, dan di mana
   selisihnya tetap dapat dilihat orang yang memeriksa?
3. Bila apa adanya: apakah `2,825601535925207120348922139444%` (30 desimal, nyata di 1 kontrak)
   dapat diterima di dalam sel grid, atau perlu batas atas yang lebih longgar — misalnya 8, yang
   sama dengan batas penyimpanan?

Jawabannya dicatat di `KEPUTUSAN-PENYELARASAN-REPO.md` beserta syarat pembalikannya, dan
mengubahnya adalah mengubah satu argumen di `labels.ts` (`DESIMAL_PERSEN` / golongan
`persenShare`).
