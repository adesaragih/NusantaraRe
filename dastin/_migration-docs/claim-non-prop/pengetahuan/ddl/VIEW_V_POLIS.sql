-- SUMBER: DDL_Script_ClaimNonProp.xls (reverse-engineer Toad)
-- UNTUK DIBACA, BUKAN UNTUK DIJALANKAN. Klausa DROP/ALTER-DROP tidak boleh dieksekusi.


  CREATE OR REPLACE FORCE EDITIONABLE VIEW "POOLDATA"."V_POLIS" ("RN", "CASEID", "POLICYNO", "STARTDATETIME", "ENDDATETIME", "SOURCEOFBUSINESS", "SOURCEOFBUSINESSNAME", "CEDINGCO", "CEDINGCONAME", "BUSINESSCODE", "BUSINESSNAME", "CUSTOMERNAME", "SOURCEOFBUSINESSROOT", "PRODUCTIONDATE", "PREMI", "DISCOUNT", "TSI", "TYPE", "PRODKE", "BUSINESSTYPE", "MARKETINGCODE", "CLIENTIDCON", "CLIENTIDORG", "QQ", "SELECTRENEWAL", "RNWSTATUS") AS 
  SELECT ROWNUM,
        a.DATA_JSON.IDNewBisnis AS CASEID,
        a.nopolis AS POLICYNO,
        CASE a.data_json.QuotationData.BusinessFac when 'F' then
                SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.StartDateTime'), 7, 2)
                || '/' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.StartDateTime'), 5, 2)
                || '/' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.StartDateTime'), 1, 4)
                || ' ' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.StartDateTime'), 10, 2)
                || ':' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.StartDateTime'), 12, 2)
                || ':' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.StartDateTime'), 14, 2)
                when 'T' then
                SUBSTR(JSON_VALUE (DATA_JSON, '$.StartDate'), 7, 2)
                || '/' || SUBSTR(JSON_VALUE (DATA_JSON, '$.StartDate'), 5, 2)
                || '/' || SUBSTR(JSON_VALUE (DATA_JSON, '$.StartDate'), 1, 4)
                || ' ' || SUBSTR(JSON_VALUE (DATA_JSON, '$.StartDate'), 10, 2)
        end AS STARTDATETIME,
        CASE a.data_json.QuotationData.BusinessFac when 'F' then
                SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.EndDateTime'), 7, 2)
                || '/' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.EndDateTime'), 5, 2)
                || '/' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.EndDateTime'), 1, 4)
                || ' ' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.EndDateTime'), 10, 2)
                || ':' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.EndDateTime'), 12, 2)
                || ':' || SUBSTR(JSON_VALUE (DATA_JSON, '$.PolicyData.EndDateTime'), 14, 2)
                when 'T' then
                SUBSTR(JSON_VALUE (DATA_JSON, '$.EndDate'), 7, 2)
                || '/' || SUBSTR(JSON_VALUE (DATA_JSON, '$.EndDate'), 5, 2)
                || '/' || SUBSTR(JSON_VALUE (DATA_JSON, '$.EndDate'), 1, 4)
                || ' ' || SUBSTR(JSON_VALUE (DATA_JSON, '$.EndDate'), 10, 2)
        end AS ENDDATETIME,
        getnewid (JSON_VALUE (DATA_JSON, '$.QuotationData.SourceOfBusiness'),'m_agent') AS SOURCEOFBUSINESS,
        get_name_agent (getnewid (JSON_VALUE (DATA_JSON, '$.QuotationData.SourceOfBusiness'),'m_agent')) AS SOURCEOFBUSINESSNAME,
        JSON_VALUE (DATA_JSON, '$.QuotationData.CedingCo') AS CEDINGCO,
        JSON_VALUE (DATA_JSON, '$.QuotationData.CedingCoName') AS CEDINGCONAME,
        getnewid (JSON_VALUE (DATA_JSON, '$.QuotationData.BusinessCode'),'M_BUSINESS') AS BUSINESSCODE,
        GET_BUSINESS_NAME (getnewid (JSON_VALUE (DATA_JSON, '$.QuotationData.BusinessCode'),'M_BUSINESS')) AS BUSINESSNAME,
        JSON_VALUE (DATA_JSON, '$.QuotationData.InsuredName') AS CUSTOMERNAME,
        "POOLDATA"."GET_ROOT_AGEN"("POOLDATA"."GETNEWID"(JSON_VALUE("DATA_JSON" FORMAT JSON , '$.QuotationData.SourceOfBusiness' RETURNING VARCHAR2(4000) NULL ON ERROR),'m_agent')) AS SOURCEOFBUSINESSROOT,
        SYSDATE AS PRODUCTIONDATE,
        a.DATA_JSON.Payment.Premium AS PREMI,
        a.DATA_JSON.Payment.Diskon AS DISCOUNT,
        a.DATA_JSON.SumOfTSI AS TSI,
        a.DATA_JSON.TypeOfPolicy AS TYPEPOLICY,
        a.prodke AS PRODKE,
        CASE
              when a.DATA_JSON.QuotationData.BusinessFac = 'F' AND A.DATA_JSON.IsFacRetro is NULL then 'FACULTATIVE IN'
              WHEN a.DATA_JSON.QuotationData.BusinessFac = 'F' AND A.DATA_JSON.IsFacRetro ='1' then 'RETROCESSION'
              when a.DATA_JSON.QuotationData.BusinessFac = 'T' then 'TREATY'
        end as BusinessType,
        getnewid (JSON_VALUE (DATA_JSON, '$.QuotationData.MarketingCode'),'M_MARKETINGOFFICER') AS MARKETINGCODE,
        GETNEWID (JSON_VALUE (DATA_JSON, '$.CIFData.Customer_P.ASMClientID'),'m_client') AS CLIENTIDCON,
        GETNEWID (TRIM (JSON_VALUE (DATA_JSON, '$.CIFData.Customer_C.ASMClientID')),'m_client') AS CLIENTIDORG,
        A.DATA_JSON.QQName AS QQ,
        'false' as SELECTRENEWAL,
        NVL (JSON_VALUE ("DATA_JSON" FORMAT JSON , '$.RenewalStatus'),' ') AS RNWSTATUS
    FROM json_polis a