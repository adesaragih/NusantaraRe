---
status: accepted
label: DECIDED-TEKNIS
---

# Presisi tinggi sepanjang rantai perhitungan, pembulatan hanya di tepi

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Sistem lama tidak punya aturan pembulatan: dari 673 ekspresi aritmetika, **518 (77%) tidak menyatakan skala sama sekali** dan bergantung pada perilaku desimal bawaan Pega, sementara 155 sisanya memakai tujuh skala berbeda. Kami menetapkan aturan baru: nilai moneter dihitung dan disimpan pada satu presisi tinggi seragam sepanjang rantai perhitungan, dan pembulatan hanya dilakukan di dua tepi — saat ditampilkan ke pengguna, dan saat dikirim ke sistem hilir.

## Considered Options

Menyalin skala per ekspresi apa adanya ditolak: skalanya tidak konsisten bahkan di dalam satu rule (`CountClaimTNP_Act` memakai 10, 20, dan 5; `CountLossAllocation_act` memakai 10 dan 20 dalam porsi berimbang), sehingga menyalinnya berarti mengabadikan ketidaksengajaan.

## Consequences

Skala ternyata **bukan semata kebiasaan penulisan**, melainkan mengikuti batas modul: inti klaim non-proporsional memakai 20, sisi treaty/premi memakai 4, pemformatan tampilan memakai 0 dan 2 (selalu lewat "bagi dengan 1"), dan tarif pajak di `SetPPNPPH` memakai 8. Dugaan awal bahwa ini murni kebiasaan developer meleset dan sudah dikoreksi.

Angka presisi final **belum ditetapkan** dan ADR ini tidak boleh menyebut angka apa pun sebagai final sebelum hasil profil data Oracle diterima. Risiko yang harus diputuskan lebih dulu: bila profil menunjukkan kolom uang di produksi benar-benar menyimpan lebih dari dua desimal, maka aturan "bulatkan dua desimal saat kirim ke hilir" justru akan mengubah angka yang selama ini diterima akuntansi.

## Catatan 18 September 2026 — tetap `proposed`, dan alasannya berubah

Status **tidak** dinaikkan. Dua hal yang perlu tercatat:

1. **`NUMBER(20,4)` pada `TREATYINPRODUCTION` diperlakukan sebagai konvensi satu modul**, bukan standar perusahaan, sampai akuntansi menyatakan sebaliknya. Ia **tidak** disebarkan ke tabel lain dan tidak dipakai sebagai dasar keputusan presisi.
2. **Alasan ADR ini masih draft sudah berubah.** Semula: datanya belum ditarik. Sekarang: **basis data memang tidak menyimpan presisi untuk nilai klaim sama sekali** — nilainya tidak punya kolom, tersimpan sebagai teks di dalam JSON, dan bahkan view `CLAIMXOL` mengeluarkannya sebagai `varchar2`. Profil data tidak akan mengubah kenyataan itu.

Yang menghambat sekarang adalah keputusan kebijakan, bukan ketersediaan data — `ASK-AKUNTANSI.md` pertanyaan 1.

## Naik ke `accepted` — 18 September 2026

Status naik lewat **AK-1** di `REGISTER-RATIFIKASI.md`, dan **bukan** karena data profil akhirnya datang — catatan di atas sudah menjelaskan bahwa data itu tidak akan pernah menjawabnya.

Yang berubah: angka presisinya kini diturunkan dari **skala yang dipakai sistem lama itu sendiri**, bukan dari profil kolom yang tidak ada. Skala 20 pada 126 ekspresi di seluruh rule inti klaim non-proporsional (`BLUEPRINT.md` §6.3).

| Kelompok | Tipe |
|---|---|
| U1 nilai uang antara · P1 persentase dan porsi · K1 kurs · L1 limit dan premi deposit | `NUMBER(38,20)` |
| P2 tarif pajak dan brokerage | `NUMBER(11,8)` |
| C1 cacah dan nomor urut | `NUMBER(9)` |

**Tidak ada kelompok untuk "nilai uang final".** Pembulatan dua desimal terjadi hanya di tepi — lapisan view kompatibilitas dan arsip muatan keluar. Alasannya `TempKasir.CARI9 = .TotalClaim - .PremiumSpreaded`: `.TotalClaim` adalah nilai antara sekaligus bahan instruksi bayar, dan satu kolom tidak dapat berskala 20 dan 2 sekaligus.

Ratifikasi akuntansi dicatat di `REGISTER-RATIFIKASI.md` dan **tidak menahan pekerjaan**.

## Bukti yang menguatkan — sapuan S10, 18 September 2026

Ditambahkan sesudah ADR ini `accepted`. Ia **menguatkan**, tidak menutup apa pun dan tidak mengubah putusan.

Dua rantai nilai uang ditelusuri dari lahir sampai keluar, bukan dicacah per ekspresi:

