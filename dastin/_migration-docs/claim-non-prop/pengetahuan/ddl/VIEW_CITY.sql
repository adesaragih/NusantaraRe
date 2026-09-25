-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE FORCE EDITIONABLE VIEW "POOLDATA"."CITY" ("ID", "PROVINCEID", "NOTE", "BRANCHID", "EMAIL", "MOID", "JABODETABEKSTATUS", "ZIPCODE", "PROVINCENAME", "BRANCHNAME") AS 
  SELECT a.ID,
          a.ProvinceID,
          a.Note,
          a.BranchID,
          a.Email,
          a.MOID,
          a.JABODETABEKSTATUS,
          c.ZipCode AS ZIPCODE,
          d.NOTE AS PROVINCENAME,
          e.NAME AS BRANCHNAME
     FROM CITYINPUT a
          LEFT JOIN DISTRICTINPUT b
             ON b.CITYID = a.ID
          LEFT JOIN RWINPUT c
             ON c.DistrictID = b.ID
          LEFT JOIN PROVINCE d
             ON a.ProvinceID = d.ID
          LEFT JOIN BRANCH e
             ON a.BranchID = e.ID