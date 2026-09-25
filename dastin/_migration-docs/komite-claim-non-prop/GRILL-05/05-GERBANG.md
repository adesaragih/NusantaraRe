> Modul  : Komite Claim Non Prop · Ronde 05 · 2026-09-20
> Peran  : interogator
> Masukan: `01-TEMUAN.md` · `02-SIDANG.md` · `03-LUBANG.md` · `REGISTER-PAGAR.md` · `KETETAPAN.md`
> Status : DITUTUP 2026-09-20
> Sifat  : TAMBAH-SAJA

# 05 · GERBANG KESIAPAN PER ALIRAN KERJA

Kesiapan dinilai **per aliran kerja**, bukan per seksi dokumen. Satu seksi
`PENGETAHUAN.md` dapat memuat dua aliran dengan tingkat risiko berbeda; membekukan seksi
berarti membekukan yang tidak perlu dan meloloskan yang perlu.

Tiga keadaan: **BOLEH MULAI** — spesifikasi boleh ditulis penuh · **TERBATAS** — boleh
ditulis dengan bagian yang dinamai ditunda · **BEKU** — tidak ditulis sama sekali.

---

## A-1a · Mekanika keputusan dan perutean jenjang

**Isi.** Gerbang lingkar jenjang, penambahan hitungan, penetapan keadaan terima/tolak,
pemilihan penerima giliran, penutupan case.

**Temuan penopang.** N-01, N-02, N-03, N-11, N-12 · F-13, F-17 (sidang) · G-17, K-02.

**Keadaan: BOLEH MULAI.**

Seluruh mekanika kini berbukti langkah: `IsKomiteLoop` terbaca utuh, penambahan hitungan
berlokasi, router terbaca sampai kode transisinya, dan urutan giliran — baris pertama yang
belum memutuskan, urut `.DEGREE ASC` — berhenti menjadi tafsir. Tidak ada bagian A-1a yang
menunggu berkas dan tidak ada yang menunggu keputusan.

**Syarat yang sudah gugur.** Pagar yang pernah menahan A-1a menunggu "parameter metode dua
activity inti" (G-01/Q-9). Router tidak memerlukan parameter metode: seluruh isinya
`Property-Set`, yang memang terekspor. Pagar itu tidak pernah menyentuh A-1a dan tidak
dipasang kembali.

---

## A-1b · Pembentukan sirkulasi

**Isi.** Kedua jalur pembuatan case anak, pembentukan roster, penetapan jumlah jenjang,
penulisan nomor komite kembali ke klaim, penanda tutup/tolak.

**Temuan penopang.** N-04, N-05, N-06, N-07, N-08, N-09, N-10 · F-13 s/d F-24 (sidang).

**Keadaan: BOLEH MULAI.**

Aliran ini dibuka ronde 4 dan dibaca penuh ronde 5. Tujuh temuan baru seluruhnya berdiri di
atas langkah yang dibaca; tak satu pun menunggu bukti. Dua di antaranya menuntut ketetapan
baru (N-07 memperluas H-5, N-08 menuntut kegagalan yang terlihat), dan ketetapan adalah
pekerjaan ronde ini, bukan penghalangnya.

**Yang harus ikut ditulis, bukan ditunda.** Pembentukan roster yang idempoten (N-06),
transaksi tunggal untuk pembuatan-dan-akibatnya (N-07), kegagalan yang terlihat (N-08).
Menunda ketiganya berarti memindahkan tiga cacat ke sistem baru dengan sengaja.

---

## A-2 · Layar dan aturan medan

**Isi.** 296 kendali, 137 baca-saja, 15 wajib; medan induk penghalang; hak baca.

**Temuan penopang.** G-14, G-20 · E-4, J-1, J-2, J-4 · F-6 s/d F-9.

**Keadaan: BOLEH MULAI.**

Tidak digrilling ronde ini dan tidak berubah karenanya. Satu tambahan kecil dari ronde ini:
`FlagProrate` bernilai empat, bukan dua (N-15), sehingga kendali yang bergantung padanya
punya empat keadaan tampilan. Itu memperjelas, bukan menahan.

