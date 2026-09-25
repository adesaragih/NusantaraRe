---
status: aktif
golongan: baru
---

# 35: Tepat satu dari persen quota share atau jumlah lines surplus terisi, ditentukan jenis treaty-nya

*Asal: `DAFTAR-PEKERJAAN.md` `P-22` · `SPEC-MODEL-DATA.md` §10.4 · `SPEC-INVARIAN.md` `INV-30`, `INV-31`.*

**What to build:** Pada tiap baris ketentuan proporsional, **tepat satu** dari `PERSEN_QUOTA_SHARE`
dan `JUMLAH_LINES_SURPLUS` terisi, dan **yang mana** ditentukan `JENIS_TREATY` baris itu. Mengisi
keduanya ditolak; mengisi yang tidak sesuai jenisnya ditolak; mengosongkan keduanya ditolak.

**Kenapa begini:** Kedua kolom itu **mengukur hal yang berbeda dengan satuan yang berbeda** —
persentase terhadap jumlah lines — dan keduanya terisi berarti barisnya **tidak dapat ditafsirkan**:
tidak ada yang tahu mana yang dipakai menghitung. Sistem lama membiarkan keduanya terisi, sebab
pemeriksaannya ada di layar saja, dan layar dapat dilewati lewat jalur simpan yang lain. Bentuk
constraint-nya **lintas baris pada satu baris** — `INV-31` menyebutnya demikian — sehingga ia dapat
ditegakkan basis data, bukan aplikasi.

**Persyaratan:** **`INV-31`** (*tepat satu* dari keduanya terisi, ditentukan `JENIS_TREATY`) ·
`INV-30` (`JENIS_TREATY` hanya `QUOTA_SHARE` atau `SURPLUS`) · `INV-41` (persentase di rentang 0–100)

**Tidak termasuk:** **Entitas `DETAIL_PROPORSIONAL`** — irisan `34` yang membuatnya beserta kedua
kolomnya.
**Baris surplus yang menuntut baris quota share** — irisan `36`. Bedanya nyata: yang di sini **satu
baris terhadap dirinya sendiri**, yang di sana **satu baris terhadap saudaranya di versi yang sama**.
**`KAPASITAS_SURPLUS`** — turunan, `INV-48`.

**Jalur gagal:** Keduanya terisi -> **ditolak** `INV-31`, pesannya menyebut jenis treaty barisnya ·
`JENIS_TREATY` `QUOTA_SHARE` dengan `JUMLAH_LINES_SURPLUS` terisi -> ditolak · Keduanya kosong ->
ditolak · `JENIS_TREATY` bernilai selain kedua nilai sah -> ditolak `INV-30`.

**Uji:** **Negatif — empat, bukan satu:** keduanya terisi; keduanya kosong; jenis `QUOTA_SHARE` dengan
kolom surplus; jenis `SURPLUS` dengan kolom quota share.
**Positif:** baris `QUOTA_SHARE` berpersen 30 -> diterima; baris `SURPLUS` berjumlah 9 lines ->
diterima.
**Positif kedua — dan ia yang menangkap constraint yang terlalu ketat:** **satu layer yang memuat
baris `QUOTA_SHARE` dan baris `SURPLUS` sekaligus** -> **diterima**. Susunan quota share bersurplus
adalah bentuk yang lazim; constraint yang ditulis per layer alih-alih per baris menolaknya dan lulus
keempat uji negatif di atas.

**Menggantikan:** tidak ada aturan lama yang digantikan — **golongannya BARU**, dan karena itu
`CARA MENYALAKANNYA` wajib. Sistem lama memeriksanya **hanya di layar**; berapa banyak baris warisan
yang memuat keduanya **tidak diketahui**, dan itu yang mode peringatan hitung.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **pemilik proses Treaty In, mingguan**, atas catatan pelanggaran. Angka pertamanya menjawab pertanyaan yang belum pernah dijawab: berapa baris warisan memuat keduanya |
| 3 | **ambang berangka** — kuasa memblokir dinyalakan ketika pelanggaran pada baris **yang lahir di sistem baru** nol selama **4 minggu berturut-turut**. Baris warisan dihitung terpisah dan tidak menahan penyalaan |
| 4 | **siapa boleh menyalakan** — pemilik proses Treaty In |

**Blocked by:** 34

**Dasar:**
```
EVIDENCED(SPEC-INVARIAN.md INV-31 - "bentuk lintas baris"; pemeriksaan sistem lama hanya di badan seksi)
        DECIDED(INV-30, INV-31)
```

- [ ] `INV-31` terpasang **di basis data**, bukan hanya di aplikasi — dan letaknya tertulis sebagai keputusan
- [ ] keempat uji negatif lulus, masing-masing dijalankan terpisah
- [ ] kedua uji positif lulus — termasuk satu layer bercampur `QUOTA_SHARE` dan `SURPLUS`
- [ ] pesan penolakan menyebut **jenis treaty barisnya**, bukan hanya "salah satu harus kosong"
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk ambang berangka
- [ ] hitungan baris warisan yang melanggar **dilaporkan terpisah** dari baris baru