| Rantai | Skala yang dinyatakan | Berkas yang dilewati |
|---|---|---|
| `.GrossValue` | **tidak ada, di satu titik pun** | `CountLossAllocation_act`, `CreateChildKomiteCNP_Act`, `SaveCNPLayerList_Act`, `SaveDataToOSAksep_Act`, `SaveToOS`, `GetSelisihActual_Act` |
| `.AdjusterFeeValue` | **20**, di setiap pembagian | rule yang sama, sebagian besar |
| `Local.URLimit` | **10** | `CountLossAllocation_act` 6038, 6065 |

**Tiga skala dalam satu modul, dua di antaranya di berkas yang sama.** Angka 77% di badan ADR ini adalah statistik; ini contoh berjalan: sebuah nilai uang yang benar-benar dikirim ke sistem hilir menempuh enam berkas tanpa satu pun pernyataan skala, sementara nilai di sebelahnya menyatakannya setiap kali.

Menguatkan pernyataan *"skala mengikuti batas modul"* di bagian Consequences dengan satu koreksi halus: pada rantai ini skala mengikuti **besaran yang dihitung**, bukan rule yang menghitungnya — `CountLossAllocation_act` memakai 10 dan 20 untuk besaran berbeda di dalam dirinya sendiri.

Tercatat sebagai **E1**, **E2**, **E12**, **B1**, **B2**, dan **D24**. Rinci di `SAPUAN-S10-S16.md` §3.

## Penutupan B1, B2, E1, E2, E12 — 19 September 2026

**Catatan status**: ADR ini **sudah** `accepted` sejak 18 September 2026 lewat AK-1; ia tidak dinaikkan lagi hari ini. Yang ditambahkan bagian ini adalah **penutupan lima butir register** yang selama ini masih menunggunya, beserta alasan mengapa menunggu lebih lama tidak akan mengubah apa pun.

### Dasarnya sudah lengkap dan tidak akan berubah oleh data

Empat lapisan diperiksa, dan keempatnya sepakat:

| Lapisan | Keadaan |
|---|---|
| Tipe kolom | **88 dari 120** kolom `NUMBER` di DDL tanpa presisi maupun skala |
| Muatan JSON dan blob | nilai uang tersimpan sebagai **teks**; view `CLAIMXOL` mengeluarkannya sebagai `varchar2` |
| Ekspresi Pega | **77%** tanpa skala; rantai `.GrossValue` menempuh enam berkas tanpa satu pun pernyataan |
| Java tertanam | **nol `double`, nol `float`, nol `BigDecimal`**; uang melintas sebagai `String` 92 kali lewat `.toString()` (**D39**) |

Lapisan keempat adalah yang menutup pertanyaannya. Selama hanya tiga lapisan terbaca, masih mungkin berharap kendali skala hidup di tempat yang belum dilihat. Java adalah tempat terakhir itu, dan ia **tidak menghitung uang sama sekali** — ia memindahkannya sebagai teks.

> **Tidak ada presisi yang dapat dipulihkan, karena tidak pernah ada presisi yang ditegakkan.**

Profil data Oracle tidak dapat mengubah kalimat itu. Karena itu **B1 dan B2 tidak lagi menunggu REQ-001**.

### Maka presisi adalah keputusan maju, bukan temuan

1. **Skala kanonik 20 sepanjang rantai hitung.** Angka ini bukan selera: ia **skala tertinggi yang benar-benar dipakai sistem lama** — inti klaim non-proporsional, dan setiap pembagian di rantai `.AdjusterFeeValue`. Memakainya karena itu **tidak dapat kehilangan apa pun yang pernah ada**.
2. **Pembulatan hanya di tepi** — saat ditampilkan, saat dikirim ke hilir, saat dibukukan. **Tidak di tengah rantai.**
3. **Migrasi mengurai teks apa adanya** dan menyimpan seluruh digit yang ada. Nilai yang **tidak terurai tidak dibulatkan dan tidak dibuang** — ia masuk jalur tiket `14` dan tercatat beserta asalnya.
4. **Skala 10 dan 4 yang terbaca di sistem lama dicatat sebagai perilaku lama di tabel perbedaan, bukan ditiru.** `Local.URLimit` memakai 10; `TREATYINPRODUCTION` memakai `NUMBER(20,4)`. Keduanya fakta tentang sistem lama, bukan standar yang diwarisi.

### Apa yang ini tutup, dan apa yang tetap berdiri

| Butir | Isinya | Penutupnya |
|---|---|---|
| **B1** | presisi dan skala final tiap kolom nilai | butir 1 di atas |
| **B2** | presisi tidak dijaga di tiga (kini empat) lapis | tidak ada yang dapat dipulihkan; keputusan maju menggantikannya |
| **E1** | dua rule mencampur skala pembagian | skala mengikuti besaran, bukan rule — tidak ditiru |
| **E2** | 77% ekspresi tanpa skala | sama |
| **E12** | tiga skala dalam satu modul | sama |

**REQ-001 tetap diminta, dan kedudukannya berubah**: ia **memverifikasi** apa yang benar-benar tersimpan di blob, dan **tidak lagi menahan**. Bila jawabannya kelak memperlihatkan digit yang lebih banyak daripada dugaan, skala 20 tetap menampungnya; bila lebih sedikit, tidak ada yang hilang.
