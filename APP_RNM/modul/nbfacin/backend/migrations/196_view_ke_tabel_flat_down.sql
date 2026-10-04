-- Jalur mundur 196 - tabel flat dikembalikan menjadi VIEW warisan, teks view PERSIS `D:\migrasi\RNM\DDL\<NAMA>.txt`
-- (POOLDATA. diganti {skema}.; rujukan M_* di dalamnya apa adanya). Isi tabel flat dibuang - view membaca M_* lagi.
-- Urutan terbalik dari maju.
DROP TABLE {skema}.ACCUMULATION
/
CREATE OR REPLACE FORCE VIEW {skema}.ACCUMULATION
(
   ID,
   ACCUMULATION,
   ACCUMULATIONNAME,
   NOTE,
   KEYWORD,
   SCOPEAREA,
   CZONE,
   CZONEID,
   PROVINCE,
   PROVINCEID,
   ZIPCODE,
   ACCUMULATIONTYPE,
   SYARIAHSTATUS
)
AS
   SELECT DISTINCT b.JSONDATA.ID,
                   b.JSONDATA.Accumulation,
                   (SELECT c.JSONDATA.AccumulationType
                      FROM m_accumulatedtype c
                     WHERE b.JSONDATA.Accumulation = c.ID)
                      AS AccumulationName,
                   UPPER (b.JSONDATA.Note),
                   b.JSONDATA.Keyword,
                   b.JSONDATA.ScopeArea,
                   b.JSONDATA.CZone,
                   b.JSONDATA.CZoneID,
                   b.JSONDATA.ProvinceName,
                   b.JSONDATA.ProvinceID,
                   b.JSONDATA.PostalCode,
                   b.JSONDATA.AccumulationType,
                   b.JSONDATA.SyariahStatus
     FROM m_accumulation b
    WHERE b.JSONDATA.IsActive IS NULL
/
DROP TABLE {skema}.CZONE
/
CREATE OR REPLACE FORCE VIEW {skema}.CZONE
(
   ID,
   CODE,
   DESCRIPTION,
   GROUPOF,
   GROUPOFNAME,
   TGLUPDATE,
   USERID
)
AS
     SELECT a.JSONDATA.ID,
            a.JSONDATA.Code,
            a.JSONDATA.Description,
            a.JSONDATA.GroupOf,
            (SELECT b.JSONDATA.Description
               FROM M_CZONE b
              WHERE b.JSONDATA.Code = a.JSONDATA.GroupOf)
               AS GroupOfName,
            a.JSONDATA.TglUpdate,
            a.JSONDATA.UserID
       FROM m_czone a
   ORDER BY GroupOf ASC
/
DROP TABLE {skema}.ACCUMULATEDTYPE
/
CREATE OR REPLACE FORCE VIEW {skema}.ACCUMULATEDTYPE
(
   ID,
   ACCUMULATIONTYPE,
   KEYWORD,
   NOTE,
   TYPE
)
AS
   SELECT a.JSONDATA.ID,
          a.JSONDATA.AccumulationType,
          a.JSONDATA.Keyword,
          a.JSONDATA.Note,
          a.JSONDATA.TYPE
     FROM m_accumulatedtype a
/
DROP TABLE {skema}.PROVINCE
/
CREATE OR REPLACE FORCE VIEW {skema}.PROVINCE
(
   ID,
   NATIONID,
   NOTE,
   NATIONNAME
)
AS
   SELECT a.JSONDATA.ID AS ID,
          a.JSONDATA.NationID AS NationID,
          a.JSONDATA.Note AS Note,
          (SELECT b.JSONDATA.Note
             FROM M_Nation b
            WHERE a.JSONDATA.NationID = b.ID)
             AS NationName
     FROM M_PROVINCE a
/
