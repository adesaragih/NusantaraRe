# Pengetahuan Penggrill — Migrasi Treaty In Adjustment

> **Isi berkas:** seluruh pengetahuan yang terkumpul di percakapan penggrill (Claude di claude.ai) selama sesi grilling modul **Treaty In Adjustment**.
> **Disusun:** 24 September 2026, menurut tanggal di sisi penggrill. Berkas milik AI grilling bertanggal 25–26 September karena jam mesinnya berbeda.
> **Sifat:** catatan kerja untuk melanjutkan sesi. **Ini bukan sumber kebenaran.** Sumber kebenaran tetap ekspor XML, `PENGETAHUAN.md`, log dan folder `GRILL-*` milik AI grilling, serta artefak induk. Fakta sistem lama di sini adalah laporan AI grilling yang sudah diperiksa sejauh percakapan ini. Bila berbeda dengan berkas proyek, **berkas proyek yang menang**.

---

## 0. Cara memakai berkas ini

- **Melanjutkan sebagai penggrill di sesi baru:** baca §1–§3 (konteks), §6 (keputusan terkunci), §7–§11 (yang masih terbuka), dan §14 (daftar periksa meninjau jawaban AI).
- **Menyerahkan ke orang lain:** §6 dan §15 cukup untuk tahu posisi terakhir.
- **Butir bertanda *(perlu dicek)*** adalah hal yang tidak dapat saya pastikan dari percakapan. Cocokkan dengan indeks `KEPUTUSAN-GRILLING-ADJUSTMENT.md`.

---

## 1. Proyek, fase, dan peran

**Migrasi:** aplikasi Pega (GISFW, Pega 8.8) → frontend **React JS**, backend **Golang**, basis data **Oracle**.

**Modul:** Treaty In Adjustment — pembuatan addendum atau revisi atas kontrak Treaty In, lengkap dengan mesin selisih nilai lama-baru dan persetujuan berjenjang.

**Tiga fase berurutan:**

1. **grilling** — fase sekarang;
2. to-spec;
3. to-ticket.

Fase grilling tidak boleh menulis spesifikasi, DDL, struct Go, komponen React, endpoint, maupun tiket. Detail rancangan yang muncul dicatat sebagai **bahan to-spec**. AI grilling tidak pindah ke to-spec sendiri: di akhir ia menulis serah terima, "Grilling selesai, siap lanjut ke to-spec bila Anda setuju", lalu berhenti.

**Peran (METODE §1.3):**

| Pihak | Peran | Tempat |
|---|---|---|
| AI grilling (Claude Code, akses D:\) | **pembedah** — membaca ekspor, mengumpulkan bukti, menjawab, mengusulkan | direktori proyek |
| Penggrill (Claude di claude.ai) | **penggrill** — menolak jawaban yang belum berbukti, mencari yang tidak ditanyakan, menjaga batas | percakapan ini |
| Pengguna | pemilik sesi, meneruskan pesan antara keduanya, mengekspor data, menghubungi bisnis | — |

**Pola kerja:** pengguna menempel jawaban AI grilling ke percakapan penggrill. Penggrill meninjaunya, lalu menyusun **prompt balasan** dalam blok markdown untuk ditempel balik.

**Prompt pembuka** yang dipakai: `PROMPT-GRILLING-ADJUSTMENT.md` (skill `grill-me` dari mattpocock/skills, pohon keputusan A–K, aturan main, format log, kriteria selesai).

---

## 2. Lokasi berkas

Sebagian lokasi di bawah berasal dari laporan AI grilling dan belum pernah saya lihat langsung.

