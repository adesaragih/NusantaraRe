-- Jalur mundur 199 - tab Inw Fac Cedant Panels: daftar cedant, Share of Ceding, dan Share Cedant Type ikut hilang.
DROP SEQUENCE {skema}.SEQ_T_CEDINGCEDANTLIST
/
DROP TABLE {skema}.T_CEDINGCEDANTLIST CASCADE CONSTRAINTS
/
ALTER TABLE {skema}.T_QUOTATIONDATA DROP (SHARE_OF_CEDING)
/
ALTER TABLE {skema}.T_GENERAL_POLIS DROP (SHARE_CEDANT_TYPE)
/
