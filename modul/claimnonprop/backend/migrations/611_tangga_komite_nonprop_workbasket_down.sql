-- Mundur 611 - baris NONPROP ber-OPERATOR_ID workbasket dikosongkan. ⚠️ Akun orang roster lama TIDAK dipulihkan
-- (nama orang dilarang di kode): OPERATOR_ID, NAME kosong - diisi ulang admin. Workbasket milik claimprop 537 tidak
-- disentuh. Untuk skema uji.
UPDATE {skema}.EMAILKOMITE SET OPERATOR_ID = NULL, NAME = NULL
 WHERE STS_KLAIM = 'NONPROP'
   AND OPERATOR_ID IN ('ReasClaimDeptHead', 'ReasClaimTechDivHead', 'ReasClaimOpsDir', 'ReasClaimTechDir')
/
