-- Mundur 640 - baris FACIN ber-OPERATOR_ID workbasket dikosongkan. ⚠️ Akun orang roster lama TIDAK dipulihkan (nama
-- orang dilarang di kode): OPERATOR_ID, NAME kosong - diisi ulang admin. ReasClaimSPVA / ReasClaimSPVB TIDAK dibuang:
-- keduanya sudah ada di DEV sebelum 640 (sisipan 640 hanya untuk skema baru); empat workbasket lain milik claimprop 537.
-- Untuk skema uji.
UPDATE {skema}.EMAILKOMITE SET OPERATOR_ID = NULL, NAME = NULL
 WHERE STS_KLAIM = 'FACIN'
   AND OPERATOR_ID IN ('ReasClaimSPVA', 'ReasClaimDeptHead', 'ReasClaimTechDivHead', 'ReasClaimOpsDir', 'ReasClaimTechDir')
/
