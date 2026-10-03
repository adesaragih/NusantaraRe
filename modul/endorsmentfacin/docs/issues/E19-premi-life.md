# E19: Perhitungan premi jiwa

**What to build:** Premi endorsement jiwa dihitung lewat **Seam 3 bentuk 1 yang terparameterisasi** —
bukan bentuk kelima.

`[terverifikasi]` Bentuknya tetap keluarga rumus dasar, tetapi dengan **tiga perbedaan
parameter**: basisnya nilai pertanggungan **setelah dikurangi bagian ceding**, rate diambil dari
**properti rate rata-rata jiwa**, dan ada **satu faktor pembebanan tambahan**.

⛔ **Properti bernama "rate" di rumus ini berperan sebagai faktor pembebanan**, bukan rate premi.
Membacanya sebagai rate premi akan salah total.

✅ Pembagi yang lebih kecil daripada lini lain adalah **konsekuensi K-018**, bukan rumus lain —
hanya satu faktor yang berkontribusi ke satuan.

**Asal (Pega).** `Activity\ReCountPremiLifeEDM` — penetapan nilai pertanggungan, nilai pertanggungan
liability, premi, dan nilai ter-spreading

**Keputusan.** K-051 · K-018 · K-010/K-012 · K-046

**Blocked by:** E12 · E18

**Status:** blocked — 01-10-2026. ⛔ `[terverifikasi]` `Activity/ReCountPremiLifeEDM.xml` langkah 4.1.4/4.1.5
memilih cabang lewat `.TSILiability<="3000000000"` / `>"3000000000"` — ambang uang tanpa mata uang,
dibandingkan sebagai teks: tepat contoh CLAUDE.md §7 yang wajib menunggu OQ-037/OQ-040/OQ-046 sebelum
angkanya dipindahkan. NB tidak punya lini Life dan tidak berencana memasukkannya ke `MesinPremiFacIn`
(sesi nusantarare-c3), jadi port kelak di modul ini — sesudah OQ itu dijawab.

- [ ] Premi jiwa dihitung lewat **Seam 3**, sebagai **bentuk 1 terparameterisasi** — bukan bentuk baru
- [ ] Basis = nilai pertanggungan **liability**, bukan nilai pertanggungan penuh
- [ ] Rate dari properti rate rata-rata jiwa, **bukan** properti rate umum
- [ ] Pembagi diturunkan **resolver K-018** dari faktor yang ikut
- [ ] Nilai uang bertipe `Money`; rasio bertipe `Ratio`
- [ ] **K-046** `K046_Life_RateSebagaiPembebanan` — properti "rate" adalah faktor pembebanan
- [ ] **K-046** `K046_Life_TreatyType1000036_SpreadingNol` — satu kode jenis treaty memaksa nilai ter-spreading menjadi nol; arti kodenya **`[pertanyaan terbuka]`**, diport apa adanya
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
