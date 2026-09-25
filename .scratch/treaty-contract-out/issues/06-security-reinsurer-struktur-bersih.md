# 06: Security reinsurer — struktur bersih

**Status:** ready-for-agent

**Blocked by:** 05 (security menggantung pada reinsurer)

## Hasil & nilai pengguna

Sebagai **admin master treaty**, saya ingin mencatat **security** di bawah seorang reinsurer beserta
porsinya, supaya eksposur berjenjang terlihat; dan sebagai **organisasi**, saya ingin baris security
punya **identitas sendiri**, supaya mengubah nama security tidak memutus rujukannya.
*(User story 13–14 di spec)*

⚠️ Inilah tiket dengan pembersihan struktur paling dalam di konteks ini.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/models` | Entitas security reinsurer |
| `internal/repository` | Tulis/baca dengan **kolom bernama** dan **PK surrogate** |
| `internal/services` | Aturan per baris |
| `internal/handlers` | Endpoint security |
| `frontend/` | Grid security di bawah baris reinsurer |

## Rule Pega sumber

| Rule | Class / Nama / Tipe | Path | Peran |
| --- | --- | --- | --- |
| `SaveSecurityReinsurer_Act` | `@BASECLASS` / `SAVESECURITYREINSURER_ACT` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/SaveSecurityReinsurer_Act.xml` | orkestrator simpan |
| `ShowEditSecurityReinsurer` | `@BASECLASS` / `SHOWEDITSECURITYREINSURER` / `RULE-OBJ-ACTIVITY` | `Treaty Contract Out/Activity/ShowEditSecurityReinsurer.xml` | muat untuk diubah |
| `InsertToMTreatySecurity` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `ASM!INSERTTOMTREATYSECURITY` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/InsertToMTreatySecurity.xml` | ⚠️ INSERT **posisional** |
| `UpdateMTreatySecurity` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/UpdateMTreatySecurity.xml` | ⚠️ kunci `trim()` |
| `DeleteSecurityReinsurer` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `RULE-CONNECT-SQL` | `Treaty Contract Out/RDBList/DeleteSecurityReinsurer.xml` | ⚠️ kunci `trim()` |
| `SelectSecurityReinsurer` | `ASM-FW-GISFW-INT-MTREATYSECURITY` / `RULE-OBJ-REPORT-DEFINITION` | `Treaty Contract Out/ReportDefinition/SelectSecurityReinsurer.xml` | daftar |

`[terverifikasi]` Keadaan existing — **satu-satunya master di modul ini yang ditulis SQL mentah**:

```
insert into mtreatysecurity values ({…}, '', '', {…}, {…PCT_SHARE}, '', {…REAS_SECURITY})
update mtreatysecurity set THN_TREATY=…, PCT_SHARE=…, REAS_SECURITY=…
 where REAS_ID = … and trim(REAS_SECURITY) = trim(…)
delete from mtreatysecurity where REAS_ID = … and trim(REAS_SECURITY) = trim(…)
```

⚠️ INSERT **tanpa daftar kolom** — tujuh nilai berposisi, **tiga di antaranya string kosong**.
⚠️ Ketiganya **tanpa `COMMIT`**.
⚠️ Kunci pembaruan dan penghapusan memakai **`trim(REAS_SECURITY)`** — nama dipakai sebagai bagian
kunci, dan `trim()` mengakui datanya bertabur spasi.

`[data DBA]` DDL existing: `THN_TREATY VARCHAR2(4) DEFAULT '1' NOT NULL`, `TOP_ID VARCHAR2(9)`,
`TP_TREATY CHAR(2)`, `REAS_ID CHAR(7) NOT NULL`, `PCT_SHARE VARCHAR2(99)`, `USER_ID CHAR(99)`,
`REAS_SECURITY CHAR(10) NOT NULL` — **tanpa primary key**. INSERT posisional mengosongkan
`TOP_ID`, `TP_TREATY`, `USER_ID`.

