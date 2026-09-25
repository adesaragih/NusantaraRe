---
status: accepted
label: DECIDED
---

# Kurs yang dipakai ikut tersimpan bersama nilai yang dikonversinya

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, empat lapisan; `pengetahuan/ddl/` 49 berkas. Sapuan S4 dan S10, dijalankan 18 September 2026 — `SAPUAN-S3-S9.md` dan `SAPUAN-S10-S16.md`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut.

> **Naik ke `accepted` 19 September 2026, dan kedua cabangnya dicabut.** Versi pertama berstatus `proposed` **justru karena** dua cabang terbuka (A dan B) yang menunggu jawaban akuntansi atas `ASK-AKUNTANSI` butir 7. Cabang itu tidak dijawab — ia **dihapuskan** oleh keputusan seragam di bawah, sehingga alasan satu-satunya ADR ini berstatus `proposed` ikut hilang. Butir 7 tetap hidup sebagai **verifikasi**, dan tidak menahan ADR ini maupun tiket mana pun.

**Setiap nilai uang, tanpa pengecualian, disimpan lengkap dengan lima hal di baris yang sama**: nilai asli, kode mata uangnya, nilai rupiahnya, kurs yang dipakai menghasilkan nilai rupiah itu, dan asal-usul kurs tersebut.

Nilai rupiah tanpa kurs tidak dapat lahir; kurs tanpa nilai rupiah boleh ada; kolom rupiah yang kosong boleh ada.

## Mengapa bukan sekadar kehati-hatian

Tiga bukti yang berdiri sendiri-sendiri.

**Pertama, sistem lama sudah melakukannya di satu tempat.** `POOLDATA.CLAIMXOL2` memuat `KursIDR` **berdampingan** dengan nilai yang dikonversinya, di satu baris. Begitu pula `SpreadingRisk` dan `CNPCurrencyList` di sisi Pega, yang membawa `KursIDR` bersama `ClaimAmountIDR`.

**Kedua, dan ini yang menentukan: kolom rupiah dapat terisi tanpa kurs pernah dipakai.** Satu kolom yang sama diisi oleh dua rule dengan perlakuan berbeda:

| Rule | Perlakuan |
|---|---|
| `CountLossAllocation_act` baris 9585 | menguji mata uang lebih dulu, lalu mengalikan kurs bila perlu |
| `SetActualPremium_ACT` baris 681 | **menyalin nilai apa adanya**, tanpa menguji mata uang, tanpa mengalikan kurs |

Bila nilai sumbernya bukan rupiah, kolom berlabel rupiah memuat angka yang bukan rupiah — dan **tidak ada apa pun di baris itu yang memperlihatkannya**. Menyimpan kursnya membuat keadaan itu terbaca alih-alih tersembunyi. (D23)

**Ketiga, rantai nilai melewati beberapa berkas tanpa menyatakan skala.** Rantai `.GrossValue` ditelusuri dari lahir sampai keluar dan tidak menyatakan skala di satu titik pun, sementara rantai `.AdjusterFeeValue` menyatakan skala di setiap pembagian. Nilai yang tidak membawa jejak perlakuannya tidak dapat direkonsiliasi belakangan. (D24)

## Considered Options

**Menyimpan nilai rupiah saja** ditolak: itu keadaan sistem lama, dan D23 memperlihatkan keadaan itu sudah menghasilkan kolom yang tidak dapat dipercaya tanpa membaca kode yang mengisinya.

**Menyimpan kurs di tabel kurs terpisah dan menghubungkannya lewat tanggal** ditolak sebagai penggantinya, meski tabel kurs bertanggal tetap dibutuhkan untuk keperluan lain. Alasannya: kurs yang **dipakai** pada satu perhitungan belum tentu kurs yang **berlaku** pada tanggal itu. Yang perlu direkonsiliasi adalah yang dipakai, bukan yang seharusnya.

## Keberlakuan seragam — keputusan 19 September 2026, menggantikan cabang A dan B

