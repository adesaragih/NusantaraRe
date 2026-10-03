# 01: Registrasi klaim dan validasi nomor polis treaty

**Status:** ready-for-agent
**Blocked by:** 00 (PREFACTOR)
**Menutup:** AC 102 · 103 · 104 · 108 · 109 · 110 · 111 · 112 *(8 AC)* — US 1–4

## Hasil & nilai pengguna

Claim Admin dapat mendaftarkan klaim baru dengan nomor polis treaty, dan sistem menolak sejak awal
nomor polis yang kosong, salah format, atau tidak ada di master treaty. Klaim ganda pada Date of Loss
yang sama **benar-benar tertahan**, bukan sekadar diperingati — sementara klaim pengganti atas klaim
yang sudah ditolak tetap dapat dibuat.

## Area codebase

Endpoint registrasi klaim · validasi masuk · lookup master polis · layar registrasi.

## Rule Pega sumber

| Rule | Step | Yang dibuktikannya |
| --- | --- | --- |
| `Activity/CheckNoPolicy.xml` | 2 | pesan format salah; ⚠️ arah gerbang **terbalik** — pesan muncul saat nomor polis **tidak** mengandung penanda format |
| `Activity/CheckNoPolicy.xml` | 4 | menarik data polis lewat `RDBList/GetDataNopolisTreatyin.xml` |
| `RDBList/GetDataNopolisTreatyin.xml` | — | membaca tabel `policyjson`; ⚠️ tanpa prefiks schema |
| `RDBList/SetPolicyTreatyProp.xml` | — | menarik nomor polis, nomor offer, kelompok treaty, sumber bisnis, tahun UW, kuartal; ⚠️ alias berbohong + literal dipalsukan jadi kolom + tabel tanpa prefiks |
| `Activity/SetMasterID.xml` | 3 | pemanggil `SetPolicyTreatyProp` |
| `RDBList/CekPolicyNumber_SQL.xml` | — | mencari id di `pooldata.TREATYINPRODUCTION` menurut nomor polis |
| `Activity/CheeckNoRNM_Act.xml` | 4 | ⚠️ arah **terbalik** — pencarian jalan saat nomor polis **tidak** kosong |
| `Activity/CheeckNoRNM_Act.xml` | 5 | menolak saat hasil pencarian **nol baris** |
| `Activity/CheckNopolicy_Act.xml` | 3 · 4 · 5 | nomor polis kosong ditolak; dua step berikutnya menuliskan syarat yang sama dalam dua bentuk berlawanan |
| `Activity/CheckDateDOL_Act.xml` | 20 | ⚠️ loop riwayat klaim; flag prakondisi **`false`** — pembebasan satu nomor polis yang di-hardcode **mati** |
| `Activity/CheckDateDOL_Act.xml` | 20.1 | menandai duplikat saat DOL sama dan nomor klaim berbeda; menyimpan kunci klaim lama |
| `Activity/CheckReportDate_Act.xml` | 11 · 14 | **BARU 2026-09-19 (ronde 6)** — ⚠️ pemeriksa Report Date terhadap **Start/End Date Policy** keduanya **ber-remark**, dan catatan pengembangnya berbunyi *"matikan protek dalam periode polis"*. **Proteksi ini dimatikan dengan sengaja** |
| `Activity/CheckDateReceived_Act.xml` | 10 · 13 | **BARU** — pola yang sama untuk **Date Received**: kedua pemeriksa terhadap periode polis **ber-remark** |
| `Activity/CheckPeriodPolicy_Act.xml` | 1 | **BARU** — ⚠️ langkah 1 ber-remark padahal ia **satu-satunya pengisi** `Local.CompareDate` dan `Local.EqualDate`; langkah 2 dan 4 tetap membacanya. `[terbuka]` butir 7 spec |
| `Activity/SetEndDate_Act.xml` | 8 | **BARU** — pemanggilan `CheckPeriodPolicy_Act` dari sini **ber-remark**; rule itu tetap dirujuk dari dua Section |
| `Activity/AddAdjustment_Act.xml` · `Section/InputAcceptation.xml` · `Section/OutstandingClaim.xml` | — | **BARU** — ⚠️ properti `.ClaimData.PeriodPolicyTBA` dan pesan *"Policy period is TBA. Please verify the dates."*. **Konsep periode polis "belum pasti" belum pernah tercatat.** `[terbuka]` butir 8 spec |
| `Activity/MakeLowercase_Act.xml` | 1 | **BARU** — lokasi kerugian dan uraian laporan **ditimpa versi huruf kecil**; nilai asli tidak disimpan |

> **Catatan lingkup tambalan ronde 6.** Tambahan di atas **tidak** menyangkut `CheckNoPolicy`.
> Butir *"pencarian polis rangkap"* yang sempat digantung ronde 4 sudah **GUGUR**
> `[data work owner 2026-09-19]`: langkah 3 dan 4 rule itu ber-remark, jadi pencariannya **tidak
> rangkap**. Baris `CheckNoPolicy` di tabel ini **tidak disentuh**.

