---
status: aktif
golongan: pelestarian
---

# 27: Termin pembayaran premi dicatat bernomor urut, per mata uang

*Asal: `DAFTAR-PEKERJAAN.md` `P-12` · `SPEC-MODEL-DATA.md` §10.16, §10.16a · `SPEC-INVARIAN.md` `INV-12`.*

**What to build:** **PK** mencatat termin pembayaran premi — angsuran bernomor urut beserta tanggal
jatuh temponya dan jumlahnya — **per mata uang**.

Artefak: entitas `TERMIN`, kunci asingnya, kolom `KODE_MATA_UANG`, dan `INV-12`.

**PEMBUAT PERTAMA** untuk `TERMIN`.

**Kenapa begini:** **Kunci alaminya sempat ditulis salah, dan koreksinya adalah alasan irisan ini ada
dalam bentuk ini.** `INV-12` semula berbunyi *"nomor termin unik di dalam satu versi"*. Aktivitas yang
**menambah barisnya** menunjukkan ada satu tingkat pengelompokan di antaranya: daftar yang tampak
"termin" ternyata **baris per mata uang**, dan termin yang sebenarnya hidup satu tingkat di bawahnya —
langkah penambahnya bahkan berlabel jujur, *"Set currency list"*. `UNIQUE` tanpa mata uang **menolak
data yang sah**: kontrak yang menagih empat angsuran dalam dua mata uang punya dua baris bernomor
satu, dan keduanya benar.

> Ini **instans ketiga** dari pola *"lingkup ditulis dari bentuk yang terlihat, bukan dari
> penulisnya"*, dan yang paling bersih. Siapa pun yang menyederhanakan kunci ini mengulanginya.

**Persyaratan:** **`INV-12`** — *"kode mata uang + nomor termin unik di dalam satu versi"*,
**sebagaimana diperbaiki 24 September 2026** · `INV-36` · `INV-44` · `INV-41` (persentase di rentang
0–100, bila terminnya dinyatakan dalam persen)

**Tidak termasuk:** **`JUMLAH_TERMIN` dan `MEMAKAI_PRORATA` di kepala versi** — keduanya kolom
`VERSI_KONTRAK`, lahir bersama irisan `14`.
**Mesin pro rata** — `GRL-15`, **sengaja tidak dibangun**. Pernyataan keputusan, bukan `TODO`.
**Rekonsiliasi jumlah termin terhadap premi** — `INV-47` lewat *materialized view*, irisan `38`.

**Jalur gagal:** Dua baris bernomor termin **dan** mata uang sama -> **ditolak** `INV-12` · Nilai
termin terisi tanpa mata uang -> ditolak `INV-36` · Nomor termin nol atau negatif -> ditolak.

**Uji:** **Negatif:** baris kembar pada sumbu penuh; nilai tanpa mata uang; nomor termin nol.
**Positif — dan ia yang menangkap kunci yang terlalu sempit, sekaligus membuktikan koreksi `INV-12`
memang diterapkan:** satu versi dengan **termin nomor 1 dalam IDR dan termin nomor 1 dalam USD** ->
**diterima**. Uji inilah yang gagal pada rumusan `INV-12` yang lama, dan ia harus dijalankan — bukan
diargumentasikan.

**Menggantikan:** `P-12` melestarikan daftar termin. Yang bergeser **kunci alaminya**, dari *"nomor
termin"* menjadi *"mata uang + nomor termin"* — dan bunyi lamanya disebut di sini justru supaya yang
sudah membacanya tahu apa yang bergeser.

**Blocked by:** 14 · 15

**Dasar:**
```
EVIDENCED(aktivitas penambah baris termin@ekspor-2026-09 - langkah "Set currency list"; SPEC-MODEL-DATA §10.16a)
        DECIDED(INV-12, GRL-15, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `TERMIN` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md` dengan `KODE_MATA_UANG` pada barisnya
- [ ] `INV-12` terpasang sebagai `UNIQUE` **bertiga kolom**, dan mata uangnya ada di dalamnya
- [ ] uji positif lulus: termin nomor 1 dalam dua mata uang **diterima** — dijalankan, bukan diargumentasikan
- [ ] ketiadaan mesin pro rata tertulis sebagai **pernyataan keputusan** yang menyebut `GRL-15`
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`
