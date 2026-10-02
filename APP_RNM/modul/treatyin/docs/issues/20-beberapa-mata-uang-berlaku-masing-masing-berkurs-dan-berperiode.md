---
status: aktif
golongan: perubahan
---

# 20: Lebih dari satu mata uang berlaku pada satu kontrak, masing-masing dengan kurs dan periodenya

*Asal: `DAFTAR-PEKERJAAN.md` `P-06` · `SPEC-MODEL-DATA.md` §10.10 · `ADR-0053`.*

**What to build:** **PK** mencatat **beberapa** mata uang yang berlaku pada satu versi kontrak,
masing-masing dengan **kursnya dan periode berlakunya**. Mata uang yang sama dua kali pada satu versi
ditolak.

Artefak: entitas `MATA_UANG_KONTRAK`, kunci asingnya ke `VERSI_KONTRAK` dan ke tabel acuan
`MATA_UANG`, dan `INV-07`.

**PEMBUAT PERTAMA** untuk `MATA_UANG_KONTRAK`.

**Kenapa begini:** Sistem lama menyimpan mata uang kontrak sebagai **satu skalar** `Currency` di
kepala. Kontrak yang benar-benar bermata uang lebih dari satu **tidak dapat dinyatakan**, dan yang
terjadi adalah orang memilih satu lalu mengonversi sisanya dengan tangan — konversi yang **kursnya
tidak tercatat di mana pun**. `ADR-0053` menjadikan mata uang sebuah **daftar**, dan itu perubahan
bentuk yang tidak dapat dipasang belakangan: skalar yang sudah terisi **tidak menyimpan baris kedua
yang dulu ada**.

**Persyaratan:** `ADR-0053` · `INV-07` (kode mata uang unik di dalam satu versi) · `INV-44` (kode
mata uang merujuk tabel acuan) · `INV-38` (kurs selalu lebih besar dari nol; tidak pernah nol, tidak
pernah dipaksa satu) · `INV-57` (tanggal kurs tidak lebih akhir dari tanggal transaksi yang
memakainya) · `INV-43` (kurs dicatat **sebagai nilai pada versi**, tidak dibaca ulang saat dibutuhkan)

**Tidak termasuk:** **Pembekuan kurs pada versi yang sudah disetujui** — `P-27`, batch 2: titik
bekunya **adalah sebuah keadaan**, dan `REV-3` dapat mengubah daftar keadaan. Irisan ini memasang
`INV-43` sebagai bentuk kolom — kursnya **tersimpan**, bukan dibaca ulang — tanpa menyatakan **kapan**
ia membeku.
**Keterangan kegagalan konversi** — irisan `21`.
**`CurrencyRelation`** — `ADR-0053` sudah memutuskan ia **pasif** dan tidak disimpan.

**Jalur gagal:** Dua baris bermata uang sama pada satu versi -> **ditolak** `INV-07`, pesannya
menyebut kodenya · Kurs bernilai nol atau kosong -> **ditolak** `INV-38` · Kode mata uang yang tidak
ada di tabel acuan -> ditolak `INV-44` · Tanggal kurs lebih akhir daripada tanggal transaksi ->
ditolak `INV-57`.

**Uji:** **Negatif:** sisipkan mata uang kembar; kurs nol; kurs satu yang dipaksakan; kode tak
terdaftar; tanggal kurs di masa depan.
**Positif — dan ia pokok irisan ini:** satu versi dengan **tiga** mata uang, masing-masing berkurs
dan berperiode berbeda, **diterima seluruhnya**. Sebuah rancangan yang diam-diam masih menyimpan mata
uang sebagai skalar lulus setiap uji negatif di atas — sebab satu baris tidak pernah bentrok dengan
dirinya sendiri.

**Menggantikan:** skalar `TreatyIn.Currency` di kepala kontrak, dan pasangan kolom kembar `Rp` / `Usd`
yang §12.4 nyatakan **dibuang** (`INV-46`). Golongannya **PERUBAHAN** — hasilnya **tidak dapat diuji
terhadap data lama**, sebab data lama hanya pernah memuat satu mata uang per kontrak dan itu justru
bentuk yang dibuang.

**Blocked by:** 14 · 15

**Dasar:**
```
EVIDENCED(TreatyIn.Currency@ekspor-2026-09 skalar tunggal; kolom kembar Rp/Usd, SPEC-MODEL-DATA §12.4)
        DECIDED(ADR-0053, INV-46, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `MATA_UANG_KONTRAK` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`, dengan kunci asing ke `VERSI_KONTRAK` dan `MATA_UANG`
- [ ] `INV-07` menolak mata uang kembar di dalam satu versi
- [ ] `INV-38` menolak kurs nol **dan** menolak kurs yang dipaksakan menjadi satu — dua uji, bukan satu
- [ ] `INV-57` terpasang: tanggal kurs tidak dapat lebih akhir daripada transaksi yang memakainya
- [ ] uji positif lulus: tiga mata uang pada satu versi diterima seluruhnya
- [ ] sapuan membuktikan **tidak ada** kolom mata uang skalar yang tertinggal di `VERSI_KONTRAK` selain `KODE_MATA_UANG_KONTRAK` yang §10.2 sebutkan
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`
