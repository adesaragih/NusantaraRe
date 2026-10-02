---
status: aktif
golongan: perubahan
---

# 15: Himpunan acuan bertambah tanpa mengubah arti apa pun yang sudah tercatat

*Asal: `DAFTAR-PEKERJAAN.md` `P-48` · `SPEC-MODEL-DATA.md` §10.22 · `ADR-0038`.*

**What to build:** Keenam himpunan yang dapat bertambah — **mata uang, jenis potongan, jenis
reasuransi, bahaya, kelompok treaty, kelas bisnis** — berdiri sebagai **tabel acuan**, dan menambah
satu baris baru **tidak menyentuh satu pun kontrak yang sudah tercatat** dan **tidak menuntut
perubahan skema**.

Artefak: keenam tabel acuan, kunci alaminya, dan penanda aktifnya.

**PEMBUAT PERTAMA** untuk `MATA_UANG`, `JENIS_POTONGAN`, `JENIS_REASURANSI`, `BAHAYA`,
`KELOMPOK_TREATY`, `KELAS_BISNIS`.

**Kenapa begini:** Sistem lama menaruh **empat bahaya sebagai empat pasang kolom di kepala kontrak** —
`Earthquake`, `FloodJab`, `FloodNation`, `RSMDLimit`, masing-masing berpasangan dengan kolom mata
uangnya. Bentuk itu menuntut **kolom kesembilan** begitu ada bahaya baru, dan bahaya baru adalah hal
yang pasti terjadi di reasuransi. `ADR-0038` membalikkannya: **yang bertambah tanpa mengubah arti
disimpan sebagai data, bukan sebagai nama kolom dan bukan sebagai `CHECK`.**

**Persyaratan:** `INV-62` (himpunan yang dapat bertambah disimpan sebagai tabel acuan, bukan `CHECK`) ·
`INV-68` (`KODE` unik di seluruh tabel, pada keenamnya) · `INV-44` (kode mata uang merujuk tabel acuan
mata uang) · `INV-01` · `INV-02` · `ADR-0038`

**Tidak termasuk:** **Pemindahan isinya dari sistem lama** — irisan `44`. Irisan ini membuat
tabelnya dan membuktikan ia dapat bertambah; mengisinya dari sumber lama pekerjaan lain.
**`JENIS_REASURANSI` bersusun** — kolom `ID_INDUK` yang menunjuk baris lain di tabel yang sama
**ikut di sini**, tetapi arti susunannya dipakai pertama kali oleh irisan `38`.
**Pemilik tabel acuan pembagian kapasitas** — `P-49`, dan `G2` menyatakannya **tegas tidak masuk**
penyerahan pertama: ia mewujudkan butir eskalasi 7 yang belum diputuskan.

**Jalur gagal:** Dua baris berkode sama di satu tabel acuan -> **ditolak** `INV-68`, pesannya
menyebut kodenya · Kode mata uang pada kontrak yang tidak ada di tabel acuan -> ditolak `INV-44` ·
Sebuah `CHECK` berisi daftar nilai bahaya -> **tidak ada**, dan ketiadaannya diperiksa, bukan
diandaikan.

**Uji:** **Negatif:** sisipkan kode ganda di tiap tabel; rujuk mata uang yang tidak terdaftar.
**Positif — dan ia pokok irisan ini:** tambahkan **bahaya kesembilan** ke tabel acuan, lalu catat
batas tanggungan untuknya pada sebuah kontrak. Keduanya **berhasil tanpa satu pun perubahan skema**.
Uji negatif saja tidak dapat membuktikannya: sebuah tabel acuan yang benar dan sebuah tabel acuan
yang `CHECK`-nya masih tertinggal di tempat lain **menolak hal yang sama**.

**Menggantikan:** `P-48` tidak menggantikan aturan bernomor cacat, melainkan **bentuk** yang §10.0b
uraikan: delapan kolom bernama bahaya di kepala kontrak. Golongannya **PERUBAHAN** — sistem lama
melakukannya, dan kita sengaja mengubahnya — sehingga hasilnya **tidak dapat diuji terhadap data
lama**: data lama justru menunjukkan bentuk yang dibuang.

**Blocked by:** None (can start immediately)

**Dasar:**
```
EVIDENCED(SPEC-MODEL-DATA.md §10.0b - delapan kolom bernama bahaya, dibaca dari kepala TreatyIn)
        DECIDED(ADR-0038, INV-62, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] keenam tabel acuan berdiri dengan bentuk seragam — kode, nama, penanda aktif — sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-68` menolak kode ganda di **keenamnya**, diperiksa satu per satu bukan disimpulkan dari satu
- [ ] **nol `CHECK` berisi daftar nilai** untuk keenam himpunan ini; diperiksa dengan sapuan, hasilnya dicetak
- [ ] `JENIS_REASURANSI` membawa `ID_INDUK` yang menunjuk barisnya sendiri, dengan perilaku hapusnya ditetapkan sadar
- [ ] uji positif lulus: bahaya kesembilan ditambahkan dan langsung dapat dipakai, **tanpa perubahan skema**
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`
