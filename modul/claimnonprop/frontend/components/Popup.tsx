// Pop-up harness Claim Non Prop (showHarness Target=popup / local action): ChooseMasterTNonProp, ViewListPolicyCNP,
// CauseofLoss_Harness, CatastrofeList, Hitung_Test, ViewOldAllocation, ViewAttachmentNP, ViewHistoryMasterID_NP. Judul =
// pyLabel / pyWindowName rule; judul kolom VERBATIM section XML; baris dibaca dari server (`GET .../pilihan/{jenis}`) dan
// pilihan dikirim balik sebagai aksi (server membaca ulang barisnya). Pola `modul/claimprop/frontend/components/Popup.tsx`
// (disalin, bukan impor).

import { useEffect, useState, type ReactNode } from 'react'

import { Gagal, Memuat, Modal } from '../../../../inti/frontend/components/ui/dasar'
import { pilihanKasus, type Baris, type Pilihan } from '../api'
import { tampilTanggal } from '../ketikTanggal'
import { CNP } from '../labels'
import { tampilAngka } from '../nilai'
import { kunciMaster, type KolomMaster } from './kunciBaris'
import { jumlahHalaman, potongHalaman } from './susun'

export type JenisPopup =
  | 'master'
  | 'daftarPolis'
  | 'sebab'
  | 'katastrofe'
  | 'selisihAktual'
  | 'alokasiLama'
  | 'lampiranBayar'
  | 'riwayatMaster'

/** Sumber popup master (`BrowseDtlMasterTNP_Act` Param.TreatyIN): IN = Outstanding / Input Acceptation, INEDM. */
export type SumberMaster = 'IN' | 'INEDM'

/** Kolom popup master: judul VERBATIM section ChooseMasterTNonProp (grid S2, `Test != 'OUT'`); `kunci` = filter. */
const KOLOM_MASTER: { label: string; kunci?: string; nilai: (b: KolomMaster) => string }[] = [
  { label: 'Treaty ID', kunci: 'treatyId', nilai: (b) => b.treatyId },
  { label: 'Contract Name', kunci: 'contractName', nilai: (b) => b.treatyContractName },
  { label: 'Source Of Business', kunci: 'sob', nilai: (b) => b.sob },
  { label: 'Insured Name', kunci: 'ceding', nilai: (b) => b.ceding },
  { label: 'Proportion Type', nilai: (b) => b.proportionType },
  { label: 'Treaty Group', kunci: 'treatyGroup', nilai: (b) => b.treatyGroup },
  { label: 'Class of Business', kunci: 'classOfBusiness', nilai: (b) => b.classOfBusiness },
  { label: 'Treaty Year', kunci: 'treatyYear', nilai: (b) => b.treatyYear },
]
const PER_HALAMAN_MASTER = 50
const BATAS_MASTER = 500 // models.BatasMaster

/** `models.BarisPolis` - Polis.pxResults (GetDataPolisNonProp_SQL). */
interface BarisPolis {
  policyNo: string
  sob: string
  cedingCo: string
  beginDate: string
  endDate: string
}

interface BarisSebab {
  id: string
  description: string
}

interface BarisKatastrofe {
  id: string
  stsKatastrofe: string
  nonKatastrofeType: string
  note: string
}

/** `models.SelisihAktual` - GetSelisihActual_Act. */
interface SelisihAktual {
  estimasiTotal: string
  estimasiSpread: string
  akseptasi: Baris[] | null
  selisihTotal: string
  selisihSpread: string
  aktualTotal: string
  aktualSpread: string
}

/** `repository.BarisLampiranBayar` - ViewAttachmentNP. */
interface BarisLampiran {
  kategori: string
  namaFile: string
  noAksep: string
  noPrekas: string
  tanggal: string
}

/** `services.RiwayatMaster` - ViewHistoryMasterID_NP. */
interface RiwayatMaster {
  judul: string
  baris: {
    caseId: string
    xol: string
    currency: string
    grossAdjustment: string
    cnpReinstatement: string
  }[]
  total: Baris[] | null
}

const JUDUL: Record<JenisPopup, string> = {
  master: CNP.popMaster,
  daftarPolis: CNP.popPolis,
  sebab: CNP.popSebab,
  katastrofe: CNP.popKatastrofe,
  selisihAktual: CNP.popSelisih,
  alokasiLama: CNP.popAlokasiLama,
  lampiranBayar: CNP.popLampiranBayar,
  riwayatMaster: CNP.popRiwayatMaster,
}

/** Sel angka rata kanan. */
function Angka({ v }: { v: string | undefined }) {
  return <td className="claimnonprop__angka">{tampilAngka(v ?? '')}</td>
}