| Lokasi | Isi |
|---|---|
| `D:\XML_NURE\Treaty In Adjustment\` | ekspor XML Pega modul Adjustment (379 berkas) |
| `D:\XML_NURE\_migration-docs\METODE-GRILLING.md` | metode grilling. Disalin dari Downloads; kini di akar karena berlaku untuk semua modul |
| `D:\XML_NURE\_migration-docs\treaty-in\` | **modul induk**: `docs/adr/` (ADR-0034 … 0055), `SPEC-MODEL-DATA.md` (bentuk awal, langkah 8–10 belum dijalankan), `SPEC-INVARIAN.md` (INV-04, INV-25, …), `4-erd-dan-tabel-datar/STRUKTUR-DATA.md`, `ERD.md`, `CONTEXT.md`, `DAFTAR-ESKALASI-MANAJEMEN.md`, `BENTUK-PENULIS-PROPERTI.md`, `SAPUAN-DAN-NAMA-TAGNYA.md`, `PETA-TELUSUR-JSON.md` ⚠, `SEAM-ADJUSTMENT.md` ⚠, `KEPUTUSAN-SAMBUNGAN-ADJUSTMENT.md` ⚠, `TEMUAN-ADJUSTMENT-DITUNDA.md` ⚠, `tools/langkah-hidup.py` |
| `D:\XML_NURE\_migration-docs\treaty-in-adjustment\` | **modul ini**: `PENGETAHUAN.md` (AS-IS), `KEPUTUSAN-GRILLING-ADJUSTMENT.md` (kini **indeks** lintas ronde), `GRILL-A/`, `GRILL-B/`, `USULAN-REVISI-ADR.md`, `PEMILAHAN-SISA-GRILLING.md`, `PETA-SUMBER-INDUK.md` (50 berkas induk), `CABANG-K-PEMETAAN-TO-SPEC.md`, `ekspor-tambahan/` (EXP-1), `tools/tulis.py` (penyapu penulis properti terkalibrasi) |
| `D:\XML_NURE\_migration-docs\komite-claim-non-prop\GRILL-05\` | **contoh format** satu folder ronde (lihat §3.3) |

⚠ = ditulis sebelum koreksi `PENGETAHUAN.md`. Faktanya dikutip lalu diperiksa, tidak diwarisi.

**Catatan praktis:** berkas yang diunduh pengguna sering tertinggal di `C:\Users\Administrator\Downloads\`. Ini terjadi pada METODE dan pada EXP-1. Ekspor berikutnya langsung ditaruh di `ekspor-tambahan/`.

---

## 3. Cara kerja sesi yang sudah disepakati

### 3.1 Aturan main grilling

- **Satu pertanyaan per giliran.** Setiap pertanyaan memuat konteks → bukti (§, TDA, UA, nama aturan) → pilihan bernomor → rekomendasi beserta alasannya → arah dampak bila salah.
- **Jangan menanyakan yang bisa dibaca.** Jangan menebak yang tidak terbaca; tandai *tidak terbaca*.
- **Pisahkan "masih bisa terjadi"** (terbaca dari aturan; boleh menjadi dasar rancangan) dari **"pernah terjadi"** (hanya data yang menjawab; tidak memblokir rancangan; menjadi UA).
- **Pertanyaan ke bisnis berbentuk daftar untuk dibantah** (DB): pernyataan konkret, dalam bahasa bisnis, tanpa nama tabel atau aturan.
- **Artefak induk tidak disunting diam-diam.** Revisi ADR ditulis sebagai draf REV. Perubahan pada artefak induk lain ditunjukkan dulu sebagai diff dan menunggu izin.
- **Setiap keputusan membawa label** PELESTARIAN, PERUBAHAN (wajib menyebut berubah dari apa), atau BARU. Pernyataan cakupan dokumen diberi keterangan "bukan perilaku — tidak diuji".
- **Nasib TDA** memakai kosakata: diperbaiki, dilestarikan, atau ditunda.
- **Keluaran ditulis ke disk begitu langkahnya selesai.**

### 3.2 Aturan bukti yang lahir di sesi ini

Pelengkap METODE. Sebagian sudah diusulkan balik sebagai TA atau IND.

1. **Salinan bukan bukti.** Isi di dalam `pyIncludedRuleXML` atau `pyRuleVersionsList` tidak dihitung dan tidak dipakai sebagai bukti. Bungkusan itu bisa memuat versi lama atau salinan seksi yang disertakan (MA-04).
2. **Kalibrasi sebelum percaya pada nol (TA-04).** Sapuan bernilai nol tidak dilaporkan sebelum perkakasnya terbukti menemukan satu kasus positif yang sudah diketahui. Dasarnya sudah ada di induk: `SAPUAN-DAN-NAMA-TAGNYA.md` — *"sebuah sapuan atas korpus yang berisi banyak jenis aturan tidak boleh bersandar pada satu nama tag"*.
3. **Bentuk penulis properti.** Ada lima bentuk, masing-masing punya tag sendiri:

   | Bentuk | Tag |
   |---|---|
   | Property-Set di Activity | `PropertiesName`/`PropertiesValue` |
   | Data Transform | `pyPropertiesName`/`pyPropertiesValue`, dengan `pyDisabled` |
   | kontrol Section | `pyPropertyTarget` |
   | pengikatan sel (IND-4) | `pyValue` pada sel yang tidak read-only |
   | Report Definition | `pyTargetProperty` |

   Titik buta: 32 langkah Java (sudah disapu sebagai teks — nol untuk properti bisnis; `adoptJSONObject` hanya memulihkan nilai tersimpan) dan 191 langkah SQL/REST (tidak pernah menulis langsung ke `TreatyIn.*`).
4. **Tiga penanda mati:**
   - `pyStepsBlockName = "//"` untuk langkah Activity — termasuk langkah bersarang di bawah langkah yang mati;
   - `pyDisabled = true` untuk Data Transform;
   - `pyStepsPreCondition = false` untuk precondition yang tertulis tetapi nonaktif.

   `langkah-hidup.py` induk buta terhadap `pyDisabled` (1.306 kemunculan, semuanya di Section/Harness/Obj-Model). Ketujuh klaim mati yang menopang keputusan sudah diaudit ulang dan semuanya bertahan.
5. **`pyDisabled` di dalam mode kontrol (`Embed-Control-Mode`) adalah setelan mode, bukan "kendali mati"** (IND-5):
   - `pyDisabledNew = always` → selalu nonaktif (146 mode);
   - `true` dengan `pyDisabledWhen` → nonaktif bersyarat (249 mode);
   - `false` → aktif.
6. **Keterlihatan = hasil perkalian semua tingkat** (sel × wadah × privilege). `pyVisible` adalah pemilih: `OTHER` berarti "pakai kondisinya", `ALWAYS` berarti "abaikan kondisi di sebelahnya" (TA-06).
7. **Parameter dibaca sebagai pasangan nama-nilai, bukan dari keberadaan nama** (TA-07). Contoh: tombol Edit mengirim `revisionstate=` kosong.
8. **Baca temuan bernomor sendiri sebelum membangun rantai sebab** (TA-08). **Baca artefak induk sebelum menyatakan temuan baru** (TA-05, berlaku untuk semua artefak induk, bukan hanya ADR).
9. **Pembalikan klaim dinyatakan terang di badan laporan**, tidak hanya sebagai perubahan artefak.

### 3.3 Format berkas grilling — mengikuti contoh GRILL-05

Satu folder per ronde, dan ronde = cabang: `GRILL-A/`, `GRILL-B/`, dan seterusnya. Isinya delapan berkas berkepala blok kutip (Modul · Ronde · tanggal, Peran, Masukan, Status, Sifat):

| Berkas | Isi di sesi ini |
|---|---|
| `00-LINGKUP` | cakupan, bukti yang dipakai dan yang sengaja tidak dipakai, aturan berhenti, aturan bukti (§5a) |
| `01-TEMUAN` | temuan baru ronde (NA-xx untuk A, NB-xx untuk B) |
| `02-SIDANG` | sidang atas klaim bahan sebelumnya: DIKUATKAN / DIPERLUAS / SALAH KAPRAH / DITUTUP |
| `03-LUBANG` | EXP, DB, UA, titipan ke cabang lain |
| `04-PARITAS` | skenario paritas (untuk A: dua, keduanya paritas data) |
| `05-GERBANG` | kesiapan untuk to-spec |
| `06-PUTUSAN` | GRL ronde itu, ditulis dengan peran penilai; koreksi metodologis MA-xx |
| `07-AUDIT` | buku besar, kesalahan penggrill (§2a), daftar untuk dibantah siap kirim (§3), titik periksa (§4) |

Fase yang tidak berlaku diisi satu kalimat alasannya, bukan pengisi. Log `KEPUTUSAN-GRILLING-ADJUSTMENT.md` menjadi **indeks**: status cabang, Lacak TDA, Koreksi atas bahan, ekspor tambahan, daftar eskalasi, dan satu baris per GRL yang merujuk berkasnya.

### 3.4 Sumber induk dipakai aktif

Instruksi pengguna: grilling dan to-spec memakai `D:\XML_NURE\_migration-docs\treaty-in` supaya cepat selesai.

- Sebelum setiap butir GRILL, dan sebelum menggolongkan apa pun sebagai REKOMENDASI, cari lewat `PETA-SUMBER-INDUK.md`. Bila sudah dijawab → KONFIRMASI dengan kutipan dan barisnya. Bila sebagian → tanyakan sisanya saja.
- **Urutan wewenang:** untuk fakta sistem lama, ekspor dan `PENGETAHUAN.md` menang (§3.7). Untuk rancangan, ADR dan spesifikasi induk menang, kecuali ada REV.
- Kasus yang sudah terjadi: bentuk simpan rujukan versi dasar (`STRUKTUR-DATA.md` §1.1), aturan sapuan (`SAPUAN-DAN-NAMA-TAGNYA.md`), `KUNCI_PADANAN` (`SEAM-ADJUSTMENT` §3), dan fakta jalur kelompok (ADR-0055 §4.1) — semuanya sudah ada di induk sebelum "ditemukan" ulang.

---

## 4. Ringkasan METODE-GRILLING

- **Pertanyaan pokok (§1.2), diulang di setiap butir:** *ini kebutuhan bisnis, atau akibat dari cara sistem lama dibangun?*
- **Prinsip payung (§2.0):** struktur yang terlihat bukan struktur yang berlaku. **§2.0a:** rujukan bukan panggilan; panggilan bukan langkah hidup. **§2.2:** penegakan yang bisa berhenti diam-diam menuntut uji negatif dan pemantau.
- **Pola berulang:**

  | § | Pola |
  |---|---|
  | 3.1 | diam bukan bukti — penyisiran harus tertutup |
  | 3.2 | klaim "tidak pernah" aman untuk menunda, berbahaya untuk membuang |
  | 3.3 | nilai bawaan menyamarkan perbedaan |
  | 3.4 | kemiripan bentuk bukan bukti jenis — versi menggantikan, salinan berdampingan |
  | 3.5 | uji yang dijawab "ya" oleh kedua kemungkinan bukan uji |
  | 3.6 | preseden tidak otomatis berlaku |
  | 3.7 | ADR tidak memutuskan fakta; bila fakta ≠ rancangan, itu PERUBAHAN |
  | 3.8 | kemampuan mati yang dinyalakan adalah kemampuan baru — rencanakan penyalaannya |
  | 3.9 | hitungan bukan daftar |

- **Cara bertanya (§4):** satu per giliran (§4.1); pilihan bernomor (§4.2); tanyakan bukti, bukan pendapat (§4.3); "apa yang harus benar supaya ini salah" (§4.4); tarik akibat sampai uang, migrasi, dan orang di hari pertama (§4.5); pilih sumber yang benar — ekspor, data, atau orang (§4.6); daftar untuk dibantah (§4.7).
- **Menilai jawaban (§5):** bukti berkas mengalahkan argumen rapi (§5.1); uji negatif (§5.2); kriteria yang tidak bisa gagal bukan kriteria (§5.3); terima koreksi dengan terang (§5.4); hentikan dan laporkan bila tidak dapat dijawab (§5.5).
- **Proses (§6):** tulis saat langkah selesai (§6.1); lubang dilaporkan, tidak ditambal (§6.2); cacat tidak diperbaiki diam-diam (§6.3); empat golongan, tidak ada yang kelima (§6.4); label (§6.5); butir eskalasi bukan pekerjaan teknis (§6.6); setiap hasil membawa alasannya (§6.7).
- **Penutupan (§7):** batas pertanyaan ditetapkan di muka (§7.1); selesai bila **"semua bisa di-spec"** (§7.2); titik periksa berkala (§7.3); aturan embargo (§7.4).
- **Bagian VIII basi di beberapa titik** — lihat §12.

---

## 5. Fakta sistem lama yang sudah terbaca

Kadar setiap klaim mengikuti laporan AI grilling. Rincian lengkap ada di `PENGETAHUAN.md` dan folder `GRILL-*`.

### 5.1 Data dan penyimpanan

- **Addendum adalah salinan utuh kontrak** di `M_TREATY_IN_EDM` (JSON) dan `TREATY_IN_EDM` (datar). Prosedur `PEGA_M_TREATY_IN_EDM` menulis ke keduanya. Pengenalnya `‹kontrak›/Rnn`, dengan `OLDID` = ID baris yang **dipilih** di picker (kontrak atau addendum).
- **Satu `JSONDATA` memuat empat pohon:** nilai baru, `OLDDATA` (potret **beku** kontrak asal, bukan rujukan), `ActualValue`, dan `ValueDifference`. Nilai lama karena itu dapat direproduksi. Selisih tersimpan belum tentu benar (TDA-04, TDA-05) dan sebagian tidak pernah tersimpan (TDA-06).
- **`ValueDifference`:** 33 akar (21 di antaranya agregat `Total*`), 165 jalur daun. Hanya tujuh akar yang ditulis aturan ber-awalan `TreatyEDM*`; sisanya ditulis aturan milik Treaty In.
- **`EDMDATE = SYSDATE`** pada setiap pembaruan (jam Oracle).
- **Irisan dengan ekspor Treaty In:** 317 berkas menurut `pzInsKey`, tidak byte-identik. Ada enam jalur berbeda versi. Permukaan khas Adjustment ±56 berkas.
- **Sisi 323/56:** sepuluh dari lima belas TDA ada di irisan — cacat Treaty In yang berjalan hari ini.
- **Modul tidak memakai objek kerja Pega** (kelas `Data-Portal`; nol Obj-Save dan penyelesaian assignment). Ada satu aturan Connect-REST (ServiceGoogle, untuk unggah dokumen).

### 5.2 Pembuatan addendum dan penomoran

- **Tombol Revisi** → `TreatyCreateEDM` menyemai `EDMState = 1`. **Tombol Penyesuaian** → `EDMState = 3`.
- **Picker Revisi** memuat dua radio yang dapat disunting tanpa syarat apa pun (NA-01):
  - `EDMState` = {1 Internal, 2 External}
  - `EDMMaterialType` = {1 Material, 2 Non Material}

  Sumbernya adalah `PromptList` properti di kelas `ASM-FW-GISFW-Int-TREATY_IN` (EXP-1, versi 01-01-56; versi yang lebih tinggi *perlu dicek*). Catatan pengembangnya: "display only not for validation". `EDMState` dibuat 16 Juni 2020; `EDMMaterialType` 26 Oktober 2020.
- **Kombinasi sah lewat layar ada lima:** (1,1), (1,2), (2,1), (2,2), dan (3,1). Nilai lain di data berarti anomali. Penyesuaian premi selalu material di kedua jalur.
- **`TreatyInSetEditPre`:**
  - 1 → komentar "Had Created Internal Edit" dan `EDMEffective ← Commencement`;
  - 2 → "Had Created External Addendum";
  - 3 → "Had Created Addendum Premium", `ActualValue.EGNPI ← EGNPI`, `AddendumPremi ← 1`.

  Jenis 3 juga membelokkan `SaveTreatyIn_EDM_Act` langkah 2.
- **Penomoran** (`TreatyInRevisi_post`): `HASIL1` dari `GetTreatyRevisionID` hanya diuji kosong atau tidak; **nilainya dibuang**. Nomor dihitung ulang dari `TreatyIn.ID` lewat `@substring(ID,10,12)`, yang meleset satu dan hanya membaca digit satuan. Akibatnya:
  - memilih `…/R03` menghasilkan `…/R04`;
  - memilih `…/R01` ketika R02 sudah ada menghasilkan R02 lagi;
  - deret R09 → R010 → R11 → R02 — bergantung pada perilaku `@substring` di luar panjang teks (tidak terbaca). Kemungkinannya: diam (menimpa) atau berisik (galat).
- **Penjaga duplikat** (`TreatyInEdmCheckDuplicate`) **mati**. Tabrakan masuk cabang `UPDATE`, menimpa addendum lain, dan melapor berhasil (TDA-01).
- **Picker** memakai `pySourceType = Property`, bersumber `TreatyLoadMasterJoinEdm`: UNION `treaty_in` dengan `treaty_in_edm` tanpa WHERE, tanpa pembeda jenis, dan tanpa saringan keadaan (TDA-11). `pyRDName` yang juga tercantum hanyalah sisa pengaturan.

### 5.3 Materialitas dan penyuntingan

- **Materialitas hanya hidup di layar** — dan ia **sakelar dua arah**, bukan satu. Hitungannya di `GRILL-D/01-TEMUAN.md` `TD-02`; **angka 220/18 dicabut 24 Sep 2026** (`MA-13`) karena tidak dapat direproduksi. `InputTreatyInOffer` tidak punya kondisi sendiri; hitungan 70 sebelumnya ternyata salinan. Nol penegakan di sisi simpan (`TDA-10`) — **ditutup `GRL-20` butir 3**.
- **Sensus sel FIELD:** 762 aktif tanpa syarat · 116 selalu nonaktif · 83 bersyarat ViewState · 17 bersyarat lain.
- **Lapisan beku di layar addendum:** Ceding, Commencement, dan Termination dapat disunting tanpa syarat, masing-masing satu kontrol, di kedua layar (juga sesudah Revision atau Force Edit). ProportionType tidak terjangkau di layar addendum. CedingID dan source of business tidak punya kontrol di pohon layar addendum.
- **NA-08:** kontrol Ceding di `InputTreatyInAdjustment` hanya menulis nama, jadi CedingID tertinggal. Penulis CedingID ada empat: autocomplete, `TreatyInMappingDataconvert`, `TreatyInSetReinsured` (keduanya menulis nama dan ID bersamaan), dan `ShowSummary` (teks bebas). Prosedur simpan menerima nama dan ID sebagai dua parameter, sehingga ketidaksinkronan ikut tersimpan.
- **ShowSummary adalah layar pengembang.** Tombol pembukanya bersyarat `pyUserName = 'ALDO SAPUTRA1'`. Apakah operator itu ada adalah pertanyaan data (UA-15).
- **`TreatyInSetTreatyYear`** (Termination = Commencement + 1 tahun) tidak terjangkau: sel pemicunya hanya-baca (NA-20).

### 5.4 Revisi-di-tempat Treaty In

- **Empat kontrol di layar penawaran** (`InputTreatyInOffer`): Edit, View, Copy, Revision. Keempatnya dibatasi peran inputor lewat `pyUserData`/`pyCondition` → `pyWorkBasketList(2) = 'ReasTreatyInAdmin'`.
- **Hanya Revision** yang mengirim `viewstate=1` dan `revisionstate=1`, dan hanya bekerja atas kontrak yang sudah Resolve Complete. Edit mengirim keduanya kosong, jadi tidak memendekkan rantai.
- **Revision → `SetTreatyIn_Act`** langkah 6–11 (precondition aktif):
  - `ViewState = 1`, `IsEditData = 1`;
  - `RevisionState = 1`;
  - `StatusAkseptasi` dikosongkan;
  - `PositionUsername` diisi;
  - komentar "Create Revision" (jam Oracle);
  - simpan.

  Pengajuan berikutnya masuk `Akseptasi_DT` cabang 2 (Admin → SecHead → selesai).
- **Jalur Adjustment tidak menulis baris kontrak.** `TreatyInEDMSetValue` langkah 2 memanggil `SetTreatyIn_Act` tanpa parameter, sehingga langkah 7–11 dilewati. **TDA-07 dicabut.**
- **Sesudah Revision, `ViewState = 1` hanya mengunci sebagian field uang.** Yang masih dapat disunting antara lain: TreatyInTabsNonProportional (36 sel, termasuk Deductible, AggregateLimit, Amount), TreatyRetroList (26, termasuk Deductible, Limit, NetPremi), Share dan ShareRetro (Value, Pct, Currency), DetailLimits (4), ditambah tujuh seksi lain. Total angsuran, 14 sel LayersEDM, dan DetailEGNPI selalu nonaktif. **Jadi bukan "revisi tanpa uang".**
- **Penulis `RevisionState` ada lima:**

  | Penulis | Nilai |
  |---|---|
  | `SetTreatyIn_Act` langkah 7 | 1 |
  | `Akseptasi_DT` 2.2.1.4 (Accept) | "" |
  | `Akseptasi_DT` 2.2.2.4 (**Reject**) | **1 lagi** |
  | `Akseptasi_DT` 2.2.3.4 (Decline) | "" |
  | `TreatyInCopy` (salin kontrak) | "" |

- **Penularan ke addendum.** Addendum mewarisi `RevisionState` dari baris kontrak yang dimuat, karena jalur addendum tidak menulisnya. Akibatnya: berantai dua tingkat, dan `TreatyInSetEdit` langkah 2 menyetel `ViewState = 1`, sehingga addendum lahir **terkunci sebagian**. Ini terjadi bila kontrak sedang direvisi-di-tempat, atau revisinya ditolak atau tidak selesai.

### 5.5 Persetujuan, penolakan, dan pintu samping

- **Rantai empat tingkat:** Admin → SecHead → DeptHead → Director.
- **Addendum yang disetujui tidak menggantikan kontrak.** `SaveTreatyIn_EDM_Act` langkah 12–13 mati.
- **Menolak addendum menghapus barisnya** (TDA-02), tanpa membersihkan detail dan lampiran (TDA-03).
- **Jalur GroupLeader tidak dapat dimasuki.** Tidak ada penulis `Position = ReasTreatyInGroupLeader` di kedua ekspor (sapuan 24 penulisan, nilainya tertutup). `Akseptasi_DT` 1.4 `pyDisabled` (dibuat 2019, lalu dinonaktifkan). Tombol di `TreatyInActionButtons` hanya membaca `Position`. `TreatyInSetValue` 2.6 menulis `PositionUsername = "BERNARD"`.
- **Nama orang tersemat di aturan** (TDA-09): "BERNARD", "IRVANDY", "YOHANESKRISTIAWAN", "NANDINA", "ALDO SAPUTRA", "Daniel Suhana".
- **`TreatyInSetToDirector` mati** (langkah 1–4 ber-blok `//`, 4.1–4.7 ikut mati karena induknya).
- **Tombol pengembang:**

  | Tombol | Terlihat oleh | Akibat |
  |---|---|---|
  | Force Edit (dev) | setiap operator divisi IT (sel `ALWAYS`) | `ViewState = 0` |
  | Force Resolve Complete(dev) | operator divisi IT bernama ALDO SAPUTRA atau Daniel Suhana | menyetel selesai disetujui tanpa penyetuju, lalu `SaveTreatyIn_Act` |
  | ReturnToInputor(dev) | idem | mengembalikan ke admin, tanpa menyimpan |
  | Save EDM(dev) | divisi IT, pada halaman addendum, **tanpa syarat status** | menyimpan addendum |

