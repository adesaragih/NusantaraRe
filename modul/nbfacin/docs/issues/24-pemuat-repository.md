# 24: Pemuat data lama — seam `repository` (himpunan baris → Oracle) + perkakas pemuatan

> ⚠️ **Disusun agent dari spec atas perintah work owner — bukan hasil `/to-tickets`.** Perintah: arahan work owner
> 02-10-2026 langkah 4 (pemuat repository + tiket 17).

**What to build:** seam **`repository`** spec 11 — himpunan baris `loader.Flatten` (tiket 22) ditulis ke tabel flat
dalam **satu transaksi** per penawaran, ditambah **pemilihan versi terakhir** fase 1 dan perkakas pemuatan di atas
repository (bukan dipanggil dari `handlers`).

**Asal** `[terverifikasi]`:
- Spec 11 *Keputusan yang mengikat* 1–7, 15; *Yang diuji di seam `repository`* 1–5; BAHAN §5a, §8.
- `D:\migrasi\RNM\DDL\JSON_POLIS.txt`: `IDPEGA VARCHAR2(50)` (PK), `DATA_JSON CLOB`, `NOPOLIS VARCHAR2(100) NOT NULL`,
  `NOENDORS`, **`PRODKE VARCHAR2(5)`**, `TGL_INPUT DATE`, `STS_KONVERSI`, `TGL_KONVERSI`, `USERNAME`; unik
  `(NOPOLIS, IDPEGA)`.
- Asal urutan: `GetProdKeOldData_SQL` (`ORDER BY TGL_INPUT DESC`), dirujuk `NB FacIn\Activity\GetData_ACT.xml` L472 dan
  `GetPolicyData_ACT.xml` L764.

**Blocked by:** 22, 23 (tabel belum ada; migrasi ditahan butir 66)

**Status:** needs-info — menunggu 23; seam dapat mulai diuji dengan tiruan basis data setelah bentuk tabel tetap

## Yang dibangun

- [ ] Pembaca `JSON_POLIS` (baca saja): `IDPEGA`, `DATA_JSON`, `NOPOLIS`, `PRODKE`, `TGL_INPUT`, `STS_KONVERSI`,
      `TGL_KONVERSI`, `USERNAME`; `T_WORK_POLIS.POSISI/NOURUT/PUTARAN` dari baris terakhir
      `HISTORYAKSEPTASIPRODUCTION` (V-48) — ⛔ `USERNAME` login orang: disalin ke kolomnya, tidak pernah dicatat di log
- [ ] **Versi terakhir fase 1:** kelompok `NOPOLIS`, urut `TGL_INPUT DESC`; periksa silang `PRODKE` tertinggi
      (dibandingkan **sebagai angka**) — menunjuk baris berbeda → **`panic`**
- [ ] **`PROD_KE`:** teks → angka; bukan angka → **`panic`** (BAHAN §3), bukan kosong, bukan nol
- [ ] `OLD_POLIS_ID` fase 1 `NULL` seluruhnya (K-071 J-5)
- [ ] Satu penawaran = **satu transaksi**; ID dari sekuens, `PARENT_ID` = ID baris `Baris.Induk`; urutan sisip =
      `Baris.Kunci` (induk selalu lebih dulu)
- [ ] **Tiket 17:** laporan per penawaran dari `Diagnostik.MataUangTakDiketahui` (sudah dihitung Flatten); baris lama
      masuk `UNKNOWN` (K-069) — ADR-0006 `panic` berlaku bagi **tulisan baru aplikasi**, bukan baris migrasi
- [ ] Laporan per jalan: baris per tabel, `TakTerpetakan`, `PenunjukBelumDikonversi`, `Dibuang` — **hitungan saja**
- [ ] Dapat diulang: dua jalan → isi tabel sama (`ROW_UID` deterministik, A49)
- [ ] **Fase 2 mengosongkan 78 tabel** lalu memuat penuh — ⛔ operasi destruktif: hanya atas perintah eksplisit work
      owner, tercatat, tidak pernah efek samping

## Keputusan work owner 02-10-2026 (butir 68.8, diteruskan sesi `nusantarare-0f`, "setuju")

- [ ] Penawaran yang **gagal dimuat** (galat `Flatten` atau galat tulis) → **dilewati**, galatnya dicatat (tanpa nilai
      data), pemuatan lanjut ke penawaran berikutnya; **ringkasan di akhir muat** (berhasil / dilewati per jenis galat).
      Berhenti keras spec 11 tetap berlaku **per penawaran** — transaksinya batal, tidak ada baris setengah jadi.
- [ ] Medan tak terpetakan masuk **penampung** (ADR-0023, `Hasil.Penampung` tiket 22). Ringkasan akhir muat menyebut
      jumlah entri penampung per jalur :: medan (tanpa nilai).
- [ ] ⚠️ **Ralat butir 68.8** (keputusan work owner 02-10-2026 butir 69, "mau"): rekomendasi lama "medan tak terpetakan
      masuk penampung, **bukan alasan menolak muat**" **kurang tepat** — ADR-0023 berbunyi "Penampung itu **wajib kosong**
      sebelum pekerjaan dinyatakan selesai" dan akibat 2 "Penampung berisi = pekerjaan belum selesai". Yang berlaku:
      **muat untuk UJI boleh** dengan penampung berisi; **muat produksi / fase 1 baru boleh dinyatakan selesai bila
      penampung kosong** — setiap medannya sudah mendapat kolom atau keputusan "dibuang" eksplisit.

## Masih terbuka

1. ~~Di mana penampung disimpan di Oracle~~ → ✅ **butir 69: tidak disimpan di Oracle** (ADR-0023: penampung bukan tempat
   simpan akhir). Usulan kolom untuk isinya: `docs/USULAN-KOLOM-PENAMPUNG.md`, menunggu work owner.
2. **J-16** penyaring "endorsement sudah jadi" — pertanyaan terbuka untuk work owner (`STS_KONVERSI` bukan jawabannya).

## Comments
