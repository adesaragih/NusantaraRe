# E18: Jalur Life sebagai alur tersendiri

**What to build:** Endorsement jiwa berjalan lewat **flow `InputEDMLife` tersendiri** — bukan cabang
kondisi di dalam alur lini umum.

`[terverifikasi]` Empat titik cabang memisahkannya: pengisian nilai lama per baris **keluar** untuk
jiwa · penomoran endorsement memakai pasangan langkah **komplementer** · jalur produksi jiwa berjalan
**sejajar**, bukan di dalam percabangan tujuh lini · penanda baris warisan punya varian sendiri.

⛔ **Jiwa melewati tangga akseptasi** (K-044). Korpus tidak memuat rule tangga akseptasi untuk jiwa.
⚠️ **`[pertanyaan terbuka]`** apakah itu memang tanpa persetujuan — milik work owner dan
Underwriting. Sampai dijawab, **diport apa adanya**.

⚠️ Tiket ini **tidak bersandar** pada klaim arsip bahwa nilai dasar akseptasi adalah selisih TSI —
klaim itu tetap **`[belum diuji]`**.

**Asal (Pega).** `When\IsLife` (enam belas cabang atas kode bisnis lama) ·
`Activity\SetOldData` langkah 1 · `Activity\SaveEDMToJsonPolicy_Act` langkah 12–13, 15, 28–29 ·
`Activity\SetOLDValueToEDMWork_LIFE` · `RDBList\GetKodeProdLife_SQL`

**Keputusan.** **K-044** · K-046 · K-006

**Blocked by:** E06 · E08 · E10

**Status:** blocked

- [ ] Flow **`InputEDMLife` tersendiri** — dimodelkan sebagai alur, bukan cabang `if`
- [ ] Predikat jiwa membaca **kode bisnis lama**, bukan jenis bisnis dari hasil query (berbeda dari predikat lini lain, E01)
- [ ] Jiwa **melewati lapis B** — gerbang keluar E11
- [ ] Penomoran endorsement jiwa lewat query kode produk jiwa; rule generator bernama-jiwa **nihil di korpus** dan itu **bukan `panic`** — penomorannya tertanam di pemanggil
- [ ] Jalur produksi jiwa berjalan **sejajar**, bukan di dalam percabangan tujuh lini
- [ ] ⛔ Jiwa **tanpa tangga akseptasi** — Seam 2 tidak dipanggil di jalur ini
- [ ] Status "tanpa persetujuan" tercatat sebagai **`[pertanyaan terbuka]`**, bukan diputuskan di kode
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
