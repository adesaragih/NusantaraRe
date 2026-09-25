# E07: Lapis A — sembilan penyalinan yang menyimpang dari pola

**What to build:** Penyalinan dokumen polis lama tidak seluruhnya berbentuk "salin field bernama
sama". Sembilan di antaranya berbeda, dan **empat halaman akan kosong** bila diabaikan.

`[terverifikasi]` Dari **54 penugasan**: **51** bersumber dokumen polis lama, **3** tidak. Dari yang
51, **45** berpola seragam dan **6 menyimpang** — menulis ke halaman **selain** agregat penawaran.

> 🔁 **Arsip lintas-siklus keliru di titik ini.** Ia menyatakan seluruh 54 berbentuk seragam. Itu
> **tidak berlaku untuk 9 dari 54**. Implementasi yang mengikuti arsip meninggalkan empat halaman
> kosong.

**Asal (Pega).** `Activity\SetValueToEDMWork` langkah 14.3

**Keputusan.** K-047 · `CLAUDE.md` §4.6

**Blocked by:** E06

**Status:** blocked

- [ ] **54 penugasan** terimplementasi, tidak diringkas
- [ ] **Enam cermin lintas-halaman** ditulis eksplisit: jenis bisnis ke halaman kutipan · tiga field pembayaran ke halaman polis · penanda B2B ke akar objek kerja · nama marketing ke halaman parameter kredit
- [ ] **Kriteria kunci:** keempat halaman itu **tidak boleh kosong** setelah lapis A berjalan
- [ ] Tiga field pembayaran ditulis **dua kali** — ke agregat penawaran **dan** ke halaman polis
- [ ] Tiga penugasan yang **bukan** dari dokumen lama tetap dijalankan
- [ ] Empat kelompok field terisi lengkap: identitas bisnis · periode · struktur share dan kapasitas · pembayaran
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
