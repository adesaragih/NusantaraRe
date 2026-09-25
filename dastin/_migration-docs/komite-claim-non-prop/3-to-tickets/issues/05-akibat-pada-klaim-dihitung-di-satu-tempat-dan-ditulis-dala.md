---
status: tertahan
---

# 05: Akibat pada klaim dihitung di satu tempat dan ditulis dalam transaksi yang sama

*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Nasib klaim dihitung dari tiga masukan — jenis sirkulasi, bersyarat, dan
hasil — oleh satu fungsi yang ditulis sekali, lalu ditulis ke klaim di dalam transaksi
keputusan yang sama. Pembaca melihat "komite menyetujui usulan penutupan" sebagai **klaim
ditutup**, tidak pernah sebagai "klaim disetujui".

Fungsi akibat berikut ketujuh barisnya; tulis-balik ke sisi Klaim; model baca yang menyajikan
nilai tampil dari fungsi yang sama.

**Persyaratan:** `S-031`, `S-032`, `S-033`, `S-034`, `S-035`, `S-046`, `S-049`.

**Tidak termasuk:** nomor akseptasi — tiket `08`. Efek ke luar — tiket `11`.

**Jalur gagal:** kombinasi yang tidak tercakup ketujuh baris → `422` sebagai galat
konfigurasi, keputusan tidak tersimpan · salah satu tulis-balik gagal → seluruh transaksi
batal, `503`.

**Uji:** BARU — **ketujuh baris** fungsi akibat, ketiga nilai tampil diperiksa **dari model
baca**, bukan dari kode · BARU untuk keputusan bersyarat mendarat pada jenjang yang memutus ·
BARU — uji DDL: nol kolom menyimpan "sedang disirkulasikan".

**Menggantikan:** `KomitePostAdjustment` — akibat jalur A, tersebar ·
`KomitePostAdjustmentCWP` — akibat jalur B · `KomitePostAdjustment`·14.12 — penanda usul tutup
mendarat di induk pada persetujuan terakhir non-bersyarat (`F-5`) ·
`KomitePostAdjustment`·20 — penolakan satu usulan menolak seluruh klaim; **dipertahankan**
sebagai baris ketiga tabel, ditandai `D-1` paritas-dengan-pertanyaan-bisnis, bukan cacat ·
`InsertChronology_DT` — kronologi, **tanpa** pengecualian per peran · `KOMITECOUNT` dan
`FLAGONGOINGCOMMITTE` sebagai kolom, **dibuang**.

**Blocked by:**

- `04` — Keputusan tercatat pada jenjang yang benar, dan orang lain ditolak


**Dasar:** DECIDED(`D-3`, `J-3`, `D-1`, `D-2`, `K-07`, keputusan beku no. 9, ADR-0016,
ADR-0032). EVIDENCED: `F-5`; GRILL-06 P6-4 — jalur bersyarat menulis ke halaman yang
berbeda, bukan hanya ke baris yang berbeda.

- [ ] Fungsi akibat ditulis **sekali**, di satu tempat; ketujuh barisnya diuji.
- [ ] Hasil tidak pernah berarti nasib klaim; keadaan sirkulasi dan akibat pada klaim
      menempati **kolom yang berbeda**.
- [ ] Model baca menyajikan ketiga nilai tampil dari fungsi yang sama; **nol** pembaca
      menafsirkan hasil sendiri.
- [ ] Seluruh tulis-balik terjadi di dalam transaksi keputusan — nol pesan antar-konteks, nol
      tautan basis data.
- [ ] Keputusan bersyarat mendarat pada **jenjang yang memutus**, bukan pada baris terakhir
      daftar dan bukan pada halaman yang berbeda. Invarian `I-3`.
- [ ] Penanda "sedang disirkulasikan" **tidak disimpan**; ia dihitung dari jenjang saat
      dibaca.
- [ ] Pengaju dan maksud tersaji **terpisah** dari akibat.

**Ketidakpastian:** Tidak ada.