/** Satu baris "label : nilai" hanya-baca (section SIMPLELAYOUT harness). */
function BarisNilai({ label, v }: { label: string; v: string }) {
  return (
    <div className="claimnonprop__baris">
      <span className="claimnonprop__label-medan">{label}</span>
      <div className="claimnonprop__isi-medan">
        <span className="claimnonprop__nilai claimnonprop__angka">{tampilAngka(v)}</span>
      </div>
    </div>
  )
}

/** Tabel hanya-baca: judul kolom VERBATIM + baris. */
function Tabel({ kolom, baris }: { kolom: readonly string[]; baris: ReactNode[] }) {
  return (
    <table className="claimnonprop__tabel">
      <thead>
        <tr>
          {kolom.map((k, i) => (
            <th key={i}>{k}</th>
          ))}
        </tr>
      </thead>
      <tbody>
        {baris.length === 0 ? (
          <tr>
            <td className="muted" colSpan={kolom.length}>
              —
            </td>
          </tr>
        ) : (
          baris
        )}
      </tbody>
    </table>
  )
}

export default function Popup({
  id,
  jenis,
  sumber = 'IN',
  pelaku,
  sts,
  mataUang,
  onPilih,
  onTutup,
}: {
  id: string
  jenis: JenisPopup
  /** Popup master: Choose Master In (IN) / Choose Master Input Acceptation sesudah akseptasi (INEDM). */
  sumber?: SumberMaster
  pelaku: string
  /** Catastrophe / NonKatastrofeType klaim (form katastrofe baru, read-only). */
  sts: { sts: string; non: string }
  /** BrowseCurrency_RD - label mata uang (kolom Currency View Old Allocation). */
  mataUang: readonly Pilihan[]
  onPilih: (aksi: string, param: string) => void
  onTutup: () => void
}) {
  const [cari, setCari] = useState('')
  const [data, setData] = useState<unknown>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [formBaru, setFormBaru] = useState(false)
  const [catatan, setCatatan] = useState('')
  const [saring, setSaring] = useState<Record<string, string>>({})
  const [hal, setHal] = useState(1)

  useEffect(() => {
    // Jawaban basi diabaikan: permintaan saringan lama tidak boleh menimpa hasil saringan terbaru.
    let aktif = true
    const t = setTimeout(() => {
      const s = jenis === 'master' ? { ...saring, sumber } : saring
      pilihanKasus<unknown>(id, jenis, 0, cari, s).then(
        (d) => {
          if (!aktif) return
          setData(d)
          setGalat(null)
        },
        (g: unknown) => {
          if (aktif) setGalat(g)
        },
      )
    }, 300)
    return () => {
      aktif = false
      clearTimeout(t)
    }
  }, [id, jenis, cari, saring, sumber])

  const labelMataUang = (v: string) => mataUang.find((m) => m.nilai === v)?.label ?? v
  const pakaiCari = jenis === 'sebab' || jenis === 'katastrofe'
  let isi
  if (galat) isi = <Gagal galat={galat} />
  else if (data === null) isi = <Memuat pesan={CNP.memuat} />
  else if (jenis === 'master') {
    const semua = data as KolomMaster[]
    const nHal = jumlahHalaman(semua.length, PER_HALAMAN_MASTER)
    const halIni = Math.min(hal, nHal)
    const awal = (halIni - 1) * PER_HALAMAN_MASTER
    isi = (
      <>
        <p className="claimnonprop__catatan">{CNP.catatanMaster}</p>
        <div className="claimnonprop__pager">
          <span>
            {CNP.menampilkan} {semua.length === 0 ? 0 : awal + 1}–{Math.min(awal + PER_HALAMAN_MASTER, semua.length)}{' '}
            {CNP.dari} {semua.length}
            {semua.length >= BATAS_MASTER && ` (${CNP.batasMaster})`}
          </span>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={halIni <= 1}
            onClick={() => setHal(halIni - 1)}
            aria-label="Previous"
          >
            ‹
          </button>
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            disabled={halIni >= nHal}
            onClick={() => setHal(halIni + 1)}
            aria-label="Next"
          >
            ›
          </button>
        </div>
        <table className="claimnonprop__tabel">
          <thead>
            <tr>
              <th />
              {KOLOM_MASTER.map((k) => (
                <th key={k.label}>{k.label}</th>
              ))}
            </tr>
            <tr>
              <th />
              {KOLOM_MASTER.map((k) => (
                <th key={k.label}>
                  {k.kunci && (
                    <input
                      className="field__input claimnonprop__input--sel"
                      placeholder={CNP.saring}
                      aria-label={`${CNP.saring} ${k.label}`}
                      value={saring[k.kunci] ?? ''}
                      onChange={(e) => {
                        const kunci = k.kunci ?? ''
                        setSaring((s) => ({ ...s, [kunci]: e.target.value }))
                        setHal(1)
                      }}
                    />
                  )}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {semua.length === 0 && (
              <tr>
                <td className="muted" colSpan={KOLOM_MASTER.length + 1}>
                  —
                </td>
              </tr>
            )}
            {potongHalaman(semua, PER_HALAMAN_MASTER, halIni).map((b) => (
              <tr key={kunciMaster(b)}>
                <td>
                  <button
                    type="button"
                    className="btn btn--sm btn--primary"
                    onClick={() =>
                      onPilih('SetValueClaimTNP', `${b.treatyId}|${b.classOfBusinessId}|${b.treatyGroup}|${sumber}`)
                    }
                  >
                    {CNP.pilih}
                  </button>
                </td>
                {KOLOM_MASTER.map((k) => (
                  <td key={k.label}>{k.nilai(b)}</td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </>
    )
  } else if (jenis === 'daftarPolis') {
    // ViewListPolicyCNP: tautan Policy No = CheckNoPolicy (NOPOLIS) lalu closeContainer.
    isi = (
      <Tabel
        kolom={['Policy No', 'Source of Business', 'Ceding Co', 'Begin Date', 'End Date']}
        baris={(data as BarisPolis[]).map((b, i) => (
          <tr key={`${b.policyNo}|${i}`}>
            <td>
              <button
                type="button"
                className="claimnonprop__tautan"
                onClick={() => onPilih('CheckNoPolicy', b.policyNo)}
              >
                {b.policyNo}
              </button>
            </td>
            <td>{b.sob}</td>
            <td>{b.cedingCo}</td>
            <td>{tampilTanggal(b.beginDate, 'tanggal')}</td>
            <td>{tampilTanggal(b.endDate, 'tanggal')}</td>
          </tr>
        ))}
      />
    )
  } else if (jenis === 'sebab') {
    isi = (
      <>
        <div className="claimnonprop__label">{CNP.judulSebab}</div>
        <table className="claimnonprop__tabel">
          <thead>
            <tr>
              <th>Cause of Loss</th>
              <th />
            </tr>
          </thead>
          <tbody>
            {(data as BarisSebab[]).map((b) => (
              <tr key={b.id}>
                <td>{b.description}</td>
                <td>
                  <button
                    type="button"
                    className="btn btn--sm btn--primary"
                    onClick={() => onPilih('GetNameCauseofLoss', b.id)}
                  >
                    {CNP.pilih}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </>
    )
  } else if (jenis === 'katastrofe') {
    isi = formBaru ? (
      <div className="form-grid">
        <label className="field">
          <span className="field__label">Catastrophe</span>
          <input className="field__input field__input--readonly" readOnly tabIndex={-1} value={sts.sts} />
        </label>
        {sts.sts === 'Non-Catastrophe' && (
          <label className="field">
            <span className="field__label">&nbsp;</span>
            <input className="field__input field__input--readonly" readOnly tabIndex={-1} value={sts.non} />
          </label>
        )}
        <label className="field field--lebar">
          <span className="field__label">{CNP.note}</span>
          <textarea
            className="field__input claimnonprop__area"
            value={catatan}
            onChange={(e) => setCatatan(e.target.value)}
          />
        </label>
        <label className="field">
          <span className="field__label">{CNP.userInput}</span>
          <input className="field__input field__input--readonly" readOnly tabIndex={-1} value={pelaku} />
        </label>
        <div className="field--lebar claimnonprop__tombol claimnonprop__tombol--akhir">
          <button type="button" className="btn btn--ghost btn--sm" onClick={() => setFormBaru(false)}>
            {CNP.cancel}
          </button>
          <button type="button" className="btn btn--sm btn--primary" onClick={() => onPilih('SaveCatasrtope', catatan)}>
            {CNP.save}
          </button>
        </div>
      </div>
    ) : (
      <>
        <div className="claimnonprop__grid-alat">
          <button type="button" className="btn btn--sm" onClick={() => setFormBaru(true)}>
            {CNP.addNew}
          </button>
        </div>
        <table className="claimnonprop__tabel">
          <thead>
            <tr>
              <th />
              <th />
              <th>{CNP.note}</th>
              <th />
            </tr>
          </thead>
          <tbody>
            {(data as BarisKatastrofe[]).map((b) => (
              <tr key={b.id}>
                <td>{b.stsKatastrofe}</td>
                <td>{b.nonKatastrofeType}</td>
                <td>{b.note}</td>
                <td>
                  <button
                    type="button"
                    className="btn btn--sm btn--primary"
                    onClick={() => onPilih('SetCatastrope', b.id)}
                  >
                    {CNP.pilih}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </>
    )
  } else if (jenis === 'selisihAktual') {
    const r = data as SelisihAktual
    isi = (
      <div className="claimnonprop__tumpuk">
        <section className="claimnonprop__sub">
          <h4 className="claimnonprop__subjudul">{CNP.nilaiEstimasi}</h4>
          <BarisNilai label={CNP.totalClaim} v={r.estimasiTotal} />
          <BarisNilai label={CNP.claimSpreadedKecil} v={r.estimasiSpread} />
        </section>
        <section className="claimnonprop__sub">
          <h4 className="claimnonprop__subjudul">{CNP.nilaiAkseptasi}</h4>
          <Tabel
            kolom={[CNP.totalClaim, CNP.claimSpreadedKecil]}
            baris={(r.akseptasi ?? []).map((b, i) => (
              <tr key={i}>
                <Angka v={b.CARI2} />
                <Angka v={b.CARI1} />
              </tr>
            ))}
          />
        </section>
        <section className="claimnonprop__sub">
          <h4 className="claimnonprop__subjudul">{CNP.nilaiSelisih}</h4>
          <BarisNilai label={CNP.totalClaim} v={r.selisihTotal} />
          <BarisNilai label={CNP.claimSpreaded} v={r.selisihSpread} />
        </section>
        <section className="claimnonprop__sub">
          <h4 className="claimnonprop__subjudul">{CNP.cadangan100}</h4>
          <BarisNilai label={CNP.totalClaim} v={r.aktualTotal} />
          <BarisNilai label={CNP.claimSpreaded} v={r.aktualSpread} />
        </section>
      </div>
    )
  } else if (jenis === 'alokasiLama') {
    isi = (
      <section className="claimnonprop__sub">
        <h4 className="claimnonprop__subjudul">{CNP.alokasiLama}</h4>
        <Tabel
          kolom={[
            'Currency',
            'Treaty Name',
            'Claim Amount',
            'RNM Share (%)',
            'Claim Amount RNM',
            'Adjuster Fee',
            'Salvage',
          ]}
          baris={(data as Baris[]).map((b, i) => (
            <tr key={i}>
              <td>{labelMataUang(b.CurrencyID ?? '')}</td>
              <td>{b.TreatyName}</td>
              <Angka v={b.ClaimAmountAdjust} />
              <Angka v={b.ClaimPercentage} />
              <Angka v={b.ClaimSpreaded} />
              <Angka v={b.AdjusterFee} />
              <Angka v={b.Salvage} />
            </tr>
          ))}
        />
      </section>
    )
  } else if (jenis === 'lampiranBayar') {
    isi = (
      <Tabel
        kolom={['Kategori', 'Nama File', 'No Akseptasi', 'No Prekas', 'Tanggal/Waktu']}
        baris={(data as BarisLampiran[]).map((b, i) => (
          <tr key={i}>
            <td>{b.kategori}</td>
            <td>{b.namaFile}</td>
            <td>{b.noAksep}</td>
            <td>{b.noPrekas}</td>
            <td>{tampilTanggal(b.tanggal, 'tanggal-waktu')}</td>
          </tr>
        ))}
      />
    )
  } else {
    // ViewHistoryMasterID_NP: label InputSpreading.CARI40, grid TempMaster lalu TempTotal. Judul kolom keempat grid
    // pertama VERBATIM XML ("Loss to Layer (100%)" dua kali; isinya CNPReinstatement - kelainan XML, PARITAS).
    const r = data as RiwayatMaster
    isi = (
      <div className="claimnonprop__tumpuk">
        {r.judul !== '' && <div className="claimnonprop__label">{r.judul}</div>}
        <Tabel
          kolom={['CLAIM ID', 'Currency', 'Loss to Layer (100%)', 'Loss to Layer (100%)', 'Layer']}
          baris={r.baris.map((b, i) => (
            <tr key={i}>
              <td>{b.caseId}</td>
              <td>{b.currency}</td>
              <Angka v={b.grossAdjustment} />
              <Angka v={b.cnpReinstatement} />
              <td>{b.xol}</td>
            </tr>
          ))}
        />
        <Tabel
          kolom={['Layer', 'Currency', 'Total Loss to Layer (100%)', 'Total Reinstatement Premium (100%)']}
          baris={(r.total ?? []).map((b, i) => (
            <tr key={i}>
              <td>{b.XOL}</td>
              <td>{b.Currency}</td>
              <Angka v={b.GrossAdjustment} />
              <Angka v={b.CNPReinstatement} />
            </tr>
          ))}
        />
      </div>
    )
  }

  return (
    <Modal
      judul={JUDUL[jenis]}
      onTutup={onTutup}
      labelBatal={CNP.tutup}
      penuh={jenis === 'master'}
      lebar={jenis !== 'master'}
    >
      {pakaiCari && !formBaru && (
        <input
          className="field__input claimnonprop__cari"
          placeholder={CNP.cariPopup}
          value={cari}
          onChange={(e) => setCari(e.target.value)}
        />
      )}
      {isi}
    </Modal>
  )
}
