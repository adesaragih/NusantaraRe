# 09: Angka uang — total penyesuaian per mata uang

**Status:** ready-for-agent

**Blocked by:** **03 (layar komite)**

## Hasil & nilai pengguna

Sebagai **penyetuju komite**, saya melihat **total nilai penyesuaian per mata uang** beserta
**totalnya dalam rupiah**, sehingga saya bisa membandingkan antar mata uang sebelum memutuskan.

*(User story 8, 32, 33 di spec)*

## Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/Activity/SetKomiteList_Act.xml` | `[terverifikasi]` **satu-satunya penghitung uang berkelas komite**. **Langkah 1** menghentikan activity bila jenis pengajuan bukan penyesuaian. **Langkah 6** membuang baris bermata-uang-sama dari daftar sementara |
| idem | `[terverifikasi]` **langkah 7.1.1**, di dalam perulangan, menjumlahkan **per mata uang**: total kotor, total bersih, dan keduanya **dikali kurs** untuk rupiah. **Langkah 7.2** menaruh kedua total pertama ke properti kasus |
| idem | `[terverifikasi]` gerbang penjumlahan: baris **berstatus tolak tidak ikut dijumlahkan**, dan penjumlahan hanya untuk mata uang yang sedang diproses |

⭐ `[terverifikasi]` **Nol perhitungan uang di potongan Java** — seluruh rumus terbaca dari penugasan
properti biasa.

## Kurs

`[keputusan work owner]` 2026-09-18 — **"ikuti aja query di xml CurrencyStandard."**

`[terverifikasi]` Kurs diambil lewat **fungsi tersimpan basis data** dengan **dua masukan saja**:
kode mata uang dan **tanggal server saat query dijalankan**. Tidak ada nomor kasus, tidak ada lini,
tidak ada tanggal kasus.

⭐ **Kurs yang terkunci adalah kurs TANGGAL PENCARIAN DIJALANKAN** — bukan tanggal kasus dibuat,
bukan tanggal komite menyetujui. ⚠️ Isi fungsi tersimpannya **tidak ada di korpus** dan
`[keputusan work owner]` **tidak diminta ke DBA**.

`[terverifikasi]` Penguncian nilainya terjadi **di luar modul ini** — pola *ambil-hanya-bila-kosong*
ada di modul Claim Prop. **Di dalam modul ini nol pola penguncian untuk uang.**

## ⚠️ TITIK YANG SENGAJA DIUBAH — ketelitian angka

`[keputusan work owner]` 2026-09-18, diralat 2026-09-19. **Ini perubahan sadar, bukan peniruan.**

`[terverifikasi]` Di Pega **tidak ada pembulatan sama sekali**, dan persen dibagi seratus dengan
**tiga cara berbeda** — dua di antaranya di activity yang sama untuk menghitung hal yang sama.

| Di sistem baru | Ketetapan |
| --- | --- |
| Hitungan | sampai **20 angka di belakang koma**, **nol pembulatan di tengah jalan** |
| Penyimpanan | **mengikuti bentuk kolom yang sudah ada — 20 digit, 8 di belakang koma** |
| Tampilan | **4 angka di belakang koma** |
| Tipe | **desimal, bukan bilangan pecahan biner** |

⚠️ **Pembulatan terjadi di batas penyimpanan, diterima sadar.** Angka persen yang di Pega dihitung
sampai 10 desimal **akan menjadi 8** saat disimpan. **Hasil sistem baru AKAN BERBEDA dari Pega di
angka belakang koma, dan itu disengaja.**

⚠️ `[keputusan work owner]` 2026-09-19 — bila nilai melewati batas digit, basis data **menolak
menyimpan**, bukan membulatkan. **Kegagalan itu harus terlihat.**

## Yang harus diuji

**Diverifikasi oleh:** spec.md AC 30 · 36 · 37 · 38 · 39 · 40 · 41 · 42 · 43 · 44 · 71

- [ ] Total per mata uang benar, dan baris **berstatus tolak tidak ikut dijumlahkan**.
- [ ] Total rupiah = total mata uang **dikali kurs** yang tersimpan pada baris itu.
- [ ] Baris bermata-uang-sama **digabung** sebelum dijumlahkan.
- [ ] Pada jalur **bukan penyesuaian**, penghitungan **tidak berjalan sama sekali**.
- [ ] Nilai uang melewati aplikasi **tanpa melewati bilangan pecahan biner**.
- [ ] Hitungan berjalan sampai **20 angka di belakang koma**; **tidak ada pembulatan di tengah**.
- [ ] Nilai yang melewati batas digit **ditolak dengan galat yang terlihat**.

## Butir `[terbuka]` yang menyentuh tiket ini

- **Ketelitian pembagian persen di Pega tidak seragam** — mana yang sah belum ditetapkan. Ini
  menentukan bagaimana **selisih terhadap data lama** diterangkan, bukan apa yang dikerjakan
  sistem baru.
- **Seberapa besar selisih** akibat persen 10 desimal menjadi 8 saat disimpan — **belum diukur**,
  dan **menumpuk lewat perulangan** sebelum sampai ke tiket 12.
- **Berapa kali perulangan berputar** — menyentuh berapa kali angka berubah.

## Seam & verifikasi

Memakai ulang seam Claim Prop.
