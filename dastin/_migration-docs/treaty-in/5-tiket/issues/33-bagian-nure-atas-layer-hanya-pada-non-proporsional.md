---
status: aktif
golongan: pelestarian
---

# 33: Bagian NuRe atas sebuah layer dicatat, dan bagian hanya dapat dibuat pada kontrak non-proporsional

*Asal: `DAFTAR-PEKERJAAN.md` `P-17` · `SPEC-MODEL-DATA.md` §10.5, §12.3 · `SPEC-INVARIAN.md` `INV-33`, `INV-64`.*

**What to build:** **PK** mencatat **bagian NuRe** atas sebuah layer — premi bruto dan minimumnya per
mata uang — dan sistem **menolak** pembuatan bagian pada kontrak yang sifat proporsinya
`PROPORSIONAL`.

Artefak: entitas `BAGIAN` dan kedua tabel anak paket uang `NILAI_PREMI_BRUTO` dan
`NILAI_PREMI_BRUTO_MINIMUM`; `INV-64`; `INV-33`.

**PEMBUAT PERTAMA** untuk `BAGIAN`, `NILAI_PREMI_BRUTO`, `NILAI_PREMI_BRUTO_MINIMUM`.

**Kenapa begini:** §12.3 **mengoreksi** §8: `BAGIAN` **bukan** satu entitas untuk kedua cabang. Di
cabang non-proporsional ia kelas tersendiri dengan 49 properti dan baris `Share[]`-nya sendiri; di
cabang proporsional **tidak ada kelas tersendiri sama sekali** — bagian NuRe di sana adalah **atribut
pada baris `DETAIL_PROPORSIONAL`**. Uji yang dipakai semula — *"berbeda arti versus tidak dipakai"* —
dijalankan pada sumbu yang salah; sumbu yang menentukan **masukan versus turunan**, dan pada sumbu itu
kedua sisi tidak sebanding.

**Persyaratan:** `INV-64` (satu `BAGIAN` per `LAYER` — `ID_LAYER` unik di `BAGIAN`) · **`INV-33`**
(`BAGIAN` hanya ada di bawah kontrak ber-`SIFAT_PROPORSI` `NON_PROPORSIONAL`), bergolongan **INDEKS
UNIK** · `INV-36` · `INV-52` (premi bruto dikurangi seluruh potongan sama dengan premi bersih) —
**dibawa sebagai invarian**, penegakannya di irisan `37`

**Tidak termasuk:** **Potongan atas premi** — irisan `37`; `POTONGAN` menggantung pada `BAGIAN`
**dan** `DETAIL_PROPORSIONAL`, sehingga ia menunggu keduanya.
**Penyebaran bagian NuRe ke susunan retro** — irisan `38`.
**Sepuluh daftar `RNMSpreadedList…XOL`** — **turunan, tidak disimpan** (`ADR-0037`, §12.3). Itu
pernyataan keputusan: siapa pun yang menyimpannya menanam kolom turunan yang `INV-58` larang.

**Jalur gagal:** Dua bagian pada satu layer -> **ditolak** `INV-64` · Bagian pada kontrak
`PROPORSIONAL` -> **ditolak** `INV-33`, dan pesannya menyebut sifat proporsi kontraknya · Premi bruto
terisi tanpa mata uang -> ditolak `INV-36`.

**Uji:** **Negatif:** dua bagian pada satu layer; bagian di bawah kontrak proporsional; premi bruto
tanpa mata uang.
**Positif:** satu bagian dengan **dua baris premi bruto bermata uang berbeda** -> diterima.
**Positif kedua — dan ia yang menangkap `INV-33` yang salah arah:** sebuah `DETAIL_PROPORSIONAL` pada
kontrak **proporsional** -> **diterima**. Indeks unik yang ditulis terbalik menolak cabang yang justru
sah, dan ia lulus setiap uji negatif di atas.

**Menggantikan:** `P-17` melestarikan baris `Share[]` cabang non-proporsional. Yang bergeser:
`BAGIAN` **tidak lagi mengaku melayani kedua cabang** — §12.3 mengoreksi §8 — dan premi brutonya
menjadi **tabel anak per mata uang** (`P-8` golongan A) alih-alih satu angka.

**Blocked by:** 31

**Dasar:**
```
EVIDENCED(Data-TreatyInShare@ekspor-2026-09 - 49 properti, baris Share[]; cabang proporsional tanpa kelas tersendiri)
        DECIDED(INV-33, INV-64, ADR-0037, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `BAGIAN`, `NILAI_PREMI_BRUTO`, `NILAI_PREMI_BRUTO_MINIMUM` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-64` terpasang; `INV-33` terpasang **sebagai indeks unik**, dan bentuk itu tertulis sebagai keputusan
- [ ] pesan penolakan `INV-33` menyebut **sifat proporsi kontraknya**
- [ ] kedua uji positif lulus — dua mata uang diterima, **dan** cabang proporsional tidak ikut tertolak
- [ ] sapuan membuktikan **nol** kolom `RNMSpreadedList…` tersimpan
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`
