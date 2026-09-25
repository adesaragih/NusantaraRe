# E16: Berkas pendukung layar endorsement

**What to build:** Ratusan berkas pendukung khas endorsement — layar, aksi layar, transformasi data,
dan aktivitas pembantu yang tidak punya padanan di New Business.

`[terverifikasi]` **350 berkas diport.** Empat berkas fitur pilih-tertanggung **tidak diport**
(K-044) meskipun ada di korpus — itu keputusan sadar, bukan kelalaian.

⛔ **Tidak satu pun kelompok menjadi modul domain baru.** Yang menyentuh perhitungan masuk ke modul
premi lewat Seam 3; sisanya lapisan tampilan dan pendukung.

⚠️ **Klasifikasi wajib dari tipe rule sebenarnya, bukan nama folder.** Seluruh berkas di folder
bernama "RDBList" ternyata bertipe rule SQL — itu aturan di korpus ini, bukan pengecualian.

**Asal (Pega).** 94 Activity · 81 Section · 54 FlowAction · 29 DataTransform EDM-only

**Keputusan.** K-051 · K-044 · K-046 · K-006 · `CLAUDE.md` §3.3, §4.6

**Blocked by:** E01 · E03

**Status:** blocked

- [ ] **350 berkas diport**; keempat berkas pilih-tertanggung **tidak** — tercatat sebagai keputusan
- [ ] Pasangan layar berakhiran penanda underwriter tetap **dua komponen terpisah** — ikut sistem lama
- [ ] Berkas yang bernama sama dengan New Business **tetapi bertipe rule berbeda** tidak memakai ulang implementasi New Business
- [ ] Klasifikasi dari tipe rule sebenarnya, **bukan** nama folder
- [ ] **Lima berkas tanpa rujukan** diport apa adanya bila terbaca — ⛔ **tidak divonis usang** (K-006); statusnya tetap **`[pertanyaan terbuka]`**
- [ ] Tiap berkas menyebut rule Pega asalnya dalam komentar (§4.6)
