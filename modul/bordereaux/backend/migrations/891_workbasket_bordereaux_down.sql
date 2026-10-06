-- Mundur 891 - pemegang workbasket Bordereaux lalu ketiga barisnya; untuk skema uji.
DELETE FROM {skema}.M_LOGIN_GO_WORKBASKET
WHERE WORKBASKET_ID IN ('ReasBordereauxAdmin', 'ReasBordereauxChecker', 'ReasBordereauxSupervisor')
/
DELETE FROM {skema}.M_WORKBASKET
WHERE WORKBASKET_ID IN ('ReasBordereauxAdmin', 'ReasBordereauxChecker', 'ReasBordereauxSupervisor')
/
