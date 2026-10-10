// Form Bordereaux - padanan harness/section `InputBordereaux`: header (Type, Type Business, Choose Master Treaty,
// periode, Reff No), tab Details (Template, Upload CSV, grid detail), Summary (total per mata uang), Submit
// (Maker > Checker > Supervisor), lampiran (`AttachmentsBdx`), History, Save, Close. Tanpa JSON: Save menulis header dan
// detail ke tabelnya.

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Memuat } from '../../../../inti/frontend/components/ui/dasar'
import { kodeStatusGalat, pesanGalat } from '../../../../inti/frontend/klien'
import { unduhTemplat } from '../../../../inti/frontend/templat/api'
import { buka, simpan, submit, unggahCsv, type Baris, type Hak, type MasterTreaty, type Pilihan, type Ringkasan, type Riwayat } from '../api'
import { barisPesan, businessUntuk, formatAngka, kunciKombinasi } from '../aturan'
import { BDX, teksStatus } from '../labels'
import DialogMasterTreaty from './DialogMasterTreaty'
import GridDetail from './GridDetail'
import KepalaForm, { ISIAN_KOSONG, type Isian } from './KepalaForm'
import PanelLampiran from './PanelLampiran'

type Tab = 'details' | 'summary' | 'submit'

const HAK_KOSONG: Hak = { ubah: false, hapus: false, submit: false, putuskan: false, lampiran: false }

