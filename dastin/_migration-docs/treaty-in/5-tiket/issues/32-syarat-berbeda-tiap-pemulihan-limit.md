---
status: aktif
golongan: baru
---

# 32: Syarat berbeda ditetapkan untuk tiap pemulihan limit pada sebuah layer

*Asal: `DAFTAR-PEKERJAAN.md` §2.2b `P-59` · `SPEC-MODEL-DATA.md` §10.3b · `TETAPAN-DI-KODE.md` §3.*

**What to build:** **PK** menetapkan **syarat yang berbeda untuk tiap pemulihan limit** pada sebuah
layer — pemulihan pertama, kedua, dan seterusnya boleh berpersentase berbeda, baik porsi limit yang
dipulihkan maupun tarif premi pemulihannya.

**Kenapa begini:** Sistem lama **tidak dapat menyatakannya.** `SetReinstatementPct.xml` menulis
`ReinstatementPct = "100"` dan `AdditionalPct = "100"` sebagai **tetapan di dalam kode**, bukan
sebagai masukan — sehingga setiap pemulihan selalu seragam, dan ketentuan pasar yang biasa
*("pemulihan pertama gratis, kedua berbayar 100%")* **tidak punya tempat**. Yang **PELESTARIAN**
hanyalah menyimpan pemulihan sebagai daftar; yang **BARU** adalah **nilainya boleh berbeda**.

> Ini **instans keempat** pola *"cacat bersembunyi di balik nilai bawaan"*: selama semuanya 100, satu
> baris dan sepuluh baris menghasilkan angka yang sama — dan tidak ada yang pernah melihat batasnya.

**Persyaratan:** `INV-49` (pemulihan adalah **besaran berulang**, dikecualikan dari partisi `INV-47`) ·
`INV-41` (persentase di rentang 0–100) · `SPEC-MODEL-DATA.md` §10.3b

**Tidak termasuk:** **Entitas `PEMULIHAN_LIMIT` itu sendiri** — irisan `31` yang membuatnya.
**Premi pemulihan yang dihitung darinya** — turunan, `ADR-0037`; dan `TANPA_HITUNG_PREMI_PEMULIHAN`
sudah menjadi kolom `LAYER` di irisan `31`.
**Baris selisih untuk perubahan persentase pemulihan** — tiket `13`, jalur addendum.

**Jalur gagal:** Persentase pemulihan di luar 0–100 -> **ditolak** `INV-41` · Dua baris pemulihan
bernomor urut sama pada satu layer -> ditolak · `INV-47` menuntut baris pemulihan berjumlah utuh ->
**tidak terjadi**; `INV-49` mengecualikannya, dan bila ia tetap menolak maka pengecualiannya belum
terpasang.

**Uji:** **Negatif:** persentase 101; persentase negatif; nomor urut pemulihan kembar.
**Positif — dan ia pokok irisan ini:** satu layer dengan **pemulihan pertama berporsi 100% bertarif
0%** dan **pemulihan kedua berporsi 100% bertarif 100%** -> **diterima**, dan keduanya terbaca kembali
berbeda. Data uji **wajib memuat nilai selain 100**; uji yang seluruhnya bernilai 100 **tidak
memisahkan apa pun** — ia menghasilkan hasil yang sama pada rancangan lama maupun baru.
**Positif kedua:** `INV-47` **tidak** menolak layer berpemulihan dua baris.

**Menggantikan:** tetapan `ReinstatementPct = "100"` dan `AdditionalPct = "100"` di
`SetReinstatementPct.xml`. **Golongannya BARU**, dan karena itu `CARA MENYALAKANNYA` wajib.

**Satu pertanyaan wawancara menyertainya, dan ia tidak menahan pembangunannya:** apakah keseragaman
100/100 di sistem lama **kehendak bisnis** atau **akibat cara sistem dibangun**? **Datanya tidak
dapat menjawab** — seluruhnya disemai 100, sehingga *"selalu seragam"* dan *"tidak pernah bisa
berbeda"* menghasilkan data yang persis sama.

**CARA MENYALAKANNYA:**

| # | Isi |
|---|---|
| 1 | `SPEC-INVARIAN.md` §3 — **CO-6, CO-7, CO-8** |
| 2 | **siapa membaca, seberapa sering** — **PK yang menyusun layer non-proporsional**, pada setiap kontrak berpemulihan. **Bukan pemantau berkala**: kegagalannya muncul saat seseorang mencoba mengisi nilai berbeda |
| 3 | **ambang berangka** — nyala penuh ketika **satu kontrak** tercatat dengan dua pemulihan berpersentase berbeda **dan lolos persetujuan**. Sebelum itu ia tersedia tetapi bernilai bawaan seragam, sehingga **tidak mengubah apa pun** |
| 4 | **siapa boleh menyalakan** — pemilik proses, setelah bisnis mengonfirmasi bahwa pemulihan bertingkat memang dipakai NuRe |

**Blocked by:** 31

**Dasar:**
```
EVIDENCED(SetReinstatementPct@ekspor-2026-09 - ReinstatementPct="100" dan AdditionalPct="100" sebagai tetapan)
        DECIDED(INV-49, SPEC-MODEL-DATA §10.3b)
```

- [ ] porsi dan tarif tiap baris pemulihan dapat **berbeda satu sama lain**, dan terbaca kembali berbeda
- [ ] `INV-41` terpasang dan diuji pada **kedua** ujung rentang
- [ ] uji positif dijalankan dengan **nilai selain 100**; data uji seluruhnya-100 **tidak diterima sebagai bukti**
- [ ] `INV-47` tidak menolak layer berpemulihan dua baris — `INV-49` terpasang sebagai pengecualian bernama
- [ ] keempat butir **CARA MENYALAKANNYA** terisi, termasuk ambang berangka
- [ ] pertanyaan wawancara *"kehendak bisnis atau akibat cara sistem dibangun"* **terdaftar**, beserta catatan bahwa datanya tidak dapat menjawabnya
