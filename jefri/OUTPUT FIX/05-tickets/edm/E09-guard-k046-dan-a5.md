# E09: Dua guard yang dipertahankan dan satu perbaikan sadar

**What to build:** Dua guard di sistem lama menguji **properti tujuan** alih-alih properti sumber.
Keduanya **diport apa adanya**. Satu hal diperbaiki — dan hanya satu.

`[terverifikasi]` Guard rate lama berperilaku demikian pada **seluruh** kemunculannya; guard premi
lama hanya pada **satu** dari delapan belas kemunculan.

⛔ **Keduanya perilaku yang BENAR menurut work owner (K-046).** Meluruskannya mengubah angka selisih
dan memutus rekonsiliasi paralel run.

⚠️ **A.5 — satu-satunya perbaikan sadar di seluruh modul before-image.** Nilai cadangan premi
Nusantara Re lama memakai **string** di sistem lama; di sistem baru ia **angka nol**, karena nilai
uang bertipe `Money` dan `Money` tidak boleh menampung string.

**Asal (Pega).** `Activity\SetOldData` — penugasan premi Nusantara Re lama (dua tempat) dan guard
rate/premi lama

**Keputusan.** **K-046** · **A.5** · K-010/K-012 · ADR-0001

**Blocked by:** E08

**Status:** blocked

- [ ] **K-046** `K046_RateOld_GuardSelfReferential_6dari6` — guard menguji properti tujuan; **perilaku benar**, bukan cacat
- [ ] **K-046** `K046_PremiumOld_SelfReferential_HanyaFire` — hanya satu dari delapan belas kemunculan
- [ ] Komentar kode menyebut rule Pega asal **dan** K-046, agar tidak "diluruskan" pembaca berikutnya
- [ ] **A.5:** nilai cadangan premi Nusantara Re lama adalah **angka nol**, bukan string — ditandai **perbaikan, bukan port**
- [ ] Pembanding rekonsiliasi **menormalkan** string dan angka sebelum membandingkan
- [ ] Selisih tipe yang muncul **dijelaskan oleh A.5** dalam prosedur rekonsiliasi — bukan ditandai cacat
- [ ] ⚠️ Cache ekspresi di korpus memuat bentuk berbeda dari yang dieksekusi; **yang mengikat ekspresi tersimpan**
