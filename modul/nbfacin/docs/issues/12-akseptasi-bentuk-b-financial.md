# 12: Akseptasi lini financial — filter ambang limit, tanpa tangga berjenjang

> ⚠️ **AMANDEMEN 1 Oktober 2026** `[keputusan work owner]` (`../KEPUTUSAN-30-09-2026.md` butir 46–47).
> Premis badan tiket di bawah — **"tidak bereskalasi", "keluaran daftar jabatan", "tanpa antrean"** — **dibantah
> korpus**: langkah 22 `GetLimitAkseptasi_ActFlow` (dan langkah 23 `GetLimitAkseptasi_Act`) naik **satu tingkat**
> menurut antrean saat ini (UW Financial → Kadiv Keuangan → Direktur Marketing → Direktur Teknik), lalu Decision23
> merutekannya. Yang tetap berlaku: satu query per jenis pertanggungan, `WHERE` apa adanya, ejaan berspasi, dua
> jalur terpisah, fixture tanpa nama. Langkah 17 dibuang (butir 47). Teks asli di bawah dipertahankan sebagai
> riwayat; status tiap kriteria ada di bab Comments.

**What to build:** Kasus lini financial menemukan jabatan berwenangnya lewat **penyaringan ambang
limit**, bukan lewat tangga bertingkat. Lini ini **tidak bereskalasi**.

⛔ **Bentuk tabelnya berbeda, dan menyamakannya membuat tangga financial macet total.** Tabel limit
financial **tidak punya** kolom jabatan atasan, tidak punya kolom antrean, dan kolom limitnya bukan
ambang tunggal melainkan **empat ambang per jenis pertanggungan**: bond, credit CL, credit NCL, trade.

⛔ **Ejaan jabatannya PAKAI SPASI**, berbeda dari bentuk standar yang tanpa spasi. Bila kode
mencocokkan dengan token tanpa spasi, **tidak ada approver financial yang pernah ditemukan**.
**Normalisasi ejaan dilarang** — menghapus spasi agar seragam adalah perubahan perilaku, bukan
migrasi.

**Yang tidak ada** adalah eskalasi berjenjang dan antrean per jabatan. **Yang tetap ada dan wajib
diport** adalah **filter ambangnya**: tiga rule SQL menyaring dengan `WHERE <kolom_limit> <= nilai`,
masing-masing untuk satu jenis pertanggungan, tanpa penggabungan tabel dan tanpa pengurutan.
Membuang `WHERE` menghapus satu-satunya kontrol wewenang yang dimiliki lini ini.

Keluarannya karena itu **daftar jabatan yang limitnya menampung nilai** — bukan satu jabatan tujuan
berikutnya.

⚠️ Kedua mekanisme **tidak disatukan di balik satu abstraksi**. Memaksa bentuk ini ke dalam bentuk
tangga berarti mengarang langkah naik yang tidak ada.

**Blocked by:** 11

**Status:** ready-for-human — diport 01-10-2026 dengan premis diamandemen (butir 46), menunggu tinjauan A26–A27

- [ ] Satu query per jenis pertanggungan (bond, credit CL, credit NCL), memakai kolom ambangnya masing-masing
- [ ] **`WHERE` dipertahankan apa adanya**; tidak ada penyaringan tambahan di luar ambang limit
- [ ] Pencocokan jabatan memakai ejaan **berspasi**; tidak ada normalisasi
- [ ] Keluaran berupa **daftar jabatan**, bukan jabatan tujuan berikutnya
- [ ] Tidak ada langkah naik dan tidak ada token antrean untuk lini ini
- [ ] Kedua bentuk tetap dua jalur terpisah di kode; kasus uji membuktikan bentuk standar **tidak** dipakai untuk lini financial dan sebaliknya
- [ ] Fixture bentuk ini juga **tanpa kolom nama**

## Comments

### 2026-10-01 — ditahan (agent)

