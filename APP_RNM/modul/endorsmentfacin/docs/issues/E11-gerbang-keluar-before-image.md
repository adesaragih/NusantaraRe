# E11: Tiga gerbang keluar lapis B

**What to build:** Pengisian nilai lama per baris **berhenti sebelum berjalan** pada tiga keadaan.
Ketiganya menghasilkan lapis B kosong — dan itu **perilaku yang benar**, bukan galat.

`[terverifikasi]` Gerbang pertama mengeluarkan kasus **jiwa**; gerbang kedua mengeluarkan kasus yang
**bukan endorsement**; gerbang ketiga mengeluarkan polis dengan **lokasi melebihi ambang**.

⛔ **Polis besar tidak mendapat nilai lama per baris sama sekali** — layar dan perhitungan per baris
menampilkan nol. Diport apa adanya; alasan ambangnya tidak diketahui dari korpus.

📌 Gerbang kedua memakai predikat yang membaca **agregat tersimpan**. Predikat kebalikannya membaca
**halaman aktif** — lihat E02; jangan diimplementasikan sebagai negasi.

**Asal (Pega).** `Activity\SetOldData` langkah 1–3

**Keputusan.** K-044 · K-047 · K-046

**Blocked by:** E08

**Status:** selesai — 01-10-2026, `backend/services/nilailama.go` `IsiNilaiLama` langkah 1–3; test
`TestGerbangKeluarLapisB` (termasuk batas tepat 100 lokasi = tidak keluar). `IsLife`/`IsEDM` diterima
sebagai masukan sampai registry E01 ada.

- [x] Kasus **jiwa** keluar sebelum blok pengisian — jiwa memakai jalurnya sendiri (E18)
- [x] Kasus **bukan endorsement** keluar
- [x] Polis dengan **lokasi melebihi ambang** keluar — `> 100`, daftar **kerja** (`.OfferFacIn.LocationList`)
- [x] Ketiganya menghasilkan lapis B **kosong tanpa galat** — bukan kegagalan
- [x] Ambang lokasi diambil dari korpus apa adanya; alasannya tetap **`[pertanyaan terbuka]`**
- [x] Menyebut rule Pega asalnya dalam komentar (§4.6)
