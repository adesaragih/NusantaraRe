---
status: aktif
golongan: perubahan
---

# 28: Skala ko-asuransi dicatat sebagai beberapa baris, bukan satu nilai

*Asal: `DAFTAR-PEKERJAAN.md` `P-13` · `SPEC-MODEL-DATA.md` §10.17 · `SPEC-INVARIAN.md` `INV-13`.*

**What to build:** **PK** mencatat skala ko-asuransi sebagai **beberapa baris** — tiap baris satu
persen limit dengan persen ko-asuransinya — bukan sebagai satu angka di kepala kontrak.

Artefak: entitas `SKALA_KOASURANSI`, kunci asingnya ke `VERSI_KONTRAK`, dan `INV-13`.

**PEMBUAT PERTAMA** untuk `SKALA_KOASURANSI`.

**Kenapa begini:** Skala ko-asuransi **adalah sebuah tangga**: persentase yang berbeda berlaku pada
lapis limit yang berbeda. Sistem lama punya daftarnya — `CoInScale[]` — **dan juga** sebuah skalar
bernama sama di kepala, `TreatyIn.CoInScale`, yang **nol penulis** dan hanya muncul di empat
harness/section serta **tidak terdeklarasi sebagai properti**. Skalar itu **penampung layar**, bukan
atribut, dan `COCOK-SILANG-CACAH-ATRIBUT.md` §4 mengadilinya demikian. Membawanya sebagai satu nilai
akan meratakan tangga menjadi satu anak tangga.

**Persyaratan:** `INV-13` (persen limit unik di dalam satu versi) · `INV-41` (persentase di rentang
0–100) · `INV-01`

**Tidak termasuk:** **Skalar `TreatyIn.CoInScale`** — **tidak dibawa**, dan itu putusan bernomor:
nol penulis atas kelima bentuk penulis, tak terdeklarasi, hanya di layar. Ini **pernyataan
keputusan**; siapa pun yang menambahkannya kembali membalikkan adjudikasi `F-2`.
**Perhitungan bagian ko-asuransi** — turunan, `ADR-0037`.

**Jalur gagal:** Dua baris berpersen limit sama pada satu versi -> **ditolak** `INV-13` · Persen di
luar 0–100 -> ditolak `INV-41` · Sebuah kolom skalar ko-asuransi di `VERSI_KONTRAK` -> **tidak ada**,
dan ketiadaannya diperiksa dengan sapuan.

**Uji:** **Negatif:** persen limit kembar; persen 101; persen negatif.
**Positif — dan ia pokok irisan ini:** satu versi dengan **empat baris skala** berpersen limit menaik
-> **diterima seluruhnya**, dan keempatnya terbaca kembali **dalam urutan yang sama**. Sebuah
rancangan yang diam-diam masih meratakannya menjadi satu nilai lulus ketiga uji negatif di atas —
satu baris tidak pernah bentrok dengan dirinya sendiri.

**Menggantikan:** skalar `TreatyIn.CoInScale` di kepala. Golongannya **PERUBAHAN**: bentuknya sengaja
diubah dari satu nilai menjadi daftar, sehingga **cocok-tidaknya terhadap data lama berarti
kebalikan** dari yang biasa.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(datar-treatyin-lama.csv baris 172 - TreatyIn.CoInScale Skalar, "(tidak dideklarasikan)", 0 penulis)
        EVIDENCED(COCOK-SILANG-CACAH-ATRIBUT.md §4 - adjudikasi kelima calon F-2)
        DECIDED(INV-13, ADR-0037, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `SKALA_KOASURANSI` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-13` terpasang atas persen limit di dalam satu versi
- [ ] `INV-41` terpasang dan diuji pada **kedua** ujung rentang
- [ ] uji positif lulus: empat baris skala diterima dan terbaca kembali berurutan
- [ ] sapuan membuktikan **tidak ada** kolom skalar ko-asuransi di `VERSI_KONTRAK`
- [ ] ketidakbawaan skalar `CoInScale` tertulis sebagai **pernyataan keputusan** yang menyebut buktinya
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`