- **Force Resolve di layar addendum** memanggil prosedur **kontrak** `PEGA_TREATY_IN` → `UPDATE … WHERE ID = 'XXXXXXX/Rnn'` → nol baris berubah, tanpa galat, dengan pesan "Data Sudah Disimpan". Kegagalan senyap murni.
- **Force Edit (dev) ditambah Save EDM(dev)** membuat operator IT dapat membuka kunci addendum yang sudah disetujui, menyuntingnya, lalu menyimpannya di tempat. `SaveTreatyIn_EDM_Act` tidak menghitung ulang selisih; mesin selisih hanya dipanggil `TreatyInSubmitEDM` (langkah 6 komentar, 7 hitung selisih, 8 simpan) dan kontrol tab. Hal ini dapat diukur (UA-16).
- **`CommentList`** hanya punya ruas OperatorName, IsApproved, Suggest, dan Date — tidak merekam tingkat penyetuju. Zona waktunya campuran: komentar persetujuan memakai jam Pega (GMT); "Create Revision" memakai jam Oracle; langkah jam Oracle di `AddCommentList_Act` langkah 1 mati. Kapan pergantian zona terjadi tidak terbaca.

### 5.6 Mesin selisih

- **Rumus:** baru − `OLDDATA`, dipadankan menurut **posisi** baris (TDA-04), tanpa memeriksa mata uang (TDA-05). Untuk `EDMState` 3 ada percabangan: share dihitung dari `ActualValue` (§5.4 PENGETAHUAN).
- **Kemampuan yang mati:**
  - pro rata — 10 dari 53 langkah mati: lima blok `//` ditambah lima langkah bersarang. `EDMEffective` hanya dibaca `TreatyCalculateProratePct`;
  - potongan share fakultatif (`TreatyEDMDifferenceDeduction` langkah 4–5, "not enable yet");
  - ringkasan share yang ditulis ke `ActualValue` lalu ditimpa (TDA-06).

