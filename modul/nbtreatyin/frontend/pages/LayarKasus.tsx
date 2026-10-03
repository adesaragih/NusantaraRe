// Layar realisasi satu kasus - flow action posisinya:
//   Admin   `InboxPolicyTreatyIn`  -> `GeneralPolicyTreatyIn` / `DetailPolicyTreatyIn`
//   Atasan  `DeptHeadTreatyIn_UW`  -> `GeneralDeptHeadTreatyIn_UW` / `DetailDeptHeadTreatyIn_UW`
// beserta `SpreadingRiskList`, jadwal angsuran, dan `ListSuggest`.
//
// ⛔ Setiap refresh berhitung dikirim ke backend (`POST .../hitung`); layar ini
// tidak menghitung apa pun. Medan wajib datang dari backend (`medanWajib`).

import { useCallback, useEffect, useState } from 'react'

import { Gagal, Memuat, Modal, Panel, type Opsi } from '../../../../inti/frontend/components/ui/dasar'
import {
  ambilAcuan,
  bukaKasus,
  daftar,
  hitung,
  kirimKasus,
  nilai,
  pilihBisnis,
  pilihSumberBisnis,
  riwayatKasus,
  setel,
  setelDaftar,
  simpanKasus,
  terbitkanNomor,
  type Acuan,
  type Baris,
  type Halaman,
  type Layar,
  type NomorPolis,
  type Riwayat,
} from '../api'
import KotakMedan from '../components/KotakMedan'
import PilihBisnis from '../components/PilihBisnis'
import PilihSumberBisnis, { tampilTombolSOB } from '../components/PilihSumberBisnis'
import {
  BAGIAN,
  JUDUL,
  JUDUL_POSISI,
  KOLOM_ANGSURAN,
  KOLOM_SPREADING,
  KOLOM_USULAN,
  KONFIRMASI_TOLAK,
  PESAN,
  PILIHAN_APPROVAL,
  PORTAL,
  TOMBOL,
} from '../labels'
import { MEDAN_ADMIN_UANG, MEDAN_ADMIN_UMUM, MEDAN_ATASAN_UANG, MEDAN_ATASAN_UMUM, medanTampil, type Medan } from '../medan'

const P = 'PolicyTreatyIn.'
const SPREADING = P + 'SpreadingRiskList'
const ANGSURAN = P + 'ListInstallment'
const USULAN = P + 'SuggestList'
const TEMPAT_TANGGAL_PRODUKSI = 'LISTSUGGEST_PRODUCTIONDATE'

function opsi(p: { nilai: string; label: string }[] | null | undefined): Opsi[] {
  return (p ?? []).map((x) => ({ value: x.nilai, label: x.label }))
}

