---
status: tertahan
---

# 02: Sirkulasi usulan pembayaran lahir lengkap dengan jenjangnya

*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Analis mengirim satu usulan pembayaran ke komite; sistem memilih kelas
kewenangan dari tabel seleksi, menyusun jenjang dari roster urut derajat menaik, dan
sirkulasi itu terbaca kembali lengkap dengan jenjangnya. Kombinasi masukan yang tidak
tercakup aturan **ditolak sebagai galat konfigurasi**, dan sirkulasi tanpa jenjang **ditolak
saat dibuat** — bukan lahir lalu diam.

`SIRKULASI` dan `JENJANG_SIRKULASI` beserta sequence dan constraint-nya; operasi
`BentukSirkulasi` dan `BacaSirkulasi`; jalur pembentukan dan jalur baca satu sirkulasi serta
seluruh sirkulasi sebuah klaim.

Ini peluru penjejak inti: ia memotong skema, perintah, kueri, dan uji sekaligus. Ia tidak
dipecah karena tiap belahannya menjadi irisan mendatar yang tidak dapat didemokan sendiri.

**Persyaratan:** `S-001`, `S-002`, `S-003`, `S-004`, `S-005`, `S-006`, `S-007`, `S-013`,
`S-014`, `S-016`, `S-044`, `S-047`, `S-064`, `S-065`, `S-067`, `S-070`.

**Tidak termasuk:** jalur menutup dan menolak klaim — tiket `07`. Jalur gagal pembentukan dan
maksud pengajuan — tiket `03`. Memutus — tiket `04`.

**Jalur gagal:** jenis sirkulasi kosong atau di luar tiga nilai → `422` menyebut medan jenis ·
kombinasi masukan tak tercakup aturan seleksi → `422` menyebut ketiga nilainya, nol sirkulasi
dan nol roster tersusun · roster terpilih kosong → `422` menyebut kelas kewenangan yang tidak
punya pemegang aktif · pemanggil tidak berhak mengajukan atas klaim itu → `403`, dapat
dibedakan dari `503` · kunci idempotensi sama dengan muatan berbeda → `409` · sirkulasi atas
usulan itu masih berjalan → `409`.

**Uji:** P5-01, P5-03 (derajat dan giliran pertama) · P5-09, P5-11 (kombinasi tak tercakup) ·
P5-10 (kelas dari nilai dan bagian) · P5-13 (nilai usulan yang diajukan, **wajib memakai klaim
dengan dua usulan beda kelas** — pada klaim berusulan tunggal deviasi ini lolos tanpa
terlihat) · P5-06 (roster idempoten) · P5-16 (hak) · BARU untuk idempotensi dan nomor urut.

**Menggantikan:** `KomiteTreaty_Flow` — mesin keadaan sirkulasi ·
`CreateChildKomiteCNP_Act`·9, 26.13–26.15 — jumlah jenjang dari cacah baris laporan, menjadi
akibat isi roster (`F-13`) · `CreateChildKomiteCNP_Act`·10 — bagian treaty ke pembanding ·
`CreateChildKomiteCNP_Act`·11 — nilai baris terakhir sebagai penentu (`F-17`) ·
`CreateChildKomiteCNP_Act`·14 — cabang yang syaratnya tidak pernah dapat benar, **dibuang**
(`F-16`) · `CreateChildKomiteCNP_Act`·28 — pembatalan hanya pada alokasi layer kosong
(`F-22`) · `KomitePostAdjustment`·16 — nomor urut sirkulasi atas satu usulan (`F-11`) ·
`KomiteCount` dan `KomiteLoop` sebagai kolom, **dibuang** — keduanya turunan (`D-2`, `K5-1`).

**Blocked by:**

- `01` — Roster dan tabel seleksi berdiri sebagai data yang dikelola


**Dasar:** DECIDED(`K5-1`, `K5-4`, `K5-5`, `K6-4`, `H-6`, `D-2`, `E-2`, keputusan beku no. 8,
ADR-0030, ADR-0031, ADR-0032). EVIDENCED: `F-11`, `F-13`, `F-14`, `F-16`, `F-17`, `F-22`,
`F-24`.

- [ ] Jenis sirkulasi datang dari **pemanggil**; nol pembacaan kolom analisis.
- [ ] Kelas kewenangan ditentukan satu pencarian ke tabel seleksi dengan tiga masukan; nol
      cabang kelas di dalam kode.
- [ ] Masukan nilai adalah nilai **usulan yang sedang diajukan**, bukan baris terakhir daftar
      dan bukan jumlah seluruh usulan.
- [ ] Derajat disalin dari roster saat pembentukan, urut menaik, dan **tidak dihitung ulang**
      sesudahnya.
- [ ] Dua jenjang berderajat sama dalam satu sirkulasi **tidak dapat ditulis** — ditolak
      constraint, bukan diperiksa kode. Ini yang mewujudkan invarian `I-2`.
- [ ] Membentuk ulang roster sebuah sirkulasi menghasilkan roster **sama panjang**, bukan
      bertambah.
- [ ] Sirkulasi tanpa jenjang **ditolak saat dibuat**; nol berkas lahir.
- [ ] Nol kolom menyimpan jumlah jenjang — dibuktikan pencarian atas katalog skema yang
      pulang kosong. Ini yang mewujudkan invarian `I-8`.
- [ ] Dua sirkulasi bernomor urut sama atas satu usulan ditolak constraint.
- [ ] Permintaan berulang berkunci sama mengembalikan sirkulasi yang sama, bukan yang kedua.
- [ ] Tali klaim–sirkulasi dibaca **dari sisi sirkulasi**; nol pembacaan yang bersandar pada
      pengenal tersimpan di sisi klaim.
- [ ] Pembaca yang berhak membuka klaim dapat membaca sirkulasinya, tanpa peran tambahan.

**Ketidakpastian:** Tidak ada.