---

## 6. Keputusan terkunci (GRL)

Isi lengkap setiap keputusan ada di `GRILL-X/06-PUTUSAN.md`. Di bawah ini intinya, beserta label dan syarat-syarat yang diminta penggrill.

### Cabang A — rekonsiliasi batas lingkup dan ADR induk (selesai)

**GRL-01 (A1) — Satu model.**
- Addendum adalah `VERSI_KONTRAK`. Entitas `PENYESUAIAN` dibuang.
- Hanya ada satu ERD, satu daftar invarian, satu daftar keadaan, dan satu rantai persetujuan — seluruhnya milik spesifikasi induk.
- Susunan berkas spesifikasi termasuk bahan to-spec.
- `NILAI_SELISIH` **bukan** entitas sementara (ralat PG-02): ia mengikat lewat ADR-0048 butir 3 dan gerbang METODE §8.3. `BESARAN_DAPAT_DISESUAIKAN` masih sementara.
- Temuan yang mengubah entitas induk dicatat sebagai usulan untuk langkah 8–10 induk.
- Embargo Adjustment dicabut per 23 September 2026. CONTEXT.md induk §1.1 dan §4 sudah diperbarui dengan izin.
- **Label:** PERUBAHAN, berubah dari "addendum berdampingan; kontrak tidak pernah digantikan".
- **Yang tersirat:** versi yang disetujui menggantikan pendahulunya.
- **Titipan:** ke D (rencana penyalaan, §3.8) dan ke I (versi berlaku kontrak warisan).
- **Menunggu:** PP-1 — siapa yang mengerjakan to-spec induk. Hanya memengaruhi arah dampak.

**GRL-02 (A2) — ADR-0036 berlaku dengan penyesuaian.**
- **Prinsip:** angka dasar persetujuan sebuah addendum — termasuk selisih yang ditampilkan kepada penyetuju — tidak berubah sesudah persetujuan, juga bila rumus diperbaiki kelak.
- **Mekanisme** diputuskan di E. Pilihannya tinggal dua: bekukan hasil beserta pemadanannya, atau versikan rumus.
- **Label dibagi tiga:**
  - niat "angka persetujuan beku" → PELESTARIAN, karena `OLDDATA` dan `ValueDifference` memang beku di JSON;
  - baris yang disetujui tidak boleh dapat ditimpa → PERUBAHAN (TDA-01);
  - setiap angka ringkasan tersimpan → PERUBAHAN (TDA-06).

  Sumber PERUBAHAN ketiga: Force Edit + Save EDM(dev) yang menyunting addendum yang sudah disetujui.
- **Menunggu:** DB-1.

**GRL-03 (A3) — ADR-0048 berlaku dengan penyesuaian.**
- **Batas ADR mati:** bukan perilaku, tidak diuji.
- **Butir 3 mengikat** — selisih disimpan di tabel tersendiri. Label PERUBAHAN, berubah dari sub-pohon `ValueDifference` di JSON yang sebagiannya tidak tersimpan.
- **Butir 2 mengikat** — "data lama tidak disimpan ulang, diambil lewat SELECT", aturan pemilik proses (METODE §8.5). Label PERUBAHAN, karena sistem lama menyalin. Yang membuat SELECT aman adalah pembekuan.
- **Butir 1 dibaca sebagai ID versi:** PELESTARIAN sementara. *Catatan penggrill:* syaratnya adalah "menjadi PERUBAHAN bila B memilih rantai linear", dan GRL-10 sudah menghapus percabangan. **Label ini perlu diperbarui** — *(perlu dicek apakah sudah)*.
- Penolakan pilihan (b) beralasan §3.7: fakta tidak membatalkan rancangan.
- **Menunggu:** DB-2 (sudah dirumuskan ulang).

**GRL-04 (A4) — Penyesuaian bukan entitas.**
- `PENYESUAIAN` tidak kembali. Label PELESTARIAN.
- Terkunci dari sisi ekspor. Satu-satunya pembatalnya ada di bisnis (METODE §8.6 Q1, diperluas ke kardinalitas): satu dokumen addendum mengubah lebih dari satu kontrak atau versi, **atau** satu revisi menggabungkan lebih dari satu dokumen. Diukur lewat DB-3 dan DB-4.
- DB-5 (addendum berlaku di tengah periode) tidak membatalkan keputusan ini, tetapi mengubah C dan E.

**GRL-05 (A5) — ADR-0040 berlaku dengan penyesuaian.**
- Butir 1 → PELESTARIAN.
- **Butir 3 (lapisan beku)** adalah **larangan**: lima hal tidak boleh berubah — cedant, source of business, ProportionType, tanggal mulai, tanggal berakhir; "perubahan atas salah satunya berarti kontrak lain". Label PERUBAHAN, karena sistem lama membuka Ceding, Commencement, dan Termination tanpa syarat.
  - UA-10 mengukur kelima field, dengan nama dan ID cedant dipisah.
  - DB-6 dan DB-7 menanyakan perpanjangan periode dan pergantian cedant lewat addendum. Bila dibantah → usulan revisi ADR-0040.
  - Tahap peringatan akan menyimpang dari ADR, jadi harus lewat REV. R-D2 dicabut.
- Butir 5 dilepas ke C1.

**GRL-06 (A6) — Klausa data warisan ADR-0049.**
- Materialitas warisan dibawa apa adanya, dengan asal-usulnya. Nilainya tidak ditimpa atau diganti.
- Menghitung pembanding untuk deteksi boleh (UA-3 membutuhkannya), asalkan hasilnya disimpan terpisah.
- Klausa ini tidak diperluas ke semua data warisan: cabang I harus menjalankan aturan baru atas baris warisan setidaknya untuk menentukan versi berlaku.
- Label PELESTARIAN. A6 ditutup dengan rujukan ke C1.

**GRL-07 (A7) — ADR-0052 berlaku.**
- Keputusan 1–3 → PELESTARIAN. Jalur kelompok tidak dapat dimasuki, jadi membuangnya bukan perubahan.
- Keputusan 4 — rantai empat tingkat, tanpa pemendekan → PERUBAHAN (TDA-08).
- DB-10 ditulis ulang secara netral: adakah kebijakan rantai pendek untuk revisi, dan apa yang sebenarnya dapat berubah sesudah Revision. DB-10b dicabut.
- UA-13 mengukur tambahan beban penyetuju. REV-2 diajukan.

**GRL-08 (A8) — ADR-0055: daftar keadaan dipakai untuk versi addendum tanpa keadaan tambahan.**

