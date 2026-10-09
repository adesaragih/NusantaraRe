# 03: Layar komite — 93 kolom, tiga wajah menurut jalur

**Status:** dibangun — wajah TT 2 saja *(RALAT 08-10-2026; status lama: `ready-for-agent`)*

> **RALAT 08-10-2026** — implementasi satu modul (prompt `_brief/PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-PROP.md` §7). Kalimat lama tetap di bawah, dikutip di sini:
>
> - *Tiga wajah* — TT 2 (ADJUSTMENT) dibangun; TT 3 (REJECT) tanpa penulis di Claim Prop (dilaporkan, tak terjangkau); TT 4 (CLOSE) ditunda **OQ-CP-06** (keputusan work owner 08-10-2026).
> - ~~Tombol "View more details" tampil **nonaktif** (harness `ViewClaimFormKomite` tidak diekspor).~~ RALAT kedua 08-10-2026 (harness diekspor work owner): tombol **aktif**, membuka berkas Claim Prop klaim induk hanya-baca di jendela di atas layar komite (`PropsRute.onLihatBerkas`; section tab harness tidak diekspor). Butir terbuka 9 (pintu masuk `ViewDetailInterest`) tetap.
> - Prompt values diekspor work owner 08-10-2026: label `KomiteAproval` (grid Committe Accept Status), `AcceptanceStatus` (History Adjustment, daftar kerja), `Payable`; "Subjectivity Note" = dropdown `SubjectivityNote.xml` (kode 1-7 disimpan, divalidasi server).
> - AC 84 (wewenang hapus kronologi) ditempel di sini: aksi itu tidak ada di Section `ShowTransfer` wajah TT 2 — nol tombol dibangun.


**Blocked by:** **01 (kasus komite lahir)**

## Hasil & nilai pengguna

Sebagai **penyetuju komite**, saya membuka kasus dan melihat **data klaim induknya** — tertanggung,
nomor polis, tanggal kejadian, penyebab, lokasi — **beserta nilai penyesuaian yang diajukan**.
Layarnya **menyesuaikan diri dengan jenis pengajuan**, jadi saya tidak dibingungkan kolom yang tidak
relevan.

*(User story 3, 4, 6, 7, 8, 9, 10, 11, 12 di spec)*

## Bentuk layar

`[terverifikasi]` **93 kolom** — **62 milik kasus komite**, **31 milik kasus klaim induk**. Angka
ini **SAH** `[keputusan work owner]` 2026-09-18.

⭐ **Tiga wajah menurut jenis pengajuan.** Dari 62 kolom komite, **24 hanya tampil pada jalur
penyesuaian**. Tiga judul bagian layar juga bergerbang jenis pengajuan.

⚠️ **Dua kolom yang wajib ikut terbawa meski di luar hitungan 93:** **"dibuat oleh"** dan
**"tanggal dibuat"**. Keduanya **tampil di layar**; keduanya di luar hitungan hanya karena aturan
pencatatan memisahkan properti bawaan Pega.

## Perilaku Pega yang ditiru

| Rule | Perilaku yang ditiru |
| --- | --- |
| `Komite Claim Prop/FlowAction/ViewTransferDtl.xml` · `ASM-FW-GCNMFW-WORK-KOMITETREATY!VIEWTRANSFERDTL` | `[terverifikasi]` **pembungkus, nol kolom sendiri**. Kontraknya: penyiap memuat → layar menampilkan → penyimpan menyimpan |
| `Komite Claim Prop/Section/ShowTransfer.xml` · `ASM-FW-GCNMFW-WORK-KOMITETREATY!SHOWTRANSFER` | `[terverifikasi]` 348 sel; **dua famili gerbang** — satu di tingkat sel, satu di tingkat layout. **12 sel di balik gerbang mati** |
| `Komite Claim Prop/Activity/SetKomiteList_Act.xml` | `[terverifikasi]` penyiap layar; **langkah 1** menghentikan activity bila jenis pengajuan bukan penyesuaian |
| `Komite Claim Prop/Section/ViewDetailInterest.xml` | `[terverifikasi]` layar kedua — **3 kolom, ketiganya hanya-baca**, dan **ketiganya sudah ada di layar utama** |

⭐ `[terverifikasi]` **Layar kedua tidak menambah satu kolom pun.**

## ⭐ Daftar penyetuju DITAMPILKAN — **perilaku BARU**

