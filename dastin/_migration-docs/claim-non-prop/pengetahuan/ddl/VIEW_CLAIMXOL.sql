-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE FORCE EDITIONABLE VIEW "POOLDATA"."CLAIMXOL" ("CASEID", "GrossAdjustment", "CNPReinstatement", "Currency", "KursIDR", "XOL", "TANGGAL") AS 
  SELECT caseid,
          TotalXOLGross AS "GrossAdjustment",
          CNPReinstatement AS "CNPReinstatement",
          Currency AS "Currency",
          KursIDR AS "KursIDR",
          XOL,
          TANGGAL
     FROM pooldata.os_akseptasi_klaim OS, json_table (
                                    data_json, '$.CNPLayerList[*]'
                                    columns(
                                    XOL varchar2 path '$.XOL',
                                    nested path '$.CNPCurrencyList[*]' columns(TotalXOLGross varchar2 path '$.TotalXOLGross',CNPReinstatement varchar2 path '$.CNPReinstatement',Currency varchar2 path '$.Currency',KursIDR varchar2 path '$.KursIDR')
                                    ))where SUBSTR(CASEID,1,25)='ASM-FW-GCNMFW-WORK CLMNP-'

UNION ALL

SELECT caseid,
          TotalXOLGross AS "GrossAdjustment",
          CNPReinstatement AS "CNPReinstatement",
          Currency AS "Currency",
          KursIDR AS "KursIDR",
          XOL,
          TANGGAL
     FROM pooldata.OS_AKSEPTASI_SUBJECTIVITY OS, json_table (
                                    data_json, '$.CNPLayerList[*]'
                                    columns(
                                    XOL varchar2 path '$.XOL',
                                    nested path '$.CNPCurrencyList[*]' columns(TotalXOLGross varchar2 path '$.TotalXOLGross',CNPReinstatement varchar2 path '$.CNPReinstatement',Currency varchar2 path '$.Currency',KursIDR varchar2 path '$.KursIDR')
                                    ))where SUBSTR(CASEID,1,25)='ASM-FW-GCNMFW-WORK CLMNP-' and STS_SUBJECTIVITY = '1'