---
status: selesai
---

# 13: Bentuk lama boleh masuk utuh, di satu tempat saja


> **SELESAI 19 September 2026, dikerjakan bersama tiket `14`.** Keduanya dua sisi satu seam: apa yang mendarat, dan apa yang tidak terurai darinya.
>
> **Satu kriteria penerimaan DIGANTI — bukan diberi catatan, tetapi ditulis ulang di daftarnya.**
>
> Bunyi lama: *"Muatan yang bukan JSON sah **ditolak**."* Itu **dicabut dan diganti tiga kriteria baru** di bawah. `CK_MIGRASI_PENDARATAN_1` tidak lagi menjaga `MUATAN IS JSON`.
>
> Dua sebab, keduanya terbaca dari sistem lama:
>
> | Sebab | Bukti |
> |---|---|
> | Bentuk muatan lama **tidak dibatasi** | **D43** — `adoptJSONObject` di tiga rule mengadopsi teks `HASIL1` tanpa memeriksa bentuknya, dan salah satunya menulis hasilnya ke halaman objek kerja. Bentuk yang tidak dibatasi tidak dapat dijaga dengan CHECK; ia hanya dapat **ditampung** |
> | Sistem lama **memang menghasilkan JSON rusak** | **D41** — muatan kasir dirakit dengan penyambungan teks; `Nett`, `Deductible`, `KaliDeduct` dikirim tanpa kutip. Nilai kosong menghasilkan `"Nett":,` — JSON tidak sah |
>
> Gerbang `IS JSON` akan menolak **persis muatan yang paling perlu tercatat**, dan menolaknya bertabrakan dengan **AK-4**: tidak ada baris yang dibuang karena "cuma sedikit". Muatan yang ditolak di pintu tidak punya rumah.
>
> **Yang menggantikannya**: muatan mendarat utuh apa pun bentuknya, dan hasil penguraiannya **dicatat sebagai fakta** — `BENTUK_TERURAI` biner, `SEBAB_GAGAL_URAI` terisi bila dan hanya bila gagal. `CK_MIGRASI_PENDARATAN_2` menolak dua keadaan yang mustahil: gagal tanpa sebab, dan berhasil tapi bersebab.
>
> Itu `SPEC-MODEL-DATA.md` §21.8 sebagai tabel.

*Asal: `T-21` di `_migration-docs/claim-non-prop/TICKETS.md`. Lapisan data saja — tidak ada Golang, React, endpoint, layar, service, repository, atau ORM.*

**What to build:** Satu-satunya tempat JSON hidup di skema ini.

`MIGRASI_PENDARATAN` beserta `PK_MIGRASI_PENDARATAN`, `CK_MIGRASI_PENDARATAN_1`, `CK_MIGRASI_PENDARATAN_2`, `IX_MIGRASI_PENDARATAN_1` — `ddl-usulan/21_MIGRASI_PENDARATAN.sql`.

**Blocked by:**

- ~~`01`~~ *(selesai)* — Skema `KLAIMNP` dan akun aplikasi berdiri, dan tidak ada akun lain yang dapat menulis


**Dasar:** DECIDED(ADR-0028). EVIDENCED: `BLUEPRINT.md` §13.2, §13.5 — nilai uang lama tersimpan sebagai **teks di dalam JSON**, dan `CLAIMXOL` mengeluarkannya sebagai `varchar2`.

- [x] **Muatan mendarat utuh, apa pun bentuknya** — tidak ada muatan yang ditolak di pintu.
- [x] **Hasil penguraian tercatat sebagai fakta**, bukan dipakai sebagai syarat masuk — `BENTUK_TERURAI`, dan `SEBAB_GAGAL_URAI` terisi bila dan hanya bila gagal (`CK_MIGRASI_PENDARATAN_2`).
- [x] Muatan yang gagal terurai **tetap tersimpan utuh** dan dapat ditunjuk dari `MIGRASI_NILAI_DITOLAK`.
- [x] Tidak ada kolom JSON di tabel kanonik mana pun — diperiksa dengan sapuan atas seluruh skema.

**Ketidakpastian:** Tidak ada REQ yang menahan. REQ-009 memverifikasi — ia memperbanyak daftar medan yang dikenali, dan tabel ini menerima apa pun yang datang entah medannya dikenali atau tidak.
