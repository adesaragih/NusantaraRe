---
status: tertahan
---

# 09: Layar persetujuan komite dengan tiga masukan dan kunci yang ditegakkan dua kali

*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Pemegang jenjang membuka layar persetujuan, membaca seluruh nilai sebagai
proyeksi baca-saja, dan mengisi **tiga hal saja** — keputusan, komentar, dan penanda usulan.
Empat penanda berkunci jenjang pertama; medan induk yang baca-saja **tidak memblokir**
keputusan.

Layar persetujuan komite; penegakan kunci sisi klien.

**Persyaratan:** `S-058`, `S-059`, `S-060`, `S-061`, `S-062`, `S-063`.

**Tidak termasuk:** kebenaran tiap aturan kunci — ia diuji **sekali** di seam perintah, bukan
dua kali. Seam layar sempit dengan sengaja: ia memeriksa bahwa penegakan pertama **ada**.

**Jalur gagal:** keputusan tanpa komentar → `422`, keputusan tidak tersimpan · penyuntingan
penanda dari jenjang bukan pertama → `422` **di server**, meski klien mengizinkannya · medan
di luar ketiga masukan terkirim → diabaikan, dan ditolak `422` bila berbeda dari nilai
tersimpan.

**Uji:** BARU — tiap kunci diperiksa **ada** di klien, lalu diperiksa **ditegakkan** di
server dengan permintaan yang melewati klien. Kunci yang hanya ada di klien **bukan kunci**,
dan uji itu yang membuktikannya.

**Menggantikan:** `ShowTransfer` — layar persetujuan · `ViewTransferDtl` — pembuka layar ·
`ResetSubjectivityNote` · medan rekening, 16 kemunculan, **seluruhnya baca-saja** (`F-1`,
`F-6`) · 16 properti nilai, **seluruhnya baca-saja** (`F-9`) · medan bersufiks dua,
**dibuang** — nol pembaca di Activity modul ini (`F-3`) · blok bersyarat yang syaratnya tidak
pernah dapat benar dan blok bersyarat "tidak pernah", **dibuang** · dua `pyDisabledWhen` yang
menunjuk kendali **yang sudah baca-saja**, **dibuang** — lihat `FAKTA-LAYAR-01.md`.

**Blocked by:**

- `04` — Keputusan tercatat pada jenjang yang benar, dan orang lain ditolak


**Dasar:** DECIDED(`E-4`, `J-2`, `J-1`, ADR-0006). EVIDENCED: `F-1`, `F-3`, `F-6`, `F-7`,
`F-8`, `F-9`; `3-to-tickets/FAKTA-LAYAR-01.md` — enam `pyDisabledWhen`, keenamnya terbaca
dari `pyValue` selnya sendiri, bukan dari kedekatan posisi.

- [ ] Satu-satunya masukan komite adalah keputusan, komentar, dan penanda usulan.
- [ ] Seluruh panel nilai baca-saja; komite **tidak dapat** mengubah satu pun angka uang dan
      **tidak pernah** menyunting rekening penerima.
- [ ] Keempat penanda berkunci jenjang pertama adalah `.IsSubjectivity`, `.SubjectivityNote`,
      `.Adjustment.IsProposeClose`, `.Adjustment.IsPropReserved` — bernama, bukan dicacah.
- [ ] Label mata uang baca-saja bagi **seluruh** jenjang.
- [ ] Medan induk baca-saja **tidak memblokir** keputusan; komentar wajib.
- [ ] Tiap blok yang dibuang **tercatat beserta syarat lamanya**, bukan dihilangkan diam-diam.
- [ ] Tiap kunci ditegakkan di klien **dan** di server; nol kunci yang hanya ada di klien.

**Ketidakpastian:** **Bentuk kendali catatan syarat.** `.SubjectivityNote` di sistem lama
adalah `pxDropdown`, bukan kotak teks, sementara `SPEC-KOMITE-01.md` menuliskannya "teks
panjang" dan kolomnya `VARCHAR2(2000 CHAR)`. **Isi daftar pilihannya tidak terbaca**, dan itu
lubang bukti yang sah: `Rule-Obj-FieldValue` memang tidak ikut diekspor —
`INVENTARIS-BUKTI.md` §1, baris "XML rule Komite", kolom "Apa yang TIDAK diliput", entri
**Field Value**. Baris itu sudah ada; tidak ada baris baru ditambahkan.
