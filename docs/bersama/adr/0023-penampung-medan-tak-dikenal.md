---
status: accepted
tanggal: 2026-09-23
sumber: spec penyimpanan NB Treaty In - `.scratch/nb-treaty-in/`, keputusan work owner
---

# Pemecah dokumen wajib punya penampung medan tak dikenal

Setiap pemecah dokumen JSON menyediakan **penampung** bagi medan yang tidak dikenal daftar
kolomnya. Medan yang tidak dikenali **masuk penampung**, tidak dibuang.

Penampung itu **wajib kosong** sebelum pekerjaan dinyatakan selesai.

## Kenapa

ADR-0017 menetapkan daftar medan disusun dari sapuan aturan. Sapuan itu jujur tetapi **tidak dapat
membuktikan kelengkapan** - cacahnya adalah **batas bawah**, bukan total. Medan yang tidak pernah
dirujuk aturan mana pun tetap dapat hadir pada dokumen tersimpan.

Penampung adalah jaring pengaman yang menggantikan kepastian yang memang tidak ada.

Tanpa penampung, medan yang tidak dikenali **hilang diam-diam**. Uji pulang-pergi tidak akan
menangkapnya, sebab yang hilang tidak pernah masuk untuk dibandingkan.

## Akibat

1. Pemecah dokumen menulis medan tak dikenal ke penampung, bukan mengabaikannya dan bukan gagal.
2. Isi penampung diperiksa sebagai bagian dari kriteria selesai. Penampung berisi = pekerjaan
   belum selesai.
3. Test yang menemukan medan hilang tanpa masuk penampung **gagal**.
4. Berlaku untuk setiap modul yang memecah dokumen JSON tersimpan.
