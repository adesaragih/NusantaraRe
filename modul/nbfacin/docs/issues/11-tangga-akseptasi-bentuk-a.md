# 11: Tangga akseptasi bentuk standar — satu keputusan, satu transisi

**What to build:** Seorang underwriter mengambil satu keputusan, dan kasus berpindah **satu langkah**
pada tangga akseptasi — persis seperti hari ini.

⛔ **Ini mesin keadaan satu langkah, bukan loop.** Jangan merancangnya sebagai proses yang menghitung
seluruh rantai approver sekaligus: bentuk itu **tidak ada di sistem lama** dan berperilaku berbeda
ketika rantai terputus di tengah.

**Tiga field state yang tidak boleh disatukan menjadi satu "status":**

| Field | Isinya |
| --- | --- |
| **antrean** | peran yang **sedang** memegang kasus (ruang nama tersendiri) |
| **kode jabatan tujuan** (`next_approver_position`) | jabatan **berikutnya** dalam tangga — ruang nama berbeda |
| **hasil keputusan** | keputusan underwriting terakhir |

⚠️ Kode jabatan tujuan disimpan di properti warisan bernama `LetterNo`, yang **tidak berisi nomor
surat**. Antrean dan kode jabatan diberi **tipe berbeda** supaya kompilator menangkap pertukarannya —
salah satu tempat sistem tipe menangkap cacat warisan.

⚠️ **Tangga selesai adalah penyelesaian normal, bukan galat.** Ketika tidak ada jabatan tujuan yang
cocok, wewenang sudah cukup dan kasus berhenti di situ.

Tabel limit **disuntikkan sebagai fixture**, dan fixture itu sekaligus **kontrak bentuk data** ke DBA.
Yang hilang dari tangga ini dulu adalah datanya, bukan logikanya — kini datanya sudah ada, dengan
ejaan jabatan **tanpa spasi** yang terbukti cocok dengan token routing di rule.

⛔ **Kolom nama dan login WAJIB dibuang** saat fixture dibangun dari berkas data — keduanya memuat
nama orang. Dibuang, bukan disamarkan.

**Blocked by:** 01

**Status:** ready-for-human — diport 01-10-2026; keputusan agent A20–A25 dikonfirmasi (butir 48)

- [ ] Satu keputusan manusia menghasilkan **tepat satu** transisi; tidak ada kasus uji yang mengharapkan rantai
- [ ] Ketiga field state terpisah dan bergerak independen
- [ ] Antrean dan kode jabatan bertipe **berbeda**; menukarnya **gagal saat kompilasi**
- [ ] Tidak ada jabatan tujuan yang cocok → **tangga selesai**, ditandai penyelesaian normal, bukan galat
- [ ] Urutan eskalasi mengikuti ambang limit menaik; baris pertama yang cocok = approver berikutnya
- [ ] Fixture memuat kolom yang **dibaca query** saja — **tanpa kolom nama dan login**
- [ ] Ejaan nilai jabatan **tanpa spasi**, sesuai data nyata
- [ ] Pemeriksaan kebocoran nama dijalankan atas fixture sebelum di-commit, memakai pencocokan **batas kata** — bukan substring
- [ ] ⚠️ Lima belas tautologi pembanding limit dan empat nomor polis literal yang ada di rule alur **direproduksi apa adanya** dan ditandai kandidat perbaikan — bukan dirapikan

## Comments

### 2026-10-01 — dipetakan, BELUM diport (agent)

Rule yang menjalankan satu langkah tangga sudah dipetakan (`GETLIMITAKSEPTASI_ACTFLOW`, query bentuk A,
pemetaan `LetterNo` → antrean di konektor `Decision23` flow Offer). Port **ditahan**, bukan dilewati,
karena tiga hal menentukan approver berikutnya dan korpus tidak menjelaskannya (`../PERTANYAAN-AKSEPTASI.md` B1–B4):

1. arti kode transisi **5** yang menggerbangi blok tangga (langkah 19) dan penimpa `LetterNo`
   (langkah 20);
2. ejaan `JABATAN`: blok tangga mencocokkan literal berspasi DAN tanpa spasi di activity yang sama;
3. 9 dari 14 query membaca `LOGIN` pengguna (guard identitas) — usul: limit pengguna jadi masukan.

`[terverifikasi]` dihitung dua cara: 15 tautologi limit di `_ActFlow` dan 15 di `_Act`; empat nomor polis
literal di langkah 17–18 (L4688, L5109). Antrean ditulis flow, bukan activity, dan tidak satu pun query
membaca `WORKBASKET` — koreksi atas spec Modul 4.

### 2026-10-01 — dua dari empat penahan terjawab (butir 31, 33)

- **B2 terjawab** — ejaan `JABATAN` memang campuran (berspasi dan tanpa spasi). Cabang berspasi
  (`"DIREKTUR TEKNIK"` L4635) dan tanpa spasi (`"KADIVFACULTATIVE"` L4929) keduanya hidup; dicocokkan
  persis per literal.
