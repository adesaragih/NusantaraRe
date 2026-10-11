package handlers

// Tab Payment kasus FIRE (tiket 51): GET / PUT /api/nbfacin/kasus/{caseId}/payment dan POST …/payment/hitung. Kontrak
// frontend modul/nbfacin/frontend/api.ts `TampilanPayment` / `MataUangPayment` / `CicilanPayment` / `AksiPayment`.
// Due Date kabel `DD-MM-YYYY`, tersimpan teks Pega `YYYYMMDD`. % cicilan berskala sampai 20 (contoh Pega
// "33.33333333333333333333"), jadi diurai dengan pola sendiri, bukan UraiDesimalIsian (8 desimal).

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// batasBadanPayment - batas bita badan POST / PUT.
const batasBadanPayment = 1 << 20

// polaPersenCicilan - 0-3 digit bulat, sampai 20 desimal.
var polaPersenCicilan = regexp.MustCompile(`^[0-9]{1,3}(\.[0-9]{1,20})?$`)

type cicilanPaymentKabel struct {
	InstallmentNo         string `json:"installmentNo"`
	DueDate               string `json:"dueDate"`
	InstallmentPercentage string `json:"installmentPercentage"`
	Premium               string `json:"premium"`
	RICommision           string `json:"riCommision"`
	Deduction2            string `json:"deduction2"`
	BrokerageFee          string `json:"brokerageFee"`
	PPh                   string `json:"pph"`
	PPN                   string `json:"ppn"`
	Discount              string `json:"discount"`
	PaymentTotal          string `json:"paymentTotal"`
}

type mataUangPaymentKabel struct {
	Currency        string                `json:"currency"`
	Premium         string                `json:"premium"`
	SumTotalPayment string                `json:"sumTotalPayment"`
	Cicilan         []cicilanPaymentKabel `json:"cicilan"`
}

type tampilanPaymentKabel struct {
	Installment string                 `json:"installment"`
	TarifPph    string                 `json:"tarifPph"`
	TarifPpn    string                 `json:"tarifPpn"`
	MataUang    []mataUangPaymentKabel `json:"mataUang"`
	Pesan       []string               `json:"pesan"`
}

// tanggalKeKabel - `YYYYMMDD` -> `DD-MM-YYYY` (teks lain apa adanya).
func tanggalKeKabel(t string) string {
	if d, err := time.Parse("20060102", t); err == nil {
		return d.Format("02-01-2006")
	}
	return t
}

func kePaymentKabel(d services.TampilanPayment) tampilanPaymentKabel {
	out := tampilanPaymentKabel{Installment: d.Installment, TarifPph: teks(d.TarifPPH), TarifPpn: teks(d.TarifPPN),
		MataUang: []mataUangPaymentKabel{}, Pesan: append([]string{}, d.Pesan...)}
	for _, m := range d.MataUang {
		mk := mataUangPaymentKabel{Currency: m.Baris.Currency, Premium: teks(m.Baris.Premium), SumTotalPayment: teks(m.SumTotalPayment),
			Cicilan: []cicilanPaymentKabel{}}
		for _, c := range m.Cicilan {
			no := ""
			if c.InstallmentNo > 0 {
				no = strconv.Itoa(c.InstallmentNo)
			}
			persen := ""
			if c.Persen != nil {
				persen = utils.FormatDecimal(c.Persen) // apa adanya, tanpa pembulatan 8 desimal
			}
			mk.Cicilan = append(mk.Cicilan, cicilanPaymentKabel{InstallmentNo: no, DueDate: tanggalKeKabel(c.DueDate),
				InstallmentPercentage: persen, Premium: teks(c.Premium), RICommision: teks(c.RICommision), Deduction2: teks(c.Deduction2),
				BrokerageFee: teks(c.BrokerageFee), PPh: teks(c.PPh), PPN: teks(c.PPN), Discount: teks(c.Discount),
				PaymentTotal: teks(c.PaymentTotal)})
		}
		out.MataUang = append(out.MataUang, mk)
	}
	return out
}

type isianCicilanKabel struct {
	InstallmentNo         string `json:"installmentNo"`
	DueDate               string `json:"dueDate"`
	InstallmentPercentage string `json:"installmentPercentage"`
}

type isianMataUangPaymentKabel struct {
	Currency string              `json:"currency"`
	Cicilan  []isianCicilanKabel `json:"cicilan"`
}

