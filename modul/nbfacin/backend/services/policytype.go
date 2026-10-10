package services

import (
	"context"
	"errors"

	"nusantarare/modul/nbfacin/backend/repository"
)

// Policy Type (`.QuotationData.PolicyType`) - keputusan work owner 09-10-2026 "COBA CEK PolicyType, YANG DISIMPAN ITU STANDARD
// VALUENYA. ( 0,1,2)", meralat butir 78.2 ("Label text for now"). `[terverifikasi]` D:\migrasi\RNM\DDL\PolicyType.xml
// (ASM-FW-GISFW-DATA-QUOTATION!POLICYTYPE), dipasangkan per rowdata Embed-PropertyPromtValues: pyStandardValue 0 = "Individual
// Policy", 1 = "Master Policy", 2 = "Declaration Policy" (pyLocalizedValue). Data lama T_QUOTATIONDATA.POLICY_TYPE berlabel
// dipetakan saat DIBACA (tanpa migrasi UPDATE).

// policyTypeMaster - standard value Master Policy (ProtectPremiPolicy_Act `PolicyType==1`).
const policyTypeMaster = "1"

// labelPolicyType - pyLocalizedValue -> pyStandardValue.
var labelPolicyType = map[string]string{"Individual Policy": "0", "Master Policy": policyTypeMaster, "Declaration Policy": policyTypeDeklarasi}

// policyTypeDeklarasi - kode Declaration Policy (When IsDeclarationPolicy `Quotation.PolicyType = 2`).
const policyTypeDeklarasi = "2"

// kodePolicyTypeSah - isian PUT …/general yang diterima ("" = kosong).
var kodePolicyTypeSah = map[string]bool{"": true, "0": true, policyTypeMaster: true, "2": true}

// policyTypeKasus - kode Policy Type case (T_QUOTATIONDATA.POLICY_TYPE lewat kodePolicyType); tanpa pembaca kasus -> "".
func (s *Service) policyTypeKasus(ctx context.Context, id string) (string, error) {
	if s.kasus == nil {
		return "", nil
	}
	k, err := s.kasus.BacaKasus(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return "", ErrKasusTidakAda
	}
	if err != nil {
		return "", err
	}
	return kodePolicyType(k.General.PolicyType), nil
}

// kodePolicyType - label -> kode; kode (atau nilai lain) apa adanya.
func kodePolicyType(s string) string {
	if k, ada := labelPolicyType[s]; ada {
		return k
	}
	return s
}