Sapuan S4 menemukan **tujuh baris, sembilan nama** properti bernilai uang yang **tidak punya padanan rupiah sama sekali**: `AdjusterFee`, `AdjusterFeeValue`, `Salvage`, `SalvageValue`, `CNPOthersFee`, `OthersFee`, `AdjustmentValue`, `TotalClaim`, `PremiumSpreaded`. Penamaan alternatif (`Rupiah`, `Idr`, `_IDR`, `Rp`) ikut disapu dan mengembalikan nol; yang ada hanya sufiks `RNM`, yaitu porsi pihak, bukan mata uang. Lima dari sembilan masuk hitungan instruksi bayar, dan rantai `.AdjusterFeeValue` terbukti tidak melewati satu pun titik konversi.

Versi pertama membaca temuan itu sebagai dua cabang — **A**: ketiadaan itu kesengajaan, sehingga ADR ini tidak berlaku untuk kelas itu; **B**: ketiadaan itu kekurangan, sehingga kelas itu memerlukan padanan rupiah. **Keduanya dicabut.** ADR ini berlaku untuk **seluruh** nilai uang, termasuk kesembilan nama di atas.

**Alasannya perbandingan biaya, bukan pengetahuan baru.** Kedua cabang tetap tidak terjawab dari XML; yang berubah adalah kesadaran bahwa keduanya **tidak berbiaya setara**:

| Bila ternyata… | dan kita menyediakan kolomnya | dan kita tidak menyediakannya |
|---|---|---|
| rupiahnya **tidak pernah** dibutuhkan | kolom kosong — tidak ada yang rugi | benar, kebetulan |
| rupiahnya **dibutuhkan** | benar | **migrasi kedua** |

Satu sisi berbiaya kolom kosong; sisi lain berbiaya migrasi kedua. Menunggu jawaban akuntansi untuk memilih di antara keduanya berarti menahan seluruh model demi menghindari biaya yang lebih kecil dari biaya menunggunya.

Ini bentuk yang sama dengan ADR-0025, yang menyediakan kolom lingkup sebelum ada yang memerlukannya.

**`ASK-AKUNTANSI` butir 7 tetap berdiri**, dan kedudukannya berubah: ia **memverifikasi** apakah kesembilan kolom itu terisi atau tetap kosong. Ia tidak lagi menahan ADR ini, tidak menahan `SPEC-MODEL-DATA.md`, dan tidak menahan tiket mana pun.

## Consequences

**Tidak ada lagi rancangan yang ditunda oleh ADR ini.** Pembatasan versi pertama — *"tidak ada tabel, kolom, aturan, atau uji yang boleh dirancang dengan mengandaikan cabang A maupun cabang B"* — dicabut bersama cabangnya.

Yang berlaku sekarang seragam: **setiap** nama bernilai uang mendapat kelima medan, baik yang hari ini sudah punya padanan rupiah — `ClaimAmountIDR`, `GrossValueIDR`, `ValueIDR`, `AmountIDR`, `TotalSumInsuredIDR` — maupun kesembilan yang hari ini tidak punya.

**Dua akibat yang harus dilihat terpisah.** Yang pertama: kolom rupiah untuk kesembilan nama itu akan **kosong pada data migrasi**, karena nilainya memang tidak pernah ada di sistem lama — kosong, bukan nol (ADR-0019). Yang kedua: struktur menerima nilai itu bila kelak ada yang menghitungnya, tanpa perubahan skema.

ADR ini juga menutup pertanyaan satuan pada `.ValueAdjustment` yang datang dari luar: nilai yang menyeberang batas **wajib membawa mata uang dan kursnya**, dan nilai tanpa mata uang **ditolak di batas**, bukan diterima lalu ditebak. Itu yang membuat H1/H2 berhenti menjadi masalah rancangan — lihat `BLUEPRINT.md` §8.7.

ADR ini juga tidak menutup D23. Kolom rupiah yang terisi tanpa konversi adalah cacat sistem lama; apakah nilai lamanya diperbaiki, dibiarkan, atau didaftar mengikuti **AK-2a** dan **AK-2b**, dan itu putusan tersendiri.