⚠️ **Penyimpangan sadar 5 — struktur dibersihkan.** `[keputusan work owner]` Tabel anak
`TREATYREINSURER` yang wajar: **PK surrogate**, seluruh kolom **bernama dan diisi eksplisit**,
`PCT_SHARE` **desimal**, `REAS_SECURITY` **atribut biasa** — bukan bagian kunci.

## ADR terkait

**ADR-0003** (persen non-float), **ADR-0007** (jejak audit), **ADR-0015** (kegagalan eksplisit).

## Acceptance criteria

- [ ] Security dicatat **di bawah seorang reinsurer** beserta porsinya. *(AC 17 spec; User story 13)*
- [ ] ⚠️ Baris security punya **primary key surrogate** sendiri; `REAS_SECURITY` adalah **atribut
      biasa**, **bukan** bagian kunci. Test yang menemukan kunci berbasis nama security **gagal**.
      *(AC 18 spec; User story 14; penyimpangan sadar 5)*
- [ ] ⚠️ Seluruh kolom baris security **diisi eksplisit dan bernama**; **tidak ada** penulisan
      berposisi. Test yang menemukan penulisan tanpa daftar kolom **gagal**. *(AC 19 spec;
      penyimpangan sadar 5)*
- [ ] ⚠️ Pencocokan baris security **tidak** memakai `trim()` atas nilai kunci. *(AC 20 spec;
      penyimpangan sadar 5)*
- [ ] ⚠️ `PCT_SHARE` bertipe **desimal**, bukan teks. *(AC 51 spec; **ADR-0003**; penyimpangan
      sadar 6)*
- [ ] Mengubah **nama security** **tidak** memutus rujukan barisnya ke reinsurer. *(User story 14)*
- [ ] Baris security dapat **ditambah dan dihapus** tanpa menyentuh baris reinsurer induknya —
      kecuali ketika penyimpanan dilakukan sebagai satu kesatuan (tiket 09).
- [ ] Menghapus seorang reinsurer **menghapus juga** baris security di bawahnya. *(lihat tiket 10
      untuk kaskade dari kontrak)*
- [ ] Penyimpanan yang gagal menghasilkan kegagalan **terang-terangan**. *(AC 38, 39 spec;
      **ADR-0015**)*
- [ ] Setiap penyimpanan mencatat **jejak audit**. *(AC 41 spec; **ADR-0007**)*

## Blocker

**Tidak ada.**

⚠️ `[terbuka]` **Tiga kolom yang dikosongkan INSERT posisional** — `TOP_ID`, `TP_TREATY`,
`USER_ID`. `[data DBA]` DDL-nya sudah diketahui, tetapi **apakah ketiganya masih punya arti** tidak
terbaca dari korpus. Tiket ini **membawanya sebagai kolom bernama**; bila migrasi (tiket 01)
membuktikan seluruh barisnya kosong, pembuangannya adalah keputusan terpisah — **jangan dibuang di
sini atas inisiatif sendiri**.

## Catatan

⚠️ **Mengapa `trim()` tidak dibawa.** Pemakaian `trim()` di kunci adalah **pengakuan bahwa datanya
kotor**, bukan aturan bisnis. Membawanya ke sistem baru berarti membawa kekotorannya, dan menutup
kemungkinan dua security yang namanya hanya berbeda spasi tetap dianggap berbeda. Migrasi (tiket 01)
memberi setiap baris identitas sendiri; sesudah itu nama tidak perlu lagi memikul beban kunci.

⚠️ `RDBList/DeleteFromTreatyReinsurer_Act.xml` menghapus `MTREATYSECURITY` lalu `TREATYREINSURER`
**tanpa `COMMIT`** — konsisten dengan temuan bahwa modul ini menyerahkan batas transaksi kepada
pemanggil (tiket 09).

## Seam & perintah verifikasi

**Seam: API HTTP** terhadap **skema uji Oracle nyata** — keutuhan rujukan setelah pemberian PK dan
hilangnya kunci berbasis nama **hanya berperilaku benar pada basis data sungguhan**.

```
go test ./internal/...
cd frontend && npm test
make check
```
