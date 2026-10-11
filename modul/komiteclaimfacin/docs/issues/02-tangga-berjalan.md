# 02: Tangga berjalan — giliran, maju, selesai, tutup seketika

**Status:** ready-for-agent
**Blocked by:** 00 · 01
**Menutup:** AC 14 · 15 · 16 · 17 · 18 · 19 · 20 · 21 · 22 *(9 AC)* — US 14–18

> ⛔ **RALAT 10-10-2026.** Kepala lamanya dikutip utuh, tidak dihapus: *"**Blocked by:** 00 · 01"* dan *"**Menutup:** AC
> 14 · 15 · 16 · 17 · 18 · 19 · 20 · 21 · 22 *(9 AC)* — US 14–18"* → tiket 01 **gugur**, jadi tiket ini hanya menunggu
> 00. **Eskalasi gugur** (prompt tahap 2 §3: nol bukti di korpus): AC 20–22 dan US 17–18 tidak dibangun. Tiket ini
> menutup AC 14–19.

## Hasil & nilai pengguna

Hari ini Tangga persetujuan **belum berjalan**: tidak ada yang menentukan siapa menerima giliran, kapan tangga maju, dan kapan ia berhenti.

Sesudah tiket ini, Giliran berpindah ke **pemegang jenjang terendah yang belum memutuskan**. ⭐ Komite **selesai** ketika seluruh jenjang menyetujui, dan **tertutup seketika** ketika satu menolak.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Penugasan berputar | router menelusuri susunan jenjang; anggota **pertama yang belum memutuskan** menerima giliran, lalu aktivitas **keluar seketika** |
| ⚠️ Pemutus lama | ⛔ di Pega tangga berhenti karena router **tidak menghasilkan penerima** — ⚠️ **perilaku diam yang tidak dapat diuji** |
| Eskalasi | kasus dapat dinaikkan **satu jenjang** bila pemegang giliran berhalangan; ⛔ turun jenjang **dilarang** |

> ⛔ **RALAT 10-10-2026.** Baris lamanya dikutip utuh, tidak dihapus: *"Penugasan berputar | router menelusuri susunan
> jenjang; anggota **pertama yang belum memutuskan** menerima giliran, lalu aktivitas **keluar seketika**"*, *"⚠️ Pemutus
> lama | ⛔ di Pega tangga berhenti karena router **tidak menghasilkan penerima** — ⚠️ **perilaku diam yang tidak dapat
> diuji**"* dan *"Eskalasi | kasus dapat dinaikkan **satu jenjang** bila pemegang giliran berhalangan; ⛔ turun jenjang
> **dilarang**"* →
>
> - **Router Fac In tidak keluar seketika.** `KomiteRouter` S6.1 (`AssignTo := .KomiteID` bila `KomiteAproval == 0`)
>   tidak punya transisi keluar, jadi di Pega `AssignTo` = baris menunggu **terakhir**; S1–S5 dan S7 ter-remark. Sistem
>   baru memakai **baris tingkat berjalan** (`KomiteCount`), selaras prompt §5 #1.
>   ⛔ **RALAT 10-10-2026 (kedua)** atas kalimat di atas, dikutip utuh: *"tidak punya transisi keluar, jadi di Pega `AssignTo` = baris menunggu **terakhir**"* → `KomiteRouter` S6.1 **punya**
>   transisi pasca-langkah (`pyStepsTransParams` `true` → 6 keluar; alat dump sebelumnya hanya mencetak prasyarat):
>   Pega menugaskan baris menunggu **pertama**. Baris tingkat berjalan (`KomiteCount`) sistem baru setara pada alur
>   berurutan; prompt §5 #1 tetap berlaku untuk `SetProteksiSubmiteKomite` L2.1 (lolos untuk baris menunggu mana saja).
> - **Pemutus lama bukan router kosong.** Sesudah `KomitePostAct`, Decision `IsKomiteLoop` (`AcceptStatus == "1"` AND
>   `KomiteCount <= KomiteLoop`) mengulang assignment atau mengakhiri alur. `KomitePost_Adjustment` S24 (pre=false)
>   **selalu** menaikkan `KomiteCount` satu, dan S14 menyetel `KomiteCount := KomiteLoop` saat ditolak, sehingga sesudah
>   tingkat akhir atau penolakan `KomiteCount = KomiteLoop + 1` dan alur selesai. Sistem baru menulis syarat ini
>   eksplisit (pola Komite Prop `MasihBerjalan`) dan memajukan `KOMITE_COUNT` + `POSITION` ke workbasket berikut, atau
>   mengosongkan `POSITION` + mengisi `STATUS_WORK` saat selesai. Hasilnya tetap K1.
> - **Eskalasi gugur**: tidak ada rule eskalasi di korpus.

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

> ⛔ **RALAT 10-10-2026.** Butir dan langkah lamanya dikutip utuh, tidak dihapus: *"Eskalasi **naik satu jenjang**
> tersedia; ⛔ **turun jenjang ditolak**"*, *"Eskalasi **terekam**: siapa memindahkan, kapan, dari jenjang mana ke
> mana"*, *"Jenjang yang **dilewati** eskalasi tercatat **tanpa keputusan** — ⛔ bukan sebagai menyetujui"* dan *"5.
> Keluarkan pengaju lebih dulu ⇒ ⭐ tangga tetap benar walau jumlah jenjang **berubah di awal**."* → ketiga butir eskalasi
> **gugur**. Langkah 5 diganti: **perluasan** tingkat 1 (KCF-02) menambah jenjang di awal, lalu `KOMITE_LOOP` dihitung
> ulang ⇒ tangga selesai sesudah jenjang terakhir hasil perluasan. Tambahan dari XML: penolakan di tingkat mana pun
> menandai sisa baris menunggu `"2"` (`KomitePost_Adjustment` S7.2.1.9.1.1) lalu kasus tutup seketika; Submit ditolak
> bila kasus tertutup atau tingkatnya sudah memutus.
