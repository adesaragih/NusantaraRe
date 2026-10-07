# Rekonsiliasi AC — EDM Treaty In (06-10-2026)

Peta setiap acceptance criteria dua spec EDM ke tiket, sumber kewenangannya, dan hasil pencocokan dengan
**tabel NB nyata** dan **XML**. ⛔ **Tanpa status implementasi** — itu milik `HASIL-IMPLEMENTASI` asisten utama.

**Sumber:** `WO` = `[keputusan work owner]` · `XML` = `[terverifikasi]` dari korpus · `DBA` = `[data DBA]` ·
`PS` = `[penyimpangan sadar]` · `+NB` = bergantung pada tabel/keputusan NB nyata.
**Catatan:** `✓` = cocok dengan NB nyata dan XML (sejauh diperiksa) · `⛔ dikoreksi #n` / `⚠️ #n` = lihat baris
*n* di `KOREKSI-DOKUMEN-2026-10-06.md` · `butir WO` = menunggu keputusan, tidak ditutup di sini.

Cacah: `spec-penyimpanan-relasional.md` **58** nomor (AC 56 ditarik ⇒ **57** berlaku) · `spec.md` **59**.

## 1 · `spec-penyimpanan-relasional.md` — 57 AC berlaku

| AC | Ringkas | Tiket | Sumber | Catatan |
| ---: | --- | ---: | --- | --- |
| 1 | baris endorsemen `PRODKE ≥ 1` | 02 | WO +NB | ✓ `PRODKE NUMBER(10) NOT NULL DEFAULT 0` (320 baris 33) |
| 2 | `OLD_POLIS_ID` terisi, menunjuk generasi sebelumnya | 02 | WO +NB | ✓ FK ke `T_WORK_POLIS` (320 baris 113) |
| 3 | dua baris tak berbagi `OLD_POLIS_ID` | 02 | WO +NB | ✓ `UQ_GP_TREATY_OLD` (320 baris 114) |
| 4 | dua endorsemen serentak: yang kedua ditolak | 02 | XML +NB | ✓ indeks unik `(NOPOLIS, PRODKE)` hanya baris bernomor (320 baris 117) — baris EDM selalu bernomor |
| 5 | nomor `NOPOLIS + "/E" + dua digit` | 02 | WO | ⛔ dikoreksi #44, #25: **sekurangnya** dua digit; kolom `NOENDORS` |
| 6 | sunting generasi berpenerus ditolak | 02 | WO +NB | ✓ generasi tertutup = punya penerus (320 baris 19–20) |
| 7 | selisih terhadap generasi tepat sebelumnya | 03 | WO | ✓ |
| 8 | generasi kehilangan `NOURUT` ditolak di `services` | 03 | WO | ✓ |
| 9 | tidak ada tabel salinan nilai lama | 03 | XML | ✓; didamaikan dengan `spec.md` AC 7–9 (#13) |
| 10 | baris rincian tidak dapat dihapus | 04 | WO | ✓ |
| 11 | `NOURUT` lama terbawa apa adanya | 04 | WO | ✓ |
| 12 | baris baru `NOURUT = maks + 1` | 04 | WO | ✓ |
| 13 | `PASANGAN_BERGESER = 1` = anomali | 04 | XML | ✓ kolom 361/362/363 |
| 14 | pembatalan = generasi baru bernilai nol | 05 | XML | ⚠️ #47: XML menolkan 16 medan + angsuran + spreading, bukan seluruhnya — **butir WO** |
| 15 | jenis endorsemen di kolomnya sendiri, dipilih di awal | 05 | WO +NB | ✓ `EDM_TYPE` (320 baris 48; `T_POLIS_QUOTATION`); label kode 1–4 tidak ada di korpus |
| 16 | satu rumus `baris_ini − baris(OLD_POLIS_ID)` | 06 | WO | ⚠️ #40 varian kedua **berpemilih** (fakta dikoreksi; keputusan tetap); #41 cakupan NonProp — **butir WO** |
| 17 | medan uang = hasil pengurangan | 06 | XML | ✓ `EDMTCalculateTreatyDifference` 1–6 |
| 18 | persen disalin, tidak dikurangi | 06 | XML | ✓ |
| 19 | kunci disalin | 06 | XML | ✓ — `Currency`/`CurrencyID` **tidak** disalin di jalur proporsional (#39) |
| 20 | `DUE_TO` dari tanda selisih | 06 | XML | ⛔ dikoreksi #37: hanya lapisan XOL |
| 21 | `DUE_TO` XOL dari jumlah sepanjang lapisan | 06 | XML | ⛔ dikoreksi #38: per lapisan; induk basi, nol pembaca |
| 22 | NonProp: angka sama sebelum/sesudah penyeragaman | 06 | XML | ⚠️ #41 premis *"sudah seragam"* keliru — hanya terpenuhi bila rumus XML NonProp dipertahankan — **butir WO** |
| 23 | rumus selisih tidak di basis data | 07 | WO | ✓ |
| 24 | tabel selisih hanya ditulis aplikasi | 07 | WO | ✓ |
| 25 | ditulis dalam transaksi generasinya | 07 | WO | ✓ |
| 26 | bangun ulang hanya `SUMBER = 'GO'` | 07 | WO +NB | ✓ `CK_POLIS_DIFFERENCE_SUMBER` (360 baris 44) |
| 27 | isi beda dari hitung ulang ⇒ tabelnya salah | 07 | WO | ✓ |
| 28 | kunci penyaring `NOPOLIS PRODKE EDM_NO IDPEGA` | 07 | WO | ✓ 360; ⚠️ #16 klaim *"gaya TREATYINPRODUCTION"* untuk `EDM_NO` tidak tepat (di sana `NOENDORS`) |
| 29 | nol tabel selisih induk XOL | 07 | XML | ✓ kesimpulan; ⛔ alasan dikoreksi #48 (induk = turunan) |
| 30 | nol tabel selisih rincian angsuran | 07 | XML | ✓ diukur ulang 06-10: 12 rujukan skalar, 0 sarang |
| 31 | `repository` ambil dua baris, tanpa agregasi | 07 | WO | ✓ |
| 32 | 21 kolom khas endorsemen kosong di polis baru | 01 | XML | ⛔ #26 bertentangan NB nyata — disusun ulang (asisten utama) |
| 33 | empat kolom khas NB kosong di endorsemen | 01 | XML +NB | ⛔ #27: **tiga** kolom (P36) |
| 34 | keadaan layar tidak tersimpan | 01 | WO | ✓ |
| 35 | spreading EDM presisi 20, NB 10 | 08 | XML | ✓ `CountSpreading_Act` EDM langkah 5.1 |
| 36 | baris pertama spreading bawaan 100 | 08 | XML | ✓ `CountSpreading_Act` langkah 4 |
| 37 | rincian angsuran bertingkat tersimpan (NonProp) | 08 | XML | ✓ `FillPaymentInstallmentEDMT` 2.2 membangun `InstallmentList` |
| 38 | rincian dari salinan master tersimpan sama | 08 | XML | ⚠️ #46 pada `EDMT-` dibangun ulang dari selisih XOL — **butir WO** |
| 39 | nilai migrasi tidak dihitung ulang | 09 | WO | ✓ |
| 40 | pemasangan antar generasi lewat `NOURUT` | 09 | WO | ✓ XML memasangkan per `.pxListSubscript` (langkah 2.1, 3.1) |
| 41 | `PASANGAN_BERGESER` tanpa mengubah angka | 09 | XML | ✓ |
| 42 | `RUMUS_BERLAPIS = 1` bila generasi lama ber-`EDMNo` | 09 | XML | ✓ kini berbukti syarat dijalankan (#40); ⚠️ #17 induk 360 tanpa kolom |
| 43 | kedua penanda hanya `SUMBER = 'PEGA'` | 09 | WO | ✓ |
| 44 | pemuat migrasi lewat antarmuka yang sama | 10 | WO +NB | ✓; NB menyerahkan `PRODKE > 0` ke tiket 10 (#36) |
| 45 | satu transaksi per generasi | 11 | WO | ✓ |
| 46 | skema `POOLDATA.` eksplisit | 11 | WO | ✓ (SQL lama 16 / 31 — K Bab 5.2) |
| 47 | kolom pelaku dari identitas login | 11 | WO | ✓ |
| 48 | nol kolom uang `float` | 11 | WO | ✓ |
| 49 | skala desimal kolom uang | 11 | WO +NB | ⛔ #28: **sepuluh** (`NUMBER(38,10)`), bukan delapan |
| 50 | kode/penanda teks, nol di depan utuh | 11 | WO | ✓ |
| 51 | teks kosong → `NULL` | 11 | WO | ✓ |
| 52 | arah ketergantungan tidak dibalik | 11 | WO | ✓ |
| 53 | satu antarmuka `repository` | 11 | WO | ✓ |
| 54 | bentuk tabel dasar tidak berubah | 01 | XML +NB | ⛔ #23: tabel dasar = NB 320–328 (`T_POLIS_SURVEY` ikut) |
| 55 | Prop 10 tabel, NonProp 14 | 01 | XML +NB | ⚠️ #24, #53 cacah **belum pasti** — **butir WO** |
| 56 | ~~sebaran tambahan~~ | ~~08~~ | — | ⛔ **ditarik** 23-09 sore; diperkuat XML #54 |
| 57 | `REMARK` ≥ 128 utuh pulang-pergi | **12** | XML +NB | ⛔ #3 pindah dari 11; #22, #30 kolom `T_GENERAL_POLIS_TREATY.REMARK VARCHAR2(128)` ada |
| 58 | kelebihan skala dibulatkan, bukan ditolak | 11 | WO +NB | ⛔ #28: dibulatkan bila > **10** desimal |

## 2 · `spec.md` — 59 AC

⚠️ **Tidak satu pun tiket EDM menutup AC `spec.md`** — `issues/00-PETA-AC.md` hanya memetakan
`spec-penyimpanan-relasional.md`. Kolom *Tiket* karena itu `—`. Pemetaannya keputusan asisten utama / WO.

| AC | Ringkas | Tiket | Sumber | Catatan |
| ---: | --- | ---: | --- | --- |
| 1 | endorsemen = jenis kasus tersendiri | — | XML | ✓ kelas `…Work-EndorsementTreaty` (`Flow/InputAddendumTreatyIn.xml`) |
| 2 | tangga tiga jenjang | — | ~~XML~~ WO · PS | ⛔ #42 XML membolehkan Sec Head selesai sendiri; tiga jenjang = keputusan |
| 3 | tolak atasan ⇒ kembali ke admin | — | WO | ✓ `Transition8/21/4` → `Assignment2`; catatan: `Transition21` tidak memulihkan `PositionNote` (XML) |
| 4 | tolak admin ⇒ berkas tertutup ditolak | — | WO | ✓ `Decision3` No → `End3`; status akhir tidak ditulis rule mana pun di korpus |
| 5 | `IsApproved` dibandingkan sebagai teks | — | WO | ✓ `DecisionTable/isApproved` (`"0"` → No) |
| 6 | wewenang lewat peran, bukan nama | — | PS | ✓ XML memakai 2 ID operator, kolom telepon, ID operator di tombol (tidak disalin) |
| 7 | data lama disalin saat berkas dibuat | — | WO +NB | ⚠️ #13, #33 sistem baru: baris generasi `OLD_POLIS_ID` yang dibekukan, tanpa JSON |
| 8 | perubahan polis induk tak mengubah nilai lama | — | XML +NB | ✓ via generasi tertutup (#13) |
| 9 | snapshot bagian berkas, bukan rujukan tertunda | — | XML | ⚠️ #13 didamaikan dengan SP ID-8 / AC 9 |
| 10 | selisih terhadap keadaan tepat sebelumnya | — | WO | ✓ |
| 11 | pemilih cara hitung = isi `OldData.EDMNo` | — | XML | ✓ dikonfirmasi langkah 1 / 4 (#40) |
| 12 | satu polis boleh diendorse berkali-kali | — | WO | ✓ (tanpa batas — SP KEPUTUSAN butir 1) |
| 13 | nomor = induk + `/E` + dua digit | — | WO | ⛔ #44 **sekurangnya** dua digit |
| 14 | nomor urut +1 dari terakhir | — | XML | ✓; ⚠️ #56 `PRODKE` `JSON_POLIS` dihitung berbeda (`COUNT`) |
| 15 | < 10 dibubuhi nol | — | XML | ✓; ≥ 100 tidak dipotong (#44) |
| 16 | adendum & premi tambahan satu jalur nomor | — | WO | ⚠️ #45 jalur kedua berpenjaga `EDMNo==""`; `[terbuka]` 9.2 #3 tetap |
| 17 | nomor urut dalam transaksi penyimpanan | — | XML | ✓ |
| 18 | pembatalan = jenis endorsemen, dipilih di awal | — | WO | ✓ `EDMType` 4 |
| 19 | pembatalan lewat jalur yang sama | — | XML | ✓; panggilan efektif `EDMChooseBusiness_Act` 5 (#47) |
| 20 | 1 baris & % kosong/0 ⇒ bawaan 100 | — | XML | ✓ `CountSpreading_Act` 4 |
| 21 | `SplitRNMSharePct / RNMShare` presisi 20 | — | XML | ✓ hanya bila `SharePercentage` kosong (5.1); maksud dagang P60 tetap `[terbuka]` |
| 22 | ketidakseragaman presisi ditiru | — | WO | ✓ |
| 23 | spreading dihitung ulang saat premi berubah | — | XML | ✓ `CountOGPONP_Act` 9.1 → `CountSpreading_Act` |
| 24 | uang tidak pernah `float` | — | WO | ✓ |
| 25 | `/1,022` bila `Inclusive` persis | — | WO | ✓ |
| 26 | PPH 2 % · PPN 2,2 % dari potongan sesudah pembagian | — | XML | ✓ |
| 27 | empat medan potongan/bagian = persentase | — | WO +NB | ⛔ #31 `DEDUCTION1/2` = **uang** (K3) |
| 28 | presisi penuh; bulat hanya di penyajian | — | XML +NB | ⛔ #14 bulat saat dimuat bila > 10 desimal |
| 29 | uang ke pencapaian sebagai desimal tetap | — | DBA | ✓ (NB menulis `ACHIEVEMENT` langsung, tanpa prosedur) |
| 30 | jenis treaty dari medan penentu sebelum urai | — | XML | ✓ — berlaku bagi pemuat dokumen lama (tiket 10 / NB) |
| 31 | `ListInstallment` bersarang terbaca | — | XML | ✓ (pemuat) |
| 32 | `ListInstallment` datar terbaca | — | XML | ✓ (pemuat) |
| 33 | kedua bentuk wajib diuji | — | XML | ✓ |
| 34 | medan hilang = kosong, bukan galat | — | DBA | ✓ |
| 35 | tiga bentuk contoh produksi terbaca | — | DBA | ✓ |
| 36 | FacOut tidak dibangun | — | WO | ⚠️ #49 XML dapat mengirim kedua bila `IsFacRetro` — **butir WO** |
| 37 | sukses menguji FacIn dan FacOut | — | XML | ✓ `When/IsSuccessHitService` |
| 38 | gagal konversi ⇒ hapus data produksi | — | WO | ⚠️ #50 XML juga bersyarat `IsPEGAPROD` — **butir WO** |
| 39 | tidak menghapus bila `STS_KONVERSI = 1` | — | DBA | ✓ `serviceInsertArasapas_act` 12 |
| 40 | penghapusan lewat satu fungsi `repository` | — | WO | ✓ |
| 41 | pengguna diberi tahu sebelum hapus | — | WO | ✓ |
| 42 | pesan galat sampai ke pengguna | — | WO | ⚠️ #50 langkah 8 / 10 XML dilewati untuk `IsTreatyIn`; AC memegang niat P52 |
| 43 | penutupan paksa tidak dibangun | — | WO | ✓ langkah 18 ber-`//` |
| 44 | tanpa autentikasi ke produksi | — | WO | ✓ `ConnectREST` `pyUseAuthentication = false` |
| 45 | alamat tujuan dari variabel lingkungan | — | PS | ✓ (XML: `M_LINK_SERVICE` → `LinkService`) |
| 46 | uji dan produksi beda basis data | — | WO | ✓ |
| 47 | satu transaksi | — | WO | ✓ |
| 48 | skema eksplisit | — | WO | ✓ |
| 49 | tanggal tutup buku satu baris global | — | DBA | ⚠️ #52 naskah tanpa `ROWNUM`; baris pertama `pxResults(1)` |
| 50 | pelaku riwayat dari identitas login | — | WO | ✓ |
| 51 | tanpa pencarian nomor urut antrean | — | WO | ✓ |
| 52 | pemantauan idempoten | — | DBA | ✓ |
| 53 | pemantauan tidak memperbarui | — | DBA | ✓ |
| 54 | anti-duplikat pencapaian tanpa nilai uang | — | PS | ✓; kunci pengenal tetap `[terbuka]` 9.2 #2 |
| 55 | pesan galat penyimpanan tanpa markah | — | PS | ✓ |
| 56 | penanda lingkungan tak mengubah pesan | — | PS | ✓ |
| 57 | ke-380 langkah `Property-Set` dibangun | — | XML | ✓ diukur ulang 06-10: 380 dua cara (korpus berubah, angka tetap) |
| 58 | rumus dari isi langkah, bukan tebakan | — | XML | ✓ |
| 59 | penelusuran sampai 6 tingkat sarang | — | XML | ✓ kedalaman 165·83·65·29·20·18 tetap |

## 3 · Ringkasan

| | SP | S |
| --- | ---: | ---: |
| AC berlaku | 57 | 59 |
| ✓ cocok *(boleh bercatatan kecil)* | 43 | 48 |
| ⛔ dikoreksi | 9 | 4 |
| ⚠️ catatan / butir WO | 5 | 7 |
| berpindah tiket | 1 (AC 57 → 12) | — |
| tanpa tiket EDM | 0 | **59** |

Cacah dihitung dari tanda pertama kolom *Catatan* (pola teks atas berkas ini; AC 56 SP tidak dihitung).
