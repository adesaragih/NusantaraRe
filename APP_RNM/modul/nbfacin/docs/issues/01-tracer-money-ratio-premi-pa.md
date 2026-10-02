# 01: Tracer — uang, rasio, dan satu jalur premi PA yang hidup ujung ke ujung

**What to build:** Satu kasus PA sederhana masuk, satu angka premi keluar — lewat tipe uang dan rasio
yang sesungguhnya, bukan angka telanjang. Ini irisan **tertipis yang lengkap**: begitu ia hijau,
seluruh rantai (parsing masukan → uang bermata-uang → rasio berskala → rumus → pembulatan → hasil)
terbukti hidup, dan tiket berikutnya tinggal menambah bentuk rumus.

Uang tidak pernah `float`. Rasio membawa satuannya sendiri. Keduanya **tipe berbeda yang tidak dapat
dijumlahkan** — satu-satunya jembatan adalah perkalian eksplisit `uang × rasio → uang`. Nilai masuk
sebagai teks berkoma desimal dan dikonversi **di batas input**, bukan tersebar di dalam perhitungan.

Bentuk rumus yang dipakai tracer ini adalah yang paling sederhana dari PA — pembagi 1.000 langsung,
tanpa pro-rata: `CalculatePremiPA_FacIn` **L1003**. Bentuk PA lainnya menyusul di tiket 04.

**Blocked by:** None (can start immediately)

**Status:** ready-for-human — diport 01-10-2026, menunggu tinjauan work owner

- [ ] Nilai uang memakai tipe desimal dengan **mata uang wajib menyertai**; tidak ada `float` di jalur uang mana pun
- [ ] Rasio memakai tipe terpisah yang **membawa skalanya sendiri** (per mille / persen)
- [ ] Menjumlahkan uang dengan rasio **gagal saat kompilasi** — dibuktikan berkas uji yang wajib tidak terkompilasi, bukan unit test runtime
- [ ] Satu-satunya jembatan uang↔rasio adalah operasi perkalian eksplisit
- [ ] Parser masukan menerima **koma sebagai pemisah desimal** dan memangkas spasi di ujung; konvensi ditetapkan eksplisit, **tidak diserahkan ke locale** (K-027)
- [ ] Resolver lini bisnis → skala berisi **PA saja** pada tiket ini
- [ ] Pintu masuk perhitungan premi menerima satu kasus PA dan mengembalikan uang
- [ ] **Rekonsiliasi eksak**: fixture PA bentuk `L1003` cocok dengan sistem lama **sampai digit terakhir**, tanpa toleransi (ADR-0001)
- [ ] Pembulatan menuliskan presisinya di tempatnya, disertai komentar yang menyebut rule Pega asal dan nomor langkahnya

## Comments

### 2026-10-01 — implementasi tracer (agent, `/implement`)

Kode: `APP_RNM/modul/nbfacin/backend/services/premium/` (`premium.go`, `premium_test.go`,
`kompilasi_test.go`, `testdata/`). Belum di-commit (permintaan work owner: git diurus sendiri).

⚠️ **Tracer ini TIDAK memakai L1003 sebagai premi akhir.** `[terverifikasi]` Diurai dari
`NB FacIn\Activity\CalculatePremiPA_FacIn.xml`: langkah 6 (L1003, `CalculateMethod_FacIn=='3'`)
ditimpa langkah 7 (L1146) yang ber-`pyStepsPreCondition=false` — menurut P-11 tetap jalan tanpa
syarat. Keputusan work owner 01-10-2026: **port apa adanya**, premi akhir PA = L1146. `[dugaan]` P-11
belum dilihat langsung pada langkah 7 ini — `../PERTANYAAN-NB01.md` butir 3.

| Kriteria | Keadaan |
| --- | --- |
| Uang desimal, mata uang menyertai, nol `float` | ✅ `inti/backend/uang.Money` (keputusan: ikuti `inti`) |
| Rasio tipe terpisah membawa skala ‰/% | ✅ tipe lokal `rasio` = `uang.Ratio` + satuan; `uang.Ratio.Scale` tidak dipakai (di `claimlife` ia jumlah desimal) |
| Uang + rasio gagal kompilasi | ✅ `TestUangTambahRasioGagalKompilasi` + berkas kontrol yang wajib terkompilasi |
| Jembatan perkalian eksplisit | ✅ satu-satunya jalan: `kaliEksak` di dalam rumus |
| Parser koma K-027, spasi dipangkas, tanpa locale | ✅ titik **ditolak** (ambigu); `70,000 ` = 70 teruji |
| Resolver lini bisnis PA saja | ✅ lini lain → `panic` |
| Pintu masuk menerima satu kasus PA, mengembalikan uang | ✅ `premium.Calculate` |
| Rekonsiliasi eksak dengan sistem lama | ⛔ **terbuka** — tidak ada kasus PA di lima berkas kasus; tes memakai contoh hitung manual (keputusan 30-09-2026). `../PERTANYAAN-NB01.md` butir 1 |
| Presisi pembulatan literal + komentar asal rule | ✅ presisi 20, L1146. ⚠️ Mode pembulatan belum terverifikasi → hasil tak eksak **ditolak** (`ErrModePembulatanBelumTerverifikasi`); butir 2 |

