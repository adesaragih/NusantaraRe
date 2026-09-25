# ADR-0038 — Aturan yang bisa berubah disimpan sebagai data

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In

## Konteks

Tiga kali dalam sesi yang sama, jawaban atas pertanyaan yang berbeda berbentuk sama:

1. **Perutean persetujuan** — sistem lama menanam nama orang di dalam aturan, dan memakai kolom
   nomor telepon sebagai kode peran.
2. **Kontrak baca ke acuan luar** — bentuk yang dipakai adapter untuk membaca susunan baku NuRe.
3. **Kelengkapan transisi** — apa yang harus sudah ada untuk berpindah dari satu keadaan ke
   keadaan berikutnya.

Ketiganya berubah tanpa mengubah arti benda yang diaturnya.

## Keputusan

Aturan yang dapat berubah **tanpa mengubah arti benda yang diaturnya** disimpan sebagai **data**,
bukan ditanam di dalam kode.

Termasuk di dalamnya: perutean dan jenjang persetujuan, batas wewenang, daftar kelengkapan per
transisi, dan pemetaan adapter ke acuan luar.

## Yang TIDAK termasuk

Invarian — aturan yang menjaga catatan tetap bermakna — **bukan** data. Ia bagian dari arti
bendanya dan ditegakkan oleh model. Contoh: sebuah layer tidak boleh punya limit tanpa mata uang;
periode berakhir tidak boleh mendahului periode mulai; potongan tidak boleh melebihi bruto yang
menjadi dasarnya; persentase penyebaran harus berjumlah sama dengan bagian NuRe.

Perbedaan keduanya:

| | Invarian | Kelengkapan transisi |
|---|---|---|
| Berlaku | kapan saja data disimpan | hanya pada perpindahan keadaan tertentu |
| Melanggarnya berarti | catatan tidak berarti apa-apa | catatan sah tapi belum siap maju |
| Disimpan sebagai | model | data |

## Konsekuensi

- Daftar kelengkapan per transisi harus bisa diubah orang bisnis tanpa rilis.
- Penugasan peran bertanggal (ADR-0044) adalah data, dan harus menyimpan **sejak kapan** sebuah
  peran melekat pada seseorang — bukan hanya keadaan hari ini.
