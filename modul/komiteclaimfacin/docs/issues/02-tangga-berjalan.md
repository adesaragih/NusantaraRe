# 02: Tangga berjalan — giliran, maju, selesai, tutup seketika

**Status:** ready-for-agent
**Blocked by:** 00 · 01
**Menutup:** AC 14 · 15 · 16 · 17 · 18 · 19 · 20 · 21 · 22 *(9 AC)* — US 14–18

## Hasil & nilai pengguna

Hari ini Tangga persetujuan **belum berjalan**: tidak ada yang menentukan siapa menerima giliran, kapan tangga maju, dan kapan ia berhenti.

Sesudah tiket ini, Giliran berpindah ke **pemegang jenjang terendah yang belum memutuskan**. ⭐ Komite **selesai** ketika seluruh jenjang menyetujui, dan **tertutup seketika** ketika satu menolak.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Penugasan berputar | router menelusuri susunan jenjang; anggota **pertama yang belum memutuskan** menerima giliran, lalu aktivitas **keluar seketika** |
| ⚠️ Pemutus lama | ⛔ di Pega tangga berhenti karena router **tidak menghasilkan penerima** — ⚠️ **perilaku diam yang tidak dapat diuji** |
| Eskalasi | kasus dapat dinaikkan **satu jenjang** bila pemegang giliran berhalangan; ⛔ turun jenjang **dilarang** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K1** — ⭐ **Komite berakhir bila SELURUH jenjang menyetujui; satu menolak ⇒ tutup seketika.** ⭐ Pemutusnya **keadaan tiap jenjang**, ⛔ **bukan pencacah**

## Yang harus diuji

- [ ] Giliran diberikan ke **pemegang jenjang terendah yang belum memutuskan**
- [ ] ⭐ Seluruh jenjang menyetujui ⇒ kasus **selesai**
- [ ] ⭐ Satu jenjang menolak ⇒ kasus **tertutup seketika**; ⛔ jenjang berikutnya **tidak** menerima giliran
- [ ] ⛔ Sebuah jenjang **hanya dapat memutuskan sekali**; percobaan kedua **ditolak**
- [ ] Eskalasi **naik satu jenjang** tersedia; ⛔ **turun jenjang ditolak**
- [ ] Eskalasi **terekam**: siapa memindahkan, kapan, dari jenjang mana ke mana
- [ ] Jenjang yang **dilewati** eskalasi tercatat **tanpa keputusan** — ⛔ bukan sebagai menyetujui

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **6** | Isi susunan jenjang | ⚠️ menahan jumlah tingkat yang diuji |

## Seam & verifikasi

**Seam:** lapisan layanan komite — masukkan susunan jenjang, jalankan urutan keputusan.
⛔ **Uji PERILAKU, bukan pencacah.** ⭐ Pencacah boleh ada sebagai tampilan; ia **bukan kebenaran**.
3. Tiga jenjang, semua setuju ⇒ ⭐ selesai sesudah yang **ketiga**.
4. Tiga jenjang, yang **kedua menolak** ⇒ ⭐ tutup **seketika**; jenjang ketiga **tidak** menerima giliran.
5. Keluarkan pengaju lebih dulu ⇒ ⭐ tangga tetap benar walau jumlah jenjang **berubah di awal**.