Tes diuji mutasi: enam mutasi (cek tak eksak, koma sebagai ribuan, pembagi 1.000, lini bawaan,
pembulatan dibuang, titik diterima) masing-masing menggagalkan tes.

### 2026-10-01 — tindak lanjut `/code-review` (agent)

Review dua sumbu (Standards dan Spec) menghasilkan 4 pelanggaran + 6 smell dan 7 temuan spec. Keputusan
work owner atas temuan yang butuh keputusan: `../KEPUTUSAN-30-09-2026.md` butir 9–12.
**Entri sebelumnya di atas basi pada tiga baris tabelnya** — yang berlaku:

| Kriteria | Keadaan sekarang |
| --- | --- |
| Rasio tipe terpisah membawa skala ‰/% | ✅ `rasio` = `uang.Ratio` + satuan; gagal-kompilasi kini diuji untuk `rasio` lokal juga (`kompilasi_gagal_rasio.go`, tag `gagalkompilasi`) |
| Jembatan perkalian eksplisit | ✅ satu metode bertipe `rasio.kali(uang.Money) (uang.Money, error)` |
| Parser koma K-027 | ✅ **hanya `.Rate`**; TSI, ProRatePercent, Discount lewat `utils.ParseDecimal` berlabel `[dugaan]` |

Perubahan lain: medan kosong = kosong, bukan nol (`uang.ErrUangKosong` / `ErrRasioKosong`, ADR-U-0022);
konversi di awal `Calculate` dicatat sebagai batas masukan sementara; komentar langkah 2–3 diperbaiki
(langkah 3 menulis `.DiscountPercentage`); syarat tersimpan `.CalculateMethod_FacIn==1` di langkah 7
(L1226) masuk komentar dan `../PERTANYAAN-NB01.md` butir 3; butir 4 baru (`.Discount` kosong); uji
mata uang kosong; `t.Skip` bila perintah `go` tidak ada. Usulan ke tim inti:
`../USULAN-PR-TIM-INTI-NB01.md`.

Uji mutasi diulang: sembilan mutasi, masing-masing menggagalkan tes yang tepat; setiap mutasi
diperiksa tetap terkompilasi (satu tangkapan palsu dari putaran pertama sudah dibuang).

### 2026-10-01 — rekonsiliasi kasus nyata (agent)

⚠️ **Dua entri di atas basi pada rumus premi akhir.** Yang berlaku (`../KEPUTUSAN-30-09-2026.md`
butir 13–17):

- **Langkah 7 (L1146) tidak dijalankan.** `CalculateMethod_FacIn` `'1'` → langkah 4 L713, `'3'` →
  langkah 6 L1003 (bentuk awal tiket ini), `'2'` → L858 belum diport (`ErrBentukBelumDiport`, tiket 04).
- Semua medan bertitik desimal; koma K-027 pindah ke pemuat `repository` kelak.
- Diskon kosong dikurangi nol; pembulatan seri ditolak (`ErrPembulatanSeri`). ⚠️ *Diganti A37 (01-10-2026, tiket 18): setengah-ke-atas, `ErrPembulatanSeri` dihapus.*

| Kriteria | Keadaan |
| --- | --- |
| **Rekonsiliasi eksak** | ✅ **4/4 kasus PA nyata nol selisih** (1 metode 1, 3 metode 3) — `rekonsiliasi_test.go`, daftar-izin lima medan, tanpa nomor kasus. Diperiksa dua cara: Python `decimal` dan Go `apd` |
| Parser koma K-027 | ⚠️ **Dipindah ke luar tiket ini**: format kolom Oracle, bukan data kerja Pega. Tiket ini kini membaca titik; parser K-027 menunggu pemuat `repository` (usulan `../USULAN-PR-TIM-INTI-NB01.md`) |
| Kriteria lain | ✅ seperti entri sebelumnya |

Uji mutasi putaran ketiga: sepuluh mutasi, semuanya tertangkap. Sisa terbuka: `../PERTANYAAN-NB01.md`
butir 5 (arti `pyStepsPreCondition=false`, bertentangan dengan P-11).

### 2026-10-01 — butir 5 dijawab (agent)

Pengecualian P-11 berlaku untuk PA saja (`../KEPUTUSAN-30-09-2026.md` butir 18). Tidak mengubah kode
tiket ini; komentar `premium.go` diperbarui. Kelima pertanyaan `../PERTANYAAN-NB01.md` kini terjawab.
