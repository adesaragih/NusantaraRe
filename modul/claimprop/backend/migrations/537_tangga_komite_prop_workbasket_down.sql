-- Mundur 537 - empat workbasket dibuang beserta pemegangnya; baris PROP ber-OPERATOR_ID workbasket dikosongkan.
-- ⚠️ Akun orang roster lama TIDAK dipulihkan (nama orang dilarang di kode): OPERATOR_ID, NAME kosong - diisi ulang
-- admin. Untuk skema uji.
UPDATE {skema}.EMAILKOMITE SET OPERATOR_ID = NULL, NAME = NULL
 WHERE STS_KLAIM = 'PROP'
   AND OPERATOR_ID IN ('ReasClaimDeptHead', 'ReasClaimTechDivHead', 'ReasClaimOpsDir', 'ReasClaimTechDir')
/
DELETE FROM {skema}.M_LOGIN_GO_WORKBASKET
 WHERE WORKBASKET_ID IN ('ReasClaimDeptHead', 'ReasClaimTechDivHead', 'ReasClaimOpsDir', 'ReasClaimTechDir')
/
DELETE FROM {skema}.M_WORKBASKET
 WHERE WORKBASKET_ID IN ('ReasClaimDeptHead', 'ReasClaimTechDivHead', 'ReasClaimOpsDir', 'ReasClaimTechDir')
/
