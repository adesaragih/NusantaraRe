# 10: Jembatan predikat lini bisnis → resolver skala

**What to build:** Lini bisnis sebuah kasus **ditentukan oleh predikat yang sesungguhnya dipakai
sistem lama**, lalu mengalir ke resolver skala — bukan diisi pemanggil sebagai parameter yang
diasumsikan benar.

Ini titik integrasi antara registry predikat dan perhitungan premi: dua bagian yang sudah bekerja
sendiri-sendiri, disambungkan di satu tempat yang dapat diuji. Predikat yang terlibat mengenali
FIRE, PA, MBU, ANEKA, MARINE CARGO, GOLF, dan BONDING.

⚠️ **Kehati-hatian yang diwarisi dari discovery.** Lini bisnis ditentukan satu properti, **tetapi dua
predikat membacanya lewat jalur berbeda di dalam agregat** — sebagian memakai dua sampai tiga jalur.
Apakah ketiga jalur selalu sinkron **belum terjawab**. Sampai terjawab, jalur-jalur itu
**dipertahankan apa adanya** dan perbedaannya dicatat saat runtime; properti itu **tidak boleh
dinormalisasi** menjadi satu field.

⚠️ Sebagian nilai lini bisnis muncul **dengan spasi di depan atau belakang**. Dipangkas di batas
input, dan dicatat sebagai kandidat perbaikan.

**Blocked by:** 08, 03

**Status:** ready-for-human — diport 01-10-2026, menunggu tinjauan work owner

- [ ] Lini bisnis diturunkan lewat registry predikat, bukan diterima mentah dari pemanggil
- [ ] Hasilnya mengalir ke resolver skala; **tidak ada** jalur yang melewati resolver
- [ ] Jalur baca ganda pada properti penentu **dipertahankan**, tidak dinormalisasi
- [ ] Ketidaksinkronan antar jalur **tercatat saat runtime** alih-alih dipilih diam-diam salah satunya
- [ ] Spasi di ujung nilai dipangkas di batas input, dicatat sebagai kandidat perbaikan
- [ ] Kasus yang lini bisnisnya tidak dikenali **gagal keras** lewat resolver (tiket 03), bukan memakai default

## Comments

### 2026-10-01 — implementasi (agent)

Kode: `premium/jembatan.go` — `LiniDariPredikat(kasus) (HasilLini, error)`.

`[terverifikasi]` Gerbang yang sungguh dipakai sistem lama untuk memilih rumus premi:
`CountGrossPremi_Act` langkah 5.1 `IsFire`, 5.2 `IsAneka`, 5.3 `isGolfInsurance`, 5.4 →
`CountGPWMarinePAMbu_Act` langkah 1.1 `IsMarineCargo`, 1.2 `IsMBU`, 1.3 `IsPA`. Gerbangnya
**berdiri sendiri** — satu kasus bisa membuka lebih dari satu — jadi hasilnya **daftar**, urut langkah
(`../KEPUTUSAN-30-09-2026.md` bab "Keputusan agent — menunggu konfirmasi").

| Kriteria | Keadaan |
| --- | --- |
| Lini lewat registry predikat | ✅ |
| Mengalir ke resolver | ✅ setiap lini membawa `Satuan` dari `SatuanRate` |
| Jalur baca ganda dipertahankan | ✅ empat jalur `BusinessType` dibaca apa adanya |
| Ketaksinkronan tercatat saat runtime | ✅ `HasilLini.Peringatan` |
| Spasi dipangkas di batas input | ✅ + dicatat sebagai peringatan (kandidat perbaikan) |
| Tak dikenali → gagal keras lewat resolver | ✅ panic |

⚠️ Temuan: (1) BONDING tidak punya gerbang sendiri — kasus Bonding membuka `IsAneka` lewat
`IsBondingAndCustomBonds`. (2) Langkah 5/1 juga bergerbang `StatusBusiness != 3` — milik alur
pemanggil, tidak diport. (3) Berkas kasus nyata menyimpan `QuotationData.BusinessType` di bawah
OfferFacIn, sedangkan `IsPA` membaca `pyWorkPage.Quotation.BusinessType` — halaman yang tidak ada di
berkas kasus (`../PERTANYAAN-LANJUTAN.md` butir 5). (4) Beberapa predikat di rantai `IsAneka` membandingkan
`BusinessCode` dengan literal angka, jadi kasus tanpa kode bisnis ditolak (butir 20). Uji mutasi: 4/4.

### 2026-10-01 — tindak lanjut review (agent)

- 🐞 Diperbaiki: bila tidak satu gerbang terbuka tetapi `BusinessType` kebetulan berisi nama lini resolver
  (mis. `"BONDING"`), jembatan dulu mengembalikan daftar kosong tanpa galat. Kini panic langsung
  (`TestLiniDariPredikatTakDikenaliWalauNamaLini`).
- Urutan evaluasi `append` diperbaiki; pemangkasan spasi kini teruji tercatat
  (`TestLiniDariPredikatPangkasDicatat`).
- ⏸ Kriteria *"bukan diterima mentah dari pemanggil"* baru sebagian: `Calculate` masih menerima
  `Input.LiniBisnis` dari pemanggil; `LiniDariPredikat` tersedia tetapi belum disambungkan — penyambungnya
  alur pemanggil (yang juga memegang gerbang `StatusBusiness != 3`, keputusan agent A2).
- Bonding → ANEKA dicatat sebagai keputusan agent A11.

### 2026-10-01 — `StatusBusiness = 3` = Endorsement (keputusan work owner, butir 29)

Gerbang `StatusBusiness != 3` di atas `CountGrossPremi_Act` langkah 5 / `CountGPWMarinePAMbu_Act`
langkah 1 berarti premi endorsement tidak lewat pemilih rumus NB. Tetap tidak diport di jembatan (A2,
dikonfirmasi): milik alur pemanggil. A11 (Bonding → ANEKA) dikonfirmasi. Pertanyaan jalur ganda
`pyWorkPage.Quotation` vs `pyWorkPage.OfferFacIn.QuotationData` tetap terbuka.

### 2026-10-01 — jalur ganda Quotation terjawab (butir 49)

`pyWorkPage.Quotation` dan `pyWorkPage.OfferFacIn.QuotationData` berisi sama; pemuat mengisi `Quotation.*` dari
`QuotationData.*`. Peringatan beda-jalur di `LiniDariPredikat` tetap dipertahankan sebagai penjaga data.
