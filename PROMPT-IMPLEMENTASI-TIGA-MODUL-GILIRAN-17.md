# PROMPT — GILIRAN 17 *(sesi baru, folder `OUTPUT_HASIL_RNM`, cabang `main` @ `6c98e89` atau lebih baru)*: **LEMBAR KEPUTUSAN TERAKHIR — 31 OQ tiga modul diterapkan sekaligus, sisanya jadi daftar serah terima**

> Hanya konteks **Claim Life, PremiumList Life, Komite Claim Life**. Keputusan work owner 29-09-2026: *"rekomendasi"* — seluruh rekomendasi
> asisten atas 31 OQ diterima. GILIRAN-3 *(A/B/C)*, 4–16 tetap rujukan; aturan berhenti GILIRAN-6; **setiap pembacaan activity mencetak
> `pyStepsBlockName`**; commit dengan jalur eksplisit *(`git commit --only`)*. **Jangan** menyentuh berkas `tco_*` maupun suntingan frontend
> sesi lain yang belum di-commit di `main` *(`index.html`, `App.tsx`, `PagarGalat.tsx`, `labels.ts`, `styles.css`, `Shell.tsx`, `dasar.tsx`,
> `hooks/useTema.ts`, `lib/tema*.ts`)* — bila paket menuntut salah satunya, **berhenti pada paket itu dan laporkan**.

## 0. KEPUTUSAN — dari lembar keputusan 29-09-2026

### Claim Life

| OQ | Keputusan | Kerjakan |
| --- | --- | --- |
| **M1** | tanggal klaim **terkunci sesudah Save to RNM pertama berhasil** | gerbang keempat isian `EditDateClaimLife_Section` *(b1000, b1313, b1550, b1788)*: tahap Outstanding + Admin **dan** belum pernah Save to RNM; penanda "sudah Save to RNM" diturunkan dari kolom yang ada *(nomor klaim atau status baris yang ditulis rute Save to RNM — pilih yang pasti, buktikan)*; uji |
| **M2** | tanya DBA | kode tidak berubah → daftar serah terima |
| **M3** | **tunda** sampai OQ-N11 terjawab | → daftar serah terima *(bergantung `CLAIM_GROSS`)* |
| **M4** | klaim XOL **tidak** membawa peserta dari berkas | Upload CSV XOL tetap tidak dibangun; tercatat sebagai keputusan, OQ ditutup |
| **M5** | alasan penolakan disimpan di **kolom komentar baru jejak klaim** | migrasi **`021_kolom_komentar_jejak.sql`** *(+ down)*: `T_CLAIMLF_JEJAK.KOMENTAR` *(tipe dan lebar dari `KomiteComment`/`Remarks` b1687 section `RejectOSClaimLife_Sec`)*; dialog Reject Outstanding dibangun *(Date b790, PIC b975, Remarks b1687, Submit b3117, VERBATIM)*; rute tolak yang ada menerima dan menyimpan alasannya; uji |
| **M6** | mencabut peserta = **penanda**, layar menyembunyikan | periksa kolom penanda yang ada di `T_CLAIMLF_PREMIUMLIST_DETAIL`; bila tidak ada → migrasi **`022`** satu kolom penanda; tombol `DELETE` `InputOSClaimLife` b17865 menandai, tidak menghapus; semua pembaca menyaring penanda; uji |
| **M7** | izin baca **`RATE_LIFE`** sempit seperti butir bh; sumber `OUTWARDRATEID` ke DBA | satu pembaca baca-saja `RATE_LIFE` berkolom tetap, penjaga master dipersempit seperti bh; `OUTWARDRATEID` tetap `[terbuka — DBA]`, Spreading tetap tidak dipanggil sampai sumbernya ada → daftar serah terima |
| **N1** | bendera simpan = **keadaan turunan**, tanpa kolom | OQ ditutup; tidak ada migrasi |
| **N2** | cermin `OS_AKSEPTASI_KLAIM_LIFE` **mengisi** `NAME_OF_INSURED`, `DOB`, `CEDINGCO` seperti Pega | ditulis **di dalam SQL** dari baris sumber `M_LIFE_PREMIUM_DETAIL` *(`INSERT … SELECT` / `UPDATE … SELECT`)*, sehingga nama dan tanggal lahir **tidak melintasi Go**, log, uji, atau dokumen; klaim ganda antarklaim baru kini tertangkap; uji `db` *(SKIP)* + penjaga statik "nol kolom nama di model Go" tetap hijau |
| **N3** | gerbang retro langkah 27 **dipertahankan** *(XML Claim Life hidup)*; cermin header tetap di transaksi simpan | OQ ditutup, kode tidak berubah |
| **N4** | pemberitahuan | → daftar serah terima pemilik ekspor |
| **N5** | keputusan Komite "cutover 7 Feb 2025 tidak dipakai" **berlaku juga** di Claim Life | penukaran retro `InsertJsonClaimLife_Act` langkah 2 cukup **dua** WHEN *(Type TP/TR b1268, security terisi b1291)*; syarat `ProdDateTime` b1314 tidak dipakai; `ErrGerbangRetroTakTerputuskan` dibuang bila tidak lagi mungkin terjadi; uji |

