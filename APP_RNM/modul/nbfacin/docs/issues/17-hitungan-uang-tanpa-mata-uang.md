# 17: Hitungan nilai uang tanpa mata uang di pemuat `repository`

**What to build:** Saat pemuat `repository` NB membaca data, ia **menghitung** berapa nilai uang yang tiba
tanpa mata uang dan melaporkannya — tidak memberi mata uang bawaan, tidak berhenti.

Asal: tiket 02 ditutup work owner 01-10-2026 (butir 34, A9) karena digantikan butir 3 (ikuti
`inti/backend/uang` apa adanya). Dua kewajibannya tidak ikut gugur dan dipindah ke sini.
`[terverifikasi]` berkas kasus PA dan MBU nyata menyimpan `Currency` kosong pada coverage (tiket 02 bab
Comments); pengukuran produksi D1 mencatat 4 dari 691.925 baris tanpa mata uang (dikutip tiket 02).

**Blocked by:** 24 (pemuat `repository`) — ~~pemuat `repository` NB (menunggu tabel flat)~~

> ⚠️ **Ralat 02-10-2026** (register bab *Ralat kelima*): tiket ini **tidak** menunggu keputusan tabel flat — K-069,
> K-073, K-074 dan F-3 sudah menutupnya. Sebab salah baca: spec pemuatan (`04-spec\11-spec-pemuatan-data-lama.md`)
> tidak ikut dipindah ke `nbfacin/docs`. Yang berlaku untuk baris lama yang dimigrasikan: kolom pembawa mata uang
> **`DEFAULT 'UNKNOWN' NOT NULL`** (K-069, 24 kolom) — baris tanpa mata uang **terekam dan terhitung**, bukan ditolak;
> `panic` ADR-0006 hanya untuk **tulisan baru aplikasi**.

**Status:** needs-info — hitungan sudah ada di `loader.Flatten` (`Diagnostik.MataUangTakDiketahui`, tiket 22);
pelaporannya menunggu tiket 24

- [ ] Pemuat menghitung nilai uang tanpa mata uang per medan, dan hitungannya terlihat (log/laporan), tanpa default diam-diam
      — ✅ bagian hitung: `loader.Flatten` per `TABEL.KOLOM`, termasuk yang mewarisi `UNKNOWN` (tiket 22,
      `TestFlattenMataUang`); ⬜ bagian laporan: tiket 24. `UNKNOWN` adalah sentinel K-069 yang **bersuara**, bukan bawaan diam-diam
- [ ] Nilai tanpa mata uang tetap dibawa apa adanya (butir 3), seperti `TestPremiPAMataUangKosongTidakDiberiBawaan`
- [ ] Jalur tulis Oracle gagal keras bila diminta menulis uang tanpa mata uang — kewajiban jalur tulis produksi, tidak dihapus

## Comments
