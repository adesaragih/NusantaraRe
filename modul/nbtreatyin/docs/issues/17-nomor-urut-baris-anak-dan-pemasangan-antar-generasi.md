# 17: Nomor urut baris anak dan pemasangan antar generasi

**Status:** sebagian *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ready-for-agent)*
**Blocked by:** **16**
**Menutup:** NB AC **8–11** *(4 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-11..ID-13

## Hasil & nilai pengguna

Setiap baris rincian — angsuran, penyebaran, lapisan — membawa nomor urut di dalam induknya.
Nomor itulah yang kelak memasangkan baris generasi baru dengan baris generasi lama, sehingga selisih
per baris berarti.

⚠️ **Kunci dagang sengaja TIDAK dipakai untuk memasangkan.** Ia turun pangkat menjadi
**pemeriksa** — bila pasangan menurut nomor urut ternyata berbeda kunci dagangnya, itu ditandai,
bukan dipakai memasangkan ulang.

## Yang dibangun

Kolom nomor urut pada **setiap** tabel anak, unik di dalam satu induk.

⭐ **Pada polis baru, baris boleh dihapus dan nomornya dinomori ulang rapat `1..n`** — aman karena
belum ada generasi untuk dipasangkan. ⚠️ Begitu generasi ditutup, nomornya **beku**.
⛔ Aturan berbeda berlaku di endorsemen — tiket **24**.

## Batas — yang TIDAK termasuk

⛔ Perilaku nomor urut di endorsemen — tiket **24**.
⛔ Penanda pasangan bergeser — tiket **28**.

## Cara mengujinya

Lewat seam `repository`. Uji utama: polis dengan tiga baris angsuran, baris kedua dihapus, lalu
dibaca kembali — nomor urutnya harus rapat, bukan berlubang.

## Acceptance criteria

- [x] **AC 8** — setiap tabel anak memiliki kolom nomor urut
- [ ] 🟡 **AC 9** — menghapus baris di tengah menghasilkan penomoran ulang yang rapat
- [x] **AC 10** — nomor urut unik di dalam satu induk
- [x] **AC 11** — pemasangan antar generasi memakai nomor urut, **bukan** kunci dagang

## ⭐ Putaran 2 — paket penyimpanan (03-10-2026)

Dasar: PROMPT-NB-TREATY-IN-PUTARAN-2 bab 0 butir 11–12, bab 2 K4/K16/K17; rincian kolom `docs/PERBANDINGAN-KOLOM-DIAGRAM.md`.

- `NOURUT` tetap ada di **setiap** tabel anak (diagram F28) — ditagih `TestTabelDanKolomMengikutiDiagramGrilling`.
  `T_POLIS_CEDING` kini unik per `(QUOTATION_ID, NOURUT)` (induknya `T_POLIS_QUOTATION`, diagram O39).
- `NOURUT` catatan usulan di tabel warisan `HISTORYAKSEPTASIPRODUCTION` = berikutnya per `IDPEGA`
  (tidak dinomori ulang — riwayat, bukan baris dokumen).
- Status tetap **sebagian** (AC 9 🟡: penomoran ulang lawan Oracle belum dijalankan, K11).