### PremiumList Life

| OQ | Keputusan | Kerjakan |
| --- | --- | --- |
| **PL-09** | tulis `M_LIFE_PREMIUM_SUMMARY` sesuai **pl2** | dalam transaksi simpan summary yang sama, kolom VERBATIM dari RDB/prosedur yang ditiru |
| **PL-10** | kolom uang kosong di `M_LIFE_PREMIUM_DETAIL` warisan = **0** seperti Pega | satu fungsi konversi di tepi repository warisan; ADR-U-0027 tetap untuk tabel `T_*`; uji |
| **PL-11** | tetap stub; tanya pemilik Arasapas | → daftar serah terima |
| **PL-12** | unggah CSV: kolom uang kosong = **0** seperti `ValidasiUploadPL_act` langkah 2 | uji |
| **PL-13** | ikut XML: angka **25** tertanam | konstanta bertanda bukti baris; OQ ditutup |
| **PL-14** | `convertJsonNusareToProduction` **ditiru** lewat outbox, pelaksana **stub** *(berkaitan PL-11)* | efek outbox baru, tanpa panggilan nyata |

### Komite Claim Life

| OQ | Keputusan | Kerjakan |
| --- | --- | --- |
| **K-04a** | DBA memastikan rentang nomor; uji keunikan tetap | → daftar serah terima |
| **K-05** | **tiru** langkah 5.1 *(hidup — `pyStepsBlockName` kosong b5795)*: tingkat akhir menolak → keputusan seluruh tingkat ditimpa `2`, komentar kosong, `DateApprove` sekarang | satu `UPDATE` bersyarat dalam transaksi keputusan; **sebelum** menimpa, keputusan lama tiap tingkat dicatat di jejak *(ADR-0007)* sehingga riwayat tangga tiket 09 tetap dapat dibaca; layar riwayat membaca jejak untuk tingkat yang tertimpa; uji |
| **K-05b** | ikut XML: Tolak di tingkat tengah menghentikan tangga tanpa menyentuh baris klaim | perilaku sekarang dipertahankan; akibatnya *(baris tetap ber-`KOMITE_ID`, tidak dapat diserahkan ulang)* dicatat jelas di tiket 05 dan panduan uji |
| **K-06** | Kasir tetap **stub**; sambungan nyata = fitur baru | ADR-0015 diberi catatan bertanggal; OQ ditutup |

## 1. DAFTAR SERAH TERIMA

Buat **`DAFTAR-SERAH-TERIMA-TIGA-MODUL.md`** di akar `OUTPUT_HASIL_RNM`, dikelompokkan per penerima, tiap butir satu pertanyaan + bukti baris +
akibat bila tidak dijawab:

| Penerima | Butir |
| --- | --- |
| **DBA** | M2, M7 *(`OUTWARDRATEID`)*, K-04a, PL-17 *(`PC_DATA_UNIQUEID` NBLF-)*, OQ-001, OQ-002, OQ-013, OQ-018, OQ-047, skema uji kosong *(G1)* |
| **Pemilik ekspor Pega** | N4, N11 *(`CLAIM_GROSS`)* |
| **Pemilik Arasapas** | PL-11 |
| **Product + UW / Finance** | OQ-032, OQ-037, OQ-060 |
| **Menunggu jawaban lain** | M3 *(sesudah N11)* |

## 2. URUTAN — satu commit per paket

| # | Paket | Commit |
| ---: | --- | --- |
| 1 | OQ yang cukup ditutup di dokumen: M4, N1, N3, PL-13, K-06 | `docs: lembar keputusan — OQ yang ditutup tanpa kode` |
| 2 | Claim Life M1, M5 *(021)*, M6 *(022 bila perlu)* | `claim-life: M1/M5/M6 — kunci tanggal, alasan tolak, cabut peserta` |
| 3 | Claim Life N2, N5, M7 *(pembaca `RATE_LIFE`)* | `claim-life: N2/N5/M7 — cermin warisan lengkap, tukar retro dua syarat, RATE_LIFE` |
| 4 | PremiumList PL-09, PL-10, PL-12, PL-14 | `premiumlist-life: PL-09/10/12/14` |
| 5 | Komite K-05, K-05b | `komite-claim-life: K-05/K-05b` |
| 6 | daftar serah terima + status tiket + panduan uji + OQ ditutup | `docs: daftar serah terima tiga modul` |

Migrasi `021`/`022` dijalankan **work owner** sesudah giliran selesai. Tiap paket: ralat bertanggal di tiket terkait; uji murni + handler +
`db` *(SKIP)*; angka uji dengan dan tanpa tag `db`.

## 3. LAPORAN

Satu pesan: tabel **OQ → commit → rule XML → kode**; OQ ditutup; butir serah terima; angka uji tiap commit; bab **TELEMETRI EKSEKUSI**.

---

*Disusun 29 September 2026 dari lembar keputusan 31 OQ (jawaban work owner "rekomendasi"), dengan PL-11 dibaca dari tiket 06 PremiumList dan K-05
dari tiket 05 Komite serta `KomitePostAdjustment.xml` langkah 5.1 (`pyStepsBlockName` kosong b5795).*
