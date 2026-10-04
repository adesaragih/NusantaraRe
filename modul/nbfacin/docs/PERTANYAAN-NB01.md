# Lima pertanyaan yang menahan tiket NB-01 — premi PA

**Untuk:** work owner (meneruskan ke DBA dan pemilik aplikasi Pega) · **Dari:** pengerjaan tiket
`issues/01-tracer-money-ratio-premi-pa.md` · **Tanggal:** 1 Oktober 2026

Kode perhitungan premi PA sudah ada dan teruji dengan contoh hitung. Yang belum bisa dinyatakan:
**angkanya sama dengan sistem lama sampai digit terakhir.** Empat jawaban di bawah yang menentukannya.

> ✅ **Status 1 Oktober 2026 — butir 1–4 dijawab work owner, satu pertanyaan baru (butir 5).**
> Keputusannya: `KEPUTUSAN-30-09-2026.md` butir 13–17.
>
> | Butir | Jawaban | Yang terjadi |
> | ---: | --- | --- |
> | 1 | Kasus PA ada di `D:\migrasi\RNM\DDL\CONTOH\` | 4 kasus NB PA; rekonsiliasi **4/4 nol selisih** |
> | 2 | "20 desimal saja" | Data menyingkirkan pemotongan; setengah-ke-atas vs setengah-ke-genap tak terbedakan → kasus seri ditolak |
> | 3 | "Hanya saat `CalculateMethod_FacIn==1`" | ⚠️ **Data membantahnya**: kasus metode 1 cocok dengan langkah 4 L713, bukan langkah 7. Work owner memilih ikuti data |
> | 4 | "Rekomendasi agar tidak bermasalah" | Diskon kosong dikurangi nol — keempat kasus tanpa tag `Discount` membuktikannya |

---

## 1. Satu kasus PA nyata — untuk DBA / work owner

**Yang kami minta:** satu kasus New Business lini **PA** dari sistem lama, berisi masukan coverage dan
premi akhirnya.

**Mengapa:** lima berkas kasus yang ada (`EDM-13445`, `NB-176005`, `NB-181231`, `NB-184233`,
`RNW-10579`) berlini FIRE, AS. KREDIT, dan MARINE CARGO — **tidak ada PA**. Tanpa kasus PA, tes hanya
membandingkan kode dengan hitungan manual atas rumus yang sama, bukan dengan sistem lama.

**Medan yang dibutuhkan per coverage:** `.TSI`, `.Rate`, `.Discount`, `.CalculateMethod_FacIn`,
`.Premium` (hasil), dan `pyWorkPage.OfferFacIn.ProRatePercent` — **dalam teks persis seperti
tersimpan** (termasuk pemisah desimalnya). K-027 hanya membuktikan koma untuk `RATE`; untuk TSI,
ProRatePercent, dan Discount kode kini menduga titik, dan kasus nyata yang memastikannya.

⛔ Kasus akan melewati de-identifikasi tiket 15 sebelum masuk repositori. Mohon **jangan** kirim nama
tertanggung, alamat, atau nama operator bila bisa dihindari.

## 2. Mode pembulatan `@Math.divide` — untuk pemilik aplikasi Pega

**Yang kami tanyakan:** saat `@Math.divide(x, y, 20)` harus membuang digit di belakang desimal ke-20,
Pega membulatkan dengan cara apa — **setengah ke atas**, **setengah ke genap**, atau **dipotong**?

**Mengapa:** korpus hanya menyebut presisinya (4 atau 20 desimal), tidak cara membulatkannya. Sampai
terjawab, kode **menolak** kasus yang perlu dibulatkan (galat `ErrModePembulatanBelumTerverifikasi`)
dan hanya menghitung kasus yang hasilnya eksak. Keputusan work owner 30-09-2026.

**Cara menjawab:** satu contoh angka dari sistem lama yang hasilnya berdigit lebih dari 20 desimal
sudah cukup untuk membedakan ketiganya.

## 3. Langkah 7 `CalculatePremiPA_FacIn` — untuk work owner, lewat UI Pega

**Yang kami minta:** buka activity `CalculatePremiPA_FacIn` (kelas `ASM-FW-GISFW-Data-Coverage`) di UI
Pega dan lihat **langkah 7** (Property-Set `.Premium`, rumus
`@Math.divide((.TSI*pyWorkPage.OfferFacIn.ProRatePercent*.Rate),100000,20)- .Discount`). Apakah
langkah itu **aktif tanpa syarat**, atau hanya saat **`.CalculateMethod_FacIn==1`**?

Yang kedua bukan dugaan kosong: langkah 7 **masih menyimpan syarat `.CalculateMethod_FacIn==1`**
(`CalculatePremiPA_FacIn.xml` L1226), hanya kotak prakondisinya yang tidak dicentang. Syarat yang sama
dipakai langkah 4 (prorata, L791 `=='1'`).

**Mengapa:** langkah 7 ber-`pyStepsPreCondition=false`. Menurut keputusan P-11 (diperiksa di UI Pega
pada langkah rule lain) tanda itu berarti **tetap jalan**. Bila benar, langkah 7 selalu menimpa premi
yang dihitung langkah 4, 5, dan 6 — sehingga ketiga bentuk rumus lain (prorata, short period, fixrate)
**tidak pernah** menjadi premi akhir. Kode saat ini mengikuti tafsir itu (keputusan work owner
01-10-2026). Bila ternyata langkah 7 nonaktif atau hanya untuk metode 1, hasil kode berubah, dan
tiket 04 ikut berubah.

## 4. `.Discount` kosong — untuk pemilik aplikasi Pega

**Yang kami tanyakan:** bila `.Discount` coverage **kosong**, apakah ekspresi `… - .Discount` di Pega
menghitungnya sebagai **nol**, atau gagal?

**Mengapa:** aturan proyek (ADR-U-0022) menyatakan teks kosong pada kolom angka adalah **kosong, bukan
nol**, jadi kode kini **menolak** menghitung premi dengan diskon kosong (`uang.ErrUangKosong`). Bila
Pega menganggapnya nol, setiap coverage tanpa diskon akan ditolak di sistem baru padahal berjalan di
sistem lama. (Untuk `@toDecimal("")` sudah ada bukti nol di modul Claim Life — tetapi itu fungsi
lain, bukan pengurangan langsung.)

---

Jawaban dicatat di `KEPUTUSAN-30-09-2026.md` (atau register K) dan diturunkan ke kode serta tes.

## 5. Arti `pyStepsPreCondition=false` — untuk work owner, lewat UI Pega *(baru, 1 Oktober 2026)*

> ✅ **Dijawab 1 Oktober 2026:** *"itu hanya berlaku di PA saja"* — P-11 tetap berlaku umum.
> `KEPUTUSAN-30-09-2026.md` butir 18, termasuk cakupan yang belum diputuskan (21 langkah di 8 rule PA).

**Yang kami tanyakan:** keputusan P-11 menyatakan langkah ber-`pyStepsPreCondition=false` **tetap
jalan**. Di `CalculatePremiPA_FacIn` langkah 7 bertanda sama, tetapi data menunjukkan langkah itu
**tidak berpengaruh**: kasus PA metode 1 (`NB-90996`) cocok eksak dengan langkah 4 dan meleset 0,0048
dari langkah 7. Mana yang benar untuk langkah lain bertanda sama?

**Mengapa penting di luar tiket ini:** `[terverifikasi]` korpus NB memuat **577 langkah** ber-tanda
itu, 404 di antaranya masih menyimpan ekspresi kondisi (`01-activity/01-inventaris-activity-nb.md`
T-5). Pemetaan alur EDM juga dinyatakan valid berdasar P-11. Bila P-11 hanya berlaku untuk sebagian,
tafsir di modul lain ikut perlu ditinjau.

**Yang sudah bisa dikatakan:** satu langkah (P-11, `SaveEDMToJsonPolicy_Act`) terlihat aktif di UI;
satu langkah (langkah 7 ini) terbukti tidak berpengaruh oleh data. `[dugaan]` Penyebabnya bisa
pengaturan langkah lain yang tidak terbaca dari ekspor, atau versi rule produksi yang berbeda dari
korpus (OQ-011).

**Cara menjawab:** buka langkah 7 `CalculatePremiPA_FacIn` di UI Pega produksi, dan bandingkan
tampilannya dengan langkah "Get ProdKe" `SaveEDMToJsonPolicy_Act` yang dipakai P-11.
