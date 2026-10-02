# E15: Selisih per baris spreading dan pemetaan kolomnya

**What to build:** Untuk setiap baris spreading, hitung selisih antara nilai sesudah dan nilai
sebelum, lalu petakan ke pasangan kolom produksi.

⛔ Tiket ini **menghasilkan nilai**; **menuliskannya ke tabel produksi milik E21** yang masih
ter-block.

`[terverifikasi]` Rumus kanoniknya: nilai baru **diprorata lebih dulu**, baru dikurangi nilai lama
yang **jenis treaty-nya sama**.

⛔ **Nilai lama diambil dari lapis A**, bukan dari nilai lama per baris. Ini titik yang paling mudah
salah diimplementasikan.

⛔ **Pasangan dicocokkan menurut jenis treaty, bukan indeks posisi.** Penambahan atau penghapusan
baris tidak boleh menggeser pasangan.

**Asal (Pega).** `Activity\SaveFacinProdEDMFire_Act` (puluhan penugasan pengurangan) dan lima cabang
lini lainnya · `RDBList\InsertTreatyProduction_Sql` untuk pemetaan kolom

**Keputusan.** K-048 · K-027 · K-046 · K-010/K-012

**Blocked by:** E12 · E06

**Status:** blocked — 01-10-2026. Penghalang E06 gugur; tetap ter-block **E12** (bentuk premi) dan
penulisannya E21.

- [ ] Selisih = (nilai baru × porsi sisa periode) − nilai lama berjenis treaty sama
- [ ] Nilai lama dari **lapis A** — kasus uji dengan lapis A dan lapis B sengaja **berbeda**
- [ ] Pencocokan menurut **jenis treaty**; kasus uji dengan baris ditambah dan dihapus di tengah
- [ ] Tujuh varian menyimpang terimplementasi: penyesuaian spreading · penyesuaian mata uang (dua arah) · penyesuaian rate · coverage hilang → negatif · jenis endorsement pembatalan → balik tanda
- [ ] Sembilan pasangan kolom terpetakan
- [ ] Nilai berkoma desimal dibaca sesuai **K-027**
- [ ] **K-046** `K046_PasanganMenjadiSelisih_TidakSetangkup` — kolom "sesudah" lebih sedikit daripada kolom "selisih"; tiga pasangan memakai kolom tanpa sufiks
- [ ] **K-046** `K046_EjaanKolomPCT_BROKERGARE_FEE_Dipertahankan` — kolom persen brokerage salah eja di skema; skema tidak berubah (§4.3)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
