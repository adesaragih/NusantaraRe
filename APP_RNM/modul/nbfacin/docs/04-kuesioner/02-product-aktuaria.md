# Kuesioner Product / Aktuaria — Satuan Rate

**Purpose:** memastikan satuan rate yang benar untuk tiga lini bisnis, di mana **label layar dan
rumus perhitungan saling bertentangan**. Jawaban Anda menentukan apakah label layar yang diperbaiki,
atau rumusnya — dan salah memilih menggeser premi satu ordo besaran 10 pada lini yang bersangkutan.

**From:** Tim Migrasi Facultative Inward · **To:** Product / Aktuaria ·
**How your answers will be used:** dicatat sebagai keputusan proyek, lalu menentukan tabel satuan
rate per lini bisnis di sistem baru.

---

> ## 🟡 TERJAWAB SEBAGIAN — 16 September 2026
>
> | Butir | Keadaan |
> | --- | --- |
> | **2** — tiga label bertentangan | ✅ **TERTUTUP** → **K-016** + **K-018**. Rumus yang menang di ketiganya; ketiga label diperbaiki, **nol angka berubah** |
> | **1** — satuan di slip | 🟡 **BELUM TERTUTUP** → menunggu **Pengukuran D2** (sampel nilai rate dari DBA) |
>
> **Mengapa butir 1 belum tertutup padahal sudah dijawab.** Jawabannya (ketiganya persen) **cocok
> untuk MBU** tetapi **bertentangan dengan rumus korpus untuk PA dan Layering**. Pertentangan itu
> sempat menjadi **Konflik K-1**, dan ditutup **K-018** — tetapi yang dikunci K-018 adalah satuan
> **yang disimpan dan dihitung**, bukan satuan **di slip**. Keduanya bisa berbeda dan tetap konsisten
> bila ada konversi di antaranya, sehingga butir 1 tetap terbuka sampai D2 masuk.
>
> ⚠️ **Ini bukan koreksi terhadap penjawab.** Lihat catatan di bawah butir 1.
>
> Naskah pertanyaan **tidak diubah** — hanya stub jawaban yang diisi.
> Penomoran: berkas ini memakai **butir 1 = slip · butir 2 = label**; ringkasan grilling memakai
> kode terbalik (**P1 = label · P2 = slip**). Isinya sama, hanya urutannya berbeda.

## Context

Aplikasi Facultative Inward sedang dipindahkan dari Pega ke sistem baru. Kami menemukan sesuatu yang
awalnya tampak seperti cacat, tetapi ternyata **memang begitu rancangannya**: properti rate yang
sama memakai **satuan berbeda tergantung lini bisnis**.

| Lini bisnis | Satuan rate menurut rumus |
| --- | --- |
| Fire, PA | **per mille (‰)** |
| Aneka, Bonding, Golf, Marine Cargo, MBU | **persen (%)** |

Ini terbukti kuat: satu berkas perhitungan memuat tiga lini berdampingan dengan dua satuan berbeda,
dan sebagian label layar bahkan mencetak tanda ‰ secara harfiah. Jadi kami **tidak** akan
menyeragamkan satuannya — itu justru akan merusak separuh portofolio.

Masalahnya, pada **tiga tempat** label layar dan rumus tidak sepakat. Karena yang menghitung uang
adalah rumusnya, bukan labelnya, angka yang keluar selama ini mengikuti rumus. Tetapi bila labelnya
yang benar, berarti ada lini yang selama ini salah hitung.

## How to answer

Tenggat: **23 September 2026**. Perkiraan waktu pengisian 15–20 menit.

Isi langsung di bawah tiap pertanyaan. Jawaban parsial dan "saya tidak tahu" tetap berguna — tandai
yang Anda ragukan daripada melewatinya.

⚠️ Kami juga meminta DBA mengambil **sampel nilai rate dari basis data**. Besaran angkanya akan
membuktikan satuannya secara langsung (rate 2,5 menunjukkan per mille; 0,25 menunjukkan persen).
Bila jawaban Anda dan data itu bertentangan, kami akan kembali kepada Anda — jadi tidak perlu
ragu menjawab berdasarkan praktik yang Anda ketahui.

---

## Satuan rate

### 1. Untuk **MBU, PA, dan Layering** — satuan rate mana yang dipakai di slip?

Ini pertanyaan intinya, dan paling cepat dijawab dari praktik sehari-hari.

Ketika underwriter menuliskan rate pada slip untuk ketiga lini di bawah, angka itu dalam satuan apa?

| Lini | Satuan yang dipakai di slip (isi: ‰ / % / lainnya) |
| --- | --- |
| **MBU** (kendaraan bermotor) | **%** |
| **PA** (personal accident) | **%** |
| **Layering** (struktur berlapis) | **%** |

_Why this matters: satu jawaban ini menutup ketiga pertentangan di pertanyaan berikutnya sekaligus.
Bila praktiknya jelas, kami tidak perlu membedah kode lebih jauh._

> **MBU, PA, dan Layering — ketiganya persen (%).**

🟡 **Status: belum tertutup — `[menunggu D2]`.**

**Bedakan dua hal yang mudah tertukar di sini.** Keduanya bisa **berbeda dan tetap benar**, asalkan
ada konversi di antaranya:

