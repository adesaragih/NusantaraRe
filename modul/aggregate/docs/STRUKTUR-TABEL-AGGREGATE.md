# Struktur Tabel — Aggregate

**Keputusan work owner 04-10-2026:** modul `aggregate` mengunggah data aggregate per zona penilaian ke tabel warisan
`POOLDATA.AGGREGATE` mengikuti langkah XML Pega folder korpus `Aggregate`. Modul ini **tidak membuat tabel baru** dan
tidak mengubah struktur tabel mana pun; satu-satunya objek baru adalah `SEQ_AGGREGATE` (migrasi 880).

Katalog dan profil DEV 04-10-2026 (agregat): 13.450 baris, diunggah Maret-Juni 2025, As At 31-12-2023 s.d.
31-12-2024, 129 Master Treaty, 56 ceding, 22 zona, Treaty Type OR/QS/SURPLUS, Coverage EQVET/RSMDCC/TSFWD; PK
`AGGREGATE_PK` atas `ID`; nol trigger; nol view atau prosedur yang memakainya.

## AGGREGATE

Tabel warisan Pega, terdaftar `Tabel warisan: dibaca, tidak dibuat` di `MODUL.md` (ditulis tanpa diubah strukturnya).
Kolom `CSV` = nomor kolom berkas CSV (`UploadCSVAggregate_Act` langkah 6, `.COLUMN1`-`.COLUMN38`).

| Kolom | Tipe | Null | Kunci | Dipakai | Sumber |
| --- | --- | --- | --- | --- | --- |
| `ID` | teks | tidak | PK | rincian | `AGG-` + `SEQ_AGGREGATE`; nomor terpakai dilewati (Pega: `pzGenerateUniqueID(tools, "AGG")`; DEV terbesar AGG-33686) |
| `TANGGAL_INPUT` | date | ya | | daftar "Tanggal Input" | `SYSDATE` saat Save, SATU nilai untuk seluruh unggahan (Pega: per baris) |
| `USER_INPUT` | teks | ya | | — | akun pelaku (`LOGIN_ID`) |
| `ASSESMENT_ZONE` | teks | ya | | grid | CSV 1; 3 huruf pertama dicocokkan ke `ASSESSMENT_ZONE.ASSESSMENT_CODE` -> `ASSESSMENT_NOTE`; tidak ada = kosong (Save menolak) |
| `TREATY_TYPE` | teks | ya | | grid, daftar, chart | CSV 2, huruf besar; Save hanya menerima OR, QS, SURPLUS |
| `COVERAGE` | teks | ya | | grid, chart | CSV 3 |
| `AS_AT` | date | ya | | grid, daftar | CSV 4 `dd/MM/yyyy` |
| `UW_YEAR` | teks | ya | | grid, daftar | CSV 5 |
| `CEDING_CODE` | teks | ya | | grid, daftar, chart | `CEDINGID` Master ID pertama (CSV 6 ditimpa) |
| `CEDING_NAME` | teks | ya | | grid, daftar, chart | `CEDING` Master ID pertama (CSV 7 ditimpa); kosong = Save menolak |
| `CURRENCY` | teks | ya | | grid | CSV 8 |
| `TO_USD` | angka | ya | | grid | `TREATYEXCHANGEYEARLY` tahun `TREATYYEAR` + mata uang: TOUSD; IDR = 1 / TOIDR baris USD, 20 desimal (CSV 9 ditimpa); tanpa kurs = kosong (Save menolak) |
| `NOR_BUILDINGS` | angka | ya | | grid | CSV 10; titik ribuan dibuang, koma = desimal |
| `BUILDINGS` | angka | ya | | grid | CSV 11 |
| `NOR_STOCKS` | angka | ya | | grid | CSV 12 |
| `STOCKS` | angka | ya | | grid | CSV 13 |
| `NOR_MACHINERY` | angka | ya | | grid | CSV 14 |
| `MACHINERY` | angka | ya | | grid | CSV 15 |
| `NOR_OTHER_CONTENTS` | angka | ya | | grid | CSV 16 |
| `OTHER_CONTENTS` | angka | ya | | grid | CSV 17 |
| `NOR_CONSEQUENTIAL_LOSS` | angka | ya | | grid | CSV 18 |
| `CONSEQUENTIAL_LOSS` | angka | ya | | grid | CSV 19 |
| `NOR_RESIDENTIAL` | angka | ya | | grid | CSV 20 |
| `RESIDENTIAL` | angka | ya | | grid | CSV 21 |
| `NOR_COMMERCIAL` | angka | ya | | grid | CSV 22 |
| `COMMERCIAL` | angka | ya | | grid | CSV 23 |
| `NOR_INDUSTRIAL` | angka | ya | | grid | CSV 24 |
| `INDUSTRIAL` | angka | ya | | grid | CSV 25 |
| `NOR_AGRICULTURE` | angka | ya | | grid | CSV 26 |
| `AGRICULTURE` | angka | ya | | grid | CSV 27 |
| `NOR_MISCELLANEOUS` | angka | ya | | grid | CSV 28 |
| `MISCELLANEOUS` | angka | ya | | grid | CSV 29 |
| `NOR_UTILITIES` | angka | ya | | grid | CSV 30 |
| `UTILITIES` | angka | ya | | grid | CSV 31 |
| `TOTAL_NO_OF_RISK` | angka | ya | | grid, total | CSV 32 apa adanya (tidak dihitung ulang) |
| `TOTAL_IN_AMOUNT` | angka | ya | | grid, total | CSV 33 apa adanya (tidak dihitung ulang) |
| `TOTAL_IN_AMOUNT_IN_USD` | angka | ya | | grid, total | `TO_USD` x `TOTAL_IN_AMOUNT` (CSV 34 ditimpa) |
| `RNM_SHARE` | angka | ya | | grid | `RNM_SHARE` Master ID pertama; beberapa Master ID = jumlahnya (CSV 35 ditimpa); kosong = Save menolak |
| `RNM_VALUE` | angka | ya | | grid, total | `TOTAL_IN_AMOUNT` x `RNM_SHARE` / 100, 4 desimal setengah-ke-atas (CSV 36 ditimpa) |
| `RNM_VALUE_IN_USD` | angka | ya | | grid, total, chart | `TO_USD` x `RNM_VALUE` share Master ID PERTAMA - beberapa Master ID: tidak dihitung ulang, seperti Pega (CSV 37 ditimpa) |
| `REMARK` | teks | ya | | grid | CSV 38 |
| `COMMENCEMENT` | date | ya | | — | **tidak diisi** (Pega tidak pernah mengisinya; keputusan work owner 04-10-2026). DEV: 13.440 baris lama terisi, dibiarkan |
| `TREATYYEAR` | teks | ya | | grid | `TREATYYEAR.TREATYYEAR` periode `STARTDATE` <= As At <= `ENDDATE`, baris pertama tanpa urutan (keputusan work owner: ikuti XML). DEV: data lama memakai tahun Master Treaty (978 dari 13.450 cocok dengan rumus ini) |
| `M_TREATY_ID` | teks | ya | | grid | `TREATYID` Master ID pertama; beberapa Master ID = digabung `;` (lebar 15 byte: gabungan yang lebih panjang ditolak Save) |

