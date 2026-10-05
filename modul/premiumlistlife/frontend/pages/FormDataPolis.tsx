// Data polis layar Input Premium Detail — tiket 03 bagian 2 (bagian 1–3 layar).
//
// Meniru tiga kolom atas `Section/ShowLifePremiumDetail.xml`: kiri data polis
// (produk, Type, R/I SLIP, Premium Payment Method, Marketing Officer, Annuity
// Interest, Premium Refund Factor), tengah pihak (Ceding, Policy Holder, Billing
// Name, Retrocessionaire), kanan tanggal-tanggal penawaran, WPC, dan status.
//
// ⛔ Data penawaran (tanggal, Age Limit, Coverage Period, Sum Insured,
// Underwriting Policy, Marketing Note, dan seterusnya) DITAMPILKAN, tidak
// diisi di sini — ia milik layar Input Offer. Insured Name dan Occupation tetap
// disembunyikan, sama seperti di layar Input Offer (keputusan work owner).
//
// ⛔ `Save Data` HANYA menyimpan dan memeriksa medan wajib (keputusan work owner
// 05-10-2026). `Calculate1_Act` tidak dijalankan; hitung summary berjalan di
// Calculate CSV, dan batas produk `SavePremiumList_Act` menjadi penolakan
// Validate CSV.

import { useCallback, useEffect, useState } from 'react'

import { Field, Gagal, Memuat, Modal, Pilih, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilDataPolis,
  ambilPenawaranPolis,
  cariCedingPolis,
  cariMarketingPolis,
  cariProdukPolis,
  cariRISlipPolis,
  isiDariDataPolis,
  simpanDataPolis,
  type DataPolis,
  type IsiDataPolis,
  type PenawaranPolis,
  type PilihanKode,
} from '../api'
import { KOLOM_PRODUK, LABEL_DATA_POLIS, RINCIAN_PRODUK, TEKS_PILIH, TEKS_TOMBOL_PILIH } from '../labels'
import AreaTeks from '../components/AreaTeks'
import IsianTanggal from '../components/IsianTanggal'
import ModalRincianProduk from '../components/ModalRincianProduk'
import { tanggalTampil } from '../tanggal'
import '../premiumlistlife.css'

/** Type yang menuntut R/I SLIP dan Billing Name (`.Type == 'TR' || 'TP'`). */
export function typeRetro(type: string): boolean {
  return type === 'TP' || type === 'TR'
}

/**
 * Kolom wajib yang masih kosong — label layar, urut seperti di layar.
 *
 * `pyRequired` korpus: Product Name (ProtectAccept), Type, Premium Payment
 * Method, Marketing Officer, Annuity Interest, Premium Refund Factor; R/I SLIP
 * dan Billing Name hanya untuk TP/TR. Server memeriksa ulang dengan aturan
 * yang sama (`models.SusunDataPolis`).
 */
export function kolomWajibDataPolis(isi: IsiDataPolis): string[] {
  const periksa: [boolean, string][] = [
    [isi.productName.trim() === '', LABEL_DATA_POLIS.productName],
    [isi.type.trim() === '', LABEL_DATA_POLIS.type],
    [typeRetro(isi.type) && isi.riSlipRnm.trim() === '', LABEL_DATA_POLIS.riSlip],
    [isi.proRateType.trim() === '', LABEL_DATA_POLIS.proRateType],
    [isi.marketingName.trim() === '', LABEL_DATA_POLIS.marketing],
    [isi.annuityInterest.trim() === '', LABEL_DATA_POLIS.annuityInterest],
    [isi.premiumRefundFactor.trim() === '', LABEL_DATA_POLIS.premiumRefundFactor],
    [typeRetro(isi.type) && isi.retroName.trim() === '', LABEL_DATA_POLIS.billing],
    [typeRetro(isi.type) && isi.securityReinsurer.trim() === '', LABEL_DATA_POLIS.retro],
    // Wajib di layar ini — keputusan work owner 02-10-2026.
    [isi.dateReceived.trim() === '', LABEL_DATA_POLIS.dateReceived],
  ]
  return periksa.filter(([kosong]) => kosong).map(([, label]) => label)
}

