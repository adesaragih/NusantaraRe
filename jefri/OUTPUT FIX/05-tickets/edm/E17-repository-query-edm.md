# E17: Lapisan query dan lookup endorsement ⛔ BLOCKED

**What to build:** Query, halaman data, dan definisi laporan khas endorsement — lapisan yang membaca
Oracle untuk kebutuhan layar dan perhitungan.

⛔ **BLOCKED.** Seam repository **belum ada**, dan bentuknya belum dapat dirancang tanpa skema yang
sudah pasti. Ini **kewajiban tertunda sejak spec New Business**, bukan kelalaian putaran ini.

`[terverifikasi]` Termasuk di dalamnya **tarif travel**, yang rule query-nya **nihil di seluruh
korpus**. Struktur tabelnya tidak dapat diketahui dari korpus dan **tidak ditebak**.

**Asal (Pega).** 22 rule SQL (folder `RDBList\`) · 14 halaman data · 28 definisi laporan EDM-only

**Keputusan.** K-051 · K-027 · `CLAUDE.md` §4.3, §4.5

**Blocked by:** ⛔ **seam repository** (belum ada)

**Status:** blocked

- [ ] ⛔ Menunggu seam repository dirancang
- [ ] Tarif travel menjadi **masukan dari repository** — bukan dihitung di modul premi (lihat E12 bentuk 3)
- [ ] Rule query tarif travel **nihil di korpus** → **`[pertanyaan terbuka]`**, jangan ditebak
- [ ] Tabel pita tarif personal accident juga **belum diketahui** strukturnya
- [ ] Seluruh nilai berkoma desimal dibaca sesuai **K-027**
- [ ] Endpoint dan host dari **konfigurasi**, tidak pernah literal (`CLAUDE.md` §4.4)
- [ ] Skema Oracle **tidak berubah** (§4.3) — termasuk nama kolom yang salah eja
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
