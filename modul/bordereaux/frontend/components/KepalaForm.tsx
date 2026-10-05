// Kepala form Bordereaux (`InputBordereaux` bagian atas): Type, Type Business, Master Treaty, periode, Reff No.
// Komponen tampilan saja - keadaan dipegang FormBordereaux. Ditata ulang atas permintaan work owner 04-10-2026
// ("RAPIHKAN ... yang enak dilihat dan user friendly", "kalo view buat readonly aja ga usah disable"):
//   - ubah: pilihan berupa chip berkisi, Master Treaty satu panel, isian satu baris yang membungkus;
//   - View / tanpa hak ubah: NILAI saja (teks), tanpa kotak isian maupun radio yang dinonaktifkan.

import type { Pilihan } from '../api'
import { keIso, keKabel } from '../aturan'
import { BDX } from '../labels'

/** Isian kepala form. */
export interface Isian {
  bdxId: string
  type: string
  business: string
  masterId: string
  cedingName: string
  sobName: string
  treatyName: string
  /** DD-MM-YYYY. */
  reportStart: string
  reportEnd: string
  reffNoSoa: string
  reffNoBdx: string
}

export const ISIAN_KOSONG: Isian = {
  bdxId: '',
  type: '',
  business: '',
  masterId: '',
  cedingName: '',
  sobName: '',
  treatyName: '',
  reportStart: '',
  reportEnd: '',
  reffNoSoa: '',
  reffNoBdx: '',
}

const strip = (s: string) => (s === '' ? '—' : s)

/** `DD-MM-YYYY` -> `DD/MM/YYYY` untuk dibaca. */
function tanggalBaca(s: string): string {
  return s === '' ? '—' : s.replace(/-/g, '/')
}

function Chip({ nama, nilai, dipilih, onPilih }: { nama: string; nilai: string; dipilih: boolean; onPilih: () => void }) {
  return (
    <label className={dipilih ? 'bordereaux__chip bordereaux__chip--aktif' : 'bordereaux__chip'}>
      <input type="radio" name={nama} checked={dipilih} onChange={onPilih} />
      <span>{nilai}</span>
    </label>
  )
}

function Wajib() {
  return (
    <span className="bordereaux__wajib" aria-hidden="true">
      *
    </span>
  )
}

/** Panel Master Treaty - sama untuk ubah dan View. */
function PanelTreaty({ isian, tombol }: { isian: Isian; tombol: React.ReactNode }) {
  const kosong = isian.masterId === ''
  return (
    <div className="bordereaux__treaty">
      <div className="bordereaux__treaty-kepala">
        <span className="bordereaux__judul-bagian">
          {BDX.masterTreaty}
          {tombol !== null && <Wajib />}
        </span>
        {tombol}
      </div>
      {kosong ? (
        <p className="muted bordereaux__treaty-kosong">{BDX.masterKosong}</p>
      ) : (
        <dl className="bordereaux__treaty-isi">
          <div>
            <dt>{BDX.masterId}</dt>
            <dd>{isian.masterId}</dd>
          </div>
          <div>
            <dt>{BDX.cedingCo}</dt>
            <dd>{strip(isian.cedingName)}</dd>
          </div>
          <div>
            <dt>{BDX.sobName}</dt>
            <dd>{strip(isian.sobName)}</dd>
          </div>
          <div>
            <dt>{BDX.treatyName}</dt>
            <dd>{strip(isian.treatyName)}</dd>
          </div>
        </dl>
      )}
    </div>
  )
}

/** View: seluruh isian sebagai nilai, tanpa kontrol. */
function KepalaBaca({ isian }: { isian: Isian }) {
  return (
    <section className="bordereaux__kartu bordereaux__kepala-form">
      <dl className="bordereaux__ringkas">
        <div>
          <dt>{BDX.type}</dt>
          <dd>
            <span className="bordereaux__nilai-chip">{strip(isian.type)}</span>
          </dd>
        </div>
        <div>
          <dt>{BDX.typeBusiness}</dt>
          <dd>
            <span className="bordereaux__nilai-chip">{strip(isian.business)}</span>
          </dd>
        </div>
        <div>
          <dt>{BDX.bdxStartEnd}</dt>
          <dd>
            {tanggalBaca(isian.reportStart)} – {tanggalBaca(isian.reportEnd)}
          </dd>
        </div>
        <div>
          <dt>{BDX.reffNoOfSoa}</dt>
          <dd>{strip(isian.reffNoSoa)}</dd>
        </div>
        <div>
          <dt>{BDX.reffNoOfBdx}</dt>
          <dd>{strip(isian.reffNoBdx)}</dd>
        </div>
      </dl>
      <PanelTreaty isian={isian} tombol={null} />
    </section>
  )
}

