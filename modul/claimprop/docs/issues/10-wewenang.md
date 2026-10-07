# 10: Wewenang — satu sumber, ditegakkan di lapisan layanan

**Status:** dibangun 07-10-2026 — workbasket Input Acceptation = `ReasKlaimTeknik` (keputusan 07-10-2026) — RALAT 07-10-2026 (semula `ready-for-agent`)
**Blocked by:** 00 (PREFACTOR)
**Menutup:** AC 54 · 55 · 56 · 57 · 58 · 59 · 60 · 61 *(8 AC)* — US 55–58

## Hasil & nilai pengguna

Tingkat wewenang setiap pelaku dibaca dari **satu** sumber data, sehingga tidak ada dua daftar yang
berbeda. Pergantian pemegang jabatan cukup dengan mengubah baris tabel — tanpa menyentuh kode. Dan
yang paling penting: gerbang wewenang **benar-benar menolak**, bukan sekadar menyembunyikan tombol.

⚠️⚠️ `[terverifikasi]` **Pega tidak punya satu pun gerbang wewenang blocking di sisi server:** nol
`pyValidateActivity` di 12 dari 12 FlowAction; nol pemeriksaan peran di Activity, 418 ekspresi gerbang
di 36 dari 36 Section, 19 ReportDefinition, dan 11 Harness. Gerbang penyerahan komite sepenuhnya
*client-side*. **Membangunnya di Go bukan penyimpangan — ini mengisi lubang yang memang tidak pernah
ditutup.**

## Area codebase

Roster komite · resolusi tingkat wewenang · penegakan di lapisan layanan · batas nilai Direktur Utama.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `DataTransform/InsertChronology_DT.xml` | 1.1.4 | tingkat dasar diset **tanpa syarat** lebih dulu |
| `DataTransform/InsertChronology_DT.xml` | 1.1.5 · 1.1.6 · 1.1.7 | ⚠️ tiga cabang bernama orang — **nilai dasar + tiga tingkat komite**, bukan empat tingkat sejajar |
| `ReportDefinition/FilterEmailKomiteWithLimit.xml` | — | roster; penyaringnya **tiga** — batas bawah pita, status klaim, status aktif |
| `RDBList/GetLimitDirekturUtama_SQL.xml` | — | ⚠️ mencari menurut **nama orang** yang di-hardcode |
| `Activity/AttachmentProtect_ACT.xml` | 5 | batas nilai Direktur Utama menggerbangi **kewajiban dokumen**, bukan tangga persetujuan |
| `Activity/AddKomiteTreatyChild_ACT.xml` · `Activity/SetKomiteTreaty_ACT.xml` | — | pemakai roster |

## ADR terkait

**ADR-0014** (keputusan komite ditegakkan di lapisan layanan — ⚠️ cakupan asli ADR-nya Komite Claim
Life; yang dipakai di sini **prinsipnya**).

## Acceptance criteria

- [ ] ⚠️ `[keputusan work owner]` Tingkat wewenang dibaca dari **satu sumber**: roster komite. Tidak ketemu → tingkat dasar Claim Admin. **Alasan menyimpang:** di Pega tingkat dasar diset tanpa syarat lalu ditimpa tiga cabang bernama orang *(AC 54 spec)*
- [ ] ⚠️ `[keputusan work owner]` **Empat nama orang yang di-hardcode dibuang** — tiga di transform jejak audit dan satu di rule batas Direktur Utama. Perilaku tidak berubah, tetapi pergantian pemegang jabatan cukup lewat baris tabel. **Alasan menyimpang:** identitas dunia nyata dipaku ke dalam kode; orangnya pindah jabatan → jejak audit salah label tanpa peringatan *(AC 55 spec)*
- [ ] `[keputusan work owner]` Batas nilai **Direktur Utama** adalah wewenang **terpisah**, bukan tingkat kelima tangga komite *(AC 56 spec)*
- [ ] ⚠️ Seluruh gerbang wewenang **ditegakkan di lapisan layanan**; test yang menemukan gerbang hanya di UI **gagal**. **Catatan pengisian lubang:** ini **bukan penyimpangan** — spec menyatakannya terang di Implementation Decision 7 (*"Ini bukan penyimpangan — melainkan mengisi lubang yang memang tidak pernah ditutup"*). Pega tidak punya satu pun gerbang blocking di sisi server, jadi tidak ada perilaku yang disimpangi *(AC 57 spec)*
- [ ] Test: panggil endpoint penyerahan komite langsung, tanpa UI → **ditolak** bila syarat tidak terpenuhi *(AC 58 spec)*
- [ ] `[terverifikasi]` Jumlah tingkat komite = **COUNT baris roster aktif** yang pita limitnya mencakup nilai klaim, dihitung saat penyerahan — **tidak** dari konstanta *(AC 59 spec)*
- [ ] Batas **atas** pita roster **tidak** dijadikan penyaring; penyaring tetap batas bawah + status klaim + status aktif. **Catatan paritas:** bila batas atas ikut menyaring, hanya satu pita lolos dan tangga berjenjang hilang *(AC 60 spec)*
- [ ] ✅ `[data DBA]` Kunci pencocokan pelaku ke roster = **kolom operator id pada roster**, dikonfirmasi sama dengan user id Pega. Kolom id pengguna baru **batal**; kolom email tidak dipakai sebagai kunci. **Nol perubahan skema** untuk bagian ini *(AC 61 spec)*

## Catatan `[terbuka]` ringan — TIDAK memblokir

⚠️ `[terbuka]` Roster memuat **enam kolom status per jenis aksi** — adjustment, adjuster, registrasi,
reject, salvage, survey — tetapi **nol** dipakai menyaring oleh rule mana pun (sensus 329 berkas).
Apakah wewenang per jenis aksi memang dimaksudkan berlaku, atau keenam kolom itu warisan yang tidak
terpakai? Pemilik: **Finance + work owner**. Bawa apa adanya; yang ditiru adalah tiga penyaring yang
berjalan.

⚠️ **Jebakan nama:** salah satu nama kolom itu juga muncul sebagai properti muatan layanan konversi
non-life. Tabrakan nama, bukan pemakaian — jangan salah simpulkan kolom roster itu terpakai.

## Perintah verifikasi

```
jalankan test "pelaku tidak ada di roster -> tingkat dasar Claim Admin"
jalankan test "pelaku ada di roster -> tingkat dari kolom tingkat roster"
jalankan test "panggil endpoint kirim-komite langsung tanpa hak -> DITOLAK di layanan"
jalankan test "jumlah tingkat komite = jumlah baris roster yang lolos, bukan konstanta"
jalankan test "batas Direktur Utama tidak menambah tingkat tangga komite"
cari nama orang yang di-hardcode di kode dan SQL          -> nihil
```