function opsi(p: PilihanKode[]): Opsi[] {
  return p.map((x) => ({ value: x.kode, label: x.nama }))
}

/** Popup yang sedang terbuka. */
type Popup = 'produk' | 'marketing' | 'rislip' | 'billing' | 'retro'

/** Satu baris popup — apa yang tampil dan apa yang diterapkan. */
interface BarisPopup {
  kunci: string
  sel: string[]
  terapkan: (isi: IsiDataPolis) => IsiDataPolis
}

export default function FormDataPolis({
  polisID,
  bernomor,
  onTersimpan,
  onBelumTersimpan,
}: {
  polisID: string
  /** PL_NUMBER sudah terbit — `Choose Product Name` tampil hanya bila belum. */
  bernomor: boolean
  /** Dipanggil sesudah tersimpan — kepala halaman (Type) ikut dimuat ulang. */
  onTersimpan: () => void
  /**
   * Dilapori `true` selama isian di layar BERBEDA dari yang terakhir tersimpan
   * — Confirm dikunci selama itu (keputusan work owner 03-10-2026).
   */
  onBelumTersimpan?: (belum: boolean) => void
}) {
  const [data, setData] = useState<DataPolis | null>(null)
  const [offer, setOffer] = useState<PenawaranPolis | null>(null)
  const [isi, setIsi] = useState<IsiDataPolis | null>(null)
  const [galatMuat, setGalatMuat] = useState<unknown>(null)
  const [galatSimpan, setGalatSimpan] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [tersimpan, setTersimpan] = useState(false)
  // Jawaban Save Data: Type / Product Name berganti → rekap summary dihapus (05-10-2026).
  const [rekapDihapus, setRekapDihapus] = useState(false)
  const [popup, setPopup] = useState<Popup | null>(null)
  // Popup isi Product Name - tombol View (05-10-2026).
  const [lihatProduk, setLihatProduk] = useState(false)
  // Isian terakhir yang TERSIMPAN (dimuat atau sesudah Save Data), sebagai teks -
  // pembanding "ada perubahan yang belum disimpan".
  const [dasar, setDasar] = useState('')
  const belumTersimpan = isi !== null && dasar !== '' && JSON.stringify(isi) !== dasar
  useEffect(() => {
    onBelumTersimpan?.(belumTersimpan)
  }, [belumTersimpan, onBelumTersimpan])
  // Layar ditinggalkan: tidak ada lagi perubahan yang menahan Confirm.
  useEffect(() => () => onBelumTersimpan?.(false), [onBelumTersimpan])

  const terima = useCallback((d: DataPolis) => {
    setData(d)
    setIsi(isiDariDataPolis(d))
    setDasar(JSON.stringify(isiDariDataPolis(d)))
  }, [])

  useEffect(() => {
    let hidup = true
    void (async () => {
      try {
        const [d, o] = await Promise.all([ambilDataPolis(polisID), ambilPenawaranPolis(polisID)])
        if (!hidup) return
        terima(d)
        setOffer(o)
      } catch (e) {
        if (hidup) setGalatMuat(e)
      }
    })()
    return () => {
      hidup = false
    }
  }, [polisID, terima])

  if (galatMuat !== null) return <Gagal galat={galatMuat} />
  if (data === null || isi === null || offer === null) return <Memuat />

  const bisa = data.bolehDisimpan
  const kurang = kolomWajibDataPolis(isi)
  const ubah = (medan: keyof IsiDataPolis) => (v: string) => {
    setIsi({ ...isi, [medan]: v })
    setTersimpan(false)
    setRekapDihapus(false)
  }
  const tampil = (label: string, nilai: string) => (
    <Field key={label} label={label} value={nilai} onChange={() => {}} readOnly />
  )

  async function simpan(): Promise<void> {
    if (isi === null || sibuk || kolomWajibDataPolis(isi).length > 0) return
    setSibuk(true)
    setGalatSimpan(null)
    setTersimpan(false)
    setRekapDihapus(false)
    try {
      const hasil = await simpanDataPolis(polisID, isi)
      terima(hasil)
      setRekapDihapus(hasil.rekapDihapus === true)
      setTersimpan(true)
      onTersimpan()
    } catch (e) {
      setGalatSimpan(e)
    } finally {
      setSibuk(false)
    }
  }

  // Tombol pilih MENEMPEL di kanan kotak isiannya (`pl-dp-pilih`), jadi kotak
  // + tombol selebar satu sel grid seperti isian lain. Teks tombol pendek
  // ("Choose"); nama lengkapnya di `aria-label` dan `title`.
  const tombol = (label: string, jenis: Popup) =>
    bisa && (
      <button
        type="button"
        className="btn btn--ghost"
        aria-label={label}
        title={label}
        onClick={() => { setPopup(jenis) }}
      >
        {TEKS_TOMBOL_PILIH}
      </button>
    )

  return (
    <section className="panel pl-datapolis">
      <h3 className="panel__title">{LABEL_DATA_POLIS.judul}</h3>
      {/*
        TIGA KOLOM BERDAMPINGAN (permintaan work owner 02-10-2026: "3 baris ke
        samping seperti sebelumnya, namun lebih rapih"): kiri Policy Data,
        tengah Parties & Coverage (+ Retrocession untuk TP/TR), kanan Dates &
        Status. SETIAP kolom satu grid DUA sel sama lebar (premiumlistlife.css);
        isian teks panjang dan isian bertombol pilih selebar kolom
        (`pl-dp-lebar`), sehingga semua kotak lurus. Urutan isian tetap urutan
        sel `ShowLifePremiumDetail`.
      */}
      <div className="pl-dp-kolom-tiga">
        <div className="pl-dp-bagian">
          <h4 className="pl-offer__subjudul">Policy Data</h4>
          <div className="pl-dp-grid">
            <div className={!bernomor && bisa ? 'pl-dp-pilih pl-dp-pilih--dua pl-dp-lebar' : 'pl-dp-pilih pl-dp-lebar'}>
              <Field label={LABEL_DATA_POLIS.productName} value={isi.productName} onChange={() => {}} readOnly required />
              {/* `Choose Product Name` tampil bila PL_NUMBER belum ada. */}
              {!bernomor && tombol(LABEL_DATA_POLIS.pilihProduk, 'produk')}
              {/* View: isi Product Name terpilih, juga sesudah bernomor (05-10-2026). */}
              <button
                type="button"
                className="btn btn--ghost"
                aria-label={RINCIAN_PRODUK.namaTombol}
                title={RINCIAN_PRODUK.namaTombol}
                disabled={isi.productNameId.trim() === ''}
                onClick={() => { setLihatProduk(true) }}
              >
                {RINCIAN_PRODUK.tombol}
              </button>
            </div>
            {tampil(LABEL_DATA_POLIS.productNameId, isi.productNameId)}
            <Pilih
              kosong={TEKS_PILIH}
              label={LABEL_DATA_POLIS.type}
              value={isi.type}
              onChange={ubah('type')}
              opsi={opsi(data.pilihan.type)}
              required
            />
            {tampil(LABEL_DATA_POLIS.typeCeding, offer.typeCedingName)}
            <Pilih
              kosong={TEKS_PILIH}
              label={LABEL_DATA_POLIS.proRateType}
              value={isi.proRateType}
              onChange={ubah('proRateType')}
              opsi={opsi(data.pilihan.proRateType)}
              required
            />
            <div className="pl-dp-pilih pl-dp-lebar">
              <Field label={LABEL_DATA_POLIS.marketing} value={isi.marketingName} onChange={() => {}} readOnly required />
              {tombol(LABEL_DATA_POLIS.marketing, 'marketing')}
            </div>
            {tampil(LABEL_DATA_POLIS.sumInsured, offer.sumInsured)}
            {offer.noOffer !== '' && tampil(LABEL_DATA_POLIS.noOffer, offer.noOffer)}
            {/* Desimal sebagai TEKS, titik sebagai pemisah — bukan input number. */}
            <Field
              label={LABEL_DATA_POLIS.annuityInterest}
              value={isi.annuityInterest}
              onChange={ubah('annuityInterest')}
              readOnly={!bisa}
              required
            />
            <Field
              label={LABEL_DATA_POLIS.premiumRefundFactor}
              value={isi.premiumRefundFactor}
              onChange={ubah('premiumRefundFactor')}
              readOnly={!bisa}
              required
            />
            <div className="pl-dp-lebar">
              <AreaTeks label={LABEL_DATA_POLIS.ketentuanUnderwriting} value={offer.ketentuanUnderwriting} readOnly />
            </div>
          </div>
        </div>

        {/* Kolom tengah: pihak polis dan cakupan (dibaca dari Input Offer), lalu Retrocession. */}
        <div>
          <div className="pl-dp-bagian">
            <h4 className="pl-offer__subjudul">Parties &amp; Coverage</h4>
            <div className="pl-dp-grid">
              <div className="pl-dp-lebar">{tampil(LABEL_DATA_POLIS.ceding, isi.cedingCoName)}</div>
              <div className="pl-dp-lebar">{tampil(LABEL_DATA_POLIS.policyHolder, isi.policyHolderName)}</div>
              {tampil(LABEL_DATA_POLIS.jenisAsuransi, offer.jenisAsuransi)}
              {tampil(LABEL_DATA_POLIS.businessCode, offer.businessName)}
              {tampil(LABEL_DATA_POLIS.batasUsia, offer.batasUsiaPeserta === null ? '' : String(offer.batasUsiaPeserta))}
              {tampil(LABEL_DATA_POLIS.periode, offer.periodePertanggungan)}
            </div>
          </div>

          {/*
            Billing Name dan Retrocessionaire HANYA untuk Type TP/TR (keputusan
            work owner 01-10-2026) — selain itu disembunyikan, dan server
            mengosongkannya saat Save Data (`models.SusunDataPolis`). R/I SLIP
            RNM No. di bawah catatan dan di atas Billing Name (permintaan work
            owner 02-10-2026).
          */}
          {typeRetro(isi.type) && (
            <div className="pl-dp-bagian">
              <h4 className="pl-offer__subjudul">Retrocession</h4>
              <p className="pl-datapolis__catatan" role="note">
                {LABEL_DATA_POLIS.catatanBilling}
              </p>
              <div className="pl-dp-grid">
                <div className="pl-dp-pilih pl-dp-lebar">
                  <Field label={LABEL_DATA_POLIS.riSlip} value={isi.riSlipRnm} onChange={() => {}} readOnly required />
                  {tombol(LABEL_DATA_POLIS.riSlip, 'rislip')}
                </div>
                <div className="pl-dp-pilih pl-dp-lebar">
                  <Field label={LABEL_DATA_POLIS.billing} value={isi.retroName} onChange={() => {}} readOnly required />
                  {tombol(LABEL_DATA_POLIS.pilihBilling, 'billing')}
                </div>
                {/* Wajib untuk TP/TR — keputusan work owner 01-10-2026. */}
                <div className="pl-dp-pilih pl-dp-lebar">
                  <Field label={LABEL_DATA_POLIS.retro} value={isi.securityReinsurer} onChange={() => {}} readOnly required />
                  {tombol(LABEL_DATA_POLIS.pilihRetro, 'retro')}
                </div>
              </div>
            </div>
          )}
        </div>

        <div className="pl-dp-bagian">
          <h4 className="pl-offer__subjudul">Dates &amp; Status</h4>
          <div className="pl-dp-grid">
            {/* DAPAT DIISI dan WAJIB di layar ini (keputusan work owner 02-10-2026). */}
            <IsianTanggal
              label={LABEL_DATA_POLIS.dateReceived}
              value={isi.dateReceived}
              onChange={ubah('dateReceived')}
              readOnly={!bisa}
              required
            />
            {tampil(LABEL_DATA_POLIS.tanggalPenawaran, tanggalTampil(offer.tanggalPenawaran))}
            {tampil(LABEL_DATA_POLIS.tanggalRespon, tanggalTampil(offer.tanggalRespon))}
            {tampil(LABEL_DATA_POLIS.tanggalKonfirmasi, tanggalTampil(offer.tanggalKonfirmasi))}
            {tampil(LABEL_DATA_POLIS.tanggalRealisasi, tanggalTampil(offer.tanggalRealisasi))}
            {tampil(LABEL_DATA_POLIS.tanggalKonfirmasiBalik, tanggalTampil(offer.tanggalKonfirmasiBalik))}
            {tampil(LABEL_DATA_POLIS.tanggalBind, tanggalTampil(offer.tanggalBind))}
            {offer.tbc !== null && tampil(LABEL_DATA_POLIS.tanggalTbc, tanggalTampil(offer.tanggalTbc))}
            {tampil(LABEL_DATA_POLIS.wpc, tanggalTampil(data.wpc))}
            {/*
              `Status` penawaran SENGAJA tidak ditampilkan di sini (keputusan
              work owner 02-10-2026) — tetap diisi dan dilihat di Input Offer
              Life.
            */}
            <div className="pl-dp-lebar">{tampil(LABEL_DATA_POLIS.statusUpdate, offer.statusUpdate)}</div>
            <div className="pl-dp-lebar"><AreaTeks label={LABEL_DATA_POLIS.keteranganMarketing} value={offer.keteranganMarketing} readOnly /></div>
          </div>
        </div>
      </div>

      {galatSimpan !== null && <Gagal galat={galatSimpan} />}
      {/*
        Type / Product Name berganti → rekap summary dihapus server; Confirm menolak polis tanpa
        rekap sampai Calculate CSV dijalankan ulang (keputusan work owner 05-10-2026).
        Panel Summary dimuat ulang lewat `onTersimpan` (versiSummary).
      */}
      {rekapDihapus && (
        <div className="alert alert--error" role="alert">
          {LABEL_DATA_POLIS.rekapDihapus}
        </div>
      )}
      {bisa && (
        <div className="pl-offer__aksi">
          <button
            type="button"
            className="btn btn--primary"
            disabled={sibuk || kurang.length > 0}
            onClick={() => { void simpan() }}
          >
            {LABEL_DATA_POLIS.simpan}
          </button>
          {kurang.length > 0 && (
            <span className="pl-offer__kurang" role="status">
              Required: {kurang.join(', ')}
            </span>
          )}
          {tersimpan && kurang.length === 0 && (
            <span className="pl-offer__tersimpan" role="status">
              Policy data saved.
            </span>
          )}
        </div>
      )}

      {lihatProduk && (
        <ModalRincianProduk produkID={isi.productNameId} onTutup={() => { setLihatProduk(false) }} />
      )}
      {popup !== null && (
        <PopupCari
          judul={judulPopup(popup)}
          kolom={kolomPopup(popup)}
          cari={(teks) => cariPopup(popup, polisID, teks)}
          onPilih={(b) => {
            setIsi(b.terapkan(isi))
            setTersimpan(false)
            setPopup(null)
          }}
          onTutup={() => { setPopup(null) }}
        />
      )}
    </section>
  )
}

