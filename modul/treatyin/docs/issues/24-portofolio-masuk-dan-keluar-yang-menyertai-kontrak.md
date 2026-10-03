---
status: aktif
golongan: pelestarian
---

# 24: Portofolio masuk dan keluar yang menyertai kontrak dicatat

*Asal: `DAFTAR-PEKERJAAN.md` `P-09` · `SPEC-MODEL-DATA.md` §10.13 · `SPEC-INVARIAN.md` `INV-66`.*

**What to build:** **PK** mencatat **portofolio masuk** dan **portofolio keluar** yang menyertai
sebuah versi kontrak — premi maupun klaim — dengan arah dan jenisnya terbedakan.

Artefak: entitas `PORTOFOLIO`, kunci asingnya ke `VERSI_KONTRAK`, dan `INV-66`.

**PEMBUAT PERTAMA** untuk `PORTOFOLIO`.

**Kenapa begini:** **Arah adalah pembedanya, bukan nama kolom.** Sistem lama menyimpan portofolio
sebagai daftar `Portfolio` dengan arah sebagai nilai; membawanya sebagai dua kolom bernama — "masuk"
dan "keluar" — akan menuntut kolom ketiga begitu ada arah lain, dan itu bentuk yang sama dengan
delapan kolom bahaya yang `ADR-0038` bongkar. Kunci alaminya **arah + jenis** di dalam satu versi,
dan itu baru bernomor pada langkah 3 to-spec (`INV-66`) — sebelumnya entitas ini berdiri **tanpa kunci
alami sama sekali**.

**Persyaratan:** `INV-66` (arah + jenis portofolio unik di dalam satu versi) · `INV-36` · `INV-01`

**Tidak termasuk:** **Perhitungan nilai portofolio** — ia turunan dari premi dan klaim periode
berjalan, dan `ADR-0037` melarang turunan disimpan. Yang dicatat **kesepakatannya**, bukan hasilnya.
**Portofolio pada cabang retro keluar** — `RETRO_KELUAR` adalah **GEL-2**, tegas di luar penyerahan
pertama.

**Jalur gagal:** Dua baris berarah **dan** berjenis sama pada satu versi -> **ditolak** `INV-66` ·
Nilai portofolio terisi tanpa mata uang -> ditolak `INV-36`.

**Uji:** **Negatif:** baris kembar pada sumbu arah + jenis.
**Positif — dan ia yang menangkap kunci yang terlalu sempit:** satu versi dengan portofolio **masuk**
dan portofolio **keluar** berjenis sama -> **diterima**. Sebuah `UNIQUE` yang hanya menyebut jenis
lulus uji negatif di atas dan menolak pasangan masuk-keluar yang justru **bentuk paling lazim**.

**Menggantikan:** `P-09` melestarikan daftar `Portfolio`. Yang bergeser: ia memperoleh **kunci alami
bernomor**, yang sebelumnya tidak ada — `SPEC-MODEL-DATA.md` §5.0a menandainya *"kunci alaminya belum
bernomor"* sampai langkah 3 to-spec menutupnya.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(TreatyIn.Portfolio@ekspor-2026-09 - daftar berarah sebagai nilai)
        DECIDED(INV-66, ADR-0037, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `PORTOFOLIO` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-66` terpasang atas **arah + jenis** di dalam satu versi
- [ ] uji positif lulus: masuk dan keluar berjenis sama **diterima**
- [ ] arah tersimpan sebagai **nilai**, bukan sebagai dua kolom bernama — diperiksa terhadap `KAMUS-KOLOM.md`
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`
