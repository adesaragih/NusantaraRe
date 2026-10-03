# Definisi view produk Life di DEV — dibaca asisten 30-09-2026 dari `ALL_VIEWS` (baca-saja)

Tiga view Oracle membaca `JSONDATA` tabel warisan `M_PRODUCT_LIFE` dan `M_PRODUCTINWARD_LIFE`. Kode Claim Life
membaca view `PRODUCTINWARD_LIFE` (ambang produk, kategori dokumen, rate, validasi tanggal). Dua prosedur Pega,
`PEGA_M_PRODUCT_LIFE` dan `PEGA_M_PRODUCT_INWARD_LIFE`, juga bergantung pada kedua tabel (`ALL_DEPENDENCIES`).

Akibatnya bagi modul Master Product Name Life:

- Produk yang disimpan aplikasi baru **wajib** tetap terbaca lewat ketiga view ini.
- Nama kunci JSON **peka huruf besar-kecil** dan harus persis: mis. `DocumentClaim[*].Document`,
  `UnderwritingLimitList`, `OutwardList[0].OVR_COMM`, `POLICYHODER` (ejaan Pega, bukan `POLICYHOLDER`).
- `OutwardList` tetap dibaca view `PRODUCT_LIFE` walau jalur pengisinya mati di Pega; kuncinya ditulis seperti Pega menulisnya.

Cacah baris DEV saat dibaca: `M_PRODUCT_LIFE` 196, `M_PRODUCTINWARD_LIFE` 197. Masing-masing punya satu constraint cek (`IS JSON`), nol PK.

```sql
=== VIEW DOCUMENTCLAIM_LIFE
select
       b.ID,
       jt.DOCUMENT
from
     pooldata.M_PRODUCT_LIFE b ,
     json_table(b.jsondata, '$'
    COLUMNS(
    NESTED PATH'$.DocumentClaim[*]'
        COLUMNS(DOCUMENT VARCHAR2 PATH '$.Document'))) JT

=== VIEW PRODUCTINWARD_LIFE
SELECT a.ID,
      a.JSONDATA.PRODUCTID,
      a.JSONDATA.INSURED,
      a.JSONDATA.CEDING,
      a.JSONDATA.TREATYNUMBER,
      a.JSONDATA.ADDENDUMWORD,
      a.JSONDATA.AMANDEMENTSCHD,
      a.JSONDATA.INWARDTREATYNM,
      a.JSONDATA.BEGIN,
      a.JSONDATA.MATURE,
      a.JSONDATA.CEDINGRETENTIONNUM,
      a.JSONDATA.CEDINGRETENTIONPCT,
      a.JSONDATA.CEDINGLIMIT,
      a.JSONDATA.CEDINGLIMITXPN,
      a.JSONDATA.MINAGE,
      a.JSONDATA.MAXAGE,
      a.JSONDATA.BIRTHDAY,
      a.JSONDATA.EXTRAPREMI,
      a.JSONDATA.CURRENCY,
      a.JSONDATA.RNMSHARE,
      a.JSONDATA.EXTRAMORTALITY,
      a.JSONDATA.RNMLIMITNUM,
      a.JSONDATA.RNMLIMITPCT,
      a.JSONDATA.LIENCLAUSE,
      a.JSONDATA.MONTHS,
      a.JSONDATA.MINSUMINSURED,
      a.JSONDATA.MAXSUMINSURED,
      a.JSONDATA.MAXSUMREASURED,
      a.JSONDATA.MAXCONTRACT,
      a.JSONDATA.PAYMENT,
      a.JSONDATA.PROPORTIONALTABLE,
      a.JSONDATA.SUBJECTTO,
      a.JSONDATA.POLICYHODER,
      a.JSONDATA.POLICYHODERNAME,
      a.JSONDATA.BROKERAGE,
      a.JSONDATA.ADDENDUMNO,
      a.JSONDATA.AMANDEMENTNO,
      a.JSONDATA.MAXDATARECEIVE,
      a.JSONDATA.MAXEXPIREDCLAIM,
      a.JSONDATA.STNC
     FROM m_productinward_life a

=== VIEW PRODUCT_LIFE
SELECT a.ID,
      a.JSONDATA.TYPE,
      a.JSONDATA.TYPE_CEDING,
      a.JSONDATA.CEDING,
      a.JSONDATA.CEDINGID,
      a.JSONDATA.SOBNAME,
      a.JSONDATA.SOBID,
      a.JSONDATA.CAUSEID,
      a.JSONDATA.GRUP,
      a.JSONDATA.PRODUCTNAME,
      a.JSONDATA.PRODUCTCODE,
      a.JSONDATA.PRODUCTTYPEID,
      a.JSONDATA.PRODUCTTYPE,
      a.JSONDATA.RIRISKID,
      a.JSONDATA.RIRISK,
      a.JSONDATA.RIRATEID,
      a.JSONDATA.RIRATE,
      a.JSONDATA.RICOMMID,
      a.JSONDATA.RICOMM,
      a.JSONDATA.INWARDNAME,
      a.JSONDATA.UnderwritingLimitList,
      a.JSONDATA.OUTWARDNAMEID,
      a.JSONDATA.OUTWARDNAME,
      a.JSONDATA.OUTWARDRATEID,
      a.JSONDATA.OUTWARDRATE,
      a.JSONDATA.OUTWARDCOMMID,
      a.JSONDATA.OUTWARDCOMM,
      a.JSONDATA.BENEFITID,
      a.JSONDATA.BENEFIT,
      a.JSONDATA.CAUSE,
      a.JSONDATA.OutwardList[0] .OVR_COMM,
      a.JSONDATA.POLICYHODER,
      a.JSONDATA.POLICYHODERNAME,
      a.JSONDATA.TREATYNUMBER,
      a.JSONDATA.CREATEOP,
      a.JSONDATA.UPDATEOP
      FROM M_PRODUCT_LIFE a
```

## Catatan 02-10-2026 — view TIDAK dibangun ulang di atas tabel flat `[keputusan work owner 02-10-2026]`

Rancangan brief `PROMPT-PINDAH-FLAT-MASTER-PRODUCT-NAME-LIFE.md` §2.3/T5 (ketiga view di atas tabel flat, berkas migrasi) ditolak penjaga
inti `TestSeluruhCreateDapatDibacaNamanya`; jawaban work owner: *"tidak ada table view yang dipake, semua simpan dan baca dari table
flat"*. Definisi di atas karena itu **tetap** definisi yang berlaku di DEV: ketiga view terus membaca `JSONDATA` kedua tabel lama, yang
berhenti diperbarui sesudah peralihan. Satu-satunya pembaca kode di luar modul ini: Claim Life `repository/ambangproduk.go`
(`PRODUCTINWARD_LIFE`) — OQ-FLAT-04. `DOCUMENTCLAIM_LIFE` dan `PRODUCT_LIFE` tidak dibaca kode mana pun (hanya disebut di komentar).

