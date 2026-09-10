import { defineStore } from 'pinia'
import axios from 'axios'

const API_BASE = '/api'

export interface LoginResult {
  success: boolean
  otpRequired?: boolean
  otpSetupRequired?: boolean
  error?: string
}

export interface TodayStatus {
  masuk?: any
  pulang?: any
  is_masuk: boolean
  is_pulang: boolean
  clock_in_time?: string
  clock_out_time?: string
  [key: string]: any
}

export const useAuthStore = defineStore('auth', {
  state: () => {
    const savedToken = localStorage.getItem('garda_token') || ''
    if (savedToken) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${savedToken}`
    }
    return {
      token: savedToken,
      user: (() => {
        try {
          const cached = localStorage.getItem('garda_user')
          return cached ? JSON.parse(cached) : null
        } catch {
          return null
        }
      })(),
      loading: false,
      error: '',
      otpRequired: false,
      tempToken: '',
      usersList: [] as any[],
      presensiHistory: [] as any[],
      todayStatus: {
        masuk: null as any,
        pulang: null as any,
        is_masuk: false,
        is_pulang: false,
        clock_in_time: undefined,
        clock_out_time: undefined
      } as TodayStatus
    }
  },

  getters: {
    isAuthenticated: (state) => !!state.token,
    isSuperAdmin: (state) => state.user?.role === 'superadmin',
    isAdmin: (state) => state.user?.role === 'superadmin' || state.user?.role === 'admin',
    isCallTaker: (state) => state.user?.role === 'call_taker'
  },

  actions: {
    setAuthToken(token: string) {
      this.token = token
      if (token) {
        localStorage.setItem('garda_token', token)
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
      } else {
        localStorage.removeItem('garda_token')
        delete axios.defaults.headers.common['Authorization']
      }
    },

    setUser(user: any) {
      this.user = user
      if (user) {
        localStorage.setItem('garda_user', JSON.stringify(user))
      } else {
        localStorage.removeItem('garda_user')
      }
    },

    async login(userIdentifier: string, passwordStr: string): Promise<LoginResult> {
      this.loading = true
      this.error = ''
      this.otpRequired = false
      this.tempToken = ''

      try {
        const response = await axios.post(`${API_BASE}/auth/login`, { identifier: userIdentifier, password: passwordStr })
        
        if (response.data.otp_required) {
          this.otpRequired = true
          this.tempToken = response.data.temp_token
          const isSetup = !!response.data.otp_setup_required
          return { success: false, otpRequired: true, otpSetupRequired: isSetup }
        }

        this.setAuthToken(response.data.token)
        this.setUser(response.data.user)
        try {
          await this.fetchTodayStatus()
        } catch (e) {
          console.warn('[AuthStore] fetchTodayStatus skipped on login:', e)
        }
        return { success: true }
      } catch (err: any) {
        console.error('[AuthStore] Login failed:', err)
        this.error = err.response?.data?.error || 'Login gagal. Periksa kembali NIP / Email Dinas dan Password Anda.'
        return { success: false, error: this.error }
      } finally {
        this.loading = false
      }
    },

    async verify2fa(code: string): Promise<boolean> {
      this.loading = true
      this.error = ''

      try {
        const response = await axios.post(`${API_BASE}/auth/verify-2fa`, {
          temp_token: this.tempToken,
          code
        })

        this.setAuthToken(response.data.token)
        this.setUser(response.data.user)
        this.otpRequired = false
        this.tempToken = ''
        try {
          await this.fetchTodayStatus()
        } catch (e) {
          console.warn('[AuthStore] fetchTodayStatus skipped on verify2fa:', e)
        }
        return true
      } catch (err: any) {
        console.error('[AuthStore] Verify 2FA failed:', err)
        this.error = err.response?.data?.error || 'Kode verifikasi 2FA salah atau expired.'
        return false
      } finally {
        this.loading = false
      }
    },

    async selfReset2fa(backupCode: string): Promise<boolean> {
      this.loading = true
      this.error = ''

      try {
        const response = await axios.post(`${API_BASE}/auth/2fa/self-reset`, {
          temp_token: this.tempToken,
          backup_code: backupCode
        })

        if (response.data.success) {
          this.otpRequired = false
          this.tempToken = ''
          return true
        }
        return false
      } catch (err: any) {
        console.error('[AuthStore] Self reset 2FA failed:', err)
        this.error = err.response?.data?.error || 'Gagal mereset 2FA secara mandiri. Pastikan kode backup Anda benar.'
        return false
      } finally {
        this.loading = false
      }
    },

    async setup2fa() {
      try {
        const response = await axios.post(`${API_BASE}/auth/2fa/setup`, {}, {
          headers: { Authorization: `Bearer ${this.token || this.tempToken}` }
        })
        return response.data
      } catch (err: any) {
        this.error = err.response?.data?.error || 'Gagal menyiapkan 2FA.'
        return null
      }
    },

    async enable2fa(code: string, secret?: string) {
      try {
        const response = await axios.post(`${API_BASE}/auth/2fa/enable`, { code, secret }, {
          headers: { Authorization: `Bearer ${this.token || this.tempToken}` }
        })
        if (this.user) {
          this.user.totp_enabled = true
          this.setUser(this.user)
        }
        return response.data
      } catch (err: any) {
        this.error = err.response?.data?.error || 'Kode verifikasi 2FA tidak valid.'
        throw err
      }
    },

    async disable2fa() {
      try {
        const response = await axios.post(`${API_BASE}/auth/2fa/disable`, {}, {
          headers: { Authorization: `Bearer ${this.token}` }
        })
        if (this.user) {
          this.user.totp_enabled = false
          this.setUser(this.user)
        }
        return response.data
      } catch (err: any) {
        this.error = err.response?.data?.error || 'Gagal menonaktifkan 2FA.'
        throw err
      }
    },

    async fetchProfile() {
      if (!this.token) return
      try {
        axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
        const response = await axios.get(`${API_BASE}/users/me`)
        this.setUser(response.data.user)
        await this.fetchTodayStatus()
      } catch (err: any) {
        console.error('[AuthStore] Fetch profile failed:', err)
        if (err.response && err.response.status === 401) {
          this.logout()
        }
      }
    },

    async fetchUsers() {
      const token = this.token || localStorage.getItem('garda_token') || ''
      try {
        const headers: any = {}
        if (token) {
          headers['Authorization'] = `Bearer ${token}`
        }
        const response = await axios.get(`${API_BASE}/admin/users`, { headers })
        if (response.data && Array.isArray(response.data.users)) {
          this.usersList = response.data.users
        }
      } catch (err: any) {
        console.error('[AuthStore] Fetch users failed:', err)
        if (err.response && err.response.status === 401) {
          this.logout()
        }
      }
    },

    async createUser(payload: any) {
      try {
        const token = this.token || localStorage.getItem('garda_token') || ''
        const headers: any = {}
        if (token) {
          headers['Authorization'] = `Bearer ${token}`
        }
        const response = await axios.post(`${API_BASE}/admin/users`, {
          nip: payload.nip || '',
          email: payload.email,
          password: payload.password || 'password123',
          name: payload.name,
          jabatan: payload.jabatan || 'SMK Negeri 1 Bulukumba',
          unit_kerja: payload.unit_kerja || 'Rekayasa Perangkat Lunak',
          role: payload.role || 'intern',
          permissions: payload.permissions || ['submit_attendance']
        }, { headers })

        if (response.data && response.data.user) {
          const newUser = response.data.user
          const idx = this.usersList.findIndex((u: any) => u.email === newUser.email || (u.nip && u.nip === newUser.nip))
          if (idx >= 0) {
            this.usersList[idx] = newUser
          } else {
            this.usersList.unshift(newUser)
          }
        }
        await this.fetchUsers()
        return response.data || { success: true }
      } catch (err: any) {
        console.error('[AuthStore] createUser error:', err)
        this.error = err.response?.data?.error || err.message || 'Gagal menambahkan user.'
        return { success: false, error: this.error }
      }
    },

    async updateUser(id: number, payload: any) {
      try {
        const response = await axios.put(`${API_BASE}/admin/users/${id}`, {
          nip: payload.nip || '',
          email: payload.email || '',
          name: payload.name || '',
          jabatan: payload.jabatan || '',
          unit_kerja: payload.unit_kerja || '',
          role: payload.role || 'admin',
          permissions: payload.permissions || [],
          password: payload.password || ''
        }, {
          headers: { Authorization: `Bearer ${this.token}` }
        })
        await this.fetchUsers()
        return response.data
      } catch (err: any) {
        this.error = err.response?.data?.error || 'Gagal memperbarui user.'
        throw err
      }
    },

    async toggleUserActive(id: number, isActive: boolean) {
      try {
        const response = await axios.put(`${API_BASE}/admin/users/${id}/toggle-active`, { is_active: isActive }, {
          headers: { Authorization: `Bearer ${this.token}` }
        })
        await this.fetchUsers()
        return response.data
      } catch (err: any) {
        this.error = err.response?.data?.error || 'Gagal merubah keaktifan.'
        throw err
      }
    },

    async resetUser2fa(id: number) {
      try {
        const response = await axios.post(`${API_BASE}/admin/users/${id}/reset-2fa`, {}, {
          headers: { Authorization: `Bearer ${this.token}` }
        })
        await this.fetchUsers()
        return response.data
      } catch (err: any) {
        this.error = err.response?.data?.error || 'Gagal mereset 2FA user.'
        throw err
      }
    },

    async fetchTodayStatus() {
      if (!this.token) return
      try {
        const response = await axios.get(`${API_BASE}/presensi/today`, {
          headers: { Authorization: `Bearer ${this.token}` }
        })
        this.todayStatus = response.data
      } catch (err) {
        console.error('[AuthStore] Fetch today status failed:', err)
      }
    },

    async fetchHistory() {
      if (!this.token) return
      try {
        const response = await axios.get(`${API_BASE}/presensi/history`, {
          headers: { Authorization: `Bearer ${this.token}` }
        })
        this.presensiHistory = response.data.history
      } catch (err) {
        console.error('[AuthStore] Fetch history failed:', err)
      }
    },

    async clockIn(lat: number, lng: number, qrToken: string) {
      try {
        const response = await axios.post(`${API_BASE}/presensi/clock-in`, {
          latitude: lat,
          longitude: lng,
          qr_token: qrToken
        }, {
          headers: { Authorization: `Bearer ${this.token}` }
        })
        await this.fetchTodayStatus()
        await this.fetchHistory()
        return response.data
      } catch (err: any) {
        this.error = err.response?.data?.error || 'Presensi Masuk gagal.'
        throw err
      }
    },

    async clockOut(lat: number, lng: number, qrToken: string) {
      try {
        const response = await axios.post(`${API_BASE}/presensi/clock-out`, {
          latitude: lat,
          longitude: lng,
          qr_token: qrToken
        }, {
          headers: { Authorization: `Bearer ${this.token}` }
        })
        await this.fetchTodayStatus()
        await this.fetchHistory()
        return response.data
      } catch (err: any) {
        this.error = err.response?.data?.error || 'Presensi Pulang gagal.'
        throw err
      }
    },

    logout() {
      this.setAuthToken('')
      this.setUser(null)
      this.otpRequired = false
      this.tempToken = ''
      this.todayStatus = { masuk: null, pulang: null, is_masuk: false, is_pulang: false }
    }
  }
})
