# 12: Kronologi komite dan penandaan data lama

**Status:** ready-for-agent
**Blocked by:** 05 · 11
**Menutup:** AC 86 · 87 · 88 · 89 · 90 · 91 · 92 · 93 · 94 · 95 *(10 AC)* — US 24 · 25 · 27

## Hasil & nilai pengguna

Hari ini Catatan kronologi komite **disimpan sebagai kalimat jadi**, sehingga ⚠️ jabatan yang keliru **membeku di dalam teks**. ⚠️ Dan catatan akseptasi lama lini **MBU dan Travel** dibangun **tanpa cabangnya** — tanpa penanda apa pun yang membedakannya.

Sesudah tiket ini, ⭐ Catatan kronologi disusun dari **medan tersimpan** saat ditampilkan, ⛔ bukan disimpan sebagai kalimat jadi. Dan ⭐ **catatan lama MBU/Travel DITANDAI** sehingga laporan **dapat menyaringnya**.

## Perilaku Pega yang ditiru

| Yang dibaca | Rule |
| --- | --- |
| Catatan kronologi | setiap keputusan komite menambah satu catatan |
| ⛔ RALAT penting | ⛔ sasarannya **catatan kronologi kasus** — ⭐ **jejak internal**, ⚠️ **bukan** dokumen akseptasi yang keluar perusahaan |
| ⚠️ Nilai bawaan | ⛔ siapa pun di luar daftar tercatat sebagai **jabatan direksi** |

⭐ Sumber: `komite-claim-facin\spec.md` · `claim-facin\STRUKTUR-TABEL-CLAIM-FACIN.md` §5.

## Keputusan work owner yang mengikat

- **K8** — ⭐ **Data lama MBU dan Travel DIBIARKAN, tetapi DITANDAI** — ⛔ tidak dibangun ulang, sebab angkanya mungkin **sudah dilaporkan keluar**
- **K13** — ⭐ Catatan disusun dari **medan tersimpan** — akun, kode jabatan, keputusan, waktu

## Yang harus diuji

- [ ] ⭐ Setiap keputusan komite menambah **catatan kronologi**
- [ ] ⭐ Catatan disusun dari **medan tersimpan**, ⛔ **bukan disimpan sebagai kalimat jadi**
- [ ] ⛔⛔ Menghapus kasus **TIDAK menghapus** kronologinya — ⭐ jejak yang ikut terhapus **berhenti menjadi jejak**
- [ ] ⛔ **RALAT:** sasaran catatan itu **jejak internal**, ⚠️ bukan dokumen keluar — ⭐ tetapi akibatnya tetap serius: **jejak audit mencatat jabatan yang keliru**
- [ ] ⭐ Catatan akseptasi lama MBU/Travel **ditandai** dengan kolom penanda ber-nilai dua keadaan
- [ ] ⭐ **Alasan penandanya KOLOM, bukan catatan prosa:** supaya **laporan dapat menyaringnya** — ⛔ bukan hanya pembaca manusia yang kebetulan membaca catatan kaki
- [ ] ⛔ Catatan **baru** tidak menerima penanda itu
- [ ] ⛔ Data lama **tidak dibangun ulang**

## Butir `[terbuka]` yang menyentuh tiket ini

| Butir | Isinya | Menahan? |
| --- | --- | --- |
| **29** | Cacah baris MBU/Travel terdampak — `[data DBA]` | ⚠️ menahan **cacahnya**, ⭐ tidak menahan cara penandaannya |

## Seam & verifikasi

**Seam:** lapisan layanan komite — penulisan kronologi dan penandaan migrasi.
1. Putuskan satu jenjang ⇒ ⭐ kronologi bertambah satu baris.
2. Hapus kasus komitenya ⇒ ⛔ kronologinya **tetap ada**.
3. Ubah tulisan sebuah jabatan di daftar induk ⇒ ⭐ catatan lama **ikut tulisan baru**, ⛔ tetapi **tidak berpindah ke jabatan yang berbeda**.
4. Periksa catatan akseptasi lama MBU ⇒ ⭐ penandanya **menyala**.
5. Terbitkan catatan **baru** ⇒ ⛔ penandanya **padam**.
6. Jalankan laporan dengan saringan penanda ⇒ ⭐ kedua kelompok **terpisah bersih**.
