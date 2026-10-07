# 45: Tab Coverage FIRE — tahap C3 Deductible

> ⚠️ **Disusun agent atas perintah sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Frontend dibangun sesi 0f;
> backend (simpan / baca) diserahkan ke sesi c3.

**What to build:** setiap coverage di tab Coverage FIRE memuat grid **Deductible** (tambah / hapus / form lipat), disimpan
bersama Save tab Coverage.

**Status:** frontend dan backend selesai 03-10-2026 (uji hijau). ⛔ Butuh migrasi **194** dijalankan sesudah 193.

> ⚠️ RALAT (sesi 0f): versi awal baris ini menyebut `PUT …/objek` memakai `DisallowUnknownFields` sehingga
> `deductibles` ditolak 400. Keliru — `simpanObjek` (`handlers/objek.go`) memakai decoder biasa; `DisallowUnknownFields`
> hanya di fungsi `urai` (`handlers.go`) untuk rute lain. Sebelum backend tiket ini, `deductibles` dibuang diam-diam.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\` dan `D:\migrasi\RNM\DDL\`

- `Section\CoverageItem.xml`: grid PageList `.DeductibleList` berjudul "Deductible" (tampil `!isClaim`); kepala kolom kosong
  kecuali "Time Excess (Days)"; Add = addRow + `addDeductible_ACT`, Delete = deleteRow; form = flow action
  `InputDtlDeductibleFire_FacIn` → section `addDeductible.xml`.
- `Section\addDeductible.xml` sel 5–14: Type Deductible (selalu); Pct Deductible (`.TypeDeductible != 7`); MinMax
  (`!= 7`, baca-saja bila Type 0 / 7); Currency (baca-saja bila Type 0); Type Deductible / Pct Deductible kedua
  (`.MinMax = 3`); Condition; Amount (selalu); Condition teks (`.Condition = 5`); Time Excess (In Days). Tidak ada medan
  wajib; teks pilihan kosong "Choose".
- Daftar pilihan = aturan properti `ASM-FW-GISFW-DATA-DEDUCTIBLE!TYPEDEDUCTIBLE`, `!TYPEDEDUCTIBLE2`, `!MINMAX`,
  `!CONDITION` (`DDL\`, dicari menurut `<pxInsName>` — lihat catatan Condition di bawah).
- Sampel kasus `DDL\P-5 *.txt` (kelas `ASM-FW-GISFW-Data-Deductible`) memuat properti `TypeDeductible`, `PctDeductible`,
  `MinMax`, `Currency`, `TypeDeductible2`, `PctDeductible2`, `Condition`, `InputCondition`, `Amount`, `TimeExcess`,
  `FlagCurrency`, `IndexProperty` / `IndexPropertyItem` / `IndexCoverage` / `IndexDeductible`. Properti kosong tidak
  diekspor (mis. satu deductible MinMax 3 tanpa `Currency` / `Amount`).
- `DDL\Condition.xml` kini berisi `DATA-DEDUCTIBLE!CONDITION`; `DATA-PROPERTYITEM!CONDITION` (Condition Object Item,
  tiket 39) dikirim ulang work owner sebagai `DDL\ConditionObjectItem.xml` (03-10-2026). Uji label mencari keduanya
  menurut `<pxInsName>`, bukan nama berkas.

## Frontend (sesi 0f)

- `components/GridDeductible.tsx` (`deductibleBaru`, `adaGalatDeductible`), dipasang di bawah `FormCoverage`; mata uang
  dari `daftarMataUang`. `TabCoverage` menolak Save bila ada angka deductible tak sah.
- `api.ts`: `CoverageObjek.deductibles?: Deductible[]`; `Deductible` = sepuluh medan teks (kode / desimal).
- `labels.ts`: `JUDUL_DEDUCTIBLE`, `FORM_DEDUCTIBLE`, `KOLOM_TIME_EXCESS`, `DEDUCTIBLE_KOSONG`, `OPSI_TYPE_DEDUCTIBLE(2)`,
  `OPSI_MINMAX`, `OPSI_KONDISI_DEDUCTIBLE`.
- Uji: `labels.test.ts` (sel 5–14, syarat tampil, grid CoverageItem, empat PromptList menurut pxInsName),
  `components/GridDeductible.test.tsx`.

## Kontrak backend (untuk sesi c3)

- `GET` / `PUT …/objek`: setiap coverage membawa `deductibles: Deductible[]` (urutan = urutan grid; `[]` bila kosong).
  Medan teks: `typeDeductible`, `minMax`, `currency`, `typeDeductible2`, `condition`, `inputCondition`; desimal (teks,
  ADR-0034): `pctDeductible`, `pctDeductible2`, `amount`, `timeExcess`. Kosong = NULL.
- 400 bila desimal tak sah (`…deductibles[d].<medan>`). Tidak ada medan wajib (sesuai addDeductible).
- Penyimpanan: tabel anak coverage (rancangan tabel diserahkan c3, uang NUMBER(38,8)); migrasi berikutnya setelah 193.

## Keputusan agent (menunggu konfirmasi)

- **R-1** Kepala kolom grid memakai label form (kepala Pega kosong kecuali "Time Excess (Days)").
- **R-2** Deductible baru bermata uang item (`addDeductible_ACT` membaca `Local.Curr`) `[dugaan]`.
- **R-3** Ditunda ke tahap berikut: penanda beda mata uang (`Activity\SaveDeductible.xml` dan `addDeductible_ACT.xml`
  merujuk `IsBedaCurr`; `SaveDeductible` mengisi `FlagCurrency` `[terverifikasi]`; kaitannya dengan flow action
  `FlowAction\ProtectCurrency.xml` `[dugaan]`), Descriptions, View Old Deductible, Copy Deductible To All / For BI.
- **R-4** Deductible tersimpan lewat Save tab Coverage, bukan tombol Save tersendiri per form.

## Acceptance criteria

- [x] Grid Deductible per coverage: Tambah membuka form, Hapus menghapus baris, baris memperlihatkan label pilihan.
- [x] Syarat tampil / baca-saja mengikuti addDeductible; Amount berformat ribuan.
- [x] Angka tak sah memblokir Save.
- [x] Backend menyimpan dan mengembalikan `deductibles` (sesi c3, 03-10-2026).
- [ ] Migrasi 194 dijalankan di DEV (work owner).

## Backend (sesi c3, 03-10-2026)

⛔ **Migrasi 194 WAJIB dijalankan di DEV sebelum backend baru** — urutan … → 193 → **194**. `GET` / `PUT …/objek` membaca
dan menulis `T_DEDUCTIBLELIST`; tanpa 194 seluruh tab Object gagal ORA-00942.

- `GET` / `PUT …/objek`: `coverages[k].deductibles` (selalu larik ke luar, 10 kunci persis kontrak `Deductible`; boleh
  absen saat masuk). Urutan = SEQ_NO. Kosong → NULL. Coverage yang dihitung ulang (tiket 43 / 44) membawa deductible-nya
  apa adanya. `POST hitung-coverage` / `hitung-net-rate` juga meneruskan `deductibles` (tidak disimpan).
- 400 ber-jalur `baris[n].items[m].coverages[k].deductibles[d].<medan>`: desimal tak sah (`pctDeductible`,
  `pctDeductible2`, `amount`, `timeExcess`); `typeDeductible` / `typeDeductible2` bukan bilangan bulat ≤ 5 digit; lebar
  `minMax` / `currency` 50, `condition` / `inputCondition` 500 byte; `timeExcess` > 30 karakter; `currency` terisi
  tetapi tidak ada di `CURRENCY` (≠ ITL, pola A133). Tanpa medan wajib (mata uang pun boleh kosong, A169). Paling banyak 500 deductible per
  coverage. Keanggotaan kode pilihan tidak diperiksa (pola butir 85).
- ⚠️ Koreksi kontrak `[terverifikasi]`: `PUT …/objek` (`handlers/objek.go` `simpanObjek`) TIDAK memakai
  `DisallowUnknownFields` — itu fungsi `urai` (`handlers.go:268`) untuk rute lain; medan tak dikenal di objek diabaikan
  (sebelum tiket ini `deductibles` dibuang diam-diam, bukan 400).
- Tabel: migrasi **194** `T_DEDUCTIBLELIST` sebagian (18 kolom; FK `PARENT_ID` → `T_COVERAGELIST.ID`, indeks `PARENT_ID`,
  `SEQ_T_DEDUCTIBLELIST`); dihapus sebelum coverage saat PUT. `repository/deductible.go`; daftar kolom bersama
  `repository/kolomdata.go` (T_COVERAGELIST ikut memakainya).
- Bukti sampel `[terverifikasi]` (dihitung dua cara: blok JSON terurai dan `grep -c` pxObjClass — keduanya **222**
  deductible di `DDL\P-5 *.txt`): TypeDeductible / TypeDeductible2 / MinMax / Condition bulat 1 digit; PctDeductible
  `9`, `99`, `9.9`; Amount bulat 3–8 digit; TimeExcess `14` / `30`; **12** deductible tanpa Currency — **ke-12-nya juga
  tanpa Amount**. Kode TypeDeductible = `pyStandardValue` 0–7 (`DDL\TypeDeductible.xml`; properti String/Text).
  Aturan properti `TimeExcess` tidak ada di korpus — tipe Pega `belum terverifikasi`; rancangan VARCHAR2(30).

### Keputusan agent (menunggu konfirmasi)

| # | Keputusan | Dasar |
| --- | --- | --- |
| A169 | `currency` boleh kosong (juga bila `amount` diisi). Kosong: kolom CURRENCY tidak disisipkan sehingga DEFAULT `'UNKNOWN'` rancangan yang mengisi (keadaan eksplisit K-012); dibaca kembali sebagai kosong. Terisi: wajib ada di `CURRENCY` (≠ ITL) | Rancangan `CURRENCY VARCHAR2(50) DEFAULT 'UNKNOWN' NOT NULL` (K-069, usulan 10) + K-012 (aplikasi tidak menulis `UNKNOWN`); addDeductible tanpa Required. ⚠️ Versi awal mewajibkan mata uang bila Amount diisi — **dicabut** sesudah review sumbu spec: kontrak "tidak ada medan wajib", dan dasarnya hanya korelasi sampel (12 deductible tanpa mata uang, semuanya tanpa Amount `[terverifikasi]`), bukan aturan; bisa menahan Save data lama. `[dugaan]` mata uang harfiah "UNKNOWN" di tabel CURRENCY akan terbaca kosong — tidak diperiksa |
| A170 | `TYPE_DEDUCTIBLE` / `TYPE_DEDUCTIBLE2` = **NUMBER(5)** (rancangan NUMBER polos); kabel teks, SQL `TO_NUMBER` / `TO_CHAR`; bentuk kode bulat ≤ 5 digit diperiksa | Kode `pyStandardValue` 0–7 `[terverifikasi]`; penjaga `presisiSah` melarang NUMBER polos, NUMBER(5) termasuk bentuk sah |
| A171 | `timeExcess` = teks desimal (≥ 0, ≤ 8 desimal) di kolom rancangan VARCHAR2(30), panjang ≤ 30 karakter | Kontrak (desimal); rancangan VARCHAR2(30); sampel bulat hari; tipe properti Pega tidak ada di korpus |