## Tabel dan view yang dibaca saja

| Objek | Kolom dibaca | Untuk |
| --- | --- | --- |
| `TREATYINDETAILJOINEDM` (view) | `ID`, `TREATYID`, `CEDINGID`, `CEDING`, `RNM_SHARE`, `TREATYYEAR`, `PROPORTIONTYPE`, `TREATYCONTRACTNAME`, `TREATYGROUP`, `SOB` | popup Master ID: `CEDING` memuat kata cari (huruf besar), (`Proportional` dan `PROPERTY`) atau `NonProportional`; ganda menurut TREATYID, CEDINGID, RNM_SHARE, TREATYYEAR, TREATYGROUP dibuang. DEV: 36.303 baris, `ID` unik |
| `ASSESSMENT_ZONE` | `ASSESSMENT_CODE`, `ASSESSMENT_NOTE` | Assessment Zone (DEV: 22 zona) |
| `TREATYYEAR` | `TREATYYEAR`, `STARTDATE`, `ENDDATE` (teks `yyyyMMdd`) | Treaty Year dari As At (DEV: 183 baris, periode per grup bertumpuk) |
| `TREATYEXCHANGEYEARLY` | `TREATYYEAR`, `CURRENCY`, `TOUSD`, `TOIDR` (teks) | To USD (DEV: 140 baris) |

## Sequence

- `SEQ_AGGREGATE` (migrasi 880): nomor `n` ID baru; mulai dari nomor AGG terbesar + 1, dihitung di basis data tempat
  migrasi berjalan. Generator Pega (`pzGenerateUniqueID`) tidak dibaca.
