-- 192 - kolom alamat risiko dilebarkan (bug DEV 03-10-2026, diteruskan sesi `nusantarare-0f`).
--
-- Save tab Object gagal "baris[0].riskLocation paling banyak 50 byte": Risk Location = rangkaian Title + Address +
-- Territory + District + City + Province + Nation (`SetRiskIdDT_FacIn`), pasti > 50. Perintah work owner butir 87
-- (*"kolomnya buatin bisa sampe 4000"*), DIRALAT butir 88: *"ASM_ADDRESS, ROAD_NAME saja dilebarin segitu, yg lainnya
-- 100 saja"* - ASM_ADDRESS (Risk Location) dan ROAD_NAME (Address) VARCHAR2(4000) = lebar sumber RISKADDRESS
-- (`DDL\RISKADDRESS.txt`); delapan kolom alamat lain VARCHAR2(100). BUILDING_NO tidak berubah. 186 SUDAH dijalankan di
-- DEV: tidak diubah, migrasi baru. Satu pernyataan per kolom. Amandemen loader `amandemenLebar`. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_RISKLOCATION MODIFY (ASM_ADDRESS VARCHAR2(4000))
/
ALTER TABLE {skema}.T_RISKLOCATION MODIFY (ASM_CITY VARCHAR2(100))
/
ALTER TABLE {skema}.T_RISKLOCATION MODIFY (ASM_DISTRICT VARCHAR2(100))
/
ALTER TABLE {skema}.T_RISKLOCATION MODIFY (ASMRW VARCHAR2(100))
/
ALTER TABLE {skema}.T_RISKLOCATION MODIFY (ASM_ZIP_CODE VARCHAR2(100))
/
ALTER TABLE {skema}.T_PROPERTY MODIFY (ROAD_NAME VARCHAR2(4000))
/
ALTER TABLE {skema}.T_PROPERTY MODIFY (ROAD_TYPE VARCHAR2(100))
/
ALTER TABLE {skema}.T_PROPERTY MODIFY (PROVINCE VARCHAR2(100))
/
ALTER TABLE {skema}.T_PROPERTY MODIFY (COUNTRY VARCHAR2(100))
/
ALTER TABLE {skema}.T_PROPERTY MODIFY (ALM_RISK_ID VARCHAR2(100))
/