| Perpindahan | Label |
|---|---|
| LAHIR, AJUKAN, SETUJUI×3, KEMBALIKAN×3 | PELESTARIAN |
| TOLAK → DITOLAK | PERUBAHAN (TDA-02) |
| DITOLAK sebagai keadaan terminal | PERUBAHAN |
| rantai empat tingkat | PERUBAHAN |
| DRAFT → DIBATALKAN, PERBAIKAN_WARISAN | BARU |

- SecHead Reject yang menulis ulang `RevisionState = 1` tidak punya padanan, dan memang tidak diperlukan.
- **Butir iii:** membuang `ViewState` membuat addendum yang terkunci (sebagian) karena warisan menjadi dapat disunting. PERUBAHAN.
- Preseden GRL-06 tidak berlaku untuk kelima penanda lama, karena kelimanya mekanisme, bukan keputusan manusia. `Position` dan `StatusAkseptasi` tetap dipakai untuk memetakan baris yang sedang dalam proses saat peralihan.
- REV-3 diajukan. TDA-02 → diperbaiki. TDA-03 → diperbaiki; sisa warisannya ke I (UA-7).

### Cabang B — model versi dan identitas (selesai)

**GRL-09 (B1) — Pengenal tampilan.**
- **Versi baru:** diturunkan saat ditampilkan, tidak disimpan. Bila kelak perlu disimpan, harus beralasan dan dijaga invarian dengan uji negatif.
- **Baris warisan:** dilestarikan apa adanya, dengan preseden GRL-06 — nomor itu sudah beredar di luar sistem.
- **Bahan to-spec B-2:** pengenal versi baru tidak boleh sama dengan pengenal warisan mana pun di kontrak yang sama.
- Awalannya mengikuti ADR-0040 butir 1, tanpa anggapan tujuh karakter.
- Melestarikan nomor tidak memulihkan isi yang tertimpa (UA-14).
- **Label:** pengenal warisan → PELESTARIAN; bentuk turunan tanpa lebar tetap → PERUBAHAN (TDA-12).
- **Menunggu:** DB-11.

**GRL-10 (B2) — Rantai versi.**
- **Tidak ada percabangan.** Versi dasar adalah versi berlaku terakhir saat versi baru dibuat — bukan n−1. PERUBAHAN, berubah dari "dasar = baris mana pun yang dipilih di picker".
- **Bentuk simpan → KONFIRMASI:** `ID_VERSI_KONTRAK_DASAR` disimpan eksplisit (`STRUKTUR-DATA.md` §1.1).
- `OLDID` warisan dibawa sebagai atribut warisan → PELESTARIAN.
- **Bahan to-spec, berlaku untuk versi non-warisan:**
  - B-3: `ID_VERSI_KONTRAK_DASAR` wajib terisi pada setiap versi penyesuaian, dan wajib kosong pada versi pertama;
  - B-4: yang ditunjuk harus versi berlaku terakhir, bukan DITOLAK atau DIBATALKAN.
- **Bahan to-spec B-1 (INV-04):** menyimpan versi dengan nomor urut yang sudah dipakai ditolak, dan pesannya menyebut nomor yang bentrok.
- Cabang B ditutup dengan KONFIRMASI untuk: picker (INV-25), revisi-di-tempat menjadi versi (ADR-0055), dan deteksi duplikat (ADR-0040).

### Cabang I — migrasi (sebagian)

**GRL-11 (I1) — Versi berlaku adalah turunan.**
- Rumusnya: versi ber-KEADAAN DISETUJUI dengan `NOMOR_URUT_VERSI` tertinggi. KONTRAK tidak membawa kolom penunjuk.
- Dari delapan keadaan ADR-0055, tidak ada keadaan "digantikan", dan DISETUJUI satu-satunya yang berarti "pernah berlaku".
- **Syaratnya:** urutan nomor = urutan persetujuan. Untuk versi baru, INV-25 menjaminnya. Untuk baris warisan, I2 yang harus menjaminnya.
- Bahan to-spec I-1 dan I-2 sudah ditulis beserta kasus yang dapat gagal.
- **Titipan ke I2:** ambiguitas arti hukum (addendum terakhir yang disetujui) versus arti operasional (baris `TREATY_IN` yang dibaca hilir, yang dapat disunting di tempat lewat Revision sesudah addendum terakhir). Diukur UA-17, akibatnya ke orang lewat DB-12.
- **Label:** PERUBAHAN, berubah dari "kontrak tidak pernah digantikan".

### Cabang C — jenis dan materialitas (sebagian)

**GRL-12 (C1) — Materialitas diturunkan dari akibat.**
- Material = versi itu memiliki sedikitnya satu baris `NILAI_SELISIH`, dibaca dari baris yang **tersimpan dan beku**, bukan dari hitung ulang.
- Tidak ada atribut untuk versi **baru**. `EDMMATERIALTYPE` warisan tetap disimpan (GRL-06).
- **Pilihan (c) ditolak** karena §3.8 (penolakan baru di hari peralihan, frekuensinya diukur UA-3) dan ADR-0037 butir 3.
- **Pilihan (b) ditolak** dengan bukti kombinasi: jenis 1 dan 2 masing-masing dapat dipasangkan dengan kedua materialitas.
- Pembatalnya DB-15: bila bisnis memakai materialitas **sebelum** perubahan dibuat, keputusan ditinjau ke arah (c).
- R-H1 (bantuan layar yang tidak disimpan) dan pemberitahuan ke pengisi dititipkan ke H. REV-4 diajukan untuk ADR-0049 (sumber turunannya).
- TDA-10 → diperbaiki. **Label:** PERUBAHAN.

**GRL-13 (C2) — Jenis addendum.** *Status: sudah dijawab penggrill di prompt terakhir; belum ada laporan bahwa AI grilling menuliskannya.*
- **Untuk versi baru:** dua nilai — `PENYESUAIAN_PREMI` dan satu nilai untuk selebihnya (nama finalnya bahan to-spec).
- DB-16 menjadi pintu ke (b+), yaitu penanda "berasal dari dokumen eksternal". Tempatnya — mungkin `DOKUMEN_KONTRAK` — diputuskan sesudah DB-16 terjawab.
- `EDMState` warisan dibawa apa adanya. REV-5 diajukan untuk himpunan jenis ADR-0049.
- **Pemasangan jenis:** tiga jenis ADR-0049 (perubahan estimasi, penyesuaian ke nilai aktual, administratif) bukan Internal/External/Premi. Sesudah GRL-12, "administratif" menjadi turunan, dan ADR-0049 menyusut ke bentuk dua nilai.
- **Catatan penggrill:**
  - EGNPI lazimnya berarti *Estimated* Gross Net Premium Income — konfirmasi lewat DB;
  - **"penyesuaian premi selalu material" wajib tetap benar sesudah GRL-12** — ini syarat untuk C3;
  - radio Internal/External hilang dari layar — titipan H dan pemberitahuan ke pengisi;
  - label: `PENYESUAIAN_PREMI` → PELESTARIAN; pembedaan Internal/External → PERUBAHAN.

---

## 7. Pemilahan anggaran dan posisi terakhir

- **Pemilahan 43 butir** (`PEMILAHAN-SISA-GRILLING.md`) memakai empat golongan: GRILL, KONFIRMASI, REKOMENDASI TO-SPEC, dan DITUNDA. Keputusan wewenang → DITUNDA, dengan pemilik manajemen dan dicatat di daftar eskalasi (ralat PG-06).
- **Anggaran:** 40 → tujuh butir GRILL (C1, C2, C3, E1, E3, I1, I2) → E1 turun ke KONFIRMASI (`SEAM-ADJUSTMENT` §3, berstatus mengikat) → **enam**.
  - Terjawab: I1, C1, C2.
  - **Sisa: C3, E3 (E3a pro rata, E3b share fakultatif, E3c ringkasan tak tersimpan), dan I2.**
  - Titik periksa berikutnya sesudah butir kelima.
- **Penggabungan:** penyesuaian premi → C2; "yang mati/hilang" → E3.
- **KONFIRMASI penting:**

  | Butir | Sumber |
  |---|---|
  | F — bentuk simpan relasional | `STRUKTUR-DATA.md`, ADR-0034 |
  | F1 — cedant | CEDANT [luar] (`STRUKTUR-DATA.md` §3) |
  | J — lampiran | `DOKUMEN_KONTRAK` — dirujuk, tidak dimiliki |
  | G — peran | ADR-0044 — peran dan penugasan bertanggal, tanpa nama orang di aturan |
  | lapisan beku | ADR-0040 butir 3 (larangan) |
  | bentuk tabel selisih | METODE §8.3 |
  | E1 — kunci padanan | `SEAM-ADJUSTMENT` §3: Bentuk A (9 entitas) vs Bentuk B (DETAIL_PROPORSIONAL, POTONGAN, PENYEBARAN — kunci induk lebih dulu) |

