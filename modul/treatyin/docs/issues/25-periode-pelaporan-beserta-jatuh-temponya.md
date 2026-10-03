---
status: aktif
golongan: pelestarian
---

# 25: Periode pelaporan dicatat beserta tanggal jatuh temponya, dan ditolak bila di luar periode kontrak

*Asal: `DAFTAR-PEKERJAAN.md` `P-10` · `SPEC-MODEL-DATA.md` §10.14 · `SPEC-INVARIAN.md` `INV-10`, `INV-55`.*

**What to build:** **PK** mencatat periode pelaporan sebuah versi beserta **batas penyerahan, batas
konfirmasi, batas pelunasan, dan hari pengingatnya**. Periode yang jatuh **di luar periode
kontraknya ditolak**.

Artefak: entitas `PERIODE_PELAPORAN`, kunci asingnya, `INV-10`, dan `INV-55`.

**PEMBUAT PERTAMA** untuk `PERIODE_PELAPORAN`.

**Kenapa begini:** Periode pelaporan yang jatuh di luar periode kontrak **tidak dapat dipenuhi
siapa pun** — tidak ada premi yang dilaporkan untuk waktu yang kontraknya belum atau sudah tidak
berlaku. Sistem lama menyimpan keempat batas hari itu sebagai skalar di kepala (`ReportingSubmission`,
`ReportingConfirmation`, `ReportingSettlement`, `ReminderDays`) **tanpa satu pun pemeriksaan
rentang**, sehingga baris semacam itu memang dapat masuk dan **tidak ada yang tahu berapa banyak**.

**Persyaratan:** `INV-10` (periode unik di dalam satu versi) · **`INV-55`** (periode pelaporan berada
di dalam periode kontraknya) — bergolongan **APLIKASI**, bukan constraint, dan alasannya ada di
`SPEC-INVARIAN.md` §3.5 · `INV-53`

**Tidak termasuk:** **`CARA_PEMBUKUAN` dan `PERIODE_PELAPORAN_KONTRAK` di kepala versi** — keduanya
kolom `VERSI_KONTRAK` yang lahir bersama irisan `14`. Yang di sini **barisnya**, bukan penandanya.
**`CARA_PEMBUKUAN_XOL`** — `P-20`, dan `G2` menyatakannya **tegas tidak masuk** penyerahan pertama:
ia menunggu **Uji X-2**.

**Jalur gagal:** Dua baris berperiode sama pada satu versi -> **ditolak** `INV-10` · Periode yang
mulai sebelum `TANGGAL_MULAI` kontrak atau berakhir sesudah `TANGGAL_BERAKHIR` -> **ditolak**
`INV-55`, pesannya menyebut periode kontraknya · Batas pelunasan lebih awal daripada batas penyerahan
-> ditolak.

**Uji:** **Negatif:** periode kembar; periode yang mulai sehari sebelum kontrak; periode yang berakhir
sehari sesudah kontrak; urutan batas hari yang terbalik.
**Positif — dan ia menguji batas inklusif:** periode yang **mulai tepat pada** `TANGGAL_MULAI` dan
**berakhir tepat pada** `TANGGAL_BERAKHIR` -> **diterima**. `ADR-0022` menetapkan kedua batas
**inklusif**; pemeriksaan yang memakai perbandingan tegas menolak periode yang justru paling lazim.

**Menggantikan:** `P-10` melestarikan daftar `ReportingPeriodList`. Yang bergeser: rentangnya
**diperiksa**, dan di sistem lama tidak.

**Blocked by:** 14

**Dasar:**
```
EVIDENCED(TreatyIn.ReportingPeriodList@ekspor-2026-09; ReportingSubmission/Confirmation/Settlement skalar di kepala)
        DECIDED(INV-10, INV-55, ADR-0022, KTV-A)
        DIASUMSIKAN-CLEAR(KTV-A)
```

- [ ] `PERIODE_PELAPORAN` berdiri sesuai `2-to-spec/KAMUS-KOLOM.md`
- [ ] `INV-10` terpasang; `INV-55` terpasang **di lapisan aplikasi**, dan letaknya itu tertulis sebagai keputusan
- [ ] pesan penolakan `INV-55` menyebut **periode kontraknya**, bukan hanya "di luar rentang"
- [ ] uji positif lulus: periode yang berimpit tepat pada kedua batas **diterima** — batas inklusif
- [ ] `CARA_PEMBUKUAN_XOL` dinyatakan di luar irisan ini, dengan `P-20` dan Uji X-2 sebagai penagihnya
- [ ] `KTV-A` tercatat di `ASUMSI-CLEAR.md`
