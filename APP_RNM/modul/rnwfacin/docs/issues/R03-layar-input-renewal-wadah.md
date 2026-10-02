# R03: Layar input renewal — wadah, bukan duplikat

**What to build:** Layar input renewal menampilkan penawaran **dengan menyisipkan komponen yang sudah
ada dari New Business** — bukan dengan menyalinnya. Underwriter melihat layar yang sama seperti yang
ia kenal, dan tim tidak memelihara dua salinan yang bisa menyimpang.

Ini kriteria yang paling mudah dilanggar tanpa sengaja, dan paling mahal bila dilanggar: menyalin
komponen NB terasa lebih cepat, lalu perbaikan pada satu sisi diam-diam tidak sampai ke sisi lain.

`[terverifikasi]` Bukti bahwa layar renewal memang wadah: `Section\InputRenewal` hanya punya
**2 properti sendiri** (`.IsShowDetail` + template), dan menyisipkan **7 section**:
`AllSummarySection` · `PeriodeRenewal` · `InputRenewalDtl` · `InputDtlObject_FacIn` ·
`InputInwardFacultativeSuggest` · `EmailSection` · `EmailSectionCeding`.
Varian `_IsUW` menyisipkan **6 section**, juga dengan 2 properti sendiri.

📌 Empat section yang disisipkan **tidak ada di folder renewal** — `InputInwardFacultativeSuggest` ·
`InwardFacIn` · `OfferFacIn_NusaRe` · `OfferFacIn_NusaRe_IsUW`. Itu **bukan kekurangan**: keempatnya
ada di New Business dan **diwarisi apa adanya** (K-040), persis seperti 1.907 berkas identik lainnya.
Bangun sekali, pakai di kedua siklus.

⚠️ **Dua komponen terpisah** — `InputRenewal` dan `InputRenewal_IsUW` masing-masing komponen sendiri.
Keputusan sadar work owner, sengaja berbeda dari rekomendasi spec model-data NB §6; alasannya
reproduksi perilaku terekam (`CLAUDE.md` §1).

**Blocked by:** R02

**Status:** needs-info — layar; lingkup RNW sebatas layanan (butir 55), menunggu keputusan membangun layar · ⏸ pengerjaan rnwfacin dihentikan work owner 01-10-2026 (`PROMPT-LANJUT-IMPLEMENT-NB-SAJA.md`)

- [ ] `InputRenewal` dan `InputRenewal_IsUW` diport sebagai **dua komponen terpisah**
- [ ] Keduanya **menyisipkan** komponen NB — **tidak ada komponen NB yang disalin** ke dalam kode renewal
- [ ] **Kriteria pembukti:** perubahan pada sebuah komponen NB **terlihat di layar renewal** tanpa perubahan kode renewal. Bila tidak terlihat, komponennya tersalin — dan tiket ini gagal
- [ ] Keempat section yang diwarisi dari NB dipakai dari sumber NB, bukan dibuat ulang
- [ ] Layar periode (R02) tersisip, bukan diduplikasi
- [ ] Komentar menyebut section Pega asal
