# 03: Keutuhan nomor urut antar generasi

**Status:** ready-for-agent
**Blocked by:** **02**
**Bergantung pada tiket NB:** **17** *(nomor urut baris anak)*
**Menutup:** AC **7–9** *(3 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-8 · ID-14 · ID-15

## Hasil & nilai pengguna

Generasi baru **wajib memuat setiap baris rincian** yang ada di generasi sebelumnya. Yang hilang
berarti endorsemen itu diam-diam menghapus komitmen yang masih berlaku — dan itu **ditolak**.

⭐ **Nilai lama tidak menjadi tabel.** Ia baris yang ditunjuk penunjuk generasi — dibaca, bukan
disalin ke tempat kedua.

## Yang dibangun

Aturan keutuhan di lapisan layanan: sebelum generasi baru disimpan, setiap nomor urut milik
generasi sebelumnya diperiksa ada. ⛔ Yang hilang **membatalkan seluruh penyimpanan**, bukan
menyimpan sebagian.

Ditambah penegasan bahwa selisih dihitung terhadap **generasi tepat sebelumnya**, bukan terhadap
polis asli — dan bahwa tidak ada tabel salinan nilai lama.

⚠️ **Susunan bersarang yang dibuat sistem lama — nilai lama di dalam nilai lama — tidak ditiru.**
`[terverifikasi]` Korpus **tidak pernah membacanya**.

## Batas — yang TIDAK termasuk

⛔ Perilaku nomor urut saat baris baru ditambahkan — tiket **04**.
⛔ Penanda pasangan bergeser — tiket **09**.

## Cara mengujinya

Lewat seam `repository`. ⭐ **Uji utama:** generasi kedua dibuat dengan satu baris rincian
dihilangkan — penyimpanan harus **ditolak seluruhnya**, dan tidak boleh ada satu baris pun
tertinggal.

## Acceptance criteria

- [ ] **AC 7** — selisih dihitung terhadap generasi **tepat sebelumnya**
- [ ] **AC 8** — generasi yang kehilangan salah satu nomor urut **ditolak**
- [ ] **AC 9** — **nol** tabel salinan nilai lama; ia dibaca lewat penunjuk generasi