export default function FormBordereaux({
  id,
  lihat,
  pilihan,
  onTutup,
}: {
  /** null = berkas baru. */
  id: string | null
  /** View (`ViewData=1`): hanya dibaca walau berhak mengubah. */
  lihat: boolean
  pilihan: Pilihan
  onTutup: (berubah: boolean) => void
}) {
  const [isian, setIsian] = useState<Isian>(ISIAN_KOSONG)
  const [baris, setBaris] = useState<Baris[]>([])
  const [ringkasan, setRingkasan] = useState<Ringkasan[]>([])
  const [riwayat, setRiwayat] = useState<Riwayat[]>([])
  const [hak, setHak] = useState<Hak>(HAK_KOSONG)
  const [status, setStatus] = useState({ position: '', statusAksep: '' })
  const [tab, setTab] = useState<Tab>('details')
  const [memuat, setMemuat] = useState(id !== null)
  const [sibuk, setSibuk] = useState<'csv' | 'simpan' | 'submit' | null>(null)
  const [galat, setGalat] = useState<unknown>(null)
  const [tolak, setTolak] = useState<string[]>([])
  const [pesan, setPesan] = useState<string | null>(null)
  const [dialog, setDialog] = useState(false)
  const [setuju, setSetuju] = useState(true)
  const [komentar, setKomentar] = useState('')
  const [berubah, setBerubah] = useState(false)

  const muat = useCallback((bdxId: string) => {
    setMemuat(true)
    buka(bdxId).then(
      (r) => {
        const h = r.header
        setIsian({
          bdxId: h.bdxId,
          type: h.type,
          business: h.typeBusiness,
          masterId: h.masterId,
          cedingName: h.cedingName,
          sobName: h.sobName,
          treatyName: h.treatyName,
          reportStart: h.reportStart,
          reportEnd: h.reportEnd,
          reffNoSoa: h.reffNoSoa,
          reffNoBdx: h.reffNoBdx,
        })
        setBaris(r.baris)
        setRingkasan(r.ringkasan)
        setRiwayat(r.riwayat)
        setHak(r.hak)
        setStatus({ position: h.position, statusAksep: h.statusAksep })
        setMemuat(false)
      },
      (g: unknown) => {
        setGalat(g)
        setMemuat(false)
      },
    )
  }, [])

  useEffect(() => {
    if (id !== null) muat(id)
  }, [id, muat])

  const baru = isian.bdxId === ''
  const bolehUbah = !lihat && (baru ? pilihan.bolehBuat : hak.ubah)
  const kunci = kunciKombinasi(isian.type, isian.business)
  const kolom = pilihan.kolom[kunci] ?? []
  const slot = pilihan.templat[kunci]
  const selesai = status.statusAksep === 'Resolve-Complete'

  const gagal = (g: unknown) => {
    const teks = kodeStatusGalat(g) === 422 ? pesanGalat(g) : undefined
    if (teks !== undefined) setTolak(barisPesan(teks))
    else setGalat(g)
  }
  const bersih = () => {
    setGalat(null)
    setTolak([])
    setPesan(null)
  }

  const ubah = (sebagian: Partial<Isian>) => setIsian((x) => ({ ...x, ...sebagian }))

  const gantiType = (type: string) => {
    ubah({ type, business: businessUntuk(type, type === isian.type ? isian.business : '') })
    setBaris([])
    setRingkasan([])
  }

  const unggah = (f: File) => {
    bersih()
    setSibuk('csv')
    f.text()
      .then((teks) => unggahCsv(isian.type, isian.business, teks))
      .then(
        (p) => {
          setSibuk(null)
          setBaris(p.baris)
          setRingkasan(p.ringkasan)
        },
        (g: unknown) => {
          setSibuk(null)
          gagal(g)
        },
      )
  }

  const permintaan = () => ({
    bdxId: isian.bdxId,
    type: isian.type,
    business: isian.business,
    masterId: isian.masterId,
    reportStart: isian.reportStart,
    reportEnd: isian.reportEnd,
    reffNoSoa: isian.reffNoSoa,
    reffNoBdx: isian.reffNoBdx,
    baris,
  })

  const jalankanSimpan = () => {
    if (sibuk !== null) return
    bersih()
    setSibuk('simpan')
    simpan(permintaan()).then(
      (r) => {
        setSibuk(null)
        setBerubah(true)
        setPesan(BDX.tersimpan(r.bdxId))
        muat(r.bdxId)
      },
      (g: unknown) => {
        setSibuk(null)
        gagal(g)
      },
    )
  }

  // Submit pembuat = Save lalu Submit (`ActionSubmit` langkah 3 memanggil `BdxSave_Act`); Checker/Supervisor = putusan.
  const jalankanSubmit = () => {
    if (sibuk !== null || isian.bdxId === '') return
    bersih()
    setSibuk('submit')
    const simpanDulu = hak.submit && bolehUbah ? simpan(permintaan()).then(() => undefined) : Promise.resolve()
    simpanDulu
      .then(() => submit(isian.bdxId, hak.submit ? true : setuju, komentar))
      .then(
        () => {
          setSibuk(null)
          setBerubah(true)
          setKomentar('')
          setPesan(BDX.terkirim)
          muat(isian.bdxId)
        },
        (g: unknown) => {
          setSibuk(null)
          gagal(g)
        },
      )
  }

  if (memuat) return <Memuat pesan={BDX.memuatBerkas} />

  return (
    <>
      <header className="inbox__kepala">
        <h2 className="inbox__judul">{BDX.judul}</h2>
        {!baru && (
          <>
            <span className="bordereaux__id" title={BDX.id}>
              {isian.bdxId}
            </span>
            <span className="badge bordereaux__badge">{teksStatus(status.statusAksep)}</span>
          </>
        )}
      </header>

      <KepalaForm
        isian={isian}
        pilihan={pilihan}
        bolehUbah={bolehUbah}
        onGantiType={gantiType}
        onGantiBusiness={(b) => {
          ubah({ business: b })
          setBaris([])
          setRingkasan([])
        }}
        onPilihMaster={() => setDialog(true)}
        onUbah={ubah}
      />

      {pesan !== null && (
        <div className="alert alert--ok" role="status">
          {pesan}
        </div>
      )}
      {tolak.length > 0 && (
        <div className="alert alert--error" role="alert">
          <ul className="bordereaux__tolak">
            {tolak.map((t) => (
              <li key={t}>{t}</li>
            ))}
          </ul>
        </div>
      )}
      {galat !== null && <Gagal galat={galat} />}

      <nav className="bordereaux__tab" role="tablist">
        {(['details', 'summary', 'submit'] as Tab[])
          .filter((t) => t !== 'submit' || !selesai)
          .map((t) => (
            <button
              key={t}
              type="button"
              role="tab"
              aria-selected={tab === t}
              className={tab === t ? 'bordereaux__tab-butir bordereaux__tab-butir--aktif' : 'bordereaux__tab-butir'}
              onClick={() => setTab(t)}
            >
              {{ details: BDX.tabDetails, summary: BDX.tabSummary, submit: BDX.tabSubmit }[t]}
            </button>
          ))}
      </nav>

      <section className="bordereaux__kartu">
        {tab === 'details' && (
          <>
            {isian.business !== '' && bolehUbah && (
              <div className="bordereaux__aksi-unggah">
                <span className="muted">{BDX.unggahPetunjuk}</span>
                <span className="toolbar__spacer" />
                {slot !== undefined && (
                  <button
                    type="button"
                    className="btn btn--ghost"
                    onClick={() => {
                      unduhTemplat(slot.kode, slot.nama).catch((g: unknown) => setGalat(g))
                    }}
                  >
                    {BDX.template}
                  </button>
                )}
                <label className="btn btn--ghost bordereaux__unggah">
                  {sibuk === 'csv' ? BDX.membaca : BDX.uploadCsv}
                  <input
                    type="file"
                    accept=".csv,text/csv"
                    disabled={sibuk !== null}
                    onChange={(e) => {
                      const f = e.target.files?.[0]
                      e.target.value = ''
                      if (f !== undefined) unggah(f)
                    }}
                  />
                </label>
              </div>
            )}
            {isian.business === '' && <p className="muted">{BDX.pilihTypeDulu}</p>}
            {isian.business !== '' && baris.length === 0 && <p className="muted">{BDX.detailKosong}</p>}
            {isian.business !== '' && baris.length > 0 && <GridDetail kolom={kolom} baris={baris} />}
          </>
        )}
        {tab === 'summary' &&
          (ringkasan.length === 0 ? (
            <p className="muted">{BDX.ringkasanKosong}</p>
          ) : (
            <table className="inbox__tabel bordereaux__grid">
              <thead>
                <tr>
                  <th>{BDX.currency}</th>
                  <th className="bordereaux__angka">{BDX.ringkasan[isian.type]?.[0] ?? ''}</th>
                  <th className="bordereaux__angka">{BDX.ringkasan[isian.type]?.[1] ?? ''}</th>
                </tr>
              </thead>
              <tbody>
                {ringkasan.map((r) => (
                  <tr key={r.currency} className="inbox__baris">
                    <td>{r.currency === '' ? '—' : r.currency}</td>
                    <td className="bordereaux__angka">{formatAngka(r.reinsurer)}</td>
                    <td className="bordereaux__angka">{formatAngka(r.rnm)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          ))}
        {tab === 'submit' && (
          <div className="bordereaux__submit">
            {baru && <p className="muted">{BDX.submitPembuat}</p>}
            {!baru && !hak.submit && !hak.putuskan && <p className="muted">{BDX.submitTidakMenunggu}</p>}
            {!baru && hak.submit && <p>{BDX.submitPembuat}</p>}
            {!baru && hak.putuskan && (
              <fieldset className="bordereaux__radio">
                <legend className="field__label">{BDX.statusPersetujuan}</legend>
                <label>
                  <input type="radio" name="bdx-setuju" checked={setuju} onChange={() => setSetuju(true)} /> {BDX.approve}
                </label>
                <label>
                  <input type="radio" name="bdx-setuju" checked={!setuju} onChange={() => setSetuju(false)} /> {BDX.reject}
                </label>
              </fieldset>
            )}
            {!baru && (hak.submit || hak.putuskan) && (
              <>
                <label className="field">
                  <span className="field__label">{BDX.comment}</span>
                  <textarea className="field__input" rows={3} maxLength={2000} value={komentar} onChange={(e) => setKomentar(e.target.value)} />
                </label>
                <div>
                  <button type="button" className="btn btn--primary" disabled={sibuk !== null} onClick={jalankanSubmit}>
                    {sibuk === 'submit' ? BDX.mengirim : BDX.submit}
                  </button>
                </div>
              </>
            )}
          </div>
        )}
      </section>

      <div className="bordereaux__kaki">
        <span className="toolbar__spacer" />
        <button type="button" className="btn btn--ghost" disabled={sibuk !== null} onClick={() => onTutup(berubah)}>
          {BDX.close}
        </button>
        {bolehUbah && (
          <button type="button" className="btn btn--primary" disabled={sibuk !== null} onClick={jalankanSimpan}>
            {sibuk === 'simpan' ? BDX.menyimpan : BDX.save}
          </button>
        )}
      </div>

      {/* `AttachmentsBdx` (b27972, sesudah Close/Save): Upload File / Delete hanya saat Edit - mode View tidak, siapa pun
          (keputusan work owner 08-10-2026; pengecualian IT Developer Pega dibuang). */}
      {!baru && <PanelLampiran bdxId={isian.bdxId} boleh={hak.lampiran && !lihat} />}

      {!baru && (
        <section className="bordereaux__kartu">
          <h3 className="bordereaux__judul-kartu">{BDX.history}</h3>
          {riwayat.length === 0 ? (
            <p className="muted">{BDX.historyKosong}</p>
          ) : (
            <table className="inbox__tabel bordereaux__grid">
              <thead>
                <tr>
                  <th>{BDX.date}</th>
                  <th>{BDX.pic}</th>
                  <th>{BDX.approval}</th>
                  <th>{BDX.comment}</th>
                </tr>
              </thead>
              <tbody>
                {riwayat.map((r, i) => (
                  <tr key={i} className="inbox__baris">
                    <td>{r.tanggal}</td>
                    <td>{r.pic}</td>
                    <td>{r.isApproved ? BDX.approve : BDX.reject}</td>
                    <td>{r.komentar === '' ? '—' : r.komentar}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </section>
      )}

      {dialog && (
        <DialogMasterTreaty
          onTutup={() => setDialog(false)}
          onPilih={(t: MasterTreaty) => {
            ubah({ masterId: t.id, cedingName: t.cedingName, sobName: t.sobName, treatyName: t.contractName })
            setDialog(false)
          }}
        />
      )}
    </>
  )
}
