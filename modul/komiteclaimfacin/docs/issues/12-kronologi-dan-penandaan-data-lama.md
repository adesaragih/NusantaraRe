# 12: Kronologi komite dan penandaan data lama

**Status:** ready-for-agent
**Blocked by:** 05 · 11
**Menutup:** AC 86 · 87 · 88 · 89 · 90 · 91 · 92 · 93 · 94 · 95 *(10 AC)* — US 24 · 25 · 27

> ⛔ **RALAT 10-10-2026.** Kepala lamanya dikutip utuh, tidak dihapus: *"**Blocked by:** 05 · 11"* dan *"**Menutup:** AC
> 86 · 87 · 88 · 89 · 90 · 91 · 92 · 93 · 94 · 95 *(10 AC)* — US 24 · 25 · 27"* → tiket 11 **gugur**, jadi tiket ini
> hanya menunggu 05. AC 92–95 dan US 27 gugur (KCF-04). AC 89–90 diganti kronologi teks (RALAT spec). AC 86–88 dan 91
> tetap.

## Hasil & nilai pengguna

Hari ini Catatan kronologi komite **disimpan sebagai kalimat jadi**, sehingga ⚠️ jabatan yang keliru **membeku di dalam teks**. ⚠️ Dan catatan akseptasi lama lini **MBU dan Travel** dibangun **tanpa cabangnya** — tanpa penanda apa pun yang membedakannya.

Sesudah tiket ini, ⭐ Catatan kronologi disusun dari **medan tersimpan** saat ditampilkan, ⛔ bukan disimpan sebagai kalimat jadi. Dan ⭐ **catatan lama MBU/Travel DITANDAI** sehingga laporan **dapat menyaringnya**.

> ⛔ **RALAT 10-10-2026.** Kalimat lamanya dikutip utuh, tidak dihapus: *"Sesudah tiket ini, ⭐ Catatan kronologi disusun
> dari **medan tersimpan** saat ditampilkan, ⛔ bukan disimpan sebagai kalimat jadi. Dan ⭐ **catatan lama MBU/Travel
> DITANDAI** sehingga laporan **dapat menyaringnya**."* →
>
> - **Kronologi = teks di `T_VIEW_SUGGEST`**, ikut XML dan pola claimfacin (OQ-CFI-19). `KomitePost_Adjustment` S3 / S4
>   (juga `KomitePost_Reject` S3 / S4, `KomitePost_CloseClaim` S2 / S3) menyusun `Data.CARI12` =
>   `"Accepted by " + jabatan + " - " + KMT-…` / `"Rejected by " + jabatan + " - " + KMT-…`, dengan jabatan =
>   `KomiteList(KomiteCount).IDKomite` (JABATAN roster). `ChronologyInsertion_DT` menambahkan satu baris ke kronologi
>   **kasus klaim induk** (`TempOpenPage.ClaimData.Chronology`): `.pyNote` = teks itu, `.ASMUser` = pemutus,
>   `.ASMNoteType` = label langkah (kosong → "Committee"). Ditulis
>   lewat kontrak di transaksi Submit. Pengecualian `pyPosition != "IT Developer"` dibuang: semua peran terekam (pola
>   claimfacin).
> - **K8 `DIBANGUN_ATURAN_LAMA` gugur** (KCF-04: kasus komite lama Pega `KMT-` tidak dimigrasi). Tidak ada kolom penanda,
>   tidak ada penandaan data lama.

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

> ⛔ **RALAT 10-10-2026.** Baris dan butir lamanya dikutip utuh, tidak dihapus: *"⚠️ Nilai bawaan | ⛔ siapa pun di luar
> daftar tercatat sebagai **jabatan direksi**"*, *"**K13** — ⭐ Catatan disusun dari **medan tersimpan** — akun, kode
> jabatan, keputusan, waktu"*, *"⭐ Catatan disusun dari **medan tersimpan**, ⛔ **bukan disimpan sebagai kalimat jadi**"*,
> *"⛔⛔ Menghapus kasus **TIDAK menghapus** kronologinya — ⭐ jejak yang ikut terhapus **berhenti menjadi jejak**"*,
> *"⭐ Catatan akseptasi lama MBU/Travel **ditandai** dengan kolom penanda ber-nilai dua keadaan"*, *"⭐ **Alasan
> penandanya KOLOM, bukan catatan prosa:** supaya **laporan dapat menyaringnya** — ⛔ bukan hanya pembaca manusia yang
> kebetulan membaca catatan kaki"* dan *"⛔ Catatan **baru** tidak menerima penanda itu"* →
>
> - Nilai bawaan direksi = sisa editor (spec RALAT lintas-ronde #2); jalur hidup memakai JABATAN roster.
> - K13 dan "bukan kalimat jadi" **gugur**: kronologi disimpan sebagai **teks jadi** di `T_VIEW_SUGGEST` (RALAT di atas).
>   Karena jabatannya diambil dari baris tangga, jabatan keliru atau bawaan tidak dapat muncul.
> - Penghapusan: XML Komite Claim FacIn tidak punya aksi hapus kasus komite, dan kronologinya milik **kasus klaim induk**,
>   bukan kasus `KMT-`. Butir ini tidak diuji di modul komite.
> - Tiga butir penanda **gugur** (K8, KCF-04). "Data lama tidak dibangun ulang" tetap: kasus komite lama tidak dimigrasi.

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

> ⛔ **RALAT 10-10-2026.** Langkah lamanya dikutip utuh, tidak dihapus: *"2. Hapus kasus komitenya ⇒ ⛔ kronologinya
> **tetap ada**."*, *"3. Ubah tulisan sebuah jabatan di daftar induk ⇒ ⭐ catatan lama **ikut tulisan baru**, ⛔ tetapi
> **tidak berpindah ke jabatan yang berbeda**."*, *"4. Periksa catatan akseptasi lama MBU ⇒ ⭐ penandanya **menyala**."*,
> *"5. Terbitkan catatan **baru** ⇒ ⛔ penandanya **padam**."* dan *"6. Jalankan laporan dengan saringan penanda ⇒ ⭐
> kedua kelompok **terpisah bersih**."* → langkah 2–6 **gugur**. Penggantinya: putuskan tingkat 1 setuju ⇒ satu baris
> `T_VIEW_SUGGEST` kasus klaim induk ber-teks `"Accepted by " + JABATAN tingkat 1 + " - " + nomor KMT`; tolak di tingkat
> 2 ⇒ teks `"Rejected by " + …`; TT3 / TT4 ⇒ teks yang sama dari `KomitePost_Reject` / `KomitePost_CloseClaim`.