`[keputusan work owner — atas rekomendasi asisten]` 2026-09-19 — **layar komite menampilkan
daftar penyetuju**: **urutan** · **siapa** · dan untuk yang **sudah memutuskan** — **keputusan,
catatan, dan tanggalnya**. Yang **belum** memutuskan tampil **tanpa isi**.

⚠️ **Di Pega ia tidak tampil sama sekali.** `[terverifikasi]` Ketiga penggerak tangga —
`KomiteLoop`, `KomiteCount`, `KomiteList` — **nol di antaranya tampil di layar**, padahal
`KomiteList` justru berisi daftar penyetuju berikut keputusan, catatan, dan tanggalnya. **Akibatnya
di Pega: setiap penyetuju memutuskan tanpa konteks penyetuju sebelumnya.**

⛔ **NOL kolom basis data baru.** Ketiganya sudah tersimpan di baris penyetuju — kesembilan kolom
yang dikunci tiket **00**. Ini **murni penambahan tampilan**, dan ia **menutup US 3 dan US 4**,
dua user story yang sebelumnya **tidak ditutup satu AC pun**.

⚠️ **Tandanya `atas rekomendasi asisten`, jadi murah dicabut** — pencabutannya menyentuh **AC 79**
dan **AC 80**, satu baris bab 11, satu keputusan bab 10, dan dua butir uji di bawah. ⛔ **Nol kolom
tersentuh.** Alasan lengkapnya di `spec.md` **bab 4**.

## Keputusan work owner yang mengikat

| Tanggal | Keputusan |
| --- | --- |
| 2026-09-18 | Properti jenis treaty **dibuang** — hanya ada di balik gerbang mati, tidak pernah tampil |
| 2026-09-18 | Nomor klaim dan dua kotak total estimasi di balik gerbang mati **dibuang** |
| 2026-09-18 | **Angka 93 · komite 62 · klaim induk 31 SAH** |
| 2026-09-18 | Data klaim **dibaca** dari Claim Prop, **tidak digarap** modul ini |

## ⚠️ TITIK YANG SENGAJA DIUBAH — ketelitian tampilan

`[keputusan work owner]` 2026-09-18 — **angka di layar ditampilkan dengan 4 angka di belakang
koma**. Pega tidak punya aturan tampilan yang seragam. **Ini perubahan sadar.**

## Yang harus diuji

**Diverifikasi oleh:** spec.md AC 2 · 18 · 19 · 20 · 21 · 22 · 26 · 27 · 28 · 41 · **79** · **80**

- [ ] Membuka kasus komite menampilkan data klaim induk dan nilai penyesuaian yang diajukan.
- [ ] Pada jalur **penyesuaian** ke-24 kolom bergerbang **tampil**; pada jalur **tolak** dan
      **tutup** ke-24 kolom itu **tidak tampil**.
- [ ] Ketiga judul bagian layar berganti sesuai jenis pengajuan.
- [ ] Kolom yang tidak boleh diubah penyetuju tampil sebagai **bacaan saja**.
- [ ] **"Dibuat oleh"** dan **"tanggal dibuat"** tampil di layar.
- [ ] Properti jenis treaty, nomor klaim di balik gerbang mati, dan dua kotak total estimasi
      **tidak dibuat sama sekali**. Test yang menemukannya **gagal**.
- [ ] Angka uang di layar tampil dengan **4 angka di belakang koma**.
- [ ] ⚠️ Layar menampilkan **daftar penyetuju** berikut **urutan** dan **siapa**-nya, untuk
      **seluruh** tingkat — bukan hanya yang sudah lewat. Test yang tidak menemukannya **gagal**
      *(AC 79 spec)*
- [ ] ⚠️ Untuk penyetuju yang **sudah memutuskan**, layar menampilkan **keputusan, catatan, dan
      tanggal**-nya; yang **belum** memutuskan tampil **tanpa isi**. Test yang menemukan kolom
      basis data baru untuk ini **gagal** — datanya sudah ada di baris penyetuju *(AC 80 spec)*

## Butir `[terbuka]` yang menyentuh tiket ini

- **Dari mana layar kedua dibuka** belum terbaca — layar kedua tidak dirujuk layar utama,
  pembungkusnya, maupun berkas alur.

## Seam & verifikasi

Memakai ulang seam Claim Prop.