- **REKOMENDASI TO-SPEC:** R-D1 (penyalaan "versi menggantikan", dan DB-12), R-F1 (penguncian optimistis dengan penanda versi baris), R-I1 (ambang migrasi: satu persen, atau satu kontrak berlaku), R-H1.
- **DITUNDA:** pemisahan pembuat dan penyetuju (DB-14), wewenang tombol "(dev)" (eskalasi), pemberitahuan ke penyetuju dan pembuat (UA-13, UA-9a), paket REV (prasyarat K).
- **Cabang yang tertutup tanpa ronde:** D dan K. F, G, H, dan J hanya berisi KONFIRMASI atau REKOMENDASI.

---

## 8. Titipan yang masih terbuka, per cabang

**C3 — arti bisnis `ActualValue`** (butir GRILL berikutnya).
- **Bahan:** PETA-TELUSUR-JSON ⚠ ("ActualValue.* (108) memakai entitas yang sama, bukan entitas baru"), PENGETAHUAN §2.2, percabangan `EDMState` 3 pada share (§5.4), dan TDA-06.
- **Wajib dijawab:** apakah perubahan `ActualValue` menghasilkan baris `NILAI_SELISIH`. Bila tidak, penyesuaian premi akan terbaca non-material dan bertentangan dengan perilaku lama.
- **Bahan to-spec:** "versi PENYESUAIAN_PREMI yang mengubah nilai aktual → material".
- **DB:** arti EGNPI.

**E — mesin selisih.**
- **Celah yang harus dicari:** di mana titipan "mekanisme pembekuan" berada di pemilahan. `KUNCI_PADANAN` tinggal di `NILAI_VERSI_KONTRAK`, yang merupakan bentuk baca turunan. Periksa apakah `NILAI_SELISIH` membawa nilai kunci atau identitas baris kedua sisinya.
  - bila ya → KONFIRMASI "bekukan hasil beserta pemadanannya";
  - bila tidak → laporkan sebagai celah, jangan ditambal.
- **E-1a:** mata uang masuk ke kunci (ADR-0053). Perubahan mata uang tampil sebagai baris dihapus ditambah baris baru.
- **UA-18:** keunikan kunci padanan pada data warisan.
- **E3a pro rata** (ADR-0037: "faktor prorata addendum = turunan dari tanggal"; bergantung DB-5; catatan `EDMEffective`), **E3b share fakultatif**, **E3c ringkasan tak tersimpan** — masing-masing baris keputusan tersendiri, dengan §3.8 untuk kemampuan yang mati.

**I2 — aturan baru atas baris warisan.** Sebelumnya baca ADR-0042, ADR-0043 (belum pernah dibaca), dan ADR-0054. Isi yang harus dijawab:
- bagaimana `KEADAAN` dan `NOMOR_URUT_VERSI` diberikan kepada baris warisan, secara kronologis, supaya turunan GRL-11 benar;
- `ID_VERSI_KONTRAK_DASAR` untuk warisan (B-3 dan B-4 hanya untuk non-warisan), dan `OLDID` mentah tetap disimpan di sampingnya;
- arti hukum versus operasional "berlaku hari ini" (UA-17, DB-12);
- kontrak yang seluruh addendumnya ditolak;
- addendum yang tertimpa (UA-14);
- nomor baru dimulai di atas nomor warisan (GRL-09 syarat 2);
- nama dan ID cedant yang tidak sinkron — mana yang benar (DB-8, DB-9);
- DRAFT saat peralihan (R-I);
- pelanggaran lapisan beku warisan (UA-10);
- pertanyaan umum: aturan baru mana yang boleh dijalankan atas warisan, dan asal-usulnya dicatat bagaimana.

**D:** rencana penyalaan "versi menggantikan"; beban penyetuju (UA-13); pemberitahuan ke pembuat addendum yang lahir terkunci (UA-9a).

**H:** R-H1; radio materialitas dan Internal/External hilang — pengisi kontrak perlu diberi tahu.

**K:** paket REV diserahkan dan ditanggapi; PP-1; `CABANG-K-PEMETAAN-TO-SPEC.md` (keputusan mana mendarat di artefak induk mana).

---

## 9. Daftar untuk dibantah (DB) — pertanyaan ke orang

Letaknya di `GRILL-A/07-AUDIT.md` §3, siap kirim, dikelompokkan per peran penjawab. Butir sesudah DB-7 ada di folder ronde masing-masing. Rumusan di bawah adalah inti yang diketahui penggrill; bunyi persisnya ada di berkas.

| # | Inti | Mengubah apa bila dibantah |
|---|---|---|
| DB-1 | Angka selisih yang dilihat penyetuju harus tetap sama bila addendum dibuka lagi kapan pun | alasan REV-1, GRL-02 |
| DB-2 | Revisi selalu dibuat dari versi terakhir; tidak pernah sengaja dari versi yang lebih lama | pemetaan warisan (I2) |
| DB-3, DB-4 | Satu dokumen addendum selalu untuk satu kontrak; dua addendum tidak pernah digabung jadi satu revisi | **pembatal GRL-04** |
| DB-5 | Addendum selalu berlaku sejak tanggal mulai kontrak, tidak di tengah periode | E3a, C |
| DB-6, DB-7 | Perpanjangan/pemendekan periode dan pergantian cedant selalu lewat kontrak baru | REV ADR-0040 butir 3 |
| DB-8, DB-9 | Pembaca hilir memakai nama atau ID cedant; mana yang benar bila tidak sinkron | I, F |
| DB-10 | Adakah kebijakan "revisi cukup disetujui kepala seksi", dan apa saja yang boleh berubah | REV ADR-0052 butir 4 |
| DB-11 | Nomor `…/Rnn` hanya dipakai di dalam aplikasi | GRL-09 bisa menyusut |
| DB-12 | Angka kontrak ke akuntansi dan pihak lawan akan mulai memperhitungkan addendum — siapa yang perlu tahu | D, I1 |
| DB-13 | *(tidak tercatat di percakapan ini — perlu dicek)* | — |
| DB-14 | Pembuat addendum tidak boleh menyetujui addendumnya sendiri | G (wewenang) |
| DB-15 | Materialitas tidak dipakai di luar layar (persetujuan, dokumen, akuntansi, laporan) | **pembatal GRL-12** |
| DB-16 | Addendum eksternal = dokumen yang ditandatangani cedant, dengan tanggal berlaku sendiri; revisi internal tidak dikirim ke luar | GRL-13 → (b+); terkait DB-3, DB-4, DB-5, DB-11 |
| DB-(baru) | Arti EGNPI (*Estimated* Gross Net Premium Income?) | C3, nama nilai jenis |

---

## 10. Data yang perlu ditarik

### 10.1 Ekspor Pega

| # | Isi | Status |
|---|---|---|
| EXP-1 | Aturan properti `EDMState` dan `EDMMaterialType` (ternyata `PromptList`, bukan Field Value) | **selesai** — di `ekspor-tambahan/`. Versi ruleset di atas 01-01-56 *perlu dicek* |
| EXP-2 | METODE-GRILLING.md | selesai |

### 10.2 Kueri produksi (seri UA)

Semua kueri ini tidak menghalangi rancangan (METODE §7.2). Hasilnya dipakai untuk migrasi, eskalasi, dan perkiraan pekerjaan. Jalankan di salinan atau replika, karena beberapa mengurai `JSONDATA` dengan `JSON_TABLE`. Teks kueri versi terakhir diminta disusun AI grilling dalam `KUERI-UA.sql` — *(perlu dicek apakah sudah dibuat)*.

| # | Yang diukur | Kadar dan catatan |
|---|---|---|
| UA-1 | bentuk pengenal rusak (panjang ≠ 11, ada tidaknya R10, ganda/melompat); (c) addendum yang `OLDID`-nya menunjuk addendum, atau bukan versi terakhir saat dibuat | (c) = percabangan yang sebenarnya |
| UA-2 | sebaran `EDMSTATE` × `EDMMATERIALTYPE` | di luar lima kombinasi sah → anomali |
| UA-3 | non-material yang nilai uangnya berubah | hanya dapat mematahkan, tidak dapat mengesahkan |
| UA-4 | kontrak yang `RevisionState`-nya tertinggal 1 | sumber penularan |
| UA-5 | panjang daftar berubah terhadap `OLDDATA` | selisih warisan salah (TDA-04) |
| UA-6 | mata uang berbeda antara lama dan baru | selisih warisan salah (TDA-05) |
| UA-7 | detail addendum yatim | TDA-03 |
| UA-8 | ukuran `JSONDATA`, `OLDDATA`, `ActualValue` | ukuran basis data baru |
| UA-9 | addendum yang mewarisi `RevisionState = 1` (tiga hitungan) | DRAFT saat peralihan |
| UA-10 | pelanggaran lapisan beku, lima field, nama vs ID cedant, nama berubah tanpa ID | angka pemicu REV ADR-0040 |
| UA-11 | kontrak yang kini belum disetujui padahal pernah disetujui | eskalasi tombol Revision |
| UA-12 | jejak jalur kelompok: (a) `Position` sekarang, (b) `PositionUsername = 'BERNARD'`, (c) persetujuan oleh operator itu di `CommentList` — keduanya berbentuk ID operator (`pxInsName`) | (b) dan (c) hanya dapat mematahkan |
| UA-13 | revisi-di-tempat per bulan (komentar "Create Revision") | beban penyetuju |
| UA-14 | jejak penimpaan: urutan `RevisionDate` terbalik; addendum bukan-draf dengan `CommentList` kosong | isi yang hilang tetap tidak terukur |
| UA-15 | operator `pyUserName = 'ALDO SAPUTRA1'` / `pyUserIdentifier = 'aldo_saputra'` | **tabel operator Pega**, bukan POOLDATA |
| UA-16 | addendum disimpan sesudah disetujui: (b) selisih basi — **tanda utama**; (a) `EDMDATE` > persetujuan terakhir + 5 menit | (a) hanya perkiraan kapan, karena zona waktu campuran; zona server Oracle *perlu dicek ke DBA* |
| UA-17 | kontrak dengan revisi-di-tempat sesudah addendum terakhir disetujui | I2 |
| UA-18 | keunikan kunci padanan per daftar di `JSONDATA` | E1 |

