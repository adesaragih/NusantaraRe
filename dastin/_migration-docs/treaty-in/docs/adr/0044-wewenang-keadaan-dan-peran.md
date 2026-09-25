# ADR-0044 — Wewenang: keadaan dan peran, terpisah tegas

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In

## Konteks

Di seluruh ekspor Pega modul ini, **tidak ada satu pun nama hak akses yang terisi** — mekanisme
peran bawaan Pega hadir 2.202 kali sebagai tempat kosong dan tidak pernah dipakai sekali pun.

Yang benar-benar mengatur siapa boleh apa hanya dua, dan keduanya bukan peran: sebuah bendera pada
kontraknya yang diperiksa sekitar 660 kali dengan delapan ejaan berbeda, dan nama orang yang
ditanam langsung di dalam aturan persetujuan.

Jadi yang kita lihat bukan keputusan bahwa wewenang milik kontrak. Kita melihat sistem yang tidak
pernah punya cara menyatakan "siapa", sehingga pertanyaan itu tidak pernah bisa diajukan.

## Keputusan

> **Keadaan** menjawab: apakah tindakan ini mungkin sekarang.
> **Peran** menjawab: apakah orang ini boleh melakukannya.
> Setiap larangan harus punya **tepat satu** sebab.

**Uji rancangannya konkret:** setiap penolakan harus bisa menyebutkan satu alasan saja. Bila sebuah
penolakan harus berbunyi "karena kontraknya sudah disetujui **atau** karena Anda bukan direktur",
rancangannya gagal — dua mekanisme sudah saling meniru. Ini uji rancangan, bukan pedoman menulis
pesan.

## Konsekuensi

- **Entitas peran dan penugasan bertanggal masuk gelombang 1.** Penugasan menyimpan sejak kapan
  sebuah peran melekat pada seseorang, bukan hanya keadaan hari ini — karena jejak perubahan harus
  mencatat peran yang berlaku **saat itu**.
- Cakupannya bukan hanya persetujuan, melainkan seluruh tindakan.
- Nama orang tidak pernah muncul di dalam aturan.
