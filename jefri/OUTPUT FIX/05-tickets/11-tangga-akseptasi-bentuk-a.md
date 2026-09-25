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

**Status:** ready-for-agent

- [ ] Satu keputusan manusia menghasilkan **tepat satu** transisi; tidak ada kasus uji yang mengharapkan rantai
- [ ] Ketiga field state terpisah dan bergerak independen
- [ ] Antrean dan kode jabatan bertipe **berbeda**; menukarnya **gagal saat kompilasi**
- [ ] Tidak ada jabatan tujuan yang cocok → **tangga selesai**, ditandai penyelesaian normal, bukan galat
- [ ] Urutan eskalasi mengikuti ambang limit menaik; baris pertama yang cocok = approver berikutnya
- [ ] Fixture memuat kolom yang **dibaca query** saja — **tanpa kolom nama dan login**
- [ ] Ejaan nilai jabatan **tanpa spasi**, sesuai data nyata
- [ ] Pemeriksaan kebocoran nama dijalankan atas fixture sebelum di-commit, memakai pencocokan **batas kata** — bukan substring
- [ ] ⚠️ Lima belas tautologi pembanding limit dan empat nomor polis literal yang ada di rule alur **direproduksi apa adanya** dan ditandai kandidat perbaikan — bukan dirapikan
