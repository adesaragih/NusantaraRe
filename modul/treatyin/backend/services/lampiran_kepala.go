package services

import "context"

// kepalaLampiran - kepala kontrak pemilik lampiran. Panel Attachment dipakai
// Treaty In DAN Adjustment (`Section/InputTreatyInAdjustment.xml` memuat
// section Attachment yang sama, lampiran menempel pada ID berkas): ID
// kontrak dicari di `TREATY_IN`, bila tidak ada di `TREATY_IN_EDM` — ID
// penyesuaian (`1000080/R01`) hanya tercatat di sana.
func (l *Layanan) kepalaLampiran(ctx context.Context, id string) (map[string]any, bool, error) {
	kepala, ada, err := l.gudang.BacaKepalaTreatyIn(ctx, id)
	if err != nil || ada {
		return kepala, ada, err
	}
	return l.gudang.BacaKepalaPenyesuaian(ctx, id)
}