type isianPaymentKabel struct {
	Aksi        *string                     `json:"aksi"`
	Currency    *string                     `json:"currency"`
	Idx         *int                        `json:"idx"`
	Installment string                      `json:"installment"`
	MataUang    []isianMataUangPaymentKabel `json:"mataUang"`
}

// uraiIsianPayment - JSON (kunci asing -> 400) + Due Date + persen; galat sudah ditulis bila ok=false.
func uraiIsianPayment(w http.ResponseWriter, r *http.Request, denganAksi bool) (isianPaymentKabel, services.IsianPayment, bool) {
	var b isianPaymentKabel
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanPayment))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&b); err != nil || b.MataUang == nil {
		galat.Tulis(w, http.StatusBadRequest, `badan permintaan harus {"installment", "mataUang": [{currency, cicilan: [{installmentNo, dueDate, installmentPercentage}]}]}`)
		return b, services.IsianPayment{}, false
	}
	aksiLengkap := b.Aksi != nil && b.Currency != nil && b.Idx != nil
	if denganAksi != aksiLengkap || (!denganAksi && (b.Aksi != nil || b.Currency != nil || b.Idx != nil)) {
		pesan := "badan PUT tidak memuat aksi / currency / idx"
		if denganAksi {
			pesan = "badan hitung wajib memuat aksi, currency, dan idx"
		}
		galat.Tulis(w, http.StatusBadRequest, pesan)
		return b, services.IsianPayment{}, false
	}
	var masalah []string
	isi := services.IsianPayment{Installment: b.Installment, MataUang: make([]services.IsianCicilanMataUang, 0, len(b.MataUang))}
	for i, m := range b.MataUang {
		mu := services.IsianCicilanMataUang{Currency: strings.TrimSpace(m.Currency), Cicilan: make([]models.CicilanPayment, 0, len(m.Cicilan))}
		for j, c := range m.Cicilan {
			p := fmt.Sprintf("mataUang[%d].cicilan[%d].", i, j)
			cp := models.CicilanPayment{}
			if no := strings.TrimSpace(c.InstallmentNo); no != "" {
				n, err := strconv.Atoi(no)
				if err != nil {
					masalah = append(masalah, p+"installmentNo harus angka bulat")
				}
				cp.InstallmentNo = n
			}
			if tgl := strings.TrimSpace(c.DueDate); tgl != "" {
				d, err := time.Parse("02-01-2006", tgl)
				if err != nil {
					masalah = append(masalah, p+"dueDate harus DD-MM-YYYY")
				}
				cp.DueDate = d.Format("20060102")
			}
			if ps := strings.TrimSpace(c.InstallmentPercentage); ps != "" {
				v, err := utils.ParseDecimal(ps)
				if !polaPersenCicilan.MatchString(ps) || err != nil {
					masalah = append(masalah, p+"installmentPercentage harus desimal 0-999 dengan sampai 20 angka di belakang titik")
				}
				cp.Persen = v
			}
			mu.Cicilan = append(mu.Cicilan, cp)
		}
		isi.MataUang = append(isi.MataUang, mu)
	}
	if len(masalah) > 0 {
		tulisGalat(w, fmt.Errorf("%w: %s", services.ErrMasukanPayment, strings.Join(masalah, "; ")))
		return b, services.IsianPayment{}, false
	}
	return b, isi, true
}

// bacaPayment - GET /api/nbfacin/kasus/{caseId}/payment.
func bacaPayment(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.TampilanPayment(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, kePaymentKabel(d))
	}
}

// hitungPayment - POST …/payment/hitung {aksi, currency, idx, …}: tanpa simpan.
func hitungPayment(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, isi, ok := uraiIsianPayment(w, r, true)
		if !ok {
			return
		}
		d, err := svc.HitungPayment(r.Context(), r.PathValue("caseId"), *b.Aksi, strings.TrimSpace(*b.Currency), *b.Idx, isi)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, kePaymentKabel(d))
	}
}

// simpanPayment - PUT …/payment (badan tanpa aksi / currency / idx): hitung, proteksi, simpan, jawab baca ulang.
func simpanPayment(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, isi, ok := uraiIsianPayment(w, r, false)
		if !ok {
			return
		}
		d, err := svc.SimpanPayment(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), isi)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, kePaymentKabel(d))
	}
}