export default function LayarKasus({ id, onKembali }: { id: string; onKembali: (pesan?: string) => void }) {
  const [layar, setLayar] = useState<Layar | null>(null)
  const [h, setH] = useState<Halaman | null>(null)
  const [acuan, setAcuan] = useState<Acuan | null>(null)
  const [riwayat, setRiwayat] = useState<Riwayat[]>([])
  const [galat, setGalat] = useState<unknown>(null)
  const [sibuk, setSibuk] = useState(false)
  const [info, setInfo] = useState('')
  const [popupBisnis, setPopupBisnis] = useState(false)
  const [popupSOB, setPopupSOB] = useState(false)
  const [konfirmasi, setKonfirmasi] = useState(false)
  const [nomor, setNomor] = useState<NomorPolis | null>(null)

  const terima = useCallback((ly: Layar) => {
    setLayar(ly)
    setH(ly.halaman)
  }, [])

  useEffect(() => {
    let dibuang = false
    Promise.all([bukaKasus(id), ambilAcuan(), riwayatKasus(id)])
      .then(([ly, a, r]) => {
        if (dibuang) return
        terima(ly)
        setAcuan(a)
        setRiwayat(r ?? [])
      })
      .catch((e: unknown) => {
        if (!dibuang) setGalat(e)
      })
    return () => {
      dibuang = true
    }
  }, [id, terima])

  if (galat !== null && layar === null) return <Gagal galat={galat} />
  if (layar === null || h === null) return <Memuat />

  const posisi = layar.kasus.positionNote
  const admin = posisi === 'ReasTreatyInAdmin'
  const boleh = layar.bolehKerja
  const wajib = new Set(layar.medanWajib ?? [])

  async function jalankan<T>(f: () => Promise<T>, lanjut: (x: T) => void) {
    setSibuk(true)
    setGalat(null)
    setInfo('')
    try {
      lanjut(await f())
    } catch (e: unknown) {
      setGalat(e)
    } finally {
      setSibuk(false)
    }
  }

  const ubah = (jalur: string, v: string) => setH((x) => (x ? setel(x, jalur, v) : x))
  const refresh = (aksi: string, param?: string, indeks?: number, halaman?: Halaman) =>
    void jalankan(() => hitung(id, { aksi, param, indeks, halaman: halaman ?? h }), terima)

  const selesai = (m: Medan, v: string) => {
    const baru = setel(h, m.jalur, v)
    setH(baru)
    if (m.aksi && boleh) refresh(m.aksi.aksi, m.aksi.param, undefined, baru)
  }

  const kotak = (m: Medan, i: number) => (
    <KotakMedan
      key={`${m.jalur}-${i}`}
      medan={m}
      halaman={h}
      wajib={wajib.has(m.jalur)}
      hanyaBaca={!boleh}
      opsiMataUang={opsi(acuan?.mataUang)}
      opsiMO={opsi(acuan?.mo)}
      onUbah={ubah}
      onSelesai={selesai}
    />
  )

  const ubahBaris = (jalur: string, i: number, kunci: string, v: string) =>
    setH((x) => (x ? setelDaftar(x, jalur, daftar(x, jalur).map((b, j) => (j === i ? { ...b, [kunci]: v } : b))) : x))

  const kirim = () =>
    void jalankan(
      () => kirimKasus(id, h),
      (hasil) => onKembali(hasil.pesanKonversi ? `${PESAN.terkirim} ${hasil.pesanKonversi}` : PESAN.terkirim),
    )

  const tombolKirim = () => {
    if (!boleh) return null
    switch (layar.tombol) {
      case 'kirim':
        return (
          <button type="button" className="btn btn--primary" disabled={sibuk} onClick={kirim}>
            {TOMBOL.submit}
          </button>
        )
      case 'konfirmasi-tolak':
        return (
          <button type="button" className="btn btn--primary" disabled={sibuk} onClick={() => setKonfirmasi(true)}>
            {TOMBOL.submit}
          </button>
        )
      case 'nomor-polis':
        return (
          <button
            type="button"
            className="btn btn--primary"
            disabled={sibuk}
            onClick={() => void jalankan(() => terbitkanNomor(id, h), setNomor)}
          >
            {TOMBOL.submit}
          </button>
        )
    }
    return null
  }

  const spreading = daftar(h, SPREADING)
  const angsuran = daftar(h, ANGSURAN)
  const usulan = daftar(h, USULAN)
  const tampilTanggalProduksi = layar.tempat?.[TEMPAT_TANGGAL_PRODUKSI] && nilai(h, P + 'IsApproved') === '1'

  return (
    <div className="inbox nbti__layar">
      <header className="inbox__kepala">
        <h2 className="inbox__judul">
          {JUDUL_POSISI[posisi] ?? JUDUL.portal} — {layar.kasus.id}
        </h2>
        <button type="button" className="btn btn--ghost" onClick={() => onKembali()}>
          {TOMBOL.kembali}
        </button>
      </header>
      {!boleh && <div className="alert alert--info">{PORTAL.hanyaBaca}</div>}
      {nilai(h, 'NBStatus') && <div className="nbti__status">{nilai(h, 'NBStatus')}</div>}
      {galat !== null && <Gagal galat={galat} />}
      {info && <div className="alert alert--ok">{info}</div>}
      {(layar.pesan ?? []).length > 0 && (
        <div className="alert alert--warn">
          {(layar.pesan ?? []).map((p) => (
            <div key={p}>{p}</div>
          ))}
        </div>
      )}

      <Panel judul={BAGIAN.umum}>
        {admin && boleh && (
          <div className="nbti__aksi">
            <button type="button" className="btn" onClick={() => setPopupBisnis(true)}>
              {TOMBOL.chooseBusiness}
            </button>
            {tampilTombolSOB(h) && (
              <button type="button" className="btn" onClick={() => setPopupSOB(true)}>
                {TOMBOL.selectSOB}
              </button>
            )}
            {nilai(h, P + 'TreatyType') === 'XOL' && (
              <button type="button" className="btn" onClick={() => refresh('TreatyEnableDisableInput')}>
                {TOMBOL.enableDisable}
              </button>
            )}
          </div>
        )}
        <div className="form-grid">{medanTampil(admin ? MEDAN_ADMIN_UMUM : MEDAN_ATASAN_UMUM, h).map(kotak)}</div>
      </Panel>

      <Panel judul={BAGIAN.uang}>
        <div className="form-grid">{medanTampil(admin ? MEDAN_ADMIN_UANG : MEDAN_ATASAN_UANG, h).map(kotak)}</div>
      </Panel>

      <Panel judul={BAGIAN.spreading}>
        {admin && boleh && (
          <button
            type="button"
            className="btn btn--sm"
            onClick={() => setH(setelDaftar(h, SPREADING, [...spreading, {} as Baris]))}
          >
            {TOMBOL.add}
          </button>
        )}
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">{KOLOM_SPREADING.treatyType}</th>
                <th scope="col">{KOLOM_SPREADING.share}</th>
                <th scope="col">{KOLOM_SPREADING.premium}</th>
                <th scope="col">{KOLOM_SPREADING.claimPct}</th>
                <th scope="col">{KOLOM_SPREADING.claim}</th>
                {admin && boleh && <th scope="col" />}
              </tr>
            </thead>
            <tbody>
              {spreading.map((b, i) => (
                <tr key={i}>
                  <td>
                    {admin && boleh ? (
                      <select
                        className="field__input"
                        value={b.TreatyType ?? ''}
                        onChange={(e) => ubahBaris(SPREADING, i, 'TreatyType', e.target.value)}
                      >
                        <option value="" />
                        {(acuan?.spreading ?? []).map((o) => (
                          <option key={o.nilai} value={o.nilai}>
                            {o.label}
                          </option>
                        ))}
                      </select>
                    ) : (
                      (b.TreatyName ?? b.TreatyType ?? '')
                    )}
                  </td>
                  {(['SharePercentage', 'PremiumSpreaded', 'ClaimPercentage', 'ClaimSpreaded'] as const).map((k) => (
                    <td key={k}>
                      {admin && boleh && (k === 'SharePercentage' || k === 'ClaimPercentage') ? (
                        <input
                          className="field__input"
                          value={b[k] ?? ''}
                          onChange={(e) => ubahBaris(SPREADING, i, k, e.target.value)}
                          onBlur={() => refresh('CountSpreading', undefined, i + 1)}
                        />
                      ) : (
                        (b[k] ?? '')
                      )}
                    </td>
                  ))}
                  {admin && boleh && (
                    <td>
                      <button
                        type="button"
                        className="btn btn--sm btn--danger"
                        onClick={() => setH(setelDaftar(h, SPREADING, spreading.filter((_, j) => j !== i)))}
                      >
                        {TOMBOL.delete}
                      </button>
                    </td>
                  )}
                </tr>
              ))}
            </tbody>
            <tfoot>
              <tr>
                <td>{KOLOM_SPREADING.totalShare}</td>
                <td>{nilai(h, P + 'TotalSharePercentagePremium')}</td>
                <td>{nilai(h, P + 'TotalPremium')}</td>
                <td>{nilai(h, P + 'TotalSharePercentageClaim')}</td>
                <td>{nilai(h, P + 'TotalClaim')}</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </Panel>

      <Panel judul={BAGIAN.angsuran}>
        <div className="nbti__aksi">
          <label className="field__label">{KOLOM_ANGSURAN.installment}</label>
          {admin && boleh ? (
            <input
              className="field__input nbti__pendek"
              value={nilai(h, P + 'Installment')}
              onChange={(e) => ubah(P + 'Installment', e.target.value)}
              onBlur={() => refresh('FillPaymentInstallment')}
            />
          ) : (
            <span>{nilai(h, P + 'Installment')}</span>
          )}
        </div>
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">{KOLOM_ANGSURAN.no}</th>
                <th scope="col">{KOLOM_ANGSURAN.dueDate}</th>
                <th scope="col">{KOLOM_ANGSURAN.pct}</th>
                <th scope="col">{KOLOM_ANGSURAN.premium}</th>
                <th scope="col">{KOLOM_ANGSURAN.total}</th>
              </tr>
            </thead>
            <tbody>
              {angsuran.map((b, i) => (
                <tr key={i}>
                  <td>{b.InstallmentNo ?? ''}</td>
                  <td>
                    {admin && boleh ? (
                      <input
                        type="date"
                        className="field__input"
                        value={(b.DueDate ?? '').slice(0, 10)}
                        onChange={(e) => ubahBaris(ANGSURAN, i, 'DueDate', e.target.value)}
                      />
                    ) : (
                      (b.DueDate ?? '')
                    )}
                  </td>
                  <td>
                    {admin && boleh ? (
                      <input
                        className="field__input"
                        value={b.InstallmentPercentage ?? ''}
                        onChange={(e) => ubahBaris(ANGSURAN, i, 'InstallmentPercentage', e.target.value)}
                        onBlur={() => refresh('SetValidateInstallment')}
                      />
                    ) : (
                      (b.InstallmentPercentage ?? '')
                    )}
                  </td>
                  <td>
                    {admin && boleh ? (
                      <input
                        className="field__input"
                        value={b.Premium ?? ''}
                        onChange={(e) => ubahBaris(ANGSURAN, i, 'Premium', e.target.value)}
                        onBlur={() => refresh('CountPctInstallment', undefined, Number(b.InstallmentNo) || i + 1)}
                      />
                    ) : (
                      (b.Premium ?? '')
                    )}
                  </td>
                  <td>{b.PaymentTotal ?? ''}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Panel>

      <Panel judul={BAGIAN.usulan}>
        {boleh && (
          <div className="form-grid">
            <div className="field">
              <span className="field__label">
                {KOLOM_USULAN.putusan}
                <span className="field__req">*</span>
              </span>
              <div className="nbti__radio">
                {PILIHAN_APPROVAL.map((o) => (
                  <label key={o.value}>
                    <input
                      type="radio"
                      name="nbti-approval"
                      checked={nilai(h, P + 'IsApproved') === o.value}
                      onChange={() => {
                        const baru = setel(h, P + 'IsApproved', o.value)
                        setH(baru)
                        refresh('SetDueTo', undefined, undefined, baru)
                      }}
                    />
                    {o.label}
                  </label>
                ))}
              </div>
            </div>
            {tampilTanggalProduksi && (
              <KotakMedan
                medan={{ jalur: P + 'ProductionDate', label: 'Production Date', jenis: 'tanggal' }}
                halaman={h}
                wajib
                hanyaBaca={false}
                opsiMataUang={[]}
                opsiMO={[]}
                onUbah={ubah}
                onSelesai={() => undefined}
              />
            )}
            <KotakMedan
              medan={{ jalur: P + 'Suggest', label: 'Suggest', jenis: 'area' }}
              halaman={h}
              wajib
              hanyaBaca={false}
              opsiMataUang={[]}
              opsiMO={[]}
              onUbah={ubah}
              onSelesai={() => undefined}
            />
          </div>
        )}
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th scope="col">{KOLOM_USULAN.tanggal}</th>
                <th scope="col">{KOLOM_USULAN.operator}</th>
                <th scope="col">{KOLOM_USULAN.putusan}</th>
                <th scope="col">{KOLOM_USULAN.catatan}</th>
              </tr>
            </thead>
            <tbody>
              {usulan.map((b, i) => (
                <tr key={i}>
                  <td>{b.Date ?? ''}</td>
                  <td>{b.OperatorName ?? ''}</td>
                  <td>{PILIHAN_APPROVAL.find((o) => o.value === b.IsApproved)?.label ?? b.IsApproved ?? ''}</td>
                  <td className="nbti__catatan">{b.Suggest ?? ''}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Panel>

      {riwayat.length > 0 && (
        <Panel judul={JUDUL.riwayat}>
          <div className="table-wrap">
            <table>
              <tbody>
                {riwayat.map((r, i) => (
                  <tr key={i}>
                    <td>{r.tglTransfer}</td>
                    <td>{r.workbasket}</td>
                    <td>{r.status}</td>
                    <td>{r.username}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Panel>
      )}

      {boleh && (
        <div className="nbti__kaki">
          {admin && (
            <button
              type="button"
              className="btn"
              disabled={sibuk}
              onClick={() =>
                void jalankan(
                  () => simpanKasus(id, h),
                  (ly) => {
                    terima(ly)
                    setInfo(PESAN.tersimpan)
                  },
                )
              }
            >
              {TOMBOL.save}
            </button>
          )}
          {tombolKirim()}
        </div>
      )}

      {popupBisnis && (
        <PilihBisnis
          onTutup={() => setPopupBisnis(false)}
          onPilih={(idDetail) => {
            setPopupBisnis(false)
            void jalankan(() => pilihBisnis(id, idDetail, h), terima)
          }}
        />
      )}
      {popupSOB && (
        <PilihSumberBisnis
          terpilih={nilai(h, 'Quotation.SourceOfBusiness')}
          sibuk={sibuk}
          onTutup={() => setPopupSOB(false)}
          onPilih={(idAgen, tutup) => {
            if (tutup) setPopupSOB(false)
            void jalankan(() => pilihSumberBisnis(id, idAgen, h), terima)
          }}
        />
      )}
      {konfirmasi && (
        <Modal
          judul={JUDUL.tolak}
          onTutup={() => setKonfirmasi(false)}
          labelBatal={TOMBOL.no}
          aksi={
            <button
              type="button"
              className="btn btn--primary"
              onClick={() => {
                setKonfirmasi(false)
                kirim()
              }}
            >
              {TOMBOL.yes}
            </button>
          }
        >
          <p>{KONFIRMASI_TOLAK}</p>
        </Modal>
      )}
      {nomor && (
        <Modal
          judul={JUDUL.nomorPolis}
          onTutup={() => setNomor(null)}
          aksi={
            <button
              type="button"
              className="btn btn--primary"
              onClick={() => {
                setNomor(null)
                kirim()
              }}
            >
              {TOMBOL.ok}
            </button>
          }
        >
          <dl className="nbti__nomor">
            <dt>pyID</dt>
            <dd>{nomor.id}</dd>
            <dt>PolicyNo</dt>
            <dd>{nomor.policyNo}</dd>
          </dl>
        </Modal>
      )}
    </div>
  )
}
