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

**Status:** sebagian — 01-10-2026, `backend/services/jalurlife.go`; test `jalurlife_test.go`. Penghalang
E06/E08/E10 gugur. ⛔ Flow `InputEDMLife` **tidak ada di korpus** (nol berkas) — Life bercabang di
`Flow/InputAddendumFacultativeIn.xml` shape `Decision19`; yang dimodelkan cabang itu (A03). Sisa: jalur
produksi Life (E21, tabel flat) dan pelaksanaan query penomoran (E17).

- [ ] Flow **`InputEDMLife` tersendiri** — ⛔ **premis dibantah korpus**; cabang `Decision19` dimodelkan (`LangkahSesudahKonfirmasiMarketing`)
- [x] Predikat jiwa membaca **kode bisnis lama** — `[terverifikasi]` `When/IsLife.xml`: `pyWorkPage.Quotation.BusinessOldId = "L1"…"L16"`. Dinilai registry NB lewat `kontrak.PenilaiPredikatFacIn`; di EDM ia masukan. ⚠️ `pyConditionViewer` rule itu memuat cache editor `BusinessCode = "10164"` — bukan yang dieksekusi
- [x] Jiwa **melewati lapis B** — gerbang keluar E11
- [x] Penomoran endorsement jiwa lewat query kode produk jiwa (`RencanaNomorEndorsemen`, `SusunNomorEndorsemen`); `GenerateEDMNoLife` nihil, bukan `panic`. ⚠️ `GenerateEndorsementNo` (non-Life) **juga** nihil di korpus
- [ ] Jalur produksi jiwa berjalan **sejajar** — E21 (tabel flat P-10)
- [x] ⛔ Jiwa **tanpa tangga akseptasi** — `LangkahCekGalatKonversi.MembukaTanggaAkseptasi() == false`
- [x] Status "tanpa persetujuan" tercatat sebagai **`[pertanyaan terbuka]`**, bukan diputuskan di kode
- [x] Menyebut rule Pega asalnya dalam komentar (§4.6)
