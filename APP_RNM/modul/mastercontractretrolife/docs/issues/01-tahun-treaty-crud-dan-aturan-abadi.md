# 01: Tahun treaty — CRUD dan aturan abadi

**Status:** ready-for-agent

**Blocked by:** **CL-01** (kerangka aplikasi + seam API — scaffolding lintas konteks, tidak dibuat
di sini)

## Hasil & nilai pengguna

Sebagai **admin master retro life**, saya dapat membuat dan mengubah **tahun treaty** — akar seluruh
hierarki retrosesi life — dan saya tidak dapat menghapusnya, sehingga riwayat kontrak tahun-tahun
lampau tidak pernah hilang. *(User story 1–4 di spec)*

⚠️ Tiket ini juga **menetapkan pola jalur simpan** yang diwarisi keempat entitas berikutnya: upsert
dikunci identitas, identitas dibuat basis data, dan hasil procedure diperiksa.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas tahun treaty |
| `internal/repository` | Pemanggilan procedure penulis; pemeriksaan `o_message` |
| `internal/services` | Wajib-isi; larangan hapus |
| `internal/handlers` | Endpoint buat/ubah/daftar tahun treaty |
| `frontend/` | Grid tahun treaty |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveMasterTreatyYear_Life_SQL` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `ASM!SAVEMASTERTREATYYEAR_LIFE_SQL` / `RULE-CONNECT-SQL` | `Master Contract Retro Life/RDBList/SaveMasterTreatyYear_Life_SQL.xml` | `POOLDATA.INSERTTREATYYEAR_LIFE` |
| `SaveTreatyYearLife_Act` | `ASM-FW-GISFW-…` / `SAVETREATYYEARLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SaveTreatyYearLife_Act.xml` | orkestrator simpan |
| `NewInputTreatyYear_Life_Act` | `ASM-FW-GISFW-…` / `NEWINPUTTREATYYEAR_LIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/NewInputTreatyYear_Life_Act.xml` | baris baru |
| `SetTreatyYearLife_Act` | `ASM-FW-GISFW-…` / `SETTREATYYEARLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `Master Contract Retro Life/Activity/SetTreatyYearLife_Act.xml` | isi form |
| `BrowseTreatyYear_Life_RD` | `ASM-FW-GISFW-INT-TREATYYEAR_LIFE` / `BROWSETREATYYEAR_LIFE_RD` / `RULE-OBJ-REPORT-DEFINITION` | `Master Contract Retro Life/ReportDefinition/BrowseTreatyYear_Life_RD.xml` | daftar |

`[terverifikasi]` Wajib-isi saat simpan (`SaveTreatyYearLife_Act`): `UNDERWRITINGYEAR`,
`TREATYYEAR`, `STARTDATE`, `ENDDATE`.

`[terverifikasi]` **Tidak ada penghapus** untuk `TREATYYEAR_LIFE` — lima penulis, hanya **empat**
penghapus di seluruh modul.

`[data DBA]` Perilaku `INSERTTREATYYEAR_LIFE`: **upsert dikunci `ID`**
(`SELECT COUNT(1) … WHERE ID = p_ID` → ada `UPDATE`, tidak ada `INSERT`); identitas baris baru
`'1' || lpad(TREATYYEAR_LIFE_SEQ.nextval, 6, '0')`; `p_ID` **diabaikan** saat INSERT;
`TGLUPDATE` = **`SYSDATE`** (parameter `p_TGLUPDATE` diabaikan); `USERID` dari aplikasi;
**`COMMIT` di dalam procedure**; **nol validasi bisnis**.

`[data DBA]` Kolom (`ddl-tables-from-dba.md`): `ID VARCHAR2(100)` (**PK**), `TREATYYEAR VARCHAR2(100)`,
`UNDERWRITINGYEAR VARCHAR2(100)`, `USERID VARCHAR2(100)`, `TGLUPDATE DATE`, `STARTDATE DATE`,
`ENDDATE DATE`. **Semua nullable** — wajib-isi ditegakkan **di Go**.
Sequence `TREATYYEAR_LIFE_SEQ` `START WITH 44` → identitas berikutnya `1000044`.

## ADR terkait

**ADR-0006** (identitas dibuat lewat procedure/basis data — aplikasi tidak menyusunnya),
**ADR-0007** (jejak audit), **ADR-0009** (migrasi penuh).

## Acceptance criteria

