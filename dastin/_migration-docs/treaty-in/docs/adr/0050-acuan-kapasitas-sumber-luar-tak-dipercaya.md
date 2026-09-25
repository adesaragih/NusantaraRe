# ADR-0050 — Acuan kapasitas adalah sumber luar yang tidak dipercaya

**Status:** diterima **tanpa verifikasi**, 23 September 2026
**Berlaku untuk:** Treaty In

## Konteks

Pembagian kapasitas NuRe diambil dari sebuah tabel acuan. Uji D dan E — siapa pemeliharanya, dan
apakah periode susunan bisa tumpang tindih — **tidak dijalankan**.

## Keputusan

Treaty In **tidak memiliki** tabel itu. Ia dibaca lewat **adapter**, dan diperlakukan sebagai
sumber luar yang tidak dipercaya:

1. Bertipe tegas di perbatasan; nilai yang tidak terurai menyebabkan **kegagalan keras**, bukan
   nilai pengganti (ADR-0035).
2. Bila **lebih dari satu susunan cocok** untuk satu tanggal, **adapter gagal**. Tidak memilih yang
   pertama, tidak memilih yang terbaru.

> Ambiguitas adalah keadaan yang **dilaporkan**, bukan yang diselesaikan diam-diam.

## Dasar — sebagian artefak, sebagian penalaran

**Dari artefak:** tidak ada satu pun aturan di seluruh ekspor Pega yang menulis ke tabel itu —
Treaty In murni pembaca. Bentuk kolomnya adalah bentuk kontrak berperiode, bukan lookup sederhana.
Seluruh kolomnya bertipe teks bebas, tanpa primary key dan tanpa indeks.

**Dari penalaran:** ketiadaan kunci dan indeks berarti tumpang tindih periode **memang mungkin**
terjadi, dan tidak ada apa pun yang mencegahnya. Adapter yang memilih diam-diam akan menghasilkan
kapasitas yang berbeda-beda tanpa ada yang tahu mana yang benar.

## Syarat pembalikan

Pemelihara tabel ditemukan **dan** menyanggupi constraint yang menjamin ketunggalan. Saat itu
adapter boleh dilonggarkan menjadi pemilihan yang dinyatakan.

**Bila terbalik:** aturan pemilihan susunan masuk ke model. Ini salah satu dari enam keputusan yang
mengubah bentuk model.