function judulPopup(p: Popup): string {
  return {
    produk: LABEL_DATA_POLIS.pilihProduk,
    marketing: LABEL_DATA_POLIS.marketing,
    rislip: LABEL_DATA_POLIS.riSlip,
    billing: LABEL_DATA_POLIS.pilihBilling,
    retro: LABEL_DATA_POLIS.pilihRetro,
  }[p]
}

function kolomPopup(p: Popup): string[] {
  if (p === 'produk') {
    return [KOLOM_PRODUK.id, KOLOM_PRODUK.inwardName, KOLOM_PRODUK.ceding, KOLOM_PRODUK.sob, KOLOM_PRODUK.policyHolder]
  }
  if (p === 'rislip') return [LABEL_DATA_POLIS.riSlip]
  return ['ID', 'Name']
}

/**
 * Pencarian per popup, beserta penerapan pilihannya:
 * produk → `SetProdNametoPolis` (produk, SOB, Ceding, Policy Holder);
 * marketing → `.ID`→MOID, `.ClientID`→MarketingCode, `.ClientName`→nama;
 * billing/retro → BrowseCedingCoLife_RD.
 */
export async function cariPopup(p: Popup, polisID: string, teks: string): Promise<BarisPopup[]> {
  switch (p) {
    case 'produk':
      return (await cariProdukPolis(polisID, teks)).map((r) => ({
        kunci: r.id,
        sel: [r.id, r.inwardName, r.ceding, r.sobName, r.policyHolderName],
        terapkan: (i) => ({
          ...i,
          productNameId: r.id,
          productName: r.inwardName,
          sourceOfBusiness: r.sobId,
          sobName: r.sobName,
          cedingCo: r.cedingId,
          cedingCoName: r.ceding,
          policyHolder: r.policyHolder,
          policyHolderName: r.policyHolderName,
        }),
      }))
    case 'marketing':
      return (await cariMarketingPolis(teks)).map((r) => ({
        kunci: r.id,
        sel: [r.id, r.nama],
        terapkan: (i) => ({ ...i, moId: r.id, marketingCode: r.kode, marketingName: r.nama }),
      }))
    case 'rislip':
      return (await cariRISlipPolis(teks)).map((r) => ({
        kunci: r.id,
        sel: [r.nama],
        terapkan: (i) => ({ ...i, riSlipRnm: r.id }),
      }))
    case 'billing':
      return (await cariCedingPolis(teks)).map((r) => ({
        kunci: r.id,
        sel: [r.id, r.nama],
        // Retrocessionaire dikosongkan; server mengisinya untuk Billing tertentu
        // (setSecurityReinsurer_act).
        terapkan: (i) => ({ ...i, retroId: r.id, retroName: r.nama, securityReinsurerId: '', securityReinsurer: '' }),
      }))
    case 'retro':
      return (await cariCedingPolis(teks)).map((r) => ({
        kunci: r.id,
        sel: [r.id, r.nama],
        terapkan: (i) => ({ ...i, securityReinsurerId: r.id, securityReinsurer: r.nama }),
      }))
  }
}

