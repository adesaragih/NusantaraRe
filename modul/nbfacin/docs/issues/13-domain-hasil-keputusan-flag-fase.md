# 13: Domain hasil keputusan, dan satu flag fase untuk jalur yang belum terverifikasi

**What to build:** Hasil keputusan underwriting punya **enam nilai yang sah dan artinya pasti**, dan
sistem berperilaku berbeda-tetapi-terkendali pada jalur yang belum terverifikasi.

**Enam nilai, artinya terverifikasi dari dekoder literal di korpus:** Accept · Reject · Ask ·
**Banding** · Decline · Revise. Dua nilai lain milik ranah ceding, hidup di properti terpisah, dan
tidak pernah menulis ke kolom ini; satu nilai lagi nol jejak di seluruh korpus.

**`Reject` dan `Decline` tidak disatukan.** Decline adalah penolakan **final tanpa hak banding**;
Reject masih **dapat dibanding**. Efek samping Reject yang terverifikasi — mematikan konfirmasi
binding dan penerimaan slip, lalu menjadi syarat munculnya jalur banding lewat hitungan riwayat
berstatus reject — adalah **perilaku yang dikehendaki**, direproduksi apa adanya.

⚠️ **Validasi domain ini adalah perilaku BARU.** Sistem lama tidak punya validasi domain pada kolom
ini — tidak punya validasi apa pun selain wajib-isi. Karena itu ia **tidak boleh** berbentuk penolakan
diam-diam.

**Satu flag fase, bukan dua basis kode:**

| Situasi | Paralel run | Produksi |
| --- | --- | --- |
| Hasil keputusan **dikenali**, baris keputusannya belum terverifikasi | gagal keras | `decline` + catatan |
| Kondisi **tidak dikenali sama sekali** | gagal keras | gagal keras |
| Nilai **di luar domain** enam nilai | gagal keras + catatan | gagal keras di jalur tulis + catatan |

`decline` dipilih di produksi karena ia **hasil default yang terekam** — memilihnya adalah reproduksi,
bukan tebakan. Transisi sebuah jalur dari gagal-keras ke `decline` **hanya** boleh setelah jalur itu
diverifikasi terhadap ekspor produksi, bukan karena ia sering muncul dan mengganggu.

**Blocked by:** 11

**Status:** ready-for-human — diport 01-10-2026, menunggu tinjauan work owner

- [ ] Enum enam nilai, masing-masing menyebut rule Pega asal dekodernya dalam komentar
- [ ] Reject dan Decline **tetap terpisah**, dengan efek samping Reject direproduksi
- [ ] **Satu** flag fase eksplisit mengendalikan perbedaan perilaku — bukan dua basis kode
- [ ] Jalur "dikenali tetapi belum terverifikasi" berperilaku sesuai fase
- [ ] Jalur "tidak dikenali sama sekali" gagal keras di **kedua** fase
- [ ] Nilai di luar domain **tidak pernah ditolak diam-diam**; selalu disertai catatan
- [ ] Dicatat sebagai kandidat perbaikan bahwa validasi domain ini **tidak ada di sistem lama**
- [ ] Bila paralel run menemukan nilai lain benar-benar ada di data produksi, itu **membatalkan premis** dan dibuka sebagai keputusan baru — **bukan** alasan melonggarkan validasi diam-diam

## Comments

### 2026-10-01 — implementasi (agent)

Kode: `APP_RNM/modul/nbfacin/backend/services/acceptance/keputusan.go`.

| Kriteria | Keadaan |
| --- | --- |
| Enum enam nilai + asal dekoder | ✅ `IsUWAccepted.xml` L302–L306 → L343–L347, default `decline` L91; penulis tiap nilai (6 DataTransform) dicek per baris |
| Reject ≠ Decline, efek Reject direproduksi | ✅ `MatikanBindingDanRISlip`; ⚠️ temuan: **Revise** punya efek yang sama (`SetReviseProposal` L179/L209) |
| Satu flag fase | ✅ `Fase` (ParalelRun / Produksi) |
| Jalur (a) dikenali-belum-terverifikasi | ✅ `JalurBelumTerverifikasi`: panic / Decline + catatan |
| Jalur (b) tidak dikenali | ✅ `KondisiTakDikenali`: panic di kedua fase |
| Luar domain tidak ditolak diam-diam | ✅ `PeriksaDomain` selalu mencatat; panic / `ErrDiLuarDomain` |
| Kandidat perbaikan: validasi domain tak ada di sistem lama | ✅ tercatat di komentar |

⚠️ Salinan `IsUWAccepted` folder Endorsment berbeda isi baris (spec menyebut identik di ketiga folder);
NB dan RNW identik byte. Pemakai `Fase` dan ketiga fungsi itu adalah tangga NB-11, yang masih ditahan.
Uji mutasi: 6/6.

### 2026-10-01 — tindak lanjut review (agent)

- 🐞 Diperbaiki: `Fase(0)` (tidak diisi) dulu diperlakukan diam-diam sebagai Produksi; kini panic
  (`TestFaseNolDitolak`).
- Efek Revise = efek Reject dicatat sebagai keputusan agent A10.
- ⏸ Efek Reject yang ketiga — syarat jalur banding lewat hitungan riwayat berstatus REJECT
  (`Protection_Act` langkah 19–20, `GetFlagReject_SQL`) — milik tangga NB-11, belum diport.
- Bukti penulis nilai kini menyebut berkas dan kelas (`NB FacIn\DataTransform\…`, ASM-FW-GISFW-Work).