---

## 11. Paket REV, IND, dan eskalasi induk

**`USULAN-REVISI-ADR.md`** — draf; ADR induk tidak disunting. Diserahkan sekaligus di akhir grilling, kecuali PP-1 menunjukkan to-spec induk sedang berjalan. "Paket REV diserahkan dan ditanggapi" menjadi prasyarat K.

| # | ADR | Isi |
|---|---|---|
| REV-1 | 0036 | Alasan ditulis ulang. Premis SELECT **tetap berdiri** (GRL-03 butir 2); ia arahan rancangan, bukan pemerian sistem lama; pembekuanlah yang membuat SELECT aman. Paruh "dihitung saat dibaca" (angka posisi terkini) tidak dicampur |
| REV-2 | 0052 | Konteks: penyimpangan yang **berjalan** hanya satu, "dalam versi aturan yang ter-ekspor". Cabang 1.4 dibuat 2019 lalu dinonaktifkan. Pada addendum, rantai pendek diwarisi. UA-12(a) dapat membantah sebagian |
| REV-3 | 0055 | §4, empat fakta: keterlihatan Force Edit (semua IT) dan Force Resolve/Return (IT dan dua nama) sebagai hasil perkalian; Force Resolve menyimpan; `SetToDirector` mati (tetap didaftar, ditandai mati); pintu samping juga ada di layar addendum; Force Resolve di layar addendum = kegagalan senyap |
| REV-4 | 0049 | Sumber turunan materialitas (dari akibat, bukan jenis) |
| REV-5 | 0049 | Himpunan jenis (dua nilai untuk versi baru) |

**Temuan untuk induk (bukan REV):**

| # | Isi | Kadar |
|---|---|---|
| IND-1 | ShowSummary membuka field lapisan beku dan CedingID sebagai teks bebas | rendah, bergantung UA-15 |
| IND-2 | Empat kontrol `SetTreatyIn_Act` dibatasi peran inputor; hanya Revision yang bekerja atas kontrak yang sudah disetujui | tinggi |
| IND-3 | Save EDM(dev) tanpa syarat status | tinggi |
| IND-4 | Bentuk kelima penulis properti: pengikatan sel | — |
| IND-5 | `pyDisabled` di dalam mode kontrol adalah setelan per mode | — |

NA-17b dan CedingID teks bebas ikut dicatat untuk induk.

**`DAFTAR-ESKALASI-MANAJEMEN.md` induk** — sudah disunting dengan izin:
- **baris 97** diralat: sebabnya bukan `OLDID`, melainkan pemadanan posisi, mata uang, dan ringkasan tak tersimpan. Istilah "dibukukan" diganti "dilaporkan". Hanya berlaku dalam tiga keadaan. Merujuk UA-5 dan UA-6;
- **butir 1:**
  - Force Edit terlihat semua orang divisi IT (dikoreksi dari "setiap pengguna");
  - tombol **Revision** — satu tombol, dibatasi peran inputor — mengosongkan status persetujuan dan menyimpan sekaligus; tombol Edit tidak memendekkan rantai (diperiksa khusus); dikenali dari komentar "Create Revision";
  - Save EDM(dev) + Force Edit memungkinkan addendum yang sudah disetujui disunting di tempat — dapat diukur (UA-16);
  - pertanyaan wewenang: siapa yang mengizinkan pengembang menyelesaikan persetujuan pada data produksi;
  - semuanya dengan batas pemeriksaan "terbaca dari aturan, belum dijalankan".

---

## 12. TDA dan bahan yang sudah dikoreksi

### 12.1 Nasib TDA

Rangkuman penggrill. Kolom "Nasib" untuk butir yang tidak dinyatakan terang di percakapan adalah **perkiraan saya** *(cocokkan dengan tabel Lacak TDA di indeks)*.

| TDA | Pokok | Nasib |
|---|---|---|
| 01 | nomor dari baris yang dipilih; penjaga duplikat mati; tabrakan jadi `UPDATE` | diperbaiki: INV-04, B-1, B-2, GRL-10 |
| 02 | penolakan menghapus baris addendum | diperbaiki (ADR-0055, GRL-08) |
| 03 | penghapusan tidak membersihkan detail dan lampiran | diperbaiki untuk sistem baru; warisan → I (UA-7) |
| 04 | selisih dipadankan menurut posisi | diperbaiki: `KUNCI_PADANAN` (E1 KONFIRMASI) *(perkiraan)* |
| 05 | selisih tanpa memeriksa mata uang | diperbaiki: mata uang masuk kunci (E-1a) *(perkiraan)* |
| 06 | sebagian selisih ditulis ke `ActualValue` lalu ditimpa | GRL-02 (setiap ringkasan tersimpan) dan E3c — *terbuka* |
| 07 | ~~membuka kontrak menulis ke baris kontrak~~ | **dicabut** (SALAH KAPRAH): jalur Adjustment melewati langkah 7–11. Temuan yang tersisa milik Treaty In: tombol Revision |
| 08 | rantai memendek lewat `RevisionState` | premis ditulis ulang (penularan); diperbaiki (GRL-07) |
| 09 | peran dari workbasket, `pyTelephone`, dan nama tersemat | diperbaiki (ADR-0044) |
| 10 | materialitas hanya di layar | diperbaiki (GRL-12) |
| 11 | picker menyatukan kontrak dan addendum tanpa pembeda dan saringan | KONFIRMASI (INV-25); tampilan → H *(perkiraan)* |
| 12 | jenis baris dari panjang teks; offset tetap | diperbaiki (GRL-09: bentuk turunan tanpa lebar tetap) |
| 13 | jenis dan materialitas dipilih bebas di radio | ditutup sebagai fakta (EXP-1); keputusannya di GRL-12 dan GRL-13 |
| 14 | nilai disimpan sebagai teks | E presisi → REKOMENDASI *(perkiraan)* |
| 15 | penggandaan pohon dalam satu `JSONDATA` | KONFIRMASI (model relasional) *(perkiraan)* |

Kolom sisi 323/56 sengaja tidak diisi di sini. Menurut laporan AI, sepuluh dari lima belas TDA ada di irisan, tetapi rinciannya per TDA ada di tabel Lacak TDA di indeks.

### 12.2 Koreksi atas bahan

Semuanya ditulis sebagai blok koreksi bertanggal di berkas asalnya.

- **`KEPUTUSAN-SAMBUNGAN` §G2** ("selisih historis tidak dapat direproduksi") — dibantah: `OLDDATA` adalah potret beku.
- **METODE Bagian VIII** — beberapa klaim tidak berlaku lagi:

  | Klaim METODE | Keadaan sebenarnya |
  |---|---|
  | §8.1 "323 berkas byte-identik, `pzInsKey` sama" | salah dua kali: tidak byte-identik, dan irisan menurut `pzInsKey` 317 dengan enam jalur berbeda versi (NA-03) |
  | §8.3 "tabrakan mustahil" | gugur sebagai fakta; bentuk `/Rnn` tetap sebagai rancangan |
  | §8.4 "`OLDID` menunjuk kontrak → selisih historis tidak dapat direproduksi" | dibantah; butir eskalasinya diganti (baris 97) |
  | §8.4 / §8.6 Q2 (`OLDID` revisi kedua) | terjawab |
  | §8.4 `ROWNUM` | bobotnya turun |
  | §8.4 "R01..R99" | salah — patah di R10 (bersyarat) |

- **`PENGETAHUAN.md`:**
  - §1.2 — irisan dan `pzInsKey`;
  - §3.2 — radio picker;
  - §4.2 — judul dan jenis 2;
  - §4.3 — tabel materialitas dihitung ulang dari badan aturan;
  - §4.4 — sebab nomor revisi dan dua kadarnya;
  - §4.5 dan §7.2 — pencabutan TDA-07;
  - §4.6 — CedingID;
  - §5.6 — daftar 33 akar;
  - §7.4–7.6 — penulis `RevisionState` dan penularannya;
  - penyamaan 31 menjadi 33 akar.