export default function KepalaForm({
  isian,
  pilihan,
  bolehUbah,
  onGantiType,
  onGantiBusiness,
  onPilihMaster,
  onUbah,
}: {
  isian: Isian
  pilihan: Pilihan
  bolehUbah: boolean
  onGantiType: (type: string) => void
  onGantiBusiness: (business: string) => void
  onPilihMaster: () => void
  onUbah: (sebagian: Partial<Isian>) => void
}) {
  if (!bolehUbah) return <KepalaBaca isian={isian} />
  const businessDaftar = pilihan.business[isian.type] ?? []
  const subrogasi = isian.type === 'SUBROGATION'
  return (
    <section className="bordereaux__kartu bordereaux__kepala-form">
      <div className="bordereaux__bagian">
        <span className="bordereaux__judul-bagian">
          {BDX.type}
          <Wajib />
        </span>
        <div className="bordereaux__chip-grup" role="radiogroup" aria-label={BDX.type}>
          {pilihan.type.map((t) => (
            <Chip key={t} nama="bdx-type" nilai={t} dipilih={isian.type === t} onPilih={() => onGantiType(t)} />
          ))}
        </div>
      </div>

      {isian.type !== '' && (
        <div className="bordereaux__bagian">
          <span className="bordereaux__judul-bagian">
            {BDX.typeBusiness}
            <Wajib />
          </span>
          {subrogasi ? (
            <p className="muted">{BDX.subrogasiBonding}</p>
          ) : (
            <div className="bordereaux__chip-kisi" role="radiogroup" aria-label={BDX.typeBusiness}>
              {businessDaftar.map((b) => (
                <Chip key={b} nama="bdx-business" nilai={b} dipilih={isian.business === b} onPilih={() => onGantiBusiness(b)} />
              ))}
            </div>
          )}
        </div>
      )}

      <PanelTreaty
        isian={isian}
        tombol={
          isian.business === '' ? (
            <span className="muted">{BDX.pilihTypeDulu}</span>
          ) : (
            <button type="button" className="btn btn--primary btn--sm" onClick={onPilihMaster}>
              {BDX.chooseMasterTreaty}
            </button>
          )
        }
      />

      <div className="bordereaux__baris-isian">
        <div className="field bordereaux__isian-periode">
          <span className="field__label">
            {BDX.bdxStartEnd}
            <Wajib />
          </span>
          <span className="bordereaux__periode">
            <input
              className="field__input"
              type="date"
              aria-label={BDX.reportStart}
              value={keIso(isian.reportStart)}
              onChange={(e) => onUbah({ reportStart: keKabel(e.target.value) })}
            />
            <span className="muted">–</span>
            <input
              className="field__input"
              type="date"
              aria-label={BDX.reportEnd}
              value={keIso(isian.reportEnd)}
              onChange={(e) => onUbah({ reportEnd: keKabel(e.target.value) })}
            />
          </span>
        </div>
        <label className="field bordereaux__isian-reff">
          <span className="field__label">{BDX.reffNoOfSoa}</span>
          <input className="field__input" maxLength={1000} value={isian.reffNoSoa} onChange={(e) => onUbah({ reffNoSoa: e.target.value })} />
        </label>
        <label className="field bordereaux__isian-reff">
          <span className="field__label">{BDX.reffNoOfBdx}</span>
          <input className="field__input" maxLength={1000} value={isian.reffNoBdx} onChange={(e) => onUbah({ reffNoBdx: e.target.value })} />
        </label>
      </div>
    </section>
  )
}