- **B3 terjawab** — limit pengguna menjadi **masukan** tangga dari model peran, bukan lewat `LOGIN`.
- ⏸ **B1 masih menahan** — arti kode transisi `5`. Lokasinya: `NB FacIn\Activity\GetLimitAkseptasi_ActFlow.xml`,
  **23 kemunculan** — 22 "bila benar" + 1 "bila salah" — di langkah 12, 13, 19 (×2), 19.1.2/.3/.5, 20,
  21.1.2–21.1.8, 22 (×4), 23 (×4). Dihitung tiga cara yang sepakat: `grep -cE` atas
  `<pyStepsPreCondParamsWhen(True|False)>5<`, hitung elemen lewat pengurai XML, dan penelusuran langkah
  bersarang. ⚠️ Catatan agent sebelumnya menyebut **24** dan "langkah 22.2.4" — keduanya keliru; cara
  hitung lama tidak tersimpan sehingga sebabnya tidak dapat direkonstruksi. Paling mudah dicek di
  langkah 12 (prakondisi 1 `IsLimitSBondKBG`, "bila benar" = 5, L3477).
- ⏸ **B4 masih menahan** — CSV tabel `M_LIMIT_*` dari DBA (butir 35).

### 2026-10-01 — B1 terjawab: kode 5 = Skip Whens (butir 39)

Dari UI Pega langkah 12: `2` = Continue Whens, `3` = Skip Step, `5` = Skip Whens. Arti tindakannya
`[dugaan kuat]` mengikuti platform Pega (Skip Whens = When sisanya dilewati, langkah dijalankan). Bila benar,
prakondisi langkah 19 tidak pernah melewati langkah. **NB-11 kini hanya ditahan B4** (CSV tabel limit).

### 2026-10-01 — B4 terjawab: CSV tabel limit diterima (butir 42)

Fixture `backend/services/acceptance/testdata/limit/` (penjaga `TestFixtureLimitTanpaNamaLogin`). Temuan:
`JABATAN` bentuk A seluruhnya tanpa spasi; langkah 19 dan 21 berlabel `//`. Bila `//` = nonaktif
(`PERTANYAAN-AKSEPTASI.md` H1), keputusan tangga bentuk A tinggal langkah 20. **Penahan tersisa: H1.**

### 2026-10-01 — H1 terjawab; tiket tidak lagi tertahan (butir 43)

Label `//` = di-remark. Langkah 15, 16, 19, 21 tidak diport; keputusan tangga bentuk A = langkah 20 atas
daftar langkah 6–10. Semua penahan (B1 kode transisi, B2 ejaan, B3 `LOGIN`, B4 CSV, H1 remark) terjawab —
**siap dikerjakan**, dengan fixture `backend/services/acceptance/testdata/limit/`.

### 2026-10-01 — diport (butir 44, 45)

Kode: `backend/services/acceptance/tangga.go` — `Next(kasus, tabel, pengguna) (Transisi, error)`.

| Kriteria | Keadaan |
| --- | --- |
| Satu keputusan → tepat satu transisi | ✅ `TestNextSatuLangkah` — tiap panggilan satu langkah; tidak ada rantai |
| Tiga field state terpisah | ✅ `Transisi` memisah `JabatanTujuan` dan `Antrean`; `Keputusan` (NB-13) tipe sendiri (butir 45) |
| Antrean dan jabatan bertipe berbeda, tukar gagal kompilasi | ✅ `TestJabatanAntreanGagalKompilasi` (+ kontrol yang terkompilasi) |
| Tidak ada tujuan → tangga selesai, bukan galat | ✅ `TestNextTanggaSelesai` |
| Eskalasi menaik; baris pertama = approver | ✅ SQL bentuk A direproduksi (`saring`); seri antar-jabatan → `ErrUrutanTakPasti` (butir 44 d) |
| Fixture hanya kolom query, tanpa nama/login | ✅ `testdata/limit/` + `TestFixtureLimitTanpaNamaLogin`; ⚠️ `WORKBASKET`/`JABATAN_ATASAN` ikut (README) |
| Ejaan jabatan tanpa spasi | ✅ data nyata bentuk A seluruhnya tanpa spasi (butir 42) |
| Pemeriksaan kebocoran nama batas kata | ⏸ tidak dapat dijalankan: `NAMA`/`LOGIN` tidak diekspor, jadi tidak ada nilai nama sebagai sumber pencocokan; penjaga diganti pemeriksaan header |
| 15 tautologi + 4 nomor polis literal direproduksi | ⛔ dicabut: tautologi ada di langkah 19/21 yang di-remark (butir 43); langkah 17 diabaikan (butir 45) |

Uji mutasi `tangga.go`: **21/21 tertangkap** (mutasi harus terkompilasi; tertangkap = `go test` keluar
non-nol). Putaran pertama 16/21 — empat celah (urutan Banding, syarat antrean langkah 11, nilai mutlak
endorsement, Preferred Commercial di Banding) ditutup tes `TestNextCelahMutasi`. Tinjauan dua sumbu
(`/code-review`) dijalankan; temuannya diterapkan, termasuk ralat antrean `ToKadivFacultative`.

### 2026-10-01 — langkah 17 dibuang (butir 47)

Sikap butir 45 ("diabaikan di NB-11") diganti butir 47: langkah 17 tidak dipakai lagi di sistem baru, untuk
kedua bentuk. Tidak ada perubahan kode bentuk A. `ErrBentukBBelumDiport` diganti `ErrBentukB` — kasus bentuk B
kini dijalankan `NextFinancial` (NB-12).
