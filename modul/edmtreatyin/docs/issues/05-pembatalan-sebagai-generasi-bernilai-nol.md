# 05: Pembatalan sebagai generasi bernilai nol

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Polis yang sudah dibatalkan **masih boleh di-endorse lagi**. Nol gerbang peran — layar cukup memberi **peringatan**, bukan penghalang.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 4.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*
~~**Blocked by:** ⛔ `[work owner]` **sesudah dibatalkan, polis masih boleh di-endorse lagi atau tidak**~~ ⛔ **penahan gugur 23-09-2026**
**Bergantung pada tiket NB:** **19** *(pemecah dokumen — kolom jenis berkas)*
**Menutup:** AC **14–15** *(2 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-19 · ID-41

## Hasil & nilai pengguna

Pembatalan polis adalah **jenis endorsemen**, dipilih di awal saat berkas dibuat — bukan alur
tersendiri dan bukan tombol terpisah.

⛔⛔ **Pembatalan BUKAN penghapusan.** Ia menghasilkan **generasi baru berisi nol**, sehingga
jejaknya tetap ada dan dapat dipertanggungjawabkan.
⛔ **Tiket yang membuat jalur hapus untuk pembatalan salah.**

## Yang dibangun

Jenis endorsemen tersimpan di kolomnya sendiri, dan salah satu nilainya berarti pembatalan.
Memilihnya menghasilkan generasi baru yang seluruh kolom uangnya **bernilai nol**, lewat jalur
penyimpanan yang **sama** dengan endorsemen lain.

⛔ **Nol jalur hapus dibangun untuk pembatalan.**

## Batas — yang TIDAK termasuk

⛔ Aturan nomor urut — tiket **04**; pembatalan **tidak** menghapus baris, jadi aturannya berlaku apa adanya.
⛔ Perhitungan selisih terhadap generasi yang dibatalkan — tiket **06**.

## Cara mengujinya

Lewat seam `repository`. ⭐ **Uji utama:** polis dibatalkan, lalu **cacah baris rincian dihitung**
— jumlahnya harus **sama** dengan generasi sebelumnya, dan seluruh kolom uangnya nol.
⛔ Baris yang berkurang berarti gagal.

## Acceptance criteria

- [ ] **AC 14** — pembatalan menghasilkan generasi baru dengan kolom uang **bernilai nol**, bukan penghapusan
- [ ] **AC 15** — jenis endorsemen tersimpan di kolomnya sendiri dan dipilih di awal

## ⛔ Kenapa tiket ini `blocked`

`[work owner]` **Belum dijawab apakah polis yang sudah dibatalkan masih boleh di-endorse lagi.**
Jawabannya menentukan apakah generasi pembatalan menjadi **ujung rantai** — sehingga penunjuk
generasi berikutnya harus ditolak — atau sekadar generasi biasa yang nilainya nol.

⭐ **Kedua AC dapat dibangun tanpa menunggu**; yang menunggu hanya **penjaga rantai sesudahnya**.