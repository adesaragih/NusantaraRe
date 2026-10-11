# 59: Underwriter mengubah Spreading (When IsUWSpread)

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 07-10-2026 (gambar layar tab
> Spreading) "di pega ada when IsUWSpread, itu kondisi untuk ubah spreading ini, jadi di uw bisa ubah % Share RNM dan add
> delete ubah spreading". Survei korpus oleh subagen (pengurai XML expat, atribut dipasangkan per `rowdata`); butir kunci
> diperiksa ulang sesi 0f (isi `When\IsUWSpread.xml`, baris pemakaian, pesan `cekSpreadingFactIn`).

**Status:** frontend selesai 08-10-2026 (uji hijau); backend diserahkan ke sesi c3.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\` (hanya varian NB FacIn)

- `When\IsUWSpread.xml` (`ASM-FW-GISFW-WORK!ISUWSPREAD`; register `_oq011` no. 317: 3 modul, 2 versi isi): `pyLogic`
  `A OR B OR C OR D OR E OR F` atas `pyWorkPage.PositionNote` = `ReasFacInUnderwriting` / `ReasFacInSeniorUnderwriting` /
  `ReasFacInUnderwritingFinancial` / `ReasFacInJuniorUnderwriting` / `ReasFacInDepHeadUnderwriting` /
  `ReasFacInJuniorUnderwritingA`. Yang diuji hanya PositionNote.
- `Section\InputInwardFacultativeDtl_IsUW.xml`, tab Spreading (Section `pyTitle=Spreading`):
  - `.OfferFacIn.PercentShare` "% Share RNM" (baris ~28020) `pyReadOnlyCondition=!IsUWSpread`; change ->
    `CountPremiAndTSINusantaraRe_ACT` `InSpreading=Spreading` (= `hitung-share`, sama dengan marketing);
  - layout "Copy Spreading" (baris ~29414) `pyContainerVisibleWhen=IsUWSpread`;
  - tombol Save (baris ~37811) `pyCondition=IsUWSpread` -> `IsThereAnyObjectLocation_Act` `IsSave=Save` (activity itu tanpa
    cabang posisi UW).
- `Section\InputDtlSpreadingCoverage_FacIn_IsUW.xml` cabang IsUWSpread (baris ~4247), grid `.SpreadingList` per coverage:
  Add (`addRow`), hapus (ikon trash; `deleteRow` bila `StatusBusiness != 3`), `.TreatyType` (dropdown `SpreadingList1` dari
  `GetTreatyName`) dan `.SharePercentage` dapat diubah; `.TSISpreaded` / `.PremiumSpreaded` hanya untuk satu user ID
  ter-hardcode (tidak disalin). Cabang `!IsUWSpread` (baris ~645) grid baca-saja tanpa tombol.
- Change `.SharePercentage` -> preDT `SetTSIPremiSpreaded_FacIn` Type=percent (TSISpreaded = share × TSINusantaraRe / 100,
  PremiumSpreaded = share × PremiNusantaraRe / 100, skala 20) + `cekSpreadingFactIn` (TSIGrossSpreaded; pesan
  "Share Percentage can't be more than 100"; "Invalid Treaty Type Spreading!" dilewati hanya untuk Marketing).
- `Activity\CountPremiAndTSINusantaraRe_ACT.xml` langkah 11 "WHEN UW" (When IsUWSpread T=2 F=3): `CountPremiCedant_Act`,
  `SetRIComm_Act`, `CountPaymentInstallment_Act`.

## Keputusan

- **K59-1** `[ADR-0032 + keputusan work owner 07-10-2026]` `ReasFacInJuniorUnderwriting` dibaca `ReasFacInJuniorUnderwritingB`
  (migrasi 211).
- **K59-2** `[keputusan agent]` Tahap UW HANYA membuka tab Spreading (yang diminta work owner); tab lain, General, Save for
  later, dan Submit tetap baca-saja / 409. Pemegang = workbasket POSISI saat itu (pola Team Leader tiket 58).
- **K59-6** `[keputusan work owner 08-10-2026]` "spreading hanya bisa di copy atau add yg masuk ke when IsUW, di marketing
  tidak bisa add spreading": Copy Spreading dan Add / ubah baris HANYA tahap UW; marketing / Team Leader hanya spreading
  otomatis dari % Share RNM (meralat K48 S-1). "when IsUW" dibaca sebagai tahap UW = When `IsUWSpread` `[dugaan agent]`:
  When `IsUW` harfiah (workbasket pertama operator ∈ ReasFacInDirector / ReasFacInGroupLeader / ReasFacInUnderwriting) di
  Pega justru me-RO isian Copy Spreading. Frontend: blok Copy Spreading dirender hanya `tahapUw`; backend diminta menolak
  `salin` / `hitung-baris` di luar UW dan tidak menerima baris badan PUT di marketing / Team Leader.
- **K59-7** `[keputusan work owner 08-10-2026]` "spreading terbentuk otomatis juga hanya terbentuk di uw. dan di when jika
  ada workbasket ReasFacInJuniorUnderwriting diubah jadi ReasFacInJuniorUnderwritingB". `[terverifikasi]` korpus sejalan:
  `Activity\CountPremiAndTSINusantaraRe_ACT.xml` langkah 4 `Call GetKapasitasTreaty` berprekondisi `Param.spread=="Spreading"`
  (T=2 F=3); param `spread` hanya dikirim `Section\InputInwardFacultativeDtl_IsUW.xml` (baris ~27622 / ~27916), section marketing
  `InputInwardFacultativeDtl.xml` hanya `InSpreading` (grep `<pyName>spread</pyName>` nihil). Di marketing / Team Leader:
  hitung-share dan Save hanya menghitung TSI / premi Nusantara Re; baris tersimpan dipertahankan dan dihitung ulang. When:
  di `NB FacIn\When` hanya `IsUWSpread.xml` memuat nama lama (satu berkas), sudah dipetakan `PetaWorkbasket` di
  `services/rules/registry_gen.go`; backend diminta memakai When registry itu sebagai sumber tahap uw.
- **K59-3** `[keputusan agent]` (sebagian diganti K59-6) Copy Spreading tetap dapat dipakai di UW (K48 S-1). Di Pega isian grid itu ber-RO When `IsUW`
  (workbasket pertama operator ∈ ReasFacInDirector / ReasFacInGroupLeader / ReasFacInUnderwriting) dan tombol Add / Delete /
  Copy To All ber-`pyCondition=1=2` dengan `pyVisible=ALWAYS` — apakah tombol itu tampil **belum terverifikasi**.
- **K59-4** `[keputusan agent]` Label tombol grid coverage "Add" / "Delete" = label Copy Spreading (Pega grid UW memakai
  tombol gambar / ikon trash, label teks tidak ada di korpus).
- **K59-5** `[keputusan agent]` Hitung ulang baris lewat endpoint baru tanpa menyimpan (`hitung-baris`, jeda 500 ms), dipanggil
  hanya bila semua baris lengkap dan Σ % Share tiap coverage <= 100; Save ditahan selama tidak sah.

## Frontend (sesi 0f)

- `api.ts`: `KasusNB.bolehUbahSpreading`, `KasusNB.ubahBarisSpreading`, `hitungBarisSpreading(caseId, isi)`.
- `components/TabSpreading.tsx`: prop `ubahBaris`; `barisSpreadingKosong`, `gantiBarisCoverage`, `barisLebih`,
  `barisLengkap`, `barisSah`.
- `pages/InwardFacultative.tsx`: `spreadingTerbuka` / `detailBacaSaja` — pembungkus detail tidak aktif hanya saat tab
  Spreading terbuka di tahap UW; banner `TEKS_INWARD.bacaSajaKecualiSpreading`.
- `labels.ts`: `TEKS_SPREADING.shareBarisLebih` (verbatim korpus), `barisBelumLengkap`, `TEKS_INWARD.bacaSajaKecualiSpreading`.
- Uji: `TabSpreading.test.tsx`, `InwardFacultative.test.tsx`, `BlokEmailType.test.tsx`.

## Backend (diminta ke sesi c3, 08-10-2026)

1. Tahap `uw` di `services/bacasaja.go` (posisi K59-1, tanpa syarat flag, workbasket = POSITION).
2. Penjaga PUT `…/spreading`, POST `…/hitung-share`, `…/salin`, `…/hitung-baris`: marketing / Team Leader ATAU uw.
3. GET kasus: `bolehUbahSpreading`, `ubahBarisSpreading` (bool, selalu ada).
4. `POST /api/nbfacin/kasus/{caseId}/spreading/hitung-baris` badan `{percentShare, lokasi}` = hitungan Save tanpa menulis.
5. "Invalid Treaty Type Spreading!" berlaku di tahap uw.
6. Verifikasi tab Cedant / Deduction / Payment mencerminkan share baru sesudah Save di uw (langkah 11 CountPremiAndTSI).

## K59-8 — Save tahap UW menolak "Invalid Treaty Type Spreading!" (penyimpangan sadar)

- `[keputusan work owner 08-10-2026 "setuju untuk dibuat", atas rekomendasi sesi 0f]` (pertanyaan pemicu: "km
  rekomendasikan dibuatin juga atau gimana")
  PUT `…/spreading` di tahap UW menolak 400 dengan pesan verbatim "Invalid Treaty Type Spreading!" dan tidak menulis apa pun.
- Penyimpangan dari Pega: Save (`IsThereAnyObjectLocation_Act`) tidak memanggil `cekSpreadingFactIn`. Alasannya, di Pega
  onchange grid UW selalu menjalankan `cekSpreadingFactIn` langkah 8-9 (Page-Set-Messages pada work page) sebelum Save dapat
  ditekan. Di aplikasi ini hitung-baris berjeda 500 ms, sehingga tanpa penjaga Save cepat dapat menyimpan treaty tidak valid.
- Marketing / Team Leader tidak terdampak (baris badan diabaikan). hitung-baris / salin tetap non-blocking (`pesan`).
- Frontend tidak berubah: galat 400 tampil lewat `Gagal` di tab Spreading.

## Terbuka

- K59-6: tafsir "when IsUW" = tahap IsUWSpread — konfirmasi work owner bila yang dimaksud When `IsUW` harfiah.
- Hapus baris case endorsement (`StatusBusiness == 3` -> `SetFlagDeleteEDM_Act`) di luar cakupan NB.

## Backend sesi c3 — selesai 08-10-2026

Termasuk dua tambahan work owner 08-10-2026 (lewat sesi 9d): "spreading hanya bisa di copy atau add yg masuk ke when IsUW, di
marketing tidak bisa add spreading" dan "spreading terbentuk otomatis juga hanya terbentuk di uw. dan di when jika ada workbasket
ReasFacInJuniorUnderwriting diubah jadi ReasFacInJuniorUnderwritingB".

- **K59-B1 Tahap uw** (`services/bacasaja.go`): `tahapUW` = `rules.Eval("ISUWSPREAD", …)` atas registry bangkitan (bukan daftar
  tulisan tangan; kondisi D sudah ReasFacInJuniorUnderwritingB lewat `bangkit.PetaWorkbasket`, ADR-0032). Hanya PositionNote,
  tanpa flag. Workbasket wajib = POSITION (403). `wajibBolehUbah` tab lain TIDAK diubah: tetap 409 di uw; `bolehUbah` /
  `bolehSubmit` tetap false.
- **K59-B2 Penjaga per endpoint:**

  | Endpoint | marketing / TL (+workbasket) | uw (+workbasket posisi) | uw tanpa workbasket | tahap lain |
  | --- | --- | --- | --- | --- |
  | PUT `…/spreading` | 200, baris badan DIABAIKAN | 200, baris badan dipakai | 403 | 409 |
  | POST `…/salin` | 409 `ErrSpreadingHanyaUW` | 200 | 403 | 409 |
  | POST `…/hitung-baris` (baru) | 409 `ErrSpreadingHanyaUW` | 200, tidak menulis | 403 | 409 |
  | POST `…/hitung-share` | 200, TANPA otomatis | 200, spreading otomatis | 200 (tidak menulis) | 200 tanpa otomatis |

  `hitung-share` sejak tiket 48 tidak dijaga (tidak menulis); dibiarkan begitu. `salin` dulu tidak dijaga - kini 401 / 403 / 409.
- **K59-B3 PUT di marketing / TL:** baris badan diabaikan; baris TERSIMPAN dipertahankan (jumlah, jenis, share) dan nilainya
  dihitung ulang dari % Share badan (hitungNR + DT percent + total). Case baru = tanpa baris. Tidak ada spreading otomatis di
  PUT mana pun. Sebelum tiket 59 baris PUT marketing berasal dari BADAN (hasil hitung-share / salin di layar).
- **K59-B4 Spreading otomatis hanya uw** `[terverifikasi]`: `Activity\CountPremiAndTSINusantaraRe_ACT.xml` langkah 4
  `Call GetKapasitasTreaty` berprekondisi `Param.spread=="Spreading"`; `<pyName>spread</pyName>` 2 kali di
  `Section\InputInwardFacultativeDtl_IsUW.xml`, 0 di `Section\InputInwardFacultativeDtl.xml`.
- **K59-B5 "Invalid Treaty Type Spreading!"** `[terverifikasi]` `Activity\cekSpreadingFactIn.xml`: onchange baris grid UW
  (`Section\InputDtlSpreadingCoverage_FacIn_IsUW.xml` ±6364 / ±6548) memanggilnya tanpa `IsPct`, jadi langkah 7 (Sum) dan 8-9
  berjalan; langkah 8.6 / 9 dilewati hanya bila PositionNote ReasFacInMarketing atau `FlagViewOffer == "1"` (tidak disimpan
  sistem baru, `belum terverifikasi`). Kini: `hitung-baris` (uw) dan `salin` (uw) memeriksanya; PUT tidak (Save
  `IsThereAnyObjectLocation_Act` tidak memanggil cekSpreadingFactIn, K48). Pesannya di `pesan`, tidak memblokir. Langkah 4.5 / 5
  "Share Percentage can't be more than 100" -> 400 seperti PUT.
- **K59-B6 Tab lain sesudah PUT spreading di uw** (langkah 11 `IsUWSpread`: CountPremiCedant_Act, SetRIComm_Act,
  CountPaymentInstallment_Act `[terverifikasi]`):
  - TAMPILAN mengikuti share baru. Cedant (TotalTSIRNM / TotalPremiRNM / dasar per mata uang) dihitung saat baca dari tabel
    spreading. Deduction memakai Commision / BrokerageFee / PPh / PPN tersimpan hanya bila premi sama; selain itu dihitung ulang
    dari persen tersimpan (K50-3). Payment menghitung cicilan dari persen.
  - ⚠️ NILAI TERSIMPAN tidak ditulis ulang. Endpoint tab itu 409 di uw dan TIDAK dibuka:
    - `T_CURRENCYLIST` PAY_PREMIUM / PAY_COMMISION / … dan `PAY_DEDUCTION2`. Deduction2 tampil = tersimpan (nominal lama, tidak
      ikut share);
    - `T_CURRENCYLIST.SUM_TOTAL_PAYMENT` / `PAY_NET_PREMIUM` dan nominal `T_LISTINSTALLMENT`;
    - `T_QUOTATIONDATA.SHARE_OF_CEDING` ("<share>%"). Di Pega pun hanya diset klik ShareCedantType (`SetShareOfCeding`;
      CountPremiCedant_Act / SetRIComm_Act / CountPaymentInstallment_Act tidak menyentuhnya).

    Pembaca DB lain melihat nilai lama sampai tab itu di-Save lagi. Usul: port langkah 11 ke dalam transaksi PUT spreading uw.
    `[pertanyaan terbuka]` work owner.
- **K59-B7 GET kasus:** `bolehUbahSpreading`, `ubahBarisSpreading` (bool, selalu ada; lunak: galat baca tahap = false; tanpa
  pembaca tahap = true).
- **K59-B8 When nama lama:** korpus `NB FacIn\When` + `DecisionTable` / `DecisionTree` / `DataTransform` dengan nama lama hanya
  `When\IsUWSpread.xml` (terpetakan) dan `DataTransform\IsJuniorUW_B.xml` (diport tangan: `posisiJuniorUWB`). Registry bangkitan
  hanya berisi When. Uji penjaga `rules.TestRegistryTanpaWorkbasketJUWLama` (dua cara: struktur terurai + teks berkas; tepat
  satu baris kode B). endorsmentfacin di luar modul - tidak diubah.
- Tanpa migrasi. Uji: `services/spreadinguw_test.go`, `handlers/spreadinguw_test.go`, `rules/workbasket_test.go`,
  `TestPenjagaBacaSaja` (spreading di uw dikecualikan, pindah ke matriks), `TestBacaKasus` (13 kunci),
  `TestKunciKasusSamaDenganApiTs` (dua arah atas api.ts KasusNB).

> **RALAT K59-B6 (08-10-2026, permintaan sesi 9d):** "NILAI TERSIMPAN tidak ditulis ulang" tidak berlaku lagi untuk Premium
> Deduction / Payment — lihat K59-B9. SHARE_OF_CEDING tetap tidak disentuh.

- **K59-B9 Langkah 11 CountPremiAndTSI di PUT uw** (`services/spreadinguw.go`). Port `[terverifikasi]`
  `Activity\CountPremiAndTSINusantaraRe_ACT.xml` langkah 11 "WHEN UW". Berlaku HANYA di tahap uw nyata, tidak saat penjaga
  tidak aktif. Ditulis di transaksi yang SAMA dengan TulisSpreading; bila gagal, seluruh Save batal.
  - CountPremiCedant_Act: `CedingCedantList(i).CurrencyList` tidak disimpan sistem baru (T_CEDINGCEDANTLIST tanpa kolom per mata
    uang, migrasi 199), jadi tidak ada yang ditulis.
  - SetRIComm_Act memakai `hitungRIComm` tiket 50:
    - persen TERSIMPAN dibangun ulang ke presisi penuh (`persenPenuh`, K50-9) dan dikalikan premi Nusantara Re baru;
    - langkah 5 `PPnCheck = true` SELALU, sehingga PPN_CHECK ditulis "true";
    - Deduction2 = nominal tersimpan per mata uang (SetRIComm tidak menyentuhnya);
    - cabang EDM (1, 6, 8.2) di luar cakupan.
  - CountPaymentInstallment_Act memakai `segarkanCicilan` K50-12:
    - % dan Due Date tetap;
    - nilai uang / Total Payment / Net Premium dihitung ulang;
    - Net mengurangkan Deduction2 (K51-1, penyimpangan sadar yang sudah ada);
    - Net negatif (K51-7) membuat cicilan tidak ditulis.
  - Bila tab Premium Deduction belum pernah di-Save, tidak ada yang ditulis (tidak ada nilai basi; GET menghitung dari
    spreading).
  - ⚠️ Share turun dapat membuat Deduction2 > Premium. SetRIComm_Act tidak memeriksanya (ProtekDeduction2 hanya di aksi
    Deduction 2). GET tab memberi pesan "Deduction 2 cannot be greater than the Premium"; tab itu baca-saja di uw.
  - Tabel yang ditulis = tabel Save Premium Deduction jalur marketing (T_GENERAL_POLIS PAY_* / PPN_CHECK, T_CURRENCYLIST,
    T_LISTINSTALLMENT, PAY_INSTALLMENT). Tanpa migrasi. Endpoint tab lain tetap 409; hitung-baris / hitung-share tetap tidak
    menulis.
- **K59-8 PUT uw menolak jenis treaty tidak sah** — penyimpangan sadar `[keputusan work owner 08-10-2026, atas rekomendasi
  sesi 0f]` (verbatim: "setuju untuk dibuat").
  - Di tahap uw, sebelum menulis apa pun (termasuk langkah 11), `jenisTreatyKasus` menjalankan `periksaJenisTreaty`. Aturan
    yang sama dipakai hitung-baris / salin.
  - Bila tidak sah: 400 `ErrMasukanSpreading` "…: Invalid Treaty Type Spreading!".
  - Pega memeriksanya di onchange grid sebelum Save dapat ditekan; di sini hitung-baris berjeda 500 ms, jadi Save memeriksa
    ulang. Save Pega sendiri tidak memanggil cekSpreadingFactIn (K48).
  - Marketing / TL tidak diperiksa (baris badan diabaikan). hitung-baris / salin tetap non-blocking.
- Uji tambahan:
  - `TestSimpanSpreadingUWLangkah11` (nilai tertulis; satu transaksi; gagal Payment = satu rollback);
  - `TestSimpanSpreadingLangkah11HanyaUW` (marketing / TL / penjaga tidak aktif / uw tanpa Premium Deduction tersimpan: hanya
    spreading);
  - `TestSimpanSpreadingUWJenisTreatyTidakSah`;
  - `TestPutSpreadingUWJenisTreatyHandler`.
- **K59-9 PUT uw menolak Deduction 2 > Premium** — penyimpangan sadar `[keputusan work owner 08-10-2026 "setuju jangan halu",
  atas rekomendasi sesi 0f]`.
  - Kapan: sesudah langkah 11 menghitung CurrencyList baru (`barisLangkah11`) dan SEBELUM transaksi.
  - Pembanding = `pesanDeduction2` yang sudah ada (tiket 50; pesan verbatim `[terverifikasi]` `Activity\ProtekDeduction2_Act.xml`
    langkah 2): per mata uang, Deduction2 tersimpan > Premium. Premium = Σ premi Nusantara Re mata uang itu, 8 desimal, nilai
    yang sama dengan yang ditulis langkah 11, jadi basis pembanding sama.
  - Hasil: 400 `ErrMasukanSpreading` "…: Deduction 2 cannot be greater than the Premium", tanpa tulisan apa pun.
  - Tidak berlaku bila Premium Deduction belum pernah di-Save. Marketing / TL tidak berubah.
  - Alasan menyimpang: SetRIComm_Act Pega tidak memeriksanya, sementara di uw tab Deduction baca-saja sehingga UW tidak dapat
    membetulkannya.
  - ⚠️ Batas Deduction 2 = Premium lolos, tetapi Net Premium dapat negatif. Menurut K51-7 cicilan Payment lalu TIDAK ditulis
    ulang (tetap nilai lama); GET Payment menghitung ulang. Uji: `TestSimpanSpreadingUWDeduction2LebihPremi`.
- **K59-10 PUT uw menolak Net Premium < 0** — penyimpangan sadar `[keputusan work owner 08-10-2026 "setuju jangan halu", atas
  rekomendasi sesi 0f]`. CountPaymentInstallment_Act Pega tidak menolak.
  - Kapan: sesudah K59-9, sebelum transaksi.
  - Pemeriksaan: untuk setiap mata uang `barisLangkah11`, `netPremium(b) < 0` memakai fungsi yang SAMA dengan `isiCicilan` /
    `segarkanCicilan` dan Save Payment. Basisnya sama: baris 8 desimal yang ditulis langkah 11. Net mengurangkan Deduction2
    (K51-1).
  - Hasil: 400 `ErrMasukanSpreading` dengan pesan K51-7 (pesan agent, bukan korpus) dari fungsi bersama `pesanNetNegatif`
    (payment.go). Satu pesan per mata uang negatif, digabung "; ". Tidak ada tulisan.
  - Net = 0 lolos.
  - Net tidak bergantung pada cicilan: mata uang tanpa Installment / cicilan tersimpan TETAP diperiksa (Net dihitung dari
    nilai CurrencyList saja). Bila hanya mata uang bercicilan yang perlu diperiksa, perlu keputusan.
  - Tidak berlaku bila Premium Deduction belum pernah di-Save. Marketing / TL tidak berubah.
  - Uji: `TestSimpanSpreadingUWNetNegatif` (teks K51-7 dipaku).
- **K59-B11 Bug DEV 08-10-2026: PUT uw -> 500 ORA-01722** (log work owner: "repository: menulis TSI / premi Nusantara Re:
  ORA-01722: invalid number").
  - Sebab: `ringkasanRNM(k)` (by value) memanggil `hitungNR(&k)` atas salinan DANGKAL. Backing array Lokasi / Items / Coverages
    dibagi dengan pemanggil, sehingga TSI / premi NR pemanggil yang sudah dibulatkan (hitungSaveSpreading) ditimpa nilai skala
    20. TulisSpreading lalu mengikat teks 20 desimal ke mask `FM…D99999999`. Jalurnya: barisLangkah11 -> premiPerMataUang,
    hanya di uw.
  - Perbaikan di sumber: `salinDalamSpreading` (Lokasi, Total, Items, Coverages, baris) dipakai sebelum mengubah di
    `ringkasanRNM` dan `tampilan`.
  - Fungsi by value yang MENGUBAH masukan:
    - `ringkasanRNM` (dan karenanya `premiPerMataUang`) - bug;
    - `tampilan` - laten: pemanggilnya tidak memakai k lagi, tetapi ikut diperbaiki.
  - Fungsi by value yang hanya membaca: `totalMataUang`, `ringkasan`, `treatyKasus`, `jenisTreatyKasus`, `waktuMulai`,
    `tanggalMulai`, `adaIsiSpreading`, `pertahankanSpreading` (`lama` hanya dibaca; baris disalin `salinBaris`). Aritmetika
    desimal paket selalu membuat nilai baru (tambah / kurang / kali / bagi / simpan8).
  - Penjaga: `repository.PolaAngkaMasuk` (pola teks bind mask angkaMasuk; `TestPolaAngkaMasuk` dua cara). Tiruan
    `TulisSpreading` / `TulisDeduction` / `TulisPayment` (services) dan `TulisSpreading` (handlers) menolak bind di luar mask.
    Dengan perbaikan dibalik, 5 uji uw merah ("250.00000000000000000000").
  - Uji: `TestRingkasanTidakMengubahMasukan` (angka DEV: NR 50000000000 / 7962413.4520277 tetap sesudah premiPerMataUang /
    ringkasanRNM / barisLangkah11 / totalMataUang) dan `TestTampilanTidakMengubahMasukan`. Tanpa migrasi.
- **K59-B12 `POST …/spreading/periksa-template`** — permintaan work owner 08-10-2026: "INI KENAPA BISA PILIH SPL TANPA ADA
  QS, DIPEGA ADA PROTEKSI NYA TIDAK BISA SPL TANPA ADA QS".
  - `[terverifikasi]` sesi 0f: `Activity\CalcultePersentageSpeading_Act.xml` dipanggil onchange `.TreatyType` grid Copy
    Spreading. Langkah 2.13–2.16 memberi "Treaty Type can't be same"; 2.17 memberi "Can not proceed spreading without QS"
    (baris pertama, NOTE REINSURANCETYPE tanpa "QS", kecuali ORS 10007 / FACOUT 10015). Pesan muncul SEKETIKA saat treaty
    dipilih.
  - Badan `{template}`, jawaban `{pesan: string[]}` (200), tidak menyimpan.
  - Isi = `proteksiTemplate` yang SAMA dengan `…/salin` (satu fungsi): pesan verbatim, urutan sama.
  - Baris ber-treatyType kosong diabaikan. Catatan: `…/salin` justru menolaknya 400 lewat `periksaBarisShare`, tidak
    mengabaikannya.
  - > 50 baris -> 400; badan rusak / tanpa `template` -> 400.
  - Penjaga = salin: 401 -> tahap uw (409 `ErrSpreadingHanyaUW` untuk marketing / TL / lainnya) -> 403 workbasket posisi.
  - `…/salin` tetap menolak (lapis kedua).
  - Uji: `TestPeriksaTemplateSpreading` (QS dulu / SPL saja / ORS / FACOUT / kembar / kosong / penjaga) dan
    `TestPeriksaTemplateSpreadingHandler` (HTTP, jawaban persis `{"pesan":[…]}`).