> **Acceptance criteria tiket ini TIDAK berubah.** Keenam baris di atas adalah **bukti tambahan**
> dan dua butir `[terbuka]` ringan yang **tidak memblokir** — perilaku yang ditiru adalah perilaku
> yang berjalan.
| `Activity/CheckDateDOL_Act.xml` | 21 · 22 | mengambil data klaim lama lewat `ReportDefinition/RejectedClaim_RD` |
| `Activity/CheckDateDOL_Act.xml` | 23.1 | **pengecualian** — klaim lama ada di tabel reject → penanda dikembalikan ke nol |
| `Activity/CheckDateDOL_Act.xml` | 24 · 25 · 26 | peringatan tampil; penanda kesalahan diset `1` — ⚠️ **tidak pernah lolos ambang `>1`** |

## ADR terkait

**ADR-0003** (uang non-float, lewat data polis).

## Acceptance criteria

- [ ] `[terverifikasi]` Nomor polis treaty wajib berformat yang dikenali; ⚠️ arah gerbangnya terbalik di Pega — aturan sebenarnya **nomor polis harus mengandung penanda format treaty** *(AC 108 spec)*
- [ ] `[terverifikasi]` Data polis ditarik dari master saat klaim didaftarkan — nomor polis, nomor offer, kelompok treaty, sumber bisnis, tahun underwriting, kuartal *(AC 109 spec)*
- [ ] ⚠️ Nama kolom dibuat sesuai isinya dan schema selalu eksplisit. **Alasan menyimpang:** `RDBList/SetPolicyTreatyProp.xml` menamai kelompok treaty sebagai *nama bisnis*, menyisipkan literal sebagai kolom, dan mengeja tabelnya **tanpa prefiks** padahal `RDBList/CekPolicyNumber_SQL.xml` mengeja tabel yang sama **dengan** prefiks *(AC 110 spec)*
- [ ] ⚠️ `[terverifikasi]` Nomor polis yang **tidak ada di master treaty** ditolak dengan pesan yang menyebut sebabnya; penolakan dipicu saat hasil pencarian **nol baris**. **Jebakan:** arah gerbang pencariannya **terbalik** di Pega — pencarian berjalan justru saat nomor polis **tidak** kosong; kedua gerbangnya flag `true`, jadi berlaku *(AC 111 spec)*
- [ ] `[terverifikasi]` Nomor polis **kosong** ditolak terpisah dari nomor polis tidak dikenal; langkah lanjutan hanya berjalan saat nomor polis tidak kosong *(AC 112 spec)*
- [ ] ⚠️ Klaim dengan Date of Loss sama pada polis yang sama **MEMBLOKIR penyimpanan** dan menampilkan peringatan. **Alasan menyimpang:** niat memblokir tertulis di keterangan step 25 (*"make disable submit if there is similar date"*), tetapi step itu memasang nilai `1` sedangkan kedua pembacanya menguji `> 1` — peringatan tampil dan penyimpanan tetap lolos *(AC 102 spec)*
- [ ] `[terverifikasi]` Pengecualian klaim yang sudah ditolak tetap berlaku — tidak ada blokir dan tidak ada peringatan; alur klaim pengganti tidak terganggu *(AC 103 spec)*
- [ ] ⚠️ Pembebasan satu nomor polis yang di-hardcode **tidak dimigrasikan**. **Alasan menyimpang:** literal identitas dunia nyata di dalam kode, sejenis dengan empat nama orang yang dibuang di AC 55; flag prakondisinya `false` sehingga pembebasan itu sudah mati hari ini *(AC 104 spec)*

## Catatan `[terbuka]` ringan — TIDAK memblokir

⚠️ `[terbuka]` Pembebasan satu nomor polis di `Activity/CheckDateDOL_Act.xml` step **20** adalah
**bekas pengecualian bisnis yang sengaja dipasang lalu dimatikan**, bukan sisa saringan uji coba —
karena itu layak ditanyakan kembali. Masih dikehendaki? Bila ya, jadikan baris data seperti keputusan
AC 55. Pemilik: **work owner**. Bawa apa adanya; **default yang ditulis adalah tidak dimigrasikan**.

## Perintah verifikasi

```
jalankan test "nomor polis kosong -> ditolak"
jalankan test "nomor polis format salah -> ditolak"
jalankan test "nomor polis tidak ada di master -> ditolak dengan sebab"
jalankan test "DOL sama pada polis sama -> penyimpanan DITOLAK"
jalankan test "DOL sama tetapi klaim lama sudah ditolak -> penyimpanan LOLOS, tanpa peringatan"
cari literal nomor polis di kode                          -> nihil
```
