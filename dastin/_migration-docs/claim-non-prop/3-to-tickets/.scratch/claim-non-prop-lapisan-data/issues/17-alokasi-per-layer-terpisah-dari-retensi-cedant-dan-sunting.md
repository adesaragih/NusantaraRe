---
status: selesai
---

# 17: Alokasi per Layer terpisah dari Retensi Cedant, dan suntingan tidak dapat tertimpa hitung ulang

> **DILEPAS DARI KETERGANTUNGAN KOMITE, 19 September 2026.** Folder `Komite Claim Non Prop` dinyatakan **tertutup untuk batch ini** — bukan ditunda. F1/F2 karena itu tinggal permanen sebagai hipotesis sejajar, dan tiket ini **tidak boleh mati menunggunya**.
>
> **Dikerjakan di atas aturan yang sama-sama benar di kedua cabang:**
>
> | Cabang | Tuntutannya | Di tiket ini |
> |---|---|---|
> | F1 — pembacanya di Komite | kopling lintas modul masuk Boundary Contract | entri ditulis **bertanda belum terselesaikan**, arah dan pemilik **dikosongkan**, bukan ditebak |
> | F2 — tidak ada pembaca | sistem baru punya kunci hasil suntingan yang sebenarnya | **dibangun tanpa syarat** — ADR-0008 menuntutnya terlepas dari apakah sistem lama pernah punya jalurnya |
>
> Pemisahan hitung/sunting karena itu **tetap dipasang**, dan daftar kandidatnya sudah berbukti: sepuluh kolom yang ditulis mesin alokasi **dan** salah satu dari `AdjClaimCNP_Act`, `AdjClaimAmount_Act`, `SetActualPremium_ACT` (E15, S15).
>
> **26 penanda yatim golongan 3b** mengikuti kebijakan E17: **dimigrasi apa adanya dan ditandai tak berpemilik** — tidak dibuang, tidak diberi perilaku.


> **SELESAI 19 September 2026.**
>
> **Kriteria kelima tidak pernah punya constraint — sampai hari ini.** *"`DISUNTING_OLEH` kosong sementara `_SUNTING` terisi **ditolak** — suntingan selalu punya pelaku"* tertulis sejak tiket ini dibuat, dan **tidak ada satu pun `CHECK`** yang menegakkannya di `ALOKASI_LAYER` maupun `RETENSI_CEDANT`.
>
> Sebabnya ADR-0008: suntingan manual bertahan **dan terlihat**. Nilai yang menimpa hasil hitung tanpa ada yang bertanggung jawab atasnya bertahan **tanpa terlihat** — dan itu setengah dari yang dijanjikan, yaitu bagian yang justru lebih mudah disalahgunakan.
>
> Dipasang: `CK_ALOKASI_LAYER_4` (lima besaran, termasuk `ALASAN_SUNTING`) dan `CK_RETENSI_CEDANT_4`. `DISUNTING_ATAS_NAMA` tetap boleh kosong — kosong berarti tidak ada perwakilan, bukan tidak diketahui (ADR-0019).
>
> **Kedua tabel juga terkena dua cacat pasangan IDR–kurs** yang ditemukan di tiket `16`; keduanya sudah diperbaiki, dan asal-usul kurs kini ditegakkan di sini juga.
>
> **Empat kriteria lainnya sudah terpenuhi oleh bentuk yang ada**: `_DIPAKAI` adalah kolom `GENERATED ALWAYS AS (COALESCE(_SUNTING, _HITUNG)) VIRTUAL` — hitung ulang menimpa `_HITUNG` tanpa menyentuh `_SUNTING`, dan menulis langsung ke `_DIPAKAI` ditolak Oracle sendiri, bukan oleh aturan yang harus diingat orang.

*Asal: `T-05` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Hitung ulang menimpa `_HITUNG` dan tidak menyentuh `_SUNTING`; `_DIPAKAI` mengikuti tanpa ditulis siapa pun.

`ALOKASI_LAYER` dan `RETENSI_CEDANT`; kolom berpasangan `_HITUNG`/`_SUNTING`/`_DIPAKAI`; parameter yang di-snapshot; `UQ_ALOKASI_LAYER_1`, `UQ_RETENSI_CEDANT_1`.

**Tidak termasuk:** Kolom jenis reasuransi pada `RETENSI_CEDANT` — **sengaja tidak ada**. EVIDENCED: baris Retensi Cedant menulis `TreatyName = "UR"` secara literal melewati lookup (FINDING-003 bagian 2); `"UR"` artefak penyatuan yang dibatalkan ADR-0010, bukan nilai domain.

**Blocked by:**

- ~~`12`~~ *(selesai)* — Klaim tersimpan sebagai aggregate root, dengan satu Tanggal Kejadian yang tetap
- ~~`08`~~ *(selesai)* — Sapuan S4: adakah padanan IDR untuk biaya penilaian, salvage, dan biaya lain
- ~~`09`~~ *(selesai)* — Sapuan S9: apa persisnya yang ditulis `EditXOLAlokasi`


**Dasar:** DECIDED(ADR-0010, ADR-0008, ADR-0011). EVIDENCED: `BLUEPRINT.md` §2.5 — `.IsEditClaim` **ditulis lima kali, tidak pernah dibaca**; FINDING-004.

- [x] Menjalankan ulang pengisian `_HITUNG` atas baris yang `_SUNTING`-nya terisi **tidak mengubah** nilai yang dibaca lewat `_DIPAKAI`.
- [x] Menulis langsung ke `_DIPAKAI` **ditolak** — ia kolom turunan.
- [x] Baris kedua dengan (klaim, keempat field layer, mata uang) sama **ditolak** — `UQ_ALOKASI_LAYER_1`.
- [x] Baris Retensi Cedant kedua untuk (klaim, mata uang) sama **ditolak** — `UQ_RETENSI_CEDANT_1`.
- [x] `DISUNTING_OLEH` kosong sementara `_SUNTING` terisi **ditolak** — suntingan selalu punya pelaku. **Baru ditegakkan 19 Sep 2026**: `CK_ALOKASI_LAYER_4`, `CK_RETENSI_CEDANT_4`.

**Ketidakpastian:** **T-36** menentukan besaran mana yang benar-benar berpasangan `_HITUNG`/`_SUNTING`; memasangnya di kolom yang tidak pernah disunting adalah beban mati. **T-35** menentukan apakah Biaya Penilaian, Salvage, dan Biaya Lain perlu pasangan IDR di tabel ini. **REQ-033** (sebaran `LayerPart`) menentukan apakah kunci halus pernah melahirkan dua baris di tempat sistem lama melihat satu — itu mengubah **view** di T-16, bukan tabel ini.
