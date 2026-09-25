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

**Status:** blocked

- [ ] Keenam lini mengisi properti nilai lama **masing-masing sesuai petanya**, dan hanya itu
- [ ] Tingkat cedant identik di ketujuh lini
- [ ] **Seluruh** penugasan nilai lama memakai pola kosong→nol; **nol pengecualian**
- [ ] Premi Nusantara Re lama muncul **hanya** di dua lini, bukan tujuh
- [ ] Nilai uang bertipe `Money` dengan mata uangnya; rate lama bertipe `Ratio` dengan **skala per lini** (K-018)
- [ ] Lapis B diisi ulang **setiap kali layar endorsement dibuka** — bukan sekali saat kasus lahir
- [ ] ⚠️ Lapis B menangani **travel tetapi bukan jiwa**; himpunannya berbeda dari lapis C
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
