---
status: aktif
golongan: pelestarian
---

# 26: Periode akumulasi dicatat, dan ditolak bila di luar periode kontrak

*Asal: `DAFTAR-PEKERJAAN.md` `P-11` · `SPEC-MODEL-DATA.md` §10.15 · `SPEC-INVARIAN.md` `INV-11`, `INV-56`.*

**What to build:** **PK** mencatat periode akumulasi sebuah versi — jendela waktu yang beberapa
kejadian di dalamnya dihitung sebagai **satu kejadian**. Periode yang jatuh di luar periode
kontraknya ditolak.

Artefak: entitas `PERIODE_AKUMULASI`, kunci asingnya, `INV-11`, dan `INV-56`.

**PEMBUAT PERTAMA** untuk `PERIODE_AKUMULASI`.

**Kenapa begini:** Periode akumulasi **menentukan berapa besar klaim yang dibayar**, sebab ia
menentukan apakah dua kejadian dijumlahkan ke satu deductible atau ke dua. Sebuah periode akumulasi
yang jatuh di luar periode kontrak **tidak dapat memuat satu kejadian pun yang tertanggung**, dan
karena itu ia selalu kekeliruan pengisian — bukan pilihan yang sah.

**Persyaratan:** `INV-11` (periode unik di dalam satu versi) · **`INV-56`** (periode akumulasi berada
di dalam periode kontraknya), bergolongan **APLIKASI** · `INV-53` · `ADR-0022` (kedua batas inklusif)

**Tidak termasuk:** **Perhitungan akumulasi kejadian itu sendiri** — ia perilaku klaim, dan modul
klaim tidak di dalam penyerahan ini. Yang dicatat **jendelanya**, bukan hasilnya.
**`PERIODE_AKUMULASI_KONTRAK` di kepala versi** — kolom `VERSI_KONTRAK`, lahir bersama irisan `14`.

**Jalur gagal:** Dua baris berperiode sama pada satu versi -> **ditolak** `INV-11` · Periode di luar
periode kontrak -> **ditolak** `INV-56`, pesannya menyebut periode kontraknya · Periode bertanggal
terbalik -> ditolak `INV-53`.

**Uji:** **Negatif:** periode kembar; periode yang mulai sehari sebelum kontrak; periode bertanggal
terbalik.
**Positif — dan ia yang menangkap kekeliruan lingkup:** satu versi dengan **dua periode akumulasi
yang saling bertumpang tindih** tetapi berperiode berbeda -> **diterima**. Tumpang tindih **tidak
dilarang** oleh satu pun invarian, dan pemeriksaan yang melarangnya akan menolak susunan yang sah.
Ini uji yang memisahkan `INV-11` — keunikan periode — dari larangan tumpang tindih yang **tidak pernah
diputuskan siapa pun**.

**Menggantikan:** `P-11` melestarikan daftar `AccumulationList`. Yang bergeser: rentangnya diperiksa.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(TreatyIn.AccumulationList@ekspor-2026-09; AccumulationPeriod skalar di kepala)
        DECIDED(INV-11, INV-56, ADR-0022, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `PERIODE_AKUMULASI` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-11` terpasang; `INV-56` terpasang di lapisan aplikasi, dan letaknya tertulis sebagai keputusan
- [ ] uji positif lulus: dua periode bertumpang tindih **diterima**
- [ ] ketiadaan larangan tumpang tindih tertulis sebagai **pernyataan keputusan** — supaya orang berikutnya tidak menambahkannya sebagai "kerapian"
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`