- **`CONTEXT.md` induk** §1.1 dan §4 — embargo dicabut; merujuk GRL-01; ADR-0048 "sedang direkonsiliasi".

---

## 13. Kesalahan yang dinyatakan terang

### 13.1 Pembedah (MA)

| # | Kesalahan |
|---|---|
| MA-01 | precondition dibaca dengan perkakas yang tidak membaca sarangnya → dua klaim dicabut |
| MA-02 | pemanggil dihitung sebagai rujukan tanpa memeriksa parameternya → TDA-07 dilaporkan lalu dicabut |
| MA-03 | sebab disimpulkan dari dua bukti yang berdampingan ("dua pengurai") |
| MA-04 | salinan di dalam `pyIncludedRuleXML` ikut dihitung (hitungan 70 di InputTreatyInOffer, 16/16 `pyActivity`) |
| MA-05 | nama tag Data Transform salah → nol dibaca sebagai jawaban |
| MA-06 | "temuan baru" yang ternyata sudah ada di ADR-0055, karena tidak dibaca lebih dulu |
| MA-07 | satu tingkat keterlihatan dibaca sebagai seluruh rantai — dua kali berturut-turut |
| MA-08 | daftar tag tidak lengkap (`pyUserData`) → pembalikan klaim hanya muncul di dalam diff; parameter Edit disalahbaca |
| MA-09 | rantai sebab dibangun tanpa membaca ulang temuan bernomor sendiri (bukti B2) |
| MA-10 | TA-05 tidak dijalankan untuk semua artefak induk (`BENTUK-PENULIS-PROPERTI.md` sudah ada) |

**Pola bersama:** klaim negatif ("nol", "satu-satunya", "terlihat semua orang") yang dibuat dari pemeriksaan yang lebih sempit daripada klaimnya.

### 13.2 Penggrill (PG) — dicatat di `07-AUDIT` §2a

| # | Kesalahan |
|---|---|
| PG-01 | METODE dinyatakan sudah ada di direktori kerja, padahal di Downloads; berulang dengan EXP-1 |
| PG-02 | `NILAI_SELISIH` diminta "sementara", dan "disimpan atau dihitung" diserahkan ke E |
| PG-03 | memerintahkan satu berkas per GRL sebelum memeriksa contoh GRILL-05 |
| PG-04 | nomor UA-9 bertabrakan |
| PG-05 | "tiga aturan Connect-REST" (sebenarnya satu) |
| PG-06 | definisi GRILL memasukkan keputusan wewenang (bertentangan dengan §6.6) |

### 13.3 Usulan tambahan untuk METODE (TA)

Disimpan di `07-AUDIT`. Diputuskan di akhir grilling; METODE tidak disunting sekarang.

| # | Usulan |
|---|---|
| TA-03 | tambahan untuk §2.0a |
| TA-04 | sapuan bernilai nol dikalibrasi dulu terhadap kasus positif |
| TA-05 | baca seluruh artefak induk sebelum menyatakan temuan baru |
| TA-06 | keterlihatan = perkalian semua tingkat, dengan `pyVisible` per tingkat |
| TA-07 | baca nilai parameter, bukan nama; pembalikan klaim dinyatakan di badan laporan |
| TA-08 | baca ulang temuan bernomor sendiri |

---

## 14. Daftar periksa penggrill untuk setiap jawaban AI grilling

Disarikan dari kesalahan yang benar-benar terjadi di sesi ini.

1. **Klaim negatif** ("nol", "tidak ada", "satu-satunya", "terlihat semua orang"): perkakasnya sudah dikalibrasi? Semua bentuk penulis sudah disapu? Semua tingkat keterlihatan sudah dikalikan? Bila klaim itu dipakai untuk membuang atau mengunci → §3.2.
2. **Hitungan:** mintalah daftar. Periksa apakah ada salinan `pyIncludedRuleXML` yang ikut terhitung. Periksa apakah jumlahnya menutup — contohnya 323 − 6 ≠ 318, dan 96 vs 83.
3. **Konsistensi dengan temuan bernomor AI sendiri** — misalnya B2 yang bertentangan dengan NA-02 dan §4.4.
4. **Sudah dijawab induk?** Cek PETA-SUMBER-INDUK. Bila dokumennya ber-⚠, periksa status dan premisnya.
5. **Fakta vs rancangan (§3.7):** ADR tidak boleh dipakai sebagai bukti fakta, dan fakta tidak boleh membatalkan rancangan.
6. **Label:** PELESTARIAN, PERUBAHAN (dari apa), atau BARU — dan apakah label lama perlu diperbarui karena keputusan baru.
7. **Akibat sampai habis (§4.5):** uang, migrasi, orang di hari pertama, dan siapa yang harus diberi tahu.
8. **Alasan berupa pendapat** (misalnya "pengguna tidak pernah sadar memilih") → ganti dengan bukti, atau jadikan pertanyaan ke orang.
9. **Preseden (§3.6):** GRL-06 (nilai warisan dibawa apa adanya) berlaku untuk fakta atau keputusan yang dilihat orang, tidak untuk mekanisme.
10. **Kemampuan yang mati (§3.8)** tidak boleh menjadi alasan untuk sebuah pembedaan.
11. **Satu sumber kebenaran:** turunan tidak disimpan, kecuali beralasan dan dijaga uji negatif.
12. **Zona waktu dan format** sebelum membandingkan tanggal Pega dengan tanggal Oracle.
13. **Hemat putaran:** kunci bersyarat ("bila pemeriksaan X lolos, terapkan langsung"), tetapi apa pun yang berpindah masuk atau keluar golongan GRILL tetap harus disetujui penggrill.

---

## 15. Posisi terakhir dan langkah berikutnya

**Prompt balasan terakhir** berisi jawaban C2, dua catatan kecil (96 vs 83; letak mekanisme pembekuan), dan arahan C3. Bila belum terkirim ke AI grilling, **kirim dulu**.

**Urutan grilling yang tersisa:**
1. **C3** — dengan syarat "penyesuaian premi selalu material";
2. **E3** — E3a, E3b, dan E3c;
3. **titik periksa** sesudah butir kelima;
4. **I2** — sesudah membaca ADR-0042, 0043, dan 0054;
5. **penutupan** sesuai METODE §7.2, lalu serah terima yang memisahkan "yang masih menghalangi" dari "yang hanya mengukur kerusakan";
6. **berhenti** — to-spec hanya dimulai atas perintah pengguna.

**Pekerjaan pengguna yang bisa berjalan paralel:**
- kirim DB-1 sampai DB-16 (dan DB EGNPI) ke pemilik proses dan bagian teknik treaty — terutama **DB-3, DB-4, dan DB-15**, pembatal keputusan yang sudah terkunci;
- jalankan `KUERI-UA.sql` di replika, dan tanyakan zona waktu server Oracle ke DBA;
- cek UA-15 di tabel operator Pega;
- jawab **PP-1**: siapa yang mengerjakan to-spec induk. Ini menentukan kapan paket REV diserahkan;
- cek apakah ada versi properti `EDMState` atau `EDMMaterialType` di atas 01-01-56.

**Prasyarat masuk to-spec:**
- semua butir GRILL terkunci;
- setiap butir lain punya golongan dan tempat;
- paket REV diserahkan dan ditanggapi;
- `CABANG-K-PEMETAAN-TO-SPEC.md` lengkap.

GRL-01 menetapkan bahwa to-spec Adjustment ditulis **ke dalam** spesifikasi induk (langkah 8–10 SPEC-MODEL-DATA, SPEC-INVARIAN, peta telusur, ERD, STRUKTUR-DATA).

---

## 16. Glosarium singkat

| Istilah | Arti di sesi ini |
|---|---|
| kontrak / versi / addendum | addendum adalah `VERSI_KONTRAK`; versi pertama adalah kontrak asal |
| versi berlaku | versi DISETUJUI dengan `NOMOR_URUT_VERSI` tertinggi (turunan) |
| versi dasar | `ID_VERSI_KONTRAK_DASAR` — versi berlaku terakhir saat versi baru dibuat |
| material | versi yang memiliki sedikitnya satu baris `NILAI_SELISIH` (turunan, untuk versi baru) |
| lapisan beku | cedant, source of business, ProportionType, tanggal mulai, tanggal berakhir — tidak boleh berubah antar versi |
| revisi-di-tempat | tombol Revision di layar penawaran Treaty In; di model baru melahirkan versi baru dari DRAFT (ADR-0055) |
| warisan | baris yang dimigrasikan dari sistem lama; nilainya dibawa apa adanya dengan asal-usulnya bila ia keputusan atau fakta yang dilihat orang |
| irisan / 56 | berkas yang sama dengan ekspor Treaty In (temuannya milik Treaty In) / permukaan khas Adjustment |
| GRL · NA/NB · MA · PG · TA · DB · UA · EXP · REV · IND | keputusan · temuan ronde A/B · kesalahan metodologis pembedah · kesalahan penggrill · usulan tambahan METODE · daftar untuk dibantah · kueri data · ekspor tambahan · usulan revisi ADR · temuan untuk modul induk |