- [ ] Tahun treaty dapat dibuat dengan tahun underwriting, tahun treaty, tanggal mulai, dan tanggal
      akhir; keempatnya **wajib** — ditegakkan **di Go**, karena basis data **nol `NOT NULL`**.
      *(AC 1 spec)*
- [ ] ⚠️ **Tahun treaty TIDAK DAPAT DIHAPUS** lewat jalur mana pun. Test yang menemukan endpoint
      hapus tahun **gagal**. *(AC 2 spec; `[fakta bisnis — work owner]` — tahun treaty **abadi**,
      disengaja)*
- [ ] Tahun treaty dapat diubah, dan perubahannya mencatat **siapa** pelakunya. *(AC 3 spec)*
- [ ] ⚠️ **Aplikasi tidak menetapkan identitas baris baru** — ia mengirim identitas kosong dan basis
      data yang membuatnya lewat sequence. Test yang menemukan pembentukan identitas di sisi aplikasi
      **gagal**. *(AC 4 spec; **ADR-0006**)*
- [ ] Menyimpan tahun yang **sudah ada** memperbarui baris itu — **upsert dikunci identitas**, bukan
      baris kedua. Dibuktikan dengan menyimpan dua kali lalu **menghitung baris**. *(AC 45 spec)*
- [ ] Setiap penyimpanan mengirim **identitas pengguna**; **cap waktu tidak dikirim aplikasi** —
      basis data yang menetapkannya. *(AC 48 spec; `[data DBA]`)*
- [ ] `o_message` **diperiksa** setelah pemanggilan procedure: **kosong/NULL = sukses**, berisi teks
      = gagal. Kegagalan **ditampilkan**, tidak ditelan. *(AC 46 spec — pembersihan HTML dan
      konformansi lintas jalur ada di **tiket 10**)*
- [ ] Daftar tahun treaty dapat dibaca lewat API dan tampil di layar.

## Blocker

**Tidak ada.** Menunggu scaffolding **CL-01** yang sudah `ready-for-agent` di konteks Claim Life.

## Catatan

⚠️ **Pola yang ditetapkan di sini diwarisi tiket 02, 05, 06, dan 07.** Keempat jalur simpan lain
berperilaku identik: upsert dikunci `ID`, identitas dari sequence, `TGLUPDATE` = `SYSDATE`, `COMMIT`
internal, nol validasi di basis data.

⚠️ **Setiap simpan adalah transaksi mandiri** — `COMMIT` berada di dalam procedure. Tidak ada
transaksi lintas-baris di sisi basis data, dan itu **tidak boleh disembunyikan** dari pengguna.
*(AC 50 spec)*

## Seam & perintah verifikasi

**Seam: API HTTP** (dipakai ulang dari CL-01) terhadap **skema uji Oracle nyata** — procedure
**tidak di-mock**: perilaku upsert dan pembuatan identitas lewat sequence tidak dapat difake dengan
jujur.

```
go test ./internal/...
cd frontend && npm test
make check
```

---

## Ralat bertanggal 30-09-2026 — sesi implementasi (paket 0)

> Sumber: `RALAT-DEV-30-09-2026.md` (K1–K8 katalog DEV, R1–R12 pembacaan ulang XML) dan `PARITAS-LAYAR-DAN-AKSI.md`. Kalimat di atas **tidak dihapus**; yang berlaku adalah ralat ini.

| Kalimat lama | Ralat |
| --- | --- |
| *"Pemanggilan procedure penulis; pemeriksaan `o_message`"* | procedure **tidak dipanggil** (keputusan o, R4): upsert, `ID` `'1' ‖ LPAD(TREATYYEAR_LIFE_SEQ, 6, '0')`, `TGLUPDATE = SYSDATE`, `USERID` pelaku ditiru di Go; tidak ada `o_message` — galat Go sampai ke layar (K8) |
| *"`ID VARCHAR2(100)` (**PK**)"* | nol PK di DEV (K1); keunikan dijamin sequence (K3) |
| — | layar: grid `BrowseTreatyYear_Life_RD` urut `ID ASC`; tombol baris **`Edit`**, **`ReinsType`**; baris baru lewat tombol berlabel VERBATIM **`End Period`** (`InputRetrocessionLife.xml` b8888); form `Input New Data` dengan `Save`/`Cancel`; wajib-isi berpesan VERBATIM `"Value cannot be empty."` (`SaveTreatyYearLife_Act` b313). Label kolom `TREATYYEAR` = **TRANSACTION YEAR** |