| | Apa yang dimaksud | Keadaan |
| --- | --- | :-: |
| **Satuan di slip / tampilan** | Angka yang ditulis underwriter di slip, yang ditanyakan butir ini | 🟡 belum dikonfirmasi |
| **Satuan yang disimpan & dihitung** | Angka yang masuk basis data dan dipakai rumus premi | ✅ **dikunci K-018** |

Untuk **satuan yang dihitung**, korpus sudah memutuskan: **MBU = %** (sejalan jawaban ini → K-016),
tetapi **PA dan Layering = ‰** (bertentangan dengan jawaban ini → K-018).

⚠️ **Pertentangan ini tidak dibaca sebagai kesalahan penjawab.** Dua pembacaan sama masuk akalnya,
dan kami belum bisa memilih: (1) slip PA/Layering memang memakai %, lalu dikonversi sebelum disimpan
— keduanya benar; atau (2) pertanyaan ini terjawab terfokus pada MBU, dan PA/Layering tidak ikut
ditinjau terpisah.

**Yang menutupnya: Pengukuran D2** (kuesioner DBA) — sampel besaran nilai `RATE` dari basis data
membuktikan satuan tersimpan tanpa perlu pendapat siapa pun. Bila ternyata slip dan penyimpanan
memang berbeda, **titik konversinya wajib ditemukan dan diport** — itu justru temuan penting, bukan
kesalahan.

### 2. Tiga tempat di mana **label layar bertentangan dengan rumusnya** — mana yang benar?

Untuk tiap baris: apakah **label** yang benar (berarti rumusnya cacat dan selama ini salah hitung),
atau **rumus** yang benar (berarti labelnya keliru dan hanya menyesatkan mata)?

| # | Layar | Label tertulis | Satuan menurut rumus | Mana yang benar? |
| ---: | --- | --- | --- | --- |
| a | Input spreading **MBU** | `Standard Rate (‰)` | **persen (%)** | |
| b | Tampilan coverage **PA** (layar endorsement) | `Rate (%)` | **per mille (‰)** | |
| c | Daftar **Layer** | `Rate (%)` | **per mille (‰)** | |

Catatan yang mungkin membantu: pada butir (a), layar input MBU yang lain justru berlabel `Rate (%)`
— sesuai rumusnya. Jadi ada dua layar MBU dengan label berbeda untuk properti yang sama. Pada butir
(b), layar **input** PA berlabel `(‰) Policy Rate` dan sesuai rumusnya; hanya layar **tampilan** yang
berbeda.

⚠️ Butir (b) kebetulan berada di layar tampilan siklus endorsement, tetapi **pertanyaannya tetap
relevan sekarang**: yang kami tanyakan adalah satuan rate PA itu sendiri, dan itu dipakai di
perhitungan New Business. Layar tempat labelnya muncul tidak mengubah jawabannya.

_Why this matters: bila rumus yang benar (dugaan kami untuk ketiganya), perbaikannya sekadar
mengganti teks label — tidak ada angka yang berubah. Bila label yang benar, berarti ada premi yang
selama ini meleset 10×, dan itu temuan yang jauh lebih besar dari sekadar migrasi._

> **Rumusnya yang benar; labelnya yang salah.**

✅ **TERTUTUP — dicatat sebagai K-016 (MBU) dan K-018 (ketiganya).** Jawaban ini diterima penuh:
**satuan mengikuti RUMUS, bukan label layar.** Ketiga label diperbaiki:

| # | Layar | Label lama | Label baru | Angka berubah? |
| ---: | --- | :-: | :-: | :-: |
| a | Input spreading **MBU** | `Standard Rate (‰)` | **`(%)`** | **tidak** |
| b | Tampilan coverage **PA** | `Rate (%)` | **`(‰)`** | **tidak** |
| c | Daftar **Layer** (`LayerListDtl`) | `Rate (%)` | **`(‰)`** | **tidak** |

Peta satuan yang dikunci **K-018**: **PA · Layering · FIRE = ‰** (pembagi 1.000) ·
**MBU · ANEKA · BONDING · GOLF · MARINE CARGO = %** (pembagi 100).

📌 **Perhatikan arah perbaikannya — dua di antaranya berlawanan dengan dugaan awal.** Butir (a)
bergerak ke **%**, tetapi butir (b) dan (c) bergerak ke **‰**. Jadi "rumusnya yang benar" **tidak**
berarti "semuanya persen": rumus PA dan Layer justru berskala ‰. Ini yang sempat menjadi
**Konflik K-1**.

**Ketiganya perbaikan label saja — nol angka berubah**, sehingga rekonsiliasi paralel run tidak
terganggu sama sekali.

🟡 **Pengukuran D2 tetap dijalankan sebagai konfirmasi silang, bukan penentu.** Satuannya sudah
ditetapkan korpus; sampel rate dari DBA hanya mengonfirmasi. Bila data produksi ternyata
bertentangan, itu **konflik baru yang dibuka tersendiri** — bukan alasan menunda implementasi
sekarang.

---

## Anything else?

Adakah lini bisnis lain yang satuan rate-nya Anda ketahui berbeda dari daftar di bagian Context?
Atau perubahan satuan yang pernah terjadi di masa lalu sehingga data lama dan data baru tidak
sebanding?

>
