# E08: Lapis B — nilai lama per baris untuk enam lini sisanya

**What to build:** Nilai lama per baris untuk golf, aneka, personal accident, marine cargo,
kendaraan bermotor, dan travel — plus tingkat cedant yang sama di ketujuh lini.

`[terverifikasi]` Pola nilainya **tunggal dan tanpa pengecualian**: bila properti sumber tidak
kosong pakai nilainya, bila kosong pakai **nol**.

⛔ **Kosong dipetakan ke nol, bukan ke nilai tak-ada.** Aritmetika selisih bergantung pada ini.

📌 Premi Nusantara Re lama **hanya ada di dua lini** — marine cargo dan kendaraan bermotor.

**Asal (Pega).** `Activity\SetOldData` blok 4.3–4.14

**Keputusan.** K-047 · K-010/K-012 · K-018 · `CLAUDE.md` §4.6

**Blocked by:** E06

**Status:** sebagian — 01-10-2026, `backend/services/nilailama.go`; test `nilailama_test.go`. 52
penugasan diurai ulang dari `SetOldData.xml` (dua cara sepakat). Sumber nilainya baris **OldData**,
dipasangkan menurut nomor urut — langkah 4 berjalan di halaman OldData. Sisa: satuan ‰/% rate (K-018,
resolver NB-03) dan pemanggil layar.

- [x] Keenam lini mengisi properti nilai lama **masing-masing sesuai petanya**, dan hanya itu
- [x] Tingkat cedant identik di ketujuh lini
- [x] **Seluruh** penugasan nilai lama memakai pola kosong→nol; **nol pengecualian**
- [x] Premi Nusantara Re lama muncul **hanya** di dua lini, bukan tujuh
- [ ] Nilai uang bertipe `Money` dengan mata uangnya ✅; rate lama bertipe `Ratio` ✅ dengan **skala per lini** (K-018) ⛔ — `RateOld` membawa skala `Rate` sumbernya; satuan ‰/% menunggu resolver NB-03. ⚠️ Nol cadangan bermata uang = mata uang nilai sumbernya (bisa kosong) — `[pertanyaan terbuka]`
- [x] Lapis B diisi ulang **setiap kali layar endorsement dibuka** — `services.IsiNilaiLama` terpisah dan idempoten; pemanggil layarnya belum ada
- [x] ⚠️ Lapis B menangani **travel tetapi bukan jiwa**; himpunannya berbeda dari lapis C. `[dugaan]` Travel: baris kerja yang belum ada dibuat di ujung daftar (lihat laporan §2.4)
- [x] Menyebut rule Pega asalnya dalam komentar (§4.6)