Bergantung pada NB-11, dan memakai tabel `M_LIMIT_FINANCIALINS` yang isinya juga berformat `.xls` biner
(`../PERTANYAAN-AKSEPTASI.md` B4). Tiga query bentuk B sudah dikutip K-026.

### 2026-10-01 — CSV `M_LIMIT_FINANCIALINS` diterima (butir 42)

`backend/services/acceptance/testdata/limit/M_LIMIT_FINANCIALINS.csv`: 4 baris, kolom `JABATAN`,
`LIMITBOND_BOTTOM`, `LIMITCREDITCL_BOTTOM`, `LIMITCREDITNCL_BOTTOM`. `JABATAN` berspasi — cocok dengan literal
`.CARI1` langkah 22.2.2–22.2.4 `GetLimitAkseptasi_ActFlow`.

### 2026-10-01 — diport, premis tiket diamandemen (butir 46, 47)

⚠️ Premis "tidak bereskalasi / keluaran daftar jabatan / tanpa antrean" **dibantah korpus**: langkah 22
`GetLimitAkseptasi_ActFlow` (dan langkah 23 `GetLimitAkseptasi_Act`) naik satu tingkat menurut antrean saat ini,
lalu Decision23 merutekannya. Work owner memutuskan ikuti korpus. Kode: `backend/services/acceptance/tangga_financial.go`
— `NextFinancial(kasus, tabel, pengguna) (Transisi, error)`.

| Kriteria (asli) | Keadaan |
| --- | --- |
| Satu query per jenis pertanggungan, kolom ambangnya masing-masing | ✅ `kolomBentukB`; Trade Credit memakai SQL Kredit CL (langkah 13) |
| `WHERE` dipertahankan apa adanya | ✅ `<kolom> <= CARID2`, inklusif (uji batas `TSI tepat di limit`) |
| Ejaan berspasi, tanpa normalisasi | ✅ `TestNextFinancialEjaanBerspasi` |
| Keluaran daftar jabatan | ⛔ diamandemen: keluaran `Transisi` satu langkah (butir 46) |
| Tidak ada langkah naik / antrean | ⛔ diamandemen: tiga tingkat + antrean Decision23 (butir 46) |
| Dua jalur terpisah; uji membuktikan tidak saling dipakai | ✅ `TestDuaJalurTerpisah` (`ErrBentukB` / `ErrBukanBentukB`) |
| Fixture tanpa kolom nama | ✅ `TestFixtureLimitTanpaNamaLogin` |

Tidak diport: langkah 17 (butir 47), langkah 20 (LetterNo-nya ditimpa 22.1), `Data.LetterNo`, `NBStatus` (teks
berisi CARI5 = `NAMA`). Uji mutasi `tangga_financial.go` 13/14; satu yang lolos mutan ekuivalen (memakai tafsir
kolom ketika kedua tafsir sama — kode sudah menolak saat berbeda).

### 2026-10-01 — tinjauan dua sumbu (`/code-review`) diterapkan

- **Urutan langkah:** anggota grup kini keluar di langkah 3 **sebelum** langkah 4–5 membaca top risk dan
  `OldData` (sebelumnya data rusak membuat galat) — berlaku juga untuk `Next` (bentuk A).
  `TestLangkah3SebelumLangkah4Dan5`.
- **A27 dicabut** (kode mati); **A28**: kolom limit kosong → `ErrLimitKosong`.
- Uji daftar kosong, batas inklusif `WHERE`, konstanta antrean, bukti `class!NAMA` + perintah audit untuk klaim
  `LIMITTRADE_BOTTOM` dan kode `BusinessOldId` saling lepas.
- Temuan "berkas ber-CRLF" dari tinjauan **positif palsu**: `grep $'\r$'` di Git Bash memberi hitungan = jumlah
  baris; pengurai byte (Python, `b"\r\n"`) dan `git ls-files --eol` (`w/lf`) sepakat semua berkas LF tanpa BOM.

Uji mutasi: `tangga_financial.go` + perubahan `tangga.go` **16/16**; `tangga.go` (bentuk A) **20/20**.