---

## A-3 · Penomoran akseptasi

**Isi.** Penerbitan nomor akseptasi, periode buku, perilaku tanggal kosong.

**Temuan penopang.** G-05, G-10, K-05, K-06 · D-4, E-2 · Q-4, Q-11.

**Keadaan: TERBATAS.**

Mekanikanya boleh ditulis penuh — prosedurnya dipegang dan terbaca. Yang **ditunda** adalah
menetapkan periode buku sebagai aturan tetap, karena dasar faktualnya bertanda
`TAFSIR (menunggu C-01)`: 49 berkas DDL adalah hasil reverse-engineer, bukan ekspor
langsung. Syarat pembukaan: **B5-7** pada `03-LUBANG.md` — pembandingan isi badan
`PROC_GENERATE_SEQUENCE_NUMBER` terhadap basis data berjalan, pemegang DBA.

---

## A-4 · Kiriman ke Kasir

**Isi.** Muatan pembayaran, penyaring jenis treaty, perangkaian JSON, enam efek luar dalam
satu commit.

**Temuan penopang.** G-03, G-04, K-10, K-12 · Q-10 · L5-2.

**Keadaan: BEKU.**

Syarat pembukaan: **PG-04** pada `REGISTER-PAGAR.md`, menunjuk `INVENTARIS-BUKTI.md` §2.5
baris 1 — nol baris `POOLDATA.DIRECTTOKASIR_LOG` dipegang. Pembukanya satu `SELECT`
baca-saja atas log kiriman itu.

Ronde ini membaca `TransferType` yang juga muncul di jalur Kasir, dan berhenti di batas:
yang dicatat hanya penulisnya, tidak akibatnya pada muatan.

---

## A-5 · Akseptasi OS dan pembalikan

**Isi.** Penulisan baris akseptasi, rincian layer, baris pembalik, cakupan pembalikan.

**Temuan penopang.** G-06, G-08, G-09, G-11 · Q-5, Q-12, Q-14 · L5-3.

**Keadaan: BEKU.**

Syarat pembukaan: **PG-05**, menunjuk `INVENTARIS-BUKTI.md` §2.1 — tiga prosedur akseptasi
dan `XOL2_AKSEP_KLAIM` tidak dipegang. Pembukanya source ketiga prosedur itu.

---

## A-6 · Isi dokumen dan surat

**Isi.** Susunan surat jenjang berikutnya, PDF akseptasi, PDF tutup/tolak, lampiran.

**Temuan penopang.** G-13, G-15 · L5-1.

**Keadaan: BEKU RINGAN.**

Yang beku adalah **isi**; yang boleh ditulis adalah **kapan** dokumen terbit dan **siapa**
penerimanya, karena keduanya milik A-1a dan A-1b. Syarat pembukaan: **PG-06**, menunjuk
`INVENTARIS-BUKTI.md` §2.3 — `PostEmailKomiteCNP` tidak ada. Pembukanya **B5-4**.

---

## Ringkasan gerbang

| Aliran | Keadaan | Berubah ronde ini? | Syarat pembukaan |
|---|---|---|---|
| A-1a | BOLEH MULAI | ya — dari TERBATAS | — |
| A-1b | BOLEH MULAI | ya — dikuatkan | — |
| A-2 | BOLEH MULAI | tidak | — |
| A-3 | TERBATAS | tidak | PG-03 → B5-7 |
| A-4 | BEKU | tidak | PG-04 → satu `SELECT` |
| A-5 | BEKU | tidak | PG-05 → source tiga prosedur |
| A-6 | BEKU RINGAN | tidak | PG-06 → B5-4 |

Dua aliran yang membentuk inti modul — mekanika keputusan dan pembentukan sirkulasi —
keduanya kini **BOLEH MULAI**. Itu hasil ronde ini yang paling menentukan: yang menahan
keduanya ternyata bukan bukti yang hilang, melainkan berkas yang belum dibaca.