function PopupCari({
  judul,
  kolom,
  cari,
  onPilih,
  onTutup,
}: {
  judul: string
  kolom: string[]
  cari: (teks: string) => Promise<BarisPopup[]>
  onPilih: (b: BarisPopup) => void
  onTutup: () => void
}) {
  const [teks, setTeks] = useState('')
  const [hasil, setHasil] = useState<BarisPopup[] | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)

  async function jalankan(): Promise<void> {
    if (sibuk) return
    setSibuk(true)
    setGalat(null)
    try {
      setHasil(await cari(teks))
    } catch (e) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  return (
    <Modal
      judul={judul}
      onTutup={onTutup}
      onKirim={() => { void jalankan() }}
      aksi={
        <button type="submit" className="btn btn--primary" disabled={sibuk}>
          Search
        </button>
      }
      lebar
    >
      <Field label="Search" value={teks} onChange={setTeks} autoFocus />
      {galat !== null && <Gagal galat={galat} />}
      {hasil !== null && hasil.length === 0 && <p>No matching results.</p>}
      {hasil !== null && hasil.length > 0 && (
        <table className="inbox__tabel">
          <thead>
            <tr>
              {kolom.map((k) => (
                <th key={k}>{k}</th>
              ))}
              <th />
            </tr>
          </thead>
          <tbody>
            {hasil.map((b) => (
              <tr key={b.kunci}>
                {b.sel.map((s, i) => (
                  <td key={i}>{s}</td>
                ))}
                <td>
                  <button type="button" className="btn btn--ghost btn--sm" onClick={() => { onPilih(b) }}>
                    Choose
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </Modal>
  )
}
